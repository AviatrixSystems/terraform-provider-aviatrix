package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
			// Resource-specific required attributes
			spokeGroupRequiredSchema(),
			// Computed attributes from group schema
			GroupComputedSchema(),
			// Optional attributes from group schema
			GroupOptionalSchema(),
			// Resource-specific optional attributes
			spokeGroupOptionalSchema(),
		),
	}
}

// spokeGroupRequiredSchema returns the required schema attributes specific to spoke group resource.
func spokeGroupRequiredSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gw_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Gateway type for the group. Valid values: SPOKE, EDGESPOKE, STANDALONE. Case-insensitive.",
			ValidateFunc: validateSpokeGwType,
			StateFunc:    normalizeGwType,
		},
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

		// ============================================================================
		// FEATURE FLAGS
		// ============================================================================
		"enable_jumbo_frame": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
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
			Description:  "Changes the Aviatrix Spoke Gateway ASN number before you set up Aviatrix Spoke Gateway connection configurations.",
			ValidateFunc: goaviatrix.ValidateASN,
		},
		"prepend_as_path": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: goaviatrix.ValidateASN,
			},
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
			Default:      defaultBgpPollingTime,
			ValidateFunc: validation.IntBetween(10, 50),
			Description:  "BGP route polling time in seconds. Valid values: 10-50.",
		},
		"bgp_neighbor_status_polling_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      defaultBgpNeighborStatusPollingTime,
			ValidateFunc: validation.IntBetween(1, 10),
			Description:  "BGP neighbor status polling time in seconds. Valid values: 1-10.",
		},
		"bgp_hold_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      defaultBgpHoldTime,
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
			Description: "Enable preemptive mode for Active-Standby, available only with Active-Standby enabled.",
		},

		// ============================================================================
		// AWS SPECIFIC
		// ============================================================================
		"insane_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable Insane Mode for spoke gateway group. Valid values: true, false. Supported for AWS/AWSGov, GCP, Azure and OCI. If insane mode is enabled, gateway size has to at least be c5 size for AWS and Standard_D3_v2 size for Azure.",
		},

		// ============================================================================
		// GCP SPECIFIC
		// ============================================================================
		"enable_global_vpc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set to true to enable global VPC. Only supported for GCP",
		},
	}
}

// ============================================================================
// Spoke Group Create Helper Functions
// ============================================================================

// buildSpokeGroupFromResourceData constructs a GatewayGroup struct from Terraform resource data.
func buildSpokeGroupFromResourceData(d *schema.ResourceData) *goaviatrix.GatewayGroup {
	spokeGroup := &goaviatrix.GatewayGroup{
		GroupName:         getString(d, "group_name"),
		CloudType:         getInt(d, "cloud_type"),
		GwType:            getString(d, "gw_type"),
		GroupInstanceSize: getString(d, "group_instance_size"),
		VpcID:             getString(d, "vpc_id"),
		AccountName:       getString(d, "account_name"),
	}

	// Optional attributes
	if _, ok := d.GetOk("customized_cidr_list"); ok {
		spokeGroup.CustomizedCidrList = getStringSet(d, "customized_cidr_list")
	}
	if v, ok := d.GetOk("vpc_region"); ok {
		spokeGroup.VpcRegion = mustString(v)
	}
	if v, ok := d.GetOk("domain"); ok {
		spokeGroup.Domain = mustString(v)
	}
	if _, ok := d.GetOk("include_cidr"); ok {
		spokeGroup.IncludeCidr = getStringSet(d, "include_cidr")
	}

	spokeGroup.EnablePrivateVpcDefaultRoute = getBool(d, "enable_private_vpc_default_route")
	spokeGroup.EnableSkipPublicRouteTableUpdate = getBool(d, "enable_skip_public_route_table_update")

	// Feature Flags
	spokeGroup.EnableNat = getBool(d, "enable_nat")
	spokeGroup.EnableIPv6 = getBool(d, "enable_ipv6")
	spokeGroup.EnableVpcDNSServer = getBool(d, "enable_vpc_dns_server")

	// BGP Configuration
	spokeGroup.EnableBgp = getBool(d, "enable_bgp")
	spokeGroup.DisableRoutePropagation = getBool(d, "disable_route_propagation")

	if _, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		spokeGroup.SpokeBgpManualAdvertiseCidrs = getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	}

	spokeGroup.EnablePreserveAsPath = getBool(d, "enable_preserve_as_path")
	spokeGroup.EnableAutoAdvertiseS2cCidrs = getBool(d, "enable_auto_advertise_s2c_cidrs")
	spokeGroup.BgpEcmp = getBool(d, "enable_bgp_ecmp")

	// BGP Timers
	spokeGroup.BgpPollingTime = getInt(d, "bgp_polling_time")
	spokeGroup.BgpNeighborStatusPollingTime = getInt(d, "bgp_neighbor_status_polling_time")
	spokeGroup.BgpHoldTime = getInt(d, "bgp_hold_time")

	// BGP Communities
	spokeGroup.BgpSendCommunities = getBool(d, "bgp_send_communities")
	spokeGroup.BgpAcceptCommunities = getBool(d, "bgp_accept_communities")

	// BGP over LAN
	spokeGroup.EnableBgpOverLan = getBool(d, "enable_bgp_over_lan")

	// Learned CIDR Approval
	spokeGroup.EnableLearnedCidrsApproval = getBool(d, "enable_learned_cidrs_approval")
	if v, ok := d.GetOk("learned_cidrs_approval_mode"); ok {
		spokeGroup.LearnedCidrsApprovalMode = mustString(v)
	}
	if _, ok := d.GetOk("approved_learned_cidrs"); ok {
		spokeGroup.ApprovedLearnedCidrs = getStringSet(d, "approved_learned_cidrs")
	}

	// Active-Standby
	spokeGroup.EnableActiveStandby = getBool(d, "enable_active_standby")
	spokeGroup.EnableActiveStandbyPreemptive = getBool(d, "enable_active_standby_preemptive")

	// AWS Specific
	spokeGroup.InsaneMode = getBool(d, "insane_mode")

	// GCP Specific
	spokeGroup.EnableGlobalVpc = getBool(d, "enable_global_vpc")

	return spokeGroup
}

