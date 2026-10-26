package goaviatrix

type PeriodicPing struct {
	Action        string `form:"action"`
	CID           string `form:"CID"`
	GwName        string `form:"gateway_name"`
	Interval      string `form:"interval"`
	IntervalAsInt int
	IP            string `form:"ip_address"`
}

type PeriodicPingStatusResp struct {
	Return bool                     `json:"return"`
	Reason string                   `json:"reason"`
	Result PeriodicPingStatusResult `json:"results"`
}

type PeriodicPingStatusResult struct {
	Status   string   `json:"status"`
	IPs      []string `json:"address,omitempty"`
	Interval int      `json:"interval,omitempty"`
}

func (c *Client) CreatePeriodicPing(pp *PeriodicPing) error {
	pp.Action = "enable_gateway_periodic_ping"
	pp.CID = c.CID

	return c.PostAPI(pp.Action, pp, BasicCheck)
}

func (c *Client) GetPeriodicPing(pp *PeriodicPing) (*PeriodicPing, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "get_gateway_periodic_ping_status",
		"gateway_name": pp.GwName,
	}

	var data PeriodicPingStatusResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if data.Result.Status != "enabled" {
		return nil, ErrNotFound
	}

	return &PeriodicPing{
		GwName:        pp.GwName,
		IntervalAsInt: data.Result.Interval,
		IP:            data.Result.IPs[0],
	}, nil
}

func (c *Client) DeletePeriodicPing(pp *PeriodicPing) error {
	pp.Action = "disable_gateway_periodic_ping"
	pp.CID = c.CID

	return c.PostAPI(pp.Action, pp, BasicCheck)
}
