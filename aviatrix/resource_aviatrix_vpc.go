package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVpcCreate,
		Read:   resourceAviatrixVpcRead,
		Update: resourceAviatrixVpcUpdate,
		Delete: resourceAviatrixVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
						"ipv6_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 CIDR of the subnet.",
						},
						"ipv6_access_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: ValidateIPv6AccessType,
							Description:  "IPv6 access type for the subnet: \"INTERNAL\" or \"EXTERNAL\". Only supported for GCP (4).",
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
						"ipv6_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 CIDR of the subnet.",
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
						"ipv6_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 CIDR of the subnet.",
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
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable IPv6 for the VPC. Only supported for AWS (1), Azure (8), GCP (4).",
			},
			"vpc_ipv6_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "IPv6 CIDR for the VPC. Required when enable_ipv6 is true for Azure (8). Optional for GCP (4).",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := mustString(val)
					ip, ipnet, err := net.ParseCIDR(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid IPv6 CIDR, got: %s", key, v))
						return
					}

					// Ensure IPv6
					if ip.To4() != nil {
						errs = append(errs, fmt.Errorf("%q must be an IPv6 CIDR, got IPv4: %s", key, v))
						return
					}

					// Ensure it's the network address, not a host address
					if !ip.Equal(ipnet.IP) {
						errs = append(errs, fmt.Errorf("%q must be a network CIDR, not a host IP (%s)", key, v))
					}

					// Reject /128 (single host)
					ones, _ := ipnet.Mask.Size()
					if ones == 128 {
						errs = append(errs, fmt.Errorf("%q cannot be /128, must be a valid IPv6 network range", key))
					}

					return
				},
			},
			"ipv6_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: ValidateIPv6AccessType,
				Description:  "IPv6 access type for the VPC: \"INTERNAL\" or \"EXTERNAL\". Only supported for GCP (4).",
			},
		},
	}
}

