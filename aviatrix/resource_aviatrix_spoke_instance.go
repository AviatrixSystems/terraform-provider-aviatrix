package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixSpokeInstanceCreate,
		ReadContext:   resourceAviatrixSpokeInstanceRead,
		UpdateContext: resourceAviatrixSpokeInstanceUpdate,
		DeleteContext: resourceAviatrixSpokeInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: MergeSchemaMaps(
			// Required attributes
			spokeInstanceRequiredSchema(),
			// Optional attributes - Basic Configuration
			spokeInstanceOptionalBasicSchema(),
			// Optional attributes - AWS Specific
			spokeInstanceOptionalAWSSchema(),
			// Optional attributes - Azure Specific
			spokeInstanceOptionalAzureSchema(),
			// Optional attributes - OCI Specific
			spokeInstanceOptionalOCISchema(),
			// Optional attributes - Edge Specific
			spokeInstanceOptionalEdgeSchema(),
			// Computed attributes
			spokeInstanceComputedSchema(),
		),
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// buildSpokeVpcFromResourceData constructs a SpokeVpc struct from Terraform resource data.
func buildSpokeVpcFromResourceData(d *schema.ResourceData, gatewayGroup *goaviatrix.GatewayGroup) (*goaviatrix.SpokeVpc, error) {
	spokeGateway := &goaviatrix.SpokeVpc{
		GroupUUID:             getString(d, "group_uuid"),
		GwName:                getString(d, "gw_name"),
		AccountName:           gatewayGroup.AccountName,
		CloudType:             gatewayGroup.CloudType,
		VpcID:                 gatewayGroup.VpcID,
		VpcRegion:             gatewayGroup.VpcRegion,
		Subnet:                getString(d, "subnet"),
		VpcSize:               getString(d, "gw_size"),
		Zone:                  getString(d, "zone"),
		SingleAzHa:            "enabled",
		EnableSpotInstance:    getBool(d, "enable_spot_instance"),
		SpotPrice:             getString(d, "spot_price"),
		DeleteSpot:            getBool(d, "delete_spot"),
		BgpOverLan:            getBool(d, "enable_bgp_over_lan"),
		BgpLanInterfacesCount: getInt(d, "bgp_lan_interfaces_count"),
		AvailabilityDomain:    getString(d, "availability_domain"),
		FaultDomain:           getString(d, "fault_domain"),
	}

	// Single AZ HA
	if !getBool(d, "single_az_ha") {
		spokeGateway.SingleAzHa = "disabled"
	}

	// EIP allocation
	allocateNewEip := getBool(d, "allocate_new_eip")
	if allocateNewEip {
		spokeGateway.ReuseEip = "off"
	} else {
		spokeGateway.ReuseEip = "on"
		spokeGateway.Eip = getString(d, "eip")
	}

	// Insane mode
	insaneMode := getBool(d, "insane_mode")
	if insaneMode {
		spokeGateway.InsaneMode = "yes"
	} else {
		spokeGateway.InsaneMode = "no"
	}

	// Tags
	if _, ok := d.GetOk("tags"); ok {
		tagsMap := mustMap(d.Get("tags"))
		tagsJSON, err := TagsMapToJson(convertTagsMapToStringMap(tagsMap))
		if err != nil {
			return nil, fmt.Errorf("failed to convert tags to JSON: %w", err)
		}
		spokeGateway.TagJson = tagsJSON
	}

	// Private mode
	if lbVpcID := getString(d, "private_mode_lb_vpc_id"); lbVpcID != "" {
		spokeGateway.LbVpcId = lbVpcID
	}

	// Encryption
	if getBool(d, "enable_encrypt_volume") {
		spokeGateway.EncVolume = "yes"
		spokeGateway.CustomerManagedKeys = getString(d, "customer_managed_keys")
	}

	// Insertion gateway
	if getBool(d, "insertion_gateway") {
		spokeGateway.InsertionGateway = true
	}

	return spokeGateway, nil
}

