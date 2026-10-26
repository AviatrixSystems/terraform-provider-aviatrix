package aviatrix

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func externalTgwPeeringConfig(remoteCidrs ...any) map[string]any {
	cfg := map[string]any{
		"tgw_name1":       testTgwPeeringOurTgw,
		"tgw_name2":       testTgwPeeringLabel,
		"peer_tgw_id":     testTgwPeeringPeerID,
		"peer_account_id": testTgwPeeringAccountID,
		"peer_region":     testTgwPeeringRegion,
	}
	if len(remoteCidrs) > 0 {
		cfg["remote_cidrs"] = remoteCidrs
	}
	return cfg
}

func TestResourceAviatrixAWSTgwPeeringCreateSetsRemoteCidrs(t *testing.T) {
	controller := fakeTgwPeeringController{
		peeredNames: []string{testTgwPeeringLabel},
		attachments: map[string]any{testTgwPeeringPeerID: map[string]any{
			"type":            "peering",
			"peer_tgw_name":   testTgwPeeringLabel,
			"external_peer":   true,
			"peer_account_id": testTgwPeeringAccountID,
			"peer_region":     testTgwPeeringRegion,
		}},
	}
	d := schema.TestResourceDataRaw(t, resourceAviatrixAWSTgwPeering().Schema,
		externalTgwPeeringConfig("10.50.0.0/16", "10.60.0.0/16"))

	diags := resourceAviatrixAWSTgwPeeringCreate(context.Background(), d, controller.client(t))
	require.False(t, diags.HasError(), "%v", diags)

	assert.Equal(t, 1, controller.creates)
	require.Len(t, controller.cidrUpdates, 1)
	update := controller.cidrUpdates[0]
	assert.Equal(t, testTgwPeeringOurTgw, update.Get("tgw_name"))
	assert.Equal(t, testTgwPeeringLabel, update.Get("peer_tgw_name"))
	assert.ElementsMatch(t, []string{"10.50.0.0/16", "10.60.0.0/16"}, strings.Split(update.Get("cidr_list"), ","))
	assert.ElementsMatch(t, []string{"10.50.0.0/16", "10.60.0.0/16"}, getStringSet(d, "remote_cidrs"))
}

func TestResourceAviatrixAWSTgwPeeringCreateWithoutRemoteCidrs(t *testing.T) {
	controller := fakeTgwPeeringController{
		peeredNames: []string{testTgwPeeringLabel},
		attachments: externalTgwPeeringAttachments(),
	}
	d := schema.TestResourceDataRaw(t, resourceAviatrixAWSTgwPeering().Schema, externalTgwPeeringConfig())

	diags := resourceAviatrixAWSTgwPeeringCreate(context.Background(), d, controller.client(t))
	require.False(t, diags.HasError(), "%v", diags)

	assert.Equal(t, 1, controller.creates)
	assert.Empty(t, controller.cidrUpdates)
}

// If the peering is created but the CIDRs can't be set, the peering must still
// reach the state, with the CIDRs it really has, so Terraform taints it rather
// than losing track of it.
func TestResourceAviatrixAWSTgwPeeringCreateRemoteCidrsFail(t *testing.T) {
	controller := fakeTgwPeeringController{
		peeredNames: []string{testTgwPeeringLabel},
		attachments: map[string]any{testTgwPeeringPeerID: map[string]any{
			"type":            "peering",
			"peer_tgw_name":   testTgwPeeringLabel,
			"external_peer":   true,
			"peer_account_id": testTgwPeeringAccountID,
			"peer_region":     testTgwPeeringRegion,
		}},
		cidrUpdateErr: "route table is busy",
	}
	client := controller.client(t)
	r := resourceAviatrixAWSTgwPeering()
	ctx := context.Background()
	diff, err := r.SimpleDiff(ctx, nil, terraform.NewResourceConfigRaw(externalTgwPeeringConfig("10.50.0.0/16")), client)
	require.NoError(t, err)

	newState, diags := r.Apply(ctx, nil, diff, client)

	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "route table is busy")
	assert.Contains(t, diags[0].Summary, "terraform untaint")
	assert.Equal(t, 1, controller.creates)
	require.NotNil(t, newState, "the created peering must be saved to state")
	assert.Equal(t, testTgwPeeringOurTgw+"~"+testTgwPeeringLabel, newState.ID)
	assert.Equal(t, testTgwPeeringPeerID, newState.Attributes["peer_tgw_id"])
	assert.Equal(t, 0, countSetAttr(newState, "remote_cidrs"), "state must hold the CIDRs the Controller has, not the planned ones")
}

