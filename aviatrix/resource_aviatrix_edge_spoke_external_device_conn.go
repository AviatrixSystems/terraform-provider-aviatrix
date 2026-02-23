package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixEdgeSpokeExternalDeviceConn() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeSpokeExternalDeviceConnCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeSpokeExternalDeviceConnRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeSpokeExternalDeviceConnUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeSpokeExternalDeviceConnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC where the BGP Spoke Gateway is located.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of the spoke external device connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the BGP Spoke Gateway.",
			},
			"bgp_local_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP local AS number.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"bgp_remote_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP remote AS number.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"local_lan_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Local LAN IP.",
			},
			"remote_lan_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Remote LAN IP.",
			},
			"connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bgp",
				ForceNew:     true,
				Description:  "Connection type. Valid values: 'bgp'. Default value: 'bgp'.",
				ValidateFunc: validation.StringInSlice([]string{"bgp"}, false),
			},
			"tunnel_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "LAN",
				ForceNew:     true,
				Description:  "Tunnel Protocol. Valid value: 'LAN'. Default value: 'LAN'. Case insensitive.",
				ValidateFunc: validation.StringInSlice([]string{"LAN"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"enable_bgp_lan_activemesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Switch to enable BGP LAN ActiveMesh. Valid values: true, false. Default: false.",
			},
			"number_of_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of retries.",
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Retry interval in seconds.",
			},
			"enable_edge_underlay": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable BGP over WAN underlay.",
			},
			"remote_cloud_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE"}, false),
				Description:  "Remote cloud type.",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Set as true if there are two external devices.",
			},
			"backup_bgp_remote_as_num": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				Description:  "Backup BGP remote ASN (Autonomous System Number). Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"backup_local_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup Local LAN IP.",
			},
			"backup_remote_lan_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup Remote LAN IP.",
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
			"manual_bgp_advertised_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional:    true,
				Description: "Configure manual BGP advertised CIDRs for this connection.",
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
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Jumbo Frame for the edge spoke external device connection. Valid values: true, false.",
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable IPv6 on this connection",
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

func marshalEdgeSpokeExternalDeviceConnInput(d *schema.ResourceData) (*goaviatrix.ExternalDeviceConn, error) {
	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:               getString(d, "site_id"),
		ConnectionName:      getString(d, "connection_name"),
		GwName:              getString(d, "gw_name"),
		ConnectionType:      getString(d, "connection_type"),
		TunnelProtocol:      strings.ToUpper(getString(d, "tunnel_protocol")),
		LocalLanIP:          getString(d, "local_lan_ip"),
		RemoteLanIP:         getString(d, "remote_lan_ip"),
		EnableEdgeUnderlay:  getBool(d, "enable_edge_underlay"),
		RemoteCloudType:     getString(d, "remote_cloud_type"),
		BackupLocalLanIP:    getString(d, "backup_local_lan_ip"),
		BackupRemoteLanIP:   getString(d, "backup_remote_lan_ip"),
		BgpMd5Key:           getString(d, "bgp_md5_key"),
		BackupBgpMd5Key:     getString(d, "backup_bgp_md5_key"),
		EnableBgpMultihop:   getBool(d, "enable_bgp_multihop"),
		EnableJumboFrame:    getBool(d, "enable_jumbo_frame"),
		EnableIpv6:          getBool(d, "enable_ipv6"),
		RemoteLanIPv6:       getString(d, "remote_lan_ipv6_ip"),
		BackupRemoteLanIPv6: getString(d, "backup_remote_lan_ipv6_ip"),
	}

	haEnabled := getBool(d, "ha_enabled")
	if haEnabled {
		externalDeviceConn.HAEnabled = "true"
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

	enableBgpLanActiveMesh := getBool(d, "enable_bgp_lan_activemesh")

	if enableBgpLanActiveMesh &&
		(externalDeviceConn.ConnectionType != "bgp" || externalDeviceConn.TunnelProtocol != "LAN") {
		return nil, fmt.Errorf("'enable_bgp_lan_activemesh' only supports 'bgp' connection " +
			"with 'LAN' tunnel protocol")
	}
	externalDeviceConn.EnableBgpLanActiveMesh = enableBgpLanActiveMesh

	return externalDeviceConn, nil
}

