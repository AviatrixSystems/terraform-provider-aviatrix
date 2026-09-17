package goaviatrix

type VpnUserXlr struct {
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	Endpoints      string `form:"endpoints,omitempty"`
	AllEndpoints   string `json:"all,omitempty"`
	FreeEndpoints  string `json:"free,omitempty"`
	InUseEndpoints string `json:"inuse,omitempty"`
}

type VpnUserXlrAPIResp struct {
	Return  bool                `json:"return"`
	Results map[string][]string `json:"results"`
	Reason  string              `json:"reason"`
}

func (c *Client) GetVpnUserAccelerator() ([]string, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_vpn_user_xlr",
	}

	var data VpnUserXlrAPIResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return data.Results["inuse"], nil
}

func (c *Client) UpdateVpnUserAccelerator(xlr *VpnUserXlr) error {
	xlr.CID = c.CID
	xlr.Action = "update_vpn_user_xlr"

	return c.PostAPI(xlr.Action, xlr, BasicCheck)
}
