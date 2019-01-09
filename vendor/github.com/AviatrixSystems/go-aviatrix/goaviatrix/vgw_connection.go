package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// VGWConn simple struct to hold VGW Connection details
type VGWConn struct {
	Action        string `form:"action,omitempty"`
	BgpLocalAsNum string `form:"bgp_local_as_num,omitempty" json:"bgp_local_as_num,omitempty"`
	BgpVGWId      string `form:"vgw_id,omitempty" json:"bgp_vgw_id,omitempty"`
	CID           string `form:"CID,omitempty"`
	ConnName      string `form:"connection_name,omitempty" json:"name,omitempty"`
	GwName        string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	VPCId         string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
}

type VGWConnListResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type VGWConnList struct {
	Return  bool      `json:"return"`
	Results []VGWConn `json:"results"`
	Reason  string    `json:"reason"`
}

func (c *Client) CreateVGWConn(vgwConn *VGWConn) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=connect_transit_gw_to_vgw&vpc_id=%s&connection_name=%s&"+
		"transit_gw=%s&vgw_id=%s&bgp_local_as_number=%s", c.CID, vgwConn.VPCId, vgwConn.ConnName, vgwConn.GwName,
		vgwConn.BgpVGWId, vgwConn.BgpLocalAsNum)
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

func (c *Client) GetVGWConn(vgwConn *VGWConn) (*VGWConn, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_vgw_connections", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	data := VGWConnListResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if !data.Return {
		return nil, errors.New(data.Reason)
	}

	vgwConnList := data.Results
	for i := range vgwConnList {
		if vgwConnList[i] == vgwConn.ConnName {
			return vgwConn, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateVGWConn(vgwConn *VGWConn) error {
	return nil
}

func (c *Client) DeleteVGWConn(vgwConn *VGWConn) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=disconnect_transit_gw_from_vgw&vpc_id=%s&connection_name=%s",
		c.CID, vgwConn.VPCId, vgwConn.ConnName)
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
