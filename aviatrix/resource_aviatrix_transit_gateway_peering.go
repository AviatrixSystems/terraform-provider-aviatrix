package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixTransitGatewayPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitGatewayPeeringCreate,
		Read:   resourceAviatrixTransitGatewayPeeringRead,
		Update: resourceAviatrixTransitGatewayPeeringUpdate,
		Delete: resourceAviatrixTransitGatewayPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"transit_gateway_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The first transit gateway name to make a peer pair.",
			},
			"transit_gateway_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The second transit gateway name to make a peer pair.",
			},
			"gateway1_excluded_network_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded network CIDRs for the first transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway1_excluded_tgw_connections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded TGW connections for the first transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway2_excluded_network_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded network CIDRs for the second transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway2_excluded_tgw_connections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded TGW connections for the second transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"prepend_as_path1": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit_gateway_name1.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"prepend_as_path2": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit_gateway_name2.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"enable_peering_over_private_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "(Optional) Advanced option. Enable peering over private network. Only appears and applies to " +
					"when the two Multi-cloud Transit Gateways are each launched in Insane Mode and in a different cloud type. " +
					"Conflicts with `enable_insane_mode_encryption_over_internet` and `tunnel_count`. " +
					"Type: Boolean. Default: false. Available in provider version R2.17.1+",
			},
			"enable_single_tunnel_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("enable_peering_over_private_network").(bool)
				},
				Description: "(Optional) Advanced option. Enable peering with Single-Tunnel mode. Only appears and applies " +
					"to when the two Multi-cloud Transit Gateways are each launched in Insane Mode and in a different cloud type. " +
					"Required with `enable_peering_over_private_network`. Conflicts with `enable_insane_mode_encryption_over_internet` and `tunnel_count`. " +
					"Type: Boolean. Default: false. Available as of provider version R2.18+.",
			},
			"enable_insane_mode_encryption_over_internet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "(Optional) Advanced option. Enable Insane Mode Encryption over Internet. Transit gateways must be in Insane Mode. " +
					"Currently, only inter-cloud connections between AWS and Azure are supported. Required with valid `tunnel_count`. " +
					"Conflicts with `enable_peering_over_private_network` and `enable_single_tunnel_mode`. Type: Boolean. Default: false. " +
					"Available as of provider version R2.19+.",
			},
			"tunnel_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(2, 20),
				Description: "(Optional) Advanced option. Number of public tunnels. Required with `enable_insane_mode_encryption_over_internet`. " +
					"Conflicts with `enable_peering_over_private_network` and `enable_single_tunnel_mode`. Type: Integer. Valid Range: 2-20. " +
					"Available as of provider version R2.19+.",
			},
			"enable_max_performance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
				Description: "Indicates whether the maximum amount of HPE tunnels will be created. " +
					"Only valid when the two transit gateways are each launched in Insane Mode and in the same cloud type. " +
					"Available as of provider version R2.22.2+.",
			},
			"over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The underlay connects over the private network for peering with Edge Transit",
			},
			"jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable jumbo frame for over private peering with Edge Transit",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable HPE mode for peering with Edge Transit",
			},
			"gateway1_logical_ifnames": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Gateway 1 logical interface names for edge gateways where the peering originates",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway2_logical_ifnames": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway 2 logical interface names for edge gateways where the peering terminates",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAviatrixTransitGatewayPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	flag := false
	transit_gateway_name1, ok := d.Get("transit_gateway_name1").(string)
	if !ok {
		return fmt.Errorf("transit_gateway_name1 is required")
	}
	transit_gateway_name2, ok := d.Get("transit_gateway_name2").(string)
	if !ok {
		return fmt.Errorf("transit_gateway_name2 is required")
	}

	gateway1 := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      transit_gateway_name1,
	}
	gateway1Details, err := client.GetGateway(gateway1)
	if err != nil {
		return fmt.Errorf("failed to get gateway: %s", err)
	}
	gateway1CloudType := gateway1Details.CloudType
	gateway2CloudType := gateway1Details.CloudType
	log.Printf("[INFO] Gateway1 details: %#v", gateway1Details)
	gateway2 := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      transit_gateway_name2,
	}
	gateway2Details, err := client.GetGateway(gateway2)
	if err != nil {
		return fmt.Errorf("failed to get gateway: %s", err)
	}
	log.Printf("[INFO] Gateway2 details: %#v", gateway2Details)

	if goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EdgeRelatedCloudTypes) || goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		edgeTransitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: transit_gateway_name1,
			TransitGatewayName2: transit_gateway_name2,
		}

		edgeTransitGatewayPeering.EnableOverPrivateNetwork, ok = d.Get("over_private_network").(bool)
		if !ok {
			return fmt.Errorf("over_private_network is required for edge gateway peering")
		}
		edgeTransitGatewayPeering.EnableJumboFrame, ok = d.Get("jumbo_frame").(bool)
		if !ok {
			return fmt.Errorf("jumbo_frame is required for edge gateway peering")
		}
		edgeTransitGatewayPeering.EnableInsaneMode, ok = d.Get("insane_mode").(bool)
		if !ok {
			return fmt.Errorf("insane_mode is required for edge gateway peering")
		}

		// get the src wan interface names from gateway1 logical interface names
		if goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EdgeRelatedCloudTypes) {
			// process the gateway logical interface names based on the cloud type
			logicalIfNames, ok := d.Get("gateway1_logical_ifnames").([]interface{})
			if !ok {
				return fmt.Errorf("gateway1_logical_ifnames is required for edge gateway peering")
			}
			reversedInterfaceNames := ReverseIfnameTranslation(gateway1Details.IfNamesTranslation)
			if goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGENEO) {
				srcWanInterfacesStr, err := SetWanInterfaces(logicalIfNames, reversedInterfaceNames)
				if err != nil {
					return fmt.Errorf("failed to set src wan interfaces to create edge peering: %s", err)
				}
				edgeTransitGatewayPeering.SrcWanInterfaces = srcWanInterfacesStr

			} else if goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EDGEMEGAPORT) {
				if d.Get("gateway1_logical_ifnames").([]interface{}) != nil {
					edgeTransitGatewayPeering.Gateway1LogicalIfNames = getStringList(d, "gateway1_logical_ifnames")
				}
			}
		}
		// get the dst wan interface names from gateway2 logical interface names
		if goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EdgeRelatedCloudTypes) {
			// process the gateway logical interface names based on the cloud type
			logicalIfNames, ok := d.Get("gateway2_logical_ifnames").([]interface{})
			if !ok {
				return fmt.Errorf("gateway2_logical_ifnames is required for edge gateway peering")
			}
			reversedInterfaceNames := ReverseIfnameTranslation(gateway2Details.IfNamesTranslation)
			if goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGENEO) {
				dstWanInterfacesStr, err := SetWanInterfaces(logicalIfNames, reversedInterfaceNames)
				if err != nil {
					return fmt.Errorf("failed to set dst wan interfaces to create edge peering: %s", err)
				}
				edgeTransitGatewayPeering.DstWanInterfaces = dstWanInterfacesStr

			} else if goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EDGEMEGAPORT) {
				if d.Get("gateway2_logical_ifnames").([]interface{}) != nil {
					edgeTransitGatewayPeering.Gateway2LogicalIfNames = getStringList(d, "gateway2_logical_ifnames")
				}
			}
		}
		d.SetId(edgeTransitGatewayPeering.TransitGatewayName1 + "~" + edgeTransitGatewayPeering.TransitGatewayName2)
		defer resourceAviatrixTransitGatewayPeeringReadIfRequired(d, meta, &flag)

		err := client.CreateTransitGatewayPeering(edgeTransitGatewayPeering)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Transit Gateway peering: %s", err)
		}
	} else {
		var gw1Cidrs []string
		for _, cidr := range d.Get("gateway1_excluded_network_cidrs").([]interface{}) {
			gw1Cidrs = append(gw1Cidrs, cidr.(string))
		}
		var gw2Cidrs []string
		for _, cidr := range d.Get("gateway2_excluded_network_cidrs").([]interface{}) {
			gw2Cidrs = append(gw2Cidrs, cidr.(string))
		}

		var gw1Tgws []string
		for _, tgw := range d.Get("gateway1_excluded_tgw_connections").([]interface{}) {
			gw1Tgws = append(gw1Tgws, tgw.(string))
		}
		var gw2Tgws []string
		for _, tgw := range d.Get("gateway2_excluded_tgw_connections").([]interface{}) {
			gw2Tgws = append(gw2Tgws, tgw.(string))
		}

		transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1:            transit_gateway_name1,
			TransitGatewayName2:            transit_gateway_name2,
			Gateway1ExcludedCIDRs:          strings.Join(gw1Cidrs, ","),
			Gateway2ExcludedCIDRs:          strings.Join(gw2Cidrs, ","),
			Gateway1ExcludedTGWConnections: strings.Join(gw1Tgws, ","),
			Gateway2ExcludedTGWConnections: strings.Join(gw2Tgws, ","),
			InsaneModeOverInternet:         d.Get("enable_insane_mode_encryption_over_internet").(bool),
			NoMaxPerformance:               !d.Get("enable_max_performance").(bool),
		}
		if d.Get("enable_peering_over_private_network").(bool) {
			transitGatewayPeering.PrivateIPPeering = "yes"
		} else {
			transitGatewayPeering.PrivateIPPeering = "no"
		}

		if transitGatewayPeering.PrivateIPPeering == "yes" && transitGatewayPeering.InsaneModeOverInternet {
			return fmt.Errorf("enable_peering_over_private_network conflicts with enable_insane_mode_encryption_over_internet")
		}

		if d.Get("enable_single_tunnel_mode").(bool) {
			if transitGatewayPeering.PrivateIPPeering == "yes" {
				transitGatewayPeering.SingleTunnel = "yes"
			} else {
				return fmt.Errorf("enable_single_tunnel_mode is only valid when enable_peering_over_private_network is set to true")
			}
		}

		tunnelCount := d.Get("tunnel_count").(int)
		if tunnelCount != 0 {
			if transitGatewayPeering.InsaneModeOverInternet {
				transitGatewayPeering.TunnelCount = tunnelCount
			} else {
				return fmt.Errorf("tunnel_count is only valid when enable_insane_mode_encryption_over_internet is set to true")
			}
		} else {
			if transitGatewayPeering.InsaneModeOverInternet {
				return fmt.Errorf("enable_insane_mode_encryption_over_internet being set to true requires valid tunnel_count")
			}
		}

		log.Printf("[INFO] Creating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

		d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
		flag := false
		defer resourceAviatrixTransitGatewayPeeringReadIfRequired(d, meta, &flag)

		err := client.CreateTransitGatewayPeering(transitGatewayPeering)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Transit Gateway peering: %s", err)
		}

		if _, ok := d.GetOk("prepend_as_path1"); ok {
			var prependASPath []string
			for _, v := range d.Get("prepend_as_path1").([]interface{}) {
				prependASPath = append(prependASPath, v.(string))
			}
			transGwPeering := &goaviatrix.TransitGatewayPeering{
				TransitGatewayName1: d.Get("transit_gateway_name1").(string),
				TransitGatewayName2: d.Get("transit_gateway_name2").(string),
			}

			err = client.EditTransitConnectionASPathPrepend(transGwPeering, prependASPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path1: %v", err)
			}
		}

		if _, ok := d.GetOk("prepend_as_path2"); ok {
			var prependASPath []string
			for _, v := range d.Get("prepend_as_path2").([]interface{}) {
				prependASPath = append(prependASPath, v.(string))
			}
			transGwPeering := &goaviatrix.TransitGatewayPeering{
				TransitGatewayName1: d.Get("transit_gateway_name2").(string),
				TransitGatewayName2: d.Get("transit_gateway_name1").(string),
			}

			err = client.EditTransitConnectionASPathPrepend(transGwPeering, prependASPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path2: %v", err)
			}
		}
	}

	return resourceAviatrixTransitGatewayPeeringReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitGatewayPeeringReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitGatewayPeeringRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitGatewayPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGwName1 := d.Get("transit_gateway_name1").(string)
	transitGwName2 := d.Get("transit_gateway_name2").(string)

	if transitGwName1 == "" || transitGwName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return fmt.Errorf("invalid import id expected transit_gateway_name1~transit_gateway_name2")
		}
		d.Set("transit_gateway_name1", parts[0])
		d.Set("transit_gateway_name2", parts[1])
		d.SetId(id)
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}

	transitGatewayPeering, err := client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get transit peering details: %v", err)
	}

	if len(transitGatewayPeering.Gateway1LogicalIfNames) > 0 || len(transitGatewayPeering.Gateway2LogicalIfNames) > 0 {
		d.Set("over_private_network", transitGatewayPeering.EnableOverPrivateNetwork)
		d.Set("jumbo_frame", transitGatewayPeering.EnableJumboFrame)
		d.Set("insane_mode", transitGatewayPeering.EnableInsaneMode)
		d.Set("gateway1_logical_ifnames", transitGatewayPeering.Gateway1LogicalIfNames)
		d.Set("gateway2_logical_ifnames", transitGatewayPeering.Gateway2LogicalIfNames)
	} else {
		gw1CidrsFromConfig := getStringList(d, "gateway1_excluded_network_cidrs")
		err = setConfigValueIfEquivalent(d, "gateway1_excluded_network_cidrs", gw1CidrsFromConfig, transitGatewayPeering.Gateway1ExcludedCIDRsSlice)
		if err != nil {
			return fmt.Errorf("could not write gateway1_excluded_network_cidrs to state: %v", err)
		}
		gw2CidrsFromConfig := getStringList(d, "gateway2_excluded_network_cidrs")
		err = setConfigValueIfEquivalent(d, "gateway2_excluded_network_cidrs", gw2CidrsFromConfig, transitGatewayPeering.Gateway2ExcludedCIDRsSlice)
		if err != nil {
			return fmt.Errorf("could not write gateway2_excluded_network_cidrs to state: %v", err)
		}
		gw1TgwsFromConfig := getStringList(d, "gateway1_excluded_tgw_connections")
		err = setConfigValueIfEquivalent(d, "gateway1_excluded_tgw_connections", gw1TgwsFromConfig, transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice)
		if err != nil {
			return fmt.Errorf("could not write gateway1_excluded_tgw_connections to state: %v", err)
		}
		gw2TgwsFromConfig := getStringList(d, "gateway2_excluded_tgw_connections")
		err = setConfigValueIfEquivalent(d, "gateway2_excluded_tgw_connections", gw2TgwsFromConfig, transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice)
		if err != nil {
			return fmt.Errorf("could not write gateway2_excluded_tgw_connections to state: %v", err)
		}

		if transitGatewayPeering.PrependAsPath1 != "" {
			var prependAsPath []string
			for _, str := range strings.Split(transitGatewayPeering.PrependAsPath1, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("prepend_as_path1", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path1: %v", err)
			}
		}
		if transitGatewayPeering.PrependAsPath2 != "" {
			var prependAsPath []string
			for _, str := range strings.Split(transitGatewayPeering.PrependAsPath2, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("prepend_as_path2", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path2: %v", err)
			}
		}
		d.Set("enable_peering_over_private_network", transitGatewayPeering.PrivateIPPeering == "yes")
		if transitGatewayPeering.PrivateIPPeering == "yes" {
			d.Set("enable_single_tunnel_mode", transitGatewayPeering.SingleTunnel == "yes")
		} else {
			d.Set("enable_single_tunnel_mode", false)
		}
		d.Set("enable_insane_mode_encryption_over_internet", transitGatewayPeering.InsaneModeOverInternet)
		if transitGatewayPeering.InsaneModeOverInternet {
			d.Set("tunnel_count", transitGatewayPeering.TunnelCount)
		}

		d.Set("enable_max_performance", !transitGatewayPeering.NoMaxPerformance)
	}

	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return nil
}

func resourceAviatrixTransitGatewayPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}
	if d.HasChange("gateway1_excluded_network_cidrs") || d.HasChange("gateway2_excluded_network_cidrs") ||
		d.HasChange("gateway1_excluded_tgw_connections") || d.HasChange("gateway2_excluded_tgw_connections") {
		var gw1Cidrs []string
		for _, cidr := range d.Get("gateway1_excluded_network_cidrs").([]interface{}) {
			gw1Cidrs = append(gw1Cidrs, cidr.(string))
		}
		var gw2Cidrs []string
		for _, cidr := range d.Get("gateway2_excluded_network_cidrs").([]interface{}) {
			gw2Cidrs = append(gw2Cidrs, cidr.(string))
		}
		var gw1Tgws []string
		for _, tgw := range d.Get("gateway1_excluded_tgw_connections").([]interface{}) {
			gw1Tgws = append(gw1Tgws, tgw.(string))
		}
		var gw2Tgws []string
		for _, tgw := range d.Get("gateway2_excluded_tgw_connections").([]interface{}) {
			gw2Tgws = append(gw2Tgws, tgw.(string))
		}

		transitGatewayPeering.Gateway1ExcludedCIDRs = strings.Join(gw1Cidrs, ",")
		transitGatewayPeering.Gateway2ExcludedCIDRs = strings.Join(gw2Cidrs, ",")
		transitGatewayPeering.Gateway1ExcludedTGWConnections = strings.Join(gw1Tgws, ",")
		transitGatewayPeering.Gateway2ExcludedTGWConnections = strings.Join(gw2Tgws, ",")

		log.Printf("[INFO] Updating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)
		err := client.UpdateTransitGatewayPeering(transitGatewayPeering)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit Gateway peering: %s", err)
		}
	}

	if d.HasChange("prepend_as_path1") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path1").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		err := client.EditTransitConnectionASPathPrepend(transitGatewayPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path1: %v", err)
		}

	}

	if d.HasChange("prepend_as_path2") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path2").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		transitGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: d.Get("transit_gateway_name2").(string),
			TransitGatewayName2: d.Get("transit_gateway_name1").(string),
		}
		err := client.EditTransitConnectionASPathPrepend(transitGwPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path2: %v", err)
		}

	}

	if d.HasChanges("gateway1_logical_ifnames") || d.HasChange("gateway2_logical_ifnames") {
		return fmt.Errorf("cannot update logical interface names for edge transit peerings")
	}

	d.Partial(false)
	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return resourceAviatrixTransitGatewayPeeringRead(d, meta)
}

func resourceAviatrixTransitGatewayPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

	err := client.DeleteTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transit Gateway peering: %s", err)
	}

	return nil
}

// reverse interface name translations from this "eth0": "wan.0" => "wan0": "eth0"
func ReverseIfnameTranslation(ifnames map[string]string) map[string]string {
	reversed := make(map[string]string)

	for orig, translated := range ifnames {
		// Replace '.' with an empty string to match expected format
		reversed[strings.ReplaceAll(translated, ".", "")] = orig
	}
	return reversed
}

// SetWanInterfaces sets the WAN interface names based on logical interface mappings.
func SetWanInterfaces(logicalIfNames []interface{}, reversedInterfaceNames map[string]string) (string, error) {
	var wanInterfaces []string
	for _, logicalIfName := range logicalIfNames {
		ifName, ok := logicalIfName.(string)
		if !ok {
			return "", fmt.Errorf("logical_ifnames must be a list of strings")
		}
		interfaceName, exists := reversedInterfaceNames[ifName]
		if !exists {
			return "", fmt.Errorf("logical interface name %s not found in translation map", ifName)
		}
		wanInterfaces = append(wanInterfaces, interfaceName)
	}

	return strings.Join(wanInterfaces, ","), nil
}
