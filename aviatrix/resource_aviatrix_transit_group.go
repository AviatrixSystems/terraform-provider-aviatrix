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

func resourceAviatrixTransitGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixTransitGroupCreate,
		ReadContext:   resourceAviatrixTransitGroupRead,
		UpdateContext: resourceAviatrixTransitGroupUpdate,
		DeleteContext: resourceAviatrixTransitGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: MergeSchemaMaps(
			// Required attributes from group schema
			GroupRequiredSchema(),
			// Resource-specific required attributes
			transitGroupRequiredSchema(),
			// Computed attributes from group schema
			GroupComputedSchema(),
			// Optional attributes from group schema
			GroupOptionalSchema(),
			// Resource-specific optional attributes
			transitGroupOptionalSchema(),
		),
	}
}

// transitGroupRequiredSchema returns the required schema attributes specific to transit group resource.
func transitGroupRequiredSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gw_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Gateway type for the group. Valid values: TRANSIT, EDGETRANSIT, STANDALONE. Case-insensitive.",
			ValidateFunc: validateTransitGwType,
			StateFunc:    normalizeGwType,
		},
	}
}

// transitGroupOptionalSchema returns the optional schema attributes specific to transit group resource.
//
//nolint:funlen
func transitGroupOptionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Description: "Enable NAT (aka single_ip_snat) for the transit group.",
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
			Description: "Enable GRO/GSO for the transit group.",
		},
		"enable_vpc_dns_server": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable VPC DNS server.",
		},
		"enable_s2c_rx_balancing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable S2C receive balancing for the transit group.",
		},

		// ============================================================================
		// TRANSIT-SPECIFIC FEATURES
		// ============================================================================
		"enable_hybrid_connection": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable hybrid connection (TGW connection readiness) for the transit group.",
		},
		"enable_connected_transit": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable connected transit for the transit group.",
		},
		"enable_firenet": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable FireNet for the transit group.",
		},
		"enable_transit_firenet": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable Transit FireNet for the transit group.",
		},
		"enable_advertise_transit_cidr": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable advertise transit CIDR.",
		},
		"enable_transit_summarize_cidr_to_tgw": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable transit summarize CIDR to TGW.",
		},
		"enable_multi_tier_transit": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable multi-tier transit for the transit group.",
		},
		"enable_segmentation": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable segmentation (LAN segmentation) for the transit group.",
		},
		"enable_gateway_load_balancer": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable AWS Gateway Load Balancer for the transit group. Only valid for AWS.",
		},

		// ============================================================================
		// BGP CONFIGURATION
		// ============================================================================
		"enable_bgp": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable BGP for the transit group.",
		},
		"local_as_number": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Changes the Aviatrix Transit Gateway ASN number before you set up Aviatrix Transit Gateway connection configurations.",
			ValidateFunc: goaviatrix.ValidateASN,
		},
		"prepend_as_path": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of AS numbers to populate BGP AS_PATH field when it advertises to VGW or peer devices.",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: goaviatrix.ValidateASN,
			},
		},
		"bgp_manual_spoke_advertise_cidrs": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Set of CIDRs to manually advertise via BGP to spoke gateways.",
		},
		"enable_preserve_as_path": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable preserve AS path.",
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
			ValidateFunc: validation.StringInSlice([]string{"gateway", "connection"}, false),
			Description:  "Learned CIDRs approval mode. Valid values: 'gateway', 'connection'. Default: 'gateway'.",
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
			Description: "Enable Insane Mode (High Performance Encryption) for transit gateway group. Valid values: true, false. Supported for AWS/AWSGov, GCP, Azure and OCI. If insane mode is enabled, gateway size has to at least be c5 size for AWS and Standard_D3_v2 size for Azure.",
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
// Transit Group Create Helper Functions
// ============================================================================

