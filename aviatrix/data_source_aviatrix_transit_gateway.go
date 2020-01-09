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
				Description: "Range of the subnet where the transit gateway is launched.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Transit Gateway created.",
			},
		},
	}
}

func dataSourceAviatrixTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
	}

	if d.Get("account_name").(string) != "" {
		gateway.AccountName = d.Get("account_name").(string)
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		return fmt.Errorf("couldn't find Aviatrix Transit Gateway: %s", err)
	}
	if gw != nil {
		if gw.CloudType == 1 || gw.CloudType == 16 || gw.CloudType == 256 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
		} else if gw.CloudType == 4 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0])
			d.Set("vpc_reg", gw.GatewayZone)
		} else if gw.CloudType == 8 {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
		}
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("gw_size", gw.GwSize)
		d.Set("subnet", gw.VpcNet)
		d.Set("public_ip", gw.PublicIP)
	}

	d.SetId(gateway.GwName)
	return nil
}
