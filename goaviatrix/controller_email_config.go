package goaviatrix

import (
	"context"
	"strconv"
)

type EmailConfiguration struct {
	CID                              string `form:"CID,omitempty"`
	Action                           string `form:"action,omitempty"`
	AdminAlertEmail                  string `form:"admin_alert_email,omitempty"`
	CriticalAlertEmail               string `form:"critical_alert_email,omitempty"`
	SecurityEventEmail               string `form:"security_event_email,omitempty"`
	StatusChangeEmail                string `form:"status_change_email,omitempty"`
	StatusChangeNotificationInterval int    `form:"status_change_notification_interval"`
	AdminAlertEmailVerified          bool
	CriticalAlertEmailVerified       bool
	SecurityEventEmailVerified       bool
	StatusChangeEmailVerified        bool
}

type ListAdminEmailAddrResp struct {
	AdminAlertEmail    EmailCell `json:"admin_alert"`
	CriticalAlertEmail EmailCell `json:"critical_alert"`
	SecurityEventEmail EmailCell `json:"security_event"`
	StatusChangeEmail  EmailCell `json:"status_change"`
}

type EmailCell struct {
	Address  string `json:"address"`
	Verified bool   `json:"verified"`
}

type ControllerEmailConfigResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) ConfigNotificationEmails(ctx context.Context, emailConfiguration *EmailConfiguration) error {
	form := map[string]interface{}{
		"CID":    c.CID,
		"action": "add_notif_email_addr",
		"notif_email_args": map[string]interface{}{
			"admin_alert": map[string]interface{}{
				"address": emailConfiguration.AdminAlertEmail,
			},
			"critical_alert": map[string]interface{}{
				"address": emailConfiguration.CriticalAlertEmail,
			},
			"security_event": map[string]interface{}{
				"address": emailConfiguration.SecurityEventEmail,
			},
			"status_change": map[string]interface{}{
				"address": emailConfiguration.StatusChangeEmail,
			},
		},
	}

	return c.PostAPIContext2(ctx, nil, form["action"].(string), form, BasicCheck)
}

func (c *Client) GetNotificationEmails(ctx context.Context) (*EmailConfiguration, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_notif_email_addr",
	}

	type Resp struct {
		Return  bool                   `json:"return"`
		Results ListAdminEmailAddrResp `json:"results"`
		Reason  string                 `json:"reason"`
	}

	var data Resp
	err := c.GetAPIContext(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	form1 := map[string]string{
		"CID":    c.CID,
		"action": "get_rate_limit_emails",
	}

	var data1 struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	err = c.GetAPIContext(ctx, &data1, form1["action"], form1, BasicCheck)
	if err != nil {
		return nil, err
	}

	interval, err := strconv.Atoi(data1.Results)
	if err != nil {
		return nil, err
	}
	return &EmailConfiguration{
		AdminAlertEmail:                  data.Results.AdminAlertEmail.Address,
		CriticalAlertEmail:               data.Results.CriticalAlertEmail.Address,
		SecurityEventEmail:               data.Results.SecurityEventEmail.Address,
		StatusChangeEmail:                data.Results.StatusChangeEmail.Address,
		StatusChangeNotificationInterval: interval,
		AdminAlertEmailVerified:          data.Results.AdminAlertEmail.Verified,
		CriticalAlertEmailVerified:       data.Results.CriticalAlertEmail.Verified,
		SecurityEventEmailVerified:       data.Results.SecurityEventEmail.Verified,
		StatusChangeEmailVerified:        data.Results.StatusChangeEmail.Verified,
	}, nil
}

func (c *Client) SetStatusChangeNotificationInterval(emailConfiguration *EmailConfiguration) error {
	form := map[string]string{
		"CID":       c.CID,
		"action":    "set_rate_limit_emails",
		"send_rate": strconv.Itoa(emailConfiguration.StatusChangeNotificationInterval),
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
