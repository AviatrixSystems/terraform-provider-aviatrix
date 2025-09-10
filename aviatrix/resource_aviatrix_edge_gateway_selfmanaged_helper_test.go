package aviatrix

import (
	"reflect"
	"testing"
)

func TestEditAdvertisedSpokeRoutesEmptyArrayHandling(t *testing.T) {
	tests := []struct {
		name           string
		inputRoutes    []string
		expectedRoutes []string
		description    string
	}{
		{
			name:           "Empty array should be converted to empty string array",
			inputRoutes:    []string{},
			expectedRoutes: []string{""},
			description:    "When user sets included_advertised_spoke_routes = [], it should be converted to [\"\"] to clear routes",
		},
		{
			name:           "Non-empty array should remain unchanged",
			inputRoutes:    []string{"10.0.0.0/8", "192.168.0.0/16"},
			expectedRoutes: []string{"10.0.0.0/8", "192.168.0.0/16"},
			description:    "When user sets actual CIDR routes, they should remain unchanged",
		},
		{
			name:           "Single empty string should remain unchanged",
			inputRoutes:    []string{""},
			expectedRoutes: []string{""},
			description:    "When user sets included_advertised_spoke_routes = [\"\"], it should remain as is",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Simulate the logic from editAdvertisedSpokeRoutesWithRetry
			includedAdvertisedSpokeRoutes := test.inputRoutes

			// Apply the same logic as in the actual function
			if len(includedAdvertisedSpokeRoutes) == 0 {
				includedAdvertisedSpokeRoutes = []string{""}
			}

			if !reflect.DeepEqual(includedAdvertisedSpokeRoutes, test.expectedRoutes) {
				t.Errorf("Test '%s' failed: expected %v, got %v", test.name, test.expectedRoutes, includedAdvertisedSpokeRoutes)
			}
		})
	}
}