// buildTransitGroupFromResourceData constructs a GatewayGroup struct from Terraform resource data.
func buildTransitGroupFromResourceData(d *schema.ResourceData) *goaviatrix.GatewayGroup {
	transitGroup := &goaviatrix.GatewayGroup{
		GroupName:         getString(d, "group_name"),
		CloudType:         getInt(d, "cloud_type"),
		GwType:            getString(d, "gw_type"),
		GroupInstanceSize: getString(d, "group_instance_size"),
		VpcID:             getString(d, "vpc_id"),
		AccountName:       getString(d, "account_name"),
	}

	// Optional attributes
	if _, ok := d.GetOk("customized_cidr_list"); ok {
		transitGroup.CustomizedCidrList = getStringSet(d, "customized_cidr_list")
	}
	if v, ok := d.GetOk("vpc_region"); ok {
		transitGroup.VpcRegion = mustString(v)
	}
	if v, ok := d.GetOk("domain"); ok {
		transitGroup.Domain = mustString(v)
	}

	// Feature Flags
	transitGroup.EnableNat = getBool(d, "enable_nat")
	transitGroup.EnableIPv6 = getBool(d, "enable_ipv6")
	transitGroup.EnableVpcDNSServer = getBool(d, "enable_vpc_dns_server")
	transitGroup.EnableS2cRxBalancing = getBool(d, "enable_s2c_rx_balancing")

	// Transit-specific features
	transitGroup.EnableHybridConnection = getBool(d, "enable_hybrid_connection")
	transitGroup.EnableConnectedTransit = getBool(d, "enable_connected_transit")
	transitGroup.EnableFirenet = getBool(d, "enable_firenet")
	transitGroup.EnableTransitFirenet = getBool(d, "enable_transit_firenet")
	transitGroup.EnableAdvertiseTransitCidr = getBool(d, "enable_advertise_transit_cidr")
	transitGroup.EnableTransitSummarizeCidrToTgw = getBool(d, "enable_transit_summarize_cidr_to_tgw")
	transitGroup.EnableMultiTierTransit = getBool(d, "enable_multi_tier_transit")
	transitGroup.EnableSegmentation = getBool(d, "enable_segmentation")
	transitGroup.EnableGatewayLoadBalancer = getBool(d, "enable_gateway_load_balancer")

	// BGP Configuration
	transitGroup.EnableBgp = getBool(d, "enable_bgp")

	// BGP AS Number Configuration
	if v, ok := d.GetOk("local_as_number"); ok {
		transitGroup.LocalAsNumber = mustString(v)
	}
	if _, ok := d.GetOk("prepend_as_path"); ok {
		transitGroup.PrependAsPath = getStringList(d, "prepend_as_path")
	}

	if _, ok := d.GetOk("bgp_manual_spoke_advertise_cidrs"); ok {
		transitGroup.SpokeBgpManualAdvertiseCidrs = getStringSet(d, "bgp_manual_spoke_advertise_cidrs")
	}

	transitGroup.EnablePreserveAsPath = getBool(d, "enable_preserve_as_path")
	transitGroup.BgpEcmp = getBool(d, "enable_bgp_ecmp")

	// BGP Timers
	transitGroup.BgpPollingTime = getInt(d, "bgp_polling_time")
	transitGroup.BgpNeighborStatusPollingTime = getInt(d, "bgp_neighbor_status_polling_time")
	transitGroup.BgpHoldTime = getInt(d, "bgp_hold_time")

	// BGP Communities
	transitGroup.BgpSendCommunities = getBool(d, "bgp_send_communities")
	transitGroup.BgpAcceptCommunities = getBool(d, "bgp_accept_communities")

	// BGP over LAN
	transitGroup.EnableBgpOverLan = getBool(d, "enable_bgp_over_lan")

	// Learned CIDR Approval
	transitGroup.EnableLearnedCidrsApproval = getBool(d, "enable_learned_cidrs_approval")
	if v, ok := d.GetOk("learned_cidrs_approval_mode"); ok {
		transitGroup.LearnedCidrsApprovalMode = mustString(v)
	}
	if _, ok := d.GetOk("approved_learned_cidrs"); ok {
		transitGroup.ApprovedLearnedCidrs = getStringSet(d, "approved_learned_cidrs")
	}

	// Active-Standby
	transitGroup.EnableActiveStandby = getBool(d, "enable_active_standby")
	transitGroup.EnableActiveStandbyPreemptive = getBool(d, "enable_active_standby_preemptive")

	// AWS Specific
	transitGroup.InsaneMode = getBool(d, "insane_mode")

	// GCP Specific
	transitGroup.EnableGlobalVpc = getBool(d, "enable_global_vpc")

	return transitGroup
}

