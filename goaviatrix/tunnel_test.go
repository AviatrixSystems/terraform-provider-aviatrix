package goaviatrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTunnelTestClient wires a Client to a server that dispatches on the
// `action` query param: list_peer_vpc_pairs returns the given pair list, while
// list_vpcs_summary resolves a gateway name to its group (used by GetTunnel to
// reconcile config gateway names against group-keyed spoke peerings).
func newTunnelTestClient(t *testing.T, pairs []map[string]any, gwToGroup map[string]string) *Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]any
		switch r.URL.Query().Get("action") {
		case "list_peer_vpc_pairs":
			resp = map[string]any{
				"return":  true,
				"results": map[string]any{"pair_list": pairs},
			}
		case "list_vpcs_summary":
			gwName := r.URL.Query().Get("gateway_name")
			resp = map[string]any{
				"return": true,
				"results": []map[string]any{
					// GwName decodes from json "vpc_name" (see Gateway struct).
					{"vpc_name": gwName, "group_name": gwToGroup[gwName]},
				},
			}
		default:
			resp = map[string]any{"return": true, "results": []any{}}
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return &Client{
		HTTPClient: server.Client(),
		CID:        "test-cid",
		baseURL:    server.URL,
	}
}

// AVX-79062: spoke-to-spoke peerings are listed by gateway-GROUP name, but the
// aviatrix_tunnel config supplies individual gateway names. GetTunnel must
// resolve the gateway names to their groups so the created peering is found on
// Read — otherwise Read clears the ID and re-apply tries to recreate it.
func TestGetTunnelMatchesSpokePeeringByGroupName(t *testing.T) {
	client := newTunnelTestClient(t,
		[]map[string]any{
			{
				"vpc_name1":         "lily-aws-spoke1",
				"vpc_name2":         "lily-aws-spoke2",
				"peering_state":     "up",
				"peering_ha_status": "Activemesh",
				"peering_link":      "",
			},
		},
		map[string]string{
			"lily-aws-spoke1-gw1-east1a": "lily-aws-spoke1",
			"lily-aws-spoke2-gw1-east1a": "lily-aws-spoke2",
		},
	)

	tun, err := client.GetTunnel(&Tunnel{
		VpcName1: "lily-aws-spoke1-gw1-east1a",
		VpcName2: "lily-aws-spoke2-gw1-east1a",
	})
	require.NoError(t, err)
	assert.Equal(t, "lily-aws-spoke1", tun.VpcName1)
	assert.Equal(t, "lily-aws-spoke2", tun.VpcName2)
}

// gw_name1/gw_name2 may each be either a gateway name or a group name, and the
// two sides are independent. A group name won't resolve as a gateway, so
// resolveGroupName falls back to the input and the group-keyed pair still matches.
func TestGetTunnelMatchesWhenConfigUsesGroupNames(t *testing.T) {
	client := newTunnelTestClient(t,
		[]map[string]any{
			{"vpc_name1": "lily-aws-spoke1", "vpc_name2": "lily-aws-spoke2"},
		},
		// No gateway is named after the group, so GetGateway fails for both.
		map[string]string{},
	)

	tun, err := client.GetTunnel(&Tunnel{
		VpcName1: "lily-aws-spoke1",
		VpcName2: "lily-aws-spoke2",
	})
	require.NoError(t, err)
	assert.Equal(t, "lily-aws-spoke1", tun.VpcName1)
	assert.Equal(t, "lily-aws-spoke2", tun.VpcName2)
}

// Mixed input: gw_name1 is an individual gateway name while gw_name2 is a group
// name. Each side is resolved independently, so the pair still matches.
func TestGetTunnelMatchesMixedGatewayAndGroupNames(t *testing.T) {
	client := newTunnelTestClient(t,
		[]map[string]any{
			{"vpc_name1": "lily-aws-spoke1", "vpc_name2": "lily-aws-spoke2"},
		},
		map[string]string{
			"lily-aws-spoke1-gw1-east1a": "lily-aws-spoke1",
		},
	)

	tun, err := client.GetTunnel(&Tunnel{
		VpcName1: "lily-aws-spoke1-gw1-east1a", // gateway name
		VpcName2: "lily-aws-spoke2",            // group name
	})
	require.NoError(t, err)
	assert.Equal(t, "lily-aws-spoke1", tun.VpcName1)
	assert.Equal(t, "lily-aws-spoke2", tun.VpcName2)
}