// validateSpokeInstanceConfiguration validates the spoke instance configuration.
//
//nolint:cyclop
func validateSpokeInstanceConfiguration(d *schema.ResourceData, cloudType int) error {
	insaneMode := getBool(d, "insane_mode")
	insaneModeAz := getString(d, "insane_mode_az")
	insertionGateway := getBool(d, "insertion_gateway")
	insertionGatewayAz := getString(d, "insertion_gateway_az")
	enableBgpOverLan := getBool(d, "enable_bgp_over_lan")
	enableEncryptVolume := getBool(d, "enable_encrypt_volume")
	customerManagedKeys := getString(d, "customer_managed_keys")
	enableMonitorGatewaySubnets := getBool(d, "enable_monitor_gateway_subnets")
	monitorExcludeList := getStringSet(d, "monitor_exclude_list")
	enablePrivateVpcDefaultRoute := getBool(d, "enable_private_vpc_default_route")
	enableSkipPublicRouteTableUpdate := getBool(d, "enable_skip_public_route_table_update")
	enableSpotInstance := getBool(d, "enable_spot_instance")
	deleteSpot := getBool(d, "delete_spot")
	rxQueueSize := getString(d, "rx_queue_size")
	availabilityDomain := getString(d, "availability_domain")
	faultDomain := getString(d, "fault_domain")
	eip := getString(d, "eip")
	allocateNewEip := getBool(d, "allocate_new_eip")
	azureEipNameResourceGroup := getString(d, "azure_eip_name_resource_group")
	zone := getString(d, "zone")
	subnet := getString(d, "subnet")
	gwSize := getString(d, "gw_size")
	interfaces := getSet(d, "interfaces").List()
	ztpFileDownloadPath := getString(d, "ztp_file_download_path")
	ztpFileType := getString(d, "ztp_file_type")

	// Edge-specific validation
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		// Interfaces are required for edge spoke gateways
		if len(interfaces) == 0 {
			return fmt.Errorf("'interfaces' is required for Edge spoke instances")
		}

		// ZTP file download path is required for Equinix, Megaport, Self-managed
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
			if ztpFileDownloadPath == "" {
				return fmt.Errorf("'ztp_file_download_path' is required for Equinix, Megaport, and Self-managed edge spoke instances")
			}
		}

		// ZTP file type is required for Self-managed
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGESELFMANAGED) && ztpFileType == "" {
			return fmt.Errorf("'ztp_file_type' is required for Self-managed edge spoke instances")
		}

		// Return early for edge gateways - CSP validations don't apply
		return nil
	}

	// CSP-specific required field validation
	// For non-Edge cloud types, subnet and gw_size are required
	if subnet == "" {
		return fmt.Errorf("'subnet' is required for CSP spoke instances (AWS, Azure, GCP, OCI, AliCloud)")
	}
	if gwSize == "" {
		return fmt.Errorf("'gw_size' is required for CSP spoke instances (AWS, Azure, GCP, OCI, AliCloud)")
	}

	// Zone Validation - only valid for Azure
	if zone != "" && !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("'zone' is only valid for Azure (8), AzureGov (32) and AzureChina (2048)")
	}

	// Monitor Gateway Subnets Validation
	if !enableMonitorGatewaySubnets && len(monitorExcludeList) > 0 {
		return fmt.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
	}

	allowed := goaviatrix.AWS | goaviatrix.AWSGov
	if enableMonitorGatewaySubnets && !goaviatrix.IsCloudType(cloudType, allowed) {
		return fmt.Errorf("'enable_monitor_gateway_subnets' is only valid for AWS (1) or AWSGov (256)")
	}

	// Encryption Validation
	if customerManagedKeys != "" && !enableEncryptVolume {
		return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
	}

	if enableEncryptVolume && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	// BGP Over LAN Validation
	if _, ok := d.GetOk("bgp_lan_interfaces_count"); ok && !enableBgpOverLan {
		return fmt.Errorf("'bgp_lan_interfaces_count' requires enable_bgp_over_lan to be true")
	}

	if enableBgpOverLan && !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("'enable_bgp_over_lan' is only valid for Azure (8), AzureGov (32) or AzureChina (2048)")
	}

	// Insertion Gateway Validation
	if insertionGateway {
		if insaneMode {
			return fmt.Errorf("insertion_gateway and insane_mode cannot both be enabled")
		}
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'insertion_gateway' is only supported for AWS related cloud types")
		}
		if insertionGatewayAz == "" {
			return fmt.Errorf("'insertion_gateway_az' is required when 'insertion_gateway' is enabled")
		}
	} else if insertionGatewayAz != "" {
		return fmt.Errorf("'insertion_gateway_az' requires 'insertion_gateway' to be true")
	}

	// Insane Mode Validation
	if insaneMode {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|
			goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			return fmt.Errorf("insane_mode is only supported for AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}

		if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) && insaneModeAz == "" {
			return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768)")
		}
	}

	// Route Configuration Validation
	if enablePrivateVpcDefaultRoute && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_private_vpc_default_route is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if enableSkipPublicRouteTableUpdate && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_skip_public_route_table_update is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	// Spot Instance Validation
	if enableSpotInstance {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("enable_spot_instance only supports AWS and Azure related cloud types")
		}

		if deleteSpot && !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("delete_spot only supports Azure")
		}
	}

	// RX Queue Size Validation
	if rxQueueSize != "" && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("rx_queue_size only supports AWS related cloud types")
	}

	// OCI Validation
	if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (availabilityDomain == "" || faultDomain == "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
	}
	if !goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (availabilityDomain != "" || faultDomain != "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
	}

	// EIP Allocation Validation
	if !allocateNewEip {
		if eip == "" {
			return fmt.Errorf("'eip' must be set when 'allocate_new_eip' is false")
		}
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && azureEipNameResourceGroup == "" {
			return fmt.Errorf("'azure_eip_name_resource_group' must be set when 'allocate_new_eip' is false and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
		}
	}

	return nil
}

