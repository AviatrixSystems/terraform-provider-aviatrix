package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Gateway simple struct to hold gateway details
type SpokeVpc struct {
	AccountName             string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                  string `form:"action,omitempty"`
	CID                     string `form:"CID,omitempty"`
	CloudType               int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer               string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                  string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                  string `form:"gw_size,omitempty"`
	VpcID                   string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VpcNet                  string `form:"vpc_net,omitempty" json:"vpc_net,omitempty"`
	VpcRegion               string `form:"vpc_reg,omitempty" json:"vpc_region,omitempty"`
	VpcSize                 string `form:"vpc_size,omitempty" json:"vpc_size,omitempty"`
	EnableNAT               string `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	HASubnet                string `form:"specific_subnet,omitempty"`
	TransitGateway          string `form:"transit_gw,omitempty"`
}

func (c *Client) LaunchSpokeVpc(gateway *SpokeVpc) (error) {
	gateway.CID=c.CID
	gateway.Action="launch_spoke_vpc"
	resp,err := c.Post(c.baseURL, gateway)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) SpokeJoinTransit(gateway *SpokeVpc) (error) {
	enable_ha := ""
	if gateway.HASubnet != "" {
		enable_ha = "yes"
	}
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=spoke_join_transit&spoke_gw=%s&transit_gw=%s&enable_ha=%s", c.CID, gateway.GwName, gateway.TransitGateway, enable_ha)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) SpokeLeaveTransit(gateway *SpokeVpc) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=spoke_leave_transit&spoke_gw=%s", c.CID, gateway.GwName)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) EnableHaSpokeVpc(gateway *SpokeVpc) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=create_peer_ha_gw&vpc_name=%s&specific_subnet=%s", c.CID, gateway.GwName, gateway.HASubnet)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}
