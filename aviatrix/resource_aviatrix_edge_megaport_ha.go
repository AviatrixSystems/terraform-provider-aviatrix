package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
				DiffSuppressFunc: func(_, old, _ string, _ *schema.ResourceData) bool {
					return old != ""
				},
				Description: "The location where the ZTP file will be stored.",
			},
			"interfaces": interfaceSchema(),
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

func interfaceSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Description: "WAN/LAN/MANAGEMENT interfaces.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"logical_ifname":          logicalIfnameSchema(),
				"enable_dhcp":             enableDhcpSchema(),
				"wan_public_ip":           optionalStringSchema("WAN interface public IP."),
				"ip_address":              optionalStringSchema("Interface static IP address."),
				"gateway_ip":              optionalStringSchema("Gateway IP."),
				"dns_server_ip":           optionalStringSchema("Primary DNS server IP."),
				"secondary_dns_server_ip": optionalStringSchema("Secondary DNS server IP."),
				"tag":                     optionalStringSchema("Tag."),
			},
		},
	}
}

func logicalIfnameSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Logical interface name e.g., wan0, lan0, mgmt0.",
		ValidateFunc: validation.StringMatch(
			regexp.MustCompile(`^(wan|lan|mgmt)[0-9]+$`),
			"Logical interface name must start with 'wan', 'lan', or 'mgmt' followed by a number (e.g., 'wan0', 'lan1', 'mgmt2').",
		),
	}
}

func enableDhcpSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Enable DHCP.",
	}
}

func optionalStringSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: description,
	}
}

func marshalEdgeMegaportHaInput(d *schema.ResourceData) (*goaviatrix.EdgeMegaportHa, error) {
	primaryGwName, ztpFileDownloadPath, managementEgressIPPrefixList, err := parseRequiredFields(d)
	if err != nil {
		return nil, err
	}

	edgeMegaportHa := &goaviatrix.EdgeMegaportHa{
		PrimaryGwName:            primaryGwName,
		ZtpFileDownloadPath:      ztpFileDownloadPath,
		ManagementEgressIPPrefix: strings.Join(managementEgressIPPrefixList, ","),
	}

	interfaces, err := parseInterfaces(d)
	if err != nil {
		return nil, err
	}
	edgeMegaportHa.InterfaceList = interfaces

	return edgeMegaportHa, nil
}

func parseRequiredFields(d *schema.ResourceData) (string, string, []string, error) {
	primaryGwName, ok := d.Get("primary_gw_name").(string)
	if !ok || primaryGwName == "" {
		return "", "", nil, fmt.Errorf("invalid or missing value for 'primary_gw_name'")
	}

	ztpFileDownloadPath, ok := d.Get("ztp_file_download_path").(string)
	if !ok || ztpFileDownloadPath == "" {
		return "", "", nil, fmt.Errorf("invalid or missing value for 'ztp_file_download_path'")
	}

	managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
	if len(managementEgressIPPrefixList) == 0 {
		return "", "", nil, fmt.Errorf("invalid or empty value for 'management_egress_ip_prefix_list'")
	}

	return primaryGwName, ztpFileDownloadPath, managementEgressIPPrefixList, nil
}

func parseInterfaces(d *schema.ResourceData) ([]*goaviatrix.EdgeMegaportInterface, error) {
	rawInterfaces, ok := d.Get("interfaces").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("failed to parse interfaces")
	}

	var interfaces []*goaviatrix.EdgeMegaportInterface
	for _, interface0 := range rawInterfaces.List() {
		interface1, ok := interface0.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse interface")
		}

		interface2 := &goaviatrix.EdgeMegaportInterface{}
		assignInterfaceFields(interface1, interface2)
		interfaces = append(interfaces, interface2)
	}

	return interfaces, nil
}

