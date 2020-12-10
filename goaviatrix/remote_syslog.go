package goaviatrix

import (
	"log"
	"strconv"
)

type RemoteSyslog struct {
	CID                 string `form:"CID,omitempty"`
	Server              string `form:"server,omitempty" json:"server"`
	Port                int    `form:"port,omitempty" json:"port"`
	Protocol            string `form:"protocol,omitempty" json:"protocol"`
	Index               int    `form:"index,omitempty" json:"index"`
	Template            string `form:"template,omitempty" json:"template"`
	CaCertificate       string `form:"ca_certificate,omitempty"`
	PublicCertificate   string `form:"public_certificate,omitempty"`
	PrivateKey          string `form:"private_key,omitempty"`
	ExcludeGatewayInput string `form:"exclude_gateway_list,omitempty"`
}

type RemoteSyslogResp struct {
	Server           string   `json:"server"`
	Port             string   `json:"port"`
	Protocol         string   `json:"protocol"`
	Index            string   `json:"index"`
	Template         string   `json:"template"`
	ExcludedGateways []string `json:"excluded_gateway"`
	Status           string   `json:"status"`
	Notls            bool     `json:"notls"`
}

func (c *Client) EnableRemoteSyslog(r *RemoteSyslog) error {
	params := map[string]string{
		"action":               "enable_remote_syslog_logging",
		"CID":                  c.CID,
		"index":                strconv.Itoa(r.Index),
		"server":               r.Server,
		"port":                 strconv.Itoa(r.Port),
		"protocol":             r.Protocol,
		"template":             r.Template,
		"exclude_gateway_list": r.ExcludeGatewayInput,
	}

	files := []File{
		{
			Path:      r.CaCertificate,
			ParamName: "ca_certificate",
		},
		{
			Path:      r.PublicCertificate,
			ParamName: "public_certificate",
		},
		{
			Path:      r.PrivateKey,
			ParamName: "private_key",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) GetRemoteSyslogStatus(idx int) (*RemoteSyslogResp, error) {
	params := map[string]string{
		"action": "get_remote_syslog_logging_status",
		"CID":    c.CID,
		"index":  strconv.Itoa(idx),
	}

	type Resp struct {
		Return  bool             `json:"return"`
		Results RemoteSyslogResp `json:"results"`
		Reason  string           `json:"reason"`
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

func (c *Client) DisableRemoteSyslog(idx int) error {
	params := map[string]string{
		"action": "disable_remote_syslog_logging",
		"CID":    c.CID,
		"index":  strconv.Itoa(idx),
	}

	log.Printf("[INFO] Deleting remote syslog index %d", idx)

	return c.PostAPI(params["action"], params, BasicCheck)
}
