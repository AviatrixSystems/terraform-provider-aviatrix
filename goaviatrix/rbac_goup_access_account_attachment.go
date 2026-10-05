package goaviatrix

import (
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
