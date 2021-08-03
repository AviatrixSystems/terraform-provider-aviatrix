package goaviatrix

import (
	log "github.com/sirupsen/logrus"
)

type AccountUser struct {
	CID         string `form:"CID,omitempty"`
	Action      string `form:"action,omitempty"`
	UserName    string `form:"username,omitempty" json:"user_name,omitempty"`
	AccountName string `form:"account_name,omitempty" json:"acct_names,omitempty"`
	Email       string `form:"email,omitempty" json:"user_email,omitempty"`
	Password    string `form:"password,omitempty" json:"password,omitempty"`
}

type AccountUserEdit struct {
	CID         string `form:"CID,omitempty"`
	Action      string `form:"action,omitempty"`
	UserName    string `form:"username,omitempty" json:"user_name,omitempty"`
	AccountName string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Email       string `form:"email,omitempty" json:"email,omitempty"`
	What        string `form:"what,omitempty" json:"what,omitempty"`
	OldPassword string `form:"old_password,omitempty" json:"old_password,omitempty"`
	NewPassword string `form:"new_password,omitempty" json:"new_password,omitempty"`
}

type AccountUserListResp struct {
	Return          bool          `json:"return"`
	AccountUserList []AccountUser `json:"results"`
	Reason          string        `json:"reason"`
}

func (c *Client) CreateAccountUser(user *AccountUser) error {
	user.CID = c.CID
	user.Action = "add_account_user"
	return c.PostAPI(user.Action, user, BasicCheck)
}

func (c *Client) GetAccountUser(user *AccountUser) (*AccountUser, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_account_users",
	}
	var data AccountUserListResp
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	users := data.AccountUserList
	for i := range users {
		if users[i].UserName == user.UserName {
			log.Infof("Found Aviatrix user account %s", user.UserName)
			return &users[i], nil
		}
	}
	log.Errorf("Couldn't find Aviatrix user account %s", user.UserName)
	return nil, ErrNotFound
}

func (c *Client) UpdateAccountUserObject(user *AccountUserEdit) error {
	user.CID = c.CID
	user.Action = "edit_account_user"
	return c.PostAPI(user.Action, user, BasicCheck)
}

func (c *Client) DeleteAccountUser(user *AccountUser) error {
	user.CID = c.CID
	user.Action = "delete_account_user"
	return c.PostAPI(user.Action, user, BasicCheck)
}
