package goaviatrix

type CloudwatchAgent struct {
	CID                   string
	RoleArn               string
	Region                string
	LogGroupName          string
	ExcludedGatewaysInput string
}

type CloudwatchAgentResp struct {
	RoleArn          string   `json:"cw_role_arn"`
	Region           string   `json:"region"`
	LogGroupName     string   `json:"log_group_name"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
}

func (c *Client) EnableCloudwatchAgent(r *CloudwatchAgent) error {
	params := map[string]string{
		"action":               "enable_cloudwatch_agent",
		"CID":                  c.CID,
		"cloudwatch_role_arn":  r.RoleArn,
		"region":               r.Region,
		"log_group_name":       r.LogGroupName,
		"exclude_gateway_list": r.ExcludedGatewaysInput,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) GetCloudwatchAgentStatus() (*CloudwatchAgentResp, error) {
	params := map[string]string{
		"action": "get_cloudwatch_agent_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool                `json:"return"`
		Results CloudwatchAgentResp `json:"results"`
		Reason  string              `json:"reason"`
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

func (c *Client) DisableCloudwatchAgent() error {
	params := map[string]string{
		"action": "disable_cloudwatch_agent",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
