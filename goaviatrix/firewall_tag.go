package goaviatrix

import (
	"encoding/json"
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
	action := "update_policy_members"
	form := map[string]string{
		"CID":      c.CID,
		"action":   action,
		"tag_name": firewall_tag.Name,
	}
	if len(firewall_tag.CIDRList) == 0 {
		firewall_tag.CIDRList = []CIDRMember{}
	}

	args, err := json.Marshal(firewall_tag.CIDRList)
	if err != nil {
		return err
	}
	form["new_policies"] = string(args)

	return c.PostAPI(action, form, BasicCheck)
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
