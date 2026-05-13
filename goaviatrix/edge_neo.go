package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"strings"
)

type EdgeNEO struct {
	Action                             string `json:"action,omitempty"`
	CID                                string `json:"CID,omitempty"`
	AccountName                        string `json:"account_name,omitempty"`
	GwName                             string `json:"name,omitempty"`
	SiteId                             string `json:"site_id,omitempty"`
	DeviceId                           string `json:"device_id,omitempty"`
	GwSize                             string `json:"gw_resource_size,omitempty"`
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip,omitempty"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network,omitempty"`
	DnsServerIp                        string `json:"dns_server_ip,omitempty"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary,omitempty"`
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
	BgpBfdPollingTime                  int      `json:"bgp_neighbor_status_polling_time,omitempty"`
	BgpHoldTime                        int      `json:"bgp_hold_time,omitempty"`
	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing,omitempty"`
	EnableJumboFrame                   bool     `json:"jumbo_frame,omitempty"`
	Latitude                           string
	Longitude                          string
	RxQueueSize                        string `json:"rx_queue_size,omitempty"`
	State                              string `json:"vpc_state,omitempty"`
	NoProgressBar                      bool   `json:"no_progress_bar,omitempty"`
	WanInterface                       string `json:"wan_ifnames,omitempty"`
	LanInterface                       string `json:"lan_ifnames,omitempty"`
	MgmtInterface                      string `json:"mgmt_ifnames,omitempty"`
	InterfaceList                      []*EdgeNEOInterface
	Interfaces                         string `json:"interfaces,omitempty"`
	VlanList                           []*EdgeNEOVlan
	Vlan                               string `json:"vlan,omitempty"`
	EnableSingleIpSnat                 bool
	EnableAutoAdvertiseLanCidrs        string `json:"auto_advertise_lan_cidrs,omitempty"`
	LanInterfaceIpPrefix               string
	DirectAttachLan                    bool     `json:"direct_attach_lan,omitempty"`
	AdvertisedCidrList                 []string `json:"advertise_cidr_list,omitempty"`
}

type EdgeNEOInterface struct {
	IfName        string         `json:"ifname"`
	Type          string         `json:"type"`
	PublicIp      string         `json:"public_ip,omitempty"`
	Tag           string         `json:"tag,omitempty"`
	Dhcp          bool           `json:"dhcp,omitempty"`
	IpAddr        string         `json:"ipaddr,omitempty"`
	GatewayIp     string         `json:"gateway_ip,omitempty"`
	DnsPrimary    string         `json:"dns_primary,omitempty"`
	DnsSecondary  string         `json:"dns_secondary,omitempty"`
	SubInterfaces []*EdgeNEOVlan `json:"subinterfaces,omitempty"`
	VrrpState     bool           `json:"vrrp_state,omitempty"`
	VirtualIp     string         `json:"virtual_ip,omitempty"`
	IPv6Addr      string         `json:"v6_ipaddr,omitempty"`
	GatewayIPv6   string         `json:"gateway_ipv6_ip,omitempty"`
}

type EdgeNEOVlan struct {
	ParentInterface string `json:"parent_interface"`
	VlanId          string `json:"vlan_id"`
	IpAddr          string `json:"ipaddr"`
	GatewayIp       string `json:"gateway_ip,omitempty"`
	PeerIpAddr      string `json:"peer_ipaddr,omitempty"`
	PeerGatewayIp   string `json:"peer_gateway_ip,omitempty"`
	VirtualIp       string `json:"virtual_ip,omitempty"`
	Tag             string `json:"tag,omitempty"`
}

