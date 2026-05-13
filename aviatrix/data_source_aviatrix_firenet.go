package aviatrix

import (
	"errors"
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
	client := mustClient(meta)

	fireNet := &goaviatrix.FireNet{
		VpcID: getString(d, "vpc_id"),
	}

	fireNetDetail, err := client.GetFireNet(fireNet)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FireNet: %w", err)
	}
	mustSet(d, "vpc_id", fireNetDetail.VpcID)
	mustSet(d, "hashing_algorithm", fireNetDetail.HashingAlgorithm)
	mustSet(d, "tgw_segmentation_for_egress_enabled", fireNetDetail.TgwSegmentationForEgress == "yes")
	mustSet(d, "egress_static_cidrs", fireNetDetail.EgressStaticCidrs)

	if fireNetDetail.Inspection == "yes" {
		mustSet(d, "inspection_enabled", true)
	} else {
		mustSet(d, "inspection_enabled", false)
	}
	if fireNetDetail.FirewallEgress == "yes" {
		mustSet(d, "egress_enabled", true)
	} else {
		mustSet(d, "egress_enabled", false)
	}

	d.SetId(fireNetDetail.VpcID)
	return nil
}
