package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerEmailExceptionNotificationConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerEmailExceptionNotificationConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerEmailExceptionNotificationConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerEmailExceptionNotificationConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerEmailExceptionNotificationConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_email_exception_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable email exception notification.",
			},
		},
	}
}

func resourceAviatrixControllerEmailExceptionNotificationConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	enableEmailExceptionNotification := getBool(d, "enable_email_exception_notification")
	if !enableEmailExceptionNotification {
		err := client.SetEmailExceptionNotification(ctx, false)
		if err != nil {
			return diag.Errorf("could not disable email exception notification: %v", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerEmailExceptionNotificationConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerEmailExceptionNotificationConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	enableEmailExceptionNotification, err := client.GetEmailExceptionNotificationStatus(ctx)
	if err != nil {
		return diag.Errorf("could not get exception email notification status: %v", err)
	}
	mustSet(d, "enable_email_exception_notification", enableEmailExceptionNotification)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerEmailExceptionNotificationConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.HasChange("enable_email_exception_notification") {
		err := client.SetEmailExceptionNotification(ctx, getBool(d, "enable_email_exception_notification"))
		if err != nil {
			return diag.Errorf("could not update email exception notification: %v", err)
		}
	}

	return resourceAviatrixControllerEmailExceptionNotificationConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerEmailExceptionNotificationConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.SetEmailExceptionNotification(ctx, true)
	if err != nil {
		return diag.Errorf("failed to enable email exception notification: %v", err)
	}

	return nil
}
