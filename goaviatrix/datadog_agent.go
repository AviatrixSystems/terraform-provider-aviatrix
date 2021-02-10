package goaviatrix

import "strconv"

type DatadogAgent struct {
	CID                   string
	ApiKey                string
	Site                  string
	ExcludedGatewaysInput string
	MetricsOnly           bool
}

type DatadogAgentResp struct {
	ApiKey           string   `json:"api_key"`
	Site             string   `json:"site"`
	ExcludedGateways []string `json:"excluded_gateway"`
	MetricsOnly      bool     `json:"metrics_only"`
	Status           string   `json:"status"`
}

func (c *Client) EnableDatadogAgent(r *DatadogAgent) error {
	params := map[string]string{
		"action":               "enable_datadog_agent_logging",
		"CID":                  c.CID,
		"api_key":              r.ApiKey,
		"site":                 r.Site,
		"exclude_gateway_list": r.ExcludedGatewaysInput,
		"metrics_only":         strconv.FormatBool(r.MetricsOnly),
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) GetDatadogAgentStatus() (*DatadogAgentResp, error) {
	params := map[string]string{
		"action": "get_datadog_agent_logging_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool             `json:"return"`
		Results DatadogAgentResp `json:"results"`
		Reason  string           `json:"reason"`
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

func (c *Client) DisableDatadogAgent() error {
	params := map[string]string{
		"action": "disable_datadog_agent_logging",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
