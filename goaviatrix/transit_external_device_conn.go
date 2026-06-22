package goaviatrix

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ExternalDeviceConn: a simple struct to hold external device connection details
const (
	Enabled  = "enabled"
	Disabled = "disabled"

	defaultBfdReceiveInterval  = 300
	defaultBfdTransmitInterval = 300
	defaultBfdMultiplier       = 3
)

// on returns "enabled" if b is true, "disabled" otherwise
func on(b bool) string {
	if b {
		return Enabled
	}
	return Disabled
}

// join2 joins two strings with a comma separator
func join2(a, b string) string {
	return strings.Join([]string{a, b}, ",")
}

var defaultBfdConfig = BgpBfdConfig{
	TransmitInterval: defaultBfdTransmitInterval,
	ReceiveInterval:  defaultBfdReceiveInterval,
	Multiplier:       defaultBfdMultiplier,
}

type ExternalDeviceConn struct {
	Action                     string `form:"action,omitempty"`
	CID                        string `form:"CID,omitempty"`
	VpcID                      string `form:"vpc_id,omitempty"`
	ConnectionName             string `form:"connection_name,omitempty"`
	GwName                     string `form:"transit_gw,omitempty"`
	ConnectionType             string `form:"routing_protocol,omitempty"`
	BgpLocalAsNum              int    `form:"bgp_local_as_number,omitempty"`
	BgpRemoteAsNum             int    `form:"external_device_as_number,omitempty"`
	BgpSendCommunities         string `json:"conn_bgp_send_communities,omitempty"`
	BgpSendCommunitiesAdditive bool   `json:"conn_bgp_send_communities_additive,omitempty"`
	BgpSendCommunitiesBlock    bool   `json:"conn_bgp_send_communities_block,omitempty"`
	RemoteGatewayIP            string `form:"external_device_ip_address,omitempty"`
	RemoteSubnet               string `form:"remote_subnet,omitempty"`
	LocalSubnet                string `form:"local_subnet,omitempty"`
	DirectConnect              string `form:"direct_connect,omitempty"`
	PreSharedKey               string `form:"pre_shared_key,omitempty"`
	LocalTunnelCidr            string `form:"local_tunnel_ip,omitempty"`
	RemoteTunnelCidr           string `form:"remote_tunnel_ip,omitempty"`
	CustomAlgorithms           bool
	Phase1Auth                 string `form:"phase1_authentication,omitempty"`
	Phase1DhGroups             string `form:"phase1_dh_groups,omitempty"`
	Phase1Encryption           string `form:"phase1_encryption,omitempty"`
	Phase2Auth                 string `form:"phase2_authentication,omitempty"`
	Phase2DhGroups             string `form:"phase2_dh_groups,omitempty"`
	Phase2Encryption           string `form:"phase2_encryption,omitempty"`
	HAEnabled                  string `form:"enable_ha,omitempty" json:"enable_ha,omitempty"`
	BackupRemoteGatewayIP      string `form:"backup_external_device_ip_address,omitempty"`
	BackupBgpRemoteAsNum       int    `form:"backup_external_device_as_number,omitempty"`
	BackupPreSharedKey         string `form:"backup_pre_shared_key,omitempty"`
	BackupLocalTunnelCidr      string `form:"backup_local_tunnel_ip,omitempty"`
	BackupRemoteTunnelCidr     string `form:"backup_remote_tunnel_ip,omitempty"`
	BackupDirectConnect        string `form:"backup_direct_connect,omitempty"`
	EnableEdgeSegmentation     string `form:"connection_policy,omitempty"`
	EnableIkev2                string `form:"enable_ikev2,omitempty"`
	ManualBGPCidrs             []string
	TunnelProtocol             string `form:"tunnel_protocol,omitempty"`
	EnableBgpLanActiveMesh     bool   `form:"bgp_lan_activemesh,omitempty"`
	PeerVnetID                 string `form:"peer_vnet_id,omitempty"`
	RemoteLanIP                string `form:"remote_lan_ip,omitempty"`
	LocalLanIP                 string `form:"local_lan_ip,omitempty"`
	BackupRemoteLanIP          string `form:"backup_remote_lan_ip,omitempty"`
	BackupLocalLanIP           string `form:"backup_local_lan_ip,omitempty"`
	EventTriggeredHA           bool
	EnableJumboFrame           bool `form:"jumbo_frame,omitempty"`
	Phase1LocalIdentifier      string
	Phase1RemoteIdentifier     string
	PrependAsPath              string
	BgpMd5Key                  string       `form:"bgp_md5_key,omitempty"`
	BackupBgpMd5Key            string       `form:"backup_bgp_md5_key,omitempty"`
	AuthType                   string       `form:"auth_type,omitempty"`
	EnableEdgeUnderlay         bool         `form:"edge_underlay,omitempty"`
	RemoteCloudType            string       `form:"remote_cloud_type,omitempty"`
	BgpMd5KeyChanged           bool         `form:"bgp_md5_key_changed,omitempty"`
	BgpBfdConfig               BgpBfdConfig `form:"bgp_bfd_params,omitempty"`
	EnableBfd                  bool         `form:"bgp_bfd_enabled,omitempty"`
	// Multihop must not use "omitempty", it defaults to true and omitempty
	// breaks that.
	EnableBgpMultihop        bool   `form:"enable_bgp_multihop"`
	DisableActivemesh        bool   `form:"disable_activemesh,omitempty" json:"disable_activemesh,omitempty"`
	ProxyIdEnabled           bool   `form:"proxy_id_enabled,omitempty"`
	TunnelSrcIP              string `form:"local_device_ip,omitempty"`
	EnableIpv6               bool   `form:"ipv6_enabled,omitempty"`
	ExternalDeviceIPv6       string `form:"external_device_ipv6,omitempty"`
	ExternalDeviceBackupIPv6 string `form:"external_device_backup_ipv6,omitempty"`
	RemoteLanIPv6            string `form:"remote_lan_ipv6_ip,omitempty"`
	BackupRemoteLanIPv6      string `form:"backup_remote_lan_ipv6_ip,omitempty"`
}

