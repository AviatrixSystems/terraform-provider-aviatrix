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

// newTestClient returns a Client wired to a test server that responds with the
// given controller payload for every request.
func newTestClient(t *testing.T, resp map[string]any) *Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

// AVX-78638: when a transit group's backing VPC document has been deleted (e.g.
// by a rolled-back instance launch), the controller returns "VPC ... not found".
// GetGatewayGroup must translate that into ErrNotFound so the resource Read can
// clear the ID and recreate the group instead of erroring on every refresh.
func TestGetGatewayGroupMapsVPCNotFoundToErrNotFound(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"return": false,
		"reason": "rest API get_gateway_group_details Get failed: " +
			"VPC sky-azure-transit-vnet-1:rg:uuid not found.",
	})

	_, err := client.GetGatewayGroup(context.Background(), "some-uuid")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

// A genuine failure that is not a missing-resource condition must propagate
// unchanged so real errors are not silently swallowed as ErrNotFound.
func TestGetGatewayGroupPropagatesOtherErrors(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"return": false,
		"reason": "internal server error",
	})

	_, err := client.GetGatewayGroup(context.Background(), "some-uuid")
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNotFound,
		"unexpected ErrNotFound for a non-not-found failure: %v", err)
}
