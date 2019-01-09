package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CIDRMember struct {
	CIDRTag string `form:"name,omitempty" json:"name,omitempty"`
	CIDR    string `form:"cidr,omitempty" json:"cidr,omitempty"`
}

// Gateway simple struct to hold firewall_tag details
type FirewallTag struct {
	CID      string       `form:"CID,omitempty"`
	Action   string       `form:"action,omitempty"`
	Name     string       `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	CIDRList []CIDRMember `form:"new_policies,omitempty" json:"members,omitempty"`
}

type FirewallTagResp struct {
	Return  bool        `json:"return"`
	Results FirewallTag `json:"results"`
	Reason  string      `json:"reason"`
}

func (c *Client) CreateFirewallTag(firewall_tag *FirewallTag) error {
	firewall_tag.CID = c.CID
	firewall_tag.Action = "add_policy_tag"
	log.Printf("[INFO] Setting Firewall Tag: %#v", firewall_tag)
	resp, err := c.Post(c.baseURL, firewall_tag)
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

func (c *Client) UpdateFirewallTag(firewall_tag *FirewallTag) error {
	firewall_tag.CID = c.CID
	firewall_tag.Action = "update_policy_members"
	// It is not easy to marshal nested struct to create POST body
	// in proper format. On trying to create body using ajg/form package in client.go
	// in usual way, it creates body like this
	// CID=CrQ6FRPrv1FgGIv139gT&action=update_policy_members&new_policies.0.cidr=10.0.0.0/24&new_policies.0.name=a1&new_policies.1.cidr=10.1.0.0/24&new_policies.1.name=b1&tag_name=ranjan3
	// which is not understandable by our controller
	// Controller expects body key/value in array format like this
	// CID=CrQ6FRPrv1FgGIv139gT&action=update_policy_members&new_policies[0][cidr]=10.0.0.0/24&new_policies[0][name]=a1&new_policies[1][cidr]=10.1.0.0/24&new_policies[1][name]=b1&tag_name=ranjan3
	// So we are constructing body(and also sending request) here itself without calling Client.Post in client.go.
	// See this for more details https://stackoverflow.com/questions/48735329/golang-form-encode-nested-struct
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&tag_name=%s", c.CID, firewall_tag.Action, firewall_tag.Name)
	for i, cidr := range firewall_tag.CIDRList {
		body = body + fmt.Sprintf("&new_policies[%d][name]=%s&new_policies[%d][cidr]=%s", i, cidr.CIDRTag, i, cidr.CIDR)
	}
	log.Printf("[TRACE] %s %s Body: %s", verb, c.baseURL, body)
	req, err := http.NewRequest(verb, c.baseURL, strings.NewReader(body))
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
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

func (c *Client) GetFirewallTag(firewall_tag *FirewallTag) (*FirewallTag, error) {
	firewall_tag.CID = c.CID
	firewall_tag.Action = "list_policy_members"

	log.Printf("[INFO] Getting Firewall Tag: %#v", firewall_tag)
	resp, err := c.Post(c.baseURL, firewall_tag)
	if err != nil {
		return nil, err
	}
	var data FirewallTagResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find Aviatrix Firewall tag %s: %s", firewall_tag.Name, data.Reason)
		return nil, ErrNotFound
	}
	return &data.Results, nil
}

func (c *Client) DeleteFirewallTag(firewall_tag *FirewallTag) error {
	firewall_tag.CID = c.CID
	firewall_tag.Action = "del_policy_tag"
	log.Printf("[INFO] Deleting Firewall Tag: %#v", firewall_tag)
	resp, err := c.Post(c.baseURL, firewall_tag)
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