type EditExternalDeviceConnDetail struct {
	VpcID                      []string      `json:"vpc_id,omitempty"`
	ConnectionName             []string      `json:"name,omitempty"`
	GwName                     string        `json:"gw_name,omitempty"`
	BgpLocalAsNum              string        `json:"bgp_local_asn_number,omitempty"`
	BgpRemoteAsNum             string        `json:"bgp_remote_asn_number,omitempty"`
	BgpStatus                  string        `json:"bgp_status,omitempty"`
	BgpSendCommunities         string        `json:"conn_bgp_send_communities,omitempty"`
	BgpSendCommunitiesAdditive bool          `json:"conn_bgp_send_communities_additive,omitempty"`
	BgpSendCommunitiesBlock    bool          `json:"conn_bgp_send_communities_block,omitempty"`
	EnableBgpLanActiveMesh     bool          `json:"bgp_lan_activemesh,omitempty"`
	RemoteGatewayIP            string        `json:"peer_ip,omitempty"`
	RemoteSubnet               string        `json:"remote_cidr,omitempty"`
	LocalSubnet                string        `json:"local_cidr,omitempty"`
	DirectConnect              bool          `json:"direct_connect_primary,omitempty"`
	LocalTunnelCidr            string        `json:"bgp_local_ip,omitempty"`
	RemoteTunnelCidr           string        `json:"bgp_remote_ip,omitempty"`
	Algorithm                  AlgorithmInfo `json:"algorithm,omitempty"`
	HAEnabled                  string        `json:"ha_status,omitempty"`
	BackupBgpRemoteAsNum       string        `json:"bgp_remote_backup_asn_number,omitempty"`
	BackupLocalTunnelCidr      string        `json:"bgp_backup_local_ip,omitempty"`
	BackupRemoteTunnelCidr     string        `json:"bgp_backup_remote_ip,omitempty"`
	BackupDirectConnect        bool          `json:"direct_connect_backup,omitempty"`
	EnableEdgeSegmentation     bool          `json:"enable_edge_segmentation,omitempty"`
	Tunnels                    []TunnelInfo  `json:"tunnels,omitempty"`
	ActiveActiveHA             string        `json:"active_active_ha,omitempty"`
	ManualBGPCidrs             []string      `json:"conn_bgp_manual_advertise_cidrs,omitempty"`
	BackupRemoteGatewayIP      string
	PreSharedKey               string
	BackupPreSharedKey         string
	IkeVer                     string         `json:"ike_ver,omitempty"`
	PeerVnetID                 string         `json:"peer_vnet_id,omitempty"`
	RemoteLanIP                string         `json:"remote_lan_ip,omitempty"`
	LocalLanIP                 string         `json:"local_lan_ip,omitempty"`
	BackupRemoteLanIP          string         `json:"backup_remote_lan_ip,omitempty"`
	BackupLocalLanIP           string         `json:"backup_local_lan_ip,omitempty"`
	EventTriggeredHA           string         `json:"event_triggered_ha,omitempty"`
	Phase1LocalIdentifier      string         `json:"ph1_identifier,omitempty"`
	Phase1RemoteIdentifier     string         `json:"phase1_remote_id,omitempty"`
	PrependAsPath              string         `json:"conn_bgp_prepend_as_path,omitempty"`
	EnableJumboFrame           bool           `json:"jumbo_frame,omitempty"`
	WanUnderlay                bool           `json:"wan_underlay,omitempty"`
	RemoteCloudType            string         `json:"remote_cloud_type,omitempty"`
	BgpBfdConfig               map[string]int `json:"bgp_bfd_params,omitempty"`
	EnableBfd                  bool           `json:"bgp_bfd_enabled,omitempty"`
	EnableBgpMultihop          bool           `json:"bgp_multihop_enabled,omitempty"`
	DisableActivemesh          bool           `json:"disable_activemesh,omitempty"`
	ProxyIdEnabled             bool           `json:"proxy_id_enabled,omitempty"`
	TunnelSrcIP                string         `json:"local_device_ip,omitempty"`
	TunnelType                 string         `json:"tunnel_type,omitempty"`
	EnableIpv6                 bool           `json:"ipv6_enabled,omitempty"`
	ExternalDeviceIPv6         string         `json:"bgp_remote_ipv6,omitempty"`
	ExternalDeviceBackupIPv6   string         `json:"bgp_backup_remote_ipv6,omitempty"`
	ExternalLocalIPv6          string         `json:"bgp_local_ipv6,omitempty"`
	ExternalBackupLocalIPv6    string         `json:"bgp_backup_local_ipv6,omitempty"`
	RemoteLanIPv6              string         `json:"remote_lan_ipv6_ip,omitempty"`
	BackupRemoteLanIPv6        string         `json:"backup_remote_lan_ipv6_ip,omitempty"`
}

