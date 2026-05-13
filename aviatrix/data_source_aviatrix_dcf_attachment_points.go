package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
	d.Set("attachment_point_id", attachmentPoint.AttachmentPointID)
	return nil
}
