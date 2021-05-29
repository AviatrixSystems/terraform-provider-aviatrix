package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixCopilotAssociation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotAssociationCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotAssociationRead,
		DeleteWithoutTimeout: resourceAviatrixCopilotAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"copilot_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "CoPilot IP Address or Hostname.",
			},
		},
	}
}

func resourceAviatrixCopilotAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	addr := d.Get("copilot_address").(string)
	err := client.EnableCopilotAssociation(ctx, addr)
	if err != nil {
		return diag.Errorf("could not associate copilot: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixCopilotAssociationRead(ctx, d, meta)
}

func resourceAviatrixCopilotAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	copilot, err := client.GetCopilotAssociationStatus(ctx)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not get copilot association status: %v", err)
	}

	d.Set("copilot_address", copilot.IP)
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixCopilotAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableCopilotAssociation(ctx)
	if err != nil {
		return diag.Errorf("could not disable copilot association: %v", err)
	}

	return nil
}