type BgpBfdConfig struct {
	TransmitInterval int `json:"tx_interval"`
	ReceiveInterval  int `json:"rx_interval"`
	Multiplier       int `json:"multiplier"`
}

type EditBgpMd5Key struct {
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	ConnectionName string `form:"conn_name,omitempty"`
	GwName         string `form:"gateway_name,omitempty"`
	BgpMd5Key      string `form:"bgp_md5_key,omitempty"`
	BgpRemoteIP    string `form:"bgp_remote_ip,omitempty"`
}

type BgpSendCommunities struct {
	Action              string `form:"action,omitempty"`
	CID                 string `form:"CID,omitempty"`
	ConnectionName      string `form:"connection_name,omitempty"`
	GwName              string `form:"gateway_name,omitempty"`
	ConnSendCommunities string `form:"connection_bgp_send_communities,omitempty"`
	ConnSendAdditive    bool   `form:"connection_bgp_send_communities_additive,omitempty"`
	ConnSendBlock       bool   `form:"connection_bgp_send_communities_block,omitempty"`
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

	return c.PostAPI(externalDeviceConn.Action, externalDeviceConn, BasicCheck)
}

func (c *Client) GetExternalDeviceConnDetail(externalDeviceConn *ExternalDeviceConn, localGateway *Gateway) (*ExternalDeviceConn, error) {
	params := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"conn_name": externalDeviceConn.ConnectionName,
		"vpc_id":    externalDeviceConn.VpcID,
	}
	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	var data ExternalDeviceConnDetailResp
	err := c.GetAPI(&data, params["action"], params, checkFunc)
	if err != nil {
		return nil, err
	}

	externalDeviceConnDetail := data.Results.Connections
	if len(externalDeviceConnDetail.ConnectionName) == 0 {
		return nil, ErrNotFound
	}

	// Extract connection details into separate helper functions
	populateBasicConnectionInfo(externalDeviceConn, externalDeviceConnDetail)
	populateConnectionTypeInfo(externalDeviceConn, externalDeviceConnDetail)
	populateAlgorithmInfo(externalDeviceConn, externalDeviceConnDetail)
	populateDirectConnectInfo(externalDeviceConn, externalDeviceConnDetail)

	backupBgpRemoteAsNumber := parseBackupBgpRemoteAsNumber(externalDeviceConnDetail.BackupBgpRemoteAsNum)
	populateBgpSendCommunitiesInfo(externalDeviceConn, externalDeviceConnDetail)

	if externalDeviceConn.TunnelProtocol != "LAN" {
		populateNonLANTunnelInfo(externalDeviceConn, externalDeviceConnDetail, localGateway, backupBgpRemoteAsNumber)
	} else {
		populateLANTunnelInfo(externalDeviceConn, externalDeviceConnDetail, backupBgpRemoteAsNumber)
	}

	populateAdditionalConnectionInfo(externalDeviceConn, externalDeviceConnDetail)

	populateIPv6ConnectionInfo(externalDeviceConn, externalDeviceConnDetail)

	return externalDeviceConn, nil
}

