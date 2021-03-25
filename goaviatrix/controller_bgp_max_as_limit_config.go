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

func (c *Client) SetControllerBgpMaxAsLimitNonRfc1918(ctx context.Context, maxAsLimitNonRfc1918 int) error {
	data := map[string]string{
		"action":       "set_bgp_max_as_limit_non_rfc1918",
		"CID":          c.CID,
		"max_as_limit": fmt.Sprint(maxAsLimitNonRfc1918),
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

func (c *Client) DisableControllerBgpMaxAsLimitNonRfc1918(ctx context.Context) error {
	data := map[string]string{
		"action":       "set_bgp_max_as_limit_non_rfc1918",
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

func (c *Client) GetControllerBgpMaxAsLimit(ctx context.Context) (int, int, error) {
	data := map[string]string{
		"action": "show_bgp_max_as_limit",
		"CID":    c.CID,
	}

	type BgpMaxAsLimitResults struct {
		MaxAsLimit           string `json:"bgp_max_as_limit_rfc1918"`
		MaxAsLimitNonRfc1918 string `json:"bgp_max_as_limit_non_rfc1918"`
	}

	type BgpMaxAsLimitResponse struct {
		Results BgpMaxAsLimitResults
	}

	var resp BgpMaxAsLimitResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return 0, 0, err
	}

	if resp.Results.MaxAsLimit == "" && resp.Results.MaxAsLimitNonRfc1918 == "" {
		return 0, 0, ErrNotFound
	}

	var maxAsLimit, maxAsLimitNonRfc1918 int
	if resp.Results.MaxAsLimit != "" {
		maxAsLimit, err = strconv.Atoi(resp.Results.MaxAsLimit)
		if err != nil {
			return 0, 0, fmt.Errorf("error converting max_as_limit to int: %v", err)
		}
	}

	if resp.Results.MaxAsLimitNonRfc1918 != "" {
		maxAsLimitNonRfc1918, err = strconv.Atoi(resp.Results.MaxAsLimitNonRfc1918)
		if err != nil {
			return 0, 0, fmt.Errorf("error converting max_as_limit_non_rfc1918 to int: %v", err)
		}
	}

	return maxAsLimit, maxAsLimitNonRfc1918, nil
}
