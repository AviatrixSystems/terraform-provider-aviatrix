package aviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixTransitInstanceCreate,
		ReadContext:   resourceAviatrixTransitInstanceRead,
		UpdateContext: resourceAviatrixTransitInstanceUpdate,
		DeleteContext: resourceAviatrixTransitInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceAviatrixTransitInstanceCustomizeDiff,

		Schema: transitInstanceSchema(),
	}
}

func resourceAviatrixTransitInstanceCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	if d.Id() == "" {
		return nil
	}

	cloudType, ok := d.Get("cloud_type").(int)
	if ok && goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
		return nil
	}

	gwSizeConfig := d.GetRawConfig().GetAttr("gw_size")
	if gwSizeConfig.IsNull() || gwSizeConfig.AsString() == "" {
		// For CSP and AEP, gw_size can be inherited from the group's group_instance_size
		if ok && goaviatrix.IsCloudType(cloudType, goaviatrix.CSPRelatedCloudTypes|goaviatrix.EDGENEO) {
			return nil
		}
		return fmt.Errorf("'gw_size' is required for this transit instance and cannot be removed from the configuration")
	}

	return nil
}

// transitInstanceConfig holds the configuration for creating a transit instance
type transitInstanceConfig struct {
	gateway              *goaviatrix.TransitVpc
	singleAZ             bool
	enableMonitorSubnets bool
	excludedInstances    []string
	rxQueueSize          string
}

func resourceAviatrixTransitInstanceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	// Fetch transit group to get cloud_type, account_name, and vpc_id
	groupUUID := getString(d, "group_uuid")
	transitGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to get transit group %s: %v", groupUUID, err)
	}

	cloudType := transitGroup.CloudType
	log.Printf("[DEBUG] Transit group %s: CloudType=%d, AccountName=%s, VpcID=%s, VpcRegion=%s, GwUUIDList=%v",
		groupUUID, cloudType, transitGroup.AccountName, transitGroup.VpcID, transitGroup.VpcRegion, transitGroup.GwUUIDList)

	// Set computed values from the transit group
	mustSet(d, "group_name", transitGroup.GroupName)
	mustSet(d, "cloud_type", cloudType)
	mustSet(d, "account_name", transitGroup.AccountName)
	mustSet(d, "vpc_id", transitGroup.VpcID)

	// Validate gw_size: not supported for Equinix/Megaport/Self-managed, required for all other types
	gwSize := getString(d, "gw_size")
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
		if gwSize != "" {
			return diag.Errorf("'gw_size' is not supported for Equinix, Megaport, or Self-managed transit instances")
		}
	} else if gwSize == "" {
		// For CSP and AEP gateways, inherit gw_size from the group's group_instance_size
		if goaviatrix.IsCloudType(cloudType, goaviatrix.CSPRelatedCloudTypes|goaviatrix.EDGENEO) && transitGroup.GroupInstanceSize != "" {
			gwSize = transitGroup.GroupInstanceSize
			mustSet(d, "gw_size", gwSize)
		} else {
			return diag.Errorf("'gw_size' is required for CSP gateways and AEP transit instances")
		}
	}

	// Create edge transit gateway for AEP, Equinix, Megaport, Self-managed
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		err := createEdgeTransitInstance(ctx, d, client, transitGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		gwName := getString(d, "gw_name")
		d.SetId(gwName)
		return resourceAviatrixTransitInstanceRead(ctx, d, meta)
	}

	// CSP transit gateways - use create_mct_gateway API
	config, diagErr := buildTransitInstanceConfig(ctx, d, client, transitGroup)
	if diagErr != nil {
		return diagErr
	}

	log.Printf("[INFO] Creating Aviatrix Transit Instance: %#v", config.gateway)

	// Use the create_mct_gateway API which handles both primary and HA based on group state
	createdGwName, err := client.LaunchTransitInstance(config.gateway)
	if err != nil {
		return diag.Errorf("failed to create Aviatrix Transit Instance: %v", err)
	}

	// Set the ID and gw_name to the returned gateway name
	if createdGwName != "" {
		d.SetId(createdGwName)
		mustSet(d, "gw_name", createdGwName)
		config.gateway.GwName = createdGwName
	} else if config.gateway.GwName != "" {
		d.SetId(config.gateway.GwName)
	} else {
		return diag.Errorf("failed to get gateway name from API response")
	}

	if diagErr := configureTransitInstancePostCreate(d, client, config); diagErr != nil {
		return diagErr
	}

	return resourceAviatrixTransitInstanceRead(ctx, d, meta)
}

// createEdgeTransitInstance creates an edge transit gateway (Equinix, AEP/NEO, Megaport, Self-managed)
func createEdgeTransitInstance(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, transitGroup *goaviatrix.GatewayGroup) error {
	cloudType := transitGroup.CloudType
	gwName := getString(d, "gw_name")

	// Get the interface config details
	interfaces := getSet(d, "interfaces").List()
	if len(interfaces) == 0 {
		return fmt.Errorf("at least one interface is required for Edge Transit Instance")
	}

	interfacesList, err := getInterfaceDetails(interfaces, cloudType)
	if err != nil {
		return fmt.Errorf("failed to get the interface details: %w", err)
	}

	// Get management egress IP prefix list
	managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
	managementEgressIPPrefix := ""
	if len(managementEgressIPPrefixList) > 0 {
		managementEgressIPPrefix = strings.Join(managementEgressIPPrefixList, ",")
	}

	ztpFileDownloadPath := getString(d, "ztp_file_download_path")
	ztpFileType := getString(d, "ztp_file_type")

	// Create edge transit gateway using create_mct_gateway API
	gateway := &goaviatrix.TransitVpc{
		GroupUUID:                transitGroup.GroupUUID,
		CloudType:                cloudType,
		AccountName:              transitGroup.AccountName,
		GwName:                   gwName,
		VpcID:                    transitGroup.VpcID,
		VpcSize:                  getString(d, "gw_size"),
		Transit:                  true,
		Interfaces:               interfacesList,
		ManagementEgressIPPrefix: managementEgressIPPrefix,
	}

	// Interface mapping and device_id are required only for AEP/NEO edge gateway
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGENEO) {
		interfaceMappingInput := getList(d, "interface_mapping")
		interfaceMapping, err := getInterfaceMappingDetails(interfaceMappingInput)
		if err != nil {
			return fmt.Errorf("failed to get the interface mapping details: %w", err)
		}
		gateway.InterfaceMapping = interfaceMapping
		gateway.DeviceID = getString(d, "device_id")
	}

	// ZTP file download path is required for Equinix, Megaport, Self-managed edge gateways
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
		gateway.ZtpFileDownloadPath = ztpFileDownloadPath

		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGESELFMANAGED) {
			gateway.ZtpFileType = ztpFileType
			// The controller generates an ISO only when gw_registration_method is
			// "iso"; without this it defaults to cloud-init. Mirrors the legacy
			// aviatrix_transit_gateway resource. AVX-79383.
			gateway.GatewayRegistrationMethod = ztpFileType
		}
	}

	log.Printf("[INFO] Creating Aviatrix Edge Transit Instance: %#v", gateway)

	createdGwName, err := client.LaunchTransitInstance(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Edge Transit Instance: %w", err)
	}

	// Set the gw_name to the returned gateway name
	if createdGwName != "" {
		mustSet(d, "gw_name", createdGwName)
	}

	return nil
}

