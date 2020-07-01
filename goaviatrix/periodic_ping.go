package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
)

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
	resp, err := c.Post(c.baseURL, pp)
	if err != nil {
		return errors.New("HTTP Post enable_gateway_periodic_ping failed: " + err.Error())
	}

	var data APIResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body enable_gateway_periodic_ping failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode enable_gateway_periodic_ping failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API enable_gateway_periodic_ping Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetPeriodicPing(pp *PeriodicPing) (*PeriodicPing, error) {
	pp.Action = "get_gateway_periodic_ping_status"
	pp.CID = c.CID
	resp, err := c.Post(c.baseURL, pp)
	if err != nil {
		return nil, errors.New("HTTP POST get_gateway_periodic_ping_status failed: " + err.Error())
	}

	var data PeriodicPingStatusResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_gateway_periodic_ping_status failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_gateway_periodic_ping_status failed: " + err.Error() +
			"\n Body: " + b.String())
	}

	if !data.Return {
		return nil, errors.New("Rest API get_gateway_periodic_ping_status Post failed: " + data.Reason)
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
	resp, err := c.Post(c.baseURL, pp)
	if err != nil {
		return errors.New("HTTP POST disable_gateway_periodic_ping failed: " + err.Error())
	}

	var data APIResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body disable_gateway_periodic_ping failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode disable_gateway_periodic_ping failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API disable_gateway_periodic_ping Post failed: " + data.Reason)
	}

	return nil
}