func (c *Client) DeleteExternalDeviceConn(externalDeviceConn *ExternalDeviceConn) error {
	externalDeviceConn.CID = c.CID
	externalDeviceConn.Action = "disconnect_transit_gw"

	return c.PostAPI(externalDeviceConn.Action, externalDeviceConn, BasicCheck)
}

func ExternalDeviceConnBgpBfdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// In the case where BFD is disabled, we need to *not* suppress any
	// diffs to "bgp_bfd", otherwise BFD config can be left in the state.
	// That can break things, as BFD config is not allowed when BFD is
	// disabled.
	enabled, ok := d.Get("enable_bfd").(bool)
	if !enabled || !ok {
		return false
	}

	// You might expect that GetChange("bgp_bfd") would return old and
	// new lists that could be compared. Unfortunately, it doesn't seem to
	// work that way. The API can report identical values for "bgp_bfd",
	// (along with HasChange() == false), while at the same time reporting
	// a change for the first element ("bgp_bfd.0"). Fortunately for us, we
	// enforce only a single element in bgp_bfd, so we can just check that.
	// terraform will auto-populate all defaults - including for an empty
	// list - so all we need to do is compare the two elements.
	// The strong consensus on the internet is that SDKv2 was simply not
	// designed for for this sort of thing and the recommended solution
	// is to migrate to the plugin framework. Unfortunately that's not
	// a small undertaking.
	o, n := d.GetChange("bgp_bfd.0")
	return reflect.DeepEqual(o, n)
}

func TransitExternalDeviceConnPh1RemoteIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("ha_enabled") {
		return false
	}

	ip, ok := d.Get("remote_gateway_ip").(string)
	if !ok {
		return false
	}
	ipList := strings.Split(ip, ",")

	haip, ok := d.Get("backup_remote_gateway_ip").(string)
	if !ok {
		return false
	}

	o, n := d.GetChange("phase1_remote_identifier")

	haEnabled, ok := d.Get("ha_enabled").(bool)
	if !ok {
		return false
	}

	oList, ok := o.([]interface{})
	if !ok {
		return false
	}
	nList, ok := n.([]interface{})
	if !ok {
		return false
	}

	ph1RemoteIdListOld := ExpandStringList(oList)
	ph1RemoteIdListNew := ExpandStringList(nList)

	if len(ph1RemoteIdListOld) != 0 && len(ph1RemoteIdListNew) != 0 {
		if haEnabled {
			if len(ph1RemoteIdListNew) != 2 || len(ph1RemoteIdListOld) != 2 {
				if len(ph1RemoteIdListNew) == 1 && len(ph1RemoteIdListOld) == 1 {
					return ph1RemoteIdListOld[0] == ipList[0] && ph1RemoteIdListNew[0] == ipList[0]
				}
				return false
			}
			return ph1RemoteIdListOld[0] == ipList[0] && ph1RemoteIdListNew[0] == ipList[0] &&
				strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && strings.TrimSpace(ph1RemoteIdListNew[1]) == haip
		} else {
			if len(ph1RemoteIdListNew) == 1 && len(ph1RemoteIdListOld) == 1 {
				return ph1RemoteIdListOld[0] == ipList[0] && ph1RemoteIdListNew[0] == ipList[0]
			} else if len(ph1RemoteIdListNew) == 2 && len(ph1RemoteIdListOld) == 2 && len(ipList) == 2 {
				return strings.TrimSpace(ph1RemoteIdListOld[0]) == strings.TrimSpace(ipList[0]) &&
					strings.TrimSpace(ph1RemoteIdListOld[1]) == strings.TrimSpace(ipList[1]) &&
					strings.TrimSpace(ph1RemoteIdListNew[0]) == strings.TrimSpace(ipList[0]) &&
					strings.TrimSpace(ph1RemoteIdListNew[1]) == strings.TrimSpace(ipList[1])
			} else {
				return false
			}
		}
	}

	return false
}

