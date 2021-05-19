package goaviatrix

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type DeviceAwsTgwAttachment struct {
	ConnectionName          string `map:"connection_name" form:"connection_name"`
	DeviceName              string `map:"device_name" form:"device_name"`
	AwsTgwName              string `map:"tgw_name" form:"tgw_name"`
	DeviceAsn               string `map:"external_device_as_number" form:"external_device_as_number"`
	SecurityDomainName      string `map:"route_domain_name" form:"route_domain_name"`
	EnableGlobalAccelerator string `map:"enable_global_accelerator" form:"enable_global_accelerator"`
	Action                  string `map:"action" form:"action"`
	CID                     string `map:"CID" form:"CID"`
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
	tgwAttachment.CID = c.CID
	tgwAttachment.Action = "list_tgw_details"
	var data TgwAttachmentResp
	err := c.GetAPI(&data, tgwAttachment.Action, toMap(tgwAttachment), BasicCheck)
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
