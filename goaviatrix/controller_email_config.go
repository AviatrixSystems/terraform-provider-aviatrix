package goaviatrix

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type EmailConfiguration struct {
	CID                string `form:"CID,omitempty"`
	Action             string `form:"action,omitempty"`
	AdminAlertEmail    string `json:"admin_alert_email,omitempty"`
	CriticalAlertEmail string `json:"critical_alert_email,omitempty"`
	SecurityEventEmail string `json:"security_event_email,omitempty"`
	StatusChangeEmail  string `json:"status_change_email,omitempty"`
}

//type EmailArgs struct {
//	admin_alert    map[string]string
//	critical_alert map[string]string
//	security_event map[string]string
//	status_change  map[string]string
//}

type EmailArgs struct {
	AdminAlertEmail    EmailArgsSingle `form:"admin_alert"`
	CriticalAlertEmail EmailArgsSingle `form:"critical_alert"`
	SecurityEventEmail EmailArgsSingle `form:"security_event"`
	StatusChangeEmail  EmailArgsSingle `form:"status_change"`
}

type EmailArgsSingle struct {
	Address string `form:"address"`
}

type ControllerEmailConfigResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

//type EmailArgs struct {
//	AdminAlertEmail    string `json:"admin_alert_email,omitempty"`
//	CriticalAlertEmail string `json:"critical_alert_email,omitempty"`
//	SecurityEventEmail int    `json:"security_event_email,omitempty"`
//	StatusChangeEmail  string `json:"status_change_email,omitempty"`
//}

func (c *Client) ConfigNotificationEmails(emailConfiguration *EmailConfiguration) error {
	//var emailArgs map[string]map[string]string
	//emailArgs := EmailArgs{
	//	AdminAlertEmail: EmailArgsSingle{
	//		Address: emailConfiguration.AdminAlertEmail,
	//	},
	//	CriticalAlertEmail: EmailArgsSingle{
	//		Address: emailConfiguration.CriticalAlertEmail,
	//	},
	//	SecurityEventEmail: EmailArgsSingle{
	//		Address: emailConfiguration.SecurityEventEmail,
	//	},
	//	StatusChangeEmail: EmailArgsSingle{
	//		Address: emailConfiguration.StatusChangeEmail,
	//	},
	//}
	adminEmail := map[string]interface{}{
		"address": emailConfiguration.AdminAlertEmail,
	}
	str, _ := json.Marshal(adminEmail)

	criticalEmail := map[string]interface{}{
		"address": emailConfiguration.CriticalAlertEmail,
	}
	str1, _ := json.Marshal(criticalEmail)

	securityEmail := map[string]interface{}{
		"address": emailConfiguration.SecurityEventEmail,
	}
	str2, _ := json.Marshal(securityEmail)

	statusEmail := map[string]interface{}{
		"address": emailConfiguration.StatusChangeEmail,
	}
	str3, _ := json.Marshal(statusEmail)

	emailArgs := map[string]string{
		"admin_alert":    string(str),
		"critical_alert": string(str1),
		"security_event": string(str2),
		"status_change":  string(str3),
	}

	//emailArgs := map[string]interface{}{
	//	"admin_alert":    adminEmail,
	//	"critical_alert": criticalEmail,
	//	"security_event": securityEmail,
	//	"status_change":  statusEmail,
	//}

	//emailArgs := map[string]interface{}{
	//	"admin_alert":    string(str),
	//	"critical_alert": string(str1),
	//	"security_event": string(str2),
	//	"status_change":  string(str3),
	//}

	log.Printf("zjin00: emailArgs is %v", emailArgs)
	str4, _ := json.Marshal(emailArgs)
	log.Printf("zjin01: emailArgs is %v", emailArgs)

	//if err != nil {
	//	return err
	//}
	log.Printf("zjin02: str is %v", string(str4))
	form := map[string]string{
		"CID":              c.CID,
		"action":           "add_notif_email_addr",
		"notif_email_args": string(str4),
	}
	log.Printf("zjin03: str is %v", form)

	var data ControllerEmailConfigResp
	return c.PostAPIWithResponse2(&data, form["action"], form, BasicCheck)
}
