package goaviatrix

import "context"

type SymmetricRoutingResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) EnableSymmetricRouting(ctx context.Context, gwName string) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_biflow",
		"gateway_name": gwName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) DisableSymmetricRouting(ctx context.Context, gwName string) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_biflow",
		"gateway_name": gwName,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) GetSymmetricRoutingStatus(ctx context.Context, gwName string) (string, error) {
	form := map[string]string{
		"action":       "get_biflow_status",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data SymmetricRoutingResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	return data.Results, nil
}
