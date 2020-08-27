package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Policy struct {
	SrcIP       string `form:"s_ip,omitempty" json:"s_ip,omitempty"`
	DstIP       string `form:"d_ip,omitempty" json:"d_ip,omitempty"`
	Protocol    string `form:"protocol,omitempty" json:"protocol,omitempty"`
	Port        string `form:"port,omitempty" json:"port,omitempty"`
	Action      string `form:"deny_allow,omitempty" json:"deny_allow,omitempty"`
	LogEnabled  string `form:"log_enable,omitempty" json:"log_enable,omitempty"`
	Description string `form:"description,omitempty" json:"description,omitempty"`
}

// Gateway simple struct to hold firewall details
type Firewall struct {
	CID            string    `form:"CID,omitempty"`
	Action         string    `form:"action,omitempty"`
	GwName         string    `form:"vpc_name,omitempty" json:"gw_name,omitempty"`
	BasePolicy     string    `form:"base_policy,omitempty" json:"base_policy,omitempty"`
	BaseLogEnabled string    `form:"base_policy_log_enable,omitempty" json:"base_policy_log_enable,omitempty"`
	PolicyList     []*Policy `json:"security_rules,omitempty"`
	NewPolicy      string    `form:"new_policy,omitempty"`
	GwOriginalName string    `json:"gw_original_name,omitempty"`
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
	setVpcBasePolicy.Add("base_policy", firewall.BasePolicy)
	setVpcBasePolicy.Add("base_policy_log_enable", firewall.BaseLogEnabled)
	Url.RawQuery = setVpcBasePolicy.Encode()
	resp, err := c.Get(Url.String(), nil)
	log.Infof("Setting Base Policy: %#v", firewall)
	if err != nil {
		return errors.New("HTTP Get set_vpc_base_policy failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode set_vpc_base_policy failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API set_vpc_base_policy Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdatePolicy(firewall *Firewall) error {
	firewall.CID = c.CID
	firewall.Action = "update_access_policy"
	// If the PolicyList is nil it will be encoded as 'null'.
	// Instead, we want to set PolicyList to an empty slice so that it is encoded as '[]'.
	if firewall.PolicyList == nil {
		firewall.PolicyList = []*Policy{}
	}
	args, err := json.Marshal(firewall.PolicyList)
	if err != nil {
		return err
	}
	firewall.NewPolicy = string(args)
	resp, err := c.Post(c.baseURL, firewall)
	if err != nil {
		return errors.New("HTTP Post update_access_policy failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode update_access_policy failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode vpc_access_policy failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		log.Errorf("Couldn't find Aviatrix Firewall policies for gateway %s: %s", firewall.GwName, data.Reason)
		return nil, ErrNotFound
	}

	return &data.Results, nil
}

func (c *Client) ValidatePolicy(policy *Policy) error {
	if policy.Action != "allow" && policy.Action != "deny" && policy.Action != "force-drop" {
		return fmt.Errorf("valid 'action' is only 'allow', 'deny' or 'force-drop'")
	}
	protocolDefaultValues := []string{"all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp"}
	protocolVal := []string{policy.Protocol}
	if policy.Protocol == "" || len(Difference(protocolVal, protocolDefaultValues)) != 0 {
		return fmt.Errorf("protocol can only be one of {'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'}")
	}
	if policy.Protocol == "all" && policy.Port != "0:65535" {
		return fmt.Errorf("port should be '0:65535' for protocol 'all'")
	}
	if policy.Protocol == "icmp" && (policy.Port != "") {
		return fmt.Errorf("port should be empty for protocol 'icmp'")
	}
	return nil
}

func PolicyToMap(p *Policy) map[string]interface{} {
	port := p.Port
	if p.Protocol == "all" && p.Port == "" {
		port = "0:65535"
	}

	logEnabled := false
	if p.LogEnabled == "on" {
		logEnabled = true
	}

	return map[string]interface{}{
		"src_ip":      p.SrcIP,
		"dst_ip":      p.DstIP,
		"protocol":    p.Protocol,
		"port":        port,
		"action":      p.Action,
		"log_enabled": logEnabled,
		"description": p.Description,
	}
}

func (c *Client) AddFirewallPolicy(fw *Firewall) error {
	action := "append_stateful_firewall_rules"

	rules, err := json.Marshal(fw.PolicyList)
	if err != nil {
		return fmt.Errorf("could not marshal firewall policies: %v", err)
	}

	return c.PostAPI(action, struct {
		Action string `form:"action"`
		CID    string `form:"CID"`
		GwName string `form:"gateway_name"`
		Rules  string `form:"rules"`
	}{
		Action: action,
		CID:    c.CID,
		GwName: fw.GwName,
		Rules:  string(rules),
	}, BasicCheck)
}

func (c *Client) DeleteFirewallPolicy(fw *Firewall) error {
	action := "delete_stateful_firewall_rules"

	rules, err := json.Marshal(fw.PolicyList)
	if err != nil {
		return fmt.Errorf("could not marshal firewall policies: %v", err)
	}

	return c.PostAPI(action, struct {
		Action string `form:"action"`
		CID    string `form:"CID"`
		GwName string `form:"gateway_name"`
		Rules  string `form:"rules"`
	}{
		Action: action,
		CID:    c.CID,
		GwName: fw.GwName,
		Rules:  string(rules),
	}, func(action string, reason string, ret bool) error {
		if !ret {
			// Tried to delete a rule that did not exist, we don't need to fail the apply.
			if strings.Contains(reason, "Empty rules in Stateful Firewall to delete") {
				return nil
			}

			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	})
}

func (c *Client) GetFirewallPolicy(fw *Firewall) (*Firewall, error) {
	foundFirewall, err := c.GetPolicy(fw)
	if err != nil {
		return nil, fmt.Errorf("could not list firewall rules: %v", err)
	}
	rule := fw.PolicyList[0]
	found := false

	for _, p := range foundFirewall.PolicyList {
		if p.SrcIP == rule.SrcIP &&
			p.DstIP == rule.DstIP &&
			p.Protocol == rule.Protocol &&
			p.Port == rule.Port &&
			p.Action == rule.Action {
			found = true
			fw.PolicyList = []*Policy{p}
			break
		}
	}
	if !found {
		return nil, ErrNotFound
	}

	return fw, nil
}