// buildTransitInstanceConfig builds and validates the transit instance configuration
func buildTransitInstanceConfig(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, transitGroup *goaviatrix.GatewayGroup) (*transitInstanceConfig, diag.Diagnostics) {
	cloudType := transitGroup.CloudType
	accountName := transitGroup.AccountName
	vpcID := transitGroup.VpcID

	gateway := &goaviatrix.TransitVpc{
		GroupUUID:                 transitGroup.GroupUUID,
		CloudType:                 cloudType,
		AccountName:               accountName,
		GwName:                    getString(d, "gw_name"),
		VpcID:                     vpcID,
		VpcSize:                   getString(d, "gw_size"),
		Subnet:                    getString(d, "subnet"),
		Transit:                   true,
		PrivateSubnetEgressTarget: getString(d, "private_subnet_egress_target"),
	}

	// Validate and configure basic settings
	if err := validateAndConfigureBasicSettings(d, gateway, cloudType); err != nil {
		return nil, err
	}

	// Append IPv6 CIDR to gateway.Subnet
	if transitGroup.EnableIPv6 {
		updatedSubnet, subnetErr := validateAndConfigureSubnetWithIPv6Cidr(d, gateway.Subnet, cloudType)
		if subnetErr != nil {
			return nil, diag.Errorf("%s", subnetErr.Error())
		}
		gateway.Subnet = updatedSubnet
	}

	// Validate and configure cloud-specific settings
	if err := validateAndConfigureCloudSpecificSettings(d, gateway, cloudType, transitGroup); err != nil {
		return nil, err
	}

	// Validate and configure the GCP Transit FireNet LAN launch params. FireNet
	// enable/disable is managed on the transit group (AVX-78640); the controller
	// applies the group's desired state when the primary instance launches, and
	// for GCP it needs these LAN params at launch time.
	if err := validateAndConfigureTransitFireNetLan(d, gateway, cloudType, transitGroup.EnableTransitFirenet); err != nil {
		return nil, err
	}

	// Validate and configure monitoring settings
	enableMonitorSubnets, excludedInstances, err := validateAndConfigureMonitoring(d, cloudType)
	if err != nil {
		return nil, err
	}

	// Validate and configure BGP over LAN
	if err := validateAndConfigureBgpOverLan(d, gateway, cloudType); err != nil {
		return nil, err
	}

	// Validate and configure spot instance
	if err := validateAndConfigureSpotInstance(d, gateway); err != nil {
		return nil, err
	}

	// Validate and configure RX queue size
	rxQueueSize := getString(d, "rx_queue_size")
	if rxQueueSize != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return nil, diag.Errorf("rx_queue_size only supports AWS related cloud types")
	}

	// Configure tags
	if err := configureTransitInstanceTags(d, gateway, cloudType); err != nil {
		return nil, err
	}

	// Configure EIP allocation
	if err := configureTransitInstanceEIP(d, gateway, transitGroup.PrivateNetwork); err != nil {
		return nil, err
	}

	return &transitInstanceConfig{
		gateway:              gateway,
		singleAZ:             getBool(d, "single_az_ha"),
		enableMonitorSubnets: enableMonitorSubnets,
		excludedInstances:    excludedInstances,
		rxQueueSize:          rxQueueSize,
	}, nil
}

// validateAndConfigureBasicSettings validates and configures basic gateway settings
func validateAndConfigureBasicSettings(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int) diag.Diagnostics {
	// Validate subnet is required for CSP
	if gateway.Subnet == "" {
		return diag.Errorf("'subnet' is required for CSP transit instance")
	}

	// Single AZ HA
	if getBool(d, "single_az_ha") {
		gateway.SingleAzHa = "enabled"
	} else {
		gateway.SingleAzHa = "disabled"
	}

	// Zone for Azure and GCP
	zone := getString(d, "zone")
	isAzure := goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes)
	isGCP := goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes)
	if zone != "" && !isAzure && !isGCP {
		return diag.Errorf("attribute 'zone' is only for use with Azure (8), Azure GOV (32), Azure CHINA (2048) and GCP (4)")
	}
	if isAzure {
		if zone == "" {
			return diag.Errorf("'zone' is required for Azure (8), Azure GOV (32) and Azure CHINA (2048)")
		}
		if _, errs := validateAzureAZ(zone, "zone"); len(errs) > 0 {
			return diag.Errorf("%s", errs[0].Error())
		}
		gateway.Subnet = fmt.Sprintf("%s~~%s~~", getString(d, "subnet"), zone)
	}
	if isGCP {
		if zone == "" {
			return diag.Errorf("'zone' is required for GCP (4), e.g., 'us-east1-b'")
		}
		if _, errs := validateGCPZone(zone, "zone"); len(errs) > 0 {
			return diag.Errorf("%s", errs[0].Error())
		}
	}

	return nil
}

// validateAndConfigureCloudSpecificSettings validates and configures cloud-specific settings
func validateAndConfigureCloudSpecificSettings(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int, transitGroup *goaviatrix.GatewayGroup) diag.Diagnostics {
	// VPC ID validation
	if goaviatrix.IsCloudType(cloudType, goaviatrix.CSPRelatedCloudTypes) {
		gateway.VpcID = getString(d, "vpc_id")
		if gateway.VpcID == "" {
			return diag.Errorf("'vpc_id' cannot be empty for creating a transit instance")
		}
	} else {
		return diag.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
	}

	// VPC Region - always derived from the transit group
	vpcRegion := transitGroup.VpcRegion
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gateway.VpcRegion = vpcRegion
	} else if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		// For GCP, send the zone in "zone" and the plain region in "vpc_region".
		// Sending the zone as vpc_region breaks the controller's per-region sizing
		// check (AVX-79942), since supported_gw_sizes is keyed by region, not zone.
		zone := getString(d, "zone")
		gateway.Zone = zone
		gateway.VpcRegion = gcpRegionFromZone(zone)
	}

	// OCI specific
	gateway.AvailabilityDomain = getString(d, "availability_domain")
	gateway.FaultDomain = getString(d, "fault_domain")
	if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain == "" || gateway.FaultDomain == "") {
		return diag.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
	}
	if !goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain != "" || gateway.FaultDomain != "") {
		return diag.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
	}

	// Insane mode
	insaneMode := getBool(d, "insane_mode")
	insaneModeAz := getString(d, "insane_mode_az")

	isAWS := goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes)
	isSupportedInsaneCloud := goaviatrix.IsCloudType(
		cloudType,
		goaviatrix.AWSRelatedCloudTypes|
			goaviatrix.GCPRelatedCloudTypes|
			goaviatrix.AzureArmRelatedCloudTypes|
			goaviatrix.OCIRelatedCloudTypes,
	)

	if insaneModeAz != "" && !isAWS {
		return diag.Errorf("'insane_mode_az' is only valid for AWS related clouds")
	}

	enableInsane := insaneMode || insaneModeAz != ""

	if enableInsane {
		if !isSupportedInsaneCloud {
			return diag.Errorf("insane_mode is only supported for AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}

		if isAWS {
			if insaneModeAz == "" {
				return diag.Errorf("insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}
			gateway.Subnet = gateway.Subnet + "~~" + insaneModeAz
		}

		gateway.InsaneMode = "yes"
	} else {
		gateway.InsaneMode = "no"
	}

	// Private subnet egress target validation
	privateSubnetEgressTarget := getString(d, "private_subnet_egress_target")
	if privateSubnetEgressTarget != "" && !enableInsane {
		return diag.Errorf("'private_subnet_egress_target' requires 'insane_mode' to be enabled")
	}

	return nil
}

// validateAndConfigureTransitFireNetLan validates the GCP Transit FireNet LAN launch
// params and folds them into the launch. FireNet enable/disable itself is a group-level
// setting on aviatrix_transit_group (AVX-78640) and is applied by the controller when
// the primary instance launches; only the GCP LAN params live on the instance because
// they must be supplied at launch time. transitFireNet is the group's desired state.
func validateAndConfigureTransitFireNetLan(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int, transitFireNet bool) diag.Diagnostics {
	lanVpcID := getString(d, "lan_vpc_id")
	lanPrivateSubnet := getString(d, "lan_private_subnet")

	if transitFireNet && goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		if lanVpcID == "" || lanPrivateSubnet == "" {
			return diag.Errorf("'lan_vpc_id' and 'lan_private_subnet' are required when 'cloud_type' = 4 (GCP) and the transit group has 'enable_transit_firenet' = true")
		}
		gateway.LanVpcID = lanVpcID
		gateway.LanPrivateSubnet = lanPrivateSubnet
	}

	if (!transitFireNet || !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes)) && (lanVpcID != "" || lanPrivateSubnet != "") {
		return diag.Errorf("'lan_vpc_id' and 'lan_private_subnet' are only valid when 'cloud_type' = 4 (GCP) and the transit group has 'enable_transit_firenet' = true")
	}

	return nil
}

// validateAndConfigureMonitoring validates and configures monitoring settings
func validateAndConfigureMonitoring(d *schema.ResourceData, cloudType int) (bool, []string, diag.Diagnostics) {
	enableMonitorSubnets := getBool(d, "enable_monitor_gateway_subnets")
	var excludedInstances []string
	for _, v := range getSet(d, "monitor_exclude_list").List() {
		excludedInstances = append(excludedInstances, mustString(v))
	}

	if enableMonitorSubnets && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes^goaviatrix.AWSChina) {
		return false, nil, diag.Errorf("'enable_monitor_gateway_subnets' is only valid for AWS (1), AWSGov (256), AWS Top Secret (16384) or AWS Secret (32768)")
	}
	if !enableMonitorSubnets && len(excludedInstances) != 0 {
		return false, nil, diag.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
	}

	return enableMonitorSubnets, excludedInstances, nil
}

