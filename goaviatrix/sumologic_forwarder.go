package goaviatrix

type SumologicForwarder struct {
	CID                   string
	AccessID              string
	AccessKey             string
	SourceCategory        string
	CustomCfg             string
	ExcludedGatewaysInput string
}

type SumologicForwarderResp struct {
	AccessID         string   `json:"acc_id"`
	AccessKey        string   `json:"acc_key"`
	SourceCategory   string   `json:"source_category"`
	CustomConfig     string   `json:"custom_cfg"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
}

func (c *Client) EnableSumologicForwarder(r *SumologicForwarder) error {
	params := map[string]string{
		"action":               "enable_sumologic_logging",
		"CID":                  c.CID,
		"access_id":            r.AccessID,
		"access_key":           r.AccessKey,
		"source_category":      r.SourceCategory,
		"custom_cfg":           r.CustomCfg,
		"exclude_gateway_list": r.ExcludedGatewaysInput,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) GetSumologicForwarderStatus() (*SumologicForwarderResp, error) {
	params := map[string]string{
		"action": "get_sumologic_logging_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool                   `json:"return"`
		Results SumologicForwarderResp `json:"results"`
		Reason  string                 `json:"reason"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	if data.Results.Status == "disabled" {
		return nil, ErrNotFound
	}

	return &data.Results, nil
}

func (c *Client) DisableSumologicForwarder() error {
	params := map[string]string{
		"action": "disable_sumologic_logging",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