func CreateBgpBfdConfig(bfd map[string]interface{}) BgpBfdConfig {
	// Set default values
	transmitInterval := defaultBfdTransmitInterval
	receiveInterval := defaultBfdReceiveInterval
	multiplier := defaultBfdMultiplier

	// Override defaults if provided in bfd1
	if value, ok := bfd["transmit_interval"].(int); ok {
		transmitInterval = value
	}
	if value, ok := bfd["receive_interval"].(int); ok {
		receiveInterval = value
	}
	if value, ok := bfd["multiplier"].(int); ok {
		multiplier = value
	}

	// Create and return BgpBfdConfig instance
	bfd2 := BgpBfdConfig{
		TransmitInterval: transmitInterval,
		ReceiveInterval:  receiveInterval,
		Multiplier:       multiplier,
	}
	return bfd2
}

func GetUpdatedBgpBfdConfig(bgpBfdConfig []interface{}) BgpBfdConfig {
	var bgpBfd BgpBfdConfig

	if len(bgpBfdConfig) > 0 {
		// get the user provided bgp bfd config
		for _, v := range bgpBfdConfig {
			bfdConfig, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			bgpBfd = CreateBgpBfdConfig(bfdConfig)
		}
	} else {
		// use default bgp bfd config
		bgpBfd = defaultBfdConfig
	}

	return bgpBfd
}

func (c *Client) EditTransitExternalDeviceConnASPathPrepend(externalDeviceConn *ExternalDeviceConn, prependASPath []string) error {
	action := "edit_transit_connection_as_path_prepend"
	return c.PostAPI(action, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		GatewayName    string `form:"gateway_name"`
		ConnectionName string `form:"connection_name"`
		PrependASPath  string `form:"connection_as_path_prepend"`
	}{
		CID:            c.CID,
		Action:         action,
		GatewayName:    externalDeviceConn.GwName,
		ConnectionName: externalDeviceConn.ConnectionName,
		PrependASPath:  strings.Join(prependASPath, ","),
	}, BasicCheck)
}

func (c *Client) EditConnectionBgpBfd(externalDeviceConn *ExternalDeviceConn) error {
	action := "edit_connection_bgp_bfd"
	data := map[string]interface{}{
		"CID":                c.CID,
		"action":             action,
		"gateway_name":       externalDeviceConn.GwName,
		"connection_name":    externalDeviceConn.ConnectionName,
		"connection_bgp_bfd": externalDeviceConn.EnableBfd,
	}
	if externalDeviceConn.EnableBfd {
		data["connection_bgp_bfd_receive_interval"] = externalDeviceConn.BgpBfdConfig.ReceiveInterval
		data["connection_bgp_bfd_transmit_interval"] = externalDeviceConn.BgpBfdConfig.TransmitInterval
		data["connection_bgp_bfd_detect_multiplier"] = externalDeviceConn.BgpBfdConfig.Multiplier
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) EditConnectionBgpMultihop(externalDeviceConn *ExternalDeviceConn) error {
	var action string
	if externalDeviceConn.EnableBgpMultihop {
		action = "enable_connection_bgp_multihop"
	} else {
		action = "disable_connection_bgp_multihop"
	}
	data := map[string]interface{}{
		"CID":             c.CID,
		"action":          action,
		"gateway_name":    externalDeviceConn.GwName,
		"connection_name": externalDeviceConn.ConnectionName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) EditBgpMd5Key(editBgpMd5Key *EditBgpMd5Key) error {
	editBgpMd5Key.CID = c.CID
	editBgpMd5Key.Action = "update_bgp_connection_md5_signature"

	return c.PostAPI(editBgpMd5Key.Action, editBgpMd5Key, BasicCheck)
}

func (c *Client) EnableJumboFrameExternalDeviceConn(externalDeviceConn *ExternalDeviceConn) error {
	params := map[string]string{
		"CID":             c.CID,
		"action":          "enable_jumbo_frame_on_connection_to_cloudn",
		"connection_name": externalDeviceConn.ConnectionName,
		"vpc_id":          externalDeviceConn.VpcID,
	}

	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "is already enabled") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	return c.PostAPI(externalDeviceConn.Action, params, checkFunc)
}

func (c *Client) DisableJumboFrameExternalDeviceConn(externalDeviceConn *ExternalDeviceConn) error {
	params := map[string]string{
		"CID":             c.CID,
		"action":          "disable_jumbo_frame_on_connection_to_cloudn",
		"connection_name": externalDeviceConn.ConnectionName,
		"vpc_id":          externalDeviceConn.VpcID,
	}

	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "is already disabled") || strings.Contains(reason, "AVXERR-SITE2CLOUD-0069") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	return c.PostAPI(externalDeviceConn.Action, params, checkFunc)
}

