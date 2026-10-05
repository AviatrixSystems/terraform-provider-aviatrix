package goaviatrix

import (
	"context"
)

type EnforcementLevel struct {
	Level string `json:"level,omitempty"`
}

func (c *Client) SetEnforcementLevel(ctx context.Context, enforcementLevel *EnforcementLevel) error {
	var endpoint string
	if enforcementLevel.Level == "Strict" {
		endpoint = "mitm/origin-cert-verify?level=ENFORCED"
	} else if enforcementLevel.Level == "Ignore" {
		endpoint = "mitm/origin-cert-verify?level=DISABLED"
	} else {
		endpoint = "mitm/origin-cert-verify?level=PERMISSIVE"
	}
	return c.PutAPIContext25(ctx, endpoint, endpoint)
}

func (c *Client) GetEnforcementLevel(ctx context.Context) (*EnforcementLevel, error) {
	endpoint := "mitm/origin-cert-verify-2"

	var data EnforcementLevel
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) UpdateEnforcementLevel(ctx context.Context, enforcementLevel *EnforcementLevel) error {
	return c.SetEnforcementLevel(ctx, enforcementLevel)
}

func (c *Client) DeleteEnforcementLevel(ctx context.Context) error {
	endpoint := "mitm/origin-cert-verify?level=PERMISSIVE"
	return c.PutAPIContext25(ctx, endpoint, nil)
}