func resourceAviatrixEdgeSpokeExternalDeviceConnCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	externalDeviceConn, err := marshalEdgeSpokeExternalDeviceConnInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if !externalDeviceConn.EnableEdgeUnderlay && externalDeviceConn.ConnectionName == "" {
		return diag.Errorf("'connection_name' is required when 'enable_edge_underlay' is false")
	}

	if externalDeviceConn.EnableEdgeUnderlay && externalDeviceConn.ConnectionName != "" {
		return diag.Errorf("please set 'connection_name' to empty when 'enable_edge_underlay' is true")
	}

	if externalDeviceConn.HAEnabled == "true" {
		if externalDeviceConn.BackupRemoteLanIP == "" {
			return diag.Errorf("ha is enabled and 'tunnel_protocol' = 'LAN', please specify 'backup_remote_lan_ip'")
		}

		if externalDeviceConn.BackupBgpRemoteAsNum == 0 {
			return diag.Errorf("ha is enabled, and 'connection_type' is 'bgp', please specify 'backup_bgp_remote_as_num'")
		}
	} else {
		if externalDeviceConn.BackupRemoteLanIP != "" || externalDeviceConn.BackupLocalLanIP != "" {
			return diag.Errorf("ha is not enabled, please set 'backup_remote_lan_ip' and 'backup_local_lan_ip' to empty")
		}
		if externalDeviceConn.BackupBgpRemoteAsNum != 0 {
			return diag.Errorf("ha is not enabled, and 'connection_type' is 'bgp', please specify 'backup_bgp_remote_as_num' to empty")
		}
		if externalDeviceConn.BackupRemoteLanIPv6 != "" {
			return diag.Errorf("ha is not enabled, please set 'backup_remote_lan_ipv6_ip' to empty")
		}
	}

	if externalDeviceConn.EnableIpv6 {
		if externalDeviceConn.RemoteLanIPv6 == "" {
			return diag.Errorf("'remote_lan_ipv6_ip' is required when 'enable_ipv6' is true")
		}
		if externalDeviceConn.HAEnabled == "true" && externalDeviceConn.BackupRemoteLanIPv6 == "" {
			return diag.Errorf("'backup_remote_lan_ipv6_ip' is required when 'enable_ipv6' is true and ha is enabled")
		}
	} else {
		if externalDeviceConn.RemoteLanIPv6 != "" {
			return diag.Errorf("IPv6 is not enabled, please set 'remote_lan_ipv6_ip' to empty")
		}
		if externalDeviceConn.BackupRemoteLanIPv6 != "" {
			return diag.Errorf("IPv6 is not enabled, please set 'backup_remote_lan_ipv6_ip' to empty")
		}
	}

	flag := false
	defer resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx, d, meta, &flag)

	numberOfRetries := getInt(d, "number_of_retries")
	retryInterval := getInt(d, "retry_interval")

	var edgeExternalDeviceConn goaviatrix.EdgeExternalDeviceConn
	if externalDeviceConn.EnableEdgeUnderlay {
		edgeExternalDeviceConn = goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)
	}

	var connName string
	for i := 0; ; i++ {
		if externalDeviceConn.EnableEdgeUnderlay {
			connName, err = client.CreateEdgeExternalDeviceConn(&edgeExternalDeviceConn)
			log.Printf("[DEBUG] Created underlay connection %s", connName)
		} else {
			err = client.CreateExternalDeviceConn(externalDeviceConn)
		}

		if err != nil {
			if !strings.Contains(err.Error(), "not ready") && !strings.Contains(err.Error(), "not up") {
				return diag.Errorf("failed to create Edge as a Spoke external device connection: %s", err)
			}
		} else {
			break
		}
		if i < numberOfRetries {
			time.Sleep(time.Duration(retryInterval) * time.Second)
		} else {
			d.SetId("")
			return diag.Errorf("failed to create Edge as a Spoke external device connection: %s", err)
		}
	}

	enableBFD := getBool(d, "enable_bfd")

	externalDeviceConn.EnableBfd = enableBFD
	bgp_bfd := getList(d, "bgp_bfd")

	// set the bgp bfd config details only if the user has enabled BFD
	if enableBFD {
		// set bgp bfd using the config details provided by the user
		if len(bgp_bfd) > 0 {
			for _, bfd0 := range bgp_bfd {
				bfd1, ok := bfd0.(map[string]interface{})
				if !ok {
					return diag.Errorf("expected bgp_bfd to be a map, but got %T", bfd0)
				}
				externalDeviceConn.BgpBfdConfig = goaviatrix.CreateBgpBfdConfig(bfd1)
			}
		} else {
			// set the bgp bfd config using the default values
			externalDeviceConn.BgpBfdConfig = defaultBfdConfig
		}
		err := client.EditConnectionBgpBfd(externalDeviceConn)
		if err != nil {
			return diag.Errorf("could not update BGP BFD config: %v", err)
		}
	} else {
		// if BFD is disabled and BGP BFD config is provided then throw an error
		if len(bgp_bfd) > 0 {
			return diag.Errorf("bgp_bfd config can't be set when BFD is disabled")
		}
	}

	if externalDeviceConn.EnableEdgeUnderlay {
		externalDeviceConn.ConnectionName = connName
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}

		err = client.EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return diag.Errorf("could not set prepend_as_path: %v", err)
		}
	}

	manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
	if len(manualBGPCidrs) > 0 {
		err = client.EditTransitConnectionBGPManualAdvertiseCIDRs(externalDeviceConn.GwName, externalDeviceConn.ConnectionName, manualBGPCidrs)
		if err != nil {
			return diag.Errorf("could not edit manual advertised BGP cidrs: %v", err)
		}
	}

	if !externalDeviceConn.EnableBgpMultihop && externalDeviceConn.ConnectionType != "bgp" {
		return diag.Errorf("multihop can only be configured for BGP connections")
	}

	d.SetId(externalDeviceConn.ConnectionName + "~" + externalDeviceConn.VpcID + "~" + externalDeviceConn.GwName)
	return resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeSpokeExternalDeviceConnRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeSpokeExternalDeviceConnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	connectionName := getString(d, "connection_name")

	vpcID := getString(d, "site_id")
	if connectionName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'site_id' received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 3 {
			return diag.Errorf("expected import ID in the form 'connection_name~site_id~gw_name' instead got %q", id)
		}
		mustSet(d, "connection_name", parts[0])
		mustSet(d, "site_id", parts[1])
		mustSet(d, "gw_name", parts[2])
		d.SetId(id)
	}

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          getString(d, "site_id"),
		ConnectionName: getString(d, "connection_name"),
		GwName:         getString(d, "gw_name"),
	}

	localGateway, err := getGatewayDetails(client, externalDeviceConn.GwName)
	if err != nil {
		return diag.Errorf("could not get local gateway details: %v", err)
	}

	conn, err := client.GetExternalDeviceConnDetail(externalDeviceConn, localGateway)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't find Edge as a Spoke external device conn: %s, %#v", err, externalDeviceConn)
	}

	// If we are dealing with a HAGW, we're interested in the BGP backup
	// parameters.
	inputGwName := externalDeviceConn.GwName
	if conn.EnableEdgeUnderlay && strings.Contains(inputGwName, "hagw") {
		conn.GwName = inputGwName
		conn.BgpRemoteAsNum = conn.BackupBgpRemoteAsNum
		conn.LocalLanIP = conn.BackupLocalLanIP
		conn.RemoteLanIP = conn.BackupRemoteLanIP
	}
	mustSet(d, "site_id", conn.VpcID)
	mustSet(d, "connection_name", conn.ConnectionName)
	mustSet(d, "gw_name", conn.GwName)
	mustSet(d, "connection_type", conn.ConnectionType)
	mustSet(d, "tunnel_protocol", conn.TunnelProtocol)
	mustSet(d, "local_lan_ip", conn.LocalLanIP)
	mustSet(d, "remote_lan_ip", conn.RemoteLanIP)
	mustSet(d, "enable_edge_underlay", conn.EnableEdgeUnderlay)
	mustSet(d, "remote_cloud_type", conn.RemoteCloudType)
	mustSet(d, "enable_ipv6", conn.EnableIpv6)
	if conn.EnableIpv6 {
		mustSet(d, "remote_lan_ipv6_ip", conn.RemoteLanIPv6)
	}

	_ = d.Set("enable_bgp_lan_activemesh", conn.EnableBgpLanActiveMesh)
	if conn.BgpLocalAsNum != 0 {
		mustSet(d, "bgp_local_as_num", strconv.Itoa(conn.BgpLocalAsNum))
	}
	if conn.BgpRemoteAsNum != 0 {
		mustSet(d, "bgp_remote_as_num", strconv.Itoa(conn.BgpRemoteAsNum))
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

	if conn.HAEnabled == "enabled" {
		if !conn.EnableEdgeUnderlay {
			mustSet(d, "ha_enabled", true)
			if conn.BackupBgpRemoteAsNum != 0 {
				mustSet(d, "backup_bgp_remote_as_num", strconv.Itoa(conn.BackupBgpRemoteAsNum))
			}
			mustSet(d, "backup_remote_lan_ip", conn.BackupRemoteLanIP)
			mustSet(d, "backup_local_lan_ip", conn.BackupLocalLanIP)
			if conn.EnableIpv6 {
				mustSet(d, "backup_remote_lan_ipv6_ip", conn.BackupRemoteLanIPv6)
			}
		}
	} else {
		mustSet(d, "ha_enabled", false)
	}

	if conn.PrependAsPath != "" {
		var prependAsPath []string
		for _, str := range strings.Split(conn.PrependAsPath, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path", prependAsPath)
		if err != nil {
			return diag.Errorf("could not set value for prepend_as_path: %v", err)
		}
	}

	err = d.Set("enable_bgp_multihop", conn.EnableBgpMultihop)
	if err != nil {
		return diag.Errorf("could not set value for enable_bgp_multihop: %v", err)
	}

	if err := d.Set("manual_bgp_advertised_cidrs", conn.ManualBGPCidrs); err != nil {
		return diag.Errorf("could not set value for manual_bgp_advertised_cidrs: %v", err)
	}
	if err := d.Set("enable_jumbo_frame", conn.EnableJumboFrame); err != nil {
		return diag.Errorf("could not set value for enable_jumbo_frame: %v", err)
	}
	d.SetId(conn.ConnectionName + "~" + conn.VpcID + "~" + conn.GwName)
	return nil
}

func resourceAviatrixEdgeSpokeExternalDeviceConnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	d.Partial(true)

	externalDeviceConn, err := marshalEdgeSpokeExternalDeviceConnInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("prepend_as_path") {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}
		err := client.EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return diag.Errorf("could not update prepend_as_path: %v", err)
		}
	}

	enableBfd := getBool(d, "enable_bfd")

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
				return diag.Errorf("could not update BGP BFD config: %v", err)
			}
		} else {
			// bgp bfd is disabled
			if len(bgpBfdConfig) > 0 {
				return diag.Errorf("bgp_bfd config can't be set when BFD is disabled")
			}
			externalDeviceConn := &goaviatrix.ExternalDeviceConn{
				GwName:         getString(d, "gw_name"),
				ConnectionName: getString(d, "connection_name"),
				EnableBfd:      enableBfd,
			}
			err := client.EditConnectionBgpBfd(externalDeviceConn)
			if err != nil {
				return diag.Errorf("could not disable BGP BFD config: %v", err)
			}
		}
	}

	if d.HasChanges("enable_bgp_multihop") {
		if err := client.EditConnectionBgpMultihop(externalDeviceConn); err != nil {
			return diag.Errorf("could not update BGP BFD config: %v", err)
		}
	}

	if externalDeviceConn.EnableEdgeUnderlay && d.HasChanges("bgp_md5_key", "backup_bgp_md5_key") {
		edgeExternalDeviceConn := goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)

		bgpMD5Key := getString(d, "bgp_md5_key")

		edgeExternalDeviceConn.BgpMd5Key = bgpMD5Key

		backupBGPMD5Key := getString(d, "backup_bgp_md5_key")

		edgeExternalDeviceConn.BackupBgpMd5Key = backupBGPMD5Key

		edgeExternalDeviceConn.BgpMd5KeyChanged = true

		_, err := client.CreateEdgeExternalDeviceConn(&edgeExternalDeviceConn)
		if err != nil {
			return diag.Errorf("could not update BGP MD5 key: %s", err)
		}
	}

	if d.HasChange("manual_bgp_advertised_cidrs") {
		manualBGPCidrs := getStringSet(d, "manual_bgp_advertised_cidrs")
		err := client.EditTransitConnectionBGPManualAdvertiseCIDRs(externalDeviceConn.GwName, externalDeviceConn.ConnectionName, manualBGPCidrs)
		if err != nil {
			return diag.Errorf("could not edit manual advertise manual cidrs: %v", err)
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		vpcID := getString(d, "site_id")

		connectionName := getString(d, "connection_name")

		gwName := getString(d, "gw_name")

		connectionType := getString(d, "connection_type")

		tunnelProtocol := getString(d, "tunnel_protocol")

		enableJumboFrame := getBool(d, "enable_jumbo_frame")

		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:            vpcID,
			ConnectionName:   connectionName,
			GwName:           gwName,
			ConnectionType:   connectionType,
			TunnelProtocol:   tunnelProtocol,
			EnableJumboFrame: enableJumboFrame,
		}
		if externalDeviceConn.EnableJumboFrame {
			if externalDeviceConn.ConnectionType != "bgp" {
				return diag.Errorf("jumbo frame is only supported on BGP connection")
			}
			if err := client.EnableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return diag.Errorf("failed to enable jumbo frame for external device connection %q: %v", externalDeviceConn.ConnectionName, err)
			}
		} else {
			if err := client.DisableJumboFrameExternalDeviceConn(externalDeviceConn); err != nil {
				return diag.Errorf("failed to disable jumbo frame for external device connection %q: %v", externalDeviceConn.ConnectionName, err)
			}
		}
	}
	d.Partial(false)
	return resourceAviatrixEdgeSpokeExternalDeviceConnRead(ctx, d, meta)
}

func resourceAviatrixEdgeSpokeExternalDeviceConnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	externalDeviceConn, err := marshalEdgeSpokeExternalDeviceConnInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if externalDeviceConn.EnableEdgeUnderlay {
		edgeExternalDeviceConn := goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)
		err := client.DeleteEdgeExternalDeviceConn(&edgeExternalDeviceConn)
		if err != nil {
			return diag.Errorf("failed to delete Edge as a Spoke external device connection: %s", err)
		}
	} else {
		err := client.DeleteExternalDeviceConn(externalDeviceConn)
		if err != nil {
			return diag.Errorf("failed to delete Edge as a Spoke external device connection: %s", err)
		}
	}

	return nil
}
