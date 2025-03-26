package aviatrix

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestGetDstWanInterfaces(t *testing.T) {
	tests := []struct {
		name                 string
		logicalIfNames       []string
		gatewayDetails       *goaviatrix.Gateway
		reversedIfNames      map[string]string
		setWanInterfacesResp string
		setWanInterfacesErr  error
		expectedResult       string
		expectErr            bool
	}{
		{
			name:           "Successful case",
			logicalIfNames: []string{"eth0", "eth1"},
			gatewayDetails: &goaviatrix.Gateway{
				IfNamesTranslation: map[string]string{
					"wan1": "eth0",
					"wan2": "eth1",
				},
			},
			reversedIfNames: map[string]string{
				"eth0": "wan1",
				"eth1": "wan2",
			},
			setWanInterfacesResp: "wan1,wan2",
			setWanInterfacesErr:  nil,
			expectedResult:       "wan1,wan2",
			expectErr:            false,
		},
		{
			name:           "SetWanInterfaces returns an error",
			logicalIfNames: []string{"eth0", "eth2"},
			gatewayDetails: &goaviatrix.Gateway{
				IfNamesTranslation: map[string]string{
					"wan1": "eth0",
					"wan2": "eth1",
				},
			},
			reversedIfNames: map[string]string{
				"eth0": "wan1",
				"eth1": "wan2",
			},
			setWanInterfacesResp: "",
			setWanInterfacesErr:  errors.New("invalid interface"),
			expectedResult:       "",
			expectErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getDstWanInterfaces(tt.logicalIfNames, tt.gatewayDetails)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetEdgeTransitLogicalIfNames(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		attachment    *goaviatrix.SpokeTransitAttachment
		gateway       *goaviatrix.Gateway
		expectErr     bool
		expectedError string
		expectedIfs   []string
		expectedDst   string
	}{
		{
			name: "Valid Edge Gateway with Logical Interfaces",
			resourceData: map[string]interface{}{
				"transit_gateway_logical_ifnames": []interface{}{"wan0", "wan1"},
			},
			attachment: &goaviatrix.SpokeTransitAttachment{},
			gateway: &goaviatrix.Gateway{
				CloudType:          goaviatrix.EDGEEQUINIX,
				IfNamesTranslation: map[string]string{"eth0": "wan.0", "eth1": "wan.1"},
			},
			expectErr:   false,
			expectedIfs: []string{"wan0", "wan1"},
			expectedDst: "eth0,eth1",
		},
		{
			name:         "Missing transit_gateway_logical_ifnames",
			resourceData: map[string]interface{}{
				// No "transit_gateway_logical_ifnames"
			},
			attachment: &goaviatrix.SpokeTransitAttachment{},
			gateway: &goaviatrix.Gateway{
				CloudType: goaviatrix.EdgeRelatedCloudTypes,
			},
			expectErr:     true,
			expectedError: "transit_gateway_logical_ifnames is required for all edge gateways",
		},
		{
			name: "Invalid transit_gateway_logical_ifnames type",
			resourceData: map[string]interface{}{
				"transit_gateway_logical_ifnames": "not-a-list",
			},
			attachment: &goaviatrix.SpokeTransitAttachment{},
			gateway: &goaviatrix.Gateway{
				CloudType: goaviatrix.EdgeRelatedCloudTypes,
			},
			expectErr:     true,
			expectedError: "transit_gateway_logical_ifnames is required for all edge gateways",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"transit_gateway_logical_ifnames": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			}, tt.resourceData)

			err := getEdgeTransitLogicalIfNames(d, tt.gateway, tt.attachment)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIfs, tt.attachment.TransitGatewayLogicalIfNames)
				assert.Equal(t, tt.expectedDst, tt.attachment.DstWanInterfaces)
			}
		})
	}
}

