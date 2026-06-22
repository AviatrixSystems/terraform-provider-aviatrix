package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

type CopilotSecurityGroupManagementConfig struct {
	Action                               string `json:"action,omitempty"`
	CID                                  string `json:"CID,omitempty"`
	CloudType                            int    `json:"cloud_type,omitempty"`
	AccountName                          string `json:"account_name,omitempty"`
	Region                               string `json:"region,omitempty"`
	Zone                                 string `json:"zone,omitempty"`
	VpcId                                string `json:"vpc_id,omitempty"`
	InstanceId                           string `json:"instance_id,omitempty"`
	InstanceIdReturn                     string `json:"inst_id,omitempty"`
	EnableCopilotSecurityGroupManagement bool
	LogEnable                            bool   `json:"log_enable,omitempty"`
	State                                string `json:"state,omitempty"`
}

func (c *Client) EnableCopilotSecurityGroupManagement(ctx context.Context, copilotSecurityGroupManagementConfig *CopilotSecurityGroupManagementConfig) error {
	copilotSecurityGroupManagementConfig.Action = "enable_copilot_sg"
	copilotSecurityGroupManagementConfig.CID = c.CID
	copilotSecurityGroupManagementConfig.LogEnable = true

	return c.PostAPIContext2(ctx, nil, copilotSecurityGroupManagementConfig.Action, copilotSecurityGroupManagementConfig, BasicCheck)
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

	return &data.Results, nil
}

func (c *Client) DisableCopilotSecurityGroupManagement(ctx context.Context) error {
	form := map[string]string{
		"action":     "disable_copilot_sg",
		"CID":        c.CID,
		"log_enable": "true",
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "not enabled") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext(ctx, form["action"], form, checkFunc)
}
