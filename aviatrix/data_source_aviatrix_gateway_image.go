package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixGatewayImage() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixGatewayImageRead,

		Schema: map[string]*schema.Schema{
			"software_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Software version.",
			},
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Type of cloud service provider.",
				ValidateFunc: validateCloudType,
			},
			"image_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compatible image version for the given software_version.",
			},
		},
	}
}

func dataSourceAviatrixGatewayImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	cloudType := getInt(d, "cloud_type")
	softwareVersion := getString(d, "software_version")
	v, err := client.GetCompatibleImageVersion(ctx, cloudType, softwareVersion)
	if err != nil {
		return diag.Errorf("could not get compatible image version: %v", err)
	}
	mustSet(d, "image_version", v)
	d.SetId(v)
	return nil
}
