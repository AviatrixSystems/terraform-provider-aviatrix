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
	Phase1Auth          string `form:"phase1_auth,omitempty"`
	Phase1DhGroups      string `form:"phase1_dh_group,omitempty"`
	Phase1Encryption    string `form:"phase1_encryption,omitempty"`
	Phase2Auth          string `form:"phase2_auth,omitempty"`
	Phase2DhGroups      string `form:"phase2_dh_group,omitempty"`
	Phase2Encryption    string `form:"phase2_encryption,omitempty"`
	CustomAlgorithms    bool
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
	VpcID               []string      `json:"vpc_id,omitempty"`
	TunnelName          []string      `json:"name,omitempty"`
	ConnType            string        `json:"type,omitempty"`
	TunnelType          []string      `json:"tunnel_type,omitempty"`
	GwName              []string      `json:"gw_name,omitempty"`
	BackupGwName        []string      `json:"backup_gateway_name,omitempty"`
	RemoteGwIP          []string      `json:"remote_gateway_ip,omitempty"`
	RemoteGwIP2         []string      `json:"backup_remote_gateway_ip,omitempty"`
	Tunnels             []TunnelInfo  `json:"tunnels,omitempty"`
	RemoteSubnet        string        `json:"real_remote_cidr,omitempty"`
	LocalSubnet         string        `json:"real_local_cidr,omitempty"`
	RemoteCidr          string        `json:"remote_cidr,omitempty"`
	LocalCidr           string        `json:"local_cidr,omitempty"`
	HAEnabled           string        `json:"ha_status,omitempty"`
	PeerType            string        `json:"peer_type,omitempty"`
	RemoteSubnetVirtual string        `json:"virt_remote_cidr,omitempty"`
	LocalSubnetVirtual  string        `json:"virt_local_cidr,omitempty"`
	Algorithm           AlgorithmInfo `json:"algorithm,omitempty"`
	//PreSharedKey        string `json:"pre_shared_key,omitempty"`
	//BackupPreSharedKey  string `json:"backup_pre_shared_key,omitempty"`
	//SslServerPool       string `json:"ssl_server_pool,omitempty"`
	//NetworkType         string `json:"network_type,omitempty"`
	//CloudSubnetCidr     string `json:"cloud_subnet_cidr,omitempty"`
}

type Site2CloudConnDetailResp struct {
	Return  bool                     `json:"return"`
	Results Site2CloudConnDetailList `json:"results"`
	Reason  string                   `json:"reason"`
}

type Site2CloudConnDetailList struct {
	Connections EditSite2CloudConnDetail `json:"connections"`
}

type TunnelInfo struct {
	Status       string `json:"status"`
	IPAddr       string `json:"ip_addr"`
	Name         string `json:"name"`
	PeerIP       string `json:"peer_ip"`
	GwName       string `json:"gw_name"`
	TunnelStatus string `json:"tunnel_status"`
}

type AlgorithmInfo struct {
	Phase1Auth      []string `json:"ph1_auth"`
	Phase1DhGroups  []string `json:"ph1_dh"`
	Phase1Encrption []string `json:"ph1_encr"`
	Phase2Auth      []string `json:"ph2_auth"`
	Phase2DhGroups  []string `json:"ph2_dh"`
	Phase2Encrption []string `json:"ph2_encr"`
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

	s2cConnDetail := data.Results.Connections
	if len(s2cConnDetail.TunnelName) != 0 {
		site2cloud.GwName = s2cConnDetail.GwName[0]
		site2cloud.ConnType = s2cConnDetail.ConnType
		site2cloud.TunnelType = s2cConnDetail.TunnelType[0]
		site2cloud.RemoteGwType = s2cConnDetail.PeerType
		if site2cloud.ConnType == "mapped" {
			site2cloud.RemoteSubnet = s2cConnDetail.RemoteSubnet
			site2cloud.LocalSubnet = s2cConnDetail.LocalSubnet
			site2cloud.RemoteSubnetVirtual = s2cConnDetail.RemoteSubnetVirtual
			site2cloud.LocalSubnetVirtual = s2cConnDetail.LocalSubnetVirtual
		} else {
			site2cloud.RemoteSubnet = s2cConnDetail.RemoteCidr
			site2cloud.LocalSubnet = s2cConnDetail.LocalCidr
		}
		site2cloud.HAEnabled = s2cConnDetail.HAEnabled
		for i := range s2cConnDetail.Tunnels {
			if s2cConnDetail.Tunnels[i].GwName == site2cloud.GwName {
				site2cloud.RemoteGwIP = s2cConnDetail.Tunnels[i].PeerIP
			} else if s2cConnDetail.Tunnels[i].GwName == site2cloud.GwName+"-hagw" {
				site2cloud.BackupGwName = s2cConnDetail.Tunnels[i].GwName
				site2cloud.RemoteGwIP2 = s2cConnDetail.Tunnels[i].PeerIP
			}
		}

		site2cloud.Phase1Auth = s2cConnDetail.Algorithm.Phase1Auth[0]
		site2cloud.Phase1DhGroups = s2cConnDetail.Algorithm.Phase1DhGroups[0]
		site2cloud.Phase1Encryption = s2cConnDetail.Algorithm.Phase1Encrption[0]
		site2cloud.Phase2Auth = s2cConnDetail.Algorithm.Phase2Auth[0]
		site2cloud.Phase2DhGroups = s2cConnDetail.Algorithm.Phase2DhGroups[0]
		site2cloud.Phase2Encryption = s2cConnDetail.Algorithm.Phase2Encrption[0]
		if s2cConnDetail.Algorithm.Phase1Auth[0] != "" {
			site2cloud.CustomAlgorithms = true
		} else {
			site2cloud.CustomAlgorithms = false
		}
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

func (c *Client) Site2CloudAlgorithmCheck(site2cloud *Site2Cloud) error {
	Phase1AuthList := []string{"SHA-1", "SHA-256", "SHA-384", "SHA-512"}
	Phase1DhGroupsList := []string{"1", "2", "5", "14", "15", "16", "17", "18"}
	Phase1EncrptionList := []string{"AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "3DES"}
	Phase2AuthList := []string{"HMAC-SHA-1", "HMAC-SHA-256", "HMAC-SHA-384", "HMAC-SHA-512", "NO-AUTH"}
	Phase2DhGroupsList := []string{"1", "2", "5", "14", "15", "16", "17", "18"}
	Phase2EncrptionList := []string{"AES-128-CBC", "AES-128-GCM-64", "AES-128-GCM-96", "AES-128-GCM-128", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "3DES", "NULL-ENCR"}

	if !Contains(Phase1AuthList, site2cloud.Phase1Auth) {
		return errors.New("invalid value for phase_1_authentication")
	}
	if !Contains(Phase1DhGroupsList, site2cloud.Phase1DhGroups) {
		return errors.New("invalid value for phase_1_dh_groups")
	}
	if !Contains(Phase1EncrptionList, site2cloud.Phase1Encryption) {
		return errors.New("invalid value for phase_1_encryption")
	}
	if !Contains(Phase2AuthList, site2cloud.Phase2Auth) {
		return errors.New("invalid value for phase_2_authentication")
	}
	if !Contains(Phase2DhGroupsList, site2cloud.Phase2DhGroups) {
		return errors.New("invalid value for phase_2_dh_groups")
	}
	if !Contains(Phase2EncrptionList, site2cloud.Phase2Encryption) {
		return errors.New("invalid value for phase_2_encryption")
	}
	return nil
}
