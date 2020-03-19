package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

//ExternalDeviceConn: a simple struct to hold external device connection details
type ExternalDeviceConn struct {
	Action                  string `form:"action,omitempty"`
	CID                     string `form:"CID,omitempty"`
	VpcID                   string `form:"vpc_id,omitempty"`
	ConnName                string `form:"connection_name,omitempty"`
	GwName                  string `form:"transit_gw,omitempty"`
	ConnType                string `form:"routing_protocol,omitempty"`
	BgpLocalAsNumber        int    `form:"bgp_local_as_number,omitempty"`
	BgpRemoteAsNumber       int    `form:"external_device_as_number,omitempty"`
	RemoteGatewayIP         string `form:"external_device_ip_address"`
	RemoteSubnet            string `form:"remote_subnet,omitempty"`
	DirectConnect           string `form:"direct_connect,omitempty"`
	PreSharedKey            string `form:"pre_shared_key,omitempty"`
	LocalTunnelIP           string `form:"local_tunnel_ip,omitempty"`
	RemoteTunnelIP          string `form:"remote_tunnel_ip,omitempty"`
	CustomAlgorithms        bool
	Phase1Auth              string `form:"phase1_auth,omitempty"`
	Phase1DhGroups          string `form:"phase1_dh_group,omitempty"`
	Phase1Encryption        string `form:"phase1_encryption,omitempty"`
	Phase2Auth              string `form:"phase2_auth,omitempty"`
	Phase2DhGroups          string `form:"phase2_dh_group,omitempty"`
	Phase2Encryption        string `form:"phase2_encryption,omitempty"`
	HAEnabled               string `form:"enable_ha,omitempty" json:"enable_ha,omitempty"`
	BackupRemoteGatewayIP   string `form:"backup_external_device_ip_address"`
	BackupBgpRemoteAsNumber int    `form:"backup_external_device_as_number,omitempty"`
	BackupPreSharedKey      string `form:"backup_pre_shared_key,omitempty"`
	BackupLocalTunnelIP     string `form:"backup_local_tunnel_ip,omitempty"`
	BackupRemoteTunnelIP    string `form:"backup_remote_tunnel_ip,omitempty"`
	BackupDirectConnect     string `form:"backup_direct_connect,omitempty"`
	EnableEdgeSegmentation  string `form:"connection_policy,omitempty"`
}

type EditExternalDeviceConnDetail struct {
	VpcID                   []string      `json:"vpc_id,omitempty"`
	ConnName                []string      `json:"name,omitempty"`
	ConnType                string        `json:"type,omitempty"`
	TunnelType              string        `json:"tunnel_type,omitempty"`
	GwName                  string        `json:"gw_name,omitempty"`
	Tunnels                 []TunnelInfo  `json:"tunnels,omitempty"`
	RemoteSubnet            string        `json:"real_remote_cidr,omitempty"`
	LocalSubnet             string        `json:"real_local_cidr,omitempty"`
	RemoteCidr              string        `json:"remote_cidr,omitempty"`
	LocalCidr               string        `json:"local_cidr,omitempty"`
	HAEnabled               string        `json:"ha_status,omitempty"`
	PeerType                string        `json:"peer_type,omitempty"`
	RemoteSubnetVirtual     string        `json:"virt_remote_cidr,omitempty"`
	LocalSubnetVirtual      string        `json:"virt_local_cidr,omitempty"`
	Algorithm               AlgorithmInfo `json:"algorithm,omitempty"`
	RouteTableList          []string      `json:"rtbls,omitempty"`
	SslServerPool           []string      `json:"ssl_server_pool,omitempty"`
	DeadPeerDetectionConfig string        `json:"dpd_config,omitempty"`
	EnableActiveActive      string        `json:"active_active_ha,omitempty"`
}

type ExternalDeviceConnDetailResp struct {
	Return  bool                         `json:"return"`
	Results ExternalDeviceConnDetailList `json:"results"`
	Reason  string                       `json:"reason"`
}

type ExternalDeviceConnDetailList struct {
	Connections EditExternalDeviceConnDetail `json:"connections"`
}