// validateAndConfigureBgpOverLan validates and configures BGP over LAN settings
func validateAndConfigureBgpOverLan(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int) diag.Diagnostics {
	bgpOverLan := getBool(d, "enable_bgp_over_lan")
	if bgpOverLan && !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCP) {
		return diag.Errorf("'enable_bgp_over_lan' is only valid for GCP (4), Azure (8), AzureGov (32) or AzureChina (2048)")
	}

	bgpLanInterfacesCount, isCountSet := d.GetOk("bgp_lan_interfaces_count")
	if isCountSet && (!bgpOverLan || !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes)) {
		return diag.Errorf("'bgp_lan_interfaces_count' is only valid for BGP over LAN enabled transit for Azure (8), AzureGov (32) or AzureChina (2048)")
	} else if !isCountSet && bgpOverLan && goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return diag.Errorf("please specify 'bgp_lan_interfaces_count' for BGP over LAN enabled Azure transit: %s", gateway.GwName)
	}

	if bgpOverLan {
		gateway.BgpOverLan = true
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			gateway.BgpLanInterfacesCount = mustInt(bgpLanInterfacesCount)
		}
	}

	return nil
}

// validateAndConfigureSpotInstance validates and configures spot instance settings
func validateAndConfigureSpotInstance(d *schema.ResourceData, gateway *goaviatrix.TransitVpc) diag.Diagnostics {
	enableSpotInstance := getBool(d, "enable_spot_instance")
	spotPrice := getString(d, "spot_price")
	deleteSpot := getBool(d, "delete_spot")

	if enableSpotInstance {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return diag.Errorf("enable_spot_instance only supports AWS and Azure related cloud types")
		}
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && deleteSpot {
			return diag.Errorf("delete_spot only supports Azure")
		}
		gateway.EnableSpotInstance = true
		gateway.SpotPrice = spotPrice
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			gateway.DeleteSpot = deleteSpot
		}
	}

	return nil
}

