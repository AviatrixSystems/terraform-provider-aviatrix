package aviatrix

import (
	"context"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixTunnelDetectionTime() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixTunnelDetectionTimeCreate,
		ReadContext:   resourceAviatrixTunnelDetectionTimeRead,
		UpdateContext: resourceAviatrixTunnelDetectionTimeUpdate,
		DeleteContext: resourceAviatrixTunnelDetectionTimeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"aviatrix_entity": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Controller",
				ForceNew:    true,
				Description: "Gateway name to change IPSec tunnel down detection time for. If empty or set to \"Controller\", all gateways share the same tunnel down detection time.",
			},
			"detection_time": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(20, 600),
				Description:  "Specify a IPSec tunnel down detection time. The minimum is 20 seconds. The maximum is 600 seconds.",
			},
		},
	}
}

func resourceAviatrixTunnelDetectionTimeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	aviatrixEntity := d.Get("aviatrix_entity").(string)
	detectionTime := d.Get("detection_time").(int)
	err := client.SetTunnelDetectionTime(ctx, detectionTime, aviatrixEntity)
	if err != nil {
		return diag.Errorf("failed to create tunnel detection time resource: %v", err)
	}

	d.SetId(aviatrixEntity)
	return nil
}

func resourceAviatrixTunnelDetectionTimeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	aviatrixEntity := d.Get("aviatrix_entity").(string)
	if aviatrixEntity == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no aviatrix_entity name received. Import Id is %s", id)
		d.SetId(id)
		aviatrixEntity = id
	}
	d.Set("aviatrix_entity", aviatrixEntity)

	detectionTime, err := client.GetTunnelDetectionTime(ctx, aviatrixEntity)
	if err != nil {
		return diag.Errorf("failed to read tunnel detection time: %v", err)
	}
	d.Set("detection_time", detectionTime)

	return nil
}

func resourceAviatrixTunnelDetectionTimeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	aviatrixEntity := d.Get("aviatrix_entity").(string)
	if d.HasChange("detection_time") {
		detectionTime := d.Get("detection_time").(int)
		err := client.SetTunnelDetectionTime(ctx, detectionTime, aviatrixEntity)
		if err != nil {
			return diag.Errorf("failed to create tunnel detection time resource: %v", err)
		}
	}

	return nil
}

func resourceAviatrixTunnelDetectionTimeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	aviatrixEntity := d.Get("aviatrix_entity").(string)
	var defaultDetectionTime int
	if aviatrixEntity == "Controller" {
		defaultDetectionTime = 60
	} else {
		var err error
		defaultDetectionTime, err = client.GetTunnelDetectionTime(ctx, "Controller")
		if err != nil {
			return diag.Errorf("failed to delete tunnel detection time resource: %v", err)
		}
	}

	err := client.SetTunnelDetectionTime(ctx, defaultDetectionTime, aviatrixEntity)
	if err != nil {
		return diag.Errorf("failed to delete tunnel detection time resource: %v", err)
	}

	return nil
}
