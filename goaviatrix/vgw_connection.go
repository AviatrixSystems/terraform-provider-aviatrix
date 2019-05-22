package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
)

// VGWConn simple struct to hold VGW Connection details
type VGWConn struct {
	Action                       string `form:"action,omitempty"`
	BgpLocalAsNum                string `form:"bgp_local_asn_num,omitempty" json:"bgp_local_asn_num,omitempty"`
	BgpVGWId                     string `form:"vgw_id,omitempty" json:"bgp_vgw_id,omitempty"`
	CID                          string `form:"CID,omitempty"`
	ConnName                     string `form:"connection_name,omitempty" json:"name,omitempty"`
	GwName                       string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	VPCId                        string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	EnableAdvertiseTransitCidr   bool
	BgpManualSpokeAdvertiseCidrs string `form:"cidr,omitempty"`
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
	ConnName                     []string   `json:"name"`
	GwName                       []string   `json:"gw_name"`
	VPCId                        []string   `json:"vpc_id"`
	BgpVGWId                     []string   `json:"bgp_vgw_id"`
	BgpLocalAsNum                []string   `json:"bgp_local_asn_number"`
	AdvertiseTransitCidr         string     `json:"advertise_transit_cidr"`
	BgpManualSpokeAdvertiseCidrs [][]string `json:"bgp_manual_spoke_advertise_cidrs"`
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

func (c *Client) GetVGWConnDetail(vgwConn *VGWConn) (*VGWConn, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for get_site2cloud_conn_detail") + err.Error())
	}
	listVgwConnections := url.Values{}
	listVgwConnections.Add("CID", c.CID)
	listVgwConnections.Add("action", "get_site2cloud_conn_detail")
	listVgwConnections.Add("vpc_id", vgwConn.VPCId)
	listVgwConnections.Add("conn_name", vgwConn.ConnName)
	Url.RawQuery = listVgwConnections.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get get_site2cloud_conn_detail failed: " + err.Error())
	}
	var data VGWConnDetailResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_site2cloud_conn_detail failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API get_site2cloud_conn_detail Get failed: " + data.Reason)
	}

	if data.Results.Connections.ConnName[0] != "" {
		vgwConn.VPCId = data.Results.Connections.VPCId[0]
		vgwConn.GwName = data.Results.Connections.GwName[0]
		vgwConn.BgpVGWId = data.Results.Connections.BgpVGWId[0]
		vgwConn.BgpLocalAsNum = data.Results.Connections.BgpLocalAsNum[0]
		if data.Results.Connections.AdvertiseTransitCidr == "yes" {
			vgwConn.EnableAdvertiseTransitCidr = true
		} else if data.Results.Connections.AdvertiseTransitCidr == "no" {
			vgwConn.EnableAdvertiseTransitCidr = false
		}
		if len(data.Results.Connections.BgpManualSpokeAdvertiseCidrs) != 0 {
			bgpMSAN := ""
			for i := range data.Results.Connections.BgpManualSpokeAdvertiseCidrs[0] {
				if i == 0 {
					bgpMSAN = bgpMSAN + data.Results.Connections.BgpManualSpokeAdvertiseCidrs[0][i]
				} else {
					bgpMSAN = bgpMSAN + "," + data.Results.Connections.BgpManualSpokeAdvertiseCidrs[0][i]
				}
			}
			vgwConn.BgpManualSpokeAdvertiseCidrs = bgpMSAN
		}
		return vgwConn, nil
	}

	return nil, ErrNotFound
}

func (c *Client) EnableAdvertiseTransitCidr(vgwConn *VGWConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_advertise_transit_cidr") + err.Error())
	}
	enableAdvertiseTransitCidr := url.Values{}
	enableAdvertiseTransitCidr.Add("CID", c.CID)
	enableAdvertiseTransitCidr.Add("action", "enable_advertise_transit_cidr")
	enableAdvertiseTransitCidr.Add("vpc_id", vgwConn.VPCId)
	enableAdvertiseTransitCidr.Add("connection_name", vgwConn.ConnName)

	Url.RawQuery = enableAdvertiseTransitCidr.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_advertise_transit_cidr failed: " + err.Error())
	}
	var data VGWConnEnableAdvertiseTransitCidrResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_advertise_transit_cidr failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API enable_advertise_transit_cidr Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableAdvertiseTransitCidr(vgwConn *VGWConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_advertise_transit_cidr") + err.Error())
	}
	disableAdvertiseTransitCidr := url.Values{}
	disableAdvertiseTransitCidr.Add("CID", c.CID)
	disableAdvertiseTransitCidr.Add("action", "disable_advertise_transit_cidr")
	disableAdvertiseTransitCidr.Add("vpc_id", vgwConn.VPCId)
	disableAdvertiseTransitCidr.Add("connection_name", vgwConn.ConnName)

	Url.RawQuery = disableAdvertiseTransitCidr.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_advertise_transit_cidr failed: " + err.Error())
	}
	var data VGWConnEnableAdvertiseTransitCidrResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_advertise_transit_cidr failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_advertise_transit_cidr Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SetBgpManualSpokeAdvertisedNetworks(vgwConn *VGWConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for set_bgp_manual_spoke_advertised_networks") + err.Error())
	}
	setBgpManualSpokeAdvertisedNetworks := url.Values{}
	setBgpManualSpokeAdvertisedNetworks.Add("CID", c.CID)
	setBgpManualSpokeAdvertisedNetworks.Add("action", "set_bgp_manual_spoke_advertised_networks")
	setBgpManualSpokeAdvertisedNetworks.Add("vpc_id", vgwConn.VPCId)
	setBgpManualSpokeAdvertisedNetworks.Add("connection_name", vgwConn.ConnName)
	setBgpManualSpokeAdvertisedNetworks.Add("cidr", vgwConn.BgpManualSpokeAdvertiseCidrs)

	Url.RawQuery = setBgpManualSpokeAdvertisedNetworks.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get set_bgp_manual_spoke_advertised_networks failed: " + err.Error())
	}
	var data VGWConnBgpManualSpokeAdvertisedNetworksResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode set_bgp_manual_spoke_advertised_networks failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API set_bgp_manual_spoke_advertised_networks Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableBgpManualSpokeAdvertisedNetworks(vgwConn *VGWConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_bgp_manual_spoke_advertised_networks") + err.Error())
	}
	disableBgpManualSpokeAdvertisedNetworks := url.Values{}
	disableBgpManualSpokeAdvertisedNetworks.Add("CID", c.CID)
	disableBgpManualSpokeAdvertisedNetworks.Add("action", "disable_bgp_manual_spoke_advertised_networks")
	disableBgpManualSpokeAdvertisedNetworks.Add("vpc_id", vgwConn.VPCId)
	disableBgpManualSpokeAdvertisedNetworks.Add("connection_name", vgwConn.ConnName)

	Url.RawQuery = disableBgpManualSpokeAdvertisedNetworks.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_bgp_manual_spoke_advertised_networks failed: " + err.Error())
	}
	var data VGWConnBgpManualSpokeAdvertisedNetworksResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_bgp_manual_spoke_advertised_networks failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_bgp_manual_spoke_advertised_networks Get failed: " + data.Reason)
	}
	return nil
}