// configureTransitInstanceTags configures tags for the transit instance
func configureTransitInstanceTags(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int) diag.Diagnostics {
	_, tagsOk := d.GetOk("tags")
	if tagsOk {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return diag.Errorf("error creating transit instance: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		tagsMap, err := extractTags(d, gateway.CloudType)
		if err != nil {
			return diag.Errorf("error creating tags for transit instance: %v", err)
		}
		tagsJSON, err := TagsMapToJson(tagsMap)
		if err != nil {
			return diag.Errorf("failed to add tags when creating transit instance: %v", err)
		}
		gateway.TagJson = tagsJSON
	}
	return nil
}

// configureTransitInstanceEIP configures EIP allocation settings
func configureTransitInstanceEIP(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, privateNetwork bool) diag.Diagnostics {
	allocateNewEip := getBool(d, "allocate_new_eip")
	if allocateNewEip || privateNetwork {
		gateway.ReuseEip = ""
		return nil
	}

	gateway.ReuseEip = "on"
	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		return diag.Errorf("failed to create transit instance: 'allocate_new_eip' can only be set to 'false' when cloud_type is AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048) or AWS Top Secret (16384)")
	}
	if _, ok := d.GetOk("eip"); !ok {
		return diag.Errorf("failed to create transit instance: 'eip' must be set when 'allocate_new_eip' is false")
	}

	azureEipName, azureEipNameOk := d.GetOk("azure_eip_name_resource_group")
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		if !azureEipNameOk {
			return diag.Errorf("failed to create transit instance: 'azure_eip_name_resource_group' must be set when 'allocate_new_eip' is false and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
		}
		gateway.Eip = fmt.Sprintf("%s:%s", mustString(azureEipName), getString(d, "eip"))
	} else {
		if azureEipNameOk {
			return diag.Errorf("failed to create transit instance: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}
		gateway.Eip = getString(d, "eip")
	}

	return nil
}

// configureTransitInstancePostCreate configures settings after the transit instance is created
func configureTransitInstancePostCreate(d *schema.ResourceData, client *goaviatrix.Client, config *transitInstanceConfig) diag.Diagnostics {
	gwName := getString(d, "gw_name")

	// Disable single AZ HA if not enabled
	if !config.singleAZ {
		if err := disableTransitInstanceSingleAZHA(client, gwName); err != nil {
			return err
		}
	}

	// FireNet / Transit FireNet is enabled at the group level (AVX-78640); the
	// controller applies the group's desired state as the primary instance launches.

	// Configure routing settings
	if err := configureTransitInstanceRouting(client, d, config.gateway); err != nil {
		return err
	}

	// Enable monitor gateway subnets
	if config.enableMonitorSubnets {
		if err := client.EnableMonitorGatewaySubnets(config.gateway.GwName, config.excludedInstances); err != nil {
			return diag.Errorf("could not enable monitor gateway subnets: %v", err)
		}
	}

	// Set tunnel detection time
	if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
		if err := client.ModifyTunnelDetectionTime(config.gateway.GwName, mustInt(detectionTime)); err != nil {
			return diag.Errorf("could not set tunnel detection time during Transit Instance creation: %v", err)
		}
	}

	// Set RX queue size
	if config.rxQueueSize != "" {
		gwRxQueueSize := &goaviatrix.Gateway{
			GwName:      gwName,
			RxQueueSize: config.rxQueueSize,
		}
		if err := client.SetRxQueueSize(gwRxQueueSize); err != nil {
			return diag.Errorf("failed to set rx queue size for transit %s: %v", config.gateway.GwName, err)
		}
	}

	return nil
}

// disableTransitInstanceSingleAZHA disables single AZ HA for the transit instance
func disableTransitInstanceSingleAZHA(client *goaviatrix.Client, gwName string) diag.Diagnostics {
	singleAZGateway := &goaviatrix.Gateway{
		GwName:   gwName,
		SingleAZ: "no",
	}
	log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
	if err := client.DisableSingleAZGateway(singleAZGateway); err != nil {
		return diag.Errorf("failed to disable single AZ GW HA: %v", err)
	}
	return nil
}

// configureTransitInstanceRouting configures routing settings for the transit instance
func configureTransitInstanceRouting(client *goaviatrix.Client, d *schema.ResourceData, gateway *goaviatrix.TransitVpc) diag.Diagnostics {
	gwName := getString(d, "gw_name")

	// BGP manual spoke advertise cidrs
	bgpManualSpokeAdvertiseCidrs := getString(d, "bgp_manual_spoke_advertise_cidrs")
	if bgpManualSpokeAdvertiseCidrs != "" {
		gateway.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
		if err := client.SetBgpManualSpokeAdvertisedNetworks(gateway); err != nil {
			return diag.Errorf("failed to set BGP Manual Spoke Advertise Cidrs: %v", err)
		}
	}

	// Customized spoke vpc routes
	if err := configureCustomizedSpokeVpcRoutes(client, d, gwName); err != nil {
		return err
	}

	// Filtered spoke vpc routes
	if err := configureFilteredSpokeVpcRoutes(client, d, gwName); err != nil {
		return err
	}

	// Excluded advertised spoke routes
	if err := configureExcludedAdvertisedSpokeRoutes(client, d, gwName); err != nil {
		return err
	}

	// Customized transit vpc routes
	var customizedTransitVpcRoutes []string
	for _, v := range getSet(d, "customized_transit_vpc_routes").List() {
		customizedTransitVpcRoutes = append(customizedTransitVpcRoutes, mustString(v))
	}
	if len(customizedTransitVpcRoutes) != 0 {
		if err := client.UpdateTransitGatewayCustomizedVpcRoute(gateway.GwName, customizedTransitVpcRoutes); err != nil {
			return diag.Errorf("couldn't update transit instance customized vpc route: %v", err)
		}
	}

	return nil
}

// configureCustomizedSpokeVpcRoutes configures customized spoke VPC routes with retry logic
func configureCustomizedSpokeVpcRoutes(client *goaviatrix.Client, d *schema.ResourceData, gwName string) diag.Diagnostics {
	customizedSpokeVpcRoutes := getString(d, "customized_spoke_vpc_routes")
	if customizedSpokeVpcRoutes == "" {
		return nil
	}

	transitGateway := &goaviatrix.Gateway{
		GwName:                   gwName,
		CustomizedSpokeVpcRoutes: strings.Split(customizedSpokeVpcRoutes, ","),
	}

	for i := 0; ; i++ {
		log.Printf("[INFO] Editing customized routes of transit instance: %s ", transitGateway.GwName)
		err := client.EditGatewayCustomRoutes(transitGateway)
		if err == nil {
			break
		}
		if i <= 10 && strings.Contains(err.Error(), "when it is down") {
			time.Sleep(10 * time.Second)
		} else {
			return diag.Errorf("failed to customize spoke vpc routes of transit instance: %s due to: %v", transitGateway.GwName, err)
		}
	}

	return nil
}

// configureFilteredSpokeVpcRoutes configures filtered spoke VPC routes with retry logic
func configureFilteredSpokeVpcRoutes(client *goaviatrix.Client, d *schema.ResourceData, gwName string) diag.Diagnostics {
	filteredSpokeVpcRoutes := getString(d, "filtered_spoke_vpc_routes")
	if filteredSpokeVpcRoutes == "" {
		return nil
	}

	transitGateway := &goaviatrix.Gateway{
		GwName:                 gwName,
		FilteredSpokeVpcRoutes: strings.Split(filteredSpokeVpcRoutes, ","),
	}

	for i := 0; ; i++ {
		log.Printf("[INFO] Editing filtered routes of transit instance: %s ", transitGateway.GwName)
		err := client.EditGatewayFilterRoutes(transitGateway)
		if err == nil {
			break
		}
		if i <= 10 && strings.Contains(err.Error(), "when it is down") {
			time.Sleep(10 * time.Second)
		} else {
			return diag.Errorf("failed to edit filtered spoke vpc routes of transit instance: %s due to: %v", transitGateway.GwName, err)
		}
	}

	return nil
}

// configureExcludedAdvertisedSpokeRoutes configures excluded advertised spoke routes with retry logic
func configureExcludedAdvertisedSpokeRoutes(client *goaviatrix.Client, d *schema.ResourceData, gwName string) diag.Diagnostics {
	advertisedSpokeRoutesExclude := getString(d, "excluded_advertised_spoke_routes")
	if advertisedSpokeRoutesExclude == "" {
		return nil
	}

	transitGateway := &goaviatrix.Gateway{
		GwName:                gwName,
		AdvertisedSpokeRoutes: strings.Split(advertisedSpokeRoutesExclude, ","),
	}

	for i := 0; ; i++ {
		log.Printf("[INFO] Editing customized routes advertisement of transit instance: %s ", transitGateway.GwName)
		err := client.EditGatewayAdvertisedCidr(transitGateway)
		if err == nil {
			break
		}
		if i <= 10 && strings.Contains(err.Error(), "when it is down") {
			time.Sleep(10 * time.Second)
		} else {
			return diag.Errorf("failed to edit advertised spoke vpc routes of transit instance: %s due to: %v", transitGateway.GwName, err)
		}
	}

	return nil
}

func resourceAviatrixTransitInstanceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	ignoreTagsConfig := client.IgnoreTagsConfig

	var isImport bool
	gwName := getString(d, "gw_name")
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		mustSet(d, "gw_name", id)
		gwName = id
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: getString(d, "account_name"),
		GwName:      gwName,
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't find Aviatrix Transit Instance: %v", err)
	}

	log.Printf("[TRACE] reading transit instance %s: %#v", getString(d, "gw_name"), gw)
	mustSet(d, "cloud_type", gw.CloudType)
	mustSet(d, "account_name", gw.AccountName)
	mustSet(d, "gw_name", gw.GwName)
	mustSet(d, "gw_size", gw.GwSize)

	// group_name is Computed; group_uuid is Required+ForceNew but not Computed,
	// so nothing re-derives it. Resolve both from the owning group's name on
	// every read so they survive terraform import (Read-only path) and cannot
	// silently drift.
	if gw.GroupName != "" {
		mustSet(d, "group_name", gw.GroupName)
	}
	if err := setGroupUUIDFromGatewayName(ctx, client, d, gw.GroupName); err != nil {
		return diag.Errorf("failed to resolve group_uuid for transit instance %s: %s", gw.GwName, err)
	}

	// Edge cloud type
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		mustSet(d, "bgp_lan_ip_list", nil)
		if gw.DeviceID != "" {
			_ = d.Set("device_id", gw.DeviceID)
		}
		// Set interfaces
		if len(gw.Interfaces) != 0 {
			userInterfaces := getSet(d, "interfaces").List()
			userInterfaceOrder, err := getUserInterfaceOrder(userInterfaces)
			if err != nil {
				return diag.Errorf("could not get user interface order: %v", err)
			}
			gwInterfaces := filterCloudManagedEdgeInterfaces(gw.Interfaces, gw.CloudType)
			interfaces := setInterfaceDetails(gwInterfaces, userInterfaceOrder)
			if err = d.Set("interfaces", interfaces); err != nil {
				return diag.Errorf("could not set interfaces into state: %v", err)
			}
		}
		// Set interface mapping
		if len(gw.InterfaceMapping) != 0 {
			interfaceMapping := setInterfaceMappingDetails(gw.InterfaceMapping)
			if err = d.Set("interface_mapping", interfaceMapping); err != nil {
				return diag.Errorf("could not set interface mapping into state: %v", err)
			}
		}
		// Set eip map
		if gw.EipMap != nil {
			log.Printf("[TRACE] eip map: %#v", gw.EipMap)
			eipMap, err := setEipMapDetails(gw.EipMap, gw.IfNamesTranslation)
			if err != nil {
				return diag.Errorf("could not set eip map details: %v", err)
			}
			if err = d.Set("eip_map", eipMap); err != nil {
				return diag.Errorf("could not set eip map into state: %v", err)
			}
		}
		// Set management egress ip prefix list
		if gw.ManagementEgressIPPrefix == "" {
			_ = d.Set("management_egress_ip_prefix_list", nil)
		} else {
			_ = d.Set("management_egress_ip_prefix_list", strings.Split(gw.ManagementEgressIPPrefix, ","))
		}
		return nil
	}

	// CSP transit instance
	mustSet(d, "eip", gw.PublicIP)
	mustSet(d, "public_ip", gw.PublicIP)
	mustSet(d, "cloud_instance_id", gw.CloudnGatewayInstID)
	mustSet(d, "security_group_id", gw.GwSecurityGroupID)
	mustSet(d, "private_ip", gw.PrivateIP)
	mustSet(d, "single_az_ha", gw.SingleAZ == "yes")
	mustSet(d, "image_version", gw.ImageVersion)
	mustSet(d, "software_version", gw.SoftwareVersion)
	mustSet(d, "rx_queue_size", gw.RxQueueSize)
	mustSet(d, "subnet", gw.VpcNet)
	mustSet(d, "tunnel_detection_time", gw.TunnelDetectionTime)

	setGatewayIPv6IPState(d, gw)
	// FireNet is group-level (AVX-78640); only the GCP Transit FireNet LAN launch
	// params live on the instance, so read them back when the gateway has it enabled.
	if gw.EnableTransitFirenet && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "lan_vpc_id", gw.BundleVpcInfo.LAN.VpcID)
		mustSet(d, "lan_private_subnet", strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0])
	}

	// BGP over LAN
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan {
		mustSet(d, "bgp_lan_interfaces_count", gw.BgpLanInterfacesCount)
	} else {
		mustSet(d, "bgp_lan_interfaces_count", nil)
	}
	mustSet(d, "enable_bgp_over_lan", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan)

	// BGP LAN IP list for Azure
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan {
		bgpLanIPInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return diag.Errorf("could not get BGP LAN IP info for Azure transit instance %s: %v", gateway.GwName, err)
		}
		if err = d.Set("bgp_lan_ip_list", bgpLanIPInfo.AzureBgpLanIpList); err != nil {
			return diag.Errorf("could not set bgp_lan_ip_list into state: %v", err)
		}
		if err = d.Set("azure_bgp_lan_ip_list", bgpLanIPInfo.AzureBgpLanIpList); err != nil {
			return diag.Errorf("could not set azure_bgp_lan_ip_list into state: %v", err)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan {
		bgpLanIPInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return diag.Errorf("could not get BGP LAN IP info for GCP transit instance %s: %v", gateway.GwName, err)
		}
		if err = d.Set("bgp_lan_ip_list", bgpLanIPInfo.BgpLanIpList); err != nil {
			return diag.Errorf("could not set bgp_lan_ip_list into state: %v", err)
		}
	} else {
		mustSet(d, "bgp_lan_ip_list", nil)
		mustSet(d, "azure_bgp_lan_ip_list", nil)
	}

	// LAN interface CIDR
	lanCidr, err := client.GetTransitGatewayLanCidr(gw.GwName)
	if err != nil && !errors.Is(err, goaviatrix.ErrNotFound) {
		log.Printf("[WARN] Error getting lan cidr for transit instance %s due to %s", gw.GwName, err)
	}
	mustSet(d, "lan_interface_cidr", lanCidr)

	// Zone for Azure
	if _, zoneIsSet := d.GetOk("zone"); goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || zoneIsSet) &&
		gw.GatewayZone != "AvailabilitySet" {
		mustSet(d, "zone", "az-"+gw.GatewayZone)
	}
	// Zone for GCP
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "zone", gw.GatewayZone)
	}

	// Azure EIP name resource group
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			mustSet(d, "azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Transit Instance %s", gw.GwName)
		}
	}

	// VPC ID and allocate_new_eip by cloud type
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
		if gw.PrivateNetwork {
			mustSet(d, "allocate_new_eip", false)
		} else if gw.AllocateNewEipRead && !gw.EnablePrivateOob {
			mustSet(d, "allocate_new_eip", true)
		} else {
			mustSet(d, "allocate_new_eip", false)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		if gw.PrivateNetwork {
			mustSet(d, "allocate_new_eip", false)
		} else {
			mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if gw.CloudType == goaviatrix.AliCloud {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
		mustSet(d, "allocate_new_eip", true)
	}

	// Insane mode
	if gw.InsaneMode == "yes" {
		mustSet(d, "insane_mode", true)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "insane_mode_az", gw.GatewayZone)
		} else {
			mustSet(d, "insane_mode_az", "")
		}
	} else {
		mustSet(d, "insane_mode", false)
		mustSet(d, "insane_mode_az", "")
	}

	// Customized spoke vpc routes
	if len(gw.CustomizedSpokeVpcRoutes) != 0 {
		if customizedRoutes := getString(d, "customized_spoke_vpc_routes"); customizedRoutes != "" {
			customizedRoutesArray := strings.Split(customizedRoutes, ",")
			if len(goaviatrix.Difference(customizedRoutesArray, gw.CustomizedSpokeVpcRoutes)) == 0 &&
				len(goaviatrix.Difference(gw.CustomizedSpokeVpcRoutes, customizedRoutesArray)) == 0 {
				mustSet(d, "customized_spoke_vpc_routes", customizedRoutes)
			} else {
				mustSet(d, "customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
			}
		} else {
			mustSet(d, "customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
		}
	} else {
		mustSet(d, "customized_spoke_vpc_routes", "")
	}

	// Filtered spoke vpc routes
	if len(gw.FilteredSpokeVpcRoutes) != 0 {
		if filteredSpokeVpcRoutes := getString(d, "filtered_spoke_vpc_routes"); filteredSpokeVpcRoutes != "" {
			filteredSpokeVpcRoutesArray := strings.Split(filteredSpokeVpcRoutes, ",")
			if len(goaviatrix.Difference(filteredSpokeVpcRoutesArray, gw.FilteredSpokeVpcRoutes)) == 0 &&
				len(goaviatrix.Difference(gw.FilteredSpokeVpcRoutes, filteredSpokeVpcRoutesArray)) == 0 {
				mustSet(d, "filtered_spoke_vpc_routes", filteredSpokeVpcRoutes)
			} else {
				mustSet(d, "filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
			}
		} else {
			mustSet(d, "filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
		}
	} else {
		mustSet(d, "filtered_spoke_vpc_routes", "")
	}

	// Excluded advertised spoke routes
	if len(gw.ExcludeCidrList) != 0 {
		if advertisedSpokeRoutes := getString(d, "excluded_advertised_spoke_routes"); advertisedSpokeRoutes != "" {
			advertisedSpokeRoutesArray := strings.Split(advertisedSpokeRoutes, ",")
			if len(goaviatrix.Difference(advertisedSpokeRoutesArray, gw.ExcludeCidrList)) == 0 &&
				len(goaviatrix.Difference(gw.ExcludeCidrList, advertisedSpokeRoutesArray)) == 0 {
				mustSet(d, "excluded_advertised_spoke_routes", advertisedSpokeRoutes)
			} else {
				mustSet(d, "excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
			}
		} else {
			mustSet(d, "excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
		}
	} else {
		mustSet(d, "excluded_advertised_spoke_routes", "")
	}

	// BGP manual spoke advertise cidrs
	var bgpManualSpokeAdvertiseCidrs []string
	if _, ok := d.GetOk("bgp_manual_spoke_advertise_cidrs"); ok {
		bgpManualSpokeAdvertiseCidrs = strings.Split(getString(d, "bgp_manual_spoke_advertise_cidrs"), ",")
	}
	if len(goaviatrix.Difference(bgpManualSpokeAdvertiseCidrs, gw.BgpManualSpokeAdvertiseCidrs)) != 0 ||
		len(goaviatrix.Difference(gw.BgpManualSpokeAdvertiseCidrs, bgpManualSpokeAdvertiseCidrs)) != 0 {
		bgpMSAN := ""
		for i := range gw.BgpManualSpokeAdvertiseCidrs {
			if i == 0 {
				bgpMSAN = bgpMSAN + gw.BgpManualSpokeAdvertiseCidrs[i]
			} else {
				bgpMSAN = bgpMSAN + "," + gw.BgpManualSpokeAdvertiseCidrs[i]
			}
		}
		mustSet(d, "bgp_manual_spoke_advertise_cidrs", bgpMSAN)
	} else {
		mustSet(d, "bgp_manual_spoke_advertise_cidrs", getString(d, "bgp_manual_spoke_advertise_cidrs"))
	}

	// Customized transit vpc routes
	mustSet(d, "customized_transit_vpc_routes", gw.CustomizedTransitVpcRoutes)

	// Monitor gateway subnets
	mustSet(d, "enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
	if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
		return diag.Errorf("setting 'monitor_exclude_list' to state: %v", err)
	}

	// Tags
	// Only overwrite state when the backend actually returns tags. list_vpcs_summary
	// does not return tags for HA gateway instances, so gw.Tags is nil for them;
	// unconditionally setting would wipe the create-time tags from state and produce
	// a perpetual "+ tags" plan diff (AVX-79035). This mirrors aviatrix_spoke_instance.
	if gw.Tags != nil && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		tags := goaviatrix.KeyValueTags(gw.Tags).IgnoreConfig(ignoreTagsConfig)
		if err := d.Set("tags", tags); err != nil {
			log.Printf("[WARN] Error setting tags for (%s): %s", d.Id(), err)
		}
	}

	// OCI specific
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		if gw.GatewayZone != "" {
			mustSet(d, "availability_domain", gw.GatewayZone)
		} else {
			mustSet(d, "availability_domain", getString(d, "availability_domain"))
		}
		mustSet(d, "fault_domain", gw.FaultDomain)
	}

	// Spot instance
	if gw.EnableSpotInstance {
		mustSet(d, "enable_spot_instance", true)
		mustSet(d, "spot_price", gw.SpotPrice)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.DeleteSpot {
			mustSet(d, "delete_spot", gw.DeleteSpot)
		}
	}

	return nil
}

func resourceAviatrixTransitInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}
	log.Printf("[INFO] Updating Aviatrix Transit Instance: %#v", gateway)

	d.Partial(true)

	// Check for non-updatable fields
	if err := validateTransitInstanceUpdateRestrictions(d); err != nil {
		return err
	}

	// Handle edge transit gateway updates
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		if err := updateEdgeTransitInstance(ctx, d, client, gateway); err != nil {
			return err
		}
	} else {
		// CSP transit gateway updates
		// Update Single AZ HA
		if err := updateTransitInstanceSingleAZHA(d, client); err != nil {
			return err
		}

		// Update GW Size (not supported for edge)
		if err := updateTransitInstanceSize(d, client, gateway); err != nil {
			return err
		}
	}

	// Update Tags (common to both CSP and edge)
	if err := updateTransitInstanceTags(d, client, gateway); err != nil {
		return err
	}

	// Update routing configuration
	if err := updateTransitInstanceRouting(d, client, gateway); err != nil {
		return err
	}

	// Update monitoring settings
	if err := updateTransitInstanceMonitoring(d, client, gateway); err != nil {
		return err
	}

	// Update tunnel detection time
	if err := updateTransitInstanceTunnelDetection(d, client, gateway); err != nil {
		return err
	}

	// Update RX queue size
	if err := updateTransitInstanceRxQueueSize(d, client); err != nil {
		return err
	}

	// Update BGP over LAN
	if err := updateTransitInstanceBgpOverLan(d, client, gateway); err != nil {
		return err
	}

	d.Partial(false)
	return resourceAviatrixTransitInstanceRead(ctx, d, meta)
}

