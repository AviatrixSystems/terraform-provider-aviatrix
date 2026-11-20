package aviatrix

import (
	"context"
	"encoding/json"
	"testing"

	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func externalConnectionTestConfig(raw map[string]any) map[string]any {
	base := map[string]any{
		"name":             "conn",
		"group_uuid":       "group-uuid",
		"routing_protocol": "bgp",
		"tunnel_protocol":  "LAN",
		"tunnels": []any{
			map[string]any{"local_gw_name": "gw-1", "remote_device_ip": "10.0.4.10"},
		},
	}
	for k, v := range raw {
		base[k] = v
	}
	return base
}

func externalConnectionTestData(t *testing.T, raw map[string]any) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceAviatrixExternalConnection().Schema, externalConnectionTestConfig(raw))
}

// externalConnectionTestDiff plans the change from a state built from oldRaw to
// the config newRaw, attaching the raw config to the prior state as the gRPC
// plan path does so CustomizeDiff can inspect it.
func externalConnectionTestDiff(t *testing.T, oldRaw, newRaw map[string]any) *terraform.InstanceDiff {
	t.Helper()
	diff, err := externalConnectionTestDiffErr(t, oldRaw, newRaw)
	require.NoError(t, err)
	return diff
}

func externalConnectionTestDiffErr(t *testing.T, oldRaw, newRaw map[string]any) (*terraform.InstanceDiff, error) {
	t.Helper()
	r := resourceAviatrixExternalConnection()
	old := externalConnectionTestData(t, oldRaw)
	old.SetId("conn-uuid")

	coreSchema := schema.InternalMap(r.Schema).CoreConfigSchema()
	cfgJSON, err := json.Marshal(externalConnectionTestConfig(newRaw))
	require.NoError(t, err)
	cfgVal, err := ctyjson.Unmarshal(cfgJSON, coreSchema.ImpliedType())
	require.NoError(t, err)

	state := old.State()
	state.RawConfig = cfgVal

	return r.Diff(context.Background(), state, terraform.NewResourceConfigShimmed(cfgVal, coreSchema), nil)
}

