package aviatrix

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestGatewayAzureBgpLanIPListJSON(t *testing.T) {
	var withList goaviatrix.Gateway
	require.NoError(t, json.Unmarshal([]byte(`{"arm_bgp_lan_ip_list": ["10.1.2.4", "10.1.2.12"]}`), &withList))
	assert.Equal(t, []string{"10.1.2.4", "10.1.2.12"}, withList.AzureBgpLanIPList)

	var emptyList goaviatrix.Gateway
	require.NoError(t, json.Unmarshal([]byte(`{"arm_bgp_lan_ip_list": []}`), &emptyList))
	assert.NotNil(t, emptyList.AzureBgpLanIPList)
	assert.Empty(t, emptyList.AzureBgpLanIPList)

	var withoutList goaviatrix.Gateway
	require.NoError(t, json.Unmarshal([]byte(`{}`), &withoutList))
	assert.Nil(t, withoutList.AzureBgpLanIPList)
}

func TestSetAzureBgpLanState(t *testing.T) {
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
			name: "primary instance uses its own IPs",
			gw: &goaviatrix.Gateway{
				GwName: "transit-1", IsHagw: "no", EnableBgpOverLan: true,
				AzureBgpLanIPList: []string{"10.1.2.4", "10.1.2.12"},
			},
			expectedIPs: []any{"10.1.2.4", "10.1.2.12"},
		},
		{
			name: "non-primary instance uses its own IPs, not the primary's",
			gw: &goaviatrix.Gateway{
				GwName: "transit-3", IsHagw: "yes", EnableBgpOverLan: true,
				AzureBgpLanIPList: []string{"10.1.2.36", "10.1.2.44"},
			},
			expectedIPs: []any{"10.1.2.36", "10.1.2.44"},
		},
		{
			name: "older controller falls back for the primary instance",
			gw: &goaviatrix.Gateway{
				GwName: "transit-1", IsHagw: "no", EnableBgpOverLan: true,
			},
			expectFallback: true,
			expectedIPs:    []any{"10.1.2.4", "10.1.2.12"},
		},
		{
			name: "older controller leaves non-primary IPs empty",
			gw: &goaviatrix.Gateway{
				GwName: "transit-3", IsHagw: "yes", EnableBgpOverLan: true,
			},
			expectedIPs: []any{},
		},
		{
			name:        "instance without BGP over LAN clears the lists",
			gw:          &goaviatrix.Gateway{GwName: "transit-2", IsHagw: "yes"},
			expectedIPs: []any{},
		},
		{
			name: "fallback error is returned",
			gw: &goaviatrix.Gateway{
				GwName: "transit-1", IsHagw: "no", EnableBgpOverLan: true,
			},
			fallbackErr:    errors.New("boom"),
			expectFallback: true,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAviatrixTransitInstance().Schema, map[string]any{})
			fallbackCalled := false
			getPrimaryIPList := func() ([]string, error) {
				fallbackCalled = true
				if tt.fallbackErr != nil {
					return nil, tt.fallbackErr
				}
				return primaryIPs, nil
			}

			err := setAzureBgpLanState(d, tt.gw, getPrimaryIPList)

			assert.Equal(t, tt.expectFallback, fallbackCalled)
			if tt.expectErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.fallbackErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedIPs, d.Get("bgp_lan_ip_list"))
			assert.Equal(t, tt.expectedIPs, d.Get("azure_bgp_lan_ip_list"))
		})
	}
}

func TestValidateAndConfigureBasicSettingsZone(t *testing.T) {
	tests := []struct {
		name           string
		cloudType      int
		zone           string
		expectErr      bool
		expectedSubnet string
	}{
		{
			name:           "Azure without zone uses availability set",
			cloudType:      goaviatrix.Azure,
			expectedSubnet: "10.99.3.0/24",
		},
		{
			name:           "Azure with zone appends it to subnet",
			cloudType:      goaviatrix.Azure,
			zone:           "az-2",
			expectedSubnet: "10.99.3.0/24~~az-2~~",
		},
		{
			name:      "Azure with invalid zone is rejected",
			cloudType: goaviatrix.Azure,
			zone:      "2",
			expectErr: true,
		},
		{
			name:      "GCP without zone is rejected",
			cloudType: goaviatrix.GCP,
			expectErr: true,
		},
		{
			name:      "AWS with zone is rejected",
			cloudType: goaviatrix.AWS,
			zone:      "az-1",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := map[string]any{"subnet": "10.99.3.0/24"}
			if tt.zone != "" {
				raw["zone"] = tt.zone
			}
			d := schema.TestResourceDataRaw(t, resourceAviatrixTransitInstance().Schema, raw)
			gateway := &goaviatrix.TransitVpc{CloudType: tt.cloudType, Subnet: getString(d, "subnet")}

			diags := validateAndConfigureBasicSettings(d, gateway, tt.cloudType)

			if tt.expectErr {
				assert.True(t, diags.HasError())
				return
			}
			assert.False(t, diags.HasError())
			assert.Equal(t, tt.expectedSubnet, gateway.Subnet)
		})
	}
}
