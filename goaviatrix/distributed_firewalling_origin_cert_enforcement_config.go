package goaviatrix

import (
	"context"
	"fmt"
)

type EnforcementLevel struct {
	Level string `json:"level,omitempty"`
}

func (c *Client) SetEnforcementLevel(ctx context.Context, enforcementLevel *EnforcementLevel) error {
	endpoint := fmt.Sprintf("mitm/origin-cert-verify?level=%s", enforcementLevel.Level)
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
	endpoint := fmt.Sprintf("mitm/origin-cert-verify?level=%s", enforcementLevel.Level)
	return c.PutAPIContext25(ctx, endpoint, enforcementLevel)
}

func (c *Client) DeleteEnforcementLevel(ctx context.Context) error {
	endpoint := "mitm/origin-cert-verify?level=PERMISSIVE"
	return c.PutAPIContext25(ctx, endpoint, nil)
}
