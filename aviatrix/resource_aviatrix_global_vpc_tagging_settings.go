package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixGlobalVpcTaggingSettings() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixGlobalVpcTaggingSettingsCreate,
		ReadWithoutTimeout:   resourceAviatrixGlobalVpcTaggingSettingsRead,
		UpdateWithoutTimeout: resourceAviatrixGlobalVpcTaggingSettingsUpdate,
		DeleteWithoutTimeout: resourceAviatrixGlobalVpcTaggingSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"service_state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"semi_automatic", "automatic", "disabled"}, false),
				Description:  "Service state.",
			},
			"enable_alert": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set to true to enable alert.",
			},
		},
	}
}

func marshalGlobalVpcTaggingSettingsInput(d *schema.ResourceData) *goaviatrix.GlobalVpcTaggingSettings {
	globalVpcTaggingSettings := &goaviatrix.GlobalVpcTaggingSettings{
		ServiceState: d.Get("service_state").(string),
		EnableAlert:  d.Get("enable_alert").(bool),
	}

	return globalVpcTaggingSettings
}

func resourceAviatrixGlobalVpcTaggingSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	globalVpcTaggingSettings := marshalGlobalVpcTaggingSettingsInput(d)

	flag := false
	defer resourceAviatrixGlobalVpcTaggingSettingsReadIfRequired(ctx, d, meta, &flag)

	err := client.UpdateGlobalVpcTaggingSettings(ctx, globalVpcTaggingSettings)
	if err != nil {
		return diag.Errorf("failed to create global vpc tagging settings: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixGlobalVpcTaggingSettingsReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixGlobalVpcTaggingSettingsReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixGlobalVpcTaggingSettingsRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixGlobalVpcTaggingSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	globalVpcTaggingSettings, err := client.GetGlobalVpcTaggingSettings(ctx)
	if err != nil {
		return diag.Errorf("failed to read global vpc tagging settings: %s", err)
	}

	d.Set("service_state", globalVpcTaggingSettings.ServiceState)
	d.Set("enable_alert", globalVpcTaggingSettings.EnableAlert)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixGlobalVpcTaggingSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChanges("service_state", "enable_alert") {
		globalVpcTaggingSettings := marshalGlobalVpcTaggingSettingsInput(d)

		err := client.UpdateGlobalVpcTaggingSettings(ctx, globalVpcTaggingSettings)
		if err != nil {
			return diag.Errorf("failed to update global vpc tagging settings: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixGlobalVpcTaggingSettingsRead(ctx, d, meta)
}

func resourceAviatrixGlobalVpcTaggingSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	globalVpcTaggingSettings := &goaviatrix.GlobalVpcTaggingSettings{
		ServiceState: "semi_automatic",
		EnableAlert:  false,
	}
	err := client.UpdateGlobalVpcTaggingSettings(ctx, globalVpcTaggingSettings)
	if err != nil {
		return diag.Errorf("failed to delete global vpc tagging settings: %v", err)
	}

	return nil
}
