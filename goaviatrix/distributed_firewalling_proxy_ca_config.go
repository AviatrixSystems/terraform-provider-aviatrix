package goaviatrix

import (
	"context"
)

type ProxyCaConfig struct {
	CaCert         string `json:"cert,omitempty"`
	CaKey          string `json:"file_key,omitempty"`
	SerialNumber   string `json:"serial_number,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	CommonName     string `json:"common_name,omitempty"`
	ExpirationDate string `json:"expire_date,omitempty"`
	SANs           string `json:"SANs,omitempty"`
	UploadInfo     string `json:"upload_info,omitempty"`
}

type ProxyCaCertInstance struct {
	SerialNumber   string `json:"serial_number,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	CommonName     string `json:"common_name,omitempty"`
	ExpirationDate string `json:"expire_date,omitempty"`
	Sans           string `json:"SANs,omitempty"`
	UploadInfo     string `json:"upload_info,omitempty"`
}

func (c *Client) SetNewCertificate(ctx context.Context, proxyCaConfig *ProxyCaConfig) error {
	endpoint := "mitm/ca"

	files := []File{
		{
			ParamName:      "file_cert",
			UseFileContent: true,
			FileName:       "ca_cert.pem",
			FileContent:    proxyCaConfig.CaCert,
		},
		{
			ParamName:      "file_key",
			UseFileContent: true,
			FileName:       "ca_key.pem",
			FileContent:    proxyCaConfig.CaKey,
		},
	}

	return c.PostFileContext25(ctx, endpoint, nil, files)
}

func (c *Client) GetCaCertificate(ctx context.Context) (*ProxyCaConfig, error) {
	endpoint := "mitm/ca-2"

	var data ProxyCaConfig
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) GetMetaCaCertificate(ctx context.Context) (*ProxyCaCertInstance, error) {
	endpoint := "mitm/ca"
	form := map[string]string{
		"info": "true",
	}
	var data ProxyCaCertInstance
	err := c.GetAPIContext25(ctx, &data, endpoint, form)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteCaCertificate(ctx context.Context) error {
	endpoint := "mitm/ca"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
