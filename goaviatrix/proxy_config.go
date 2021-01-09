package goaviatrix

import (
	"strings"
)

type ProxyConfig struct {
	HttpProxy          string `json:"http_proxy,omitempty"`
	HttpsProxy         string `json:"https_proxy,omitempty"`
	ProxyCaCertificate string
}

type ProxyConfigResp struct {
	Return  bool
	Results ProxyConfig
	Reason  string
}

func (c *Client) CreateProxyConfig(proxyConfig *ProxyConfig) error {
	action := "apply_proxy_config"
	params := map[string]string{
		"CID":         c.CID,
		"action":      action,
		"http_proxy":  proxyConfig.HttpProxy,
		"https_proxy": proxyConfig.HttpsProxy,
	}

	if proxyConfig.ProxyCaCertificate != "" {
		var files []File
		ca := File{
			ParamName:      "server_ca_cert",
			UseFileContent: true,
			FileName:       "ca.pem", // fake name for ca
			FileContent:    proxyConfig.ProxyCaCertificate,
		}
		files = append(files, ca)
		return c.PostFileAPI(params, files, BasicCheck)
	} else {
		return c.PostAPI(action, params, BasicCheck)
	}
}

func (c *Client) GetProxyConfig() (*ProxyConfig, error) {
	formData := map[string]string{
		"action": "show_proxy_config",
		"CID":    c.CID,
	}
	var data ProxyConfigResp
	err := c.GetAPI(&data, formData["action"], formData, BasicCheck)
	if err != nil {
		return nil, err
	}
	if data.Results.HttpProxy != "" && data.Results.HttpsProxy != "" {
		return &ProxyConfig{
			HttpProxy:  strings.TrimPrefix(data.Results.HttpProxy, "http://"),
			HttpsProxy: strings.TrimPrefix(data.Results.HttpsProxy, "http://"),
		}, nil
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteProxyConfig() error {
	action := "delete_proxy_config"
	data := map[string]interface{}{
		"action": action,
		"CID":    c.CID,
	}
	return c.PostAPI(action, data, BasicCheck)
}
