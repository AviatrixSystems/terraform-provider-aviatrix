package aviatrix

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeNEO() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeNEOCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeNEORead,
		UpdateWithoutTimeout: resourceAviatrixEdgeNEOUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeNEODelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP account name.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP name.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"device_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge NEO device ID.",
			},
			"gw_size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"small", "medium", "large", "x-large"}, false),
				Description:  "Gateway size (CPU and Memory).",
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
				Description: "State of Edge as a Spoke.",
			},
			"wan_interface_names": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of WAN interface names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lan_interface_names": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of LAN interface names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"management_interface_names": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of management interface names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		DeprecationMessage: "Since V3.1.1+, please use resource aviatrix_edge_platform instead. Resource " +
			"aviatrix_edge_neo will be deprecated in the V3.2.0 release.",
	}
}

func marshalEdgeNEOInput(d *schema.ResourceData) *goaviatrix.EdgeNEO {
	edgeNEO := &goaviatrix.EdgeNEO{
		AccountName:                        d.Get("account_name").(string),
		GwName:                             d.Get("gw_name").(string),
		SiteId:                             d.Get("site_id").(string),
		DeviceId:                           d.Get("device_id").(string),
		GwSize:                             d.Get("gw_size").(string),
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
		WanInterface:                       strings.Join(getStringList(d, "wan_interface_names"), ","),
		LanInterface:                       strings.Join(getStringList(d, "lan_interface_names"), ","),
		MgmtInterface:                      strings.Join(getStringList(d, "management_interface_names"), ","),
		DnsProfileName:                     d.Get("dns_profile_name").(string),
		EnableSingleIpSnat:                 d.Get("enable_single_ip_snat").(bool),
	}

	interfaces := d.Get("interfaces").(*schema.Set).List()
	for _, interface0 := range interfaces {
		interface1 := interface0.(map[string]interface{})

		interface2 := &goaviatrix.EdgeNEOInterface{
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

		edgeNEO.InterfaceList = append(edgeNEO.InterfaceList, interface2)
	}

	vlan := d.Get("vlan").(*schema.Set).List()
	for _, vlan0 := range vlan {
		vlan1 := vlan0.(map[string]interface{})

		vlan2 := &goaviatrix.EdgeNEOVlan{
			ParentInterface: vlan1["parent_interface_name"].(string),
			IpAddr:          vlan1["ip_address"].(string),
			GatewayIp:       vlan1["gateway_ip"].(string),
			PeerIpAddr:      vlan1["peer_ip_address"].(string),
			PeerGatewayIp:   vlan1["peer_gateway_ip"].(string),
			VirtualIp:       vlan1["vrrp_virtual_ip"].(string),
			Tag:             vlan1["tag"].(string),
		}

		vlan2.VlanId = strconv.Itoa(vlan1["vlan_id"].(int))

		edgeNEO.VlanList = append(edgeNEO.VlanList, vlan2)
	}

	if d.Get("enable_auto_advertise_lan_cidrs").(bool) {
		edgeNEO.EnableAutoAdvertiseLanCidrs = "enable"
	} else {
		edgeNEO.EnableAutoAdvertiseLanCidrs = "disable"
	}

	return edgeNEO
}

func resourceAviatrixEdgeNEOCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeNEO := marshalEdgeNEOInput(d)

	// checks before creation
	if !edgeNEO.EnableEdgeActiveStandby && edgeNEO.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeNEO.EnableLearnedCidrsApproval && len(edgeNEO.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeNEO.PrependAsPath) != 0 {
		if edgeNEO.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeNEO.Latitude != "" && edgeNEO.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeNEO.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeNEO.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	// create
	d.SetId(edgeNEO.GwName)
	flag := false
	defer resourceAviatrixEdgeNEOReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeNEO(ctx, edgeNEO); err != nil {
		return diag.Errorf("could not create Edge NEO: %v", err)
	}

	// advanced configs
	// use following variables to reuse functions for transit, spoke, gateway and EaaS
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeNEO.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeNEO.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeNEO.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeNEO.GwName,
	}

	if edgeNEO.LocalAsNumber != "" {
		err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeNEO.LocalAsNumber)
		if err != nil {
			return diag.Errorf("could not set 'local_as_number' after Edge NEO creation: %v", err)
		}
	}

	if len(edgeNEO.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeNEO.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge NEO creation: %v", err)
		}
	}

	if len(edgeNEO.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeNEO.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge NEO creation: %v", err)
		}
	}

	if len(edgeNEO.SpokeBgpManualAdvertisedCidrs) != 0 {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeNEO.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.EnablePreserveAsPath {
		err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
		if err != nil {
			return diag.Errorf("could not enable spoke preserve as path after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.BgpPollingTime >= 10 && edgeNEO.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeNEO.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.BgpHoldTime >= 12 && edgeNEO.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeNEO.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeNEO.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.Latitude != "" || edgeNEO.Longitude != "" {
		gatewayForEaasFunctions.Latitude = edgeNEO.Latitude
		gatewayForEaasFunctions.Longitude = edgeNEO.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeNEO.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.EnableSingleIpSnat {
		gatewayForGatewayFunctions.GatewayName = edgeNEO.GwName
		err := client.EnableSNat(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("failed to enable single IP SNAT: %s", err)
		}
	}

	if edgeNEO.EnableAutoAdvertiseLanCidrs == "disable" {
		err := client.UpdateEdgeNEO(ctx, edgeNEO)
		if err != nil {
			return diag.Errorf("could not disable auto advertise LAN CIDRs after Edge NEO creation: %v", err)
		}
	}

	if edgeNEO.EnableEdgeActiveStandby || edgeNEO.EnableEdgeActiveStandbyPreemptive {
		err := client.UpdateEdgeNEO(ctx, edgeNEO)
		if err != nil {
			return diag.Errorf("could not update Edge active standby or Edge active standby preemptive after Edge NEO creation: %v", err)
		}
	}

	return resourceAviatrixEdgeNEOReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeNEOReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeNEORead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeNEORead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// handle import
	if d.Get("gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	edgeNEOResp, err := client.GetEdgeNEO(ctx, d.Get("gw_name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge NEO: %v", err)
	}

	d.Set("account_name", edgeNEOResp.AccountName)
	d.Set("gw_name", edgeNEOResp.GwName)
	d.Set("site_id", edgeNEOResp.SiteId)
	d.Set("device_id", edgeNEOResp.DeviceId)
	d.Set("gw_size", edgeNEOResp.GwSize)
	d.Set("enable_management_over_private_network", edgeNEOResp.EnableManagementOverPrivateNetwork)
	d.Set("dns_server_ip", edgeNEOResp.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeNEOResp.SecondaryDnsServerIp)
	d.Set("local_as_number", edgeNEOResp.LocalAsNumber)
	d.Set("prepend_as_path", edgeNEOResp.PrependAsPath)
	d.Set("enable_edge_active_standby", edgeNEOResp.EnableEdgeActiveStandby)
	d.Set("enable_edge_active_standby_preemptive", edgeNEOResp.EnableEdgeActiveStandbyPreemptive)
	d.Set("enable_learned_cidrs_approval", edgeNEOResp.EnableLearnedCidrsApproval)

	if edgeNEOResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeNEOResp.ManagementEgressIpPrefix, ","))
	}

	if edgeNEOResp.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: edgeNEOResp.GwName})
		if err != nil {
			return diag.Errorf("could not get advanced config for Edge NEO: %v", err)
		}

		err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not set approved_learned_cidrs into state: %v", err)
		}
	} else {
		d.Set("approved_learned_cidrs", nil)
	}

	spokeBgpManualAdvertisedCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	if len(goaviatrix.Difference(spokeBgpManualAdvertisedCidrs, edgeNEOResp.SpokeBgpManualAdvertisedCidrs)) != 0 ||
		len(goaviatrix.Difference(edgeNEOResp.SpokeBgpManualAdvertisedCidrs, spokeBgpManualAdvertisedCidrs)) != 0 {
		d.Set("spoke_bgp_manual_advertise_cidrs", edgeNEOResp.SpokeBgpManualAdvertisedCidrs)
	} else {
		d.Set("spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertisedCidrs)
	}

	d.Set("enable_preserve_as_path", edgeNEOResp.EnablePreserveAsPath)
	d.Set("bgp_polling_time", edgeNEOResp.BgpPollingTime)
	d.Set("bgp_hold_time", edgeNEOResp.BgpHoldTime)
	d.Set("enable_edge_transitive_routing", edgeNEOResp.EnableEdgeTransitiveRouting)
	d.Set("enable_jumbo_frame", edgeNEOResp.EnableJumboFrame)
	if edgeNEOResp.Latitude != 0 || edgeNEOResp.Longitude != 0 {
		d.Set("latitude", fmt.Sprintf("%.6f", edgeNEOResp.Latitude))
		d.Set("longitude", fmt.Sprintf("%.6f", edgeNEOResp.Longitude))
	} else {
		d.Set("latitude", "")
		d.Set("longitude", "")
	}

	d.Set("rx_queue_size", edgeNEOResp.RxQueueSize)
	d.Set("state", edgeNEOResp.State)
	d.Set("wan_interface_names", edgeNEOResp.WanInterface)
	d.Set("lan_interface_names", edgeNEOResp.LanInterface)
	d.Set("management_interface_names", edgeNEOResp.MgmtInterface)

	var interfaces []map[string]interface{}
	var vlan []map[string]interface{}
	for _, interface0 := range edgeNEOResp.InterfaceList {
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

	d.Set("dns_profile_name", edgeNEOResp.DnsProfileName)
	d.Set("enable_single_ip_snat", edgeNEOResp.EnableNat == "yes" && edgeNEOResp.SnatMode == "primary")
	d.Set("enable_auto_advertise_lan_cidrs", edgeNEOResp.EnableAutoAdvertiseLanCidrs)

	d.SetId(edgeNEOResp.GwName)
	return nil
}

func resourceAviatrixEdgeNEOUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeNEO := marshalEdgeNEOInput(d)

	// checks before update
	if !edgeNEO.EnableEdgeActiveStandby && edgeNEO.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeNEO.EnableLearnedCidrsApproval && len(edgeNEO.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeNEO.PrependAsPath) != 0 {
		if edgeNEO.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeNEO.Latitude != "" && edgeNEO.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeNEO.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeNEO.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	d.Partial(true)

	// update configs
	// use following variables to reuse functions for transit, spoke and gateway
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeNEO.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeNEO.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeNEO.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeNEO.GwName,
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeNEO.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge NEO update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeNEO.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge NEO update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeNEO.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeNEO.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge NEO update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeNEO.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge NEO update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge NEO update: %v", err)
			}
		}
	}

	if edgeNEO.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeNEO.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge NEO update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeNEO.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge NEO update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeNEO.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge NEO update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge NEO update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeNEO.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge NEO update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeNEO.GwName, edgeNEO.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge NEO update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeNEO.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeNEO.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge NEO update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeNEO.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge NEO update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeNEO.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge NEO update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge NEO update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		gatewayForEaasFunctions.Latitude = edgeNEO.Latitude
		gatewayForEaasFunctions.Longitude = edgeNEO.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge NEO update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeNEO.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge NEO update: %v", err)
		}
	}

	if d.HasChanges("management_egress_ip_prefix_list", "interfaces", "vlan", "dns_profile_name",
		"enable_auto_advertise_lan_cidrs", "enable_edge_active_standby", "enable_edge_active_standby_preemptive") {
		err := client.UpdateEdgeNEO(ctx, edgeNEO)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list, WAN/LAN/VLAN interfaces, "+
				"DNS profile name, auto advertise LAN CIDRs, Edge active standby or Edge active standby preemptive "+
				"during Edge NEO update: %v", err)
		}
	}

	if d.HasChange("enable_single_ip_snat") {
		gatewayForGatewayFunctions.GatewayName = edgeNEO.GwName

		if edgeNEO.EnableSingleIpSnat {
			err := client.EnableSNat(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("failed to enable single IP SNAT during update: %s", err)
			}
		} else {
			err := client.DisableSNat(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("failed to disable single IP SNAT during update: %s", err)
			}
		}

	}

	d.Partial(false)

	return resourceAviatrixEdgeNEORead(ctx, d, meta)
}

func resourceAviatrixEdgeNEODelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	gwName := d.Get("gw_name").(string)

	err := client.DeleteEdgeNEO(ctx, accountName, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge NEO: %v", err)
	}

	return nil
}
