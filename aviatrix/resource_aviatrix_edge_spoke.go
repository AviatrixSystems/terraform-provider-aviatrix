package aviatrix

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeSpoke() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeSpokeCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeSpokeRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeSpokeUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeSpokeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge as a Spoke name.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
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
				Description: "WAN interface IP / prefix.",
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
				Description: "LAN interface IP / prefix.",
			},
			"management_egress_ip_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Management egress gateway IP / prefix.",
			},
			"enable_over_private_network": {
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
				Description: "Management interface IP / prefix.",
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
				Description: "The location where the Edge as a CaaG ZTP file will be stored.",
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
			"enable_active_standby": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enables Active-Standby Mode.",
			},
			"enable_active_standby_preemptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.",
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
				Description: "Approved learned CIDRs for BGP Spoke Gateway. Available as of provider version R2.21+.",
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
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.FloatBetween(-90, 90),
				Description:  "The latitude of the Edge as a Spoke.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, n := d.GetChange("latitude")
					return n.(float64) == 0
				},
			},
			"longitude": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.FloatBetween(-180, 180),
				Description:  "The longitude of the Edge as a Spoke.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, n := d.GetChange("longitude")
					return n.(float64) == 0
				},
			},
			"wan_public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "WAN interface public IP.",
			},
		},
	}
}

func marshalEdgeSpokeInput(d *schema.ResourceData) *goaviatrix.EdgeSpoke {
	edgeSpoke := &goaviatrix.EdgeSpoke{
		GwName:                        d.Get("gw_name").(string),
		SiteId:                        d.Get("site_id").(string),
		ManagementInterfaceConfig:     d.Get("management_interface_config").(string),
		ManagementEgressIpPrefix:      d.Get("management_egress_ip_prefix").(string),
		EnableOverPrivateNetwork:      d.Get("enable_over_private_network").(bool),
		WanInterfaceIpPrefix:          d.Get("wan_interface_ip_prefix").(string),
		WanDefaultGatewayIp:           d.Get("wan_default_gateway_ip").(string),
		LanInterfaceIpPrefix:          d.Get("lan_interface_ip_prefix").(string),
		ManagementInterfaceIpPrefix:   d.Get("management_interface_ip_prefix").(string),
		ManagementDefaultGatewayIp:    d.Get("management_default_gateway_ip").(string),
		DnsServerIp:                   d.Get("dns_server_ip").(string),
		SecondaryDnsServerIp:          d.Get("secondary_dns_server_ip").(string),
		ZtpFileType:                   d.Get("ztp_file_type").(string),
		ZtpFileDownloadPath:           d.Get("ztp_file_download_path").(string),
		EnableActiveStandby:           d.Get("enable_active_standby").(bool),
		EnableActiveStandbyPreemptive: d.Get("enable_active_standby_preemptive").(bool),
		LocalAsNumber:                 d.Get("local_as_number").(string),
		PrependAsPath:                 getStringList(d, "prepend_as_path"),
		EnableLearnedCidrsApproval:    d.Get("enable_learned_cidrs_approval").(bool),
		ApprovedLearnedCidrs:          getStringSet(d, "approved_learned_cidrs"),
		SpokeBgpManualAdvertisedCidrs: getStringSet(d, "spoke_bgp_manual_advertise_cidrs"),
		EnablePreserveAsPath:          d.Get("enable_preserve_as_path").(bool),
		BgpPollingTime:                d.Get("bgp_polling_time").(int),
		BgpHoldTime:                   d.Get("bgp_hold_time").(int),
		EnableEdgeTransitiveRouting:   d.Get("enable_edge_transitive_routing").(bool),
		EnableJumboFrame:              d.Get("enable_jumbo_frame").(bool),
		Latitude:                      d.Get("latitude").(float64),
		Longitude:                     d.Get("longitude").(float64),
		WanPublicIp:                   d.Get("wan_public_ip").(string),
	}

	return edgeSpoke
}

func resourceAviatrixEdgeSpokeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeSpoke := marshalEdgeSpokeInput(d)

	// checks before creation
	if edgeSpoke.ManagementInterfaceConfig == "DHCP" && (edgeSpoke.ManagementInterfaceIpPrefix != "" || edgeSpoke.ManagementDefaultGatewayIp != "" ||
		edgeSpoke.DnsServerIp != "" || edgeSpoke.SecondaryDnsServerIp != "") {
		return diag.Errorf("'management_interface_ip', 'management_default_gateway_ip', 'dns_server_ip' and 'secondary_dns_server_ip' are only valid when 'management_interface_config' is Static")
	}

	if edgeSpoke.ManagementInterfaceConfig == "Static" && (edgeSpoke.ManagementInterfaceIpPrefix == "" || edgeSpoke.ManagementDefaultGatewayIp == "" ||
		edgeSpoke.DnsServerIp == "" || edgeSpoke.SecondaryDnsServerIp == "") {
		return diag.Errorf("'management_interface_ip', 'management_default_gateway_ip', 'dns_server_ip' and 'secondary_dns_server_ip' are required when 'management_interface_config' is Static")
	}

	if !edgeSpoke.EnableActiveStandby && edgeSpoke.EnableActiveStandbyPreemptive {
		return diag.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	} else if edgeSpoke.EnableActiveStandby && !edgeSpoke.EnableActiveStandbyPreemptive {
		edgeSpoke.ActiveStandby = "non-preemptive"
	} else if edgeSpoke.EnableActiveStandbyPreemptive {
		edgeSpoke.ActiveStandby = "preemptive"
	}

	if !edgeSpoke.EnableLearnedCidrsApproval && len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		if edgeSpoke.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
		}
	}

	// create
	d.SetId(edgeSpoke.GwName)
	flag := false
	defer resourceAviatrixEdgeSpokeReadIfRequired(ctx, d, meta, &flag)

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
			return diag.Errorf("could not set 'local_as_number' after Edge as a Spoke creation: %v", err)
		}
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		err := client.SetPrependASPath(gatewayForTransitFunctions, edgeSpoke.PrependAsPath)
		if err != nil {
			return diag.Errorf("could not set 'prepend_as_path' after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.EnableLearnedCidrsApproval {
		err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not enable learned CIDRs approval after Edge as a Spoke creation: %v", err)
		}
	}

	if len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeSpoke.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved CIDRs after Edge as a Spoke creation: %v", err)
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
			return diag.Errorf("could not enable spoke preserve as path after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.BgpPollingTime >= 10 && edgeSpoke.BgpPollingTime != defaultBgpPollingTime {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeSpoke.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.BgpHoldTime >= 12 && edgeSpoke.BgpHoldTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gatewayForSpokeFunctions.GwName, edgeSpoke.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change BGP Hold Time after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.EnableEdgeTransitiveRouting {
		err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.EnableJumboFrame {
		err := client.EnableJumboFrame(gatewayForGatewayFunctions)
		if err != nil {
			return diag.Errorf("could not disable jumbo frame after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.Latitude != 0 || edgeSpoke.Longitude != 0 {
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not enable Edge transitive routing after Edge as a Spoke creation: %v", err)
		}
	}

	if edgeSpoke.WanPublicIp != "" {
		err := client.UpdateEdgeSpokeIpConfigurations(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not config WAN public IP after Edge as a Spoke creation: %v", err)
		}
	}

	return resourceAviatrixEdgeSpokeReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeSpokeReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeSpokeRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeSpokeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("could not read Edge as a Spoke: %v", err)
	}

	d.Set("gw_name", edgeSpoke.GwName)
	d.Set("enable_over_private_network", edgeSpoke.EnableOverPrivateNetwork)
	d.Set("management_egress_ip_prefix", edgeSpoke.ManagementEgressIpPrefix)
	d.Set("wan_interface_ip_prefix", edgeSpoke.WanInterfaceIpPrefix)
	d.Set("wan_default_gateway_ip", edgeSpoke.WanDefaultGatewayIp)
	d.Set("lan_interface_ip_prefix", edgeSpoke.LanInterfaceIpPrefix)
	d.Set("management_default_gateway_ip", edgeSpoke.ManagementDefaultGatewayIp)
	d.Set("dns_server_ip", edgeSpoke.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeSpoke.SecondaryDnsServerIp)

	if edgeSpoke.Dhcp {
		d.Set("management_interface_config", "DHCP")
	} else {
		d.Set("management_interface_config", "Static")
		d.Set("management_interface_ip_prefix", edgeSpoke.ManagementInterfaceIpPrefix)
	}

	d.Set("local_as_number", edgeSpoke.LocalAsNumber)
	d.Set("prepend_as_path", edgeSpoke.PrependAsPath)
	d.Set("enable_active_standby", edgeSpoke.EnableActiveStandby)
	d.Set("enable_active_standby_preemptive", edgeSpoke.EnableActiveStandbyPreemptive)

	d.Set("enable_learned_cidrs_approval", edgeSpoke.EnableLearnedCidrsApproval)

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
	d.Set("latitude", edgeSpoke.Latitude)
	d.Set("longitude", edgeSpoke.Longitude)
	d.Set("wan_public_ip", edgeSpoke.WanPublicIp)

	d.SetId(edgeSpoke.GwName)
	return nil
}

func resourceAviatrixEdgeSpokeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	// read configs
	edgeSpoke := marshalEdgeSpokeInput(d)

	// checks before update
	if !edgeSpoke.EnableLearnedCidrsApproval && len(edgeSpoke.ApprovedLearnedCidrs) != 0 {
		return diag.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if len(edgeSpoke.PrependAsPath) != 0 {
		if edgeSpoke.LocalAsNumber == "" {
			return diag.Errorf("'prepend_as_path' must be empty if 'local_as_number' is not set")
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

	if d.HasChanges("management_egress_ip_prefix", "wan_interface_ip_prefix", "wan_default_gateway_ip", "lan_interface_ip_prefix", "wan_public_ip") {
		err := client.UpdateEdgeSpokeIpConfigurations(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update IP configurations during Edge as a Spoke update: %v", err)
		}
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(edgeSpoke.PrependAsPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gatewayForTransitFunctions, nil)
			if err != nil {
				return diag.Errorf("could not delete prepend_as_path during Edge as a Spoke update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			err := client.SetLocalASNumber(gatewayForTransitFunctions, edgeSpoke.LocalAsNumber)
			if err != nil {
				return diag.Errorf("could not set local_as_number during Edge as a Spoke update: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(edgeSpoke.PrependAsPath) > 0 {
			err := client.SetPrependASPath(gatewayForTransitFunctions, edgeSpoke.PrependAsPath)
			if err != nil {
				return diag.Errorf("could not set prepend_as_path during Edge as a Spoke update: %v", err)
			}
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		if edgeSpoke.EnableLearnedCidrsApproval {
			err := client.EnableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not enable learned cidrs approval during Edge as a Spoke update: %v", err)
			}
		} else {
			err := client.DisableTransitLearnedCidrsApproval(gatewayForTransitFunctions)
			if err != nil {
				return diag.Errorf("could not disable learned cidrs approval during Edge as a Spoke update: %v", err)
			}
		}
	}

	if edgeSpoke.EnableLearnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gatewayForTransitFunctions.ApprovedLearnedCidrs = edgeSpoke.ApprovedLearnedCidrs
		err := client.UpdateTransitPendingApprovedCidrs(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not update approved learned CIDRs during Edge as a Spoke update: %v", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		gatewayForTransitFunctions.BgpManualSpokeAdvertiseCidrs = strings.Join(edgeSpoke.SpokeBgpManualAdvertisedCidrs, ",")
		err := client.SetBgpManualSpokeAdvertisedNetworks(gatewayForTransitFunctions)
		if err != nil {
			return diag.Errorf("could not set spoke BGP manual advertised CIDRs during Edge as a Spoke update: %v", err)
		}
	}

	if d.HasChange("enable_preserve_as_path") {
		if edgeSpoke.EnablePreserveAsPath {
			err := client.EnableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not enable preserve as path during Edge as a Spoke update: %v", err)
			}
		} else {
			err := client.DisableSpokePreserveAsPath(gatewayForSpokeFunctions)
			if err != nil {
				return diag.Errorf("could not disable preserve as path during Edge as a Spoke update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		err := client.SetBgpPollingTimeSpoke(gatewayForSpokeFunctions, strconv.Itoa(edgeSpoke.BgpPollingTime))
		if err != nil {
			return diag.Errorf("could not set bgp polling time during Edge as a Spoke update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(edgeSpoke.GwName, edgeSpoke.BgpHoldTime)
		if err != nil {
			return diag.Errorf("could not change bgp hold time during Edge as a Spoke update: %v", err)
		}
	}

	if d.HasChange("enable_edge_transitive_routing") {
		if edgeSpoke.EnableEdgeTransitiveRouting {
			err := client.EnableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
			if err != nil {
				return diag.Errorf("could not enable transitive routing during Edge as a Spoke update: %v", err)
			}
		} else {
			err := client.DisableEdgeSpokeTransitiveRouting(ctx, edgeSpoke.GwName)
			if err != nil {
				return diag.Errorf("could not disable transitive routing during Edge as a Spoke update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if edgeSpoke.EnableJumboFrame {
			err := client.EnableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during Edge as a Spoke update: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gatewayForGatewayFunctions)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during Edge as a Spoke update: %v", err)
			}
		}
	}

	if d.HasChanges("latitude", "longitude") {
		err := client.UpdateEdgeSpokeGeoCoordinate(ctx, edgeSpoke)
		if err != nil {
			return diag.Errorf("could not update geo coordinate during Edge as a Spoke update: %v", err)
		}
	}

	d.Partial(false)

	return resourceAviatrixEdgeSpokeRead(ctx, d, meta)
}

func resourceAviatrixEdgeSpokeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	siteId := d.Get("site_id").(string)

	ztpFileDownloadPath := d.Get("ztp_file_download_path").(string)
	ztpFileType := d.Get("ztp_file_type").(string)

	err := client.DeleteEdgeSpoke(ctx, gwName)
	if err != nil {
		return diag.Errorf("could not delete Edge as a Spoke: %v", err)
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
