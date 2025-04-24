package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixDistributedFirewallingZeroTrustRule() *schema.Resource {

	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingZeroTrustRuleCreate,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingZeroTrustRuleUpdate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingZeroTrustRuleRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingZeroTrustRuleDelete,

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
				Description: "Boolean value to enable or disable logging for the zero trust rule.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*goaviatrix.Client)

	zeroTrustRuleConfig := &goaviatrix.DistributedFirewallingZeroTrustRule{
		Action:  d.Get("action").(string),
		Logging: d.Get("logging").(bool),
	}

	if err := client.UpdateDistributedFirewallingZeroTrust(ctx, zeroTrustRuleConfig); err != nil {
		return diag.Errorf("failed to update the zero trust rule: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingZeroTrustRuleRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*goaviatrix.Client)

	zeroTrustRuleConfig := &goaviatrix.DistributedFirewallingZeroTrustRule{
		Action:  d.Get("action").(string),
		Logging: d.Get("logging").(bool),
	}

	if err := client.UpdateDistributedFirewallingZeroTrust(ctx, zeroTrustRuleConfig); err != nil {
		return diag.Errorf("failed to update the zero trust rule: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingZeroTrustRuleRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	postRuleList, err := client.GetDistributedFirewallingZeroTrustRule(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read the zero trust rule: %s", err)
	}

	d.Set("action", postRuleList["action"])
	d.Set("logging", postRuleList["logging"])
	d.Set("uuid", postRuleList["uuid"])

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// restore to the original zero trust rule
	zeroTrustRuleConfig := &goaviatrix.DistributedFirewallingZeroTrustRule{
		Action:  "PERMIT",
		Logging: false,
	}

	if err := client.UpdateDistributedFirewallingZeroTrust(ctx, zeroTrustRuleConfig); err != nil {
		return diag.Errorf("failed to update the zero trust rule: %v", err)
	}

	return nil
}