// validateSpokeGroupConfiguration validates the spoke group configuration for cloud-type and feature dependencies.
func validateSpokeGroupConfiguration(spokeGroup *goaviatrix.GatewayGroup) error {
	if spokeGroup.EnablePrivateVpcDefaultRoute && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_private_vpc_default_route is only valid for AWS related cloud types")
	}

	if spokeGroup.EnableSkipPublicRouteTableUpdate && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_skip_public_route_table_update is only valid for AWS related cloud types")
	}

	if spokeGroup.EnableIPv6 && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWS|goaviatrix.Azure) {
		return fmt.Errorf("enable_ipv6 is only valid for AWS (1) and Azure (8)")
	}

	if spokeGroup.EnableBgpOverLan && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("enable_bgp_over_lan is only valid for Azure related cloud types")
	}

	if spokeGroup.InsaneMode && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("insane_mode is only valid for AWS related cloud types")
	}

	if spokeGroup.EnableGlobalVpc && !goaviatrix.IsCloudType(spokeGroup.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("enable_global_vpc is only valid for GCP related cloud types")
	}

	if !spokeGroup.EnableBgp && spokeGroup.LocalAsNumber != "" {
		return fmt.Errorf("local_as_number can only be set when enable_bgp is true")
	}

	if len(spokeGroup.PrependAsPath) > 0 && spokeGroup.LocalAsNumber == "" {
		return fmt.Errorf("prepend_as_path can only be set when local_as_number is set")
	}

	if spokeGroup.DisableRoutePropagation && !spokeGroup.EnableBgp {
		return fmt.Errorf("disable_route_propagation requires enable_bgp to be true")
	}

	if spokeGroup.EnableActiveStandbyPreemptive && !spokeGroup.EnableActiveStandby {
		return fmt.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !spokeGroup.EnableLearnedCidrsApproval && len(spokeGroup.ApprovedLearnedCidrs) > 0 {
		return fmt.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	return nil
}

// applyBgpCommunities applies BGP communities settings (accept/send) to the spoke group.
func applyBgpCommunities(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	// Default for both is false, so only set if true
	acceptComm := getBool(d, "bgp_accept_communities")
	sendComm := getBool(d, "bgp_send_communities")

	if acceptComm {
		if err := client.SetGatewayGroupBgpCommunitiesAccept(ctx, groupName, acceptComm); err != nil {
			return fmt.Errorf("failed to set accept BGP communities: %w", err)
		}
	}

	if sendComm {
		if err := client.SetGatewayGroupBgpCommunitiesSend(ctx, groupName, sendComm); err != nil {
			return fmt.Errorf("failed to set send BGP communities: %w", err)
		}
	}

	return nil
}

// applyFeatureFlags applies feature flags (jumbo frame, GRO/GSO, NAT, VPC DNS Server, IPv6) to the spoke group.
func applyFeatureFlags(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	if getBool(d, "enable_nat") {
		log.Printf("[INFO] Enabling NAT for spoke group: %s", groupName)
		if err := client.EnableGatewayGroupSNat(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable NAT: %w", err)
		}
	}

	if getBool(d, "enable_vpc_dns_server") {
		log.Printf("[INFO] Enabling VPC DNS Server for spoke group: %s", groupName)
		if err := client.EnableGatewayGroupVpcDNSServer(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable VPC DNS Server: %w", err)
		}
	}

	if getBool(d, "enable_ipv6") {
		log.Printf("[INFO] Enabling IPv6 for spoke group: %s", groupName)
		if err := client.EnableIPv6GatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable IPv6: %w", err)
		}
	}

	return nil
}

