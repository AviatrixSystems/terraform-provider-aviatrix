package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixTransitGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixTransitGatewayRead,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Transit Gateway name. This can be used for getting gateway.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of a Cloud-Account in Aviatrix controller.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC ID.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of cloud provider.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance type.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Range of the subnet where the transit gateway is launched.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.",
			},
			"allocate_new_eip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the eip is newly allocated or not.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Transit Gateway created.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Subnet. Required for enabling HA for AWS/Azure transit gateway.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Zone. Required if enabling HA for GCP.",
			},
			"ha_insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Transit HA Gateway. Required if insane_mode is enabled and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set).",
			},
			"ha_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA Transit Gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable this feature.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable Source NAT feature in 'single_ip' mode for this container.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Instance tag of cloud provider. Only supported for AWS provider.",
			},
			"enable_hybrid_connection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Sign of readiness for TGW connection.",
			},
			"connected_transit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Connected Transit status.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable Insane Mode for Spoke Gateway.",
			},
			"enable_firenet": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether firenet interfaces is enabled.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable vpc_dns_server for Gateway. Valid values: true, false.",
			},
			"enable_advertise_transit_cidr": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable advertise transit VPC network CIDR.",
			},
			"bgp_manual_spoke_advertise_cidrs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intended CIDR list to advertise to VGW.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable encrypt gateway EBS volume. Only supported for AWS provider.",
			},
			"customized_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, " +
					"it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. " +
					"It applies to this spoke gateway only."},
			"filtered_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, " +
					"filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s " +
					"routing table. It applies to this spoke gateway only.",
			},
			"excluded_advertised_spoke_routes": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of comma separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'. " +
					"When configured, it inspects all the advertised CIDRs from its spoke gateways and " +
					"remove those included in the 'Excluded CIDR List'.",
			},
			"customized_transit_vpc_routes": {
				Type:     schema.TypeSet,
				Computed: true,
				Description: "A list of CIDRs to be customized for the transit VPC routes. " +
					"When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs." +
					"To be effective, `enable_advertise_transit_cidr` or firewall management access for a transit firenet gateway must be enabled.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_transit_firenet": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Switch to enable/disable transit firenet interfaces for transit gateway.",
			},
			"enable_egress_transit_firenet": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specify whether to enable egress transit firenet interfaces or not.",
			},
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Switch to enable/disable encrypted transit approval for transit gateway.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the transit gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the transit gateway.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the transit gateway created.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of HA transit gateway.",
			},
			"ha_gw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix transit gateway unique name of HA transit gateway.",
			},
			"ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of HA transit gateway.",
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
				Description: "A map of tags assigned to the transit gateway.",
			},
			"enable_multi_tier_transit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Multi-tier Transit mode on transit gateway.",
			},
			"tunnel_detection_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IPSec tunnel down detection time for the transit gateway.",
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

			// NEW ATTRIBUTES
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'.",
			},
			"enable_gateway_load_balancer": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Enable firenet interfaces with AWS Gateway Load Balancer. Only set when `enable_firenet` or `enable_transit_firenet`" +
					" are set to true and `cloud_type` = 1 (AWS). Currently AWS Gateway Load Balancer is only supported " +
					"in AWS regions us-west-2 and us-east-1.",
			},
			"lan_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LAN VPC ID. Only used for GCP Transit FireNet.",
			},
			"lan_private_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LAN Private Subnet. Only used for GCP Transit FireNet.",
			},
			"learned_cidrs_approval_mode": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Set the learned CIDRs approval mode. Only set when 'enable_learned_cidrs_approval' is " +
					"set to true. If set to 'gateway', learned CIDR approval applies to ALL connections. If set to " +
					"'connection', learned CIDR approval is configured on a per connection basis. When configuring per " +
					"connection, use the enable_learned_cidrs_approval attribute within the connection resource to " +
					"toggle learned CIDR approval.",
			},
			"approved_learned_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Computed:    true,
				Description: "Approved learned CIDRs. Available as of provider version R2.21+.",
			},
			"bgp_polling_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BGP route polling time. Unit is in seconds.",
			},
			"prepend_as_path": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				Computed:    true,
				Description: "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.",
			},
			"bgp_ecmp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Equal Cost Multi Path (ECMP) routing for the next hop.",
			},
			"enable_segmentation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable segmentation to allow association of transit gateway to security domains.",
			},
			"enable_active_standby": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enables Active-Standby Mode, available only with HA enabled.",
			},
			"enable_active_standby_preemptive": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.",
			},
			"enable_monitor_gateway_subnets": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Enable [monitor gateway subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet). " +
					"Only valid for cloud_type = 1 (AWS) or 256 (AWSGov).",
			},
			"monitor_exclude_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A set of monitored instance ids. Only valid when 'enable_monitor_gateway_subnets' = true.",
			},
			"enable_bgp_over_lan": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Pre-allocate a network interface(eth4) for \"BGP over LAN\" functionality. Only valid for cloud_type = 4 (GCP) and 8 (Azure). Available as of provider version R2.18+",
			},
			"bgp_lan_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP Transit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC-ID of GCP cloud provider.",
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
							Description:  "Subnet Info.",
						},
					},
				},
			},
			"ha_bgp_lan_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP HA Transit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC-ID of GCP cloud provider.",
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
							Description:  "Subnet Info.",
						},
					},
				},
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable jumbo frame support for transit gateway.",
			},
			"bgp_hold_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "BGP Hold Time.",
			},
			"enable_transit_summarize_cidr_to_tgw": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable summarize CIDR to TGW.",
			},
			"enable_spot_instance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable spot instance. NOT supported for production deployment.",
			},
			"spot_price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Price for spot instance. NOT supported for production deployment.",
			},
			"azure_eip_name_resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the public IP address and its resource group in Azure to assign to this Transit Gateway.",
			},
			"ha_azure_eip_name_resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the public IP address and its resource group in Azure to assign to the HA Transit Gateway.",
			},
			"local_as_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Changes the Aviatrix Transit Gateway ASN number before you setup Aviatrix Transit Gateway connection configurations.",
			},
			"ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA security group used for the transit gateway.",
			},
			"lan_interface_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit gateway lan interface cidr.",
			},
			"ha_lan_interface_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit gateway lan interface cidr for the HA gateway.",
			},
			"bgp_lan_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "List of available BGP LAN interface IPs for transit external device connection creation. " +
					"Only supports GCP. Available as of provider version R2.21.0+.",
			},
			"ha_bgp_lan_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "List of available BGP LAN interface IPs for transit external device HA connection creation. " +
					"Only supports GCP. Available as of provider version R2.21.0+.",
			},
		},
	}
}

func dataSourceAviatrixTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit Gateway: %s", err)
	}
	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("subnet", gw.VpcNet)

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

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			_, zoneIsSet := d.GetOk("zone")
			if zoneIsSet && gw.GatewayZone != "AvailabilitySet" {
				d.Set("zone", "az-"+gw.GatewayZone)
			}
		}

		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("public_ip", gw.PublicIP)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)
		d.Set("enable_multi_tier_transit", gw.EnableMultitierTransit)
		d.Set("image_version", gw.ImageVersion)
		d.Set("software_version", gw.SoftwareVersion)

		d.Set("enable_private_oob", gw.EnablePrivateOob)
		if gw.EnablePrivateOob {
			d.Set("oob_management_subnet", strings.Split(gw.OobManagementSubnet, "~~")[0])
			d.Set("oob_availability_zone", gw.GatewayZone)
		}

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			d.Set("single_ip_snat", true)
		} else {
			d.Set("single_ip_snat", false)
		}

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			d.Set("enable_hybrid_connection", gw.EnableHybridConnection)
		} else {
			d.Set("enable_hybrid_connection", false)
		}

		if gw.ConnectedTransit == "yes" {
			d.Set("connected_transit", true)
		} else {
			d.Set("connected_transit", false)
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
		if len(gw.ExcludeCidrList) != 0 {
			d.Set("excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
		} else {
			d.Set("excluded_advertised_spoke_routes", "")
		}

		gwDetail, err := client.GetGatewayDetail(gw)
		if err != nil {
			return fmt.Errorf("couldn't get Aviatrix Transit Gateway: %s", err)
		}

		d.Set("enable_firenet", gwDetail.EnableFireNet)
		d.Set("enable_transit_firenet", gwDetail.EnableTransitFireNet)
		d.Set("enable_egress_transit_firenet", gwDetail.EnableEgressTransitFireNet)
		d.Set("customized_transit_vpc_routes", gwDetail.CustomizedTransitVpcRoutes)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled" {
			d.Set("enable_vpc_dns_server", true)
		} else {
			d.Set("enable_vpc_dns_server", false)
		}

		if gwDetail.EnableAdvertiseTransitCidr == "yes" {
			d.Set("enable_advertise_transit_cidr", true)
		} else {
			d.Set("enable_advertise_transit_cidr", false)
		}

		if gwDetail.LearnedCidrsApproval == "yes" {
			d.Set("enable_learned_cidrs_approval", true)
		} else {
			d.Set("enable_learned_cidrs_approval", false)
		}

		bgpMSAN := ""
		for i := range gwDetail.BgpManualSpokeAdvertiseCidrs {
			if i == 0 {
				bgpMSAN = bgpMSAN + gwDetail.BgpManualSpokeAdvertiseCidrs[i]
			} else {
				bgpMSAN = bgpMSAN + "," + gwDetail.BgpManualSpokeAdvertiseCidrs[i]
			}
		}
		d.Set("bgp_manual_spoke_advertise_cidrs", bgpMSAN)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			tags := &goaviatrix.Tags{
				ResourceType: "gw",
				ResourceName: d.Get("gw_name").(string),
				CloudType:    gw.CloudType,
			}

			tagList, err := client.GetTags(tags)
			if err != nil {
				log.Printf("[WARN] Failed to get tags for transit gateway %s: %v", tags.ResourceName, err)
			}
			if len(tags.Tags) > 0 {
				if err := d.Set("tags", tags.Tags); err != nil {
					log.Printf("[WARN] Error setting tags for transit gateway %s: %v", tags.ResourceName, err)
				}
			}
			if len(tagList) > 0 {
				d.Set("tag_list", tagList)
			}
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
			d.Set("ha_security_group_id", gw.HaGw.GwSecurityGroupID)

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

			if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
				if len(azureEip) == 3 {
					d.Set("ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
				} else {
					log.Printf("[WARN] could not get Azure EIP name and resource group for the HA Gateway %s", gw.GwName)
				}
			}

			lanCidr, err := client.GetTransitGatewayLanCidr(gw.HaGw.GwName)
			if err != nil && err != goaviatrix.ErrNotFound {
				log.Printf("[WARN] Error getting lan cidr for HA transit gateway %s due to %s", gw.HaGw.GwName, err)
			}
			d.Set("ha_lan_interface_cidr", lanCidr)
		}

		d.Set("tunnel_detection_time", gw.TunnelDetectionTime)
		d.Set("enable_gateway_load_balancer", gw.EnableGatewayLoadBalancer)
		d.Set("learned_cidrs_approval_mode", gw.LearnedCidrsApprovalMode)
		d.Set("enable_jumbo_frame", gw.JumboFrame)
		d.Set("enable_private_oob", gw.EnablePrivateOob)
		d.Set("bgp_polling_time", strconv.Itoa(gw.BgpPollingTime))
		d.Set("bgp_hold_time", gw.BgpHoldTime)
		d.Set("local_as_number", gw.LocalASNumber)
		d.Set("bgp_ecmp", gw.BgpEcmp)
		d.Set("enable_segmentation", gw.EnableSegmentation)
		d.Set("enable_active_standby", gw.EnableActiveStandby)
		d.Set("enable_active_standby_preemptive", gw.EnableActiveStandbyPreemptive)
		d.Set("enable_transit_summarize_cidr_to_tgw", gw.EnableTransitSummarizeCidrToTgw)

		if gw.EnableTransitFirenet && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("lan_vpc_id", gw.BundleVpcInfo.LAN.VpcID)
			d.Set("lan_private_subnet", strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0])
		}

		if gw.EnableLearnedCidrsApproval {
			transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: gw.GwName})
			if err != nil {
				return fmt.Errorf("could not get advanced config for transit gateway: %v", err)
			}

			if err = d.Set("approved_learned_cidrs", transitAdvancedConfig.ApprovedLearnedCidrs); err != nil {
				return fmt.Errorf("could not set approved_learned_cidrs into state: %v", err)
			}
		} else {
			d.Set("approved_learned_cidrs", nil)
		}

		var prependAsPath []string
		for _, p := range strings.Split(gw.PrependASPath, " ") {
			if p != "" {
				prependAsPath = append(prependAsPath, p)
			}
		}
		err = d.Set("prepend_as_path", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %v", err)
		}

		d.Set("enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
		if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
			return fmt.Errorf("setting 'monitor_exclude_list' to state: %v", err)
		}

		d.Set("enable_bgp_over_lan", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan {
			if len(gw.BgpLanInterfaces) != 0 {
				var interfaces []map[string]interface{}
				for _, bgpLanInterface := range gw.BgpLanInterfaces {
					interfaceDict := make(map[string]interface{})
					interfaceDict["vpc_id"] = bgpLanInterface.VpcID
					interfaceDict["subnet"] = bgpLanInterface.Subnet
					interfaces = append(interfaces, interfaceDict)
				}
				if err = d.Set("bgp_lan_interfaces", interfaces); err != nil {
					return fmt.Errorf("could not set bgp_lan_interfaces into state: %v", err)
				}
			}

			if len(gw.HaGw.HaBgpLanInterfaces) != 0 {
				var haInterfaces []map[string]interface{}
				for _, haBgpLanInterface := range gw.HaGw.HaBgpLanInterfaces {
					interfaceDict := make(map[string]interface{})
					interfaceDict["vpc_id"] = haBgpLanInterface.VpcID
					interfaceDict["subnet"] = haBgpLanInterface.Subnet
					haInterfaces = append(haInterfaces, interfaceDict)
				}
				if err = d.Set("ha_bgp_lan_interfaces", haInterfaces); err != nil {
					return fmt.Errorf("could not set ha_bgp_lan_interfaces into state: %v", err)
				}
			}

			bgpLanIpInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not get BGP LAN IP info for GCP transit gateway %s: %v", gateway.GwName, err)
			}
			if err = d.Set("bgp_lan_ip_list", bgpLanIpInfo.BgpLanIpList); err != nil {
				return fmt.Errorf("could not set bgp_lan_ip_list into state: %v", err)
			}
			if len(bgpLanIpInfo.HaBgpLanIpList) != 0 {
				if err = d.Set("ha_bgp_lan_ip_list", bgpLanIpInfo.HaBgpLanIpList); err != nil {
					return fmt.Errorf("could not set ha_bgp_lan_ip_list into tate: %v", err)
				}
			} else {
				d.Set("ha_bgp_lan_ip_list", nil)
			}
		} else {
			d.Set("bgp_lan_ip_list", nil)
			d.Set("ha_bgp_lan_ip_list", nil)
		}

		if gw.EnableSpotInstance {
			d.Set("enable_spot_instance", true)
			d.Set("spot_price", gw.SpotPrice)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.ReuseEip, ":")
			if len(azureEip) == 3 {
				d.Set("azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the Transit Gateway %s", gw.GwName)
			}
		}

		lanCidr, err := client.GetTransitGatewayLanCidr(gw.GwName)
		if err != nil && err != goaviatrix.ErrNotFound {
			log.Printf("[WARN] Error getting lan cidr for transit gateway %s due to %s", gw.GwName, err)
		}
		d.Set("lan_interface_cidr", lanCidr)
	}

	d.SetId(gateway.GwName)
	return nil
}
