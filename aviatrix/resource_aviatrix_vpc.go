package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region where this gateway will be launched.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the VPC to be created.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet of the VPC to be created.",
			},
			"aviatrix_transit_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specify the VPC as Aviatrix Transit VPC or not.",
			},
			"aviatrix_firenet_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specify the VPC as Aviatrix FireNet VPC or not.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VPC created.",
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subnet of the VPC to be created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Optional:    true,
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
							Optional:    true,
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
							Optional:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		Name:        d.Get("name").(string),
		Cidr:        d.Get("cidr").(string),
	}

	if vpc.Region == "" {
		return fmt.Errorf("region can not be empty")
	}

	aviatrixTransitVpc := d.Get("aviatrix_transit_vpc").(bool)
	aviatrixFireNetVpc := d.Get("aviatrix_firenet_vpc").(bool)

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
	d.Set("region", vC.Region)
	d.Set("name", vC.Name)
	d.Set("cidr", vC.Cidr)
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

	d.Set("vpc_id", vC.VpcID)

	var subnetList []map[string]string
	var publicSubnetList []map[string]string
	var privateSubnetList []map[string]string
	for _, subnet := range vC.Subnets {
		sub := make(map[string]string)
		sub["cidr"] = subnet.Cidr
		sub["name"] = subnet.Name
		sub["subnet_id"] = subnet.SubnetID

		subnetList = append(subnetList, sub)
		if strings.Contains(subnet.Name, "Public") {
			publicSubnetList = append(publicSubnetList, sub)
		}
		if strings.Contains(subnet.Name, "Private") {
			privateSubnetList = append(privateSubnetList, sub)
		}
	}

	if err := d.Set("subnets", subnetList); err != nil {
		log.Printf("[WARN] Error setting subnets for (%s): %s", d.Id(), err)
	}
	if err := d.Set("public_subnets", publicSubnetList); err != nil {
		log.Printf("[WARN] Error setting public subnets for (%s): %s", d.Id(), err)
	}
	if err := d.Set("private_subnets", privateSubnetList); err != nil {
		log.Printf("[WARN] Error setting private subnets for (%s): %s", d.Id(), err)
	}

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
