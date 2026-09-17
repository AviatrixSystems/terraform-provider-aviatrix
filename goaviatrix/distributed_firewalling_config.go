package goaviatrix

import (
	"context"
	"fmt"
)

type DistributedFirewallingConfig struct {
	EnableDistributedFirewalling bool `json:"enable_distributed_firewalling"`
}

func (c *Client) EnableDistributedFirewalling(ctx context.Context) error {
	action := "enable_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": "microseg",
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) DisableDistributedFirewalling(ctx context.Context) error {
	action := "disable_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": "microseg",
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) GetDistributedFirewallingStatus(ctx context.Context) (*DistributedFirewallingConfig, error) {
	action := "get_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": "microseg",
	}
	dFWConfig := &DistributedFirewallingConfig{}

	type ControllerSingleFeatureStatus struct {
		Feature string `json:"feature"`
		Enabled bool   `json:"enabled"`
	}

	type DistributedFirewallingConfigResp struct {
		Results ControllerSingleFeatureStatus `json:"results"`
	}

	var resp DistributedFirewallingConfigResp
	err := c.PostAPIContext2(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	dFWConfig.EnableDistributedFirewalling = resp.Results.Enabled
	return dFWConfig, nil
}