// validateTransitInstanceUpdateRestrictions checks for non-updatable fields
func validateTransitInstanceUpdateRestrictions(d *schema.ResourceData) diag.Diagnostics {
	if d.HasChange("allocate_new_eip") {
		return diag.Errorf("updating allocate_new_eip is not allowed")
	}
	if d.HasChange("eip") {
		return diag.Errorf("updating eip is not allowed")
	}
	if d.HasChange("azure_eip_name_resource_group") {
		return diag.Errorf("failed to update transit instance: changing 'azure_eip_name_resource_group' is not allowed")
	}
	if d.HasChange("lan_vpc_id") {
		return diag.Errorf("updating lan_vpc_id is not allowed")
	}
	if d.HasChange("lan_private_subnet") {
		return diag.Errorf("updating lan_private_subnet is not allowed")
	}
	return nil
}

// updateTransitInstanceSingleAZHA updates Single AZ HA setting
func updateTransitInstanceSingleAZHA(d *schema.ResourceData, client *goaviatrix.Client) diag.Diagnostics {
	if !d.HasChange("single_az_ha") {
		return nil
	}

	singleAZGateway := &goaviatrix.Gateway{
		GwName: getString(d, "gw_name"),
	}
	singleAZ := getBool(d, "single_az_ha")

	if singleAZ {
		singleAZGateway.SingleAZ = "yes"
		log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
		if err := client.EnableSingleAZGateway(singleAZGateway); err != nil {
			return diag.Errorf("failed to enable single AZ GW HA for %s: %v", singleAZGateway.GwName, err)
		}
	} else {
		singleAZGateway.SingleAZ = "no"
		log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
		if err := client.DisableSingleAZGateway(singleAZGateway); err != nil {
			return diag.Errorf("failed to disable single AZ GW HA for %s: %v", singleAZGateway.GwName, err)
		}
	}

	return nil
}

