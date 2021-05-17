package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//ExternalDeviceConn: a simple struct to hold external device connection details
type ExternalDeviceConn struct {
	Action                 string `form:"action,omitempty"`
	CID                    string `form:"CID,omitempty"`
	VpcID                  string `form:"vpc_id,omitempty"`
	ConnectionName         string `form:"connection_name,omitempty"`
	GwName                 string `form:"transit_gw,omitempty"`
	ConnectionType         string `form:"routing_protocol,omitempty"`
	BgpLocalAsNum          int    `form:"bgp_local_as_number,omitempty"`
	BgpRemoteAsNum         int    `form:"external_device_as_number,omitempty"`
	RemoteGatewayIP        string `form:"external_device_ip_address"`
	RemoteSubnet           string `form:"remote_subnet,omitempty"`
	DirectConnect          string `form:"direct_connect,omitempty"`
	PreSharedKey           string `form:"pre_shared_key,omitempty"`
	LocalTunnelCidr        string `form:"local_tunnel_ip,omitempty"`
	RemoteTunnelCidr       string `form:"remote_tunnel_ip,omitempty"`
	CustomAlgorithms       bool
	Phase1Auth             string `form:"phase1_authentication,omitempty"`
	Phase1DhGroups         string `form:"phase1_dh_groups,omitempty"`
	Phase1Encryption       string `form:"phase1_encryption,omitempty"`
	Phase2Auth             string `form:"phase2_authentication,omitempty"`
	Phase2DhGroups         string `form:"phase2_dh_groups,omitempty"`
	Phase2Encryption       string `form:"phase2_encryption,omitempty"`
	HAEnabled              string `form:"enable_ha,omitempty" json:"enable_ha,omitempty"`
	BackupRemoteGatewayIP  string `form:"backup_external_device_ip_address"`
	BackupBgpRemoteAsNum   int    `form:"backup_external_device_as_number,omitempty"`
	BackupPreSharedKey     string `form:"backup_pre_shared_key,omitempty"`
	BackupLocalTunnelCidr  string `form:"backup_local_tunnel_ip,omitempty"`
	BackupRemoteTunnelCidr string `form:"backup_remote_tunnel_ip,omitempty"`
	BackupDirectConnect    string `form:"backup_direct_connect,omitempty"`
	EnableEdgeSegmentation string `form:"connection_policy,omitempty"`
	EnableIkev2            string `form:"enable_ikev2,omitempty"`
	ManualBGPCidrs         []string
	TunnelProtocol         string `form:"tunnel_protocol,omitempty"`
	PeerVnetId             string `form:"peer_vnet_id,omitempty"`
	RemoteLanIP            string `form:"remote_lan_ip,omitempty"`
	LocalLanIP             string `form:"local_lan_ip,omitempty"`
	BackupRemoteLanIP      string `form:"backup_remote_lan_ip,omitempty"`
	BackupLocalLanIP       string `form:"backup_local_lan_ip,omitempty"`
	EventTriggeredHA       bool
	Phase1RemoteIdentifier string
}

