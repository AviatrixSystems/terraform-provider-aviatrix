package goaviatrix

import "strings"

func (c *Client) EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn *ExternalDeviceConn, prependASPath []string) error {
	action := "edit_spoke_connection_as_path_prepend"
	return c.PostAPI(action, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		GatewayName    string `form:"gateway_name"`
		ConnectionName string `form:"connection_name"`
		PrependASPath  string `form:"connection_as_path_prepend"`
	}{
		CID:            c.CID,
		Action:         action,
		GatewayName:    externalDeviceConn.GwName,
		ConnectionName: externalDeviceConn.ConnectionName,
		PrependASPath:  strings.Join(prependASPath, ","),
	}, BasicCheck)
}
