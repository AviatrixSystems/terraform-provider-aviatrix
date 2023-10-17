package goaviatrix

import "context"

type EdgeResp struct {
	GwName        string       `json:"gw_name"`
	InterfaceList []*Interface `json:"interfaces"`
}

type EdgeListResp struct {
	Return  bool       `json:"return"`
	Results []EdgeResp `json:"results"`
	Reason  string     `json:"reason"`
}

func (c *Client) GetEdgeGatewayWanInterfaces(ctx context.Context, gwName string) (*EdgeResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeList := data.Results
	for _, edge := range edgeList {
		if edge.GwName == gwName {
			return &edge, nil
		}
	}

	return nil, ErrNotFound
}
