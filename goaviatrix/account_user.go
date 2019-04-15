package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
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
		return errors.New("HTTP Post add_account_user failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_account_user failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_account_user Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAccountUser(user *AccountUser) (*AccountUser, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_account_users ") + err.Error())
	}
	listAccountUsers := url.Values{}
	listAccountUsers.Add("CID", c.CID)
	listAccountUsers.Add("action", "list_account_users")
	Url.RawQuery = listAccountUsers.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_account_users failed: " + err.Error())
	}
	var data AccountUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_account_users failed: " + err.Error())
	}

	if !data.Return {
		return nil, errors.New("Rest API enable_transit_ha Get failed: " + data.Reason)
	}
	users := data.AccountUserList
	for i := range users {
		if users[i].UserName == user.UserName {
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
		return errors.New("HTTP Post edit_account_user failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode edit_account_user failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API edit_account_user Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteAccountUser(user *AccountUser) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_account_user ") + err.Error())
	}
	deleteAccountUsers := url.Values{}
	deleteAccountUsers.Add("CID", c.CID)
	deleteAccountUsers.Add("action", "delete_account_user")
	deleteAccountUsers.Add("username", user.UserName)
	Url.RawQuery = deleteAccountUsers.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get delete_account_user failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_account_user failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_account_user Get failed: " + data.Reason)
	}
	return nil
}
