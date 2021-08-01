package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
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

	return c.PostAPI(firewall_tag.Action, firewall_tag, BasicCheck)
}

func (c *Client) UpdateFirewallTag(firewall_tag *FirewallTag) error {
	// TODO: use PostAPI - tags need special processing
	firewall_tag.CID = c.CID
	firewall_tag.Action = "update_policy_members"
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&tag_name=%s", c.CID, firewall_tag.Action, firewall_tag.Name)
	for i, cidr := range firewall_tag.CIDRList {
		body = body + fmt.Sprintf("&new_policies[%d][name]=%s&new_policies[%d][cidr]=%s", i, cidr.CIDRTag, i, cidr.CIDR)
	}
	log.Tracef("%s %s Body: %s", verb, c.baseURL, body)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode update_policy_members failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API update_policy_members Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetFirewallTag(firewall_tag *FirewallTag) (*FirewallTag, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_policy_members",
		"tag_name": firewall_tag.Name,
	}

	var data FirewallTagResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			return ErrNotFound
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) DeleteFirewallTag(firewall_tag *FirewallTag) error {
	firewall_tag.CID = c.CID
	firewall_tag.Action = "del_policy_tag"

	return c.PostAPI(firewall_tag.Action, firewall_tag, BasicCheck)
}
