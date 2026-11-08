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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// fakePeeringFilterController answers the calls the transit peering Update
// makes and records every edit_inter_transit_gateway_peering request. Gateways
// resolve as plain gateways (not groups), and the trailing Read sees the
// peering as gone so Update returns without needing a full details payload.
type fakePeeringFilterController struct {
	edits []url.Values
}

func (f *fakePeeringFilterController) client(t *testing.T) *goaviatrix.Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		body, _ := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if form, err := url.ParseQuery(string(body)); err == nil {
			for k, v := range form {
				params[k] = v
			}
		}
		assert.NoError(t, json.NewEncoder(w).Encode(f.respond(t, params)))
	}))
	t.Cleanup(server.Close)
	client, err := goaviatrix.NewClient("admin", "password", strings.TrimPrefix(server.URL, "https://"), server.Client(), nil)
	require.NoError(t, err)
	return client
}

func (f *fakePeeringFilterController) respond(t *testing.T, params url.Values) map[string]any {
	switch action := params.Get("action"); action {
	case "get_api_token":
		return map[string]any{"return": true, "results": map[string]any{"api_token": "token"}}
	case "login":
		return map[string]any{"return": true, "CID": "cid"}
	case "get_gateway_group_details_by_name":
		return map[string]any{"return": false, "reason": "not found"}
	case "list_vpcs_summary":
		return map[string]any{"return": true, "results": []map[string]any{{"vpc_name": params.Get("gateway_name")}}}
	case "edit_inter_transit_gateway_peering":
		f.edits = append(f.edits, params)
		return map[string]any{"return": true, "results": "Success"}
	case "get_inter_transit_gateway_peering_details":
		return map[string]any{"return": false, "reason": "does not exist"}
	default:
		t.Errorf("unexpected controller action %q", action)
		return map[string]any{"return": false, "reason": "unexpected action " + action}
	}
}

// TestTransitGatewayPeeringUpdateFilters drives the resource Update for filter
// edits (AVX-83481). Clearing a filter must reach the controller as a
// present-but-empty key — the controller treats an absent key as "preserve" —
// and an edit to one filter field must resend every other filter field
// unchanged, since the controller replaces the configured set.
func TestTransitGatewayPeeringUpdateFilters(t *testing.T) {
	filterKeys := map[string]string{
		"gateway1_excluded_network_cidrs":   "src_filter_list",
		"gateway2_excluded_network_cidrs":   "dst_filter_list",
		"gateway1_to_gateway2_filter_cidrs": "src_to_dst_filter_cidrs",
		"gateway2_to_gateway1_filter_cidrs": "dst_to_src_filter_cidrs",
		"gateway1_excluded_tgw_connections": "source_exclude_connections",
		"gateway2_excluded_tgw_connections": "destination_exclude_connections",
	}
	stored := map[string][]string{
		"gateway1_excluded_network_cidrs":   {"10.0.0.0/16"},
		"gateway2_excluded_network_cidrs":   {"10.1.0.0/16"},
		"gateway1_to_gateway2_filter_cidrs": {"100.81.0.0/16", "100.83.0.0/16"},
		"gateway2_to_gateway1_filter_cidrs": {"100.82.0.0/16"},
		"gateway1_excluded_tgw_connections": {"conn-1"},
		"gateway2_excluded_tgw_connections": {"conn-2"},
	}

	tests := []struct {
		name     string
		set      map[string][]string // new config values, applied over the stored state
		wantEdit map[string]string   // expected filter keys on the edit; nil = no edit
	}{
		{
			name: "clearing every filter sends every key empty",
			set: map[string][]string{
				"gateway1_excluded_network_cidrs":   {},
				"gateway2_excluded_network_cidrs":   {},
				"gateway1_to_gateway2_filter_cidrs": {},
				"gateway2_to_gateway1_filter_cidrs": {},
				"gateway1_excluded_tgw_connections": {},
				"gateway2_excluded_tgw_connections": {},
			},
			wantEdit: map[string]string{
				"src_filter_list":                 "",
				"dst_filter_list":                 "",
				"src_to_dst_filter_cidrs":         "",
				"dst_to_src_filter_cidrs":         "",
				"source_exclude_connections":      "",
				"destination_exclude_connections": "",
			},
		},
		{
			name: "clearing one excluded cidr list resends the rest",
			set: map[string][]string{
				"gateway1_excluded_network_cidrs": {},
			},
			wantEdit: map[string]string{
				"src_filter_list":                 "",
				"dst_filter_list":                 "10.1.0.0/16",
				"src_to_dst_filter_cidrs":         "100.81.0.0/16,100.83.0.0/16",
				"dst_to_src_filter_cidrs":         "100.82.0.0/16",
				"source_exclude_connections":      "conn-1",
				"destination_exclude_connections": "conn-2",
			},
		},
		{
			name: "no filter change sends no edit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resourceAviatrixTransitGatewayPeering()
			state := planData(t, r, nil, toConfig(stored)).State()
			next := toConfig(stored)
			for k, v := range tt.set {
				next[k] = toAnyList(v)
			}
			d := planData(t, r, state, next)

			fake := &fakePeeringFilterController{}
			require.NoError(t, resourceAviatrixTransitGatewayPeeringUpdate(d, fake.client(t)))

			if tt.wantEdit == nil {
				assert.Empty(t, fake.edits)
				return
			}
			require.Len(t, fake.edits, 1)
			edit := fake.edits[0]
			assert.Equal(t, "transit-1", edit.Get("gateway1"))
			assert.Equal(t, "transit-2", edit.Get("gateway2"))
			for _, key := range filterKeys {
				want, ok := tt.wantEdit[key]
				require.True(t, ok, "test case missing expectation for %q", key)
				assert.Contains(t, edit, key, "filter key %q must be sent", key)
				assert.Equal(t, want, edit.Get(key), "key %q", key)
			}
			_, sentAz := edit["enable_az_affinity"]
			assert.False(t, sentAz, "filter edit must not carry enable_az_affinity")
		})
	}
}

// toConfig builds a raw resource config for the two transit gateways plus the
// given filter lists.
func toConfig(filters map[string][]string) map[string]any {
	cfg := map[string]any{
		"transit_gateway_name1": "transit-1",
		"transit_gateway_name2": "transit-2",
	}
	for k, v := range filters {
		cfg[k] = toAnyList(v)
	}
	return cfg
}

func toAnyList(v []string) []any {
	out := make([]any, len(v))
	for i, s := range v {
		out[i] = s
	}
	return out
}

// planData returns the ResourceData Terraform hands to Create/Update when
// moving from state to cfg, so HasChange reflects a real plan diff.
func planData(t *testing.T, r *schema.Resource, state *terraform.InstanceState, cfg map[string]any) *schema.ResourceData {
	t.Helper()
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
	require.NoError(t, err)
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	require.NoError(t, err)
	if state == nil {
		d.SetId("transit-1~transit-2")
	}
	return d
}
