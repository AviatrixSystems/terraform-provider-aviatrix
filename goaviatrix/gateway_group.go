package goaviatrix

import (
	"context"
	"fmt"
)

// GatewayGroup represents a Gateway Group resource
type GatewayGroup struct {
	Action string `form:"action,omitempty" json:"action,omitempty"`
	CID    string `form:"CID,omitempty" json:"CID,omitempty"`
	// Required
	GroupName         string `form:"group_name,omitempty" json:"group_name,omitempty"` // Changed json tag from "name" to "group_name"
	GroupUUID         string `form:"group_uuid,omitempty" json:"uuid,omitempty"`
	CloudType         int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	GwType            string `form:"gw_type,omitempty" json:"gw_type,omitempty"`
	GroupInstanceSize string `form:"group_instance_size,omitempty" json:"group_instance_size,omitempty"`
	VpcID             string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	AccountName       string `form:"account_name,omitempty" json:"account_name,omitempty"`

	// Optional
	CustomizedCidrList               []string `form:"customized_cidr_list,omitempty" json:"customized_cidr_list,omitempty"`
	S2cRxBalancing                   bool     `form:"s2c_rx_balancing,omitempty" json:"s2c_rx_balancing,omitempty"`
	ExplicitlyCreated                bool     `form:"explicitly_created,omitempty" json:"explicitly_created,omitempty"`
	Subnet                           string   `form:"subnet,omitempty" json:"subnet,omitempty"`
	VpcRegion                        string   `form:"vpc_region,omitempty" json:"vpc_region,omitempty"`
	Domain                           string   `form:"domain,omitempty" json:"domain,omitempty"`
	IncludeCidr                      []string `form:"include_cidr,omitempty" json:"include_cidr,omitempty"`
	EnablePrivateVpcDefaultRoute     bool     `form:"private_vpc_default_enabled,omitempty" json:"private_vpc_default_enabled,omitempty"`
	EnableSkipPublicRouteTableUpdate bool     `form:"skip_public_vpc_update_enabled,omitempty" json:"skip_public_vpc_update_enabled,omitempty"`

	// Feature Flags
	// Note: EnableJumboFrame and EnableGroGso do not have form tags because they must be set
	// via separate enable/disable API calls after creation, not during the create call
	EnableJumboFrame   bool `form:"group_jumbo_frame,omitempty" json:"group_jumbo_frame,omitempty"`
	EnableNat          bool `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	EnableIPv6         bool `form:"is_ipv6_enabled,omitempty" json:"is_ipv6_enabled,omitempty"`
	EnableGroGso       bool `form:"gro_gso,omitempty" json:"gro_gso,omitempty"`
	EnableVpcDNSServer bool `form:"use_vpc_dns_server,omitempty" json:"use_vpc_dns_server,omitempty"`

	// BGP Configuration
	EnableBgp                    bool     `form:"bgp_enabled,omitempty" json:"bgp_enabled,omitempty"`
	LocalAsNumber                string   `form:"bgp_local_as_num,omitempty" json:"bgp_local_as_num,omitempty"`
	PrependAsPath                []string `form:"bgp_prepend_as_path,omitempty" json:"bgp_prepend_as_path,omitempty"`
	DisableRoutePropagation      bool     `form:"disable_route_propagation,omitempty" json:"disable_route_propagation,omitempty"`
	SpokeBgpManualAdvertiseCidrs []string `form:"bgp_manual_spoke_advertise_cidrs,omitempty" json:"bgp_manual_spoke_advertise_cidrs,omitempty"`
	EnablePreserveAsPath         bool     `form:"preserve_as_path,omitempty" json:"preserve_as_path,omitempty"`
	EnableAutoAdvertiseS2cCidrs  bool     `form:"auto_advertise_s2c_cidrs,omitempty" json:"auto_advertise_s2c_cidrs,omitempty"`
	BgpEcmp                      bool     `form:"bgp_ecmp,omitempty" json:"bgp_ecmp,omitempty"`

	// BGP Timers
	BgpPollingTime               int `form:"bgp_polling_time,omitempty" json:"bgp_polling_time,omitempty"`
	BgpNeighborStatusPollingTime int `form:"bgp_neighbor_status_polling_time,omitempty" json:"bgp_neighbor_status_polling_time,omitempty"`
	BgpHoldTime                  int `form:"bgp_hold_time,omitempty" json:"bgp_hold_time,omitempty"`

	// BGP Communities
	BgpSendCommunities   bool `form:"send_communities,omitempty" json:"send_communities,omitempty"`
	BgpAcceptCommunities bool `form:"accept_communities,omitempty" json:"accept_communities,omitempty"`

	// BGP over LAN
	EnableBgpOverLan bool `form:"enable_bgp_over_lan,omitempty" json:"enable_bgp_over_lan,omitempty"`

	// Learned CIDR Approval
	EnableLearnedCidrsApproval bool     `form:"learned_cidrs_approval,omitempty" json:"learned_cidrs_approval,omitempty"`
	LearnedCidrsApprovalMode   string   `form:"learned_cidrs_approval_mode,omitempty" json:"learned_cidrs_approval_mode,omitempty"`
	ApprovedLearnedCidrs       []string `form:"approved_learned_cidrs,omitempty" json:"approved_learned_cidrs,omitempty"`

	// Active-Standby
	EnableActiveStandby           bool `form:"active_standby,omitempty" json:"active_standby,omitempty"`
	EnableActiveStandbyPreemptive bool `form:"active_standby_preemptive,omitempty" json:"active_standby_preemptive,omitempty"`

	// AWS Specific
	InsaneMode          bool   `form:"group_hpe_enabled,omitempty" json:"group_hpe_enabled,omitempty"`
	EnableEncryptVolume bool   `form:"gw_enc,omitempty" json:"gw_enc,omitempty"`
	CustomerManagedKeys string `form:"customer_managed_keys,omitempty" json:"customer_managed_keys,omitempty"`

	// GCP Specific
	EnableGlobalVpc bool `form:"global_vpc,omitempty" json:"global_vpc,omitempty"`

	// Transit-specific fields
	EnableConnectedTransit          bool `form:"connected_transit,omitempty" json:"connected_transit,omitempty"`
	EnableFirenet                   bool `form:"firenet_enabled,omitempty" json:"firenet_enabled,omitempty"`
	EnableTransitFirenet            bool `form:"transit_firenet_enabled,omitempty" json:"transit_firenet_enabled,omitempty"`
	EnableAdvertiseTransitCidr      bool `form:"advertise_transit_cidr,omitempty" json:"advertise_transit_cidr,omitempty"`
	EnableSegmentation              bool `form:"segmentation_enabled,omitempty" json:"segmentation_enabled,omitempty"`
	EnableHybridConnection          bool `form:"enable_hybrid_connection,omitempty" json:"enable_hybrid_connection,omitempty"`
	EnableTransitSummarizeCidrToTgw bool `form:"transit_summarize_cidr_to_tgw,omitempty" json:"transit_summarize_cidr_to_tgw,omitempty"`
	EnableMultiTierTransit          bool `form:"multi_tier_transit,omitempty" json:"multi_tier_transit,omitempty"`
	EnableS2cRxBalancing            bool `form:"enable_s2c_rx_balancing,omitempty" json:"enable_s2c_rx_balancing,omitempty"`
	EnableGatewayLoadBalancer       bool `form:"enable_gateway_load_balancer,omitempty" json:"enable_gateway_load_balancer,omitempty"`

	// Computed (read-only)
	GwUUIDList []string `json:"gw_uuid_list,omitempty"`
	VpcUUID    string   `json:"vpc_uuid,omitempty"`
	VendorName string   `json:"vendor_name,omitempty"`
}

// GatewayGroupResp represents the API response for get gateway group
type GatewayGroupResp struct {
	Return  bool         `json:"return"`
	Results GatewayGroup `json:"results"`
	Reason  string       `json:"reason"`
}

// CreateGatewayGroup creates a new gateway group
func (c *Client) CreateGatewayGroup(ctx context.Context, spokeGroup *GatewayGroup) error {
	spokeGroup.Action = "create_gateway_group"
	spokeGroup.CID = c.CID

	type Resp struct {
		Return    bool   `json:"return"`
		Results   string `json:"results"`
		Reason    string `json:"reason"`
		GroupName string `json:"group_name"`
		GroupUUID string `json:"group_uuid"`
	}

	var resp Resp
	err := c.PostAPIContext2(ctx, &resp, spokeGroup.Action, spokeGroup, BasicCheck)
	if err != nil {
		return err
	}

	// Capture the UUID returned from the create operation
	spokeGroup.GroupUUID = resp.GroupUUID

	return nil
}

// GetGatewayGroup retrieves gateway group details by UUID
func (c *Client) GetGatewayGroup(ctx context.Context, groupUUID string) (*GatewayGroup, error) {
	form := map[string]string{
		"action":     "get_gateway_group_details",
		"CID":        c.CID,
		"group_uuid": groupUUID,
	}

	var resp GatewayGroupResp
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if !resp.Return {
		return nil, fmt.Errorf("failed to get gateway group: %s", resp.Reason)
	}

	return &resp.Results, nil
}

// UpdateGatewayGroup updates an existing gateway group
func (c *Client) UpdateGatewayGroup(ctx context.Context, gatewayGroup *GatewayGroup) error {
	gatewayGroup.Action = "update_gateway_group"
	gatewayGroup.CID = c.CID

	return c.PostAPIContext2(ctx, nil, gatewayGroup.Action, gatewayGroup, BasicCheck)
}

// DeleteGatewayGroup deletes a gateway group by UUID
func (c *Client) DeleteGatewayGroup(ctx context.Context, groupUUID string) error {
	form := map[string]string{
		"action":     "delete_gateway_group",
		"CID":        c.CID,
		"group_uuid": groupUUID,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// ============================================================================
// Gateway Group Update APIs
// ============================================================================

// UpdateGatewayGroupSize updates the gateway group instance size (edit_gw_config)
func (c *Client) UpdateGatewayGroupSize(ctx context.Context, groupName, instanceSize string) error {
	form := map[string]string{
		"action":              "edit_gw_config",
		"CID":                 c.CID,
		"gateway_name":        groupName,
		"group_instance_size": instanceSize,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// SetGatewayGroupBgpCommunitiesAccept sets the BGP accept communities for a gateway group
func (c *Client) SetGatewayGroupBgpCommunitiesAccept(ctx context.Context, groupName string, accept bool) error {
	action := "set_gateway_accept_bgp_communities_override"
	form := map[string]interface{}{
		"action":             action,
		"CID":                c.CID,
		"gateway_name":       groupName,
		"accept_communities": accept,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// SetGatewayGroupBgpCommunitiesSend sets the BGP send communities for a gateway group
func (c *Client) SetGatewayGroupBgpCommunitiesSend(ctx context.Context, groupName string, send bool) error {
	action := "set_gateway_send_bgp_communities_override"
	form := map[string]interface{}{
		"action":           action,
		"CID":              c.CID,
		"gateway_name":     groupName,
		"send_communities": send,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// EnableGatewayGroupSNat enables SNAT for a gateway group
func (c *Client) EnableGatewayGroupSNat(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_snat",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableGatewayGroupSNat disables SNAT for a gateway group
func (c *Client) DisableGatewayGroupSNat(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_snat",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableGatewayGroupVpcDNSServer enables VPC DNS Server for a gateway group
func (c *Client) EnableGatewayGroupVpcDNSServer(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_vpc_dns_server",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableGatewayGroupVpcDNSServer disables VPC DNS Server for a gateway group
func (c *Client) DisableGatewayGroupVpcDNSServer(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_vpc_dns_server",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// SetBgpPollingTimeGatewayGroup sets the BGP polling time for a gateway group
func (c *Client) SetBgpPollingTimeGatewayGroup(ctx context.Context, groupName string, pollingTime int) error {
	action := "change_bgp_polling_time"
	form := map[string]interface{}{
		"action":           action,
		"CID":              c.CID,
		"group_name":       groupName,
		"bgp_polling_time": pollingTime,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// SetBgpBfdPollingTimeGatewayGroup sets the BGP neighbor status polling time for a gateway group
func (c *Client) SetBgpBfdPollingTimeGatewayGroup(ctx context.Context, groupName string, pollingTime int) error {
	action := "change_bgp_neighbor_status_polling_time"
	form := map[string]interface{}{
		"action":                           action,
		"CID":                              c.CID,
		"group_name":                       groupName,
		"bgp_neighbor_status_polling_time": pollingTime,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// ChangeBgpHoldTimeGatewayGroup changes the BGP hold time for a gateway group
func (c *Client) ChangeBgpHoldTimeGatewayGroup(ctx context.Context, groupName string, holdTime int) error {
	action := "change_bgp_hold_time"
	form := map[string]interface{}{
		"action":        action,
		"CID":           c.CID,
		"group_name":    groupName,
		"bgp_hold_time": holdTime,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// SetLocalASNumberGatewayGroup sets the local AS number for a gateway group
func (c *Client) SetLocalASNumberGatewayGroup(ctx context.Context, groupName, localAsNumber string) error {
	form := map[string]string{
		"action":          "edit_transit_local_as_number",
		"CID":             c.CID,
		"group_name":      groupName,
		"local_as_number": localAsNumber,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// SetPrependASPathGatewayGroup sets the prepend AS path for a gateway group
func (c *Client) SetPrependASPathGatewayGroup(ctx context.Context, groupName string, prependAsPath []string) error {
	action := "edit_aviatrix_transit_advanced_config"
	form := map[string]interface{}{
		"action":          action,
		"CID":             c.CID,
		"group_name":      groupName,
		"prepend_as_path": prependAsPath,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// SetBgpEcmpGatewayGroup enables or disables BGP ECMP for a gateway group
func (c *Client) SetBgpEcmpGatewayGroup(ctx context.Context, groupName string, enable bool) error {
	action := "disable_bgp_ecmp"
	if enable {
		action = "enable_bgp_ecmp"
	}

	form := map[string]string{
		"action":       action,
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableActiveStandbyGatewayGroup enables Active-Standby mode for a gateway group
func (c *Client) EnableActiveStandbyGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_active_standby",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableActiveStandbyPreemptiveGatewayGroup enables Active-Standby Preemptive mode for a gateway group
func (c *Client) EnableActiveStandbyPreemptiveGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_active_standby",
		"CID":          c.CID,
		"gateway_name": groupName,
		"preemptive":   "true",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableActiveStandbyGatewayGroup disables Active-Standby mode for a gateway group
func (c *Client) DisableActiveStandbyGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_active_standby",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableJumboFrameGatewayGroup enables jumbo frame support for a gateway group
func (c *Client) EnableJumboFrameGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_jumbo_frame",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableJumboFrameGatewayGroup disables jumbo frame support for a gateway group
func (c *Client) DisableJumboFrameGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_jumbo_frame",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableGroGsoGatewayGroup enables GRO/GSO for a gateway group
func (c *Client) EnableGroGsoGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_gro_gso",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableGroGsoGatewayGroup disables GRO/GSO for a gateway group
func (c *Client) DisableGroGsoGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_gro_gso",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableIPv6GatewayGroup enables IPv6 for a gateway group
func (c *Client) EnableIPv6GatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_ipv6",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableIPv6GatewayGroup disables IPv6 for a gateway group
func (c *Client) DisableIPv6GatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_ipv6",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableSpokePreserveAsPathGatewayGroup enables preserve AS path for a spoke gateway group
func (c *Client) EnableSpokePreserveAsPathGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_spoke_preserve_as_path",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableSpokePreserveAsPathGatewayGroup disables preserve AS path for a spoke gateway group
func (c *Client) DisableSpokePreserveAsPathGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_spoke_preserve_as_path",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableSpokeLearnedCidrsApprovalGatewayGroup enables learned CIDRs approval for a spoke gateway group
func (c *Client) EnableSpokeLearnedCidrsApprovalGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":                 "enable_bgp_gateway_cidr_approval",
		"CID":                    c.CID,
		"group_name":             groupName,
		"learned_cidrs_approval": "on",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableSpokeLearnedCidrsApprovalGatewayGroup disables learned CIDRs approval for a spoke gateway group
func (c *Client) DisableSpokeLearnedCidrsApprovalGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":                 "disable_bgp_gateway_cidr_approval",
		"CID":                    c.CID,
		"group_name":             groupName,
		"learned_cidrs_approval": "off",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokePendingApprovedCidrsGatewayGroup updates the approved learned CIDRs for a spoke gateway group
func (c *Client) UpdateSpokePendingApprovedCidrsGatewayGroup(ctx context.Context, groupName string, approvedCidrs []string) error {
	action := "set_bgp_gateway_approved_cidr_rules"
	form := map[string]interface{}{
		"action":                 action,
		"CID":                    c.CID,
		"group_name":             groupName,
		"approved_learned_cidrs": approvedCidrs,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// EnablePrivateVpcDefaultRouteGatewayGroup enables private VPC default route for a gateway group
func (c *Client) EnablePrivateVpcDefaultRouteGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_private_vpc_default_route",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisablePrivateVpcDefaultRouteGatewayGroup disables private VPC default route for a gateway group
func (c *Client) DisablePrivateVpcDefaultRouteGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_private_vpc_default_route",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableSkipPublicRouteUpdateGatewayGroup enables skip public route table update for a gateway group
func (c *Client) EnableSkipPublicRouteUpdateGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_skip_public_route_table_update",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableSkipPublicRouteUpdateGatewayGroup disables skip public route table update for a gateway group
func (c *Client) DisableSkipPublicRouteUpdateGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_skip_public_route_table_update",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableAutoAdvertiseS2CCidrsGatewayGroup enables auto advertise S2C CIDRs for a gateway group
func (c *Client) EnableAutoAdvertiseS2CCidrsGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_auto_advertise_s2c_cidrs",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableAutoAdvertiseS2CCidrsGatewayGroup disables auto advertise S2C CIDRs for a gateway group
func (c *Client) DisableAutoAdvertiseS2CCidrsGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_auto_advertise_s2c_cidrs",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// SetSpokeBgpManualAdvertisedNetworksGatewayGroup sets the spoke BGP manual advertise CIDRs for a gateway group.
// The cidrs parameter should be a comma-separated list of CIDR strings (e.g., "10.0.0.0/16,10.1.0.0/16").
func (c *Client) SetSpokeBgpManualAdvertisedNetworksGatewayGroup(ctx context.Context, groupName, cidrs string) error {
	form := map[string]string{
		"action":                           "edit_aviatrix_spoke_advanced_config",
		"CID":                              c.CID,
		"group_name":                       groupName,
		"bgp_manual_spoke_advertise_cidrs": cidrs,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableSpokeOnpremRoutePropagationGatewayGroup disables spoke on-prem route propagation for a gateway group
func (c *Client) DisableSpokeOnpremRoutePropagationGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_spoke_onprem_route_propagation",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableSpokeOnpremRoutePropagationGatewayGroup enables spoke on-prem route propagation for a gateway group
func (c *Client) EnableSpokeOnpremRoutePropagationGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_spoke_onprem_route_propagation",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableGlobalVpcGatewayGroup enables global VPC for a gateway group (GCP)
func (c *Client) EnableGlobalVpcGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_global_vpc",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableGlobalVpcGatewayGroup disables global VPC for a gateway group (GCP)
func (c *Client) DisableGlobalVpcGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_global_vpc",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// ============================================================================
// Transit-Specific Gateway Group APIs
// ============================================================================

// EnableTransitPreserveAsPathGatewayGroup enables preserve AS path for a transit gateway group
func (c *Client) EnableTransitPreserveAsPathGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_transit_preserve_as_path",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableTransitPreserveAsPathGatewayGroup disables preserve AS path for a transit gateway group
func (c *Client) DisableTransitPreserveAsPathGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_transit_preserve_as_path",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableTransitLearnedCidrsApprovalGatewayGroup enables learned CIDRs approval for a transit gateway group
func (c *Client) EnableTransitLearnedCidrsApprovalGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":                 "enable_transit_learned_cidrs_approval",
		"CID":                    c.CID,
		"gateway_name":           groupName,
		"learned_cidrs_approval": "on",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableTransitLearnedCidrsApprovalGatewayGroup disables learned CIDRs approval for a transit gateway group
func (c *Client) DisableTransitLearnedCidrsApprovalGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":                 "disable_transit_learned_cidrs_approval",
		"CID":                    c.CID,
		"gateway_name":           groupName,
		"learned_cidrs_approval": "off",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateTransitPendingApprovedCidrsGatewayGroup updates the approved learned CIDRs for a transit gateway group
func (c *Client) UpdateTransitPendingApprovedCidrsGatewayGroup(ctx context.Context, groupName string, approvedCidrs []string) error {
	action := "set_transit_gateway_approved_cidr_rules"
	form := map[string]interface{}{
		"action":                 action,
		"CID":                    c.CID,
		"gateway_name":           groupName,
		"approved_learned_cidrs": approvedCidrs,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

// EnableConnectedTransitGatewayGroup enables connected transit for a gateway group
func (c *Client) EnableConnectedTransitGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_connected_transit",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableConnectedTransitGatewayGroup disables connected transit for a gateway group
func (c *Client) DisableConnectedTransitGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_connected_transit",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableSegmentationGatewayGroup enables segmentation for a gateway group
func (c *Client) EnableSegmentationGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_segmentation",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableSegmentationGatewayGroup disables segmentation for a gateway group
func (c *Client) DisableSegmentationGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_segmentation",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableAdvertiseTransitCidrGatewayGroup enables advertise transit CIDR for a gateway group
func (c *Client) EnableAdvertiseTransitCidrGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_advertise_transit_cidr",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableAdvertiseTransitCidrGatewayGroup disables advertise transit CIDR for a gateway group
func (c *Client) DisableAdvertiseTransitCidrGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_advertise_transit_cidr",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// SetTransitBgpManualAdvertisedNetworksGatewayGroup sets the transit BGP manual advertise CIDRs for a gateway group.
// The cidrs parameter should be a comma-separated list of CIDR strings (e.g., "10.0.0.0/16,10.1.0.0/16").
func (c *Client) SetTransitBgpManualAdvertisedNetworksGatewayGroup(ctx context.Context, groupName, cidrs string) error {
	form := map[string]string{
		"action":                           "edit_aviatrix_transit_advanced_config",
		"CID":                              c.CID,
		"gateway_name":                     groupName,
		"bgp_manual_spoke_advertise_cidrs": cidrs,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableHybridConnectionGatewayGroup enables hybrid connection for a gateway group
func (c *Client) EnableHybridConnectionGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_hybrid_connection",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableHybridConnectionGatewayGroup disables hybrid connection for a gateway group
func (c *Client) DisableHybridConnectionGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_hybrid_connection",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableTransitSummarizeCidrToTgwGatewayGroup enables transit summarize CIDR to TGW for a gateway group
func (c *Client) EnableTransitSummarizeCidrToTgwGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_transit_summarize_cidr_to_tgw",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableTransitSummarizeCidrToTgwGatewayGroup disables transit summarize CIDR to TGW for a gateway group
func (c *Client) DisableTransitSummarizeCidrToTgwGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_transit_summarize_cidr_to_tgw",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableMultiTierTransitGatewayGroup enables multi-tier transit for a gateway group
func (c *Client) EnableMultiTierTransitGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_multi_tier_transit",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableMultiTierTransitGatewayGroup disables multi-tier transit for a gateway group
func (c *Client) DisableMultiTierTransitGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_multi_tier_transit",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableS2cRxBalancingGatewayGroup enables S2C RX balancing for a gateway group
func (c *Client) EnableS2cRxBalancingGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_s2c_rx_balancing",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableS2cRxBalancingGatewayGroup disables S2C RX balancing for a gateway group
func (c *Client) DisableS2cRxBalancingGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_s2c_rx_balancing",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableGatewayLoadBalancerGatewayGroup enables AWS Gateway Load Balancer for a gateway group
func (c *Client) EnableGatewayLoadBalancerGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_gateway_load_balancer",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableGatewayLoadBalancerGatewayGroup disables AWS Gateway Load Balancer for a gateway group
func (c *Client) DisableGatewayLoadBalancerGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_gateway_load_balancer",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableFireNetGatewayGroup enables FireNet for a gateway group
func (c *Client) EnableFireNetGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_gateway_firenet_interfaces",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableFireNetGatewayGroup disables FireNet for a gateway group
func (c *Client) DisableFireNetGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_gateway_firenet_interfaces",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableTransitFireNetGatewayGroup enables Transit FireNet for a gateway group
func (c *Client) EnableTransitFireNetGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_gateway_for_transit_firenet",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// EnableTransitFireNetWithGWLBGatewayGroup enables Transit FireNet with AWS Gateway Load Balancer for a gateway group
func (c *Client) EnableTransitFireNetWithGWLBGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "enable_gateway_for_transit_firenet",
		"CID":          c.CID,
		"gateway_name": groupName,
		"mode":         "gwlb",
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// DisableTransitFireNetGatewayGroup disables Transit FireNet for a gateway group
func (c *Client) DisableTransitFireNetGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":       "disable_gateway_for_transit_firenet",
		"CID":          c.CID,
		"gateway_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
