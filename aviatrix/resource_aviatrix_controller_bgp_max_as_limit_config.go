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
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 254),
				Description:  "The maximum AS path limit allowed by transit gateways when handling BGP/Peering route propagation for RFC 1918 CIDRs.",
			},
			"max_as_limit_non_rfc1918": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 254),
				Description:  "The maximum AS path limit allowed by transit gateways when handling BGP/Peering route propagation for non-RFC 1918 CIDRs.",
			},
		},
	}
}

func resourceAviatrixControllerBgpMaxAsLimitConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	maxAsLimit, maxAsLimitOk := d.GetOk("max_as_limit")
	maxAsLimitNonRfc1918, maxAslimitNonRfc1918Ok := d.GetOk("max_as_limit_non_rfc1918")

	if !maxAsLimitOk && !maxAslimitNonRfc1918Ok {
		return diag.Errorf("error creating controller BGP max AS limit config: at least one of max_as_limit and max_as_limit_non_rfc1918 must be provided")
	}

	if maxAsLimitOk {
		err := client.SetControllerBgpMaxAsLimit(ctx, maxAsLimit.(int))
		if err != nil {
			return diag.Errorf("failed to create controller BGP max AS limit config for RGC 1918 CIDRs: %v", err)
		}
	}

	if maxAslimitNonRfc1918Ok {
		err := client.SetControllerBgpMaxAsLimitNonRfc1918(ctx, maxAsLimitNonRfc1918.(int))
		if err != nil {
			return diag.Errorf("failed to create controller BGP max AS limit config for non-RFC 1918 CIDRs: %v", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpMaxAsLimitConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpMaxAsLimitConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	maxAsLimit, maxAsLimitNonRfc1918, err := client.GetControllerBgpMaxAsLimit(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get controller BGP max AS limit config: %v", err)
	}

	d.Set("max_as_limit", maxAsLimit)
	d.Set("max_as_limit_non_rfc1918", maxAsLimitNonRfc1918)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerBgpMaxAsLimitConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("max_as_limit") {
		_, maxAsLimitOk := d.GetOk("max_as_limit")
		if maxAsLimitOk {
			maxAsLimit := d.Get("max_as_limit").(int)
			err := client.SetControllerBgpMaxAsLimit(ctx, maxAsLimit)
			if err != nil {
				return diag.Errorf("failed to update controller BGP max AS limit config: %v", err)
			}
		} else {
			err := client.DisableControllerBgpMaxAsLimit(ctx)
			if err != nil {
				return diag.Errorf("failed to update controller BGP max AS limit config: %v", err)
			}
		}
	}

	if d.HasChange("max_as_limit_non_rfc1918") {
		_, maxAsLimitNonRfc1918Ok := d.GetOk("max_as_limit_non_rfc1918")
		if maxAsLimitNonRfc1918Ok {
			maxAsLimitNonRfc1918 := d.Get("max_as_limit_non_rfc1918").(int)
			err := client.SetControllerBgpMaxAsLimitNonRfc1918(ctx, maxAsLimitNonRfc1918)
			if err != nil {
				return diag.Errorf("failed to update controller BGP max AS limit config for non-RFC 1918 CIDRs: %v", err)
			}
		} else {
			err := client.DisableControllerBgpMaxAsLimitNonRfc1918(ctx)
			if err != nil {
				return diag.Errorf("failed to update controller BGP max AS limit config for non-RFC 1918 CIDRs: %v", err)
			}
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

	err = client.DisableControllerBgpMaxAsLimitNonRfc1918(ctx)
	if err != nil {
		return diag.Errorf("failed to delete controller BGP max AS limit config for non-RFC 1918 CIDRs: %v", err)
	}

	return nil
}
