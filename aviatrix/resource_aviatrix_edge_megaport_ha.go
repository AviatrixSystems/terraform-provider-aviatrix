package aviatrix

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixEdgeMegaportHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeMegaportHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeMegaportHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeMegaportHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeMegaportHaDelete,
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
						"logical_ifname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logical interface name e.g., wan0, lan0, mgmt0.",
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^(wan|lan|mgmt)[0-9]+$`),
								"Logical interface name must start with 'wan', 'lan', or 'mgmt' followed by a number (e.g., 'wan0', 'lan1', 'mgmt2').",
							),
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
				Description: "Edge Megaport account name.",
			},
		},
	}
}

func marshalEdgeMegaportHaInput(d *schema.ResourceData) *goaviatrix.EdgeMegaportHa {
	edgeMegaportHa := &goaviatrix.EdgeMegaportHa{
		PrimaryGwName:            d.Get("primary_gw_name").(string),
		ZtpFileDownloadPath:      d.Get("ztp_file_download_path").(string),
		ManagementEgressIpPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, interface0 := range interfaces {
		interface1 := interface0.(map[string]interface{})

		interface2 := &goaviatrix.EdgeMegaportInterface{
			LogicalInterfaceName: interface1["logical_ifname"].(string),
			PublicIp:             interface1["wan_public_ip"].(string),
			Tag:                  interface1["tag"].(string),
			Dhcp:                 interface1["enable_dhcp"].(bool),
			IpAddr:               interface1["ip_address"].(string),
			GatewayIp:            interface1["gateway_ip"].(string),
			DnsPrimary:           interface1["dns_server_ip"].(string),
			DnsSecondary:         interface1["secondary_dns_server_ip"].(string),
		}

		edgeMegaportHa.InterfaceList = append(edgeMegaportHa.InterfaceList, interface2)
	}

	return edgeMegaportHa
}

func resourceAviatrixEdgeMegaportHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeMegaportHa := marshalEdgeMegaportHaInput(d)

	edgeMegaportHaName, err := client.CreateEdgeMegaportHa(ctx, edgeMegaportHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Megaport HA: %s", err)
	}

	d.SetId(edgeMegaportHaName)
	return resourceAviatrixEdgeMegaportHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeMegaportHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("primary_gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeMegaportHaResp, err := client.GetEdgeMegaportHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Megaport HA: %v", err)
	}

	d.Set("primary_gw_name", edgeMegaportHaResp.PrimaryGwName)
	d.Set("account_name", edgeMegaportHaResp.AccountName)

	if edgeMegaportHaResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeMegaportHaResp.ManagementEgressIpPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, interface0 := range edgeMegaportHaResp.InterfaceList {
		interface1 := make(map[string]interface{})
		interface1["logical_ifname"] = interface0.LogicalInterfaceName
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

	d.SetId(edgeMegaportHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeMegaportHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeMegaportHa := marshalEdgeMegaportHaInput(d)

	d.Partial(true)

	gatewayForEdgeMegaportFunctions := &goaviatrix.EdgeMegaport{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeMegaportFunctions.InterfaceList = edgeMegaportHa.InterfaceList
		gatewayForEdgeMegaportFunctions.ManagementEgressIpPrefix = edgeMegaportHa.ManagementEgressIpPrefix

		err := client.UpdateEdgeMegaportHa(ctx, gatewayForEdgeMegaportFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge Megaport HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeMegaportHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeMegaportHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeMegaportHa := marshalEdgeMegaportHaInput(d)
	accountName := d.Get("account_name").(string)

	err := client.DeleteEdgeMegaport(ctx, accountName, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge Megaport HA: %v", err)
	}

	fileName := edgeMegaportHa.ZtpFileDownloadPath + "/" + edgeMegaportHa.PrimaryGwName + "-hagw-cloud-init.txt"

	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