// convertTagsMapToStringMap converts map[string]interface{} to map[string]string.
func convertTagsMapToStringMap(tagsMap map[string]interface{}) map[string]string {
	result := make(map[string]string, len(tagsMap))
	for k, v := range tagsMap {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

// ============================================================================
// CRUD Operations
// ============================================================================

func resourceAviatrixSpokeInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	groupUUID := getString(d, "group_uuid")
	gwName := getString(d, "gw_name")

	// Get the gateway group to retrieve cloud_type, account_name, vpc_id, region for the gateway
	gatewayGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to get gateway group %s: %s", groupUUID, err)
	}

	cloudType := gatewayGroup.CloudType

	// Validate configuration
	if err := validateSpokeInstanceConfiguration(d, cloudType); err != nil {
		return diag.FromErr(err)
	}

	// Determine if this is a primary gateway (gw_count == 0) or HA gateway (gw_count > 0)
	gwCount := len(gatewayGroup.GwUUIDList)
	isPrimaryGateway := gwCount == 0

	// Handle edge spoke gateway creation
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		edgeGwName, err := createEdgeSpokeInstance(ctx, d, client, gatewayGroup, isPrimaryGateway)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(edgeGwName)
		return resourceAviatrixSpokeInstanceRead(ctx, d, meta)
	}

	// Build the spoke gateway from resource data for CSP
	spokeGateway, err := buildSpokeVpcFromResourceData(d, gatewayGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle insane mode AZ for AWS
	if getBool(d, "insane_mode") && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		insaneModeAz := getString(d, "insane_mode_az")
		spokeGateway.Subnet = spokeGateway.Subnet + "~~" + insaneModeAz
	}

	// Handle insertion gateway AZ for AWS
	if getBool(d, "insertion_gateway") && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		insertionGatewayAz := getString(d, "insertion_gateway_az")
		spokeGateway.Subnet = spokeGateway.Subnet + "~~" + insertionGatewayAz
	}

	// Handle Azure zone
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		zone := getString(d, "zone")
		if zone != "" {
			spokeGateway.Subnet = spokeGateway.Subnet + "~~" + zone + "~~"
		}
	}

	// Handle private mode subnet zone
	privateModeSubnetZone := getString(d, "private_mode_subnet_zone")
	if privateModeSubnetZone != "" && spokeGateway.LbVpcId != "" {
		spokeGateway.Subnet = spokeGateway.Subnet + "~~" + privateModeSubnetZone
	}

	// Handle Azure EIP
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && !getBool(d, "allocate_new_eip") {
		azureEipName := getString(d, "azure_eip_name_resource_group")
		spokeGateway.Eip = azureEipName + ":" + spokeGateway.Eip
	}

	// Handle encryption - explicitly set to "no" for AWS if not enabled (same as spoke_gateway)
	if !getBool(d, "enable_encrypt_volume") && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		spokeGateway.EncVolume = "no"
	}

	// Create the spoke instance using the appropriate API based on gateway count
	if isPrimaryGateway {
		log.Printf("[INFO] Creating Primary Spoke Instance: %#v", spokeGateway)
		err = client.LaunchSpokeVpc(spokeGateway)
		if err != nil {
			return diag.Errorf("failed to create primary spoke instance: %s", err)
		}
		d.SetId(gwName)
	} else {
		// Build HA gateway struct with proper fields for create_multicloud_ha_gateway API
		spokeHaGateway := &goaviatrix.SpokeHaGateway{
			GroupUUID:  groupUUID,
			GwName:     gwName,
			GwSize:     getString(d, "gw_size"),
			Subnet:     spokeGateway.Subnet,
			Zone:       getString(d, "zone"),
			InsaneMode: spokeGateway.InsaneMode,
			TagJSON:    spokeGateway.TagJson,
		}

		// Auto-generate HA gateway name if not provided
		if gwName == "" {
			spokeHaGateway.AutoGenHaGwName = "yes"
		}

		// Set EIP only if not allocating new one
		if !getBool(d, "allocate_new_eip") {
			spokeHaGateway.Eip = getString(d, "eip")
		}

		// OCI specific
		spokeHaGateway.AvailabilityDomain = getString(d, "availability_domain")
		spokeHaGateway.FaultDomain = getString(d, "fault_domain")

		// Insertion gateway
		spokeHaGateway.InsertionGateway = getBool(d, "insertion_gateway")

		log.Printf("[INFO] Creating HA Spoke Instance: %#v", spokeHaGateway)
		haGwName, err := client.CreateSpokeHaGw(spokeHaGateway)
		if err != nil {
			return diag.Errorf("failed to create HA spoke instance: %s", err)
		}

		// Set the resource ID after successful creation
		if haGwName != "" {
			d.SetId(haGwName)
			gwName = haGwName
		} else if gwName != "" {
			d.SetId(gwName)
		} else {
			return diag.Errorf("failed to determine HA spoke instance gateway name")
		}
	}

	// Apply post-creation configuration
	if diags := configureSpokeInstancePostCreate(d, client, gwName); diags != nil {
		return diags
	}

	return resourceAviatrixSpokeInstanceRead(ctx, d, meta)
}

