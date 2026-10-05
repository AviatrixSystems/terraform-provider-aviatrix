package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixControllerMetadata() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixControllerMetadataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud type.",
			},
		},
	}
}

func dataSourceAviatrixControllerMetadataRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	controllerMetadata, err := client.GetControllerMetadata(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get controller metadata: %s", err)
	}

	d.Set("region", controllerMetadata.Region)
	d.Set("vpc_id", controllerMetadata.VpcId)
	d.Set("instance_id", controllerMetadata.InstanceId)
	d.Set("cloud_type", controllerMetadata.CloudType)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}
