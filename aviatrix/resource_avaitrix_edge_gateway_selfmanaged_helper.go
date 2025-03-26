package aviatrix

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	pciRegex := `^[0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}\.[0-9a-fA-F]$`
	if matched, _ := regexp.MatchString(pciRegex, value); matched {
		return
	}
	errs = append(errs, fmt.Errorf("%q must be a valid MAC address, PCI ID, or 'auto', got: %s", key, value))
	return warns, errs
}

func getCustomInterfaceMapDetails(customInterfaceMap []interface{}) (map[string][]goaviatrix.CustomInterfaceMap, error) {
	// Create a map to structure the Custom interface map data
	customInterfaceMapStructured := make(map[string][]goaviatrix.CustomInterfaceMap)

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
		identifierValue, ok := customInterface["idenitifer_value"].(string)
		if !ok {
			return nil, fmt.Errorf("identifier value must be a string")
		}

		// Append the EIP entry to the corresponding interface
		customInterfaceEntry := goaviatrix.CustomInterfaceMap{
			IdentifierType:  identifierType,
			IdentifierValue: identifierValue,
		}
		customInterfaceMapStructured[logicalIfName] = append(customInterfaceMapStructured[logicalIfName], customInterfaceEntry)
	}

	return customInterfaceMapStructured, nil
}

func setCustomInterfaceMapping(customInterfaceMap map[string][]goaviatrix.CustomInterfaceMap, userCustomInterfaceOrder []string) ([]interface{}, error) {
	var result []interface{}

	// Iterate over the user-provided order
	for _, logicalIfName := range userCustomInterfaceOrder {
		mappings, exists := customInterfaceMap[logicalIfName]
		if !exists {
			return nil, fmt.Errorf("logical interface name %s not found in custom interface map", logicalIfName)
		}

		for _, mapping := range mappings {
			if mapping.IdentifierType == "" {
				return nil, fmt.Errorf("identifier type cannot be empty for logical interface: %s", logicalIfName)
			}
			if mapping.IdentifierValue == "" {
				return nil, fmt.Errorf("identifier value cannot be empty for logical interface: %s", logicalIfName)
			}

			entry := map[string]interface{}{
				"logical_ifname":   logicalIfName,
				"identifier_type":  mapping.IdentifierType,
				"idenitifer_value": mapping.IdentifierValue,
			}
			result = append(result, entry)
		}
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
		if2 := &goaviatrix.EdgeSpokeInterface{
			IfName:    if1["name"].(string),
			Type:      if1["type"].(string),
			Dhcp:      if1["enable_dhcp"].(bool),
			PublicIp:  if1["wan_public_ip"].(string),
			IpAddr:    if1["ip_address"].(string),
			GatewayIp: if1["gateway_ip"].(string),
			Tag:       if1["tag"].(string),
		}

		if if1["type"].(string) == "LAN" {
			if2.VrrpState = if1["enable_vrrp"].(bool)
			if2.VirtualIp = if1["vrrp_virtual_ip"].(string)
		}

		edgeSpoke.InterfaceList = append(edgeSpoke.InterfaceList, if2)
	}
	return nil
}

func populateVlans(d *schema.ResourceData, edgeSpoke *goaviatrix.EdgeSpoke) error {
	vlanRaw, ok := d.Get("vlan").(*schema.Set)
	if !ok {
		return fmt.Errorf("failed to get vlan: expected *schema.Set, got %T", d.Get("vlan"))
	}
	vlan := vlanRaw.List()
	for _, vlan0 := range vlan {
		vlan1 := vlan0.(map[string]interface{})
		vlan2 := &goaviatrix.EdgeSpokeVlan{
			ParentInterface: vlan1["parent_interface_name"].(string),
			IpAddr:          vlan1["ip_address"].(string),
			GatewayIp:       vlan1["gateway_ip"].(string),
			PeerIpAddr:      vlan1["peer_ip_address"].(string),
			PeerGatewayIp:   vlan1["peer_gateway_ip"].(string),
			VirtualIp:       vlan1["vrrp_virtual_ip"].(string),
			Tag:             vlan1["tag"].(string),
		}
		vlan2.VlanId = strconv.Itoa(vlan1["vlan_id"].(int))
		edgeSpoke.VlanList = append(edgeSpoke.VlanList, vlan2)
	}
	return nil
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
