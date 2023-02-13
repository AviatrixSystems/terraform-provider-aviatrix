package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
			"compute_node_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"management_interface_config": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Management interface configuration. Valid values: 'DHCP' and 'Static'.",
				ValidateFunc: validation.StringInSlice([]string{"DHCP", "Static"}, false),
			},
			"lan_interface_ip_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "LAN interface IP/prefix.",
			},
			"interfaces": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ifname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "",
						},
						"ipaddr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"dns_primary": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"dns_secondary": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"admin_state": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "",
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
		PrimaryGwName:             d.Get("primary_gw_name").(string),
		ComputeNodeUuid:           d.Get("compute_node_uuid").(string),
		ManagementInterfaceConfig: d.Get("management_interface_config").(string),
		LanInterfaceIpPrefix:      d.Get("lan_interface_ip_prefix").(string),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, if0 := range interfaces {
		if1 := if0.(map[string]interface{})

		if2 := &goaviatrix.Interface{
			IfName:       if1["ifname"].(string),
			Type:         if1["type"].(string),
			Bandwidth:    if1["bandwidth"].(int),
			PublicIp:     if1["public_ip"].(string),
			Tag:          if1["tag"].(string),
			Dhcp:         if1["dhcp"].(bool),
			IpAddr:       if1["ipaddr"].(string),
			GatewayIp:    if1["gateway_ip"].(string),
			DnsPrimary:   if1["dns_primary"].(string),
			DnsSecondary: if1["dns_secondary"].(string),
		}

		if if1["admin_state"].(bool) {
			if2.AdminState = "enabled"
		} else {
			if2.AdminState = "disabled"
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

		err = client.UpdateEdgeCSP(ctx, gatewayForEdgeCSPFunctions)
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

	if edgeCSPHaResp.Dhcp {
		d.Set("management_interface_config", "DHCP")
	} else {
		d.Set("management_interface_config", "Static")
	}

	var interfaces []map[string]interface{}
	for _, if0 := range edgeCSPHaResp.InterfaceList {
		if if0.Type != "MANAGEMENT" {
			if1 := make(map[string]interface{})
			if1["ifname"] = if0.IfName
			if1["type"] = if0.Type
			if1["bandwidth"] = if0.Bandwidth
			if1["public_ip"] = if0.PublicIp
			if1["tag"] = if0.Tag
			if1["dhcp"] = if0.Dhcp
			if1["ipaddr"] = if0.IpAddr
			if1["gateway_ip"] = if0.GatewayIp
			if1["dns_primary"] = if0.DnsPrimary
			if1["dns_secondary"] = if0.DnsSecondary

			if if0.AdminState == "enabled" {
				if1["admin_state"] = true
			} else {
				if1["admin_state"] = false
			}

			interfaces = append(interfaces, if1)
		}
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

	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: d.Id(),
	}
	gatewayForEdgeCSPFunctions := &goaviatrix.EdgeCSP{
		GwName: d.Id(),
	}

	if d.HasChanges("lan_interface_ip_prefix") {
		gatewayForEaasFunctions.LanInterfaceIpPrefix = edgeCSPHa.LanInterfaceIpPrefix

		err := client.UpdateEdgeSpokeIpConfigurations(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update IP configurations during Edge CSP HA update: %v", err)
		}
	}

	if d.HasChange("interfaces") {
		gatewayForEdgeCSPFunctions.InterfaceList = edgeCSPHa.InterfaceList

		err := client.UpdateEdgeCSP(ctx, gatewayForEdgeCSPFunctions)
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
