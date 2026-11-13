package goaviatrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newEdcEditClient returns a Client whose controller records each request's
// path and form and answers with the given return value.
func newEdcEditClient(t *testing.T, ok bool) (*Client, *[]string, *[]map[string]string) {
	t.Helper()
	var paths []string
	var forms []map[string]string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm failed: %v", err)
			return
		}
		form := make(map[string]string)
		for k, v := range r.Form {
			if len(v) > 0 {
				form[k] = v[0]
			}
		}
		paths = append(paths, r.URL.Path)
		forms = append(forms, form)
		resp := map[string]any{"return": ok, "results": "done", "reason": "refused"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("Encode failed: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return &Client{
		HTTPClient: server.Client(),
		CID:        "test-cid",
		baseURL:    server.URL + "/v2/api",
	}, &paths, &forms
}

func TestEditExternalDeviceConnEnableIkev2(t *testing.T) {
	client, paths, forms := newEdcEditClient(t, true)

	edit := &EditExternalDeviceConnEnableIkev2{GwName: "transit1", ConnectionName: "conn1"}
	require.NoError(t, client.EditExternalDeviceConnEnableIkev2(edit))

	require.Len(t, *forms, 1)
	assert.Equal(t, "/v2/api", (*paths)[0])
	assert.Equal(t, map[string]string{
		"action":       "edit_external_device_conn",
		"CID":          "test-cid",
		"gateway_name": "transit1",
		"conn_name":    "conn1",
		"enable_ikev2": "true",
	}, (*forms)[0])
}

func TestEditExternalDeviceConnEnableIkev2Error(t *testing.T) {
	client, _, _ := newEdcEditClient(t, false)

	err := client.EditExternalDeviceConnEnableIkev2(&EditExternalDeviceConnEnableIkev2{
		GwName: "transit1", ConnectionName: "conn1",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "refused")
}
