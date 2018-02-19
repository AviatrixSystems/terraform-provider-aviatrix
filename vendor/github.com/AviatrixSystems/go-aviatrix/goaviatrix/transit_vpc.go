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
	Subnet                  string `form:"public_subnet,omitempty" json:"vpc_net,omitempty"`
	HASubnet                string `form:"ha_subnet,omitempty"`
	VpcRegion               string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                 string `form:"gw_size,omitempty" json:"gw_size,omitempty"`
	TagList                 string `form:"tags,omitempty"`
}

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) (error) {
	gateway.CID=c.CID
	gateway.Action="create_transit_gw"
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
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=enable_transit_ha&gw_name=%s&public_subnet=%s", c.CID, gateway.GwName, gateway.HASubnet)
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
