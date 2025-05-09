package goaviatrix

import (
	"fmt"
)

func (c *Client) SetGatewayBgpCommunitiesAccept(gwName string, sendComm bool) error {
	data := map[string]string{
		"action":             "set_gateway_accept_bgp_communities_override",
		"CID":                c.CID,
		"gateway_name":       gwName,
		"accept_communities": fmt.Sprint(sendComm),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) SetGatewayBgpCommunitiesSend(gwName string, sendComm bool) error {
	data := map[string]string{
		"action":           "set_gateway_send_bgp_communities_override",
		"CID":              c.CID,
		"gateway_name":     gwName,
		"send_communities": fmt.Sprint(sendComm),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) GetGatewayBgpCommunities(gwName string) (bool, bool, error) {
	data := map[string]string{
		"action":       "show_bgp_communities_gateway_overrides",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	type BgpCommunitiesGatewayResults struct {
		BgpCommunitiesGatewayAccept bool   `json:"accept_communities"`
		BgpCommunitiesGatewaySend   bool   `json:"send_communities"`
		BgpCommunitiesGatewayText   string `json:"text"`
	}

	type BgpCommunitiesGatewayResponse struct {
		Results BgpCommunitiesGatewayResults
	}

	var resp BgpCommunitiesGatewayResponse
	err := c.GetAPI(&resp, data["action"], data, BasicCheck)
	if err != nil {
		return false, false, err
	}
	return resp.Results.BgpCommunitiesGatewayAccept, resp.Results.BgpCommunitiesGatewaySend, nil
}
