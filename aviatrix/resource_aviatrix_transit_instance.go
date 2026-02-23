package aviatrix

import (
	"context"
	b64 "encoding/base64"
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

		Schema: transitInstanceSchema(),
	}
}

// transitInstanceConfig holds the configuration for creating a transit instance
type transitInstanceConfig struct {
	gateway                   *goaviatrix.TransitVpc
	singleAZ                  bool
	enableFireNet             bool
	enableTransitFireNet      bool
	enableGatewayLoadBalancer bool
	enableMonitorSubnets      bool
	excludedInstances         []string
	rxQueueSize               string
}

func resourceAviatrixTransitInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// Fetch transit group to get cloud_type, account_name, and vpc_id
	groupUUID := getString(d, "group_uuid")
	transitGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to get transit group %s: %v", groupUUID, err)
	}

	cloudType := transitGroup.CloudType
	gwName := getString(d, "gw_name")

	// Determine if this is a primary gateway (gw_count == 0) or HA gateway (gw_count > 0)
	gwCount := len(transitGroup.GwUUIDList)
	isPrimaryGateway := gwCount == 0

	log.Printf("[DEBUG] Transit group %s: CloudType=%d, AccountName=%s, VpcID=%s, VpcRegion=%s, GwUUIDList=%v, isPrimaryGateway=%t",
		groupUUID, cloudType, transitGroup.AccountName, transitGroup.VpcID, transitGroup.VpcRegion, transitGroup.GwUUIDList, isPrimaryGateway)

	// Use group_name as gw_name if not provided (only for primary gateway)
	// For HA gateway, allow auto-generation by leaving gw_name empty
	if gwName == "" && isPrimaryGateway {
		gwName = transitGroup.GroupName
		mustSet(d, "gw_name", gwName)
	}

	// Set computed values from the transit group
	mustSet(d, "group_name", transitGroup.GroupName)
	mustSet(d, "cloud_type", cloudType)
	mustSet(d, "account_name", transitGroup.AccountName)
	mustSet(d, "vpc_id", transitGroup.VpcID)

	// Create edge transit gateway for AEP, Equinix, Megaport, Self-managed
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		err := createEdgeTransitInstance(ctx, d, client, transitGroup, isPrimaryGateway)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(gwName)
		return resourceAviatrixTransitInstanceRead(ctx, d, meta)
	}

	// CSP transit gateways
	if isPrimaryGateway {
		// Build and validate gateway configuration for primary CSP transit gateways
		config, diagErr := buildTransitInstanceConfig(ctx, d, client, transitGroup)
		if diagErr != nil {
			return diagErr
		}

		log.Printf("[INFO] Creating Primary Aviatrix Transit Instance: %#v", config.gateway)

		d.SetId(config.gateway.GwName)
		err = client.LaunchTransitVpc(config.gateway)
		if err != nil {
			return diag.Errorf("failed to create Aviatrix Transit Instance: %v", err)
		}

		// Configure post-creation settings
		if diagErr = configureTransitInstancePostCreate(ctx, d, client, config); diagErr != nil {
			return diagErr
		}
	} else {
		// Create HA transit gateway using group_uuid
		transitHaGateway := &goaviatrix.TransitHaGateway{
			GroupUUID: groupUUID,
			GwName:    gwName,
			GwSize:    getString(d, "gw_size"),
			Subnet:    getString(d, "subnet"),
		}

		// Zone for Azure
		zone := getString(d, "zone")
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && zone != "" {
			transitHaGateway.Subnet = fmt.Sprintf("%s~~%s~~", getString(d, "subnet"), zone)
		}

		// Zone for GCP
		if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			transitHaGateway.Zone = transitGroup.VpcRegion
		}

		// OCI specific
		if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) {
			transitHaGateway.AvailabilityDomain = getString(d, "availability_domain")
			transitHaGateway.FaultDomain = getString(d, "fault_domain")
		}

		// Insane mode
		insaneModeAz := getString(d, "insane_mode_az")
		if insaneModeAz != "" {
			transitHaGateway.Subnet = strings.Join([]string{transitHaGateway.Subnet, insaneModeAz}, "~~")
			transitHaGateway.InsaneMode = "yes"
		}

		// EIP
		if !getBool(d, "allocate_new_eip") {
			transitHaGateway.Eip = getString(d, "eip")
		}

		// Tags
		if _, tagsOk := d.GetOk("tags"); tagsOk {
			tagsMap, err := extractTags(d, cloudType)
			if err == nil {
				tagsJSON, err := TagsMapToJson(tagsMap)
				if err == nil {
					transitHaGateway.TagJSON = tagsJSON
				}
			}
		}

		// Auto-generate HA gateway name if not provided
		if gwName == "" {
			transitHaGateway.AutoGenHaGwName = "yes"
		}

		log.Printf("[INFO] Creating HA Aviatrix Transit Instance: %#v", transitHaGateway)

		haGwName, err := client.CreateTransitHaGw(transitHaGateway)
		if err != nil {
			return diag.Errorf("failed to create HA Aviatrix Transit Instance: %v", err)
		}

		// Set the ID and gw_name to the returned HA gateway name or the provided name
		if haGwName != "" {
			d.SetId(haGwName)
			mustSet(d, "gw_name", haGwName)
		} else if gwName != "" {
			d.SetId(gwName)
		} else {
			return diag.Errorf("failed to get HA gateway name from API response")
		}
	}

	return resourceAviatrixTransitInstanceRead(ctx, d, meta)
}

