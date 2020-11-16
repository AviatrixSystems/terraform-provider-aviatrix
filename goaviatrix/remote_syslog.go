package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

type RemoteSyslog struct {
	CID                 string   `form:"CID,omitempty"`
	Server              string   `form:"server,omitempty" json:"server"`
	Port                int      `form:"port,omitempty" json:"port"`
	Protocol            string   `form:"protocol,omitempty" json:"protocol"`
	Index               int      `form:"index,omitempty" json:"index"`
	Template            string   `form:"template,omitempty" json:"template"`
	CaCertificate       string   `form:"ca_certificate,omitempty"`
	PublicCertificate   string   `form:"public_certificate,omitempty"`
	PrivateKey          string   `form:"private_key,omitempty"`
	ExcludeGatewayInput string   `form:"exclude_gateway_list,omitempty"`
	ExcludedGateway     []string `json:"excluded_gateway"`
	Status              string   `json:"status"`
	Notls               bool     `json:"notls"`
}

type RemoteSyslogResp struct {
	Server          string   `json:"server"`
	Port            string   `json:"port"`
	Protocol        string   `json:"protocol"`
	Index           string   `json:"index"`
	Template        string   `json:"template"`
	ExcludedGateway []string `json:"excluded_gateway"`
	Status          string   `json:"status"`
	Notls           bool     `json:"notls"`
}

func (c *Client) EnableRemoteSyslog(r *RemoteSyslog) error {
	params := map[string]string{
		"action":   "enable_remote_syslog_logging",
		"CID":      c.CID,
		"index":    strconv.Itoa(r.Index),
		"server":   r.Server,
		"port":     strconv.Itoa(r.Port),
		"protocol": r.Protocol,
		"template": r.Template,
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

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post enable_remote_syslog_logging failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}

	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body enable_remote_syslog_logging failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode enable_remote_syslog_logging failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API enable_remote_syslog_logging Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetRemoteSyslogStatus(idx int) (*RemoteSyslogResp, error) {
	params := map[string]string{
		"action": "get_remote_syslog_logging_status",
		"CID":    c.CID,
		"index":  strconv.Itoa(idx),
	}

	type Resp struct {
		Return  bool             `json:"return,omitempty"`
		Results RemoteSyslogResp `json:"results,omitempty"`
		Reason  string           `json:"reason,omitempty"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) DisableRemoteSyslog(idx int) error {
	params := map[string]string{
		"action": "disable_remote_syslog_logging",
		"CID":    c.CID,
		"index":  strconv.Itoa(idx),
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}
