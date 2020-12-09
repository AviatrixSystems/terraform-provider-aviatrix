package goaviatrix

import "strconv"

type NetflowAgent struct {
	CID                   string
	ServerIp              string
	Port                  int
	Version               int
	ExcludedGatewaysInput string
}

type NetflowAgentResp struct {
	ServerIp         string   `json:"server_ip"`
	Port             string   `json:"port"`
	Version          string   `json:"version"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
}

func (c *Client) EnableNetflowAgent(r *NetflowAgent) error {
	params := map[string]string{
		"action":               "enable_netflow_agent",
		"CID":                  c.CID,
		"server_ip":            r.ServerIp,
		"port":                 strconv.Itoa(r.Port),
		"version":              strconv.Itoa(r.Version),
		"exclude_gateway_list": r.ExcludedGatewaysInput,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) GetNetflowAgentStatus() (*NetflowAgentResp, error) {
	params := map[string]string{
		"action": "get_netflow_agent",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool             `json:"return"`
		Results NetflowAgentResp `json:"results"`
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

func (c *Client) DisableNetflowAgent() error {
	params := map[string]string{
		"action": "disable_netflow_agent",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
