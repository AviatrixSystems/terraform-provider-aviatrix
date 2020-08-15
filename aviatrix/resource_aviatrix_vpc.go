package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVpcCreate,
		Read:   resourceAviatrixVpcRead,
		Delete: resourceAviatrixVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixVpcMigrateState,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account name. This account will be used to create an Aviatrix VPC.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the VPC to be created.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Region of cloud provider. Required to be empty for GCP provider, and non-empty for other providers.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Subnet of the VPC to be created. Required to be empty for GCP provider, and non-empty for other providers.",
			},
			"subnet_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Subnet size.",
			},
			"num_of_subnet_pairs": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of public subnet and private subnet pair to be created.",
			},
			"aviatrix_transit_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specify the VPC as Aviatrix Transit VPC or not. Required to be false for GCP provider.",
			},
			"aviatrix_firenet_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specify the VPC as Aviatrix FireNet VPC or not. Required to be false for GCP provider.",
			},
			"subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "List of subnet of the VPC to be created. Required to be non-empty for GCP provider, and empty for other providers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet region.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet cidr.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
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
			"private_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of private subnet of the VPC to be created.",
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
		},
	}
}

func resourceAviatrixVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		CloudType:        d.Get("cloud_type").(int),
		AccountName:      d.Get("account_name").(string),
		Region:           d.Get("region").(string),
		Name:             d.Get("name").(string),
		Cidr:             d.Get("cidr").(string),
		SubnetSize:       d.Get("subnet_size").(int),
		NumOfSubnetPairs: d.Get("num_of_subnet_pairs").(int),
	}
	if vpc.Region == "" && vpc.CloudType != goaviatrix.GCP {
		return fmt.Errorf("please specifiy 'region'")
	} else if vpc.Region != "" && vpc.CloudType == goaviatrix.GCP {
		return fmt.Errorf("please specify 'region' in 'subnets' for GCP provider")
	}

	if vpc.Cidr == "" && vpc.CloudType != goaviatrix.GCP {
		return fmt.Errorf("please specify 'cidr'")
	} else if vpc.Cidr != "" && vpc.CloudType == goaviatrix.GCP {
		return fmt.Errorf("please specify 'cidr' in 'subnets' for GCP provider")
	}

	if vpc.SubnetSize != 0 && vpc.NumOfSubnetPairs != 0 {
		if vpc.CloudType != goaviatrix.AWS && vpc.CloudType != goaviatrix.AZURE {
			return fmt.Errorf("advanced option('subnet_size' and 'num_of_subnet_pairs') is only supported for AWS and Azure provider")
		}
	} else if vpc.SubnetSize != 0 || vpc.NumOfSubnetPairs != 0 {
		if vpc.CloudType == goaviatrix.AWS || vpc.CloudType == goaviatrix.AZURE {
			return fmt.Errorf("please specify both 'subnet_size' and 'num_of_subnet_pairs' to enable advanced options")
		} else {
			return fmt.Errorf("advanced option('subnet_size' and 'num_of_subnet_pairs') is only supported for AWS and Azure provider")
		}
	}

	aviatrixTransitVpc := d.Get("aviatrix_transit_vpc").(bool)
	aviatrixFireNetVpc := d.Get("aviatrix_firenet_vpc").(bool)

	if aviatrixTransitVpc && vpc.CloudType != goaviatrix.AWS {
		return fmt.Errorf("currently 'aviatrix_transit_vpc' is only supported for AWS provider")
	}
	if aviatrixFireNetVpc && vpc.CloudType != goaviatrix.AWS && vpc.CloudType != goaviatrix.AZURE {
		return fmt.Errorf("currently'aviatrix_firenet_vpc' is only supported for AWS and AZURE provider")
	}
	if aviatrixTransitVpc && aviatrixFireNetVpc {
		return fmt.Errorf("vpc cannot be aviatrix transit vpc and aviatrix firenet vpc at the same time")
	}
	if aviatrixTransitVpc {
		vpc.AviatrixTransitVpc = "yes"
		log.Printf("[INFO] Creating a new Aviatrix Transit VPC: %#v", vpc)
	}
	if aviatrixFireNetVpc {
		vpc.AviatrixFireNetVpc = "yes"
		log.Printf("[INFO] Creating a new Aviatrix FireNet VPC: %#v", vpc)
	}
	if !aviatrixTransitVpc && !aviatrixFireNetVpc {
		log.Printf("[INFO] Creating a new VPC: %#v", vpc)
	}

	if vpc.CloudType == goaviatrix.GCP {
		if _, ok := d.GetOk("subnets"); ok {
			subnets := d.Get("subnets").([]interface{})
			for _, subnet := range subnets {
				sub := subnet.(map[string]interface{})
				subnetInfo := goaviatrix.SubnetInfo{
					Name:   sub["name"].(string),
					Region: sub["region"].(string),
					Cidr:   sub["cidr"].(string),
				}
				vpc.Subnets = append(vpc.Subnets, subnetInfo)
			}
		} else {
			return fmt.Errorf("subnets is required to be non-empty for GCP provider")
		}
	} else if _, ok := d.GetOk("subnets"); ok {
		return fmt.Errorf("subnets is required to be empty for providers other than GCP")
	}

	err := client.CreateVpc(vpc)
	if err != nil {
		if vpc.AviatrixTransitVpc == "yes" {
			return fmt.Errorf("failed to create a new Aviatrix Transit VPC: %s", err)
		}
		return fmt.Errorf("failed to create a new VPC: %s", err)
	}

	d.SetId(vpc.Name)
	return resourceAviatrixVpcRead(d, meta)
}

func resourceAviatrixVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcName := d.Get("name").(string)
	if vpcName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		d.Set("name", id)
		d.SetId(id)
		return resourceAviatrixVpcRead(d, meta)
	}

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

	log.Printf("[INFO] Found VPC: %#v", vpc)

	d.Set("cloud_type", vC.CloudType)
	d.Set("account_name", vC.AccountName)
	d.Set("name", vC.Name)
	if vC.CloudType != goaviatrix.GCP {
		d.Set("region", vC.Region)
		d.Set("cidr", vC.Cidr)
	}
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
		if vC.CloudType == goaviatrix.AZURE {
			d.Set("resource_group", strings.Split(vC.VpcID, ":")[1])
		}
	}

	subnetsMap := make(map[string]map[string]interface{})
	var subnetsKeyArray []string
	for _, subnet := range vC.Subnets {
		subnetInfo := make(map[string]interface{})
		if vC.CloudType == goaviatrix.GCP {
			subnetInfo["region"] = subnet.Region
		}
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		if vC.CloudType != goaviatrix.GCP {
			subnetInfo["subnet_id"] = subnet.SubnetID
		}

		var key string
		if vC.CloudType == goaviatrix.GCP {
			key = subnet.Region + "~" + subnet.Cidr + "~" + subnet.Name
		} else {
			key = subnet.Cidr + "~" + subnet.Name + "~" + subnet.SubnetID
		}
		subnetsMap[key] = subnetInfo
		subnetsKeyArray = append(subnetsKeyArray, key)
	}

	var subnetsFromFile []map[string]interface{}
	if vC.CloudType == goaviatrix.GCP {
		subnets := d.Get("subnets").([]interface{})
		for _, subnet := range subnets {
			sub := subnet.(map[string]interface{})
			subnetInfo := &goaviatrix.SubnetInfo{
				Cidr:   sub["cidr"].(string),
				Name:   sub["name"].(string),
				Region: sub["region"].(string),
			}

			key := subnetInfo.Region + "~" + subnetInfo.Cidr + "~" + subnetInfo.Name
			if val, ok := subnetsMap[key]; ok {
				if goaviatrix.CompareMapOfInterface(sub, val) {
					subnetsFromFile = append(subnetsFromFile, sub)
					delete(subnetsMap, key)
				}
			}
		}
	}
	if len(subnetsKeyArray) != 0 {
		for i := 0; i < len(subnetsKeyArray); i++ {
			if subnetsMap[subnetsKeyArray[i]] != nil {
				subnetsFromFile = append(subnetsFromFile, subnetsMap[subnetsKeyArray[i]])
			}
		}
	}

	if err := d.Set("subnets", subnetsFromFile); err != nil {
		log.Printf("[WARN] Error setting 'subnets' for (%s): %s", d.Id(), err)
	}

	var privateSubnets []map[string]interface{}
	for _, subnet := range vC.PrivateSubnets {
		subnetInfo := make(map[string]interface{})
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		subnetInfo["subnet_id"] = subnet.SubnetID

		privateSubnets = append(privateSubnets, subnetInfo)
	}
	if err := d.Set("private_subnets", privateSubnets); err != nil {
		log.Printf("[WARN] Error setting 'private_subnets' for (%s): %s", d.Id(), err)
	}

	var publicSubnets []map[string]interface{}
	for _, subnet := range vC.PublicSubnets {
		subnetInfo := make(map[string]interface{})
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		subnetInfo["subnet_id"] = subnet.SubnetID

		publicSubnets = append(publicSubnets, subnetInfo)
	}
	if err := d.Set("public_subnets", publicSubnets); err != nil {
		log.Printf("[WARN] Error setting 'public_subnets' for (%s): %s", d.Id(), err)
	}

	d.SetId(vpcName)
	return nil
}

func resourceAviatrixVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("name").(string),
	}

	log.Printf("[INFO] Deleting VPC: %#v", vpc)

	err := client.DeleteVpc(vpc)
	if err != nil {
		return fmt.Errorf("failed to delete VPC: %s", err)
	}

	return nil
}
