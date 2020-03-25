package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

type RbacGroupAccessAccountAttachment struct {
	CID               string `form:"CID,omitempty"`
	Action            string `form:"action,omitempty"`
	GroupName         string `form:"group_name,omitempty" json:"group_name,omitempty"`
	AccessAccountName string `form:"accounts,omitempty" json:"accounts,omitempty"`
}

type RbacGroupAccessAccountAttachmentListResp struct {
	Return                               bool     `json:"return"`
	RbacGroupAccessAccountAttachmentList []string `json:"results"`
	Reason                               string   `json:"reason"`
}

func (c *Client) CreateRbacGroupAccessAccountAttachment(rbacGroupAccessAccountAttachment *RbacGroupAccessAccountAttachment) error {
	rbacGroupAccessAccountAttachment.CID = c.CID
	rbacGroupAccessAccountAttachment.Action = "add_access_accounts_to_rbac_group"
	resp, err := c.Post(c.baseURL, rbacGroupAccessAccountAttachment)
	if err != nil {
		return errors.New("HTTP Post 'add_access_accounts_to_rbac_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_access_accounts_to_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_access_accounts_to_rbac_group' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetRbacGroupAccessAccountAttachment(rbacGroupAccessAccountAttachment *RbacGroupAccessAccountAttachment) (*RbacGroupAccessAccountAttachment, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_access_accounts_in_rbac_group': ") + err.Error())
	}
	listAccessAccountsInRbacGroup := url.Values{}
	listAccessAccountsInRbacGroup.Add("CID", c.CID)
	listAccessAccountsInRbacGroup.Add("action", "list_access_accounts_in_rbac_group")
	listAccessAccountsInRbacGroup.Add("group_name", rbacGroupAccessAccountAttachment.GroupName)
	Url.RawQuery = listAccessAccountsInRbacGroup.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_access_accounts_in_rbac_group' failed: " + err.Error())
	}
	var data RbacGroupAccessAccountAttachmentListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_access_accounts_in_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'list_access_accounts_in_rbac_group' Get failed: " + data.Reason)
	}
	attachments := data.RbacGroupAccessAccountAttachmentList
	for i := range attachments {
		if attachments[i] == rbacGroupAccessAccountAttachment.AccessAccountName {
			log.Printf("[INFO] Found Aviatrix RBAC group access account attachment: %s",
				rbacGroupAccessAccountAttachment.GroupName+"~"+rbacGroupAccessAccountAttachment.AccessAccountName)
			return rbacGroupAccessAccountAttachment, nil
		}
	}

	log.Printf("Couldn't find Aviatrix RBAC group access account attachment: %s",
		rbacGroupAccessAccountAttachment.GroupName+"~"+rbacGroupAccessAccountAttachment.AccessAccountName)
	return nil, ErrNotFound
}

func (c *Client) DeleteRbacGroupAccessAccountAttachment(rbacGroupAccessAccountAttachment *RbacGroupAccessAccountAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'delete_access_accounts_from_rbac_group': ") + err.Error())
	}
	deleteAccessAccountsFromRbacGroups := url.Values{}
	deleteAccessAccountsFromRbacGroups.Add("CID", c.CID)
	deleteAccessAccountsFromRbacGroups.Add("action", "delete_access_accounts_from_rbac_group")
	deleteAccessAccountsFromRbacGroups.Add("group_name", rbacGroupAccessAccountAttachment.GroupName)
	deleteAccessAccountsFromRbacGroups.Add("accounts", rbacGroupAccessAccountAttachment.AccessAccountName)
	Url.RawQuery = deleteAccessAccountsFromRbacGroups.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'delete_access_accounts_from_rbac_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_access_accounts_from_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_access_accounts_from_rbac_group' Get failed: " + data.Reason)
	}
	return nil
}
