package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixSpokeGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixSpokeGatewayRead,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Spoke gateway name. This can be used for getting spoke gateway.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC-ID/VNet-Name of cloud provider.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of cloud provider.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Size of the gateway instance.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public Subnet Info.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specify whether to enable Source NAT feature in 'single_ip' mode on the gateway or not.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Required when allocate_new_eip is false. It uses specified EIP for this gateway.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Subnet. Required if enabling HA for AWS/Azure.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Zone. Required if enabling HA for GCP.",
			},
			"ha_insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Spoke HA Gateway. Required if insane_mode is true and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Gateway Size.",
			},
			"ha_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA Spoke Gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"transit_gw": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit Gateways this spoke has joined.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Instance tag of cloud provider.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Insane Mode for Spoke Gateway. Valid values: true, false. If insane mode is enabled, gateway size has to at least be c5 size.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable vpc_dns_server for Gateway. Valid values: true, false.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.",
			},
			"customized_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, " +
					"it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. " +
					"It applies to this spoke gateway only.",
			},
			"filtered_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, " +
					"filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s " +
					"routing table. It applies to this spoke gateway only.",
			},
			"included_advertised_spoke_routes": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. " +
					"When configured, it will replace all advertised routes from this VPC.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the spoke gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the spoke gateway created.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of HA spoke gateway.",
			},
			"ha_gw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix spoke gateway unique name of HA spoke gateway.",
			},
			"ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of HA spoke gateway.",
			},
			"enable_private_oob": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable private OOB.",
			},
			"oob_management_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OOB management subnet.",
			},
			"oob_availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OOB subnet availability zone.",
			},
			"ha_oob_management_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OOB HA management subnet.",
			},
			"ha_oob_availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OOB HA availability zone.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A map of tags assigned to the spoke gateway.",
			},
			"tunnel_detection_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IPSec tunnel down detection time for the spoke gateway.",
			},
			"availability_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability domain for OCI.",
			},
			"fault_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fault domain for OCI.",
			},
			"ha_availability_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA availability domain for OCI.",
			},
			"ha_fault_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA fault domain for OCI.",
			},
			"software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software version of the gateway.",
			},
			"ha_software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software version of the HA gateway.",
			},
			"image_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image version of the gateway.",
			},
			"ha_image_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image version of the HA gateway.",
			},
			"all_subnets_for_inspection": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subnets for inspection.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAviatrixSpokeGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
	}

	if d.Get("account_name").(string) != "" {
		gateway.AccountName = d.Get("account_name").(string)
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix spoke gateway: %s", err)
	}
	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)

			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0])
			d.Set("vpc_reg", gw.GatewayZone)

			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		} else if gw.CloudType == goaviatrix.AliCloud {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		}

		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("public_ip", gw.PublicIP)

		d.Set("subnet", gw.VpcNet)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)
		d.Set("image_version", gw.ImageVersion)
		d.Set("software_version", gw.SoftwareVersion)

		d.Set("enable_private_oob", gw.EnablePrivateOob)
		if gw.EnablePrivateOob {
			d.Set("oob_management_subnet", strings.Split(gw.OobManagementSubnet, "~~")[0])
			d.Set("oob_availability_zone", gw.GatewayZone)
		}

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				d.Set("insane_mode_az", gw.GatewayZone)
			} else {
				d.Set("insane_mode_az", "")
			}
		} else {
			d.Set("insane_mode", false)
			d.Set("insane_mode_az", "")
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled" {
			d.Set("enable_vpc_dns_server", true)
		} else {
			d.Set("enable_vpc_dns_server", false)
		}

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			d.Set("single_ip_snat", true)
		} else {
			d.Set("single_ip_snat", false)
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			d.Set("customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
		} else {
			d.Set("customized_spoke_vpc_routes", "")
		}

		if len(gw.FilteredSpokeVpcRoutes) != 0 {
			d.Set("filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
		} else {
			d.Set("filtered_spoke_vpc_routes", "")
		}

		if len(gw.AdvertisedSpokeRoutes) != 0 {
			d.Set("included_advertised_spoke_routes", strings.Join(gw.AdvertisedSpokeRoutes, ","))
		} else {
			d.Set("included_advertised_spoke_routes", "")
		}

		if gw.SpokeVpc == "yes" {
			var transitGws []string
			if gw.TransitGwName != "" {
				transitGws = append(transitGws, gw.TransitGwName)
			}
			if gw.EgressTransitGwName != "" {
				transitGws = append(transitGws, gw.EgressTransitGwName)
			}
			d.Set("transit_gw", strings.Join(transitGws, ","))
		} else {
			d.Set("transit_gw", "")
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			d.Set("availability_domain", gw.GatewayZone)
			d.Set("fault_domain", gw.FaultDomain)
		}

		haGateway := &goaviatrix.Gateway{
			AccountName: d.Get("account_name").(string),
			GwName:      d.Get("gw_name").(string) + "-hagw",
		}
		haGw, _ := client.GetGateway(haGateway)
		if haGw != nil {
			if goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
				d.Set("ha_subnet", haGw.VpcNet)
				d.Set("ha_zone", "")
			} else if goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
				d.Set("ha_zone", haGw.GatewayZone)
				d.Set("ha_subnet", "")
			}

			d.Set("ha_public_ip", haGw.PublicIP)
			d.Set("ha_gw_size", haGw.GwSize)
			d.Set("ha_cloud_instance_id", haGw.CloudnGatewayInstID)
			d.Set("ha_gw_name", haGw.GwName)
			d.Set("ha_private_ip", haGw.PrivateIP)
			d.Set("ha_image_version", haGw.ImageVersion)
			d.Set("ha_software_version", haGw.SoftwareVersion)
			if haGw.InsaneMode == "yes" && goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				d.Set("ha_insane_mode_az", haGw.GatewayZone)
			} else {
				d.Set("ha_insane_mode_az", "")
			}

			if goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
				d.Set("ha_availability_domain", haGw.GatewayZone)
				d.Set("ha_fault_domain", haGw.FaultDomain)
			}

			if haGw.EnablePrivateOob {
				d.Set("ha_oob_management_subnet", strings.Split(haGw.OobManagementSubnet, "~~")[0])
				d.Set("ha_oob_availability_zone", haGw.GatewayZone)
			}
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			tags := &goaviatrix.Tags{
				ResourceType: "gw",
				ResourceName: d.Get("gw_name").(string),
				CloudType:    gw.CloudType,
			}

			tagList, err := client.GetTags(tags)
			if err != nil {
				log.Printf("[WARN] Failed to get tags for spoke gateway %s: %v", tags.ResourceName, err)
			}
			if len(tags.Tags) > 0 {
				if err := d.Set("tags", tags.Tags); err != nil {
					log.Printf("[WARN] Error setting tags for spoke gateway %s: %v", tags.ResourceName, err)
				}
			}
			if len(tagList) > 0 {
				d.Set("tag_list", tagList)
			}
		}

		d.Set("tunnel_detection_time", gw.TunnelDetectionTime)
	}

	allSubnetsForInspection, err := client.GetAllSubnetsForInspection(gateway.GwName)
	if err != nil {
		return fmt.Errorf("couldn't get all subnets for inspection: %s", err)
	}
	d.Set("all_subnets_for_inspection", allSubnetsForInspection)

	d.SetId(gateway.GwName)
	return nil
}
