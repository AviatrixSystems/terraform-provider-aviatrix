package goaviatrix

import (
	"errors"
	"fmt"
	"strconv"
)

type GetBgpMaxAsLimitResult struct {
	MaxAsLimit string `json:"bgp_max_hop"`
}

type SetBgpMaxAsLimitResponse struct {
	Return    bool                   `json:"return"`
	Reason    string                 `json:"reason"`
	Results   GetBgpMaxAsLimitResult `json:"results"`
	ErrorCode int                    `json:"errorcode"`
}

func (c *Client) SetControllerBgpMaxAsLimit(maxAsLimit string) error {
	if maxAsLimit != "" {
		maxAsLimitInt, err := strconv.Atoi(maxAsLimit)
		if err != nil {
			return fmt.Errorf("error converting max_as_limit to int: %v", err)
		} else if maxAsLimitInt < 1 || maxAsLimitInt > 254 {
			return fmt.Errorf("max_as_limit must be between 1 and 254")
		}
	}

	data := map[string]string{
		"action":       "set_bgp_max_as_limit",
		"CID":          c.CID,
		"max_as_limit": maxAsLimit,
	}

	checkFunc := func(action, reason string, ret bool) error {
		if !ret {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) GetControllerBgpMaxAsLimit() (string, error) {
	data := map[string]string{
		"action": "show_bgp_max_as_limit",
		"CID":    c.CID,
	}

	var resp SetBgpMaxAsLimitResponse
	err := c.GetAPI(&resp, data["action"], data, BasicCheck)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("rest API set_bgp_max_as_limit failed with error code %d: %s", resp.ErrorCode, resp.Reason)
	} else if !resp.Return {
		return "", errors.New("Rest API set_bgp_max_as_limit Post failed: " + resp.Reason)
	}

	return resp.Results.MaxAsLimit, nil
}