func (c *Client) CreateExternalDeviceConn(externalDeviceConn *ExternalDeviceConn) error {
	externalDeviceConn.CID = c.CID
	externalDeviceConn.Action = "connect_transit_gw_to_external_device"
	resp, err := c.Post(c.baseURL, externalDeviceConn)
	if err != nil {
		return errors.New("HTTP Post 'connect_transit_gw_to_external_device' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'connect_transit_gw_to_external_device' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'connect_transit_gw_to_external_device' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetExternalDeviceConnDetail(externalDeviceConn *ExternalDeviceConn) (*ExternalDeviceConn, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'GetExternalDeviceConnDetail': ") + err.Error())
	}
	getExternalDeviceConnDetail := url.Values{}
	getExternalDeviceConnDetail.Add("CID", c.CID)
	getExternalDeviceConnDetail.Add("action", "get_site2cloud_conn_detail")
	getExternalDeviceConnDetail.Add("conn_name", externalDeviceConn.ConnName)
	getExternalDeviceConnDetail.Add("vpc_id", externalDeviceConn.VpcID)
	Url.RawQuery = getExternalDeviceConnDetail.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'get_site2cloud_conn_detail' failed: " + err.Error())
	}
	var data ExternalDeviceConnDetailResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'get_site2cloud_conn_detail' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API 'get_site2cloud_conn_detail' Get failed: " + data.Reason)
	}

	externalDeviceConnDetail := data.Results.Connections
	if len(externalDeviceConnDetail.ConnName) != 0 {
		externalDeviceConn.GwName = externalDeviceConnDetail.GwName
		externalDeviceConn.ConnType = externalDeviceConnDetail.ConnType
		//if externalDeviceConn.ConnType == "mapped" {
		//	externalDeviceConn.RemoteSubnet = externalDeviceConnDetail.RemoteSubnet
		//	externalDeviceConn.LocalSubnet = externalDeviceConnDetail.LocalSubnet
		//	externalDeviceConn.RemoteSubnetVirtual = externalDeviceConnDetail.RemoteSubnetVirtual
		//	externalDeviceConn.LocalSubnetVirtual = externalDeviceConnDetail.LocalSubnetVirtual
		//} else {
		//	externalDeviceConn.RemoteSubnet = externalDeviceConnDetail.RemoteCidr
		//	externalDeviceConn.LocalSubnet = externalDeviceConnDetail.LocalCidr
		//}
		//externalDeviceConn.HAEnabled = externalDeviceConnDetail.HAEnabled
		//for i := range s2cConnDetail.Tunnels {
		//	if externalDeviceConnDetail.Tunnels[i].GwName == externalDeviceConn.GwName {
		//		externalDeviceConn.RemoteGwIP = externalDeviceConnDetail.Tunnels[i].PeerIP
		//	} else {
		//		externalDeviceConn.BackupGwName = externalDeviceConnDetail.Tunnels[i].GwName
		//		externalDeviceConn.RemoteGwIP2 = externalDeviceConnDetail.Tunnels[i].PeerIP
		//	}
		//}
		if externalDeviceConnDetail.Algorithm.Phase1Auth[0] == Phase1AuthDefault &&
			externalDeviceConnDetail.Algorithm.Phase2Auth[0] == Phase2AuthDefault &&
			externalDeviceConnDetail.Algorithm.Phase1DhGroups[0] == Phase1DhGroupDefault &&
			externalDeviceConnDetail.Algorithm.Phase2DhGroups[0] == Phase2DhGroupDefault &&
			externalDeviceConnDetail.Algorithm.Phase1Encrption[0] == Phase1EncryptionDefault &&
			externalDeviceConnDetail.Algorithm.Phase2Encrption[0] == Phase2EncryptionDefault {
			externalDeviceConn.CustomAlgorithms = false
			externalDeviceConn.Phase1Auth = ""
			externalDeviceConn.Phase2Auth = ""
			externalDeviceConn.Phase1DhGroups = ""
			externalDeviceConn.Phase2DhGroups = ""
			externalDeviceConn.Phase1Encryption = ""
			externalDeviceConn.Phase2Encryption = ""
		} else {
			externalDeviceConn.CustomAlgorithms = true
			externalDeviceConn.Phase1Auth = externalDeviceConnDetail.Algorithm.Phase1Auth[0]
			externalDeviceConn.Phase2Auth = externalDeviceConnDetail.Algorithm.Phase2Auth[0]
			externalDeviceConn.Phase1DhGroups = externalDeviceConnDetail.Algorithm.Phase1DhGroups[0]
			externalDeviceConn.Phase2DhGroups = externalDeviceConnDetail.Algorithm.Phase2DhGroups[0]
			externalDeviceConn.Phase1Encryption = externalDeviceConnDetail.Algorithm.Phase1Encrption[0]
			externalDeviceConn.Phase2Encryption = externalDeviceConnDetail.Algorithm.Phase2Encrption[0]
		}

		return externalDeviceConn, nil
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteExternalDeviceConn(externalDeviceConn *ExternalDeviceConn) error {
	externalDeviceConn.CID = c.CID
	externalDeviceConn.Action = "disconnect_transit_gw"
	resp, err := c.Post(c.baseURL, externalDeviceConn)
	if err != nil {
		return errors.New("HTTP Post 'disconnect_transit_gw' failed: " + err.Error())
	}
	log.Printf("zjin030 nothing is wrong here")

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disconnect_transit_gw' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disconnect_transit_gw' Post failed: " + data.Reason)
	}
	return nil
}
