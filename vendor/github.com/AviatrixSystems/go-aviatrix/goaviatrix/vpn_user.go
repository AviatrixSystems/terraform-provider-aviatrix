package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

// VPNUser simple struct to hold vpn_user details
type VPNUser struct {
	Action       string `form:"action,omitempty" json:"action,omitempty"`
	CID          string `form:"CID,omitempty" json:"CID,omitempty"`
	SamlEndpoint string `form:"saml_endpoint,omitempty" json:"saml_endpoint,omitempty"`
	VpcID        string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	GwName       string `form:"lb_name,omitempty" json:"lb_name,omitempty"`
	UserName     string `form:"username" json:"_id,omitempty"`
	UserEmail    string `form:"user_email,omitempty" json:"email,omitempty"`
}

type VPNUserListResp struct {
	Return  bool      `json:"return"`
	Results []VPNUser `json:"results"`
	Reason  string    `json:"reason"`
}

func (c *Client) CreateVPNUser(vpnUser *VPNUser) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for add_vpn_user") + err.Error())
	}
	addVpnUser := url.Values{}
	addVpnUser.Add("CID", c.CID)
	addVpnUser.Add("action", "add_vpn_user")
	addVpnUser.Add("vpc_id", vpnUser.VpcID)
	addVpnUser.Add("username", vpnUser.UserName)
	addVpnUser.Add("user_email", vpnUser.UserEmail)
	addVpnUser.Add("lb_name", vpnUser.GwName)
	addVpnUser.Add("saml_endpoint", vpnUser.SamlEndpoint)
	Url.RawQuery = addVpnUser.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get add_vpn_user failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_vpn_user failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_vpn_user Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetVPNUser(vpnUser *VPNUser) (*VPNUser, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_vpn_users") + err.Error())
	}
	listVpnUsers := url.Values{}
	listVpnUsers.Add("CID", c.CID)
	listVpnUsers.Add("action", "list_vpn_users")
	Url.RawQuery = listVpnUsers.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_vpn_users failed: " + err.Error())
	}
	var data VPNUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpn_users failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpn_users Get failed: " + data.Reason)
	}
	vuList := data.Results
	for i := range vuList {
		if vuList[i].UserName == vpnUser.UserName {
			return &vuList[i], nil
		}
	}
	log.Printf("VPNUser %s not found", vpnUser.UserName)
	return nil, ErrNotFound
}

func (c *Client) DeleteVPNUser(vpnUser *VPNUser) error {
	vpnUser.Action = "delete_vpn_user"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s&vpc_id=%s&username=%s", c.CID, vpnUser.Action,
		vpnUser.VpcID, vpnUser.UserName)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return errors.New("HTTP Get delete_vpn_user failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_vpn_user failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_vpn_user Get failed: " + data.Reason)
	}
	return nil
}
