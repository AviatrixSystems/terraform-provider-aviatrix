package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDistributedFirewallingConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingConfigRead,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_distributed_firewalling": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable Distributed-firewalling.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enableDFW := d.Get("enable_distributed_firewalling").(bool)
	if enableDFW {
		err := client.EnableDistributedFirewalling(ctx)
		if err != nil {
			return diag.Errorf("failed to enable Distributed-firewalling: %s", err)
		}
	} else {
		err := client.DisableDistributedFirewalling(ctx)
		if err != nil {
			return diag.Errorf("failed to disable Distributed-firewalling: %s", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingConfigRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	distributedFirewalling, err := client.GetDistributedFirewallingStatus(ctx)
	if err != nil {
		return diag.Errorf("failed to read Distributed-firewalling status: %s", err)
	}
	d.Set("enable_distributed_firewalling", distributedFirewalling.EnableDistributedFirewalling)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("enable_distributed_firewalling") {
		distributedFirewalling := d.Get("enable_distributed_firewalling").(bool)

		if distributedFirewalling {
			err := client.EnableDistributedFirewalling(ctx)
			if err != nil {
				return diag.Errorf("failed to enable Distributed-firewalling during update: %s", err)
			}
		} else {
			err := client.DisableDistributedFirewalling(ctx)
			if err != nil {
				return diag.Errorf("failed to disable Distributed-firewalling during update: %s", err)
			}
		}
	}

	return resourceAviatrixDistributedFirewallingConfigRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableDistributedFirewalling(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling config: %s", err)
	}

	return nil
}
