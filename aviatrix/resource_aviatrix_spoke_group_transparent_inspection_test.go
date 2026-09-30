package aviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// newTransparentInspectionStatusClient returns a Client whose v2.5 gateway
// insertion status endpoint reports the given enabled value and counts calls.
func newTransparentInspectionStatusClient(t *testing.T, enabled bool) (*goaviatrix.Client, *int) {
	t.Helper()
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		resp := []map[string]any{{"vpc_id": r.URL.Query().Get("vpc_id"), "enabled": enabled}}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return &goaviatrix.Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
	}, &calls
}

func TestSpokeGroupTransparentInspectionEnabled(t *testing.T) {
	awsSpokeGroup := goaviatrix.GatewayGroup{
		GroupName: "spoke-group",
		CloudType: goaviatrix.AWS,
		GwType:    "GwGroupType.SPOKE",
		VpcID:     "vpc-1",
		EnableNat: true,
	}

	tests := []struct {
		name          string
		mutate        func(g *goaviatrix.GatewayGroup)
		statusEnabled bool
		want          bool
		wantCalls     int
	}{
		{
			name:          "enabled on AWS spoke group with NAT on",
			statusEnabled: true,
			want:          true,
			wantCalls:     1,
		},
		{
			name:          "not enabled on AWS spoke group with NAT on",
			statusEnabled: false,
			want:          false,
			wantCalls:     1,
		},
		{
			name:          "NAT off skips the status call",
			mutate:        func(g *goaviatrix.GatewayGroup) { g.EnableNat = false },
			statusEnabled: true,
			want:          false,
			wantCalls:     0,
		},
		{
			name:          "non-AWS group skips the status call",
			mutate:        func(g *goaviatrix.GatewayGroup) { g.CloudType = goaviatrix.Azure },
			statusEnabled: true,
			want:          false,
			wantCalls:     0,
		},
		{
			name:          "edge spoke group skips the status call",
			mutate:        func(g *goaviatrix.GatewayGroup) { g.GwType = "GwGroupType.EDGESPOKE" },
			statusEnabled: true,
			want:          false,
			wantCalls:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := awsSpokeGroup
			if tt.mutate != nil {
				tt.mutate(&group)
			}
			client, calls := newTransparentInspectionStatusClient(t, tt.statusEnabled)

			got, err := spokeGroupTransparentInspectionEnabled(context.Background(), client, &group)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, *calls)
		})
	}
}
