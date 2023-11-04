package aviatrix

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		PrimaryGwName:            d.Get("primary_gw_name").(string),
		SiteId:                   d.Get("site_id").(string),
		ZtpFileType:              d.Get("ztp_file_type").(string),
		ZtpFileDownloadPath:      d.Get("ztp_file_download_path").(string),
		ManagementEgressIpPrefix: strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, if0 := range interfaces {
		if1 := if0.(map[string]interface{})

		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:    if1["name"].(string),
			Type:      if1["type"].(string),
			PublicIp:  if1["wan_public_ip"].(string),
			Dhcp:      if1["enable_dhcp"].(bool),
			IpAddr:    if1["ip_address"].(string),
			GatewayIp: if1["gateway_ip"].(string),
		}

		edgeVmSelfmanagedHa.InterfaceList = append(edgeVmSelfmanagedHa.InterfaceList, if2)
	}

	return edgeVmSelfmanagedHa
}

func resourceAviatrixEdgeVmSelfmanagedHaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeVmSelfmanagedHa := marshalEdgeVmSelfmanagedHaInput(d)

	edgeVmSelfmanagedHaName, err := client.CreateEdgeVmSelfmanagedHa(ctx, edgeVmSelfmanagedHa)
	if err != nil {
		return diag.Errorf("failed to create Edge VM Selfmanaged HA: %s", err)
	}

	d.SetId(edgeVmSelfmanagedHaName)
	return resourceAviatrixEdgeVmSelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeVmSelfmanagedHaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("primary_gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "-hagw")
		d.Set("primary_gw_name", parts[0])
		d.SetId(id)
	}

	edgeVmSelfmanagedHaResp, err := client.GetEdgeVmSelfmanagedHa(ctx, d.Get("primary_gw_name").(string)+"-hagw")
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge VM Selfmanaged HA: %v", err)
	}

	d.Set("primary_gw_name", edgeVmSelfmanagedHaResp.PrimaryGwName)
	d.Set("site_id", edgeVmSelfmanagedHaResp.SiteId)

	if edgeVmSelfmanagedHaResp.ZtpFileType == "iso" || edgeVmSelfmanagedHaResp.ZtpFileType == "cloud-init" {
		d.Set("ztp_file_type", edgeVmSelfmanagedHaResp.ZtpFileType)
	}

	if edgeVmSelfmanagedHaResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeVmSelfmanagedHaResp.ManagementEgressIpPrefix, ","))
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
	client := meta.(*goaviatrix.Client)

	edgeVmSelfmanagedHa := marshalEdgeVmSelfmanagedHaInput(d)

	d.Partial(true)

	gatewayForEdgeVmSelfmanagedFunctions := &goaviatrix.EdgeSpoke{
		GwName: d.Id(),
	}

	if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
		gatewayForEdgeVmSelfmanagedFunctions.InterfaceList = edgeVmSelfmanagedHa.InterfaceList
		gatewayForEdgeVmSelfmanagedFunctions.ManagementEgressIpPrefix = edgeVmSelfmanagedHa.ManagementEgressIpPrefix

		err := client.UpdateEdgeVmSelfmanagedHa(ctx, gatewayForEdgeVmSelfmanagedFunctions)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list or WAN/LAN/VLAN interfaces during Edge VM Selfmanaged HA update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeVmSelfmanagedHaRead(ctx, d, meta)
}

func resourceAviatrixEdgeVmSelfmanagedHaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

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
