package aviatrix

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"custom_interface_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of custom interface mappings containing logical interfaces mapped to mac addresses or pci id's.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logical_ifname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logical interface name e.g., wan0, mgmt0, lan0.",
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^(wan|mgmt|lan)[0-9]+$`),
								"Logical interface name must start with 'wan', 'lan' or 'mgmt' followed by a number (e.g., 'wan0', 'lan0', 'mgmt0').",
							),
						},
						"identifier_type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Type of identifier used to map the logical interface to the physical interface e.g., mac, pci, system-assigned.",
							ValidateFunc: validation.StringInSlice([]string{"mac", "pci", "system-assigned"}, false),
						},
						"identifier_value": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Value of the identifier used to map the logical interface to the physical interface. Can be a MAC address, PCI ID, or auto if system-assigned.",
							ValidateFunc: validateIdentifierValue,
						},
					},
				},
			},
		},
	}
}

func marshalEdgeGatewaySelfmanagedHaInput(d *schema.ResourceData) *goaviatrix.EdgeVmSelfmanagedHa {
	edgeGatewaySelfmanagedHa := &goaviatrix.EdgeVmSelfmanagedHa{
		PrimaryGwName:            d.Get("primary_gw_name").(string),
		SiteId:                   d.Get("site_id").(string),
		ZtpFileType:              d.Get("ztp_file_type").(string),
		ZtpFileDownloadPath:      d.Get("ztp_file_download_path").(string),
		DnsServerIp:              d.Get("dns_server_ip").(string),
		SecondaryDnsServerIp:     d.Get("secondary_dns_server_ip").(string),
		ManagementEgressIPPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, if0 := range interfaces {
		if1 := if0.(map[string]interface{})

		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:       if1["name"].(string),
			Type:         if1["type"].(string),
			PublicIp:     if1["wan_public_ip"].(string),
			Dhcp:         if1["enable_dhcp"].(bool),
			IpAddr:       if1["ip_address"].(string),
			GatewayIp:    if1["gateway_ip"].(string),
			DNSPrimary:   if1["dns_server_ip"].(string),
			DNSSecondary: if1["secondary_dns_server_ip"].(string),
		}

		edgeGatewaySelfmanagedHa.InterfaceList = append(edgeGatewaySelfmanagedHa.InterfaceList, if2)
	}

	return edgeGatewaySelfmanagedHa
}

func resourceAviatrixEdgeGatewaySelfmanagedHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeGatewaySelfmanagedHa := marshalEdgeGatewaySelfmanagedHaInput(d)

	customInterfaceMapping, ok := d.Get("custom_interface_mapping").([]interface{})
	if ok {
		customInterfaceMap, err := getCustomInterfaceMapDetails(customInterfaceMapping)
		if err != nil {
			return diag.Errorf("failed to get custom interface mapping details: %s", err)
		}
		edgeGatewaySelfmanagedHa.CustomInterfaceMapping = customInterfaceMap
	}

	edgeGatewaySelfmanagedHaName, err := client.CreateEdgeVmSelfmanagedHa(ctx, edgeGatewaySelfmanagedHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Gateway Selfmanaged HA: %s", err)
	}

	d.SetId(edgeGatewaySelfmanagedHaName)
	return resourceAviatrixEdgeGatewaySelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeGatewaySelfmanagedHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("primary_gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeGatewaySelfmanagedHaResp, err := client.GetEdgeVmSelfmanagedHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if err == goaviatrix.ErrNotFound {
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
		d.Set("ztp_file_type", edgeGatewaySelfmanagedHaResp.ZtpFileType)
	}

	if edgeGatewaySelfmanagedHaResp.ZtpFileType == "cloud_init" {
		d.Set("ztp_file_type", "cloud-init")
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

		interfaces = append(interfaces, if1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	if len(edgeGatewaySelfmanagedHaResp.CustomInterfaceMapping) != 0 {
		// get the order of custom interface mapping
		userCustomInterfaceMapping, ok := d.Get("custom_interface_mapping").([]interface{})
		if !ok {
			return diag.Errorf("failed to get custom_interface_mapping")
		}
		userCustomInterfaceOrder, err := getCustomInterfaceOrder(userCustomInterfaceMapping)
		if err != nil {
			return diag.Errorf("failed to get custom_interface_order: %s\n", err)
		}
		customInterfaceMapping, err := setCustomInterfaceMapping(edgeGatewaySelfmanagedHaResp.CustomInterfaceMapping, userCustomInterfaceOrder)
		if !ok {
			return diag.Errorf("failed to get custom_interface_mapping: %s\n", err)
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("custom_interface_mapping", customInterfaceMapping); err != nil {
			return diag.Errorf("failed to set custom_interface_mapping: %s\n", err)
		}
	}

	d.SetId(edgeGatewaySelfmanagedHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeGatewaySelfmanagedHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

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

	if d.HasChange("custom_interface_mapping") {
		return diag.Errorf("updating custom interface mapping after the selfmanaged Edge Gateway creation is not supported")
	}

	d.Partial(false)
	return resourceAviatrixEdgeGatewaySelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeGatewaySelfmanagedHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

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
