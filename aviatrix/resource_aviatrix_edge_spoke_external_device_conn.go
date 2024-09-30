package aviatrix

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
					return strings.ToUpper(old) == strings.ToUpper(new)
				},
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("enable_edge_underlay").(bool)
				},
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "BGP BFD configuration details applied to a BGP session.",
				MaxItems:    1,
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
		},
	}
}

func marshalEdgeSpokeExternalDeviceConnInput(d *schema.ResourceData) *goaviatrix.ExternalDeviceConn {
	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:              d.Get("site_id").(string),
		ConnectionName:     d.Get("connection_name").(string),
		GwName:             d.Get("gw_name").(string),
		ConnectionType:     d.Get("connection_type").(string),
		TunnelProtocol:     strings.ToUpper(d.Get("tunnel_protocol").(string)),
		LocalLanIP:         d.Get("local_lan_ip").(string),
		RemoteLanIP:        d.Get("remote_lan_ip").(string),
		EnableEdgeUnderlay: d.Get("enable_edge_underlay").(bool),
		RemoteCloudType:    d.Get("remote_cloud_type").(string),
		BackupLocalLanIP:   d.Get("backup_local_lan_ip").(string),
		BackupRemoteLanIP:  d.Get("backup_remote_lan_ip").(string),
		BgpMd5Key:          d.Get("bgp_md5_key").(string),
		BackupBgpMd5Key:    d.Get("backup_bgp_md5_key").(string),
	}

	haEnabled := d.Get("ha_enabled").(bool)
	if haEnabled {
		externalDeviceConn.HAEnabled = "true"
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

	return externalDeviceConn
}

