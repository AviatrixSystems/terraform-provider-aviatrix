package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixEdgeVmSelfmanaged() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeVmSelfmanagedCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeVmSelfmanagedRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeVmSelfmanagedUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeVmSelfmanagedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge VM Selfmanaged name.",
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
				Description: "State of Edge VM Selfmanaged.",
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
		},
		DeprecationMessage: "Since V3.1.2+, please use resource aviatrix_edge_gateway_selfmanaged instead. Resource " +
			"aviatrix_edge_vm_selfmanaged will be deprecated in the V3.2.0 release.",
	}
}

func marshalEdgeVmSelfmanagedInput(d *schema.ResourceData) *goaviatrix.EdgeSpoke {
	edgeSpoke := &goaviatrix.EdgeSpoke{
		GwName:                             getString(d, "gw_name"),
		SiteId:                             getString(d, "site_id"),
		ManagementEgressIpPrefix:           strings.Join(getStringSet(d, "management_egress_ip_prefix_list"), ","),
		EnableManagementOverPrivateNetwork: getBool(d, "enable_management_over_private_network"),
		DnsServerIp:                        getString(d, "dns_server_ip"),
		SecondaryDnsServerIp:               getString(d, "secondary_dns_server_ip"),
		ZtpFileType:                        getString(d, "ztp_file_type"),
		ZtpFileDownloadPath:                getString(d, "ztp_file_download_path"),
		EnableEdgeActiveStandby:            getBool(d, "enable_edge_active_standby"),
		EnableEdgeActiveStandbyPreemptive:  getBool(d, "enable_edge_active_standby_preemptive"),
		LocalAsNumber:                      getString(d, "local_as_number"),
		PrependAsPath:                      getStringList(d, "prepend_as_path"),
		EnableLearnedCidrsApproval:         getBool(d, "enable_learned_cidrs_approval"),
		ApprovedLearnedCidrs:               getStringSet(d, "approved_learned_cidrs"),
		SpokeBgpManualAdvertisedCidrs:      getStringSet(d, "spoke_bgp_manual_advertise_cidrs"),
		EnablePreserveAsPath:               getBool(d, "enable_preserve_as_path"),
		BgpPollingTime:                     getInt(d, "bgp_polling_time"),
		BgpHoldTime:                        getInt(d, "bgp_hold_time"),
		EnableEdgeTransitiveRouting:        getBool(d, "enable_edge_transitive_routing"),
		EnableJumboFrame:                   getBool(d, "enable_jumbo_frame"),
		Latitude:                           getString(d, "latitude"),
		Longitude:                          getString(d, "longitude"),
		RxQueueSize:                        getString(d, "rx_queue_size"),
	}

	interfaces := getSet(d, "interfaces").List()
	for _, if0 := range interfaces {
		if1 := mustMap(if0)

		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:       mustString(if1["name"]),
			Type:         mustString(if1["type"]),
			Dhcp:         mustBool(if1["enable_dhcp"]),
			PublicIp:     mustString(if1["wan_public_ip"]),
			IpAddr:       mustString(if1["ip_address"]),
			GatewayIp:    mustString(if1["gateway_ip"]),
			DNSPrimary:   mustString(if1["dns_server_ip"]),
			DNSSecondary: mustString(if1["secondary_dns_server_ip"]),
		}

		edgeSpoke.InterfaceList = append(edgeSpoke.InterfaceList, if2)
	}

	return edgeSpoke
}