// createEdgeSpokeInstance creates an edge spoke gateway (Equinix, AEP/NEO, Megaport, Self-managed)
func createEdgeSpokeInstance(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, gatewayGroup *goaviatrix.GatewayGroup, isPrimaryGateway bool) (string, error) {
	cloudType := gatewayGroup.CloudType
	gwName := getString(d, "gw_name")

	// Get the interface config details
	interfaces := getSet(d, "interfaces").List()
	if len(interfaces) == 0 {
		return "", fmt.Errorf("at least one interface is required for Edge Spoke Instance")
	}

	var interfaceList []*goaviatrix.EdgeSpokeInterface
	for _, if0 := range interfaces {
		if1 := mustMap(if0)
		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:    mustString(if1["logical_ifname"]),
			Type:      mustString(if1["type"]),
			Dhcp:      mustBool(if1["dhcp"]),
			PublicIp:  mustString(if1["public_ip"]),
			IpAddr:    mustString(if1["ip_address"]),
			GatewayIp: mustString(if1["gateway_ip"]),
		}
		interfaceList = append(interfaceList, if2)
	}

	// Management egress IP prefix list
	managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
	managementEgressIPPrefix := ""
	if len(managementEgressIPPrefixList) > 0 {
		managementEgressIPPrefix = strings.Join(managementEgressIPPrefixList, ",")
	}

	// ZTP file settings
	ztpFileDownloadPath := getString(d, "ztp_file_download_path")
	ztpFileType := getString(d, "ztp_file_type")

	if isPrimaryGateway {
		// Create primary edge spoke gateway
		edgeSpoke := &goaviatrix.EdgeSpoke{
			GwName:                   gwName,
			SiteId:                   gatewayGroup.VpcID,
			InterfaceList:            interfaceList,
			ManagementEgressIpPrefix: managementEgressIPPrefix,
		}

		// ZTP file download path is required for Equinix, Megaport, Self-managed edge gateways
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
			edgeSpoke.ZtpFileDownloadPath = ztpFileDownloadPath

			if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGESELFMANAGED) {
				edgeSpoke.ZtpFileType = ztpFileType
			}
		}

		log.Printf("[INFO] Creating Primary Aviatrix Edge Spoke Instance: %#v", edgeSpoke)

		err := client.CreateEdgeSpoke(ctx, edgeSpoke)
		if err != nil {
			return "", fmt.Errorf("failed to create primary Aviatrix Edge Spoke Instance: %w", err)
		}
		return gwName, nil
	}

	{
		// Create HA edge spoke gateway using group_uuid
		edgeSpokeHa := &goaviatrix.EdgeSpokeHa{
			GroupUUID:                gatewayGroup.GroupUUID,
			SiteID:                   gatewayGroup.VpcID,
			InterfaceList:            interfaceList,
			ManagementEgressIPPrefix: managementEgressIPPrefix,
		}

		// ZTP file settings for HA
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT|goaviatrix.EDGESELFMANAGED) {
			edgeSpokeHa.ZtpFileDownloadPath = ztpFileDownloadPath

			if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGESELFMANAGED) {
				edgeSpokeHa.ZtpFileType = ztpFileType
			}
		}

		log.Printf("[INFO] Creating HA Aviatrix Edge Spoke Instance: %#v", edgeSpokeHa)

		haGwName, err := client.CreateEdgeSpokeHa(ctx, edgeSpokeHa)
		if err != nil {
			return "", fmt.Errorf("failed to create HA Aviatrix Edge Spoke Instance: %w", err)
		}
		return haGwName, nil
	}
}

