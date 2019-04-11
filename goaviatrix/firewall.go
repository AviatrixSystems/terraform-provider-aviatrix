package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_spoke_to_transit_gw") + err.Error())
	}
	setVpcBasePolicy := url.Values{}
	setVpcBasePolicy.Add("CID", c.CID)
	setVpcBasePolicy.Add("action", "set_vpc_base_policy")
	setVpcBasePolicy.Add("vpc_name", firewall.GwName)
	setVpcBasePolicy.Add("base_policy", firewall.BaseAllowDeny)
	setVpcBasePolicy.Add("base_policy_log_enable", firewall.BaseLogEnable)
	Url.RawQuery = setVpcBasePolicy.Encode()
	resp, err := c.Get(Url.String(), nil)

	log.Printf("[INFO] Setting Base Policy: %#v", firewall)

	if err != nil {
		return errors.New("HTTP Get set_vpc_base_policy failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode set_vpc_base_policy failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API set_vpc_base_policy Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdatePolicy(firewall *Firewall) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for update_access_policy") + err.Error())
	}
	updateAccessPolicy := url.Values{}
	updateAccessPolicy.Add("CID", c.CID)
	updateAccessPolicy.Add("action", "update_access_policy")
	updateAccessPolicy.Add("vpc_name", firewall.GwName)

	log.Printf("[INFO] Updating Aviatrix firewall for gateway: %#v", firewall)

	args, err := json.Marshal(firewall.PolicyList)
	if err != nil {
		return err
	}
	updateAccessPolicy.Add("new_policy", string(args))
	Url.RawQuery = updateAccessPolicy.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get update_access_policy failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode update_access_policy failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API update_access_policy Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetPolicy(firewall *Firewall) (*Firewall, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for vpc_access_policy") + err.Error())
	}
	vpcAccessPolicy := url.Values{}
	vpcAccessPolicy.Add("CID", c.CID)
	vpcAccessPolicy.Add("action", "vpc_access_policy")
	vpcAccessPolicy.Add("vpc_name", firewall.GwName)
	Url.RawQuery = vpcAccessPolicy.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get vpc_access_policy failed: " + err.Error())
	}
	var data FirewallResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode vpc_access_policy failed: " + err.Error())
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find Aviatrix Firewall policies for gateway %s: %s", firewall.GwName,
			data.Reason)
		return nil, ErrNotFound
	}

	return &data.Results, nil
}

func (c *Client) ValidatePolicy(policy *Policy) error {
	if policy.AllowDeny != "allow" && policy.AllowDeny != "deny" {
		return fmt.Errorf("valid AllowDeny is only 'allow' or 'deny'")
	}
	protocolDefaultValues := []string{"all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp"}
	protocolVal := []string{policy.Protocol}
	if policy.Protocol == "" || len(Difference(protocolVal, protocolDefaultValues)) != 0 {
		return fmt.Errorf("protocal can only be one of {'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'}")
	}
	if policy.Protocol == "all" && policy.Port != "0:65535" {
		return fmt.Errorf("port should be '0:65535' for protocal 'all'")
	}
	if policy.Protocol == "icmp" && (policy.Port != "") {
		return fmt.Errorf("port should be empty for protocal 'icmp'")
	}
	return nil
}