func resourceAviatrixEdgeSpokeExternalDeviceConnCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := marshalEdgeSpokeExternalDeviceConnInput(d)

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
	}

	if externalDeviceConn.EnableEdgeUnderlay && externalDeviceConn.HAEnabled == "true" {
		return diag.Errorf("please use a separate edge_spoke_external_device_conn to create WAN underlay connection for Edge HA")
	}

	flag := false
	defer resourceAviatrixEdgeSpokeExternalDeviceConnReadIfRequired(ctx, d, meta, &flag)

	numberOfRetries := d.Get("number_of_retries").(int)
	retryInterval := d.Get("retry_interval").(int)

	var edgeExternalDeviceConn goaviatrix.EdgeExternalDeviceConn
	if externalDeviceConn.EnableEdgeUnderlay {
		edgeExternalDeviceConn = goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)
	}

	var err error
	var result string
	for i := 0; ; i++ {
		if externalDeviceConn.EnableEdgeUnderlay {
			result, err = client.CreateEdgeExternalDeviceConn(&edgeExternalDeviceConn)
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

	enableBFD, ok := d.Get("enable_bfd").(bool)
	if !ok {
		return diag.Errorf("expected enable_bfd to be a boolean, but got %T", d.Get("enable_bfd"))
	}
	externalDeviceConn.EnableBfd = enableBFD
	// set the bgp bfd config details only if the user has enabled BFD
	if enableBFD {
		bgp_bfd, ok := d.Get("bgp_bfd").([]interface{})
		if !ok {
			return diag.Errorf("expected bgp_bfd to be a list of maps, but got %T", d.Get("bgp_bfd"))
		}
		// set bgp bfd using the config details provided by the user
		if len(bgp_bfd) > 0 {
			for _, bfd0 := range bgp_bfd {
				bfd1, ok := bfd0.(map[string]interface{})
				if !ok {
					return diag.Errorf("expected bgp_bfd to be a map, but got %T", bfd0)
				}
				transmitInterval := defaultBfdTransmitInterval
				receiveInterval := defaultBfdReceiveInterval
				multiplier := defaultBfdMultiplier
				if value, ok := bfd1["transmit_interval"].(int); ok {
					transmitInterval = value
				}
				if value, ok := bfd1["receive_interval"].(int); ok {
					receiveInterval = value
				}
				if value, ok := bfd1["multiplier"].(int); ok {
					multiplier = value
				}
				bfd2 := &goaviatrix.BgpBfdConfig{
					TransmitInterval: transmitInterval,
					ReceiveInterval:  receiveInterval,
					Multiplier:       multiplier,
				}
				externalDeviceConn.BgpBfdConfig = append(externalDeviceConn.BgpBfdConfig, bfd2)
			}
		} else {
			// set the bgp bfd config using the default values
			bfd := &goaviatrix.BgpBfdConfig{
				TransmitInterval: defaultBfdTransmitInterval,
				ReceiveInterval:  defaultBfdReceiveInterval,
				Multiplier:       defaultBfdMultiplier,
			}
			externalDeviceConn.BgpBfdConfig = append(externalDeviceConn.BgpBfdConfig, bfd)
		}
		err := client.EditConnectionBgpBfd(externalDeviceConn)
		if err != nil {
			return diag.Errorf("could not update BGP BFD config: %v", err)
		}
	}

	if externalDeviceConn.EnableEdgeUnderlay {
		re := regexp.MustCompile(`underlay BGP connection (.*) (?:in|on)`)
		match := re.FindStringSubmatch(result)
		if len(match) < 2 {
			return diag.Errorf("could not get underlay BGP connection name")
		}
		connName := match[1]
		d.Set("connection_name", connName)
		externalDeviceConn.ConnectionName = connName
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
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

	d.SetId(d.Get("connection_name").(string) + "~" + externalDeviceConn.VpcID + "~" + externalDeviceConn.GwName)
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
	client := meta.(*goaviatrix.Client)

	vpcID := d.Get("site_id").(string)
	if vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'site_id' received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 3 {
			return diag.Errorf("expected import ID in the form 'connection_name~site_id~gw_name' instead got %q", id)
		}
		d.Set("connection_name", parts[0])
		d.Set("site_id", parts[1])
		d.Set("gw_name", parts[2])
		d.SetId(id)
	}

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          d.Get("site_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		GwName:         d.Get("gw_name").(string),
	}

	conn, err := client.GetEdgeExternalDeviceConnDetail(externalDeviceConn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't find Edge as a Spoke external device conn: %s, %#v", err, externalDeviceConn)
	}

	d.Set("site_id", conn.VpcID)
	d.Set("connection_name", conn.ConnectionName)
	d.Set("gw_name", conn.GwName)
	d.Set("connection_type", conn.ConnectionType)
	d.Set("tunnel_protocol", conn.TunnelProtocol)
	d.Set("local_lan_ip", conn.LocalLanIP)
	d.Set("remote_lan_ip", conn.RemoteLanIP)
	d.Set("enable_edge_underlay", conn.EnableEdgeUnderlay)
	d.Set("remote_cloud_type", conn.RemoteCloudType)

	if conn.BgpLocalAsNum != 0 {
		d.Set("bgp_local_as_num", strconv.Itoa(conn.BgpLocalAsNum))
	}
	if conn.BgpRemoteAsNum != 0 {
		d.Set("bgp_remote_as_num", strconv.Itoa(conn.BgpRemoteAsNum))
	}

	enable_bfd, ok := d.Get("enable_bfd").(bool)
	if !ok {
		return diag.Errorf("expected enable_bfd to be a boolean, but got %T", d.Get("enable_bfd"))
	}
	d.Set("enable_bfd", enable_bfd)
	if conn.EnableBfd && len(conn.BgpBfdConfig) > 0 {
		var bgpBfdConfig []map[string]interface{}
		for _, bfd := range conn.BgpBfdConfig {
			bfdMap := make(map[string]interface{})
			bfdMap["transmit_interval"] = defaultBfdTransmitInterval
			bfdMap["receive_interval"] = defaultBfdReceiveInterval
			bfdMap["multiplier"] = defaultBfdMultiplier
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
		}
		d.Set("bgp_bfd", bgpBfdConfig)
	}

	if conn.HAEnabled == "enabled" {
		if !conn.EnableEdgeUnderlay {
			d.Set("ha_enabled", true)
			if conn.BackupBgpRemoteAsNum != 0 {
				d.Set("backup_bgp_remote_as_num", strconv.Itoa(conn.BackupBgpRemoteAsNum))
			}
			d.Set("backup_remote_lan_ip", conn.BackupRemoteLanIP)
			d.Set("backup_local_lan_ip", conn.BackupLocalLanIP)
		}
	} else {
		d.Set("ha_enabled", false)
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

	if err := d.Set("manual_bgp_advertised_cidrs", conn.ManualBGPCidrs); err != nil {
		return diag.Errorf("could not set value for manual_bgp_advertised_cidrs: %v", err)
	}

	d.SetId(conn.ConnectionName + "~" + conn.VpcID + "~" + conn.GwName)
	return nil
}

func resourceAviatrixEdgeSpokeExternalDeviceConnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	d.Partial(true)

	externalDeviceConn := marshalEdgeSpokeExternalDeviceConnInput(d)

	if d.HasChange("prepend_as_path") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		err := client.EditSpokeExternalDeviceConnASPathPrepend(externalDeviceConn, prependASPath)
		if err != nil {
			return diag.Errorf("could not update prepend_as_path: %v", err)
		}
	}

	enableBfd, ok := d.Get("enable_bfd").(bool)
	if !ok {
		return diag.Errorf("expected enable_bfd to be a boolean, but got %T", d.Get("enable_bfd"))
	}
	if enableBfd {
		// get the new BGP BFD config
		bgpBfdConfig, ok := d.Get("bgp_bfd").([]interface{})
		if !ok {
			return diag.Errorf("expected bgp_bfd to be a list of maps, but got %T", d.Get("bgp_bfd"))
		}
		var bgpBfdConfigList []*goaviatrix.BgpBfdConfig
		// Update the BGP BFD config if bfd is enabled and config has changed
		if len(bgpBfdConfig) > 0 && d.HasChange("bgp_bfd") {
			for _, v := range bgpBfdConfig {
				bfdConfig := v.(map[string]interface{})
				transmitInterval, ok := bfdConfig["transmit_interval"].(int)
				if !ok {
					transmitInterval = defaultBfdTransmitInterval
				}
				receiveInterval, ok := bfdConfig["receive_interval"].(int)
				if !ok {
					receiveInterval = defaultBfdReceiveInterval
				}
				multiplier, ok := bfdConfig["multiplier"].(int)
				if !ok {
					multiplier = defaultBfdMultiplier
				}
				bgpBfdConfigList = append(bgpBfdConfigList, &goaviatrix.BgpBfdConfig{
					TransmitInterval: transmitInterval,
					ReceiveInterval:  receiveInterval,
					Multiplier:       multiplier,
				})
			}
		} else {
			// set the bgp bfd config using the default values
			bgpBfdConfigList = append(bgpBfdConfigList, &goaviatrix.BgpBfdConfig{
				TransmitInterval: defaultBfdTransmitInterval,
				ReceiveInterval:  defaultBfdReceiveInterval,
				Multiplier:       defaultBfdMultiplier,
			})
		}
		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			GwName:         d.Get("gw_name").(string),
			ConnectionName: d.Get("connection_name").(string),
			EnableBfd:      d.Get("enable_bfd").(bool),
			BgpBfdConfig:   bgpBfdConfigList,
		}
		err := client.EditConnectionBgpBfd(externalDeviceConn)
		if err != nil {
			return diag.Errorf("could not update BGP BFD config: %v", err)
		}
	} else {
		if d.HasChange("enable_bfd") {
			externalDeviceConn := &goaviatrix.ExternalDeviceConn{
				GwName:         d.Get("gw_name").(string),
				ConnectionName: d.Get("connection_name").(string),
				EnableBfd:      d.Get("enable_bfd").(bool),
			}
			err := client.EditConnectionBgpBfd(externalDeviceConn)
			if err != nil {
				return diag.Errorf("could not disable BGP BFD config: %v", err)
			}
		}
	}

	if externalDeviceConn.EnableEdgeUnderlay && d.HasChanges("bgp_md5_key", "backup_bgp_md5_key") {
		edgeExternalDeviceConn := goaviatrix.EdgeExternalDeviceConn(*externalDeviceConn)

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

	d.Partial(false)
	return resourceAviatrixEdgeSpokeExternalDeviceConnRead(ctx, d, meta)
}

func resourceAviatrixEdgeSpokeExternalDeviceConnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := marshalEdgeSpokeExternalDeviceConnInput(d)

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
