package goaviatrix

import (
	"fmt"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// Site2Cloud simple struct to hold site2cloud details
type Site2Cloud struct {
	Action                  string `form:"action,omitempty"`
	CID                     string `form:"CID,omitempty"`
	VpcID                   string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	ConnName                string `form:"name,omitempty" json:"name,omitempty"`
	RemoteGwType            string `form:"remote_gw_type,omitempty" json:"peer_type,omitempty"`
	TunnelType              string `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	GwName                  string `form:"gw_name"`
	RemoteGwIP              string `form:"peer_ip,omitempty" json:"peer_ip,omitempty"`
	PreSharedKey            string `form:"presk,omitempty"`
	RemoteSubnet            string `form:"remote_cidr,omitempty" json:"remote_cidr,omitempty"`
	LocalSubnet             string `form:"cloud_subnet,omitempty" json:"local_cidr,omitempty"`
}

type Site2CloudResp struct {
	Return  bool   `json:"return"`
	Results Site2CloudConnList `json:"results"`
	Reason  string `json:"reason"`
}

type Site2CloudConnList struct {
	Connections []Site2Cloud `json:"connections"`
}

func (c *Client) CreateSite2Cloud(site2cloud *Site2Cloud) (error) {
	site2cloud.CID=c.CID
	site2cloud.Action="add_site2cloud_conn"
	resp,err := c.Post(c.baseURL, site2cloud)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetSite2Cloud(site2cloud *Site2Cloud) (*Site2Cloud, error) {
	site2cloud.Action="list_site2cloud_conn"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s&conn_name=%s", c.CID, site2cloud.Action, site2cloud.ConnName)
	resp,err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data Site2CloudResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}

	return &data.Results.Connections[0], nil
}

func (c *Client) UpdateSite2Cloud(site2cloud *Site2Cloud) (error) {
	site2cloud.CID = c.CID
	site2cloud.Action = "edit_site2cloud_conn"
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&vpc_id=%s&conn_name=%s&cloud_subnet_cidr=%s&remote_cidr=%s", c.CID, site2cloud.Action, site2cloud.VpcID ,site2cloud.ConnName ,site2cloud.LocalSubnet, site2cloud.RemoteSubnet)
	log.Printf("[TRACE] %s %s Body: %s", verb, c.baseURL, body)
	req, err := http.NewRequest(verb, c.baseURL, strings.NewReader(body))
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DeleteSite2Cloud(site2cloud *Site2Cloud) (error) {
	site2cloud.CID=c.CID
	site2cloud.Action="delete_site2cloud_conn"
	resp,err := c.Post(c.baseURL, site2cloud)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}
