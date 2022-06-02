package goaviatrix

import (
	"context"
)

type CopilotSecurityGroupManagement struct {
	Action      string `form:"action,omitempty"`
	CID         string `form:"CID,omitempty"`
	CloudType   int    `form:"cloud_type,omitempty" json:"cloud_type"`
	AccountName string `form:"account_name,omitempty" json:"account_name"`
	Region      string `form:"region,omitempty" json:"region"`
	Zone        string `form:"zone,omitempty" json:"zone"`
	VpcId       string `form:"vpc_id,omitempty" json:"vpc_id"`
	InstanceID  string `form:"instance_id,omitempty" json:"inst_id"`
	LogEnable   bool   `form:"log_enable"`
	State       string `json:"state"`
}

func (c *Client) EnableCopilotSecurityGroupManagement(ctx context.Context, copilotSecurityGroupManagement *CopilotSecurityGroupManagement) error {
	copilotSecurityGroupManagement.Action = "enable_copilot_sg"
	copilotSecurityGroupManagement.CID = c.CID
	copilotSecurityGroupManagement.LogEnable = true

	err := c.PostAPIContext(ctx, copilotSecurityGroupManagement.Action, copilotSecurityGroupManagement, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCopilotSecurityGroupManagement(ctx context.Context) (*CopilotSecurityGroupManagement, error) {
	form := map[string]string{
		"action": "get_copilot_sg",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool                           `json:"return"`
		Results CopilotSecurityGroupManagement `json:"results"`
		Reason  string                         `json:"reason"`
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