func TestExternalConnectionCustomizeDiffApprovedCidrs(t *testing.T) {
	approvalOn := map[string]any{
		"conn_learned_cidrs_approval": true,
		"approved_cidrs":              []any{"192.168.102.0/24"},
	}

	t.Run("disabling approval without approved_cidrs recomputes it", func(t *testing.T) {
		diff := externalConnectionTestDiff(t, approvalOn, map[string]any{"conn_learned_cidrs_approval": false})
		require.Contains(t, diff.Attributes, "approved_cidrs.#")
		assert.True(t, diff.Attributes["approved_cidrs.#"].NewComputed)
	})

	t.Run("configured approved_cidrs are planned as written", func(t *testing.T) {
		diff := externalConnectionTestDiff(t, map[string]any{"conn_learned_cidrs_approval": false}, approvalOn)
		require.Contains(t, diff.Attributes, "approved_cidrs.#")
		assert.False(t, diff.Attributes["approved_cidrs.#"].NewComputed)
		assert.Equal(t, "1", diff.Attributes["approved_cidrs.#"].New)
	})

	t.Run("unchanged approval keeps approved_cidrs", func(t *testing.T) {
		diff := externalConnectionTestDiff(t, approvalOn, map[string]any{
			"conn_learned_cidrs_approval":     true,
			"connection_bgp_send_communities": "111:111",
		})
		assert.NotContains(t, diff.Attributes, "approved_cidrs.#")
	})

	t.Run("disabling approval with approved_cidrs fails at plan time", func(t *testing.T) {
		_, err := externalConnectionTestDiffErr(t, approvalOn, map[string]any{
			"conn_learned_cidrs_approval": false,
			"approved_cidrs":              []any{"192.168.102.0/24"},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "'approved_cidrs' requires 'conn_learned_cidrs_approval' to be true")
	})

	t.Run("empty approved_cidrs is allowed without approval", func(t *testing.T) {
		_, err := externalConnectionTestDiffErr(t, approvalOn, map[string]any{
			"conn_learned_cidrs_approval": false,
			"approved_cidrs":              []any{},
		})
		require.NoError(t, err)
	})

	t.Run("approved_cidrs with approval passes", func(t *testing.T) {
		_, err := externalConnectionTestDiffErr(t, map[string]any{"conn_learned_cidrs_approval": false}, approvalOn)
		require.NoError(t, err)
	})
}

// externalConnectionTestUpdateData builds the ResourceData Update sees when
// applying the plan from oldRaw to newRaw.
func externalConnectionTestUpdateData(t *testing.T, oldRaw, newRaw map[string]any) *schema.ResourceData {
	t.Helper()
	r := resourceAviatrixExternalConnection()
	state := externalConnectionTestData(t, oldRaw)
	state.SetId("conn-uuid")
	diff := externalConnectionTestDiff(t, oldRaw, newRaw)

	coreSchema := schema.InternalMap(r.Schema).CoreConfigSchema()
	cfgJSON, err := json.Marshal(externalConnectionTestConfig(newRaw))
	require.NoError(t, err)
	diff.RawConfig, err = ctyjson.Unmarshal(cfgJSON, coreSchema.ImpliedType())
	require.NoError(t, err)

	d, err := schema.InternalMap(r.Schema).Data(state.State(), diff)
	require.NoError(t, err)
	return d
}

func TestExternalConnectionApprovalEnabledWithApprovedCidrs(t *testing.T) {
	off := map[string]any{"conn_learned_cidrs_approval": false}
	on := map[string]any{"conn_learned_cidrs_approval": true, "approved_cidrs": []any{"10.10.0.0/16"}}
	tests := []struct {
		name   string
		oldRaw map[string]any
		newRaw map[string]any
		want   bool
	}{
		{
			name:   "enabling approval with empty approved_cidrs",
			oldRaw: off,
			newRaw: map[string]any{"conn_learned_cidrs_approval": true, "approved_cidrs": []any{}},
			want:   true,
		},
		{
			name:   "enabling approval without approved_cidrs",
			oldRaw: off,
			newRaw: map[string]any{"conn_learned_cidrs_approval": true},
			want:   false,
		},
		{
			name:   "disabling approval",
			oldRaw: on,
			newRaw: map[string]any{"conn_learned_cidrs_approval": false, "approved_cidrs": []any{}},
			want:   false,
		},
		{
			name:   "approval unchanged",
			oldRaw: on,
			newRaw: map[string]any{
				"conn_learned_cidrs_approval":     true,
				"approved_cidrs":                  []any{"10.10.0.0/16"},
				"connection_bgp_send_communities": "111:111",
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := externalConnectionTestUpdateData(t, test.oldRaw, test.newRaw)

			assert.Equal(t, test.want, externalConnectionApprovalEnabledWithApprovedCidrs(d))
		})
	}

	t.Run("empty approved_cidrs is sent as an empty list", func(t *testing.T) {
		d := externalConnectionTestUpdateData(t, off,
			map[string]any{"conn_learned_cidrs_approval": true, "approved_cidrs": []any{}})

		assert.False(t, d.HasChange("approved_cidrs"))
		v := externalConnectionStringSet(d, "approved_cidrs")
		assert.NotNil(t, v)
		assert.Empty(t, v)
	})
}

func TestMarshalExternalConnectionInputBgpRouteSettings(t *testing.T) {
	d := externalConnectionTestData(t, map[string]any{
		"conn_learned_cidrs_approval":              true,
		"approved_cidrs":                           []any{"10.10.0.0/16"},
		"manual_bgp_advertised_cidrs":              []any{"10.20.0.0/16", "10.30.0.0/16"},
		"prepend_as_path":                          []any{"65000", "65000"},
		"connection_bgp_send_communities":          "111:111 444:444",
		"connection_bgp_send_communities_additive": true,
	})

	conn, err := marshalExternalConnectionInput(d)

	require.NoError(t, err)
	assert.Equal(t, []string{"10.10.0.0/16"}, conn.ApprovedCidrs)
	assert.ElementsMatch(t, []string{"10.20.0.0/16", "10.30.0.0/16"}, conn.ManualBgpAdvertisedCidrs)
	assert.Equal(t, []int{65000, 65000}, conn.PrependAsPath)
	assert.Equal(t, "111:111 444:444", conn.BgpSendCommunities)
	assert.True(t, conn.BgpSendCommunitiesAdd)
	assert.False(t, conn.BgpSendCommunitiesBlock)
}

func TestMarshalExternalConnectionInputValidation(t *testing.T) {
	tests := []struct {
		name   string
		raw    map[string]any
		errMsg string
	}{
		{
			name: "block with communities",
			raw: map[string]any{
				"connection_bgp_send_communities":       "111:111",
				"connection_bgp_send_communities_block": true,
			},
			errMsg: "cannot be set when 'connection_bgp_send_communities_block' is true",
		},
		{
			name: "block with additive",
			raw: map[string]any{
				"connection_bgp_send_communities_additive": true,
				"connection_bgp_send_communities_block":    true,
			},
			errMsg: "cannot be set when 'connection_bgp_send_communities_block' is true",
		},
		{
			name:   "approved cidrs without approval",
			raw:    map[string]any{"approved_cidrs": []any{"10.10.0.0/16"}},
			errMsg: "'approved_cidrs' requires 'conn_learned_cidrs_approval' to be true",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := marshalExternalConnectionInput(externalConnectionTestData(t, test.raw))

			require.Error(t, err)
			assert.Contains(t, err.Error(), test.errMsg)
		})
	}
}

func TestExternalConnectionStringSetEmptyIsNotNil(t *testing.T) {
	d := externalConnectionTestData(t, nil)

	v := externalConnectionStringSet(d, "manual_bgp_advertised_cidrs")

	assert.NotNil(t, v)
	assert.Empty(t, v)
}

func TestExpandExternalConnectionPrependASPathEmpty(t *testing.T) {
	d := externalConnectionTestData(t, nil)

	v, err := expandExternalConnectionPrependASPath(d)

	require.NoError(t, err)
	assert.NotNil(t, v)
	assert.Empty(t, v)
}
