package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Site2Cloud simple struct to hold site2cloud details
type Site2Cloud struct {
	Action              string `form:"action,omitempty"`
	CID                 string `form:"CID,omitempty"`
	VpcID               string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	TunnelName          string `form:"connection_name" json:"name,omitempty"`
	RemoteGwType        string `form:"remote_gateway_type,omitempty"`
	ConnType            string `form:"connection_type,omitempty" json:"type,omitempty"`
	TunnelType          string `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	GwName              string `form:"primary_cloud_gateway_name,omitempty" json:"gw_name,omitempty"`
	BackupGwName        string `form:"backup_gateway_name,omitempty"`
	RemoteGwIP          string `form:"remote_gateway_ip,omitempty" json:"peer_ip,omitempty"`
	RemoteGwIP2         string `form:"backup_remote_gateway_ip,omitempty"`
	PreSharedKey        string `form:"pre_shared_key,omitempty"`
	BackupPreSharedKey  string `form:"backup_pre_shared_key,omitempty"`
	RemoteSubnet        string `form:"remote_subnet_cidr,omitempty" json:"remote_cidr,omitempty"`
	LocalSubnet         string `form:"local_subnet_cidr,omitempty" json:"local_cidr,omitempty"`
	HAEnabled           string `form:"ha_enabled,omitempty" json:"ha_status,omitempty"`
	PeerType            string `form:"peer_type,omitempty"`
	SslServerPool       string `form:"ssl_server_pool,omitempty"`
	NetworkType         string `form:"network_type,omitempty"`
	CloudSubnetCidr     string `form:"cloud_subnet_cidr,omitempty"`
	RemoteCidr          string `form:"remote_cidr,omitempty"`
	RemoteSubnetVirtual string `form:"virtual_remote_subnet_cidr,omitempty" json:"virtual_remote_subnet_cidr,omitempty"`
	LocalSubnetVirtual  string `form:"virtual_local_subnet_cidr,omitempty" json:"virtual_local_subnet_cidr,omitempty"`
}

type EditSite2Cloud struct {
	Action          string `form:"action,omitempty"`
	CID             string `form:"CID,omitempty"`
	VpcID           string `form:"vpc_id,omitempty"`
	ConnName        string `form:"conn_name"`
	GwName          string `form:"primary_cloud_gateway_name,omitempty"`
	NetworkType     string `form:"network_type,omitempty"`
	CloudSubnetCidr string `form:"cloud_subnet_cidr,omitempty"`
}

type Site2CloudResp struct {
	Return  bool               `json:"return"`
	Results Site2CloudConnList `json:"results"`
	Reason  string             `json:"reason"`
}

type Site2CloudConnList struct {
	Connections []Site2Cloud `json:"connections"`
}

type EditSite2CloudConnDetail struct {
	VpcID        []string `json:"vpc_id,omitempty"`
	TunnelName   []string `json:"name,omitempty"`
	ConnType     string   `json:"type,omitempty"`
	TunnelType   string   `json:"tunnel_type,omitempty"`
	GwName       []string `json:"gw_name,omitempty"`
	BackupGwName []string `json:"backup_gateway_name,omitempty"`
	RemoteGwIP   []string `json:"remote_gateway_ip,omitempty"`
	RemoteGwIP2  []string `json:"backup_remote_gateway_ip,omitempty"`
	//PreSharedKey        string `json:"pre_shared_key,omitempty"`
	//BackupPreSharedKey  string `json:"backup_pre_shared_key,omitempty"`
	RemoteSubnet string `json:"real_remote_cidr,omitempty"`
	LocalSubnet  string `json:"real_local_cidr,omitempty"`
	HAEnabled    string `json:"ha_status,omitempty"`
	PeerType     string `json:"peer_type,omitempty"`
	//SslServerPool       string `json:"ssl_server_pool,omitempty"`
	//NetworkType         string `json:"network_type,omitempty"`
	//CloudSubnetCidr     string `json:"cloud_subnet_cidr,omitempty"`
	RemoteSubnetVirtual string `json:"virt_remote_cidr,omitempty"`
	LocalSubnetVirtual  string `json:"virt_local_cidr,omitempty"`
}

type Site2CloudConnDetailResp struct {
	Return  bool                     `json:"return"`
	Results Site2CloudConnDetailList `json:"results"`
	Reason  string                   `json:"reason"`
}

type Site2CloudConnDetailList struct {
	Connections EditSite2CloudConnDetail `json:"connections"`
}

func (c *Client) CreateSite2Cloud(site2cloud *Site2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "add_site2cloud"
	resp, err := c.Post(c.baseURL, site2cloud)
	if err != nil {
		return errors.New("HTTP Post add_site2cloud failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_site2cloud failed: " + err.Error())
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find s2c connection %s: %s", site2cloud.TunnelName, data.Reason)
		return errors.New("Rest API add_site2cloud Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSite2Cloud(site2cloud *Site2Cloud) (*Site2Cloud, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_site2cloud_conn") + err.Error())
	}
	listSite2CloudConn := url.Values{}
	listSite2CloudConn.Add("CID", c.CID)
	listSite2CloudConn.Add("action", "list_site2cloud_conn")
	listSite2CloudConn.Add("connection_name", site2cloud.TunnelName)
	Url.RawQuery = listSite2CloudConn.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_site2cloud_conn failed: " + err.Error())
	}
	var data Site2CloudResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_site2cloud_conn failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_site2cloud_conn Get failed: " + data.Reason)
	}
	for i := 0; i < len(data.Results.Connections); i++ {
		conn := data.Results.Connections[i]
		if site2cloud.VpcID == conn.VpcID && site2cloud.TunnelName == conn.TunnelName {
			return &conn, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) GetSite2CloudConnDetail(site2cloud *Site2Cloud) (*Site2Cloud, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for get_site2cloud_conn_detail") + err.Error())
	}
	getSite2CloudConnDetail := url.Values{}
	getSite2CloudConnDetail.Add("CID", c.CID)
	getSite2CloudConnDetail.Add("action", "get_site2cloud_conn_detail")
	getSite2CloudConnDetail.Add("conn_name", site2cloud.TunnelName)
	getSite2CloudConnDetail.Add("vpc_id", site2cloud.VpcID)
	Url.RawQuery = getSite2CloudConnDetail.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get get_site2cloud_conn_detail failed: " + err.Error())
	}
	var data Site2CloudConnDetailResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_site2cloud_conn_detail failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API get_site2cloud_conn_detail Get failed: " + data.Reason)
	}
	if len(data.Results.Connections.TunnelName) != 0 {
		site2cloud.GwName = data.Results.Connections.GwName[0]
		site2cloud.ConnType = data.Results.Connections.ConnType
		site2cloud.TunnelType = data.Results.Connections.TunnelType
		site2cloud.RemoteGwType = data.Results.Connections.PeerType
		site2cloud.RemoteGwIP = data.Results.Connections.RemoteGwIP[0]
		site2cloud.RemoteSubnet = data.Results.Connections.RemoteSubnet
		site2cloud.LocalSubnet = data.Results.Connections.LocalSubnet
		site2cloud.HAEnabled = data.Results.Connections.HAEnabled
		site2cloud.RemoteSubnetVirtual = data.Results.Connections.RemoteSubnetVirtual
		site2cloud.LocalSubnetVirtual = data.Results.Connections.LocalSubnetVirtual
		site2cloud.BackupGwName = data.Results.Connections.BackupGwName[0]
		site2cloud.RemoteGwIP2 = data.Results.Connections.RemoteGwIP2[0]
		return site2cloud, nil
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateSite2Cloud(site2cloud *EditSite2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "edit_site2cloud_conn"
	resp, err := c.Post(c.baseURL, site2cloud)
	if err != nil {
		return errors.New("HTTP Post edit_site2cloud_conn failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode edit_site2cloud_conn failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API edit_site2cloud_conn Post failed: " + data.Reason)
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
		return errors.New("HTTP Post NewRequest delete_site2cloud_connection failed: " + err.Error())
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.New("HTTP Post delete_site2cloud_connection failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_site2cloud_connection failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_site2cloud_connection Post failed: " + data.Reason)
	}
	return nil
}
