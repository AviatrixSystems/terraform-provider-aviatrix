package aviatrix

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func validateIdentifierValue(val interface{}, key string) (warns []string, errs []error) {
	value, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string, got: %T", key, val))
		return warns, errs
	}
	// Check if the value is "auto"
	if value == "auto" {
		return
	}
	// Check if the value is a valid MAC address
	macRegex := `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`
	if matched, _ := regexp.MatchString(macRegex, value); matched {
		return
	}
	// Check if the value is a valid PCI ID
	pciRegex := `^(pci@)?[0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}\.[0-9a-fA-F]$`
	if matched, _ := regexp.MatchString(pciRegex, value); matched {
		return
	}
	errs = append(errs, fmt.Errorf("%q must be a valid MAC address, PCI ID, or 'auto', got: %s", key, value))
	return warns, errs
}

func getCustomInterfaceMapDetails(customInterfaceMap []interface{}) (map[string]goaviatrix.CustomInterfaceMap, error) {
	// Create a map to structure the Custom interface map data
	customInterfaceMapStructured := make(map[string]goaviatrix.CustomInterfaceMap)

	// Populate the structured map
	for _, customInterfaceMap := range customInterfaceMap {
		customInterface, ok := customInterfaceMap.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type: expected map[string]interface{}, got %T", customInterfaceMap)
		}
		logicalIfName, ok := customInterface["logical_ifname"].(string)
		if !ok {
			return nil, fmt.Errorf("logical interface name must be a string")
		}

		identifierType, ok := customInterface["identifier_type"].(string)
		if !ok {
			return nil, fmt.Errorf("identifier type must be a string")
		}
		identifierValue, ok := customInterface["identifier_value"].(string)
		if !ok {
			return nil, fmt.Errorf("identifier value must be a string")
		}

		// Append the EIP entry to the corresponding interface
		customInterfaceEntry := goaviatrix.CustomInterfaceMap{
			IdentifierType:  identifierType,
			IdentifierValue: identifierValue,
		}
		customInterfaceMapStructured[logicalIfName] = customInterfaceEntry
	}

	return customInterfaceMapStructured, nil
}

func setCustomInterfaceMapping(customInterfaceMap map[string]goaviatrix.CustomInterfaceMap, userCustomInterfaceOrder []string) ([]interface{}, error) {
	var result []interface{}

	// Iterate over the user-provided order
	for _, logicalIfName := range userCustomInterfaceOrder {
		logicalName := strings.ToUpper(logicalIfName)
		mapping, exists := customInterfaceMap[logicalName]
		if !exists {
			return nil, fmt.Errorf("logical interface name %s not found in custom interface map", logicalIfName)
		}

		if mapping.IdentifierType == "" {
			return nil, fmt.Errorf("identifier type cannot be empty for logical interface: %s", logicalIfName)
		}
		if mapping.IdentifierValue == "" {
			return nil, fmt.Errorf("identifier value cannot be empty for logical interface: %s", logicalIfName)
		}

		entry := map[string]interface{}{
			"logical_ifname":   logicalIfName,
			"identifier_type":  mapping.IdentifierType,
			"identifier_value": mapping.IdentifierValue,
		}
		result = append(result, entry)
	}
	return result, nil
}

func getCustomInterfaceOrder(userCustomInterfaceMapping []interface{}) ([]string, error) {
	var order []string

	for _, mapping := range userCustomInterfaceMapping {
		mappingMap, ok := mapping.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type: expected map[string]interface{}, got %T", mapping)
		}

		logicalIfName, ok := mappingMap["logical_ifname"].(string)
		if !ok || logicalIfName == "" {
			return nil, fmt.Errorf("logical_ifname must be a non-empty string")
		}

		order = append(order, logicalIfName)
	}

	return order, nil
}

func populateInterfaces(d *schema.ResourceData, edgeSpoke *goaviatrix.EdgeSpoke) error {
	interfacesRaw, ok := d.Get("interfaces").(*schema.Set)
	if !ok {
		return fmt.Errorf("failed to get interfaces: expected *schema.Set, got %T", d.Get("interfaces"))
	}

	interfaces := interfacesRaw.List()
	for _, if0 := range interfaces {
		if1, ok := if0.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to get interface: expected map[string]interface{}, got %T", if0)
		}

		if2, err := buildEdgeSpokeInterface(if1)
		if err != nil {
			return err
		}

		edgeSpoke.InterfaceList = append(edgeSpoke.InterfaceList, if2)
	}
	return nil
}

