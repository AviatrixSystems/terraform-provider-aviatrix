package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"strings"
)

type EdgeCSP struct {
	Action                             string `json:"action,omitempty"`
	CID                                string `json:"CID,omitempty"`
	AccountName                        string `json:"account_name,omitempty"`
	GwName                             string `json:"name,omitempty"`
	SiteId                             string `json:"site_id,omitempty"`
	ProjectUuid                        string `json:"project_uuid,omitempty"`
	ComputeNodeUuid                    string `json:"compute_node_uuid,omitempty"`
	TemplateUuid                       string `json:"template_uuid,omitempty"`
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip,omitempty"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network,omitempty"`
	DnsServerIp                        string `json:"dns_server_ip,omitempty"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary,omitempty"`
	Dhcp                               bool   `json:"dhcp,omitempty"`
	EnableEdgeActiveStandby            bool   `json:"enable_active_standby,omitempty"`
	DisableEdgeActiveStandby           bool   `json:"disable_active_standby,omitempty"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"enable_active_standby_preemptive,omitempty"`
	DisableEdgeActiveStandbyPreemptive bool   `json:"disable_active_standby_preemptive,omitempty"`
	LocalAsNumber                      string `json:"local_as_number,omitempty"`
	PrependAsPath                      []string
	PrependAsPathReturn                string   `json:"prepend_as_path,omitempty"`
	IncludeCidrList                    []string `json:"include_cidr_list,omitempty"`
	EnableLearnedCidrsApproval         bool     `json:"enable_learned_cidrs_approval,omitempty"`
	ApprovedLearnedCidrs               []string `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string `json:"bgp_manual_spoke_advertise_cidrs,omitempty"`
	EnablePreserveAsPath               bool     `json:"preserve_as_path,omitempty"`
	BgpPollingTime                     int      `json:"bgp_polling_time,omitempty"`
	BgpHoldTime                        int      `json:"bgp_hold_time,omitempty"`
	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing,omitempty"`
	EnableJumboFrame                   bool     `json:"jumbo_frame,omitempty"`
	Latitude                           string
	Longitude                          string
	RxQueueSize                        string `json:"rx_queue_size,omitempty"`
	State                              string `json:"vpc_state,omitempty"`
	NoProgressBar                      bool   `json:"no_progress_bar,omitempty"`
	WanInterface                       string `json:"wan_ifname,omitempty"`
	LanInterface                       string `json:"lan_ifname,omitempty"`
	MgmtInterface                      string `json:"mgmt_ifname,omitempty"`
	InterfaceList                      []*Interface
	Interfaces                         string `json:"interfaces,omitempty"`
	VlanList                           []*Vlan
	Vlan                               string `json:"vlan,omitempty"`
	DnsProfileName                     string `json:"dns_profile_name,omitempty"`
	EnableSingleIpSnat                 bool
	EnableAutoAdvertiseLanCidrs        string `json:"auto_advertise_lan_cidrs,omitempty"`
	LanInterfaceIpPrefix               string
}

type Interface struct {
	IfName        string  `json:"ifname"`
	Type          string  `json:"type"`
	Bandwidth     int     `json:"bandwidth"`
	PublicIp      string  `json:"public_ip"`
	Tag           string  `json:"tag"`
	Dhcp          bool    `json:"dhcp"`
	IpAddr        string  `json:"ipaddr"`
	GatewayIp     string  `json:"gateway_ip"`
	DnsPrimary    string  `json:"dns_primary"`
	DnsSecondary  string  `json:"dns_secondary"`
	SubInterfaces []*Vlan `json:"subinterfaces"`
	VrrpState     bool    `json:"vrrp_state"`
	VirtualIp     string  `json:"virtual_ip"`
}

type Vlan struct {
	ParentInterface string `json:"parent_interface"`
	VlanId          string `json:"vlan_id"`
	IpAddr          string `json:"ipaddr"`
	GatewayIp       string `json:"gateway_ip"`
	PeerIpAddr      string `json:"peer_ipaddr"`
	PeerGatewayIp   string `json:"peer_gateway_ip"`
	VirtualIp       string `json:"virtual_ip"`
	Tag             string `json:"tag"`
}

