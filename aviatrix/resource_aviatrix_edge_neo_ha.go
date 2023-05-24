package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeNEOHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeNEOHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeNEOHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeNEOHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeNEOHaDelete,
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
			"device_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge NEO device ID.",
			},
			"interfaces": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "WAN/LAN/MANAGEMENT interfaces.",
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
							Description: "The rate of data can be moved through the interface, requires an integer value. Unit is in Mb/s.",
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
			"management_egress_ip_prefix_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of management egress gateway IP/prefix.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Edge NEO account name.",
			},
		},
		DeprecationMessage: "Since V3.1.1+, please use resource aviatrix_edge_platform_ha instead. Resource " +
			"aviatrix_edge_neo_ha will be deprecated in the V3.2.0 release.",
	}
}

func marshalEdgeNEOHaInput(d *schema.ResourceData) *goaviatrix.EdgeNEOHa {
	edgeNEOHa := &goaviatrix.EdgeNEOHa{
		PrimaryGwName:            d.Get("primary_gw_name").(string),
		DeviceId:                 d.Get("device_id").(string),
		ManagementEgressIpPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, interface0 := range interfaces {
		interface1 := interface0.(map[string]interface{})

		interface2 := &goaviatrix.EdgeNEOInterface{
			IfName:       interface1["name"].(string),
			Type:         interface1["type"].(string),
			Bandwidth:    interface1["bandwidth"].(int),
			PublicIp:     interface1["wan_public_ip"].(string),
			Tag:          interface1["tag"].(string),
			Dhcp:         interface1["enable_dhcp"].(bool),
			IpAddr:       interface1["ip_address"].(string),
			GatewayIp:    interface1["gateway_ip"].(string),
			DnsPrimary:   interface1["dns_server_ip"].(string),
			DnsSecondary: interface1["secondary_dns_server_ip"].(string),
		}

		edgeNEOHa.InterfaceList = append(edgeNEOHa.InterfaceList, interface2)
	}

	return edgeNEOHa
}

func resourceAviatrixEdgeNEOHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeNEOHa := marshalEdgeNEOHaInput(d)

	edgeNEOHaName, err := client.CreateEdgeNEOHa(ctx, edgeNEOHa)
	if err != nil {
		return diag.Errorf("failed to create Edge NEO HA: %s", err)
	}

	d.SetId(edgeNEOHaName)
	return resourceAviatrixEdgeNEOHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeNEOHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("primary_gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeNEOHaResp, err := client.GetEdgeNEOHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge NEO HA: %v", err)
	}

	d.Set("primary_gw_name", edgeNEOHaResp.PrimaryGwName)
	d.Set("device_id", edgeNEOHaResp.DeviceId)
	d.Set("account_name", edgeNEOHaResp.AccountName)

	if edgeNEOHaResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeNEOHaResp.ManagementEgressIpPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, if0 := range edgeNEOHaResp.InterfaceList {
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

	d.SetId(edgeNEOHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeNEOHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeNEOHa := marshalEdgeNEOHaInput(d)

	d.Partial(true)

	gatewayForEdgeNEOFunctions := &goaviatrix.EdgeNEO{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeNEOFunctions.InterfaceList = edgeNEOHa.InterfaceList
		gatewayForEdgeNEOFunctions.ManagementEgressIpPrefix = edgeNEOHa.ManagementEgressIpPrefix

		err := client.UpdateEdgeNEOHa(ctx, gatewayForEdgeNEOFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge NEO HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeNEOHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeNEOHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)

	err := client.DeleteEdgeNEO(ctx, accountName, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge NEO: %v", err)
	}

	return nil
}
