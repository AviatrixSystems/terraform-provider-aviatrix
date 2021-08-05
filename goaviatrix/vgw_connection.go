package goaviatrix

import (
	"strings"
)

// VGWConn simple struct to hold VGW Connection details
type VGWConn struct {
	Action           string `form:"action,omitempty"`
	BgpLocalAsNum    string `form:"bgp_local_asn_num,omitempty" json:"bgp_local_asn_num,omitempty"`
	BgpVGWId         string `form:"vgw_id,omitempty" json:"bgp_vgw_id,omitempty"`
	BgpVGWAccount    string `form:"bgp_vgw_account_name,omitempty" json:"bgp_vgw_account,omitempty"`
	BgpVGWRegion     string `form:"bgp_vgw_region,omitempty" json:"bgp_vgw_region,omitempty"`
	CID              string `form:"CID,omitempty"`
	ConnName         string `form:"connection_name,omitempty" json:"name,omitempty"`
	GwName           string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	VPCId            string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	ManualBGPCidrs   []string
	EventTriggeredHA bool
	PrependAsPath    string
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
	ConnName         []string `json:"name"`
	GwName           string   `json:"gw_name"`
	VPCId            []string `json:"vpc_id"`
	BgpVGWId         string   `json:"bgp_vgw_id"`
	BgpVGWAccount    string   `json:"bgp_vgw_account"`
	BgpVGWRegion     string   `json:"bgp_vgw_region"`
	BgpLocalAsNum    string   `json:"bgp_local_asn_number"`
	ManualBGPCidrs   []string `json:"conn_bgp_manual_advertise_cidrs"`
	EventTriggeredHA string   `json:"event_triggered_ha"`
	PrependAsPath    string   `json:"conn_bgp_prepend_as_path"`
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
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "connect_transit_gw_to_vgw",
		"vpc_id":               vgwConn.VPCId,
		"connection_name":      vgwConn.ConnName,
		"transit_gw":           vgwConn.GwName,
		"vgw_id":               vgwConn.BgpVGWId,
		"bgp_vgw_account_name": vgwConn.BgpVGWAccount,
		"bgp_vgw_region":       vgwConn.BgpVGWRegion,
		"bgp_local_as_number":  vgwConn.BgpLocalAsNum,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetVGWConn(vgwConn *VGWConn) (*VGWConn, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_vgw_connections",
	}

	data := VGWConnListResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
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
	form := map[string]string{
		"CID":             c.CID,
		"action":          "disconnect_transit_gw_from_vgw",
		"vpc_id":          vgwConn.VPCId,
		"connection_name": vgwConn.ConnName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
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
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
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
		vgwConn.EventTriggeredHA = data.Results.Connections.EventTriggeredHA == "enabled"
		vgwConn.PrependAsPath = data.Results.Connections.PrependAsPath
		return vgwConn, nil
	}
	return nil, ErrNotFound
}

func (c *Client) EditVgwConnectionASPathPrepend(vgwConn *VGWConn, prependASPath []string) error {
	action := "edit_transit_connection_as_path_prepend"
	return c.PostAPI(action, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		GatewayName    string `form:"gateway_name"`
		ConnectionName string `form:"connection_name"`
		PrependASPath  string `form:"connection_as_path_prepend"`
	}{
		CID:            c.CID,
		Action:         action,
		GatewayName:    vgwConn.GwName,
		ConnectionName: vgwConn.ConnName,
		PrependASPath:  strings.Join(prependASPath, ","),
	}, BasicCheck)
}