type EdgeNEOResp struct {
	AccountName                        string `json:"account_name"`
	GwName                             string `json:"gw_name"`
	SiteId                             string `json:"vpc_id"`
	DeviceId                           string `json:"edge_csp_device_id"`
	GwSize                             string `json:"edge_csp_gateway_size"`
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network"`
	DnsServerIp                        string `json:"dns_server_ip"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary"`
	ActiveStandby                      string `json:"active_standby"`
	EnableEdgeActiveStandby            bool   `json:"edge_active_standby"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"edge_active_standby_preemptive"`
	LocalAsNumber                      string `json:"local_as_number"`
	PrependAsPath                      []string
	PrependAsPathReturn                string              `json:"prepend_as_path"`
	IncludeCidrList                    []string            `json:"include_cidr_list"`
	EnableLearnedCidrsApproval         bool                `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs               []string            `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string            `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath               bool                `json:"preserve_as_path"`
	BgpPollingTime                     int                 `json:"bgp_polling_time"`
	BgpBfdPollingTime                  int                 `json:"bgp_neighbor_status_polling_time"`
	BgpHoldTime                        int                 `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting        bool                `json:"edge_transitive_routing"`
	EnableJumboFrame                   bool                `json:"jumbo_frame"`
	Latitude                           float64             `json:"latitude"`
	Longitude                          float64             `json:"longitude"`
	WanPublicIp                        string              `json:"public_ip"`
	PrivateIP                          string              `json:"private_ip"`
	RxQueueSize                        string              `json:"rx_queue_size"`
	State                              string              `json:"vpc_state"`
	WanInterface                       []string            `json:"edge_csp_wan_ifname"`
	LanInterface                       []string            `json:"edge_csp_lan_ifname"`
	MgmtInterface                      []string            `json:"edge_csp_mgmt_ifname"`
	InterfaceList                      []*EdgeNEOInterface `json:"interfaces"`
	EnableNat                          string              `json:"enable_nat"`
	SnatMode                           string              `json:"snat_target"`
	EnableAutoAdvertiseLanCidrs        bool                `json:"auto_advertise_lan_cidrs"`
	AdvertisedCidrList                 []string            `json:"advertise_cidr_list,omitempty"`
}

type EdgeNEOListResp struct {
	Return  bool          `json:"return"`
	Results []EdgeNEOResp `json:"results"`
	Reason  string        `json:"reason"`
}

func (c *Client) CreateEdgeNEO(ctx context.Context, edgeNEO *EdgeNEO) error {
	edgeNEO.Action = "create_edge_csp_gateway"
	edgeNEO.CID = c.CID
	edgeNEO.NoProgressBar = true
	edgeNEO.DirectAttachLan = true

	interfaces, err := json.Marshal(edgeNEO.InterfaceList)
	if err != nil {
		return err
	}

	edgeNEO.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if len(edgeNEO.VlanList) == 0 {
		edgeNEO.VlanList = []*EdgeNEOVlan{}
	}

	vlan, err := json.Marshal(edgeNEO.VlanList)
	if err != nil {
		return err
	}

	edgeNEO.Vlan = b64.StdEncoding.EncodeToString(vlan)

	err = c.PostAPIContext2(ctx, nil, edgeNEO.Action, edgeNEO, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeNEO(ctx context.Context, gwName string) (*EdgeNEOResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeNEOListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeNEOList := data.Results
	for _, edgeNEO := range edgeNEOList {
		if edgeNEO.GwName == gwName {
			for _, p := range strings.Split(edgeNEO.PrependAsPathReturn, " ") {
				if p != "" {
					edgeNEO.PrependAsPath = append(edgeNEO.PrependAsPath, p)
				}
			}

			return &edgeNEO, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteEdgeNEO(ctx context.Context, accountName, name string) error {
	form := map[string]string{
		"action":       "delete_edge_csp_gateway",
		"CID":          c.CID,
		"account_name": accountName,
		"name":         name,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) UpdateEdgeNEO(ctx context.Context, edgeNEO *EdgeNEO) error {
	edgeNEO.Action = "update_edge_gateway"
	edgeNEO.CID = c.CID

	interfaces, err := json.Marshal(edgeNEO.InterfaceList)
	if err != nil {
		return err
	}

	edgeNEO.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if len(edgeNEO.VlanList) == 0 {
		edgeNEO.VlanList = []*EdgeNEOVlan{}
	}

	vlan, err := json.Marshal(edgeNEO.VlanList)
	if err != nil {
		return err
	}

	edgeNEO.Vlan = b64.StdEncoding.EncodeToString(vlan)

	return c.PostAPIContext2(ctx, nil, edgeNEO.Action, edgeNEO, BasicCheck)
}