// applyBgpTimers applies BGP timer settings (polling time, neighbor status polling, hold time) to the spoke group.
func applyBgpTimers(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	bgpPollingTime := getInt(d, "bgp_polling_time")
	if bgpPollingTime != defaultBgpPollingTime {
		if err := client.SetBgpPollingTimeGatewayGroup(ctx, groupName, bgpPollingTime); err != nil {
			return fmt.Errorf("failed to set BGP polling time: %w", err)
		}
	}

	bgpNeighborStatusPollingTime := getInt(d, "bgp_neighbor_status_polling_time")
	if bgpNeighborStatusPollingTime != defaultBgpNeighborStatusPollingTime {
		if err := client.SetBgpBfdPollingTimeGatewayGroup(ctx, groupName, bgpNeighborStatusPollingTime); err != nil {
			return fmt.Errorf("failed to set BGP neighbor status polling time: %w", err)
		}
	}

	bgpHoldTime := getInt(d, "bgp_hold_time")
	if bgpHoldTime != defaultBgpHoldTime {
		if err := client.ChangeBgpHoldTimeGatewayGroup(ctx, groupName, bgpHoldTime); err != nil {
			return fmt.Errorf("failed to set BGP hold time: %w", err)
		}
	}

	return nil
}

// applyBgpConfiguration applies BGP configuration (local AS number, prepend AS path, BGP ECMP) to the spoke group.
func applyBgpConfiguration(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	if val, ok := d.GetOk("local_as_number"); ok {
		if err := client.SetLocalASNumberGatewayGroup(ctx, groupName, mustString(val)); err != nil {
			return fmt.Errorf("failed to set local AS number: %w", err)
		}

		if _, ok := d.GetOk("prepend_as_path"); ok {
			prependPath := getStringList(d, "prepend_as_path")
			if len(prependPath) > 0 {
				if err := client.SetPrependASPathGatewayGroup(ctx, groupName, prependPath); err != nil {
					return fmt.Errorf("failed to set prepend AS path: %w", err)
				}
			}
		}
	}

	if val, ok := d.GetOk("enable_bgp_ecmp"); ok {
		if err := client.SetBgpEcmpGatewayGroup(ctx, groupName, mustBool(val)); err != nil {
			return fmt.Errorf("failed to set BGP ECMP: %w", err)
		}
	}

	return nil
}

// applyActiveStandby applies Active-Standby settings to the spoke group.
func applyActiveStandby(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	enableActiveStandby := getBool(d, "enable_active_standby")
	if !enableActiveStandby {
		return nil
	}

	enableActiveStandbyPreemptive := getBool(d, "enable_active_standby_preemptive")
	if enableActiveStandbyPreemptive {
		if err := client.EnableActiveStandbyPreemptiveGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable Active-Standby Preemptive mode: %w", err)
		}
	} else {
		if err := client.EnableActiveStandbyGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable Active-Standby: %w", err)
		}
	}

	return nil
}

