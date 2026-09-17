package goaviatrix

type HTTPSCertConfig struct {
	CACertificateFilePath     string
	ServerCertificateFilePath string
	ServerPrivateKeyFilePath  string
	CACertificateFile         string
	ServerCertificateFile     string
	ServerPrivateKeyFile      string
}

func (c *Client) ImportNewHTTPSCerts(certConfig *HTTPSCertConfig) error {
	data := map[string]string{
		"action": "import_new_https_certs",
		"CID":    c.CID,
	}

	var files []File

	if certConfig.CACertificateFilePath != "" {
		files = []File{
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
	} else {
		files = []File{
			{
				ParamName:      "ca_cert",
				UseFileContent: true,
				FileName:       "ca_cert.pem", // fake name for ca cert
				FileContent:    certConfig.CACertificateFile,
			},
			{
				ParamName:      "server_cert",
				UseFileContent: true,
				FileName:       "server_cert.pem", // fake name for server cert
				FileContent:    certConfig.ServerCertificateFile,
			},
			{
				ParamName:      "private_key",
				UseFileContent: true,
				FileName:       "private_key.pem", // fake name for private key
				FileContent:    certConfig.ServerPrivateKeyFile,
			},
		}
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
