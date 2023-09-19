package aviatrix

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeEquinix() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeEquinixCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeEquinixRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeEquinixUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeEquinixDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Equinix account name.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Equinix name.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
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
				Description:  "BGP route polling time for BGP spoke gateway in seconds. Valid values are between 12 and 360.",
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
				Description: "State of Edge as a Spoke.",
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
			"dns_profile_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS profile to be associated with gateway, select an existing template.",
			},
			"enable_single_ip_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Single IP SNAT.",
			},
			"enable_auto_advertise_lan_cidrs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable auto advertise LAN CIDRs.",
			},
		},
	}
}

func marshalEdgeEquinixInput(d *schema.ResourceData) *goaviatrix.EdgeEquinix {
	edgeEquinix := &goaviatrix.EdgeEquinix{
		AccountName:                        d.Get("account_name").(string),
		GwName:                             d.Get("gw_name").(string),
		SiteId:                             d.Get("site_id").(string),
		ZtpFileDownloadPath:                d.Get("ztp_file_download_path").(string),
		ManagementEgressIpPrefix:           strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
		EnableManagementOverPrivateNetwork: d.Get("enable_management_over_private_network").(bool),
		DnsServerIp:                        d.Get("dns_server_ip").(string),
		SecondaryDnsServerIp:               d.Get("secondary_dns_server_ip").(string),
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
		DnsProfileName:                     d.Get("dns_profile_name").(string),
		EnableSingleIpSnat:                 d.Get("enable_single_ip_snat").(bool),
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
			VrrpState:    interface1["enable_vrrp"].(bool),
			VirtualIp:    interface1["vrrp_virtual_ip"].(string),
		}

		edgeEquinix.InterfaceList = append(edgeEquinix.InterfaceList, interface2)
	}

	vlan := d.Get("vlan").(*schema.Set).List()
	for _, vlan0 := range vlan {
		vlan1 := vlan0.(map[string]interface{})

		vlan2 := &goaviatrix.EdgeEquinixVlan{
			ParentInterface: vlan1["parent_interface_name"].(string),
			IpAddr:          vlan1["ip_address"].(string),
			GatewayIp:       vlan1["gateway_ip"].(string),
			PeerIpAddr:      vlan1["peer_ip_address"].(string),
			PeerGatewayIp:   vlan1["peer_gateway_ip"].(string),
			VirtualIp:       vlan1["vrrp_virtual_ip"].(string),
			Tag:             vlan1["tag"].(string),
		}

		vlan2.VlanId = strconv.Itoa(vlan1["vlan_id"].(int))

		edgeEquinix.VlanList = append(edgeEquinix.VlanList, vlan2)
	}

	if d.Get("enable_auto_advertise_lan_cidrs").(bool) {
		edgeEquinix.EnableAutoAdvertiseLanCidrs = "enable"
	} else {
		edgeEquinix.EnableAutoAdvertiseLanCidrs = "disable"
	}

	return edgeEquinix
}

func resourceAviatrixEdgeEquinixCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeEquinix := marshalEdgeEquinixInput(d)

	// checks before creation
	if !edgeEquinix.EnableEdgeActiveStandby && edgeEquinix.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not enable Preemptive Mode when Active-Standby is disabled")
	}

	if !edgeEquinix.EnableLearnedCidrsApproval && len(edgeEquinix.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeEquinix.PrependAsPath) != 0 {
		if edgeEquinix.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeEquinix.Latitude != "" && edgeEquinix.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeEquinix.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeEquinix.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	// create
	d.SetId(edgeEquinix.GwName)
	flag := false
	defer resourceAviatrixEdgeEquinixReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeEquinix(ctx, edgeEquinix); err != nil {
		return diag.Errorf("could not create Edge Equinix %s: %v", edgeEquinix.GwName, err)
	}

	// advanced configs
	// use following variables to reuse functions for transit, spoke, gateway and EaaS
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeEquinix.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeEquinix.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeEquinix.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeEquinix.GwName,
	}

	if edgeEquinix.LocalAsNumber != "" {
		err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeEquinix.LocalAsNumber)
		if err != nil {
			return diag.Errorf("could not set 'local_as_number' after Edge Equinix creation: %v", err)
		}
	}

	if len(edgeEquinix.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeEquinix.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge Equinix creation: %v", err)
		}
	}

	if len(edgeEquinix.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeEquinix.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge Equinix creation: %v", err)
		}
	}

	if len(edgeEquinix.SpokeBgpManualAdvertisedCidrs) != 0 {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeEquinix.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.EnablePreserveAsPath {
		err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
		if err != nil {
			return diag.Errorf("could not enable spoke preserve as path after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.BgpPollingTime >= 10 && edgeEquinix.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeEquinix.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.BgpHoldTime >= 12 && edgeEquinix.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeEquinix.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeEquinix.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.Latitude != "" || edgeEquinix.Longitude != "" {
		gatewayForEaasFunctions.Latitude = edgeEquinix.Latitude
		gatewayForEaasFunctions.Longitude = edgeEquinix.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeEquinix.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.EnableSingleIpSnat {
		gatewayForGatewayFunctions.GatewayName = edgeEquinix.GwName
		err := client.EnableSNat(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("failed to enable single IP SNAT: %s", err)
		}
	}

	if edgeEquinix.EnableAutoAdvertiseLanCidrs == "disable" {
		err := client.UpdateEdgeEquinix(ctx, edgeEquinix)
		if err != nil {
			return diag.Errorf("could not disable auto advertise LAN CIDRs after Edge Equinix creation: %v", err)
		}
	}

	if edgeEquinix.EnableEdgeActiveStandby || edgeEquinix.EnableEdgeActiveStandbyPreemptive {
		err := client.UpdateEdgeEquinix(ctx, edgeEquinix)
		if err != nil {
			return diag.Errorf("could not update Edge active standby or Edge active standby preemptive after Edge Equinix creation: %v", err)
		}
	}

	return resourceAviatrixEdgeEquinixReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeEquinixReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeEquinixRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeEquinixRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// handle import
	if d.Get("gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	edgeEquinixResp, err := client.GetEdgeEquinix(ctx, d.Get("gw_name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Equinix: %v", err)
	}

	d.Set("account_name", edgeEquinixResp.AccountName)
	d.Set("gw_name", edgeEquinixResp.GwName)
	d.Set("site_id", edgeEquinixResp.SiteId)
	d.Set("enable_management_over_private_network", edgeEquinixResp.EnableManagementOverPrivateNetwork)
	d.Set("dns_server_ip", edgeEquinixResp.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeEquinixResp.SecondaryDnsServerIp)
	d.Set("local_as_number", edgeEquinixResp.LocalAsNumber)
	d.Set("prepend_as_path", edgeEquinixResp.PrependAsPath)
	d.Set("enable_edge_active_standby", edgeEquinixResp.EnableEdgeActiveStandby)
	d.Set("enable_edge_active_standby_preemptive", edgeEquinixResp.EnableEdgeActiveStandbyPreemptive)
	d.Set("enable_learned_cidrs_approval", edgeEquinixResp.EnableLearnedCidrsApproval)

	if edgeEquinixResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeEquinixResp.ManagementEgressIpPrefix, ","))
	}

	if edgeEquinixResp.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: edgeEquinixResp.GwName})
		if err != nil {
			return diag.Errorf("could not get advanced config for Edge Equinix: %v", err)
		}

		err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not set approved_learned_cidrs into state: %v", err)
		}
	} else {
		d.Set("approved_learned_cidrs", nil)
	}

	spokeBgpManualAdvertisedCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	if len(goaviatrix.Difference(spokeBgpManualAdvertisedCidrs, edgeEquinixResp.SpokeBgpManualAdvertisedCidrs)) != 0 ||
		len(goaviatrix.Difference(edgeEquinixResp.SpokeBgpManualAdvertisedCidrs, spokeBgpManualAdvertisedCidrs)) != 0 {
		d.Set("spoke_bgp_manual_advertise_cidrs", edgeEquinixResp.SpokeBgpManualAdvertisedCidrs)
	} else {
		d.Set("spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertisedCidrs)
	}

	d.Set("enable_preserve_as_path", edgeEquinixResp.EnablePreserveAsPath)
	d.Set("bgp_polling_time", edgeEquinixResp.BgpPollingTime)
	d.Set("bgp_hold_time", edgeEquinixResp.BgpHoldTime)
	d.Set("enable_edge_transitive_routing", edgeEquinixResp.EnableEdgeTransitiveRouting)
	d.Set("enable_jumbo_frame", edgeEquinixResp.EnableJumboFrame)
	if edgeEquinixResp.Latitude != 0 || edgeEquinixResp.Longitude != 0 {
		d.Set("latitude", fmt.Sprintf("%.6f", edgeEquinixResp.Latitude))
		d.Set("longitude", fmt.Sprintf("%.6f", edgeEquinixResp.Longitude))
	} else {
		d.Set("latitude", "")
		d.Set("longitude", "")
	}

	d.Set("rx_queue_size", edgeEquinixResp.RxQueueSize)
	d.Set("state", edgeEquinixResp.State)

	var interfaces []map[string]interface{}
	var vlan []map[string]interface{}
	for _, interface0 := range edgeEquinixResp.InterfaceList {
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
		interface1["vrrp_virtual_ip"] = interface0.VirtualIp

		if interface0.Type == "LAN" {
			interface1["enable_vrrp"] = interface0.VrrpState
		}

		if interface0.Type == "LAN" && interface0.SubInterfaces != nil {
			for _, vlan0 := range interface0.SubInterfaces {
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

		interfaces = append(interfaces, interface1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	if err = d.Set("vlan", vlan); err != nil {
		return diag.Errorf("failed to set vlan: %s\n", err)
	}

	d.Set("dns_profile_name", edgeEquinixResp.DnsProfileName)
	d.Set("enable_single_ip_snat", edgeEquinixResp.EnableNat == "yes" && edgeEquinixResp.SnatMode == "primary")
	d.Set("enable_auto_advertise_lan_cidrs", edgeEquinixResp.EnableAutoAdvertiseLanCidrs)

	d.SetId(edgeEquinixResp.GwName)
	return nil
}

func resourceAviatrixEdgeEquinixUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeEquinix := marshalEdgeEquinixInput(d)

	// checks before update
	if !edgeEquinix.EnableEdgeActiveStandby && edgeEquinix.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeEquinix.EnableLearnedCidrsApproval && len(edgeEquinix.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeEquinix.PrependAsPath) != 0 {
		if edgeEquinix.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeEquinix.Latitude != "" && edgeEquinix.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeEquinix.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeEquinix.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	d.Partial(true)

	// update configs
	// use following variables to reuse functions for transit, spoke and gateway
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeEquinix.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeEquinix.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeEquinix.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeEquinix.GwName,
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeEquinix.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge Equinix update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeEquinix.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge Equinix update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeEquinix.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeEquinix.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge Equinix update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeEquinix.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge Equinix update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge Equinix update: %v", err)
			}
		}
	}

	if edgeEquinix.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeEquinix.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge Equinix update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeEquinix.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge Equinix update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeEquinix.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge Equinix update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge Equinix update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeEquinix.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge Equinix update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeEquinix.GwName, edgeEquinix.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge Equinix update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeEquinix.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeEquinix.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge Equinix update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeEquinix.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge Equinix update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeEquinix.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge Equinix update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge Equinix update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		gatewayForEaasFunctions.Latitude = edgeEquinix.Latitude
		gatewayForEaasFunctions.Longitude = edgeEquinix.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge Equinix update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeEquinix.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge Equinix update: %v", err)
		}
	}

	if d.HasChanges("management_egress_ip_prefix_list", "interfaces", "vlan", "dns_profile_name",
		"enable_auto_advertise_lan_cidrs", "enable_edge_active_standby", "enable_edge_active_standby_preemptive") {
		err := client.UpdateEdgeEquinix(ctx, edgeEquinix)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list, WAN/LAN/VLAN interfaces, "+
				"DNS profile name, auto advertise LAN CIDRs, Edge active standby or Edge active standby preemptive "+
				"during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("enable_single_ip_snat") {
		gatewayForGatewayFunctions.GatewayName = edgeEquinix.GwName

		if edgeEquinix.EnableSingleIpSnat {
			err := client.EnableSNat(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("failed to enable single IP SNAT during Edge Equinix update: %s", err)
			}
		} else {
			err := client.DisableSNat(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("failed to disable single IP SNAT during Edge Equinix update: %s", err)
			}
		}

	}

	d.Partial(false)

	return resourceAviatrixEdgeEquinixRead(ctx, d, meta)
}

func resourceAviatrixEdgeEquinixDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	gwName := d.Get("gw_name").(string)
	siteId := d.Get("site_id").(string)
	ztpFileDownloadPath := d.Get("ztp_file_download_path").(string)

	err := client.DeleteEdgeEquinix(ctx, accountName, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge Equinix: %v", err)
	}

	fileName := ztpFileDownloadPath + "/" + gwName + "-" + siteId + "-cloud-init.txt"

	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
