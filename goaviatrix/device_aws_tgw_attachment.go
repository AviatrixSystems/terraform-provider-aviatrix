package goaviatrix

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type DeviceAwsTgwAttachment struct {
	ConnectionName          string `form:"connection_name"`
	DeviceName              string `form:"device_name"`
	AwsTgwName              string `form:"tgw_name"`
	DeviceAsn               string `form:"external_device_as_number"`
	SecurityDomainName      string `form:"route_domain_name"`
	EnableGlobalAccelerator string `form:"enable_global_accelerator"`
	Action                  string `form:"action"`
	CID                     string `form:"CID"`
}

func (b *DeviceAwsTgwAttachment) ID() string {
	return b.ConnectionName + "~" + b.DeviceName + "~" + b.AwsTgwName
}

func (c *Client) CreateDeviceAwsTgwAttachment(attachment *DeviceAwsTgwAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_aws_tgw"
	attachment.CID = c.CID
	return c.PostAPI(attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetDeviceAwsTgwAttachment(tgwAttachment *DeviceAwsTgwAttachment) (*DeviceAwsTgwAttachment, error) {
	form := map[string]string{
		"action":                    "list_tgw_details",
		"CID":                       c.CID,
		"connection_name":           tgwAttachment.ConnectionName,
		"device_name":               tgwAttachment.DeviceName,
		"tgw_name":                  tgwAttachment.AwsTgwName,
		"external_device_as_number": tgwAttachment.DeviceAsn,
		"route_domain_name":         tgwAttachment.SecurityDomainName,
		"enable_global_accelerator": tgwAttachment.EnableGlobalAccelerator,
	}
	var data TgwAttachmentResp
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	var deviceAttachment AttachmentInfo
	var found bool
	for _, attachment := range data.Results.Attachments {
		if attachment.VpcName == tgwAttachment.ConnectionName {
			deviceAttachment = attachment
			found = true
		}
	}
	if !found {
		return nil, ErrNotFound
	}

	// aws_side_asn can return as either string or int from API
	if len(deviceAttachment.AwsSideAsnRaw) != 0 {
		// First try as string
		var asnString string
		var asnInt int
		err = json.Unmarshal(deviceAttachment.AwsSideAsnRaw, &asnString)
		if err != nil {
			// String failed, must be int
			err = json.Unmarshal(deviceAttachment.AwsSideAsnRaw, &asnInt)
			if err != nil {
				return nil, fmt.Errorf("json decode list_tgw_details aws_side_asn field failed: aws_side_asn = %s: %v", string(deviceAttachment.AwsSideAsnRaw), err)
			}
			asnString = strconv.Itoa(asnInt)
		}
		deviceAttachment.AwsSideAsn = asnString
	}

	return &DeviceAwsTgwAttachment{
		ConnectionName:          tgwAttachment.ConnectionName,
		DeviceName:              tgwAttachment.DeviceName,
		AwsTgwName:              tgwAttachment.AwsTgwName,
		DeviceAsn:               deviceAttachment.AwsSideAsn,
		SecurityDomainName:      deviceAttachment.SecurityDomainName,
		EnableGlobalAccelerator: strconv.FormatBool(deviceAttachment.EnableGlobalAccelerator),
	}, nil
}
