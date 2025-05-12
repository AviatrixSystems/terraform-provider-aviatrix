package goaviatrix

import (
	"context"
	"fmt"
	"strconv"
)

type DistributedFirewallingDefaultActionPolicy struct {
	Action  string `json:"action"`
	Logging bool   `json:"logging"`
}

func (c *Client) UpdateDistributedFirewallingDefaultActionPolicy(ctx context.Context, request *DistributedFirewallingDefaultActionPolicy) error {
	endpoint := "microseg/default-action-policy"
	return c.PutAPIContext25(ctx, endpoint, request)
}

func (c *Client) GetDistributedFirewallingDefaultActionPolicy(ctx context.Context) (map[string]string, error) {

	endpoint := "microseg/default-action-policy"

	var response DistributedFirewallingDefaultActionPolicy

	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get distributed firewalling default action policy: %w", err)
	}

	result := map[string]string{
		"action":  response.Action,
		"logging": strconv.FormatBool(response.Logging),
	}

	return result, nil
}
