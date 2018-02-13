 package goaviatrix

import (
	"fmt"
	"encoding/json"
	"errors"
)

// VGWConn simple struct to hold VGW Connection details
type VGWConn struct {
	Action                  string `form:"action,omitempty"`
	BgpLocalAsNum           string `form:"bgp_local_as_num,omitempty" json:"bgp_local_as_num,omitempty"`
	BgpVGWId                string `form:"bgp_vgw_id,omitempty" json:"bgp_vgw_id,omitempty"`
	CID                     string `form:"CID,omitempty"`
	ConnName                string `form:"connection_name,omitempty" json:"name,omitempty"`
	GwName                  string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	VPCId                   string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
}

type VGWConnListResp struct {
	Return  bool   `json:"return"`
	Results []VGWConn `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateVGWConn(vgw_conn *VGWConn) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=connect_transit_gw_to_vgw&vpc_alias=%s&connection_name=%s&transit_gw=%s&vgw_alias=%s&bgp_local_as_number=%s", c.CID, vgw_conn.VPCId, vgw_conn.ConnName, vgw_conn.GwName, vgw_conn.BgpVGWId, vgw_conn.BgpLocalAsNum)
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

func (c *Client) GetVGWConn(vgw_conn *VGWConn) (*VGWConn, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_vgw_connections", c.CID)
	resp,err := c.Get(path, nil)
		if err != nil {
		return nil, err
	}
	var data VGWConnListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	vgw_connlist:= data.Results
	for i := range vgw_connlist {
		if vgw_connlist[i].ConnName == vgw_conn.ConnName {
			return &vgw_connlist[i], nil
		}
	}
	return nil, fmt.Errorf("VGW Connection %s not found", vgw_conn.ConnName)
}

func (c *Client) UpdateVGWConn(vgw_conn *VGWConn) (error) {
	return nil
}

func (c *Client) DeleteVGWConn(vgw_conn *VGWConn) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=disconnect_transit_gw_from_vgw&vpc_alias=%s&connection_name=%s", c.CID, vgw_conn.VPCId, vgw_conn.ConnName)
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

