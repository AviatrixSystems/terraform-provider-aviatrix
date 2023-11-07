package aviatrix

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeGatewaySelfmanaged() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeGatewaySelfmanagedCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeGatewaySelfmanagedRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeGatewaySelfmanagedUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeGatewaySelfmanagedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge gateway selfmanaged name.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"management_egress_ip_prefix_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of management egress gateway IP/prefix.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_management_over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable management over private network.",
			},
			"dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "DNS server IP.",
				ValidateFunc: validation.IsIPAddress,
			},
			"secondary_dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Secondary DNS server IP.",
				ValidateFunc: validation.IsIPAddress,
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
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Local AS number.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of AS numbers to prepend gateway BGP AS_Path field.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"enable_edge_active_standby": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Edge Active-Standby Mode.",
			},
			"enable_edge_active_standby_preemptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Preemptive Mode for Edge Active-Standby, available only with Active-Standby enabled.",
			},
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable learned CIDR approval for BGP Spoke Gateway. Valid values: true, false.",
			},
			"approved_learned_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional:    true,
				Description: "Approved learned CIDRs for BGP Spoke Gateway.",
			},
			"spoke_bgp_manual_advertise_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Intended CIDR list to be advertised to external BGP router.",
			},
			"enable_preserve_as_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable preserve as path when advertising manual summary CIDRs on BGP spoke gateway.",
			},
			"bgp_polling_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpPollingTime,
				ValidateFunc: validation.IntBetween(10, 50),
				Description:  "BGP route polling time for BGP Spoke Gateway. Unit is in seconds. Valid values are between 10 and 50.",
			},
			"bgp_hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpHoldTime,
				ValidateFunc: validation.IntBetween(12, 360),
				Description:  "BGP Hold Time for BGP Spoke Gateway. Unit is in seconds. Valid values are between 12 and 360.",
			},
			"enable_edge_transitive_routing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Edge transitive routing.",
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable jumbo frame.",
			},
			"latitude": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     goaviatrix.ValidateEdgeSpokeLatitude,
				Description:      "The latitude of the Edge as a Spoke.",
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncEdgeSpokeCoordinate,
			},
			"longitude": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     goaviatrix.ValidateEdgeSpokeLongitude,
				Description:      "The longitude of the Edge as a Spoke.",
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncEdgeSpokeCoordinate,
			},
			"rx_queue_size": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1K", "2K", "4K"}, false),
				Description:  "Ethernet interface RX queue size.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of Edge gateway selfmanaged.",
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
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Interface type.",
							ValidateFunc: validation.StringInSlice([]string{"WAN", "LAN", "MANAGEMENT"}, false),
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
						"enable_vrrp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable VRRP.",
						},
						"vrrp_virtual_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VRRP virtual IP.",
						},
						"tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag.",
						},
					},
				},
			},
			"vlan": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "VLAN configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_interface_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parent interface name.",
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "VLAN ID.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "LAN sub-interface IP address.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "LAN sub-interface gateway IP.",
						},
						"peer_ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "LAN sub-interface IP address on HA gateway.",
						},
						"peer_gateway_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "LAN sub-interface gateway IP on HA gateway.",
						},
						"vrrp_virtual_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "LAN sub-interface virtual IP.",
						},
						"tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag.",
						},
					},
				},
			},
		},
	}
}

