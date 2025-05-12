package aviatrix

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixDistributedFirewallingDefaultActionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingDefaultActionPolicyCreate,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingDefaultActionPolicyUpdate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingDefaultActionPolicyRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingDefaultActionPolicyDelete,

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
				ForceNew:    true,
				Description: "Boolean value to enable or disable logging for the default action policy.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingDefaultActionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	defaultActionPolicyConfig := &goaviatrix.DistributedFirewallingDefaultActionPolicy{
		Action:  action,
		Logging: logging,
	}

	if err := client.UpdateDistributedFirewallingDefaultActionPolicy(ctx, defaultActionPolicyConfig); err != nil {
		return diag.Errorf("failed to update the default action policy: %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))
	return resourceAviatrixDistributedFirewallingDefaultActionPolicyRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingDefaultActionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	defaultActionPolicyConfig := &goaviatrix.DistributedFirewallingDefaultActionPolicy{
		Action:  action,
		Logging: logging,
	}

	if err := client.UpdateDistributedFirewallingDefaultActionPolicy(ctx, defaultActionPolicyConfig); err != nil {
		return diag.Errorf("failed to update the default action policy: %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))

	return resourceAviatrixDistributedFirewallingDefaultActionPolicyRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingDefaultActionPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.Id() != strings.ReplaceAll(client.ControllerIP, ".", "-") {
		return diag.Errorf("ID: %s does not match controller IP %q: please provide correct ID for importing", d.Id(), client.ControllerIP)
	}

	defaultActionPolicy, err := client.GetDistributedFirewallingDefaultActionPolicy(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read the default action policy: %s", err)
	}

	if err := d.Set("action", defaultActionPolicy["action"]); err != nil {
		return diag.Errorf("failed to set 'action': %v", err)
	}

	logging, err := strconv.ParseBool(defaultActionPolicy["logging"])
	if err != nil {
		return diag.Errorf("failed to parse 'logging' as bool: %v", err)
	}
	if err := d.Set("logging", logging); err != nil {
		return diag.Errorf("failed to set 'logging': %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))
	return nil
}

func resourceAviatrixDistributedFirewallingDefaultActionPolicyDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}


	defaultActionPolicyConfig := &goaviatrix.DistributedFirewallingDefaultActionPolicy{
		Action:  "PERMIT",
		Logging: false,
	}

	if err := client.UpdateDistributedFirewallingDefaultActionPolicy(ctx, defaultActionPolicyConfig); err != nil {
		return diag.Errorf("failed to update the default action policy: %v", err)
	}

	return nil
}