func (c *Client) ConnectionBGPSendCommunities(bgpSendCommunities *BgpSendCommunities) error {
	bgpSendCommunities.CID = c.CID
	bgpSendCommunities.Action = "edit_connection_bgp_send_communities"

	params := map[string]string{
		"CID":                             c.CID,
		"action":                          "edit_connection_bgp_send_communities",
		"gateway_name":                    bgpSendCommunities.GwName,
		"connection_name":                 bgpSendCommunities.ConnectionName,
		"connection_bgp_send_communities": bgpSendCommunities.ConnSendCommunities,
		"connection_bgp_send_communities_additive": fmt.Sprint(bgpSendCommunities.ConnSendAdditive),
		"connection_bgp_send_communities_block":    fmt.Sprint(bgpSendCommunities.ConnSendBlock),
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

// configureHAForTwoDevices configures HA settings when there are two external devices
func configureHAForTwoDevices(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail, localGateway *Gateway, remoteIP []string, backupBgpRemoteAsNumber int) {
	// Check if this is an edge transit gateway that supports proper HA
	isEdgeTransitGateway := localGateway != nil && localGateway.EdgeGateway && localGateway.TransitVpc == "yes"

	if isEdgeTransitGateway {
		// Edge transit gateways support true HA with separate tunnel configurations
		externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr
		externalDeviceConn.BackupLocalTunnelCidr = externalDeviceConnDetail.BackupLocalTunnelCidr
		externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr
		externalDeviceConn.BackupRemoteTunnelCidr = externalDeviceConnDetail.BackupRemoteTunnelCidr
		externalDeviceConn.BackupRemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[1]
		externalDeviceConn.HAEnabled = Enabled
		externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
	} else {
		// Non-edge gateways treat dual devices as combined configuration (no true HA)
		externalDeviceConn.LocalTunnelCidr = join2(externalDeviceConnDetail.LocalTunnelCidr, externalDeviceConnDetail.BackupLocalTunnelCidr)
		externalDeviceConn.RemoteTunnelCidr = join2(externalDeviceConnDetail.RemoteTunnelCidr, externalDeviceConnDetail.BackupRemoteTunnelCidr)
		externalDeviceConn.RemoteGatewayIP = join2(remoteIP[0], remoteIP[1])
		externalDeviceConn.HAEnabled = Disabled
	}
}

// populateBasicConnectionInfo populates basic connection information
func populateBasicConnectionInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
	if len(externalDeviceConnDetail.VpcID) != 0 {
		externalDeviceConn.VpcID = externalDeviceConnDetail.VpcID[0]
	}
	externalDeviceConn.ConnectionName = externalDeviceConnDetail.ConnectionName[0]
	externalDeviceConn.GwName = externalDeviceConnDetail.GwName
	externalDeviceConn.RemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[0]
	externalDeviceConn.BgpSendCommunities = externalDeviceConnDetail.BgpSendCommunities
}

// populateConnectionTypeInfo populates connection type and BGP information
func populateConnectionTypeInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
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
	externalDeviceConn.LocalSubnet = externalDeviceConnDetail.LocalSubnet
}

// populateAlgorithmInfo populates encryption algorithm information
func populateAlgorithmInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
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
}

// populateDirectConnectInfo populates direct connect information
func populateDirectConnectInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
	externalDeviceConn.DirectConnect = on(externalDeviceConnDetail.DirectConnect)
	externalDeviceConn.BackupDirectConnect = on(externalDeviceConnDetail.BackupDirectConnect)
}

// parseBackupBgpRemoteAsNumber parses backup BGP remote AS number
func parseBackupBgpRemoteAsNumber(backupBgpRemoteAsNumStr string) int {
	if backupBgpRemoteAsNumStr != "" {
		backupBgpRemoteAsNumberRead, _ := strconv.Atoi(backupBgpRemoteAsNumStr)
		return backupBgpRemoteAsNumberRead
	}
	return 0
}

// populateBgpSendCommunitiesInfo populates BGP send communities information
func populateBgpSendCommunitiesInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
	// get_site2cloud_conn_detail API returns one field for communities, namely, conn_bgp_send_communities
	// Example1: "conn_bgp_send_communities": "additive 444:444"
	// Example2: "conn_bgp_send_communities": "block"
	// We need to parse this field to set the BgpSendCommunities, BgpSendCommunitiesAdditive and BgpSendCommunitiesBlock fields
	if externalDeviceConnDetail.BgpSendCommunities == "" {
		return
	}
	parts := strings.Fields(externalDeviceConnDetail.BgpSendCommunities)
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case "block":
		externalDeviceConn.BgpSendCommunitiesBlock = true
	case "additive":
		externalDeviceConn.BgpSendCommunitiesAdditive = true
		externalDeviceConn.BgpSendCommunities = strings.Join(parts[1:], " ")
	}
}

