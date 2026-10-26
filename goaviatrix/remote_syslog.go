package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RemoteSyslog struct {
	CID                 string `form:"CID,omitempty"`
	Server              string `form:"server,omitempty" json:"server"`
	Port                int    `form:"port,omitempty" json:"port"`
	Protocol            string `form:"protocol,omitempty" json:"protocol"`
	Index               int    `form:"index,omitempty" json:"index"`
	Name                string `form:"name,omitempty"`
	Template            string `form:"template,omitempty" json:"template"`
	CaCertificate       string `form:"ca_certificate,omitempty"`
	PublicCertificate   string `form:"public_certificate,omitempty"`
	PrivateKey          string `form:"private_key,omitempty"`
	ExcludeGatewayInput string `form:"exclude_gateway_list,omitempty"`
}

type PortWrapper string

func (w *PortWrapper) UnmarshalJSON(data []byte) (err error) {
	if port, err := strconv.Atoi(string(data)); err == nil {
		str := strconv.Itoa(port)
		*w = PortWrapper(str)
		return nil
	}
	var str string
	err = myUnmarshal(data, &str)
	if err != nil {
		return err
	}
	return myUnmarshal([]byte(str), w)
}

func myUnmarshal(input []byte, target interface{}) error {
	if len(input) == 0 {
		return nil
	}
	return json.Unmarshal(input, target)
}

type RemoteSyslogResp struct {
	Server           string      `json:"server"`
	Port             PortWrapper `json:"port"`
	Protocol         string      `json:"protocol"`
	Index            string      `json:"index"`
	Name             string      `json:"name"`
	Template         string      `json:"template"`
	ExcludedGateways []string    `json:"excluded_gateway"`
	Status           string      `json:"status"`
	Notls            bool        `json:"notls"`
}

func (c *Client) EnableRemoteSyslog(r *RemoteSyslog) error {
	params := map[string]string{
		"action":               "enable_remote_syslog_logging",
		"CID":                  c.CID,
		"index":                strconv.Itoa(r.Index),
		"name":                 r.Name,
		"server":               r.Server,
		"port":                 strconv.Itoa(r.Port),
		"protocol":             r.Protocol,
		"template":             r.Template,
		"exclude_gateway_list": r.ExcludeGatewayInput,
	}

	var files []File

	if r.CaCertificate != "" {
		ca := File{
			ParamName:      "ca_certificate",
			UseFileContent: true,
			FileName:       "ca.pem", // fake name for ca
			FileContent:    r.CaCertificate,
		}
		files = append(files, ca)
	}

	if r.PublicCertificate != "" {
		ca := File{
			ParamName:      "public_certificate",
			UseFileContent: true,
			FileName:       "public.pem", // fake name for public certificate
			FileContent:    r.PublicCertificate,
		}
		files = append(files, ca)
	}

	if r.PrivateKey != "" {
		ca := File{
			ParamName:      "private_key",
			UseFileContent: true,
			FileName:       "private.pem", // fake name for private key
			FileContent:    r.PrivateKey,
		}
		files = append(files, ca)
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

	err := c.GetAPIRemoteSyslog(&data, params["action"], params, BasicCheck)
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

func (c *Client) GetAPIRemoteSyslog(v interface{}, action string, d map[string]string, checkFunc CheckAPIResponseFunc) error {
	Url, err := c.urlEncode(d)
	if err != nil {
		return fmt.Errorf("could not url encode values for action %q: %v", action, err)
	}

	try, maxTries, backoff := 0, 5, 500*time.Millisecond
	var resp *http.Response
	for {
		try++
		resp, err = c.GetContext(context.Background(), Url, nil)
		if err == nil {
			break
		}

		if try == maxTries {
			return fmt.Errorf("HTTP Get %s failed: %v", action, err)
		}
		time.Sleep(backoff)
		// Double the backoff time after each failed try
		backoff *= 2
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	var data APIResp
	if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&data); err != nil {
		return fmt.Errorf("Json Decode into standard format failed: %v\n Body: %s", err, bodyString)
	}
	if err := checkFunc(action, "Get", data.Reason, data.Return); err != nil {
		return err
	}

	err = myUnmarshal(buf.Bytes(), &v)
	if err != nil {
		return fmt.Errorf("json unmarshal failed: %v\n Body: %s", err, buf.String())
	}

	return nil
}
