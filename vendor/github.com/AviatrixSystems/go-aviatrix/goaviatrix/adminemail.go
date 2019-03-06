package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type AdminEmailRequest struct {
	APIRequest
	Email string `form:"admin_email" url:"admin_email"`
}
type LoginProcRequest struct {
	Action   string `form:"action" url:"action"`
	Username string `form:"username" url:"username"`
	Password string `form:"password" url:"password"`
}
type AdminEmailResponse struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

type LoginProcResponse struct {
	AdminEmail   string `json:"admin_email"`
	InitialSetup bool   `json:"initial_setup"`
}

func (c *Client) SetAdminEmail(adminEmail string) error {
	log.Printf("[TRACE] Setting admin email to '%s'", adminEmail)
	admin := new(AdminEmailRequest)
	admin.Email = adminEmail
	admin.Action = "add_admin_email_addr"
	admin.CID = c.CID
	_, _, err := c.Do("GET", admin)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetAdminEmail(username string, password string) (string, error) {
	log.Printf("[TRACE] Getting admin email")
	path := fmt.Sprintf("https://%s/v1/backend1", c.ControllerIP)
	admin := new(LoginProcRequest)
	admin.Action = "login_proc"
	admin.Username = username
	admin.Password = password
	resp, err := c.Post(path, admin)
	if err != nil {
		return "", errors.New("HTTP Post login_proc failed: " + err.Error())
	}
	var data LoginProcResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", errors.New("Json Decode login_proc failed: " + err.Error())
	}
	return data.AdminEmail, nil
}