// populateNonLANTunnelInfo populates information for non-LAN tunnels including HA logic
func populateNonLANTunnelInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail, localGateway *Gateway, backupBgpRemoteAsNumber int) {
	externalDeviceConn.DisableActivemesh = externalDeviceConnDetail.DisableActivemesh
	externalDeviceConn.TunnelSrcIP = externalDeviceConnDetail.TunnelSrcIP

	// extrnalDeviceConnDetail.HAEnabled returned from API indicate whether connection's local gateway has HA enabled
	// We need to put whether remote HA is enabled in externalDeviceConn.HAEnabled
	if externalDeviceConnDetail.HAEnabled == "enabled" {
		configureLocalGatewayHAEnabled(externalDeviceConn, externalDeviceConnDetail, localGateway, backupBgpRemoteAsNumber)
	} else {
		configureLocalGatewayHADisabled(externalDeviceConn, externalDeviceConnDetail, backupBgpRemoteAsNumber)
	}
}

// configureLocalGatewayHAEnabled configures HA when local gateway has HA enabled
func configureLocalGatewayHAEnabled(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail, localGateway *Gateway, backupBgpRemoteAsNumber int) {
	nTunnels := len(externalDeviceConnDetail.Tunnels)
	switch nTunnels {
	case 2:
		remoteIP := strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")
		nIPs := len(remoteIP)
		switch nIPs {
		case 2:
			if remoteIP[0] == remoteIP[1] {
				// one external device, no remote HA
				externalDeviceConn.LocalTunnelCidr = join2(externalDeviceConnDetail.LocalTunnelCidr, externalDeviceConnDetail.BackupLocalTunnelCidr)
				externalDeviceConn.RemoteTunnelCidr = join2(externalDeviceConnDetail.RemoteTunnelCidr, externalDeviceConnDetail.BackupRemoteTunnelCidr)
				externalDeviceConn.HAEnabled = Disabled
			} else {
				// two external devices, remote has HA
				// activemesh is disabled, 2 straight tunnels only
				configureHAForTwoDevices(externalDeviceConn, externalDeviceConnDetail, localGateway, remoteIP, backupBgpRemoteAsNumber)
			}
		case 4:
			if remoteIP[0] == remoteIP[2] && remoteIP[1] == remoteIP[3] {
				externalDeviceConn.LocalTunnelCidr = join2(externalDeviceConnDetail.LocalTunnelCidr, externalDeviceConnDetail.BackupLocalTunnelCidr)
				externalDeviceConn.RemoteTunnelCidr = join2(externalDeviceConnDetail.RemoteTunnelCidr, externalDeviceConnDetail.BackupRemoteTunnelCidr)
				externalDeviceConn.RemoteGatewayIP = join2(remoteIP[0], remoteIP[1])
				externalDeviceConn.HAEnabled = Disabled
			}
		}
	case 4:
		// activemesh is enabled, 4 tunnels, 2 straight tunnels and 2 crossing tunnels
		externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr
		externalDeviceConn.BackupLocalTunnelCidr = externalDeviceConnDetail.BackupLocalTunnelCidr
		externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr
		externalDeviceConn.BackupRemoteTunnelCidr = externalDeviceConnDetail.BackupRemoteTunnelCidr
		externalDeviceConn.BackupRemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[1]
		externalDeviceConn.HAEnabled = Enabled
		externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
	}
}

// configureLocalGatewayHADisabled configures HA when local gateway has HA disabled
func configureLocalGatewayHADisabled(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail, backupBgpRemoteAsNumber int) {
	// local gateway no HA
	externalDeviceConn.LocalTunnelCidr = externalDeviceConnDetail.LocalTunnelCidr
	externalDeviceConn.RemoteTunnelCidr = externalDeviceConnDetail.RemoteTunnelCidr
	if len(externalDeviceConnDetail.Tunnels) == 2 {
		// two external devices, remote has HA
		externalDeviceConn.BackupLocalTunnelCidr = externalDeviceConnDetail.BackupLocalTunnelCidr
		externalDeviceConn.BackupRemoteTunnelCidr = externalDeviceConnDetail.BackupRemoteTunnelCidr
		externalDeviceConn.BackupRemoteGatewayIP = strings.Split(externalDeviceConnDetail.RemoteGatewayIP, ",")[1]
		externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
		externalDeviceConn.HAEnabled = Enabled
	} else {
		// one external device, no remote HA
		externalDeviceConn.HAEnabled = Disabled
	}
}

