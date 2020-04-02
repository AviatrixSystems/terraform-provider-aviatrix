package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type RbacGroupUserAttachment struct {
	CID       string `form:"CID,omitempty"`
	Action    string `form:"action,omitempty"`
	GroupName string `form:"group_name,omitempty" json:"group_name,omitempty"`
	UserName  string `form:"users,omitempty" json:"users,omitempty"`
}

type RbacGroupUserAttachmentListResp struct {
	Return                      bool     `json:"return"`
	RbacGroupUserAttachmentList []string `json:"results"`
	Reason                      string   `json:"reason"`
}

func (c *Client) CreateRbacGroupUserAttachment(rbacGroupUserAttachment *RbacGroupUserAttachment) error {
	rbacGroupUserAttachment.CID = c.CID
	rbacGroupUserAttachment.Action = "add_users_to_rbac_group"
	resp, err := c.Post(c.baseURL, rbacGroupUserAttachment)
	if err != nil {
		return errors.New("HTTP Post 'add_users_to_rbac_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_users_to_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_users_to_rbac_group' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetRbacGroupUserAttachment(rbacGroupUserAttachment *RbacGroupUserAttachment) (*RbacGroupUserAttachment, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_users_in_rbac_group': ") + err.Error())
	}
	listUsersInRbacGroup := url.Values{}
	listUsersInRbacGroup.Add("CID", c.CID)
	listUsersInRbacGroup.Add("action", "list_users_in_rbac_group")
	listUsersInRbacGroup.Add("group_name", rbacGroupUserAttachment.GroupName)
	Url.RawQuery = listUsersInRbacGroup.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_users_in_rbac_group' failed: " + err.Error())
	}
	var data RbacGroupUserAttachmentListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_users_in_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'list_users_in_rbac_group' Get failed: " + data.Reason)
	}
	attachments := data.RbacGroupUserAttachmentList
	for i := range attachments {
		if attachments[i] == rbacGroupUserAttachment.UserName {
			log.Infof("Found Aviatrix RBAC group user attachment: %s",
				rbacGroupUserAttachment.GroupName+"~"+rbacGroupUserAttachment.UserName)
			return rbacGroupUserAttachment, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix RBAC group user attachment: %s",
		rbacGroupUserAttachment.GroupName+"~"+rbacGroupUserAttachment.UserName)
	return nil, ErrNotFound
}

func (c *Client) DeleteRbacGroupUserAttachment(rbacGroupUserAttachment *RbacGroupUserAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'delete_users_from_rbac_group': ") + err.Error())
	}
	deleteUsersFromRbacGroups := url.Values{}
	deleteUsersFromRbacGroups.Add("CID", c.CID)
	deleteUsersFromRbacGroups.Add("action", "delete_users_from_rbac_group")
	deleteUsersFromRbacGroups.Add("group_name", rbacGroupUserAttachment.GroupName)
	deleteUsersFromRbacGroups.Add("users", rbacGroupUserAttachment.UserName)
	Url.RawQuery = deleteUsersFromRbacGroups.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'delete_users_from_rbac_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_users_from_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_users_from_rbac_group' Get failed: " + data.Reason)
	}
	return nil
}
