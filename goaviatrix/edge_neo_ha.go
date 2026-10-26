package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
)

type EdgeNEOHa struct {
	Action                   string `json:"action"`
	CID                      string `json:"CID"`
	PrimaryGwName            string `json:"primary_gw_name"`
	DeviceId                 string `json:"device_id"`
	InterfaceList            []*EdgeNEOInterface
	Interfaces               string `json:"interfaces"`
	NoProgressBar            bool   `json:"no_progress_bar,omitempty"`
	ManagementEgressIpPrefix string `json:"mgmt_egress_ip,omitempty"`
	DirectAttachLan          bool   `json:"direct_attach_lan"`
}

type EdgeNEOHaResp struct {
	AccountName              string              `json:"account_name"`
	PrimaryGwName            string              `json:"primary_gw_name"`
	GwName                   string              `json:"gw_name"`
	DeviceId                 string              `json:"edge_csp_device_id"`
	InterfaceList            []*EdgeNEOInterface `json:"interfaces"`
	ManagementEgressIpPrefix string              `json:"mgmt_egress_ip"`
}

type EdgeNEOHaListResp struct {
	Return  bool            `json:"return"`
	Results []EdgeNEOHaResp `json:"results"`
	Reason  string          `json:"reason"`
}

func (c *Client) CreateEdgeNEOHa(ctx context.Context, edgeNEOHa *EdgeNEOHa) (string, error) {
	edgeNEOHa.CID = c.CID
	edgeNEOHa.Action = "create_multicloud_ha_gateway"
	edgeNEOHa.NoProgressBar = true
	edgeNEOHa.DirectAttachLan = true

	interfaces, err := json.Marshal(edgeNEOHa.InterfaceList)
	if err != nil {
		return "", err
	}

	edgeNEOHa.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2HaGw(ctx, nil, edgeNEOHa.Action, edgeNEOHa, BasicCheck)
}

func (c *Client) GetEdgeNEOHa(ctx context.Context, gwName string) (*EdgeNEOHaResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeNEOHaListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeNEOHaList := data.Results
	for _, edgeNEOHa := range edgeNEOHaList {
		if edgeNEOHa.GwName == gwName {
			return &edgeNEOHa, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeNEOHa(ctx context.Context, edgeNEO *EdgeNEO) error {
	form := map[string]string{
		"action":         "update_edge_gateway",
		"CID":            c.CID,
		"name":           edgeNEO.GwName,
		"mgmt_egress_ip": edgeNEO.ManagementEgressIpPrefix,
	}

	interfaces, err := json.Marshal(edgeNEO.InterfaceList)
	if err != nil {
		return err
	}

	form["interfaces"] = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
