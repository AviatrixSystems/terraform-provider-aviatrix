package goaviatrix

import (
	"encoding/json"
	"errors"
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
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpn_user_xlr failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpn_user_xlr Get failed: " + data.Reason)
	}

	elbList := make([]string, 0)
	elbList = data.Results["inuse"]
	return elbList, nil
}

func (c *Client) UpdateVpnUserAccelerator(xlr *VpnUserXlr) error {
	xlr.CID = c.CID
	xlr.Action = "update_vpn_user_xlr"
	resp, err := c.Post(c.baseURL, xlr)

	if err != nil {
		return errors.New("HTTP Post update_vpn_user_xlr failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode update_vpn_user_xlr failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API update_vpn_user_xlr Get failed: " + data.Reason)
	}
	return nil
}