// validateTransitGroupConfiguration validates the transit group configuration for cloud-type and feature dependencies.
func validateTransitGroupConfiguration(transitGroup *goaviatrix.GatewayGroup) error {
	if transitGroup.EnableIPv6 && !goaviatrix.IsCloudType(transitGroup.CloudType, goaviatrix.AWS|goaviatrix.Azure) {
		return fmt.Errorf("enable_ipv6 is only valid for AWS (1) and Azure (8)")
	}

	if transitGroup.EnableBgpOverLan && !goaviatrix.IsCloudType(transitGroup.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("enable_bgp_over_lan is only valid for Azure related cloud types")
	}

	if transitGroup.InsaneMode && !goaviatrix.IsCloudType(transitGroup.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCI) {
		return fmt.Errorf("insane_mode is only valid for AWS, GCP, Azure and OCI cloud types")
	}

	if transitGroup.EnableGatewayLoadBalancer && !goaviatrix.IsCloudType(transitGroup.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_gateway_load_balancer is only valid for AWS related cloud types")
	}

	if transitGroup.EnableGlobalVpc && !goaviatrix.IsCloudType(transitGroup.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("enable_global_vpc is only valid for GCP related cloud types")
	}

	if !transitGroup.EnableBgp && transitGroup.LocalAsNumber != "" {
		return fmt.Errorf("local_as_number can only be set when enable_bgp is true")
	}

	if len(transitGroup.PrependAsPath) > 0 && transitGroup.LocalAsNumber == "" {
		return fmt.Errorf("prepend_as_path can only be set when local_as_number is set")
	}

	if transitGroup.EnableActiveStandbyPreemptive && !transitGroup.EnableActiveStandby {
		return fmt.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !transitGroup.EnableLearnedCidrsApproval && len(transitGroup.ApprovedLearnedCidrs) > 0 {
		return fmt.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	if transitGroup.EnableFirenet && transitGroup.EnableTransitFirenet {
		return fmt.Errorf("can't enable firenet and transit firenet at the same time")
	}

	if transitGroup.EnableGatewayLoadBalancer && !transitGroup.EnableFirenet && !transitGroup.EnableTransitFirenet {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
	}

	return nil
}

// applyTransitBgpCommunities applies BGP communities settings (accept/send) to the transit group.
func applyTransitBgpCommunities(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
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

// applyTransitFeatureFlags applies feature flags (jumbo frame, GRO/GSO, NAT, VPC DNS Server, IPv6) to the transit group.
func applyTransitFeatureFlags(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	if getBool(d, "enable_nat") {
		log.Printf("[INFO] Enabling NAT for transit group: %s", groupName)
		if err := client.EnableGatewayGroupSNat(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable NAT: %w", err)
		}
	}

	if getBool(d, "enable_vpc_dns_server") {
		log.Printf("[INFO] Enabling VPC DNS Server for transit group: %s", groupName)
		if err := client.EnableGatewayGroupVpcDNSServer(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable VPC DNS Server: %w", err)
		}
	}

	if getBool(d, "enable_ipv6") {
		log.Printf("[INFO] Enabling IPv6 for transit group: %s", groupName)
		if err := client.EnableIPv6GatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable IPv6: %w", err)
		}
	}

	// Features with default=true: Disable if user explicitly sets to false
	if !getBool(d, "enable_jumbo_frame") {
		log.Printf("[INFO] Disabling Jumbo Frame for transit group: %s", groupName)
		if err := client.DisableJumboFrameGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to disable jumbo frame: %w", err)
		}
	}

	if !getBool(d, "enable_gro_gso") {
		log.Printf("[INFO] Disabling GRO/GSO for transit group: %s", groupName)
		if err := client.DisableGroGsoGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to disable GRO/GSO: %w", err)
		}
	}

	return nil
}

// applyTransitBgpTimers applies BGP timer settings (polling time, neighbor status polling, hold time) to the transit group.
func applyTransitBgpTimers(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
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

// applyTransitBgpConfiguration applies BGP configuration (local AS number, prepend AS path, BGP ECMP) to the transit group.
func applyTransitBgpConfiguration(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
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

// applyTransitActiveStandby applies Active-Standby settings to the transit group.
func applyTransitActiveStandby(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
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

// applyTransitFireNetSettings applies FireNet and Transit FireNet settings to the transit group.
func applyTransitFireNetSettings(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	enableFireNet := getBool(d, "enable_firenet")
	enableTransitFireNet := getBool(d, "enable_transit_firenet")
	enableGatewayLoadBalancer := getBool(d, "enable_gateway_load_balancer")

	if enableFireNet {
		log.Printf("[INFO] Enabling FireNet for transit group: %s", groupName)
		if err := client.EnableFireNetGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable FireNet: %w", err)
		}
	}

	if enableTransitFireNet {
		log.Printf("[INFO] Enabling Transit FireNet for transit group: %s", groupName)
		if enableGatewayLoadBalancer {
			if err := client.EnableTransitFireNetWithGWLBGatewayGroup(ctx, groupName); err != nil {
				return fmt.Errorf("failed to enable Transit FireNet with Gateway Load Balancer: %w", err)
			}
		} else {
			if err := client.EnableTransitFireNetGatewayGroup(ctx, groupName); err != nil {
				return fmt.Errorf("failed to enable Transit FireNet: %w", err)
			}
		}
	}

	return nil
}

// applyTransitSpecificSettings applies transit-specific settings (preserve AS path, learned CIDRs, connected transit, etc.) to the transit group.
func applyTransitSpecificSettings(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, groupName string) error {
	if getBool(d, "enable_preserve_as_path") {
		log.Printf("[INFO] Enabling Preserve AS Path for transit group: %s", groupName)
		if err := client.EnableTransitPreserveAsPathGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable Preserve AS Path: %w", err)
		}
	}

	if getBool(d, "enable_learned_cidrs_approval") {
		log.Printf("[INFO] Enabling learned CIDRs approval for transit group: %s", groupName)
		if err := client.EnableTransitLearnedCidrsApprovalGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable learned CIDRs approval: %w", err)
		}
	}

	approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")
	if len(approvedLearnedCidrs) != 0 {
		if err := client.UpdateTransitPendingApprovedCidrsGatewayGroup(ctx, groupName, approvedLearnedCidrs); err != nil {
			return fmt.Errorf("failed to update approved learned CIDRs: %w", err)
		}
	}

	if getBool(d, "enable_connected_transit") {
		log.Printf("[INFO] Enabling connected transit for transit group: %s", groupName)
		if err := client.EnableConnectedTransitGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable connected transit: %w", err)
		}
	}

	if getBool(d, "enable_segmentation") {
		log.Printf("[INFO] Enabling segmentation for transit group: %s", groupName)
		if err := client.EnableSegmentationGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable segmentation: %w", err)
		}
	}

	if getBool(d, "enable_advertise_transit_cidr") {
		log.Printf("[INFO] Enabling advertise transit CIDR for transit group: %s", groupName)
		if err := client.EnableAdvertiseTransitCidrGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable advertise transit CIDR: %w", err)
		}
	}

	if _, ok := d.GetOk("bgp_manual_spoke_advertise_cidrs"); ok {
		bgpManualSpokeAdvertiseCidrs := getStringSet(d, "bgp_manual_spoke_advertise_cidrs")
		cidrs := strings.Join(bgpManualSpokeAdvertiseCidrs, ",")

		log.Printf("[INFO] Setting BGP manual spoke advertise CIDRs for transit group: %s", groupName)
		if err := client.SetTransitBgpManualAdvertisedNetworksGatewayGroup(ctx, groupName, cidrs); err != nil {
			return fmt.Errorf("failed to set BGP manual spoke advertise CIDRs: %w", err)
		}
	}

	if getBool(d, "enable_global_vpc") {
		log.Printf("[INFO] Enabling global VPC for transit group: %s", groupName)
		if err := client.EnableGlobalVpcGatewayGroup(ctx, groupName); err != nil {
			return fmt.Errorf("failed to enable global VPC: %w", err)
		}
	}

	return nil
}

