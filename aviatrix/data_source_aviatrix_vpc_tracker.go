package aviatrix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixVpcTracker() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixVpcTrackerRead,
		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Get VPCs from a single cloud provider. " +
					"For example, if cloud_type = 4, only GCP VPCs will be returned.",
				ValidateFunc: validation.IntInSlice([]int{1, 4, 8, 16, 256}),
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Get VPCs that match the given CIDR.",
				ValidateFunc: validation.IsCIDR,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Get VPCs that match the given region.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Get VPCs that match the given access account name.",
			},
			"vpc_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of VPCs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VPC cloud type, for example '1' (AWS).",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC id, for example 'vpc-069eb82201c8456e3' (AWS).",
						},
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Aviatrix access account name associated with the VPC.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC region, for example 'us-west-1' (AWS).",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC name, for example 'controller_vpc(vpc-069eb82201c8456e3)' (AWS).",
						},
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC cidr, for example '10.0.0.1/24'. Set for AWS/AZURE only.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of instances running in the VPC.",
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of subnets associated with the VPC, GCP only.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet region.",
									},
									"cidr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet CIDR.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet name",
									},
									"gw_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway ip in the subnet",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixVpcTrackerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vpcTracker, err := client.GetVpcTracker()
	if err != nil {
		return fmt.Errorf("could not get vpc list: %s", err)
	}
	vpcTracker = filterVpcTrackerResult(d, vpcTracker)

	var vpcList []map[string]interface{}
	for _, vpc := range vpcTracker {
		vpcList = append(vpcList, map[string]interface{}{
			"cloud_type":     vpc.CloudType,
			"vpc_id":         vpc.VpcID,
			"account_name":   vpc.AccountName,
			"region":         vpc.Region,
			"name":           vpc.Name,
			"cidr":           vpc.Cidr,
			"instance_count": vpc.InstanceCount,
			"subnets":        vpcTrackerSubnetsToMaps(vpc.Subnets),
		})
	}
	err = d.Set("vpc_list", vpcList)
	if err != nil {
		return fmt.Errorf("could not set vpc list: %s", err)
	}

	// vpc_tracker takes no inputs so always has the same id
	d.SetId("vpc_tracker")

	return nil
}

func filterVpcTrackerResult(d *schema.ResourceData, vpcList []*goaviatrix.VpcTracker) []*goaviatrix.VpcTracker {
	cloudType := d.Get("cloud_type").(int)
	cidr := d.Get("cidr").(string)
	region := d.Get("region").(string)
	accountName := d.Get("account_name").(string)

	var filteredList []*goaviatrix.VpcTracker
	for _, vpc := range vpcList {
		if cloudType != 0 && vpc.CloudType != cloudType {
			continue
		}
		if cidr != "" && vpc.Cidr != cidr {
			continue
		}
		if region != "" && vpc.Region != region {
			continue
		}
		if accountName != "" && vpc.AccountName != accountName {
			continue
		}
		filteredList = append(filteredList, vpc)
	}

	return filteredList
}

func vpcTrackerSubnetsToMaps(s []goaviatrix.VPCTrackerSubnet) []map[string]interface{} {
	var m []map[string]interface{}
	for _, sn := range s {
		m = append(m, map[string]interface{}{
			"region": sn.Region,
			"name":   sn.Name,
			"cidr":   sn.Cidr,
			"gw_ip":  sn.GatewayIP,
		})
	}
	return m
}
