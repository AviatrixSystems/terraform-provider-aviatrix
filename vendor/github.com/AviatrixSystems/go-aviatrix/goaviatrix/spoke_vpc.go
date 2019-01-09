package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	TransitGateway string `form:"transit_gw,omitempty"`
	TagList        string `form:"tags,omitempty"`
}

func (c *Client) LaunchSpokeVpc(spoke *SpokeVpc) error {
	spoke.CID = c.CID
	spoke.Action = "create_spoke_gw"
	resp, err := c.Post(c.baseURL, spoke)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) SpokeJoinTransit(spoke *SpokeVpc) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=attach_spoke_to_transit_gw&spoke_gw=%s&transit_gw=%s",
		c.CID, spoke.GwName, spoke.TransitGateway)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) SpokeLeaveTransit(spoke *SpokeVpc) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=detach_spoke_from_transit_gw&spoke_gw=%s", c.CID,
		spoke.GwName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		if strings.Contains(data.Reason, "has not joined to any transit") {
			log.Printf("[INFO] spoke VPC is already left from transit VPC %s", data.Reason)
			return nil
		}
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) EnableHaSpokeVpc(spoke *SpokeVpc) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=enable_spoke_ha&gw_name=%s&public_subnet=%s", c.CID,
		spoke.GwName, spoke.HASubnet)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		if strings.Contains(data.Reason, "HA GW already exists") {
			log.Printf("[INFO] HA is already enabled %s", data.Reason)
			return nil
		}
		log.Printf("[ERROR] Enabling HA failed with error %s", data.Reason)
		return errors.New(data.Reason)
	}
	return nil
}