func assignInterfaceFields(interface1 map[string]interface{}, interface2 *goaviatrix.EdgeMegaportInterface) {
	if logicalIfname, ok := interface1["logical_ifname"].(string); ok {
		interface2.LogicalInterfaceName = logicalIfname
	}
	if publicIP, ok := interface1["wan_public_ip"].(string); ok {
		interface2.PublicIP = publicIP
	}
	if tag, ok := interface1["tag"].(string); ok {
		interface2.Tag = tag
	}
	if dhcp, ok := interface1["enable_dhcp"].(bool); ok {
		interface2.Dhcp = dhcp
	}
	if ipAddr, ok := interface1["ip_address"].(string); ok {
		interface2.IPAddr = ipAddr
	}
	if gatewayIP, ok := interface1["gateway_ip"].(string); ok {
		interface2.GatewayIP = gatewayIP
	}
	if dnsPrimary, ok := interface1["dns_server_ip"].(string); ok {
		interface2.DNSPrimary = dnsPrimary
	}
	if dnsSecondary, ok := interface1["secondary_dns_server_ip"].(string); ok {
		interface2.DNSSecondary = dnsSecondary
	}
}

func resourceAviatrixEdgeMegaportHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok || client == nil {
		return diag.Errorf("failed to cast meta to *goaviatrix.Client")
	}

	edgeMegaportHa, err := marshalEdgeMegaportHaInput(d)
	if err != nil {
		return diag.Errorf("failed to marshal Edge Megaport HA input: %s", err)
	}

	edgeMegaportHaName, err := client.CreateEdgeMegaportHa(ctx, edgeMegaportHa)
	if err != nil {
		return diag.Errorf("failed to create Edge Megaport HA: %s", err)
	}

	d.SetId(edgeMegaportHaName)
	return resourceAviatrixEdgeMegaportHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeMegaportHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok || client == nil {
		return diag.Errorf("failed to cast meta to *goaviatrix.Client")
	}

	if primaryGwName, ok := d.Get("primary_gw_name").(string); ok && primaryGwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		_ = d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeMegaportHaResp, err := client.GetEdgeMegaportHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Megaport HA: %v", err)
	}

	_ = d.Set("primary_gw_name", edgeMegaportHaResp.PrimaryGwName)
	_ = d.Set("account_name", edgeMegaportHaResp.AccountName)

	if edgeMegaportHaResp.ManagementEgressIPPrefix == "" {
		_ = d.Set("management_egress_ip_prefix_list", nil)
	} else {
		_ = d.Set("management_egress_ip_prefix_list", strings.Split(edgeMegaportHaResp.ManagementEgressIPPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, interface0 := range edgeMegaportHaResp.InterfaceList {
		interface1 := make(map[string]interface{})
		interface1["logical_ifname"] = interface0.LogicalInterfaceName
		interface1["wan_public_ip"] = interface0.PublicIP
		interface1["tag"] = interface0.Tag
		interface1["enable_dhcp"] = interface0.Dhcp
		interface1["ip_address"] = interface0.IPAddr
		interface1["gateway_ip"] = interface0.GatewayIP
		interface1["dns_server_ip"] = interface0.DNSPrimary
		interface1["secondary_dns_server_ip"] = interface0.DNSSecondary

		interfaces = append(interfaces, interface1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	d.SetId(edgeMegaportHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeMegaportHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok || client == nil {
		return diag.Errorf("failed to cast meta to *goaviatrix.Client")
	}

	edgeMegaportHa, err := marshalEdgeMegaportHaInput(d)
	if err != nil {
		return diag.Errorf("failed to marshal Edge Megaport HA input: %s", err)
	}

	d.Partial(true)

	gatewayForEdgeMegaportFunctions := &goaviatrix.EdgeMegaport{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeMegaportFunctions.InterfaceList = edgeMegaportHa.InterfaceList
		gatewayForEdgeMegaportFunctions.ManagementEgressIPPrefix = edgeMegaportHa.ManagementEgressIPPrefix

		err := client.UpdateEdgeMegaportHa(ctx, gatewayForEdgeMegaportFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge Megaport HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeMegaportHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeMegaportHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok || client == nil {
		return diag.Errorf("failed to cast meta to *goaviatrix.Client")
	}

	edgeMegaportHa, err := marshalEdgeMegaportHaInput(d)
	if err != nil {
		return diag.Errorf("failed to marshal Edge Megaport HA input: %s", err)
	}
	accountName, ok := d.Get("account_name").(string)
	if !ok {
		return diag.Errorf("failed to get account name")
	}

	err = client.DeleteEdgeMegaport(ctx, accountName, d.Id())
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
