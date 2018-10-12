package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type ProfileRule struct {
        Protocol    string `form:"proto,omitempty" json:"protocol,omitempty"`
		Target      string `form:"target,omitempty" json:"target,omitempty"`
        Port        string `form:"port,omitempty" json:"port,omitempty"`
		Action      string `form:"action,omitempty" json:"action,omitempty"`
}

// Gateway simple struct to hold profile details
type Profile struct {
	Action                  string `form:"action,omitempty"`
	CID                     string `form:"CID,omitempty"`
	Name                 	string `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	BaseRule              	string `form:"base_rule,omitempty" json:"status,omitempty"`
	Policy              	[]ProfileRule `form:"domain_names[],omitempty" json:"domain_names,omitempty"`
	UserList                []string `form:"user_names,omitempty" json:"user_names,omitempty"`
}


type ProfilePolicyListRule struct {
	Return  bool   `json:"return"`
	Results []string `json:"results"`
	Reason  string `json:"reason"`
}

type ProfilePolicyListResp struct {
	Return  bool   `json:"return"`
	Results []ProfileRule `json:"results"`
	Reason  string `json:"reason"`
}


type ProfileUserListResp struct {
	Return  bool   `json:"return"`
	Results map[string][]string`json:"results"`
	Reason  string `json:"reason"`
}


func (c *Client) CreateProfile(profile *Profile) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=add_user_profile&profile_name=%s&base_policy=%s",
		c.CID, profile.Name, profile.BaseRule)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return{
		return errors.New(data.Reason)
	}
	policyStr, _ := json.Marshal(profile.Policy)
	path = c.baseURL + fmt.Sprintf("?CID=%s&action=update_profile_policy&profile_name=%s&policy=%s",
		c.CID, profile.Name, policyStr)
	log.Printf("[INFO] Creating Aviatrix Profile with Policy: %v", path)
	resp,err = c.Get(path, nil)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return{
		return errors.New(data.Reason)
	}

	path = c.baseURL + fmt.Sprintf("?CID=%s&action=update_profile_policy&profile_name=%s&policy=%s",
		c.CID, profile.Name, policyStr)
	resp,err = c.Get(path, nil)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return{
		return errors.New(data.Reason)
	}
		for _, user := range profile.UserList {
		path = c.baseURL + fmt.Sprintf("?CID=%s&action=add_profile_member&profile_name=%s&username=%s",
			c.CID, profile.Name, user)
		resp, err = c.Get(path, nil)
		if err != nil {
			return err
		}

		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if !data.Return {
			return errors.New(data.Reason)
		}
	}
	return nil
}


func (c *Client) GetProfile(profile *Profile) (*Profile, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_profile_policies&profile_name=%s", c.CID, profile.Name)
	resp,err := c.Get(path, nil)

	if err != nil {
		return nil, err
	}
	var data ProfilePolicyListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if !data.Return {
		return nil, errors.New(data.Reason)
		log.Printf("Couldn't find Aviatrix profile %s", profile.Name)
	}
	profile.Policy = data.Results
	log.Printf("[TRACE] Profile policy %s", profile.Policy)

	path = c.baseURL + fmt.Sprintf("?CID=%s&action=list_user_profile_names", c.CID)
	resp,err = c.Get(path, nil)

	if err != nil {
		return nil, err
	}
	var data2 ProfileUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data2); err != nil {
		return nil, err
	}

	profile.UserList = data2.Results[profile.Name]

	log.Printf("[TRACE] Profile list of users %s", profile.UserList)


	return profile, nil
}



func (c *Client) UpdateProfilePolicy(profile *Profile) (error) {
	log.Printf("[TRACE] Updating Profile Policy %#v",profile)
	policyStr, _ := json.Marshal(profile.Policy)
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=update_profile_policy&profile_name=%s&policy=%s",
		c.CID, profile.Name, policyStr)

	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return{
		return errors.New(data.Reason)
	}
	return nil

}


func (c *Client) AttachUsers(profile *Profile) (error) {
	log.Printf("[TRACE] Attaching users %s",profile.UserList)
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=add_profile_member&profile_name=%s", c.CID, profile.Name)
	for i := range profile.UserList {
		newPath := path + fmt.Sprintf("&username=%s", profile.UserList[i])
		resp,err := c.Get(newPath, nil)
		if err != nil {
			return err
		}
		var data APIResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if !data.Return{
			return errors.New(data.Reason)
		}
	}
	return nil
}

func (c *Client) DetachUsers(profile *Profile) (error) {
	log.Printf("[TRACE] Attaching users %s",profile.UserList)
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=del_profile_member&profile_name=%s", c.CID, profile.Name)
	for i := range profile.UserList {
		newPath := path + fmt.Sprintf("&username=%s", profile.UserList[i])
		resp,err := c.Get(newPath, nil)
		if err != nil {
			return err
		}
		var data APIResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if !data.Return{
			return errors.New(data.Reason)
		}
	}
	return nil
}


func (c *Client) DeleteProfile(profile *Profile) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=del_user_profile&profile_name=%s", c.CID, profile.Name)
	resp,err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return{
		return errors.New(data.Reason)
	}
	return nil
}