// updateTransitInstanceTags updates tags for the transit instance
func updateTransitInstanceTags(d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	if !d.HasChange("tags") {
		return nil
	}

	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return diag.Errorf("failed to update transit instance: adding tags is only supported for AWS (1), Azure (8), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	tags := &goaviatrix.Tags{
		ResourceType: "gw",
		ResourceName: getString(d, "gw_name"),
		CloudType:    gateway.CloudType,
	}

	tagsMap, err := extractTags(d, gateway.CloudType)
	if err != nil {
		return diag.Errorf("failed to update tags for transit instance: %v", err)
	}
	tags.Tags = tagsMap

	tagsJSON, err := TagsMapToJson(tagsMap)
	if err != nil {
		return diag.Errorf("failed to update tags for transit instance: %v", err)
	}
	tags.TagJson = tagsJSON

	if err := client.UpdateTags(tags); err != nil {
		return diag.Errorf("failed to update tags for transit instance: %v", err)
	}

	return nil
}

// updateTransitInstanceSize updates the gateway size
func updateTransitInstanceSize(d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	if !d.HasChange("gw_size") {
		return nil
	}

	gateway.VpcSize = getString(d, "gw_size")
	if err := client.UpdateGateway(gateway); err != nil {
		return diag.Errorf("failed to update Aviatrix Transit Instance size: %v", err)
	}

	return nil
}

