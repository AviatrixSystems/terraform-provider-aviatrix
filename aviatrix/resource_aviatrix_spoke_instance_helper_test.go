package aviatrix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestConfigureSpokeInstanceEIP(t *testing.T) {
	tests := []struct {
		name           string
		cloudType      int
		privateNetwork bool
		raw            map[string]any
		expectedEip    string
	}{
		{
			name:           "Azure private network with allocate_new_eip false sends nothing",
			cloudType:      goaviatrix.Azure,
			privateNetwork: true,
			raw: map[string]any{
				"allocate_new_eip": false,
			},
			expectedEip: "",
		},
		{
			name:           "Azure private network ignores a configured eip",
			cloudType:      goaviatrix.Azure,
			privateNetwork: true,
			raw: map[string]any{
				"allocate_new_eip":              false,
				"eip":                           "192.0.2.10",
				"azure_eip_name_resource_group": "test-ip:test-rg",
			},
			expectedEip: "",
		},
		{
			name:           "Azure reuses the combined name, resource group and address",
			cloudType:      goaviatrix.Azure,
			privateNetwork: false,
			raw: map[string]any{
				"allocate_new_eip":              false,
				"eip":                           "192.0.2.10",
				"azure_eip_name_resource_group": "test-ip:test-rg",
			},
			expectedEip: "test-ip:test-rg:192.0.2.10",
		},
		{
			name:           "Azure with allocate_new_eip true sends nothing",
			cloudType:      goaviatrix.Azure,
			privateNetwork: false,
			raw: map[string]any{
				"allocate_new_eip": true,
				"eip":              "192.0.2.10",
			},
			expectedEip: "",
		},
		{
			name:           "AWS reuses the address alone",
			cloudType:      goaviatrix.AWS,
			privateNetwork: false,
			raw: map[string]any{
				"allocate_new_eip": false,
				"eip":              "192.0.2.10",
			},
			expectedEip: "192.0.2.10",
		},
		{
			name:           "AWS private network ignores a configured eip",
			cloudType:      goaviatrix.AWS,
			privateNetwork: true,
			raw: map[string]any{
				"allocate_new_eip": false,
				"eip":              "192.0.2.10",
			},
			expectedEip: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAviatrixSpokeInstance().Schema, tt.raw)
			spokeGateway := &goaviatrix.SpokeVpc{CloudType: tt.cloudType}

			configureSpokeInstanceEIP(d, spokeGateway, tt.privateNetwork)

			assert.Equal(t, tt.expectedEip, spokeGateway.Eip)
		})
	}
}
