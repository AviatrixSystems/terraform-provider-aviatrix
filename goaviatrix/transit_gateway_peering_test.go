package goaviatrix

import (
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
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "success"}))
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
			require.NoError(t, client.UpdateEdgeSpokeTransitPeeringTunnelCount("spoke-gw", "transit-gw", tt.tunnelCount))

			assert.Equal(t, "edit_inter_transit_gateway_peering", captured["action"])
			assert.Equal(t, "spoke-gw", captured["gateway1"])
			assert.Equal(t, "transit-gw", captured["gateway2"])
			assert.Equal(t, "test-cid", captured["CID"])
			// The literal string value must be sent; for 0 this is the crux of
			// the fix (an int 0 would arrive as an empty, controller-ignored value).
			assert.Equal(t, "yes", captured["__tunnel_count_present"], "tunnel_count key must be present")
			assert.Equal(t, tt.wantValue, captured["tunnel_count"])
		})
	}
}