// configureSpokeInstancePostCreate applies configuration that must be done after gateway creation.
func configureSpokeInstancePostCreate(d *schema.ResourceData, client *goaviatrix.Client, gwName string) diag.Diagnostics {
	// Handle single AZ HA
	if !getBool(d, "single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   gwName,
			SingleAZ: "no",
		}
		err := client.DisableSingleAZGateway(singleAZGateway)
		if err != nil {
			return diag.Errorf("failed to disable single AZ GW HA: %s", err)
		}
	}

	// Handle filtered spoke VPC routes
	if filteredRoutes := getString(d, "filtered_spoke_vpc_routes"); filteredRoutes != "" {
		gateway := &goaviatrix.Gateway{
			GwName:                 gwName,
			FilteredSpokeVpcRoutes: strings.Split(filteredRoutes, ","),
		}
		err := client.EditGatewayFilterRoutes(gateway)
		if err != nil {
			return diag.Errorf("failed to set filtered_spoke_vpc_routes: %s", err)
		}
	}

	// Handle included advertised spoke routes
	if includedRoutes := getString(d, "included_advertised_spoke_routes"); includedRoutes != "" {
		gateway := &goaviatrix.Gateway{
			GwName:                gwName,
			AdvertisedSpokeRoutes: strings.Split(includedRoutes, ","),
		}
		err := client.EditGatewayAdvertisedCidr(gateway)
		if err != nil {
			return diag.Errorf("failed to set included_advertised_spoke_routes: %s", err)
		}
	}

	// Handle enable_private_vpc_default_route
	if getBool(d, "enable_private_vpc_default_route") {
		gateway := &goaviatrix.Gateway{GwName: gwName}
		err := client.EnablePrivateVpcDefaultRoute(gateway)
		if err != nil {
			return diag.Errorf("failed to enable private VPC default route: %s", err)
		}
	}

	// Handle enable_skip_public_route_table_update
	if getBool(d, "enable_skip_public_route_table_update") {
		gateway := &goaviatrix.Gateway{GwName: gwName}
		err := client.EnableSkipPublicRouteUpdate(gateway)
		if err != nil {
			return diag.Errorf("failed to enable skip public route table update: %s", err)
		}
	}

	// Handle monitor gateway subnets
	if getBool(d, "enable_monitor_gateway_subnets") {
		excludeList := getStringSet(d, "monitor_exclude_list")
		err := client.EnableMonitorGatewaySubnets(gwName, excludeList)
		if err != nil {
			return diag.Errorf("failed to enable monitor gateway subnets: %s", err)
		}
	}

	// Handle tunnel detection time
	if _, ok := d.GetOk("tunnel_detection_time"); ok {
		tunnelDetectionTime := getInt(d, "tunnel_detection_time")
		err := client.ModifyTunnelDetectionTime(gwName, tunnelDetectionTime)
		if err != nil {
			return diag.Errorf("failed to set tunnel detection time: %s", err)
		}
	}

	// Handle RX queue size
	if rxQueueSize := getString(d, "rx_queue_size"); rxQueueSize != "" {
		gateway := &goaviatrix.Gateway{
			GwName:      gwName,
			RxQueueSize: rxQueueSize,
		}
		err := client.SetRxQueueSize(gateway)
		if err != nil {
			return diag.Errorf("failed to set rx_queue_size: %s", err)
		}
	}

	return nil
}

func resourceAviatrixSpokeInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	gwName := d.Id()
	if gwName == "" {
		return diag.Errorf("resource ID (gateway name) is empty")
	}

	log.Printf("[INFO] Reading Spoke Instance: %s", gwName)

	gateway, err := client.GetGateway(&goaviatrix.Gateway{GwName: gwName})
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read spoke instance: %s", err)
	}

	// Check if this is an edge spoke gateway
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		return readEdgeSpokeInstance(ctx, d, client, gateway)
	}

	// Required attributes
	mustSet(d, "gw_name", gateway.GwName)
	mustSet(d, "subnet", gateway.VpcNet)
	mustSet(d, "gw_size", gateway.GwSize)

	// Basic optional attributes
	mustSet(d, "allocate_new_eip", gateway.AllocateNewEipRead)
	mustSet(d, "single_az_ha", gateway.SingleAZ == "yes")
	mustSet(d, "private_mode_lb_vpc_id", gateway.LbVpcId)
	mustSet(d, "insane_mode", gateway.InsaneMode == "yes")
	mustSet(d, "tunnel_detection_time", gateway.TunnelDetectionTime)
	mustSet(d, "enable_private_vpc_default_route", gateway.PrivateVpcDefaultEnabled)
	mustSet(d, "enable_skip_public_route_table_update", gateway.SkipPublicVpcUpdateEnabled)

	// Spot instance
	mustSet(d, "enable_spot_instance", gateway.EnableSpotInstance)
	mustSet(d, "spot_price", gateway.SpotPrice)
	mustSet(d, "delete_spot", gateway.DeleteSpot)

	// BGP over LAN
	mustSet(d, "enable_bgp_over_lan", gateway.EnableBgpOverLan)
	if gateway.EnableBgpOverLan {
		mustSet(d, "bgp_lan_interfaces_count", gateway.BgpLanInterfacesCount)
	}

	// Azure BGP LAN IP list
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gateway.EnableBgpOverLan {
		bgpLanIPInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return diag.Errorf("could not get BGP LAN IP info for Azure spoke instance %s: %s", gateway.GwName, err)
		}
		mustSet(d, "azure_bgp_lan_ip_list", bgpLanIPInfo.AzureBgpLanIpList)
	} else {
		mustSet(d, "azure_bgp_lan_ip_list", nil)
	}

	// Monitor gateway subnets
	mustSet(d, "enable_monitor_gateway_subnets", gateway.MonitorSubnetsAction == "enable")
	if gateway.MonitorSubnetsAction == "enable" {
		mustSet(d, "monitor_exclude_list", gateway.MonitorExcludeGWList)
	}

	// Encryption
	mustSet(d, "enable_encrypt_volume", gateway.EnableEncryptVolume)

	// OCI specific
	mustSet(d, "availability_domain", gateway.AvailabilityDomain)
	mustSet(d, "fault_domain", gateway.FaultDomain)

	// Filtered and advertised routes
	if len(gateway.FilteredSpokeVpcRoutes) > 0 {
		mustSet(d, "filtered_spoke_vpc_routes", strings.Join(gateway.FilteredSpokeVpcRoutes, ","))
	} else {
		mustSet(d, "filtered_spoke_vpc_routes", "")
	}

	if len(gateway.AdvertisedSpokeRoutes) > 0 {
		mustSet(d, "included_advertised_spoke_routes", strings.Join(gateway.AdvertisedSpokeRoutes, ","))
	} else {
		mustSet(d, "included_advertised_spoke_routes", "")
	}

	// Computed attributes
	mustSet(d, "security_group_id", gateway.GwSecurityGroupID)
	mustSet(d, "cloud_instance_id", gateway.CloudnGatewayInstID)
	mustSet(d, "private_ip", gateway.PrivateIP)
	mustSet(d, "public_ip", gateway.PublicIP)
	mustSet(d, "eip", gateway.Eip)
	mustSet(d, "software_version", gateway.SoftwareVersion)
	mustSet(d, "image_version", gateway.ImageVersion)

	// AWS specific
	mustSet(d, "rx_queue_size", gateway.RxQueueSize)

	// Tags
	if gateway.Tags != nil {
		mustSet(d, "tags", gateway.Tags)
	}

	return nil
}

// readEdgeSpokeInstance reads edge spoke gateway attributes
func readEdgeSpokeInstance(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, gateway *goaviatrix.Gateway) diag.Diagnostics {
	gwName := gateway.GwName

	// Get edge spoke details
	edgeSpoke, err := client.GetEdgeSpoke(ctx, gwName)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read edge spoke instance: %s", err)
	}

	// Basic attributes
	mustSet(d, "gw_name", edgeSpoke.GwName)

	// Management egress IP prefix list
	prefix := edgeSpoke.ManagementEgressIpPrefix
	if prefix == "" {
		mustSet(d, "management_egress_ip_prefix_list", []string{})
	} else {
		mustSet(d, "management_egress_ip_prefix_list", strings.Split(prefix, ","))
	}

	// ZTP file type
	mustSet(d, "ztp_file_type", edgeSpoke.ZtpFileType)

	// Interfaces
	var interfaces []map[string]interface{}
	for _, iface := range edgeSpoke.InterfaceList {
		ifaceMap := map[string]interface{}{
			"logical_ifname": iface.IfName,
			"type":           iface.Type,
			"dhcp":           iface.Dhcp,
			"public_ip":      iface.PublicIp,
			"ip_address":     iface.IpAddr,
			"gateway_ip":     iface.GatewayIp,
		}
		interfaces = append(interfaces, ifaceMap)
	}
	mustSet(d, "interfaces", interfaces)

	// Tunnel detection time (from base gateway)
	mustSet(d, "tunnel_detection_time", gateway.TunnelDetectionTime)

	// Tags (from base gateway)
	if gateway.Tags != nil {
		mustSet(d, "tags", gateway.Tags)
	}

	// Computed attributes from base gateway
	mustSet(d, "software_version", gateway.SoftwareVersion)
	mustSet(d, "image_version", gateway.ImageVersion)

	return nil
}

