package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
)

// VGWConn simple struct to hold VGW Connection details
type VGWConn struct {
	Action        string `form:"action,omitempty"`
	BgpLocalAsNum string `form:"bgp_local_asn_num,omitempty" json:"bgp_local_asn_num,omitempty"`
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for connect_transit_gw_to_vgw") + err.Error())
	}
	connectTransitGwToVgw := url.Values{}
	connectTransitGwToVgw.Add("CID", c.CID)
	connectTransitGwToVgw.Add("action", "connect_transit_gw_to_vgw")
	connectTransitGwToVgw.Add("vpc_id", vgwConn.VPCId)
	connectTransitGwToVgw.Add("connection_name", vgwConn.ConnName)
	connectTransitGwToVgw.Add("transit_gw", vgwConn.GwName)
	connectTransitGwToVgw.Add("vgw_id", vgwConn.BgpVGWId)
	connectTransitGwToVgw.Add("bgp_local_as_number", vgwConn.BgpLocalAsNum)
	Url.RawQuery = connectTransitGwToVgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get connect_transit_gw_to_vgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode connect_transit_gw_to_vgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API connect_transit_gw_to_vgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetVGWConn(vgwConn *VGWConn) (*VGWConn, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_vgw_connections") + err.Error())
	}
	listVgwConnections := url.Values{}
	listVgwConnections.Add("CID", c.CID)
	listVgwConnections.Add("action", "list_vgw_connections")
	Url.RawQuery = listVgwConnections.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_vgw_connections failed: " + err.Error())
	}
	data := VGWConnListResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vgw_connections failed: " + err.Error())
	}

	if !data.Return {
		return nil, errors.New("Rest API list_vgw_connections Get failed: " + data.Reason)
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disconnect_transit_gw_from_vgw") + err.Error())
	}
	disconnectTransitGwFromVgw := url.Values{}
	disconnectTransitGwFromVgw.Add("CID", c.CID)
	disconnectTransitGwFromVgw.Add("action", "disconnect_transit_gw_from_vgw")
	disconnectTransitGwFromVgw.Add("vpc_id", vgwConn.VPCId)
	disconnectTransitGwFromVgw.Add("connection_name", vgwConn.ConnName)
	Url.RawQuery = disconnectTransitGwFromVgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disconnect_transit_gw_from_vgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disconnect_transit_gw_from_vgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disconnect_transit_gw_from_vgw Get failed: " + data.Reason)
	}
	return nil
}
