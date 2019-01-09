package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	resp, err := c.Post(c.baseURL, user)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetAccountUser(user *AccountUser) (*AccountUser, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_account_users", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data AccountUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if !data.Return {
		return nil, errors.New(data.Reason)
	}
	users := data.AccountUserList
	for i := range users {
		if users[i].UserName == user.UserName && users[i].AccountName == user.AccountName {
			log.Printf("[INFO] Found Aviatrix user account %s", user.UserName)
			return &users[i], nil
		}
	}
	log.Printf("Couldn't find Aviatrix user account %s", user.UserName)
	return nil, ErrNotFound

}

func (c *Client) UpdateAccountUserObject(user *AccountUserEdit) error {
	user.CID = c.CID
	user.Action = "edit_account_user"
	resp, err := c.Post(c.baseURL, user)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DeleteAccountUser(user *AccountUser) error {
	path := c.baseURL + fmt.Sprintf("?action=delete_account_user&CID=%s&username=%s", c.CID, user.UserName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
