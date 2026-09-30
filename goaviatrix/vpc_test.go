package goaviatrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetVpcAzureLiveSubnetNamesParsesList verifies the live subnet names from
// list_vpc_all_subnets are returned as a name-keyed set, so the Read path can
// reconcile DB subnets against live cloud state (AVX-67843). Blank names are
// dropped so they can't match a real subnet during filtering.
func TestGetVpcAzureLiveSubnetNamesParsesList(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"return": true,
		"results": []map[string]any{
			{"cidr": "10.0.0.0/28", "name": "vnet-Public-FW-ingress-egress-1", "str": "10.0.0.0/28~~vnet-Public-FW-ingress-egress-1"},
			{"cidr": "10.0.0.0/28", "name": "vnet-Public-FW-ingress-egress-1", "str": "10.0.0.0/28~~az-1~~vnet-Public-FW-ingress-egress-1"},
			{"cidr": "10.0.0.32/28", "name": "vnet-Public-gateway-and-firewall-mgmt-1", "str": "10.0.0.32/28~~vnet-Public-gateway-and-firewall-mgmt-1"},
			{"name": ""},
		},
	})

	got, err := client.GetVpcAzureLiveSubnetNames(&Vpc{VpcID: "vnet:rg:uuid"})
	require.NoError(t, err)
	assert.Equal(t, map[string]bool{
		"vnet-Public-FW-ingress-egress-1":         true,
		"vnet-Public-gateway-and-firewall-mgmt-1": true,
	}, got, "blank names must be dropped, real names kept")
}

// TestGetVpcAzureLiveSubnetNamesEmptyList verifies an empty subnet_list yields a
// non-nil, empty map. The Read path treats an empty result as ambiguous and
// skips filtering, so this must not be conflated with an error.
func TestGetVpcAzureLiveSubnetNamesEmptyList(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"return":  true,
		"results": []map[string]any{},
	})

	got, err := client.GetVpcAzureLiveSubnetNames(&Vpc{VpcID: "vnet:rg:uuid"})
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Empty(t, got)
}

// TestGetVpcAzureLiveSubnetNamesPropagatesError verifies a controller failure is
// surfaced rather than returned as an empty set, so the Read path errors out
// instead of silently wiping every subnet from state as false drift.
func TestGetVpcAzureLiveSubnetNamesPropagatesError(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"return": false,
		"reason": "internal server error",
	})

	_, err := client.GetVpcAzureLiveSubnetNames(&Vpc{VpcID: "vnet:rg:uuid"})
	require.Error(t, err)
}

// TestGetVpcAzureLiveSubnetNamesRequestsJSONFormat verifies the request asks
// for json_format, since without it the controller returns "cidr~~name"
// strings rather than objects carrying a name field.
func TestGetVpcAzureLiveSubnetNamesRequestsJSONFormat(t *testing.T) {
	var gotJSONFormat string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotJSONFormat = r.URL.Query().Get("json_format")
		if err := json.NewEncoder(w).Encode(map[string]any{"return": true, "results": []map[string]any{}}); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	client := &Client{HTTPClient: server.Client(), CID: "test-cid", baseURL: server.URL}

	_, err := client.GetVpcAzureLiveSubnetNames(&Vpc{VpcID: "vnet:rg:uuid"})
	require.NoError(t, err)
	assert.Equal(t, "true", gotJSONFormat)
}
