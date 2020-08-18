package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
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
	resp, err := c.Post(c.baseURL, attachment)
	if err != nil {
		return errors.New("HTTP Post attach_cloudwan_device_to_aws_tgw failed: " + err.Error())
	}

	var data APIResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body attach_cloudwan_device_to_aws_tgw failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode attach_cloudwan_device_to_aws_tgw failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API attach_cloudwan_device_to_aws_tgw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetDeviceAwsTgwAttachment(tgwAttachment *DeviceAwsTgwAttachment) (*DeviceAwsTgwAttachment, error) {
	tgwAttachment.CID = c.CID
	tgwAttachment.Action = "list_tgw_details"

	resp, err := c.Post(c.baseURL, tgwAttachment)
	if err != nil {
		return nil, errors.New("HTTP Post list_tgw_details failed: " + err.Error())
	}

	var data TgwAttachmentResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body list_tgw_details failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_tgw_details failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_tgw_details Post failed: " + data.Reason)
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

	return &DeviceAwsTgwAttachment{
		ConnectionName:          tgwAttachment.ConnectionName,
		DeviceName:              tgwAttachment.DeviceName,
		AwsTgwName:              tgwAttachment.AwsTgwName,
		DeviceAsn:               deviceAttachment.AwsSideAsn,
		SecurityDomainName:      deviceAttachment.SecurityDomainName,
		EnableGlobalAccelerator: strconv.FormatBool(deviceAttachment.EnableGlobalAccelerator),
	}, nil
}
