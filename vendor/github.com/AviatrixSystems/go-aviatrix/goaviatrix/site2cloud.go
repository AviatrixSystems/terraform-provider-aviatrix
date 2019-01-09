package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Site2Cloud simple struct to hold site2cloud details
type Site2Cloud struct {
	Action             string `form:"action,omitempty"`
	CID                string `form:"CID,omitempty"`
	VpcID              string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	TunnelName         string `form:"connection_name" json:"name,omitempty"`
	RemoteGwType       string `form:"remote_gateway_type,omitempty" json:"peer_type,omitempty"`
	ConnType           string `form:"connection_type,omitempty" json:""connection_type,omitempty"`
	TunnelType         string `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	GwName             string `form:"primary_cloud_gateway_name,omitempty" json:"gw_name,omitempty"`
	BackupGwName       string `form:"backup_gateway_name,omitempty"`
	RemoteGwIP         string `form:"remote_gateway_ip,omitempty" json:"peer_ip,omitempty"`
	RemoteGwIP2        string `form:"backup_remote_gateway_ip,omitempty"`
	PreSharedKey       string `form:"pre_shared_key,omitempty"`
	BackupPreSharedKey string `form:"backup_pre_shared_key,omitempty"`
	RemoteSubnet       string `form:"remote_subnet_cidr,omitempty" json:"remote_cidr,omitempty"`
	LocalSubnet        string `form:"local_subnet_cidr,omitempty" json:"local_cidr,omitempty"`
	HAEnabled          string `form:"ha_enabled,omitempty" json:"ha_status,omitempty"`
}

type Site2CloudResp struct {
	Return  bool               `json:"return"`
	Results Site2CloudConnList `json:"results"`
	Reason  string             `json:"reason"`
}

type Site2CloudConnList struct {
	Connections []Site2Cloud `json:"connections"`
}

func (c *Client) CreateSite2Cloud(site2cloud *Site2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "add_site2cloud"
	resp, err := c.Post(c.baseURL, site2cloud)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find s2c connection %s: %s", site2cloud.TunnelName, data.Reason)
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetSite2Cloud(site2cloud *Site2Cloud) (*Site2Cloud, error) {
	site2cloud.Action = "list_site2cloud_conn"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s&connection_name=%s", c.CID, site2cloud.Action,
		site2cloud.TunnelName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data Site2CloudResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}
	for i := 0; i < len(data.Results.Connections); i++ {
		conn := data.Results.Connections[i]
		if site2cloud.VpcID == conn.VpcID {
			return &conn, nil
		}
	}
	return nil, ErrNotFound

}

func (c *Client) UpdateSite2Cloud(site2cloud *Site2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "edit_site2cloud_conn"
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&vpc_id=%s&conn_name=%s&local_subnet_cidr=%s&remote_subnet_cidr=%s",
		c.CID, site2cloud.Action, site2cloud.VpcID, site2cloud.TunnelName, site2cloud.LocalSubnet,
		site2cloud.RemoteSubnet)
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
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DeleteSite2Cloud(site2cloud *Site2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "delete_site2cloud_connection"
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&vpc_id=%s&connection_name=%s", c.CID, site2cloud.Action,
		site2cloud.VpcID, site2cloud.TunnelName)

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
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
