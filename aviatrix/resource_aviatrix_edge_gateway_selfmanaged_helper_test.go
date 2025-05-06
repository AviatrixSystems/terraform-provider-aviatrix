package aviatrix

import (
	"reflect"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func TestBuildEdgeSpokeInterface(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		expected      *goaviatrix.EdgeSpokeInterface
		expectErr     bool
		expectedError string
	}{
		{
			name: "Valid input for non-LAN interface",
			input: map[string]interface{}{
				"name":          "eth0",
				"type":          "WAN",
				"enable_dhcp":   true,
				"wan_public_ip": "192.168.1.1",
				"ip_address":    "192.168.1.2",
				"gateway_ip":    "192.168.1.254",
				"tag":           "tag1",
			},
			expected: &goaviatrix.EdgeSpokeInterface{
				IfName:    "eth0",
				Type:      "WAN",
				Dhcp:      true,
				PublicIp:  "192.168.1.1",
				IpAddr:    "192.168.1.2",
				GatewayIp: "192.168.1.254",
				Tag:       "tag1",
			},
			expectErr: false,
		},
		{
			name: "Valid input for LAN interface",
			input: map[string]interface{}{
				"name":            "eth1",
				"type":            "LAN",
				"enable_dhcp":     false,
				"wan_public_ip":   "192.168.2.1",
				"ip_address":      "192.168.2.2",
				"gateway_ip":      "192.168.2.254",
				"tag":             "tag2",
				"enable_vrrp":     true,
				"vrrp_virtual_ip": "192.168.2.100",
			},
			expected: &goaviatrix.EdgeSpokeInterface{
				IfName:    "eth1",
				Type:      "LAN",
				Dhcp:      false,
				PublicIp:  "192.168.2.1",
				IpAddr:    "192.168.2.2",
				GatewayIp: "192.168.2.254",
				Tag:       "tag2",
				VrrpState: true,
				VirtualIp: "192.168.2.100",
			},
			expectErr: false,
		},
		{
			name: "Missing required field: name",
			input: map[string]interface{}{
				"type":          "WAN",
				"enable_dhcp":   true,
				"wan_public_ip": "192.168.1.1",
				"ip_address":    "192.168.1.2",
				"gateway_ip":    "192.168.1.254",
				"tag":           "tag1",
			},
			expectErr:     true,
			expectedError: "invalid type for 'name': expected string, got <nil>",
		},
		{
			name: "Missing required field: type",
			input: map[string]interface{}{
				"name":          "eth0",
				"enable_dhcp":   true,
				"wan_public_ip": "192.168.1.1",
				"ip_address":    "192.168.1.2",
				"gateway_ip":    "192.168.1.254",
				"tag":           "tag1",
			},
			expectErr:     true,
			expectedError: "invalid type for 'type': expected string, got <nil>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := buildEdgeSpokeInterface(test.input)

			if test.expectErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				} else if err.Error() != test.expectedError {
					t.Errorf("expected error: %s, got: %s", test.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %s", err)
				}
				if !reflect.DeepEqual(result, test.expected) {
					t.Errorf("expected result: %+v, got: %+v", test.expected, result)
				}
			}
		})
	}
}

