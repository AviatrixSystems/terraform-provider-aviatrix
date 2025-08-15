package goaviatrix

import (
	"context"
	"strings"
)

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
