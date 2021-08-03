package goaviatrix

import (
	"fmt"
	"strings"
	"time"
)

type DeviceTag struct {
	Name          string `form:"tag_name,omitempty"`
	Config        string `form:"custom_cfg,omitempty"`
	Devices       []string
	DevicesString string `form:"include_device_list,omitempty"`
	CID           string `form:"CID"`
	Action        string `form:"action"`
}

func (c *Client) CreateDeviceTag(deviceTag *DeviceTag) error {
	// Create the tag
	deviceTag.CID = c.CID
	deviceTag.Action = "add_cloudwan_configtag"
	err := c.PostAPI(deviceTag.Action, deviceTag, BasicCheck)
	if err != nil {
		return err
	}

	// Set the tag config
	if err := c.UpdateDeviceTagConfig(deviceTag); err != nil {
		return err
	}

	// Attach the devices to the tag
	if err := c.AttachDeviceTag(deviceTag); err != nil {
		return err
	}

	// Commit the tag config to the devices
	if err := c.CommitDeviceTag(deviceTag); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDeviceTag(brt *DeviceTag) (*DeviceTag, error) {
	// Check if a tag exists with the given name
	form := map[string]string{
		"action":              "list_cloudwan_configtag_names",
		"CID":                 c.CID,
		"tag_name":            brt.Name,
		"custom_cfg":          brt.Config,
		"include_device_list": brt.DevicesString,
	}
	type Resp struct {
		Return  bool     `json:"return,omitempty"`
		Results []string `json:"results,omitempty"`
		Reason  string   `json:"reason,omitempty"`
	}
	var data Resp
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if !Contains(data.Results, brt.Name) {
		return nil, ErrNotFound
	}

	// Get the details for the tag
	form["action"] = "get_cloudwan_configtag_details"
	type DetailsResults struct {
		TagName         string   `json:"gtag_name"`
		AttachedDevices []string `json:"rgw_name"`
		Config          string   `json:"custom_cfg"`
	}
	type DetailsResp struct {
		Return  bool           `json:"return,omitempty"`
		Results DetailsResults `json:"results,omitempty"`
		Reason  string         `json:"reason,omitempty"`
	}
	var detailsData DetailsResp
	err = c.GetAPI(&detailsData, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	brt.Devices = detailsData.Results.AttachedDevices
	brt.Config = detailsData.Results.Config
	return brt, nil
}

func (c *Client) UpdateDeviceTagConfig(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "edit_cloudwan_configtag"
	return c.PostAPI(brt.Action, brt, BasicCheck)
}

func (c *Client) AttachDeviceTag(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "attach_devices_to_cloudwan_configtag"
	brt.DevicesString = strings.Join(brt.Devices, ", ")
	return c.PostAPI(brt.Action, brt, BasicCheck)
}

func (c *Client) CommitDeviceTag(brt *DeviceTag) error {
	tries, maxTries := 0, 5
	backoff := 15 * time.Second
	var err error

	for tries < maxTries {
		err = c.commitDeviceTagOnce(brt)
		if err != nil {
			tries++
			if tries < maxTries {
				time.Sleep(backoff)
				backoff *= 2
			}
			continue
		}
		break
	}

	if err != nil {
		return fmt.Errorf("tried to commit device tag %d times but could not succeed: %v", maxTries, err)
	}

	return nil
}

func (c *Client) commitDeviceTagOnce(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "commit_cloudwan_configtag_to_devices"
	return c.PostAPI(brt.Action, brt, BasicCheck)
}

func (c *Client) DeleteDeviceTag(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "delete_cloudwan_configtag"
	return c.PostAPI(brt.Action, brt, BasicCheck)
}
