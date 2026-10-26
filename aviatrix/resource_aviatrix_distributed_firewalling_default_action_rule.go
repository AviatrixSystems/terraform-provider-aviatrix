package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDistributedFirewallingDefaultActionRule() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingDefaultActionRuleCreate,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingDefaultActionRuleUpdate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingDefaultActionRuleRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingDefaultActionRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DENY", "PERMIT"}, true),
				Description: "Action for the specified source and destination Smart Groups." +
					"Must be one of PERMIT or DENY.",
			},
			"logging": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Boolean value to enable or disable logging for the default action rule.",
			},
			"log_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Logging profile UUID for the default action rule.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingDefaultActionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	action, ok := d.Get("action").(string)
	if !ok {
		return diag.Errorf("failed to assert 'action' as string")
	}

	logging, ok := d.Get("logging").(bool)
	if !ok {
		return diag.Errorf("failed to assert 'logging' as bool")
	}

	defaultActionRuleConfig := &goaviatrix.DistributedFirewallingDefaultActionRule{
		Action:  action,
		Logging: logging,
	}

	if logProfile, ok := d.GetOk("log_profile"); ok {
		defaultActionRuleConfig.LogProfile = logProfile.(string)
	}

	if err := client.UpdateDistributedFirewallingDefaultActionRule(ctx, defaultActionRuleConfig); err != nil {
		return diag.Errorf("failed to update the default action rule: %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))
	return resourceAviatrixDistributedFirewallingDefaultActionRuleRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingDefaultActionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	action, ok := d.Get("action").(string)
	if !ok {
		return diag.Errorf("failed to assert 'action' as string")
	}

	logging, ok := d.Get("logging").(bool)
	if !ok {
		return diag.Errorf("failed to assert 'logging' as bool")
	}

	defaultActionRuleConfig := &goaviatrix.DistributedFirewallingDefaultActionRule{
		Action:  action,
		Logging: logging,
	}

	if logProfile, ok := d.GetOk("log_profile"); ok {
		defaultActionRuleConfig.LogProfile = logProfile.(string)
	}

	if err := client.UpdateDistributedFirewallingDefaultActionRule(ctx, defaultActionRuleConfig); err != nil {
		return diag.Errorf("failed to update the default action rule: %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))

	return resourceAviatrixDistributedFirewallingDefaultActionRuleRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingDefaultActionRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.Id() != strings.ReplaceAll(client.ControllerIP, ".", "-") {
		return diag.Errorf("ID: %s does not match controller IP %q: please provide correct ID for importing", d.Id(), client.ControllerIP)
	}

	defaultActionRule, err := client.GetDistributedFirewallingDefaultActionRule(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read the default action rule: %s", err)
	}

	if err := d.Set("action", defaultActionRule.Action); err != nil {
		return diag.Errorf("failed to set 'action': %v", err)
	}

	if err := d.Set("logging", defaultActionRule.Logging); err != nil {
		return diag.Errorf("failed to set 'logging': %v", err)
	}

	// Only update log_profile if the API returns a non-empty value, it's empty by default
	if defaultActionRule.LogProfile != "" {
		if err := d.Set("log_profile", defaultActionRule.LogProfile); err != nil {
			return diag.Errorf("failed to set 'log_profile': %v", err)
		}
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))
	return nil
}

func resourceAviatrixDistributedFirewallingDefaultActionRuleDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	defaultActionRuleConfig := &goaviatrix.DistributedFirewallingDefaultActionRule{
		Action:     "PERMIT",
		Logging:    false,
		LogProfile: "",
	}

	if err := client.UpdateDistributedFirewallingDefaultActionRule(ctx, defaultActionRuleConfig); err != nil {
		return diag.Errorf("failed to update the default action rule: %v", err)
	}

	return nil
}
