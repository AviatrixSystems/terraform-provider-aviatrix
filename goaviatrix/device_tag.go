package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
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
	resp, err := c.Post(c.baseURL, deviceTag)

	if err != nil {
		return errors.New("HTTP Post add_cloudwan_configtag failed: " + err.Error())
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
		return errors.New("Reading response body add_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode add_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API add_cloudwan_configtag Post failed: " + data.Reason)
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
	brt.CID = c.CID
	brt.Action = "list_cloudwan_configtag_names"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return nil, errors.New("HTTP Post get_cloudwan_configtag_details failed: " + err.Error())
	}

	type Resp struct {
		Return  bool     `json:"return,omitempty"`
		Results []string `json:"results,omitempty"`
		Reason  string   `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body list_cloudwan_configtag_names failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_cloudwan_configtag_names failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_cloudwan_configtag_names Post failed: " + data.Reason)
	}

	if !Contains(data.Results, brt.Name) {
		return nil, ErrNotFound
	}

	// Get the details for the tag
	brt.Action = "get_cloudwan_configtag_details"
	resp, err = c.Post(c.baseURL, brt)

	if err != nil {
		return nil, errors.New("HTTP Post get_cloudwan_configtag_details failed: " + err.Error())
	}

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
	b = bytes.Buffer{}
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_cloudwan_configtag_details failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&detailsData); err != nil {
		return nil, errors.New("Json Decode get_cloudwan_configtag_details failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !detailsData.Return {
		return nil, errors.New("Rest API get_cloudwan_configtag_details Post failed: " + detailsData.Reason)
	}

	brt.Devices = detailsData.Results.AttachedDevices
	brt.Config = detailsData.Results.Config
	return brt, nil
}

func (c *Client) UpdateDeviceTagConfig(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "edit_cloudwan_configtag"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post edit_cloudwan_configtag failed: " + err.Error())
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
		return errors.New("Reading response body edit_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode edit_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API edit_cloudwan_configtag Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) AttachDeviceTag(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "attach_devices_to_cloudwan_configtag"
	brt.DevicesString = strings.Join(brt.Devices, ", ")
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post attach_devices_to_cloudwan_configtag failed: " + err.Error())
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
		return errors.New("Reading response body attach_devices_to_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode attach_devices_to_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API attach_devices_to_cloudwan_configtag Post failed: " + data.Reason)
	}

	return nil
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
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post commit_cloudwan_configtag_to_devices failed: " + err.Error())
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
		return errors.New("Reading response body commit_cloudwan_configtag_to_devices failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode commit_cloudwan_configtag_to_devices failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API commit_cloudwan_configtag_to_devices Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DeleteDeviceTag(brt *DeviceTag) error {
	brt.CID = c.CID
	brt.Action = "delete_cloudwan_configtag"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post delete_cloudwan_configtag failed: " + err.Error())
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
		return errors.New("Reading response body delete_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode delete_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API delete_cloudwan_configtag Post failed: " + data.Reason)
	}

	return nil
}
