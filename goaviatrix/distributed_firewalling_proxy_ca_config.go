package goaviatrix

import (
	"context"
	"fmt"
)

type ProxyCaConfig struct {
	CaCert string `json:"cert,omitempty"`
	CaKey string `json:"ca_key,omitempty"`
}

func (c *Client) SetEnforcementLevel(ctx context.Context, proxyCaConfig *ProxyCaConfig) error {
	endpoint := fmt.Sprintf("mitm/ca")
	return c.PutAPIContext25(ctx, endpoint, proxyCaConfig)
}

//func (c *Client) GetEnforcementLevel(ctx context.Context) (*ProxyCaConfig, error) {
//	endpoint := "mitm/origin-cert-verify-2"
//
//	var data EnforcementLevel
//	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return &data, nil
//}
//
//func (c *Client) UpdateEnforcementLevel(ctx context.Context, enforcementLevel *ProxyCaConfig) error {
//	endpoint := fmt.Sprintf("mitm/origin-cert-verify?level=%s", enforcementLevel.Level)
//	return c.PutAPIContext25(ctx, endpoint, enforcementLevel)
//}
//
//func (c *Client) DeleteEnforcementLevel(ctx context.Context) error {
//	endpoint := "mitm/origin-cert-verify?level=PERMISSIVE"
//	return c.PutAPIContext25(ctx, endpoint, nil)
//}