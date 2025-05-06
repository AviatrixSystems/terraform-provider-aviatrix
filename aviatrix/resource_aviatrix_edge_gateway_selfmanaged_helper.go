package aviatrix

import (
	"fmt"
	"strconv"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

	if2 := &goaviatrix.EdgeSpokeInterface{
		IfName:    ifName,
		Type:      ifType,
		Dhcp:      enableDhcp,
		PublicIp:  publicIP,
		IpAddr:    ipAddr,
		GatewayIp: gatewayIP,
		Tag:       tag,
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
