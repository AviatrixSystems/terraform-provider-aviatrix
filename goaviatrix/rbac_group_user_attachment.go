package goaviatrix

import (
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
