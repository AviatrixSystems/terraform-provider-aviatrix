package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixControllerEmailConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerEmailConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerEmailConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerEmailConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerEmailConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"admin_alert_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email to receive important account and certification information.",
			},
			"critical_alert_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email to receive field notices and critical notices.",
			},
			"security_event_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email to receive security and CVE (Common Vulnerabilities and Exposures) notification emails.",
			},
			"status_change_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email to receive system/tunnel status notification emails.",
			},
			"status_change_notification_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Status change notification interval in seconds.",
			},
			"admin_alert_email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether admin alert notification email is verified.",
			},
			"critical_alert_email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether critical alert notification email is verified.",
			},
			"security_event_email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether security event notification email is verified.",
			},
			"status_change_email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether status change notification email is verified.",
			},
		},
	}
}

func resourceAviatrixControllerEmailConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	emailConfiguration := &goaviatrix.EmailConfiguration{
		AdminAlertEmail:                  getString(d, "admin_alert_email"),
		CriticalAlertEmail:               getString(d, "critical_alert_email"),
		SecurityEventEmail:               getString(d, "security_event_email"),
		StatusChangeEmail:                getString(d, "status_change_email"),
		StatusChangeNotificationInterval: getInt(d, "status_change_notification_interval"),
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixControllerEmailConfigReadIfRequired(ctx, d, meta, &flag)

	err := client.ConfigNotificationEmails(ctx, emailConfiguration)
	if err != nil {
		return diag.Errorf("could not config controller emails: %v", err)
	}

	if emailConfiguration.StatusChangeNotificationInterval != 60 {
		err := client.SetStatusChangeNotificationInterval(emailConfiguration)
		if err != nil {
			return diag.Errorf("could not set status change notification interval: %v", err)
		}
	}

	return resourceAviatrixControllerEmailConfigReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixControllerEmailConfigReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixControllerEmailConfigRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixControllerEmailConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	emailConfiguration, err := client.GetNotificationEmails(ctx)
	if err != nil {
		return diag.Errorf("could not get notification emails: %v", err)
	}
	mustSet(d, "admin_alert_email", emailConfiguration.AdminAlertEmail)
	mustSet(d, "critical_alert_email", emailConfiguration.CriticalAlertEmail)
	mustSet(d, "security_event_email", emailConfiguration.SecurityEventEmail)
	mustSet(d, "status_change_email", emailConfiguration.StatusChangeEmail)
	mustSet(d, "status_change_notification_interval", emailConfiguration.StatusChangeNotificationInterval)
	mustSet(d, "admin_alert_email_verified", emailConfiguration.AdminAlertEmailVerified)
	mustSet(d, "critical_alert_email_verified", emailConfiguration.CriticalAlertEmailVerified)
	mustSet(d, "security_event_email_verified", emailConfiguration.SecurityEventEmailVerified)
	mustSet(d, "status_change_email_verified", emailConfiguration.StatusChangeEmailVerified)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerEmailConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.HasChanges("admin_alert_email", "critical_alert_email", "security_event_email", "status_change_email") {
		emailConfiguration := &goaviatrix.EmailConfiguration{}
		if d.HasChange("admin_alert_email") {
			emailConfiguration.AdminAlertEmail = getString(d, "admin_alert_email")
		}
		if d.HasChange("critical_alert_email") {
			emailConfiguration.CriticalAlertEmail = getString(d, "critical_alert_email")
		}
		if d.HasChange("security_event_email") {
			emailConfiguration.SecurityEventEmail = getString(d, "security_event_email")
		}
		if d.HasChange("status_change_email") {
			emailConfiguration.StatusChangeEmail = getString(d, "status_change_email")
		}

		err := client.ConfigNotificationEmails(ctx, emailConfiguration)
		if err != nil {
			return diag.Errorf("could not config controller emails: %v", err)
		}
	}

	if d.HasChange("status_change_notification_interval") {
		emailConfiguration := &goaviatrix.EmailConfiguration{
			StatusChangeNotificationInterval: getInt(d, "status_change_notification_interval"),
		}

		err := client.SetStatusChangeNotificationInterval(emailConfiguration)
		if err != nil {
			return diag.Errorf("could not update status change notification interval: %v", err)
		}
	}

	return resourceAviatrixControllerEmailConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerEmailConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	emailConfiguration := &goaviatrix.EmailConfiguration{
		StatusChangeNotificationInterval: 60,
	}

	err := client.SetStatusChangeNotificationInterval(emailConfiguration)
	if err != nil {
		return diag.Errorf("could not set status change notification interval to default value: %v", err)
	}

	return nil
}
