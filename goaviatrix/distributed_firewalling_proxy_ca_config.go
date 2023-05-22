package goaviatrix

import (
	"context"
	"fmt"
)

type ProxyCaConfig struct {
	CaCert         string `json:"cert,omitempty"`
	CaKey          string `json:"ca_key,omitempty"`
	SerialNumber   string `json:"serial_number,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	CommonName     string `json:"common_name,omitempty"`
	ExpirationDate string `json:"expire_date,omitempty"`
	SANs           string `json:"SANs,omitempty"`
	UploadInfo     string `json:"upload_info,omitempty"`
}

func (c *Client) SetNewCertificate(ctx context.Context, proxyCaConfig *ProxyCaConfig) error {
	endpoint := fmt.Sprintf("mitm/ca")
	return c.PutAPIContext25(ctx, endpoint, proxyCaConfig)
}

func (c *Client) GetCaCertificate(ctx context.Context) (*ProxyCaConfig, error) {
	endpoint := "mitm/ca?info"

	var data ProxyCaConfig
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteEnforcementLevel(ctx context.Context) error {
	endpoint := "mitm/ca"
	return c.PutAPIContext25(ctx, endpoint, nil)
}