func TestPopulateLANFields(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		expected      *goaviatrix.EdgeSpokeInterface
		expectErr     bool
		expectedError string
	}{
		{
			name: "Valid input",
			input: map[string]interface{}{
				"enable_vrrp":     true,
				"vrrp_virtual_ip": "192.168.1.100",
			},
			expected: &goaviatrix.EdgeSpokeInterface{
				VrrpState: true,
				VirtualIp: "192.168.1.100",
			},
			expectErr: false,
		},
		{
			name: "Missing enable_vrrp field",
			input: map[string]interface{}{
				"vrrp_virtual_ip": "192.168.1.100",
			},
			expectErr:     true,
			expectedError: "invalid type for 'enable_vrrp': expected bool, got <nil>",
		},
		{
			name: "Missing vrrp_virtual_ip field",
			input: map[string]interface{}{
				"enable_vrrp": true,
			},
			expectErr:     true,
			expectedError: "invalid type for 'vrrp_virtual_ip': expected string, got <nil>",
		},
		{
			name: "Invalid type for enable_vrrp",
			input: map[string]interface{}{
				"enable_vrrp":     "true",
				"vrrp_virtual_ip": "192.168.1.100",
			},
			expectErr:     true,
			expectedError: "invalid type for 'enable_vrrp': expected bool, got string",
		},
		{
			name: "Invalid type for vrrp_virtual_ip",
			input: map[string]interface{}{
				"enable_vrrp":     true,
				"vrrp_virtual_ip": 12345,
			},
			expectErr:     true,
			expectedError: "invalid type for 'vrrp_virtual_ip': expected string, got int",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if2 := &goaviatrix.EdgeSpokeInterface{}
			err := populateLANFields(test.input, if2)

			if test.expectErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				} else if err.Error() != test.expectedError {
					t.Errorf("expected error: %s, got: %s", test.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %s", err)
				}
				if !reflect.DeepEqual(if2, test.expected) {
					t.Errorf("expected result: %+v, got: %+v", test.expected, if2)
				}
			}
		})
	}
}

func TestBuildEdgeSpokeVlan(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		expected      *goaviatrix.EdgeSpokeVlan
		expectErr     bool
		expectedError string
	}{
		{
			name: "Valid input",
			input: map[string]interface{}{
				"parent_interface_name": "eth0",
				"ip_address":            "192.168.1.2",
				"gateway_ip":            "192.168.1.254",
				"peer_ip_address":       "192.168.1.3",
				"peer_gateway_ip":       "192.168.1.253",
				"vrrp_virtual_ip":       "192.168.1.100",
				"tag":                   "tag1",
				"vlan_id":               100,
			},
			expected: &goaviatrix.EdgeSpokeVlan{
				ParentInterface: "eth0",
				IpAddr:          "192.168.1.2",
				GatewayIp:       "192.168.1.254",
				PeerIpAddr:      "192.168.1.3",
				PeerGatewayIp:   "192.168.1.253",
				VirtualIp:       "192.168.1.100",
				Tag:             "tag1",
				VlanId:          "100",
			},
			expectErr: false,
		},
		{
			name: "Missing parent_interface_name",
			input: map[string]interface{}{
				"ip_address":      "192.168.1.2",
				"gateway_ip":      "192.168.1.254",
				"peer_ip_address": "192.168.1.3",
				"peer_gateway_ip": "192.168.1.253",
				"vrrp_virtual_ip": "192.168.1.100",
				"tag":             "tag1",
				"vlan_id":         100,
			},
			expectErr:     true,
			expectedError: "invalid type for 'parent_interface_name': expected string, got <nil>",
		},
		{
			name: "Invalid type for vlan_id",
			input: map[string]interface{}{
				"parent_interface_name": "eth0",
				"ip_address":            "192.168.1.2",
				"gateway_ip":            "192.168.1.254",
				"peer_ip_address":       "192.168.1.3",
				"peer_gateway_ip":       "192.168.1.253",
				"vrrp_virtual_ip":       "192.168.1.100",
				"tag":                   "tag1",
				"vlan_id":               "invalid",
			},
			expectErr:     true,
			expectedError: "invalid type for 'vlan_id': expected int, got string",
		},
		{
			name: "Missing ip_address",
			input: map[string]interface{}{
				"parent_interface_name": "eth0",
				"gateway_ip":            "192.168.1.254",
				"peer_ip_address":       "192.168.1.3",
				"peer_gateway_ip":       "192.168.1.253",
				"vrrp_virtual_ip":       "192.168.1.100",
				"tag":                   "tag1",
				"vlan_id":               100,
			},
			expectErr:     true,
			expectedError: "invalid type for 'ip_address': expected string, got <nil>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := buildEdgeSpokeVlan(test.input)

			if test.expectErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				} else if err.Error() != test.expectedError {
					t.Errorf("expected error: %s, got: %s", test.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %s", err)
				}
				if !reflect.DeepEqual(result, test.expected) {
					t.Errorf("expected result: %+v, got: %+v", test.expected, result)
				}
			}
		})
	}
}
