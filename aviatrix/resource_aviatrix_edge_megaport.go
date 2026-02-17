package aviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixEdgeMegaport() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeMegaportCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeMegaportRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeMegaportUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeMegaportDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Megaport account name.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Megaport name.",
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
				Description: "The location where the ZTP file will be stored locally.",
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
			"bgp_neighbor_status_polling_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpNeighborStatusPollingTime,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "BGP neighbor status polling time for BGP Spoke Gateway. Unit is in seconds. Valid values are between 1 and 10.",
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
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of WAN/LAN/MANAGEMENT interfaces, each represented as a map.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logical_ifname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logical interface name e.g., wan0, lan0, mgmt0.",
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^(wan|lan|mgmt)[0-9]+$`),
								"Logical interface name must start with 'wan', 'lan', or 'mgmt' followed by a number (e.g., 'wan0', 'lan1', 'mgmt2').",
							),
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "VLAN configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_logical_interface_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parent logical interface name e.g. lan0",
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
				Deprecated:  "DNS profile support has been removed.",
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
			"interface_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of interface names mapped to interface types and indices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Interface name (e.g., 'eth0', 'eth1').",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Interface type (e.g., 'wan', 'lan', 'mgmt').",
							ValidateFunc: validation.StringInSlice([]string{"wan", "lan", "mgmt"}, false),
						},
						"index": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Interface index (e.g., 0, 1).",
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
			"included_advertised_spoke_routes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A list of CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func marshalEdgeMegaportInput(d *schema.ResourceData) (*goaviatrix.EdgeMegaport, error) {
	edgeMegaport := &goaviatrix.EdgeMegaport{
		AccountName:                        getString(d, "account_name"),
		GwName:                             getString(d, "gw_name"),
		SiteID:                             getString(d, "site_id"),
		ZtpFileDownloadPath:                getString(d, "ztp_file_download_path"),
		ManagementEgressIPPrefix:           strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
		EnableManagementOverPrivateNetwork: getBool(d, "enable_management_over_private_network"),
		DNSServerIP:                        getString(d, "dns_server_ip"),
		SecondaryDNSServerIP:               getString(d, "secondary_dns_server_ip"),
		EnableEdgeActiveStandby:            getBool(d, "enable_edge_active_standby"),
		EnableEdgeActiveStandbyPreemptive:  getBool(d, "enable_edge_active_standby_preemptive"),
		LocalAsNumber:                      getString(d, "local_as_number"),
		PrependAsPath:                      getStringList(d, "prepend_as_path"),
		EnableLearnedCidrsApproval:         getBool(d, "enable_learned_cidrs_approval"),
		ApprovedLearnedCidrs:               getStringSet(d, "approved_learned_cidrs"),
		SpokeBgpManualAdvertisedCidrs:      getStringSet(d, "spoke_bgp_manual_advertise_cidrs"),
		EnablePreserveAsPath:               getBool(d, "enable_preserve_as_path"),
		BgpPollingTime:                     getInt(d, "bgp_polling_time"),
		BgpBfdPollingTime:                  getInt(d, "bgp_neighbor_status_polling_time"),
		BgpHoldTime:                        getInt(d, "bgp_hold_time"),
		EnableEdgeTransitiveRouting:        getBool(d, "enable_edge_transitive_routing"),
		EnableJumboFrame:                   getBool(d, "enable_jumbo_frame"),
		Latitude:                           getString(d, "latitude"),
		Longitude:                          getString(d, "longitude"),
		RxQueueSize:                        getString(d, "rx_queue_size"),
		EnableSingleIpSnat:                 getBool(d, "enable_single_ip_snat"),
	}

	interfaces := getList(d, "interfaces")
	for _, interface0 := range interfaces {
		interface1 := mustMap(interface0)

		interface2 := &goaviatrix.EdgeMegaportInterface{
			LogicalInterfaceName: mustString(interface1["logical_ifname"]),
			PublicIP:             mustString(interface1["wan_public_ip"]),
			Tag:                  mustString(interface1["tag"]),
			Dhcp:                 mustBool(interface1["enable_dhcp"]),
			IPAddr:               mustString(interface1["ip_address"]),
			GatewayIP:            mustString(interface1["gateway_ip"]),
			DNSPrimary:           mustString(interface1["dns_server_ip"]),
			DNSSecondary:         mustString(interface1["secondary_dns_server_ip"]),
		}

		// vrrp_state and virtual_ip are only applicable for LAN interfaces
		logicalIfname := mustString(interface1["logical_ifname"])
		enableVrrp := mustBool(interface1["enable_vrrp"])
		if strings.HasPrefix(logicalIfname, "lan") && enableVrrp {
			interface2.VrrpState = enableVrrp
			interface2.VirtualIP = mustString(interface1["vrrp_virtual_ip"])
		}

		edgeMegaport.InterfaceList = append(edgeMegaport.InterfaceList, interface2)
	}

	vlan := getList(d, "vlan")
	for _, vlan0 := range vlan {
		vlan1 := mustMap(vlan0)

		vlan2 := &goaviatrix.EdgeMegaportVlan{
			ParentLogicalInterface: mustString(vlan1["parent_logical_interface_name"]),
			IPAddr:                 mustString(vlan1["ip_address"]),
			GatewayIP:              mustString(vlan1["gateway_ip"]),
			PeerIPAddr:             mustString(vlan1["peer_ip_address"]),
			PeerGatewayIP:          mustString(vlan1["peer_gateway_ip"]),
			VirtualIP:              mustString(vlan1["vrrp_virtual_ip"]),
			Tag:                    mustString(vlan1["tag"]),
		}

		vlandID := mustInt(vlan1["vlan_id"])
		vlan2.VlanID = strconv.Itoa(vlandID)

		edgeMegaport.VlanList = append(edgeMegaport.VlanList, vlan2)
	}

	interfaceMapping := map[string][]string{}
	interfaceMappingList := getList(d, "interface_mapping")

	if len(interfaceMappingList) > 0 {
		// get the user provided interface mapping
		for _, value := range interfaceMappingList {
			mappingMap := mustMap(value)
			interfaceName := mustString(mappingMap["name"])
			interfaceType := mustString(mappingMap["type"])
			interfaceIndex := mustInt(mappingMap["index"])
			interfaceMapping[interfaceName] = []string{interfaceType, strconv.Itoa(interfaceIndex)}
		}
	} else {
		// get the default interface mapping for megaport
		interfaceMapping = map[string][]string{
			"eth0": {"wan", "0"},
			"eth1": {"lan", "0"},
			"eth2": {"mgmt", "0"},
			"eth3": {"wan", "1"},
			"eth4": {"wan", "2"},
		}
	}
	interfaceMappingJSON, err := json.Marshal(interfaceMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal interface mapping to json: %w", err)
	}
	edgeMegaport.InterfaceMapping = string(interfaceMappingJSON)

	if getBool(d, "enable_auto_advertise_lan_cidrs") {
		edgeMegaport.EnableAutoAdvertiseLanCidrs = "enable"
	} else {
		edgeMegaport.EnableAutoAdvertiseLanCidrs = "disable"
	}

	return edgeMegaport, nil
}

func resourceAviatrixEdgeMegaportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// read configs
	edgeMegaport, err := marshalEdgeMegaportInput(d)
	if err != nil {
		return diag.Errorf("error reading Edge Megaport configuration: %s", err)
	}

	// checks before creation
	if !edgeMegaport.EnableEdgeActiveStandby && edgeMegaport.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not enable Preemptive Mode when Active-Standby is disabled")
	}

	if !edgeMegaport.EnableLearnedCidrsApproval && len(edgeMegaport.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeMegaport.PrependAsPath) != 0 {
		if edgeMegaport.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeMegaport.Latitude != "" && edgeMegaport.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeMegaport.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeMegaport.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	// create
	d.SetId(edgeMegaport.GwName)
	flag := false
	defer resourceAviatrixEdgeMegaportReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeMegaport(ctx, edgeMegaport); err != nil {
		return diag.Errorf("could not create Edge Megaport %s: %v", edgeMegaport.GwName, err)
	}

	// advanced configs
	// use following variables to reuse functions for transit, spoke, gateway and EaaS
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeMegaport.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeMegaport.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeMegaport.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeMegaport.GwName,
	}

	if edgeMegaport.LocalAsNumber != "" {
		err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeMegaport.LocalAsNumber)
		if err != nil {
			return diag.Errorf("could not set 'local_as_number' after Edge Megaport creation: %v", err)
		}
	}

	if len(edgeMegaport.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeMegaport.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge Megaport creation: %v", err)
		}
	}

	if len(edgeMegaport.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeMegaport.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge Megaport creation: %v", err)
		}
	}

	if len(edgeMegaport.SpokeBgpManualAdvertisedCidrs) != 0 {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeMegaport.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.EnablePreserveAsPath {
		err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
		if err != nil {
			return diag.Errorf("could not enable spoke preserve as path after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.BgpPollingTime >= 10 && edgeMegaport.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, edgeMegaport.BgpPollingTime)
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.BgpBfdPollingTime >= 1 && edgeMegaport.BgpBfdPollingTime != defaultBgpNeighborStatusPollingTime {
		err := client.SetBgpBfdPollingTimeSpoke(gatewayForSpokeFunctions, edgeMegaport.BgpBfdPollingTime)
		if err != nil {
			return diag.Errorf("could not set bgp neighbor status polling time after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.BgpHoldTime >= 12 && edgeMegaport.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeMegaport.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeMegaport.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.Latitude != "" || edgeMegaport.Longitude != "" {
		gatewayForEaasFunctions.Latitude = edgeMegaport.Latitude
		gatewayForEaasFunctions.Longitude = edgeMegaport.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeMegaport.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.EnableSingleIpSnat {
		gatewayForGatewayFunctions.GatewayName = edgeMegaport.GwName
		err := client.EnableSNat(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("failed to enable single IP SNAT: %s", err)
		}
	}

	if edgeMegaport.EnableAutoAdvertiseLanCidrs == "disable" {
		err := client.UpdateEdgeMegaport(ctx, edgeMegaport)
		if err != nil {
			return diag.Errorf("could not disable auto advertise LAN CIDRs after Edge Megaport creation: %v", err)
		}
	}

	if edgeMegaport.EnableEdgeActiveStandby || edgeMegaport.EnableEdgeActiveStandbyPreemptive {
		err := client.UpdateEdgeMegaport(ctx, edgeMegaport)
		if err != nil {
			return diag.Errorf("could not update Edge active standby or Edge active standby preemptive after Edge Megaport creation: %v", err)
		}
	}

	err = editAdvertisedSpokeRoutesWithRetry(client, gatewayForGatewayFunctions, d)
	if err != nil {
		return diag.Errorf("failed to edit advertised spoke vpc routes of spoke gateway: %q: %s", gatewayForGatewayFunctions.GwName, err)
	}

	return resourceAviatrixEdgeMegaportReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeMegaportReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeMegaportRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeMegaportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// handle import
	if getString(d, "gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		_ = d.Set("gw_name", id)
		d.SetId(id)
	}

	edgeMegaportResp, err := client.GetEdgeMegaport(ctx, getString(d, "gw_name"))
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge Megaport: %v", err)
	}

	edgeMegaportFields := map[string]interface{}{
		"account_name":                           edgeMegaportResp.AccountName,
		"gw_name":                                edgeMegaportResp.GwName,
		"site_id":                                edgeMegaportResp.SiteId,
		"enable_management_over_private_network": edgeMegaportResp.EnableManagementOverPrivateNetwork,
		"dns_server_ip":                          edgeMegaportResp.DNSServerIP,
		"secondary_dns_server_ip":                edgeMegaportResp.SecondaryDNSServerIP,
		"local_as_number":                        edgeMegaportResp.LocalAsNumber,
		"prepend_as_path":                        edgeMegaportResp.PrependAsPath,
		"enable_edge_active_standby":             edgeMegaportResp.EnableEdgeActiveStandby,
		"enable_edge_active_standby_preemptive":  edgeMegaportResp.EnableEdgeActiveStandbyPreemptive,
		"enable_learned_cidrs_approval":          edgeMegaportResp.EnableLearnedCidrsApproval,
		"enable_preserve_as_path":                edgeMegaportResp.EnablePreserveAsPath,
		"bgp_polling_time":                       edgeMegaportResp.BgpPollingTime,
		"bgp_neighbor_status_polling_time":       edgeMegaportResp.BgpBfdPollingTime,
		"bgp_hold_time":                          edgeMegaportResp.BgpHoldTime,
		"enable_edge_transitive_routing":         edgeMegaportResp.EnableEdgeTransitiveRouting,
		"enable_jumbo_frame":                     edgeMegaportResp.EnableJumboFrame,
		"rx_queue_size":                          edgeMegaportResp.RxQueueSize,
		"state":                                  edgeMegaportResp.State,
		"dns_profile_name":                       edgeMegaportResp.DNSProfileName,
		"enable_single_ip_snat":                  edgeMegaportResp.EnableNat == "yes" && edgeMegaportResp.SnatMode == "primary",
		"enable_auto_advertise_lan_cidrs":        edgeMegaportResp.EnableAutoAdvertiseLanCidrs,
	}

	for key, value := range edgeMegaportFields {
		if err := d.Set(key, value); err != nil {
			log.Printf("[WARN] Failed to set %s: %v", key, err)
		}
	}

	if edgeMegaportResp.ManagementEgressIPPrefix == "" {
		_ = d.Set("management_egress_ip_prefix_list", nil)
	} else {
		_ = d.Set("management_egress_ip_prefix_list", strings.Split(edgeMegaportResp.ManagementEgressIPPrefix, ","))
	}

	if edgeMegaportResp.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: edgeMegaportResp.GwName})
		if err != nil {
			return diag.Errorf("could not get advanced config for Edge Megaport: %v", err)
		}

		err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs)
		if err != nil {
			return diag.Errorf("could not set approved_learned_cidrs into state: %v", err)
		}
	} else {
		_ = d.Set("approved_learned_cidrs", nil)
	}

	if len(edgeMegaportResp.AdvertisedCidrList) > 0 {
		_ = d.Set("included_advertised_spoke_routes", edgeMegaportResp.AdvertisedCidrList)
	}

	spokeBgpManualAdvertisedCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	if len(goaviatrix.Difference(spokeBgpManualAdvertisedCidrs, edgeMegaportResp.SpokeBgpManualAdvertisedCidrs)) != 0 ||
		len(goaviatrix.Difference(edgeMegaportResp.SpokeBgpManualAdvertisedCidrs, spokeBgpManualAdvertisedCidrs)) != 0 {
		_ = d.Set("spoke_bgp_manual_advertise_cidrs", edgeMegaportResp.SpokeBgpManualAdvertisedCidrs)
	} else {
		_ = d.Set("spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertisedCidrs)
	}

	if edgeMegaportResp.Latitude != 0 || edgeMegaportResp.Longitude != 0 {
		_ = d.Set("latitude", fmt.Sprintf("%.6f", edgeMegaportResp.Latitude))
		_ = d.Set("longitude", fmt.Sprintf("%.6f", edgeMegaportResp.Longitude))
	} else {
		_ = d.Set("latitude", "")
		_ = d.Set("longitude", "")
	}

	var interfaces []map[string]interface{}
	var vlan []map[string]interface{}
	// get user interface order
	userInterfaces := getList(d, "interfaces")

	userInterfaceOrder, err := getUserInterfaceOrder(userInterfaces)
	if err != nil {
		return diag.Errorf("failed to get user interface order: %s", err)
	}
	interfaceList := sortSpokeInterfacesByCustomOrder(edgeMegaportResp.InterfaceList, userInterfaceOrder)
	for _, interface0 := range interfaceList {
		interface1 := make(map[string]interface{})
		interface1["logical_ifname"] = interface0.LogicalInterfaceName
		if interface0.PublicIP != "" {
			interface1["wan_public_ip"] = interface0.PublicIP
		}
		if interface0.Dhcp {
			interface1["enable_dhcp"] = interface0.Dhcp
		}
		if interface0.IPAddr != "" {
			interface1["ip_address"] = interface0.IPAddr
		}
		if interface0.GatewayIP != "" {
			interface1["gateway_ip"] = interface0.GatewayIP
		}
		if interface0.DNSPrimary != "" {
			interface1["dns_server_ip"] = interface0.DNSPrimary
		}
		if interface0.DNSSecondary != "" {
			interface1["secondary_dns_server_ip"] = interface0.DNSSecondary
		}

		if strings.HasPrefix(interface0.LogicalInterfaceName, "lan") {
			interface1["enable_vrrp"] = interface0.VrrpState
			interface1["vrrp_virtual_ip"] = interface0.VirtualIP
		}

		if interface0.Tag != "" {
			interface1["tag"] = interface0.Tag
		}

		if strings.HasPrefix(interface0.LogicalInterfaceName, "lan") && interface0.SubInterfaces != nil {
			for _, vlan0 := range interface0.SubInterfaces {
				vlan1 := make(map[string]interface{})
				vlan1["parent_logical_interface_name"] = vlan0.ParentLogicalInterface
				vlan1["ip_address"] = vlan0.IPAddr
				vlan1["gateway_ip"] = vlan0.GatewayIP
				vlan1["peer_ip_address"] = vlan0.PeerIPAddr
				vlan1["peer_gateway_ip"] = vlan0.PeerGatewayIP
				vlan1["vrrp_virtual_ip"] = vlan0.VirtualIP
				vlan1["tag"] = vlan0.Tag
				vlanID, _ := strconv.Atoi(vlan0.VlanID)
				vlan1["vlan_id"] = vlanID

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

	// set interface mapping for megaport
	var interfaceMapping []map[string]interface{}
	for _, interfaceMap := range edgeMegaportResp.InterfaceMapping {
		interfaceMapping1 := make(map[string]interface{})
		interfaceMapping1["name"] = interfaceMap.Name
		interfaceMapping1["type"] = interfaceMap.Type
		interfaceMapping1["index"] = interfaceMap.Index
		interfaceMapping = append(interfaceMapping, interfaceMapping1)
	}
	if err = d.Set("interface_mapping", interfaceMapping); err != nil {
		return diag.Errorf("failed to set interface mapping: %s\n", err)
	}

	d.SetId(edgeMegaportResp.GwName)
	return nil
}

func resourceAviatrixEdgeMegaportUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// read configs
	edgeMegaport, err := marshalEdgeMegaportInput(d)
	if err != nil {
		return diag.Errorf("error reading Edge Megaport configuration: %s", err)
	}

	// checks before update
	if !edgeMegaport.EnableEdgeActiveStandby && edgeMegaport.EnableEdgeActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	if !edgeMegaport.EnableLearnedCidrsApproval && len(edgeMegaport.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeMegaport.PrependAsPath) != 0 {
		if edgeMegaport.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	if edgeMegaport.Latitude != "" && edgeMegaport.Longitude != "" {
		latitude, _ := strconv.ParseFloat(edgeMegaport.Latitude, 64)
		longitude, _ := strconv.ParseFloat(edgeMegaport.Longitude, 64)
		if latitude == 0 && longitude == 0 {
			return diag.Errorf("latitude and longitude must not be zero at the same time")
		}
	}

	d.Partial(true)

	// update configs
	// use following variables to reuse functions for transit, spoke and gateway
	gatewayForTransitFunctions := &goaviatrix.TransitVpc{
		GwName: edgeMegaport.GwName,
	}
	gatewayForSpokeFunctions := &goaviatrix.SpokeVpc{
		GwName: edgeMegaport.GwName,
	}
	gatewayForGatewayFunctions := &goaviatrix.Gateway{
		GwName: edgeMegaport.GwName,
	}
	gatewayForEaasFunctions := &goaviatrix.EdgeSpoke{
		GwName: edgeMegaport.GwName,
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeMegaport.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge Megaport update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeMegaport.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge Megaport update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeMegaport.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeMegaport.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge Megaport update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeMegaport.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge Megaport update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge Megaport update: %v", err)
			}
		}
	}

	if edgeMegaport.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeMegaport.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge Megaport update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeMegaport.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge Megaport update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeMegaport.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge Megaport update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge Megaport update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, edgeMegaport.BgpPollingTime)
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge Megaport update: %v", err)
		}
	}

	if d.HasChange("bgp_neighbor_status_polling_time") {
		err := client.SetBgpBfdPollingTimeSpoke(gatewayForSpokeFunctions, edgeMegaport.BgpBfdPollingTime)
		if err != nil {
			return diag.Errorf("could not set bgp neighbor status polling time during Edge Megaport update: %v", err)
		}
	}

	if d.HasChange("included_advertised_spoke_routes") {
		err := editAdvertisedSpokeRoutesWithRetry(client, gatewayForGatewayFunctions, d)
		if err != nil {
			return diag.Errorf("could not update included advertised spoke routes during Edge Megaport Gateway update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeMegaport.GwName, edgeMegaport.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge Megaport update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeMegaport.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeMegaport.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge Megaport update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeMegaport.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge Megaport update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeMegaport.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge Megaport update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge Megaport update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		gatewayForEaasFunctions.Latitude = edgeMegaport.Latitude
		gatewayForEaasFunctions.Longitude = edgeMegaport.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge Megaport update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeMegaport.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge Megaport update: %v", err)
		}
	}

	if d.HasChanges("management_egress_ip_prefix_list", "interfaces", "vlan",
		"enable_auto_advertise_lan_cidrs", "enable_edge_active_standby", "enable_edge_active_standby_preemptive") {
		err := client.UpdateEdgeMegaport(ctx, edgeMegaport)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list, WAN/LAN/VLAN interfaces, "+
				"DNS profile name, auto advertise LAN CIDRs, Edge active standby or Edge active standby preemptive "+
				"during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("interface_mapping") {
		// interface mapping is configured only during the creation of the gateway. Any updates to the interface mapping are not supported
		return diag.Errorf("interface mapping cannot be updated after the Edge Megaport is created")
	}

	if d.HasChange("enable_single_ip_snat") {
		gatewayForGatewayFunctions.GatewayName = edgeMegaport.GwName

		if edgeMegaport.EnableSingleIpSnat {
			err := client.EnableSNat(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("failed to enable single IP SNAT during Edge Megaport update: %s", err)
			}
		} else {
			err := client.DisableSNat(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("failed to disable single IP SNAT during Edge Megaport update: %s", err)
			}
		}
	}

	d.Partial(false)

	return resourceAviatrixEdgeMegaportRead(ctx, d, meta)
}

func resourceAviatrixEdgeMegaportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	accountName := getString(d, "account_name")
	gwName := getString(d, "gw_name")
	siteId := getString(d, "site_id")
	ztpFileDownloadPath := getString(d, "ztp_file_download_path")

	err := client.DeleteEdgeMegaport(ctx, accountName, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge Megaport: %v", err)
	}

	fileName := ztpFileDownloadPath + "/" + gwName + "-" + siteId + "-cloud-init.txt"
	err = os.Remove(fileName)
	if err != nil {
		log.Printf("[WARN] could not remove the ztp file: %v", err)
	}

	return nil
}
