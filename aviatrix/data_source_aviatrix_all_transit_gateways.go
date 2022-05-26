package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixListAllTransitGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixAllTransitGatewayRead,

		Schema: map[string]*schema.Schema{
			"transit_gateway_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all Transit Gateways.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gw_name": {
							Type:        schema.TypeString,
							Computed:    true,
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
						"enable_vpc_dns_server": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable/Disable vpc_dns_server for Gateway. Valid values: true, false.",
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
						"software_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Software version of the gateway.",
						},
						"image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image version of the gateway.",
						},
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
						"bgp_polling_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP route polling time. Unit is in seconds.",
						},
						"prepend_as_path": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
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
										Type:        schema.TypeString,
										Required:    true,
										Description: "Subnet Info.",
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
										Type:        schema.TypeString,
										Required:    true,
										Description: "Subnet Info.",
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
						"local_as_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Changes the Aviatrix Transit Gateway ASN number before you setup Aviatrix Transit Gateway connection configurations.",
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
				},
			},
		},
	}
}

func dataSourceAviatrixAllTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	TransitGatewayList, err := client.GetTransitGatewayList()
	if err != nil {
		return err
	}
	var result []map[string]interface{}
	for i := range TransitGatewayList {
		gw := TransitGatewayList[i]
		transitGateway := make(map[string]interface{})
		transitGateway["cloud_type"] = gw.CloudType
		transitGateway["account_name"] = gw.AccountName
		transitGateway["gw_name"] = gw.GwName
		transitGateway["subnet"] = gw.VpcNet
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			transitGateway["vpc_id"] = strings.Split(gw.VpcID, "~~")[0]
			transitGateway["vpc_reg"] = gw.VpcRegion
			transitGateway["allocate_new_eip"] = gw.AllocateNewEipRead
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			transitGateway["vpc_id"] = gw.VpcID
			transitGateway["vpc_reg"] = gw.GatewayZone
			transitGateway["allocate_new_eip"] = gw.AllocateNewEipRead
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			transitGateway["vpc_id"] = gw.VpcID
			transitGateway["vpc_reg"] = gw.VpcRegion
			transitGateway["allocate_new_eip"] = true
		} else if gw.CloudType == goaviatrix.AliCloud {
			transitGateway["vpc_id"] = strings.Split(gw.VpcID, "~~")[0]
			transitGateway["vpc_reg"] = gw.VpcRegion
			transitGateway["allocate_new_eip"] = true
		}
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if gw.GatewayZone != "AvailabilitySet" {
				transitGateway["zone"] = "az-" + gw.GatewayZone
			}
		}
		transitGateway["enable_encrypt_volume"] = gw.EnableEncryptVolume
		transitGateway["public_ip"] = gw.PublicIP
		transitGateway["gw_size"] = gw.GwSize
		transitGateway["cloud_instance_id"] = gw.CloudnGatewayInstID
		transitGateway["security_group_id"] = gw.GwSecurityGroupID
		transitGateway["private_ip"] = gw.PrivateIP
		transitGateway["enable_multi_tier_transit"] = gw.EnableMultitierTransit
		transitGateway["image_version"] = gw.ImageVersion
		transitGateway["software_version"] = gw.SoftwareVersion
		transitGateway["enable_private_oob"] = gw.EnablePrivateOob
		if gw.EnablePrivateOob {
			transitGateway["oob_management_subnet"] = strings.Split(gw.OobManagementSubnet, "~~")[0]
			transitGateway["oob_availability_zone"] = gw.GatewayZone
		}
		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			transitGateway["single_ip_snat"] = true
		} else {
			transitGateway["single_ip_snat"] = false
		}
		if gw.SingleAZ == "yes" {
			transitGateway["single_az_ha"] = true
		} else {
			transitGateway["single_az_ha"] = false
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			transitGateway["enable_hybrid_connection"] = gw.EnableHybridConnection
		} else {
			transitGateway["enable_hybrid_connection"] = false
		}

		if gw.ConnectedTransit == "yes" {
			transitGateway["connected_transit"] = true
		} else {
			transitGateway["connected_transit"] = false
		}

		if gw.InsaneMode == "yes" {
			transitGateway["insane_mode"] = true
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				transitGateway["insane_mode_az"] = gw.GatewayZone
			} else {
				transitGateway["insane_mode_az"] = ""
			}
		} else {
			transitGateway["insane_mode"] = false
			transitGateway["insane_mode_az"] = ""
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			transitGateway["customized_spoke_vpc_routes"] = strings.Join(gw.CustomizedSpokeVpcRoutes, ",")
		} else {
			transitGateway["customized_spoke_vpc_routes"] = ""
		}

		if len(gw.FilteredSpokeVpcRoutes) != 0 {
			transitGateway["filtered_spoke_vpc_routes"] = strings.Join(gw.FilteredSpokeVpcRoutes, ",")
		} else {
			transitGateway["filtered_spoke_vpc_routes"] = ""
		}

		if len(gw.ExcludeCidrList) != 0 {
			transitGateway["excluded_advertised_spoke_routes"] = strings.Join(gw.ExcludeCidrList, ",")
		} else {
			transitGateway["excluded_advertised_spoke_routes"] = ""
		}
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled" {
			transitGateway["enable_vpc_dns_server"] = true
		} else {
			transitGateway["enable_vpc_dns_server"] = false
		}
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			transitGateway["availability_domain"] = gw.GatewayZone
			transitGateway["fault_domain"] = gw.FaultDomain
		}
		transitGateway["tunnel_detection_time"] = gw.TunnelDetectionTime
		transitGateway["enable_gateway_load_balancer"] = gw.EnableGatewayLoadBalancer
		transitGateway["learned_cidrs_approval_mode"] = gw.LearnedCidrsApprovalMode
		transitGateway["enable_jumbo_frame"] = gw.JumboFrame
		transitGateway["enable_private_oob"] = gw.EnablePrivateOob
		transitGateway["bgp_polling_time"] = strconv.Itoa(gw.BgpPollingTime)
		transitGateway["bgp_hold_time"] = gw.BgpHoldTime
		transitGateway["local_as_number"] = gw.LocalASNumber
		transitGateway["bgp_ecmp"] = gw.BgpEcmp
		transitGateway["enable_segmentation"] = gw.EnableSegmentation
		transitGateway["enable_active_standby"] = gw.EnableActiveStandby
		transitGateway["enable_active_standby_preemptive"] = gw.EnableActiveStandbyPreemptive
		transitGateway["enable_transit_summarize_cidr_to_tgw"] = gw.EnableTransitSummarizeCidrToTgw

		if gw.EnableTransitFirenet && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			transitGateway["lan_vpc_id"] = gw.BundleVpcInfo.LAN.VpcID
			transitGateway["lan_private_subnet"] = strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0]
		}

		var prependAsPath []string
		for _, p := range strings.Split(gw.PrependASPath, " ") {
			if p != "" {
				prependAsPath = append(prependAsPath, p)
			}
		}

		transitGateway["prepend_as_path"] = prependAsPath
		transitGateway["enable_monitor_gateway_subnets"] = gw.MonitorSubnetsAction == "enable"
		transitGateway["monitor_exclude_list"] = gw.MonitorExcludeGWList

		transitGateway["enable_bgp_over_lan"] = goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan {
			if len(gw.BgpLanInterfaces) != 0 {
				var interfaces []map[string]interface{}
				for _, bgpLanInterface := range gw.BgpLanInterfaces {
					interfaceDict := make(map[string]interface{})
					interfaceDict["vpc_id"] = bgpLanInterface.VpcID
					interfaceDict["subnet"] = bgpLanInterface.Subnet
					interfaces = append(interfaces, interfaceDict)
				}
				transitGateway["bgp_lan_interfaces"] = interfaces
			}

			if len(gw.HaGw.HaBgpLanInterfaces) != 0 {
				var haInterfaces []map[string]interface{}
				for _, haBgpLanInterface := range gw.HaGw.HaBgpLanInterfaces {
					interfaceDict := make(map[string]interface{})
					interfaceDict["vpc_id"] = haBgpLanInterface.VpcID
					interfaceDict["subnet"] = haBgpLanInterface.Subnet
					haInterfaces = append(haInterfaces, interfaceDict)
				}
				transitGateway["ha_bgp_lan_interfaces"] = haInterfaces
			}
		} else {
			transitGateway["bgp_lan_ip_list"] = nil
			transitGateway["ha_bgp_lan_ip_list"] = nil
		}

		if gw.EnableSpotInstance {
			transitGateway["enable_spot_instance"] = true
			transitGateway["spot_price"] = gw.SpotPrice
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.ReuseEip, ":")
			if len(azureEip) == 3 {
				transitGateway["azure_eip_name_resource_group"] = fmt.Sprintf("%s:%s", azureEip[0], azureEip[1])
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the Transit Gateway %s", gw.GwName)
			}
		}

		result = append(result, transitGateway)
	}

	if err = d.Set("transit_gateway_list", result); err != nil {
		return fmt.Errorf("couldn't set transit_gateway_list: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil

}
