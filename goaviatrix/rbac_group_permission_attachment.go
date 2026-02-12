package goaviatrix

import (
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

	return c.PostAPI(rbacGroupPermissionAttachment.Action, rbacGroupPermissionAttachment, BasicCheck)
}

func (c *Client) GetRbacGroupPermissionAttachment(rbacGroupPermissionAttachment *RbacGroupPermissionAttachment) (*RbacGroupPermissionAttachment, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "list_rbac_group_permissions",
		"group_name": rbacGroupPermissionAttachment.GroupName,
	}

	var data RbacGroupPermissionAttachmentListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
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
	form := map[string]string{
		"CID":         c.CID,
		"action":      "delete_permissions_from_rbac_group",
		"group_name":  rbacGroupPermissionAttachment.GroupName,
		"permissions": rbacGroupPermissionAttachment.PermissionName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
