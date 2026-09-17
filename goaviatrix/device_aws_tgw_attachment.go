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
	Async                   bool   `form:"async,omitempty"`
}

func (b *DeviceAwsTgwAttachment) ID() string {
	return b.ConnectionName + "~" + b.DeviceName + "~" + b.AwsTgwName
}

func (c *Client) CreateDeviceAwsTgwAttachment(attachment *DeviceAwsTgwAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_aws_tgw"
	attachment.CID = c.CID
	attachment.Async = true
	return c.PostAsyncAPI(attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetDeviceAwsTgwAttachment(tgwAttachment *DeviceAwsTgwAttachment) (*DeviceAwsTgwAttachment, error) {
	tgwVpcAttachment := &AwsTgwVpcAttachment{
		TgwName: tgwAttachment.AwsTgwName,
		VpcID:   tgwAttachment.ConnectionName,
	}
	tgwAttachmentInfo, err := c.GetAwsTgwAttachmentInfo(tgwVpcAttachment)
	if err != nil {
		return nil, err
	}
	// aws_side_asn can return as either string or int from API
	if len(tgwAttachmentInfo.AwsSideAsnRaw) != 0 {
		// First try as string
		var asnString string
		var asnInt int
		err = json.Unmarshal(tgwAttachmentInfo.AwsSideAsnRaw, &asnString)
		if err != nil {
			// String failed, must be int
			err = json.Unmarshal(tgwAttachmentInfo.AwsSideAsnRaw, &asnInt)
			if err != nil {
				return nil, fmt.Errorf("json decode get_tgw_attachment_details aws_side_asn field failed: aws_side_asn = %s: %w",
					string(tgwAttachmentInfo.AwsSideAsnRaw), err)
			}
			asnString = strconv.Itoa(asnInt)
		}
		tgwAttachmentInfo.AwsSideAsn = asnString
	}

	return &DeviceAwsTgwAttachment{
		ConnectionName:          tgwAttachment.ConnectionName,
		DeviceName:              tgwAttachment.DeviceName,
		AwsTgwName:              tgwAttachment.AwsTgwName,
		DeviceAsn:               tgwAttachmentInfo.AwsSideAsn,
		SecurityDomainName:      tgwAttachmentInfo.SecurityDomainName,
		EnableGlobalAccelerator: strconv.FormatBool(tgwAttachmentInfo.EnableGlobalAccelerator),
	}, nil
}
