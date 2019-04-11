package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixGatewayRead,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account name. This can be used for logging in to CloudN console or UserConnect controller.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway name. This can be used for getting gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS VPC ID.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS VPC Region.",
			},
			"vpc_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance type.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Gateway created.",
			},
		},
	}
}

func dataSourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
	}
	if d.Get("account_name").(string) != "" {
		gateway.AccountName = d.Get("account_name").(string)
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		return fmt.Errorf("couldn't find Aviatrix Gateway: %s", err)
	}
	if gw != nil {
		index := strings.Index(gw.VpcID, "~~")
		if index > 0 {
			gw.VpcID = gw.VpcID[:index]
		}
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_size", gw.GwSize)
		d.Set("vpc_net", gw.VpcNet)
		d.Set("public_ip", gw.PublicIP)
	}
	d.SetId(gateway.GwName)
	return nil
}
