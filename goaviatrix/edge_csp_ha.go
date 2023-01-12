package goaviatrix

import (
	"context"
)

type EdgeCSPHa struct {
	Action          string `json:"action"`
	CID             string `json:"CID"`
	PrimaryGwName   string `json:"primary_gw_name"`
	Dhcp            bool   `json:"dhcp"`
	ComputeNodeUuid string `json:"compute_node_uuid"`
	LanIp           string `json:"lan_ip"`
}

type EdgeCSPHaResp struct {
	AccountName     string `json:"account_name"`
	PrimaryGwName   string `json:"primary_gw_name"`
	GwName          string `json:"gw_name"`
	Dhcp            bool   `json:"dhcp"`
	ComputeNodeUuid string `json:"edge_csp_compute_node_uuid"`
	LanIp           string `json:"lan_ip"`
}

type EdgeCSPHaListResp struct {
	Return  bool            `json:"return"`
	Results []EdgeCSPHaResp `json:"results"`
	Reason  string          `json:"reason"`
}

func (c *Client) CreateEdgeCSPHa(ctx context.Context, edgeCSPHa *EdgeCSPHa) (string, error) {
	edgeCSPHa.CID = c.CID
	edgeCSPHa.Action = "create_multicloud_ha_gateway"

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

func (c *Client) DeleteEdgeCSPHa(ctx context.Context, name string) error {
	form := map[string]string{
		"action": "delete_edge_csp_instance",
		"CID":    c.CID,
		//"account_name": accountName,
		"name": name,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
