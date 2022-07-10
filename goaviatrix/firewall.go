package goaviatrix

import (
	"encoding/json"
	"fmt"
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
	GwName         string    `form:"vpc_name,omitempty" json:"vpc_name,omitempty"`
	BasePolicy     string    `form:"base_policy,omitempty" json:"base_policy,omitempty"`
	BaseLogEnabled string    `form:"base_policy_log_enable,omitempty" json:"base_policy_log_enable,omitempty"`
	PolicyList     []*Policy `json:"security_rules,omitempty"`
	NewPolicy      string    `form:"new_policy,omitempty"`
}

type FirewallResp struct {
	Return  bool     `json:"return"`
	Results Firewall `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) SetBasePolicy(firewall *Firewall) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "set_vpc_base_policy",
		"vpc_name":               firewall.GwName,
		"base_policy":            firewall.BasePolicy,
		"base_policy_log_enable": firewall.BaseLogEnabled,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
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

	return c.PostAPI(firewall.Action, firewall, BasicCheck)
}

func (c *Client) GetPolicy(firewall *Firewall) (*Firewall, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "vpc_access_policy",
		"vpc_name": firewall.GwName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				log.Errorf("Couldn't find Aviatrix Firewall policies for gateway %s: %s", firewall.GwName, reason)
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}
	var data FirewallResp
	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
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
	rules, err := json.Marshal(fw.PolicyList)
	if err != nil {
		return fmt.Errorf("could not marshal firewall policies: %v", err)
	}

	form := map[string]string{
		"CID":          c.CID,
		"action":       "append_stateful_firewall_rules",
		"gateway_name": fw.GwName,
		"rules":        string(rules),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DeleteFirewallPolicy(fw *Firewall) error {
	rules, err := json.Marshal(fw.PolicyList)
	if err != nil {
		return fmt.Errorf("could not marshal firewall policies: %v", err)
	}

	form := map[string]string{
		"CID":          c.CID,
		"action":       "delete_stateful_firewall_rules",
		"gateway_name": fw.GwName,
		"rules":        string(rules),
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			// Tried to delete a rule that did not exist, we don't need to fail the apply.
			if strings.Contains(reason, "Empty rules in Stateful Firewall to delete") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) GetFirewallPolicy(fw *Firewall) (*Firewall, error) {
	foundFirewall, err := c.GetPolicy(fw)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
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
