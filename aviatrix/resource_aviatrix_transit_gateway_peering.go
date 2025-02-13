package aviatrix

import (
	"context"
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
				Type:        schema.TypeList,
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
	transitGatewayName1, ok := d.Get("transit_gateway_name1").(string)
	if !ok {
		return fmt.Errorf("transit_gateway_name1 is required")
	}

	transitGatewayName2, ok := d.Get("transit_gateway_name2").(string)
	if !ok {
		return fmt.Errorf("transit_gateway_name2 is required")
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: transitGatewayName1,
		TransitGatewayName2: transitGatewayName2,
	}

	transitGatewayPeering.EnableOverPrivateNetwork, ok = d.Get("enable_peering_over_private_network").(bool)
	if !ok {
		return fmt.Errorf("enable_peering_over_private_network is required for edge gateway peering")
	}
	transitGatewayPeering.EnableJumboFrame, ok = d.Get("jumbo_frame").(bool)
	if !ok {
		return fmt.Errorf("jumbo_frame is required for edge gateway peering")
	}
	transitGatewayPeering.EnableInsaneMode, ok = d.Get("insane_mode").(bool)
	if !ok {
		return fmt.Errorf("insane_mode is required for edge gateway peering")
	}

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
	if !goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EdgeRelatedCloudTypes) && !goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EdgeRelatedCloudTypes) {
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

	err = client.CreateTransitGatewayPeering(context.Background(), transitGatewayPeering)
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

	if err := d.Set("enable_peering_over_private_network", transitGatewayPeering.EnableOverPrivateNetwork); err != nil {
		return fmt.Errorf("failed to set enable_peering_over_private_network: %w", err)
	}
	if err := d.Set("jumbo_frame", transitGatewayPeering.EnableJumboFrame); err != nil {
		return fmt.Errorf("failed to set jumbo_frame: %w", err)
	}
	if err := d.Set("insane_mode", transitGatewayPeering.EnableInsaneMode); err != nil {
		return fmt.Errorf("failed to set insane_mode: %w", err)
	}
	if err := d.Set("gateway1_logical_ifnames", transitGatewayPeering.Gateway1LogicalIfNames); err != nil {
		return fmt.Errorf("failed to set gateway1_logical_ifnames: %w", err)
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

func setWanInterfaceNames(
	logicalIfNames []string,
	cloudType int,
	gatewayDetails *goaviatrix.Gateway,
	gatewayPrefix string,
	transitGatewayPeering *goaviatrix.TransitGatewayPeering,
) error {
	if len(logicalIfNames) == 0 {
		return nil
	}

	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGENEO) {
		// Process logical interface names for Equinix/NEO cloud types
		reversedInterfaceNames := ReverseIfnameTranslation(gatewayDetails.IfNamesTranslation)
		wanInterfacesStr, err := SetWanInterfaces(convertToInterfaceSlice(logicalIfNames), reversedInterfaceNames)
		if err != nil {
			return fmt.Errorf("failed to set %s WAN interfaces to create edge peering: %w", gatewayPrefix, err)
		}

		if gatewayPrefix == "gateway1" {
			transitGatewayPeering.SrcWanInterfaces = wanInterfacesStr
		} else {
			transitGatewayPeering.DstWanInterfaces = wanInterfacesStr
		}
	} else if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEMEGAPORT) {
		// Process logical interface names for Megaport cloud type
		if gatewayPrefix == "gateway1" {
			transitGatewayPeering.Gateway1LogicalIfNames = logicalIfNames
			log.Printf("[INFO] Gateway1 Logical Interface Names: %#v", transitGatewayPeering.Gateway1LogicalIfNames)
		} else {
			transitGatewayPeering.Gateway2LogicalIfNames = logicalIfNames
			log.Printf("[INFO] Gateway2 Logical Interface Names: %#v", transitGatewayPeering.Gateway2LogicalIfNames)
		}
	}

	return nil
}

func convertToInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, v := range strs {
		result[i] = v
	}
	return result
}

func getStringListFromResource(d *schema.ResourceData, key string) ([]string, error) {
	var result []string
	value, ok := d.GetOk(key)
	// don't set the value if it's not present.
	if !ok {
		return nil, nil
	}
	interfaceList, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s is not a list of strings", key)
	}
	for _, v := range interfaceList {
		strValue, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%s contains non-string elements", key)
		}
		result = append(result, strValue)
	}
	return result, nil
}

func setASPathPrepend(d *schema.ResourceData, client *goaviatrix.Client, prependField, transitGatewayName1, transitGatewayName2 string) error {
	if _, ok := d.GetOk(prependField); ok {
		prependASPath, err := getStringListFromResource(d, prependField)
		if err != nil {
			return fmt.Errorf("%s is not a list of strings", prependField)
		}

		// Set up the Transit Gateway Peering struct with the appropriate names.
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: transitGatewayName1,
			TransitGatewayName2: transitGatewayName2,
		}

		// Call the client method to edit the transit connection.
		err = client.EditTransitConnectionASPathPrepend(transGwPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set %s: %w", prependField, err)
		}
	}
	return nil
}