// ============================================================================
// Transit Group CRUD Operations
// ============================================================================

func resourceAviatrixTransitGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// Build the transit group from resource data
	transitGroup := buildTransitGroupFromResourceData(d)

	// Validate configuration
	if err := validateTransitGroupConfiguration(transitGroup); err != nil {
		return diag.FromErr(err)
	}

	// Create the transit group
	log.Printf("[INFO] Creating Transit Group: %#v", transitGroup)
	if err := client.CreateGatewayGroup(ctx, transitGroup); err != nil {
		return diag.Errorf("failed to create transit group: %s", err)
	}

	// Use GroupUUID as the resource ID
	d.SetId(transitGroup.GroupUUID)
	groupName := getString(d, "group_name")

	// Apply post-creation settings
	if err := applyTransitBgpCommunities(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply BGP communities: %s", err)
	}

	if err := applyTransitFeatureFlags(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply feature flags: %s", err)
	}

	if err := applyTransitBgpTimers(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply BGP timers: %s", err)
	}

	if err := applyTransitBgpConfiguration(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply BGP configuration: %s", err)
	}

	if err := applyTransitActiveStandby(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply Active-Standby: %s", err)
	}

	if err := applyTransitFireNetSettings(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply FireNet settings: %s", err)
	}

	if err := applyTransitSpecificSettings(ctx, d, client, groupName); err != nil {
		return diag.Errorf("failed to apply transit-specific settings: %s", err)
	}

	return resourceAviatrixTransitGroupRead(ctx, d, meta)
}

