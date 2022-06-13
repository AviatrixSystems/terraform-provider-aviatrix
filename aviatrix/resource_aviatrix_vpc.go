package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVpcCreate,
		Read:   resourceAviatrixVpcRead,
		Update: resourceAviatrixVpcUpdate,
		Delete: resourceAviatrixVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixVpcMigrateState,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "Type of cloud service provider.",
				ValidateFunc: validateCloudType,
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
			"enable_private_oob_subnet": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Switch to enable private oob subnet. Only supported for AWS/AWSGov provider. Valid values: true, false. Default value: false.",
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
			"enable_native_gwlb": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable Native AWS GWLB for FireNet Function. Only valid with cloud_type = 1 (AWS). " +
					"Valid values: true or false. Default value: false. Available as of provider version R2.18+.",
			},
			"private_mode_subnets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Switch to only launch private subnets. Only available when Private Mode is enabled on the Controller.",
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
							ForceNew:    true,
							Description: "Subnet region.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Subnet cidr.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
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
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Resource group of the Azure VPC created.",
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
			"azure_vnet_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure vnet resource ID.",
			},
			"route_tables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of route table ids associated with this VPC.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"availability_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of OCI availability domains.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"fault_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of OCI fault domains.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAviatrixVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		CloudType:              d.Get("cloud_type").(int),
		AccountName:            d.Get("account_name").(string),
		Region:                 d.Get("region").(string),
		Name:                   d.Get("name").(string),
		Cidr:                   d.Get("cidr").(string),
		SubnetSize:             d.Get("subnet_size").(int),
		NumOfSubnetPairs:       d.Get("num_of_subnet_pairs").(int),
		EnablePrivateOobSubnet: d.Get("enable_private_oob_subnet").(bool),
		PrivateModeSubnets:     d.Get("private_mode_subnets").(bool),
	}
	if vpc.Region == "" && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("please specifiy 'region'")
	} else if vpc.Region != "" && vpc.CloudType == goaviatrix.GCP {
		return fmt.Errorf("please specify 'region' in 'subnets' for GCP provider")
	}

	if vpc.Cidr == "" && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("please specify 'cidr'")
	} else if vpc.Cidr != "" && goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("please specify 'cidr' in 'subnets' for GCP provider")
	}
	if vpc.SubnetSize != 0 && vpc.NumOfSubnetPairs != 0 {
		if !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("advanced option('subnet_size' and 'num_of_subnet_pairs') is only supported by AWS (1), Azure (8), AWSGov (256), AWSChina (1024) and AzureChina (2048)")
		}
	} else if vpc.SubnetSize != 0 || vpc.NumOfSubnetPairs != 0 {
		if goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("please specify both 'subnet_size' and 'num_of_subnet_pairs' to enable advanced options")
		} else {
			return fmt.Errorf("advanced option('subnet_size' and 'num_of_subnet_pairs') is only supported by AWS (1), Azure (8), AWSGov (256), AWSChina (1024) and AzureChina (2048)")
		}
	}
	if vpc.EnablePrivateOobSubnet {
		if !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("advanced option('enable_private_oob_subnet') is only supported by AWS (1), AWSGov (256), and AWSChina (1024)")
		}
	}

	aviatrixTransitVpc := d.Get("aviatrix_transit_vpc").(bool)
	aviatrixFireNetVpc := d.Get("aviatrix_firenet_vpc").(bool)

	if aviatrixTransitVpc && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		return fmt.Errorf("currently 'aviatrix_transit_vpc' is only supported by AWS (1), AWSGov (256), AWSChina (1024) and Alibaba Cloud (8192)")
	}
	if aviatrixFireNetVpc && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		return fmt.Errorf("currently'aviatrix_firenet_vpc' is only supported by AWS (1), Azure (8), OCI (16), AWSGov (256), AWSChina (1024) and AzureChina (2048)")
	}
	if aviatrixTransitVpc && aviatrixFireNetVpc {
		return fmt.Errorf("vpc cannot be aviatrix transit vpc and aviatrix firenet vpc at the same time")
	}

	nativeGwlb := d.Get("enable_native_gwlb").(bool)
	if nativeGwlb && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWS) {
		return fmt.Errorf("'enable_native_gwlb' is only valid with cloud_type = 1 (AWS)")
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

	if goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
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

	if resourceGroup, ok := d.GetOk("resource_group"); ok {
		if !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("error creating vpc: resource_group is required to be empty for providers other than Azure (8), AzureGov (32) and AzureChina (2048)")
		}
		vpc.ResourceGroup = resourceGroup.(string)
	}

	err := client.CreateVpc(vpc)
	if err != nil {
		if vpc.AviatrixTransitVpc == "yes" {
			return fmt.Errorf("failed to create a new Aviatrix Transit VPC: %s", err)
		}
		return fmt.Errorf("failed to create a new VPC: %s", err)
	}

	d.SetId(vpc.Name)

	// We intentionally call the Read function early here to load 'vpc_id' for later use
	err = resourceAviatrixVpcRead(d, meta)
	if err != nil {
		return err
	}
	vpc.VpcID = d.Get("vpc_id").(string)

	if nativeGwlb {
		err = client.EnableNativeAwsGwlbFirenet(vpc)
		if err != nil {
			return fmt.Errorf("could not enable native AWS Gwlb: %v", err)
		}
	}

	return nil
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
	if !goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("region", vC.Region)
		d.Set("cidr", vC.Cidr)
	}
	if vC.SubnetSize != 0 {
		d.Set("subnet_size", vC.SubnetSize)
	}
	if vC.NumOfSubnetPairs != 0 {
		d.Set("num_of_subnet_pairs", vC.NumOfSubnetPairs)
	}
	d.Set("enable_private_oob_subnet", vC.EnablePrivateOobSubnet)
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

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		//d.Set("vpc_id", strings.Split(vC.VpcID, "~-~")[0])
		d.Set("vpc_id", vC.VpcID)
	} else {
		d.Set("vpc_id", vC.VpcID)
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
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

		var resourceGroup string
		if vC.ResourceGroup != "" {
			resourceGroup = vC.ResourceGroup
		} else {
			resourceGroup = strings.Split(vC.VpcID, ":")[1]
		}

		azureVnetResourceId := "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vC.Name
		d.Set("resource_group", resourceGroup)
		d.Set("azure_vnet_resource_id", azureVnetResourceId)
	}

	subnetsMap := make(map[string]map[string]interface{})
	var subnetsKeyArray []string
	for _, subnet := range vC.Subnets {
		subnetInfo := make(map[string]interface{})
		if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			subnetInfo["region"] = subnet.Region
		}
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		if !goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			subnetInfo["subnet_id"] = subnet.SubnetID
		}

		var key string
		if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			key = subnet.Region + "~" + subnet.Cidr + "~" + subnet.Name
		} else {
			key = subnet.Cidr + "~" + subnet.Name + "~" + subnet.SubnetID
		}
		subnetsMap[key] = subnetInfo
		subnetsKeyArray = append(subnetsKeyArray, key)
	}

	var subnetsFromFile []map[string]interface{}
	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
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

	d.SetId(vC.Name)

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.AWS) {
		firenetDetail, err := client.GetFireNet(&goaviatrix.FireNet{VpcID: vC.VpcID})
		if err == goaviatrix.ErrNotFound {
			d.Set("enable_native_gwlb", false)
		} else if err != nil {
			return fmt.Errorf("could not get FireNet details to read enable_native_gwlb: %v", err)
		} else {
			d.Set("enable_native_gwlb", firenetDetail.NativeGwlb)
		}
	} else {
		d.Set("enable_native_gwlb", false)
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		var rtbs []string
		rtbs, err = getAllRouteTables(vpc, client)

		if err != nil {
			return fmt.Errorf("could not get vpc route table ids: %v", err)
		}

		if err := d.Set("route_tables", rtbs); err != nil {
			log.Printf("[WARN] Error setting route tables for (%s): %s", d.Id(), err)
		}
	} else {
		d.Set("route_tables", []string{})
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		availabilityDomains, err := client.ListOciVpcAvailabilityDomains(vC)
		if err != nil {
			return fmt.Errorf("could not get OCI availability domains: %v", err)
		}
		d.Set("availability_domains", availabilityDomains)

		faultDomains, err := client.ListOciVpcFaultDomains(vC)
		if err != nil {
			return fmt.Errorf("could not get OCI fault domains: %v", err)
		}
		d.Set("fault_domains", faultDomains)
	} else {
		d.Set("availability_domains", nil)
		d.Set("fault_domains", nil)
	}

	d.Set("private_mode_subnets", vC.PrivateModeSubnets)

	return nil
}

func resourceAviatrixVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		AccountName: d.Get("account_name").(string),
		VpcID:       d.Get("vpc_id").(string),
		Region:      d.Get("region").(string),
	}

	if d.HasChange("enable_native_gwlb") {
		nativeGwlb := d.Get("enable_native_gwlb").(bool)
		if nativeGwlb {
			err := client.EnableNativeAwsGwlbFirenet(vpc)
			if err != nil {
				return fmt.Errorf("could not enable native AWS gwlb firenet: %v", err)
			}
		} else {
			err := client.DisableNativeAwsGwlbFirenet(vpc)
			if err != nil {
				return fmt.Errorf("could not disable native AWS gwlb firenet: %v", err)
			}
		}
	}

	return nil
}

func resourceAviatrixVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpc := &goaviatrix.Vpc{
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("name").(string),
		VpcID:       d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting VPC: %#v", vpc)

	if d.Get("enable_native_gwlb").(bool) {
		err := client.DisableNativeAwsGwlbFirenet(vpc)
		if err != nil {
			return fmt.Errorf("could not disable native AWS gwlb: %v", err)
		}
	}

	err := client.DeleteVpc(vpc)
	if err != nil {
		return fmt.Errorf("failed to delete VPC: %s", err)
	}

	return nil
}