// createEdgeTransitInstance creates an edge transit gateway (Equinix, AEP/NEO, Megaport, Self-managed)
func createEdgeTransitInstance(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, transitGroup *goaviatrix.GatewayGroup, isPrimaryGateway bool) error {
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

	if isPrimaryGateway {
		// Create primary edge transit gateway
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
			}
		}

		log.Printf("[INFO] Creating Primary Aviatrix Edge Transit Instance: %#v", gateway)

		err = client.LaunchTransitVpc(gateway)
		if err != nil {
			return fmt.Errorf("failed to create primary Aviatrix Edge Transit Instance: %w", err)
		}
	} else {
		// Create HA edge transit gateway using group_uuid
		// Encode interfaces for HA gateway
		interfacesJSON, err := json.Marshal(interfacesList)
		if err != nil {
			return fmt.Errorf("failed to marshal interfaces: %w", err)
		}
		interfacesEncoded := b64.StdEncoding.EncodeToString(interfacesJSON)

		transitHaGateway := &goaviatrix.TransitHaGateway{
			GroupUUID:                transitGroup.GroupUUID,
			CloudType:                cloudType,
			GwName:                   gwName,
			GwSize:                   getString(d, "gw_size"),
			VpcID:                    transitGroup.VpcID,
			Interfaces:               interfacesEncoded,
			ManagementEgressIPPrefix: managementEgressIPPrefix,
		}

		// Interface mapping and device_id are required only for AEP/NEO edge gateway
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGENEO) {
			interfaceMappingInput := getList(d, "interface_mapping")
			interfaceMapping, err := getInterfaceMappingDetails(interfaceMappingInput)
			if err != nil {
				return fmt.Errorf("failed to get the interface mapping details: %w", err)
			}
			transitHaGateway.InterfaceMapping = interfaceMapping
			transitHaGateway.DeviceID = getString(d, "device_id")
		}

		// ZTP file settings for HA
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
			transitHaGateway.ZtpFileDownloadPath = ztpFileDownloadPath

			if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGESELFMANAGED) {
				transitHaGateway.ZtpFileType = ztpFileType
			}
		}

		log.Printf("[INFO] Creating HA Aviatrix Edge Transit Instance: %#v", transitHaGateway)

		_, err = client.CreateTransitHaGw(transitHaGateway)
		if err != nil {
			return fmt.Errorf("failed to create HA Aviatrix Edge Transit Instance: %w", err)
		}
	}

	return nil
}