// applySpokeSpecificSettings applies spoke-specific settings (preserve AS path, learned CIDRs, route propagation, etc.) to the spoke group.
func applySpokeSpecificSettings(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	if getBool(d, "enable_preserve_as_path") {
		log.Printf("[INFO] Enabling Preserve AS Path for spoke group: %s", groupName)
		if err := client.EnableSpokePreserveAsPathGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable Preserve AS Path: %w", err)
		}
	}

	if getBool(d, "enable_learned_cidrs_approval") {
		log.Printf("[INFO] Enabling learned CIDRs approval for spoke group: %s", groupName)
		if err := client.EnableSpokeLearnedCidrsApprovalGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable learned CIDRs approval: %w", err)
		}
	}

	approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")
	if len(approvedLearnedCidrs) != 0 {
		if err := client.UpdateSpokePendingApprovedCidrsGatewayGroup(ctx, groupName, approvedLearnedCidrs); err != nil {
			return fmt.Errorf("failed to update approved learned CIDRs: %w", err)
		}
	}

	if getBool(d, "enable_private_vpc_default_route") {
		log.Printf("[INFO] Enabling private VPC default route for spoke group: %s", groupName)
		if err := client.EnablePrivateVpcDefaultRouteGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable private VPC default route: %w", err)
		}
	}

	if getBool(d, "enable_skip_public_route_table_update") {
		log.Printf("[INFO] Enabling skip public route update for spoke group: %s", groupName)
		if err := client.EnableSkipPublicRouteUpdateGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable skip public route update: %w", err)
		}
	}

	if getBool(d, "enable_auto_advertise_s2c_cidrs") {
		log.Printf("[INFO] Enabling auto advertise s2c cidrs for spoke group: %s", groupName)
		if err := client.EnableAutoAdvertiseS2CCidrsGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable auto advertise s2c CIDRs: %w", err)
		}
	}

	if _, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		spokeBgpManualAdvertiseCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
		cidrs := strings.Join(spokeBgpManualAdvertiseCidrs, ",")

		log.Printf("[INFO] Setting spoke BGP manual advertise CIDRs for spoke group: %s", groupName)
		if err := client.SetSpokeBgpManualAdvertisedNetworksGatewayGroup(ctx, groupName, cidrs); err != nil {
			return fmt.Errorf("failed to set spoke BGP manual advertise CIDRs: %w", err)
		}
	}

	if getBool(d, "disable_route_propagation") {
		log.Printf("[INFO] Disabling route propagation for spoke group: %s", groupName)
		if err := client.DisableSpokeOnpremRoutePropagationGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to disable route propagation: %w", err)
		}
	}

	if getBool(d, "enable_global_vpc") {
		log.Printf("[INFO] Enabling global VPC for spoke group: %s", groupName)
		if err := client.EnableGlobalVpcGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable global VPC: %w", err)
		}
	}

	return nil
}

// ============================================================================
// Spoke Group CRUD Operations
// ============================================================================

func resourceAviatrixSpokeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// Build the spoke group from resource data
	spokeGroup := buildSpokeGroupFromResourceData(d)

	// Validate configuration
	if err := validateSpokeGroupConfiguration(spokeGroup); err != nil {
		return diag.FromErr(err)
	}

	// Create the spoke group
	log.Printf("[INFO] Creating Spoke Group: %#v", spokeGroup)
	if err := client.CreateGatewayGroup(ctx, spokeGroup); err != nil {
		return diag.Errorf("failed to create spoke group: %s", err)
	}

	// Use GroupUUID as the resource ID
	d.SetId(spokeGroup.GroupUUID)
	groupName := getString(d, "group_name")

	// Apply post-creation settings
	if err := applyBgpCommunities(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply BGP communities: %s", err)
	}

	if err := applyFeatureFlags(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply feature flags: %s", err)
	}

	if err := applyBgpTimers(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply BGP timers: %s", err)
	}

	if err := applyBgpConfiguration(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply BGP configuration: %s", err)
	}

	if err := applyActiveStandby(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply Active-Standby: %s", err)
	}

	if err := applySpokeSpecificSettings(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply spoke-specific settings: %s", err)
	}

	return resourceAviatrixSpokeGroupRead(ctx, d, meta)
}

