package goaviatrix

type SplunkLoggingResp struct {
	Server           string   `json:"server_ip"`
	Port             string   `json:"server_port"`
	CustomConfig     string   `json:"custom_input_cfg"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
}

func (c *Client) GetSplunkLoggingStatus() (*SplunkLoggingResp, error) {
	params := map[string]string{
		"action": "get_splunk_logging_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool              `json:"return,omitempty"`
		Results SplunkLoggingResp `json:"results,omitempty"`
		Reason  string            `json:"reason,omitempty"`
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

func (c *Client) DisableSplunkLogging() error {
	params := map[string]string{
		"action": "disable_splunk_logging",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
