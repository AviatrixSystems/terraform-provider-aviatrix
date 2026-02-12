package goaviatrix

import (
	"strings"

	log "github.com/sirupsen/logrus"
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

	return c.PostAPI(rbacGroupAccessAccountAttachment.Action, rbacGroupAccessAccountAttachment, BasicCheck)
}

func (c *Client) GetRbacGroupAccessAccountAttachment(rbacGroupAccessAccountAttachment *RbacGroupAccessAccountAttachment) (*RbacGroupAccessAccountAttachment, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "list_access_accounts_in_rbac_group",
		"group_name": rbacGroupAccessAccountAttachment.GroupName,
	}

	var data RbacGroupAccessAccountAttachmentListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	attachments := data.RbacGroupAccessAccountAttachmentList
	for i := range attachments {
		if attachments[i] == rbacGroupAccessAccountAttachment.AccessAccountName {
			log.Infof("Found Aviatrix RBAC group access account attachment: %s",
				rbacGroupAccessAccountAttachment.GroupName+"~"+rbacGroupAccessAccountAttachment.AccessAccountName)
			return rbacGroupAccessAccountAttachment, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix RBAC group access account attachment: %s",
		rbacGroupAccessAccountAttachment.GroupName+"~"+rbacGroupAccessAccountAttachment.AccessAccountName)
	return nil, ErrNotFound
}

func (c *Client) DeleteRbacGroupAccessAccountAttachment(rbacGroupAccessAccountAttachment *RbacGroupAccessAccountAttachment) error {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "delete_access_accounts_from_rbac_group",
		"group_name": rbacGroupAccessAccountAttachment.GroupName,
		"accounts":   rbacGroupAccessAccountAttachment.AccessAccountName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

// ListRbacGroupAccessAccounts retrieves a list of all access accounts assigned to the specified
// RBAC group.
func (c *Client) ListRbacGroupAccessAccounts(groupName string) ([]string, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "list_access_accounts_in_rbac_group",
		"group_name": groupName,
	}
	var data RbacGroupAccessAccountAttachmentListResp
	if err := c.GetAPI(&data, form["action"], form, BasicCheck); err != nil {
		return nil, err
	}
	return data.RbacGroupAccessAccountAttachmentList, nil
}

// AddRbacGroupAccessAccounts adds one or more access accounts to the specified RBAC group.
func (c *Client) AddRbacGroupAccessAccounts(groupName string, accessAccounts []string) error {
	if len(accessAccounts) == 0 {
		return nil
	}
	payload := &RbacGroupAccessAccountAttachment{
		CID:               c.CID,
		Action:            "add_access_accounts_to_rbac_group",
		GroupName:         groupName,
		AccessAccountName: strings.Join(accessAccounts, ","), // API expects comma-separated in "accounts"
	}
	return c.PostAPI(payload.Action, payload, BasicCheck)
}

// DeleteRbacGroupAccessAccounts removes one or more access accounts from the specified RBAC
// group. Takes a group name and a slice of accessAccounts to remove
func (c *Client) DeleteRbacGroupAccessAccounts(groupName string, accessAccounts []string) error {
	if len(accessAccounts) == 0 {
		return nil
	}
	form := map[string]string{
		"CID":        c.CID,
		"action":     "delete_access_accounts_from_rbac_group",
		"group_name": groupName,
		"accounts":   strings.Join(accessAccounts, ","),
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

// SetRbacGroupAccessAccounts sets the exact membership of an RBAC group by comparing
// current access accounts with desired access accounts, then adding missing access accounts and removing
// extra access accounts. Ensures the group contains exactly the specified access accounts.
func (c *Client) SetRbacGroupAccessAccounts(groupName string, desired []string) error {
	current, err := c.ListRbacGroupAccessAccounts(groupName)
	if err != nil {
		return err
	}
	toAdd, toDel := diffStrings(current, desired)
	if err := c.AddRbacGroupAccessAccounts(groupName, toAdd); err != nil {
		return err
	}
	if err := c.DeleteRbacGroupAccessAccounts(groupName, toDel); err != nil {
		return err
	}
	return nil
}
