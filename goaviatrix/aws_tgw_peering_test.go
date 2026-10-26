package goaviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPeeringTgwName   = "tgwo-east"
	testPeerLabel        = "partner-tgw"
	testPeerTgwID        = "tgw-0a1b2c3d4e5f60718"
	testPeerAccountID    = "123456789012"
	testPeerRegion       = "us-west-2"
	testOtherPeerTgwID   = "tgw-0f0e0d0c0b0a09080"
	testManagedPeerTgwID = "tgw-0123456789abcdef0"
)

// newAwsTgwPeeringStub answers list_tgw_details with the given attachments and
// records the form of each add or delete submit, keyed by action.
func newAwsTgwPeeringStub(t *testing.T, attachments map[string]any) (*Client, map[string]url.Values) {
	t.Helper()
	submitted := map[string]url.Values{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		assert.NoError(t, r.ParseForm())

		var results any
		switch action := r.Form.Get("action"); action {
		case "check_task_status":
			results = "done"
		case "add_tgw_peering", "delete_tgw_peering":
			submitted[action] = r.Form
			results = "req-1"
		case "list_tgw_details":
			assert.Equal(t, testPeeringTgwName, r.Form.Get("tgw_name"))
			results = map[string]any{"name": testPeeringTgwName, "attachments": attachments}
		default:
			t.Errorf("unexpected action %q", action)
			return
		}
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": results}))
	}))
	t.Cleanup(server.Close)
	return &Client{HTTPClient: server.Client(), CID: "test-cid", baseURL: server.URL}, submitted
}

func externalPeeringAttachment(label string) map[string]any {
	return map[string]any{
		"type":            "peering",
		"peer_tgw_name":   label,
		"external_peer":   true,
		"peer_account_id": testPeerAccountID,
		"peer_region":     testPeerRegion,
	}
}

func TestCreateAwsTgwPeeringSendsExternalPeer(t *testing.T) {
	tests := []struct {
		name    string
		peering AwsTgwPeering
		want    map[string]string
	}{
		{
			name: "external peer",
			peering: AwsTgwPeering{
				TgwName1:      testPeeringTgwName,
				TgwName2:      testPeerLabel,
				PeerTgwID:     testPeerTgwID,
				PeerAccountID: testPeerAccountID,
				PeerRegion:    testPeerRegion,
			},
			want: map[string]string{
				"peer_tgw_id":     testPeerTgwID,
				"peer_account_id": testPeerAccountID,
				"peer_region":     testPeerRegion,
			},
		},
		{
			name:    "managed peer sends no peer fields",
			peering: AwsTgwPeering{TgwName1: testPeeringTgwName, TgwName2: "tgwo-west"},
			want:    map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, submitted := newAwsTgwPeeringStub(t, nil)

			require.NoError(t, client.CreateAwsTgwPeering(context.Background(), &tt.peering))

			form := submitted["add_tgw_peering"]
			require.NotNil(t, form, "add_tgw_peering was not submitted")
			assert.Equal(t, tt.peering.TgwName1, form.Get("tgw_name1"))
			assert.Equal(t, tt.peering.TgwName2, form.Get("tgw_name2"))
			for _, key := range []string{"peer_tgw_id", "peer_account_id", "peer_region"} {
				want, wantSent := tt.want[key]
				_, sent := form[key]
				assert.Equalf(t, wantSent, sent, "%s sent", key)
				assert.Equal(t, want, form.Get(key))
			}
		})
	}
}

func TestDeleteAwsTgwPeeringSendsOnlyPeerTgwID(t *testing.T) {
	client, submitted := newAwsTgwPeeringStub(t, nil)

	require.NoError(t, client.DeleteAwsTgwPeering(context.Background(), &AwsTgwPeering{
		TgwName1:      testPeeringTgwName,
		TgwName2:      testPeerLabel,
		PeerTgwID:     testPeerTgwID,
		PeerAccountID: testPeerAccountID,
		PeerRegion:    testPeerRegion,
	}))

	form := submitted["delete_tgw_peering"]
	require.NotNil(t, form, "delete_tgw_peering was not submitted")
	assert.Equal(t, testPeerTgwID, form.Get("peer_tgw_id"))
	assert.NotContains(t, form, "peer_account_id")
	assert.NotContains(t, form, "peer_region")
}

func TestGetAwsTgwExternalPeering(t *testing.T) {
	tests := []struct {
		name         string
		attachments  map[string]any
		label        string
		knownTgwID   string
		wantExternal bool
		wantTgwID    string
		wantErr      string
	}{
		{
			name:         "external peer",
			attachments:  map[string]any{testPeerTgwID: externalPeeringAttachment(testPeerLabel)},
			label:        testPeerLabel,
			wantExternal: true,
			wantTgwID:    testPeerTgwID,
		},
		{
			name: "managed peer",
			attachments: map[string]any{testManagedPeerTgwID: map[string]any{
				"type":          "peering",
				"peer_tgw_name": "tgwo-west",
			}},
			label: "tgwo-west",
		},
		{
			name:        "label not found",
			attachments: map[string]any{testPeerTgwID: externalPeeringAttachment("other-label")},
			label:       testPeerLabel,
		},
		{
			name: "shared label without a peer TGW ID",
			attachments: map[string]any{
				testPeerTgwID:      externalPeeringAttachment(testPeerLabel),
				testOtherPeerTgwID: externalPeeringAttachment(testPeerLabel),
			},
			label:   testPeerLabel,
			wantErr: "more than one external peering",
		},
		{
			name: "shared label with a peer TGW ID",
			attachments: map[string]any{
				testPeerTgwID:      externalPeeringAttachment(testPeerLabel),
				testOtherPeerTgwID: externalPeeringAttachment(testPeerLabel),
			},
			label:        testPeerLabel,
			knownTgwID:   testOtherPeerTgwID,
			wantExternal: true,
			wantTgwID:    testOtherPeerTgwID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newAwsTgwPeeringStub(t, tt.attachments)
			peering := &AwsTgwPeering{TgwName1: testPeeringTgwName, TgwName2: tt.label, PeerTgwID: tt.knownTgwID}

			external, err := client.GetAwsTgwExternalPeering(context.Background(), peering)

			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantExternal, external)
			if !tt.wantExternal {
				assert.Equal(t, tt.knownTgwID, peering.PeerTgwID)
				assert.Empty(t, peering.PeerAccountID)
				assert.Empty(t, peering.PeerRegion)
				return
			}
			assert.Equal(t, tt.wantTgwID, peering.PeerTgwID)
			assert.Equal(t, testPeerAccountID, peering.PeerAccountID)
			assert.Equal(t, testPeerRegion, peering.PeerRegion)
		})
	}
}
