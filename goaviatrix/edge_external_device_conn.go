package goaviatrix

import "context"

type EdgeExternalDeviceConn struct {
	Action                 string `json:"action,omitempty"`
	CID                    string `json:"CID,omitempty"`
	VpcID                  string `json:"vpc_id,omitempty"`
	ConnectionName         string `json:"conn_name,omitempty"`
	GwName                 string `json:"gw_name,omitempty"`
	ConnectionType         string `json:"routing_protocol,omitempty"`
	BgpLocalAsNum          int    `json:"local_asn,omitempty"`
	BgpRemoteAsNum         int    `json:"external_device_asn,omitempty"`
	RemoteGatewayIP        string `json:"external_device_ip_address,omitempty"`
	RemoteSubnet           string `json:"remote_subnet,omitempty"`
	DirectConnect          string `json:"direct_connect,omitempty"`
	PreSharedKey           string `json:"pre_shared_key,omitempty"`
	LocalTunnelCidr        string `json:"local_tunnel_ip,omitempty"`
	RemoteTunnelCidr       string `json:"remote_tunnel_ip,omitempty"`
	CustomAlgorithms       bool
	Phase1Auth             string `json:"phase1_authentication,omitempty"`
	Phase1DhGroups         string `json:"phase1_dh_groups,omitempty"`
	Phase1Encryption       string `json:"phase1_encryption,omitempty"`
	Phase2Auth             string `json:"phase2_authentication,omitempty"`
	Phase2DhGroups         string `json:"phase2_dh_groups,omitempty"`
	Phase2Encryption       string `json:"phase2_encryption,omitempty"`
	HAEnabled              string `json:"enable_ha,omitempty"`
	BackupRemoteGatewayIP  string `json:"backup_external_device_ip_address,omitempty"`
	BackupBgpRemoteAsNum   int    `json:"external_device_backup_asn,omitempty"`
	BackupPreSharedKey     string `json:"backup_pre_shared_key,omitempty"`
	BackupLocalTunnelCidr  string `json:"backup_local_tunnel_ip,omitempty"`
	BackupRemoteTunnelCidr string `json:"backup_remote_tunnel_ip,omitempty"`
	BackupDirectConnect    string `json:"backup_direct_connect,omitempty"`
	EnableEdgeSegmentation string `json:"connection_policy,omitempty"`
	EnableIkev2            string `json:"enable_ikev2,omitempty"`
	ManualBGPCidrs         []string
	TunnelProtocol         string `json:"tunnel_protocol,omitempty"`
	EnableBgpLanActiveMesh bool   `json:"bgp_lan_activemesh,omitempty"`
	PeerVnetId             string `json:"peer_vnet_id,omitempty"`
	RemoteLanIP            string `json:"remote_lan_ip,omitempty"`
	LocalLanIP             string `json:"local_lan_ip,omitempty"`
	BackupRemoteLanIP      string `json:"backup_remote_lan_ip,omitempty"`
	BackupLocalLanIP       string `json:"backup_local_lan_ip,omitempty"`
	EventTriggeredHA       bool
	EnableJumboFrame       bool
	Phase1LocalIdentifier  string
	Phase1RemoteIdentifier string
	PrependAsPath          string
	BgpMd5Key              string       `json:"bgp_md5_key,omitempty"`
	BackupBgpMd5Key        string       `json:"backup_bgp_md5_key,omitempty"`
	AuthType               string       `json:"auth_type,omitempty"`
	EnableEdgeUnderlay     bool         `json:"edge_underlay,omitempty"`
	RemoteCloudType        string       `json:"remote_cloud_type,omitempty"`
	BgpMd5KeyChanged       bool         `json:"bgp_md5_key_changed,omitempty"`
	BgpBfdConfig           BgpBfdConfig `json:"bgp_bfd_params,omitempty"`
	EnableBfd              bool         `json:"bgp_bfd_enabled,omitempty"`
	// Multihop must not use "omitempty"; It defaults to true and omitempty
	// breaks that.
	EnableBgpMultihop bool `form:"enable_bgp_multihop"`
	DisableActivemesh bool
	TunnelSrcIP       string
}

func (c *Client) CreateEdgeExternalDeviceConn(edgeExternalDeviceConn *EdgeExternalDeviceConn) (string, error) {
	edgeExternalDeviceConn.CID = c.CID
	edgeExternalDeviceConn.Action = "transit_connect_external_device"

	return c.PostAPIContext2WithResult(context.Background(), nil, edgeExternalDeviceConn.Action, edgeExternalDeviceConn, BasicCheck)
}

func (c *Client) DeleteEdgeExternalDeviceConn(edgeExternalDeviceConn *EdgeExternalDeviceConn) error {
	edgeExternalDeviceConn.CID = c.CID
	edgeExternalDeviceConn.Action = "transit_disconnect"

	return c.PostAPIContext2(context.Background(), nil, edgeExternalDeviceConn.Action, edgeExternalDeviceConn, BasicCheck)
}