func resourceAviatrixVpcCreate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	vpc := &goaviatrix.Vpc{
		CloudType:              getInt(d, "cloud_type"),
		AccountName:            getString(d, "account_name"),
		Region:                 getString(d, "region"),
		Name:                   getString(d, "name"),
		Cidr:                   getString(d, "cidr"),
		SubnetSize:             getInt(d, "subnet_size"),
		NumOfSubnetPairs:       getInt(d, "num_of_subnet_pairs"),
		EnablePrivateOobSubnet: getBool(d, "enable_private_oob_subnet"),
	}
	if vpc.Region == "" && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("please specify 'region'")
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

	aviatrixTransitVpc := getBool(d, "aviatrix_transit_vpc")
	aviatrixFireNetVpc := getBool(d, "aviatrix_firenet_vpc")

	if aviatrixTransitVpc && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		return fmt.Errorf("currently 'aviatrix_transit_vpc' is only supported by AWS (1), AWSGov (256), AWSChina (1024) and Alibaba Cloud (8192)")
	}
	if aviatrixFireNetVpc && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		return fmt.Errorf("currently'aviatrix_firenet_vpc' is only supported by AWS (1), Azure (8), OCI (16), AWSGov (256), AWSChina (1024) and AzureChina (2048)")
	}
	if aviatrixTransitVpc && aviatrixFireNetVpc {
		return fmt.Errorf("vpc cannot be aviatrix transit vpc and aviatrix firenet vpc at the same time")
	}

	nativeGwlb := getBool(d, "enable_native_gwlb")
	if nativeGwlb && !goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AWS) {
		return fmt.Errorf("'enable_native_gwlb' is only valid with cloud_type = 1 (AWS)")
	}

	if aviatrixTransitVpc {
		vpc.AviatrixTransitVpc = "yes"
		log.Printf("[INFO] Creating a new Aviatrix Transit VPC: %#v", vpc)
	} else {
		vpc.AviatrixTransitVpc = "no"
	}
	if aviatrixFireNetVpc {
		vpc.AviatrixFireNetVpc = "yes"
		log.Printf("[INFO] Creating a new Aviatrix FireNet VPC: %#v", vpc)
	} else {
		vpc.AviatrixFireNetVpc = "no"
	}
	if !aviatrixTransitVpc && !aviatrixFireNetVpc {
		log.Printf("[INFO] Creating a new VPC: %#v", vpc)
	}

	ipv6Enabled := getBool(d, "enable_ipv6")

	if goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		if _, ok := d.GetOk("subnets"); ok {
			subnets := getList(d, "subnets")
			for _, subnet := range subnets {
				sub := mustMap(subnet)
				subnetInfo := goaviatrix.SubnetInfo{
					Name:   mustString(sub["name"]),
					Region: mustString(sub["region"]),
					Cidr:   mustString(sub["cidr"]),
				}
				if ipv6Enabled {
					if ipv6Cidr, ok := sub["ipv6_cidr"]; ok {
						subnetInfo.IPv6Cidr = mustString(ipv6Cidr)
					}
					if ipv6AccessType, ok := sub["ipv6_access_type"]; ok {
						subnetInfo.IPv6AccessType = mustString(ipv6AccessType)
					}
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
		vpc.ResourceGroup = mustString(resourceGroup)
	}

	// Handle IPv6 fields
	if ipv6Enabled {
		if err := IPv6SupportedOnCloudType(vpc.CloudType); err != nil {
			return fmt.Errorf("error creating vpc: enable_ipv6 is not supported, %w", err)
		}
		vpc.EnableIpv6 = true
		log.Printf("[INFO] Enabling IPv6 in VPC: %#v", vpc)

		// Handle ipv6_access_type for Azure
		if goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if vpcIpv6Cidr, ok := d.GetOk("vpc_ipv6_cidr"); ok {
				vpc.VpcIpv6Cidr = mustString(vpcIpv6Cidr)
			} else {
				return fmt.Errorf("error creating vpc: valid vpc_ipv6_cidr is required when enable_ipv6 is true for Azure")
			}
		}

		if goaviatrix.IsCloudType(vpc.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			if ipv6AccessType, ok := d.GetOk("ipv6_access_type"); ok {
				vpc.Ipv6AccessType = mustString(ipv6AccessType)
			}
			if vpcIpv6Cidr, ok := d.GetOk("vpc_ipv6_cidr"); ok {
				vpc.VpcIpv6Cidr = mustString(vpcIpv6Cidr)
			}
		}
	}

	err := client.CreateVpc(vpc)
	if err != nil {
		if vpc.AviatrixTransitVpc == "yes" {
			return fmt.Errorf("failed to create a new Aviatrix Transit VPC: %w", err)
		}
		return fmt.Errorf("failed to create a new VPC: %w", err)
	}

	d.SetId(vpc.Name)

	// We intentionally call the Read function early here to load 'vpc_id' for later use
	err = resourceAviatrixVpcRead(d, meta)
	if err != nil {
		return err
	}
	vpc.VpcID = getString(d, "vpc_id")

	if nativeGwlb {
		err = client.EnableNativeAwsGwlbFirenet(vpc)
		if err != nil {
			return fmt.Errorf("could not enable native AWS Gwlb: %w", err)
		}
	}

	return nil
}

