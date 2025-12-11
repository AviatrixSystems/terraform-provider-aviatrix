package aviatrix

import (
	"context"
	"errors"
	"log"
	"strings"

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
			// Required attributes from group schema
			GroupRequiredSchema(),
			// Computed attributes from group schema
			GroupComputedSchema(),
			// Azure computed attributes from group schema
			GroupAzureComputedSchema(),
			// Resource-specific optional attributes
			spokeGroupOptionalSchema(),
		),
	}
}

// spokeGroupOptionalSchema returns the optional schema attributes specific to spoke group resource.
//
//nolint:funlen
func spokeGroupOptionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// ============================================================================
		// OPTIONAL ATTRIBUTES - Basic Configuration
		// ============================================================================
		"customized_cidr_list": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Set of customized CIDRs for the spoke group.",
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
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
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
		"enable_jumbo_frame": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
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
		"local_as_number": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "Changes the Aviatrix Spoke Gateway ASN number before you set up Aviatrix Spoke Gateway connection configurations.",
			ValidateFunc: goaviatrix.ValidateASN,
		},
		"prepend_as_path": {
			Type:         schema.TypeList,
			Optional:     true,
			RequiredWith: []string{"local_as_number"},
			Description:  "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.",
			Elem:         &schema.Schema{Type: schema.TypeString},
		},
		"disable_route_propagation": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable route propagation to transit.",
		},
		"spoke_bgp_manual_advertise_cidrs": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Set of CIDRs to manually advertise via BGP.",
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
		"enable_bgp_ecmp": {
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
			Description: "Number of interfaces that will be created for BGP over LAN enabled Azure spoke. " +
				"Only valid for 8 (Azure), 32 (AzureGov) or AzureChina (2048). Default value: 1. ",
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

//nolint:cyclop, funlen
func resourceAviatrixSpokeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	spokeGroup := &goaviatrix.GatewayGroup{
		GroupName:         d.Get("group_name").(string),
		CloudType:         d.Get("cloud_type").(int),
		GwType:            d.Get("gw_type").(string),
		GroupInstanceSize: d.Get("group_instance_size").(string),
		VpcID:             d.Get("vpc_id").(string),
		AccountName:       d.Get("account_name").(string),
	}

	// Optional attributes
	if v, ok := d.GetOk("customized_cidr_list"); ok {
		spokeGroup.CustomizedCidrList = getStringSet(d, "customized_cidr_list")
		_ = v
	}
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
		spokeGroup.IncludeCidr = getStringSet(d, "include_cidr")
		_ = v
	}

	spokeGroup.EnablePrivateVpcDefaultRoute = d.Get("enable_private_vpc_default_route").(bool)
	spokeGroup.EnableSkipPublicRouteTableUpdate = d.Get("enable_skip_public_route_table_update").(bool)
	spokeGroup.Edge = d.Get("edge").(bool)

	// Feature Flags
	spokeGroup.EnableJumboFrame = d.Get("enable_jumbo_frame").(bool)
	spokeGroup.EnableNat = d.Get("enable_nat").(bool)
	spokeGroup.EnableIPv6 = d.Get("enable_ipv6").(bool)
	spokeGroup.EnableGroGso = d.Get("enable_gro_gso").(bool)
	spokeGroup.EnableVpcDNSServer = d.Get("enable_vpc_dns_server").(bool)

	// BGP Configuration
	spokeGroup.EnableBgp = d.Get("enable_bgp").(bool)
	spokeGroup.DisableRoutePropagation = d.Get("disable_route_propagation").(bool)

	if v, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		spokeGroup.SpokeBgpManualAdvertiseCidrs = getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
		_ = v
	}

	spokeGroup.EnablePreserveAsPath = d.Get("enable_preserve_as_path").(bool)
	spokeGroup.EnableAutoAdvertiseS2cCidrs = d.Get("enable_auto_advertise_s2c_cidrs").(bool)
	spokeGroup.BgpEcmp = d.Get("enable_bgp_ecmp").(bool)

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

	err := client.CreateGatewayGroup(ctx, spokeGroup)
	if err != nil {
		return diag.Errorf("failed to create spoke group: %s", err)
	}

	d.SetId(spokeGroup.GroupName)
	return resourceAviatrixSpokeGroupRead(ctx, d, meta)
}

