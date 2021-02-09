package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixVpcRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the VPC created.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account name of the VPC created.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the VPC created.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet of the VPC created.",
			},
			"subnet_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Subnet size.",
			},
			"num_of_subnet_pairs": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of public subnet and private subnet pair created.",
			},
			"aviatrix_transit_vpc": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Switch if the VPC created is an Aviatrix Transit VPC or not.",
			},
			"aviatrix_firenet_vpc": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Switch if the VPC created is an Aviatrix FireNet VPC or not.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VPC created.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group of the Azure VPC created.",
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subnet of the VPC to be created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet cidr.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},
			"public_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of public subnet of the VPC to be created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public subnet cidr.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public subnet name.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public subnet ID.",
						},
					},
				},
			},
			"private_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of private subnet of the VPC to be created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private subnet cidr.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private subnet name.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private subnet ID.",
						},
					},
				},
			},
			"route_tables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of AWS route table ids associated with this VPC. Only populated for AWS vpc.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"route_tables_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Description: "Filters the route_tables list to contain only public or private route tables. " +
					"Valid values are 'private' or 'public'. If not set then route_tables are not filtered.",
			},
			"azure_vnet_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure vnet resource ID.",
			},
		},
	}
}

func dataSourceAviatrixVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		Name: d.Get("name").(string),
	}

	vC, err := client.GetVpc(vpc)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find VPC: %s", err)
	}

	d.Set("cloud_type", vC.CloudType)
	d.Set("account_name", vC.AccountName)
	d.Set("region", vC.Region)
	d.Set("name", vC.Name)
	d.Set("cidr", vC.Cidr)
	if vC.SubnetSize != 0 {
		d.Set("subnet_size", vC.SubnetSize)
	}
	if vC.NumOfSubnetPairs != 0 {
		d.Set("num_of_subnet_pairs", vC.NumOfSubnetPairs)
	}
	if vC.AviatrixTransitVpc == "yes" {
		d.Set("aviatrix_transit_vpc", true)
	} else {
		d.Set("aviatrix_transit_vpc", false)
	}

	if vC.AviatrixFireNetVpc == "yes" {
		d.Set("aviatrix_firenet_vpc", true)
	} else {
		d.Set("aviatrix_firenet_vpc", false)
	}

	if vC.CloudType == goaviatrix.GCP {
		d.Set("vpc_id", strings.Split(vC.VpcID, "~-~")[0])
	} else {
		d.Set("vpc_id", vC.VpcID)
	}

	if vC.CloudType == goaviatrix.AZURE {
		account := &goaviatrix.Account{
			AccountName: d.Get("account_name").(string),
		}

		acc, err := client.GetAccount(account)
		if err != nil {
			if err != goaviatrix.ErrNotFound {
				return fmt.Errorf("aviatrix Account: %s", err)
			}
		}

		var subscriptionId string
		if acc != nil {
			subscriptionId = acc.ArmSubscriptionId
		}

		resourceGroup := strings.Split(vC.VpcID, ":")[1]
		azureVnetResourceId := "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vC.Name
		d.Set("resource_group", resourceGroup)
		d.Set("azure_vnet_resource_id", azureVnetResourceId)
	}

	var subnetList []map[string]string
	for _, subnet := range vC.Subnets {
		sub := make(map[string]string)
		sub["cidr"] = subnet.Cidr
		sub["name"] = subnet.Name
		sub["subnet_id"] = subnet.SubnetID

		subnetList = append(subnetList, sub)
	}
	if err := d.Set("subnets", subnetList); err != nil {
		log.Printf("[WARN] Error setting subnets for (%s): %s", d.Id(), err)
	}

	var privateSubnetList []map[string]string
	var publicSubnetList []map[string]string
	for _, subnet := range vC.PrivateSubnets {
		sub := make(map[string]string)
		sub["cidr"] = subnet.Cidr
		sub["name"] = subnet.Name
		sub["subnet_id"] = subnet.SubnetID

		privateSubnetList = append(privateSubnetList, sub)
	}
	for _, subnet := range vC.PublicSubnets {
		sub := make(map[string]string)
		sub["cidr"] = subnet.Cidr
		sub["name"] = subnet.Name
		sub["subnet_id"] = subnet.SubnetID

		publicSubnetList = append(publicSubnetList, sub)
	}
	if err := d.Set("private_subnets", privateSubnetList); err != nil {
		log.Printf("[WARN] Error setting 'private_subnets' for (%s): %s", d.Id(), err)
	}
	if err := d.Set("public_subnets", publicSubnetList); err != nil {
		log.Printf("[WARN] Error setting 'public_subnets' for (%s): %s", d.Id(), err)
	}

	if vC.CloudType == goaviatrix.AWS || vC.CloudType == goaviatrix.AWSGOV || vC.CloudType == goaviatrix.AZURE {
		var rtbs []string
		routeTableFilter := d.Get("route_tables_filter")
		if routeTableFilter == "private" {
			rtbs, err = getPrivateRouteTables(vpc, client)
		} else if routeTableFilter == "public" {
			rtbs, err = getPublicRouteTables(vpc, client)
		} else {
			rtbs, err = getAllRouteTables(vpc, client)
		}

		if err != nil {
			return fmt.Errorf("could not get vpc route table ids: %v", err)
		}

		if err := d.Set("route_tables", rtbs); err != nil {
			log.Printf("[WARN] Error setting route tables for (%s): %s", d.Id(), err)
		}
	}

	d.SetId(vC.Name)
	return nil
}

// To find all the private route tables we will remove the public route tables
// from the list of all route tables.
func getPrivateRouteTables(vpc *goaviatrix.Vpc, client *goaviatrix.Client) ([]string, error) {
	all, err := getAllRouteTables(vpc, client)
	if err != nil {
		return nil, err
	}

	public, err := getPublicRouteTables(vpc, client)
	if err != nil {
		return nil, err
	}

	var rtbs []string
	for _, rtb := range all {
		if !goaviatrix.Contains(public, rtb) {
			rtbs = append(rtbs, rtb)
		}
	}
	return rtbs, nil
}

func getPublicRouteTables(vpc *goaviatrix.Vpc, client *goaviatrix.Client) ([]string, error) {
	vpc.PublicRoutesOnly = true
	rtbs, err := client.GetVpcRouteTableIDs(vpc)
	if err != nil {
		return nil, err
	}
	return rtbs, nil
}

func getAllRouteTables(vpc *goaviatrix.Vpc, client *goaviatrix.Client) ([]string, error) {
	vpc.PublicRoutesOnly = false
	rtbs, err := client.GetVpcRouteTableIDs(vpc)
	if err != nil {
		return nil, err
	}
	return rtbs, nil
}
