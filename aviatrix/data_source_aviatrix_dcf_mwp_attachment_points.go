package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixDcfMwpAttachmentPoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDcfMwpAttachmentPointsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the attachment point.",
			},
			"attachment_point_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the attachment point.",
			},
		},
	}
}

func dataSourceAviatrixDcfMwpAttachmentPointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("client must be of type *goaviatrix.Client")
	}

	name := d.Get("name").(string)

	attachmentPoint, err := client.GetDCFAttachmentPoint(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(attachmentPoint.AttachmentPointID)
	d.Set("name", name)
	return nil
}