//nolint:funlen,cyclop
func resourceAviatrixSpokeInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	gwName := d.Id()

	log.Printf("[INFO] Updating Spoke Instance: %s", gwName)

	groupUUID := getString(d, "group_uuid")

	// Get the gateway group to retrieve cloud_type for validation
	gatewayGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to get gateway group for validation: %s", err)
	}
	cloudType := gatewayGroup.CloudType

	if err := validateSpokeInstanceConfiguration(d, cloudType); err != nil {
		return diag.FromErr(err)
	}

	// Common updates for both CSP and edge spoke gateways
	// Tags
	if d.HasChange("tags") {
		tagsMap := mustMap(d.Get("tags"))
		tagsJSON, err := TagsMapToJson(convertTagsMapToStringMap(tagsMap))
		if err != nil {
			return diag.Errorf("failed to convert tags to JSON: %s", err)
		}
		err = client.UpdateTags(&goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: gwName,
			TagJson:      tagsJSON,
			CloudType:    cloudType,
		})
		if err != nil {
			return diag.Errorf("failed to update tags: %s", err)
		}
	}

	// Tunnel Detection Time
	if d.HasChange("tunnel_detection_time") {
		tunnelDetectionTime := getInt(d, "tunnel_detection_time")
		if err := client.ModifyTunnelDetectionTime(gwName, tunnelDetectionTime); err != nil {
			return diag.Errorf("failed to update tunnel_detection_time: %s", err)
		}
	}

	// Handle edge spoke gateway updates
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		if diags := updateEdgeSpokeInstance(ctx, d, client, cloudType, gwName); diags != nil {
			return diags
		}
		return resourceAviatrixSpokeInstanceRead(ctx, d, meta)
	}

	// Gateway Size (CSP only - not supported for edge)
	if d.HasChange("gw_size") {
		gateway := &goaviatrix.Gateway{
			CloudType: cloudType,
			GwName:    gwName,
			VpcSize:   getString(d, "gw_size"),
		}
		err := client.UpdateGateway(gateway)
		if err != nil {
			return diag.Errorf("failed to update gw_size: %s", err)
		}
	}

	// Filtered Spoke VPC Routes
	if d.HasChange("filtered_spoke_vpc_routes") {
		filteredRoutes := getString(d, "filtered_spoke_vpc_routes")
		var routes []string
		if filteredRoutes != "" {
			routes = strings.Split(filteredRoutes, ",")
		}
		gateway := &goaviatrix.Gateway{
			GwName:                 gwName,
			FilteredSpokeVpcRoutes: routes,
		}
		err := client.EditGatewayFilterRoutes(gateway)
		if err != nil {
			return diag.Errorf("failed to update filtered_spoke_vpc_routes: %s", err)
		}
	}

	// Included Advertised Spoke Routes
	if d.HasChange("included_advertised_spoke_routes") {
		includedRoutes := getString(d, "included_advertised_spoke_routes")
		var routes []string
		if includedRoutes != "" {
			routes = strings.Split(includedRoutes, ",")
		}
		gateway := &goaviatrix.Gateway{
			GwName:                gwName,
			AdvertisedSpokeRoutes: routes,
		}
		err := client.EditGatewayAdvertisedCidr(gateway)
		if err != nil {
			return diag.Errorf("failed to update included_advertised_spoke_routes: %s", err)
		}
	}

	// Private VPC Default Route
	if d.HasChange("enable_private_vpc_default_route") {
		gateway := &goaviatrix.Gateway{GwName: gwName}
		if getBool(d, "enable_private_vpc_default_route") {
			err := client.EnablePrivateVpcDefaultRoute(gateway)
			if err != nil {
				return diag.Errorf("failed to enable Private VPC Default Route: %s", err)
			}
		} else {
			err := client.DisablePrivateVpcDefaultRoute(gateway)
			if err != nil {
				return diag.Errorf("failed to disable Private VPC Default Route: %s", err)
			}
		}
	}

	// Skip Public Route Table Update
	if d.HasChange("enable_skip_public_route_table_update") {
		gateway := &goaviatrix.Gateway{GwName: gwName}
		if getBool(d, "enable_skip_public_route_table_update") {
			err := client.EnableSkipPublicRouteUpdate(gateway)
			if err != nil {
				return diag.Errorf("failed to enable Skip Public Route Table Update: %s", err)
			}
		} else {
			err := client.DisableSkipPublicRouteUpdate(gateway)
			if err != nil {
				return diag.Errorf("failed to disable Skip Public Route Table Update: %s", err)
			}
		}
	}

	// Monitor Gateway Subnets
	if d.HasChange("enable_monitor_gateway_subnets") {
		if getBool(d, "enable_monitor_gateway_subnets") {
			excludeList := getStringSet(d, "monitor_exclude_list")
			err := client.EnableMonitorGatewaySubnets(gwName, excludeList)
			if err != nil {
				return diag.Errorf("failed to enable Monitor Gateway Subnets: %s", err)
			}
		} else {
			err := client.DisableMonitorGatewaySubnets(gwName)
			if err != nil {
				return diag.Errorf("failed to disable Monitor Gateway Subnets: %s", err)
			}
		}
	} else if d.HasChange("monitor_exclude_list") && getBool(d, "enable_monitor_gateway_subnets") {
		excludeList := getStringSet(d, "monitor_exclude_list")
		err := client.EnableMonitorGatewaySubnets(gwName, excludeList)
		if err != nil {
			return diag.Errorf("failed to update monitor_exclude_list: %s", err)
		}
	}

	// RX Queue Size
	if d.HasChange("rx_queue_size") {
		rxQueueSize := getString(d, "rx_queue_size")
		gateway := &goaviatrix.Gateway{
			GwName:      gwName,
			RxQueueSize: rxQueueSize,
		}
		err := client.SetRxQueueSize(gateway)
		if err != nil {
			return diag.Errorf("failed to update rx_queue_size: %s", err)
		}
	}

	// Single AZ HA
	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName: gwName,
		}

		if getBool(d, "single_az_ha") {
			singleAZGateway.SingleAZ = "yes"
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return diag.Errorf("failed to enable single AZ GW HA: %s", err)
			}
		} else {
			singleAZGateway.SingleAZ = "no"
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return diag.Errorf("failed to disable single AZ GW HA: %s", err)
			}
		}
	}

	// Encrypt Volume
	if d.HasChange("enable_encrypt_volume") {
		if getBool(d, "enable_encrypt_volume") {
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
				return diag.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              gwName,
				CustomerManagedKeys: getString(d, "customer_managed_keys"),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return diag.Errorf("failed to enable encrypt gateway volume: %s", err)
			}
		} else {
			return diag.Errorf("cannot disable encrypt volume for gateway: %s", gwName)
		}
	} else if d.HasChange("customer_managed_keys") {
		return diag.Errorf("updating customer_managed_keys only is not allowed")
	}

	return resourceAviatrixSpokeInstanceRead(ctx, d, meta)
}

