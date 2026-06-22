package aviatrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixSpokeGateways() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixSpokeGatewaysRead,

		Schema: map[string]*schema.Schema{
			"gateway_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all Spoke Gateways.",
				Elem: &schema.Resource{
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
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'.",
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
						"enable_monitor_gateway_subnets": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: "Enable [monitor gateway subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet). " +
								"Only valid for cloud_type = 1 (AWS) or 256 (AWSGov). Valid values: true, false. Default value: false.",
						},
						"monitor_exclude_list": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A set of monitored instance ids. Only set when 'enable_monitor_gateway_subnets' = true.",
						},
						"enable_jumbo_frame": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable jumbo frame support for spoke gateway.",
						},
						"enable_private_vpc_default_route": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Config Private VPC Default Route.",
						},
						"enable_skip_public_route_table_update": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Skip Public Route Table Update.",
						},
						"enable_auto_advertise_s2c_cidrs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Automatically advertise remote CIDR to Aviatrix Transit Gateway when route based Site2Cloud Tunnel is created.",
						},
						"spoke_bgp_manual_advertise_cidrs": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Intended CIDR list to be advertised to external BGP router.",
						},
						"enable_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable BGP.",
						},
						"enable_learned_cidrs_approval": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Switch to enable/disable encrypted transit approval for BGP Spoke Gateway.",
						},
						"learned_cidrs_approval_mode": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "Set the learned CIDRs approval mode for BGP Spoke Gateway. Only valid when 'enable_learned_cidrs_approval' is " +
								"set to true. Currently, only 'gateway' is supported: learned CIDR approval applies to " +
								"ALL connections. Default value: 'gateway'.",
						},
						"bgp_ecmp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable Equal Cost Multi Path (ECMP) routing for the next hop for BGP Spoke Gateway.",
						},
						"enable_active_standby": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enables Active-Standby Mode, available only with HA enabled for BGP Spoke Gateway.",
						},
						"enable_active_standby_preemptive": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.",
						},
						"disable_route_propagation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Disables route propagation on BGP Spoke to attached Transit Gateway.",
						},
						"local_as_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Changes the Aviatrix BGP Spoke Gateway ASN number before you setup Aviatrix BGP Spoke Gateway connection configurations.",
						},
						"prepend_as_path": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices. Only valid for BGP Spoke Gateway",
						},
						"bgp_polling_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP route polling time for BGP Spoke Gateway. Unit is in seconds.",
						},
						"bgp_hold_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP Hold Time for BGP Spoke Gateway. Unit is in seconds.",
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
							Description: "The name of the public IP address and its resource group in Azure to assign to this Spoke Gateway.",
						},
						"eip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP address of the Spoke Gateway.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixSpokeGatewaysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	SpokeGatewayList, err := client.GetSpokeGatewayList(ctx)
	if err != nil {
		return diag.Errorf("could not get Aviatrix Spoke Gateway List: %s", err)
	}
	var result []map[string]interface{}
	for i := range SpokeGatewayList {
		gw := SpokeGatewayList[i]
		spokeGateway := make(map[string]interface{})

		spokeGateway["gw_name"] = gw.GwName
		spokeGateway["cloud_type"] = gw.CloudType
		spokeGateway["account_name"] = gw.AccountName

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			spokeGateway["vpc_id"] = strings.Split(gw.VpcID, "~~")[0]
			spokeGateway["vpc_reg"] = gw.VpcRegion

			if gw.AllocateNewEipRead {
				spokeGateway["allocate_new_eip"] = true
			} else {
				spokeGateway["allocate_new_eip"] = false
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			spokeGateway["vpc_id"] = gw.VpcID
			spokeGateway["vpc_reg"] = gw.GatewayZone

			if gw.AllocateNewEipRead {
				spokeGateway["allocate_new_eip"] = true
			} else {
				spokeGateway["allocate_new_eip"] = false
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			spokeGateway["vpc_id"] = gw.VpcID
			spokeGateway["vpc_reg"] = gw.VpcRegion
			spokeGateway["allocate_new_eip"] = true
		} else if gw.CloudType == goaviatrix.AliCloud {
			spokeGateway["vpc_id"] = strings.Split(gw.VpcID, "~~")[0]
			spokeGateway["vpc_reg"] = gw.VpcRegion
			spokeGateway["allocate_new_eip"] = true
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			spokeGateway["zone"] = "az-" + gw.GatewayZone
		}

		spokeGateway["enable_encrypt_volume"] = gw.EnableEncryptVolume
		spokeGateway["public_ip"] = gw.PublicIP
		spokeGateway["subnet"] = gw.VpcNet
		spokeGateway["gw_size"] = gw.GwSize
		spokeGateway["cloud_instance_id"] = gw.CloudnGatewayInstID
		spokeGateway["security_group_id"] = gw.GwSecurityGroupID
		spokeGateway["private_ip"] = gw.PrivateIP
		spokeGateway["image_version"] = gw.ImageVersion
		spokeGateway["software_version"] = gw.SoftwareVersion
		spokeGateway["eip"] = gw.PublicIP
		spokeGateway["enable_private_oob"] = gw.EnablePrivateOob
		if gw.EnablePrivateOob {
			spokeGateway["oob_management_subnet"] = strings.Split(gw.OobManagementSubnet, "~~")[0]
			spokeGateway["oob_availability_zone"] = gw.GatewayZone
		}

		if gw.SingleAZ == "yes" {
			spokeGateway["single_az_ha"] = true
		} else {
			spokeGateway["single_az_ha"] = false
		}

		if gw.InsaneMode == "yes" {
			spokeGateway["insane_mode"] = true
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				spokeGateway["insane_mode_az"] = gw.GatewayZone
			} else {
				spokeGateway["insane_mode_az"] = ""
			}
		} else {
			spokeGateway["insane_mode"] = false
			spokeGateway["insane_mode_az"] = ""
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled" {
			spokeGateway["enable_vpc_dns_server"] = true
		} else {
			spokeGateway["enable_vpc_dns_server"] = false
		}

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			spokeGateway["single_ip_snat"] = true
		} else {
			spokeGateway["single_ip_snat"] = false
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			spokeGateway["customized_spoke_vpc_routes"] = strings.Join(gw.CustomizedSpokeVpcRoutes, ",")
		} else {
			spokeGateway["customized_spoke_vpc_routes"] = ""
		}

		if len(gw.FilteredSpokeVpcRoutes) != 0 {
			spokeGateway["filtered_spoke_vpc_routes"] = strings.Join(gw.FilteredSpokeVpcRoutes, ",")
		} else {
			spokeGateway["filtered_spoke_vpc_routes"] = ""
		}

		if len(gw.AdvertisedSpokeRoutes) != 0 {
			spokeGateway["included_advertised_spoke_routes"] = strings.Join(gw.AdvertisedSpokeRoutes, ",")
		} else {
			spokeGateway["included_advertised_spoke_routes"] = ""
		}

		if gw.SpokeVpc == "yes" {
			var transitGws []string
			if gw.TransitGwName != "" {
				transitGws = append(transitGws, gw.TransitGwName)
			}
			if gw.EgressTransitGwName != "" {
				transitGws = append(transitGws, gw.EgressTransitGwName)
			}
			spokeGateway["transit_gw"] = strings.Join(transitGws, ",")
		} else {
			spokeGateway["transit_gw"] = ""
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			spokeGateway["availability_domain"] = gw.GatewayZone
			spokeGateway["fault_domain"] = gw.FaultDomain
		}

		spokeGateway["tunnel_detection_time"] = gw.TunnelDetectionTime
		spokeGateway["enable_jumbo_frame"] = gw.JumboFrame
		spokeGateway["enable_private_vpc_default_route"] = gw.PrivateVpcDefaultEnabled
		spokeGateway["enable_skip_public_route_table_update"] = gw.SkipPublicVpcUpdateEnabled
		spokeGateway["enable_auto_advertise_s2c_cidrs"] = gw.AutoAdvertiseCidrsEnabled
		spokeGateway["spoke_bgp_manual_advertise_cidrs"] = gw.BgpManualSpokeAdvertiseCidrs
		spokeGateway["enable_bgp"] = gw.EnableBgp
		spokeGateway["enable_learned_cidrs_approval"] = gw.EnableLearnedCidrsApproval
		spokeGateway["bgp_ecmp"] = gw.BgpEcmp
		spokeGateway["enable_active_standby"] = gw.EnableActiveStandby
		spokeGateway["enable_active_standby_preemptive"] = gw.EnableActiveStandbyPreemptive
		spokeGateway["disable_route_propagation"] = gw.DisableRoutePropagation
		spokeGateway["local_as_number"] = gw.LocalASNumber
		var prependAsPath []string
		for _, p := range strings.Split(gw.PrependASPath, " ") {
			if p != "" {
				prependAsPath = append(prependAsPath, p)
			}
		}
		spokeGateway["prepend_as_path"] = prependAsPath
		spokeGateway["enable_monitor_gateway_subnets"] = gw.MonitorSubnetsAction == "enable"
		spokeGateway["monitor_exclude_list"] = gw.MonitorExcludeGWList

		if gw.EnableBgp {
			spokeGateway["learned_cidrs_approval_mode"] = gw.LearnedCidrsApprovalMode
			spokeGateway["bgp_polling_time"] = gw.BgpPollingTime
			spokeGateway["bgp_hold_time"] = gw.BgpHoldTime
		} else {
			spokeGateway["learned_cidrs_approval_mode"] = "gateway"
			spokeGateway["bgp_polling_time"] = 50
			spokeGateway["bgp_hold_time"] = 180
		}

		if gw.EnableSpotInstance {
			spokeGateway["enable_spot_instance"] = true
			spokeGateway["spot_price"] = gw.SpotPrice
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.ReuseEip, ":")
			if len(azureEip) == 3 {
				spokeGateway["azure_eip_name_resource_group"] = fmt.Sprintf("%s:%s", azureEip[0], azureEip[1])
			}
		}

		result = append(result, spokeGateway)
	}

	if err = d.Set("gateway_list", result); err != nil {
		return diag.Errorf("couldn't set gateway_list: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}
