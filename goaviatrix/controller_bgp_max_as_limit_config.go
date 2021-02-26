package goaviatrix

import (
	"fmt"
	"strconv"
	"strings"
)

func (c *Client) SetControllerBgpMaxAsLimit(maxAsLimit int) error {
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
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) DisableControllerBgpMaxAsLimit() error {
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
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) GetControllerBgpMaxAsLimit() (int, error) {
	data := map[string]string{
		"action": "show_bgp_max_as_limit",
		"CID":    c.CID,
	}

	type BgpMaxAsLimitResponse struct {
		Return bool   `json:"return"`
		Reason string `json:"reason"`
		Result string `json:"results"`
	}

	var resp BgpMaxAsLimitResponse
	err := c.GetAPI(&resp, data["action"], data, BasicCheck)
	if err != nil {
		return -1, err
	}

	if !resp.Return {
		return -1, fmt.Errorf("rest API set_bgp_max_as_limit Post failed: %s", resp.Reason)
	} else if resp.Result == "" {
		return -1, nil
	}

	maxAsLimit, err := strconv.Atoi(resp.Result)
	if err != nil {
		return -1, fmt.Errorf("error converting max_as_limit to int: %v", err)
	} else if maxAsLimit < 1 || maxAsLimit > 254 {
		return -1, fmt.Errorf("rest API show_bgp_max_as_limit returned invalid value for max_as_limit: %d. It must be an integer in the range of [1-254]", maxAsLimit)
	}
	return maxAsLimit, nil
}
