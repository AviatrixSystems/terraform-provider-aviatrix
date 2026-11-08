package goaviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpdateEdgeSpokeTransitPeeringTunnelCount exercises the real client method
// behind the aviatrix_edge_spoke_transit_attachment tunnel-count update
// (AVX-55065) against a stub controller, asserting the exact wire form.
//
// tunnel_count=0 means "use the max tunnels for the instance size". The
// controller edit handler gates on `if tunnel_count:` and only the string "0"
// is truthy there, so the request MUST carry a literal tunnel_count=0. The
// original bug encoded an int 0, which ajg/form serializes to an empty value
// (tunnel_count=) that the controller reads as falsy and silently ignores. If
// the fix ever regresses to the int path, the count==0 case below fails.
func TestUpdateEdgeSpokeTransitPeeringTunnelCount(t *testing.T) {
	var captured map[string]string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cap the body: this is a stub, and gosec (G120) flags ParseForm on an
		// unbounded body. Use assert (not require) inside the handler goroutine,
		// since require's FailNow must run on the test goroutine (testifylint go-require).
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		assert.NoError(t, r.ParseForm())

		// The edit runs async: the initial POST submits the task and gets a
		// request_id back in "results"; the client then polls check_task_status
		// until it returns. Report the poll as done immediately so the method
		// completes, and only capture the initial submit's form for assertions.
		if r.Form.Get("action") == "check_task_status" {
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "done"}))
			return
		}

		captured = make(map[string]string)
		for k, v := range r.Form {
			if len(v) > 0 {
				captured[k] = v[0]
			}
		}
		// Report presence explicitly so an empty value is distinguishable from
		// an absent key (r.Form drops neither, but the map above would).
		if _, ok := r.Form["tunnel_count"]; ok {
			captured["__tunnel_count_present"] = "yes"
		}
		// "results" carries the async request_id used for the follow-up poll.
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "req-123"}))
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		CID:        "test-cid",
		baseURL:    server.URL,
	}

	tests := []struct {
		name        string
		tunnelCount int
		wantValue   string
	}{
		{name: "zero means max tunnels", tunnelCount: 0, wantValue: "0"},
		{name: "explicit count", tunnelCount: 5, wantValue: "5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captured = nil
			require.NoError(t, client.UpdateEdgeSpokeTransitPeeringTunnelCount(context.Background(), "spoke-gw", "transit-gw", tt.tunnelCount))

			assert.Equal(t, "edit_inter_transit_gateway_peering", captured["action"])
			assert.Equal(t, "spoke-gw", captured["gateway1"])
			assert.Equal(t, "transit-gw", captured["gateway2"])
			assert.Equal(t, "test-cid", captured["CID"])
			// The edit must be submitted asynchronously so the multi-gateway
			// drain/rebuild does not overrun the request timeout (AVX-55065).
			assert.Equal(t, "true", captured["async"], "edit must be submitted async")
			// The literal string value must be sent; for 0 this is the crux of
			// the fix (an int 0 would arrive as an empty, controller-ignored value).
			assert.Equal(t, "yes", captured["__tunnel_count_present"], "tunnel_count key must be present")
			assert.Equal(t, tt.wantValue, captured["tunnel_count"])
		})
	}
}

// TestUpdateTransitGatewayPeeringFilters asserts the wire form of the peering
// filter edit (AVX-83481). The controller treats an absent filter cidr key as
// "preserve" and an empty value as "clear", so every filter key must reach the
// controller on every filter edit, empty ones included. The whole form is
// compared so a dropped or unexpected key fails the test.
func TestUpdateTransitGatewayPeeringFilters(t *testing.T) {
	var form map[string][]string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		assert.NoError(t, r.ParseForm())
		form = r.Form
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "Success"}))
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		CID:        "test-cid",
		baseURL:    server.URL,
	}

	// wantForm builds the expected form: the fixed keys plus every filter key,
	// defaulting to "" and overridden by the case's non-empty values.
	wantForm := func(set map[string]string) map[string][]string {
		f := map[string][]string{
			"CID":                             {"test-cid"},
			"action":                          {"edit_inter_transit_gateway_peering"},
			"gateway1":                        {"transit-1"},
			"gateway2":                        {"transit-2"},
			"src_filter_list":                 {""},
			"dst_filter_list":                 {""},
			"src_to_dst_filter_cidrs":         {""},
			"dst_to_src_filter_cidrs":         {""},
			"source_exclude_connections":      {""},
			"destination_exclude_connections": {""},
		}
		for k, v := range set {
			f[k] = []string{v}
		}
		return f
	}

	tests := []struct {
		name    string
		peering TransitGatewayPeering
		set     map[string]string
	}{
		{
			name: "empty lists are sent as empty values",
			peering: TransitGatewayPeering{
				TransitGatewayName1: "transit-1",
				TransitGatewayName2: "transit-2",
			},
		},
		{
			name: "one side cleared, the other kept",
			peering: TransitGatewayPeering{
				TransitGatewayName1:            "transit-1",
				TransitGatewayName2:            "transit-2",
				Gateway2ExcludedCIDRs:          "10.1.0.0/16",
				Gateway2ToGateway1FilterCIDRs:  "100.82.0.0/16",
				Gateway2ExcludedTGWConnections: "conn-2",
			},
			set: map[string]string{
				"dst_filter_list":                 "10.1.0.0/16",
				"dst_to_src_filter_cidrs":         "100.82.0.0/16",
				"destination_exclude_connections": "conn-2",
			},
		},
		{
			name: "non-empty lists are sent as-is",
			peering: TransitGatewayPeering{
				TransitGatewayName1:            "transit-1",
				TransitGatewayName2:            "transit-2",
				Gateway1ExcludedCIDRs:          "10.0.0.0/16",
				Gateway2ExcludedCIDRs:          "10.1.0.0/16",
				Gateway1ToGateway2FilterCIDRs:  "100.81.0.0/16,100.83.0.0/16",
				Gateway2ToGateway1FilterCIDRs:  "100.82.0.0/16",
				Gateway1ExcludedTGWConnections: "conn-1",
				Gateway2ExcludedTGWConnections: "conn-2",
			},
			set: map[string]string{
				"src_filter_list":                 "10.0.0.0/16",
				"dst_filter_list":                 "10.1.0.0/16",
				"src_to_dst_filter_cidrs":         "100.81.0.0/16,100.83.0.0/16",
				"dst_to_src_filter_cidrs":         "100.82.0.0/16",
				"source_exclude_connections":      "conn-1",
				"destination_exclude_connections": "conn-2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form = nil
			require.NoError(t, client.UpdateTransitGatewayPeeringFilters(&tt.peering))
			assert.Equal(t, wantForm(tt.set), form)
		})
	}
}