func marshalEdgeGatewaySelfmanagedInput(d *schema.ResourceData) *goaviatrix.EdgeSpoke {
	edgeSpoke := &goaviatrix.EdgeSpoke{
		GwName:                             d.Get("gw_name").(string),
		SiteId:                             d.Get("site_id").(string),
		ManagementEgressIpPrefix:           strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
		EnableManagementOverPrivateNetwork: d.Get("enable_management_over_private_network").(bool),
		DnsServerIp:                        d.Get("dns_server_ip").(string),
		SecondaryDnsServerIp:               d.Get("secondary_dns_server_ip").(string),
		ZtpFileType:                        d.Get("ztp_file_type").(string),
		ZtpFileDownloadPath:                d.Get("ztp_file_download_path").(string),
		EnableEdgeActiveStandby:            d.Get("enable_edge_active_standby").(bool),
		EnableEdgeActiveStandbyPreemptive:  d.Get("enable_edge_active_standby_preemptive").(bool),
		LocalAsNumber:                      d.Get("local_as_number").(string),
		PrependAsPath:                      getStringList(d, "prepend_as_path"),
		EnableLearnedCidrsApproval:         d.Get("enable_learned_cidrs_approval").(bool),
		ApprovedLearnedCidrs:               getStringSet(d, "approved_learned_cidrs"),
		SpokeBgpManualAdvertisedCidrs:      getStringSet(d, "spoke_bgp_manual_advertise_cidrs"),
		EnablePreserveAsPath:               d.Get("enable_preserve_as_path").(bool),
		BgpPollingTime:                     d.Get("bgp_polling_time").(int),
		BgpHoldTime:                        d.Get("bgp_hold_time").(int),
		EnableEdgeTransitiveRouting:        d.Get("enable_edge_transitive_routing").(bool),
		EnableJumboFrame:                   d.Get("enable_jumbo_frame").(bool),
		Latitude:                           d.Get("latitude").(string),
		Longitude:                          d.Get("longitude").(string),
		RxQueueSize:                        d.Get("rx_queue_size").(string),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, if0 := range interfaces {
		if1 := if0.(map[string]interface{})

		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:    if1["name"].(string),
			Type:      if1["type"].(string),
			Dhcp:      if1["enable_dhcp"].(bool),
			PublicIp:  if1["wan_public_ip"].(string),
			IpAddr:    if1["ip_address"].(string),
			GatewayIp: if1["gateway_ip"].(string),
			VrrpState: if1["enable_vrrp"].(bool),
			VirtualIp: if1["vrrp_virtual_ip"].(string),
			Tag:       if1["tag"].(string),
		}

		edgeSpoke.InterfaceList = append(edgeSpoke.InterfaceList, if2)
	}

	vlan := d.Get("vlan").(*schema.Set).List()
	for _, vlan0 := range vlan {
		vlan1 := vlan0.(map[string]interface{})

		vlan2 := &goaviatrix.EdgeSpokeVlan{
			ParentInterface: vlan1["parent_interface_name"].(string),
			IpAddr:          vlan1["ip_address"].(string),
			GatewayIp:       vlan1["gateway_ip"].(string),
			PeerIpAddr:      vlan1["peer_ip_address"].(string),
			PeerGatewayIp:   vlan1["peer_gateway_ip"].(string),
			VirtualIp:       vlan1["vrrp_virtual_ip"].(string),
			Tag:             vlan1["tag"].(string),
		}

		vlan2.VlanId = strconv.Itoa(vlan1["vlan_id"].(int))

		edgeSpoke.VlanList = append(edgeSpoke.VlanList, vlan2)
	}

	return edgeSpoke
}

