package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type RbacGroupPermissionAttachment struct {
	CID            string `form:"CID,omitempty"`
	Action         string `form:"action,omitempty"`
	GroupName      string `form:"group_name,omitempty" json:"group_name,omitempty"`
	PermissionName string `form:"permissions,omitempty" json:"permissions,omitempty"`
}

type RbacGroupPermissionAttachmentListResp struct {
	Return                            bool                       `json:"return"`
	RbacGroupPermissionAttachmentList []PermissionAttachmentInfo `json:"results"`
	Reason                            string                     `json:"reason"`
}

type PermissionAttachmentInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func (c *Client) CreateRbacGroupPermissionAttachment(rbacGroupPermissionAttachment *RbacGroupPermissionAttachment) error {
	rbacGroupPermissionAttachment.CID = c.CID
	rbacGroupPermissionAttachment.Action = "add_permissions_to_rbac_group"
	resp, err := c.Post(c.baseURL, rbacGroupPermissionAttachment)
	if err != nil {
		return errors.New("HTTP Post 'add_permissions_to_rbac_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_permissions_to_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_permissions_to_rbac_group' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetRbacGroupPermissionAttachment(rbacGroupPermissionAttachment *RbacGroupPermissionAttachment) (*RbacGroupPermissionAttachment, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_rbac_group_permissions': ") + err.Error())
	}
	listRbacGroupPermissions := url.Values{}
	listRbacGroupPermissions.Add("CID", c.CID)
	listRbacGroupPermissions.Add("action", "list_rbac_group_permissions")
	listRbacGroupPermissions.Add("group_name", rbacGroupPermissionAttachment.GroupName)
	Url.RawQuery = listRbacGroupPermissions.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_rbac_group_permissions' failed: " + err.Error())
	}
	var data RbacGroupPermissionAttachmentListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_rbac_group_permissions' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'list_rbac_group_permissions' Get failed: " + data.Reason)
	}
	attachments := data.RbacGroupPermissionAttachmentList
	for i := range attachments {
		if attachments[i].Name == rbacGroupPermissionAttachment.PermissionName {
			log.Infof("Found Aviatrix RBAC group permission attachment: %s",
				rbacGroupPermissionAttachment.GroupName+"~"+rbacGroupPermissionAttachment.PermissionName)
			return rbacGroupPermissionAttachment, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix RBAC group permission attachment: %s",
		rbacGroupPermissionAttachment.GroupName+"~"+rbacGroupPermissionAttachment.PermissionName)
	return nil, ErrNotFound
}

func (c *Client) DeleteRbacGroupPermissionAttachment(rbacGroupPermissionAttachment *RbacGroupPermissionAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'delete_permissions_from_rbac_group': ") + err.Error())
	}
	deletePermissionsFromRbacGroup := url.Values{}
	deletePermissionsFromRbacGroup.Add("CID", c.CID)
	deletePermissionsFromRbacGroup.Add("action", "delete_permissions_from_rbac_group")
	deletePermissionsFromRbacGroup.Add("group_name", rbacGroupPermissionAttachment.GroupName)
	deletePermissionsFromRbacGroup.Add("permissions", rbacGroupPermissionAttachment.PermissionName)
	Url.RawQuery = deletePermissionsFromRbacGroup.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'delete_permissions_from_rbac_group' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_permissions_from_rbac_group' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_permissions_from_rbac_group' Get failed: " + data.Reason)
	}
	return nil
}