type EdgeCSPResp struct {
	AccountName                        string `json:"account_name"`
	GwName                             string `json:"gw_name"`
	SiteId                             string `json:"vpc_id"`
	ProjectUuid                        string `json:"edge_csp_project_uuid"`
	ComputeNodeUuid                    string `json:"edge_csp_compute_node_uuid"`
	TemplateUuid                       string `json:"edge_csp_template_uuid"`
	ManagementInterfaceConfig          string
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network"`
	LanInterfaceIpPrefix               string `json:"lan_ip"`
	ManagementInterfaceIpPrefix        string `json:"mgmt_ip"`
	ManagementDefaultGatewayIp         string `json:"mgmt_default_gateway"`
	DnsServerIp                        string `json:"dns_server_ip"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary"`
	Dhcp                               bool   `json:"dhcp"`
	ActiveStandby                      string `json:"active_standby"`
	EnableEdgeActiveStandby            bool   `json:"edge_active_standby"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"edge_active_standby_preemptive"`
	LocalAsNumber                      string `json:"local_as_number"`
	PrependAsPath                      []string
	PrependAsPathReturn                string       `json:"prepend_as_path"`
	IncludeCidrList                    []string     `json:"include_cidr_list"`
	EnableLearnedCidrsApproval         bool         `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs               []string     `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string     `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath               bool         `json:"preserve_as_path"`
	BgpPollingTime                     int          `json:"bgp_polling_time"`
	BgpHoldTime                        int          `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting        bool         `json:"edge_transitive_routing"`
	EnableJumboFrame                   bool         `json:"jumbo_frame"`
	Latitude                           float64      `json:"latitude"`
	Longitude                          float64      `json:"longitude"`
	WanPublicIp                        string       `json:"public_ip"`
	PrivateIP                          string       `json:"private_ip"`
	RxQueueSize                        string       `json:"rx_queue_size"`
	State                              string       `json:"vpc_state"`
	WanInterface                       []string     `json:"edge_csp_wan_ifname"`
	LanInterface                       []string     `json:"edge_csp_lan_ifname"`
	MgmtInterface                      []string     `json:"edge_csp_mgmt_ifname"`
	InterfaceList                      []*Interface `json:"interfaces"`
	DnsProfileName                     string       `json:"dns_profile_name"`
	EnableNat                          string       `json:"enable_nat"`
	SnatMode                           string       `json:"snat_target"`
	EnableAutoAdvertiseLanCidrs        bool         `json:"auto_advertise_lan_cidrs"`
}

type EdgeCSPListResp struct {
	Return  bool          `json:"return"`
	Results []EdgeCSPResp `json:"results"`
	Reason  string        `json:"reason"`
}

func (c *Client) CreateEdgeCSP(ctx context.Context, edgeCSP *EdgeCSP) error {
	edgeCSP.Action = "create_edge_csp_instance"
	edgeCSP.CID = c.CID
	edgeCSP.NoProgressBar = true

	interfaces, err := json.Marshal(edgeCSP.InterfaceList)
	if err != nil {
		return err
	}

	edgeCSP.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if edgeCSP.VlanList == nil || len(edgeCSP.VlanList) == 0 {
		edgeCSP.VlanList = []*Vlan{}
	}

	vlan, err := json.Marshal(edgeCSP.VlanList)
	if err != nil {
		return err
	}

	edgeCSP.Vlan = b64.StdEncoding.EncodeToString(vlan)

	err = c.PostAPIContext2(ctx, nil, edgeCSP.Action, edgeCSP, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeCSP(ctx context.Context, gwName string) (*EdgeCSPResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeCSPListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeCSPList := data.Results
	for _, edgeCSP := range edgeCSPList {
		if edgeCSP.GwName == gwName {
			for _, p := range strings.Split(edgeCSP.PrependAsPathReturn, " ") {
				if p != "" {
					edgeCSP.PrependAsPath = append(edgeCSP.PrependAsPath, p)
				}
			}

			return &edgeCSP, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteEdgeCSP(ctx context.Context, accountName, name string) error {
	form := map[string]string{
		"action":       "delete_edge_csp_instance",
		"CID":          c.CID,
		"account_name": accountName,
		"name":         name,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) UpdateEdgeCSP(ctx context.Context, edgeCSP *EdgeCSP) error {
	edgeCSP.Action = "update_edge_gateway"
	edgeCSP.CID = c.CID

	interfaces, err := json.Marshal(edgeCSP.InterfaceList)
	if err != nil {
		return err
	}

	edgeCSP.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if edgeCSP.VlanList == nil || len(edgeCSP.VlanList) == 0 {
		edgeCSP.VlanList = []*Vlan{}
	}

	vlan, err := json.Marshal(edgeCSP.VlanList)
	if err != nil {
		return err
	}

	edgeCSP.Vlan = b64.StdEncoding.EncodeToString(vlan)

	return c.PostAPIContext2(ctx, nil, edgeCSP.Action, edgeCSP, BasicCheck)
}
