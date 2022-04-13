package aviatrix

import (
	"fmt"
	"sort"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixFirewallInstanceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFirewallInstanceImagesRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"firewall_images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of firewall instances associated with fireNet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"firewall_image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the firewall image.",
						},
						"firewall_image_version": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Versions of the firewall image.",
						},
						"firewall_size": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Instance sizes of the firewall image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixFirewallInstanceImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Get("vpc_id").(string)

	firewallInstanceImages, err := client.GetFirewallInstanceImages(vpcId)
	if err != nil {
		return fmt.Errorf("couldn't get firewall instance images: %s", err)
	}

	var images []map[string]interface{}
	for _, image := range *firewallInstanceImages {
		fI := make(map[string]interface{})
		fI["firewall_image"] = image.Image
		versionList := image.Version
		sort.Slice(versionList, func(i, j int) bool {
			return sortVersion(versionList, i, j)
		})
		fI["firewall_image_version"] = reverseArray(versionList)
		sizeList := image.Size
		sort.Slice(sizeList, func(i, j int) bool {
			return sortSize(sizeList, i, j)
		})
		fI["firewall_size"] = sizeList
		images = append(images, fI)
	}

	if err = d.Set("firewall_images", images); err != nil {
		return fmt.Errorf("couldn't set firewall_images: %s", err)
	}

	d.SetId(vpcId)
	return nil
}
