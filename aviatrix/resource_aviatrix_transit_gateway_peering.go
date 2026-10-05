package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitGatewayPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitGatewayPeeringCreate,
		Read:   resourceAviatrixTransitGatewayPeeringRead,
		Update: resourceAviatrixTransitGatewayPeeringUpdate,
		Delete: resourceAviatrixTransitGatewayPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
					return !getBool(d, "enable_peering_over_private_network")
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
			"jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable jumbo frame for over private peering with Edge Transit",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Gateway 2 logical interface names for edge gateways where the peering terminates",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"disable_activemesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Disable ActiveMesh, no crossing tunnels",
			},
		},
	}
}

func resourceAviatrixTransitGatewayPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	flag := false
	transitGatewayName1 := getString(d, "transit_gateway_name1")

	transitGatewayName2 := getString(d, "transit_gateway_name2")

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: transitGatewayName1,
		TransitGatewayName2: transitGatewayName2,
		DisableActivemesh:   getBool(d, "disable_activemesh"),
	}

	transitGatewayPeering.EnableOverPrivateNetwork = getBool(d, "enable_peering_over_private_network")

	transitGatewayPeering.EnableJumboFrame = getBool(d, "jumbo_frame")

	// insane_mode is optional for edge gateway peering
	transitGatewayPeering.EnableInsaneMode = getBool(d, "insane_mode")

	gateway1Details, err := getGatewayDetails(client, transitGatewayName1)
	if err != nil {
		return err
	}
	gateway1CloudType := gateway1Details.CloudType
	gateway2Details, err := getGatewayDetails(client, transitGatewayName2)
	if err != nil {
		return err
	}
	gateway2CloudType := gateway2Details.CloudType
	// Set source WAN interface names for gateway1
	if goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		transit1InterfaceRaw, ok := d.GetOk("gateway1_logical_ifnames")
		if !ok {
			return fmt.Errorf("gateway1_logical_ifnames is required for edge gateway peering")
		}
		if _, ok := transit1InterfaceRaw.([]interface{}); !ok {
			return fmt.Errorf("gateway1_logical_ifnames must be a list of strings")
		}
		gw1LogicalIfNames := getStringList(d, "gateway1_logical_ifnames")
		if err := setWanInterfaceNames(gw1LogicalIfNames, gateway1CloudType, gateway1Details, "gateway1", transitGatewayPeering); err != nil {
			return err
		}
	}

	// Set destination WAN interface names for gateway2
	if goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		transit2InterfaceRaw, ok := d.GetOk("gateway2_logical_ifnames")
		if !ok {
			return fmt.Errorf("gateway2_logical_ifnames is required for edge gateway peering")
		}
		if _, ok := transit2InterfaceRaw.([]interface{}); !ok {
			return fmt.Errorf("gateway2_logical_ifnames must be a list of strings")
		}
		gw2LogicalIfNames := getStringList(d, "gateway2_logical_ifnames")
		if err := setWanInterfaceNames(gw2LogicalIfNames, gateway2CloudType, gateway2Details, "gateway2", transitGatewayPeering); err != nil {
			return err
		}
	}

	if err := setExcludedResources(d, transitGatewayPeering); err != nil {
		return err
	}

	// options only supported for non EAT peerings
	if setNonEATPeering(gateway1CloudType, gateway2CloudType) {
		log.Printf("[INFO] Setting non EAT peering options")
		if err := setNonEATPeeringOptions(d, transitGatewayPeering); err != nil {
			return err
		}
	}

	log.Printf("[INFO] Creating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)
	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	defer func() {
		if err := resourceAviatrixTransitGatewayPeeringReadIfRequired(d, meta, &flag); err != nil {
			log.Printf("[ERROR] Failed to read Aviatrix Transit Gateway peering: %v", err)
		}
	}()

	const timeoutDuration = 180 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	err = client.CreateTransitGatewayPeering(ctx, transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit Gateway peering: %w", err)
	}

	// Set AS Path prepend for gateway1 (prepend_as_path1)
	if err := setASPathPrepend(d, client, "prepend_as_path1", transitGatewayName1, transitGatewayName2); err != nil {
		return err
	}

	// Set AS Path prepend for gateway2 (prepend_as_path2)
	if err := setASPathPrepend(d, client, "prepend_as_path2", transitGatewayName2, transitGatewayName1); err != nil {
		return err
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
	client := mustClient(meta)

	transitGwName1 := getString(d, "transit_gateway_name1")
	transitGwName2 := getString(d, "transit_gateway_name2")

	if transitGwName1 == "" || transitGwName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return fmt.Errorf("invalid import id expected transit_gateway_name1~transit_gateway_name2")
		}
		mustSet(d, "transit_gateway_name1", parts[0])
		mustSet(d, "transit_gateway_name2", parts[1])
		d.SetId(id)
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: getString(d, "transit_gateway_name1"),
		TransitGatewayName2: getString(d, "transit_gateway_name2"),
	}

	transitGatewayPeering, err := client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get transit peering details: %w", err)
	}

	if err := d.Set("enable_peering_over_private_network", transitGatewayPeering.EnableOverPrivateNetwork); err != nil {
		return fmt.Errorf("failed to set enable_peering_over_private_network: %w", err)
	}
	if err := d.Set("jumbo_frame", transitGatewayPeering.EnableJumboFrame); err != nil {
		return fmt.Errorf("failed to set jumbo_frame: %w", err)
	}
	gateway1Details, err := getGatewayDetails(client, transitGatewayPeering.TransitGatewayName1)
	if err != nil {
		return fmt.Errorf("failed to get gateway1 details: %w", err)
	}
	gateway1CloudType := gateway1Details.CloudType
	gateway2Details, err := getGatewayDetails(client, transitGatewayPeering.TransitGatewayName2)
	if err != nil {
		return fmt.Errorf("failed to get gateway2 details: %w", err)
	}
	gateway2CloudType := gateway2Details.CloudType
	// Set insane mode only for the edge gateways and only if user explicitly set it
	if goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EdgeRelatedCloudTypes) || goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		// Only set insane_mode in state if user explicitly provided it in configuration
		if err := d.Set("insane_mode", transitGatewayPeering.EnableInsaneMode); err != nil {
			return fmt.Errorf("failed to set insane_mode: %w", err)
		}
	}

	// Set the gateway1 logical interface names for the edge gateways
	if len(transitGatewayPeering.Gateway1LogicalIfNames) > 0 {
		logicalIfNames, err := getLogicalIfNames(gateway1Details, transitGatewayPeering.Gateway1LogicalIfNames)
		if err != nil {
			return fmt.Errorf("failed to set gateway1_logical_ifnames: %w", err)
		}
		_ = d.Set("gateway1_logical_ifnames", logicalIfNames)
	}

	// Set the gateway2 logical interface names for the edge gateways
	if len(transitGatewayPeering.Gateway2LogicalIfNames) > 0 {
		logicalIfNames, err := getLogicalIfNames(gateway2Details, transitGatewayPeering.Gateway2LogicalIfNames)
		if err != nil {
			return fmt.Errorf("failed to set gateway2 logical ifnames: %w", err)
		}
		_ = d.Set("gateway2_logical_ifnames", logicalIfNames)
	}
	if err := d.Set("gateway2_logical_ifnames", transitGatewayPeering.Gateway2LogicalIfNames); err != nil {
		return fmt.Errorf("failed to set gateway2_logical_ifnames: %w", err)
	}

	gw1CidrsFromConfig := getStringList(d, "gateway1_excluded_network_cidrs")
	err = setConfigValueIfEquivalent(d, "gateway1_excluded_network_cidrs", gw1CidrsFromConfig, transitGatewayPeering.Gateway1ExcludedCIDRsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway1_excluded_network_cidrs to state: %w", err)
	}
	gw2CidrsFromConfig := getStringList(d, "gateway2_excluded_network_cidrs")
	err = setConfigValueIfEquivalent(d, "gateway2_excluded_network_cidrs", gw2CidrsFromConfig, transitGatewayPeering.Gateway2ExcludedCIDRsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway2_excluded_network_cidrs to state: %w", err)
	}
	gw1TgwsFromConfig := getStringList(d, "gateway1_excluded_tgw_connections")
	err = setConfigValueIfEquivalent(d, "gateway1_excluded_tgw_connections", gw1TgwsFromConfig, transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway1_excluded_tgw_connections to state: %w", err)
	}
	gw2TgwsFromConfig := getStringList(d, "gateway2_excluded_tgw_connections")
	err = setConfigValueIfEquivalent(d, "gateway2_excluded_tgw_connections", gw2TgwsFromConfig, transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway2_excluded_tgw_connections to state: %w", err)
	}

	if transitGatewayPeering.PrependAsPath1 != "" {
		var prependAsPath []string
		for _, str := range strings.Split(transitGatewayPeering.PrependAsPath1, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path1", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path1: %w", err)
		}
	}
	if transitGatewayPeering.PrependAsPath2 != "" {
		var prependAsPath []string
		for _, str := range strings.Split(transitGatewayPeering.PrependAsPath2, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path2", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path2: %w", err)
		}
	}
	if err := d.Set("enable_peering_over_private_network", transitGatewayPeering.PrivateIPPeering == "yes"); err != nil {
		return fmt.Errorf("failed to set enable_peering_over_private_network: %w", err)
	}

	enableSingleTunnel := transitGatewayPeering.PrivateIPPeering == "yes" && transitGatewayPeering.SingleTunnel == "yes"
	if err := d.Set("enable_single_tunnel_mode", enableSingleTunnel); err != nil {
		return fmt.Errorf("failed to set enable_single_tunnel_mode: %w", err)
	}
	if err := d.Set("enable_insane_mode_encryption_over_internet", transitGatewayPeering.InsaneModeOverInternet); err != nil {
		return fmt.Errorf("failed to set enable_insane_mode_encryption_over_internet: %w", err)
	}

	if transitGatewayPeering.InsaneModeOverInternet {
		if err := d.Set("tunnel_count", transitGatewayPeering.TunnelCount); err != nil {
			return fmt.Errorf("failed to set tunnel_count: %w", err)
		}
	}

	if err := d.Set("enable_max_performance", !transitGatewayPeering.NoMaxPerformance); err != nil {
		return fmt.Errorf("failed to set enable_max_performance: %w", err)
	}
	if err := d.Set("disable_activemesh", transitGatewayPeering.DisableActivemesh); err != nil {
		return fmt.Errorf("error setting disable_activemesh: %w", err)
	}
	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return nil
}

func resourceAviatrixTransitGatewayPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	d.Partial(true)

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: getString(d, "transit_gateway_name1"),
		TransitGatewayName2: getString(d, "transit_gateway_name2"),
	}
	if d.HasChange("gateway1_excluded_network_cidrs") || d.HasChange("gateway2_excluded_network_cidrs") ||
		d.HasChange("gateway1_excluded_tgw_connections") || d.HasChange("gateway2_excluded_tgw_connections") {
		var gw1Cidrs []string
		for _, cidr := range getList(d, "gateway1_excluded_network_cidrs") {
			gw1Cidrs = append(gw1Cidrs, mustString(cidr))
		}
		var gw2Cidrs []string
		for _, cidr := range getList(d, "gateway2_excluded_network_cidrs") {
			gw2Cidrs = append(gw2Cidrs, mustString(cidr))
		}
		var gw1Tgws []string
		for _, tgw := range getList(d, "gateway1_excluded_tgw_connections") {
			gw1Tgws = append(gw1Tgws, mustString(tgw))
		}
		var gw2Tgws []string
		for _, tgw := range getList(d, "gateway2_excluded_tgw_connections") {
			gw2Tgws = append(gw2Tgws, mustString(tgw))
		}

		transitGatewayPeering.Gateway1ExcludedCIDRs = strings.Join(gw1Cidrs, ",")
		transitGatewayPeering.Gateway2ExcludedCIDRs = strings.Join(gw2Cidrs, ",")
		transitGatewayPeering.Gateway1ExcludedTGWConnections = strings.Join(gw1Tgws, ",")
		transitGatewayPeering.Gateway2ExcludedTGWConnections = strings.Join(gw2Tgws, ",")

		log.Printf("[INFO] Updating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)
		err := client.UpdateTransitGatewayPeering(transitGatewayPeering)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit Gateway peering: %w", err)
		}
	}

	if d.HasChange("prepend_as_path1") {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path1") {
			prependASPath = append(prependASPath, mustString(v))
		}

		err := client.EditTransitConnectionASPathPrepend(transitGatewayPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path1: %w", err)
		}

	}

	if d.HasChange("prepend_as_path2") {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path2") {
			prependASPath = append(prependASPath, mustString(v))
		}
		transitGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: getString(d, "transit_gateway_name2"),
			TransitGatewayName2: getString(d, "transit_gateway_name1"),
		}
		err := client.EditTransitConnectionASPathPrepend(transitGwPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path2: %w", err)
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
	client := mustClient(meta)

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: getString(d, "transit_gateway_name1"),
		TransitGatewayName2: getString(d, "transit_gateway_name2"),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

	err := client.DeleteTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transit Gateway peering: %w", err)
	}

	return nil
}
