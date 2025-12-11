package goaviatrix

import (
	"context"
	"fmt"
)

// GatewayGroup represents a Gateway Group resource
type GatewayGroup struct {
	Action string `form:"action,omitempty"`
	CID    string `form:"CID,omitempty"`

	// Required
	GroupName         string `form:"name,omitempty" json:"name,omitempty"`
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
	EnableSkipPublicRouteTableUpdate bool     `form:"skip_public_vpc_update_enabled,omitempty" json:"skip_public_vpc_update_enabled,omitempty"` // part of the spoke config in group schema
	Edge                             bool     `form:"edge,omitempty" json:"edge,omitempty"`

	// Feature Flags
	EnableJumboFrame   bool `form:"jumbo_frame,omitempty" json:"jumbo_frame,omitempty"`
	EnableNat          bool `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	EnableIPv6         bool `form:"enable_ipv6,omitempty" json:"enable_ipv6,omitempty"`
	EnableGroGso       bool `form:"enable_gro_gso,omitempty" json:"enable_gro_gso,omitempty"` // remove this, set using a get call client.GetGroGsoStatus(gw)
	EnableVpcDNSServer bool `form:"use_vpc_dns_server,omitempty" json:"use_vpc_dns_server,omitempty"`

	// BGP Configuration
	EnableBgp                    bool     `form:"bgp_enabled,omitempty" json:"bgp_enabled,omitempty"`
	LocalAsNumber                string   `form:"local_as_number,omitempty" json:"local_as_number,omitempty"`
	PrependAsPath                []string `form:"prepend_as_path,omitempty" json:"prepend_as_path,omitempty"`
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
	BgpSendCommunities   bool `form:"bgp_send_communities,omitempty" json:"bgp_send_communities,omitempty"`     // remove this, set using a get call client.GetGatewayBgpCommunities(gateway.GwName)
	BgpAcceptCommunities bool `form:"bgp_accept_communities,omitempty" json:"bgp_accept_communities,omitempty"` // remove this, set using a get call client.GetGatewayBgpCommunities(gateway.GwName)

	// BGP over LAN
	EnableBgpOverLan      bool `form:"enable_bgp_over_lan,omitempty" json:"enable_bgp_over_lan,omitempty"`
	BgpLanInterfacesCount int  `form:"bgp_over_lan_intf_cnt,omitempty" json:"bgp_over_lan_intf_cnt,omitempty"`

	// Learned CIDR Approval
	EnableLearnedCidrsApproval bool     `form:"enable_learned_cidrs_approval,omitempty" json:"enable_learned_cidrs_approval,omitempty"`
	LearnedCidrsApprovalMode   string   `form:"learned_cidrs_approval_mode,omitempty" json:"learned_cidrs_approval_mode,omitempty"`
	ApprovedLearnedCidrs       []string `form:"approved_learned_cidrs,omitempty" json:"approved_learned_cidrs,omitempty"`

	// Active-Standby
	EnableActiveStandby           bool `form:"enable_active_standby,omitempty" json:"enable_active_standby,omitempty"`
	EnableActiveStandbyPreemptive bool `form:"enable_active_standby_preemptive,omitempty" json:"enable_active_standby_preemptive,omitempty"`

	// AWS Specific
	InsaneMode          bool   `form:"insane_mode,omitempty" json:"insane_mode,omitempty"` // rename is EnableGroupHpe
	InsaneModeAz        string `form:"gateway_zone,omitempty" json:"gateway_zone,omitempty"`
	EnableEncryptVolume bool   `form:"gw_enc,omitempty" json:"gw_enc,omitempty"`
	CustomerManagedKeys string `form:"customer_managed_keys,omitempty" json:"customer_managed_keys,omitempty"`

	// GCP Specific
	EnableGlobalVpc bool `form:"global_vpc,omitempty" json:"global_vpc,omitempty"`

	// Computed (read-only)
	GwUUIDList                []string `json:"gw_uuid_list,omitempty"`
	VpcUUID                   string   `json:"vpc_uuid,omitempty"`
	VendorName                string   `json:"vendor_name,omitempty"`
	SoftwareVersion           string   `json:"gw_software_version,omitempty"`
	ImageVersion              string   `json:"gw_image_name,omitempty"`
	AzureEipNameResourceGroup string   `json:"azure_eip_name_resource_group,omitempty"` // remove this, set using azureEip
	BgpLanIPList              []string `json:"AzureBgpLanIpList,omitempty"`
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

	return c.PostAPIContext2(ctx, nil, spokeGroup.Action, spokeGroup, BasicCheck)
}

// GetSpokeGroup retrieves spoke group details
func (c *Client) GetGatewayGroup(ctx context.Context, groupName string) (*GatewayGroup, error) {
	form := map[string]string{
		"action":     "get_gateway_group_details",
		"CID":        c.CID,
		"group_name": groupName,
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

// DeleteGatewayGroup deletes a gateway group
func (c *Client) DeleteGatewayGroup(ctx context.Context, groupName string) error {
	form := map[string]string{
		"action":     "delete_gateway_group",
		"CID":        c.CID,
		"group_name": groupName,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
