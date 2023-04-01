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
				Description: "Primary gateway name.",
			},
			"compute_node_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Compute node UUID.",
			},
			"interfaces": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Interface name.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Interface type.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Bandwidth.",
						},
						"enable_dhcp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable DHCP.",
						},
						"wan_public_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "WAN interface public IP.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface static IP address.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Gateway IP.",
						},
						"dns_server_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Primary DNS server IP.",
						},
						"secondary_dns_server_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Secondary DNS server IP.",
						},
						"tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag.",
						},
					},
				},
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Edge CSP account name.",
			},
		},
	}
}

func marshalEdgeCSPHaInput(d *schema.ResourceData) *goaviatrix.EdgeCSPHa {
	edgeCSPHa := &goaviatrix.EdgeCSPHa{
		PrimaryGwName:   d.Get("primary_gw_name").(string),
		ComputeNodeUuid: d.Get("compute_node_uuid").(string),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, if0 := range interfaces {
		if1 := if0.(map[string]interface{})

		if2 := &goaviatrix.Interface{
			IfName:       if1["name"].(string),
			Type:         if1["type"].(string),
			Bandwidth:    if1["bandwidth"].(int),
			PublicIp:     if1["wan_public_ip"].(string),
			Tag:          if1["tag"].(string),
			Dhcp:         if1["enable_dhcp"].(bool),
			IpAddr:       if1["ip_address"].(string),
			GatewayIp:    if1["gateway_ip"].(string),
			DnsPrimary:   if1["dns_server_ip"].(string),
			DnsSecondary: if1["secondary_dns_server_ip"].(string),
		}

		edgeCSPHa.InterfaceList = append(edgeCSPHa.InterfaceList, if2)
	}

	return edgeCSPHa
}

func resourceAviatrixEdgeCSPHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeCSPHa := marshalEdgeCSPHaInput(d)

	edgeCSPHaName, err := client.CreateEdgeCSPHa(ctx, edgeCSPHa)
	if err != nil {
		return diag.Errorf("failed to create Edge CSP HA: %s", err)
	}

	d.SetId(edgeCSPHaName)

	gatewayForEdgeCSPFunctions := &goaviatrix.EdgeCSP{
		GwName: edgeCSPHaName,
	}

	if len(edgeCSPHa.InterfaceList) != 0 {
		gatewayForEdgeCSPFunctions.InterfaceList = edgeCSPHa.InterfaceList

		err = client.UpdateEdgeCSPHa(ctx, gatewayForEdgeCSPFunctions)
		if err != nil {
			return diag.Errorf("could not config WAN/LAN interfaces after Edge CSP HA creation: %v", err)
		}
	}

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
		return diag.Errorf("could not read Edge CSP HA: %v", err)
	}

	d.Set("primary_gw_name", edgeCSPHaResp.PrimaryGwName)
	d.Set("compute_node_uuid", edgeCSPHaResp.ComputeNodeUuid)
	d.Set("account_name", edgeCSPHaResp.AccountName)

	var interfaces []map[string]interface{}
	for _, if0 := range edgeCSPHaResp.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
		if1["bandwidth"] = if0.Bandwidth
		if1["wan_public_ip"] = if0.PublicIp
		if1["tag"] = if0.Tag
		if1["enable_dhcp"] = if0.Dhcp
		if1["ip_address"] = if0.IpAddr
		if1["gateway_ip"] = if0.GatewayIp
		if1["dns_server_ip"] = if0.DnsPrimary
		if1["secondary_dns_server_ip"] = if0.DnsSecondary

		interfaces = append(interfaces, if1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	d.SetId(edgeCSPHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeCSPHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeCSPHa := marshalEdgeCSPHaInput(d)

	d.Partial(true)

	gatewayForEdgeCSPFunctions := &goaviatrix.EdgeCSP{
		GwName: d.Id(),
	}

	if d.HasChange("interfaces") {
		gatewayForEdgeCSPFunctions.InterfaceList = edgeCSPHa.InterfaceList

		err := client.UpdateEdgeCSPHa(ctx, gatewayForEdgeCSPFunctions)
		if err != nil {
			return diag.Errorf("could not update WAN/LAN interfaces during Edge CSP HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeCSPHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeCSPHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)

	err := client.DeleteEdgeCSP(ctx, accountName, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge CSP: %v", err)
	}

	return nil
}
