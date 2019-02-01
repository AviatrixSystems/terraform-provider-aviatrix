package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
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
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=modify_split_tunnel&command=get&vpc_id=%s&lb_name=%s",
		c.CID, splitTunnel.VpcID, splitTunnel.ElbName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data SplitTunnelResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}
	return &data.Results, nil
}

func (c *Client) ModifySplitTunnel(splitTunnel *SplitTunnel) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=modify_split_tunnel&command=modify&vpc_id=%s&lb_name=%s"+
		"&split_tunnel=%s&additional_cidrs=%s&nameservers=%s&search_domains=%s", c.CID, splitTunnel.VpcID, splitTunnel.ElbName,
		splitTunnel.SplitTunnel, splitTunnel.AdditionalCidrs, splitTunnel.NameServers, splitTunnel.SearchDomains)
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