//nolint:cyclop,funlen
func resourceAviatrixSpokeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// The resource ID is the group UUID
	groupUUID := d.Id()
	if groupUUID == "" {
		return diag.Errorf("resource ID (group UUID) is empty")
	}

	log.Printf("[INFO] Reading Spoke Group: %s", groupUUID)

	spokeGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read spoke group: %s", err)
	}

	// Set required attributes
	mustSet(d, "group_name", spokeGroup.GroupName)
	mustSet(d, "cloud_type", spokeGroup.CloudType)
	// Normalize gw_type by removing "GwGroupType." prefix if present
	gwType := strings.TrimPrefix(spokeGroup.GwType, "GwGroupType.")
	mustSet(d, "gw_type", gwType)
	mustSet(d, "group_instance_size", spokeGroup.GroupInstanceSize)
	mustSet(d, "vpc_id", spokeGroup.VpcID)
	mustSet(d, "account_name", spokeGroup.AccountName)

	// Store group_uuid in state (used as resource ID)
	if spokeGroup.GroupUUID != "" {
		mustSet(d, "group_uuid", spokeGroup.GroupUUID)
	}

	mustSet(d, "customized_cidr_list", spokeGroup.CustomizedCidrList)
	mustSet(d, "explicitly_created", spokeGroup.ExplicitlyCreated)
	mustSet(d, "vpc_region", spokeGroup.VpcRegion)
	mustSet(d, "domain", spokeGroup.Domain)
	mustSet(d, "include_cidr", spokeGroup.IncludeCidr)
	mustSet(d, "enable_private_vpc_default_route", spokeGroup.EnablePrivateVpcDefaultRoute)
	mustSet(d, "enable_skip_public_route_table_update", spokeGroup.EnableSkipPublicRouteTableUpdate)

	// Feature Flags
	mustSet(d, "enable_jumbo_frame", spokeGroup.EnableJumboFrame)
	mustSet(d, "enable_nat", spokeGroup.EnableNat)
	mustSet(d, "enable_ipv6", spokeGroup.EnableIPv6)
	mustSet(d, "enable_gro_gso", spokeGroup.EnableGroGso)
	mustSet(d, "enable_vpc_dns_server", spokeGroup.EnableVpcDNSServer)

	// BGP Configuration
	mustSet(d, "enable_bgp", spokeGroup.EnableBgp)
	mustSet(d, "local_as_number", spokeGroup.LocalAsNumber)
	mustSet(d, "prepend_as_path", spokeGroup.PrependAsPath)
	mustSet(d, "disable_route_propagation", spokeGroup.DisableRoutePropagation)
	mustSet(d, "spoke_bgp_manual_advertise_cidrs", spokeGroup.SpokeBgpManualAdvertiseCidrs)
	mustSet(d, "enable_preserve_as_path", spokeGroup.EnablePreserveAsPath)
	mustSet(d, "enable_auto_advertise_s2c_cidrs", spokeGroup.EnableAutoAdvertiseS2cCidrs)
	mustSet(d, "enable_bgp_ecmp", spokeGroup.BgpEcmp)

	// BGP Timers - Set to schema defaults if API returns 0
	if spokeGroup.BgpPollingTime == 0 {
		mustSet(d, "bgp_polling_time", defaultBgpPollingTime)
	} else {
		mustSet(d, "bgp_polling_time", spokeGroup.BgpPollingTime)
	}
	if spokeGroup.BgpNeighborStatusPollingTime == 0 {
		mustSet(d, "bgp_neighbor_status_polling_time", defaultBgpNeighborStatusPollingTime)
	} else {
		mustSet(d, "bgp_neighbor_status_polling_time", spokeGroup.BgpNeighborStatusPollingTime)
	}
	if spokeGroup.BgpHoldTime == 0 {
		mustSet(d, "bgp_hold_time", defaultBgpHoldTime)
	} else {
		mustSet(d, "bgp_hold_time", spokeGroup.BgpHoldTime)
	}

	// BGP Communities
	mustSet(d, "bgp_send_communities", spokeGroup.BgpSendCommunities)
	mustSet(d, "bgp_accept_communities", spokeGroup.BgpAcceptCommunities)

	// BGP over LAN
	mustSet(d, "enable_bgp_over_lan", spokeGroup.EnableBgpOverLan)

	// Learned CIDR Approval
	mustSet(d, "enable_learned_cidrs_approval", spokeGroup.EnableLearnedCidrsApproval)
	// Set to schema default if API returns empty string
	if spokeGroup.LearnedCidrsApprovalMode == "" {
		mustSet(d, "learned_cidrs_approval_mode", "gateway")
	} else {
		mustSet(d, "learned_cidrs_approval_mode", spokeGroup.LearnedCidrsApprovalMode)
	}
	mustSet(d, "approved_learned_cidrs", spokeGroup.ApprovedLearnedCidrs)

	// Active-Standby
	mustSet(d, "enable_active_standby", spokeGroup.EnableActiveStandby)
	mustSet(d, "enable_active_standby_preemptive", spokeGroup.EnableActiveStandbyPreemptive)

	// AWS Specific
	mustSet(d, "insane_mode", spokeGroup.InsaneMode)

	// GCP Specific
	mustSet(d, "enable_global_vpc", spokeGroup.EnableGlobalVpc)

	// Computed attributes
	mustSet(d, "gw_uuid_list", spokeGroup.GwUUIDList)
	mustSet(d, "vpc_uuid", spokeGroup.VpcUUID)

	// Use GroupUUID as the resource ID
	if spokeGroup.GroupUUID != "" {
		d.SetId(spokeGroup.GroupUUID)
	}
	return nil
}

