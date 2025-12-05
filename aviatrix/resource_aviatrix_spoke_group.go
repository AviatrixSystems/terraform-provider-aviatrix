package aviatrix

import (
	"context"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixSpokeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixSpokeGroupCreate,
		ReadContext:   resourceAviatrixSpokeGroupRead,
		UpdateContext: resourceAviatrixSpokeGroupUpdate,
		DeleteContext: resourceAviatrixSpokeGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: MergeSchemaMaps(
			// Required attributes from shared schema
			SpokeGroupRequiredSchema(),
			// Computed attributes from shared schema
			SpokeGroupComputedSchema(),
			// Azure computed attributes from shared schema
			SpokeGroupAzureComputedSchema(),
			// Resource-specific optional attributes
			spokeGroupOptionalSchema(),
		),
	}
}

// spokeGroupOptionalSchema returns the optional schema attributes specific to spoke group resource.
func spokeGroupOptionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// ============================================================================
		// OPTIONAL ATTRIBUTES - Basic Configuration
		// ============================================================================
		"customized_cidr_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of customized CIDRs for the spoke group.",
		},
		"s2c_rx_balancing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable S2C receive packet CPU re-balancing.",
		},
		"explicitly_created": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicates if the group was explicitly created.",
		},
		"subnet": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
			Description:  "Subnet CIDR. Required for CSP.",
		},
		"vpc_region": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Region of cloud provider. Required for CSP.",
		},
		"domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Network domain for the spoke group.",
		},
		"include_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Include CIDR for the spoke group.",
		},
		"enable_private_vpc_default_route": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable private VPC default route. Only valid for AWS.",
		},
		"enable_skip_public_route_table_update": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Skip public route table update. Only valid for AWS.",
		},
		"edge": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicates if this is an edge spoke group.",
		},

		// ============================================================================
		// FEATURE FLAGS
		// ============================================================================
		"enable_group_hpe": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable High Performance Encryption (HPE) for the group.",
		},
		"enable_jumbo_frame": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable jumbo frame support.",
		},
		"enable_nat": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable NAT (aka single_ip_snat) for the spoke group.",
		},
		"enable_ipv6": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable IPv6. Only valid for AWS and Azure.",
		},
		"enable_gro_gso": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable GRO/GSO for the spoke group.",
		},
		"enable_vpc_dns_server": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable VPC DNS server.",
		},

		// ============================================================================
		// BGP CONFIGURATION
		// ============================================================================
		"enable_bgp": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable BGP for the spoke group.",
		},
		"disable_route_propagation": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable route propagation to transit.",
		},
		"spoke_bgp_manual_advertise_cidrs": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of CIDRs to manually advertise via BGP.",
		},
		"enable_preserve_as_path": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable preserve AS path.",
		},
		"enable_auto_advertise_s2c_cidrs": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable auto advertise S2C CIDRs.",
		},
		"bgp_ecmp": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable BGP ECMP routing.",
		},

		// BGP Timers
		"bgp_polling_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      50,
			ValidateFunc: validation.IntBetween(10, 50),
			Description:  "BGP route polling time in seconds. Valid values: 10-50.",
		},
		"bgp_neighbor_status_polling_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      5,
			ValidateFunc: validation.IntBetween(1, 10),
			Description:  "BGP neighbor status polling time in seconds. Valid values: 1-10.",
		},
		"bgp_hold_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      180,
			ValidateFunc: validation.IntBetween(12, 360),
			Description:  "BGP hold time in seconds. Valid values: 12-360.",
		},

		// BGP Communities
		"bgp_send_communities": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable BGP send communities.",
		},
		"bgp_accept_communities": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable BGP accept communities.",
		},

		// BGP over LAN
		"enable_bgp_over_lan": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable BGP over LAN. Only valid for Azure.",
		},
		"bgp_lan_interfaces_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Number of BGP LAN interfaces. Only valid for Azure.",
		},

		// Learned CIDR Approval
		"enable_learned_cidrs_approval": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable learned CIDR approval.",
		},
		"learned_cidrs_approval_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "gateway",
			ValidateFunc: validation.StringInSlice([]string{"gateway"}, false),
			Description:  "Learned CIDRs approval mode. Only 'gateway' is supported for spoke.",
		},
		"approved_learned_cidrs": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Set of approved learned CIDRs.",
		},

		// Active-Standby
		"enable_active_standby": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable Active-Standby mode.",
		},
		"enable_active_standby_preemptive": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable preemptive mode for Active-Standby.",
		},

		// ============================================================================
		// AWS SPECIFIC
		// ============================================================================
		"insane_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable Insane Mode (HPE). Only valid for AWS.",
		},
		"insane_mode_az": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Availability zone for Insane Mode. Required when insane_mode is enabled for AWS.",
		},
		"enable_encrypt_volume": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable EBS volume encryption. Only valid for AWS.",
		},
		"customer_managed_keys": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Customer managed key ID for EBS volume encryption.",
		},

		// ============================================================================
		// GCP SPECIFIC
		// ============================================================================
		"enable_global_vpc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable global VPC. Only valid for GCP.",
		},
	}
}

func resourceAviatrixSpokeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	spokeGroup := &goaviatrix.SpokeGroup{
		GroupName:         d.Get("group_name").(string),
		CloudType:         d.Get("cloud_type").(int),
		GwType:            d.Get("gw_type").(string),
		GroupInstanceSize: d.Get("group_instance_size").(string),
		VpcID:             d.Get("vpc_id").(string),
		AccountName:       d.Get("account_name").(string),
	}

	// Optional attributes
	if v, ok := d.GetOk("customized_cidr_list"); ok {
		spokeGroup.CustomizedCidrList = getStringList(d, "customized_cidr_list")
		_ = v
	}
	spokeGroup.S2cRxBalancing = d.Get("s2c_rx_balancing").(bool)
	spokeGroup.ExplicitlyCreated = d.Get("explicitly_created").(bool)

	if v, ok := d.GetOk("subnet"); ok {
		spokeGroup.Subnet = v.(string)
	}
	if v, ok := d.GetOk("vpc_region"); ok {
		spokeGroup.VpcRegion = v.(string)
	}
	if v, ok := d.GetOk("domain"); ok {
		spokeGroup.Domain = v.(string)
	}
	if v, ok := d.GetOk("include_cidr"); ok {
		spokeGroup.IncludeCidr = v.(string)
	}

	spokeGroup.EnablePrivateVpcDefaultRoute = d.Get("enable_private_vpc_default_route").(bool)
	spokeGroup.EnableSkipPublicRouteTableUpdate = d.Get("enable_skip_public_route_table_update").(bool)
	spokeGroup.Edge = d.Get("edge").(bool)

	// Feature Flags
	spokeGroup.EnableGroupHpe = d.Get("enable_group_hpe").(bool)
	spokeGroup.EnableJumboFrame = d.Get("enable_jumbo_frame").(bool)
	spokeGroup.EnableNat = d.Get("enable_nat").(bool)
	spokeGroup.EnableIPv6 = d.Get("enable_ipv6").(bool)
	spokeGroup.EnableGroGso = d.Get("enable_gro_gso").(bool)
	spokeGroup.EnableVpcDnsServer = d.Get("enable_vpc_dns_server").(bool)

	// BGP Configuration
	spokeGroup.EnableBgp = d.Get("enable_bgp").(bool)
	spokeGroup.DisableRoutePropagation = d.Get("disable_route_propagation").(bool)

	if v, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		spokeGroup.SpokeBgpManualAdvertiseCidrs = getStringList(d, "spoke_bgp_manual_advertise_cidrs")
		_ = v
	}

	spokeGroup.EnablePreserveAsPath = d.Get("enable_preserve_as_path").(bool)
	spokeGroup.EnableAutoAdvertiseS2cCidrs = d.Get("enable_auto_advertise_s2c_cidrs").(bool)
	spokeGroup.BgpEcmp = d.Get("bgp_ecmp").(bool)

	// BGP Timers
	spokeGroup.BgpPollingTime = d.Get("bgp_polling_time").(int)
	spokeGroup.BgpNeighborStatusPollingTime = d.Get("bgp_neighbor_status_polling_time").(int)
	spokeGroup.BgpHoldTime = d.Get("bgp_hold_time").(int)

	// BGP Communities
	spokeGroup.BgpSendCommunities = d.Get("bgp_send_communities").(bool)
	spokeGroup.BgpAcceptCommunities = d.Get("bgp_accept_communities").(bool)

	// BGP over LAN
	spokeGroup.EnableBgpOverLan = d.Get("enable_bgp_over_lan").(bool)
	if v, ok := d.GetOk("bgp_lan_interfaces_count"); ok {
		spokeGroup.BgpLanInterfacesCount = v.(int)
	}

	// Learned CIDR Approval
	spokeGroup.EnableLearnedCidrsApproval = d.Get("enable_learned_cidrs_approval").(bool)
	if v, ok := d.GetOk("learned_cidrs_approval_mode"); ok {
		spokeGroup.LearnedCidrsApprovalMode = v.(string)
	}
	if v, ok := d.GetOk("approved_learned_cidrs"); ok {
		spokeGroup.ApprovedLearnedCidrs = getStringSet(d, "approved_learned_cidrs")
		_ = v
	}

	// Active-Standby
	spokeGroup.EnableActiveStandby = d.Get("enable_active_standby").(bool)
	spokeGroup.EnableActiveStandbyPreemptive = d.Get("enable_active_standby_preemptive").(bool)

	// AWS Specific
	spokeGroup.InsaneMode = d.Get("insane_mode").(bool)
	if v, ok := d.GetOk("insane_mode_az"); ok {
		spokeGroup.InsaneModeAz = v.(string)
	}
	spokeGroup.EnableEncryptVolume = d.Get("enable_encrypt_volume").(bool)
	if v, ok := d.GetOk("customer_managed_keys"); ok {
		spokeGroup.CustomerManagedKeys = v.(string)
	}

	// GCP Specific
	spokeGroup.EnableGlobalVpc = d.Get("enable_global_vpc").(bool)

	// Validations
	if spokeGroup.EnablePrivateVpcDefaultRoute && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return diag.Errorf("enable_private_vpc_default_route is only valid for AWS related cloud types")
	}

	if spokeGroup.EnableSkipPublicRouteTableUpdate && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return diag.Errorf("enable_skip_public_route_table_update is only valid for AWS related cloud types")
	}

	if spokeGroup.EnableIPv6 && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWS|goaviatrix.Azure) {
		return diag.Errorf("enable_ipv6 is only valid for AWS (1) and Azure (8)")
	}

	if spokeGroup.EnableBgpOverLan && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return diag.Errorf("enable_bgp_over_lan is only valid for Azure related cloud types")
	}

	if spokeGroup.InsaneMode && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return diag.Errorf("insane_mode is only valid for AWS related cloud types")
	}

	if spokeGroup.InsaneMode && goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) && spokeGroup.InsaneModeAz == "" {
		return diag.Errorf("insane_mode_az is required when insane_mode is enabled for AWS")
	}

	if spokeGroup.EnableEncryptVolume && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return diag.Errorf("enable_encrypt_volume is only valid for AWS related cloud types")
	}

	if spokeGroup.EnableGlobalVpc && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return diag.Errorf("enable_global_vpc is only valid for GCP related cloud types")
	}

	if spokeGroup.DisableRoutePropagation && !spokeGroup.EnableBgp {
		return diag.Errorf("disable_route_propagation requires enable_bgp to be true")
	}

	if spokeGroup.EnableActiveStandbyPreemptive && !spokeGroup.EnableActiveStandby {
		return diag.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !spokeGroup.EnableLearnedCidrsApproval && len(spokeGroup.ApprovedLearnedCidrs) > 0 {
		return diag.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	log.Printf("[INFO] Creating Spoke Group: %#v", spokeGroup)

	err := client.CreateSpokeGroup(ctx, spokeGroup)
	if err != nil {
		return diag.Errorf("failed to create spoke group: %s", err)
	}

	d.SetId(spokeGroup.GroupName)
	return resourceAviatrixSpokeGroupRead(ctx, d, meta)
}

func resourceAviatrixSpokeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	groupName := d.Id()
	if groupName == "" {
		groupName = d.Get("group_name").(string)
	}

	log.Printf("[INFO] Reading Spoke Group: %s", groupName)

	spokeGroup, err := client.GetSpokeGroup(ctx, groupName)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read spoke group: %s", err)
	}

	// Set required attributes
	d.Set("group_name", spokeGroup.GroupName)
	d.Set("cloud_type", spokeGroup.CloudType)
	d.Set("gw_type", spokeGroup.GwType)
	d.Set("group_instance_size", spokeGroup.GroupInstanceSize)
	d.Set("vpc_id", spokeGroup.VpcID)
	d.Set("account_name", spokeGroup.AccountName)

	// Set optional attributes
	if err := d.Set("customized_cidr_list", spokeGroup.CustomizedCidrList); err != nil {
		return diag.Errorf("failed to set customized_cidr_list: %s", err)
	}
	d.Set("s2c_rx_balancing", spokeGroup.S2cRxBalancing)
	d.Set("explicitly_created", spokeGroup.ExplicitlyCreated)
	d.Set("subnet", spokeGroup.Subnet)
	d.Set("vpc_region", spokeGroup.VpcRegion)
	d.Set("domain", spokeGroup.Domain)
	d.Set("include_cidr", spokeGroup.IncludeCidr)
	d.Set("enable_private_vpc_default_route", spokeGroup.EnablePrivateVpcDefaultRoute)
	d.Set("enable_skip_public_route_table_update", spokeGroup.EnableSkipPublicRouteTableUpdate)
	d.Set("edge", spokeGroup.Edge)

	// Feature Flags
	d.Set("enable_group_hpe", spokeGroup.EnableGroupHpe)
	d.Set("enable_jumbo_frame", spokeGroup.EnableJumboFrame)
	d.Set("enable_nat", spokeGroup.EnableNat)
	d.Set("enable_ipv6", spokeGroup.EnableIPv6)
	d.Set("enable_gro_gso", spokeGroup.EnableGroGso)
	d.Set("enable_vpc_dns_server", spokeGroup.EnableVpcDnsServer)

	// BGP Configuration
	d.Set("enable_bgp", spokeGroup.EnableBgp)
	d.Set("disable_route_propagation", spokeGroup.DisableRoutePropagation)
	if err := d.Set("spoke_bgp_manual_advertise_cidrs", spokeGroup.SpokeBgpManualAdvertiseCidrs); err != nil {
		return diag.Errorf("failed to set spoke_bgp_manual_advertise_cidrs: %s", err)
	}
	d.Set("enable_preserve_as_path", spokeGroup.EnablePreserveAsPath)
	d.Set("enable_auto_advertise_s2c_cidrs", spokeGroup.EnableAutoAdvertiseS2cCidrs)
	d.Set("bgp_ecmp", spokeGroup.BgpEcmp)

	// BGP Timers
	d.Set("bgp_polling_time", spokeGroup.BgpPollingTime)
	d.Set("bgp_neighbor_status_polling_time", spokeGroup.BgpNeighborStatusPollingTime)
	d.Set("bgp_hold_time", spokeGroup.BgpHoldTime)

	// BGP Communities
	d.Set("bgp_send_communities", spokeGroup.BgpSendCommunities)
	d.Set("bgp_accept_communities", spokeGroup.BgpAcceptCommunities)

	// BGP over LAN
	d.Set("enable_bgp_over_lan", spokeGroup.EnableBgpOverLan)
	d.Set("bgp_lan_interfaces_count", spokeGroup.BgpLanInterfacesCount)

	// Learned CIDR Approval
	d.Set("enable_learned_cidrs_approval", spokeGroup.EnableLearnedCidrsApproval)
	d.Set("learned_cidrs_approval_mode", spokeGroup.LearnedCidrsApprovalMode)
	if err := d.Set("approved_learned_cidrs", spokeGroup.ApprovedLearnedCidrs); err != nil {
		return diag.Errorf("failed to set approved_learned_cidrs: %s", err)
	}

	// Active-Standby
	d.Set("enable_active_standby", spokeGroup.EnableActiveStandby)
	d.Set("enable_active_standby_preemptive", spokeGroup.EnableActiveStandbyPreemptive)

	// AWS Specific
	d.Set("insane_mode", spokeGroup.InsaneMode)
	d.Set("insane_mode_az", spokeGroup.InsaneModeAz)
	d.Set("enable_encrypt_volume", spokeGroup.EnableEncryptVolume)
	// Note: customer_managed_keys is sensitive, only set if already in state
	if _, ok := d.GetOk("customer_managed_keys"); ok {
		d.Set("customer_managed_keys", spokeGroup.CustomerManagedKeys)
	}

	// GCP Specific
	d.Set("enable_global_vpc", spokeGroup.EnableGlobalVpc)

	// Computed attributes
	if err := d.Set("gw_uuid_list", spokeGroup.GwUuidList); err != nil {
		return diag.Errorf("failed to set gw_uuid_list: %s", err)
	}
	d.Set("vpc_uuid", spokeGroup.VpcUuid)
	d.Set("vendor_name", spokeGroup.VendorName)
	d.Set("software_version", spokeGroup.SoftwareVersion)
	d.Set("image_version", spokeGroup.ImageVersion)

	// Azure Computed
	d.Set("azure_eip_name_resource_group", spokeGroup.AzureEipNameResourceGroup)
	if err := d.Set("bgp_lan_ip_list", spokeGroup.BgpLanIpList); err != nil {
		return diag.Errorf("failed to set bgp_lan_ip_list: %s", err)
	}

	d.SetId(spokeGroup.GroupName)
	return nil
}

func resourceAviatrixSpokeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Spoke Group: %s", d.Id())

	spokeGroup := &goaviatrix.SpokeGroup{
		GroupName:         d.Get("group_name").(string),
		CloudType:         d.Get("cloud_type").(int),
		GwType:            d.Get("gw_type").(string),
		GroupInstanceSize: d.Get("group_instance_size").(string),
		VpcID:             d.Get("vpc_id").(string),
		AccountName:       d.Get("account_name").(string),
	}

	// Build the full spokeGroup object with current values
	if v, ok := d.GetOk("customized_cidr_list"); ok {
		spokeGroup.CustomizedCidrList = getStringList(d, "customized_cidr_list")
		_ = v
	}
	spokeGroup.S2cRxBalancing = d.Get("s2c_rx_balancing").(bool)
	spokeGroup.ExplicitlyCreated = d.Get("explicitly_created").(bool)

	if v, ok := d.GetOk("subnet"); ok {
		spokeGroup.Subnet = v.(string)
	}
	if v, ok := d.GetOk("vpc_region"); ok {
		spokeGroup.VpcRegion = v.(string)
	}
	if v, ok := d.GetOk("domain"); ok {
		spokeGroup.Domain = v.(string)
	}
	if v, ok := d.GetOk("include_cidr"); ok {
		spokeGroup.IncludeCidr = v.(string)
	}

	spokeGroup.EnablePrivateVpcDefaultRoute = d.Get("enable_private_vpc_default_route").(bool)
	spokeGroup.EnableSkipPublicRouteTableUpdate = d.Get("enable_skip_public_route_table_update").(bool)
	spokeGroup.Edge = d.Get("edge").(bool)

	// Feature Flags
	spokeGroup.EnableGroupHpe = d.Get("enable_group_hpe").(bool)
	spokeGroup.EnableJumboFrame = d.Get("enable_jumbo_frame").(bool)
	spokeGroup.EnableNat = d.Get("enable_nat").(bool)
	spokeGroup.EnableIPv6 = d.Get("enable_ipv6").(bool)
	spokeGroup.EnableGroGso = d.Get("enable_gro_gso").(bool)
	spokeGroup.EnableVpcDnsServer = d.Get("enable_vpc_dns_server").(bool)

	// BGP Configuration
	spokeGroup.EnableBgp = d.Get("enable_bgp").(bool)
	spokeGroup.DisableRoutePropagation = d.Get("disable_route_propagation").(bool)

	if v, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		spokeGroup.SpokeBgpManualAdvertiseCidrs = getStringList(d, "spoke_bgp_manual_advertise_cidrs")
		_ = v
	}

	spokeGroup.EnablePreserveAsPath = d.Get("enable_preserve_as_path").(bool)
	spokeGroup.EnableAutoAdvertiseS2cCidrs = d.Get("enable_auto_advertise_s2c_cidrs").(bool)
	spokeGroup.BgpEcmp = d.Get("bgp_ecmp").(bool)

	// BGP Timers
	spokeGroup.BgpPollingTime = d.Get("bgp_polling_time").(int)
	spokeGroup.BgpNeighborStatusPollingTime = d.Get("bgp_neighbor_status_polling_time").(int)
	spokeGroup.BgpHoldTime = d.Get("bgp_hold_time").(int)

	// BGP Communities
	spokeGroup.BgpSendCommunities = d.Get("bgp_send_communities").(bool)
	spokeGroup.BgpAcceptCommunities = d.Get("bgp_accept_communities").(bool)

	// BGP over LAN
	spokeGroup.EnableBgpOverLan = d.Get("enable_bgp_over_lan").(bool)
	if v, ok := d.GetOk("bgp_lan_interfaces_count"); ok {
		spokeGroup.BgpLanInterfacesCount = v.(int)
	}

	// Learned CIDR Approval
	spokeGroup.EnableLearnedCidrsApproval = d.Get("enable_learned_cidrs_approval").(bool)
	if v, ok := d.GetOk("learned_cidrs_approval_mode"); ok {
		spokeGroup.LearnedCidrsApprovalMode = v.(string)
	}
	if v, ok := d.GetOk("approved_learned_cidrs"); ok {
		spokeGroup.ApprovedLearnedCidrs = getStringSet(d, "approved_learned_cidrs")
		_ = v
	}

	// Active-Standby
	spokeGroup.EnableActiveStandby = d.Get("enable_active_standby").(bool)
	spokeGroup.EnableActiveStandbyPreemptive = d.Get("enable_active_standby_preemptive").(bool)

	// AWS Specific
	spokeGroup.InsaneMode = d.Get("insane_mode").(bool)
	if v, ok := d.GetOk("insane_mode_az"); ok {
		spokeGroup.InsaneModeAz = v.(string)
	}
	spokeGroup.EnableEncryptVolume = d.Get("enable_encrypt_volume").(bool)
	if v, ok := d.GetOk("customer_managed_keys"); ok {
		spokeGroup.CustomerManagedKeys = v.(string)
	}

	// GCP Specific
	spokeGroup.EnableGlobalVpc = d.Get("enable_global_vpc").(bool)

	// Validations for update
	if spokeGroup.EnableActiveStandbyPreemptive && !spokeGroup.EnableActiveStandby {
		return diag.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !spokeGroup.EnableLearnedCidrsApproval && len(spokeGroup.ApprovedLearnedCidrs) > 0 {
		return diag.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	err := client.UpdateSpokeGroup(ctx, spokeGroup)
	if err != nil {
		return diag.Errorf("failed to update spoke group: %s", err)
	}

	return resourceAviatrixSpokeGroupRead(ctx, d, meta)
}

func resourceAviatrixSpokeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	groupName := d.Id()

	log.Printf("[INFO] Deleting Spoke Group: %s", groupName)

	err := client.DeleteSpokeGroup(ctx, groupName)
	if err != nil {
		return diag.Errorf("failed to delete spoke group: %s", err)
	}

	return nil
}
