package goaviatrix

import (
	"context"
	"strings"
)

func (c *Client) EnableControllerBgpCommunitiesGlobal(ctx context.Context) error {
	data := map[string]string{
		"action": "enable_bgp_communities_global",
		"CID":    c.CID,
	}
	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) DisableControllerBgpCommunitiesGlobal(ctx context.Context) error {
	data := map[string]string{
		"action": "disable_bgp_communities_global",
		"CID":    c.CID,
	}
	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) GetControllerBgpCommunitiesGlobal(ctx context.Context) (bool, error) {
	data := map[string]string{
		"action": "show_bgp_communities_global",
		"CID":    c.CID,
	}

	type BgpCommunitiesGlobalResults struct {
		BgpCommunitiesGlobal string `json:"result"`
	}

	type BgpCommunitiesGlobalResponse struct {
		Results BgpCommunitiesGlobalResults
	}

	var resp BgpCommunitiesGlobalResponse
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return false, err
	}
	if strings.Contains(resp.Results.BgpCommunitiesGlobal, "enabled") {
		return true, nil
	} else if strings.Contains(resp.Results.BgpCommunitiesGlobal, "disabled") {
		return false, nil
	}
	return false, ErrNotFound
}
