package goaviatrix

import (
	"context"
)

type CopilotSecurityGroupManagementConfig struct {
	Action           string `json:"action,omitempty"`
	CID              string `json:"CID,omitempty"`
	CloudType        int    `json:"cloud_type,omitempty"`
	AccountName      string `json:"account_name,omitempty"`
	Region           string `json:"region,omitempty"`
	Zone             string `json:"zone,omitempty"`
	VpcId            string `json:"vpc_id,omitempty"`
	InstanceID       string `json:"instance_id,omitempty"`
	InstanceIDReturn string `json:"inst_id,omitempty"`
	LogEnable        bool   `json:"log_enable,omitempty"`
	State            string `json:"state,omitempty"`
}

func (c *Client) EnableCopilotSecurityGroupManagement(ctx context.Context, copilotSecurityGroupManagementConfig *CopilotSecurityGroupManagementConfig) error {
	copilotSecurityGroupManagementConfig.Action = "enable_copilot_sg"
	copilotSecurityGroupManagementConfig.CID = c.CID
	copilotSecurityGroupManagementConfig.LogEnable = true

	err := c.PostAPIContext2(ctx, nil, copilotSecurityGroupManagementConfig.Action, copilotSecurityGroupManagementConfig, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCopilotSecurityGroupManagementConfig(ctx context.Context) (*CopilotSecurityGroupManagementConfig, error) {
	form := map[string]string{
		"action": "get_copilot_sg",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool                                 `json:"return"`
		Results CopilotSecurityGroupManagementConfig `json:"results"`
		Reason  string                               `json:"reason"`
	}

	var data Resp

	err := c.GetAPIContext(ctx, &data, form["action"], form, BasicCheck)

	if err != nil {
		return nil, err
	}

	if data.Results.State == "Disabled" {
		return nil, ErrNotFound
	}

	return &data.Results, nil
}

func (c *Client) DisableCopilotSecurityGroupManagement(ctx context.Context) error {
	form := map[string]string{
		"action":     "disable_copilot_sg",
		"CID":        c.CID,
		"log_enable": "true",
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}