func buildEdgeSpokeInterface(if1 map[string]interface{}) (*goaviatrix.EdgeSpokeInterface, error) {
	ifName, err := getStringFromMap(if1, "name")
	if err != nil {
		return nil, err
	}

	ifType, err := getStringFromMap(if1, "type")
	if err != nil {
		return nil, err
	}

	enableDhcp, err := getBoolFromMap(if1, "enable_dhcp")
	if err != nil {
		return nil, err
	}

	publicIP, err := getStringFromMap(if1, "wan_public_ip")
	if err != nil {
		return nil, err
	}

	ipAddr, err := getStringFromMap(if1, "ip_address")
	if err != nil {
		return nil, err
	}

	gatewayIP, err := getStringFromMap(if1, "gateway_ip")
	if err != nil {
		return nil, err
	}

	tag, err := getStringFromMap(if1, "tag")
	if err != nil {
		return nil, err
	}

	ipv6Addr, err := getStringFromMap(if1, "ipv6_address")
	if err != nil {
		return nil, err
	}

	gatewayIPv6, err := getStringFromMap(if1, "gateway_ipv6")
	if err != nil {
		return nil, err
	}

	if2 := &goaviatrix.EdgeSpokeInterface{
		IfName:      ifName,
		Type:        ifType,
		Dhcp:        enableDhcp,
		PublicIp:    publicIP,
		IpAddr:      ipAddr,
		GatewayIp:   gatewayIP,
		Tag:         tag,
		IPv6Addr:    ipv6Addr,
		GatewayIPv6: gatewayIPv6,
	}

	if ifType == "LAN" {
		if err := populateLANFields(if1, if2); err != nil {
			return nil, err
		}
	}

	return if2, nil
}

func populateLANFields(if1 map[string]interface{}, if2 *goaviatrix.EdgeSpokeInterface) error {
	enableVrrp, err := getBoolFromMap(if1, "enable_vrrp")
	if err != nil {
		return err
	}

	virtualIP, err := getStringFromMap(if1, "vrrp_virtual_ip")
	if err != nil {
		return err
	}

	if2.VrrpState = enableVrrp
	if2.VirtualIp = virtualIP
	return nil
}

func populateVlans(d *schema.ResourceData, edgeSpoke *goaviatrix.EdgeSpoke) error {
	vlanRaw, ok := d.Get("vlan").(*schema.Set)
	if !ok {
		return fmt.Errorf("failed to get vlan: expected *schema.Set, got %T", d.Get("vlan"))
	}

	vlan := vlanRaw.List()
	for _, vlan0 := range vlan {
		vlan1, ok := vlan0.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to get vlan entry: expected map[string]interface{}, got %T", vlan0)
		}

		vlan2, err := buildEdgeSpokeVlan(vlan1)
		if err != nil {
			return err
		}

		edgeSpoke.VlanList = append(edgeSpoke.VlanList, vlan2)
	}
	return nil
}

func buildEdgeSpokeVlan(vlan1 map[string]interface{}) (*goaviatrix.EdgeSpokeVlan, error) {
	parentInterface, err := getStringFromMap(vlan1, "parent_interface_name")
	if err != nil {
		return nil, err
	}

	ipAddr, err := getStringFromMap(vlan1, "ip_address")
	if err != nil {
		return nil, err
	}

	gatewayIP, err := getStringFromMap(vlan1, "gateway_ip")
	if err != nil {
		return nil, err
	}

	peerIPAddr, err := getStringFromMap(vlan1, "peer_ip_address")
	if err != nil {
		return nil, err
	}

	peerGatewayIP, err := getStringFromMap(vlan1, "peer_gateway_ip")
	if err != nil {
		return nil, err
	}

	virtualIP, err := getStringFromMap(vlan1, "vrrp_virtual_ip")
	if err != nil {
		return nil, err
	}

	tag, err := getStringFromMap(vlan1, "tag")
	if err != nil {
		return nil, err
	}

	vlanID, err := getIntFromMap(vlan1, "vlan_id")
	if err != nil {
		return nil, err
	}

	return &goaviatrix.EdgeSpokeVlan{
		ParentInterface: parentInterface,
		IpAddr:          ipAddr,
		GatewayIp:       gatewayIP,
		PeerIpAddr:      peerIPAddr,
		PeerGatewayIp:   peerGatewayIP,
		VirtualIp:       virtualIP,
		Tag:             tag,
		VlanId:          strconv.Itoa(vlanID),
	}, nil
}

func populateCustomInterfaceMapping(d *schema.ResourceData, edgeSpoke *goaviatrix.EdgeSpoke) error {
	customInterfaceMapping, ok := d.Get("custom_interface_mapping").([]interface{})
	if ok {
		customInterfaceMap, err := getCustomInterfaceMapDetails(customInterfaceMapping)
		if err != nil {
			return err
		}
		edgeSpoke.CustomInterfaceMapping = customInterfaceMap
	}
	return nil
}

func getStringFromMap(data map[string]interface{}, key string) (string, error) {
	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("invalid type for '%s': expected string, got %T", key, data[key])
	}
	return value, nil
}

func getBoolFromMap(data map[string]interface{}, key string) (bool, error) {
	value, ok := data[key].(bool)
	if !ok {
		return false, fmt.Errorf("invalid type for '%s': expected bool, got %T", key, data[key])
	}
	return value, nil
}

func getIntFromMap(data map[string]interface{}, key string) (int, error) {
	value, ok := data[key].(int)
	if !ok {
		return 0, fmt.Errorf("invalid type for '%s': expected int, got %T", key, data[key])
	}
	return value, nil
}