//nolint:cyclop,funlen
func resourceAviatrixSpokeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	groupName := getString(d, "group_name")
	cloudType := getInt(d, "cloud_type")
	groupUUID := d.Id()

	log.Printf("[INFO] Updating Spoke Group: %s (UUID: %s)", groupName, groupUUID)

	// Validations for update
	if getBool(d, "enable_active_standby_preemptive") && !getBool(d, "enable_active_standby") {
		return diag.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !getBool(d, "enable_learned_cidrs_approval") && len(getStringSet(d, "approved_learned_cidrs")) > 0 {
		return diag.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	// ============================================================================
	// Gateway Size - API: edit_gw_config
	// ============================================================================
	if d.HasChange("group_instance_size") {
		instanceSize := getString(d, "group_instance_size")
		err := client.UpdateGatewayGroupSize(ctx, groupName, instanceSize)
		if err != nil {
			return diag.Errorf("failed to update group_instance_size: %s", err)
		}
	}

	// ============================================================================
	// BGP Communities - API: set_gateway_accept/send_bgp_communities_override
	// ============================================================================
	if d.HasChange("bgp_accept_communities") {
		acceptComm := getBool(d, "bgp_accept_communities")
		err := client.SetGatewayGroupBgpCommunitiesAccept(ctx, groupName, acceptComm)
		if err != nil {
			return diag.Errorf("failed to update accept BGP communities for group %s: %s", groupName, err)
		}
	}

	if d.HasChange("bgp_send_communities") {
		sendComm := getBool(d, "bgp_send_communities")
		err := client.SetGatewayGroupBgpCommunitiesSend(ctx, groupName, sendComm)
		if err != nil {
			return diag.Errorf("failed to update send BGP communities for group %s: %s", groupName, err)
		}
	}

	// ============================================================================
	// NAT (SNAT) - API: enable_snat / disable_snat
	// ============================================================================
	if d.HasChange("enable_nat") {
		enableSNat := getBool(d, "enable_nat")
		if enableSNat {
			err := client.EnableGatewayGroupSNat(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable NAT for spoke group: %s", err)
			}
		} else {
			err := client.DisableGatewayGroupSNat(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable NAT for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// VPC DNS Server - API: enable_vpc_dns_server / disable_vpc_dns_server
	// ============================================================================
	if d.HasChange("enable_vpc_dns_server") {
		enableVpcDNSServer := getBool(d, "enable_vpc_dns_server")
		if enableVpcDNSServer {
			err := client.EnableGatewayGroupVpcDNSServer(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable VPC DNS Server for spoke group: %s", err)
			}
		} else {
			err := client.DisableGatewayGroupVpcDNSServer(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable VPC DNS Server for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// BGP Polling Time - API: change_bgp_polling_time
	// ============================================================================
	if d.HasChange("bgp_polling_time") {
		bgpPollingTime := getInt(d, "bgp_polling_time")
		err := client.SetBgpPollingTimeGatewayGroup(ctx, groupName, bgpPollingTime)
		if err != nil {
			return diag.Errorf("could not update bgp polling time during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// BGP Neighbor Status Polling Time - API: change_bgp_neighbor_status_polling_time
	// ============================================================================
	if d.HasChange("bgp_neighbor_status_polling_time") {
		bgpBfdPollingTime := getInt(d, "bgp_neighbor_status_polling_time")
		err := client.SetBgpBfdPollingTimeGatewayGroup(ctx, groupName, bgpBfdPollingTime)
		if err != nil {
			return diag.Errorf("could not update bgp neighbor status polling time during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// BGP Hold Time - API: change_bgp_hold_time
	// ============================================================================
	if d.HasChange("bgp_hold_time") {
		bgpHoldTime := getInt(d, "bgp_hold_time")
		err := client.ChangeBgpHoldTimeGatewayGroup(ctx, groupName, bgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// Local AS Number - API: edit_transit_local_as_number
	// Prepend AS Path - API: edit_aviatrix_transit_advanced_config
	// ============================================================================
	if d.HasChanges("local_as_number", "prepend_as_path") {
		localAsNumber := getString(d, "local_as_number")
		err := client.SetLocalASNumberGatewayGroup(ctx, groupName, localAsNumber)
		if err != nil {
			return diag.Errorf("could not set local_as_number for spoke group: %s", err)
		}
		prependASPath := getStringList(d, "prepend_as_path")
		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err = client.SetPrependASPathGatewayGroup(ctx, groupName, prependASPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// BGP ECMP - API: enable_bgp_ecmp / disable_bgp_ecmp
	// ============================================================================
	if d.HasChange("enable_bgp_ecmp") {
		enableBgpEcmp := getBool(d, "enable_bgp_ecmp")
		err := client.SetBgpEcmpGatewayGroup(ctx, groupName, enableBgpEcmp)
		if err != nil {
			return diag.Errorf("could not enable bgp ecmp during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// Active-Standby - API: enable_active_standby / disable_active_standby
	// ============================================================================
	if d.HasChange("enable_active_standby") || d.HasChange("enable_active_standby_preemptive") {
		if getBool(d, "enable_active_standby") {
			if getBool(d, "enable_active_standby_preemptive") {
				if err := client.EnableActiveStandbyPreemptiveGatewayGroup(ctx, groupName); err != nil {
					return diag.Errorf("could not enable Preemptive Mode for Active-Standby during spoke group update: %s", err)
				}
			} else {
				if err := client.EnableActiveStandbyGatewayGroup(ctx, groupName); err != nil {
					return diag.Errorf("could not enable Active-Standby during spoke group update: %s", err)
				}
			}
		} else {
			if getBool(d, "enable_active_standby_preemptive") {
				return diag.Errorf("could not enable Preemptive Mode with Active-Standby disabled")
			}
			if err := client.DisableActiveStandbyGatewayGroup(ctx, groupName); err != nil {
				return diag.Errorf("could not disable Active-Standby during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Jumbo Frame - API: enable_jumbo_frame / disable_jumbo_frame
	// ============================================================================
	if d.HasChange("enable_jumbo_frame") {
		enableJumboFrame := getBool(d, "enable_jumbo_frame")
		if enableJumboFrame {
			err := client.EnableJumboFrameGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during spoke group update: %s", err)
			}
		} else {
			err := client.DisableJumboFrameGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// GRO/GSO - API: enable_gro_gso / disable_gro_gso
	// ============================================================================
	if d.HasChange("enable_gro_gso") {
		enableGroGso := getBool(d, "enable_gro_gso")
		if enableGroGso {
			err := client.EnableGroGsoGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable gro gso during spoke group update: %s", err)
			}
		} else {
			err := client.DisableGroGsoGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable gro gso during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// IPv6 - API: enable_ipv6 / disable_ipv6
	// ============================================================================
	if d.HasChange("enable_ipv6") {
		enableIPv6 := getBool(d, "enable_ipv6")
		if enableIPv6 {
			err := client.EnableIPv6GatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable ipv6 during spoke group update: %s", err)
			}
		} else {
			err := client.DisableIPv6GatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable ipv6 during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Preserve AS Path - API: enable_spoke_preserve_as_path / disable_spoke_preserve_as_path
	// ============================================================================
	if d.HasChange("enable_preserve_as_path") {
		enableBgp := getBool(d, "enable_bgp")
		enableSpokePreserveAsPath := getBool(d, "enable_preserve_as_path")
		if enableSpokePreserveAsPath && !enableBgp {
			return diag.Errorf("enable_preserve_as_path is not supported for Non-BGP spoke group during group update")
		}
		if !enableSpokePreserveAsPath {
			err := client.DisableSpokePreserveAsPathGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable Preserve AS Path during spoke group update: %s", err)
			}
		} else {
			err := client.EnableSpokePreserveAsPathGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable Preserve AS Path during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Learned CIDRs Approval - API: enable_bgp_gateway_cidr_approval / disable_bgp_gateway_cidr_approval
	// ============================================================================
	learnedCidrsApproval := getBool(d, "enable_learned_cidrs_approval")
	if d.HasChange("enable_learned_cidrs_approval") {
		if learnedCidrsApproval {
			err := client.EnableSpokeLearnedCidrsApprovalGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable learned cidrs approval for spoke group: %s", err)
			}
		} else {
			err := client.DisableSpokeLearnedCidrsApprovalGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable learned cidrs approval for spoke group: %s", err)
			}
		}
	}

	// ============================================================================
	// Approved Learned CIDRs - API: set_bgp_gateway_approved_cidr_rules
	// ============================================================================
	if learnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")
		err := client.UpdateSpokePendingApprovedCidrsGatewayGroup(ctx, groupName, approvedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs: %s", err)
		}
	}

	// ============================================================================
	// Private VPC Default Route - API: enable_private_vpc_default_route / disable_private_vpc_default_route
	// ============================================================================
	if d.HasChange("enable_private_vpc_default_route") {
		enablePrivateVpcDefaultRoute := getBool(d, "enable_private_vpc_default_route")
		if enablePrivateVpcDefaultRoute {
			err := client.EnablePrivateVpcDefaultRouteGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable private vpc default route during spoke group update: %s", err)
			}
		} else {
			err := client.DisablePrivateVpcDefaultRouteGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable private vpc default route during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Skip Public Route Table Update - API: enable_skip_public_route_table_update / disable_skip_public_route_table_update
	// ============================================================================
	if d.HasChange("enable_skip_public_route_table_update") {
		enableSkipPublicRouteTableUpdate := getBool(d, "enable_skip_public_route_table_update")
		if enableSkipPublicRouteTableUpdate {
			err := client.EnableSkipPublicRouteUpdateGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable skip public route update during spoke group update: %s", err)
			}
		} else {
			err := client.DisableSkipPublicRouteUpdateGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable skip public route update during spoke group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Auto Advertise S2C CIDRs - API: enable_auto_advertise_s2c_cidrs / disable_auto_advertise_s2c_cidrs
	// ============================================================================
	if d.HasChange("enable_auto_advertise_s2c_cidrs") {
		enableAutoAdvertiseS2CCidrs := getBool(d, "enable_auto_advertise_s2c_cidrs")
		if enableAutoAdvertiseS2CCidrs {
			err := client.EnableAutoAdvertiseS2CCidrsGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable auto advertise s2c cidrs during spoke group update: %s", err)
			}
		} else {
			err := client.DisableAutoAdvertiseS2CCidrsGatewayGroup(ctx, groupName)
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
		cidrs := strings.Join(spokeBgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetSpokeBgpManualAdvertisedNetworksGatewayGroup(ctx, groupName, cidrs)
		if err != nil {
			return diag.Errorf("failed to set spoke bgp manual advertise CIDRs during spoke group update: %s", err)
		}
	}

	// ============================================================================
	// Route Propagation - API: enable_spoke_onprem_route_propagation / disable_spoke_onprem_route_propagation
	// ============================================================================
	if d.HasChange("disable_route_propagation") {
		disableRoutePropagation := getBool(d, "disable_route_propagation")
		enableBgp := getBool(d, "enable_bgp")
		if disableRoutePropagation && !enableBgp {
			return diag.Errorf("disable route propagation is not supported for Non-BGP Spoke during spoke group update")
		}
		if disableRoutePropagation {
			err := client.DisableSpokeOnpremRoutePropagationGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable route propagation for spoke group %s during spoke group update: %v", groupName, err)
			}
		} else {
			err := client.EnableSpokeOnpremRoutePropagationGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable route propagation for spoke group %s during spoke group update: %v", groupName, err)
			}
		}
	}

	// ============================================================================
	// Global VPC (GCP) - API: enable_global_vpc / disable_global_vpc
	// ============================================================================
	if d.HasChange("enable_global_vpc") {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			return diag.Errorf("enable_global_vpc is only valid for GCP related cloud types")
		}
		enableGlobalVpc := getBool(d, "enable_global_vpc")
		if enableGlobalVpc {
			err := client.EnableGlobalVpcGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable global vpc during spoke group update: %s", err)
			}
		} else {
			err := client.DisableGlobalVpcGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable global vpc during spoke group update: %s", err)
			}
		}
	}

	return resourceAviatrixSpokeGroupRead(ctx, d, meta)
}

func resourceAviatrixSpokeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// The resource ID is the group UUID
	groupUUID := d.Id()
	groupName := getString(d, "group_name")

	log.Printf("[INFO] Deleting Spoke Group: %s (UUID: %s)", groupName, groupUUID)

	err := client.DeleteGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to delete spoke group: %s", err)
	}

	return nil
}
