package aviatrix

import (
	"reflect"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
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
