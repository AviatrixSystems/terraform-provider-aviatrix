package aviatrix

import (
	"context"
	"fmt"
	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerGatewayKeepaliveConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceControllerGatewayKeepaliveConfigCreate,
		ReadContext:   resourceControllerGatewayKeepaliveConfigRead,
		UpdateContext: resourceControllerGatewayKeepaliveConfigUpdate,
		DeleteContext: resourceControllerGatewayKeepaliveConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"keep_alive_speed": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateControllerGatewayKeepaliveSpeed,
				Description:  "Gateway keepalive speed.",
			},
		},
	}
}

func resourceControllerGatewayKeepaliveConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	speed := d.Get("keep_alive_speed").(string)
	err := client.SetGatewayKeepaliveConfig(ctx, speed)
	if err != nil {
		return diag.Errorf("could not create Controller Gateway Keepalive Config: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceControllerGatewayKeepaliveConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	speed, err := client.GetGatewayKeepaliveConfig(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Controller Gateway Keepalive Config: %v", err)
	}

	d.Set("keep_alive_speed", speed)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceControllerGatewayKeepaliveConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("keep_alive_speed") {
		speed := d.Get("keep_alive_speed").(string)
		err := client.SetGatewayKeepaliveConfig(ctx, speed)
		if err != nil {
			return diag.Errorf("could not update Controller Gateway Keepalive Config: %v", err)
		}
	}

	return nil
}

func resourceControllerGatewayKeepaliveConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.SetGatewayKeepaliveConfig(ctx, "medium")
	if err != nil {
		return diag.Errorf("could not destroy Controller Gateway Keepalive Config: %v", err)
	}
	return nil
}

func validateControllerGatewayKeepaliveSpeed(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if !stringInSlice(v, []string{"slow", "medium", "fast"}) {
		errs = append(errs, fmt.Errorf("%s must be one of slow, medium or fast", key))
	}
	return
}
