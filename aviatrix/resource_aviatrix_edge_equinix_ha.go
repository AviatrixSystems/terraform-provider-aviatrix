package aviatrix

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixEdgeEquinixHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeEquinixHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeEquinixHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeEquinixHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeEquinixHaDelete,
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
			"ztp_file_download_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
				},
				Description: "The location where the ZTP file will be stored.",
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
						"ipv6_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface static IPv6 address.",
						},
						"gateway_ipv6": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Gateway IPv6 IP.",
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
				Description: "Edge Equinix account name.",
			},
		},
	}
}

func marshalEdgeEquinixHaInput(d *schema.ResourceData) *goaviatrix.EdgeEquinixHa {
	edgeEquinixHa := &goaviatrix.EdgeEquinixHa{
		PrimaryGwName:            getString(d, "primary_gw_name"),
		ZtpFileDownloadPath:      getString(d, "ztp_file_download_path"),
		ManagementEgressIpPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := getSet(d, "interfaces").List()
	for _, interface0 := range interfaces {
		interface1 := mustMap(interface0)

		interface2 := &goaviatrix.EdgeEquinixInterface{
			IfName:       mustString(interface1["name"]),
			Type:         mustString(interface1["type"]),
			PublicIp:     mustString(interface1["wan_public_ip"]),
			Tag:          mustString(interface1["tag"]),
			Dhcp:         mustBool(interface1["enable_dhcp"]),
			IpAddr:       mustString(interface1["ip_address"]),
			GatewayIp:    mustString(interface1["gateway_ip"]),
			DnsPrimary:   mustString(interface1["dns_server_ip"]),
			DnsSecondary: mustString(interface1["secondary_dns_server_ip"]),
		}
		if v, ok := interface1["ipv6_address"]; ok && v != nil {
			ip := mustString(v)
			if ip != "" {
				interface2.IPv6Addr = ip

				// gateway_ipv6 only makes sense if ipv6_address is set
				if gwv, ok := interface1["gateway_ipv6"]; ok && gwv != nil {
					gw := mustString(gwv)
					if gw != "" {
						interface2.GatewayIPv6IP = gw
					}
				}
			}
		}

		edgeEquinixHa.InterfaceList = append(edgeEquinixHa.InterfaceList, interface2)
	}

	return edgeEquinixHa
}

func resourceAviatrixEdgeEquinixHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeEquinixHa := marshalEdgeEquinixHaInput(d)

	edgeEquinixHaName, err := client.CreateEdgeEquinixHa(ctx, edgeEquinixHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Equinix HA: %s", err)
	}

	d.SetId(edgeEquinixHaName)
	return resourceAviatrixEdgeEquinixHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeEquinixHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if getString(d, "primary_gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		mustSet(d, "primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeEquinixHaResp, err := client.GetEdgeEquinixHa(ctx, getString(d, "primary_gw_name")+"-hagw")
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Equinix HA: %v", err)
	}
	mustSet(d, "primary_gw_name", edgeEquinixHaResp.PrimaryGwName)
	mustSet(d, "account_name", edgeEquinixHaResp.AccountName)

	if edgeEquinixHaResp.ManagementEgressIpPrefix == "" {
		mustSet(d, "management_egress_ip_prefix_list", nil)
	} else {
		mustSet(d, "management_egress_ip_prefix_list", strings.Split(edgeEquinixHaResp.ManagementEgressIpPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, interface0 := range edgeEquinixHaResp.InterfaceList {
		interface1 := make(map[string]interface{})
		interface1["name"] = interface0.IfName
		interface1["type"] = interface0.Type
		interface1["wan_public_ip"] = interface0.PublicIp
		interface1["tag"] = interface0.Tag
		interface1["enable_dhcp"] = interface0.Dhcp
		interface1["ip_address"] = interface0.IpAddr
		interface1["gateway_ip"] = interface0.GatewayIp
		interface1["dns_server_ip"] = interface0.DnsPrimary
		interface1["secondary_dns_server_ip"] = interface0.DnsSecondary

		interfaces = append(interfaces, interface1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	d.SetId(edgeEquinixHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeEquinixHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeEquinixHa := marshalEdgeEquinixHaInput(d)

	d.Partial(true)

	gatewayForEdgeEquinixFunctions := &goaviatrix.EdgeEquinix{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeEquinixFunctions.InterfaceList = edgeEquinixHa.InterfaceList
		gatewayForEdgeEquinixFunctions.ManagementEgressIpPrefix = edgeEquinixHa.ManagementEgressIpPrefix

		err := client.UpdateEdgeEquinixHa(ctx, gatewayForEdgeEquinixFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge Equinix HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeEquinixHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeEquinixHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeEquinixHa := marshalEdgeEquinixHaInput(d)
	accountName := getString(d, "account_name")

	err := client.DeleteEdgeEquinix(ctx, accountName, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge Equinix HA: %v", err)
	}

	fileName := edgeEquinixHa.ZtpFileDownloadPath + "/" + edgeEquinixHa.PrimaryGwName + "-hagw-cloud-init.txt"

	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
