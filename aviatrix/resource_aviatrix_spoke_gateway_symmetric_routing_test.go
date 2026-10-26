package aviatrix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// fakeSymmetricRoutingController answers the calls setSpokeGatewaySymmetricRouting
// makes and records the update_gateway_group request.
type fakeSymmetricRoutingController struct {
	groupUUID     string
	enabled       bool
	updateRequest map[string]any
}

func (f *fakeSymmetricRoutingController) client(t *testing.T) *goaviatrix.Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := map[string]any{}
		for k, v := range r.URL.Query() {
			params[k] = v[0]
		}
		body, _ := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			assert.NoError(t, json.Unmarshal(body, &params))
		} else if form, err := url.ParseQuery(string(body)); err == nil {
			for k, v := range form {
				params[k] = v[0]
			}
		}
		assert.NoError(t, json.NewEncoder(w).Encode(f.respond(t, params)))
	}))
	t.Cleanup(server.Close)
	client, err := goaviatrix.NewClient("admin", "password", strings.TrimPrefix(server.URL, "https://"), server.Client(), nil)
	require.NoError(t, err)
	return client
}

func (f *fakeSymmetricRoutingController) respond(t *testing.T, params map[string]any) map[string]any {
	switch action := params["action"]; action {
	case "get_api_token":
		return map[string]any{"return": true, "results": map[string]any{"api_token": "token"}}
	case "login":
		return map[string]any{"return": true, "CID": "cid"}
	case "list_vpcs_summary":
		return map[string]any{"return": true, "results": []map[string]any{{
			"vpc_name":                 "spoke-gw",
			"group_uuid":               f.groupUUID,
			"enable_symmetric_routing": f.enabled,
		}}}
	case "update_gateway_group":
		f.updateRequest = params
		return map[string]any{"return": true, "results": "updated"}
	default:
		t.Errorf("unexpected action %q", action)
		return map[string]any{"return": false, "reason": "unexpected action"}
	}
}

func TestSetSpokeGatewaySymmetricRoutingUpdatesOwningGroup(t *testing.T) {
	for _, enable := range []bool{true, false} {
		fake := &fakeSymmetricRoutingController{groupUUID: "group-uuid-1"}
		client := fake.client(t)

		require.NoError(t, setSpokeGatewaySymmetricRouting(context.Background(), client, "spoke-gw", enable))
		require.NotNil(t, fake.updateRequest, "update_gateway_group was not called")
		assert.Equal(t, "group-uuid-1", fake.updateRequest["group_uuid"])
		assert.Equal(t, enable, fake.updateRequest["enable_symmetric_routing"])
	}
}

// Without a group UUID there is nothing to update, so fail instead of sending an empty group_uuid.
func TestSetSpokeGatewaySymmetricRoutingRequiresGroupUUID(t *testing.T) {
	fake := &fakeSymmetricRoutingController{}
	client := fake.client(t)

	err := setSpokeGatewaySymmetricRouting(context.Background(), client, "spoke-gw", true)
	require.ErrorContains(t, err, "has no gateway group")
	assert.Nil(t, fake.updateRequest)
}

// Read sets enable_symmetric_routing from the group fields list_vpcs_summary returns.
func TestGetGatewayDecodesSymmetricRouting(t *testing.T) {
	fake := &fakeSymmetricRoutingController{groupUUID: "group-uuid-1", enabled: true}
	client := fake.client(t)

	gw, err := client.GetGateway(&goaviatrix.Gateway{GwName: "spoke-gw"})
	require.NoError(t, err)
	assert.Equal(t, "group-uuid-1", gw.GroupUUID)
	assert.True(t, gw.EnableSymmetricRouting)
}

func TestSpokeGatewayEnableSymmetricRoutingSchema(t *testing.T) {
	field, ok := resourceAviatrixSpokeGateway().Schema["enable_symmetric_routing"]
	require.True(t, ok)
	assert.True(t, field.Optional)
	assert.False(t, field.ForceNew, "toggling symmetric routing must not replace the gateway")
	// No default: an unset field must not disable symmetric routing enabled outside Terraform.
	assert.True(t, field.Computed)
	assert.Nil(t, field.Default)
}
