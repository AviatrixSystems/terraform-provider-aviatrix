package goaviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newDCFPolicyScopeTestClient creates a v2.5 client backed by a test server.
func newDCFPolicyScopeTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)
	return &Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
	}
}

// TestCreateDCFPolicyScopeOmitsExcluded verifies request omission and response decoding.
func TestCreateDCFPolicyScopeOmitsExcluded(t *testing.T) {
	var body map[string]any
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v2.5/api/dcf/policy-scopes", r.URL.Path)
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"name":"production",
			"uuid":"123e4567-e89b-12d3-a456-426614174000",
			"included":{"selectors":[{"csps":["aws"]}]},
			"excluded":{"selectors":[{"preset":["none"]}]}
		}`))
	}))

	scope, err := client.CreateDCFPolicyScope(context.Background(), &DCFPolicyScopeInput{
		Name: "production",
		Included: DCFPolicyScopeSelectors{Selectors: []DCFPolicyScopeSelector{
			{"csps": {"aws"}},
		}},
	})

	require.NoError(t, err)
	assert.NotContains(t, body, "excluded")
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", scope.UUID)
	require.NotNil(t, scope.Excluded)
	assert.Equal(t, []string{"none"}, scope.Excluded.Selectors[0]["preset"])
}

// TestDCFPolicyScopeCRUDPathsAndBodies verifies UUID endpoints and update payloads.
func TestDCFPolicyScopeCRUDPathsAndBodies(t *testing.T) {
	const scopeUUID = "123e4567-e89b-12d3-a456-426614174000"
	var updateBody DCFPolicyScopeInput
	requests := 0
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		assert.Equal(t, "/v2.5/api/dcf/policy-scopes/"+scopeUUID, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{
				"name":"production",
				"uuid":"` + scopeUUID + `",
				"included":{"selectors":[{"vpc_ids":["vpc-1","vpc-2"]}]},
				"excluded":{"selectors":[{"preset":["none"]}]}
			}`))
		case http.MethodPut:
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&updateBody))
			_, _ = w.Write([]byte(`{}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	scope, err := client.GetDCFPolicyScope(context.Background(), scopeUUID)
	require.NoError(t, err)
	assert.Equal(t, []string{"vpc-1", "vpc-2"}, scope.Included.Selectors[0]["vpc_ids"])

	excluded := DCFPolicyScopeSelectors{Selectors: []DCFPolicyScopeSelector{{"gateway_group_ids": {"group-1"}}}}
	err = client.UpdateDCFPolicyScope(context.Background(), scopeUUID, &DCFPolicyScopeInput{
		Name:     "updated",
		Included: DCFPolicyScopeSelectors{Selectors: []DCFPolicyScopeSelector{{"preset": {"all"}}}},
		Excluded: &excluded,
	})
	require.NoError(t, err)
	assert.Equal(t, "updated", updateBody.Name)
	require.NotNil(t, updateBody.Excluded)
	assert.Equal(t, []string{"group-1"}, updateBody.Excluded.Selectors[0]["gateway_group_ids"])

	require.NoError(t, client.DeleteDCFPolicyScope(context.Background(), scopeUUID))
	assert.Equal(t, 3, requests)
}

// TestGetDCFPolicyScopeByName verifies lookup by name through the list endpoint.
func TestGetDCFPolicyScopeByName(t *testing.T) {
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2.5/api/dcf/policy-scopes", r.URL.Path)
		_, _ = w.Write([]byte(`{"scopes":[
			{"name":"Global Scope","uuid":"def000ad-8000-0000-0000-000000000001",
			 "included":{"selectors":[{"preset":["all"]}]},"excluded":{"selectors":[{"preset":["none"]}]}},
			{"name":"production","uuid":"123e4567-e89b-12d3-a456-426614174000",
			 "included":{"selectors":[{"vpc_ids":["vpc-1"]}]},"excluded":{"selectors":[{"preset":["none"]}]}}
		]}`))
	}))

	scope, err := client.GetDCFPolicyScopeByName(context.Background(), "production")
	require.NoError(t, err)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", scope.UUID)
	assert.Equal(t, []string{"vpc-1"}, scope.Included.Selectors[0]["vpc_ids"])

	_, err = client.GetDCFPolicyScopeByName(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrNotFound)
}

// TestDCFPolicyScopeNotFound verifies GET and DELETE map HTTP 404 to ErrNotFound.
func TestDCFPolicyScopeNotFound(t *testing.T) {
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Policy scope not found"}`))
	}))

	_, err := client.GetDCFPolicyScope(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrNotFound)
	err = client.DeleteDCFPolicyScope(context.Background(), "missing")
	assert.True(t, errors.Is(err, ErrNotFound))
}
