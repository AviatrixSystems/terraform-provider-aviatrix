package aviatrix

import (
	"context"
	"log"
	"strconv"
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
				Type:             schema.TypeList,
				Required:         true,
				Description:      "",
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncInterfaces,
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
						"vrrp_state": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "",
						},
					},
				},
			},
			"vlan": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "",
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncVlan,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_interface": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Required:    true,
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
						"admin_state": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "",
						},
						"peer_ipaddr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"peer_gateway_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"virtual_ip": {
							Type:        schema.TypeString,
							Required:    true,
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

	interfaces := d.Get("interfaces").([]interface{})
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
			VrrpState:    if1["vrrp_state"].(bool),
		}

		if if1["admin_state"].(bool) {
			if2.AdminState = "enabled"
		} else {
			if2.AdminState = "disabled"
		}

		edgeCSPHa.InterfaceList = append(edgeCSPHa.InterfaceList, if2)
	}

	vlan := d.Get("vlan").([]interface{})
	for _, v0 := range vlan {
		v1 := v0.(map[string]interface{})

		v2 := &goaviatrix.Vlan{
			ParentInterface: v1["parent_interface"].(string),
			IpAddr:          v1["ipaddr"].(string),
			GatewayIp:       v1["gateway_ip"].(string),
			PeerIpAddr:      v1["peer_ipaddr"].(string),
			PeerGatewayIp:   v1["peer_gateway_ip"].(string),
			VirtualIp:       v1["virtual_ip"].(string),
		}

		v2.VlanId = strconv.Itoa(v1["vlan_id"].(int))

		if v1["admin_state"].(bool) {
			v2.AdminState = "enabled"
		} else {
			v2.AdminState = "disabled"
		}

		edgeCSPHa.VlanList = append(edgeCSPHa.VlanList, v2)
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

	if len(edgeCSPHa.InterfaceList) != 0 || len(edgeCSPHa.VlanList) != 0 {
		gatewayForEdgeCSPFunctions.InterfaceList = edgeCSPHa.InterfaceList
		gatewayForEdgeCSPFunctions.VlanList = edgeCSPHa.VlanList

		err = client.UpdateEdgeCSP(ctx, gatewayForEdgeCSPFunctions)
		if err != nil {
			return diag.Errorf("could not config WAN/LAN/VLAN after Edge CSP HA creation: %v", err)
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
	var vlan []map[string]interface{}
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

			if if0.Type == "LAN" {
				if1["vrrp_state"] = if0.VrrpState
			}

			if if0.Type == "LAN" && if0.SubInterfaces != nil {
				for _, v0 := range if0.SubInterfaces {
					v1 := make(map[string]interface{})
					v1["parent_interface"] = v0.ParentInterface
					v1["ipaddr"] = v0.IpAddr
					v1["gateway_ip"] = v0.GatewayIp
					v1["peer_ipaddr"] = v0.PeerIpAddr
					v1["peer_gateway_ip"] = v0.PeerGatewayIp
					v1["virtual_ip"] = v0.VirtualIp

					vlandid, _ := strconv.Atoi(v0.VlanId)
					v1["vlan_id"] = vlandid

					if v0.AdminState == "enabled" {
						v1["admin_state"] = true
					} else {
						v1["admin_state"] = false
					}

					vlan = append(vlan, v1)
				}
			}

			interfaces = append(interfaces, if1)
		}
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	if err = d.Set("vlan", vlan); err != nil {
		return diag.Errorf("failed to set vlan: %s\n", err)
	}

	d.SetId(edgeCSPHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeCSPHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeCSPHa := marshalEdgeCSPHaInput(d)

	if d.HasChange("vlan") {
		vlanOld, _ := d.GetChange("vlan")
		if len(vlanOld.([]interface{})) != 0 {
			return diag.Errorf("vlan is not allowed to be updated")
		}
	}

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

	if d.HasChange("interfaces") || d.HasChange("vlan") {
		gatewayForEdgeCSPFunctions.InterfaceList = edgeCSPHa.InterfaceList
		gatewayForEdgeCSPFunctions.VlanList = edgeCSPHa.VlanList

		err := client.UpdateEdgeCSP(ctx, gatewayForEdgeCSPFunctions)
		if err != nil {
			return diag.Errorf("could not update WAN/LAN/VLAN interfaces during Edge CSP HA update: %v", err)
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
