package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Optional:    true,
				Default:     true,
				Description: "Enable email exception notification.",
			},
			"critical_alert_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     true,
				Description: "Enable email exception notification.",
			},
			"security_event_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     true,
				Description: "Enable email exception notification.",
			},
			"status_change_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     true,
				Description: "Enable email exception notification.",
			},
			//"status_change_notification_interval": {
			//	Type:        schema.TypeInt,
			//	Optional:    true,
			//	Default:     true,
			//	Description: "Enable email exception notification.",
			//},
		},
	}
}

func resourceAviatrixControllerEmailConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	emailConfiguration := &goaviatrix.EmailConfiguration{
		AdminAlertEmail:    d.Get("admin_alert_email").(string),
		CriticalAlertEmail: d.Get("critical_alert_email").(string),
		SecurityEventEmail: d.Get("security_event_email").(string),
		StatusChangeEmail:  d.Get("status_change_email").(string),
	}

	err := client.ConfigNotificationEmails(emailConfiguration)
	if err != nil {
		return diag.Errorf("could not config controller emails: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
	//return resourceAviatrixControllerEmailConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerEmailConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	enableEmailExceptionNotification, err := client.GetEmailExceptionNotificationStatus(ctx)
	if err != nil {
		return diag.Errorf("could not get exception email notification status: %v", err)
	}
	d.Set("enable_email_exception_notification", enableEmailExceptionNotification)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerEmailConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("enable_email_exception_notification") {
		err := client.SetEmailExceptionNotification(ctx, d.Get("enable_email_exception_notification").(bool))
		if err != nil {
			return diag.Errorf("could not update email exception notification: %v", err)
		}
	}

	return resourceAviatrixControllerEmailExceptionNotificationConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerEmailConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := meta.(*goaviatrix.Client)
	//
	//err := client.SetEmailExceptionNotification(ctx, true)
	//if err != nil {
	//	return diag.Errorf("failed to enable email exception notification: %v", err)
	//}

	return nil
}
