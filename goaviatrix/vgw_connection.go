package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

// VGWConn simple struct to hold VGW Connection details
type VGWConn struct {
	Action         string `form:"action,omitempty"`
	BgpLocalAsNum  string `form:"bgp_local_asn_num,omitempty" json:"bgp_local_asn_num,omitempty"`
	BgpVGWId       string `form:"vgw_id,omitempty" json:"bgp_vgw_id,omitempty"`
	BgpVGWAccount  string `form:"bgp_vgw_account_name,omitempty" json:"bgp_vgw_account,omitempty"`
	BgpVGWRegion   string `form:"bgp_vgw_region,omitempty" json:"bgp_vgw_region,omitempty"`
	CID            string `form:"CID,omitempty"`
	ConnName       string `form:"connection_name,omitempty" json:"name,omitempty"`
	GwName         string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	VPCId          string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	ManualBGPCidrs []string
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

type VGWConnDetailResp struct {
	Return  bool          `json:"return"`
	Results VGWConnDetail `json:"results"`
	Reason  string        `json:"reason"`
}

type VGWConnDetail struct {
	Connections ConnectionDetail `json:"connections"`
}

type ConnectionDetail struct {
	ConnName       []string `json:"name"`
	GwName         string   `json:"gw_name"`
	VPCId          []string `json:"vpc_id"`
	BgpVGWId       string   `json:"bgp_vgw_id"`
	BgpVGWAccount  string   `json:"bgp_vgw_account"`
	BgpVGWRegion   string   `json:"bgp_vgw_region"`
	BgpLocalAsNum  string   `json:"bgp_local_asn_number"`
	ManualBGPCidrs []string `json:"conn_bgp_manual_advertise_cidrs"`
}

type VGWConnEnableAdvertiseTransitCidrResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

type VGWConnBgpManualSpokeAdvertisedNetworksResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
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
	connectTransitGwToVgw.Add("bgp_vgw_account_name", vgwConn.BgpVGWAccount)
	connectTransitGwToVgw.Add("bgp_vgw_region", vgwConn.BgpVGWRegion)
	connectTransitGwToVgw.Add("bgp_local_as_number", vgwConn.BgpLocalAsNum)
	Url.RawQuery = connectTransitGwToVgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get connect_transit_gw_to_vgw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode connect_transit_gw_to_vgw failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vgw_connections failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disconnect_transit_gw_from_vgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disconnect_transit_gw_from_vgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetVGWConnDetail(vgwConn *VGWConn) (*VGWConn, error) {
	params := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"vpc_id":    vgwConn.VPCId,
		"conn_name": vgwConn.ConnName,
	}
	var data VGWConnDetailResp
	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}
	if data.Results.Connections.ConnName[0] != "" {
		vgwConn.VPCId = data.Results.Connections.VPCId[0]
		vgwConn.GwName = data.Results.Connections.GwName
		vgwConn.BgpVGWId = data.Results.Connections.BgpVGWId
		vgwConn.BgpVGWAccount = data.Results.Connections.BgpVGWAccount
		vgwConn.BgpVGWRegion = data.Results.Connections.BgpVGWRegion
		vgwConn.BgpLocalAsNum = data.Results.Connections.BgpLocalAsNum
		vgwConn.ManualBGPCidrs = data.Results.Connections.ManualBGPCidrs
		return vgwConn, nil
	}
	return nil, ErrNotFound
}
