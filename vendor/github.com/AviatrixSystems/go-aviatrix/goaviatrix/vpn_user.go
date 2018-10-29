package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (c *Client) CreateVPNUser(vpn_user *VPNUser) error {
	vpn_user.Action = "add_vpn_user"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s&vpc_id=%s&username=%s&user_email=%s&lb_name=%s&saml_endpoint=%s", c.CID, vpn_user.Action, vpn_user.VpcID, vpn_user.UserName, vpn_user.UserEmail, vpn_user.GwName, vpn_user.SamlEndpoint)

	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetVPNUser(vpn_user *VPNUser) (*VPNUser, error) {
	vpn_user.Action = "list_vpn_users"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s", c.CID, vpn_user.Action)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data VPNUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}
	vulist := data.Results
	for i := range vulist {
		if vulist[i].UserName == vpn_user.UserName {
			return &vulist[i], nil
		}
	}
	log.Printf("VPNUser %s not found", vpn_user.UserName)
	return nil, ErrNotFound
}

func (c *Client) DeleteVPNUser(vpn_user *VPNUser) error {
	vpn_user.Action = "delete_vpn_user"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s&vpc_id=%s&username=%s", c.CID, vpn_user.Action, vpn_user.VpcID, vpn_user.UserName)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
