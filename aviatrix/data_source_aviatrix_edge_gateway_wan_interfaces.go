package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixEdgeGatewayWanInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixEdgeGatewayWanInterfacesRead,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge gateway name.",
			},
			"wan_interfaces": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Set of Edge WAN interfaces.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAviatrixEdgeGatewayWanInterfacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)

	edgeResp, err := client.GetEdgeGatewayWanInterfaces(ctx, gwName)
	if err != nil {
		return diag.Errorf("couldn't get wan interfaces for edge gateway %s: %s", gwName, err)
	}

	var wanInterfaces []string
	for _, if0 := range edgeResp.InterfaceList {
		if if0.Type == "WAN" {
			wanInterfaces = append(wanInterfaces, if0.IfName)
		}
	}

	if err = d.Set("wan_interfaces", wanInterfaces); err != nil {
		return diag.Errorf("failed to set wan_interfaces: %s\n", err)
	}

	d.SetId(gwName)
	return nil
}
