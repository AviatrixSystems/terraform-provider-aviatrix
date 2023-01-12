package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeCSPHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeCSPHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeCSPHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeCSPHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeCSPHaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"primary_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the primary gateway.",
			},
			"dhcp": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"compute_node_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"lan_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func resourceAviatrixEdgeCSPHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeCSPHa := &goaviatrix.EdgeCSPHa{
		PrimaryGwName:   d.Get("primary_gw_name").(string),
		Dhcp:            d.Get("dhcp").(bool),
		ComputeNodeUuid: d.Get("compute_node_uuid").(string),
		LanIp:           d.Get("lan_ip").(string),
	}

	edgeCSPHaName, err := client.CreateEdgeCSPHa(ctx, edgeCSPHa)
	if err != nil {
		return diag.Errorf("failed to create Aviatrix Spoke HA Gateway: %s", err)
	}

	d.SetId(edgeCSPHaName)
	return resourceAviatrixEdgeCSPHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeCSPHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("primary_gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeCSPHaResp, err := client.GetEdgeCSPHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge CSP: %v", err)
	}

	d.Set("primary_gw_name", edgeCSPHaResp.PrimaryGwName)
	d.Set("dhcp", edgeCSPHaResp.Dhcp)
	d.Set("compute_node_uuid", edgeCSPHaResp.ComputeNodeUuid)
	d.Set("lan_ip", edgeCSPHaResp.LanIp)
	d.Set("account_name", edgeCSPHaResp.AccountName)

	d.SetId(edgeCSPHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeCSPHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := meta.(*goaviatrix.Client)
	//
	//gateway := &goaviatrix.Gateway{
	//	CloudType: d.Get("cloud_type").(int),
	//	GwName:    d.Get("gw_name").(string),
	//}
	//
	//if d.HasChange("gw_size") {
	//	gateway.GwName = d.Get("gw_name").(string)
	//	gateway.GwSize = d.Get("gw_size").(string)
	//	err := client.UpdateGateway(gateway)
	//	if err != nil {
	//		return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway %s: %s", gateway.GwName, err)
	//	}
	//}
	//
	//d.Partial(false)
	//d.SetId(gateway.GwName)
	//return resourceAviatrixSpokeHaGatewayRead(d, meta)
	return nil
}

func resourceAviatrixEdgeCSPHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	//primaryGwName := d.Get("primary_gw_name").(string)

	err := client.DeleteEdgeCSP(ctx, accountName, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge CSP: %v", err)
	}

	return nil
}