func resourceAviatrixVpcRead(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	vpcName := getString(d, "name")
	if vpcName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		mustSet(d, "name", id)
		d.SetId(id)
		return resourceAviatrixVpcRead(d, meta)
	}

	vpc := &goaviatrix.Vpc{
		Name: getString(d, "name"),
	}

	vC, err := client.GetVpc(vpc)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find VPC: %w", err)
	}

	log.Printf("[INFO] Found VPC: %#v", vpc)
	mustSet(d, "cloud_type", vC.CloudType)
	mustSet(d, "account_name", vC.AccountName)
	mustSet(d, "name", vC.Name)
	if !goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "region", vC.Region)
		mustSet(d, "cidr", vC.Cidr)
	}
	if vC.SubnetSize != 0 {
		mustSet(d, "subnet_size", vC.SubnetSize)
	}
	if vC.NumOfSubnetPairs != 0 {
		mustSet(d, "num_of_subnet_pairs", vC.NumOfSubnetPairs)
	}
	mustSet(d, "enable_private_oob_subnet", vC.EnablePrivateOobSubnet)
	if vC.AviatrixTransitVpc == "yes" {
		mustSet(d, "aviatrix_transit_vpc", true)
	} else {
		mustSet(d, "aviatrix_transit_vpc", false)
	}
	if vC.AviatrixFireNetVpc == "yes" {
		mustSet(d, "aviatrix_firenet_vpc", true)
	} else {
		mustSet(d, "aviatrix_firenet_vpc", false)
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(
			// d.Set("vpc_id", strings.Split(vC.VpcID, "~-~")[0])
			d, "vpc_id", vC.VpcID)
	} else {
		mustSet(d, "vpc_id", vC.VpcID)
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		account := &goaviatrix.Account{
			AccountName: getString(d, "account_name"),
		}

		acc, err := client.GetAccount(account)
		if err != nil {
			if !errors.Is(err, goaviatrix.ErrNotFound) {
				return fmt.Errorf("aviatrix Account: %w", err)
			}
		}

		subscriptionId := acc.ArmSubscriptionId

		var resourceGroup string
		if vC.ResourceGroup != "" {
			resourceGroup = vC.ResourceGroup
		} else {
			resourceGroup = strings.Split(vC.VpcID, ":")[1]
		}

		azureVnetResourceId := "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vC.Name
		mustSet(d, "resource_group", resourceGroup)
		mustSet(d, "azure_vnet_resource_id", azureVnetResourceId)
	}

	subnetsForState := buildSubnetsForState(vC.CloudType, vC.Subnets, getList(d, "subnets"))
	if err := d.Set("subnets", subnetsForState); err != nil {
		log.Printf("[WARN] Error setting 'subnets' for (%s): %s", d.Id(), err)
	}

	var privateSubnets []map[string]any
	for _, subnet := range vC.PrivateSubnets {
		subnetInfo := make(map[string]any)
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		subnetInfo["subnet_id"] = subnet.SubnetID
		subnetInfo["ipv6_cidr"] = subnet.IPv6Cidr

		privateSubnets = append(privateSubnets, subnetInfo)
	}
	if err := d.Set("private_subnets", privateSubnets); err != nil {
		log.Printf("[WARN] Error setting 'private_subnets' for (%s): %s", d.Id(), err)
	}

	var publicSubnets []map[string]any
	for _, subnet := range vC.PublicSubnets {
		subnetInfo := make(map[string]any)
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		subnetInfo["subnet_id"] = subnet.SubnetID
		subnetInfo["ipv6_cidr"] = subnet.IPv6Cidr

		publicSubnets = append(publicSubnets, subnetInfo)
	}
	if err := d.Set("public_subnets", publicSubnets); err != nil {
		log.Printf("[WARN] Error setting 'public_subnets' for (%s): %s", d.Id(), err)
	}

	d.SetId(vC.Name)

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.AWS) {
		firenetDetail, err := client.GetFireNet(&goaviatrix.FireNet{VpcID: vC.VpcID})
		if errors.Is(err, goaviatrix.ErrNotFound) {
			mustSet(d, "enable_native_gwlb", false)
		} else if err != nil {
			return fmt.Errorf("could not get FireNet details to read enable_native_gwlb: %w", err)
		} else {
			mustSet(d, "enable_native_gwlb", firenetDetail.NativeGwlb)
		}
	} else {
		mustSet(d, "enable_native_gwlb", false)
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		var rtbs []string
		rtbs, err = getAllRouteTables(vpc, client)
		if err != nil {
			return fmt.Errorf("could not get vpc route table ids: %w", err)
		}

		if err := d.Set("route_tables", rtbs); err != nil {
			log.Printf("[WARN] Error setting route tables for (%s): %s", d.Id(), err)
		}
	} else {
		mustSet(d, "route_tables", []string{})
	}

	if goaviatrix.IsCloudType(vC.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		availabilityDomains, err := client.ListOciVpcAvailabilityDomains(vC)
		if err != nil {
			return fmt.Errorf("could not get OCI availability domains: %w", err)
		}
		mustSet(d, "availability_domains", availabilityDomains)

		faultDomains, err := client.ListOciVpcFaultDomains(vC)
		if err != nil {
			return fmt.Errorf("could not get OCI fault domains: %w", err)
		}
		mustSet(d, "fault_domains", faultDomains)
	} else {
		mustSet(d, "availability_domains", nil)
		mustSet(d, "fault_domains", nil)
	}
	mustSet(d, "enable_ipv6", vC.EnableIpv6)
	if vC.VpcIpv6Cidr != "" {
		mustSet(d, "vpc_ipv6_cidr", vC.VpcIpv6Cidr)
	}

	return nil
}