// The stored pair is keyed with gateway names sorted, so a config whose
// gw_name1/gw_name2 order is reversed relative to the stored pair must still match.
func TestGetTunnelMatchesReversedOrder(t *testing.T) {
	client := newTunnelTestClient(t,
		[]map[string]any{
			{"vpc_name1": "gw-a", "vpc_name2": "gw-b"},
		},
		map[string]string{},
	)

	tun, err := client.GetTunnel(&Tunnel{VpcName1: "gw-b", VpcName2: "gw-a"})
	require.NoError(t, err)
	assert.Equal(t, "gw-a", tun.VpcName1)
}

// A standalone encrypted peering keyed by individual gateway names must still
// match directly, without needing a group resolution.
func TestGetTunnelMatchesGatewayNameDirectly(t *testing.T) {
	client := newTunnelTestClient(t,
		[]map[string]any{
			{"vpc_name1": "gw1", "vpc_name2": "gw2"},
		},
		map[string]string{},
	)

	tun, err := client.GetTunnel(&Tunnel{VpcName1: "gw1", VpcName2: "gw2"})
	require.NoError(t, err)
	assert.Equal(t, "gw1", tun.VpcName1)
}

// matchesTunnelPair is the contract GetTunnel relies on to reconcile a listed
// peering against the configured gateways. These cases pin that contract
// directly so a change to either the matcher or its callers can't silently
// break the pairing. The scenarios mirror the GetTunnel cases above:
// direct-name, group-resolved, mixed, and reversed orientation.
func TestMatchesTunnelPair(t *testing.T) {
	tests := []struct {
		name                   string
		pair1, pair2           string // as listed by list_peer_vpc_pairs
		cfg1, cfg2, grp1, grp2 string // configured names and their resolved groups
		want                   bool
	}{
		{
			name:  "direct gateway names, same order",
			pair1: "gw1", pair2: "gw2",
			cfg1: "gw1", cfg2: "gw2", grp1: "gw1", grp2: "gw2",
			want: true,
		},
		{
			name:  "direct gateway names, reversed order",
			pair1: "gw-a", pair2: "gw-b",
			cfg1: "gw-b", cfg2: "gw-a", grp1: "gw-b", grp2: "gw-a",
			want: true,
		},
		{
			name:  "spoke-to-spoke listed by group name",
			pair1: "spoke1", pair2: "spoke2",
			cfg1: "spoke1-gw1", cfg2: "spoke2-gw1", grp1: "spoke1", grp2: "spoke2",
			want: true,
		},
		{
			name:  "mixed gateway name and group name",
			pair1: "spoke1", pair2: "spoke2",
			cfg1: "spoke1-gw1", cfg2: "spoke2", grp1: "spoke1", grp2: "spoke2",
			want: true,
		},
		{
			name:  "no match",
			pair1: "other1", pair2: "other2",
			cfg1: "gw1", cfg2: "gw2", grp1: "gw1", grp2: "gw2",
			want: false,
		},
		{
			name:  "only one side matches",
			pair1: "gw1", pair2: "other2",
			cfg1: "gw1", cfg2: "gw2", grp1: "gw1", grp2: "gw2",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesTunnelPair(tt.pair1, tt.pair2, tt.cfg1, tt.cfg2, tt.grp1, tt.grp2)
			assert.Equal(t, tt.want, got)
		})
	}
}

// A tunnel that does not exist must return ErrNotFound so Read can clear state.
func TestGetTunnelNotFound(t *testing.T) {
	client := newTunnelTestClient(t,
		[]map[string]any{
			{"vpc_name1": "other1", "vpc_name2": "other2"},
		},
		map[string]string{},
	)

	_, err := client.GetTunnel(&Tunnel{VpcName1: "gw1", VpcName2: "gw2"})
	assert.ErrorIs(t, err, ErrNotFound)
}
