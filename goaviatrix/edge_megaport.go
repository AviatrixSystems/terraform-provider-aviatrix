package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"os"
	"strings"
)

type EdgeMegaport struct {
	Action                             string                   `json:"action,omitempty"`
	CID                                string                   `json:"CID,omitempty"`
	AccountName                        string                   `json:"account_name,omitempty"`
	GwName                             string                   `json:"name,omitempty"`
	SiteID                             string                   `json:"site_id,omitempty"`
	ZtpFileDownloadPath                string                   `json:"-"`
	ManagementEgressIPPrefix           string                   `json:"mgmt_egress_ip,omitempty"`
	EnableManagementOverPrivateNetwork bool                     `json:"mgmt_over_private_network,omitempty"`
	DNSServerIP                        string                   `json:"dns_server_ip,omitempty"`
	SecondaryDNSServerIP               string                   `json:"dns_server_ip_secondary,omitempty"`
	Dhcp                               bool                     `json:"dhcp,omitempty"`
	EnableEdgeActiveStandby            bool                     `json:"enable_active_standby,omitempty"`
	DisableEdgeActiveStandby           bool                     `json:"disable_active_standby,omitempty"`
	EnableEdgeActiveStandbyPreemptive  bool                     `json:"enable_active_standby_preemptive,omitempty"`
	DisableEdgeActiveStandbyPreemptive bool                     `json:"disable_active_standby_preemptive,omitempty"`
	LocalAsNumber                      string                   `json:"local_as_number,omitempty"`
	PrependAsPath                      []string                 `json:"-"`
	PrependAsPathReturn                string                   `json:"prepend_as_path,omitempty"`
	IncludeCidrList                    []string                 `json:"include_cidr_list,omitempty"`
	EnableLearnedCidrsApproval         bool                     `json:"enable_learned_cidrs_approval,omitempty"`
	ApprovedLearnedCidrs               []string                 `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string                 `json:"bgp_manual_spoke_advertise_cidrs,omitempty"`
	EnablePreserveAsPath               bool                     `json:"preserve_as_path,omitempty"`
	BgpPollingTime                     int                      `json:"bgp_polling_time,omitempty"`
	BgpBfdPollingTime                  int                      `json:"bgp_neighbor_status_polling_time,omitempty"`
	BgpHoldTime                        int                      `json:"bgp_hold_time,omitempty"`
	EnableEdgeTransitiveRouting        bool                     `json:"edge_transitive_routing,omitempty"`
	EnableJumboFrame                   bool                     `json:"jumbo_frame,omitempty"`
	Latitude                           string                   `json:"-"`
	Longitude                          string                   `json:"-"`
	RxQueueSize                        string                   `json:"rx_queue_size,omitempty"`
	State                              string                   `json:"vpc_state,omitempty"`
	NoProgressBar                      bool                     `json:"no_progress_bar,omitempty"`
	InterfaceList                      []*EdgeMegaportInterface `json:"-"`
	Interfaces                         string                   `json:"interfaces,omitempty"`
	VlanList                           []*EdgeMegaportVlan      `json:"-"`
	Vlan                               string                   `json:"vlan,omitempty"`
	EnableSingleIpSnat                 bool                     `json:"-"`
	EnableAutoAdvertiseLanCidrs        string                   `json:"auto_advertise_lan_cidrs,omitempty"`
	LanInterfaceIpPrefix               string                   `json:"-"`
	InterfaceMapping                   string                   `json:"interface_mapping,omitempty"`
	CloudType                          int                      `json:"cloud_type,omitempty"`
}

type EdgeMegaportInterface struct {
	LogicalInterfaceName string              `json:"logical_ifname"`
	PublicIP             string              `json:"public_ip,omitempty"`
	Tag                  string              `json:"tag,omitempty"`
	Dhcp                 bool                `json:"dhcp,omitempty"`
	IPAddr               string              `json:"ipaddr,omitempty"`
	GatewayIP            string              `json:"gateway_ip,omitempty"`
	DNSPrimary           string              `json:"dns_primary,omitempty"`
	DNSSecondary         string              `json:"dns_secondary,omitempty"`
	SubInterfaces        []*EdgeMegaportVlan `json:"subinterfaces,omitempty"`
	VrrpState            bool                `json:"vrrp_state,omitempty"`
	VirtualIP            string              `json:"virtual_ip,omitempty"`
}

type EdgeMegaportVlan struct {
	ParentLogicalInterface string `json:"parent_logical_interface"`
	VlanID                 string `json:"vlan_id"`
	IPAddr                 string `json:"ipaddr"`
	GatewayIP              string `json:"gateway_ip,omitempty"`
	PeerIPAddr             string `json:"peer_ipaddr,omitempty"`
	PeerGatewayIP          string `json:"peer_gateway_ip,omitempty"`
	VirtualIP              string `json:"virtual_ip,omitempty"`
	Tag                    string `json:"tag,omitempty"`
}

type EdgeMegaportResp struct {
	AccountName                        string              `json:"account_name"`
	GwName                             string              `json:"gw_name"`
	SiteId                             string              `json:"vpc_id"`
	ManagementInterfaceConfig          string              `json:"-"`
	ManagementEgressIPPrefix           string              `json:"mgmt_egress_ip"`
	EnableManagementOverPrivateNetwork bool                `json:"mgmt_over_private_network"`
	LanInterfaceIpPrefix               string              `json:"lan_ip"`
	ManagementInterfaceIPPrefix        string              `json:"mgmt_ip"`
	ManagementDefaultGatewayIP         string              `json:"mgmt_default_gateway"`
	DNSServerIP                        string              `json:"dns_server_ip"`
	SecondaryDNSServerIP               string              `json:"dns_server_ip_secondary"`
	Dhcp                               bool                `json:"dhcp"`
	ActiveStandby                      string              `json:"active_standby"`
	EnableEdgeActiveStandby            bool                `json:"edge_active_standby"`
	EnableEdgeActiveStandbyPreemptive  bool                `json:"edge_active_standby_preemptive"`
	LocalAsNumber                      string              `json:"local_as_number"`
	PrependAsPath                      []string            `json:"-"`
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
	WanPublicIP                        string              `json:"public_ip"`
	PrivateIP                          string              `json:"private_ip"`
	RxQueueSize                        string              `json:"rx_queue_size"`
	State                              string              `json:"vpc_state"`
	InterfaceList                      []MegaportInterface `json:"interfaces"`
	DNSProfileName                     string              `json:"dns_profile_name"`
	EnableNat                          string              `json:"enable_nat"`
	SnatMode                           string              `json:"snat_target"`
	EnableAutoAdvertiseLanCidrs        bool                `json:"auto_advertise_lan_cidrs"`
	InterfaceMapping                   []InterfaceMapping  `json:"interface_mapping"`
	AdvertisedCidrList                 []string            `json:"advertise_cidr_list,omitempty"`
}

type EdgeMegaportListResp struct {
	Return  bool               `json:"return"`
	Results []EdgeMegaportResp `json:"results"`
	Reason  string             `json:"reason"`
}

type MegaportInterface struct {
	LogicalInterfaceName string              `json:"logical_ifname"`
	Name                 string              `json:"ifname,omitempty"`
	PublicIP             string              `json:"public_ip,omitempty"`
	Tag                  string              `json:"tag,omitempty"`
	Dhcp                 bool                `json:"dhcp,omitempty"`
	IPAddr               string              `json:"ipaddr,omitempty"`
	GatewayIP            string              `json:"gateway_ip,omitempty"`
	DNSPrimary           string              `json:"dns_primary,omitempty"`
	DNSSecondary         string              `json:"dns_secondary,omitempty"`
	SubInterfaces        []*EdgeMegaportVlan `json:"subinterfaces,omitempty"`
	VrrpState            bool                `json:"vrrp_state,omitempty"`
	VirtualIP            string              `json:"virtual_ip,omitempty"`
}

type CreateEdgeMegaportResp struct {
	Return bool   `json:"return"`
	Result string `json:"results"`
	Reason string `json:"reason"`
}

func (c *Client) CreateEdgeMegaport(ctx context.Context, edgeMegaport *EdgeMegaport) error {
	edgeMegaport.Action = "create_megaport_instance"
	edgeMegaport.CID = c.CID
	edgeMegaport.NoProgressBar = true

	interfaces, err := json.Marshal(edgeMegaport.InterfaceList)
	if err != nil {
		return err
	}

	edgeMegaport.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if len(edgeMegaport.VlanList) != 0 {
		vlan, err := json.Marshal(edgeMegaport.VlanList)
		if err != nil {
			return err
		}
		edgeMegaport.Vlan = b64.StdEncoding.EncodeToString(vlan)
	}

	var data CreateEdgeMegaportResp

	err = c.PostAPIContext2(ctx, &data, edgeMegaport.Action, edgeMegaport, BasicCheck)
	if err != nil {
		return err
	}

	fileName := edgeMegaport.ZtpFileDownloadPath + "/" + edgeMegaport.GwName + "-" + edgeMegaport.SiteID + "-cloud-init.txt"

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = outFile.WriteString(data.Result)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeMegaport(ctx context.Context, gwName string) (*EdgeMegaportResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeMegaportListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeMegaportList := data.Results
	for _, edgeMegaport := range edgeMegaportList {
		if edgeMegaport.GwName == gwName {
			for _, p := range strings.Split(edgeMegaport.PrependAsPathReturn, " ") {
				if p != "" {
					edgeMegaport.PrependAsPath = append(edgeMegaport.PrependAsPath, p)
				}
			}

			return &edgeMegaport, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteEdgeMegaport(ctx context.Context, accountName, name string) error {
	form := map[string]string{
		"action":       "delete_megaport_instance",
		"CID":          c.CID,
		"account_name": accountName,
		"name":         name,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) UpdateEdgeMegaport(ctx context.Context, edgeMegaport *EdgeMegaport) error {
	edgeMegaport.Action = "update_edge_gateway"
	edgeMegaport.CID = c.CID
	edgeMegaport.CloudType = EDGEMEGAPORT

	interfaces, err := json.Marshal(edgeMegaport.InterfaceList)
	if err != nil {
		return err
	}

	edgeMegaport.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if len(edgeMegaport.VlanList) != 0 {
		vlan, err := json.Marshal(edgeMegaport.VlanList)
		if err != nil {
			return err
		}
		edgeMegaport.Vlan = b64.StdEncoding.EncodeToString(vlan)
	}

	return c.PostAPIContext2(ctx, nil, edgeMegaport.Action, edgeMegaport, BasicCheck)
}