// updateEdgeSpokeInstance handles updates specific to edge spoke gateways
func updateEdgeSpokeInstance(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, cloudType int, gwName string) diag.Diagnostics {
	// Validate non-updatable edge fields
	if d.HasChange("device_id") {
		return diag.Errorf("updating device_id is not supported for edge spoke instance")
	}
	if d.HasChange("gw_size") {
		return diag.Errorf("updating gw_size is not supported for edge spoke instance")
	}

	// Update interfaces and management egress IP prefix list
	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		interfaceList := getSet(d, "interfaces").List()
		interfacesEncoded, err := getInterfaceDetails(interfaceList, cloudType)
		if err != nil {
			return diag.Errorf("failed to get interface details: %v", err)
		}

		gateway := &goaviatrix.TransitVpc{
			GwName:     gwName,
			Interfaces: interfacesEncoded,
		}

		managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
		gateway.ManagementEgressIPPrefix = strings.Join(managementEgressIPPrefixList, ",")

		if err := client.UpdateEdgeGateway(gateway); err != nil {
			return diag.Errorf("failed to update edge spoke instance interfaces: %v", err)
		}
	}

	return nil
}

func resourceAviatrixSpokeInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	gwName := d.Id()
	log.Printf("[INFO] Deleting Spoke Instance: %s", gwName)

	groupUUID := getString(d, "group_uuid")

	// Get the gateway group to retrieve cloud_type for deletion
	gatewayGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to get gateway group: %s", err)
	}

	// Use appropriate delete API based on cloud type
	if goaviatrix.IsCloudType(gatewayGroup.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		err = client.DeleteEdgeSpoke(ctx, gwName)
		if err != nil {
			return diag.Errorf("failed to delete edge spoke instance: %s", err)
		}
	} else {
		err = client.DeleteGateway(&goaviatrix.Gateway{
			CloudType: gatewayGroup.CloudType,
			GwName:    gwName,
		})
		if err != nil {
			return diag.Errorf("failed to delete spoke instance: %s", err)
		}
	}

	return nil
}
