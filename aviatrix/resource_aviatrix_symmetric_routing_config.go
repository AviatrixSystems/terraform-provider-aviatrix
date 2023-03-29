package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSymmetricRoutingConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixSymmetricRoutingConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixSymmetricRoutingConfigRead,
		UpdateWithoutTimeout: resourceAviatrixSymmetricRoutingConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixSymmetricRoutingConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_symmetric_routing": {
				Type:        schema.TypeBool,
				Required:    true,
				Default:     false,
				Description: "Specify whether to enable symmetric routing for a given spoke gateway or not.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Spoke gateway name.",
			},
		},
	}
}

func resourceAviatrixSymmetricRoutingConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	enableSymmetricRouting := d.Get("enable_symmetric_routing").(bool)

	if enableSymmetricRouting {
		status, _ := client.GetSymmetricRoutingStatus(ctx, gwName)
		if strings.Contains(status, "enabled") {
			log.Printf("[INFO] Symmetric routing is already enabled")
		} else {
			err := client.EnableSymmetricRouting(ctx, gwName)
			if err != nil {
				return diag.Errorf("failed to enable symmetric routing for spoke gateway %s: %s", gwName, err)
			}
		}
	} else {
		status, _ := client.GetSymmetricRoutingStatus(ctx, gwName)
		if strings.Contains(status, "disabled") {
			log.Printf("[INFO] Symmetric routing is already disabled")
		} else {
			err := client.DisableSymmetricRouting(ctx, gwName)
			if err != nil {
				return diag.Errorf("failed to disable symmetric routing for spoke gateway %s: %s", gwName, err)
			}
		}
	}

	d.SetId(gwName)
	return resourceAviatrixSymmetricRoutingConfigRead(ctx, d, meta)
}

func resourceAviatrixSymmetricRoutingConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	status, err := client.GetSymmetricRoutingStatus(ctx, d.Get("gw_name").(string))
	if err != nil {
		return diag.Errorf("could not read symmetric routing status: %s", err)
	}

	if strings.Contains(status, "enabled") {
		d.Set("enable_symmetric_routing", true)
	} else if strings.Contains(status, "disabled") {
		d.Set("enable_symmetric_routing", false)
	}

	d.SetId(d.Get("gw_name").(string))
	return nil
}

func resourceAviatrixSymmetricRoutingConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("enable_symmetric_routing") {
		gwName := d.Get("gw_name").(string)
		enableSymmetricRouting := d.Get("enable_symmetric_routing").(bool)

		if enableSymmetricRouting {
			err := client.EnableSymmetricRouting(ctx, gwName)
			if err != nil {
				return diag.Errorf("failed to enable symmetric routing for spoke gateway %s during update: %s", gwName, err)
			}
		} else {
			err := client.DisableSymmetricRouting(ctx, gwName)
			if err != nil {
				return diag.Errorf("failed to disable symmetric routing for spoke gateway %s during update: %s", gwName, err)
			}
		}
	}

	return resourceAviatrixSymmetricRoutingConfigRead(ctx, d, meta)
}

func resourceAviatrixSymmetricRoutingConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)

	err := client.DisableSymmetricRouting(ctx, gwName)
	if err != nil {
		return diag.Errorf("failed to disable symmetric routing for spoke gateway %s during deletion: %s", gwName, err)
	}

	return nil
}
