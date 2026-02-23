package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
			"approved_learned_cidrs": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Approved learned CIDRs for BGP Spoke Gateway. Available as of provider version R2.21+.",
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
			"ha_azure_eip_name_resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the public IP address and its resource group in Azure to assign to the HA Spoke Gateway.",
			},
			"ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA security group used for the spoke gateway.",
			},
			"eip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The EIP address of the Spoke Gateway.",
			},
			"ha_eip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The EIP address of the HA Spoke Gateway.",
			},
		},
	}
}

func dataSourceAviatrixSpokeGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		GwName: getString(d, "gw_name"),
	}

	if getString(d, "account_name") != "" {
		gateway.AccountName = getString(d, "account_name")
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix spoke gateway: %w", err)
	}
	if gw != nil {
		mustSet(d, "cloud_type", gw.CloudType)
		mustSet(d, "account_name", gw.AccountName)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
			mustSet(d, "vpc_reg", gw.VpcRegion)

			if gw.AllocateNewEipRead {
				mustSet(d, "allocate_new_eip", true)
			} else {
				mustSet(d, "allocate_new_eip", false)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			mustSet(d, "vpc_id", gw.VpcID)
			mustSet(d, "vpc_reg", gw.GatewayZone)

			if gw.AllocateNewEipRead {
				mustSet(d, "allocate_new_eip", true)
			} else {
				mustSet(d, "allocate_new_eip", false)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			mustSet(d, "vpc_id", gw.VpcID)
			mustSet(d, "vpc_reg", gw.VpcRegion)
			mustSet(d, "allocate_new_eip", true)
		} else if gw.CloudType == goaviatrix.AliCloud {
			mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
			mustSet(d, "vpc_reg", gw.VpcRegion)
			mustSet(d, "allocate_new_eip", true)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			_, zoneIsSet := d.GetOk("zone")
			if zoneIsSet && gw.GatewayZone != "AvailabilitySet" {
				mustSet(d, "zone", "az-"+gw.GatewayZone)
			}
		}
		mustSet(d, "enable_encrypt_volume", gw.EnableEncryptVolume)
		mustSet(d, "public_ip", gw.PublicIP)
		mustSet(d, "subnet", gw.VpcNet)
		mustSet(d, "gw_size", gw.GwSize)
		mustSet(d, "cloud_instance_id", gw.CloudnGatewayInstID)
		mustSet(d, "security_group_id", gw.GwSecurityGroupID)
		mustSet(d, "private_ip", gw.PrivateIP)
		mustSet(d, "image_version", gw.ImageVersion)
		mustSet(d, "software_version", gw.SoftwareVersion)
		mustSet(d, "eip", gw.PublicIP)
		mustSet(d, "ha_eip", gw.HaGw.PublicIP)
		mustSet(d, "enable_private_oob", gw.EnablePrivateOob)
		if gw.EnablePrivateOob {
			mustSet(d, "oob_management_subnet", strings.Split(gw.OobManagementSubnet, "~~")[0])
			mustSet(d, "oob_availability_zone", gw.GatewayZone)
		}

		if gw.SingleAZ == "yes" {
			mustSet(d, "single_az_ha", true)
		} else {
			mustSet(d, "single_az_ha", false)
		}

		if gw.InsaneMode == "yes" {
			mustSet(d, "insane_mode", true)
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				mustSet(d, "insane_mode_az", gw.GatewayZone)
			} else {
				mustSet(d, "insane_mode_az", "")
			}
		} else {
			mustSet(d, "insane_mode", false)
			mustSet(d, "insane_mode_az", "")
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled" {
			mustSet(d, "enable_vpc_dns_server", true)
		} else {
			mustSet(d, "enable_vpc_dns_server", false)
		}

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			mustSet(d, "single_ip_snat", true)
		} else {
			mustSet(d, "single_ip_snat", false)
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			mustSet(d, "customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
		} else {
			mustSet(d, "customized_spoke_vpc_routes", "")
		}

		if len(gw.FilteredSpokeVpcRoutes) != 0 {
			mustSet(d, "filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
		} else {
			mustSet(d, "filtered_spoke_vpc_routes", "")
		}

		if len(gw.AdvertisedSpokeRoutes) != 0 {
			mustSet(d, "included_advertised_spoke_routes", strings.Join(gw.AdvertisedSpokeRoutes, ","))
		} else {
			mustSet(d, "included_advertised_spoke_routes", "")
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			mustSet(d, "availability_domain", gw.GatewayZone)
			mustSet(d, "fault_domain", gw.FaultDomain)
		}

		haGateway := &goaviatrix.Gateway{
			AccountName: getString(d, "account_name"),
			GwName:      getString(d, "gw_name") + "-hagw",
		}
		haGw, _ := client.GetGateway(haGateway)
		if haGw != nil {
			if goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
				mustSet(d, "ha_subnet", haGw.VpcNet)
				mustSet(d, "ha_zone", "")
			} else if goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
				mustSet(d, "ha_zone", haGw.GatewayZone)
				mustSet(d, "ha_subnet", "")
			}
			mustSet(d, "ha_public_ip", haGw.PublicIP)
			mustSet(d, "ha_gw_size", haGw.GwSize)
			mustSet(d, "ha_cloud_instance_id", haGw.CloudnGatewayInstID)
			mustSet(d, "ha_gw_name", haGw.GwName)
			mustSet(d, "ha_private_ip", haGw.PrivateIP)
			mustSet(d, "ha_image_version", haGw.ImageVersion)
			mustSet(d, "ha_software_version", haGw.SoftwareVersion)
			mustSet(d, "ha_security_group_id", gw.HaGw.GwSecurityGroupID)
			if haGw.InsaneMode == "yes" && goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				mustSet(d, "ha_insane_mode_az", haGw.GatewayZone)
			} else {
				mustSet(d, "ha_insane_mode_az", "")
			}

			if goaviatrix.IsCloudType(haGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
				mustSet(d, "ha_availability_domain", haGw.GatewayZone)
				mustSet(d, "ha_fault_domain", haGw.FaultDomain)
			}

			if haGw.EnablePrivateOob {
				mustSet(d, "ha_oob_management_subnet", strings.Split(haGw.OobManagementSubnet, "~~")[0])
				mustSet(d, "ha_oob_availability_zone", haGw.GatewayZone)
			}

			if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
				if len(azureEip) == 3 {
					mustSet(d, "ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
				} else {
					log.Printf("[WARN] could not get Azure EIP name and resource group for the HA Gateway %s", gw.GwName)
				}
			}
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			tags := &goaviatrix.Tags{
				ResourceType: "gw",
				ResourceName: getString(d, "gw_name"),
				CloudType:    gw.CloudType,
			}

			_, err := client.GetTags(tags)
			if err != nil {
				log.Printf("[WARN] Failed to get tags for spoke gateway %s: %v", tags.ResourceName, err)
			}
			if len(tags.Tags) > 0 {
				if err := d.Set("tags", tags.Tags); err != nil {
					log.Printf("[WARN] Error setting tags for spoke gateway %s: %v", tags.ResourceName, err)
				}
			}
		}
		mustSet(d, "tunnel_detection_time", gw.TunnelDetectionTime)
		mustSet(d, "enable_jumbo_frame", gw.JumboFrame)
		mustSet(d, "enable_private_vpc_default_route", gw.PrivateVpcDefaultEnabled)
		mustSet(d, "enable_skip_public_route_table_update", gw.SkipPublicVpcUpdateEnabled)
		mustSet(d, "enable_auto_advertise_s2c_cidrs", gw.AutoAdvertiseCidrsEnabled)
		mustSet(d, "spoke_bgp_manual_advertise_cidrs", gw.BgpManualSpokeAdvertiseCidrs)
		mustSet(d, "enable_bgp", gw.EnableBgp)
		mustSet(d, "enable_learned_cidrs_approval", gw.EnableLearnedCidrsApproval)
		if gw.EnableLearnedCidrsApproval {
			spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: gw.GwName})
			if err != nil {
				return fmt.Errorf("could not get advanced config for spoke gateway: %w", err)
			}

			if err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs); err != nil {
				return fmt.Errorf("could not set approved_learned_cidrs into state: %w", err)
			}
		} else {
			mustSet(d, "approved_learned_cidrs", nil)
		}
		mustSet(d, "bgp_ecmp", gw.BgpEcmp)
		mustSet(d, "enable_active_standby", gw.EnableActiveStandby)
		mustSet(d, "enable_active_standby_preemptive", gw.EnableActiveStandbyPreemptive)
		mustSet(d, "disable_route_propagation", gw.DisableRoutePropagation)
		var prependAsPath []string
		for _, p := range strings.Split(gw.PrependASPath, " ") {
			if p != "" {
				prependAsPath = append(prependAsPath, p)
			}
		}
		err = d.Set("prepend_as_path", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %w", err)
		}
		mustSet(d, "enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
		if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
			return fmt.Errorf("setting 'monitor_exclude_list' to state: %w", err)
		}

		if gw.EnableBgp {
			mustSet(d, "learned_cidrs_approval_mode", gw.LearnedCidrsApprovalMode)
			mustSet(d, "bgp_polling_time", gw.BgpPollingTime)
			mustSet(d, "bgp_hold_time", gw.BgpHoldTime)
		} else {
			mustSet(d, "learned_cidrs_approval_mode", "gateway")
			mustSet(d, "bgp_polling_time", 50)
			mustSet(d, "bgp_hold_time", 180)
		}

		if gw.EnableSpotInstance {
			mustSet(d, "enable_spot_instance", true)
			mustSet(d, "spot_price", gw.SpotPrice)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.ReuseEip, ":")
			if len(azureEip) == 3 {
				mustSet(d, "azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the Spoke Gateway %s", gw.GwName)
			}
		}
	}

	d.SetId(gateway.GwName)
	return nil
}
