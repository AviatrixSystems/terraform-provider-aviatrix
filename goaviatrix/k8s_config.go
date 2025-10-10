package goaviatrix

import (
	"context"
	"fmt"
)

type K8sConfig struct {
	EnableK8s         bool `json:"enable_k8s"`
	EnableDcfPolicies bool `json:"enable_dcf_policies"`
}

const (
	FeatureK8s            = "k8s"
	FeatureK8sDcfPolicies = "k8s_dcf_policies"

	actionEnableControllerFeature  = "enable_controller_feature"
	actionDisableControllerFeature = "disable_controller_feature"
)

func (c *Client) ToggleControllerFeature(ctx context.Context, feature string, enabled bool) error {
	action := actionEnableControllerFeature
	if !enabled {
		action = actionDisableControllerFeature
	}
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": feature,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPIContext2(ctx, nil, action, form, checkFunc)
}

func (c *Client) GetK8sStatus(ctx context.Context) (*K8sConfig, error) {
	k8sConfig := &K8sConfig{}

	// Get k8s feature status
	action := "get_controller_feature"
	form := map[string]string{
		"CID":     c.CID,
		"action":  action,
		"feature": FeatureK8s,
	}

	type ControllerSingleFeatureStatus struct {
		Feature string `json:"feature"`
		Enabled bool   `json:"enabled"`
	}

	type K8sConfigResp struct {
		Results ControllerSingleFeatureStatus `json:"results"`
	}

	var resp K8sConfigResp
	err := c.PostAPIContext2(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	k8sConfig.EnableK8s = resp.Results.Enabled

	// Get k8s_dcf_policies feature status
	form["feature"] = FeatureK8sDcfPolicies
	var dcfResp K8sConfigResp
	err = c.PostAPIContext2(ctx, &dcfResp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	k8sConfig.EnableDcfPolicies = dcfResp.Results.Enabled
	return k8sConfig, nil
}
