package goaviatrix

import (
	"errors"
	"fmt"
	"strings"
)

// VPNUser simple struct to hold vpn_user details
type VPNUser struct {
	Action       string   `form:"action,omitempty" json:"action,omitempty"`
	CID          string   `form:"CID,omitempty" json:"CID,omitempty"`
	SamlEndpoint string   `form:"saml_endpoint,omitempty" json:"saml_endpoint,omitempty"`
	VpcID        string   `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	GwName       string   `form:"lb_name,omitempty" json:"lb_name,omitempty"`
	DnsName      string   `json:"dns,omitempty"`
	DnsEnabled   bool     `json:"dns_enabled,omitempty"`
	UserName     string   `form:"username" json:"_id,omitempty"`
	UserEmail    string   `form:"user_email,omitempty" json:"email,omitempty"`
	Profiles     []string `json:"profiles,omitempty"`
}

type VPNUserResp struct {
	Return  bool        `json:"return"`
	Results VPNUserInfo `json:"results"`
	Reason  string      `json:"reason"`
}

type VPNUserInfo struct {
	VpnUser VPNUser `json:"vpn_user"`
}

func (c *Client) CreateVPNUser(vpnUser *VPNUser) error {
	form := map[string]string{
		"CID":           c.CID,
		"action":        "add_vpn_user",
		"username":      vpnUser.UserName,
		"user_email":    vpnUser.UserEmail,
		"saml_endpoint": vpnUser.SamlEndpoint,
	}

	if vpnUser.DnsEnabled {
		form["dns"] = "true"
		form["lb_name"] = vpnUser.DnsName
	} else {
		form["vpc_id"] = vpnUser.VpcID
		form["lb_name"] = vpnUser.GwName
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Sending VPN certificates to email") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) GetVPNUser(vpnUser *VPNUser) (*VPNUser, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "get_vpn_user_by_name",
		"username": vpnUser.UserName,
	}

	var data VPNUserResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Invalid VPN username") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	if data.Results.VpnUser.UserName != "" {
		if data.Results.VpnUser.UserName == vpnUser.UserName {
			return &data.Results.VpnUser, nil
		} else {
			return nil, errors.New("VPN user name does not match from response")
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteVPNUser(vpnUser *VPNUser) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "delete_vpn_user",
		"username": vpnUser.UserName,
	}

	if vpnUser.DnsEnabled {
		form["dns"] = "true"
		form["vpc_id"] = vpnUser.DnsName
	} else {
		form["vpc_id"] = vpnUser.VpcID
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
