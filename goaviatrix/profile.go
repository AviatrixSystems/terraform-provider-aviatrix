package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for add_user_profile") + err.Error())
	}
	addUserProfile := url.Values{}
	addUserProfile.Add("CID", c.CID)
	addUserProfile.Add("action", "add_user_profile")
	addUserProfile.Add("profile_name", profile.Name)
	addUserProfile.Add("base_policy", profile.BaseRule)
	Url.RawQuery = addUserProfile.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get add_user_profile failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_user_profile failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_user_profile Get failed: " + data.Reason)
	}

	policyStr, _ := json.Marshal(profile.Policy)
	updateProfilePolicy := url.Values{}
	updateProfilePolicy.Add("CID", c.CID)
	updateProfilePolicy.Add("action", "update_profile_policy")
	updateProfilePolicy.Add("profile_name", profile.Name)
	updateProfilePolicy.Add("policy", string(policyStr))
	Url.RawQuery = updateProfilePolicy.Encode()

	log.Printf("[INFO] Creating Aviatrix Profile with Policy: %v", Url.String())

	resp, err = c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get update_profile_policy failed: " + err.Error())
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode update_profile_policy failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API update_profile_policy failed: " + data.Reason)
	}

	for _, user := range profile.UserList {
		addProfileMember := url.Values{}
		addProfileMember.Add("CID", c.CID)
		addProfileMember.Add("action", "add_profile_member")
		addProfileMember.Add("profile_name", profile.Name)
		addProfileMember.Add("username", user)
		Url.RawQuery = addProfileMember.Encode()

		resp, err = c.Get(Url.String(), nil)
		if err != nil {
			return errors.New("HTTP Get add_profile_member failed: " + err.Error())
		}
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.New("Json Decode add_profile_member failed: " + err.Error())
		}
		if !data.Return {
			return errors.New("API Get add_profile_member failed: " + data.Reason)
		}
	}
	return nil
}

func (c *Client) GetProfile(profile *Profile) (*Profile, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_profile_policies") + err.Error())
	}
	listProfilePolicy := url.Values{}
	listProfilePolicy.Add("CID", c.CID)
	listProfilePolicy.Add("action", "list_profile_policies")
	listProfilePolicy.Add("profile_name", profile.Name)
	Url.RawQuery = listProfilePolicy.Encode()

	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_profile_policies failed: " + err.Error())
	}
	var data ProfilePolicyListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_profile_policies failed: " + err.Error())
	}

	if !data.Return {
		log.Printf("Couldn't find Aviatrix profile %s", profile.Name)
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, errors.New("API Get list_profile_policies failed: " + data.Reason)
	}
	profile.Policy = data.Results
	log.Printf("[TRACE] Profile policy %s", profile.Policy)

	listUserProfileNames := url.Values{}
	listUserProfileNames.Add("CID", c.CID)
	listUserProfileNames.Add("action", "list_user_profile_names")
	Url.RawQuery = listUserProfileNames.Encode()

	resp, err = c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_user_profile_names failed: " + err.Error())
	}
	var data2 ProfileUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data2); err != nil {
		return nil, errors.New("Json Decode list_user_profile_names failed: " + err.Error())
	}

	//profile.BaseRule = data2.Results[profile.Name]
	profile.UserList = data2.Results[profile.Name]

	log.Printf("[TRACE] Profile list of users %s", profile.UserList)

	return profile, nil
}

func (c *Client) UpdateProfilePolicy(profile *Profile) error {
	log.Printf("[TRACE] Updating Profile Policy %#v", profile)

	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for update_profile_policy") + err.Error())
	}
	policyStr, _ := json.Marshal(profile.Policy)
	updateProfilePolicy := url.Values{}
	updateProfilePolicy.Add("CID", c.CID)
	updateProfilePolicy.Add("action", "update_profile_policy")
	updateProfilePolicy.Add("profile_name", profile.Name)
	updateProfilePolicy.Add("policy", string(policyStr))
	Url.RawQuery = updateProfilePolicy.Encode()

	var data APIResp
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get update_profile_policy failed: " + err.Error())
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode update_profile_policy failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API update_profile_policy failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachUsers(profile *Profile) error {
	log.Printf("[TRACE] Attaching users %s", profile.UserList)
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		panic("boom")
	}
	for _, user := range profile.UserList {
		addProfileMember := url.Values{}
		addProfileMember.Add("CID", c.CID)
		addProfileMember.Add("action", "add_profile_member")
		addProfileMember.Add("profile_name", profile.Name)
		addProfileMember.Add("username", user)
		Url.RawQuery = addProfileMember.Encode()

		var data APIResp
		resp, err := c.Get(Url.String(), nil)
		if err != nil {
			return errors.New("HTTP Get add_profile_member failed: " + err.Error())
		}
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.New("Json Decode add_profile_member failed: " + err.Error())
		}
		if !data.Return {
			return errors.New("API Get add_profile_member failed: " + data.Reason)
		}
	}
	return nil
}

func (c *Client) DetachUsers(profile *Profile) error {
	log.Printf("[TRACE] Detaching users %s", profile.UserList)
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for del_profile_member") + err.Error())
	}
	for _, user := range profile.UserList {
		delProfileMember := url.Values{}
		delProfileMember.Add("CID", c.CID)
		delProfileMember.Add("action", "del_profile_member")
		delProfileMember.Add("profile_name", profile.Name)
		delProfileMember.Add("username", user)
		Url.RawQuery = delProfileMember.Encode()

		var data APIResp
		resp, err := c.Get(Url.String(), nil)
		if err != nil {
			return errors.New("HTTP Get del_profile_member failed: " + err.Error())
		}
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.New("Json Decode del_profile_member failed: " + err.Error())
		}
		if !data.Return {
			return errors.New("API Get del_profile_member failed: " + data.Reason)
		}
	}

	return nil
}

func (c *Client) DeleteProfile(profile *Profile) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for del_user_profile") + err.Error())
	}
	delUserProfile := url.Values{}
	delUserProfile.Add("CID", c.CID)
	delUserProfile.Add("action", "del_user_profile")
	delUserProfile.Add("profile_name", profile.Name)
	Url.RawQuery = delUserProfile.Encode()

	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get del_user_profile failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode del_user_profile failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("API Get del_user_profile failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetProfileBasePolicy(profile *Profile) (*Profile, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for get_profile_base_policy") + err.Error())
	}
	getProfileBasePolicy := url.Values{}
	getProfileBasePolicy.Add("CID", c.CID)
	getProfileBasePolicy.Add("action", "get_profile_base_policy")
	getProfileBasePolicy.Add("profile_name", profile.Name)
	Url.RawQuery = getProfileBasePolicy.Encode()

	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get get_profile_base_policy failed: " + err.Error())
	}
	var data ProfileBasePolicyResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_profile_base_policy failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("API Get get_profile_base_policy failed: " + data.Reason)
	} else {
		if strings.Contains(data.Results, "allow all") {
			profile.BaseRule = "allow_all"
		} else if strings.Contains(data.Results, "deny all") {
			profile.BaseRule = "deny_all"
		}
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
