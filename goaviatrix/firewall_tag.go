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
		return errors.New("HTTP Post update_policy_members failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode update_policy_members failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API update_policy_members Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetFirewallTag(firewall_tag *FirewallTag) (*FirewallTag, error) {
	firewall_tag.CID = c.CID
	firewall_tag.Action = "list_policy_members"

	log.Printf("[INFO] Getting Firewall Tag: %#v", firewall_tag)
	resp, err := c.Post(c.baseURL, firewall_tag)
	if err != nil {
		return nil, errors.New("HTTP Post list_policy_members failed: " + err.Error())
	}
	var data FirewallTagResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_policy_members failed: " + err.Error())
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
		return errors.New("HTTP Post del_policy_tag failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode del_policy_tag failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API del_policy_tag Post failed: " + data.Reason)
	}
	return nil
}
