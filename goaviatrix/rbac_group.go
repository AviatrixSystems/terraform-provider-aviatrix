package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type RbacGroup struct {
	CID       string `form:"CID,omitempty"`
	Action    string `form:"action,omitempty"`
	GroupName string `form:"group_name,omitempty" json:"group_name,omitempty"`
}

type RbacGroupListResp struct {
	Return        bool     `json:"return"`
	RbacGroupList []string `json:"results"`
	Reason        string   `json:"reason"`
}

func (c *Client) CreatePermissionGroup(rbacGroup *RbacGroup) error {
	rbacGroup.CID = c.CID
	rbacGroup.Action = "add_permission_group"
	resp, err := c.Post(c.baseURL, rbacGroup)
	if err != nil {
		return errors.New("HTTP Post 'add_permission_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_permission_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_permission_group' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetPermissionGroup(rbacGroup *RbacGroup) (*RbacGroup, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_permission_groups': ") + err.Error())
	}
	listPermissionGroups := url.Values{}
	listPermissionGroups.Add("CID", c.CID)
	listPermissionGroups.Add("action", "list_permission_groups")
	Url.RawQuery = listPermissionGroups.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_permission_groups' failed: " + err.Error())
	}
	var data RbacGroupListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_permission_groups' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'list_permission_groups' Get failed: " + data.Reason)
	}
	groups := data.RbacGroupList
	for i := range groups {
		if groups[i] == rbacGroup.GroupName {
			log.Infof("Found Aviatrix RBAC group: %s", rbacGroup.GroupName)
			return rbacGroup, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix RBAC group: %s", rbacGroup.GroupName)
	return nil, ErrNotFound
}

func (c *Client) DeletePermissionGroup(rbacGroup *RbacGroup) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'delete_permission_group': ") + err.Error())
	}
	deletePermissionGroups := url.Values{}
	deletePermissionGroups.Add("CID", c.CID)
	deletePermissionGroups.Add("action", "delete_permission_group")
	deletePermissionGroups.Add("group_name", rbacGroup.GroupName)
	Url.RawQuery = deletePermissionGroups.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'delete_permission_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_permission_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_permission_group' Get failed: " + data.Reason)
	}
	return nil
}
