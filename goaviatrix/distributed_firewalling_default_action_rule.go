package goaviatrix

import (
	"context"
	"fmt"
)

const distributedFirewallingDefaultActionRuleEndpoint = "microseg/default-action-policy"

type DistributedFirewallingDefaultActionRule struct {
	Action     string `json:"action"`
	Logging    bool   `json:"logging"`
	LogProfile string `json:"log_profile"`
}

func (c *Client) UpdateDistributedFirewallingDefaultActionRule(ctx context.Context, request *DistributedFirewallingDefaultActionRule) error {
	return c.PutAPIContext25(ctx, distributedFirewallingDefaultActionRuleEndpoint, request)
}

func (c *Client) GetDistributedFirewallingDefaultActionRule(ctx context.Context) (*DistributedFirewallingDefaultActionRule, error) {
	var response DistributedFirewallingDefaultActionRule

	err := c.GetAPIContext25(ctx, &response, distributedFirewallingDefaultActionRuleEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get distributed firewalling default action rule: %w", err)
	}

	result := &DistributedFirewallingDefaultActionRule{
		Action:     response.Action,
		Logging:    response.Logging,
		LogProfile: response.LogProfile,
	}

	return result, nil
}
