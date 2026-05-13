package goaviatrix

type SumologicForwarderResp struct {
	AccessID         string   `json:"acc_id"`
	AccessKey        string   `json:"acc_key"`
	SourceCategory   string   `json:"source_category"`
	CustomConfig     string   `json:"custom_cfg"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
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