func resourceAviatrixEdgeVmSelfmanagedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// read configs
	edgeSpoke := marshalEdgeVmSelfmanagedInput(d)

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
	defer resourceAviatrixEdgeVmSelfmanagedReadIfRequired(ctx, d, meta, &flag)

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
			return diag.Errorf("could not set 'local_as_number' after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeSpoke.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeSpoke.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge VM Selfmanaged creation: %v", err)
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
			return diag.Errorf("could not enable spoke preserve as path after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.BgpPollingTime >= 10 && edgeSpoke.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, edgeSpoke.BgpPollingTime)
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.BgpHoldTime >= 12 && edgeSpoke.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeSpoke.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.Latitude != "" || edgeSpoke.Longitude != "" {
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeSpoke.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge VM Selfmanaged creation: %v", err)
		}
	}

	if edgeSpoke.EnableEdgeActiveStandby || edgeSpoke.EnableEdgeActiveStandbyPreemptive {
		err := client.UpdateEdgeSpoke(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update Edge active standby or Edge active standby preemptive after Edge Gateway creation: %v", err)
		}
	}

	return resourceAviatrixEdgeVmSelfmanagedReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeVmSelfmanagedReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeVmSelfmanagedRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeVmSelfmanagedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// handle import
	if getString(d, "gw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		mustSet(d, "gw_name", id)
		d.SetId(id)
	}

	edgeSpoke, err := client.GetEdgeSpoke(ctx, getString(d, "gw_name"))
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge VM Selfmanaged: %v", err)
	}
	mustSet(d, "gw_name", edgeSpoke.GwName)
	mustSet(d, "site_id", edgeSpoke.SiteId)
	mustSet(d, "enable_management_over_private_network", edgeSpoke.EnableManagementOverPrivateNetwork)
	mustSet(d, "dns_server_ip", edgeSpoke.DnsServerIp)
	mustSet(d, "secondary_dns_server_ip", edgeSpoke.SecondaryDnsServerIp)
	mustSet(d, "local_as_number", edgeSpoke.LocalAsNumber)
	mustSet(d, "prepend_as_path", edgeSpoke.PrependAsPath)
	mustSet(d, "enable_edge_active_standby", edgeSpoke.EnableEdgeActiveStandby)
	mustSet(d, "enable_edge_active_standby_preemptive", edgeSpoke.EnableEdgeActiveStandbyPreemptive)
	mustSet(d, "enable_learned_cidrs_approval", edgeSpoke.EnableLearnedCidrsApproval)

	if edgeSpoke.ZtpFileType == "iso" || edgeSpoke.ZtpFileType == "cloud-init" {
		mustSet(d, "ztp_file_type", edgeSpoke.ZtpFileType)
	}

	if edgeSpoke.ZtpFileType == "cloud_init" {
		mustSet(d, "ztp_file_type", "cloud-init")
	}

	if edgeSpoke.ManagementEgressIpPrefix == "" {
		mustSet(d, "management_egress_ip_prefix_list", nil)
	} else {
		mustSet(d, "management_egress_ip_prefix_list", strings.Split(edgeSpoke.ManagementEgressIpPrefix, ","))
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
		mustSet(d, "approved_learned_cidrs", nil)
	}

	spokeBgpManualAdvertisedCidrs := getStringSet(d, "spoke_bgp_manual_advertise_cidrs")
	if len(goaviatrix.Difference(spokeBgpManualAdvertisedCidrs, edgeSpoke.SpokeBgpManualAdvertisedCidrs)) != 0 ||
		len(goaviatrix.Difference(edgeSpoke.SpokeBgpManualAdvertisedCidrs, spokeBgpManualAdvertisedCidrs)) != 0 {
		mustSet(d, "spoke_bgp_manual_advertise_cidrs", edgeSpoke.SpokeBgpManualAdvertisedCidrs)
	} else {
		mustSet(d, "spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertisedCidrs)
	}
	mustSet(d, "enable_preserve_as_path", edgeSpoke.EnablePreserveAsPath)
	mustSet(d, "bgp_polling_time", edgeSpoke.BgpPollingTime)
	mustSet(d, "bgp_hold_time", edgeSpoke.BgpHoldTime)
	mustSet(d, "enable_edge_transitive_routing", edgeSpoke.EnableEdgeTransitiveRouting)
	mustSet(d, "enable_jumbo_frame", edgeSpoke.EnableJumboFrame)
	if edgeSpoke.Latitude != 0 || edgeSpoke.Longitude != 0 {
		mustSet(d, "latitude", fmt.Sprintf("%.6f", edgeSpoke.Latitude))
		mustSet(d, "longitude", fmt.Sprintf("%.6f", edgeSpoke.Longitude))
	} else {
		mustSet(d, "latitude", "")
		mustSet(d, "longitude", "")
	}
	mustSet(d, "rx_queue_size", edgeSpoke.RxQueueSize)
	mustSet(d, "state", edgeSpoke.State)

	var interfaces []map[string]interface{}

	for _, if0 := range edgeSpoke.InterfaceList {
		if1 := make(map[string]interface{})
		if1["name"] = if0.IfName
		if1["type"] = if0.Type
		if1["enable_dhcp"] = if0.Dhcp
		if1["wan_public_ip"] = if0.PublicIp
		if1["ip_address"] = if0.IpAddr
		if1["gateway_ip"] = if0.GatewayIp
		if1["dns_server_ip"] = if0.DNSPrimary
		if1["secondary_dns_server_ip"] = if0.DNSSecondary

		interfaces = append(interfaces, if1)
	}

	if err = d.Set("interfaces", interfaces); err != nil {
		return diag.Errorf("failed to set interfaces: %s\n", err)
	}

	d.SetId(edgeSpoke.GwName)
	return nil
}

func resourceAviatrixEdgeVmSelfmanagedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// read configs
	edgeSpoke := marshalEdgeVmSelfmanagedInput(d)

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
				return diag.Errorf("could not delete prepend_as_path during Edge VM Selfmanaged update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeSpoke.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge VM Selfmanaged update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeSpoke.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeSpoke.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge VM Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeSpoke.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge VM Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge VM Selfmanaged update: %v", err)
			}
		}
	}

	if edgeSpoke.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeSpoke.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge VM Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeSpoke.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge VM Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeSpoke.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge VM Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge VM Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, edgeSpoke.BgpPollingTime)
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge VM Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeSpoke.GwName, edgeSpoke.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge VM Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeSpoke.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge VM Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge VM Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeSpoke.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge VM Selfmanaged update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge VM Selfmanaged update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge VM Selfmanaged update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeSpoke.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge VM Selfmanaged update: %v", err)
		}
	}

	if d.HasChanges("management_egress_ip_prefix_list", "interfaces", "enable_edge_active_standby",
		"enable_edge_active_standby_preemptive") {
		err := client.UpdateEdgeSpoke(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update management egress ip prefix list, WAN/LAN/MANAGEMENT interfaces, "+
				"Edge active standby or Edge active standby preemptive during Edge as a Spoke update: %v", err)
		}
	}

	d.Partial(false)

	return resourceAviatrixEdgeVmSelfmanagedRead(ctx, d, meta)
}

func resourceAviatrixEdgeVmSelfmanagedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	gwName := getString(d, "gw_name")
	siteId := getString(d, "site_id")

	ztpFileDownloadPath := getString(d, "ztp_file_download_path")
	ztpFileType := getString(d, "ztp_file_type")

	err := client.DeleteEdgeSpoke(ctx, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge VM Selfmanaged: %v", err)
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
