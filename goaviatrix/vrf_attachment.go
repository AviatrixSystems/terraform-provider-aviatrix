package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

// VrfAttachmentStatus mirrors a single entry of the controller's
// get_attachment_vrf_status response.
type VrfAttachmentStatus struct {
	PeeringName          string `json:"peering_name"`
	Gateway1             string `json:"gateway1"`
	Gateway2             string `json:"gateway2"`
	AttachmentType       string `json:"attachment_type"`
	VrfAttachmentEnabled bool   `json:"vrf_attachment_enabled"`
}

type vrfAttachmentUpdateRequest struct {
	CID                   string `json:"CID"`
	Action                string `json:"action"`
	Gateway1              string `json:"gateway1,omitempty"`
	Gateway2              string `json:"gateway2,omitempty"`
	All                   bool   `json:"all,omitempty"`
	EnableVrfOnAttachment string `json:"enable,omitempty"`
}

type vrfAttachmentStatusRequest struct {
	CID      string `json:"CID"`
	Action   string `json:"action"`
	Gateway1 string `json:"gateway1,omitempty"`
	Gateway2 string `json:"gateway2,omitempty"`
	All      bool   `json:"all,omitempty"`
}

type vrfAttachmentStatusResp struct {
	Return  bool                  `json:"return"`
	Results []VrfAttachmentStatus `json:"results"`
	Reason  string                `json:"reason"`
}

// UpdateVrfOnAttachment toggles vrf_attachment_enabled on a single eligible
// peering pair, or on all eligible peerings when all=true. The enable arg is
// translated to "yes"/"no" before being sent to the controller.
func (c *Client) UpdateVrfOnAttachment(ctx context.Context, gw1, gw2 string, enable, all bool) error {
	enableStr := "no"
	if enable {
		enableStr = "yes"
	}
	req := vrfAttachmentUpdateRequest{
		CID:                   c.CID,
		Action:                "update_vrf_on_attachment",
		EnableVrfOnAttachment: enableStr,
	}
	if all {
		req.All = true
	} else {
		req.Gateway1 = gw1
		req.Gateway2 = gw2
	}
	return c.PostAPIContext2(ctx, nil, req.Action, req, BasicCheck)
}

// GetAttachmentVrfStatus returns vrf_attachment_enabled status for one peering
// pair or for all eligible peerings when all=true.
//
// Controller returns "Peering not found or not eligible for VRF attachment"
// when a specific pair is missing/ineligible; this is mapped to ErrNotFound.
func (c *Client) GetAttachmentVrfStatus(ctx context.Context, gw1, gw2 string, all bool) ([]VrfAttachmentStatus, error) {
	req := vrfAttachmentStatusRequest{
		CID:    c.CID,
		Action: "get_attachment_vrf_status",
	}
	if all {
		req.All = true
	} else {
		req.Gateway1 = gw1
		req.Gateway2 = gw2
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "not eligible") || strings.Contains(reason, "not found") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	var resp vrfAttachmentStatusResp
	if err := c.PostAPIContext2(ctx, &resp, req.Action, req, check); err != nil {
		return nil, err
	}
	return resp.Results, nil
}
