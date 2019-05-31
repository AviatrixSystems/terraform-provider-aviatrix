package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

// Spoke gateway simple struct to hold spoke details
type SpokeVpc struct {
	AccountName    string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	CloudType      int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer      string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName         string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize         string `form:"gw_size,omitempty"`
	VpcID          string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VnetRsrcGrp    string `form:"vnet_and_resource_group_names,omitempty"`
	Subnet         string `form:"public_subnet,omitempty" json:"public_subnet,omitempty"`
	VpcRegion      string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize        string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	EnableNAT      string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	HASubnet       string `form:"ha_subnet,omitempty"`
	SingleAzHa     string `form:"single_az_ha,omitempty"`
	TransitGateway string `form:"transit_gw,omitempty"`
	TagList        string `form:"tags,omitempty"`
}

func (c *Client) LaunchSpokeVpc(spoke *SpokeVpc) error {
	spoke.CID = c.CID
	spoke.Action = "create_spoke_gw"
	resp, err := c.Post(c.baseURL, spoke)
	if err != nil {
		return errors.New("HTTP Post create_spoke_gw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode create_spoke_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API create_spoke_gw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SpokeJoinTransit(spoke *SpokeVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_spoke_to_transit_gw") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "attach_spoke_to_transit_gw")
	attachSpokeToTransitGw.Add("spoke_gw", spoke.GwName)
	attachSpokeToTransitGw.Add("transit_gw", spoke.TransitGateway)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get attach_spoke_to_transit_gw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode attach_spoke_to_transit_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API attach_spoke_to_transit_gw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SpokeLeaveTransit(spoke *SpokeVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for detach_spoke_from_transit_gw") + err.Error())
	}
	detachSpokeFromTransitGw := url.Values{}
	detachSpokeFromTransitGw.Add("CID", c.CID)
	detachSpokeFromTransitGw.Add("action", "detach_spoke_from_transit_gw")
	detachSpokeFromTransitGw.Add("spoke_gw", spoke.GwName)
	Url.RawQuery = detachSpokeFromTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get detach_spoke_from_transit_gw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode detach_spoke_from_transit_gw failed: " + err.Error())
	}
	if !data.Return {
		if strings.Contains(data.Reason, "has not joined to any transit") {
			log.Printf("[INFO] spoke VPC is already left from transit VPC %s", data.Reason)
			return nil
		}
		return errors.New("Rest API detach_spoke_from_transit_gw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableHaSpokeVpc(spoke *SpokeVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_spoke_ha") + err.Error())
	}
	enableSpokeHa := url.Values{}
	enableSpokeHa.Add("CID", c.CID)
	enableSpokeHa.Add("action", "enable_spoke_ha")
	enableSpokeHa.Add("gw_name", spoke.GwName)
	if spoke.CloudType == 1 || spoke.CloudType == 8 {
		enableSpokeHa.Add("public_subnet", spoke.HASubnet)
	} else if spoke.CloudType == 4 {
		enableSpokeHa.Add("new_zone", spoke.HASubnet)
	} else {
		return errors.New("Invalid cloud type")
	}
	Url.RawQuery = enableSpokeHa.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_spoke_ha failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_spoke_ha failed: " + err.Error())
	}
	if !data.Return {
		if strings.Contains(data.Reason, "HA GW already exists") {
			log.Printf("[INFO] HA is already enabled %s", data.Reason)
			return nil
		}
		log.Printf("[ERROR] Enabling HA failed with error %s", data.Reason)
		return errors.New("Rest API enable_spoke_ha Get failed: " + data.Reason)
	}
	return nil
}
