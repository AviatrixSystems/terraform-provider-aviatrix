package aviatrix

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixEdgeZededaHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeZededaHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeZededaHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeZededaHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeZededaHaDelete,
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
							Deprecated:  "Bandwidth will be removed in a future release.",
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
				Description: "Edge Zededa account name.",
			},
		},
	}
}

func marshalEdgeZededaHaInput(d *schema.ResourceData) *goaviatrix.EdgeCSPHa {
	edgeCSPHa := &goaviatrix.EdgeCSPHa{
		PrimaryGwName:            getString(d, "primary_gw_name"),
		ComputeNodeUuid:          getString(d, "compute_node_uuid"),
		ManagementEgressIpPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := getSet(d, "interfaces").List()
	for _, if0 := range interfaces {
		if1 := mustMap(if0)

		if2 := &goaviatrix.Interface{
			IfName:       mustString(if1["name"]),
			Type:         mustString(if1["type"]),
			PublicIp:     mustString(if1["wan_public_ip"]),
			Tag:          mustString(if1["tag"]),
			Dhcp:         mustBool(if1["enable_dhcp"]),
			IpAddr:       mustString(if1["ip_address"]),
			GatewayIp:    mustString(if1["gateway_ip"]),
			DnsPrimary:   mustString(if1["dns_server_ip"]),
			DnsSecondary: mustString(if1["secondary_dns_server_ip"]),
		}

		edgeCSPHa.InterfaceList = append(edgeCSPHa.InterfaceList, if2)
	}

	return edgeCSPHa
}

func resourceAviatrixEdgeZededaHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeCSPHa := marshalEdgeZededaHaInput(d)

	edgeCSPHaName, err := client.CreateEdgeCSPHa(ctx, edgeCSPHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Zededa HA: %s", err)
	}

	d.SetId(edgeCSPHaName)
	return resourceAviatrixEdgeZededaHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeZededaHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if getString(d, "primary_gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		mustSet(d, "primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeCSPHaResp, err := client.GetEdgeCSPHa(ctx, getString(d, "primary_gw_name")+"-hagw")
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Zededa HA: %v", err)
	}
	mustSet(d, "primary_gw_name", edgeCSPHaResp.PrimaryGwName)
	mustSet(d, "compute_node_uuid", edgeCSPHaResp.ComputeNodeUuid)
	mustSet(d, "account_name", edgeCSPHaResp.AccountName)

	if edgeCSPHaResp.ManagementEgressIpPrefix == "" {
		mustSet(d, "management_egress_ip_prefix_list", nil)
	} else {
		mustSet(d, "management_egress_ip_prefix_list", strings.Split(edgeCSPHaResp.ManagementEgressIpPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, if0 := range edgeCSPHaResp.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
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

func resourceAviatrixEdgeZededaHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeCSPHa := marshalEdgeZededaHaInput(d)

	d.Partial(true)

	gatewayForEdgeCSPFunctions := &goaviatrix.EdgeCSP{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeCSPFunctions.InterfaceList = edgeCSPHa.InterfaceList
		gatewayForEdgeCSPFunctions.ManagementEgressIpPrefix = edgeCSPHa.ManagementEgressIpPrefix

		err := client.UpdateEdgeCSPHa(ctx, gatewayForEdgeCSPFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge Zededa HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeZededaHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeZededaHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	accountName := getString(d, "account_name")

	err := client.DeleteEdgeCSP(ctx, accountName, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge Zededa HA %s: %v", d.Id(), err)
	}

	return nil
}
