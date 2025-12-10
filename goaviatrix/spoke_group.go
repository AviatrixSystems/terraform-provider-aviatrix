package goaviatrix

import (
	"context"
	"fmt"
)

// SpokeGroup represents a Spoke Gateway Group resource
type SpokeGroup struct {
	Action string `form:"action,omitempty"`
	CID    string `form:"CID,omitempty"`

	// Required
	GroupName         string `form:"group_name,omitempty" json:"group_name,omitempty"`
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
	IncludeCidr                      string   `form:"include_cidr,omitempty" json:"include_cidr,omitempty"`
	EnablePrivateVpcDefaultRoute     bool     `form:"enable_private_vpc_default_route,omitempty" json:"enable_private_vpc_default_route,omitempty"`
	EnableSkipPublicRouteTableUpdate bool     `form:"enable_skip_public_route_table_update,omitempty" json:"enable_skip_public_route_table_update,omitempty"`
	Edge                             bool     `form:"edge,omitempty" json:"edge,omitempty"`

	// Feature Flags
	EnableGroupHpe     bool `form:"enable_group_hpe,omitempty" json:"enable_group_hpe,omitempty"`
	EnableJumboFrame   bool `form:"enable_jumbo_frame,omitempty" json:"enable_jumbo_frame,omitempty"`
	EnableNat          bool `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	EnableIPv6         bool `form:"enable_ipv6,omitempty" json:"enable_ipv6,omitempty"`
	EnableGroGso       bool `form:"enable_gro_gso,omitempty" json:"enable_gro_gso,omitempty"`
	EnableVpcDnsServer bool `form:"enable_vpc_dns_server,omitempty" json:"enable_vpc_dns_server,omitempty"`

	// BGP Configuration
	EnableBgp                    bool     `form:"enable_bgp,omitempty" json:"enable_bgp,omitempty"`
	LocalAsNumber                string   `form:"local_as_number,omitempty" json:"local_as_number,omitempty"`
	PrependAsPath                []string `form:"prepend_as_path,omitempty" json:"prepend_as_path,omitempty"`
	DisableRoutePropagation      bool     `form:"disable_route_propagation,omitempty" json:"disable_route_propagation,omitempty"`
	SpokeBgpManualAdvertiseCidrs []string `form:"spoke_bgp_manual_advertise_cidrs,omitempty" json:"spoke_bgp_manual_advertise_cidrs,omitempty"`
	EnablePreserveAsPath         bool     `form:"enable_preserve_as_path,omitempty" json:"enable_preserve_as_path,omitempty"`
	EnableAutoAdvertiseS2cCidrs  bool     `form:"enable_auto_advertise_s2c_cidrs,omitempty" json:"enable_auto_advertise_s2c_cidrs,omitempty"`
	BgpEcmp                      bool     `form:"bgp_ecmp,omitempty" json:"bgp_ecmp,omitempty"`

	// BGP Timers
	BgpPollingTime               int `form:"bgp_polling_time,omitempty" json:"bgp_polling_time,omitempty"`
	BgpNeighborStatusPollingTime int `form:"bgp_neighbor_status_polling_time,omitempty" json:"bgp_neighbor_status_polling_time,omitempty"`
	BgpHoldTime                  int `form:"bgp_hold_time,omitempty" json:"bgp_hold_time,omitempty"`

	// BGP Communities
	BgpSendCommunities   bool `form:"bgp_send_communities,omitempty" json:"bgp_send_communities,omitempty"`
	BgpAcceptCommunities bool `form:"bgp_accept_communities,omitempty" json:"bgp_accept_communities,omitempty"`

	// BGP over LAN
	EnableBgpOverLan      bool `form:"enable_bgp_over_lan,omitempty" json:"enable_bgp_over_lan,omitempty"`
	BgpLanInterfacesCount int  `form:"bgp_lan_interfaces_count,omitempty" json:"bgp_lan_interfaces_count,omitempty"`

	// Learned CIDR Approval
	EnableLearnedCidrsApproval bool     `form:"enable_learned_cidrs_approval,omitempty" json:"enable_learned_cidrs_approval,omitempty"`
	LearnedCidrsApprovalMode   string   `form:"learned_cidrs_approval_mode,omitempty" json:"learned_cidrs_approval_mode,omitempty"`
	ApprovedLearnedCidrs       []string `form:"approved_learned_cidrs,omitempty" json:"approved_learned_cidrs,omitempty"`

	// Active-Standby
	EnableActiveStandby           bool `form:"enable_active_standby,omitempty" json:"enable_active_standby,omitempty"`
	EnableActiveStandbyPreemptive bool `form:"enable_active_standby_preemptive,omitempty" json:"enable_active_standby_preemptive,omitempty"`

	// AWS Specific
	InsaneMode          bool   `form:"insane_mode,omitempty" json:"insane_mode,omitempty"`
	InsaneModeAz        string `form:"insane_mode_az,omitempty" json:"insane_mode_az,omitempty"`
	EnableEncryptVolume bool   `form:"enable_encrypt_volume,omitempty" json:"enable_encrypt_volume,omitempty"`
	CustomerManagedKeys string `form:"customer_managed_keys,omitempty" json:"customer_managed_keys,omitempty"`

	// GCP Specific
	EnableGlobalVpc bool `form:"enable_global_vpc,omitempty" json:"enable_global_vpc,omitempty"`

	// Computed (read-only)
	GwUuidList                []string `json:"gw_uuid_list,omitempty"`
	VpcUuid                   string   `json:"vpc_uuid,omitempty"`
	VendorName                string   `json:"vendor_name,omitempty"`
	SoftwareVersion           string   `json:"software_version,omitempty"`
	ImageVersion              string   `json:"image_version,omitempty"`
	AzureEipNameResourceGroup string   `json:"azure_eip_name_resource_group,omitempty"`
	BgpLanIpList              []string `json:"bgp_lan_ip_list,omitempty"`
}

// SpokeGroupResp represents the API response for get spoke group
type SpokeGroupResp struct {
	Return  bool       `json:"return"`
	Results SpokeGroup `json:"results"`
	Reason  string     `json:"reason"`
}

// CreateSpokeGroup creates a new spoke gateway group
func (c *Client) CreateSpokeGroup(ctx context.Context, spokeGroup *SpokeGroup) error {
	spokeGroup.Action = "create_gateway_group"
	spokeGroup.CID = c.CID

	return c.PostAPIContext2(ctx, nil, spokeGroup.Action, spokeGroup, BasicCheck)
}

// GetSpokeGroup retrieves spoke group details
func (c *Client) GetSpokeGroup(ctx context.Context, groupName string) (*SpokeGroup, error) {
	form := map[string]string{
		"action":     "get_gateway_group_details",
		"CID":        c.CID,
		"group_name": groupName,
	}

	var resp SpokeGroupResp
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if !resp.Return {
		return nil, fmt.Errorf("failed to get spoke group: %s", resp.Reason)
	}

	return &resp.Results, nil
}

// UpdateSpokeGroup updates an existing spoke gateway group
func (c *Client) UpdateSpokeGroup(ctx context.Context, spokeGroup *SpokeGroup) error {
	spokeGroup.Action = "update_gateway_group"
	spokeGroup.CID = c.CID

	return c.PostAPIContext2(ctx, nil, spokeGroup.Action, spokeGroup, BasicCheck)
}

// DeleteSpokeGroup deletes a spoke gateway group
func (c *Client) DeleteSpokeGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":     "delete_gateway_group",
		"CID":        c.CID,
		"group_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// ============================================================================
// Individual Update APIs for Spoke Group
// ============================================================================

// UpdateSpokeGroupGwSize updates gateway size - API: edit_gw_config
func (c *Client) UpdateSpokeGroupGwSize(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":     "edit_gw_config",
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
		"gw_size":    spokeGroup.GroupInstanceSize,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupBgpAcceptCommunities updates BGP accept communities - API: set_gateway_accept_bgp_communities_override
func (c *Client) UpdateSpokeGroupBgpAcceptCommunities(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":     "set_gateway_accept_bgp_communities_override",
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
		"enable":     fmt.Sprintf("%t", spokeGroup.BgpAcceptCommunities),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupBgpSendCommunities updates BGP send communities - API: set_gateway_send_bgp_communities_override
func (c *Client) UpdateSpokeGroupBgpSendCommunities(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":     "set_gateway_send_bgp_communities_override",
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
		"enable":     fmt.Sprintf("%t", spokeGroup.BgpSendCommunities),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupSnat updates NAT/SNAT - API: enable_snat / disable_snat
func (c *Client) UpdateSpokeGroupSnat(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_snat"
	if spokeGroup.EnableNat {
		action = "enable_snat"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupVpcDnsServer updates VPC DNS server - API: enable_vpc_dns_server / disable_vpc_dns_server
func (c *Client) UpdateSpokeGroupVpcDnsServer(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_vpc_dns_server"
	if spokeGroup.EnableVpcDnsServer {
		action = "enable_vpc_dns_server"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupBgpPollingTime updates BGP polling time - API: change_bgp_polling_time
func (c *Client) UpdateSpokeGroupBgpPollingTime(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":           "change_bgp_polling_time",
		"CID":              c.CID,
		"group_name":       spokeGroup.GroupName,
		"bgp_polling_time": fmt.Sprintf("%d", spokeGroup.BgpPollingTime),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupBgpNeighborStatusPollingTime updates BGP neighbor status polling time - API: change_bgp_neighbor_status_polling_time
func (c *Client) UpdateSpokeGroupBgpNeighborStatusPollingTime(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":                           "change_bgp_neighbor_status_polling_time",
		"CID":                              c.CID,
		"group_name":                       spokeGroup.GroupName,
		"bgp_neighbor_status_polling_time": fmt.Sprintf("%d", spokeGroup.BgpNeighborStatusPollingTime),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupBgpHoldTime updates BGP hold time - API: change_bgp_hold_time
func (c *Client) UpdateSpokeGroupBgpHoldTime(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":        "change_bgp_hold_time",
		"CID":           c.CID,
		"group_name":    spokeGroup.GroupName,
		"bgp_hold_time": fmt.Sprintf("%d", spokeGroup.BgpHoldTime),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupLocalAsNumber updates local AS number - API: edit_transit_local_as_number
func (c *Client) UpdateSpokeGroupLocalAsNumber(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":          "edit_transit_local_as_number",
		"CID":             c.CID,
		"group_name":      spokeGroup.GroupName,
		"local_as_number": spokeGroup.LocalAsNumber,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupPrependAsPath updates prepend AS path - API: edit_aviatrix_transit_advanced_config
func (c *Client) UpdateSpokeGroupPrependAsPath(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]interface{}{
		"action":          "edit_aviatrix_transit_advanced_config",
		"CID":             c.CID,
		"group_name":      spokeGroup.GroupName,
		"prepend_as_path": spokeGroup.PrependAsPath,
	}
	return c.PostAPIContext2(ctx, nil, form["action"].(string), form, BasicCheck)
}

// UpdateSpokeGroupBgpEcmp updates BGP ECMP - API: enable_bgp_ecmp / disable_bgp_ecmp
func (c *Client) UpdateSpokeGroupBgpEcmp(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_bgp_ecmp"
	if spokeGroup.BgpEcmp {
		action = "enable_bgp_ecmp"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupActiveStandby updates active-standby - API: enable_active_standby / disable_active_standby
func (c *Client) UpdateSpokeGroupActiveStandby(ctx context.Context, spokeGroup *SpokeGroup) error {
	if spokeGroup.EnableActiveStandby {
		form := map[string]string{
			"action":     "enable_active_standby",
			"CID":        c.CID,
			"group_name": spokeGroup.GroupName,
			"preemptive": fmt.Sprintf("%t", spokeGroup.EnableActiveStandbyPreemptive),
		}
		return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
	}
	form := map[string]string{
		"action":     "disable_active_standby",
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupJumboFrame updates jumbo frame - API: enable_jumbo_frame / disable_jumbo_frame
func (c *Client) UpdateSpokeGroupJumboFrame(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_jumbo_frame"
	if spokeGroup.EnableJumboFrame {
		action = "enable_jumbo_frame"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupGroGso updates GRO/GSO - API: enable_gro_gso / disable_gro_gso
func (c *Client) UpdateSpokeGroupGroGso(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_gro_gso"
	if spokeGroup.EnableGroGso {
		action = "enable_gro_gso"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupS2cRxBalancing updates S2C RX balancing - API: enable_s2c_rx_balancing / disable_s2c_rx_balancing
func (c *Client) UpdateSpokeGroupS2cRxBalancing(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_s2c_rx_balancing"
	if spokeGroup.S2cRxBalancing {
		action = "enable_s2c_rx_balancing"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupIPv6 updates IPv6 - API: enable_ipv6 / disable_ipv6
func (c *Client) UpdateSpokeGroupIPv6(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_ipv6"
	if spokeGroup.EnableIPv6 {
		action = "enable_ipv6"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupPreserveAsPath updates preserve AS path - API: enable_spoke_preserve_as_path / disable_spoke_preserve_as_path
func (c *Client) UpdateSpokeGroupPreserveAsPath(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_spoke_preserve_as_path"
	if spokeGroup.EnablePreserveAsPath {
		action = "enable_spoke_preserve_as_path"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupLearnedCidrsApproval updates learned CIDRs approval - API: enable_bgp_gateway_cidr_approval / disable_bgp_gateway_cidr_approval
func (c *Client) UpdateSpokeGroupLearnedCidrsApproval(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_bgp_gateway_cidr_approval"
	if spokeGroup.EnableLearnedCidrsApproval {
		action = "enable_bgp_gateway_cidr_approval"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupApprovedLearnedCidrs updates approved learned CIDRs - API: set_bgp_gateway_approved_cidr_rules
func (c *Client) UpdateSpokeGroupApprovedLearnedCidrs(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]interface{}{
		"action":                 "set_bgp_gateway_approved_cidr_rules",
		"CID":                    c.CID,
		"group_name":             spokeGroup.GroupName,
		"approved_learned_cidrs": spokeGroup.ApprovedLearnedCidrs,
	}
	return c.PostAPIContext2(ctx, nil, form["action"].(string), form, BasicCheck)
}

// UpdateSpokeGroupPrivateVpcDefaultRoute updates private VPC default route - API: enable_private_vpc_default_route / disable_private_vpc_default_route
func (c *Client) UpdateSpokeGroupPrivateVpcDefaultRoute(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_private_vpc_default_route"
	if spokeGroup.EnablePrivateVpcDefaultRoute {
		action = "enable_private_vpc_default_route"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupSkipPublicRouteTableUpdate updates skip public route table update - API: enable_skip_public_route_table_update / disable_skip_public_route_table_update
func (c *Client) UpdateSpokeGroupSkipPublicRouteTableUpdate(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_skip_public_route_table_update"
	if spokeGroup.EnableSkipPublicRouteTableUpdate {
		action = "enable_skip_public_route_table_update"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupAutoAdvertiseS2cCidrs updates auto advertise S2C CIDRs - API: enable_auto_advertise_s2c_cidrs / disable_auto_advertise_s2c_cidrs
func (c *Client) UpdateSpokeGroupAutoAdvertiseS2cCidrs(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_auto_advertise_s2c_cidrs"
	if spokeGroup.EnableAutoAdvertiseS2cCidrs {
		action = "enable_auto_advertise_s2c_cidrs"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupManualAdvertiseCidrs updates manual advertise CIDRs - API: edit_aviatrix_spoke_advanced_config
func (c *Client) UpdateSpokeGroupManualAdvertiseCidrs(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]interface{}{
		"action":                           "edit_aviatrix_spoke_advanced_config",
		"CID":                              c.CID,
		"group_name":                       spokeGroup.GroupName,
		"spoke_bgp_manual_advertise_cidrs": spokeGroup.SpokeBgpManualAdvertiseCidrs,
	}
	return c.PostAPIContext2(ctx, nil, form["action"].(string), form, BasicCheck)
}

// UpdateSpokeGroupRoutePropagation updates route propagation - API: enable_spoke_onprem_route_propagation / disable_spoke_onprem_route_propagation
func (c *Client) UpdateSpokeGroupRoutePropagation(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "enable_spoke_onprem_route_propagation"
	if spokeGroup.DisableRoutePropagation {
		action = "disable_spoke_onprem_route_propagation"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupGlobalVpc updates global VPC - API: enable_global_vpc / disable_global_vpc
func (c *Client) UpdateSpokeGroupGlobalVpc(ctx context.Context, spokeGroup *SpokeGroup) error {
	action := "disable_global_vpc"
	if spokeGroup.EnableGlobalVpc {
		action = "enable_global_vpc"
	}
	form := map[string]string{
		"action":     action,
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupEncryptVolume updates encrypt volume - API: encrypt_gateway_volume
func (c *Client) UpdateSpokeGroupEncryptVolume(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":     "encrypt_gateway_volume",
		"CID":        c.CID,
		"group_name": spokeGroup.GroupName,
	}
	if spokeGroup.CustomerManagedKeys != "" {
		form["customer_managed_keys"] = spokeGroup.CustomerManagedKeys
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupBgpLanInterfacesCount updates BGP LAN interfaces count - API: change_bgp_over_lan_intf_cnt
func (c *Client) UpdateSpokeGroupBgpLanInterfacesCount(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":                   "change_bgp_over_lan_intf_cnt",
		"CID":                      c.CID,
		"group_name":               spokeGroup.GroupName,
		"bgp_lan_interfaces_count": fmt.Sprintf("%d", spokeGroup.BgpLanInterfacesCount),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

// UpdateSpokeGroupHpe updates HPE (High Performance Encryption) - API: TBD
func (c *Client) UpdateSpokeGroupHpe(ctx context.Context, spokeGroup *SpokeGroup) error {
	form := map[string]string{
		"action":           "update_gateway_group",
		"CID":              c.CID,
		"group_name":       spokeGroup.GroupName,
		"enable_group_hpe": fmt.Sprintf("%t", spokeGroup.EnableGroupHpe),
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
