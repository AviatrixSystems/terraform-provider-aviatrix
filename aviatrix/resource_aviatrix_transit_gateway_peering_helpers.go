package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
func SetWanInterfaces(logicalIfNames []string, reversedInterfaceNames map[string]string) (string, error) {
	var wanInterfaces []string
	for _, logicalIfName := range logicalIfNames {
		interfaceName, exists := reversedInterfaceNames[logicalIfName]
		if !exists {
			return "", fmt.Errorf("logical interface name %s not found in translation map", logicalIfName)
		}
		wanInterfaces = append(wanInterfaces, interfaceName)
	}

	return strings.Join(wanInterfaces, ","), nil
}

// setWanInterfaceNames sets the WAN interface names and logical interface names based on the cloud type.
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

	/* the cloud type here is a bit mask and we set the src/dst wan interfaces based on the cloud type
	 * if the cloud type is Equinix or NEO, we need to convert the logical interface names from wan0, wan1 to eth0, eth1. The api for these edge types does not support logical interface names.
	 * if the cloud type is Megaport, we need to set the logical interface names as the api for this edge type supports logical interface names.
	 */
	switch {
	case goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGENEO):
		// Process logical interface names for Equinix/NEO cloud types
		reversedInterfaceNames := ReverseIfnameTranslation(gatewayDetails.IfNamesTranslation)
		wanInterfacesStr, err := SetWanInterfaces(logicalIfNames, reversedInterfaceNames)
		if err != nil {
			return fmt.Errorf("failed to set %s WAN interfaces to create edge peering: %w", gatewayPrefix, err)
		}
		if gatewayPrefix == "gateway1" {
			transitGatewayPeering.SrcWanInterfaces = wanInterfacesStr
		} else {
			transitGatewayPeering.DstWanInterfaces = wanInterfacesStr
		}
	case goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEMEGAPORT):
		if gatewayPrefix == "gateway1" {
			transitGatewayPeering.Gateway1LogicalIfNames = logicalIfNames
			log.Printf("[INFO] Gateway1 Logical Interface Names: %#v", transitGatewayPeering.Gateway1LogicalIfNames)
		} else {
			transitGatewayPeering.Gateway2LogicalIfNames = logicalIfNames
			log.Printf("[INFO] Gateway2 Logical Interface Names: %#v", transitGatewayPeering.Gateway2LogicalIfNames)
		}
	default:
		return nil
	}

	return nil
}

// getStringListFromResource gets a list of strings from a resource data.
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

// setASPathPrepend sets the AS Path Prepend for the transit connection.
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

// getNonEATPeeringOptions gets the non-EAT peering options from the resource data.
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

// setNonEATPeeringOptions sets the non-EAT peering options for the transit gateway peering.
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

// setExcludedResources sets the excluded resources for the transit gateway peering.
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

// getGatewayDetails gets the gateway details from the client.
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

// setNonEATPeering sets the non-EAT peering options based on the cloud types of the gateways.
func setNonEATPeering(gateway1CloudType, gateway2CloudType int) bool {
	return !goaviatrix.IsCloudType(gateway1CloudType, goaviatrix.EdgeRelatedCloudTypes) &&
		!goaviatrix.IsCloudType(gateway2CloudType, goaviatrix.EdgeRelatedCloudTypes)
}