//nolint:cyclop,funlen
func resourceAviatrixSpokeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	groupName := d.Id()
	if groupName == "" {
		groupName = d.Get("group_name").(string)
	}

	log.Printf("[INFO] Reading Spoke Group: %s", groupName)

	spokeGroup, err := client.GetGatewayGroup(ctx, groupName)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
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

	if err := d.Set("customized_cidr_list", spokeGroup.CustomizedCidrList); err != nil {
		return diag.Errorf("failed to set customized_cidr_list: %s", err)
	}
	d.Set("explicitly_created", spokeGroup.ExplicitlyCreated)
	d.Set("subnet", spokeGroup.Subnet)
	d.Set("vpc_region", spokeGroup.VpcRegion)
	d.Set("domain", spokeGroup.Domain)
	if err := d.Set("include_cidr", spokeGroup.IncludeCidr); err != nil {
		return diag.Errorf("failed to set include_cidr: %s", err)
	}
	d.Set("enable_private_vpc_default_route", spokeGroup.EnablePrivateVpcDefaultRoute)
	d.Set("enable_skip_public_route_table_update", spokeGroup.EnableSkipPublicRouteTableUpdate)
	d.Set("edge", spokeGroup.Edge)

	// Feature Flags
	d.Set("enable_jumbo_frame", spokeGroup.EnableJumboFrame)
	d.Set("enable_nat", spokeGroup.EnableNat)
	d.Set("enable_ipv6", spokeGroup.EnableIPv6)
	d.Set("enable_gro_gso", spokeGroup.EnableGroGso)
	d.Set("enable_vpc_dns_server", spokeGroup.EnableVpcDNSServer)

	// BGP Configuration
	d.Set("enable_bgp", spokeGroup.EnableBgp)
	d.Set("disable_route_propagation", spokeGroup.DisableRoutePropagation)
	if err := d.Set("spoke_bgp_manual_advertise_cidrs", spokeGroup.SpokeBgpManualAdvertiseCidrs); err != nil {
		return diag.Errorf("failed to set spoke_bgp_manual_advertise_cidrs: %s", err)
	}
	d.Set("enable_preserve_as_path", spokeGroup.EnablePreserveAsPath)
	d.Set("enable_auto_advertise_s2c_cidrs", spokeGroup.EnableAutoAdvertiseS2cCidrs)
	d.Set("enable_bgp_ecmp", spokeGroup.BgpEcmp)

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
	if err := d.Set("gw_uuid_list", spokeGroup.GwUUIDList); err != nil {
		return diag.Errorf("failed to set gw_uuid_list: %s", err)
	}
	d.Set("vpc_uuid", spokeGroup.VpcUUID)
	d.Set("vendor_name", spokeGroup.VendorName)
	d.Set("software_version", spokeGroup.SoftwareVersion)
	d.Set("image_version", spokeGroup.ImageVersion)

	// Azure Computed
	d.Set("azure_eip_name_resource_group", spokeGroup.AzureEipNameResourceGroup)
	if err := d.Set("bgp_lan_ip_list", spokeGroup.BgpLanIPList); err != nil {
		return diag.Errorf("failed to set bgp_lan_ip_list: %s", err)
	}

	d.SetId(spokeGroup.GroupName)
	return nil
}

