package goaviatrix

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SpokeExternalDeviceConnPh1RemoteIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("ha_enabled") {
		return false
	}

	ip := d.Get("remote_gateway_ip").(string)
	ipList := strings.Split(ip, ",")
	haip := d.Get("backup_remote_gateway_ip").(string)
	o, n := d.GetChange("phase1_remote_identifier")
	haEnabled := d.Get("ha_enabled").(bool)

	ph1RemoteIdListOld := ExpandStringList(o.([]interface{}))
	ph1RemoteIdListNew := ExpandStringList(n.([]interface{}))

	if len(ph1RemoteIdListOld) != 0 && len(ph1RemoteIdListNew) != 0 {
		if haEnabled {
			if len(ph1RemoteIdListNew) != 2 || len(ph1RemoteIdListOld) != 2 {
				if len(ph1RemoteIdListNew) == 1 && len(ph1RemoteIdListOld) == 1 {
					return ph1RemoteIdListOld[0] == ipList[0] && ph1RemoteIdListNew[0] == ipList[0]
				}
				return false
			}
			return ph1RemoteIdListOld[0] == ipList[0] && ph1RemoteIdListNew[0] == ipList[0] &&
				strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && strings.TrimSpace(ph1RemoteIdListNew[1]) == haip
		} else {
			if len(ph1RemoteIdListNew) == 1 && len(ph1RemoteIdListOld) == 1 {
				return ph1RemoteIdListOld[0] == ipList[0] && ph1RemoteIdListNew[0] == ipList[0]
			} else if len(ph1RemoteIdListNew) == 2 && len(ph1RemoteIdListOld) == 2 && len(ipList) == 2 {
				return strings.TrimSpace(ph1RemoteIdListOld[0]) == strings.TrimSpace(ipList[0]) &&
					strings.TrimSpace(ph1RemoteIdListOld[1]) == strings.TrimSpace(ipList[1]) &&
					strings.TrimSpace(ph1RemoteIdListNew[0]) == strings.TrimSpace(ipList[0]) &&
					strings.TrimSpace(ph1RemoteIdListNew[1]) == strings.TrimSpace(ipList[1])
			} else {
				return false
			}
		}
	}

	if !haEnabled {
		if len(ph1RemoteIdListOld) == 1 && ph1RemoteIdListOld[0] == ipList[0] && len(ph1RemoteIdListNew) == 0 {
			return true
		}
		if len(ph1RemoteIdListOld) == 2 && strings.TrimSpace(ph1RemoteIdListOld[0]) == strings.TrimSpace(ipList[0]) &&
			strings.TrimSpace(ph1RemoteIdListOld[1]) == strings.TrimSpace(ipList[1]) &&
			len(ph1RemoteIdListNew) == 0 {
			return true
		}
	}

	if haEnabled {
		if len(ph1RemoteIdListOld) == 2 && ph1RemoteIdListOld[0] == ipList[0] && strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && len(ph1RemoteIdListNew) == 0 {
			return true
		}
		if len(ph1RemoteIdListOld) == 1 && ph1RemoteIdListOld[0] == ipList[0] && len(ph1RemoteIdListNew) == 0 {
			return true
		}
	}

	return false
}

func (c *Client) EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn *ExternalDeviceConn, prependASPath []string) error {
	action := "edit_transit_connection_as_path_prepend"
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
