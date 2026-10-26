package goaviatrix

import (
	"context"
	"strings"
)

// This file contains the API calls related to controller BGP MED to SDN metric global config
// Details on this feature and related APIs can be found here:
// Design doc: https://docs.google.com/document/d/1h-gxgwZ6OxNuLNLFgKymldqdoRU3kY1zEshbZs13C0A/edit?usp=sharing
// APIs: https://aviatrix.atlassian.net/wiki/spaces/AVXENG/pages/3136946177/BGP+MED+to+SDN+Metric+APIs

func (c *Client) EnableControllerBgpMedToSdnMetricGlobal(ctx context.Context) error {
	data := map[string]string{
		"action": "enable_bgp_med_to_sdn_metric_global",
		"CID":    c.CID,
	}
	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) DisableControllerBgpMedToSdnMetricGlobal(ctx context.Context) error {
	data := map[string]string{
		"action": "disable_bgp_med_to_sdn_metric_global",
		"CID":    c.CID,
	}
	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) GetControllerBgpMedToSdnMetricGlobal(ctx context.Context) (bool, error) {
	data := map[string]string{
		"action": "show_bgp_med_to_sdn_metric_global",
		"CID":    c.CID,
	}

	type BgpMedToSdnMetricGlobalResults struct {
		BgpMedToSdnMetricGlobal string `json:"result"`
	}

	type BgpMedToSdnMetricGlobalResponse struct {
		Results BgpMedToSdnMetricGlobalResults
	}

	var resp BgpMedToSdnMetricGlobalResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return false, err
	}
	if strings.Contains(resp.Results.BgpMedToSdnMetricGlobal, "enabled") {
		return true, nil
	} else if strings.Contains(resp.Results.BgpMedToSdnMetricGlobal, "disabled") {
		return false, nil
	}
	return false, ErrNotFound
}
