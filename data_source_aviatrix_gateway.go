package aviatrix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
)

func dataSourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixGatewayRead,

		Schema: map[string]*schema.Schema {
			"cloud_type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_reg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_net": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		AccountName:  d.Get("account_name").(string),
		GwName:       d.Get("gw_name").(string),
	}
	gw, err := client.GetGateway(gateway)
	if err != nil {
		return fmt.Errorf("Couldn't find Aviatrix Gateway: %s", err)
	}
	if gw != nil {
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_size", gw.VpcSize)
		d.Set("vpc_net", gw.VpcNet)
	}
	return nil
}
