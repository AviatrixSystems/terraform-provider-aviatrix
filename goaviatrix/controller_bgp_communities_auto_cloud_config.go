package goaviatrix

import (
	"context"
	"fmt"
)

func (c *Client) SetControllerBgpCommunitiesAutoCloud(ctx context.Context, commPrefix int) error {
	data := map[string]string{
		"action":           "enable_auto_cloud_bgp_community",
		"CID":              c.CID,
		"community_prefix": fmt.Sprint(commPrefix),
	}

	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) DisableControllerBgpCommunitiesAutoCloud(ctx context.Context) error {
	data := map[string]string{
		"action":           "disable_auto_cloud_bgp_community",
		"CID":              c.CID,
		"community_prefix": "",
	}

	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) GetControllerBgpCommunitiesAutoCloud(ctx context.Context) (int, error) {
	data := map[string]string{
		"action": "show_auto_cloud_bgp_community",
		"CID":    c.CID,
	}

	type BgpCommunitiesAutoCloudResults struct {
		BgpCommunitiesAutoCloudEnabled bool `json:"enabled"`
		BgpCommunitiesAutoCloudPrefix  int  `json:"community_prefix"`
	}

	type BgpCommunitiesAutoCloudResponse struct {
		Results BgpCommunitiesAutoCloudResults
	}

	var resp BgpCommunitiesAutoCloudResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return 0, err
	}
	if resp.Results.BgpCommunitiesAutoCloudEnabled {
		return resp.Results.BgpCommunitiesAutoCloudPrefix, nil
	}
	return 0, nil
}
