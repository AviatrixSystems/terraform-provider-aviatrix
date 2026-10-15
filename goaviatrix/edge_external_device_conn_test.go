package goaviatrix

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEdgeExternalDeviceConnOmitsBgpAdminShutdown guards the request payload.
// This struct is sent as JSON, so a field carrying only a form tag - or none -
// is still marshaled, under its Go name. BgpAdminShutdown exists here only to
// mirror ExternalDeviceConn field for field and is never set on an edge
// connection, so it must not reach the controller.
func TestEdgeExternalDeviceConnOmitsBgpAdminShutdown(t *testing.T) {
	body, err := json.Marshal(&EdgeExternalDeviceConn{
		Action:         "transit_connect_external_device",
		ConnectionName: "conn1",
		GwName:         "gw1",
	})
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(body, &payload))
	assert.NotContains(t, payload, "BgpAdminShutdown")
	assert.NotContains(t, payload, "bgp_admin_shutdown")
}
