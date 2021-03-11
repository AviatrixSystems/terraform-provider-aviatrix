package goaviatrix

import (
	"context"
	"fmt"
)

func (c *Client) SetTunnelDetectionTime(ctx context.Context, detectionTime int, aviatrixEntity string) error {
	data := map[string]string{
		"action":         "modify_detection_time",
		"CID":            c.CID,
		"detection_time": fmt.Sprint(detectionTime),
		"entity":         aviatrixEntity,
	}

	return c.PostAPIContext(ctx, data["action"], data, BasicCheck)
}

func (c *Client) GetTunnelDetectionTime(ctx context.Context, aviatrixEntity string) (int, error) {
	data := map[string]string{
		"action": "show_tunnel_status_change_detection_time",
		"CID":    c.CID,
		"entity": aviatrixEntity,
	}

	type TunnelAPIResults struct {
		DetectionTime int `json:"detection_time"`
	}

	type TunnelAPIResp struct {
		Return  bool             `json:"return"`
		Reason  string           `json:"reason"`
		Results TunnelAPIResults `json:"results"`
	}

	var resp TunnelAPIResp
	err := c.GetAPIContext(ctx, &resp, data["action"], data, BasicCheck)
	if err != nil {
		return 0, err
	}

	return resp.Results.DetectionTime, nil
}
