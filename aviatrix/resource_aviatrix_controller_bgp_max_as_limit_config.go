package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixControllerBgpMaxAsLimitConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixControllerBgpMaxAsLimitConfigCreate,
		ReadContext:   resourceAviatrixControllerBgpMaxAsLimitConfigRead,
		UpdateContext: resourceAviatrixControllerBgpMaxAsLimitConfigUpdate,
		DeleteContext: resourceAviatrixControllerBgpMaxAsLimitConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"max_as_limit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 254),
				Description:  "The maximum AS path limit allowed by transit gateways when handling BGP/Peering route propagation.",
			},
		},
	}
}

func resourceAviatrixControllerBgpMaxAsLimitConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	maxAsLimit := d.Get("max_as_limit").(int)
	err := client.SetControllerBgpMaxAsLimit(ctx, maxAsLimit)
	if err != nil {
		return diag.Errorf("failed to create controller BGP max AS limit config: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpMaxAsLimitConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpMaxAsLimitConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	maxAsLimit, err := client.GetControllerBgpMaxAsLimit(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get controller BGP max AS limit config: %v", err)
	}

	d.Set("max_as_limit", maxAsLimit)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("max_as_limit") {
		maxAsLimit := d.Get("max_as_limit").(int)
		err := client.SetControllerBgpMaxAsLimit(ctx, maxAsLimit)
		if err != nil {
			return diag.Errorf("failed to update controller BGP max AS limit config: %v", err)
		}
	}

	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableControllerBgpMaxAsLimit(ctx)
	if err != nil {
		return diag.Errorf("failed to delete controller BGP max AS limit config: %v", err)
	}

	return nil
}
