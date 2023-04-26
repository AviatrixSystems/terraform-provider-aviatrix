package goaviatrix

import (
	"context"
	"io"
	"os"
)

type EdgeNEODevice struct {
	Action                 string                  `json:"action,omitempty"`
	CID                    string                  `json:"CID,omitempty"`
	AccountName            string                  `json:"account_name,omitempty"`
	DeviceName             string                  `json:"device_name,omitempty"`
	SerialNumber           string                  `json:"serial,omitempty"`
	HardwareModel          string                  `json:"hardware_model,omitempty"`
	Network                []*EdgeNEODeviceNetwork `json:"network,omitempty"`
	DownloadConfigFile     bool
	ConfigFileDownloadPath string
}

type EdgeNEODeviceNetwork struct {
	InterfaceName string   `json:"interface,omitempty"`
	EnableDhcp    bool     `json:"dhcp,omitempty"`
	GatewayIp     string   `json:"gateway,omitempty"`
	Ipv4Cidr      string   `json:"ipv4cidr,omitempty"`
	DnsServerIps  []string `json:"dns,omitempty"`
	ProxyServerIp string   `json:"proxy,omitempty"`
}

type EdgeNEODeviceResp struct {
	AccountName      string                  `json:"account_name"`
	DeviceName       string                  `json:"name"`
	DeviceId         string                  `json:"deviceId"`
	SerialNumber     string                  `json:"serial"`
	HardwareModel    string                  `json:"hardwareId"`
	Network          []*EdgeNEODeviceNetwork `json:"interfaces"`
	ConnectionStatus string                  `json:"connectionStatus"`
}

type EdgeNEODeviceListResp struct {
	Return  bool                `json:"return"`
	Results []EdgeNEODeviceResp `json:"results"`
	Reason  string              `json:"reason"`
}

func (c *Client) OnboardEdgeNEODevice(ctx context.Context, edgeNEODevice *EdgeNEODevice) error {
	edgeNEODevice.Action = "onboard_edge_csp_device"
	edgeNEODevice.CID = c.CID

	err := c.PostAPIContext2(ctx, nil, edgeNEODevice.Action, edgeNEODevice, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeNEODevice(ctx context.Context, accountName, deviceName string) (*EdgeNEODeviceResp, error) {
	form := map[string]string{
		"action":       "list_edge_csp_devices",
		"CID":          c.CID,
		"account_name": accountName,
	}

	var data EdgeNEODeviceListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeNEODeviceList := data.Results
	for _, edgeNEODevice := range edgeNEODeviceList {
		if edgeNEODevice.DeviceName == deviceName {
			return &edgeNEODevice, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeNEODevice(ctx context.Context, edgeNEODevice *EdgeNEODevice) error {
	edgeNEODevice.Action = "update_edge_csp_device"
	edgeNEODevice.CID = c.CID

	err := c.PostAPIContext2(ctx, nil, edgeNEODevice.Action, edgeNEODevice, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteEdgeNEODevice(ctx context.Context, accountName, serialNumber string) error {
	form := map[string]string{
		"action":       "delete_edge_csp_device",
		"CID":          c.CID,
		"account_name": accountName,
		"serial":       serialNumber,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) DownloadEdgeNEOConfigFile(ctx context.Context, edgeNEODevice *EdgeNEODevice) error {
	edgeNEODevice.Action = "get_edge_csp_bootstrap_usb"
	edgeNEODevice.CID = c.CID

	resp, err := c.PostAPIContext2Download(ctx, edgeNEODevice.Action, edgeNEODevice, BasicCheck)
	if err != nil {
		return err
	}

	fileName := edgeNEODevice.ConfigFileDownloadPath + "/" + edgeNEODevice.SerialNumber + "-bootstrap-config.img"

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, resp)
	if err != nil {
		return err
	}

	return nil
}