func getNonEATPeeringOptions(d *schema.ResourceData) (map[string]bool, error) {
	transitPeering := make(map[string]bool)
	enableMaxPerformance, err := getBooleanValue(d, "enable_max_performance")
	if err != nil {
		return nil, err
	}
	transitPeering["enable_max_performance"] = enableMaxPerformance
	insaneMode, err := getBooleanValue(d, "enable_insane_mode_encryption_over_internet")
	if err != nil {
		return nil, err
	}
	transitPeering["insane_mode"] = insaneMode
	peeringOverPrivate, err := getBooleanValue(d, "enable_peering_over_private_network")
	if err != nil {
		return nil, err
	}
	transitPeering["enable_peering_over_private_network"] = peeringOverPrivate
	singleTunnelMode, err := getBooleanValue(d, "enable_single_tunnel_mode")
	if err != nil {
		return nil, err
	}
	transitPeering["enable_single_tunnel_mode"] = singleTunnelMode
	return transitPeering, nil
}

func setNonEATPeeringOptions(d *schema.ResourceData, transitGatewayPeering *goaviatrix.TransitGatewayPeering) error {
	// Validate and set boolean fields
	transitPeering, err := getNonEATPeeringOptions(d)
	if err != nil {
		return err
	}
	tunnelCount, ok := d.Get("tunnel_count").(int)
	if !ok {
		return fmt.Errorf("tunnel_count must be an integer")
	}
	if err := validateTunnelCount(transitPeering["insane_mode"], tunnelCount); err != nil {
		return err
	}
	if tunnelCount != 0 {
		transitGatewayPeering.TunnelCount = tunnelCount
	}

	// Set No max Performance and Insane Mode
	transitGatewayPeering.NoMaxPerformance = !transitPeering["enable_max_performance"]
	transitGatewayPeering.InsaneModeOverInternet = transitPeering["insane_mode"]
	// set private ip peering
	if transitPeering["enable_peering_over_private_network"] {
		transitGatewayPeering.PrivateIPPeering = "yes"
		if transitPeering["insane_mode"] {
			return fmt.Errorf("enable_peering_over_private_network conflicts with enable_insane_mode_encryption_over_internet")
		}
	} else {
		transitGatewayPeering.PrivateIPPeering = "no"
	}

	// Validate and set Single Tunnel Mode
	if transitPeering["enable_single_tunnel_mode"] {
		if transitGatewayPeering.PrivateIPPeering == "no" {
			return fmt.Errorf("enable_single_tunnel_mode is only valid when enable_peering_over_private_network is set to true")
		}
		transitGatewayPeering.SingleTunnel = "yes"
	}

	return nil
}

// Helper function to validate Tunnel Count
func validateTunnelCount(insaneMode bool, tunnelCount int) error {
	if (insaneMode && tunnelCount == 0) || (!insaneMode && tunnelCount != 0) {
		return fmt.Errorf("tunnel_count is only valid when enable_insane_mode_encryption_over_internet is set to true and must be > 0")
	}
	return nil
}

// Helper function to get multiple boolean values and validate them
func getBooleanValue(d *schema.ResourceData, key string) (bool, error) {
	value, ok := d.Get(key).(bool)
	if !ok {
		return false, fmt.Errorf("%s is required and must be a boolean", key)
	}
	return value, nil
}

func setExcludedResources(d *schema.ResourceData, transitGatewayPeering *goaviatrix.TransitGatewayPeering) error {
	// Set Gateway1 Excluded Network CIDRs
	gw1Cidrs, err := getStringListFromResource(d, "gateway1_excluded_network_cidrs")
	if err != nil {
		return err
	}
	transitGatewayPeering.Gateway1ExcludedCIDRs = strings.Join(gw1Cidrs, ",")

	// Set Gateway2 Excluded Network CIDRs
	gw2Cidrs, err := getStringListFromResource(d, "gateway2_excluded_network_cidrs")
	if err != nil {
		return err
	}
	transitGatewayPeering.Gateway2ExcludedCIDRs = strings.Join(gw2Cidrs, ",")

	// Set Gateway1 Excluded TGW Connections
	gw1Tgws, err := getStringListFromResource(d, "gateway1_excluded_tgw_connections")
	if err != nil {
		return err
	}
	transitGatewayPeering.Gateway1ExcludedTGWConnections = strings.Join(gw1Tgws, ",")

	// Set Gateway2 Excluded TGW Connections
	gw2Tgws, err := getStringListFromResource(d, "gateway2_excluded_tgw_connections")
	if err != nil {
		return err
	}
	transitGatewayPeering.Gateway2ExcludedTGWConnections = strings.Join(gw2Tgws, ",")

	return nil
}

func getGatewayDetails(client *goaviatrix.Client, gatewayName string) (*goaviatrix.Gateway, error) {
	gateway := &goaviatrix.Gateway{
		GwName: gatewayName,
	}
	gatewayDetails, err := client.GetGateway(gateway)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway %s: %w", gatewayName, err)
	}
	return gatewayDetails, nil
}
