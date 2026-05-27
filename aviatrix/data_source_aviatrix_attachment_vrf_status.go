package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixAttachmentVrfStatus() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixAttachmentVrfStatusRead,

		Schema: map[string]*schema.Schema{
			"gateway1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Source gateway or gateway-group name. Required unless `all = true`.",
			},
			"gateway2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Destination gateway or gateway-group name. Required unless `all = true`.",
			},
			"all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, operate on all eligible peerings on the controller. `gateway1`/`gateway2` must be empty.",
			},
			"enable_vrf_attachment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringInSlice([]string{"", "yes", "no"}, false),
				Description:  "Set to `yes` or `no` to toggle vrf_attachment_enabled before reading status. Empty (default) reads status without mutating.",
			},
			"attachments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Per-peering VRF attachment status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peering_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`transit-transit` or `edge-spoke-transit`.",
						},
						"vrf_attachment_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixAttachmentVrfStatusRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	all := getBool(d, "all")
	gw1 := getString(d, "gateway1")
	gw2 := getString(d, "gateway2")
	enable := getString(d, "enable_vrf_attachment")

	if all {
		if gw1 != "" || gw2 != "" {
			return diag.Errorf("`gateway1` and `gateway2` must be empty when `all = true`")
		}
	} else {
		if gw1 == "" || gw2 == "" {
			return diag.Errorf("`gateway1` and `gateway2` are required unless `all = true`")
		}
		resolved1, err := resolveGatewayName(ctx, client, gw1)
		if err != nil {
			return diag.Errorf("could not resolve gateway1 %q: %v", gw1, err)
		}
		resolved2, err := resolveGatewayName(ctx, client, gw2)
		if err != nil {
			return diag.Errorf("could not resolve gateway2 %q: %v", gw2, err)
		}
		gw1, gw2 = resolved1, resolved2
	}

	if enable == "yes" || enable == "no" {
		if err := client.UpdateVrfOnAttachment(ctx, gw1, gw2, enable == "yes", all); err != nil {
			return diag.Errorf("failed to update vrf_on_attachment: %v", err)
		}
	}

	attachments, err := client.GetAttachmentVrfStatus(ctx, gw1, gw2, all)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get attachment vrf status: %v", err)
	}

	result := make([]map[string]any, 0, len(attachments))
	for _, a := range attachments {
		result = append(result, map[string]any{
			"peering_name":           a.PeeringName,
			"gateway1":               a.Gateway1,
			"gateway2":               a.Gateway2,
			"attachment_type":        a.AttachmentType,
			"vrf_attachment_enabled": a.VrfAttachmentEnabled,
		})
	}
	if err := d.Set("attachments", result); err != nil {
		return diag.Errorf("failed to set attachments: %v", err)
	}

	if all {
		d.SetId("all")
	} else {
		d.SetId(gw1 + "~" + gw2)
	}
	return nil
}
