package goaviatrix

import (
	"strconv"
)

type FilebeatForwarder struct {
	CID                   string
	Server                string
	Port                  int
	TrustedCAFile         string
	ConfigFile            string
	ExcludedGatewaysInput string
}

type FilebeatForwarderResp struct {
	Server           string   `json:"server"`
	Port             string   `json:"port"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
}

func (c *Client) EnableFilebeatForwarder(r *FilebeatForwarder) error {
	params := map[string]string{
		"action":               "enable_logstash_logging",
		"CID":                  c.CID,
		"server_ip":            r.Server,
		"port":                 strconv.Itoa(r.Port),
		"exclude_gateway_list": r.ExcludedGatewaysInput,
		"server_type":          "other",
		"forwarder_type":       "filebeat",
	}

	var files []File

	if r.TrustedCAFile != "" {
		ca := File{
			ParamName:      "trusted_ca",
			UseFileContent: true,
			FileName:       "ca.pem", // fake name for ca
			FileContent:    r.TrustedCAFile,
		}
		files = append(files, ca)
	}

	if r.ConfigFile != "" {
		config := File{
			ParamName:      "config_file",
			UseFileContent: true,
			FileName:       "config.txt", // fake name for config file
			FileContent:    r.ConfigFile,
		}
		files = append(files, config)
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) GetFilebeatForwarderStatus() (*FilebeatForwarderResp, error) {
	params := map[string]string{
		"action": "get_logstash_logging_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool                  `json:"return,omitempty"`
		Results FilebeatForwarderResp `json:"results,omitempty"`
		Reason  string                `json:"reason,omitempty"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	if data.Results.Status == "disabled" {
		return nil, ErrNotFound
	}

	return &data.Results, nil
}

func (c *Client) DisableFilebeatForwarder() error {
	params := map[string]string{
		"action": "disable_logstash_logging",
		"CID":    c.CID,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
