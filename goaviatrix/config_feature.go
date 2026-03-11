package goaviatrix

import (
	"context"
	"fmt"
	"slices"
	"strings"
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

// GetAllFeatureNames returns the list of valid feature names from the controller API.
func (c *Client) GetAllFeatureNames(ctx context.Context) ([]string, error) {
	action := "get_all_controller_features"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}

	type ControllerSingleFeatureStatus struct {
		Feature string `json:"feature"`
		Enabled bool   `json:"enabled"`
	}

	type FeatureStatusResp struct {
		Results []ControllerSingleFeatureStatus `json:"results"`
	}

	var resp FeatureStatusResp
	err := c.PostAPIContext2(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	featureNames := make([]string, 0, len(resp.Results))
	for _, feature := range resp.Results {
		featureNames = append(featureNames, feature.Feature)
	}
	return featureNames, nil
}

func ValidateFeatureName(ctx context.Context, c *Client, featureName string) error {
	featureNames, err := c.GetAllFeatureNames(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(featureNames, featureName) {
		return fmt.Errorf("invalid feature name: %s. Valid feature names are: %s", featureName, strings.Join(featureNames, ", "))
	}
	return nil
}
