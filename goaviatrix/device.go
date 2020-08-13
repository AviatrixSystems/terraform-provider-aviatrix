package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Device represents a branch router used in CloudWAN
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
	BranchState        string               `form:"-" map:"-" json:"registered"`
	PrimaryInterface   string               `form:"-" map:"-" json:"wan_if_primary"`
	PrimaryInterfaceIP string               `form:"-" map:"-" json:"wan_if_primary_public_ip"`
	ConnectionName     string               `form:"-" map:"-" json:"conn_name"`
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
	resp, err := c.PostFile(c.baseURL, d.toMap(), files)
	if err != nil {
		return errors.New("HTTP Post register_cloudwan_device failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body register_cloudwan_device failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode register_cloudwan_device failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API register_cloudwan_device Post failed: " + data.Reason)
	}
	return nil
}

// toMap converts the struct to a map[string]string
// The 'map' tags on the struct tell us what the key name should be.
func (br *Device) toMap() map[string]string {
	out := make(map[string]string)
	v := reflect.ValueOf(br).Elem()
	tag := "map"
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" && tagv != "-" {
			out[tagv] = v.Field(i).String()
		}
	}
	return out
}

func (c *Client) GetDevice(d *Device) (*Device, error) {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "list_cloudwan_devices_summary",
	})
	if err != nil {
		return nil, errors.New("HTTP POST list_cloudwan_devices_summary failed: " + err.Error())
	}

	type Resp struct {
		Return  bool     `json:"return"`
		Results []Device `json:"results"`
		Reason  string   `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body list_cloudwan_devices_summary failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_cloudwan_devices_summary failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_cloudwan_devices_summary Post failed: " + data.Reason)
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
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "list_cloudwan_devices_summary",
	})
	if err != nil {
		return "", errors.New("HTTP POST list_cloudwan_devices_summary failed: " + err.Error())
	}

	type Resp struct {
		Return  bool     `json:"return"`
		Results []Device `json:"results"`
		Reason  string   `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return "", errors.New("Reading response body list_cloudwan_devices_summary failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return "", errors.New("Json Decode list_cloudwan_devices_summary failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return "", errors.New("Rest API list_cloudwan_devices_summary Post failed: " + data.Reason)
	}

	for _, device := range data.Results {
		if device.ConnectionName == connName {
			return device.Name, nil
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
	resp, err := c.PostFile(c.baseURL, d.toMap(), files)
	if err != nil {
		return errors.New("HTTP Post update_cloudwan_device_info failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body update_cloudwan_device_info failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode update_cloudwan_device_info failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API update_cloudwan_device_info Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeregisterDevice(d *Device) error {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
		Name   string `form:"device_name"`
	}{
		CID:    c.CID,
		Action: "deregister_cloudwan_device",
		Name:   d.Name,
	})
	if err != nil {
		return errors.New("HTTP POST deregister_cloudwan_device failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body deregister_cloudwan_device failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode deregister_cloudwan_device failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API deregister_cloudwan_device Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) ConfigureDeviceInterfaces(config *DeviceInterfaceConfig) error {
	availableInterfaces, err := c.GetBranchRouterInterfaces(&Device{Name: config.DeviceName})
	if err != nil {
		return err
	}

	if !Contains(availableInterfaces, config.PrimaryInterface) {
		return fmt.Errorf("branch router does not have the given primary interface '%s'. "+
			"Possible interfaces are [%s]", config.PrimaryInterface, strings.Join(availableInterfaces, ", "))
	}

	resp, err := c.Post(c.baseURL, struct {
		CID                string `form:"CID"`
		Action             string `form:"action"`
		Name               string `form:"device_name"`
		PrimaryInterface   string `form:"wan_primary_if"`
		PrimaryInterfaceIP string `form:"wan_primary_ip"`
	}{
		CID:                c.CID,
		Action:             "config_cloudwan_device_wan_interfaces",
		Name:               config.DeviceName,
		PrimaryInterface:   config.PrimaryInterface,
		PrimaryInterfaceIP: config.PrimaryInterfaceIP,
	})
	if err != nil {
		return errors.New("HTTP POST config_cloudwan_device_wan_interfaces failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body config_cloudwan_device_wan_interfaces failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode config_cloudwan_device_wan_interfaces failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API config_cloudwan_device_wan_interfaces Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetBranchRouterInterfaces(device *Device) ([]string, error) {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
		Name   string `form:"device_name"`
	}{
		CID:    c.CID,
		Action: "get_cloudwan_device_wan_interfaces",
		Name:   device.Name,
	})
	if err != nil {
		return nil, errors.New("HTTP POST get_cloudwan_device_wan_interfaces failed: " + err.Error())
	}

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
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_cloudwan_device_wan_interfaces failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_cloudwan_device_wan_interfaces failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API get_cloudwan_device_wan_interfaces Post failed: " + data.Reason)
	}
	var interfaces []string

	for _, v := range data.Results {
		interfaces = append(interfaces, v.Interface)
	}

	return interfaces, nil
}
