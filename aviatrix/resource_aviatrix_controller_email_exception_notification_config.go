package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerEmailExceptionNotificationConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixControllerEmailExceptionNotificationConfigCreate,
		ReadContext:   resourceAviatrixControllerEmailExceptionNotificationConfigRead,
		UpdateContext: resourceAviatrixControllerEmailExceptionNotificationConfigUpdate,
		DeleteContext: resourceAviatrixControllerEmailExceptionNotificationConfigDelete,
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
	client := meta.(*goaviatrix.Client)

	enableEmailExceptionNotification := d.Get("enable_email_exception_notification").(bool)
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

func resourceAviatrixControllerEmailExceptionNotificationConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("enable_email_exception_notification") {
		err := client.SetEmailExceptionNotification(ctx, d.Get("enable_email_exception_notification").(bool))
		if err != nil {
			return diag.Errorf("could not update email exception notification: %v", err)
		}
	}

	return resourceAviatrixControllerEmailExceptionNotificationConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerEmailExceptionNotificationConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.SetEmailExceptionNotification(ctx, true)
	if err != nil {
		return diag.Errorf("failed to enable email exception notification: %v", err)
	}

	return nil
}
