package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
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
			"eip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Transit Gateway created.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA Subnet. Required for enabling HA for AWS/AZURE transit gateway.",
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
			"ha_eip": {
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
			"enable_active_mesh": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable active mesh mode for Transit Gateway.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable vpc_dns_server for Gateway. Only supports AWS.",
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
					"It applies to this spoke gateway only​."},
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
					"remove those included in the 'Excluded CIDR List'.​",
			},
			"enable_transit_firenet": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Switch to enable/disable transit firenet interfaces for transit gateway.",
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
		return fmt.Errorf("couldn't find Aviatrix Transit Gateway: %s", err)
	}
	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("subnet", gw.VpcNet)

		if gw.CloudType == 1 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == 4 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0])
			d.Set("vpc_reg", gw.GatewayZone)
			d.Set("allocate_new_eip", true)
		} else if gw.CloudType == 8 || gw.CloudType == 16 {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		}

		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("eip", gw.PublicIP)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)

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

		if gw.CloudType == 1 {
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
			if gw.CloudType == 1 {
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

		if gw.EnableActiveMesh == "yes" {
			d.Set("enable_active_mesh", true)
		} else {
			d.Set("enable_active_mesh", false)
		}

		if gw.CloudType == 1 && gw.EnableVpcDnsServer == "Enabled" {
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
	}

	if gw.CloudType == 1 {
		tags := &goaviatrix.Tags{
			CloudType:    1,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
		}
		tagList, _ := client.GetTags(tags)
		if tagList != nil && len(tagList) != 0 {
			d.Set("tag_list", tagList)
		}
	}

	haGateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string) + "-hagw",
	}
	haGw, err := client.GetGateway(haGateway)
	if haGw != nil {
		if haGw.CloudType == 1 || haGw.CloudType == 8 || haGw.CloudType == 16 {
			d.Set("ha_subnet", haGw.VpcNet)
			d.Set("ha_zone", "")
		} else if haGw.CloudType == 4 {
			d.Set("ha_zone", haGw.GatewayZone)
			d.Set("ha_subnet", "")
		}
		d.Set("ha_eip", haGw.PublicIP)
		d.Set("ha_gw_size", haGw.GwSize)
		d.Set("ha_cloud_instance_id", haGw.CloudnGatewayInstID)
		d.Set("ha_gw_name", haGw.GwName)

		if haGw.InsaneMode == "yes" && haGw.CloudType == 1 {
			d.Set("ha_insane_mode_az", haGw.GatewayZone)
		} else {
			d.Set("ha_insane_mode_az", "")
		}
	}

	d.SetId(gateway.GwName)
	return nil
}