// buildTransitInstanceConfig builds and validates the transit instance configuration
func buildTransitInstanceConfig(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, transitGroup *goaviatrix.GatewayGroup) (*transitInstanceConfig, diag.Diagnostics) {
	cloudType := transitGroup.CloudType
	accountName := transitGroup.AccountName
	vpcID := transitGroup.VpcID

	gateway := &goaviatrix.TransitVpc{
		GroupUUID:   transitGroup.GroupUUID,
		CloudType:   cloudType,
		AccountName: accountName,
		GwName:      getString(d, "gw_name"),
		VpcID:       vpcID,
		VpcSize:     getString(d, "gw_size"),
		Subnet:      getString(d, "subnet"),
		Transit:     true,
	}

	// Validate and configure basic settings
	if err := validateAndConfigureBasicSettings(d, gateway, cloudType); err != nil {
		return nil, err
	}

	// Validate and configure cloud-specific settings
	if err := validateAndConfigureCloudSpecificSettings(d, gateway, cloudType, transitGroup); err != nil {
		return nil, err
	}

	// Validate and configure feature flags
	enableFireNet, enableTransitFireNet, enableGatewayLoadBalancer, err := validateAndConfigureFeatureFlags(d, gateway, cloudType)
	if err != nil {
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

	// Configure private mode
	if err := configureTransitInstancePrivateMode(ctx, d, gateway, cloudType, client); err != nil {
		return nil, err
	}

	// Configure EIP allocation
	privateModeInfo, _ := client.GetPrivateModeInfo(ctx)
	if err := configureTransitInstanceEIP(d, gateway, privateModeInfo); err != nil {
		return nil, err
	}

	return &transitInstanceConfig{
		gateway:                   gateway,
		singleAZ:                  getBool(d, "single_az_ha"),
		enableFireNet:             enableFireNet,
		enableTransitFireNet:      enableTransitFireNet,
		enableGatewayLoadBalancer: enableGatewayLoadBalancer,
		enableMonitorSubnets:      enableMonitorSubnets,
		excludedInstances:         excludedInstances,
		rxQueueSize:               rxQueueSize,
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

	// Zone for Azure
	zone := getString(d, "zone")
	if !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && zone != "" {
		return diag.Errorf("attribute 'zone' is only for use with Azure (8), Azure GOV (32) and Azure CHINA (2048)")
	}
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && zone == "" {
		return diag.Errorf("'zone' is required for Azure (8), Azure GOV (32) and Azure CHINA (2048)")
	}
	if zone != "" {
		gateway.Subnet = fmt.Sprintf("%s~~%s~~", getString(d, "subnet"), zone)
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
		gateway.Zone = vpcRegion
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
	insaneModeAz := getString(d, "insane_mode_az")
	if insaneModeAz != "" {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			return diag.Errorf("'insane_mode_az' is only valid for AWS related clouds")
		}
		gateway.Subnet = strings.Join([]string{gateway.Subnet, insaneModeAz}, "~~")
		gateway.InsaneMode = "yes"
	}

	return nil
}

// validateAndConfigureFeatureFlags validates and configures feature flags (FireNet, Transit FireNet, GWLB)
func validateAndConfigureFeatureFlags(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int) (bool, bool, bool, diag.Diagnostics) {
	enableFireNet := getBool(d, "enable_firenet")
	enableTransitFireNet := getBool(d, "enable_transit_firenet")
	enableGatewayLoadBalancer := getBool(d, "enable_gateway_load_balancer")
	lanVpcID := getString(d, "lan_vpc_id")
	lanPrivateSubnet := getString(d, "lan_private_subnet")

	if enableFireNet && enableTransitFireNet {
		return false, false, false, diag.Errorf("can't enable firenet function and transit firenet function at the same time")
	}

	if enableTransitFireNet {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			return false, false, false, diag.Errorf("'enable_transit_firenet' is only supported in AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), Azure China (2048)")
		}
		if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			gateway.EnableTransitFireNet = true
		}
		if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			if lanVpcID == "" || lanPrivateSubnet == "" {
				return false, false, false, diag.Errorf("'lan_vpc_id' and 'lan_private_subnet' are required when 'cloud_type' = 4 (GCP) and 'enable_transit_firenet' = true")
			}
			gateway.LanVpcID = lanVpcID
			gateway.LanPrivateSubnet = lanPrivateSubnet
		}
	}

	if (!enableTransitFireNet || !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes)) && (lanVpcID != "" || lanPrivateSubnet != "") {
		return false, false, false, diag.Errorf("'lan_vpc_id' and 'lan_private_subnet' are only valid when 'cloud_type' = 4 (GCP) and 'enable_transit_firenet' = true")
	}

	if enableGatewayLoadBalancer && !enableFireNet && !enableTransitFireNet {
		return false, false, false, diag.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
	}
	if enableGatewayLoadBalancer && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWS) {
		return false, false, false, diag.Errorf("'enable_gateway_load_balancer' is only supported by AWS (1)")
	}

	return enableFireNet, enableTransitFireNet, enableGatewayLoadBalancer, nil
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
	if bgpOverLan && !(goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCP)) {
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

// configureTransitInstancePrivateMode configures private mode settings
func configureTransitInstancePrivateMode(ctx context.Context, d *schema.ResourceData, gateway *goaviatrix.TransitVpc, cloudType int, client *goaviatrix.Client) diag.Diagnostics {
	privateModeInfo, _ := client.GetPrivateModeInfo(ctx)

	if privateModeInfo.EnablePrivateMode {
		if privateModeSubnetZone, ok := d.GetOk("private_mode_subnet_zone"); ok {
			gateway.Subnet = fmt.Sprintf("%s~~%s", gateway.Subnet, mustString(privateModeSubnetZone))
		} else {
			if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
				return diag.Errorf("%q must be set when creating a Transit Instance in AWS with Private Mode enabled on the Controller", "private_mode_subnet_zone")
			}
		}

		if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
				return diag.Errorf("private mode is only supported in AWS and Azure. %q must be empty", "private_mode_lb_vpc_id")
			}
			gateway.LbVpcID = getString(d, "private_mode_lb_vpc_id")
		}
	} else {
		if _, ok := d.GetOk("private_mode_subnet_zone"); ok {
			return diag.Errorf("%q is only valid when Private Mode is enabled on the Controller", "private_mode_subnet_zone")
		}
		if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
			return diag.Errorf("%q is only valid on when Private Mode is enabled", "private_mode_lb_vpc_id")
		}
	}

	return nil
}

