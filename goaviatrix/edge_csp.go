package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	ManagementInterfaceConfig          string
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip,omitempty"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network,omitempty"`
	LanInterfaceIpPrefix               string `json:"lan_ip,omitempty"`
	ManagementInterfaceIpPrefix        string `json:"mgmt_ip,omitempty"`
	ManagementDefaultGatewayIp         string `json:"mgmt_default_gateway,omitempty"`
	DnsServerIp                        string `json:"dns_server_ip,omitempty"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary,omitempty"`
	Dhcp                               bool   `json:"dhcp,omitempty"`
	EnableEdgeActiveStandby            bool   `json:"enable_active_standby"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"enable_active_standby_preemptive"`
	LocalAsNumber                      string `json:"local_as_number"`
	PrependAsPath                      []string
	PrependAsPathReturn                string   `json:"prepend_as_path"`
	IncludeCidrList                    []string `json:"include_cidr_list"`
	EnableLearnedCidrsApproval         bool     `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs               []string `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath               bool     `json:"preserve_as_path"`
	BgpPollingTime                     int      `json:"bgp_polling_time"`
	BgpHoldTime                        int      `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing"`
	EnableJumboFrame                   bool     `json:"jumbo_frame"`
	Latitude                           string
	Longitude                          string
	LatitudeReturn                     float64      `json:"latitude"`
	LongitudeReturn                    float64      `json:"longitude"`
	WanPublicIp                        string       `json:"wan_discovery_ip"`
	PrivateIP                          string       `json:"private_ip"`
	RxQueueSize                        string       `json:"rx_queue_size"`
	State                              string       `json:"vpc_state"`
	NoProgressBar                      bool         `json:"no_progress_bar,omitempty"`
	WanInterface                       string       `json:"wan_ifname"`
	LanInterface                       string       `json:"lan_ifname"`
	MgmtInterface                      string       `json:"mgmt_ifname"`
	InterfaceList                      []*Interface `json:"interfaces"`
	VlanList                           []*Vlan      `json:"vlan"`
	DnsProfileName                     string       `json:"dns_profile_name"`
	SingleIpSnat                       bool
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
	AdminState    string  `json:"admin_state"`
	SubInterfaces []*Vlan `json:"subinterfaces"`
	VrrpState     bool    `json:"vrrp_state"`
	VirtualIp     string  `json:"virtual_ip"`
}

type Vlan struct {
	ParentInterface string `json:"parent_interface"`
	VlanId          string `json:"vlan_id"`
	IpAddr          string `json:"ipaddr"`
	GatewayIp       string `json:"gateway_ip"`
	AdminState      string `json:"admin_state"`
	PeerIpAddr      string `json:"peer_ipaddr"`
	PeerGatewayIp   string `json:"peer_gateway_ip"`
	VirtualIp       string `json:"virtual_ip"`
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
	PrependAsPathReturn                string   `json:"prepend_as_path"`
	IncludeCidrList                    []string `json:"include_cidr_list"`
	EnableLearnedCidrsApproval         bool     `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs               []string `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath               bool     `json:"preserve_as_path"`
	BgpPollingTime                     int      `json:"bgp_polling_time"`
	BgpHoldTime                        int      `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing"`
	EnableJumboFrame                   bool     `json:"jumbo_frame"`
	Latitude                           string
	Longitude                          string
	LatitudeReturn                     float64      `json:"latitude"`
	LongitudeReturn                    float64      `json:"longitude"`
	WanPublicIp                        string       `json:"public_ip"`
	PrivateIP                          string       `json:"private_ip"`
	RxQueueSize                        string       `json:"rx_queue_size"`
	State                              string       `json:"vpc_state"`
	WanInterface                       []string     `json:"edge_csp_wan_ifname"`
	LanInterface                       []string     `json:"edge_csp_lan_ifname"`
	MgmtInterface                      []string     `json:"edge_csp_mgmt_ifname"`
	InterfaceList                      []*Interface `json:"interfaces"`
	DnsProfileName                     string       `json:"dns_profile_name"`
	SingleIpSnat                       bool         `json:"nat_enabled"`
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

	if edgeCSP.ManagementInterfaceConfig == "DHCP" {
		edgeCSP.Dhcp = true
	}

	err := c.PostAPIContext2(ctx, nil, edgeCSP.Action, edgeCSP, BasicCheck)
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
	form := map[string]string{
		"action": "update_edge_gateway",
		"CID":    c.CID,
		"name":   edgeCSP.GwName,
	}

	if edgeCSP.InterfaceList != nil && len(edgeCSP.InterfaceList) != 0 {
		interfaces, err := json.Marshal(edgeCSP.InterfaceList)
		if err != nil {
			return err
		}

		form["interfaces"] = b64.StdEncoding.EncodeToString(interfaces)
	}

	if edgeCSP.VlanList != nil && len(edgeCSP.VlanList) != 0 {
		vlan, err := json.Marshal(edgeCSP.VlanList)
		if err != nil {
			return err
		}

		form["vlan"] = b64.StdEncoding.EncodeToString(vlan)
	}

	if edgeCSP.DnsProfileName != "" {
		form["dns_profile_name"] = edgeCSP.DnsProfileName
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func DiffSuppressFuncInterfaces(k, old, new string, d *schema.ResourceData) bool {
	ifOld, ifNew := d.GetChange("interfaces")
	var interfacesOld []map[string]interface{}

	for _, if0 := range ifOld.([]interface{}) {
		if1 := if0.(map[string]interface{})
		interfacesOld = append(interfacesOld, if1)
	}

	var interfacesNew []map[string]interface{}

	for _, if0 := range ifNew.([]interface{}) {
		if1 := if0.(map[string]interface{})
		interfacesNew = append(interfacesNew, if1)
	}

	sort.Slice(interfacesOld, func(i, j int) bool {
		return interfacesOld[i]["ifname"].(string) < interfacesOld[j]["ifname"].(string)
	})

	sort.Slice(interfacesNew, func(i, j int) bool {
		return interfacesNew[i]["ifname"].(string) < interfacesNew[j]["ifname"].(string)
	})

	return reflect.DeepEqual(interfacesOld, interfacesNew)
}

func DiffSuppressFuncVlan(k, old, new string, d *schema.ResourceData) bool {
	vOld, vNew := d.GetChange("vlan")
	var vlanOld []map[string]interface{}

	for _, v0 := range vOld.([]interface{}) {
		v1 := v0.(map[string]interface{})
		vlanOld = append(vlanOld, v1)
	}

	var vlanNew []map[string]interface{}

	for _, v0 := range vNew.([]interface{}) {
		v1 := v0.(map[string]interface{})
		vlanNew = append(vlanNew, v1)
	}

	sort.Slice(vlanOld, func(i, j int) bool {
		return vlanOld[i]["parent_interface"].(string) < vlanOld[j]["parent_interface"].(string)
	})

	sort.Slice(vlanNew, func(i, j int) bool {
		return vlanNew[i]["parent_interface"].(string) < vlanNew[j]["parent_interface"].(string)
	})

	return reflect.DeepEqual(vlanOld, vlanNew)
}
