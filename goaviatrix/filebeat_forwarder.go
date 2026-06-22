package goaviatrix

type FilebeatForwarderResp struct {
	Server           string   `json:"server"`
	Port             string   `json:"port"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
}

func (c *Client) GetFilebeatForwarderStatus() (*FilebeatForwarderResp, error) {
	params := map[string]string{
		"action": "get_logstash_logging_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool                  `json:"return,omitempty"`
		Results FilebeatForwarderResp `json:"results,omitempty"`
		Reason  string                `json:"reason,omitempty"`
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

func (c *Client) DisableFilebeatForwarder() error {
	params := map[string]string{
		"action": "disable_logstash_logging",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