type EditExternalDeviceConnDetail struct {
	VpcID                  []string      `json:"vpc_id,omitempty"`
	ConnectionName         []string      `json:"name,omitempty"`
	GwName                 string        `json:"gw_name,omitempty"`
	BgpLocalAsNum          string        `json:"bgp_local_asn_number,omitempty"`
	BgpRemoteAsNum         string        `json:"bgp_remote_asn_number,omitempty"`
	BgpStatus              string        `json:"bgp_status,omitempty"`
	RemoteGatewayIP        string        `json:"peer_ip,omitempty"`
	RemoteSubnet           string        `json:"remote_cidr,omitempty"`
	DirectConnect          bool          `json:"direct_connect_primary,omitempty"`
	LocalTunnelCidr        string        `json:"bgp_local_ip,omitempty"`
	RemoteTunnelCidr       string        `json:"bgp_remote_ip,omitempty"`
	Algorithm              AlgorithmInfo `json:"algorithm,omitempty"`
	HAEnabled              string        `json:"ha_status,omitempty"`
	BackupBgpRemoteAsNum   string        `json:"bgp_remote_backup_asn_number,omitempty"`
	BackupLocalTunnelCidr  string        `json:"bgp_backup_local_ip,omitempty"`
	BackupRemoteTunnelCidr string        `json:"bgp_backup_remote_ip,omitempty"`
	BackupDirectConnect    bool          `json:"direct_connect_backup,omitempty"`
	EnableEdgeSegmentation bool          `json:"enable_edge_segmentation,omitempty"`
	Tunnels                []TunnelInfo  `json:"tunnels,omitempty"`
	ActiveActiveHA         string        `json:"active_active_ha,omitempty"`
	ManualBGPCidrs         []string      `json:"conn_bgp_manual_advertise_cidrs"`
	BackupRemoteGatewayIP  string
	PreSharedKey           string
	BackupPreSharedKey     string
	IkeVer                 string `json:"ike_ver"`
	PeerVnetId             string `json:"peer_vnet_id"`
	RemoteLanIP            string `json:"remote_lan_ip"`
	LocalLanIP             string `json:"local_lan_ip"`
	BackupRemoteLanIP      string `json:"backup_remote_lan_ip"`
	BackupLocalLanIP       string `json:"backup_local_lan_ip"`
	EventTriggeredHA       string `json:"event_triggered_ha"`
	Phase1RemoteIdentifier string `json:"phase1_remote_id"`
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
	params := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"conn_name": externalDeviceConn.ConnectionName,
		"vpc_id":    externalDeviceConn.VpcID,
	}
	checkFunc := func(action, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return errors.New("Rest API 'get_site2cloud_conn_detail' Get failed: " + reason)
		}
		return nil
	}
	var data ExternalDeviceConnDetailResp
	err := c.GetAPI(&data, params["action"], params, checkFunc)
	if err != nil {
		return nil, err
	}

	externalDeviceConnDetail := data.Results.Connections
	if len(externalDeviceConnDetail.ConnectionName) != 0 {
		if len(externalDeviceConnDetail.VpcID) != 0 {
			externalDeviceConn.VpcID = externalDeviceConnDetail.VpcID[0]
		}

		externalDeviceConn.ConnectionName = externalDeviceConnDetail.ConnectionName[0]
		externalDeviceConn.GwName = externalDeviceConnDetail.GwName

		if externalDeviceConnDetail.BgpStatus == "enabled" || externalDeviceConnDetail.BgpStatus == "Enabled" {
			bgpLocalAsNumber, _ := strconv.Atoi(externalDeviceConnDetail.BgpLocalAsNum)
			externalDeviceConn.BgpLocalAsNum = bgpLocalAsNumber
			bgpRemoteAsNumber, _ := strconv.Atoi(externalDeviceConnDetail.BgpRemoteAsNum)
			externalDeviceConn.BgpRemoteAsNum = bgpRemoteAsNumber
			externalDeviceConn.ConnectionType = "bgp"
			if len(externalDeviceConnDetail.Tunnels) >= 1 {
				tunnelProtocol := externalDeviceConnDetail.Tunnels[0].TunnelProtocol
				// LAN tunnel protocol is defined in the backend as "N/A(LAN)".
				// Here we clean that up to be just "LAN" for Terraform users.
				if strings.Contains(tunnelProtocol, "LAN") {
					tunnelProtocol = "LAN"
				}
				externalDeviceConn.TunnelProtocol = tunnelProtocol
			}
		} else {
			externalDeviceConn.RemoteSubnet = externalDeviceConnDetail.RemoteSubnet
			externalDeviceConn.ConnectionType = "static"
		}
		externalDeviceConn.RemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[0]

		// GRE and LAN tunnels cannot set Algorithms
		if externalDeviceConn.TunnelProtocol != "GRE" && externalDeviceConn.TunnelProtocol != "LAN" {
			if externalDeviceConnDetail.Algorithm.Phase1Auth[0] == "SHA-256" &&
				externalDeviceConnDetail.Algorithm.Phase2Auth[0] == "HMAC-SHA-256" &&
				externalDeviceConnDetail.Algorithm.Phase1DhGroups[0] == "14" &&
				externalDeviceConnDetail.Algorithm.Phase2DhGroups[0] == "14" &&
				externalDeviceConnDetail.Algorithm.Phase1Encrption[0] == "AES-256-CBC" &&
				externalDeviceConnDetail.Algorithm.Phase2Encrption[0] == "AES-256-CBC" {
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
		}

		if externalDeviceConnDetail.DirectConnect {
			externalDeviceConn.DirectConnect = "enabled"
		} else {
			externalDeviceConn.DirectConnect = "disabled"
		}

		backupBgpRemoteAsNumber := 0
		if externalDeviceConnDetail.BackupBgpRemoteAsNum != "" {
			backupBgpRemoteAsNumberRead, _ := strconv.Atoi(externalDeviceConnDetail.BackupBgpRemoteAsNum)
			backupBgpRemoteAsNumber = backupBgpRemoteAsNumberRead
		}
		if externalDeviceConnDetail.HAEnabled == "enabled" {
			if len(externalDeviceConnDetail.Tunnels) == 2 {
				remoteIP := strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")
				if len(remoteIP) == 2 {
					if remoteIP[0] == remoteIP[1] {
						externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr + "," + externalDeviceConnDetail.BackupLocalTunnelCidr
						externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr + "," + externalDeviceConnDetail.BackupRemoteTunnelCidr
						externalDeviceConn.HAEnabled = "disabled"
					} else {
						externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr
						externalDeviceConn.BackupLocalTunnelCidr = externalDeviceConnDetail.BackupLocalTunnelCidr
						externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr
						externalDeviceConn.BackupRemoteTunnelCidr = externalDeviceConnDetail.BackupRemoteTunnelCidr
						externalDeviceConn.BackupRemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[1]
						externalDeviceConn.HAEnabled = "enabled"
						externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
					}
				} else if len(remoteIP) == 4 {
					if remoteIP[0] == remoteIP[2] && remoteIP[1] == remoteIP[3] {
						externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr + "," + externalDeviceConnDetail.BackupLocalTunnelCidr
						externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr + "," + externalDeviceConnDetail.BackupRemoteTunnelCidr
						externalDeviceConn.RemoteGatewayIP = remoteIP[0] + "," + remoteIP[1]
						externalDeviceConn.HAEnabled = "disabled"
					}
				}
			} else if len(externalDeviceConnDetail.Tunnels) == 4 {
				externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr
				externalDeviceConn.BackupLocalTunnelCidr = externalDeviceConnDetail.BackupLocalTunnelCidr
				externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr
				externalDeviceConn.BackupRemoteTunnelCidr = externalDeviceConnDetail.BackupRemoteTunnelCidr
				externalDeviceConn.BackupRemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[1]
				externalDeviceConn.HAEnabled = "enabled"
				externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
			}
		} else {
			externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr
			externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr
			if len(externalDeviceConnDetail.Tunnels) == 2 {
				externalDeviceConn.BackupLocalTunnelCidr = externalDeviceConnDetail.BackupLocalTunnelCidr
				externalDeviceConn.BackupRemoteTunnelCidr = externalDeviceConnDetail.BackupRemoteTunnelCidr
				externalDeviceConn.BackupRemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[1]
				externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
				externalDeviceConn.HAEnabled = "enabled"
			} else {
				externalDeviceConn.HAEnabled = "disabled"
			}
		}

		if externalDeviceConnDetail.BackupDirectConnect {
			externalDeviceConn.BackupDirectConnect = "enabled"
		} else {
			externalDeviceConn.BackupDirectConnect = "disabled"
		}
		if externalDeviceConnDetail.EnableEdgeSegmentation {
			externalDeviceConn.EnableEdgeSegmentation = "enabled"
		} else {
			externalDeviceConn.EnableEdgeSegmentation = "disabled"
		}
		externalDeviceConn.ManualBGPCidrs = externalDeviceConnDetail.ManualBGPCidrs

		if externalDeviceConnDetail.IkeVer == "2" {
			externalDeviceConn.EnableIkev2 = "enabled"
		} else {
			externalDeviceConn.EnableIkev2 = "disabled"
		}
		externalDeviceConn.EventTriggeredHA = externalDeviceConnDetail.EventTriggeredHA == "enabled"
		externalDeviceConn.PeerVnetId = externalDeviceConnDetail.PeerVnetId
		externalDeviceConn.RemoteLanIP = externalDeviceConnDetail.RemoteLanIP
		externalDeviceConn.LocalLanIP = externalDeviceConnDetail.LocalLanIP
		externalDeviceConn.BackupRemoteLanIP = externalDeviceConnDetail.BackupRemoteLanIP
		externalDeviceConn.BackupLocalLanIP = externalDeviceConnDetail.BackupLocalLanIP
		externalDeviceConn.Phase1RemoteIdentifier = externalDeviceConnDetail.Phase1RemoteIdentifier

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

func TransitExternalDeviceConnPh1RemoteIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("ha_enabled") {
		return false
	}

	ip := d.Get("remote_gateway_ip").(string)
	haip := d.Get("backup_remote_gateway_ip").(string)
	o, n := d.GetChange("phase1_remote_identifier")
	haEnabled := d.Get("ha_enabled").(bool)

	ph1RemoteIdListOld := ExpandStringList(o.([]interface{}))
	ph1RemoteIdListNew := ExpandStringList(n.([]interface{}))

	if len(ph1RemoteIdListOld) != 0 && len(ph1RemoteIdListNew) != 0 {
		if haEnabled {
			if len(ph1RemoteIdListNew) != 2 || len(ph1RemoteIdListOld) != 2 {
				return false
			}
			return ph1RemoteIdListOld[0] == ip && ph1RemoteIdListNew[0] == ip &&
				strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && strings.TrimSpace(ph1RemoteIdListNew[1]) == haip
		} else {
			if len(ph1RemoteIdListNew) != 1 {
				return false
			}
			return ph1RemoteIdListOld[0] == ip && ph1RemoteIdListNew[0] == ip
		}
	}

	if !haEnabled {
		if len(ph1RemoteIdListOld) == 1 && ph1RemoteIdListOld[0] == ip && len(ph1RemoteIdListNew) == 0 {
			return true
		}
	}

	if haEnabled {
		if len(ph1RemoteIdListOld) == 2 && ph1RemoteIdListOld[0] == ip && strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && len(ph1RemoteIdListNew) == 0 {
			return true
		}
	}

	return false
}
