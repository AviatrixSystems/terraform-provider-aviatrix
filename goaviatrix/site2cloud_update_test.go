package goaviatrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateSite2CloudSendsOnlyRequestedFields(t *testing.T) {
	var capturedRequests []map[string]string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm failed: %v", err)
			return
		}
		params := make(map[string]string)
		for k, v := range r.Form {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
		capturedRequests = append(capturedRequests, params)
		resp := map[string]any{"return": true, "results": "success"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		CID:        "test-cid",
		baseURL:    server.URL,
	}

	// Simulate updating local subnet only
	localEdit := &EditSite2Cloud{
		VpcID:              "vpc-123",
		ConnName:           "test-conn",
		GwName:             "test-gw",
		CloudSubnetCidr:    "10.1.0.0/16",
		CloudSubnetVirtual: "10.2.0.0/16",
	}
	require.NoError(t, client.UpdateSite2Cloud(localEdit), "UpdateSite2Cloud (local) failed")

	// Simulate updating remote subnet only
	remoteEdit := &EditSite2Cloud{
		VpcID:               "vpc-123",
		ConnName:            "test-conn",
		GwName:              "test-gw",
		RemoteSubnet:        "172.1.0.0/16",
		RemoteSubnetVirtual: "172.2.0.0/16",
	}
	require.NoError(t, client.UpdateSite2Cloud(remoteEdit), "UpdateSite2Cloud (remote) failed")

	require.Len(t, capturedRequests, 2, "expected 2 API calls, got %d", len(capturedRequests))

	// First call should have local subnet fields but NOT remote subnet fields
	localReq := capturedRequests[0]
	assert.Equal(t, "10.1.0.0/16", localReq["cloud_subnet_cidr"], "expected cloud_subnet_cidr=10.1.0.0/16, got %q", localReq["cloud_subnet_cidr"])
	assert.Empty(t, localReq["remote_cidr"], "remote_cidr should not be present")

	// Second call should have remote subnet fields but NOT local subnet fields
	remoteReq := capturedRequests[1]
	assert.Equal(t, "172.1.0.0/16", remoteReq["remote_cidr"], "expected remote_cidr=172.1.0.0/16, got %q", remoteReq["remote_cidr"])
	assert.Empty(t, remoteReq["cloud_subnet_cidr"], "cloud_subnet_cidr should not be present")
}
