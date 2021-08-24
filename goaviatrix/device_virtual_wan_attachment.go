package goaviatrix

import (
	"fmt"
	"strings"
)

type DeviceVirtualWanAttachment struct {
	ConnectionName string `form:"connection_name"`
	DeviceName     string `form:"device_name"`
	AccountName    string `form:"account_name"`
	ResourceGroup  string `form:"arm_resource_group"`
	HubName        string `form:"virtual_wan_hub_name"`
	DeviceAsn      string `form:"external_device_as_number"`
	Action         string `form:"action"`
	CID            string `form:"CID"`
}

func (c *Client) CreateDeviceVirtualWanAttachment(attachment *DeviceVirtualWanAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_virtual_wan"
	attachment.CID = c.CID
	return c.PostAPI(attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetDeviceVirtualWanAttachment(attachment *DeviceVirtualWanAttachment) (*DeviceVirtualWanAttachment, error) {
	deviceName, err := c.GetDeviceName(attachment.ConnectionName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not get device name: %v", err)
	}

	vpcID, err := c.GetDeviceAttachmentVpcID(attachment.ConnectionName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not get device attachment VPC id: %v", err)
	}

	form := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"vpc_id":    vpcID,
		"conn_name": attachment.ConnectionName,
	}

	var data Site2CloudConnDetailResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err = c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	return &DeviceVirtualWanAttachment{
		ConnectionName: attachment.ConnectionName,
		DeviceName:     deviceName,
		AccountName:    data.Results.Connections.AzureAccountName,
		ResourceGroup:  data.Results.Connections.AzureResourceGroup,
		HubName:        data.Results.Connections.AzureVhubName,
		DeviceAsn:      data.Results.Connections.BgpRemoteASN,
	}, nil
}
