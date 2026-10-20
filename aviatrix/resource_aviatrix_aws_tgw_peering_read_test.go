package aviatrix

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	testTgwPeeringOurTgw    = "tgwo-east"
	testTgwPeeringLabel     = "partner-tgw"
	testTgwPeeringPeerID    = "tgw-0a1b2c3d4e5f60718"
	testTgwPeeringAccountID = "123456789012"
	testTgwPeeringRegion    = "us-west-2"
)

// fakeTgwPeeringController answers the calls the peering Read makes and
// counts list_tgw_details calls.
type fakeTgwPeeringController struct {
	peeredNames  []string
	attachments  map[string]any
	detailsCalls int
}

func (f *fakeTgwPeeringController) client(t *testing.T) *goaviatrix.Client {
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

func (f *fakeTgwPeeringController) respond(t *testing.T, params url.Values) map[string]any {
	switch action := params.Get("action"); action {
	case "get_api_token":
		return map[string]any{"return": true, "results": map[string]any{"api_token": "token"}}
	case "login":
		return map[string]any{"return": true, "CID": "cid"}
	case "list_peered_tgw_names":
		return map[string]any{"return": true, "results": f.peeredNames}
	case "list_tgw_details":
		f.detailsCalls++
		return map[string]any{"return": true, "results": map[string]any{
			"name":        params.Get("tgw_name"),
			"attachments": f.attachments,
		}}
	default:
		t.Errorf("unexpected action %q", action)
		return map[string]any{"return": false, "reason": "unexpected action"}
	}
}

func externalTgwPeeringAttachments() map[string]any {
	return map[string]any{testTgwPeeringPeerID: map[string]any{
		"type":            "peering",
		"peer_tgw_name":   testTgwPeeringLabel,
		"external_peer":   true,
		"peer_account_id": testTgwPeeringAccountID,
		"peer_region":     testTgwPeeringRegion,
	}}
}

func TestResourceAviatrixAWSTgwPeeringRead(t *testing.T) {
	externalState := map[string]any{
		"tgw_name1":       testTgwPeeringOurTgw,
		"tgw_name2":       testTgwPeeringLabel,
		"peer_tgw_id":     testTgwPeeringPeerID,
		"peer_account_id": testTgwPeeringAccountID,
		"peer_region":     testTgwPeeringRegion,
	}
	externalID := testTgwPeeringOurTgw + "~" + testTgwPeeringLabel
	tests := []struct {
		name            string
		id              string
		state           map[string]any
		controller      fakeTgwPeeringController
		wantID          string
		wantPeer        map[string]string
		wantDetailsCall bool
	}{
		{
			name:  "import of an external peering reads the peer details",
			id:    externalID,
			state: map[string]any{},
			controller: fakeTgwPeeringController{
				peeredNames: []string{testTgwPeeringLabel},
				attachments: externalTgwPeeringAttachments(),
			},
			wantID:          externalID,
			wantPeer:        map[string]string{"peer_tgw_id": testTgwPeeringPeerID, "peer_account_id": testTgwPeeringAccountID, "peer_region": testTgwPeeringRegion},
			wantDetailsCall: true,
		},
		{
			name:  "refresh of an external peering reads the peer details",
			id:    externalID,
			state: externalState,
			controller: fakeTgwPeeringController{
				peeredNames: []string{testTgwPeeringLabel},
				attachments: externalTgwPeeringAttachments(),
			},
			wantID:          externalID,
			wantPeer:        map[string]string{"peer_tgw_id": testTgwPeeringPeerID, "peer_account_id": testTgwPeeringAccountID, "peer_region": testTgwPeeringRegion},
			wantDetailsCall: true,
		},
		{
			name:            "refresh of a managed peering skips the details lookup",
			id:              testTgwPeeringOurTgw + "~tgwo-west",
			state:           map[string]any{"tgw_name1": testTgwPeeringOurTgw, "tgw_name2": "tgwo-west"},
			controller:      fakeTgwPeeringController{peeredNames: []string{"tgwo-west"}},
			wantID:          testTgwPeeringOurTgw + "~tgwo-west",
			wantPeer:        map[string]string{"peer_tgw_id": "", "peer_account_id": "", "peer_region": ""},
			wantDetailsCall: false,
		},
		{
			name:  "label now peered with a different external TGW",
			id:    externalID,
			state: externalState,
			controller: fakeTgwPeeringController{
				peeredNames: []string{testTgwPeeringLabel},
				attachments: map[string]any{"tgw-0ffffffffffffffff": externalTgwPeeringAttachments()[testTgwPeeringPeerID]},
			},
			wantID:          "",
			wantDetailsCall: true,
		},
		{
			name:  "label now names a managed peer",
			id:    externalID,
			state: externalState,
			controller: fakeTgwPeeringController{
				peeredNames: []string{testTgwPeeringLabel},
				attachments: map[string]any{},
			},
			wantID:          "",
			wantDetailsCall: true,
		},
		{
			name:            "peering gone from the controller",
			id:              externalID,
			state:           externalState,
			controller:      fakeTgwPeeringController{peeredNames: []string{"someone-else"}},
			wantID:          "",
			wantDetailsCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAviatrixAWSTgwPeering().Schema, tt.state)
			d.SetId(tt.id)

			require.NoError(t, resourceAviatrixAWSTgwPeeringRead(d, tt.controller.client(t)))

			assert.Equal(t, tt.wantID, d.Id())
			assert.Equal(t, tt.wantDetailsCall, tt.controller.detailsCalls > 0)
			for key, want := range tt.wantPeer {
				assert.Equalf(t, want, getString(d, key), "%s", key)
			}
		})
	}
}

func TestResourceAviatrixAWSTgwPeeringImportRejectsBadID(t *testing.T) {
	for _, id := range []string{"tgwo-east", "tgwo-east~", "~partner-tgw", "a~b~c"} {
		t.Run(id, func(t *testing.T) {
			controller := fakeTgwPeeringController{}
			d := schema.TestResourceDataRaw(t, resourceAviatrixAWSTgwPeering().Schema, map[string]any{})
			d.SetId(id)

			err := resourceAviatrixAWSTgwPeeringRead(d, controller.client(t))

			require.ErrorContains(t, err, "invalid import ID")
		})
	}
}
