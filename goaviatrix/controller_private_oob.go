package goaviatrix

import (
	"fmt"
	"strings"
)

type PrivateOobResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) EnablePrivateOob() error {
	data := map[string]string{
		"action": "enable_private_oob",
		"CID":    c.CID,
	}
	checkFunc := func(action, reason string, ret bool) error {
		if !ret && !strings.HasPrefix(reason, "enable already") {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) GetPrivateOobState() (bool, error) {
	var data PrivateOobResp
	form := map[string]string{
		"CID":    c.CID,
		"action": "get_private_oob_state",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
	}
	if data.Results == "Enabled" || data.Results == "enabled" {
		return true, nil
	} else if data.Results == "Disabled" || data.Results == "disabled" {
		return false, nil
	}
	return false, ErrNotFound
}

func (c *Client) DisablePrivateOob() error {
	data := map[string]string{
		"action": "disable_private_oob",
		"CID":    c.CID,
	}
	checkFunc := func(action, reason string, ret bool) error {
		if !ret && !strings.HasPrefix(reason, "disable already") {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}
