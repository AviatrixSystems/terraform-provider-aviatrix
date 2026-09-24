package aviatrix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestConfigureTransitInstanceEncryption(t *testing.T) {
	tests := []struct {
		name        string
		cloudType   int
		raw         map[string]any
		expectErr   bool
		expectedEnc string
		expectedCmk string
	}{
		{
			name:      "AWS encryption enabled with customer managed key",
			cloudType: goaviatrix.AWS,
			raw: map[string]any{
				"enable_encrypt_volume": true,
				"customer_managed_keys": "test-key",
			},
			expectedEnc: "yes",
			expectedCmk: "test-key",
		},
		{
			name:      "AWS encryption disabled explicitly sends no",
			cloudType: goaviatrix.AWS,
			raw: map[string]any{
				"enable_encrypt_volume": false,
			},
			expectedEnc: "no",
		},
		{
			name:      "non-AWS with encryption disabled sends nothing",
			cloudType: goaviatrix.Azure,
			raw: map[string]any{
				"enable_encrypt_volume": false,
			},
			expectedEnc: "",
		},
		{
			name:      "non-AWS with encryption enabled is rejected",
			cloudType: goaviatrix.Azure,
			raw: map[string]any{
				"enable_encrypt_volume": true,
			},
			expectErr: true,
		},
		{
			name:      "customer managed key without encryption is rejected",
			cloudType: goaviatrix.AWS,
			raw: map[string]any{
				"enable_encrypt_volume": false,
				"customer_managed_keys": "test-key",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAviatrixTransitInstance().Schema, tt.raw)
			gateway := &goaviatrix.TransitVpc{CloudType: tt.cloudType}

			diags := configureTransitInstanceEncryption(d, gateway, tt.cloudType)

			if tt.expectErr {
				assert.True(t, diags.HasError())
				return
			}
			assert.False(t, diags.HasError())
			assert.Equal(t, tt.expectedEnc, gateway.EncVolume)
			assert.Equal(t, tt.expectedCmk, gateway.CustomerManagedKeys)
		})
	}
}