//nolint:cyclop,funlen
func resourceAviatrixSpokeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	groupName := d.Get("group_name").(string)
	spokeGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	spokeVpc := &goaviatrix.SpokeVpc{
		GwName: d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Updating Spoke Group: %s", groupName)

	// Validations for update
	if d.Get("enable_active_standby_preemptive").(bool) && !d.Get("enable_active_standby").(bool) {
		return diag.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !d.Get("enable_learned_cidrs_approval").(bool) && len(getStringSet(d, "approved_learned_cidrs")) > 0 {
		return diag.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	// ============================================================================
	// Gateway Size - API: edit_gw_config
	// ============================================================================
	if d.HasChange("group_instance_size") {
		spokeGateway.VpcSize = d.Get("group_instance_size").(string)
		err := client.UpdateGateway(spokeGateway)
		if err != nil {
			return diag.Errorf("failed to update group_instance_size: %s", err)
		}
	}

	// ============================================================================
	// BGP Communities - API: set_gateway_accept/send_bgp_communities_override
	// ============================================================================
	commSendCurr, commAcceptCurr, err := client.GetGatewayBgpCommunities(spokeGateway.GwName) // update this to use primary gateway name
	if d.HasChange("bgp_accept_communities") {
		acceptComm, ok := d.Get("bgp_accept_communities").(bool)
		if ok && acceptComm != commAcceptCurr || err != nil {
			err := client.SetGatewayBgpCommunitiesAccept(spokeGateway.GwName, acceptComm)
			if err != nil {
				return diag.Errorf("failed to update accept BGP communities for group %s: %s", spokeGateway.GwName, err)
			}
		}
	}

	if d.HasChange("bgp_send_communities") {
		sendComm, ok := d.Get("bgp_send_communities").(bool)
		if ok && sendComm != commSendCurr || err != nil {
			err := client.SetGatewayBgpCommunitiesSend(spokeGateway.GwName, sendComm)
			if err != nil {
				return diag.Errorf("failed to update send BGP communities for gateway %s: %s", spokeGateway.GwName, err)
			}
		}
	}

	// ============================================================================
	// NAT (SNAT) - API: enable_snat / disable_snat
	// ============================================================================
	if d.HasChange("enable_nat") {
		enableSNat := d.Get("enable_nat").(bool)
		if enableSNat {
			err := client.EnableSNat(spokeGateway)
			if err != nil {
				return diag.Errorf("failed to enable NAT for spoke group: %s", err)
			}
		} else {
			err := client.DisableSNat(spokeGateway)
			if err != nil {
				return diag.Errorf("failed to enable NAT for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// VPC DNS Server - API: enable_vpc_dns_server / disable_vpc_dns_server
	// ============================================================================
	if d.HasChange("enable_vpc_dns_server") {
		enableVpcDNSServer := d.Get("enable_vpc_dns_server").(bool)
		if enableVpcDNSServer {
			err := client.EnableVpcDNSServer(spokeGateway)
			if err != nil {
				return diag.Errorf("failed to enable VPC DNS Server for spoke group: %s", err)
			}
		} else {
			err := client.DisableVpcDNSServer(spokeGateway)
			if err != nil {
				return diag.Errorf("failed to disable VPC DNS Server for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// BGP Polling Time - API: change_bgp_polling_time
	// ============================================================================
	if d.HasChange("bgp_polling_time") {
		bgpPollingTime := d.Get("bgp_polling_time").(int)
		err := client.SetBgpPollingTimeSpoke(spokeVpc, bgpPollingTime)
		if err != nil {
			return diag.Errorf("could not update bgp polling time during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// BGP Neighbor Status Polling Time - API: change_bgp_neighbor_status_polling_time
	// ============================================================================
	if d.HasChange("bgp_neighbor_status_polling_time") {
		bgpBfdPollingTime := d.Get("bgp_neighbor_status_polling_time").(int)
		err := client.SetBgpBfdPollingTimeSpoke(spokeVpc, bgpBfdPollingTime)
		if err != nil {
			return diag.Errorf("could not update bgp neighbor status polling time during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// BGP Hold Time - API: change_bgp_hold_time
	// ============================================================================
	if d.HasChange("bgp_hold_time") {
		bgpHoldTime := d.Get("bgp_hold_time").(int)
		err := client.ChangeBgpHoldTime(spokeGateway.GwName, bgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// Local AS Number - API: edit_transit_local_as_number
	// Prepend AS Path - API: edit_aviatrix_transit_advanced_config
	// ============================================================================
	if d.HasChanges("local_as_number", "prepend_as_path") {
		localAsNumber := d.Get("local_as_number").(string)
		err := client.SetLocalASNumberSpoke(spokeVpc, localAsNumber)
		if err != nil {
			return diag.Errorf("could not set local_as_number for spoke group: %s", err)
		}
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err = client.SetPrependASPathSpoke(spokeVpc, prependASPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// BGP ECMP - API: enable_bgp_ecmp / disable_bgp_ecmp
	// ============================================================================
	if d.HasChange("enable_bgp_ecmp") {
		enableBgpEcmp := d.Get("bgp_ecmp").(bool)
		err := client.SetBgpEcmpSpoke(spokeVpc, enableBgpEcmp)
		if err != nil {
			return diag.Errorf("could not enable bgp ecmp during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// Active-Standby - API: enable_active_standby / disable_active_standby
	// ============================================================================
	if d.HasChange("enable_active_standby") || d.HasChange("enable_active_standby_preemptive") {
		if d.Get("enable_active_standby").(bool) {
			if d.Get("enable_active_standby_preemptive").(bool) {
				if err := client.EnableActiveStandbyPreemptiveSpoke(spokeVpc); err != nil {
					return diag.Errorf("could not enable Preemptive Mode for Active-Standby during spoke group update: %s", err)
				}
			} else {
				if err := client.EnableActiveStandbySpoke(spokeVpc); err != nil {
					return diag.Errorf("could not enable Active-Standby during spoke group update: %s", err)
				}
			}
		} else {
			if d.Get("enable_active_standby_preemptive").(bool) {
				return diag.Errorf("could not enable Preemptive Mode with Active-Standby disabled")
			}
			if err := client.DisableActiveStandbySpoke(spokeVpc); err != nil {
				return diag.Errorf("could not disable Active-Standby during Spoke Gateway update: %s", err)
			}
		}
	}

	// ============================================================================
	// Jumbo Frame - API: enable_jumbo_frame / disable_jumbo_frame
	// ============================================================================
	if d.HasChange("enable_jumbo_frame") {
		enableJumboFrame := d.Get("enable_jumbo_frame").(bool)
		if enableJumboFrame {
			err := client.EnableJumboFrame(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during spoke group update: %s", err)
			}
		} else {
			err := client.DisableJumboFrame(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// GRO/GSO - API: enable_gro_gso / disable_gro_gso
	// ============================================================================
	if d.HasChange("enable_gro_gso") {
		enableGroGso := d.Get("enable_gro_gso").(bool)
		if enableGroGso {
			err := client.EnableGroGso(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable gro gso during spoke group update: %s", err)
			}
		} else {
			err := client.DisableGroGso(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable gro gso during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// IPv6 - API: enable_ipv6 / disable_ipv6
	// ============================================================================
	if d.HasChange("enable_ipv6") {
		enableIPv6 := d.Get("enable_ipv6").(bool)
		if enableIPv6 {
			err := client.EnableIPv6(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable ipv6 during spoke group update: %s", err)
			}
		} else {
			err := client.DisableIPv6(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable ipv6 during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Preserve AS Path - API: enable_spoke_preserve_as_path / disable_spoke_preserve_as_path
	// ============================================================================
	if d.HasChange("enable_preserve_as_path") {
		enableBgp := d.Get("enable_bgp").(bool)
		enableSpokePreserveAsPath := d.Get("enable_preserve_as_path").(bool)
		if enableSpokePreserveAsPath && !enableBgp {
			return diag.Errorf("enable_preserve_as_path is not supported for Non-BGP spoke group during group update")
		}
		if !enableSpokePreserveAsPath {
			err := client.DisableSpokePreserveAsPath(spokeVpc)
			if err != nil {
				return diag.Errorf("could not disable Preserve AS Path during spoke group update: %s", err)
			}
		} else {
			err := client.EnableSpokePreserveAsPath(spokeVpc)
			if err != nil {
				return diag.Errorf("could not enable Preserve AS Path during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Learned CIDRs Approval - API: enable_bgp_gateway_cidr_approval / disable_bgp_gateway_cidr_approval
	// ============================================================================
	learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if d.HasChange("enable_learned_cidrs_approval") {
		if learnedCidrsApproval {
			spokeVpc.LearnedCidrsApproval = "on"
			err := client.EnableSpokeLearnedCidrsApproval(spokeVpc)
			if err != nil {
				return diag.Errorf("failed to enable learned cidrs approval for spoke group: %s", err)
			}
		} else {
			spokeVpc.LearnedCidrsApproval = "off"
			err := client.DisableSpokeLearnedCidrsApproval(spokeVpc)
			if err != nil {
				return diag.Errorf("failed to disable learned cidrs approval for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// Approved Learned CIDRs - API: set_bgp_gateway_approved_cidr_rules
	// ============================================================================
	if learnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		spokeVpc.ApprovedLearnedCidrs = getStringSet(d, "approved_learned_cidrs")
		err := client.UpdateSpokePendingApprovedCidrs(spokeVpc)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs: %s", err)
		}
	}

	// ============================================================================
	// Private VPC Default Route - API: enable_private_vpc_default_route / disable_private_vpc_default_route
	// ============================================================================
	if d.HasChange("enable_private_vpc_default_route") {
		enablePrivateVpcDefaultRoute := d.Get("enable_private_vpc_default_route").(bool)
		if enablePrivateVpcDefaultRoute {
			err := client.EnablePrivateVpcDefaultRoute(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable private vpc default route during spoke group update: %s", err)
			}
		} else {
			err := client.DisablePrivateVpcDefaultRoute(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable private vpc default route during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Skip Public Route Table Update - API: enable_skip_public_route_table_update / disable_skip_public_route_table_update
	// ============================================================================
	if d.HasChange("enable_skip_public_route_table_update") {
		enableSkipPublicRouteTableUpdate := d.Get("enable_skip_public_route_table_update").(bool)
		if enableSkipPublicRouteTableUpdate {
			err := client.EnableSkipPublicRouteUpdate(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable skip public route update during spoke group update: %s", err)
			}
		} else {
			err := client.DisableSkipPublicRouteUpdate(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable skip public route update during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Auto Advertise S2C CIDRs - API: enable_auto_advertise_s2c_cidrs / disable_auto_advertise_s2c_cidrs
	// ============================================================================
	if d.HasChange("enable_auto_advertise_s2c_cidrs") {
		enableAutoAdvertiseS2CCidrs := d.Get("enable_auto_advertise_s2c_cidrs").(bool)
		if enableAutoAdvertiseS2CCidrs {
			err := client.EnableAutoAdvertiseS2CCidrs(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable auto advertise s2c cidrs during spoke group update: %s", err)
			}
		} else {
			err := client.DisableAutoAdvertiseS2CCidrs(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable auto advertise s2c cidrs during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Spoke BGP Manual Advertise CIDRs - API: edit_aviatrix_spoke_advanced_config
	// ============================================================================
	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		spokeBgpManualSpokeAdvertiseCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
		spokeVpc.BgpManualSpokeAdvertiseCidrs = strings.Join(spokeBgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetSpokeBgpManualAdvertisedNetworks(spokeVpc)
		if err != nil {
			return diag.Errorf("failed to set spoke bgp manual advertise CIDRs during Spoke Gateway update: %s", err)
		}
	}

	// ============================================================================
	// Route Propagation - API: enable_spoke_onprem_route_propagation / disable_spoke_onprem_route_propagation
	// ============================================================================
	if d.HasChange("disable_route_propagation") {
		disableRoutePropagation := d.Get("disable_route_propagation").(bool)
		enableBgp := d.Get("enable_bgp").(bool)
		if disableRoutePropagation && !enableBgp {
			return diag.Errorf("disable route propagation is not supported for Non-BGP Spoke during spoke group update")
		}
		if disableRoutePropagation {
			err := client.DisableSpokeOnpremRoutePropagation(spokeVpc)
			if err != nil {
				return diag.Errorf("failed to disable route propagation for Spoke %s during spoke group update: %v", spokeVpc.GwName, err)
			}
		} else {
			err := client.EnableSpokeOnpremRoutePropagation(spokeVpc)
			if err != nil {
				return diag.Errorf("failed to enable route propagation for Spoke %s during spoke group update: %v", spokeVpc.GwName, err)
			}
		}
	}

	// ============================================================================
	// Global VPC (GCP) - API: enable_global_vpc / disable_global_vpc
	// ============================================================================
	if d.HasChange("enable_global_vpc") {
		enableGlobalVpc := d.Get("enable_global_vpc").(bool)
		if enableGlobalVpc {
			err := client.EnableGlobalVpc(spokeGateway)
			if err != nil {
				return diag.Errorf("could not enable global vpc during spoke group update: %s", err)
			}
		} else {
			err := client.DisableGlobalVpc(spokeGateway)
			if err != nil {
				return diag.Errorf("could not disable global vpc during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Encrypt Volume (AWS) - API: encrypt_gateway_volume
	// ============================================================================
	if d.HasChange("enable_encrypt_volume") {
		cloudType := d.Get("cloud_type").(int)
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			return diag.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
		}
		spokeGateway.CustomerManagedKeys = d.Get("customer_managed_keys").(string)
		err := client.EnableEncryptVolume(spokeGateway)
		if err != nil {
			return diag.Errorf("failed to enable encrypt gateway volume for %s due to %s", spokeGateway.GwName, err)
		}
	}

	if d.HasChange("customer_managed_keys") {
		return diag.Errorf("updating customer_managed_keys only is not allowed")
	}

	return resourceAviatrixSpokeGroupRead(ctx, d, meta)
}

func resourceAviatrixSpokeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	groupName := d.Id()

	log.Printf("[INFO] Deleting Spoke Group: %s", groupName)

	err := client.DeleteGatewayGroup(ctx, groupName)
	if err != nil {
		return diag.Errorf("failed to delete spoke group: %s", err)
	}

	return nil
}
