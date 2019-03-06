package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
)

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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for modify_split_tunnel") + err.Error())
	}
	modifySplitTunnel := url.Values{}
	modifySplitTunnel.Add("CID", c.CID)
	modifySplitTunnel.Add("action", "modify_split_tunnel")
	modifySplitTunnel.Add("command", "get")
	modifySplitTunnel.Add("vpc_id", splitTunnel.VpcID)
	modifySplitTunnel.Add("lb_name", splitTunnel.ElbName)
	Url.RawQuery = modifySplitTunnel.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get modify_split_tunnel(get) failed: " + err.Error())
	}
	var data SplitTunnelResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode modify_split_tunnel(get) failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API modify_split_tunnel(get) Get failed: " + data.Reason)
	}
	return &data.Results, nil
}

func (c *Client) ModifySplitTunnel(splitTunnel *SplitTunnel) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for modify_split_tunnel") + err.Error())
	}
	modifySplitTunnel := url.Values{}
	modifySplitTunnel.Add("CID", c.CID)
	modifySplitTunnel.Add("action", "modify_split_tunnel")
	modifySplitTunnel.Add("command", "modify")
	modifySplitTunnel.Add("vpc_id", splitTunnel.VpcID)
	modifySplitTunnel.Add("lb_name", splitTunnel.ElbName)
	modifySplitTunnel.Add("split_tunnel", splitTunnel.SplitTunnel)
	modifySplitTunnel.Add("additional_cidrs", splitTunnel.AdditionalCidrs)
	modifySplitTunnel.Add("nameservers", splitTunnel.NameServers)
	modifySplitTunnel.Add("search_domains", splitTunnel.SearchDomains)
	Url.RawQuery = modifySplitTunnel.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get modify_split_tunnel(modify) failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode modify_split_tunnel(modify) failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API modify_split_tunnel(modify) Get failed: " + data.Reason)
	}
	return nil
}
