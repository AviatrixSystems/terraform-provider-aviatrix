package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) SetGatewayKeepaliveConfig(ctx context.Context, speed string) error {
	data := map[string]string{
		"action": "set_keep_alive_speed",
		"CID":    c.CID,
		"speed":  speed,
	}

	checkFunc := func(action, reason string, ret bool) error {
		if ret || strings.Contains(reason, "Template is already set to") {
			return nil
		}
		return fmt.Errorf("rest API %s Post failed: %s", action, reason)
	}

	return c.PostAPIContext(ctx, data["action"], data, checkFunc)
}

func (c *Client) GetGatewayKeepaliveConfig(ctx context.Context) (string, error) {
	data := map[string]string{
		"action": "get_keep_alive_speed",
		"CID":    c.CID,
	}

	type GatewayKeepaliveResults struct {
		Speed string `json:"template"`
	}

	type GatewayKeepaliveResponse struct {
		Results GatewayKeepaliveResults
	}

	var resp GatewayKeepaliveResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return "", err
	}
	if resp.Results.Speed == "" {
		return "", ErrNotFound
	}

	return resp.Results.Speed, nil
}
