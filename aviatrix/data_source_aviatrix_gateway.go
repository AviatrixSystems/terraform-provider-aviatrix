package aviatrix

import (
	"fmt"
	"strings"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixGatewayRead,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_reg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_size": {
				Type:     schema.TypeString,
				Computed: true,
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
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_size", gw.GwSize)
		d.Set("vpc_net", gw.VpcNet)
	}
	d.SetId(gateway.GwName)
	return nil
}
