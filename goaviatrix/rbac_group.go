package goaviatrix

import (
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

type RbacGroupResponse struct {
	LocalLogin bool   `json:"local_login"`
	GroupName  string `json:"name"`
}

type RbacGroupListDetailsResp struct {
	Return        bool                `json:"return"`
	RbacGroupList []RbacGroupResponse `json:"results"`
	Reason        string              `json:"reason"`
}

func (c *Client) CreatePermissionGroup(rbacGroup *RbacGroup) error {
	rbacGroup.CID = c.CID
	rbacGroup.Action = "add_permission_group"

	return c.PostAPI(rbacGroup.Action, rbacGroup, BasicCheck)
}

func (c *Client) GetPermissionGroup(rbacGroup *RbacGroup) (*RbacGroup, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_permission_groups",
	}

	var data RbacGroupListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
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
	form := map[string]string{
		"CID":        c.CID,
		"action":     "delete_permission_group",
		"group_name": rbacGroup.GroupName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableLocalLoginForRBACGroup(GroupName string) error {
	data := map[string]string{
		"action":     "enable_local_login",
		"CID":        c.CID,
		"group_name": GroupName,
	}
	return c.PostAPI("disable_local_login", data, BasicCheck)
}

func (c *Client) DisableLocalLoginForRBACGroup(GroupName string) error {
	data := map[string]string{
		"action":     "disable_local_login",
		"CID":        c.CID,
		"group_name": GroupName,
	}
	return c.PostAPI("disable_local_login", data, BasicCheck)
}

func (c *Client) GetPermissionGroupDetails(GroupName string) (*RbacGroupResponse, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_permission_group_details",
	}

	var data RbacGroupListDetailsResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	groups := data.RbacGroupList
	for i := range groups {
		if groups[i].GroupName == GroupName {
			log.Infof("Found Aviatrix RBAC group: %s", GroupName)
			return &groups[i], nil
		}
	}
	log.Errorf("Couldn't find Aviatrix RBAC group: %s", GroupName)
	return nil, ErrNotFound
}
