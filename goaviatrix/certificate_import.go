package goaviatrix

type HTTPSCertConfig struct {
	CACertificateFilePath     string
	ServerCertificateFilePath string
	ServerPrivateKeyFilePath  string
}

func (c *Client) ImportNewHTTPSCerts(certConfig *HTTPSCertConfig) error {
	data := map[string]string{
		"action": "import_new_https_certs",
		"CID":    c.CID,
	}
	files := []File{
		{
			Path:      certConfig.CACertificateFilePath,
			ParamName: "ca_cert",
		},
		{
			Path:      certConfig.ServerCertificateFilePath,
			ParamName: "server_cert",
		},
		{
			Path:      certConfig.ServerPrivateKeyFilePath,
			ParamName: "private_key",
		},
	}
	return c.PostFileAPI(data, files, BasicCheck)
}

func (c *Client) DisableImportedHTTPSCerts() error {
	data := map[string]string{
		"action": "disable_imported_certificate",
		"CID":    c.CID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

type GetHTTPSCertsStatusResp struct {
	Return  bool
	Reason  string
	Results GetHTTPSCertsStatus
}
type GetHTTPSCertsStatus struct {
	Method string
}

func (c *Client) GetHTTPSCertsStatus() (bool, error) {
	data := map[string]string{
		"action": "get_https_certs_status",
		"CID":    c.CID,
	}
	var respData GetHTTPSCertsStatusResp
	err := c.GetAPI(&respData, data["action"], data, BasicCheck)
	if err != nil {
		return false, err
	}
	return respData.Results.Method == "standard", nil
}
