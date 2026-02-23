package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeExternalDeviceConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeExternalDeviceConnCreate,
		Read:   resourceAviatrixSpokeExternalDeviceConnRead,
		Update: resourceAviatrixSpokeExternalDeviceConnUpdate,
		Delete: resourceAviatrixSpokeExternalDeviceConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC where the BGP Spoke Gateway is located.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the spoke external device connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the BGP Spoke Gateway.",
			},
			"remote_gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote Gateway IP.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldIpList := strings.Split(old, ",")
					newIpList := strings.Split(new, ",")
					if len(oldIpList) == len(newIpList) {
						for i := range oldIpList {
							if strings.TrimSpace(oldIpList[i]) != strings.TrimSpace(newIpList[i]) {
								return false
							}
						}
						return true
					}
					return false
				},
			},
			"connection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "bgp",
				ForceNew:    true,
				Description: "Connection type. Valid values: 'bgp', 'static'. Default value: 'bgp'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := mustString(val)
					if v != "bgp" && v != "static" {
						errs = append(errs, fmt.Errorf("%q must be either 'bgp' or 'static', got: %s", key, val))
					}
					return
				},
			},
			"tunnel_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPsec",
				ForceNew:     true,
				Description:  "Tunnel Protocol. Valid values: 'IPsec', 'GRE' or 'LAN'. Default value: 'IPsec'. Case insensitive.",
				ValidateFunc: validation.StringInSlice([]string{"IPsec", "GRE", "LAN"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"enable_bgp_lan_activemesh": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "Switch to enable BGP LAN ActiveMesh. Only valid for GCP and Azure with Remote Gateway HA enabled. " +
					"Requires Azure Remote Gateway insane mode enabled. Valid values: true, false. Default: false. " +
					"Available as of provider version R3.0.2+.",
			},
			"bgp_local_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "BGP local ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"bgp_remote_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "BGP remote ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"remote_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a 'static' type connection.",
			},
			"local_subnet": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Local CIDRs joined as a string with ','. Optional for a 'static' type connection with proxy ID enabled",
			},
			"proxy_id_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable proxy ID for spoke static route based external device connection.",
			},
			"direct_connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set true for private network infrastructure.",
			},
			"pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "If left blank, the pre-shared key will be auto generated.",
			},
			"local_tunnel_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceOnlyInString,
				Description:      "Source CIDR for the tunnel from the Aviatrix spoke gateway.",
			},
			"remote_tunnel_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceOnlyInString,
				Description:      "Destination CIDR for the tunnel to the external device.",
			},
			"custom_algorithms": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Switch to enable custom/non-default algorithms for IPSec Authentication/Encryption.",
			},
			"phase_1_authentication": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Phase one Authentication. Valid values: 'SHA-1', 'SHA-256', 'SHA-384' and 'SHA-512'.",
				ValidateFunc: validation.StringInSlice([]string{
					"SHA-1", "SHA-256", "SHA-384", "SHA-512",
				}, false),
			},
			"phase_2_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase two Authentication. Valid values: 'NO-AUTH', 'HMAC-SHA-1', 'HMAC-SHA-256', " +
					"'HMAC-SHA-384' and 'HMAC-SHA-512'.",
				ValidateFunc: validation.StringInSlice([]string{
					"NO-AUTH", "HMAC-SHA-1", "HMAC-SHA-256", "HMAC-SHA-384", "HMAC-SHA-512",
				}, false),
			},
			"phase_1_dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Phase one DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17', '18', '19', '20' and '21'.",
				ValidateFunc: validation.StringInSlice([]string{
					"1", "2", "5", "14", "15", "16", "17", "18", "19", "20", "21",
				}, false),
			},
			"phase_2_dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Phase two DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17', '18', '19', '20' and '21'.",
				ValidateFunc: validation.StringInSlice([]string{
					"1", "2", "5", "14", "15", "16", "17", "18", "19", "20", "21",
				}, false),
			},
			"phase_1_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and 'AES-256-CBC', " +
					"'AES-128-GCM-64', 'AES-128-GCM-96', 'AES-128-GCM-128', 'AES-256-GCM-64', 'AES-256-GCM-96', and 'AES-256-GCM-128'.",
				ValidateFunc: validation.StringInSlice([]string{
					"3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "AES-128-GCM-64", "AES-128-GCM-96",
					"AES-128-GCM-128", "AES-256-GCM-64", "AES-256-GCM-96", "AES-256-GCM-128",
				}, false),
			},
			"phase_2_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', 'AES-256-CBC', " +
					"'AES-128-GCM-64', 'AES-128-GCM-96', 'AES-128-GCM-128', 'AES-256-GCM-64', 'AES-256-GCM-96', 'AES-256-GCM-128', and 'NULL-ENCR'.",
				ValidateFunc: validation.StringInSlice([]string{
					"3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "AES-128-GCM-64", "AES-128-GCM-96",
					"AES-128-GCM-128", "AES-256-GCM-64", "AES-256-GCM-96", "AES-256-GCM-128", "NULL-ENCR",
				}, false),
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set as true if there are two external devices.",
			},
			"backup_remote_gateway_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				Description:  "Backup remote gateway IP.",
				ValidateFunc: validation.IsIPv4Address,
			},
			"backup_bgp_remote_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				Description:  "Backup BGP remote ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"backup_pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				ForceNew:    true,
				Description: "Backup pre shared key.",
			},
			"backup_local_tunnel_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceOnlyInString,
				Description:      "Source CIDR for the tunnel from the backup Aviatrix spoke gateway.",
			},
			"backup_remote_tunnel_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceOnlyInString,
				Description:      "Destination CIDR for the tunnel to the backup external device.",
			},
			"backup_direct_connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Backup direct connect for backup external device.",
			},
			"enable_learned_cidrs_approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable learned CIDR approval for the connection. Only valid with 'connection_type' = 'bgp'." +
					" Requires the spoke_gateway's 'learned_cidrs_approval_mode' attribute be set to 'connection'. " +
					"Valid values: true, false. Default value: false.",
			},
			"enable_ikev2": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set as true if use IKEv2.",
			},
			"enable_event_triggered_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Event Triggered HA.",
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Jumbo Frame for the transit external device connection. Valid values: true, false.",
			},
			"manual_bgp_advertised_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional: true,
				Description: "Configure manual BGP advertised CIDRs for this connection. Only valid with 'connection_type'" +
					" = 'bgp'.",
				DiffSuppressFunc: func(_, oldStr, newStr string, _ *schema.ResourceData) bool {
					// Suppress diff if old is null ("" or "<nil>") and new is an empty set/list ("[]")
					return (oldStr == "" || oldStr == "<nil>") && (newStr == "[]" || newStr == "")
				},
			},
			"remote_vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Name of the remote VPC for a LAN BGP connection. Only valid when 'connection_type' = 'bgp' " +
					"and tunnel_protocol' = 'LAN' with an Azure spoke gateway. Must be in the form " +
					"\"<VNET-name>:<resource-group-name>\". Available as of provider version R3.0.2+.",
			},
			"remote_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote LAN IP.",
			},
			"backup_remote_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup Remote LAN IP.",
			},
			"phase1_local_identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "public_ip",
				ValidateFunc: validation.StringInSlice([]string{"public_ip", "private_ip"}, false),
				Description:  "By default, gateway’s public IP is configured as the Local Identifier.",
			},
			"phase1_remote_identifier": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.StringCanBeEmptyButCannotBeWhiteSpace,
				},
				DiffSuppressFunc: goaviatrix.TransitExternalDeviceConnPh1RemoteIdDiffSuppressFunc,
				Description:      "List of phase 1 remote identifier of the IPsec tunnel. This can be configured as a list of any string, including empty string.",
			},
			"prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Connection AS Path Prepend customized by specifying AS PATH for a BGP connection.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"bgp_md5_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "BGP MD5 authentication key.",
			},
			"backup_bgp_md5_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Backup BGP MD5 authentication key.",
			},
			"approved_cidrs": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Set of approved cidrs. Requires 'enable_learned_cidrs_approval' to be true. Type: Set(String).",
			},
			"local_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Local LAN IP.",
			},
			"backup_local_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup Local LAN IP.",
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable BGP BFD connection.",
			},
			"bgp_bfd": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "BGP BFD configuration details applied to a BGP session.",
				MaxItems:         1,
				DiffSuppressFunc: goaviatrix.ExternalDeviceConnBgpBfdDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transmit_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "BFD transmit interval in milliseconds.",
							ValidateFunc: validation.IntBetween(10, 60000),
							Default:      defaultBfdTransmitInterval,
						},
						"receive_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "BFD receive interval in milliseconds.",
							ValidateFunc: validation.IntBetween(10, 60000),
							Default:      defaultBfdReceiveInterval,
						},
						"multiplier": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "BFD detection multiplier.",
							ValidateFunc: validation.IntBetween(2, 255),
							Default:      defaultBfdMultiplier,
						},
					},
				},
			},
			"enable_bgp_multihop": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable multihop on BGP connection.",
			},
			"connection_bgp_send_communities": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Connection based additional BGP communities to be sent. E.g. 111:111, 444:444",
			},
			"connection_bgp_send_communities_additive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do additive operation instead of replacement operation",
			},
			"connection_bgp_send_communities_block": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Block advertisement of any BGP communities on this connection",
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable IPv6 on this connection",
			},
			"external_device_ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "External device IPv6 address.",
			},
			"external_device_backup_ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup External device IPv6 address.",
			},
			"remote_lan_ipv6_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Remote LAN IPv6 address.",
			},
			"backup_remote_lan_ipv6_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup Remote LAN IPv6 address.",
			},
		},
	}
}

func resourceAviatrixSpokeExternalDeviceConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	var bgpSendCommunities *goaviatrix.BgpSendCommunities

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:                    getString(d, "vpc_id"),
		ConnectionName:           getString(d, "connection_name"),
		GwName:                   getString(d, "gw_name"),
		ConnectionType:           getString(d, "connection_type"),
		RemoteGatewayIP:          getString(d, "remote_gateway_ip"),
		RemoteSubnet:             getString(d, "remote_subnet"),
		LocalSubnet:              getString(d, "local_subnet"),
		PreSharedKey:             getString(d, "pre_shared_key"),
		LocalTunnelCidr:          getString(d, "local_tunnel_cidr"),
		RemoteTunnelCidr:         getString(d, "remote_tunnel_cidr"),
		Phase1Auth:               getString(d, "phase_1_authentication"),
		Phase1DhGroups:           getString(d, "phase_1_dh_groups"),
		Phase1Encryption:         getString(d, "phase_1_encryption"),
		Phase2Auth:               getString(d, "phase_2_authentication"),
		Phase2DhGroups:           getString(d, "phase_2_dh_groups"),
		Phase2Encryption:         getString(d, "phase_2_encryption"),
		BackupRemoteGatewayIP:    getString(d, "backup_remote_gateway_ip"),
		BackupPreSharedKey:       getString(d, "backup_pre_shared_key"),
		PeerVnetID:               getString(d, "remote_vpc_name"),
		RemoteLanIP:              getString(d, "remote_lan_ip"),
		LocalLanIP:               getString(d, "local_lan_ip"),
		BackupRemoteLanIP:        getString(d, "backup_remote_lan_ip"),
		BackupLocalLanIP:         getString(d, "backup_local_lan_ip"),
		BackupLocalTunnelCidr:    getString(d, "backup_local_tunnel_cidr"),
		BackupRemoteTunnelCidr:   getString(d, "backup_remote_tunnel_cidr"),
		BgpMd5Key:                getString(d, "bgp_md5_key"),
		BackupBgpMd5Key:          getString(d, "backup_bgp_md5_key"),
		EnableJumboFrame:         getBool(d, "enable_jumbo_frame"),
		EnableBgpMultihop:        getBool(d, "enable_bgp_multihop"),
		ProxyIdEnabled:           getBool(d, "proxy_id_enabled"),
		EnableIpv6:               getBool(d, "enable_ipv6"),
		ExternalDeviceIPv6:       getString(d, "external_device_ipv6"),
		ExternalDeviceBackupIPv6: getString(d, "external_device_backup_ipv6"),
		RemoteLanIPv6:            getString(d, "remote_lan_ipv6_ip"),
		BackupRemoteLanIPv6:      getString(d, "backup_remote_lan_ipv6_ip"),
	}

	sendComm := getString(d, "connection_bgp_send_communities")

	blockComm := getBool(d, "connection_bgp_send_communities_block")

	setPerConnCommunity := false
	if sendComm != "" || blockComm {
		connName := getString(d, "connection_name")
		setPerConnCommunity = true

		gwName := getString(d, "gw_name")

		sendAdditive := getBool(d, "connection_bgp_send_communities_additive")

		bgpSendCommunities = &goaviatrix.BgpSendCommunities{
			ConnectionName:      connName,
			GwName:              gwName,
			ConnSendCommunities: sendComm,
			ConnSendAdditive:    sendAdditive,
			ConnSendBlock:       blockComm,
		}
	}

	tunnelProtocol := strings.ToUpper(getString(d, "tunnel_protocol"))
	if tunnelProtocol == "IPSEC" {
		externalDeviceConn.TunnelProtocol = "IPsec"
	} else {
		externalDeviceConn.TunnelProtocol = tunnelProtocol
	}

	if (externalDeviceConn.RemoteGatewayIP != "" ||
		externalDeviceConn.LocalTunnelCidr != "" ||
		externalDeviceConn.BackupRemoteGatewayIP != "" ||
		externalDeviceConn.BackupLocalTunnelCidr != "") && externalDeviceConn.TunnelProtocol == "LAN" {
		return fmt.Errorf("'remote_gateway_ip', 'local_tunnel_cidr', 'backup_remote_gateway_ip' and 'backup_local_tunnel_cidr' " +
			"cannot be set with 'tunnel_protocol' = 'LAN'. Please use the appropriate LAN attributes instead")
	}
	if (externalDeviceConn.RemoteLanIP != "" ||
		externalDeviceConn.LocalLanIP != "" ||
		externalDeviceConn.BackupRemoteLanIP != "" ||
		externalDeviceConn.BackupLocalLanIP != "") && externalDeviceConn.TunnelProtocol != "LAN" {
		return fmt.Errorf("'remote_lan_ip', 'local_lan_ip', 'backup_remote_lan_ip' and 'backup_local_lan_ip' " +
			"can only be set with 'tunnel_protocol' = 'LAN'")
	}
	if externalDeviceConn.RemoteLanIP == "" && externalDeviceConn.TunnelProtocol == "LAN" {
		return fmt.Errorf("'remote_lan_ip' is required when 'tunnel_protocol' = 'LAN'")
	}
	if externalDeviceConn.EnableIpv6 && externalDeviceConn.RemoteLanIPv6 == "" && externalDeviceConn.TunnelProtocol == "LAN" {
		return fmt.Errorf("'remote_lan_ipv6_ip' is required when 'tunnel_protocol' = 'LAN'")
	}

	if externalDeviceConn.RemoteGatewayIP == "" && externalDeviceConn.TunnelProtocol != "LAN" {
		return fmt.Errorf("'remote_gateway_ip' is required when 'tunnel_protocol' != 'LAN'")
	}

	bgpLocalAsNum, err := strconv.Atoi(getString(d, "bgp_local_as_num"))
	if err == nil {
		externalDeviceConn.BgpLocalAsNum = bgpLocalAsNum
	}
	bgpRemoteAsNum, err := strconv.Atoi(getString(d, "bgp_remote_as_num"))
	if err == nil {
		externalDeviceConn.BgpRemoteAsNum = bgpRemoteAsNum
	}
	backupBgpLocalAsNum, err := strconv.Atoi(getString(d, "backup_bgp_remote_as_num"))
	if err == nil {
		externalDeviceConn.BackupBgpRemoteAsNum = backupBgpLocalAsNum
	}

	directConnect := getBool(d, "direct_connect")
	if directConnect {
		externalDeviceConn.DirectConnect = "true"
	}

	haEnabled := getBool(d, "ha_enabled")
	if haEnabled {
		externalDeviceConn.HAEnabled = "true"
	}

	backupDirectConnect := getBool(d, "backup_direct_connect")
	if backupDirectConnect {
		externalDeviceConn.BackupDirectConnect = "true"
	}

	if externalDeviceConn.ConnectionType == "bgp" && externalDeviceConn.RemoteSubnet != "" {
		return fmt.Errorf("'remote_subnet' is needed for connection type of 'static' not 'bgp'")
	} else if externalDeviceConn.ConnectionType == "static" && (externalDeviceConn.BgpLocalAsNum != 0 || externalDeviceConn.BgpRemoteAsNum != 0) {
		return fmt.Errorf("'bgp_local_as_num' and 'bgp_remote_as_num' are needed for connection type of 'bgp' not 'static'")
	}

	customAlgorithms := getBool(d, "custom_algorithms")
	if customAlgorithms {
		if externalDeviceConn.Phase1Auth == "" ||
			externalDeviceConn.Phase2Auth == "" ||
			externalDeviceConn.Phase1DhGroups == "" ||
			externalDeviceConn.Phase2DhGroups == "" ||
			externalDeviceConn.Phase1Encryption == "" ||
			externalDeviceConn.Phase2Encryption == "" {
			return fmt.Errorf("custom_algorithms is enabled, please set all of the algorithm parameters")
		} else if externalDeviceConn.Phase1Auth == goaviatrix.Phase1AuthDefault &&
			externalDeviceConn.Phase2Auth == goaviatrix.Phase2AuthDefault &&
			externalDeviceConn.Phase1DhGroups == goaviatrix.Phase1DhGroupDefault &&
			externalDeviceConn.Phase2DhGroups == goaviatrix.Phase2DhGroupDefault &&
			externalDeviceConn.Phase1Encryption == goaviatrix.Phase1EncryptionDefault &&
			externalDeviceConn.Phase2Encryption == goaviatrix.Phase2EncryptionDefault {
			return fmt.Errorf("custom_algorithms is enabled, cannot use default values for " +
				"all six algorithm parameters. Please change the value of at least one of the six algorithm parameters")
		}
	} else {
		if externalDeviceConn.Phase1Auth != "" || externalDeviceConn.Phase1DhGroups != "" ||
			externalDeviceConn.Phase1Encryption != "" || externalDeviceConn.Phase2Auth != "" ||
			externalDeviceConn.Phase2DhGroups != "" || externalDeviceConn.Phase2Encryption != "" {
			return fmt.Errorf("custom_algorithms is not enabled, all algorithm fields should be left empty")
		}
	}

	if haEnabled {
		if externalDeviceConn.TunnelProtocol == "LAN" {
			if externalDeviceConn.BackupRemoteLanIP == "" {
				return fmt.Errorf("ha is enabled and 'tunnel_protocol' = 'LAN', please specify 'backup_remote_lan_ip'")
			}
			if externalDeviceConn.EnableIpv6 && externalDeviceConn.BackupRemoteLanIPv6 == "" {
				return fmt.Errorf("ha is enabled, 'tunnel_protocol' = 'LAN' and 'enable_ipv6' is true, please specify 'backup_remote_lan_ipv6_ip'")
			}
		} else {
			if externalDeviceConn.BackupRemoteGatewayIP == "" {
				return fmt.Errorf("ha is enabled, please specify 'backup_remote_gateway_ip'")
			}
			remoteIP := strings.Split(externalDeviceConn.RemoteGatewayIP, ",")
			if len(remoteIP) > 1 {
				return fmt.Errorf("expected 'remote_gateway_ip' to contain only one valid IPv4 address, got: %s", externalDeviceConn.RemoteGatewayIP)
			}
			ip := net.ParseIP(externalDeviceConn.RemoteGatewayIP)
			if four := ip.To4(); four == nil {
				return fmt.Errorf("expected 'remote_gateway_ip' to contain a valid IPv4 address, got: %s", externalDeviceConn.RemoteGatewayIP)
			}
			if externalDeviceConn.BackupRemoteGatewayIP == externalDeviceConn.RemoteGatewayIP {
				return fmt.Errorf("expected 'backup_remote_gateway_ip' to contain a different valid IPv4 address than 'remote_gateway_ip'")
			}
			if externalDeviceConn.EnableIpv6 && externalDeviceConn.ExternalDeviceIPv6 != "" && externalDeviceConn.ExternalDeviceBackupIPv6 == "" {
				return fmt.Errorf("ha is enabled and 'enable_ipv6' is true, please specify 'external_device_backup_ipv6'")
			}
		}
		if externalDeviceConn.BackupBgpRemoteAsNum == 0 && externalDeviceConn.ConnectionType == "bgp" {
			return fmt.Errorf("ha is enabled, and 'connection_type' is 'bgp', please specify 'backup_bgp_remote_as_num'")
		}
		if externalDeviceConn.BackupRemoteGatewayIP != "" {
			externalDeviceConn.RemoteGatewayIP = externalDeviceConn.RemoteGatewayIP + "," + externalDeviceConn.BackupRemoteGatewayIP
			externalDeviceConn.BackupRemoteGatewayIP = ""
		}
	} else {
		if backupDirectConnect {
			return fmt.Errorf("ha is not enabled, please set 'backup_direct_connect' to false")
		}
		if externalDeviceConn.BackupPreSharedKey != "" || externalDeviceConn.BackupLocalTunnelCidr != "" ||
			externalDeviceConn.BackupRemoteTunnelCidr != "" || externalDeviceConn.BackupRemoteGatewayIP != "" ||
			externalDeviceConn.BackupRemoteLanIP != "" || externalDeviceConn.BackupLocalLanIP != "" {
			return fmt.Errorf("ha is not enabled, please set 'backup_pre_shared_key', 'backup_local_tunnel_cidr', " +
				"'backup_remote_gateway_ip', 'backup_remote_tunnel_cidr', 'backup_remote_lan_ip' and 'backup_local_lan_ip' to empty")
		}
		if externalDeviceConn.BackupBgpRemoteAsNum != 0 && externalDeviceConn.ConnectionType == "bgp" {
			return fmt.Errorf("ha is not enabled, and 'connection_type' is 'bgp', please specify 'backup_bgp_remote_as_num' to empty")
		}
	}

	enableLearnedCIDRApproval := getBool(d, "enable_learned_cidrs_approval")
	if externalDeviceConn.ConnectionType != "bgp" && enableLearnedCIDRApproval {
		return fmt.Errorf("'connection_type' must be 'bgp' if 'enable_learned_cidrs_approval' is set to true")
	}
	manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
	if externalDeviceConn.ConnectionType != "bgp" && len(manualBGPCidrs) != 0 {
		return fmt.Errorf("'connection_type' must be 'bgp' if 'manual_bgp_advertised_cidrs' is not empty")
	}

	approvedCidrs := getStringSet(d, "approved_cidrs")
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 {
		return fmt.Errorf("creating spoke external device conn: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	enableIkev2 := getBool(d, "enable_ikev2")
	if enableIkev2 {
		externalDeviceConn.EnableIkev2 = "true"
	}

	if externalDeviceConn.ConnectionType != "bgp" && externalDeviceConn.TunnelProtocol != "IPsec" {
		return fmt.Errorf("'tunnel_protocol' can not be set unless 'connection_type' is 'bgp'")
	}
	greOrLan := externalDeviceConn.TunnelProtocol == "GRE" || externalDeviceConn.TunnelProtocol == "LAN"
	if greOrLan && customAlgorithms {
		return fmt.Errorf("custom algorithm parameters are not valid with 'tunnel_protocol' = GRE or LAN")
	}
	if greOrLan && enableIkev2 {
		return fmt.Errorf("enable_ikev2 is not supported with 'tunnel_protocol' = GRE or LAN")
	}
	if greOrLan && externalDeviceConn.PreSharedKey != "" {
		return fmt.Errorf("'pre_shared_key' is not valid with 'tunnel_protocol' = GRE or LAN")
	}
	if externalDeviceConn.PeerVnetID != "" && (externalDeviceConn.ConnectionType != "bgp" || externalDeviceConn.TunnelProtocol != "LAN") {
		return fmt.Errorf("'remote_vpc_name' is only valid for 'connection_type' = 'bgp' and 'tunnel_protocol' = 'LAN'")
	}
	if externalDeviceConn.TunnelProtocol == "LAN" {
		if externalDeviceConn.DirectConnect == "true" || externalDeviceConn.BackupDirectConnect == "true" {
			return fmt.Errorf("enabling 'direct_connect' or 'backup_direct_connect' is not allowed for BGP over LAN connections")
		}
	}

	phase1RemoteIdentifier := getList(d, "phase1_remote_identifier")
	ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)
	if haEnabled && len(phase1RemoteIdentifier) != 0 && len(phase1RemoteIdentifier) != 2 {
		return fmt.Errorf("please either set two phase 1 remote IDs or none, when HA is enabled")
	} else if !haEnabled && len(phase1RemoteIdentifier) > 1 {
		return fmt.Errorf("please either set one phase 1 remote ID or none, when HA is disabled")
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		if externalDeviceConn.ConnectionType != "bgp" {
			return fmt.Errorf("'prepend_as_path' only supports 'bgp' connection. Please update 'connection_type' to 'bgp'")
		}
	}

	if getBool(d, "enable_bgp_lan_activemesh") {
		if externalDeviceConn.ConnectionType != "bgp" || externalDeviceConn.TunnelProtocol != "LAN" {
			return fmt.Errorf("'enable_bgp_lan_activemesh' only supports 'bgp' connection with 'LAN' tunnel protocol")
		}
		if externalDeviceConn.HAEnabled != "true" {
			return fmt.Errorf("'enable_bgp_lan_activemesh' can only be enabled with Remote Gateway HA enabled")
		}
		externalDeviceConn.EnableBgpLanActiveMesh = true
	}

	if externalDeviceConn.BgpMd5Key != "" || externalDeviceConn.BackupBgpMd5Key != "" {
		if externalDeviceConn.ConnectionType != "bgp" {
			return fmt.Errorf("BGP MD5 authentication key is only supported for BGP connection")
		}
		if externalDeviceConn.BackupBgpMd5Key != "" && !haEnabled {
			return fmt.Errorf("couldn't configure backup BGP MD5 authentication key since HA is not enabled for BGP connection: %s", externalDeviceConn.ConnectionName)
		}

		if externalDeviceConn.BgpMd5Key != "" {
			md5KeyList := strings.Split(externalDeviceConn.BgpMd5Key, ",")
			var bgpRemoteIp []string
			if strings.ToUpper(externalDeviceConn.TunnelProtocol) == "LAN" {
				bgpRemoteIp = strings.Split(externalDeviceConn.RemoteLanIP, ",")
			} else {
				bgpRemoteIp = strings.Split(externalDeviceConn.RemoteTunnelCidr, ",")
			}
			if len(md5KeyList) != len(bgpRemoteIp) {
				return fmt.Errorf("can't apply BGP MD5 authentication key since it is not set correctly for BGP connection: %s", externalDeviceConn.ConnectionName)
			}
		}

		if externalDeviceConn.BackupBgpMd5Key != "" {
			backupMd5KeyList := strings.Split(externalDeviceConn.BackupBgpMd5Key, ",")
			var backupBgpRemoteIp []string
			if strings.ToUpper(externalDeviceConn.TunnelProtocol) == "LAN" {
				backupBgpRemoteIp = strings.Split(externalDeviceConn.BackupRemoteLanIP, ",")
			} else {
				backupBgpRemoteIp = strings.Split(externalDeviceConn.BackupRemoteTunnelCidr, ",")
			}
			if len(backupMd5KeyList) != len(backupBgpRemoteIp) {
				return fmt.Errorf("can't apply Backup BGP MD5 authentication key since it is not set correctly for BGP connection: %s", externalDeviceConn.ConnectionName)
			}
		}
	}

	enableJumboFrame := getBool(d, "enable_jumbo_frame")
	if enableJumboFrame {
		if externalDeviceConn.ConnectionType != "bgp" {
			return fmt.Errorf("jumbo frame is only supported on bgp connection")
		}
	}

	if externalDeviceConn.PreSharedKey != "" {
		externalDeviceConn.AuthType = "psk"
	}

	if !externalDeviceConn.EnableBgpMultihop && externalDeviceConn.ConnectionType != "bgp" {
		return fmt.Errorf("multihop can only be configured for BGP connections")
	}

	d.SetId(externalDeviceConn.ConnectionName + "~" + externalDeviceConn.VpcID)
	flag := false
	defer func() { _ = resourceAviatrixSpokeExternalDeviceConnReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err = client.CreateExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix external device connection: %w", err)
	}

	if setPerConnCommunity {
		err = client.ConnectionBGPSendCommunities(bgpSendCommunities)
	}
	if err != nil {
		return fmt.Errorf("failed to set/block BGP communities for connection %s: %w", externalDeviceConn.ConnectionName, err)
	}

	enableBFD := getBool(d, "enable_bfd")

	if enableBFD && externalDeviceConn.ConnectionType != "bgp" {
		return fmt.Errorf("BFD is only supported for BGP connection type")
	}
	externalDeviceConn.EnableBfd = enableBFD
	bgp_bfd := getList(d, "bgp_bfd")

	// set the bgp bfd config details only if the user has enabled BFD
	if enableBFD {
		// set bgp bfd using the config details provided by the user
		if len(bgp_bfd) > 0 {
			for _, bfd0 := range bgp_bfd {
				bfd1, ok := bfd0.(map[string]interface{})
				if !ok {
					return fmt.Errorf("expected bgp_bfd to be a map, but got %T", bfd0)
				}
				externalDeviceConn.BgpBfdConfig = goaviatrix.CreateBgpBfdConfig(bfd1)
			}
		} else {
			// set the bgp bfd config using the default values
			externalDeviceConn.BgpBfdConfig = defaultBfdConfig
		}
		err := client.EditConnectionBgpBfd(externalDeviceConn)
		if err != nil {
			return fmt.Errorf("could not update BGP BFD config: %w", err)
		}
	} else {
		// if BFD is disabled and BGP BFD config is provided then throw an error
		if len(bgp_bfd) > 0 {
			return fmt.Errorf("bgp_bfd config can't be set when BFD is disabled")
		}
	}

	if getBool(d, "enable_event_triggered_ha") {
		if err := client.EnableSite2CloudEventTriggeredHA(externalDeviceConn.VpcID, externalDeviceConn.ConnectionName); err != nil {
			return fmt.Errorf("could not enable event triggered HA for external device conn after create: %w", err)
		}
	}

	if enableJumboFrame {
		if err := client.EnableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
			return fmt.Errorf("could not enable jumbo frame for external device conn: %v after create: %w", externalDeviceConn.ConnectionName, err)
		}
	} else {
		if externalDeviceConn.ConnectionType == "bgp" {
			if err := client.DisableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return fmt.Errorf("could not disable jumbo frame for external device conn: %v after create: %w", externalDeviceConn.ConnectionName, err)
			}
		}
	}

	if enableLearnedCIDRApproval {
		err = client.EnableSpokeConnectionLearnedCIDRApproval(externalDeviceConn.GwName, externalDeviceConn.ConnectionName)
		if err != nil {
			return fmt.Errorf("could not enable learned cidr approval: %w", err)
		}
		if len(approvedCidrs) > 0 {
			err = client.UpdateSpokeConnectionPendingApprovedCidrs(externalDeviceConn.GwName, externalDeviceConn.ConnectionName, approvedCidrs)
			if err != nil {
				return fmt.Errorf("could not update spoke external device conn approved cidrs after creation: %w", err)
			}
		}
	}

	if len(manualBGPCidrs) > 0 {
		err = client.EditSpokeConnectionBGPManualAdvertiseCIDRs(externalDeviceConn.GwName, externalDeviceConn.ConnectionName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual advertised BGP cidrs: %w", err)
		}
	}

	if len(phase1RemoteIdentifier) == 1 {
		var ph1RemoteId string

		if phase1RemoteIdentifier[0] == nil {
			ph1RemoteId = "\"\""
		} else {
			ph1RemoteId = mustString(phase1RemoteIdentifier[0])
		}

		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:                 externalDeviceConn.GwName,
			VpcID:                  externalDeviceConn.VpcID,
			ConnName:               externalDeviceConn.ConnectionName,
			Phase1RemoteIdentifier: ph1RemoteId,
		}

		err = client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update phase 1 remote identifier: %w", err)
		}
	}

	if len(phase1RemoteIdentifier) == 2 {
		var ph1RemoteId string

		if phase1RemoteIdentifier[0] == nil && phase1RemoteIdentifier[1] != nil {
			ph1RemoteId = "\"\"" + "," + mustString(phase1RemoteIdentifier[1])
		} else if phase1RemoteIdentifier[0] != nil && phase1RemoteIdentifier[1] == nil {
			ph1RemoteId = mustString(phase1RemoteIdentifier[0]) + "," + "\"\""
		} else if phase1RemoteIdentifier[0] == nil && phase1RemoteIdentifier[1] == nil {
			ph1RemoteId = "\"\", \"\""
		} else {
			ph1RemoteId = strings.Join(ph1RemoteIdList, ",")
		}

		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:                 externalDeviceConn.GwName,
			VpcID:                  externalDeviceConn.VpcID,
			ConnName:               externalDeviceConn.ConnectionName,
			Phase1RemoteIdentifier: ph1RemoteId,
		}

		err = client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update phase 1 remote identifier: %w", err)
		}
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}

		err = client.EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %w", err)
		}
	}

	if phase1LocalIdentifier, ok := d.GetOk("phase1_local_identifier"); ok {
		s2c := &goaviatrix.EditSite2Cloud{
			VpcID:    getString(d, "vpc_id"),
			ConnName: getString(d, "connection_name"),
		}
		if phase1LocalIdentifier == "private_ip" {
			s2c.Phase1LocalIdentifier = "private_ip"
			err = client.EditSite2CloudPhase1LocalIdentifier(s2c)
			if err != nil {
				return fmt.Errorf("could not set phase1 local identificer to private_ip for connection: %s: %w", s2c.ConnName, err)
			}
		}
	}

	return resourceAviatrixSpokeExternalDeviceConnReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSpokeExternalDeviceConnReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeExternalDeviceConnRead(d, meta)
	}
	return nil
}

