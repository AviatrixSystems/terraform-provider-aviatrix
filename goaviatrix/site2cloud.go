package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const Phase1AuthDefault = "SHA-1"
const Phase1DhGroupDefault = "2"
const Phase1EncryptionDefault = "AES-256-CBC"
const Phase2AuthDefault = "HMAC-SHA-1"
const Phase2DhGroupDefault = "2"
const Phase2EncryptionDefault = "AES-256-CBC"
const SslServerPoolDefault = "192.168.44.0/24"

// Site2Cloud simple struct to hold site2cloud details
type Site2Cloud struct {
	Action                  string   `form:"action,omitempty"`
	CID                     string   `form:"CID,omitempty"`
	VpcID                   string   `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	TunnelName              string   `form:"connection_name" json:"name,omitempty"`
	RemoteGwType            string   `form:"remote_gateway_type,omitempty"`
	ConnType                string   `form:"connection_type,omitempty" json:"type,omitempty"`
	TunnelType              string   `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	GwName                  string   `form:"primary_cloud_gateway_name,omitempty" json:"gw_name,omitempty"`
	BackupGwName            string   `form:"backup_gateway_name,omitempty"`
	RemoteGwIP              string   `form:"remote_gateway_ip,omitempty" json:"peer_ip,omitempty"`
	RemoteGwIP2             string   `form:"backup_remote_gateway_ip,omitempty"`
	PreSharedKey            string   `form:"pre_shared_key,omitempty"`
	BackupPreSharedKey      string   `form:"backup_pre_shared_key,omitempty"`
	RemoteSubnet            string   `form:"remote_subnet_cidr,omitempty" json:"remote_cidr,omitempty"`
	LocalSubnet             string   `form:"local_subnet_cidr,omitempty" json:"local_cidr,omitempty"`
	HAEnabled               string   `form:"ha_enabled,omitempty" json:"ha_status,omitempty"`
	PeerType                string   `form:"peer_type,omitempty"`
	SslServerPool           string   `form:"ssl_server_pool,omitempty"`
	NetworkType             string   `form:"network_type,omitempty"`
	CloudSubnetCidr         string   `form:"cloud_subnet_cidr,omitempty"`
	RemoteCidr              string   `form:"remote_cidr,omitempty"`
	RemoteSubnetVirtual     string   `form:"virtual_remote_subnet_cidr,omitempty" json:"virtual_remote_subnet_cidr,omitempty"`
	LocalSubnetVirtual      string   `form:"virtual_local_subnet_cidr,omitempty" json:"virtual_local_subnet_cidr,omitempty"`
	Phase1Auth              string   `form:"phase1_auth,omitempty"`
	Phase1DhGroups          string   `form:"phase1_dh_group,omitempty"`
	Phase1Encryption        string   `form:"phase1_encryption,omitempty"`
	Phase2Auth              string   `form:"phase2_auth,omitempty"`
	Phase2DhGroups          string   `form:"phase2_dh_group,omitempty"`
	Phase2Encryption        string   `form:"phase2_encryption,omitempty"`
	PrivateRouteEncryption  string   `form:"private_route_encryption,omitempty"`
	RemoteGwLatitude        float64  `form:"remote_gateway_latitude,omitempty"`
	RemoteGwLongitude       float64  `form:"remote_gateway_longitude,omitempty"`
	BackupRemoteGwLatitude  float64  `form:"backup_remote_gateway_latitude,omitempty"`
	BackupRemoteGwLongitude float64  `form:"backup_remote_gateway_longitude,omitempty"`
	RouteTableList          []string `form:"route_table_list,omitempty"`
	CustomAlgorithms        bool
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
	RouteTableList      []string      `json:"rtbls,omitempty"`
	SslServerPool       []string      `json:"ssl_server_pool,omitempty"`

	//PreSharedKey        string `json:"pre_shared_key,omitempty"`
	//BackupPreSharedKey  string `json:"backup_pre_shared_key,omitempty"`
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
	Phase1Auth      []string `json:"ph1_auth,omitempty"`
	Phase1DhGroups  []string `json:"ph1_dh,omitempty"`
	Phase1Encrption []string `json:"ph1_encr,omitempty"`
	Phase2Auth      []string `json:"ph2_auth,omitempty"`
	Phase2DhGroups  []string `json:"ph2_dh,omitempty"`
	Phase2Encrption []string `json:"ph2_encr,omitempty"`
}

