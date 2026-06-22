package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixDcfAttachmentPoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDcfAttachmentPointsRead,

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

func dataSourceAviatrixDcfAttachmentPointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	name := getString(d, "name")

	attachmentPoint, err := client.GetDCFAttachmentPoint(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(attachmentPoint.AttachmentPointID)
	mustSet(d, "name", name)
	mustSet(d, "attachment_point_id", attachmentPoint.AttachmentPointID)
	return nil
}
