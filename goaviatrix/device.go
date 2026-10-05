package goaviatrix

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
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

type DeviceWanInterface struct {
	Interface string `json:"interface"`
	IP        string `json:"ip"`
}

type DeviceWanInterfaceResp struct {
	Return  bool                 `json:"return"`
	Results []DeviceWanInterface `json:"results"`
	Reason  string               `json:"reason"`
}

type DeviceInterfaceConfig struct {
	DeviceName         string
	PrimaryInterface   string
	PrimaryInterfaceIP string
}

func (c *Client) GetDevice(d *Device) (*Device, error) {
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
		return nil, err
	}
	var foundDevice *Device
	for _, device := range data.Results {
		if device.Name == d.Name {
			foundDevice = &device
			break
		}
	}
	if foundDevice == nil {
		log.Errorf("Could not find Aviatrix device %s", d.Name)
		return nil, ErrNotFound
	}

	return foundDevice, nil
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

func (c *Client) GetDeviceInterfaces(deviceName string) (*[]DeviceWanInterface, error) {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "get_cloudwan_device_wan_interfaces",
		"device_name": deviceName,
	}

	var data DeviceWanInterfaceResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) ConfigureDeviceInterfaces(config *DeviceInterfaceConfig) error {
	deviceInterfaces, err := c.GetDeviceInterfaces(config.DeviceName)
	if err != nil {
		return err
	}

	var availableInterfaces []string
	for _, v := range *deviceInterfaces {
		availableInterfaces = append(availableInterfaces, v.Interface)
	}

	if !Contains(availableInterfaces, config.PrimaryInterface) {
		return fmt.Errorf("device does not have the given primary interface '%s'. "+
			"Possible interfaces are [%s]", config.PrimaryInterface, strings.Join(availableInterfaces, ", "))
	}

	form := map[string]string{
		"CID":            c.CID,
		"action":         "config_cloudwan_device_wan_interfaces",
		"device_name":    config.DeviceName,
		"wan_primary_if": config.PrimaryInterface,
		"wan_primary_ip": config.PrimaryInterfaceIP,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}
