package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Gateway simple struct to hold gateway details
type TransitVpc struct {
	AccountName             string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                  string `form:"action,omitempty"`
	CID                     string `form:"CID,omitempty"`
	CloudType               int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer               string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                  string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                  string `form:"gw_size,omitempty"`
	VpcID                   string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VpcNet                  string `form:"vpc_net,omitempty" json:"vpc_net,omitempty"`
	HASubnet                string `form:"specific_subnet,omitempty"`
	VpcRegion               string `form:"vpc_reg,omitempty" json:"vpc_region,omitempty"`
	VpcSize                 string `form:"vpc_size,omitempty" json:"vpc_size,omitempty"`
}

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) (error) {
	gateway.CID=c.CID
	gateway.Action="launch_transit_vpc"
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

func (c *Client) EnableHaTransitVpc(gateway *TransitVpc) (error) {
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