func resourceAviatrixSpokeExternalDeviceConnRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	connectionName := getString(d, "connection_name")
	vpcID := getString(d, "vpc_id")
	if connectionName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'connection_name' or 'vpc_id' received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return fmt.Errorf("expected import ID in the form 'connection_name~vpc_id' instead got %q", id)
		}
		mustSet(d, "connection_name", parts[0])
		mustSet(d, "vpc_id", parts[1])
		d.SetId(id)
	}

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          getString(d, "vpc_id"),
		ConnectionName: getString(d, "connection_name"),
		GwName:         getString(d, "gw_name"),
	}

	localGateway, err := getGatewayDetails(client, externalDeviceConn.GwName)
	if err != nil {
		return fmt.Errorf("could not get local gateway details: %w", err)
	}

	conn, err := client.GetExternalDeviceConnDetail(externalDeviceConn, localGateway)
	log.Printf("[TRACE] Reading Aviatrix external device conn: %s : %#v", getString(d, "connection_name"), externalDeviceConn)

	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix external device conn: %w, %#v", err, externalDeviceConn)
	}

	if conn != nil {
		mustSet(d, "vpc_id", conn.VpcID)
		mustSet(d, "connection_name", conn.ConnectionName)
		mustSet(d, "gw_name", conn.GwName)
		mustSet(d, "connection_type", conn.ConnectionType)
		mustSet(d, "remote_tunnel_cidr", conn.RemoteTunnelCidr)
		mustSet(d, "enable_event_triggered_ha", conn.EventTriggeredHA)
		mustSet(d, "enable_jumbo_frame", conn.EnableJumboFrame)
		mustSet(d, "enable_ipv6", conn.EnableIpv6)
		if conn.TunnelProtocol == "LAN" {
			mustSet(d, "remote_lan_ip", conn.RemoteLanIP)
			mustSet(d, "local_lan_ip", conn.LocalLanIP)
			mustSet(d, "enable_bgp_lan_activemesh", conn.EnableBgpLanActiveMesh)
		} else {
			mustSet(d, "remote_gateway_ip", conn.RemoteGatewayIP)
			mustSet(d, "local_tunnel_cidr", conn.LocalTunnelCidr)
			mustSet(d, "enable_bgp_lan_activemesh", false)
		}
		if conn.ConnectionType == "bgp" {
			if conn.BgpLocalAsNum != 0 {
				mustSet(d, "bgp_local_as_num", strconv.Itoa(conn.BgpLocalAsNum))
			}
			if conn.BgpRemoteAsNum != 0 {
				mustSet(d, "bgp_remote_as_num", strconv.Itoa(conn.BgpRemoteAsNum))
			}
			if conn.BackupBgpRemoteAsNum != 0 {
				mustSet(d, "backup_bgp_remote_as_num", strconv.Itoa(conn.BackupBgpRemoteAsNum))
			}
		} else {
			mustSet(d, "remote_subnet", conn.RemoteSubnet)
		}
		mustSet(d, "local_subnet", conn.LocalSubnet)
		mustSet(d, "proxy_id_enabled", conn.ProxyIdEnabled)
		if conn.DirectConnect == "enabled" {
			mustSet(d, "direct_connect", true)
		} else {
			mustSet(d, "direct_connect", false)
		}
		mustSet(d, "phase1_local_identifier", conn.Phase1LocalIdentifier)

		if conn.CustomAlgorithms {
			mustSet(d, "custom_algorithms", true)
			mustSet(d, "phase_1_authentication", conn.Phase1Auth)
			mustSet(d, "phase_2_authentication", conn.Phase2Auth)
			mustSet(d, "phase_1_dh_groups", conn.Phase1DhGroups)
			mustSet(d, "phase_2_dh_groups", conn.Phase2DhGroups)
			mustSet(d, "phase_1_encryption", conn.Phase1Encryption)
			mustSet(d, "phase_2_encryption", conn.Phase2Encryption)
		} else {
			mustSet(d, "custom_algorithms", false)
		}

		if conn.HAEnabled == "enabled" {
			mustSet(d, "ha_enabled", true)
			mustSet(d, "backup_remote_tunnel_cidr", conn.BackupRemoteTunnelCidr)
			if conn.TunnelProtocol == "LAN" {
				mustSet(d, "backup_remote_lan_ip", conn.BackupRemoteLanIP)
				mustSet(d, "backup_local_lan_ip", conn.BackupLocalLanIP)
			} else {
				mustSet(d, "backup_remote_gateway_ip", conn.BackupRemoteGatewayIP)
				mustSet(d, "backup_local_tunnel_cidr", conn.BackupLocalTunnelCidr)
			}
			if conn.BackupDirectConnect == "enabled" {
				mustSet(d, "backup_direct_connect", true)
			} else {
				mustSet(d, "backup_direct_connect", false)
			}
		} else {
			mustSet(d, "ha_enabled", false)
			mustSet(d, "backup_direct_connect", false)
		}

		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: externalDeviceConn.GwName})
		if err != nil {
			return fmt.Errorf("could not get advanced config for spoke gateway: %w", err)
		}

		for _, v := range spokeAdvancedConfig.ConnectionLearnedCIDRApprovalInfo {
			if v.ConnName == externalDeviceConn.ConnectionName {
				mustSet(d, "enable_learned_cidrs_approval", v.EnabledApproval == "yes")
				err := d.Set("approved_cidrs", v.ApprovedLearnedCidrs)
				if err != nil {
					return fmt.Errorf("could not set 'approved_cidrs' in state: %w", err)
				}
				break
			}
		}
		if len(spokeAdvancedConfig.ConnectionLearnedCIDRApprovalInfo) == 0 {
			mustSet(d, "enable_learned_cidrs_approval", false)
			mustSet(d, "approved_cidrs", nil)
		}
		mustSet(d, "enable_bfd", conn.EnableBfd)
		if conn.EnableBfd {
			var bgpBfdConfig []map[string]interface{}
			bfd := conn.BgpBfdConfig
			bfdMap := make(map[string]interface{})
			if bfd.TransmitInterval != 0 {
				bfdMap["transmit_interval"] = bfd.TransmitInterval
			}
			if bfd.ReceiveInterval != 0 {
				bfdMap["receive_interval"] = bfd.ReceiveInterval
			}
			if bfd.Multiplier != 0 {
				bfdMap["multiplier"] = bfd.Multiplier
			}
			bgpBfdConfig = append(bgpBfdConfig, bfdMap)
			mustSet(d, "bgp_bfd", bgpBfdConfig)
		}

		if conn.EnableIkev2 == "enabled" {
			mustSet(d, "enable_ikev2", true)
		} else {
			mustSet(d, "enable_ikev2", false)
		}

		if err := d.Set("manual_bgp_advertised_cidrs", conn.ManualBGPCidrs); err != nil {
			return fmt.Errorf("setting 'manual_bgp_advertised_cidrs' into state: %w", err)
		}
		if conn.TunnelProtocol == "" {
			mustSet(d, "tunnel_protocol", "IPsec")
		} else {
			mustSet(d, "tunnel_protocol", conn.TunnelProtocol)
		}
		if conn.TunnelProtocol == "LAN" {
			err = d.Set("remote_vpc_name", conn.PeerVnetID)
			if err != nil {
				return fmt.Errorf("could not set value for remote_vpc_name: %w", err)
			}
		}

		if conn.Phase1RemoteIdentifier != "" {
			ph1RemoteId := strings.Split(conn.Phase1RemoteIdentifier, ",")
			for i, v := range ph1RemoteId {
				ph1RemoteId[i] = strings.TrimSpace(v)
			}

			haEnabled := getBool(d, "ha_enabled")

			if haEnabled && len(ph1RemoteId) == 1 && ph1RemoteId[0] == "" {
				ph1RemoteId = append(ph1RemoteId, "")
			}
			mustSet(d, "phase1_remote_identifier", ph1RemoteId)
		}

		if conn.PrependAsPath != "" {
			var prependAsPath []string
			for _, str := range strings.Split(conn.PrependAsPath, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("prepend_as_path", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set value for prepend_as_path: %w", err)
			}
		}
		err = d.Set("enable_bgp_multihop", conn.EnableBgpMultihop)
		if err != nil {
			return fmt.Errorf("could not set value for enable_bgp_multihop: %w", err)
		}

		err = d.Set("connection_bgp_send_communities", conn.BgpSendCommunities)
		if err != nil {
			return fmt.Errorf("could not set value for connection_bgp_send_communities: %w", err)
		}
		err = d.Set("connection_bgp_send_communities_additive", conn.BgpSendCommunitiesAdditive)
		if err != nil {
			return fmt.Errorf("could not set value for connection_bgp_send_communities: %w", err)
		}
		err = d.Set("connection_bgp_send_communities_block", conn.BgpSendCommunitiesBlock)
		if err != nil {
			return fmt.Errorf("could not set value for connection_bgp_send_communities: %w", err)
		}

		if conn.EnableIpv6 {
			if conn.TunnelProtocol == "LAN" {
				mustSet(d, "remote_lan_ipv6", conn.RemoteLanIPv6)
			}

			if conn.TunnelProtocol == "IPsec" || conn.TunnelProtocol == "GRE" {
				mustSet(d, "external_device_ipv6", conn.ExternalDeviceIPv6)
				if conn.HAEnabled == "enabled" {
					mustSet(d, "external_device_backup_ipv6", conn.ExternalDeviceBackupIPv6)
				}
			}
		}
	}

	d.SetId(conn.ConnectionName + "~" + conn.VpcID)
	return nil
}

func resourceAviatrixSpokeExternalDeviceConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	d.Partial(true)

	approvedCidrs := getStringSet(d, "approved_cidrs")
	enableLearnedCIDRApproval := getBool(d, "enable_learned_cidrs_approval")
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 {
		return fmt.Errorf("updating spoke external device conn: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	gwName := getString(d, "gw_name")
	connName := getString(d, "connection_name")
	connType := getString(d, "connection_type")

	if d.HasChange("enable_learned_cidrs_approval") {
		enableLearnedCIDRApproval := getBool(d, "enable_learned_cidrs_approval")
		if enableLearnedCIDRApproval {
			err := client.EnableSpokeConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not enable learned cidr approval: %w", err)
			}
		} else {
			err := client.DisableSpokeConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not disable learned cidr approval: %w", err)
			}
		}
	}

	if d.HasChange("manual_bgp_advertised_cidrs") {
		manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
		err := client.EditSpokeConnectionBGPManualAdvertiseCIDRs(gwName, connName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual advertise manual cidrs: %w", err)
		}
	}

	if d.HasChange("enable_event_triggered_ha") {
		vpcID := getString(d, "vpc_id")
		if getBool(d, "enable_event_triggered_ha") {
			err := client.EnableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not enable event triggered HA for external device conn during update: %w", err)
			}
		} else {
			err := client.DisableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not disable event triggered HA for external device conn during update: %w", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:            getString(d, "vpc_id"),
			ConnectionName:   getString(d, "connection_name"),
			GwName:           getString(d, "gw_name"),
			ConnectionType:   getString(d, "connection_type"),
			TunnelProtocol:   getString(d, "tunnel_protocol"),
			EnableJumboFrame: getBool(d, "enable_jumbo_frame"),
		}
		if externalDeviceConn.EnableJumboFrame {
			if externalDeviceConn.ConnectionType != "bgp" {
				return fmt.Errorf("jumbo frame is only supported on BGP connection")
			}
			if err := client.EnableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return fmt.Errorf("failed to enable jumbo frame for external device connection %q: %w", externalDeviceConn.ConnectionName, err)
			}
		} else {
			if err := client.DisableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return fmt.Errorf("failed to disable jumbo frame for external device connection %q: %w", externalDeviceConn.ConnectionName, err)
			}
		}
	}

	if d.HasChange("approved_cidrs") {
		err := client.UpdateSpokeConnectionPendingApprovedCidrs(gwName, connName, approvedCidrs)
		if err != nil {
			return fmt.Errorf("could not update spoke external device conn learned cidrs during update: %w", err)
		}
	}

	if d.HasChange("phase1_remote_identifier") {
		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:   gwName,
			VpcID:    getString(d, "vpc_id"),
			ConnName: connName,
		}

		haEnabled := getBool(d, "ha_enabled")
		phase1RemoteIdentifier := getList(d, "phase1_remote_identifier")
		ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)
		if haEnabled && len(phase1RemoteIdentifier) != 2 {
			return fmt.Errorf("please either set two phase 1 remote IDs, when HA is enabled")
		} else if !haEnabled && len(phase1RemoteIdentifier) != 1 {
			return fmt.Errorf("please set one phase 1 remote ID, when HA is disabled")
		}

		var ph1RemoteId string

		if len(phase1RemoteIdentifier) == 1 {
			if mustString(phase1RemoteIdentifier[0]) == "" {
				ph1RemoteId = "\"\""
			} else {
				ph1RemoteId = mustString(phase1RemoteIdentifier[0])
			}
		}

		if len(phase1RemoteIdentifier) == 2 {
			if mustString(phase1RemoteIdentifier[0]) == "" && mustString(phase1RemoteIdentifier[1]) != "" {
				ph1RemoteId = "\"\"" + "," + mustString(phase1RemoteIdentifier[1])
			} else if mustString(phase1RemoteIdentifier[0]) != "" && mustString(phase1RemoteIdentifier[1]) == "" {
				ph1RemoteId = mustString(phase1RemoteIdentifier[0]) + "," + "\"\""
			} else if mustString(phase1RemoteIdentifier[0]) == "" && mustString(phase1RemoteIdentifier[1]) == "" {
				ph1RemoteId = "\"\", \"\""
			} else {
				ph1RemoteId = strings.Join(ph1RemoteIdList, ",")
			}
		}

		editSite2cloud.Phase1RemoteIdentifier = ph1RemoteId

		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update phase 1 remote identifier: %w", err)
		}
	}

	if d.HasChange("prepend_as_path") {
		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			ConnectionName: connName,
			GwName:         gwName,
			ConnectionType: getString(d, "connection_type"),
		}
		if externalDeviceConn.ConnectionType != "bgp" {
			return fmt.Errorf("'prepend_as_path' only supports 'bgp' connection. Can't update 'prepend_as_path' for '%s' connection", externalDeviceConn.ConnectionType)
		}

		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}
		err := client.EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path: %w", err)
		}
	}

	enableBfd := getBool(d, "enable_bfd")

	if connType != "bgp" && enableBfd {
		return fmt.Errorf("cannot enable BFD for non-BGP connection")
	}
	// get the BGP BFD config
	bgpBfdConfig := getList(d, "bgp_bfd")

	if d.HasChanges("enable_bfd", "bgp_bfd") {
		// bgp bfd is enabled
		if enableBfd {
			bgpBfd := goaviatrix.GetUpdatedBgpBfdConfig(bgpBfdConfig)
			externalDeviceConn := &goaviatrix.ExternalDeviceConn{
				GwName:         getString(d, "gw_name"),
				ConnectionName: getString(d, "connection_name"),
				EnableBfd:      enableBfd,
				BgpBfdConfig:   bgpBfd,
			}
			err := client.EditConnectionBgpBfd(externalDeviceConn)
			if err != nil {
				return fmt.Errorf("could not update BGP BFD config: %w", err)
			}
		} else {
			// bgp bfd is disabled
			if len(bgpBfdConfig) > 0 {
				return fmt.Errorf("bgp_bfd config can't be set when BFD is disabled")
			}
			externalDeviceConn := &goaviatrix.ExternalDeviceConn{
				GwName:         getString(d, "gw_name"),
				ConnectionName: getString(d, "connection_name"),
				EnableBfd:      enableBfd,
			}
			err := client.EditConnectionBgpBfd(externalDeviceConn)
			if err != nil {
				return fmt.Errorf("could not disable BGP BFD config: %w", err)
			}
		}
	}

	if d.HasChange("bgp_md5_key") {
		if getString(d, "connection_type") != "bgp" {
			return fmt.Errorf("can't update BGP MD5 authentication key since it is only supported for BGP connection")
		}

		oldKey, newKey := d.GetChange("bgp_md5_key")
		oldKeyStr := mustString(oldKey)
		newKeyStr := mustString(newKey)
		oldKeyList := strings.Split(oldKeyStr, ",")
		newKeyList := strings.Split(newKeyStr, ",")
		var bgpRemoteIp []string
		if strings.ToUpper(getString(d, "tunnel_protocol")) == "LAN" {
			bgpRemoteIp = strings.Split(getString(d, "remote_lan_ip"), ",")
		} else {
			bgpRemoteIp = strings.Split(getString(d, "remote_tunnel_cidr"), ",")
		}
		if newKeyStr != "" && len(newKeyList) != len(bgpRemoteIp) {
			return fmt.Errorf("can't update BGP MD5 authentication key since it is not set correctly for BGP connection: %s", getString(d, "connection_name"))
		}
		for i, v := range bgpRemoteIp {
			bgpMd5Key := ""
			if newKeyStr != "" {
				bgpMd5Key = newKeyList[i]
			}
			if newKeyStr != "" && oldKeyStr != "" && strings.TrimSpace(newKeyList[i]) == strings.TrimSpace(oldKeyList[i]) {
				continue
			}
			editBgpMd5Key := &goaviatrix.EditBgpMd5Key{
				GwName:         gwName,
				ConnectionName: connName,
				BgpRemoteIP:    v,
				BgpMd5Key:      bgpMd5Key,
			}
			err := client.EditBgpMd5Key(editBgpMd5Key)
			if err != nil {
				return fmt.Errorf("failed to update BGP MD5 authentication key: %w", err)
			}
		}
	}

	if d.HasChanges("enable_bgp_multihop") {
		enableMultihop := getBool(d, "enable_bgp_multihop")

		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			GwName:            gwName,
			ConnectionName:    connName,
			EnableBgpMultihop: enableMultihop,
		}
		if err := client.EditConnectionBgpMultihop(externalDeviceConn); err != nil {
			return fmt.Errorf("could not update multihop: %w", err)
		}
	}

	if d.HasChange("backup_bgp_md5_key") {
		if getString(d, "connection_type") != "bgp" {
			return fmt.Errorf("can't update backup BGP MD5 authentication key since it is only supported for BGP connection")
		}
		if !getBool(d, "ha_enabled") {
			return fmt.Errorf("can't update BGP backup MD5 authentication key since ha is not enabled")
		}

		oldKey, newKey := d.GetChange("backup_bgp_md5_key")
		oldKeyStr := mustString(oldKey)
		newKeyStr := mustString(newKey)
		oldKeyList := strings.Split(oldKeyStr, ",")
		newKeyList := strings.Split(newKeyStr, ",")
		var bgpRemoteIp []string
		if strings.ToUpper(getString(d, "tunnel_protocol")) == "LAN" {
			bgpRemoteIp = strings.Split(getString(d, "backup_remote_lan_ip"), ",")
		} else {
			bgpRemoteIp = strings.Split(getString(d, "backup_remote_tunnel_cidr"), ",")
		}
		if newKeyStr != "" && len(newKeyList) != len(bgpRemoteIp) {
			return fmt.Errorf("can't update backup BGP MD5 authentication key since it is not set correctly for BGP connection: %s", getString(d, "connection_name"))
		}
		for i, v := range bgpRemoteIp {
			bgpMd5Key := ""
			if newKeyStr != "" {
				bgpMd5Key = newKeyList[i]
			}
			if newKeyStr != "" && oldKeyStr != "" && strings.TrimSpace(newKeyList[i]) == strings.TrimSpace(oldKeyList[i]) {
				continue
			}
			editBgpMd5Key := &goaviatrix.EditBgpMd5Key{
				GwName:         gwName,
				ConnectionName: connName,
				BgpRemoteIP:    v,
				BgpMd5Key:      bgpMd5Key,
			}
			err := client.EditBgpMd5Key(editBgpMd5Key)
			if err != nil {
				return fmt.Errorf("failed to update backup BGP MD5 authentication key: %w", err)
			}
		}
	}

	if d.HasChange("phase1_local_identifier") {
		s2c := &goaviatrix.EditSite2Cloud{
			VpcID:                 getString(d, "vpc_id"),
			ConnName:              getString(d, "connection_name"),
			Phase1LocalIdentifier: getString(d, "phase1_local_identifier"),
		}
		err := client.EditSite2CloudPhase1LocalIdentifier(s2c)
		if err != nil {
			return fmt.Errorf("could not update phase1 local identificer for connection: %s: %w", s2c.ConnName, err)
		}
	}

	if d.HasChanges("connection_bgp_send_communities", "connection_bgp_send_communities_additive", "connection_bgp_send_communities_block") {
		// Detect whether the user wants to change the set of BGP communities sent on a given connection
		// if so, update the connection with the new set of communities, either additively or as a replacement
		// or block the communities entirely, depending on the user's choice
		sendComm := getString(d, "connection_bgp_send_communities")

		sendAdditive := getBool(d, "connection_bgp_send_communities_additive")

		sendBlock := getBool(d, "connection_bgp_send_communities_block")

		bgpSendCommunities := &goaviatrix.BgpSendCommunities{
			ConnectionName:      connName,
			GwName:              gwName,
			ConnSendCommunities: sendComm,
			ConnSendAdditive:    sendAdditive,
			ConnSendBlock:       sendBlock,
		}
		if err := client.ConnectionBGPSendCommunities(bgpSendCommunities); err != nil {
			return fmt.Errorf("failed to update bgp connection based communities for connection %q", bgpSendCommunities.ConnectionName)
		}
	}

	if d.HasChange("local_subnet") {
		vpcID := getString(d, "vpc_id")
		localSubnet := getString(d, "local_subnet")
		err := client.EditTransitConnectionLocalSubnet(vpcID, connName, localSubnet)
		if err != nil {
			return fmt.Errorf("could not update spoke external device conn local subnet: %w", err)
		}
	}

	if d.HasChange("proxy_id_enabled") {
		editSite2cloud := &goaviatrix.EditSite2Cloud{
			VpcID:    getString(d, "vpc_id"),
			ConnName: connName,
		}
		if getBool(d, "proxy_id_enabled") {
			editSite2cloud.ProxyIdEnabled = "true"
		} else {
			editSite2cloud.ProxyIdEnabled = "false"
		}
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update proxy_id_enabled: %w", err)
		}
	}

	d.Partial(false)

	return resourceAviatrixSpokeExternalDeviceConnRead(d, meta)
}

func resourceAviatrixSpokeExternalDeviceConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          getString(d, "vpc_id"),
		ConnectionName: getString(d, "connection_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix external device connection: %#v", externalDeviceConn)

	err := client.DeleteExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix external device connection: %w", err)
	}

	return nil
}