//func (c *Client) CreateSite2Cloud(site2cloud *Site2Cloud) error {
//	site2cloud.CID = c.CID
//	site2cloud.Action = "add_site2cloud"
//	resp, err := c.Post(c.baseURL, site2cloud)
//	if err != nil {
//		return errors.New("HTTP Post add_site2cloud failed: " + err.Error())
//	}
//	var data APIResp
//	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
//		return errors.New("Json Decode add_site2cloud failed: " + err.Error())
//	}
//	if !data.Return {
//		log.Printf("[INFO] Couldn't find s2c connection %s: %s", site2cloud.TunnelName, data.Reason)
//		return errors.New("Rest API add_site2cloud Post failed: " + data.Reason)
//	}
//	return nil
//}

func (c *Client) CreateSite2Cloud(site2cloud *Site2Cloud) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for add_site2cloud " + err.Error())
	}
	addSite2cloud := url.Values{}
	addSite2cloud.Add("CID", c.CID)
	addSite2cloud.Add("action", "add_site2cloud")
	addSite2cloud.Add("vpc_id", site2cloud.VpcID)
	addSite2cloud.Add("connection_name", site2cloud.TunnelName)
	addSite2cloud.Add("connection_type", site2cloud.ConnType)
	addSite2cloud.Add("remote_gateway_type", site2cloud.RemoteGwType)
	addSite2cloud.Add("tunnel_type", site2cloud.TunnelType)

	addSite2cloud.Add("ha_enabled", site2cloud.HAEnabled)
	addSite2cloud.Add("backup_gateway_name", site2cloud.BackupGwName)
	addSite2cloud.Add("backup_remote_gateway_ip", site2cloud.RemoteGwIP2)

	addSite2cloud.Add("phase1_auth", site2cloud.Phase1Auth)
	addSite2cloud.Add("phase1_dh_group", site2cloud.Phase1DhGroups)
	addSite2cloud.Add("phase1_encryption", site2cloud.Phase1Encryption)
	addSite2cloud.Add("phase2_auth", site2cloud.Phase2Auth)
	addSite2cloud.Add("phase2_dh_group", site2cloud.Phase2DhGroups)
	addSite2cloud.Add("phase2_encryption", site2cloud.Phase2Encryption)

	if site2cloud.TunnelType == "tcp" {
		addSite2cloud.Add("ssl_server_pool", site2cloud.SslServerPool)
	}

	if site2cloud.PrivateRouteEncryption == "true" {
		addSite2cloud.Add("private_route_encryption", site2cloud.PrivateRouteEncryption)
		if len(site2cloud.RouteTableList) != 0 {
			for i := range site2cloud.RouteTableList {
				addSite2cloud.Add("route_table_list["+strconv.Itoa(i)+"]", site2cloud.RouteTableList[i])
			}
		}
		latitude := fmt.Sprintf("%f", site2cloud.RemoteGwLatitude)
		longitude := fmt.Sprintf("%f", site2cloud.RemoteGwLongitude)
		addSite2cloud.Add("remote_gateway_latitude", latitude)
		addSite2cloud.Add("remote_gateway_longitude", longitude)
		if site2cloud.HAEnabled == "yes" {
			backupLatitude := fmt.Sprintf("%f", site2cloud.BackupRemoteGwLatitude)
			backupLongitude := fmt.Sprintf("%f", site2cloud.BackupRemoteGwLongitude)
			addSite2cloud.Add("remote_gateway_latitude", backupLatitude)
			addSite2cloud.Add("remote_gateway_longitude", backupLongitude)
		}
	}

	addSite2cloud.Add("primary_cloud_gateway_name", site2cloud.GwName)
	addSite2cloud.Add("remote_gateway_ip", site2cloud.RemoteGwIP)
	addSite2cloud.Add("remote_subnet_cidr", site2cloud.RemoteSubnet)
	addSite2cloud.Add("local_subnet_cidr", site2cloud.LocalSubnet)
	addSite2cloud.Add("virtual_remote_subnet_cidr", site2cloud.RemoteSubnetVirtual)
	addSite2cloud.Add("virtual_local_subnet_cidr", site2cloud.LocalSubnetVirtual)

	addSite2cloud.Add("pre_shared_key", site2cloud.PreSharedKey)
	addSite2cloud.Add("backup_pre_shared_key", site2cloud.BackupPreSharedKey)

	Url.RawQuery = addSite2cloud.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get add_site2cloud failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_site2cloud failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_site2cloud Get failed: " + data.Reason)
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
		if s2cConnDetail.Algorithm.Phase1Auth[0] == Phase1AuthDefault &&
			s2cConnDetail.Algorithm.Phase2Auth[0] == Phase2AuthDefault &&
			s2cConnDetail.Algorithm.Phase1DhGroups[0] == Phase1DhGroupDefault &&
			s2cConnDetail.Algorithm.Phase2DhGroups[0] == Phase2DhGroupDefault &&
			s2cConnDetail.Algorithm.Phase1Encrption[0] == Phase1EncryptionDefault &&
			s2cConnDetail.Algorithm.Phase2Encrption[0] == Phase2EncryptionDefault {
			site2cloud.CustomAlgorithms = false
			site2cloud.Phase1Auth = ""
			site2cloud.Phase2Auth = ""
			site2cloud.Phase1DhGroups = ""
			site2cloud.Phase2DhGroups = ""
			site2cloud.Phase1Encryption = ""
			site2cloud.Phase2Encryption = ""
		} else {
			site2cloud.CustomAlgorithms = true
			site2cloud.Phase1Auth = s2cConnDetail.Algorithm.Phase1Auth[0]
			site2cloud.Phase2Auth = s2cConnDetail.Algorithm.Phase2Auth[0]
			site2cloud.Phase1DhGroups = s2cConnDetail.Algorithm.Phase1DhGroups[0]
			site2cloud.Phase2DhGroups = s2cConnDetail.Algorithm.Phase2DhGroups[0]
			site2cloud.Phase1Encryption = s2cConnDetail.Algorithm.Phase1Encrption[0]
			site2cloud.Phase2Encryption = s2cConnDetail.Algorithm.Phase2Encrption[0]
		}
		if len(s2cConnDetail.RouteTableList) > 0 {
			site2cloud.RouteTableList = s2cConnDetail.RouteTableList
			site2cloud.PrivateRouteEncryption = "true"
		} else {
			site2cloud.PrivateRouteEncryption = "false"
		}
		if s2cConnDetail.SslServerPool[0] != "192.168.44.0/24" {
			site2cloud.SslServerPool = s2cConnDetail.SslServerPool[0]
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
	Phase1EncryptionList := []string{"AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "3DES"}
	Phase2AuthList := []string{"HMAC-SHA-1", "HMAC-SHA-256", "HMAC-SHA-384", "HMAC-SHA-512", "NO-AUTH"}
	Phase2DhGroupsList := []string{"1", "2", "5", "14", "15", "16", "17", "18"}
	Phase2EncryptionList := []string{"AES-128-CBC", "AES-128-GCM-64", "AES-128-GCM-96", "AES-128-GCM-128",
		"AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "3DES", "NULL-ENCR"}

	if !Contains(Phase1AuthList, site2cloud.Phase1Auth) {
		return errors.New("invalid value for phase_1_authentication")
	}
	if !Contains(Phase1DhGroupsList, site2cloud.Phase1DhGroups) {
		return errors.New("invalid value for phase_1_dh_groups")
	}
	if !Contains(Phase1EncryptionList, site2cloud.Phase1Encryption) {
		return errors.New("invalid value for phase_1_encryption")
	}
	if !Contains(Phase2AuthList, site2cloud.Phase2Auth) {
		return errors.New("invalid value for phase_2_authentication")
	}
	if !Contains(Phase2DhGroupsList, site2cloud.Phase2DhGroups) {
		return errors.New("invalid value for phase_2_dh_groups")
	}
	if !Contains(Phase2EncryptionList, site2cloud.Phase2Encryption) {
		return errors.New("invalid value for phase_2_encryption")
	}
	return nil
}
