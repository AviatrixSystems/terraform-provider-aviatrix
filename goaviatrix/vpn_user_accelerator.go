package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

type VpnUserXlr struct {
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	Endpoints      string `form:"endpoints,omitempty"`
	AllEndpoints   string `json:"all,omitempty"`
	FreeEndpoints  string `json:"free,omitempty"`
	InUseEndpoints string `json:"inuse,omitempty"`
}

type VpnUserXlrAPIResp struct {
	Return  bool                `json:"return"`
	Results map[string][]string `json:"results"`
	Reason  string              `json:"reason"`
}

func (c *Client) GetVpnUserAccelerator() ([]string, error) {
	xlr := VpnUserXlr{}
	xlr.CID = c.CID
	xlr.Action = "list_vpn_user_xlr"
	resp, err := c.Post(c.baseURL, xlr)

	if err != nil {
		return nil, errors.New("HTTP Get list_vpn_user_xlr failed: " + err.Error())
	}
	var data VpnUserXlrAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpn_user_xlr failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpn_user_xlr Get failed: " + data.Reason)
	}

	return data.Results["inuse"], nil
}

func (c *Client) UpdateVpnUserAccelerator(xlr *VpnUserXlr) error {
	xlr.CID = c.CID
	xlr.Action = "update_vpn_user_xlr"
	resp, err := c.Post(c.baseURL, xlr)
	if err != nil {
		return errors.New("HTTP Post update_vpn_user_xlr failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode update_vpn_user_xlr failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API update_vpn_user_xlr Get failed: " + data.Reason)
	}
	return nil
}
