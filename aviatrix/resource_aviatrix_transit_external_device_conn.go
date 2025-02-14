package aviatrix

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultBfdReceiveInterval  = 300
	defaultBfdTransmitInterval = 300
	defaultBfdMultiplier       = 3
)

var defaultBfdConfig = goaviatrix.BgpBfdConfig{
	TransmitInterval: defaultBfdTransmitInterval,
	ReceiveInterval:  defaultBfdReceiveInterval,
	Multiplier:       defaultBfdMultiplier,
}

func resourceAviatrixTransitExternalDeviceConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitExternalDeviceConnCreate,
		Read:   resourceAviatrixTransitExternalDeviceConnRead,
		Update: resourceAviatrixTransitExternalDeviceConnUpdate,
		Delete: resourceAviatrixTransitExternalDeviceConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC where the Transit Gateway is located. For GCP BGP over LAN connection, it is in the format of 'vpc_name~-~account_name'.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of the transit external device connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Transit Gateway.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bgp",
				ForceNew:     true,
				Description:  "Connection type. Valid values: 'bgp', 'static'. Default value: 'bgp'.",
				ValidateFunc: validation.StringInSlice([]string{"bgp", "static"}, false),
			},
			"tunnel_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPsec",
				ForceNew:     true,
				Description:  "Tunnel Protocol. Valid values: 'IPsec', 'GRE' or 'LAN'. Default value: 'IPsec'. Case insensitive.",
				ValidateFunc: validation.StringInSlice([]string{"IPsec", "GRE", "LAN"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToUpper(old) == strings.ToUpper(new)
				},
			},
			"enable_bgp_lan_activemesh": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "Switch to enable BGP LAN ActiveMesh. Only valid for GCP and Azure with Remote Gateway HA enabled. " +
					"Requires Azure Remote Gateway insane mode enabled. Valid values: true, false. Default: false. " +
					"Available as of provider version R2.21+.",
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Remote CIDRs joined as a string with ','. Required for a 'static' type connection.",
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
				Description:      "Source CIDR for the tunnel from the Aviatrix transit gateway.",
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
				Description:      "Source CIDR for the tunnel from the backup Aviatrix transit gateway.",
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
			"enable_edge_segmentation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Switch to allow this connection to communicate with a Network Domain via Connection Policy.",
			},
			"switch_to_ha_standby_gateway": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Only valid for Transit Gateway's with Active-Standby Mode enabled. Valid values: true, false. Default: false.",
			},
			"enable_learned_cidrs_approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable learned CIDR approval for the connection. Only valid with 'connection_type' = 'bgp'." +
					" Requires the transit_gateway's 'learned_cidrs_approval_mode' attribute be set to 'connection'. " +
					"Valid values: true, false. Default value: false. Available as of provider version R2.18+.",
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
					" = 'bgp'. Available as of provider version R2.18+.",
			},
			"remote_vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the remote VPC for a LAN BGP connection. Only valid when 'connection_type' = 'bgp' and tunnel_protocol' = 'LAN' with an Azure transit gateway. Must be in the form \"<VNET-name>:<resource-group-name>\". Available as of provider version R2.18+.",
			},
			"remote_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote LAN IP.",
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
				Description: "Local LAN IP. Required for GCP BGP over LAN Connection.",
			},
			"backup_remote_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup Remote LAN IP.",
			},
			"backup_local_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup Local LAN IP. Required for GCP BGP over LAN Connection with HA enabled.",
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
			"enable_edge_underlay": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable BGP over WAN underlay.",
			},
			"disable_activemesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Disable ActiveMesh, no crossing tunnels",
			},
			"tunnel_src_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Description:      "Local gateway tunnel source IP",
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceOnlyInString,
			},
		},
	}
}

func resourceAviatrixTransitExternalDeviceConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:                  d.Get("vpc_id").(string),
		ConnectionName:         d.Get("connection_name").(string),
		GwName:                 d.Get("gw_name").(string),
		ConnectionType:         d.Get("connection_type").(string),
		RemoteGatewayIP:        d.Get("remote_gateway_ip").(string),
		RemoteSubnet:           d.Get("remote_subnet").(string),
		PreSharedKey:           d.Get("pre_shared_key").(string),
		LocalTunnelCidr:        d.Get("local_tunnel_cidr").(string),
		RemoteTunnelCidr:       d.Get("remote_tunnel_cidr").(string),
		Phase1Auth:             d.Get("phase_1_authentication").(string),
		Phase1DhGroups:         d.Get("phase_1_dh_groups").(string),
		Phase1Encryption:       d.Get("phase_1_encryption").(string),
		Phase2Auth:             d.Get("phase_2_authentication").(string),
		Phase2DhGroups:         d.Get("phase_2_dh_groups").(string),
		Phase2Encryption:       d.Get("phase_2_encryption").(string),
		BackupRemoteGatewayIP:  d.Get("backup_remote_gateway_ip").(string),
		BackupPreSharedKey:     d.Get("backup_pre_shared_key").(string),
		BackupLocalTunnelCidr:  d.Get("backup_local_tunnel_cidr").(string),
		BackupRemoteTunnelCidr: d.Get("backup_remote_tunnel_cidr").(string),
		PeerVnetId:             d.Get("remote_vpc_name").(string),
		RemoteLanIP:            d.Get("remote_lan_ip").(string),
		LocalLanIP:             d.Get("local_lan_ip").(string),
		BackupRemoteLanIP:      d.Get("backup_remote_lan_ip").(string),
		BackupLocalLanIP:       d.Get("backup_local_lan_ip").(string),
		BgpMd5Key:              d.Get("bgp_md5_key").(string),
		BackupBgpMd5Key:        d.Get("backup_bgp_md5_key").(string),
		EnableJumboFrame:       d.Get("enable_jumbo_frame").(bool),
		EnableBgpMultihop:      d.Get("enable_bgp_multihop").(bool),
		EnableEdgeUnderlay:     d.Get("enable_edge_underlay").(bool),
		DisableActivemesh:      d.Get("disable_activemesh").(bool),
	}

	if !externalDeviceConn.EnableEdgeUnderlay && externalDeviceConn.ConnectionName == "" {
		return fmt.Errorf("'connection_name' is required when 'enable_edge_underlay' is false")
	}

	if externalDeviceConn.EnableEdgeUnderlay && externalDeviceConn.ConnectionName != "" {
		return fmt.Errorf("please set 'connection_name' to empty when 'enable_edge_underlay' is true")
	}

	tunnelProtocol := strings.ToUpper(d.Get("tunnel_protocol").(string))
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
	if externalDeviceConn.RemoteGatewayIP == "" && externalDeviceConn.TunnelProtocol != "LAN" {
		return fmt.Errorf("'remote_gateway_ip' is required when 'tunnel_protocol' != 'LAN'")
	}

	bgpLocalAsNum, err := strconv.Atoi(d.Get("bgp_local_as_num").(string))
	if err == nil {
		externalDeviceConn.BgpLocalAsNum = bgpLocalAsNum
	}
	bgpRemoteAsNum, err := strconv.Atoi(d.Get("bgp_remote_as_num").(string))
	if err == nil {
		externalDeviceConn.BgpRemoteAsNum = bgpRemoteAsNum
	}
	backupBgpLocalAsNum, err := strconv.Atoi(d.Get("backup_bgp_remote_as_num").(string))
	if err == nil {
		externalDeviceConn.BackupBgpRemoteAsNum = backupBgpLocalAsNum
	}

	directConnect := d.Get("direct_connect").(bool)
	if directConnect {
		externalDeviceConn.DirectConnect = "true"
	}

	haEnabled := d.Get("ha_enabled").(bool)
	if haEnabled {
		externalDeviceConn.HAEnabled = "true"
	}

	backupDirectConnect := d.Get("backup_direct_connect").(bool)
	if backupDirectConnect {
		externalDeviceConn.BackupDirectConnect = "true"
	}

	enableEdgeSegmentation := d.Get("enable_edge_segmentation").(bool)
	if enableEdgeSegmentation {
		externalDeviceConn.EnableEdgeSegmentation = "true"
	}

	if externalDeviceConn.ConnectionType == "bgp" && externalDeviceConn.RemoteSubnet != "" {
		return fmt.Errorf("'remote_subnet' is needed for connection type of 'static' not 'bgp'")
	} else if externalDeviceConn.ConnectionType == "static" && (externalDeviceConn.BgpLocalAsNum != 0 || externalDeviceConn.BgpRemoteAsNum != 0) {
		return fmt.Errorf("'bgp_local_as_num' and 'bgp_remote_as_num' are needed for connection type of 'bgp' not 'static'")
	}

	customAlgorithms := d.Get("custom_algorithms").(bool)
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

	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if externalDeviceConn.ConnectionType != "bgp" && enableLearnedCIDRApproval {
		return fmt.Errorf("'connection_type' must be 'bgp' if 'enable_learned_cidrs_approval' is set to true")
	}
	manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
	if externalDeviceConn.ConnectionType != "bgp" && len(manualBGPCidrs) != 0 {
		return fmt.Errorf("'connection_type' must be 'bgp' if 'manual_bgp_advertised_cidrs' is not empty")
	}

	approvedCidrs := getStringSet(d, "approved_cidrs")
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 {
		return fmt.Errorf("creating transit external device conn: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	enableIkev2 := d.Get("enable_ikev2").(bool)
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
	if externalDeviceConn.PeerVnetId != "" && (externalDeviceConn.ConnectionType != "bgp" || externalDeviceConn.TunnelProtocol != "LAN") {
		return fmt.Errorf("'remote_vpc_name' is only valid for 'connection_type' = 'bgp' and 'tunnel_protocol' = 'LAN'")
	}
	if externalDeviceConn.TunnelProtocol == "LAN" {
		if externalDeviceConn.DirectConnect == "true" || externalDeviceConn.BackupDirectConnect == "true" {
			return fmt.Errorf("enabling 'direct_connect' or 'backup_direct_connect' is not allowed for BGP over LAN connections")
		}
	}

	phase1RemoteIdentifier := d.Get("phase1_remote_identifier").([]interface{})
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

	if d.Get("enable_bgp_lan_activemesh").(bool) {
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

	enableJumboFrame := d.Get("enable_jumbo_frame").(bool)
	if enableJumboFrame {
		if externalDeviceConn.ConnectionType != "bgp" {
			return fmt.Errorf("jumbo frame is only supported on bgp connection")
		}
	}
	if val, ok := d.Get("tunnel_src_ip").(string); ok && val != "" {
		externalDeviceConn.TunnelSrcIP = val
	}
	if externalDeviceConn.PreSharedKey != "" {
		externalDeviceConn.AuthType = "psk"
	}

	if !externalDeviceConn.EnableBgpMultihop && externalDeviceConn.ConnectionType != "bgp" {
		return fmt.Errorf("multihop can only be configured for BGP connections")
	}

	flag := false
	defer resourceAviatrixTransitExternalDeviceConnReadIfRequired(d, meta, &flag)

	var edgeExternalDeviceConn goaviatrix.EdgeExternalDeviceConn
	if externalDeviceConn.EnableEdgeUnderlay {
		edgeExternalDeviceConn = goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)
	}

	var result string
	try, maxTries, backoff := 0, 8, 1000*time.Millisecond
	for {
		try++
		if externalDeviceConn.EnableEdgeUnderlay {
			result, err = client.CreateEdgeExternalDeviceConn(&edgeExternalDeviceConn)
		} else {
			err = client.CreateExternalDeviceConn(externalDeviceConn)
		}
		if err != nil {
			if strings.Contains(err.Error(), "is not up") {
				if try == maxTries {
					return fmt.Errorf("couldn't create Aviatrix transit external device connection: %s", err)
				}
				time.Sleep(backoff)
				// Double the backoff time after each failed try
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to create Aviatrix transit external device connection: %s", err)
		}
		break
	}

	if externalDeviceConn.EnableEdgeUnderlay {
		re := regexp.MustCompile(`underlay BGP connection (.*) (?:in|on)`)
		match := re.FindStringSubmatch(result)
		if len(match) < 2 {
			return fmt.Errorf("could not get underlay BGP connection name")
		}
		connName := match[1]
		if err := d.Set("connection_name", connName); err != nil {
			return fmt.Errorf("error setting connection_name: %w", err)
		}
		externalDeviceConn.ConnectionName = connName
	}
	d.SetId(externalDeviceConn.ConnectionName + "~" + externalDeviceConn.VpcID)

	enableBFD, ok := d.Get("enable_bfd").(bool)
	if !ok {
		return fmt.Errorf("expected enable_bfd to be a boolean, but got %T", d.Get("enable_bfd"))
	}

	if d.Get("connection_type").(string) != "bgp" && enableBFD {
		return fmt.Errorf("BFD is only supported for BGP connection")
	}
	externalDeviceConn.EnableBfd = enableBFD
	bgp_bfd, ok := d.Get("bgp_bfd").([]interface{})
	if !ok {
		return fmt.Errorf("expected bgp_bfd to be a list of maps, but got %T", d.Get("bgp_bfd"))
	}
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
			return fmt.Errorf("could not enable BGP BFD connection: %v", err)
		}
	} else {
		// if BFD is disabled and BGP BFD config is provided then throw an error
		if len(bgp_bfd) > 0 {
			return fmt.Errorf("bgp_bfd config can't be set when BFD is disabled")
		}
	}

	transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: externalDeviceConn.GwName})
	if err != nil {
		return fmt.Errorf("could not get advanced config for transit gateway: %v", err)
	}
	if d.Get("switch_to_ha_standby_gateway").(bool) && !transitAdvancedConfig.ActiveStandbyEnabled {
		return fmt.Errorf("can not set 'switch_to_ha_standby_gateway' unless Active-Standby Mode is enabled on " +
			"the transit gateway. Please enable Active-Standby Mode on the transit gateway and try again")
	}

	if d.Get("switch_to_ha_standby_gateway").(bool) {
		if err := client.SwitchActiveTransitGateway(externalDeviceConn.GwName, externalDeviceConn.ConnectionName); err != nil {
			return fmt.Errorf("could not switch active transit gateway to HA: %v", err)
		}
	}

	if d.Get("enable_event_triggered_ha").(bool) {
		if err := client.EnableSite2CloudEventTriggeredHA(externalDeviceConn.VpcID, externalDeviceConn.ConnectionName); err != nil {
			return fmt.Errorf("could not enable event triggered HA for external device conn after create: %v", err)
		}
	}

	if enableLearnedCIDRApproval {
		err = client.EnableTransitConnectionLearnedCIDRApproval(externalDeviceConn.GwName, externalDeviceConn.ConnectionName)
		if err != nil {
			return fmt.Errorf("could not enable learned cidr approval: %v", err)
		}
		if len(approvedCidrs) > 0 {
			err = client.UpdateTransitConnectionPendingApprovedCidrs(externalDeviceConn.GwName, externalDeviceConn.ConnectionName, approvedCidrs)
			if err != nil {
				return fmt.Errorf("could not update transit external device conn approved cidrs after creation: %v", err)
			}
		}
	}

	if len(manualBGPCidrs) > 0 {
		err = client.EditTransitConnectionBGPManualAdvertiseCIDRs(externalDeviceConn.GwName, externalDeviceConn.ConnectionName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual advertised BGP cidrs: %v", err)
		}
	}

	if len(phase1RemoteIdentifier) == 1 {
		var ph1RemoteId string

		if phase1RemoteIdentifier[0] == nil {
			ph1RemoteId = "\"\""
		} else {
			ph1RemoteId = phase1RemoteIdentifier[0].(string)
		}

		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:                 externalDeviceConn.GwName,
			VpcID:                  externalDeviceConn.VpcID,
			ConnName:               externalDeviceConn.ConnectionName,
			Phase1RemoteIdentifier: ph1RemoteId,
		}

		err = client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update phase 1 remote identifier: %s", err)
		}
	}

	if len(phase1RemoteIdentifier) == 2 {
		var ph1RemoteId string

		if phase1RemoteIdentifier[0] == nil && phase1RemoteIdentifier[1] != nil {
			ph1RemoteId = "\"\"" + "," + phase1RemoteIdentifier[1].(string)
		} else if phase1RemoteIdentifier[0] != nil && phase1RemoteIdentifier[1] == nil {
			ph1RemoteId = phase1RemoteIdentifier[0].(string) + "," + "\"\""
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
			return fmt.Errorf("failed to update phase 1 remote identifier: %s", err)
		}
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		err = client.EditTransitExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %v", err)
		}
	}

	if phase1LocalIdentifier, ok := d.GetOk("phase1_local_identifier"); ok {
		s2c := &goaviatrix.EditSite2Cloud{
			VpcID:    d.Get("vpc_id").(string),
			ConnName: d.Get("connection_name").(string),
		}
		if phase1LocalIdentifier == "private_ip" {
			s2c.Phase1LocalIdentifier = "private_ip"
			err = client.EditSite2CloudPhase1LocalIdentifier(s2c)
			if err != nil {
				return fmt.Errorf("could not set phase1 local identificer to private_ip for connection: %s: %v", s2c.ConnName, err)
			}
		}
	}

	return resourceAviatrixTransitExternalDeviceConnReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitExternalDeviceConnReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitExternalDeviceConnRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitExternalDeviceConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	vpcID := d.Get("vpc_id").(string)
	if connectionName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'connection_name' or 'vpc_id' received. Import Id is %s", id)
		parts := strings.SplitN(id, "~", 2)
		if len(parts) != 2 {
			return fmt.Errorf("expected import ID in the form 'connection_name~vpc_id' instead got %q", id)
		}
		d.Set("connection_name", parts[0])
		d.Set("vpc_id", parts[1])
		d.SetId(id)
	}

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          d.Get("vpc_id").(string),
		ConnectionName: d.Get("connection_name").(string),
	}

	conn, err := client.GetExternalDeviceConnDetail(externalDeviceConn)
	log.Printf("[TRACE] Reading Aviatrix external device conn: %s : %#v", d.Get("connection_name").(string), externalDeviceConn)

	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix external device conn: %s, %#v", err, externalDeviceConn)
	}

	if conn != nil {
		d.Set("vpc_id", conn.VpcID)
		d.Set("connection_name", conn.ConnectionName)
		d.Set("gw_name", conn.GwName)
		d.Set("connection_type", conn.ConnectionType)
		d.Set("remote_tunnel_cidr", conn.RemoteTunnelCidr)
		d.Set("enable_event_triggered_ha", conn.EventTriggeredHA)
		d.Set("enable_jumbo_frame", conn.EnableJumboFrame)
		d.Set("phase1_local_identifier", conn.Phase1LocalIdentifier)

		if conn.TunnelSrcIP != "" {
			if err := d.Set("tunnel_src_ip", conn.TunnelSrcIP); err != nil {
				return fmt.Errorf("error setting tunnel_src_ip: %w", err)
			}
		}

		if conn.TunnelProtocol == "LAN" {
			d.Set("remote_lan_ip", conn.RemoteLanIP)
			d.Set("local_lan_ip", conn.LocalLanIP)
			d.Set("enable_bgp_lan_activemesh", conn.EnableBgpLanActiveMesh)
			if err := d.Set("enable_edge_underlay", conn.EnableEdgeUnderlay); err != nil {
				return fmt.Errorf("error setting enable_edge_underlay: %w", err)
			}
		} else {
			d.Set("remote_gateway_ip", conn.RemoteGatewayIP)
			d.Set("local_tunnel_cidr", conn.LocalTunnelCidr)
			if err := d.Set("disable_activemesh", conn.DisableActivemesh); err != nil {
				return fmt.Errorf("error setting disable_activemesh: %w", err)
			}
		}
		if conn.ConnectionType == "bgp" {
			if conn.BgpLocalAsNum != 0 {
				d.Set("bgp_local_as_num", strconv.Itoa(conn.BgpLocalAsNum))
			}
			if conn.BgpRemoteAsNum != 0 {
				d.Set("bgp_remote_as_num", strconv.Itoa(conn.BgpRemoteAsNum))
			}
			if conn.BackupBgpRemoteAsNum != 0 {
				d.Set("backup_bgp_remote_as_num", strconv.Itoa(conn.BackupBgpRemoteAsNum))
			}
		} else {
			d.Set("remote_subnet", conn.RemoteSubnet)
		}
		if conn.DirectConnect == "enabled" {
			d.Set("direct_connect", true)
		} else {
			d.Set("direct_connect", false)
		}

		if conn.CustomAlgorithms {
			d.Set("custom_algorithms", true)
			d.Set("phase_1_authentication", conn.Phase1Auth)
			d.Set("phase_2_authentication", conn.Phase2Auth)
			d.Set("phase_1_dh_groups", conn.Phase1DhGroups)
			d.Set("phase_2_dh_groups", conn.Phase2DhGroups)
			d.Set("phase_1_encryption", conn.Phase1Encryption)
			d.Set("phase_2_encryption", conn.Phase2Encryption)
		} else {
			d.Set("custom_algorithms", false)
		}

		if conn.HAEnabled == "enabled" {
			d.Set("ha_enabled", true)

			d.Set("backup_remote_tunnel_cidr", conn.BackupRemoteTunnelCidr)
			if conn.TunnelProtocol == "LAN" {
				d.Set("backup_remote_lan_ip", conn.BackupRemoteLanIP)
				d.Set("backup_local_lan_ip", conn.BackupLocalLanIP)
			} else {
				d.Set("backup_remote_gateway_ip", conn.BackupRemoteGatewayIP)
				d.Set("backup_local_tunnel_cidr", conn.BackupLocalTunnelCidr)
			}
			if conn.BackupDirectConnect == "enabled" {
				d.Set("backup_direct_connect", true)
			} else {
				d.Set("backup_direct_connect", false)
			}
		} else {
			d.Set("ha_enabled", false)
			d.Set("backup_direct_connect", false)
		}

		if conn.EnableEdgeSegmentation == "enabled" {
			d.Set("enable_edge_segmentation", true)
		} else {
			d.Set("enable_edge_segmentation", false)
		}

		transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: externalDeviceConn.GwName})
		if err != nil {
			return fmt.Errorf("could not get advanced config for transit gateway: %v", err)
		}

		activeGatewayType := "Primary"
		for _, v := range transitAdvancedConfig.ActiveStandbyConnections {
			if v.ConnectionName != externalDeviceConn.ConnectionName {
				continue
			}
			activeGatewayType = v.ActiveGatewayType
		}
		d.Set("switch_to_ha_standby_gateway", activeGatewayType == "HA")

		for _, v := range transitAdvancedConfig.ConnectionLearnedCIDRApprovalInfo {
			if v.ConnName == externalDeviceConn.ConnectionName {
				d.Set("enable_learned_cidrs_approval", v.EnabledApproval == "yes")
				err := d.Set("approved_cidrs", v.ApprovedLearnedCidrs)
				if err != nil {
					return fmt.Errorf("could not set 'approved_cidrs' in state: %v", err)
				}
				break
			}
		}
		if len(transitAdvancedConfig.ConnectionLearnedCIDRApprovalInfo) == 0 {
			d.Set("enable_learned_cidrs_approval", false)
			d.Set("approved_cidrs", nil)
		}

		if conn.EnableIkev2 == "enabled" {
			d.Set("enable_ikev2", true)
		} else {
			d.Set("enable_ikev2", false)
		}

		if err := d.Set("manual_bgp_advertised_cidrs", conn.ManualBGPCidrs); err != nil {
			return fmt.Errorf("setting 'manual_bgp_advertised_cidrs' into state: %v", err)
		}
		if conn.TunnelProtocol == "" {
			d.Set("tunnel_protocol", "IPsec")
		} else {
			d.Set("tunnel_protocol", conn.TunnelProtocol)
		}
		if conn.TunnelProtocol == "LAN" {
			d.Set("remote_vpc_name", conn.PeerVnetId)
		}

		if conn.Phase1RemoteIdentifier != "" {
			ph1RemoteId := strings.Split(conn.Phase1RemoteIdentifier, ",")
			for i, v := range ph1RemoteId {
				ph1RemoteId[i] = strings.TrimSpace(v)
			}

			haEnabled := d.Get("ha_enabled").(bool)

			if haEnabled && len(ph1RemoteId) == 1 && ph1RemoteId[0] == "" {
				ph1RemoteId = append(ph1RemoteId, "")
			}

			d.Set("phase1_remote_identifier", ph1RemoteId)
		}

		d.Set("enable_bfd", conn.EnableBfd)
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

			d.Set("bgp_bfd", bgpBfdConfig)
		}

		if conn.PrependAsPath != "" {
			var prependAsPath []string
			for _, str := range strings.Split(conn.PrependAsPath, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("prepend_as_path", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set value for prepend_as_path: %v", err)
			}
		}
		err = d.Set("enable_bgp_multihop", conn.EnableBgpMultihop)
		if err != nil {
			return fmt.Errorf("could not set value for enable_bgp_multihop: %w", err)
		}
	}

	d.SetId(conn.ConnectionName + "~" + conn.VpcID)
	return nil
}

func resourceAviatrixTransitExternalDeviceConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	d.Partial(true)

	approvedCidrs := getStringSet(d, "approved_cidrs")
	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 {
		return fmt.Errorf("updating transit external device conn: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	gwName := d.Get("gw_name").(string)
	connName := d.Get("connection_name").(string)
	connType := d.Get("connection_type").(string)
	transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: gwName})
	if err != nil {
		return fmt.Errorf("could not get advanced config for transit gateway: %v", err)
	}
	if d.Get("switch_to_ha_standby_gateway").(bool) && !transitAdvancedConfig.ActiveStandbyEnabled {
		return fmt.Errorf("can not set 'switch_to_ha_standby_gateway' unless Active-Standby Mode is enabled on " +
			"the transit gateway. Please enable Active-Standby Mode on the transit gateway and try again")
	}

	if d.HasChange("switch_to_ha_standby_gateway") {
		if err := client.SwitchActiveTransitGateway(gwName, connName); err != nil {
			return fmt.Errorf("could not switch active transit gateway: %v", err)
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
		if enableLearnedCIDRApproval {
			err = client.EnableTransitConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not enable learned cidr approval: %v", err)
			}
		} else {
			err = client.DisableTransitConnectionLearnedCIDRApproval(gwName, connName)
			if err != nil {
				return fmt.Errorf("could not disable learned cidr approval: %v", err)
			}
		}
	}

	if d.HasChange("manual_bgp_advertised_cidrs") {
		manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
		err := client.EditTransitConnectionBGPManualAdvertiseCIDRs(gwName, connName, manualBGPCidrs)
		if err != nil {
			return fmt.Errorf("could not edit manual advertise manual cidrs: %v", err)
		}
	}

	if d.HasChange("enable_event_triggered_ha") {
		vpcID := d.Get("vpc_id").(string)
		if d.Get("enable_event_triggered_ha").(bool) {
			err := client.EnableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not enable event triggered HA for external device conn during update: %v", err)
			}
		} else {
			err := client.DisableSite2CloudEventTriggeredHA(vpcID, connName)
			if err != nil {
				return fmt.Errorf("could not disable event triggered HA for external device conn during update: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:            d.Get("vpc_id").(string),
			ConnectionName:   d.Get("connection_name").(string),
			GwName:           d.Get("gw_name").(string),
			ConnectionType:   d.Get("connection_type").(string),
			TunnelProtocol:   d.Get("tunnel_protocol").(string),
			EnableJumboFrame: d.Get("enable_jumbo_frame").(bool),
		}
		if externalDeviceConn.EnableJumboFrame {
			if externalDeviceConn.ConnectionType != "bgp" {
				return fmt.Errorf("jumbo frame is only supported on BGP connection")
			}
			if err := client.EnableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return fmt.Errorf("could not enable jumbo frame for external device conn: %v during update: %v", externalDeviceConn.ConnectionName, err)
			}
		} else {
			if err := client.DisableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return fmt.Errorf("could not disable jumbo frame for external device conn: %v during update: %v", externalDeviceConn.ConnectionName, err)
			}
		}
	}

	if d.HasChange("approved_cidrs") {
		err := client.UpdateTransitConnectionPendingApprovedCidrs(gwName, connName, approvedCidrs)
		if err != nil {
			return fmt.Errorf("could not update transit external device conn learned cidrs during update: %v", err)
		}
	}

	if d.HasChange("phase1_remote_identifier") {
		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:   gwName,
			VpcID:    d.Get("vpc_id").(string),
			ConnName: connName,
		}

		haEnabled := d.Get("ha_enabled").(bool)
		phase1RemoteIdentifier := d.Get("phase1_remote_identifier").([]interface{})
		ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)
		if haEnabled && len(phase1RemoteIdentifier) != 2 {
			return fmt.Errorf("please either set two phase 1 remote IDs, when HA is enabled")
		} else if !haEnabled && len(phase1RemoteIdentifier) != 1 {
			return fmt.Errorf("please set one phase 1 remote ID, when HA is disabled")
		}

		var ph1RemoteId string

		if len(phase1RemoteIdentifier) == 1 {
			if phase1RemoteIdentifier[0].(string) == "" {
				ph1RemoteId = "\"\""
			} else {
				ph1RemoteId = phase1RemoteIdentifier[0].(string)
			}
		}

		if len(phase1RemoteIdentifier) == 2 {
			if phase1RemoteIdentifier[0].(string) == "" && phase1RemoteIdentifier[1].(string) != "" {
				ph1RemoteId = "\"\"" + "," + phase1RemoteIdentifier[1].(string)
			} else if phase1RemoteIdentifier[0].(string) != "" && phase1RemoteIdentifier[1].(string) == "" {
				ph1RemoteId = phase1RemoteIdentifier[0].(string) + "," + "\"\""
			} else if phase1RemoteIdentifier[0].(string) == "" && phase1RemoteIdentifier[1].(string) == "" {
				ph1RemoteId = "\"\", \"\""
			} else {
				ph1RemoteId = strings.Join(ph1RemoteIdList, ",")
			}
		}

		editSite2cloud.Phase1RemoteIdentifier = ph1RemoteId

		err = client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update phase 1 remote identifier: %s", err)
		}
	}

	if d.HasChange("prepend_as_path") {
		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			ConnectionName: connName,
			GwName:         gwName,
			ConnectionType: d.Get("connection_type").(string),
		}
		if externalDeviceConn.ConnectionType != "bgp" {
			return fmt.Errorf("'prepend_as_path' only supports 'bgp' connection. Can't update 'prepend_as_path' for '%s' connection", externalDeviceConn.ConnectionType)
		}

		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		err = client.EditTransitExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path: %v", err)
		}
	}

	if d.HasChange("remote_subnet") {
		vpcID := d.Get("vpc_id").(string)
		remoteSubnet := d.Get("remote_subnet").(string)
		err = client.EditTransitConnectionRemoteSubnet(vpcID, connName, remoteSubnet)
		if err != nil {
			return fmt.Errorf("could not update transit external device conn remote subnet: %v", err)
		}
	}

	if d.HasChange("bgp_md5_key") {
		if d.Get("connection_type").(string) != "bgp" {
			return fmt.Errorf("can't update BGP MD5 authentication key since it is only supported for BGP connection")
		}

		oldKey, newKey := d.GetChange("bgp_md5_key")
		oldKeyList := strings.Split(oldKey.(string), ",")
		newKeyList := strings.Split(newKey.(string), ",")
		var bgpRemoteIp []string
		if strings.ToUpper(d.Get("tunnel_protocol").(string)) == "LAN" {
			bgpRemoteIp = strings.Split(d.Get("remote_lan_ip").(string), ",")
		} else {
			bgpRemoteIp = strings.Split(d.Get("remote_tunnel_cidr").(string), ",")
		}
		if newKey.(string) != "" && len(newKeyList) != len(bgpRemoteIp) {
			return fmt.Errorf("can't update BGP MD5 authentication key since it is not set correctly for BGP connection: %s", d.Get("connection_name").(string))
		}
		for i, v := range bgpRemoteIp {
			bgpMd5Key := ""
			if newKey.(string) != "" {
				bgpMd5Key = newKeyList[i]
			}
			if newKey.(string) != "" && oldKey.(string) != "" && strings.TrimSpace(newKeyList[i]) == strings.TrimSpace(oldKeyList[i]) {
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
				return fmt.Errorf("failed to update BGP MD5 authentication key: %v", err)
			}
		}
	}

	enableBfd, ok := d.Get("enable_bfd").(bool)
	if !ok {
		return fmt.Errorf("expected enable_bfd to be a boolean, but got %T", d.Get("enable_bfd"))
	}
	if connType != "bgp" && enableBfd {
		return fmt.Errorf("cannot enable BFD for non-BGP connection")
	}
	// get the BGP BFD config
	bgpBfdConfig, ok := d.Get("bgp_bfd").([]interface{})
	if !ok {
		return fmt.Errorf("expected bgp_bfd to be a list of maps, but got %T", d.Get("bgp_bfd"))
	}
	if d.HasChanges("enable_bfd", "bgp_bfd") {
		// bgp bfd is enabled
		if enableBfd {
			bgpBfd := goaviatrix.GetUpdatedBgpBfdConfig(bgpBfdConfig)
			externalDeviceConn := &goaviatrix.ExternalDeviceConn{
				GwName:         d.Get("gw_name").(string),
				ConnectionName: d.Get("connection_name").(string),
				EnableBfd:      enableBfd,
				BgpBfdConfig:   bgpBfd,
			}
			err := client.EditConnectionBgpBfd(externalDeviceConn)
			if err != nil {
				return fmt.Errorf("could not update BGP BFD config: %v", err)
			}
		} else {
			// bgp bfd is disabled
			if len(bgpBfdConfig) > 0 {
				return fmt.Errorf("bgp_bfd config can't be set when BFD is disabled")
			}
			externalDeviceConn := &goaviatrix.ExternalDeviceConn{
				GwName:         d.Get("gw_name").(string),
				ConnectionName: d.Get("connection_name").(string),
				EnableBfd:      enableBfd,
			}
			err := client.EditConnectionBgpBfd(externalDeviceConn)
			if err != nil {
				return fmt.Errorf("could not disable BGP BFD config: %v", err)
			}
		}
	}

	if d.HasChanges("enable_bgp_multihop") {
		enableMultihop, ok := d.Get("enable_bgp_multihop").(bool)
		if !ok {
			return fmt.Errorf("Can't make bool from enable_bgp_multihop value: '%v'", d.Get("enable_bgp_multihop"))
		}
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
		if d.Get("connection_type").(string) != "bgp" {
			return fmt.Errorf("can't update backup BGP MD5 authentication key since it is only supported for BGP connection")
		}
		if !d.Get("ha_enabled").(bool) {
			return fmt.Errorf("can't update BGP backup MD5 authentication key since ha is not enabled")
		}

		oldKey, newKey := d.GetChange("backup_bgp_md5_key")
		oldKeyList := strings.Split(oldKey.(string), ",")
		newKeyList := strings.Split(newKey.(string), ",")
		var bgpRemoteIp []string
		if strings.ToUpper(d.Get("tunnel_protocol").(string)) == "LAN" {
			bgpRemoteIp = strings.Split(d.Get("backup_remote_lan_ip").(string), ",")
		} else {
			bgpRemoteIp = strings.Split(d.Get("backup_remote_tunnel_cidr").(string), ",")
		}
		if newKey.(string) != "" && len(newKeyList) != len(bgpRemoteIp) {
			return fmt.Errorf("can't update backup BGP MD5 authentication key since it is not set correctly for BGP connection: %s", d.Get("connection_name").(string))
		}
		for i, v := range bgpRemoteIp {
			bgpMd5Key := ""
			if newKey.(string) != "" {
				bgpMd5Key = newKeyList[i]
			}
			if newKey.(string) != "" && oldKey.(string) != "" && strings.TrimSpace(newKeyList[i]) == strings.TrimSpace(oldKeyList[i]) {
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
				return fmt.Errorf("failed to update backup BGP MD5 authentication key: %v", err)
			}
		}
	}

	if d.HasChange("phase1_local_identifier") {
		s2c := &goaviatrix.EditSite2Cloud{
			VpcID:                 d.Get("vpc_id").(string),
			ConnName:              d.Get("connection_name").(string),
			Phase1LocalIdentifier: d.Get("phase1_local_identifier").(string),
		}
		err := client.EditSite2CloudPhase1LocalIdentifier(s2c)
		if err != nil {
			return fmt.Errorf("could not update phase1 local identificer for connection: %s: %v", s2c.ConnName, err)
		}
	}

	d.Partial(false)

	return resourceAviatrixTransitExternalDeviceConnRead(d, meta)
}

func resourceAviatrixTransitExternalDeviceConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:              d.Get("vpc_id").(string),
		ConnectionName:     d.Get("connection_name").(string),
		EnableEdgeUnderlay: d.Get("enable_edge_underlay").(bool),
	}

	log.Printf("[INFO] Deleting Aviatrix external device connection: %#v", externalDeviceConn)
	if externalDeviceConn.EnableEdgeUnderlay {
		edgeExternalDeviceConn := goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)
		if val, ok := d.Get("local_lan_ip").(string); ok {
			edgeExternalDeviceConn.LocalLanIP = val
		} else {
			return fmt.Errorf("failed to get 'local_lan_ip' as string")
		}
		if gwName, ok := d.Get("gw_name").(string); ok {
			edgeExternalDeviceConn.GwName = gwName
		} else {
			return fmt.Errorf("failed to get 'gw_name' as string")
		}
		err := client.DeleteEdgeExternalDeviceConn(&edgeExternalDeviceConn)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix external device connection: %w", err)
		}
	} else {
		err := client.DeleteExternalDeviceConn(externalDeviceConn)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix external device connection: %w", err)
		}
	}

	return nil
}
