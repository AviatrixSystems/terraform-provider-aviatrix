package goaviatrix

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ProfileRule struct {
	Protocol string `form:"proto,omitempty" json:"protocol,omitempty"`
	Target   string `form:"target,omitempty" json:"target,omitempty"`
	Port     string `form:"port,omitempty" json:"port,omitempty"`
	Action   string `form:"action,omitempty" json:"action,omitempty"`
}

// Gateway simple struct to hold profile details
type Profile struct {
	Action   string        `form:"action,omitempty"`
	CID      string        `form:"CID,omitempty"`
	Name     string        `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	BaseRule string        `form:"base_rule,omitempty" json:"status,omitempty"`
	Policy   []ProfileRule `form:"domain_names[],omitempty" json:"domain_names,omitempty"`
	UserList []string      `form:"user_names,omitempty" json:"user_names,omitempty"`
}

type ProfilePolicyListRule struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type ProfilePolicyListResp struct {
	Return  bool          `json:"return"`
	Results []ProfileRule `json:"results"`
	Reason  string        `json:"reason"`
}

type ProfileUserListResp struct {
	Return  bool                `json:"return"`
	Results map[string][]string `json:"results"`
	Reason  string              `json:"reason"`
}

type ProfileBasePolicyResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateProfile(profile *Profile) error {
	form1 := map[string]string{
		"CID":          c.CID,
		"action":       "add_user_profile",
		"profile_name": profile.Name,
		"base_policy":  profile.BaseRule,
	}

	err := c.PostAPI(form1["action"], form1, BasicCheck)
	if err != nil {
		return err
	}

	policyStr, _ := json.Marshal(profile.Policy)
	form2 := map[string]string{
		"CID":          c.CID,
		"action":       "update_profile_policy",
		"profile_name": profile.Name,
		"policy":       string(policyStr),
	}

	err = c.PostAPI(form2["action"], form2, BasicCheck)
	if err != nil {
		return err
	}

	for _, user := range profile.UserList {
		form := map[string]string{
			"CID":          c.CID,
			"action":       "add_profile_member",
			"profile_name": profile.Name,
			"username":     user,
		}

		err = c.PostAPI(form["action"], form, BasicCheck)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) GetProfile(profile *Profile) (*Profile, error) {
	form1 := map[string]string{
		"CID":          c.CID,
		"action":       "list_profile_policies",
		"profile_name": profile.Name,
	}

	var data1 ProfilePolicyListResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data1, form1["action"], form1, checkFunc)
	if err != nil {
		return nil, err
	}

	profile.Policy = data1.Results
	log.Tracef("Profile policy %s", profile.Policy)

	form2 := map[string]string{
		"CID":    c.CID,
		"action": "list_user_profile_names",
	}

	var data2 ProfileUserListResp

	err = c.GetAPI(&data2, form2["action"], form2, BasicCheck)
	if err != nil {
		return nil, err
	}

	profile.UserList = data2.Results[profile.Name]
	log.Tracef("Profile list of users %s", profile.UserList)
	return profile, nil
}

func (c *Client) UpdateProfilePolicy(profile *Profile) error {
	log.Tracef("Updating Profile Policy %#v", profile)

	policyStr, _ := json.Marshal(profile.Policy)
	form := map[string]string{
		"CID":          c.CID,
		"action":       "update_profile_policy",
		"profile_name": profile.Name,
		"policy":       string(policyStr),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) AttachUsers(profile *Profile) error {
	log.Tracef("Attaching users %s", profile.UserList)

	for _, user := range profile.UserList {
		form := map[string]string{
			"CID":          c.CID,
			"action":       "add_profile_member",
			"profile_name": profile.Name,
			"username":     user,
		}

		err := c.PostAPI(form["action"], form, BasicCheck)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) DetachUsers(profile *Profile) error {
	log.Tracef("Detaching users %s", profile.UserList)

	for _, user := range profile.UserList {
		form := map[string]string{
			"CID":          c.CID,
			"action":       "del_profile_member",
			"profile_name": profile.Name,
			"username":     user,
		}

		err := c.PostAPI(form["action"], form, BasicCheck)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) DeleteProfile(profile *Profile) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "del_user_profile",
		"profile_name": profile.Name,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetProfileBasePolicy(profile *Profile) (*Profile, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "get_profile_base_policy",
		"profile_name": profile.Name,
	}

	var data ProfileBasePolicyResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if strings.Contains(data.Results, "allow all") {
		profile.BaseRule = "allow_all"
	} else if strings.Contains(data.Results, "deny all") {
		profile.BaseRule = "deny_all"
	}

	return profile, nil
}

func (c *Client) ValidateProfileRule(profileRule *ProfileRule) error {
	if profileRule.Action != "allow" && profileRule.Action != "deny" {
		return fmt.Errorf("valid action is only 'allow' or 'deny'")
	}
	protocolDefaultValues := []string{"all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp"}
	protocolVal := []string{profileRule.Protocol}
	if profileRule.Protocol == "" || len(Difference(protocolVal, protocolDefaultValues)) != 0 {
		return fmt.Errorf("proto can only be one of {'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'}")
	}
	if (profileRule.Protocol == "all" || profileRule.Protocol == "icmp") && (profileRule.Port != "0:65535") {
		return fmt.Errorf("port should be '0:65535' for protocal 'all' or 'icmp'")
	}
	return nil
}
