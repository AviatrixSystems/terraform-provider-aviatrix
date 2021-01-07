package goaviatrix

import (
	"strings"
)

type ClientProxyConfig struct {
	HttpProxy          string `json:"http_proxy,omitempty"`
	HttpsProxy         string `json:"https_proxy,omitempty"`
	ProxyCaCertificate string
}

type ClientProxyConfigResp struct {
	Return  bool
	Results ClientProxyConfig
	Reason  string
}

func (c *Client) CreateClientProxyConfig(clientProxyConfig *ClientProxyConfig) error {
	action := "apply_proxy_config"
	if clientProxyConfig.ProxyCaCertificate != "" {
		params := map[string]string{
			"CID":         c.CID,
			"action":      action,
			"http_proxy":  clientProxyConfig.HttpProxy,
			"https_proxy": clientProxyConfig.HttpsProxy,
		}

		var files []File

		if clientProxyConfig.ProxyCaCertificate != "" {
			ca := File{
				ParamName:      "server_ca_cert",
				UseFileContent: true,
				FileName:       "ca.pem", // fake name for ca
				FileContent:    clientProxyConfig.ProxyCaCertificate,
			}
			files = append(files, ca)
		}
		return c.PostFileAPI(params, files, BasicCheck)
	} else {
		data := map[string]interface{}{
			"CID":         c.CID,
			"action":      action,
			"http_proxy":  clientProxyConfig.HttpProxy,
			"https_proxy": clientProxyConfig.HttpsProxy,
		}
		return c.PostAPI(action, data, BasicCheck)
	}
}

func (c *Client) GetClientProxyConfig() (*ClientProxyConfig, error) {
	formData := map[string]string{
		"action": "show_proxy_config",
		"CID":    c.CID,
	}
	var data ClientProxyConfigResp
	err := c.GetAPI(&data, formData["action"], formData, BasicCheck)
	if err != nil {
		return nil, err
	}
	if data.Results.HttpProxy != "" && data.Results.HttpsProxy != "" {
		return &ClientProxyConfig{
			HttpProxy:  strings.TrimPrefix(data.Results.HttpProxy, "http://"),
			HttpsProxy: strings.TrimPrefix(data.Results.HttpsProxy, "http://"),
		}, nil
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteClientProxyConfig() error {
	action := "delete_proxy_config"
	data := map[string]interface{}{
		"action": action,
		"CID":    c.CID,
	}
	return c.PostAPI(action, data, BasicCheck)
}
