package goaviatrix

import (
	"strings"
)

// Device represents a device used in CloudWAN
type Device struct {
	Action             string `form:"action,omitempty" json:"-"`
	CID                string `form:"CID,omitempty" json:"-"`
	Name               string `form:"device_name,omitempty" json:"rgw_name"`
	PublicIP           string `form:"public_ip,omitempty" json:"hostname"`
	Username           string `form:"username,omitempty" json:"username"`
	HostOS             string `form:"host_os,omitempty" json:"host_os"`
	Description        string `form:"description,omitempty" json:"description"`
	CheckReason        string `form:"-" json:"check_reason"`
	PrimaryInterface   string `form:"-" json:"wan_if_primary"`
	PrimaryInterfaceIP string `form:"-" json:"wan_if_primary_public_ip"`
	ConnectionName     string `form:"-" json:"conn_name"`
	SoftwareVersion    string `form:"-" json:"software_version"`
	IsCaag             bool   `form:"-" json:"is_caag"`
}

func (c *Client) GetDeviceName(connName string) (string, error) {
	type Resp struct {
		Return  bool     `json:"return"`
		Results []Device `json:"results"`
		Reason  string   `json:"reason"`
	}
	var data Resp
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_cloudwan_devices_summary",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	for _, device := range data.Results {
		// ConnectionName is actually a CSV list of connection names
		conns := strings.Split(device.ConnectionName, ",")
		for _, c := range conns {
			if c == connName {
				return device.Name, nil
			}
		}
	}

	return "", ErrNotFound
}
