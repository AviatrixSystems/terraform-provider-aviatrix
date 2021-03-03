package goaviatrix

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (c *Client) SetControllerBgpMaxAsLimit(ctx context.Context, maxAsLimit int) error {
	data := map[string]string{
		"action":       "set_bgp_max_as_limit",
		"CID":          c.CID,
		"max_as_limit": fmt.Sprint(maxAsLimit),
	}

	checkFunc := func(action, reason string, ret bool) error {
		if !ret && !strings.Contains(reason, "Configured BGP maximum AS limit is not changed") {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPIContext(ctx, data["action"], data, checkFunc)
}

func (c *Client) DisableControllerBgpMaxAsLimit(ctx context.Context) error {
	data := map[string]string{
		"action":       "set_bgp_max_as_limit",
		"CID":          c.CID,
		"max_as_limit": "",
	}

	checkFunc := func(action, reason string, ret bool) error {
		if !ret && !strings.Contains(reason, "Configured BGP maximum AS limit is not changed") {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPIContext(ctx, data["action"], data, checkFunc)
}

func (c *Client) GetControllerBgpMaxAsLimit(ctx context.Context) (int, error) {
	data := map[string]string{
		"action": "show_bgp_max_as_limit",
		"CID":    c.CID,
	}

	type BgpMaxAsLimitResponse struct {
		Results string
	}

	var resp BgpMaxAsLimitResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return 0, err
	}
	if resp.Results == "" {
		return 0, ErrNotFound
	}

	maxAsLimit, err := strconv.Atoi(resp.Results)
	if err != nil {
		return 0, fmt.Errorf("error converting max_as_limit to int: %v", err)
	}
	if maxAsLimit < 1 || maxAsLimit > 254 {
		return 0, fmt.Errorf("rest API show_bgp_max_as_limit returned invalid value for max_as_limit: %d. It must be an integer in the range of [1-254]", maxAsLimit)
	}
	return maxAsLimit, nil
}