// populateLANTunnelInfo populates information for LAN tunnels
func populateLANTunnelInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail, backupBgpRemoteAsNumber int) {
	externalDeviceConn.EnableBgpLanActiveMesh = externalDeviceConnDetail.EnableBgpLanActiveMesh
	if len(externalDeviceConnDetail.Tunnels) == 2 || len(externalDeviceConnDetail.Tunnels) == 4 {
		externalDeviceConn.HAEnabled = Enabled
		externalDeviceConn.BackupBgpRemoteAsNum = backupBgpRemoteAsNumber
		externalDeviceConn.BackupRemoteLanIP = externalDeviceConnDetail.BackupRemoteLanIP
		externalDeviceConn.BackupLocalLanIP = externalDeviceConnDetail.BackupLocalLanIP
	} else {
		externalDeviceConn.HAEnabled = Disabled
	}
	externalDeviceConn.RemoteLanIP = externalDeviceConnDetail.RemoteLanIP
	externalDeviceConn.LocalLanIP = externalDeviceConnDetail.LocalLanIP
}

// populateAdditionalConnectionInfo populates additional connection information
func populateAdditionalConnectionInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
	externalDeviceConn.EnableEdgeSegmentation = on(externalDeviceConnDetail.EnableEdgeSegmentation)
	externalDeviceConn.ManualBGPCidrs = externalDeviceConnDetail.ManualBGPCidrs

	externalDeviceConn.EnableIkev2 = on(externalDeviceConnDetail.IkeVer == "2")
	externalDeviceConn.EventTriggeredHA = externalDeviceConnDetail.EventTriggeredHA == "enabled"
	externalDeviceConn.EnableJumboFrame = externalDeviceConnDetail.EnableJumboFrame
	externalDeviceConn.PeerVnetID = externalDeviceConnDetail.PeerVnetID
	externalDeviceConn.Phase1RemoteIdentifier = externalDeviceConnDetail.Phase1RemoteIdentifier
	externalDeviceConn.PrependAsPath = externalDeviceConnDetail.PrependAsPath
	externalDeviceConn.EnableEdgeUnderlay = externalDeviceConnDetail.WanUnderlay
	externalDeviceConn.RemoteCloudType = externalDeviceConnDetail.RemoteCloudType
	externalDeviceConn.Phase1LocalIdentifier = externalDeviceConnDetail.Phase1LocalIdentifier
	externalDeviceConn.EnableBfd = externalDeviceConnDetail.EnableBfd
	if externalDeviceConn.EnableBfd {
		externalDeviceConn.BgpBfdConfig.TransmitInterval = externalDeviceConnDetail.BgpBfdConfig["tx_interval"]
		externalDeviceConn.BgpBfdConfig.ReceiveInterval = externalDeviceConnDetail.BgpBfdConfig["rx_interval"]
		externalDeviceConn.BgpBfdConfig.Multiplier = externalDeviceConnDetail.BgpBfdConfig["multiplier"]
	}
	externalDeviceConn.EnableBgpMultihop = externalDeviceConnDetail.EnableBgpMultihop
	externalDeviceConn.ProxyIdEnabled = externalDeviceConnDetail.ProxyIdEnabled
}

// populateIPv6ConnectionInfo populates IPv6 related connection information
func populateIPv6ConnectionInfo(externalDeviceConn *ExternalDeviceConn, externalDeviceConnDetail EditExternalDeviceConnDetail) {
	externalDeviceConn.EnableIpv6 = externalDeviceConnDetail.EnableIpv6
	if !externalDeviceConn.EnableIpv6 {
		return
	}
	if externalDeviceConn.TunnelProtocol == "LAN" {
		externalDeviceConn.RemoteLanIPv6 = externalDeviceConnDetail.RemoteLanIPv6
		if externalDeviceConnDetail.HAEnabled == "enabled" {
			externalDeviceConn.BackupRemoteLanIPv6 = externalDeviceConnDetail.BackupRemoteLanIPv6
		}
	} else {
		// for ipsec and gre tunnels
		externalDeviceConn.ExternalDeviceIPv6 = externalDeviceConnDetail.ExternalDeviceIPv6
		if externalDeviceConnDetail.HAEnabled == "enabled" {
			externalDeviceConn.ExternalDeviceBackupIPv6 = externalDeviceConnDetail.ExternalDeviceBackupIPv6
		}
	}
}
