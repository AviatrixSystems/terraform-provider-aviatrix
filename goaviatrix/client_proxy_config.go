package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type ClientProxyConfig struct {
	HttpProxy          string `form:"http_proxy,omitempty" json:"http_proxy,omitempty"`
	HttpsProxy         string `form:"https_proxy,omitempty" json:"https_proxy,omitempty"`
	ProxyCaCertificate string `form:"server_ca_cert,omitempty" json:"server_ca_cert,omitempty"`
}

type ClientProxyConfigResp struct {
	Return  bool              `json:"return"`
	Results ClientProxyConfig `json:"results"`
	Reason  string            `json:"reason"`
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for show_proxy_config") + err.Error())
	}
	showProxyConfig := url.Values{}
	showProxyConfig.Add("CID", c.CID)
	showProxyConfig.Add("action", "show_proxy_config")
	Url.RawQuery = showProxyConfig.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get show_proxy_config failed: " + err.Error())
	}
	var data ClientProxyConfigResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode show_proxy_config failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API show_proxy_config Get failed: " + data.Reason)
	}
	if data.Results.HttpProxy != "" && data.Results.HttpsProxy != "" {
		httpProxy := data.Results.HttpProxy
		if strings.HasPrefix(httpProxy, "http://") {
			httpProxy = httpProxy[7:]
		}
		httpsProxy := data.Results.HttpsProxy
		if strings.HasPrefix(httpsProxy, "http://") {
			httpsProxy = httpsProxy[7:]
		}
		return &ClientProxyConfig{
			HttpProxy:  httpProxy,
			HttpsProxy: httpsProxy,
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
