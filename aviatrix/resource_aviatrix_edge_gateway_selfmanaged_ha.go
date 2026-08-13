package aviatrix

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixEdgeGatewaySelfmanagedHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeGatewaySelfmanagedHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeGatewaySelfmanagedHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeGatewaySelfmanagedHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeGatewaySelfmanagedHaDelete,
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
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"ztp_file_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "ZTP file type.",
				ValidateFunc: validation.StringInSlice([]string{"iso", "cloud-init"}, false),
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
			"dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "DNS server IP.",
				ValidateFunc: validation.IsIPAddress,
				Deprecated:   "DNS server ip attribute will be removed in the future release.",
			},
			"secondary_dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Secondary DNS server IP.",
				ValidateFunc: validation.IsIPAddress,
				Deprecated:   "Secondary DNS server ip attribute will be removed in the future release.",
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
		},
	}
}

func marshalEdgeGatewaySelfmanagedHaInput(d *schema.ResourceData) *goaviatrix.EdgeVmSelfmanagedHa {
	edgeGatewaySelfmanagedHa := &goaviatrix.EdgeVmSelfmanagedHa{
		PrimaryGwName:            getString(d, "primary_gw_name"),
		SiteId:                   getString(d, "site_id"),
		ZtpFileType:              getString(d, "ztp_file_type"),
		ZtpFileDownloadPath:      getString(d, "ztp_file_download_path"),
		DnsServerIp:              getString(d, "dns_server_ip"),
		SecondaryDnsServerIp:     getString(d, "secondary_dns_server_ip"),
		ManagementEgressIPPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := getSet(d, "interfaces").List()
	for _, if0 := range interfaces {
		if1 := mustMap(if0)

		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:       mustString(if1["name"]),
			Type:         mustString(if1["type"]),
			PublicIp:     mustString(if1["wan_public_ip"]),
			Dhcp:         mustBool(if1["enable_dhcp"]),
			IpAddr:       mustString(if1["ip_address"]),
			GatewayIp:    mustString(if1["gateway_ip"]),
			DNSPrimary:   mustString(if1["dns_server_ip"]),
			DNSSecondary: mustString(if1["secondary_dns_server_ip"]),
			IPv6Addr:     mustString(if1["ipv6_address"]),
			GatewayIPv6:  mustString(if1["gateway_ipv6"]),
		}

		edgeGatewaySelfmanagedHa.InterfaceList = append(edgeGatewaySelfmanagedHa.InterfaceList, if2)
	}

	return edgeGatewaySelfmanagedHa
}

func resourceAviatrixEdgeGatewaySelfmanagedHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeGatewaySelfmanagedHa := marshalEdgeGatewaySelfmanagedHaInput(d)

	edgeGatewaySelfmanagedHaName, err := client.CreateEdgeVmSelfmanagedHa(ctx, edgeGatewaySelfmanagedHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Gateway Selfmanaged HA: %s", err)
	}

	d.SetId(edgeGatewaySelfmanagedHaName)
	return resourceAviatrixEdgeGatewaySelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeGatewaySelfmanagedHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if getString(d, "primary_gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		mustSet(d, "primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeGatewaySelfmanagedHaResp, err := client.GetEdgeVmSelfmanagedHa(ctx, getString(d, "primary_gw_name")+"-hagw")
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Gateway Selfmanaged HA: %v", err)
	}

	_ = d.Set("primary_gw_name", edgeGatewaySelfmanagedHaResp.PrimaryGwName)
	_ = d.Set("site_id", edgeGatewaySelfmanagedHaResp.SiteID)
	_ = d.Set("dns_server_ip", edgeGatewaySelfmanagedHaResp.DNSServerIP)
	_ = d.Set("secondary_dns_server_ip", edgeGatewaySelfmanagedHaResp.SecondaryDNSServerIP)

	if edgeGatewaySelfmanagedHaResp.ZtpFileType == "iso" || edgeGatewaySelfmanagedHaResp.ZtpFileType == "cloud-init" {
		mustSet(d, "ztp_file_type", edgeGatewaySelfmanagedHaResp.ZtpFileType)
	}

	if edgeGatewaySelfmanagedHaResp.ZtpFileType == "cloud_init" {
		mustSet(d, "ztp_file_type", "cloud-init")
	}

	if edgeGatewaySelfmanagedHaResp.ManagementEgressIPPrefix == "" {
		_ = d.Set("management_egress_ip_prefix_list", nil)
	} else {
		_ = d.Set("management_egress_ip_prefix_list", strings.Split(edgeGatewaySelfmanagedHaResp.ManagementEgressIPPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, if0 := range edgeGatewaySelfmanagedHaResp.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
		if1["wan_public_ip"] = if0.PublicIp
		if1["enable_dhcp"] = if0.Dhcp
		if1["ip_address"] = if0.IpAddr
		if1["gateway_ip"] = if0.GatewayIp
		if1["dns_server_ip"] = if0.DNSPrimary
		if1["secondary_dns_server_ip"] = if0.DNSSecondary
		if1["ipv6_address"] = if0.IPv6Addr
		if1["gateway_ipv6"] = if0.GatewayIPv6

		interfaces = append(interfaces, if1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	d.SetId(edgeGatewaySelfmanagedHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeGatewaySelfmanagedHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeGatewaySelfmanagedHa := marshalEdgeGatewaySelfmanagedHaInput(d)

	d.Partial(true)

	gatewayForEdgeGatewaySelfmanagedFunctions := &goaviatrix.EdgeSpoke{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeGatewaySelfmanagedFunctions.InterfaceList = edgeGatewaySelfmanagedHa.InterfaceList
		gatewayForEdgeGatewaySelfmanagedFunctions.ManagementEgressIpPrefix = edgeGatewaySelfmanagedHa.ManagementEgressIPPrefix

		err := client.UpdateEdgeVmSelfmanagedHa(ctx, gatewayForEdgeGatewaySelfmanagedFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge Gateway Selfmanaged HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeGatewaySelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeGatewaySelfmanagedHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteEdgeSpoke(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge Gateway Selfmanaged HA %s: %v", d.Id(), err)
	}

	edgeGatewaySelfmanagedHa := marshalEdgeGatewaySelfmanagedHaInput(d)

	var fileName string
	if edgeGatewaySelfmanagedHa.ZtpFileType == "iso" {
		fileName = edgeGatewaySelfmanagedHa.ZtpFileDownloadPath + "/" + edgeGatewaySelfmanagedHa.PrimaryGwName + "-" + edgeGatewaySelfmanagedHa.SiteId + "-ha.iso"
	} else {
		fileName = edgeGatewaySelfmanagedHa.ZtpFileDownloadPath + "/" + edgeGatewaySelfmanagedHa.PrimaryGwName + "-" + edgeGatewaySelfmanagedHa.SiteId + "-ha-cloud-init.txt"
	}

	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
