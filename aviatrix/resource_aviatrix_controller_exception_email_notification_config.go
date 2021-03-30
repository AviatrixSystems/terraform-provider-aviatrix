package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerExceptionEmailNotificationConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixControllerExceptionEmailNotificationConfigCreate,
		ReadContext:   resourceAviatrixControllerExceptionEmailNotificationConfigRead,
		UpdateContext: resourceAviatrixControllerExceptionEmailNotificationConfigUpdate,
		DeleteContext: resourceAviatrixControllerExceptionEmailNotificationConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_exception_email_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable exception email notification.",
			},
		},
	}
}

func resourceAviatrixControllerExceptionEmailNotificationConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enableExceptionEmailNotification := d.Get("enable_exception_email_notification").(bool)
	if !enableExceptionEmailNotification {
		err := client.SetExceptionEmailNotification(ctx, false)
		if err != nil {
			return diag.Errorf("could not disable exception email notification: %v", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerExceptionEmailNotificationConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	enableExceptionEmailNotification, err := client.GetExceptionEmailNotificationStatus(ctx)
	if err != nil {
		return diag.Errorf("could not get exception email notification status: %v", err)
	}
	d.Set("enable_exception_email_notification", enableExceptionEmailNotification)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerExceptionEmailNotificationConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("enable_exception_email_notification") {
		err := client.SetExceptionEmailNotification(ctx, d.Get("enable_exception_email_notification").(bool))
		if err != nil {
			return diag.Errorf("could not update exception email notification: %v", err)
		}
	}

	return nil
}

func resourceAviatrixControllerExceptionEmailNotificationConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.SetExceptionEmailNotification(ctx, true)
	if err != nil {
		return diag.Errorf("failed to enable exception email notification: %v", err)
	}

	return nil
}
