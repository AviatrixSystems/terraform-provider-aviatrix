package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

type ClientProxyConfig struct {
	HttpProxy          string `form:"http_proxy,omitempty" json:"http_proxy,omitempty"`
	HttpsProxy         string `form:"https_proxy,omitempty" json:"https_proxy,omitempty"`
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

		files := []File{
			{
				Path:      clientProxyConfig.ProxyCaCertificate,
				ParamName: "server_ca_cert",
			},
		}
		resp, err := c.PostFile(c.baseURL, params, files)
		if err != nil {
			return errors.New("HTTP Post apply_proxy_config failed: " + err.Error())
		}
		var data APIResp
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		bodyString := buf.String()
		bodyIoCopy := strings.NewReader(bodyString)
		if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
			return errors.New("Json Decode apply_proxy_config failed: " + err.Error() + "\n Body: " + bodyString)
		}
		if !data.Return {
			return errors.New("Rest API apply_proxy_config Post failed: " + data.Reason)
		}
		return nil
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