// updateTransitInstanceRouting updates all routing-related settings
func updateTransitInstanceRouting(d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	gwName := getString(d, "gw_name")

	// BGP manual spoke advertise cidrs
	if d.HasChange("bgp_manual_spoke_advertise_cidrs") {
		transitGateway := &goaviatrix.TransitVpc{
			GwName:                       gwName,
			BgpManualSpokeAdvertiseCidrs: getString(d, "bgp_manual_spoke_advertise_cidrs"),
		}
		if err := client.SetBgpManualSpokeAdvertisedNetworks(transitGateway); err != nil {
			return diag.Errorf("failed to set BGP Manual Spoke Advertise Cidrs: %v", err)
		}
	}

	// Customized spoke vpc routes
	if d.HasChange("customized_spoke_vpc_routes") {
		if err := updateTransitInstanceCustomizedSpokeRoutes(d, client, gwName); err != nil {
			return err
		}
	}

	// Filtered spoke vpc routes
	if d.HasChange("filtered_spoke_vpc_routes") {
		if err := updateTransitInstanceFilteredSpokeRoutes(d, client, gwName); err != nil {
			return err
		}
	}

	// Excluded advertised spoke routes
	if d.HasChange("excluded_advertised_spoke_routes") {
		if err := updateTransitInstanceExcludedAdvertisedRoutes(d, client, gwName); err != nil {
			return err
		}
	}

	// Customized transit vpc routes
	if d.HasChange("customized_transit_vpc_routes") {
		var customizedTransitVpcRoutes []string
		for _, v := range getSet(d, "customized_transit_vpc_routes").List() {
			customizedTransitVpcRoutes = append(customizedTransitVpcRoutes, mustString(v))
		}
		if err := client.UpdateTransitGatewayCustomizedVpcRoute(gateway.GwName, customizedTransitVpcRoutes); err != nil {
			return diag.Errorf("couldn't update transit instance customized vpc route: %v", err)
		}
	}

	return nil
}

// updateTransitInstanceCustomizedSpokeRoutes updates customized spoke VPC routes
func updateTransitInstanceCustomizedSpokeRoutes(d *schema.ResourceData, client *goaviatrix.Client, gwName string) diag.Diagnostics {
	transitGateway := &goaviatrix.Gateway{
		GwName:                   gwName,
		CustomizedSpokeVpcRoutes: strings.Split(getString(d, "customized_spoke_vpc_routes"), ","),
	}
	if getString(d, "customized_spoke_vpc_routes") == "" {
		transitGateway.CustomizedSpokeVpcRoutes = []string{""}
	}
	if err := client.EditGatewayCustomRoutes(transitGateway); err != nil {
		return diag.Errorf("failed to update customized spoke vpc routes: %v", err)
	}
	return nil
}

// updateTransitInstanceFilteredSpokeRoutes updates filtered spoke VPC routes
func updateTransitInstanceFilteredSpokeRoutes(d *schema.ResourceData, client *goaviatrix.Client, gwName string) diag.Diagnostics {
	transitGateway := &goaviatrix.Gateway{
		GwName:                 gwName,
		FilteredSpokeVpcRoutes: strings.Split(getString(d, "filtered_spoke_vpc_routes"), ","),
	}
	if getString(d, "filtered_spoke_vpc_routes") == "" {
		transitGateway.FilteredSpokeVpcRoutes = []string{""}
	}
	if err := client.EditGatewayFilterRoutes(transitGateway); err != nil {
		return diag.Errorf("failed to update filtered spoke vpc routes: %v", err)
	}
	return nil
}

// updateTransitInstanceExcludedAdvertisedRoutes updates excluded advertised spoke routes
func updateTransitInstanceExcludedAdvertisedRoutes(d *schema.ResourceData, client *goaviatrix.Client, gwName string) diag.Diagnostics {
	transitGateway := &goaviatrix.Gateway{
		GwName:                gwName,
		AdvertisedSpokeRoutes: strings.Split(getString(d, "excluded_advertised_spoke_routes"), ","),
	}
	if getString(d, "excluded_advertised_spoke_routes") == "" {
		transitGateway.AdvertisedSpokeRoutes = []string{""}
	}
	if err := client.EditGatewayAdvertisedCidr(transitGateway); err != nil {
		return diag.Errorf("failed to update excluded advertised spoke routes: %v", err)
	}
	return nil
}

// updateTransitInstanceMonitoring updates monitoring settings
func updateTransitInstanceMonitoring(d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	if d.HasChange("enable_monitor_gateway_subnets") {
		if getBool(d, "enable_monitor_gateway_subnets") {
			excludedInstances := getMonitorExcludeList(d)
			if err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances); err != nil {
				return diag.Errorf("could not enable monitor gateway subnets: %v", err)
			}
		} else {
			if err := client.DisableMonitorGatewaySubnets(gateway.GwName); err != nil {
				return diag.Errorf("could not disable monitor gateway subnets: %v", err)
			}
		}
	} else if d.HasChange("monitor_exclude_list") && getBool(d, "enable_monitor_gateway_subnets") {
		// Need to disable and re-enable to update the exclude list
		excludedInstances := getMonitorExcludeList(d)
		if err := client.DisableMonitorGatewaySubnets(gateway.GwName); err != nil {
			return diag.Errorf("could not disable monitor gateway subnets: %v", err)
		}
		if err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances); err != nil {
			return diag.Errorf("could not enable monitor gateway subnets: %v", err)
		}
	}

	return nil
}

