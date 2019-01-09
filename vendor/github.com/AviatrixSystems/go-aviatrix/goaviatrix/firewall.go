package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Policy struct {
	SrcIP     string `form:"s_ip,omitempty" json:"s_ip,omitempty"`
	DstIP     string `form:"d_ip,omitempty" json:"d_ip,omitempty"`
	Protocol  string `form:"protocol,omitempty" json:"protocol,omitempty"`
	Port      string `form:"port,omitempty" json:"port,omitempty"`
	AllowDeny string `form:"deny_allow,omitempty" json:"deny_allow,omitempty"`
	LogEnable string `form:"log_enable,omitempty" json:"log_enable,omitempty"`
}

// Gateway simple struct to hold firewall details
type Firewall struct {
	CID           string    `form:"CID,omitempty"`
	Action        string    `form:"action,omitempty"`
	GwName        string    `form:"vpc_name,omitempty" json:"vpc_name,omitempty"`
	BaseAllowDeny string    `form:"base_policy,omitempty" json:"base_policy,omitempty"`
	BaseLogEnable string    `form:"base_policy_log_enable,omitempty" json:"base_policy_log_enable,omitempty"`
	PolicyList    []*Policy `form:"new_policy[],omitempty" json:"security_rules,omitempty"`
}

type FirewallResp struct {
	Return  bool     `json:"return"`
	Results Firewall `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) SetBasePolicy(firewall *Firewall) error {
	firewall.CID = c.CID
	firewall.Action = "set_vpc_base_policy"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=set_vpc_base_policy&"+
		"vpc_name=%s&base_policy=%s&base_policy_log_enable=%s", c.CID, firewall.GwName, firewall.BaseAllowDeny,
		firewall.BaseLogEnable)
	log.Printf("[INFO] Setting Base Policy: %#v", firewall)

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

func (c *Client) UpdatePolicy(firewall *Firewall) error {
	firewall.CID = c.CID
	firewall.Action = "update_access_policy"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=update_access_policy&vpc_name=%s&new_policy=", c.CID,
		firewall.GwName)
	log.Printf("[INFO] Updating Aviatrix firewall for gateway: %#v", firewall)

	args, err := json.Marshal(firewall.PolicyList)
	if err != nil {
		return err
	}
	path = path + string(args)
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

func (c *Client) GetPolicy(firewall *Firewall) (*Firewall, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=vpc_access_policy&vpc_name=%s", c.CID, firewall.GwName)
	log.Printf("[INFO] Getting Policy: %#v", firewall)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data FirewallResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find Aviatrix Firewall policies for gateway %s: %s", firewall.GwName,
			data.Reason)
		return nil, ErrNotFound
	}
	if data.Results.BaseAllowDeny == "allow-all" {
		data.Results.BaseAllowDeny = "allow"
	} else {
		data.Results.BaseAllowDeny = "deny"
	}

	return &data.Results, nil
}