func TestSetCustomInterfaceMapping(t *testing.T) {
	tests := []struct {
		name                     string
		customInterfaceMap       map[string][]goaviatrix.CustomInterfaceMap
		userCustomInterfaceOrder []string
		expected                 []interface{}
		expectErr                bool
		expectedError            string
	}{
		{
			name: "Valid input with correct order",
			customInterfaceMap: map[string][]goaviatrix.CustomInterfaceMap{
				"wan0": {
					{
						IdentifierType:  "mac",
						IdentifierValue: "00:1A:2B:3C:4D:5E",
					},
				},
				"mgmt0": {
					{
						IdentifierType:  "pci",
						IdentifierValue: "0000:00:1f.2",
					},
				},
			},
			userCustomInterfaceOrder: []string{"mgmt0", "wan0"},
			expected: []interface{}{
				map[string]interface{}{
					"logical_ifname":   "mgmt0",
					"identifier_type":  "pci",
					"idenitifer_value": "0000:00:1f.2",
				},
				map[string]interface{}{
					"logical_ifname":   "wan0",
					"identifier_type":  "mac",
					"idenitifer_value": "00:1A:2B:3C:4D:5E",
				},
			},
			expectErr: false,
		},
		{
			name: "Logical interface name not found in custom interface map",
			customInterfaceMap: map[string][]goaviatrix.CustomInterfaceMap{
				"wan0": {
					{
						IdentifierType:  "mac",
						IdentifierValue: "00:1A:2B:3C:4D:5E",
					},
				},
			},
			userCustomInterfaceOrder: []string{"mgmt0"},
			expectErr:                true,
			expectedError:            "logical interface name mgmt0 not found in custom interface map",
		},
		{
			name: "Empty identifier type",
			customInterfaceMap: map[string][]goaviatrix.CustomInterfaceMap{
				"wan0": {
					{
						IdentifierType:  "",
						IdentifierValue: "00:1A:2B:3C:4D:5E",
					},
				},
			},
			userCustomInterfaceOrder: []string{"wan0"},
			expectErr:                true,
			expectedError:            "identifier type cannot be empty for logical interface: wan0",
		},
		{
			name: "Empty identifier value",
			customInterfaceMap: map[string][]goaviatrix.CustomInterfaceMap{
				"wan0": {
					{
						IdentifierType:  "mac",
						IdentifierValue: "",
					},
				},
			},
			userCustomInterfaceOrder: []string{"wan0"},
			expectErr:                true,
			expectedError:            "identifier value cannot be empty for logical interface: wan0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := setCustomInterfaceMapping(test.customInterfaceMap, test.userCustomInterfaceOrder)

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

func TestGetCustomInterfaceOrder(t *testing.T) {
	tests := []struct {
		name                       string
		userCustomInterfaceMapping []interface{}
		expected                   []string
		expectErr                  bool
		expectedError              string
	}{
		{
			name: "Valid input",
			userCustomInterfaceMapping: []interface{}{
				map[string]interface{}{
					"logical_ifname": "wan0",
				},
				map[string]interface{}{
					"logical_ifname": "mgmt0",
				},
			},
			expected:  []string{"wan0", "mgmt0"},
			expectErr: false,
		},
		{
			name: "Invalid type in mapping",
			userCustomInterfaceMapping: []interface{}{
				"invalid_type",
			},
			expectErr:     true,
			expectedError: "invalid type: expected map[string]interface{}, got string",
		},
		{
			name: "Missing logical_ifname",
			userCustomInterfaceMapping: []interface{}{
				map[string]interface{}{
					"identifier_type": "mac",
				},
			},
			expectErr:     true,
			expectedError: "logical_ifname must be a non-empty string",
		},
		{
			name: "Empty logical_ifname",
			userCustomInterfaceMapping: []interface{}{
				map[string]interface{}{
					"logical_ifname": "",
				},
			},
			expectErr:     true,
			expectedError: "logical_ifname must be a non-empty string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := getCustomInterfaceOrder(test.userCustomInterfaceMapping)

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
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", test.expected) {
					t.Errorf("expected result: %v, got: %v", test.expected, result)
				}
			}
		})
	}
}
