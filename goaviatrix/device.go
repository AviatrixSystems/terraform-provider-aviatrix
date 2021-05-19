package goaviatrix

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Device represents a device used in CloudWAN
type Device struct {
	Action             string               `form:"action,omitempty" map:"action" json:"-"`
	CID                string               `form:"CID,omitempty" map:"CID" json:"-"`
	Name               string               `form:"device_name,omitempty" map:"device_name" json:"rgw_name"`
	PublicIP           string               `form:"public_ip,omitempty" map:"public_ip" json:"hostname"`
	Username           string               `form:"username,omitempty" map:"username" json:"username"`
	KeyFile            string               `form:"-" map:"-" json:"-"`
	Password           string               `form:"password,omitempty" map:"password" json:"-"`
	HostOS             string               `form:"host_os,omitempty" map:"host_os" json:"host_os"`
	SshPort            int                  `form:"-" map:"-" json:"ssh_port"`
	SshPortStr         string               `form:"port,omitempty" map:"port" json:"-"`
	Address1           string               `form:"addr_1,omitempty" map:"addr_1" json:"-"`
	Address2           string               `form:"addr_2,omitempty" map:"addr_2" json:"-"`
	City               string               `form:"city,omitempty" map:"city" json:"-"`
	State              string               `form:"state,omitempty" map:"state" json:"-"`
	Country            string               `form:"country,omitempty" map:"country" json:"-"`
	ZipCode            string               `form:"zipcode,omitempty" map:"zipcode" json:"-"`
	Description        string               `form:"description,omitempty" map:"description" json:"description"`
	Address            GetDeviceRespAddress `form:"-" map:"-" json:"address"`
	CheckReason        string               `form:"-" map:"-" json:"check_reason"`
	DeviceState        string               `form:"-" map:"-" json:"registered"`
	PrimaryInterface   string               `form:"-" map:"-" json:"wan_if_primary"`
	PrimaryInterfaceIP string               `form:"-" map:"-" json:"wan_if_primary_public_ip"`
	ConnectionName     string               `form:"-" map:"-" json:"conn_name"`
	SoftwareVersion    string               `form:"-" map:"-" json:"software_version"`
	IsCaag             bool                 `form:"-" map:"-" json:"is_caag"`
}

type DeviceInterfaceConfig struct {
	DeviceName         string
	PrimaryInterface   string
	PrimaryInterfaceIP string
}

type GetDeviceRespAddress struct {
	Address1 string `json:"addr_1"`
	Address2 string `json:"addr_2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	ZipCode  string `json:"zipcode"`
}

func (c *Client) RegisterDevice(d *Device) error {
	d.Action = "register_cloudwan_device"
	d.CID = c.CID
	files := []File{
		{
			Path:      d.KeyFile,
			ParamName: "private_key_file",
		},
	}
	return c.PostFileAPI(toMap(d), files, BasicCheck)
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

	foundDevice.Address1 = foundDevice.Address.Address1
	foundDevice.Address2 = foundDevice.Address.Address2
	foundDevice.City = foundDevice.Address.City
	foundDevice.State = foundDevice.Address.State
	foundDevice.Country = foundDevice.Address.Country
	foundDevice.ZipCode = foundDevice.Address.ZipCode

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

func (c *Client) UpdateDevice(d *Device) error {
	d.Action = "update_cloudwan_device_info"
	d.CID = c.CID

	files := []File{
		{
			Path:      d.KeyFile,
			ParamName: "private_key_file",
		},
	}
	return c.PostFileAPI(toMap(d), files, BasicCheck)
}

func (c *Client) DeregisterDevice(d *Device) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "deregister_cloudwan_device",
		"device_name": d.Name,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) ConfigureDeviceInterfaces(config *DeviceInterfaceConfig) error {
	availableInterfaces, err := c.GetDeviceInterfaces(&Device{Name: config.DeviceName})
	if err != nil {
		return err
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

func (c *Client) GetDeviceInterfaces(device *Device) ([]string, error) {
	type Result struct {
		Interface string `json:"interface"`
		IP        string `json:"ip"`
	}
	type Resp struct {
		Return  bool     `json:"return"`
		Results []Result `json:"results"`
		Reason  string   `json:"reason"`
	}
	var data Resp
	form := map[string]string{
		"CID":         c.CID,
		"action":      "get_cloudwan_device_wan_interfaces",
		"device_name": device.Name,
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	var interfaces []string
	for _, v := range data.Results {
		interfaces = append(interfaces, v.Interface)
	}
	return interfaces, nil
}
