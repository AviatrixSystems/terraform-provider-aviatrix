package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
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
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Description: "Range of the subnet where the spoke gateway is launched.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the spoke gateway created.",
			},
			"allocate_new_eip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the eip is newly allocated or not.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable this feature.",
			},
			"transit_gw": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The transit gateway that the spoke gateway is attached to.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Instance tag of cloud provider. Only supported for AWS provider.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable Insane Mode for Spoke Gateway.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.",
			},
			"enable_active_mesh": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable Active Mesh Mode for Spoke Gateway.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disalbe vpc_dns_server for Gateway.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider.",
			},
			"customized_routes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of comma separated CIDRs to be customized for the spoke VPC.",
			},
			"filtered_routes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table.",
			},
			"customized_routes_advertisement": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of comma separated CIDRs to be excluded from being advertised to.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
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
		return fmt.Errorf("couldn't find Aviatrix spoke gateway: %s", err)
	}
	if gw != nil {
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
		} else if gw.CloudType == 8 {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		}
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("gw_size", gw.GwSize)
		d.Set("subnet", gw.VpcNet)
		d.Set("public_ip", gw.PublicIP)
		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if gw.SpokeVpc == "yes" {
			d.Set("transit_gw", gw.TransitGwName)
		} else {
			d.Set("transit_gw", "")
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

		if len(gw.CustomizedRoutes) != 0 {
			d.Set("customized_routes", strings.Join(gw.CustomizedRoutes, ","))
		} else {
			d.Set("customized_routes", "")
		}

		if len(gw.FilteredRoutes) != 0 {
			d.Set("filtered_routes", strings.Join(gw.FilteredRoutes, ","))
		} else {
			d.Set("filtered_routes", "")
		}

		if len(gw.CustomizedRoutesAdvertisement) != 0 {
			d.Set("customized_routes_advertisement", strings.Join(gw.CustomizedRoutesAdvertisement, ","))
		} else {
			d.Set("customized_routes_advertisement", "")
		}

		if gw.CloudType == 1 {
			tags := &goaviatrix.Tags{
				CloudType:    1,
				ResourceType: "gw",
				ResourceName: d.Get("gw_name").(string),
			}
			tagList, _ := client.GetTags(tags)
			d.Set("tag_list", tagList)
		}
	}

	d.SetId(gateway.GwName)
	return nil
}
