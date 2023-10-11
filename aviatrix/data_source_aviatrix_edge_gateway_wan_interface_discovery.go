package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixEdgeGatewayWanInterfaceDiscovery() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixEdgeGatewayWanInterfaceDiscoveryRead,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge gateway name.",
			},
			"wan_interface_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the WAN interface to be discovered.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP of the Edge gateway's WAN interface.",
			},
		},
	}
}

func dataSourceAviatrixEdgeGatewayWanInterfaceDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	wanInterfaceName := d.Get("wan_interface_name").(string)

	ip, err := client.GetEdgeGatewayWanIp(ctx, gwName, wanInterfaceName)
	if err != nil {
		return diag.Errorf("couldn't get wan interface ip for edge gateway %s: %s", gwName, err)
	}

	d.Set("ip_address", ip)

	d.SetId(gwName + "~" + wanInterfaceName)
	return nil
}