func resourceAviatrixEdgeGatewaySelfmanagedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeSpoke := marshalEdgeGatewaySelfmanagedInput(d)

	// checks before creation
	if !edgeSpoke.EnableEdgeActiveStandby && edgeSpoke.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeSpoke.EnableLearnedCidrsApproval && len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		if edgeSpoke.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeSpoke.Latitude != "" && edgeSpoke.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeSpoke.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeSpoke.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	// create
	d.SetId(edgeSpoke.GwName)
	flag := false
	defer resourceAviatrixEdgeGatewaySelfmanagedReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeSpoke(ctx, edgeSpoke); err != nil {
		return diag.Errorf("could not create Edge as a Spoke: %v", err)
	}

	// advanced configs
	// use following variables to reuse functions for transit, spoke and gateway
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeSpoke.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeSpoke.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeSpoke.GwName,
	}

	if edgeSpoke.LocalAsNumber != "" {
		err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeSpoke.LocalAsNumber)
		if err != nil {
			return diag.Errorf("could not set 'local_as_number' after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeSpoke.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeSpoke.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if len(edgeSpoke.SpokeBgpManualAdvertisedCidrs) != 0 {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeSpoke.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.EnablePreserveAsPath {
		err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
		if err != nil {
			return diag.Errorf("could not enable spoke preserve as path after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.BgpPollingTime >= 10 && edgeSpoke.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeSpoke.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.BgpHoldTime >= 12 && edgeSpoke.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeSpoke.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.Latitude != "" || edgeSpoke.Longitude != "" {
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeSpoke.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableEdgeActiveStandby || edgeSpoke.EnableEdgeActiveStandbyPreemptive {
		err := client.UpdateEdgeSpoke(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update Edge active standby or Edge active standby preemptive after Edge Gateway Selfmanaged creation: %v", err)
		}
	}

	return resourceAviatrixEdgeGatewaySelfmanagedReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeGatewaySelfmanagedReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeGatewaySelfmanagedRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeGatewaySelfmanagedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// handle import
	if d.Get("gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	edgeSpoke, err := client.GetEdgeSpoke(ctx, d.Get("gw_name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Gateway Selfmanaged: %v", err)
	}

	d.Set("gw_name", edgeSpoke.GwName)
	d.Set("site_id", edgeSpoke.SiteId)
	d.Set("enable_management_over_private_network", edgeSpoke.EnableManagementOverPrivateNetwork)
	d.Set("dns_server_ip", edgeSpoke.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeSpoke.SecondaryDnsServerIp)
	d.Set("local_as_number", edgeSpoke.LocalAsNumber)
	d.Set("prepend_as_path", edgeSpoke.PrependAsPath)
	d.Set("enable_edge_active_standby", edgeSpoke.EnableEdgeActiveStandby)
	d.Set("enable_edge_active_standby_preemptive", edgeSpoke.EnableEdgeActiveStandbyPreemptive)
	d.Set("enable_learned_cidrs_approval", edgeSpoke.EnableLearnedCidrsApproval)

	if edgeSpoke.ZtpFileType == "iso" || edgeSpoke.ZtpFileType == "cloud-init" {
		d.Set("ztp_file_type", edgeSpoke.ZtpFileType)
	}

	if edgeSpoke.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeSpoke.ManagementEgressIpPrefix, ","))
	}

	if edgeSpoke.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: edgeSpoke.GwName})
		if err != nil {
			return diag.Errorf("could not get advanced config for Edge as a Spoke: %v", err)
		}

		err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not set approved_learned_cidrs into state: %v", err)
		}
	} else {
		d.Set("approved_learned_cidrs", nil)
	}

	spokeBgpManualAdvertisedCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	if len(goaviatrix.Difference(spokeBgpManualAdvertisedCidrs, edgeSpoke.SpokeBgpManualAdvertisedCidrs)) != 0 ||
		len(goaviatrix.Difference(edgeSpoke.SpokeBgpManualAdvertisedCidrs, spokeBgpManualAdvertisedCidrs)) != 0 {
		d.Set("spoke_bgp_manual_advertise_cidrs", edgeSpoke.SpokeBgpManualAdvertisedCidrs)
	} else {
		d.Set("spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertisedCidrs)
	}

	d.Set("enable_preserve_as_path", edgeSpoke.EnablePreserveAsPath)
	d.Set("bgp_polling_time", edgeSpoke.BgpPollingTime)
	d.Set("bgp_hold_time", edgeSpoke.BgpHoldTime)
	d.Set("enable_edge_transitive_routing", edgeSpoke.EnableEdgeTransitiveRouting)
	d.Set("enable_jumbo_frame", edgeSpoke.EnableJumboFrame)
	if edgeSpoke.Latitude != 0 || edgeSpoke.Longitude != 0 {
		d.Set("latitude", fmt.Sprintf("%.6f", edgeSpoke.Latitude))
		d.Set("longitude", fmt.Sprintf("%.6f", edgeSpoke.Longitude))
	} else {
		d.Set("latitude", "")
		d.Set("longitude", "")
	}

	d.Set("rx_queue_size", edgeSpoke.RxQueueSize)
	d.Set("state", edgeSpoke.State)

	var interfaces []map[string]interface{}
	var vlan []map[string]interface{}
	for _, if0 := range edgeSpoke.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
		if1["enable_dhcp"] = if0.Dhcp
		if1["wan_public_ip"] = if0.PublicIp
		if1["ip_address"] = if0.IpAddr
		if1["gateway_ip"] = if0.GatewayIp
		if1["vrrp_virtual_ip"] = if0.VirtualIp
		if1["tag"] = if0.Tag

		if if0.Type == "LAN" {
			if1["enable_vrrp"] = if0.VrrpState
		}

		if if0.Type == "LAN" && if0.SubInterfaces != nil {
			for _, vlan0 := range if0.SubInterfaces {
				vlan1 := make(map[string]interface{})
				vlan1["parent_interface_name"] = vlan0.ParentInterface
				vlan1["ip_address"] = vlan0.IpAddr
				vlan1["gateway_ip"] = vlan0.GatewayIp
				vlan1["peer_ip_address"] = vlan0.PeerIpAddr
				vlan1["peer_gateway_ip"] = vlan0.PeerGatewayIp
				vlan1["vrrp_virtual_ip"] = vlan0.VirtualIp
				vlan1["tag"] = vlan0.Tag

				vlanId, _ := strconv.Atoi(vlan0.VlanId)
				vlan1["vlan_id"] = vlanId

				vlan = append(vlan, vlan1)
			}
		}

		interfaces = append(interfaces, if1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	if err = d.Set("vlan", vlan); err != nil {
		return diag.Errorf("failed to set vlan: %s\n", err)
	}

	d.SetId(edgeSpoke.GwName)
	return nil
}

func resourceAviatrixEdgeGatewaySelfmanagedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeSpoke := marshalEdgeGatewaySelfmanagedInput(d)

	// checks before update
	if !edgeSpoke.EnableEdgeActiveStandby && edgeSpoke.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeSpoke.EnableLearnedCidrsApproval && len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		if edgeSpoke.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeSpoke.Latitude != "" && edgeSpoke.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeSpoke.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeSpoke.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	d.Partial(true)

	// update configs
	// use following variables to reuse functions for transit, spoke and gateway
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeSpoke.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeSpoke.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeSpoke.GwName,
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeSpoke.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge Gateway Selfmanaged update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeSpoke.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge Gateway Selfmanaged update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeSpoke.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeSpoke.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge Gateway Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeSpoke.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge Gateway Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge Gateway Selfmanaged update: %v", err)
			}
		}
	}

	if edgeSpoke.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeSpoke.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge Gateway Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeSpoke.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge Gateway Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeSpoke.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge Gateway Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge Gateway Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeSpoke.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge Gateway Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeSpoke.GwName, edgeSpoke.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge Gateway Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeSpoke.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge Gateway Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge Gateway Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeSpoke.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge Gateway Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge Gateway Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge Gateway Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeSpoke.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge Gateway Selfmanaged update: %v", err)
		}
	}

	if d.HasChanges("management_egress_ip_prefix_list", "interfaces", "vlan", "enable_edge_active_standby",
		"enable_edge_active_standby_preemptive") {
		err := client.UpdateEdgeSpoke(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list, WAN/LAN/MANAGEMENT/VLAN interfaces, "+
				"Edge active standby or Edge active standby preemptive during Edge as a Spoke update: %v", err)
		}
	}

	d.Partial(false)

	return resourceAviatrixEdgeGatewaySelfmanagedRead(ctx, d, meta)
}

func resourceAviatrixEdgeGatewaySelfmanagedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	siteId := d.Get("site_id").(string)

	ztpFileDownloadPath := d.Get("ztp_file_download_path").(string)
	ztpFileType := d.Get("ztp_file_type").(string)

	err := client.DeleteEdgeSpoke(ctx, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge Gateway Selfmanaged: %v", err)
	}

	var fileName string
	if ztpFileType == "iso" {
		fileName = ztpFileDownloadPath + "/" + gwName + "-" + siteId + ".iso"
	} else {
		fileName = ztpFileDownloadPath + "/" + gwName + "-" + siteId + "-cloud-init.txt"
	}

	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