// buildSubnetsForState converts the controller's view of a VPC's subnets into
// the shape Terraform stores in the "subnets" attribute.
//
// "subnets" is a TypeList, which Terraform compares element by element by
// index. The controller does not guarantee it returns subnets in the order
// they were declared in the configuration (for GCP it commonly groups them by
// region), so echoing the controller's order back into state produces a
// spurious diff that force-replaces the VPC on every plan. To avoid this,
// subnets present in the configuration are emitted first, in their configured
// order; any remaining subnets (e.g. after an import, where there is no prior
// configuration) follow in the controller's order.
func buildSubnetsForState(cloudType int, apiSubnets []goaviatrix.SubnetInfo, configuredSubnets []any) []map[string]any {
	isGCP := goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes)

	// Index the controller's subnets by their ForceNew identity (region, cidr
	// and name uniquely identify a subnet within a VPC) so they can be looked
	// up in configured order.
	subnetsByKey := make(map[string]map[string]any, len(apiSubnets))
	orderedKeys := make([]string, 0, len(apiSubnets))
	for _, subnet := range apiSubnets {
		subnetInfo := make(map[string]any)
		if isGCP {
			subnetInfo["region"] = subnet.Region
		}
		subnetInfo["cidr"] = subnet.Cidr
		subnetInfo["name"] = subnet.Name
		if !isGCP {
			subnetInfo["subnet_id"] = subnet.SubnetID
		}
		subnetInfo["ipv6_cidr"] = subnet.IPv6Cidr
		subnetInfo["ipv6_access_type"] = subnet.IPv6AccessType

		key := subnetIdentityKey(subnet.Region, subnet.Cidr, subnet.Name)
		subnetsByKey[key] = subnetInfo
		orderedKeys = append(orderedKeys, key)
	}

	var subnetsForState []map[string]any

	// Emit configured subnets first, in configured order, so the TypeList
	// indices line up with the config and Terraform sees no diff. Only GCP
	// accepts user-defined subnets; other clouds leave "subnets" computed.
	if isGCP {
		for _, subnet := range configuredSubnets {
			sub := mustMap(subnet)
			key := subnetIdentityKey(mustString(sub["region"]), mustString(sub["cidr"]), mustString(sub["name"]))
			if subnetInfo, ok := subnetsByKey[key]; ok {
				subnetsForState = append(subnetsForState, subnetInfo)
				delete(subnetsByKey, key)
			}
		}
	}

	// Append any subnets the controller returned that were not in the config
	// (e.g. on import), preserving the controller's order for a stable result.
	for _, key := range orderedKeys {
		if subnetInfo, ok := subnetsByKey[key]; ok {
			subnetsForState = append(subnetsForState, subnetInfo)
			delete(subnetsByKey, key)
		}
	}

	return subnetsForState
}

func subnetIdentityKey(region, cidr, name string) string {
	return region + "~" + cidr + "~" + name
}

func resourceAviatrixVpcUpdate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	vpc := &goaviatrix.Vpc{
		AccountName: getString(d, "account_name"),
		VpcID:       getString(d, "vpc_id"),
		Region:      getString(d, "region"),
	}

	if d.HasChange("enable_native_gwlb") {
		nativeGwlb := getBool(d, "enable_native_gwlb")
		if nativeGwlb {
			err := client.EnableNativeAwsGwlbFirenet(vpc)
			if err != nil {
				return fmt.Errorf("could not enable native AWS gwlb firenet: %w", err)
			}
		} else {
			err := client.DisableNativeAwsGwlbFirenet(vpc)
			if err != nil {
				return fmt.Errorf("could not disable native AWS gwlb firenet: %w", err)
			}
		}
	}

	return nil
}

func resourceAviatrixVpcDelete(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	vpc := &goaviatrix.Vpc{
		AccountName: getString(d, "account_name"),
		Name:        getString(d, "name"),
		VpcID:       getString(d, "vpc_id"),
	}

	log.Printf("[INFO] Deleting VPC: %#v", vpc)

	if getBool(d, "enable_native_gwlb") {
		err := client.DisableNativeAwsGwlbFirenet(vpc)
		if err != nil {
			return fmt.Errorf("could not disable native AWS gwlb: %w", err)
		}
	}

	err := client.DeleteVpc(vpc)
	if err != nil {
		return fmt.Errorf("failed to delete VPC: %w", err)
	}

	return nil
}
