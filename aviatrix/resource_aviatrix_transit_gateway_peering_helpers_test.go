package aviatrix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestReverseIfnameTranslation(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "Basic translation reversal",
			input: map[string]string{
				"eth0": "wan.0",
				"eth1": "wan.1",
			},
			expected: map[string]string{
				"wan0": "eth0",
				"wan1": "eth1",
			},
		},
		{
			name:     "Handles empty map",
			input:    map[string]string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseIfnameTranslation(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Test %s failed. Expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSetWanInterfaces(t *testing.T) {
	tests := []struct {
		name                   string
		logicalIfNames         []string
		reversedInterfaceNames map[string]string
		expected               string
		expectErr              bool
	}{
		{
			name:           "Valid input - single interface",
			logicalIfNames: []string{"wan0"},
			reversedInterfaceNames: map[string]string{
				"wan0": "eth0",
			},
			expected:  "eth0",
			expectErr: false,
		},
		{
			name:           "Valid input - multiple interfaces",
			logicalIfNames: []string{"wan0", "wan1"},
			reversedInterfaceNames: map[string]string{
				"wan0": "eth0",
				"wan1": "eth1",
			},
			expected:  "eth0,eth1",
			expectErr: false,
		},
		{
			name:           "Interface name not found in map",
			logicalIfNames: []string{"wan2"},
			reversedInterfaceNames: map[string]string{
				"wan0": "eth0",
			},
			expected:  "",
			expectErr: true,
		},
		{
			name:           "Empty logicalIfNames",
			logicalIfNames: []string{},
			reversedInterfaceNames: map[string]string{
				"wan0": "eth0",
			},
			expected:  "",
			expectErr: false, // Empty input should return an empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SetWanInterfaces(tt.logicalIfNames, tt.reversedInterfaceNames)
			if (err != nil) != tt.expectErr {
				t.Errorf("Test %s failed: expected error=%v, got error=%v", tt.name, tt.expectErr, err != nil)
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Test %s failed: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSetWanInterfaceNames(t *testing.T) {
	tests := []struct {
		name           string
		logicalIfNames []string
		cloudType      int
		gatewayDetails *goaviatrix.Gateway
		gatewayPrefix  string
		expectedError  bool
	}{
		{
			name:           "Valid logical interfaces for Equinix cloud type",
			logicalIfNames: []string{"wan0", "wan1"},
			cloudType:      goaviatrix.EDGEEQUINIX,
			gatewayDetails: &goaviatrix.Gateway{IfNamesTranslation: map[string]string{"eth0": "wan.0", "eth1": "wan.1"}},
			gatewayPrefix:  "gateway1",
			expectedError:  false,
		},
		{
			name:           "Valid logical interfaces for Megaport cloud type",
			logicalIfNames: []string{"wan0", "wan1"},
			cloudType:      goaviatrix.EDGEMEGAPORT,
			gatewayDetails: &goaviatrix.Gateway{},
			gatewayPrefix:  "gateway1",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitGatewayPeering := &goaviatrix.TransitGatewayPeering{}

			// Call the function with actual cloud type
			err := setWanInterfaceNames(tt.logicalIfNames, tt.cloudType, tt.gatewayDetails, tt.gatewayPrefix, transitGatewayPeering)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.cloudType == goaviatrix.EDGEEQUINIX || tt.cloudType == goaviatrix.EDGENEO {
					assert.NotEmpty(t, transitGatewayPeering.SrcWanInterfaces)
				} else if tt.cloudType == goaviatrix.EDGEMEGAPORT {
					// Check the logical interfaces set for Megaport
					if tt.gatewayPrefix == "gateway1" {
						assert.Equal(t, tt.logicalIfNames, transitGatewayPeering.Gateway1LogicalIfNames)
					} else {
						assert.Equal(t, tt.logicalIfNames, transitGatewayPeering.Gateway2LogicalIfNames)
					}
				}
			}
		})
	}
}

func TestSetExcludedResources(t *testing.T) {
	tests := []struct {
		name                  string
		resourceData          map[string]interface{}
		expectedGateway1CIDRs string
		expectedGateway2CIDRs string
		expectedGateway1TGWs  string
		expectedGateway2TGWs  string
		expectError           bool
	}{
		{
			name: "Valid excluded resources",
			resourceData: map[string]interface{}{
				"gateway1_excluded_network_cidrs":   []interface{}{"192.168.1.0/24", "192.168.2.0/24"},
				"gateway2_excluded_network_cidrs":   []interface{}{"10.0.1.0/24"},
				"gateway1_excluded_tgw_connections": []interface{}{"tgw-123", "tgw-456"},
				"gateway2_excluded_tgw_connections": []interface{}{"tgw-789"},
			},
			expectedGateway1CIDRs: "192.168.1.0/24,192.168.2.0/24",
			expectedGateway2CIDRs: "10.0.1.0/24",
			expectedGateway1TGWs:  "tgw-123,tgw-456",
			expectedGateway2TGWs:  "tgw-789",
			expectError:           false,
		},
		{
			name: "Empty excluded resources",
			resourceData: map[string]interface{}{
				"gateway1_excluded_network_cidrs":   []interface{}{},
				"gateway2_excluded_network_cidrs":   []interface{}{},
				"gateway1_excluded_tgw_connections": []interface{}{},
				"gateway2_excluded_tgw_connections": []interface{}{},
			},
			expectedGateway1CIDRs: "",
			expectedGateway2CIDRs: "",
			expectedGateway1TGWs:  "",
			expectedGateway2TGWs:  "",
			expectError:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"gateway1_excluded_network_cidrs":   {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
				"gateway2_excluded_network_cidrs":   {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
				"gateway1_excluded_tgw_connections": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
				"gateway2_excluded_tgw_connections": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			}, tt.resourceData)

			transitGatewayPeering := &goaviatrix.TransitGatewayPeering{}
			err := setExcludedResources(d, transitGatewayPeering)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedGateway1CIDRs, transitGatewayPeering.Gateway1ExcludedCIDRs)
				assert.Equal(t, tt.expectedGateway2CIDRs, transitGatewayPeering.Gateway2ExcludedCIDRs)
				assert.Equal(t, tt.expectedGateway1TGWs, transitGatewayPeering.Gateway1ExcludedTGWConnections)
				assert.Equal(t, tt.expectedGateway2TGWs, transitGatewayPeering.Gateway2ExcludedTGWConnections)
			}
		})
	}
}

func TestGetBooleanValue(t *testing.T) {
	tests := []struct {
		name          string
		inputData     map[string]interface{}
		key           string
		expectedValue bool
		expectedError bool
	}{
		{
			name:          "Valid boolean value (true)",
			inputData:     map[string]interface{}{"enable_feature": true},
			key:           "enable_feature",
			expectedValue: true,
			expectedError: false,
		},
		{
			name:          "Valid boolean value (false)",
			inputData:     map[string]interface{}{"enable_feature": false},
			key:           "enable_feature",
			expectedValue: false,
			expectedError: false,
		},
		{
			name:          "Missing key (should return error)",
			inputData:     map[string]interface{}{},
			key:           "non_existent_key",
			expectedValue: false,
			expectedError: true,
		},
		{
			name:          "Invalid type (string instead of bool)",
			inputData:     map[string]interface{}{"invalid_type": "not_a_boolean"},
			key:           "invalid_type",
			expectedValue: false,
			expectedError: true,
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"enable_feature": {Type: schema.TypeBool, Optional: true},
		"invalid_type":   {Type: schema.TypeString, Optional: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, tt.inputData)
			val, err := getBooleanValue(d, tt.key)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, val)
			}
		})
	}
}

