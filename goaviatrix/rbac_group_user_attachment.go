package goaviatrix

import (
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

	return c.PostAPI(rbacGroupUserAttachment.Action, rbacGroupUserAttachment, BasicCheck)
}

func (c *Client) GetRbacGroupUserAttachment(rbacGroupUserAttachment *RbacGroupUserAttachment) (*RbacGroupUserAttachment, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "list_users_in_rbac_group",
		"group_name": rbacGroupUserAttachment.GroupName,
	}

	var data RbacGroupUserAttachmentListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
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
	form := map[string]string{
		"CID":        c.CID,
		"action":     "delete_users_from_rbac_group",
		"group_name": rbacGroupUserAttachment.GroupName,
		"users":      rbacGroupUserAttachment.UserName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

// ListRbacGroupUsers retrieves a list of all users assigned to the specified
// RBAC group.
func (c *Client) ListRbacGroupUsers(groupName string) ([]string, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "list_users_in_rbac_group",
		"group_name": groupName,
	}
	var data RbacGroupUserAttachmentListResp
	if err := c.GetAPI(&data, form["action"], form, BasicCheck); err != nil {
		return nil, err
	}
	return data.RbacGroupUserAttachmentList, nil
}

// AddRbacGroupUsers adds one or more users to the specified RBAC group.
func (c *Client) AddRbacGroupUsers(groupName string, users []string) error {
	if len(users) == 0 {
		return nil
	}
	payload := &RbacGroupUserAttachment{
		CID:       c.CID,
		Action:    "add_users_to_rbac_group",
		GroupName: groupName,
		UserName:  strings.Join(users, ","), // API expects comma-separated in "users"
	}
	return c.PostAPI(payload.Action, payload, BasicCheck)
}

// DeleteRbacGroupUsers removes one or more users from the specified RBAC
// group. Takes a group name and a slice of usernames to remove
func (c *Client) DeleteRbacGroupUsers(groupName string, users []string) error {
	if len(users) == 0 {
		return nil
	}
	form := map[string]string{
		"CID":        c.CID,
		"action":     "delete_users_from_rbac_group",
		"group_name": groupName,
		"users":      strings.Join(users, ","),
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

// SetRbacGroupUsers sets the exact membership of an RBAC group by comparing
// current users with desired users, then adding missing users and removing
// extra users. Ensures the group contains exactly the specified users.
func (c *Client) SetRbacGroupUsers(groupName string, desired []string) error {
	current, err := c.ListRbacGroupUsers(groupName)
	if err != nil {
		return err
	}
	toAdd, toDel := diffStrings(current, desired)
	if err := c.AddRbacGroupUsers(groupName, toAdd); err != nil {
		return err
	}
	if err := c.DeleteRbacGroupUsers(groupName, toDel); err != nil {
		return err
	}
	return nil
}
