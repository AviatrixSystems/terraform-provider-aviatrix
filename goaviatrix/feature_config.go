package goaviatrix

import (
	"context"
	"fmt"
)

type FeatureStatus struct {
	Enabled bool `json:"enabled"`
}

func (c *Client) EnableFeature(ctx context.Context, featureName string) error {
	action := "enable_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": featureName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) DisableFeature(ctx context.Context, featureName string) error {
	action := "disable_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": featureName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) GetFeatureStatus(ctx context.Context, featureName string) (*FeatureStatus, error) {
	action := "get_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": featureName,
	}

	type ControllerSingleFeatureStatus struct {
		Feature string `json:"feature"`
		Enabled bool   `json:"enabled"`
	}

	type FeatureStatusResp struct {
		Results ControllerSingleFeatureStatus `json:"results"`
	}

	var resp FeatureStatusResp
	err := c.PostAPIContext2(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	featureStatus := &FeatureStatus{
		Enabled: resp.Results.Enabled,
	}
	return featureStatus, nil
}
