package goaviatrix

import (
	"context"
)

type AttachmentPointResp struct {
	AttachmentPointID string `json:"attachment_point_id,omitempty"`
	Name              string `json:"name,omitempty"`
}

func (c *Client) GetDCFAttachmentPoint(ctx context.Context, name string) (*AttachmentPointResp, error) {
	endpoint := "microseg/attachment-points"
	var attachmentPoint []AttachmentPointResp

	params := map[string]string{"name": name}
	err := c.GetAPIContext25(ctx, &attachmentPoint, endpoint, params)
	if err != nil {
		return nil, err
	}

	return &attachmentPoint[0], nil
}