// configureTransitInstanceEIP configures EIP allocation settings
func configureTransitInstanceEIP(d *schema.ResourceData, gateway *goaviatrix.TransitVpc, privateModeInfo *goaviatrix.ControllerPrivateModeConfig) diag.Diagnostics {
	if privateModeInfo.EnablePrivateMode {
		return nil
	}

	allocateNewEip := getBool(d, "allocate_new_eip")
	if allocateNewEip {
		gateway.ReuseEip = "off"
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
func configureTransitInstancePostCreate(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, config *transitInstanceConfig) diag.Diagnostics {
	cloudType := getInt(d, "cloud_type")
	gwName := getString(d, "gw_name")

	// Disable single AZ HA if not enabled
	if !config.singleAZ {
		if err := disableTransitInstanceSingleAZHA(client, gwName); err != nil {
			return err
		}
	}

	// Enable FireNet
	if config.enableFireNet {
		if err := enableTransitInstanceFireNet(client, config.gateway, config.enableGatewayLoadBalancer); err != nil {
			return err
		}
	}

	// Enable Transit FireNet for AWS/OCI
	if config.enableTransitFireNet && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		if err := enableTransitInstanceTransitFireNet(client, gwName, config.enableGatewayLoadBalancer); err != nil {
			return err
		}
	}

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

// enableTransitInstanceFireNet enables FireNet for the transit instance
func enableTransitInstanceFireNet(client *goaviatrix.Client, gateway *goaviatrix.TransitVpc, enableGatewayLoadBalancer bool) diag.Diagnostics {
	if enableGatewayLoadBalancer {
		if err := client.EnableGatewayFireNetInterfacesWithGWLB(gateway); err != nil {
			return diag.Errorf("failed to enable transit GW for FireNet Interfaces with Gateway Load Balancer enabled: %v", err)
		}
	} else {
		if err := client.EnableGatewayFireNetInterfaces(gateway); err != nil {
			return diag.Errorf("failed to enable transit GW for FireNet Interfaces: %v", err)
		}
	}
	return nil
}

// enableTransitInstanceTransitFireNet enables Transit FireNet for the transit instance
func enableTransitInstanceTransitFireNet(client *goaviatrix.Client, gwName string, enableGatewayLoadBalancer bool) diag.Diagnostics {
	gwTransitFireNet := &goaviatrix.Gateway{
		GwName: gwName,
	}
	if enableGatewayLoadBalancer {
		if err := client.EnableTransitFireNetWithGWLB(gwTransitFireNet); err != nil {
			return diag.Errorf("failed to enable transit firenet with Gateway Load Balancer enabled: %v", err)
		}
	} else {
		if err := client.EnableTransitFireNet(gwTransitFireNet); err != nil {
			return diag.Errorf("failed to enable transit firenet for %s due to %v", gwTransitFireNet.GwName, err)
		}
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

func resourceAviatrixTransitInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			interfaces := setInterfaceDetails(gw.Interfaces, userInterfaceOrder)
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
	mustSet(d, "enable_firenet", gw.EnableFirenet)
	mustSet(d, "enable_gateway_load_balancer", gw.EnableGatewayLoadBalancer)
	mustSet(d, "enable_transit_firenet", gw.EnableTransitFirenet)

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
		gw.GatewayZone != "AvailabilitySet" && gw.LbVpcId == "" {
		mustSet(d, "zone", "az-"+gw.GatewayZone)
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
		if gw.AllocateNewEipRead && !gw.EnablePrivateOob {
			mustSet(d, "allocate_new_eip", true)
		} else {
			mustSet(d, "allocate_new_eip", false)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if gw.CloudType == goaviatrix.AliCloud {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
		mustSet(d, "allocate_new_eip", true)
	}

	// Insane mode
	if gw.InsaneMode == "yes" {
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "insane_mode_az", gw.GatewayZone)
		} else {
			mustSet(d, "insane_mode_az", "")
		}
	} else {
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
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
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

	// Private mode
	mustSet(d, "private_mode_lb_vpc_id", gw.LbVpcId)
	if gw.LbVpcId != "" && gw.GatewayZone != "AvailabilitySet" {
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "private_mode_subnet_zone", gw.GatewayZone)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			mustSet(d, "private_mode_subnet_zone", "az-"+gw.GatewayZone)
		}
	} else {
		mustSet(d, "private_mode_subnet_zone", nil)
	}

	return nil
}

func resourceAviatrixTransitInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}
	log.Printf("[INFO] Updating Aviatrix Transit Instance: %#v", gateway)

	d.Partial(true)

	// Check for non-updatable fields
	if err := validateTransitInstanceUpdateRestrictions(d, gateway); err != nil {
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
func validateTransitInstanceUpdateRestrictions(d *schema.ResourceData, gateway *goaviatrix.Gateway) diag.Diagnostics {
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
	if d.HasChange("enable_firenet") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSChina|goaviatrix.AzureChina) {
		return diag.Errorf("editing 'enable_firenet' in AWSChina (1024) and AzureChina (2048) is not supported")
	}
	if d.HasChange("enable_transit_firenet") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return diag.Errorf("editing 'enable_transit_firenet' in GCP (4), Azure (8), AzureGov (32) and AzureChina (2048) is not supported")
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

func resourceAviatrixTransitInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
