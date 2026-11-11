package aviatrix

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestSetSpokeAzureBgpLanState(t *testing.T) {
	primaryIPs := []string{"10.1.2.4", "10.1.2.12"}

	tests := []struct {
		name           string
		gw             *goaviatrix.Gateway
		fallbackErr    error
		expectErr      bool
		expectFallback bool
		expectedIPs    []any
	}{
		{
			name: "non-primary instance uses its own IPs, not the primary's",
			gw: &goaviatrix.Gateway{
				GwName: "spoke-3", IsHagw: "yes", EnableBgpOverLan: true,
				AzureBgpLanIPList: []string{"10.1.2.36", "10.1.2.44"},
			},
			expectedIPs: []any{"10.1.2.36", "10.1.2.44"},
		},
		{
			name: "older controller falls back for the primary instance",
			gw: &goaviatrix.Gateway{
				GwName: "spoke-1", IsHagw: "no", EnableBgpOverLan: true,
			},
			expectFallback: true,
			expectedIPs:    []any{"10.1.2.4", "10.1.2.12"},
		},
		{
			name: "older controller leaves non-primary IPs empty",
			gw: &goaviatrix.Gateway{
				GwName: "spoke-3", IsHagw: "yes", EnableBgpOverLan: true,
			},
			expectedIPs: []any{},
		},
		{
			name: "fallback error is returned",
			gw: &goaviatrix.Gateway{
				GwName: "spoke-1", IsHagw: "no", EnableBgpOverLan: true,
			},
			fallbackErr:    errors.New("boom"),
			expectFallback: true,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAviatrixSpokeInstance().Schema, map[string]any{})
			fallbackCalled := false
			getPrimaryIPList := func() ([]string, error) {
				fallbackCalled = true
				if tt.fallbackErr != nil {
					return nil, tt.fallbackErr
				}
				return primaryIPs, nil
			}

			err := setSpokeAzureBgpLanState(d, tt.gw, getPrimaryIPList)

			assert.Equal(t, tt.expectFallback, fallbackCalled)
			if tt.expectErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.fallbackErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedIPs, d.Get("azure_bgp_lan_ip_list"))
		})
	}
}