//nolint:cyclop,funlen
func resourceAviatrixTransitGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// The resource ID is the group UUID
	groupUUID := d.Id()
	if groupUUID == "" {
		return diag.Errorf("resource ID (group UUID) is empty")
	}

	log.Printf("[INFO] Reading Transit Group: %s", groupUUID)

	transitGroup, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read transit group: %s", err)
	}

	// Set required attributes
	mustSet(d, "group_name", transitGroup.GroupName)
	mustSet(d, "cloud_type", transitGroup.CloudType)
	// Normalize gw_type by removing "GwGroupType." prefix if present
	gwType := strings.TrimPrefix(transitGroup.GwType, "GwGroupType.")
	mustSet(d, "gw_type", gwType)
	mustSet(d, "group_instance_size", transitGroup.GroupInstanceSize)
	mustSet(d, "vpc_id", transitGroup.VpcID)
	mustSet(d, "account_name", transitGroup.AccountName)

	// Store group_uuid in state (used as resource ID)
	if transitGroup.GroupUUID != "" {
		mustSet(d, "group_uuid", transitGroup.GroupUUID)
	}

	mustSet(d, "customized_cidr_list", transitGroup.CustomizedCidrList)
	mustSet(d, "explicitly_created", transitGroup.ExplicitlyCreated)
	mustSet(d, "vpc_region", transitGroup.VpcRegion)
	mustSet(d, "domain", transitGroup.Domain)

	// Feature Flags
	mustSet(d, "enable_jumbo_frame", transitGroup.EnableJumboFrame)
	mustSet(d, "enable_nat", transitGroup.EnableNat)
	mustSet(d, "enable_ipv6", transitGroup.EnableIPv6)
	mustSet(d, "enable_gro_gso", transitGroup.EnableGroGso)
	mustSet(d, "enable_vpc_dns_server", transitGroup.EnableVpcDNSServer)
	mustSet(d, "enable_s2c_rx_balancing", transitGroup.EnableS2cRxBalancing)

	// Transit-specific features
	mustSet(d, "enable_hybrid_connection", transitGroup.EnableHybridConnection)
	mustSet(d, "enable_connected_transit", transitGroup.EnableConnectedTransit)
	mustSet(d, "enable_firenet", transitGroup.EnableFirenet)
	mustSet(d, "enable_transit_firenet", transitGroup.EnableTransitFirenet)
	mustSet(d, "enable_advertise_transit_cidr", transitGroup.EnableAdvertiseTransitCidr)
	mustSet(d, "enable_transit_summarize_cidr_to_tgw", transitGroup.EnableTransitSummarizeCidrToTgw)
	mustSet(d, "enable_multi_tier_transit", transitGroup.EnableMultiTierTransit)
	mustSet(d, "enable_segmentation", transitGroup.EnableSegmentation)
	mustSet(d, "enable_gateway_load_balancer", transitGroup.EnableGatewayLoadBalancer)

	// BGP Configuration
	mustSet(d, "enable_bgp", transitGroup.EnableBgp)
	mustSet(d, "local_as_number", transitGroup.LocalAsNumber)
	mustSet(d, "prepend_as_path", transitGroup.PrependAsPath)
	mustSet(d, "bgp_manual_spoke_advertise_cidrs", transitGroup.SpokeBgpManualAdvertiseCidrs)
	mustSet(d, "enable_preserve_as_path", transitGroup.EnablePreserveAsPath)
	mustSet(d, "enable_bgp_ecmp", transitGroup.BgpEcmp)

	// BGP Timers - Set to schema defaults if API returns 0
	if transitGroup.BgpPollingTime == 0 {
		mustSet(d, "bgp_polling_time", defaultBgpPollingTime)
	} else {
		mustSet(d, "bgp_polling_time", transitGroup.BgpPollingTime)
	}
	if transitGroup.BgpNeighborStatusPollingTime == 0 {
		mustSet(d, "bgp_neighbor_status_polling_time", defaultBgpNeighborStatusPollingTime)
	} else {
		mustSet(d, "bgp_neighbor_status_polling_time", transitGroup.BgpNeighborStatusPollingTime)
	}
	if transitGroup.BgpHoldTime == 0 {
		mustSet(d, "bgp_hold_time", defaultBgpHoldTime)
	} else {
		mustSet(d, "bgp_hold_time", transitGroup.BgpHoldTime)
	}

	// BGP Communities
	mustSet(d, "bgp_send_communities", transitGroup.BgpSendCommunities)
	mustSet(d, "bgp_accept_communities", transitGroup.BgpAcceptCommunities)

	// BGP over LAN
	mustSet(d, "enable_bgp_over_lan", transitGroup.EnableBgpOverLan)

	// Learned CIDR Approval
	mustSet(d, "enable_learned_cidrs_approval", transitGroup.EnableLearnedCidrsApproval)
	// Set to schema default if API returns empty string
	if transitGroup.LearnedCidrsApprovalMode == "" {
		mustSet(d, "learned_cidrs_approval_mode", "gateway")
	} else {
		mustSet(d, "learned_cidrs_approval_mode", transitGroup.LearnedCidrsApprovalMode)
	}
	mustSet(d, "approved_learned_cidrs", transitGroup.ApprovedLearnedCidrs)

	// Active-Standby
	mustSet(d, "enable_active_standby", transitGroup.EnableActiveStandby)
	mustSet(d, "enable_active_standby_preemptive", transitGroup.EnableActiveStandbyPreemptive)

	// AWS Specific
	mustSet(d, "insane_mode", transitGroup.InsaneMode)

	// GCP Specific
	mustSet(d, "enable_global_vpc", transitGroup.EnableGlobalVpc)

	// Computed attributes
	mustSet(d, "gw_uuid_list", transitGroup.GwUUIDList)
	mustSet(d, "vpc_uuid", transitGroup.VpcUUID)
	mustSet(d, "vendor_name", transitGroup.VendorName)

	// Use GroupUUID as the resource ID
	if transitGroup.GroupUUID != "" {
		d.SetId(transitGroup.GroupUUID)
	}
	return nil
}

