package goaviatrix

import (
	"context"
)

type WanInterface struct {
	Interface string `json:"interface"`
	Ip        string `json:"ip"`
}

type WanInterfaceResp struct {
	Return  bool           `json:"return"`
	Results []WanInterface `json:"results"`
	Reason  string         `json:"reason"`
}

func (c *Client) GetEdgeGatewayWanIp(ctx context.Context, gwName, wanInterfaceName string) (string, error) {
	form := map[string]string{
		"action":        "discovery_edge_gateway_wan_ip",
		"CID":           c.CID,
		"name":          gwName,
		"wan_interface": wanInterfaceName,
	}

	var data WanInterfaceResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	return data.Results[0].Ip, nil
}
