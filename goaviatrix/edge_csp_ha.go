package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
)

type EdgeCSPHa struct {
	Action                    string `json:"action"`
	CID                       string `json:"CID"`
	PrimaryGwName             string `json:"primary_gw_name"`
	ComputeNodeUuid           string `json:"compute_node_uuid"`
	Dhcp                      bool   `json:"dhcp,omitempty"`
	ManagementInterfaceConfig string
	LanInterfaceIpPrefix      string `json:"lan_ip"`
	InterfaceList             []*Interface
	Interfaces                string `json:"interfaces"`
	NoProgressBar             bool   `json:"no_progress_bar,omitempty"`
	ManagementEgressIpPrefix  string `json:"mgmt_egress_ip,omitempty"`
}

type EdgeCSPHaResp struct {
	AccountName              string       `json:"account_name"`
	PrimaryGwName            string       `json:"primary_gw_name"`
	GwName                   string       `json:"gw_name"`
	Dhcp                     bool         `json:"dhcp"`
	ComputeNodeUuid          string       `json:"edge_csp_compute_node_uuid"`
	LanInterfaceIpPrefix     string       `json:"lan_ip"`
	InterfaceList            []*Interface `json:"interfaces"`
	ManagementEgressIpPrefix string       `json:"mgmt_egress_ip"`
}

type EdgeCSPHaListResp struct {
	Return  bool            `json:"return"`
	Results []EdgeCSPHaResp `json:"results"`
	Reason  string          `json:"reason"`
}

func (c *Client) CreateEdgeCSPHa(ctx context.Context, edgeCSPHa *EdgeCSPHa) (string, error) {
	edgeCSPHa.CID = c.CID
	edgeCSPHa.Action = "create_multicloud_ha_gateway"
	edgeCSPHa.NoProgressBar = true

	interfaces, err := json.Marshal(edgeCSPHa.InterfaceList)
	if err != nil {
		return "", err
	}

	edgeCSPHa.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2HaGw(ctx, nil, edgeCSPHa.Action, edgeCSPHa, BasicCheck)
}

func (c *Client) GetEdgeCSPHa(ctx context.Context, gwName string) (*EdgeCSPHaResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeCSPHaListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeCSPHaList := data.Results
	for _, edgeCSPHa := range edgeCSPHaList {
		if edgeCSPHa.GwName == gwName {
			return &edgeCSPHa, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeCSPHa(ctx context.Context, edgeCSP *EdgeCSP) error {
	form := map[string]string{
		"action":         "update_edge_gateway",
		"CID":            c.CID,
		"name":           edgeCSP.GwName,
		"mgmt_egress_ip": edgeCSP.ManagementEgressIpPrefix,
	}

	interfaces, err := json.Marshal(edgeCSP.InterfaceList)
	if err != nil {
		return err
	}

	form["interfaces"] = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