//nolint:cyclop,funlen
func resourceAviatrixTransitGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	groupName := getString(d, "group_name")
	cloudType := getInt(d, "cloud_type")
	groupUUID := d.Id()

	log.Printf("[INFO] Updating Transit Group: %s (UUID: %s)", groupName, groupUUID)

	// Validations for update
	if getBool(d, "enable_active_standby_preemptive") && !getBool(d, "enable_active_standby") {
		return diag.Errorf("enable_active_standby_preemptive requires enable_active_standby to be true")
	}

	if !getBool(d, "enable_learned_cidrs_approval") && len(getStringSet(d, "approved_learned_cidrs")) > 0 {
		return diag.Errorf("approved_learned_cidrs requires enable_learned_cidrs_approval to be true")
	}

	if getBool(d, "enable_firenet") && getBool(d, "enable_transit_firenet") {
		return diag.Errorf("can't enable firenet and transit firenet at the same time")
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
				return diag.Errorf("failed to enable NAT for transit group: %s", err)
			}
		} else {
			err := client.DisableGatewayGroupSNat(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable NAT for transit group: %s", err)
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
				return diag.Errorf("failed to enable VPC DNS Server for transit group: %s", err)
			}
		} else {
			err := client.DisableGatewayGroupVpcDNSServer(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable VPC DNS Server for transit group: %s", err)
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
			return diag.Errorf("could not update bgp polling time during transit group update: %s", err)
		}
	}

	// ============================================================================
	// BGP Neighbor Status Polling Time - API: change_bgp_neighbor_status_polling_time
	// ============================================================================
	if d.HasChange("bgp_neighbor_status_polling_time") {
		bgpBfdPollingTime := getInt(d, "bgp_neighbor_status_polling_time")
		err := client.SetBgpBfdPollingTimeGatewayGroup(ctx, groupName, bgpBfdPollingTime)
		if err != nil {
			return diag.Errorf("could not update bgp neighbor status polling time during transit group update: %s", err)
		}
	}

	// ============================================================================
	// BGP Hold Time - API: change_bgp_hold_time
	// ============================================================================
	if d.HasChange("bgp_hold_time") {
		bgpHoldTime := getInt(d, "bgp_hold_time")
		err := client.ChangeBgpHoldTimeGatewayGroup(ctx, groupName, bgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time during transit group update: %s", err)
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
			return diag.Errorf("could not set local_as_number for transit group: %s", err)
		}
		prependASPath := getStringList(d, "prepend_as_path")
		if d.HasChange("prepend_as_path") {
			err = client.SetPrependASPathGatewayGroup(ctx, groupName, prependASPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path for transit group: %s", err)
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
			return diag.Errorf("could not enable bgp ecmp during transit group update: %s", err)
		}
	}

	// ============================================================================
	// Active-Standby - API: enable_active_standby / disable_active_standby
	// ============================================================================
	if d.HasChange("enable_active_standby") || d.HasChange("enable_active_standby_preemptive") {
		if getBool(d, "enable_active_standby") {
			if getBool(d, "enable_active_standby_preemptive") {
				if err := client.EnableActiveStandbyPreemptiveGatewayGroup(ctx, groupName); err != nil {
					return diag.Errorf("could not enable Preemptive Mode for Active-Standby during transit group update: %s", err)
				}
			} else {
				if err := client.EnableActiveStandbyGatewayGroup(ctx, groupName); err != nil {
					return diag.Errorf("could not enable Active-Standby during transit group update: %s", err)
				}
			}
		} else {
			if getBool(d, "enable_active_standby_preemptive") {
				return diag.Errorf("could not enable Preemptive Mode with Active-Standby disabled")
			}
			if err := client.DisableActiveStandbyGatewayGroup(ctx, groupName); err != nil {
				return diag.Errorf("could not disable Active-Standby during transit group update: %s", err)
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
				return diag.Errorf("could not enable jumbo frame during transit group update: %s", err)
			}
		} else {
			err := client.DisableJumboFrameGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during transit group update: %s", err)
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
				return diag.Errorf("could not enable gro gso during transit group update: %s", err)
			}
		} else {
			err := client.DisableGroGsoGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable gro gso during transit group update: %s", err)
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
				return diag.Errorf("could not enable ipv6 during transit group update: %s", err)
			}
		} else {
			err := client.DisableIPv6GatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable ipv6 during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Preserve AS Path - API: enable_transit_preserve_as_path / disable_transit_preserve_as_path
	// ============================================================================
	if d.HasChange("enable_preserve_as_path") {
		enableBgp := getBool(d, "enable_bgp")
		enablePreserveAsPath := getBool(d, "enable_preserve_as_path")
		if enablePreserveAsPath && !enableBgp {
			return diag.Errorf("enable_preserve_as_path is not supported for Non-BGP transit group during group update")
		}
		if !enablePreserveAsPath {
			err := client.DisableTransitPreserveAsPathGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable Preserve AS Path during transit group update: %s", err)
			}
		} else {
			err := client.EnableTransitPreserveAsPathGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable Preserve AS Path during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Learned CIDRs Approval - API: enable_transit_learned_cidrs_approval / disable_transit_learned_cidrs_approval
	// ============================================================================
	learnedCidrsApproval := getBool(d, "enable_learned_cidrs_approval")
	if d.HasChange("enable_learned_cidrs_approval") {
		if learnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApprovalGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable learned cidrs approval for transit group: %s", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApprovalGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable learned cidrs approval for transit group: %s", err)
			}
		}
	}

	// ============================================================================
	// Approved Learned CIDRs - API: set_bgp_gateway_approved_cidr_rules
	// ============================================================================
	if learnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")
		err := client.UpdateTransitPendingApprovedCidrsGatewayGroup(ctx, groupName, approvedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs: %s", err)
		}
	}

	// ============================================================================
	// Connected Transit - API: enable_connected_transit / disable_connected_transit
	// ============================================================================
	if d.HasChange("enable_connected_transit") {
		enableConnectedTransit := getBool(d, "enable_connected_transit")
		if enableConnectedTransit {
			err := client.EnableConnectedTransitGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable connected transit during transit group update: %s", err)
			}
		} else {
			err := client.DisableConnectedTransitGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable connected transit during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Segmentation - API: enable_segmentation / disable_segmentation
	// ============================================================================
	if d.HasChange("enable_segmentation") {
		enableSegmentation := getBool(d, "enable_segmentation")
		if enableSegmentation {
			err := client.EnableSegmentationGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable segmentation during transit group update: %s", err)
			}
		} else {
			err := client.DisableSegmentationGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable segmentation during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Advertise Transit CIDR - API: enable_advertise_transit_cidr / disable_advertise_transit_cidr
	// ============================================================================
	if d.HasChange("enable_advertise_transit_cidr") {
		enableAdvertiseTransitCidr := getBool(d, "enable_advertise_transit_cidr")
		if enableAdvertiseTransitCidr {
			err := client.EnableAdvertiseTransitCidrGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable advertise transit CIDR during transit group update: %s", err)
			}
		} else {
			err := client.DisableAdvertiseTransitCidrGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable advertise transit CIDR during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// BGP Manual Spoke Advertise CIDRs - API: edit_aviatrix_transit_advanced_config
	// ============================================================================
	if d.HasChange("bgp_manual_spoke_advertise_cidrs") {
		bgpManualSpokeAdvertiseCidrs := getStringSet(d, "bgp_manual_spoke_advertise_cidrs")
		cidrs := strings.Join(bgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetTransitBgpManualAdvertisedNetworksGatewayGroup(ctx, groupName, cidrs)
		if err != nil {
			return diag.Errorf("failed to set BGP manual spoke advertise CIDRs during transit group update: %s", err)
		}
	}

	// ============================================================================
	// Hybrid Connection - API: enable_hybrid_connection / disable_hybrid_connection
	// ============================================================================
	if d.HasChange("enable_hybrid_connection") {
		enableHybridConnection := getBool(d, "enable_hybrid_connection")
		if enableHybridConnection {
			err := client.EnableHybridConnectionGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable hybrid connection during transit group update: %s", err)
			}
		} else {
			err := client.DisableHybridConnectionGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable hybrid connection during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Transit Summarize CIDR to TGW - API: enable_transit_summarize_cidr_to_tgw / disable_transit_summarize_cidr_to_tgw
	// ============================================================================
	if d.HasChange("enable_transit_summarize_cidr_to_tgw") {
		enableTransitSummarizeCidrToTgw := getBool(d, "enable_transit_summarize_cidr_to_tgw")
		if enableTransitSummarizeCidrToTgw {
			err := client.EnableTransitSummarizeCidrToTgwGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable transit summarize CIDR to TGW during transit group update: %s", err)
			}
		} else {
			err := client.DisableTransitSummarizeCidrToTgwGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable transit summarize CIDR to TGW during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Multi-Tier Transit - API: enable_multi_tier_transit / disable_multi_tier_transit
	// ============================================================================
	if d.HasChange("enable_multi_tier_transit") {
		enableMultiTierTransit := getBool(d, "enable_multi_tier_transit")
		if enableMultiTierTransit {
			err := client.EnableMultiTierTransitGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable multi-tier transit during transit group update: %s", err)
			}
		} else {
			err := client.DisableMultiTierTransitGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable multi-tier transit during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// S2C RX Balancing - API: enable_s2c_rx_balancing / disable_s2c_rx_balancing
	// ============================================================================
	if d.HasChange("enable_s2c_rx_balancing") {
		enableS2cRxBalancing := getBool(d, "enable_s2c_rx_balancing")
		if enableS2cRxBalancing {
			err := client.EnableS2cRxBalancingGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable S2C RX balancing during transit group update: %s", err)
			}
		} else {
			err := client.DisableS2cRxBalancingGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable S2C RX balancing during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// Gateway Load Balancer (AWS) - API: enable_gateway_load_balancer / disable_gateway_load_balancer
	// ============================================================================
	if d.HasChange("enable_gateway_load_balancer") {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			return diag.Errorf("enable_gateway_load_balancer is only valid for AWS related cloud types")
		}
		enableGatewayLoadBalancer := getBool(d, "enable_gateway_load_balancer")
		if enableGatewayLoadBalancer {
			err := client.EnableGatewayLoadBalancerGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not enable gateway load balancer during transit group update: %s", err)
			}
		} else {
			err := client.DisableGatewayLoadBalancerGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable gateway load balancer during transit group update: %s", err)
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
				return diag.Errorf("could not enable global vpc during transit group update: %s", err)
			}
		} else {
			err := client.DisableGlobalVpcGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("could not disable global vpc during transit group update: %s", err)
			}
		}
	}

	// ============================================================================
	// FireNet and Transit FireNet - API: enable_gateway_firenet_interfaces / enable_gateway_for_transit_firenet
	// ============================================================================
	enableFireNet := getBool(d, "enable_firenet")
	enableTransitFireNet := getBool(d, "enable_transit_firenet")
	enableGatewayLoadBalancer := getBool(d, "enable_gateway_load_balancer")

	if enableFireNet && enableTransitFireNet {
		return diag.Errorf("can't enable firenet and transit firenet at the same time")
	}

	if enableGatewayLoadBalancer && !enableFireNet && !enableTransitFireNet {
		return diag.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
	}

	if d.HasChange("enable_firenet") && d.HasChange("enable_transit_firenet") {
		// Both changed at once - disable firenet first, then enable transit_firenet if needed
		if !enableFireNet {
			err := client.DisableFireNetGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable FireNet during transit group update: %s", err)
			}
		}

		if !enableTransitFireNet {
			err := client.DisableTransitFireNetGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable Transit FireNet during transit group update: %s", err)
			}
		}

		if enableFireNet {
			err := client.EnableFireNetGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable FireNet during transit group update: %s", err)
			}
		}

		if enableTransitFireNet {
			if enableGatewayLoadBalancer {
				err := client.EnableTransitFireNetWithGWLBGatewayGroup(ctx, groupName)
				if err != nil {
					return diag.Errorf("failed to enable Transit FireNet with Gateway Load Balancer during transit group update: %s", err)
				}
			} else {
				err := client.EnableTransitFireNetGatewayGroup(ctx, groupName)
				if err != nil {
					return diag.Errorf("failed to enable Transit FireNet during transit group update: %s", err)
				}
			}
		}
	} else if d.HasChange("enable_firenet") {
		if enableFireNet {
			err := client.EnableFireNetGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to enable FireNet during transit group update: %s", err)
			}
		} else {
			err := client.DisableFireNetGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable FireNet during transit group update: %s", err)
			}
		}
	} else if d.HasChange("enable_transit_firenet") {
		if enableTransitFireNet {
			if enableGatewayLoadBalancer {
				err := client.EnableTransitFireNetWithGWLBGatewayGroup(ctx, groupName)
				if err != nil {
					return diag.Errorf("failed to enable Transit FireNet with Gateway Load Balancer during transit group update: %s", err)
				}
			} else {
				err := client.EnableTransitFireNetGatewayGroup(ctx, groupName)
				if err != nil {
					return diag.Errorf("failed to enable Transit FireNet during transit group update: %s", err)
				}
			}
		} else {
			err := client.DisableTransitFireNetGatewayGroup(ctx, groupName)
			if err != nil {
				return diag.Errorf("failed to disable Transit FireNet during transit group update: %s", err)
			}
		}
	} else if d.HasChange("enable_gateway_load_balancer") {
		// In this branch we know that neither 'enable_transit_firenet' or 'enable_firenet' HasChange.
		// Due to the backend design it is not possible to disable or enable 'enable_gateway_load_balancer' without
		// also disabling or enabling FireNet, so we force the user to disable or enable both at the same time.
		if enableGatewayLoadBalancer {
			return diag.Errorf("can not enable 'enable_gateway_load_balancer' when 'enable_firenet' or 'enable_transit_firenet' is " +
				"already enabled. Changing from non-GWLB FireNet to GWLB FireNet requires 2 separate " +
				"`terraform apply` steps, once to disable non-GWLB FireNet, then again to enable GWLB FireNet")
		} else {
			return diag.Errorf("can not disable 'enable_gateway_load_balancer' when 'enable_firenet' or 'enable_transit_firenet' is " +
				"still enabled. Changing from GWLB FireNet to non-GWLB FireNet requires 2 separate " +
				"`terraform apply` steps, once to disable GWLB FireNet, then again to enable non-GWLB FireNet")
		}
	}

	return resourceAviatrixTransitGroupRead(ctx, d, meta)
}

func resourceAviatrixTransitGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// The resource ID is the group UUID
	groupUUID := d.Id()
	groupName := getString(d, "group_name")

	log.Printf("[INFO] Deleting Transit Group: %s (UUID: %s)", groupName, groupUUID)

	err := client.DeleteGatewayGroup(ctx, groupUUID)
	if err != nil {
		return diag.Errorf("failed to delete transit group: %s", err)
	}

	return nil
}
