package aviatrix

import (
	"context"
	"errors"
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
	client := mustClient(meta)

	controllerMetadata, err := client.GetControllerMetadata(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get controller metadata: %s", err)
	}
	mustSet(d, "region", controllerMetadata.Region)
	mustSet(d, "vpc_id", controllerMetadata.VpcId)
	mustSet(d, "instance_id", controllerMetadata.InstanceId)
	mustSet(d, "cloud_type", controllerMetadata.CloudType)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}
