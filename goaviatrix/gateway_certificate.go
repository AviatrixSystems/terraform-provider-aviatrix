package goaviatrix

import "context"

type GatewayCertificate struct {
	CaCertificate string
	CaPrivateKey  string
}

type GatewayCertificateStatusResp struct {
	APIResp
	Results GatewayCertificateStatus
}

type GatewayCertificateStatus struct {
	Status string
}

func (c *Client) ConfigureGatewayCertificate(ctx context.Context, gwCert *GatewayCertificate) error {
	data := map[string]string{
		"action": "import_gateway_ca_certificate",
		"CID":    c.CID,
	}
	files := []File{
		{
			ParamName:      "ca_cert",
			FileName:       "ca_root_cert.pem", // fake name for file
			FileContent:    gwCert.CaCertificate,
			UseFileContent: true,
		},
		{
			ParamName:      "private_key",
			FileName:       "ca_private.key", // fake name for file
			FileContent:    gwCert.CaPrivateKey,
			UseFileContent: true,
		},
	}
	return c.PostFileAPIContext(ctx, data, files, BasicCheck)
}

func (c *Client) DisableGatewayCertificate(ctx context.Context) error {
	params := map[string]string{
		"action": "disable_certificate_checking",
		"CID":    c.CID,
	}
	return c.PostAPIContext(ctx, params["action"], params, BasicCheck)
}

func (c *Client) GetGatewayCertificateStatus(ctx context.Context) (string, error) {
	formData := map[string]string{
		"action": "get_gateway_ca_certificate_status",
		"CID":    c.CID,
	}
	var data GatewayCertificateStatusResp
	err := c.GetAPIContext(ctx, &data, formData["action"], formData, BasicCheck)
	if err != nil {
		return "", err
	}
	return data.Results.Status, nil
}
