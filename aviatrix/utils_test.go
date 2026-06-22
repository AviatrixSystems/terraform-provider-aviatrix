package aviatrix

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestSortSpokeInterfacesByCustomOrder(t *testing.T) {
	// Define test cases
	tests := []struct {
		name               string
		interfaces         []goaviatrix.MegaportInterface
		userInterfaceOrder []string
		expected           []goaviatrix.MegaportInterface
	}{
		{
			name: "Sort interfaces based on custom order",
			interfaces: []goaviatrix.MegaportInterface{
				{LogicalInterfaceName: "wan2"},
				{LogicalInterfaceName: "wan0"},
				{LogicalInterfaceName: "wan1"},
			},
			userInterfaceOrder: []string{"wan0", "wan1", "wan2"},
			expected: []goaviatrix.MegaportInterface{
				{LogicalInterfaceName: "wan0"},
				{LogicalInterfaceName: "wan1"},
				{LogicalInterfaceName: "wan2"},
			},
		},
		{
			name: "Unordered interfaces with missing custom order",
			interfaces: []goaviatrix.MegaportInterface{
				{LogicalInterfaceName: "wan3"},
				{LogicalInterfaceName: "wan1"},
				{LogicalInterfaceName: "wan2"},
			},
			userInterfaceOrder: []string{"wan1", "wan2"},
			expected: []goaviatrix.MegaportInterface{
				{LogicalInterfaceName: "wan1"},
				{LogicalInterfaceName: "wan2"},
				{LogicalInterfaceName: "wan3"},
			},
		},
		{
			name: "Empty custom order",
			interfaces: []goaviatrix.MegaportInterface{
				{LogicalInterfaceName: "wan1"},
				{LogicalInterfaceName: "wan2"},
			},
			userInterfaceOrder: []string{},
			expected: []goaviatrix.MegaportInterface{
				{LogicalInterfaceName: "wan1"},
				{LogicalInterfaceName: "wan2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortSpokeInterfacesByCustomOrder(tt.interfaces, tt.userInterfaceOrder)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("sortSpokeInterfacesByCustomOrder() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateCIDRRule(t *testing.T) {
	testCases := []struct {
		name          string
		rule          string
		expectedError bool
		errorContains string
	}{
		// Valid cases
		{
			name:          "simple CIDR",
			rule:          "10.1.0.0/16",
			expectedError: false,
		},
		{
			name:          "CIDR with ge equal to prefix",
			rule:          "10.1.0.0/16 ge 16",
			expectedError: false,
		},
		{
			name:          "CIDR with ge greater than prefix",
			rule:          "10.1.0.0/16 ge 24",
			expectedError: false,
		},
		{
			name:          "CIDR with le only",
			rule:          "10.1.0.0/16 le 24",
			expectedError: false,
		},
		{
			name:          "CIDR with ge and le",
			rule:          "10.0.0.0/8 ge 16 le 24",
			expectedError: false,
		},
		{
			name:          "CIDR with equal ge and le",
			rule:          "10.1.0.0/16 ge 24 le 24",
			expectedError: false,
		},

		// Invalid cases
		{
			name:          "empty string",
			rule:          "",
			expectedError: true,
			errorContains: "invalid number of fields",
		},
		{
			name:          "invalid CIDR format",
			rule:          "10.1.0.0.0/16",
			expectedError: true,
			errorContains: "invalid IPv4 CIDR",
		},
		{
			name:          "CIDR with invalid prefix",
			rule:          "10.1.0.0/33",
			expectedError: true,
			errorContains: "invalid IPv4 CIDR",
		},
		{
			name:          "incorrect number of parts",
			rule:          "10.1.0.0/16 ge",
			expectedError: true,
			errorContains: "invalid number of fields",
		},
		{
			name:          "unknown qualifier",
			rule:          "10.1.0.0/16 gt 24",
			expectedError: true,
			errorContains: "unknown qualifier",
		},
		{
			name:          "non-numeric value for ge",
			rule:          "10.1.0.0/16 ge abc",
			expectedError: true,
			errorContains: "invalid value \"abc\"",
		},
		{
			name:          "duplicate ge qualifier",
			rule:          "10.1.0.0/16 ge 16 ge 24",
			expectedError: true,
			errorContains: "duplicate 'ge' qualifier",
		},
		{
			name:          "duplicate le qualifier",
			rule:          "10.1.0.0/16 le 16 le 24",
			expectedError: true,
			errorContains: "duplicate 'le' qualifier",
		},
		{
			name:          "ge less than prefix",
			rule:          "10.1.0.0/24 ge 16",
			expectedError: true,
			errorContains: "length 16 out of range",
		},
		{
			name:          "ge greater than 32",
			rule:          "10.1.0.0/16 ge 33",
			expectedError: true,
			errorContains: "length 33 out of range",
		},
		{
			name:          "wrong order of qualifiers",
			rule:          "10.1.0.0/16 le 24 ge 16",
			expectedError: true,
			errorContains: "'ge' must come before 'le'",
		},
		{
			name:          "non-numeric value for le",
			rule:          "10.1.0.0/16 ge 24 le xyz",
			expectedError: true,
			errorContains: "invalid value \"xyz\"",
		},
		{
			name:          "le less than ge",
			rule:          "10.1.0.0/16 ge 24 le 20",
			expectedError: true,
			errorContains: "ge length 24 > le length 20",
		},
		{
			name:          "too many parts",
			rule:          "10.1.0.0/16 ge 24 le 28 extra part",
			expectedError: true,
			errorContains: "invalid number of fields",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			warnings, errors := ValidateCIDRRule(tc.rule, "test_property")

			// Warnings should always be empty in our implementation
			assert.Empty(t, warnings, "Expected no warnings")

			if tc.expectedError {
				assert.NotEmpty(t, errors, "Expected validation errors, but got none")
				if tc.errorContains != "" {
					errorFound := false
					for _, err := range errors {
						if assert.Error(t, err) {
							if contains(err.Error(), tc.errorContains) {
								errorFound = true
								break
							}
						}
					}
					assert.True(t, errorFound, "Expected error to contain '%s', but it didn't", tc.errorContains)
				}
			} else {
				assert.Empty(t, errors, "Expected no validation errors, but got: %v", errors)
			}
		})
	}
}

func TestValidateIPv6CIDR(t *testing.T) {
	testCases := []struct {
		name          string
		input         interface{}
		key           string
		expectedError bool
		errorContains string
	}{
		// Valid IPv6 CIDR cases
		{
			name:          "valid IPv6 CIDR /64",
			input:         "2001:db8::/64",
			key:           "test_property",
			expectedError: false,
		},
		{
			name:          "valid IPv6 CIDR /128",
			input:         "2001:db8::1/128",
			key:           "test_property",
			expectedError: false,
		},
		{
			name:          "valid loopback IPv6 CIDR",
			input:         "::1/128",
			key:           "test_property",
			expectedError: false,
		},

		// Invalid input type cases
		{
			name:          "non-string input",
			input:         123,
			key:           "test_property",
			expectedError: true,
			errorContains: "expected type of \"test_property\" to be string",
		},

		// Invalid CIDR format cases
		{
			name:          "invalid CIDR format - no prefix",
			input:         "2001:db8::",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain a valid IPv6 CIDR",
		},
		{
			name:          "invalid CIDR format - malformed IPv6",
			input:         "2001:db8::xyz/64",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain a valid IPv6 CIDR",
		},

		// IPv4 CIDR cases (should be rejected)
		{
			name:          "IPv4 CIDR /24",
			input:         "192.168.1.0/24",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain an IPv6 CIDR, got IPv4",
		},
		{
			name:          "IPv4-mapped IPv6 address",
			input:         "::ffff:192.168.1.1/128",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain an IPv6 CIDR, got IPv4",
		},

		// Host address instead of network address (should be rejected)
		{
			name:          "host address in CIDR /64",
			input:         "2001:db8::5/64",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain a network CIDR, got host address",
		},
		{
			name:          "host address in CIDR /48",
			input:         "2001:db8:1::1/48",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain a network CIDR, got host address",
		},
		{
			name:          "non-zero host bits in CIDR",
			input:         "fd00:1234:5678:abcd::/56",
			key:           "test_property",
			expectedError: true,
			errorContains: "expected test_property to contain a network CIDR, got host address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			warnings, errors := validateIPv6CIDR(tc.input, tc.key)

			// Warnings should always be empty in our implementation
			assert.Empty(t, warnings, "Expected no warnings")

			if tc.expectedError {
				assert.NotEmpty(t, errors, "Expected validation errors, but got none")
				if tc.errorContains != "" {
					errorFound := false
					for _, err := range errors {
						if assert.Error(t, err) {
							if contains(err.Error(), tc.errorContains) {
								errorFound = true
								break
							}
						}
					}
					assert.True(t, errorFound, "Expected error to contain '%s', but got errors: %v", tc.errorContains, errors)
				}
			} else {
				assert.Empty(t, errors, "Expected no validation errors, but got: %v", errors)
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
