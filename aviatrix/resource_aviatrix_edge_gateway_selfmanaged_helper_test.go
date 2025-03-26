package aviatrix

import (
	"reflect"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func TestValidateIdentifierValue(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
	}{
		{
			name:      "Valid auto value",
			input:     "auto",
			expectErr: false,
		},
		{
			name:      "Valid MAC address",
			input:     "00:1A:2B:3C:4D:5E",
			expectErr: false,
		},
		{
			name:      "Valid PCI ID",
			input:     "0000:00:1f.2",
			expectErr: false,
		},
		{
			name:      "Invalid MAC address",
			input:     "00:1A:2B:3C:4D",
			expectErr: true,
		},
		{
			name:      "Invalid PCI ID",
			input:     "0000:00:1f",
			expectErr: true,
		},
		{
			name:      "Invalid random string",
			input:     "invalid_value",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, errs := validateIdentifierValue(test.input, "identifier_value")
			if test.expectErr && len(errs) == 0 {
				t.Errorf("expected an error but got none for input: %v", test.input)
			}
			if !test.expectErr && len(errs) > 0 {
				t.Errorf("did not expect an error but got: %v for input: %v", errs, test.input)
			}
		})
	}
}

func TestGetCustomInterfaceMapDetails(t *testing.T) {
	tests := []struct {
		name          string
		input         []interface{}
		expected      map[string][]goaviatrix.CustomInterfaceMap
		expectErr     bool
		expectedError string
	}{
		{
			name: "Valid input",
			input: []interface{}{
				map[string]interface{}{
					"logical_ifname":   "wan0",
					"identifier_type":  "mac",
					"idenitifer_value": "00:1A:2B:3C:4D:5E",
				},
				map[string]interface{}{
					"logical_ifname":   "mgmt0",
					"identifier_type":  "pci",
					"idenitifer_value": "0000:00:1f.2",
				},
			},
			expected: map[string][]goaviatrix.CustomInterfaceMap{
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
			expectErr: false,
		},
		{
			name: "Invalid input type",
			input: []interface{}{
				"invalid_type",
			},
			expectErr:     true,
			expectedError: "invalid type: expected map[string]interface{}, got string",
		},
		{
			name: "Missing logical_ifname",
			input: []interface{}{
				map[string]interface{}{
					"identifier_type":  "mac",
					"idenitifer_value": "00:1A:2B:3C:4D:5E",
				},
			},
			expectErr:     true,
			expectedError: "logical interface name must be a string",
		},
		{
			name: "Missing identifier_type",
			input: []interface{}{
				map[string]interface{}{
					"logical_ifname":   "wan0",
					"idenitifer_value": "00:1A:2B:3C:4D:5E",
				},
			},
			expectErr:     true,
			expectedError: "identifier type must be a string",
		},
		{
			name: "Missing identifier_value",
			input: []interface{}{
				map[string]interface{}{
					"logical_ifname":  "wan0",
					"identifier_type": "mac",
				},
			},
			expectErr:     true,
			expectedError: "identifier value must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := getCustomInterfaceMapDetails(test.input)

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
