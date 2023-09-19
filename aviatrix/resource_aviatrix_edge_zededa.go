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

func resourceAviatrixEdgeZededa() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeZededaCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeZededaRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeZededaUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeZededaDelete,
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
			"project_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP project UUID.",
			},
			"compute_node_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP compute node UUID.",
			},
			"template_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP template UUID.",
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
				Optional:    true,
				ForceNew:    true,
				Description: "List of WAN interface names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lan_interface_names": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of LAN interface names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"management_interface_names": {
				Type:        schema.TypeList,
				Optional:    true,
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
	}
}

func marshalEdgeZededaInput(d *schema.ResourceData) *goaviatrix.EdgeCSP {
	edgeCSP := &goaviatrix.EdgeCSP{
		AccountName:                        d.Get("account_name").(string),
		GwName:                             d.Get("gw_name").(string),
		SiteId:                             d.Get("site_id").(string),
		ProjectUuid:                        d.Get("project_uuid").(string),
		ComputeNodeUuid:                    d.Get("compute_node_uuid").(string),
		TemplateUuid:                       d.Get("template_uuid").(string),
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
	for _, if0 := range interfaces {
		if1 := if0.(map[string]interface{})

		if2 := &goaviatrix.Interface{
			IfName:       if1["name"].(string),
			Type:         if1["type"].(string),
			Bandwidth:    if1["bandwidth"].(int),
			PublicIp:     if1["wan_public_ip"].(string),
			Tag:          if1["tag"].(string),
			Dhcp:         if1["enable_dhcp"].(bool),
			IpAddr:       if1["ip_address"].(string),
			GatewayIp:    if1["gateway_ip"].(string),
			DnsPrimary:   if1["dns_server_ip"].(string),
			DnsSecondary: if1["secondary_dns_server_ip"].(string),
			VrrpState:    if1["enable_vrrp"].(bool),
			VirtualIp:    if1["vrrp_virtual_ip"].(string),
		}

		edgeCSP.InterfaceList = append(edgeCSP.InterfaceList, if2)
	}

	vlan := d.Get("vlan").(*schema.Set).List()
	for _, v0 := range vlan {
		v1 := v0.(map[string]interface{})

		v2 := &goaviatrix.Vlan{
			ParentInterface: v1["parent_interface_name"].(string),
			IpAddr:          v1["ip_address"].(string),
			GatewayIp:       v1["gateway_ip"].(string),
			PeerIpAddr:      v1["peer_ip_address"].(string),
			PeerGatewayIp:   v1["peer_gateway_ip"].(string),
			VirtualIp:       v1["vrrp_virtual_ip"].(string),
			Tag:             v1["tag"].(string),
		}

		v2.VlanId = strconv.Itoa(v1["vlan_id"].(int))

		edgeCSP.VlanList = append(edgeCSP.VlanList, v2)
	}

	if d.Get("enable_auto_advertise_lan_cidrs").(bool) {
		edgeCSP.EnableAutoAdvertiseLanCidrs = "enable"
	} else {
		edgeCSP.EnableAutoAdvertiseLanCidrs = "disable"
	}

	return edgeCSP
}

func resourceAviatrixEdgeZededaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeCSP := marshalEdgeZededaInput(d)

	// checks before creation
	if !edgeCSP.EnableEdgeActiveStandby && edgeCSP.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeCSP.EnableLearnedCidrsApproval && len(edgeCSP.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeCSP.PrependAsPath) != 0 {
		if edgeCSP.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeCSP.Latitude != "" && edgeCSP.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeCSP.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeCSP.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	// create
	d.SetId(edgeCSP.GwName)
	flag := false
	defer resourceAviatrixEdgeZededaReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeCSP(ctx, edgeCSP); err != nil {
		return diag.Errorf("could not create Edge Zededa: %v", err)
	}

	// advanced configs
	// use following variables to reuse functions for transit, spoke, gateway and EaaS
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeCSP.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeCSP.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeCSP.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeCSP.GwName,
	}

	if edgeCSP.LocalAsNumber != "" {
		err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeCSP.LocalAsNumber)
		if err != nil {
			return diag.Errorf("could not set 'local_as_number' after Edge Zededa creation: %v", err)
		}
	}

	if len(edgeCSP.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeCSP.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge Zededa creation: %v", err)
		}
	}

	if len(edgeCSP.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeCSP.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge Zededa creation: %v", err)
		}
	}

	if len(edgeCSP.SpokeBgpManualAdvertisedCidrs) != 0 {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeCSP.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.EnablePreserveAsPath {
		err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
		if err != nil {
			return diag.Errorf("could not enable spoke preserve as path after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.BgpPollingTime >= 10 && edgeCSP.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeCSP.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.BgpHoldTime >= 12 && edgeCSP.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeCSP.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeCSP.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.Latitude != "" || edgeCSP.Longitude != "" {
		gatewayForEaasFunctions.Latitude = edgeCSP.Latitude
		gatewayForEaasFunctions.Longitude = edgeCSP.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeCSP.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.EnableSingleIpSnat {
		gatewayForGatewayFunctions.GatewayName = edgeCSP.GwName
		err := client.EnableSNat(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("failed to enable single IP SNAT: %s", err)
		}
	}

	if edgeCSP.EnableAutoAdvertiseLanCidrs == "disable" {
		err := client.UpdateEdgeCSP(ctx, edgeCSP)
		if err != nil {
			return diag.Errorf("could not disable auto advertise LAN CIDRs after Edge Zededa creation: %v", err)
		}
	}

	if edgeCSP.EnableEdgeActiveStandby || edgeCSP.EnableEdgeActiveStandbyPreemptive {
		err := client.UpdateEdgeCSP(ctx, edgeCSP)
		if err != nil {
			return diag.Errorf("could not update Edge active standby or Edge active standby preemptive after Edge Zededa creation: %v", err)
		}
	}

	return resourceAviatrixEdgeZededaReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeZededaReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeZededaRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeZededaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// handle import
	if d.Get("gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	edgeCSPResp, err := client.GetEdgeCSP(ctx, d.Get("gw_name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Zededa: %v", err)
	}

	d.Set("account_name", edgeCSPResp.AccountName)
	d.Set("gw_name", edgeCSPResp.GwName)
	d.Set("site_id", edgeCSPResp.SiteId)
	d.Set("project_uuid", edgeCSPResp.ProjectUuid)
	d.Set("compute_node_uuid", edgeCSPResp.ComputeNodeUuid)
	d.Set("template_uuid", edgeCSPResp.TemplateUuid)
	d.Set("enable_management_over_private_network", edgeCSPResp.EnableManagementOverPrivateNetwork)
	d.Set("dns_server_ip", edgeCSPResp.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeCSPResp.SecondaryDnsServerIp)
	d.Set("local_as_number", edgeCSPResp.LocalAsNumber)
	d.Set("prepend_as_path", edgeCSPResp.PrependAsPath)
	d.Set("enable_edge_active_standby", edgeCSPResp.EnableEdgeActiveStandby)
	d.Set("enable_edge_active_standby_preemptive", edgeCSPResp.EnableEdgeActiveStandbyPreemptive)
	d.Set("enable_learned_cidrs_approval", edgeCSPResp.EnableLearnedCidrsApproval)

	if edgeCSPResp.ManagementEgressIpPrefix == "" {
		d.Set("management_egress_ip_prefix_list", nil)
	} else {
		d.Set("management_egress_ip_prefix_list", strings.Split(edgeCSPResp.ManagementEgressIpPrefix, ","))
	}

	if edgeCSPResp.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: edgeCSPResp.GwName})
		if err != nil {
			return diag.Errorf("could not get advanced config for Edge Zededa: %v", err)
		}

		err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not set approved_learned_cidrs into state: %v", err)
		}
	} else {
		d.Set("approved_learned_cidrs", nil)
	}

	spokeBgpManualAdvertisedCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	if len(goaviatrix.Difference(spokeBgpManualAdvertisedCidrs, edgeCSPResp.SpokeBgpManualAdvertisedCidrs)) != 0 ||
		len(goaviatrix.Difference(edgeCSPResp.SpokeBgpManualAdvertisedCidrs, spokeBgpManualAdvertisedCidrs)) != 0 {
		d.Set("spoke_bgp_manual_advertise_cidrs", edgeCSPResp.SpokeBgpManualAdvertisedCidrs)
	} else {
		d.Set("spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertisedCidrs)
	}

	d.Set("enable_preserve_as_path", edgeCSPResp.EnablePreserveAsPath)
	d.Set("bgp_polling_time", edgeCSPResp.BgpPollingTime)
	d.Set("bgp_hold_time", edgeCSPResp.BgpHoldTime)
	d.Set("enable_edge_transitive_routing", edgeCSPResp.EnableEdgeTransitiveRouting)
	d.Set("enable_jumbo_frame", edgeCSPResp.EnableJumboFrame)
	if edgeCSPResp.Latitude != 0 || edgeCSPResp.Longitude != 0 {
		d.Set("latitude", fmt.Sprintf("%.6f", edgeCSPResp.Latitude))
		d.Set("longitude", fmt.Sprintf("%.6f", edgeCSPResp.Longitude))
	} else {
		d.Set("latitude", "")
		d.Set("longitude", "")
	}

	d.Set("rx_queue_size", edgeCSPResp.RxQueueSize)
	d.Set("state", edgeCSPResp.State)
	d.Set("wan_interface_names", edgeCSPResp.WanInterface)
	d.Set("lan_interface_names", edgeCSPResp.LanInterface)
	d.Set("management_interface_names", edgeCSPResp.MgmtInterface)

	var interfaces []map[string]interface{}
	var vlan []map[string]interface{}
	for _, if0 := range edgeCSPResp.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
		if1["bandwidth"] = if0.Bandwidth
		if1["wan_public_ip"] = if0.PublicIp
		if1["tag"] = if0.Tag
		if1["enable_dhcp"] = if0.Dhcp
		if1["ip_address"] = if0.IpAddr
		if1["gateway_ip"] = if0.GatewayIp
		if1["dns_server_ip"] = if0.DnsPrimary
		if1["secondary_dns_server_ip"] = if0.DnsSecondary
		if1["vrrp_virtual_ip"] = if0.VirtualIp

		if if0.Type == "LAN" {
			if1["enable_vrrp"] = if0.VrrpState
		}

		if if0.Type == "LAN" && if0.SubInterfaces != nil {
			for _, v0 := range if0.SubInterfaces {
				v1 := make(map[string]interface{})
				v1["parent_interface_name"] = v0.ParentInterface
				v1["ip_address"] = v0.IpAddr
				v1["gateway_ip"] = v0.GatewayIp
				v1["peer_ip_address"] = v0.PeerIpAddr
				v1["peer_gateway_ip"] = v0.PeerGatewayIp
				v1["vrrp_virtual_ip"] = v0.VirtualIp
				v1["tag"] = v0.Tag

				vlanId, _ := strconv.Atoi(v0.VlanId)
				v1["vlan_id"] = vlanId

				vlan = append(vlan, v1)
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

	d.Set("dns_profile_name", edgeCSPResp.DnsProfileName)
	d.Set("enable_single_ip_snat", edgeCSPResp.EnableNat == "yes" && edgeCSPResp.SnatMode == "primary")
	d.Set("enable_auto_advertise_lan_cidrs", edgeCSPResp.EnableAutoAdvertiseLanCidrs)

	d.SetId(edgeCSPResp.GwName)
	return nil
}

func resourceAviatrixEdgeZededaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeCSP := marshalEdgeZededaInput(d)

	// checks before update
	if !edgeCSP.EnableEdgeActiveStandby && edgeCSP.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeCSP.EnableLearnedCidrsApproval && len(edgeCSP.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeCSP.PrependAsPath) != 0 {
		if edgeCSP.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeCSP.Latitude != "" && edgeCSP.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeCSP.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeCSP.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	d.Partial(true)

	// update configs
	// use following variables to reuse functions for transit, spoke and gateway
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeCSP.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeCSP.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeCSP.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeCSP.GwName,
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeCSP.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge Zededa update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeCSP.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge Zededa update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeCSP.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeCSP.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge Zededa update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeCSP.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge Zededa update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge Zededa update: %v", err)
			}
		}
	}

	if edgeCSP.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeCSP.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge Zededa update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeCSP.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge Zededa update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeCSP.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge Zededa update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge Zededa update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeCSP.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge Zededa update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeCSP.GwName, edgeCSP.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge Zededa update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeCSP.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeCSP.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge Zededa update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeCSP.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge Zededa update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeCSP.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge Zededa update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge Zededa update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		gatewayForEaasFunctions.Latitude = edgeCSP.Latitude
		gatewayForEaasFunctions.Longitude = edgeCSP.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge Zededa update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeCSP.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge Zededa update: %v", err)
		}
	}

	if d.HasChanges("management_egress_ip_prefix_list", "interfaces", "vlan", "dns_profile_name",
		"enable_auto_advertise_lan_cidrs", "enable_edge_active_standby", "enable_edge_active_standby_preemptive") {
		err := client.UpdateEdgeCSP(ctx, edgeCSP)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list, WAN/LAN/VLAN interfaces, "+
				"DNS profile name, auto advertise LAN CIDRs, Edge active standby or Edge active standby preemptive "+
				"during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("enable_single_ip_snat") {
		gatewayForGatewayFunctions.GatewayName = edgeCSP.GwName

		if edgeCSP.EnableSingleIpSnat {
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

	return resourceAviatrixEdgeZededaRead(ctx, d, meta)
}

func resourceAviatrixEdgeZededaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	gwName := d.Get("gw_name").(string)

	err := client.DeleteEdgeCSP(ctx, accountName, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge Zededa: %v", err)
	}

	return nil
}