// getMonitorExcludeList extracts the monitor exclude list from resource data
func getMonitorExcludeList(d *schema.ResourceData) []string {
	var excludedInstances []string
	for _, v := range getSet(d, "monitor_exclude_list").List() {
		excludedInstances = append(excludedInstances, mustString(v))
	}
	return excludedInstances
}

// updateTransitInstanceTunnelDetection updates tunnel detection time
func updateTransitInstanceTunnelDetection(d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	if !d.HasChange("tunnel_detection_time") {
		return nil
	}

	if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
		if err := client.ModifyTunnelDetectionTime(gateway.GwName, mustInt(detectionTime)); err != nil {
			return diag.Errorf("could not update tunnel detection time: %v", err)
		}
	}

	return nil
}

// updateTransitInstanceRxQueueSize updates RX queue size
func updateTransitInstanceRxQueueSize(d *schema.ResourceData, client *goaviatrix.Client) diag.Diagnostics {
	if !d.HasChange("rx_queue_size") {
		return nil
	}

	gwRxQueueSize := &goaviatrix.Gateway{
		GwName:      getString(d, "gw_name"),
		RxQueueSize: getString(d, "rx_queue_size"),
	}
	if err := client.SetRxQueueSize(gwRxQueueSize); err != nil {
		return diag.Errorf("could not update rx queue size: %v", err)
	}

	return nil
}

// updateTransitInstanceBgpOverLan updates BGP over LAN settings
func updateTransitInstanceBgpOverLan(d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	if !d.HasChanges("enable_bgp_over_lan", "bgp_lan_interfaces_count") {
		return nil
	}

	// Validate BGP over LAN enable change
	if d.HasChange("enable_bgp_over_lan") {
		if !getBool(d, "enable_bgp_over_lan") {
			return diag.Errorf("disabling BGP over LAN during update is not supported for transit: %s", gateway.GwName)
		}
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return diag.Errorf("enabling BGP over LAN during update is only supported for Azure transit")
		}
		if _, ok := d.GetOk("bgp_lan_interfaces_count"); !ok {
			return diag.Errorf("please specify 'bgp_lan_interfaces_count' to enable BGP over LAN during update for Azure transit: %s", gateway.GwName)
		}
	}

	// Validate BGP LAN interfaces count change
	if d.HasChange("bgp_lan_interfaces_count") {
		if !getBool(d, "enable_bgp_over_lan") || !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return diag.Errorf("could not update BGP over LAN interface count since it only supports BGP over LAN enabled transit for Azure (8), AzureGov (32) or AzureChina (2048)")
		}
		oldCount, newCount := d.GetChange("bgp_lan_interfaces_count")
		if mustInt(oldCount) > mustInt(newCount) {
			return diag.Errorf("deleting BGP over LAN interface during update is not supported for transit: %s", gateway.GwName)
		}
	}

	// Apply the change
	gw := &goaviatrix.Gateway{
		GwName:                gateway.GwName,
		BgpLanInterfacesCount: getInt(d, "bgp_lan_interfaces_count"),
	}
	if err := client.ChangeBgpOverLanIntfCnt(gw); err != nil {
		return diag.Errorf("could not modify BGP over LAN interface count for transit: %s during gateway update: %v", gw.GwName, err)
	}

	return nil
}

// updateEdgeTransitInstance handles updates specific to edge transit gateways
func updateEdgeTransitInstance(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	cloudType := gateway.CloudType
	gwName := getString(d, "gw_name")

	// Get WAN count from interfaces
	interfaceList := getSet(d, "interfaces").List()
	wanCount, err := countInterfaceTypes(interfaceList)
	if err != nil {
		return diag.Errorf("failed to get wan interface count: %v", err)
	}

	// Validate non-updatable edge fields
	if err := validateEdgeTransitInstanceUpdateRestrictions(d); err != nil {
		return err
	}

	// Update interfaces and management egress IP prefix list
	if err := updateEdgeTransitInstanceInterfaces(d, client, cloudType, gwName); err != nil {
		return err
	}

	// Update EIP map
	if err := updateEdgeTransitInstanceEipMap(ctx, d, client, cloudType, gwName, wanCount); err != nil {
		return err
	}

	return nil
}

// validateEdgeTransitInstanceUpdateRestrictions checks for non-updatable edge fields
func validateEdgeTransitInstanceUpdateRestrictions(d *schema.ResourceData) diag.Diagnostics {
	if d.HasChange("device_id") {
		return diag.Errorf("updating device_id is not supported for edge transit instance")
	}
	if d.HasChange("gw_size") {
		return diag.Errorf("updating gw_size is not supported for edge transit instance")
	}
	return nil
}

// updateEdgeTransitInstanceInterfaces updates edge transit gateway interfaces
func updateEdgeTransitInstanceInterfaces(d *schema.ResourceData, client *goaviatrix.Client, cloudType int, gwName string) diag.Diagnostics {
	if !d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		return nil
	}

	interfaceList := getSet(d, "interfaces").List()
	interfaceList = filterCloudManagedInterfaces(interfaceList, cloudType)
	interfaces, err := getInterfaceDetails(interfaceList, cloudType)
	if err != nil {
		return diag.Errorf("failed to get interface details: %v", err)
	}

	gateway := &goaviatrix.TransitVpc{
		GwName:     gwName,
		Interfaces: interfaces,
	}

	managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
	if len(managementEgressIPPrefixList) > 0 {
		gateway.ManagementEgressIPPrefix = strings.Join(managementEgressIPPrefixList, ",")
	}

	if err := client.UpdateEdgeGateway(gateway); err != nil {
		return diag.Errorf("failed to update edge transit instance interfaces: %v", err)
	}

	return nil
}

// updateEdgeTransitInstanceEipMap updates EIP mapping for edge transit gateway
func updateEdgeTransitInstanceEipMap(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, cloudType int, gwName string, wanCount int) diag.Diagnostics {
	if !d.HasChange("eip_map") {
		return nil
	}

	eipMap := getList(d, "eip_map")
	if len(eipMap) == 0 {
		return nil
	}

	eipMapList, err := getEipMapDetails(eipMap, wanCount, cloudType)
	if err != nil {
		return diag.Errorf("failed to get the eip map details: %v", err)
	}

	gateway := &goaviatrix.TransitVpc{
		GwName: gwName,
	}

	if cloudType == goaviatrix.EDGEMEGAPORT {
		log.Printf("[INFO] EIP Map for Edge Mega Port: %#v", eipMapList)
		gateway.LogicalEipMap = eipMapList
		gateway.CloudType = cloudType
		updateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := client.UpdateEdgeGatewayV2(updateCtx, gateway); err != nil {
			return diag.Errorf("failed to update logical eip map for edge transit instance: %v", err)
		}
	} else {
		eipMapJSON, err := json.Marshal(eipMapList)
		if err != nil {
			return diag.Errorf("failed to marshal eip_map to JSON: %v", err)
		}
		gateway.EipMap = string(eipMapJSON)
		if err := client.UpdateEdgeGateway(gateway); err != nil {
			return diag.Errorf("failed to update eip map for edge transit instance: %v", err)
		}
	}

	return nil
}

func resourceAviatrixTransitInstanceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Instance: %#v", gateway)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return diag.Errorf("failed to delete Aviatrix Transit Instance: %v", err)
	}

	return nil
}
