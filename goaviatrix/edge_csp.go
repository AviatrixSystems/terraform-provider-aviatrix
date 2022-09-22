package goaviatrix

import (
	"context"
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
	ManagementInterfaceConfig          string
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip,omitempty"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network,omitempty"`
	WanInterfaceIpPrefix               string `json:"wan_ip,omitempty"`
	WanDefaultGatewayIp                string `json:"wan_default_gateway,omitempty"`
	LanInterfaceIpPrefix               string `json:"lan_ip,omitempty"`
	ManagementInterfaceIpPrefix        string `json:"mgmt_ip,omitempty"`
	ManagementDefaultGatewayIp         string `json:"mgmt_default_gateway,omitempty"`
	DnsServerIp                        string `json:"dns_server_ip,omitempty"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary,omitempty"`
	Dhcp                               bool   `json:"dhcp,omitempty"`
	//ActiveStandby                      string `json:"active_standby,omitempty"`
	EnableEdgeActiveStandby           bool   `json:"enable_active_standby"`
	EnableEdgeActiveStandbyPreemptive bool   `json:"enable_active_standby_preemptive"`
	LocalAsNumber                     string `json:"local_as_number"`
	PrependAsPath                     []string
	PrependAsPathReturn               string   `json:"prepend_as_path"`
	IncludeCidrList                   []string `json:"include_cidr_list"`
	EnableLearnedCidrsApproval        bool     `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs              []string `json:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs     []string `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath              bool     `json:"preserve_as_path"`
	BgpPollingTime                    int      `json:"bgp_polling_time"`
	BgpHoldTime                       int      `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting       bool     `json:"edge_transitive_routing"`
	EnableJumboFrame                  bool     `json:"jumbo_frame"`
	Latitude                          string
	Longitude                         string
	LatitudeReturn                    float64 `json:"latitude"`
	LongitudeReturn                   float64 `json:"longitude"`
	WanPublicIp                       string  `json:"wan_discovery_ip"`
	PrivateIP                         string  `json:"private_ip"`
	RxQueueSize                       string  `json:"rx_queue_size"`
	State                             string  `json:"vpc_state"`
	NoProgressBar                     bool    `json:"no_progress_bar,omitempty"`
}

//type EdgeSpokeZededa struct {
//	Action                             string `form:"action,omitempty"`
//	CID                                string `form:"CID,omitempty"`
//	AccountName                        string `form:"account_name,omitempty"`
//	GwName                             string `form:"name,omitempty" json:"gw_name"`
//	SiteId                             string `form:"site_id,omitempty" json:"vpc_id"`
//	ProjectUuid                        string `form:"project_uuid,omitempty"`
//	ComputeNodeUuid                    string `form:"compute_node_uuid,omitempty"`
//	TemplateUuid                       string `form:"template_uuid,omitempty"`
//	ManagementInterfaceConfig          string
//	ManagementEgressIpPrefix           string `form:"mgmt_egress_ip,omitempty" json:"mgmt_egress_ip"`
//	EnableManagementOverPrivateNetwork bool   `form:"mgmt_over_private_network,omitempty" json:"mgmt_over_private_network"`
//	WanInterfaceIpPrefix               string `form:"wan_ip,omitempty" json:"wan_ip"`
//	WanDefaultGatewayIp                string `form:"wan_default_gateway,omitempty" json:"wan_default_gateway"`
//	LanInterfaceIpPrefix               string `form:"lan_ip,omitempty" json:"lan_ip"`
//	ManagementInterfaceIpPrefix        string `form:"mgmt_ip,omitempty" json:"mgmt_ip"`
//	ManagementDefaultGatewayIp         string `form:"mgmt_default_gateway,omitempty" json:"mgmt_default_gateway"`
//	DnsServerIp                        string `form:"dns_server_ip,omitempty" json:"dns_server_ip"`
//	SecondaryDnsServerIp               string `form:"dns_server_ip_secondary,omitempty" json:"dns_server_ip_secondary"`
//	Dhcp                               bool   `form:"dhcp,omitempty" json:"dhcp"`
//	ActiveStandby                      string `form:"active_standby,omitempty"`
//	EnableEdgeActiveStandby            bool   `json:"edge_active_standby"`
//	EnableEdgeActiveStandbyPreemptive  bool   `json:"edge_active_standby_preemptive"`
//	LocalAsNumber                      string `json:"local_as_number"`
//	PrependAsPath                      []string
//	PrependAsPathReturn                string   `json:"prepend_as_path"`
//	IncludeCidrList                    []string `json:"include_cidr_list"`
//	EnableLearnedCidrsApproval         bool     `json:"enable_learned_cidrs_approval"`
//	ApprovedLearnedCidrs               []string `form:"approved_learned_cidrs,omitempty"`
//	SpokeBgpManualAdvertisedCidrs      []string `json:"bgp_manual_spoke_advertise_cidrs"`
//	EnablePreserveAsPath               bool     `json:"preserve_as_path"`
//	BgpPollingTime                     int      `json:"bgp_polling_time"`
//	BgpHoldTime                        int      `json:"bgp_hold_time"`
//	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing"`
//	EnableJumboFrame                   bool     `json:"jumbo_frame"`
//	Latitude                           string
//	Longitude                          string
//	LatitudeReturn                     float64 `json:"latitude"`
//	LongitudeReturn                    float64 `json:"longitude"`
//	WanPublicIp                        string  `form:"wan_discovery_ip,omitempty" json:"public_ip"`
//	PrivateIP                          string  `json:"private_ip"`
//	RxQueueSize                        string  `json:"rx_queue_size"`
//	State                              string  `json:"vpc_state"`
//}

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
	WanInterfaceIpPrefix               string `json:"wan_ip"`
	WanDefaultGatewayIp                string `json:"wan_default_gateway"`
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
	LatitudeReturn                     float64 `json:"latitude"`
	LongitudeReturn                    float64 `json:"longitude"`
	WanPublicIp                        string  `json:"public_ip"`
	PrivateIP                          string  `json:"private_ip"`
	RxQueueSize                        string  `json:"rx_queue_size"`
	State                              string  `json:"vpc_state"`
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
