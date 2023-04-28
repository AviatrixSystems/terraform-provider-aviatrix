package aviatrix

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "Edge Equinix account name.",
			},
		},
	}
}

func marshalEdgeEquinixHaInput(d *schema.ResourceData) *goaviatrix.EdgeEquinixHa {
	edgeEquinixHa := &goaviatrix.EdgeEquinixHa{
		PrimaryGwName:            d.Get("primary_gw_name").(string),
		ZtpFileDownloadPath:      d.Get("ztp_file_download_path").(string),
		ManagementEgressIpPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, interface0 := range interfaces {
		interface1 := interface0.(map[string]interface{})

		interface2 := &goaviatrix.EdgeEquinixInterface{
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

		edgeEquinixHa.InterfaceList = append(edgeEquinixHa.InterfaceList, interface2)
	}

	return edgeEquinixHa
}

func resourceAviatrixEdgeEquinixHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeEquinixHa := marshalEdgeEquinixHaInput(d)

	edgeEquinixHaName, err := client.CreateEdgeEquinixHa(ctx, edgeEquinixHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Equinix HA: %s", err)
	}

	d.SetId(edgeEquinixHaName)
	return resourceAviatrixEdgeEquinixHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeEquinixHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("primary_gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeEquinixHaResp, err := client.GetEdgeEquinixHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Equinix HA: %v", err)
	}

	d.Set("primary_gw_name", edgeEquinixHaResp.PrimaryGwName)
	d.Set("account_name", edgeEquinixHaResp.AccountName)

	if edgeEquinixHaResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeEquinixHaResp.ManagementEgressIpPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, interface0 := range edgeEquinixHaResp.InterfaceList {
		interface1 := make(map[string]interface{})
		interface1["name"] = interface0.IfName
		interface1["type"] = interface0.Type
		interface1["bandwidth"] = interface0.Bandwidth
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
	client := meta.(*goaviatrix.Client)

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
	client := meta.(*goaviatrix.Client)

	edgeEquinixHa := marshalEdgeEquinixHaInput(d)
	accountName := d.Get("account_name").(string)

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
