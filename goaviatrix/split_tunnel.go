package goaviatrix

type SplitTunnel struct {
	Action          string `form:"action,omitempty"`
	CID             string `form:"CID,omitempty"`
	Command         string `form:"command,omitempty"`
	VpcID           string `form:"vpc_id,omitempty"`
	ElbName         string `form:"lb_name,omitempty"`
	SplitTunnel     string `form:"split_tunnel,omitempty"`
	AdditionalCidrs string `form:"additional_cidrs,omitempty"`
	NameServers     string `form:"nameservers,omitempty"`
	SearchDomains   string `form:"search_domains,omitempty"`
	SaveTemplate    string `form:"save_template,omitempty"`
	Dns             string `form:"dns,omitempty"`
}

type SplitTunnelResp struct {
	Return  bool            `json:"return"`
	Results SplitTunnelUnit `json:"results"`
	Reason  string          `json:"reason"`
}

type SplitTunnelUnit struct {
	NameServers     string `json:"name_servers"`
	SplitTunnel     string `json:"split_tunnel"`
	SearchDomains   string `json:"search_domains"`
	AdditionalCidrs string `json:"additional_cidrs"`
}

func (c *Client) GetSplitTunnel(splitTunnel *SplitTunnel) (*SplitTunnelUnit, error) {
	form := map[string]string{
		"CID":     c.CID,
		"action":  "modify_split_tunnel",
		"command": "get",
		"vpc_id":  splitTunnel.VpcID,
		"lb_name": splitTunnel.ElbName,
	}

	var data SplitTunnelResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) ModifySplitTunnel(splitTunnel *SplitTunnel) error {
	form := map[string]string{
		"CID":              c.CID,
		"action":           "modify_split_tunnel",
		"command":          "modify",
		"vpc_id":           splitTunnel.VpcID,
		"lb_name":          splitTunnel.ElbName,
		"split_tunnel":     splitTunnel.SplitTunnel,
		"additional_cidrs": splitTunnel.AdditionalCidrs,
		"nameservers":      splitTunnel.NameServers,
		"search_domains":   splitTunnel.SearchDomains,
	}

	if splitTunnel.Dns == "true" {
		form["dns"] = "true"
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
