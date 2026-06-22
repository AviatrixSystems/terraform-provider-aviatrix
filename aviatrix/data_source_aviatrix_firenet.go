package aviatrix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixFireNet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFireNetRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"inspection_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable traffic inspection.",
			},
			"egress_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable egress through firewall.",
			},
			"hashing_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hashing algorithm to load balance traffic across the firewall.",
			},
			"tgw_segmentation_for_egress_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable TGW segmentation for egress.",
			},
			"egress_static_cidrs": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of egress static cidrs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAviatrixFireNetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fireNet := &goaviatrix.FireNet{
		VpcID: d.Get("vpc_id").(string),
	}

	fireNetDetail, err := client.GetFireNet(fireNet)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FireNet: %s", err)
	}

	d.Set("vpc_id", fireNetDetail.VpcID)
	d.Set("hashing_algorithm", fireNetDetail.HashingAlgorithm)
	d.Set("tgw_segmentation_for_egress_enabled", fireNetDetail.TgwSegmentationForEgress == "yes")
	d.Set("egress_static_cidrs", fireNetDetail.EgressStaticCidrs)

	if fireNetDetail.Inspection == "yes" {
		d.Set("inspection_enabled", true)
	} else {
		d.Set("inspection_enabled", false)
	}
	if fireNetDetail.FirewallEgress == "yes" {
		d.Set("egress_enabled", true)
	} else {
		d.Set("egress_enabled", false)
	}

	d.SetId(fireNetDetail.VpcID)
	return nil
}
