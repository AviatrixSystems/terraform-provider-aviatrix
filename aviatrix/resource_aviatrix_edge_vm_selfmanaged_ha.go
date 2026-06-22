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

func resourceAviatrixEdgeVmSelfmanagedHa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeVmSelfmanagedHaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeVmSelfmanagedHaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeVmSelfmanagedHaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeVmSelfmanagedHaDelete,
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
		DeprecationMessage: "Since V3.1.2+, please use resource aviatrix_edge_gateway_selfmanaged_ha instead. Resource " +
			"aviatrix_edge_vm_selfmanaged_ha will be deprecated in the V3.2.0 release.",
	}
}

func marshalEdgeVmSelfmanagedHaInput(d *schema.ResourceData) *goaviatrix.EdgeVmSelfmanagedHa {
	edgeVmSelfmanagedHa := &goaviatrix.EdgeVmSelfmanagedHa{
		PrimaryGwName:            getString(d, "primary_gw_name"),
		SiteId:                   getString(d, "site_id"),
		ZtpFileType:              getString(d, "ztp_file_type"),
		ZtpFileDownloadPath:      getString(d, "ztp_file_download_path"),
		ManagementEgressIPPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := getSet(d, "interfaces").List()
	for _, if0 := range interfaces {
		if1 := mustMap(if0)

		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:    mustString(if1["name"]),
			Type:      mustString(if1["type"]),
			PublicIp:  mustString(if1["wan_public_ip"]),
			Dhcp:      mustBool(if1["enable_dhcp"]),
			IpAddr:    mustString(if1["ip_address"]),
			GatewayIp: mustString(if1["gateway_ip"]),
		}

		edgeVmSelfmanagedHa.InterfaceList = append(edgeVmSelfmanagedHa.InterfaceList, if2)
	}

	return edgeVmSelfmanagedHa
}

func resourceAviatrixEdgeVmSelfmanagedHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeVmSelfmanagedHa := marshalEdgeVmSelfmanagedHaInput(d)

	edgeVmSelfmanagedHaName, err := client.CreateEdgeVmSelfmanagedHa(ctx, edgeVmSelfmanagedHa)
	if err != nil {
		return diag.Errorf("failed to create Edge VM Selfmanaged HA: %s", err)
	}

	d.SetId(edgeVmSelfmanagedHaName)
	return resourceAviatrixEdgeVmSelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeVmSelfmanagedHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if getString(d, "primary_gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		mustSet(d, "primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeVmSelfmanagedHaResp, err := client.GetEdgeVmSelfmanagedHa(ctx, getString(d, "primary_gw_name")+"-hagw")
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge VM Selfmanaged HA: %v", err)
	}

	_ = d.Set("primary_gw_name", edgeVmSelfmanagedHaResp.PrimaryGwName)
	_ = d.Set("site_id", edgeVmSelfmanagedHaResp.SiteID)

	if edgeVmSelfmanagedHaResp.ZtpFileType == "iso" || edgeVmSelfmanagedHaResp.ZtpFileType == "cloud-init" {
		mustSet(d, "ztp_file_type", edgeVmSelfmanagedHaResp.ZtpFileType)
	}

	if edgeVmSelfmanagedHaResp.ZtpFileType == "cloud_init" {
		mustSet(d, "ztp_file_type", "cloud-init")
	}

	if edgeVmSelfmanagedHaResp.ManagementEgressIPPrefix == "" {
		_ = d.Set("management_egress_ip_prefix_list", nil)
	} else {
		_ = d.Set("management_egress_ip_prefix_list", strings.Split(edgeVmSelfmanagedHaResp.ManagementEgressIPPrefix, ","))
	}

	var interfaces []map[string]interface{}
	for _, if0 := range edgeVmSelfmanagedHaResp.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
		if1["wan_public_ip"] = if0.PublicIp
		if1["enable_dhcp"] = if0.Dhcp
		if1["ip_address"] = if0.IpAddr
		if1["gateway_ip"] = if0.GatewayIp

		interfaces = append(interfaces, if1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	d.SetId(edgeVmSelfmanagedHaResp.GwName)
	return nil
}

func resourceAviatrixEdgeVmSelfmanagedHaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	edgeVmSelfmanagedHa := marshalEdgeVmSelfmanagedHaInput(d)

	d.Partial(true)

	gatewayForEdgeVmSelfmanagedFunctions := &goaviatrix.EdgeSpoke{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeVmSelfmanagedFunctions.InterfaceList = edgeVmSelfmanagedHa.InterfaceList
		gatewayForEdgeVmSelfmanagedFunctions.ManagementEgressIpPrefix = edgeVmSelfmanagedHa.ManagementEgressIPPrefix

		err := client.UpdateEdgeVmSelfmanagedHa(ctx, gatewayForEdgeVmSelfmanagedFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge VM Selfmanaged HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeVmSelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeVmSelfmanagedHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteEdgeSpoke(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not delete Edge VM Selfmanaged HA %s: %v", d.Id(), err)
	}

	edgeVmSelfmanagedHa := marshalEdgeVmSelfmanagedHaInput(d)

	var fileName string
	if edgeVmSelfmanagedHa.ZtpFileType == "iso" {
		fileName = edgeVmSelfmanagedHa.ZtpFileDownloadPath + "/" + edgeVmSelfmanagedHa.PrimaryGwName + "-" + edgeVmSelfmanagedHa.SiteId + "-ha.iso"
	} else {
		fileName = edgeVmSelfmanagedHa.ZtpFileDownloadPath + "/" + edgeVmSelfmanagedHa.PrimaryGwName + "-" + edgeVmSelfmanagedHa.SiteId + "-ha-cloud-init.txt"
	}

	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
