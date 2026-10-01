package goaviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newSmartGatewayTestClient records each request's action and group and answers with resp.
func newSmartGatewayTestClient(t *testing.T, resp map[string]any, calls *[]string) *Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		params := map[string]string{"action": query.Get("action"), "gateway_group_name": query.Get("gateway_group_name")}
		if r.Method == http.MethodPost {
			assert.NoError(t, json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&params))
		}
		*calls = append(*calls, strings.TrimSpace(params["action"]+" "+params["gateway_group_name"]))
		assert.NoError(t, json.NewEncoder(w).Encode(resp))
	}))
	t.Cleanup(server.Close)
	return &Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
		baseURL:      server.URL,
	}
}

func TestSmartGatewayClientCalls(t *testing.T) {
	var calls []string
	client := newSmartGatewayTestClient(t, map[string]any{"return": true, "results": []map[string]any{
		{"gateway_group_name": "transit-1", "name": "transit-1", "gw_type": "TRANSIT", "smart_gw_underlay_enabled": true},
	}}, &calls)
	ctx := context.Background()

	groups, err := client.GetSmartGatewayFabricStatus(ctx)
	require.NoError(t, err)
	assert.Equal(t, []SmartGatewayGroupStatus{{GatewayGroupName: "transit-1", SmartGatewayUnderlayEnabled: true}}, groups)

	types, err := client.ListGatewayGroupTypes(ctx)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"transit-1": "TRANSIT"}, types)

	require.NoError(t, client.SetSmartGatewayUnderlay(ctx, "transit-1", true))
	require.NoError(t, client.SetSmartGatewayUnderlay(ctx, "transit-1", false))
	require.NoError(t, client.SetSmartGatewayResolver(ctx, "transit-1", true))
	require.NoError(t, client.SetSmartGatewayResolver(ctx, "transit-1", false))

	assert.Equal(t, []string{
		"get_smart_gw_fabric_status",
		"list_gateway_groups",
		"enable_smart_gw_underlay transit-1",
		"disable_smart_gw_underlay transit-1",
		"enable_smart_gw_resolver transit-1",
		"disable_smart_gw_resolver transit-1",
	}, calls)
}

func TestGetSmartGatewayResolverStatus(t *testing.T) {
	tests := map[string]bool{
		"Gateway group transit-1 smart_gw_resolver effective=enabled (fabric=enabled, underlay=enabled, resolver=enabled)":   true,
		"Gateway group transit-1 smart_gw_resolver effective=disabled (fabric=enabled, underlay=disabled, resolver=enabled)": true,
		"Gateway group transit-1 smart_gw_resolver effective=disabled (fabric=enabled, underlay=enabled, resolver=disabled)": false,
		"Gateway group spoke-1 smart_gw_resolver effective=disabled (fabric=enabled, underlay=disabled, resolver=enabled) " +
			"[cloud spoke: BGP-LU not applicable]": false,
	}
	for text, want := range tests {
		var calls []string
		client := newSmartGatewayTestClient(t, map[string]any{"return": true, "results": text}, &calls)
		got, err := client.GetSmartGatewayResolverStatus(context.Background(), "transit-1")
		require.NoError(t, err)
		assert.Equal(t, want, got, text)
		assert.Equal(t, []string{"get_smart_gw_resolver_status transit-1"}, calls)
	}
}

func TestIsSmartGatewayEligibleType(t *testing.T) {
	assert.True(t, IsSmartGatewayEligibleType("GwGroupType.EDGETRANSIT"))
	assert.True(t, IsSmartGatewayEligibleType("SPOKE"))
	assert.False(t, IsSmartGatewayEligibleType("GwGroupType.STANDALONE"))
}
