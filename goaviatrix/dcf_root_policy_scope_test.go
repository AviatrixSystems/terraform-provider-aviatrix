package goaviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPatchDCFPolicyScopeConfig verifies the PATCH endpoint, request, and response.
func TestPatchDCFPolicyScopeConfig(t *testing.T) {
	const rootUUID = "defa11a1-3000-0000-0000-000000000000"
	var body DCFPolicyScopeConfigPatch
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/v2.5/api/microseg/policy-list3/"+rootUUID, r.URL.Path)
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		_, _ = w.Write([]byte(`{
			"name":"System Root Policy Block",
			"uuid":"` + rootUUID + `",
			"system_resource":true,
			"sub_policies":[],
			"scope_config":{"scope_uuids":["scope-1","scope-2"]}
		}`))
	}))

	policyBlock, err := client.PatchDCFPolicyScopeConfig(
		context.Background(),
		rootUUID,
		[]string{"scope-1", "scope-2"},
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"scope-1", "scope-2"}, body.ScopeConfig.ScopeUUIDs)
	require.NotNil(t, policyBlock.ScopeConfig)
	assert.Equal(t, []string{"scope-1", "scope-2"}, policyBlock.ScopeConfig.ScopeUUIDs)
}

// TestGetDCFPolicyBlockScopeConfig verifies GET decodes root policy scopes.
func TestGetDCFPolicyBlockScopeConfig(t *testing.T) {
	const rootUUID = "defa11a1-3000-0000-0000-000000000000"
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2.5/api/microseg/policy-list3/"+rootUUID, r.URL.Path)
		_, _ = w.Write([]byte(`{
			"name":"System Root Policy Block",
			"uuid":"` + rootUUID + `",
			"system_resource":true,
			"sub_policies":[],
			"scope_config":{"scope_uuids":["scope-1"]}
		}`))
	}))

	policyBlock, err := client.GetDCFPolicyBlock(context.Background(), rootUUID)

	require.NoError(t, err)
	require.NotNil(t, policyBlock.ScopeConfig)
	assert.Equal(t, []string{"scope-1"}, policyBlock.ScopeConfig.ScopeUUIDs)
}

// TestGetDCFPolicyBlockNotFound verifies GET maps HTTP 404 to ErrNotFound.
func TestGetDCFPolicyBlockNotFound(t *testing.T) {
	client := newDCFPolicyScopeTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Policy block not found"}`))
	}))

	_, err := client.GetDCFPolicyBlock(context.Background(), "missing")

	assert.ErrorIs(t, err, ErrNotFound)
}
