package aviatrix

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestValidateFirenetGwName(t *testing.T) {
	azure := strconv.Itoa(goaviatrix.Azure)
	azureGov := strconv.Itoa(goaviatrix.AzureGov)
	aws := strconv.Itoa(goaviatrix.AWS)

	tests := []struct {
		name          string
		gwName        string
		firenetDetail *goaviatrix.FireNetDetail
		wantErr       string
	}{
		{
			name:          "azure horizontal scaling without gateway",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: azure, HorizontalScaling: true},
		},
		{
			name:          "azure gov horizontal scaling without gateway",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: azureGov, HorizontalScaling: true},
		},
		{
			name:          "azure horizontal scaling with gateway",
			gwName:        "transit-gw",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: azure, HorizontalScaling: true},
		},
		{
			name:          "azure without horizontal scaling requires gateway",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: azure},
			wantErr:       "'firenet_gw_name' is required",
		},
		{
			name:          "aws horizontal scaling still requires gateway",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: aws, HorizontalScaling: true},
			wantErr:       "'firenet_gw_name' is required",
		},
		{
			name:    "unknown firenet requires gateway",
			wantErr: "'firenet_gw_name' is required",
		},
		{
			name:          "non gwlb with gateway",
			gwName:        "transit-gw",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: aws},
		},
		{
			name:          "native gwlb without gateway",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: aws, NativeGwlb: true},
		},
		{
			name:          "native gwlb rejects gateway",
			gwName:        "transit-gw",
			firenetDetail: &goaviatrix.FireNetDetail{CloudType: aws, NativeGwlb: true},
			wantErr:       "Native GWLB enabled but a 'firenet_gw_name' was provided",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFirenetGwName(tt.gwName, "vpc-1", tt.firenetDetail)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestFirenetGwNameReadBackNoReplacement guards against a regression where a
// firewall created without 'firenet_gw_name' (Azure horizontal scaling) was
// replaced on the next plan because the controller reports the transit group
// name on read and the config has no value.
func TestFirenetGwNameReadBackNoReplacement(t *testing.T) {
	resources := map[string]*schema.Resource{
		"aviatrix_firewall_instance":             resourceAviatrixFirewallInstance(),
		"aviatrix_firewall_instance_association": resourceAviatrixFirewallInstanceAssociation(),
	}
	for name, r := range resources {
		t.Run(name, func(t *testing.T) {
			s := r.Schema["firenet_gw_name"]
			require.NotNil(t, s)
			assert.True(t, s.Optional)
			assert.True(t, s.Computed)

			state := &terraform.InstanceState{
				ID:         "id",
				Attributes: map[string]string{"firenet_gw_name": "azure-transit-gw"},
			}
			config := terraform.NewResourceConfigRaw(map[string]any{})
			diff, err := schema.InternalMap(map[string]*schema.Schema{"firenet_gw_name": s}).
				Diff(context.Background(), state, config, nil, nil, true)
			require.NoError(t, err)
			if diff != nil {
				assert.NotContains(t, diff.Attributes, "firenet_gw_name")
				assert.False(t, diff.RequiresNew())
			}
		})
	}
}