func TestValidateTunnelCount(t *testing.T) {
	tests := []struct {
		name          string
		insaneMode    bool
		tunnelCount   int
		expectedError bool
	}{
		{
			name:          "Valid tunnel count when insaneMode is true and tunnelCount is 1",
			insaneMode:    true,
			tunnelCount:   1,
			expectedError: false,
		},
		{
			name:          "Invalid tunnel count when insaneMode is true and tunnelCount is 0",
			insaneMode:    true,
			tunnelCount:   0,
			expectedError: true,
		},
		{
			name:          "Valid tunnel count when insaneMode is false and tunnelCount is 0",
			insaneMode:    false,
			tunnelCount:   0,
			expectedError: false,
		},
		{
			name:          "Invalid tunnel count when insaneMode is false and tunnelCount is greater than 0",
			insaneMode:    false,
			tunnelCount:   1,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTunnelCount(tt.insaneMode, tt.tunnelCount)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetNonEATPeeringOptions(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		expected      map[string]bool
		expectErr     bool
		expectedError string
	}{
		{
			name: "Valid input",
			resourceData: map[string]interface{}{
				"enable_max_performance":                      true,
				"enable_insane_mode_encryption_over_internet": false,
				"enable_peering_over_private_network":         true,
				"enable_single_tunnel_mode":                   false,
			},
			expected: map[string]bool{
				"enable_max_performance":              true,
				"insane_mode":                         false,
				"enable_peering_over_private_network": true,
				"enable_single_tunnel_mode":           false,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_max_performance":                      {Type: schema.TypeBool},
				"enable_insane_mode_encryption_over_internet": {Type: schema.TypeBool},
				"enable_peering_over_private_network":         {Type: schema.TypeBool},
				"enable_single_tunnel_mode":                   {Type: schema.TypeBool},
			}, tt.resourceData)

			result, err := getNonEATPeeringOptions(d)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSetNonEATPeeringOptions(t *testing.T) {
	tests := []struct {
		name                           string
		resourceData                   map[string]interface{}
		expectedError                  bool
		expectedNoMaxPerformance       bool
		expectedInsaneModeOverInternet bool
		expectedTunnelCount            int
		expectedPrivateIPPeering       string
		expectedSingleTunnelMode       string
	}{
		{
			name: "Valid Max Performance and Insane Mode",
			resourceData: map[string]interface{}{
				"enable_max_performance":                      true,
				"enable_insane_mode_encryption_over_internet": true,
				"enable_peering_over_private_network":         false,
				"enable_single_tunnel_mode":                   false,
				"tunnel_count":                                2,
			},
			expectedError:                  false,
			expectedNoMaxPerformance:       false,
			expectedInsaneModeOverInternet: true,
			expectedTunnelCount:            2,
			expectedPrivateIPPeering:       "no",
			expectedSingleTunnelMode:       "",
		},
		{
			name: "Insane Mode with Peering Over Private Network Error",
			resourceData: map[string]interface{}{
				"enable_max_performance":                      true,
				"enable_insane_mode_encryption_over_internet": true,
				"enable_peering_over_private_network":         true,
				"enable_single_tunnel_mode":                   false,
				"tunnel_count":                                2,
			},
			expectedError:                  true,
			expectedNoMaxPerformance:       false,
			expectedInsaneModeOverInternet: true,
			expectedTunnelCount:            2,
			expectedPrivateIPPeering:       "yes",
			expectedSingleTunnelMode:       "",
		},
		{
			name: "Valid Single Tunnel Mode with Private Peering",
			resourceData: map[string]interface{}{
				"enable_max_performance":                      false,
				"enable_insane_mode_encryption_over_internet": false,
				"enable_peering_over_private_network":         true,
				"enable_single_tunnel_mode":                   true,
				"tunnel_count":                                0,
			},
			expectedError:                  false,
			expectedNoMaxPerformance:       true,
			expectedInsaneModeOverInternet: false,
			expectedTunnelCount:            0,
			expectedPrivateIPPeering:       "yes",
			expectedSingleTunnelMode:       "yes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_max_performance":                      {Type: schema.TypeBool},
				"enable_insane_mode_encryption_over_internet": {Type: schema.TypeBool},
				"enable_peering_over_private_network":         {Type: schema.TypeBool},
				"enable_single_tunnel_mode":                   {Type: schema.TypeBool},
				"tunnel_count":                                {Type: schema.TypeInt},
			}, tt.resourceData)

			transitGatewayPeering := &goaviatrix.TransitGatewayPeering{}
			err := setNonEATPeeringOptions(d, transitGatewayPeering)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedNoMaxPerformance, transitGatewayPeering.NoMaxPerformance)
			assert.Equal(t, tt.expectedInsaneModeOverInternet, transitGatewayPeering.InsaneModeOverInternet)
			assert.Equal(t, tt.expectedTunnelCount, transitGatewayPeering.TunnelCount)
			assert.Equal(t, tt.expectedPrivateIPPeering, transitGatewayPeering.PrivateIPPeering)
			assert.Equal(t, tt.expectedSingleTunnelMode, transitGatewayPeering.SingleTunnel)
		})
	}
}

func TestGetStringListFromResource(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		key           string
		expectedList  []string
		expectedError error
	}{
		{
			name: "Valid string list",
			resourceData: map[string]interface{}{
				"valid_list": []interface{}{"one", "two", "three"},
			},
			key:           "valid_list",
			expectedList:  []string{"one", "two", "three"},
			expectedError: nil,
		},
		{
			name:          "Key not present",
			resourceData:  map[string]interface{}{},
			key:           "missing_key",
			expectedList:  nil,
			expectedError: nil,
		},
		{
			name: "Key exists but is not a list",
			resourceData: map[string]interface{}{
				"not_a_list": "string_value",
			},
			key:           "not_a_list",
			expectedList:  nil,
			expectedError: fmt.Errorf("not_a_list is not a list of strings"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"valid_list": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
				"not_a_list": {Type: schema.TypeString},
				"mixed_list": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			}, tt.resourceData)

			result, err := getStringListFromResource(d, tt.key)
			assert.Equal(t, tt.expectedList, result)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
