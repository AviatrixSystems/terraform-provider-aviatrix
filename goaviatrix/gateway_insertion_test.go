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

type recordedRequest struct {
	method string
	path   string
	query  string
	body   map[string]any
}

// newGatewayInsertionTestClient returns a v2.5 Client wired to a test server that
// records each request and responds with the given status code and payload.
func newGatewayInsertionTestClient(t *testing.T, code int, resp any) (*Client, *recordedRequest) {
	t.Helper()
	rec := &recordedRequest{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.path = r.URL.Path
		rec.query = r.URL.RawQuery
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&rec.body)
		}
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return &Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
	}, rec
}

func TestEnableGatewayInsertionSendsVpcConfig(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusCreated, map[string]any{
		"vpc_id": "vpc-1", "status": "success", "err_message": "",
	})

	err := client.EnableGatewayInsertion(context.Background(), &GatewayInsertionVpcConfig{
		VpcID:      "vpc-1",
		EgressMode: GatewayInsertionEgressModeLocalEgress,
	})
	require.NoError(t, err)

	assert.Equal(t, "PUT", rec.method)
	assert.Equal(t, "/v2.5/api/gateway-insertion/vpc", rec.path)
	assert.Equal(t, map[string]any{
		"vpc_config": map[string]any{
			"vpc_id":       "vpc-1",
			"route_tables": []any{},
			"egress_mode":  "local-egress",
		},
	}, rec.body)
}

// The controller returns HTTP 201 even when the enablement workflow fails, so the
// status field of the response body is the only failure signal.
func TestEnableGatewayInsertionFailsOnWorkflowFailure(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusCreated, map[string]any{
		"vpc_id": "vpc-1", "status": "failed", "err_message": "No gateways or no spoke gateways found for VPC vpc-1.",
	})

	err := client.EnableGatewayInsertion(context.Background(), &GatewayInsertionVpcConfig{VpcID: "vpc-1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "No gateways or no spoke gateways found")
}

func TestEnableGatewayInsertionFailsOnHTTPError(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusNotImplemented, map[string]any{
		"message": "Gateway insertion feature is not enabled",
	})

	err := client.EnableGatewayInsertion(context.Background(), &GatewayInsertionVpcConfig{VpcID: "vpc-1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Gateway insertion feature is not enabled")
}

func TestGetGatewayInsertionStatusReturnsMatchingVpc(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusOK, []map[string]any{
		{"vpc_id": "vpc-1", "enabled": true, "route_table_config": []string{"rtb-1"}, "insertion_mode": "local-egress"},
	})

	status, err := client.GetGatewayInsertionStatus(context.Background(), "vpc-1")
	require.NoError(t, err)

	assert.Equal(t, "GET", rec.method)
	assert.Equal(t, "vpc_id=vpc-1", rec.query)
	assert.Equal(t, &GatewayInsertionStatus{
		VpcID:            "vpc-1",
		Enabled:          true,
		RouteTableConfig: []string{"rtb-1"},
		InsertionMode:    "local-egress",
	}, status)
}

func TestGetGatewayInsertionStatusFailsWhenVpcMissing(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusOK, []map[string]any{})

	_, err := client.GetGatewayInsertionStatus(context.Background(), "vpc-1")
	require.Error(t, err)
}

func TestDisableGatewayInsertionSendsVpcIDQuery(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusOK, map[string]any{
		"vpc_id": "vpc-1", "status": "success", "err_message": "",
	})

	require.NoError(t, client.DisableGatewayInsertion(context.Background(), "vpc-1"))
	assert.Equal(t, "DELETE", rec.method)
	assert.Equal(t, "/v2.5/api/gateway-insertion/vpc", rec.path)
	assert.Equal(t, "vpc_id=vpc-1", rec.query)
}

func TestDisableGatewayInsertionFailsOnWorkflowFailure(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusOK, map[string]any{
		"vpc_id": "vpc-1", "status": "failed", "err_message": "No gateway insertion config found for VPC vpc-1.",
	})

	err := client.DisableGatewayInsertion(context.Background(), "vpc-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "No gateway insertion config found")
}

func TestIsGatewayInsertionEnabled(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		resp    any
		want    bool
		wantErr bool
	}{
		{
			name: "enabled",
			code: http.StatusOK,
			resp: []map[string]any{{"vpc_id": "vpc-1", "enabled": true}},
			want: true,
		},
		{
			name: "not enabled",
			code: http.StatusOK,
			resp: []map[string]any{{"vpc_id": "vpc-1", "enabled": false}},
			want: false,
		},
		{
			name: "feature disabled on controller",
			code: http.StatusNotImplemented,
			resp: map[string]any{"message": "Gateway insertion feature is not enabled"},
			want: false,
		},
		{
			name: "endpoint not served by controller",
			code: http.StatusNotFound,
			resp: map[string]any{"message": "not found"},
			want: false,
		},
		{
			name:    "server error",
			code:    http.StatusInternalServerError,
			resp:    map[string]any{"message": "boom"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newGatewayInsertionTestClient(t, tt.code, tt.resp)
			got, err := client.IsGatewayInsertionEnabled(context.Background(), "vpc-1")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
