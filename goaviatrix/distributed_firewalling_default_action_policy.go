package goaviatrix

import (
	"context"
	"fmt"
)

const distributedFirewallingDefaultActionPolicyEndpoint = "microseg/default-action-policy"

type DistributedFirewallingDefaultActionPolicy struct {
	Action  string `json:"action"`
	Logging bool   `json:"logging"`
}

func (c *Client) UpdateDistributedFirewallingDefaultActionPolicy(ctx context.Context, request *DistributedFirewallingDefaultActionPolicy) error {
	return c.PutAPIContext25(ctx, distributedFirewallingDefaultActionPolicyEndpoint, request)
}

func (c *Client) GetDistributedFirewallingDefaultActionPolicy(ctx context.Context) (*DistributedFirewallingDefaultActionPolicy, error) {
	var response DistributedFirewallingDefaultActionPolicy

	err := c.GetAPIContext25(ctx, &response, distributedFirewallingDefaultActionPolicyEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get distributed firewalling default action policy: %w", err)
	}

	result := &DistributedFirewallingDefaultActionPolicy{
		Action:  response.Action,
		Logging: response.Logging,
	}

	return result, nil
}
