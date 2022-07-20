package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerPrivateModeConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerPrivateModeConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerPrivateModeConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerPrivateModeConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerPrivateModeConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_private_mode": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable Private Mode on the Controller.",
			},
			"copilot_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Copilot instance ID to associate with the Controller for Private Mode.",
			},
			"proxies": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of proxies.",
			},
		},
	}
}

func resourceAviatrixControllerPrivateModeConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enablePrivateMode := d.Get("enable_private_mode").(bool)
	if !enablePrivateMode {
		for _, key := range []string{"copilot_instance_id", "proxies"} {
			if _, ok := d.GetOk(key); ok {
				return diag.Errorf("failed to create Controller Private Mode config: %q and %q must be empty when %q is false", "copilot_instance_id", "proxies", "enable_private_mode")
			}
		}
	}

	flag := false
	defer resourceAviatrixControllerPrivateModeConfigReadIfRequired(ctx, d, meta, &flag)

	if enablePrivateMode {
		err := client.EnablePrivateMode(ctx)
		if err != nil {
			return diag.Errorf("failed to enable Private Mode: %s", err)
		}
	} else {
		err := client.DisablePrivateMode(ctx)
		if err != nil {
			return diag.Errorf("failed to disable Private Mode: %s", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))

	if _, ok := d.GetOk("copilot_instance_id"); ok {
		copilotInstanceId := d.Get("copilot_instance_id").(string)
		err := client.UpdatePrivateModeCopilot(ctx, copilotInstanceId)
		if err != nil {
			return diag.Errorf("failed to set Copilot instance ID: %s", err)
		}
	}

	if _, ok := d.GetOk("proxies"); ok {
		proxiesSet := d.Get("proxies").(*schema.Set)
		proxies := make([]string, proxiesSet.Len())
		for i, proxy := range proxiesSet.List() {
			proxies[i] = proxy.(string)
		}

		err := client.UpdatePrivateModeControllerProxies(ctx, proxies)
		if err != nil {
			return diag.Errorf("failed to set Controller proxies: %s", err)
		}
	}

	return resourceAviatrixControllerPrivateModeConfigReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixControllerPrivateModeConfigReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixControllerPrivateModeConfigRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixControllerPrivateModeConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	controllerPrivateModeConfig, err := client.GetPrivateModeInfo(ctx)
	if err != nil {
		return diag.Errorf("failed to read Controller Private Mode Config: %s", err)
	}

	d.Set("enable_private_mode", controllerPrivateModeConfig.EnablePrivateMode)
	d.Set("copilot_instance_id", controllerPrivateModeConfig.CopilotInstanceID)

	if len(controllerPrivateModeConfig.Proxies) > 0 {
		err = d.Set("proxies", controllerPrivateModeConfig.Proxies)
		if err != nil {
			return diag.Errorf("failed to set Controller proxies during Controller PrivateMode Config read: %s", err)
		}
	} else {
		d.Set("proxies", nil)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerPrivateModeConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enablePrivateMode := d.Get("enable_private_mode").(bool)
	if d.HasChanges("enable_private_mode", "copilot_instance_id", "proxies") && !enablePrivateMode {
		for _, key := range []string{"copilot_instance_id", "proxies"} {
			if _, ok := d.GetOk(key); ok {
				return diag.Errorf("failed to update Controller Private Mode config: %q and %q must be empty when %q is false", "copilot_instance_id", "proxies", "enable_private_mode")
			}
		}
	}

	if d.HasChange("enable_private_mode") {
		if enablePrivateMode {
			err := client.EnablePrivateMode(ctx)
			if err != nil {
				return diag.Errorf("failed to enable Private Mode during update: %s", err)
			}
		} else {
			err := client.DisablePrivateMode(ctx)
			if err != nil {
				return diag.Errorf("failed to disable Private Mode during update: %s", err)
			}
		}
	}

	if d.HasChange("copilot_instance_id") {
		copilotInstanceId := d.Get("copilot_instance_id").(string)
		err := client.UpdatePrivateModeCopilot(ctx, copilotInstanceId)
		if err != nil {
			return diag.Errorf("failed to set Copilot instance ID during update: %s", err)
		}
	}

	if d.HasChange("proxies") {
		proxiesSet := d.Get("proxies").(*schema.Set)
		proxies := make([]string, proxiesSet.Len())
		for i, proxy := range proxiesSet.List() {
			proxies[i] = proxy.(string)
		}

		err := client.UpdatePrivateModeControllerProxies(ctx, proxies)
		if err != nil {
			return diag.Errorf("failed to set Controller proxies during update: %s", err)
		}
	}

	return resourceAviatrixControllerPrivateModeConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerPrivateModeConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisablePrivateMode(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Controller Private Mode config: %s", err)
	}

	return nil
}