func TestResourceAviatrixAWSTgwPeeringUpdateRemoteCidrs(t *testing.T) {
	tests := []struct {
		name     string
		oldCidrs []any
		newCidrs []any
		wantList []string
	}{
		{
			name:     "replace some CIDRs",
			oldCidrs: []any{"10.50.0.0/16", "10.60.0.0/16"},
			newCidrs: []any{"10.60.0.0/16", "10.70.0.0/16"},
			wantList: []string{"10.60.0.0/16", "10.70.0.0/16"},
		},
		{
			name:     "add the first CIDRs",
			newCidrs: []any{"10.50.0.0/16"},
			wantList: []string{"10.50.0.0/16"},
		},
		{
			name:     "remove every CIDR",
			oldCidrs: []any{"10.50.0.0/16"},
			wantList: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := fakeTgwPeeringController{
				peeredNames: []string{testTgwPeeringLabel},
				attachments: externalTgwPeeringAttachments(),
			}
			client := controller.client(t)
			r := resourceAviatrixAWSTgwPeering()
			old := schema.TestResourceDataRaw(t, r.Schema, externalTgwPeeringConfig(tt.oldCidrs...))
			old.SetId(testTgwPeeringOurTgw + "~" + testTgwPeeringLabel)
			oldState := old.State()

			ctx := context.Background()
			diff, err := r.SimpleDiff(ctx, oldState, terraform.NewResourceConfigRaw(externalTgwPeeringConfig(tt.newCidrs...)), client)
			require.NoError(t, err)
			require.NotNil(t, diff)
			assert.False(t, diff.RequiresNew(), "changing remote_cidrs must not replace the peering")

			newState, diags := r.Apply(ctx, oldState, diff, client)
			require.False(t, diags.HasError(), "%v", diags)

			require.Len(t, controller.cidrUpdates, 1)
			update := controller.cidrUpdates[0]
			require.True(t, update.Has("cidr_list"), "cidr_list must be sent even when empty")
			if tt.wantList == nil {
				assert.Empty(t, update.Get("cidr_list"))
			} else {
				assert.ElementsMatch(t, tt.wantList, strings.Split(update.Get("cidr_list"), ","))
			}
			assert.Equal(t, testTgwPeeringOurTgw+"~"+testTgwPeeringLabel, newState.ID)
			assert.Equal(t, len(tt.wantList), countSetAttr(newState, "remote_cidrs"))
		})
	}
}

func countSetAttr(s *terraform.InstanceState, key string) int {
	n := 0
	for k := range s.Attributes {
		if strings.HasPrefix(k, key+".") && k != key+".#" {
			n++
		}
	}
	return n
}

func TestResourceAviatrixAWSTgwPeeringRemoteCidrsValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{
			name:   "valid CIDRs on an external peering",
			config: externalTgwPeeringConfig("10.50.0.0/16", "192.168.1.0/24"),
		},
		{
			name: "remote_cidrs on a managed peering",
			config: map[string]any{
				"tgw_name1":    testTgwPeeringOurTgw,
				"tgw_name2":    "tgwo-west",
				"remote_cidrs": []any{"10.50.0.0/16"},
			},
			wantErr: "peer_tgw_id",
		},
		{
			name:    "host bits set",
			config:  externalTgwPeeringConfig("10.50.0.1/16"),
			wantErr: "remote_cidrs",
		},
		{
			name:    "not a CIDR",
			config:  externalTgwPeeringConfig("10.50.0.0"),
			wantErr: "remote_cidrs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := resourceAviatrixAWSTgwPeering().Validate(terraform.NewResourceConfigRaw(tt.config))
			if tt.wantErr == "" {
				assert.False(t, diags.HasError(), "%v", diags)
				return
			}
			require.True(t, diags.HasError())
			var msgs []string
			for _, d := range diags {
				msgs = append(msgs, d.Summary+" "+d.Detail)
			}
			assert.Contains(t, strings.Join(msgs, "\n"), tt.wantErr)
		})
	}
}
