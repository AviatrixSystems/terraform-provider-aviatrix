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

func resourceAviatrixEdgeCSP() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeCSPCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeCSPRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeCSPUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeCSPDelete,
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
			"management_interface_config": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Management interface configuration. Valid values: 'DHCP' and 'Static'.",
				ValidateFunc: validation.StringInSlice([]string{"DHCP", "Static"}, false),
			},
			"wan_interface_ip_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "WAN interface IP/prefix.",
			},
			"wan_default_gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "WAN default gateway IP.",
				ValidateFunc: validation.IsIPAddress,
			},
			"lan_interface_ip_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "LAN interface IP/prefix.",
			},
			"management_egress_ip_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Management egress gateway IP/prefix.",
			},
			"enable_management_over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable management over private network.",
			},
			"management_interface_ip_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Management interface IP/prefix.",
			},
			"management_default_gateway_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Management default gateway IP.",
				ValidateFunc: validation.IsIPAddress,
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
				ForceNew:    true,
				Default:     false,
				Description: "Enables Edge Active-Standby Mode.",
			},
			"enable_edge_active_standby_preemptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
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
			"wan_public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "WAN interface public IP.",
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
			"wan_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "eth0",
				ForceNew:    true,
				Description: "WAN interface name.",
			},
			"lan_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "eth1",
				ForceNew:    true,
				Description: "LAN interface name.",
			},
			"management_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "eth2",
				ForceNew:    true,
				Description: "Management interface name.",
			},
		},
	}
}

func marshalEdgeCSPInput(d *schema.ResourceData) *goaviatrix.EdgeCSP {
	edgeCSP := &goaviatrix.EdgeCSP{
		AccountName:                        d.Get("account_name").(string),
		GwName:                             d.Get("gw_name").(string),
		SiteId:                             d.Get("site_id").(string),
		ProjectUuid:                        d.Get("project_uuid").(string),
		ComputeNodeUuid:                    d.Get("compute_node_uuid").(string),
		TemplateUuid:                       d.Get("template_uuid").(string),
		ManagementInterfaceConfig:          d.Get("management_interface_config").(string),
		ManagementEgressIpPrefix:           d.Get("management_egress_ip_prefix").(string),
		EnableManagementOverPrivateNetwork: d.Get("enable_management_over_private_network").(bool),
		WanInterfaceIpPrefix:               d.Get("wan_interface_ip_prefix").(string),
		WanDefaultGatewayIp:                d.Get("wan_default_gateway_ip").(string),
		LanInterfaceIpPrefix:               d.Get("lan_interface_ip_prefix").(string),
		ManagementInterfaceIpPrefix:        d.Get("management_interface_ip_prefix").(string),
		ManagementDefaultGatewayIp:         d.Get("management_default_gateway_ip").(string),
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
		WanPublicIp:                        d.Get("wan_public_ip").(string),
		RxQueueSize:                        d.Get("rx_queue_size").(string),
		WanInterface:                       d.Get("wan_interface_name").(string),
		LanInterface:                       d.Get("lan_interface_name").(string),
		MgmtInterface:                      d.Get("management_interface_name").(string),
	}

	return edgeCSP
}

func resourceAviatrixEdgeCSPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeCSP := marshalEdgeCSPInput(d)

	// checks before creation
	if edgeCSP.ManagementInterfaceConfig == "DHCP" && (edgeCSP.ManagementInterfaceIpPrefix != "" || edgeCSP.ManagementDefaultGatewayIp != "" ||
		edgeCSP.DnsServerIp != "" || edgeCSP.SecondaryDnsServerIp != "") {
		return diag.Errorf("'management_interface_ip', 'management_default_gateway_ip', 'dns_server_ip' and 'secondary_dns_server_ip' are only valid when 'management_interface_config' is Static")
	}

	if edgeCSP.ManagementInterfaceConfig == "Static" && (edgeCSP.ManagementInterfaceIpPrefix == "" || edgeCSP.ManagementDefaultGatewayIp == "" ||
		edgeCSP.DnsServerIp == "" || edgeCSP.SecondaryDnsServerIp == "") {
		return diag.Errorf("'management_interface_ip', 'management_default_gateway_ip', 'dns_server_ip' and 'secondary_dns_server_ip' are required when 'management_interface_config' is Static")
	}

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
	defer resourceAviatrixEdgeCSPReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeCSP(ctx, edgeCSP); err != nil {
		return diag.Errorf("could not create Edge CSP: %v", err)
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
			return diag.Errorf("could not set 'local_as_number' after Edge CSP creation: %v", err)
		}
	}

	if len(edgeCSP.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeCSP.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge CSP creation: %v", err)
		}
	}

	if len(edgeCSP.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeCSP.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge CSP creation: %v", err)
		}
	}

	if len(edgeCSP.SpokeBgpManualAdvertisedCidrs) != 0 {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeCSP.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.EnablePreserveAsPath {
		err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
		if err != nil {
			return diag.Errorf("could not enable spoke preserve as path after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.BgpPollingTime >= 10 && edgeCSP.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeCSP.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.BgpHoldTime >= 12 && edgeCSP.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeCSP.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeCSP.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.Latitude != "" || edgeCSP.Longitude != "" {
		gatewayForEaasFunctions.Latitude = edgeCSP.Latitude
		gatewayForEaasFunctions.Longitude = edgeCSP.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.WanPublicIp != "" {
		gatewayForEaasFunctions.WanInterfaceIpPrefix = edgeCSP.WanInterfaceIpPrefix
		gatewayForEaasFunctions.WanDefaultGatewayIp = edgeCSP.WanDefaultGatewayIp
		gatewayForEaasFunctions.LanInterfaceIpPrefix = edgeCSP.LanInterfaceIpPrefix
		gatewayForEaasFunctions.ManagementEgressIpPrefix = edgeCSP.ManagementEgressIpPrefix
		gatewayForEaasFunctions.WanPublicIp = edgeCSP.WanPublicIp
		err := client.UpdateEdgeSpokeIpConfigurations(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not config WAN public IP after Edge CSP creation: %v", err)
		}
	}

	if edgeCSP.RxQueueSize != "" {
		gatewayForGatewayFunctions.RxQueueSize = edgeCSP.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not set rx queue size after Edge CSP creation: %v", err)
		}
	}

	return resourceAviatrixEdgeCSPReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeCSPReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeCSPRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeCSPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("could not read Edge CSP: %v", err)
	}

	d.Set("account_name", edgeCSPResp.AccountName)
	d.Set("gw_name", edgeCSPResp.GwName)
	d.Set("site_id", edgeCSPResp.SiteId)
	d.Set("project_uuid", edgeCSPResp.ProjectUuid)
	d.Set("compute_node_uuid", edgeCSPResp.ComputeNodeUuid)
	d.Set("template_uuid", edgeCSPResp.TemplateUuid)
	d.Set("enable_management_over_private_network", edgeCSPResp.EnableManagementOverPrivateNetwork)
	d.Set("management_egress_ip_prefix", edgeCSPResp.ManagementEgressIpPrefix)
	d.Set("wan_interface_ip_prefix", edgeCSPResp.WanInterfaceIpPrefix)
	d.Set("wan_default_gateway_ip", edgeCSPResp.WanDefaultGatewayIp)
	d.Set("lan_interface_ip_prefix", edgeCSPResp.LanInterfaceIpPrefix)
	d.Set("management_default_gateway_ip", edgeCSPResp.ManagementDefaultGatewayIp)
	d.Set("dns_server_ip", edgeCSPResp.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeCSPResp.SecondaryDnsServerIp)

	if edgeCSPResp.Dhcp {
		d.Set("management_interface_config", "DHCP")
	} else {
		d.Set("management_interface_config", "Static")
		d.Set("management_interface_ip_prefix", edgeCSPResp.ManagementInterfaceIpPrefix)
	}

	d.Set("local_as_number", edgeCSPResp.LocalAsNumber)
	d.Set("prepend_as_path", edgeCSPResp.PrependAsPath)
	d.Set("enable_edge_active_standby", edgeCSPResp.EnableEdgeActiveStandby)
	d.Set("enable_edge_active_standby_preemptive", edgeCSPResp.EnableEdgeActiveStandbyPreemptive)

	d.Set("enable_learned_cidrs_approval", edgeCSPResp.EnableLearnedCidrsApproval)

	if edgeCSPResp.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: edgeCSPResp.GwName})
		if err != nil {
			return diag.Errorf("could not get advanced config for Edge CSP: %v", err)
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
	if edgeCSPResp.LatitudeReturn != 0 || edgeCSPResp.LongitudeReturn != 0 {
		d.Set("latitude", fmt.Sprintf("%.6f", edgeCSPResp.LatitudeReturn))
		d.Set("longitude", fmt.Sprintf("%.6f", edgeCSPResp.LongitudeReturn))
	} else {
		d.Set("latitude", "")
		d.Set("longitude", "")
	}
	d.Set("wan_public_ip", edgeCSPResp.WanPublicIp)
	d.Set("rx_queue_size", edgeCSPResp.RxQueueSize)
	d.Set("state", edgeCSPResp.State)
	d.Set("wan_interface_name", edgeCSPResp.WanInterface)
	d.Set("lan_interface_name", edgeCSPResp.LanInterface)
	d.Set("management_interface_name", edgeCSPResp.MgmtInterface)

	d.SetId(edgeCSPResp.GwName)
	return nil
}

func resourceAviatrixEdgeCSPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeCSP := marshalEdgeCSPInput(d)

	// checks before update
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

	if d.HasChanges("management_egress_ip_prefix", "wan_interface_ip_prefix", "wan_default_gateway_ip", "lan_interface_ip_prefix", "wan_public_ip") {
		gatewayForEaasFunctions.WanInterfaceIpPrefix = edgeCSP.WanInterfaceIpPrefix
		gatewayForEaasFunctions.WanDefaultGatewayIp = edgeCSP.WanDefaultGatewayIp
		gatewayForEaasFunctions.LanInterfaceIpPrefix = edgeCSP.LanInterfaceIpPrefix
		gatewayForEaasFunctions.ManagementEgressIpPrefix = edgeCSP.ManagementEgressIpPrefix
		gatewayForEaasFunctions.WanPublicIp = edgeCSP.WanPublicIp

		err := client.UpdateEdgeSpokeIpConfigurations(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update IP configurations during Edge CSP update: %v", err)
		}
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeCSP.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge CSP update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeCSP.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge CSP update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeCSP.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeCSP.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge CSP update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeCSP.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge CSP update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge CSP update: %v", err)
			}
		}
	}

	if edgeCSP.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeCSP.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeCSP.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeCSP.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge CSP update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge CSP update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeCSP.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeCSP.GwName, edgeCSP.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeCSP.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeCSP.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge CSP update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeCSP.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge CSP update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeCSP.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge CSP update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge CSP update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		gatewayForEaasFunctions.Latitude = edgeCSP.Latitude
		gatewayForEaasFunctions.Longitude = edgeCSP.Longitude
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, gatewayForEaasFunctions)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge CSP update: %v", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		gatewayForGatewayFunctions.RxQueueSize = edgeCSP.RxQueueSize
		err := client.SetRxQueueSize(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not update rx queue size during Edge CSP update: %v", err)
		}
	}

	d.Partial(false)

	return resourceAviatrixEdgeCSPRead(ctx, d, meta)
}

func resourceAviatrixEdgeCSPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	gwName := d.Get("gw_name").(string)

	err := client.DeleteEdgeCSP(ctx, accountName, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge CSP: %v", err)
	}

	return nil
}
