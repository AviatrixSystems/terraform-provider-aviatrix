package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitExternalDeviceConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitExternalDeviceConnCreate,
		Read:   resourceAviatrixTransitExternalDeviceConnRead,
		Delete: resourceAviatrixTransitExternalDeviceConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC-ID where the Transit Gateway is located.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of for Transit GW to VGW connection connection which is going to be created.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Transit Gateway.",
			},
			"remote_gateway_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an AWS TGW.",
			},
			"connection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "bgp",
				ForceNew:    true,
				Description: "Connection type. Valid values: 'bpg', 'static'. Default value: 'bgp'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "bgp" && v != "static" {
						errs = append(errs, fmt.Errorf("%q must be either 'bgp' or 'static', got: %s", key, val))
					}
					return
				},
			},
			"bgp_local_as_number": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "AWS side as a number. Integer between 1-65535. Example: '12'. Required for a dynamic VPN connection.",
			},
			"bgp_remote_as_number": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"remote_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"direct_connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"local_tunnel_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"remote_tunnel_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
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
			},
			"phase_2_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase two Authentication. Valid values: 'NO-AUTH', 'HMAC-SHA-1', 'HMAC-SHA-256', " +
					"'HMAC-SHA-384' and 'HMAC-SHA-512'.",
			},
			"phase_1_dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Phase one DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'.",
			},
			"phase_2_dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Phase two DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'.",
			},
			"phase_1_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and " +
					"'AES-256-CBC'.",
			},
			"phase_2_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', " +
					"'AES-256-CBC', 'AES-128-GCM-64', 'AES-128-GCM-96' and 'AES-128-GCM-128'.",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_remote_gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_bgp_remote_as_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_local_tunnel_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_remote_tunnel_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_direct_connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "ID of the vpn connection.",
			},
			"enable_edge_segmentation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "ID of the vpn connection.",
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
	}

	bgpLocalAsNumber, err := strconv.Atoi(d.Get("bgp_local_as_number").(string))
	if err == nil {
		externalDeviceConn.BgpLocalAsNumber = bgpLocalAsNumber
	}
	bgpRemoteAsNumber, err := strconv.Atoi(d.Get("bgp_remote_as_number").(string))
	if err == nil {
		externalDeviceConn.BgpRemoteAsNumber = bgpRemoteAsNumber
	}
	backupBgpLocalAsNumber, err := strconv.Atoi(d.Get("backup_bgp_remote_as_number").(string))
	if err == nil {
		externalDeviceConn.BackupBgpRemoteAsNumber = backupBgpLocalAsNumber
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
		return fmt.Errorf("'remote_subnet' is needed for connection type of 'static' not 'bpg'")
	} else if externalDeviceConn.ConnectionType == "static" && (externalDeviceConn.BgpLocalAsNumber != 0 || externalDeviceConn.BgpRemoteAsNumber != 0) {
		return fmt.Errorf("'bgp_local_as_number' and 'bgp_remote_as_number' are needed for connection type of 'bgp' not 'static'")
	}

	customAlgorithms := d.Get("custom_algorithms").(bool)
	if !customAlgorithms && (externalDeviceConn.Phase1Auth != "" || externalDeviceConn.Phase1DhGroups != "" ||
		externalDeviceConn.Phase1Encryption != "" || externalDeviceConn.Phase2Auth != "" ||
		externalDeviceConn.Phase2DhGroups != "" || externalDeviceConn.Phase2Encryption != "") {
		return fmt.Errorf("'custom_algorithms' is not enabled, all algorithms fields should be left empty")
	}

	if haEnabled {
		if externalDeviceConn.BackupRemoteGatewayIP == "" {
			return fmt.Errorf("ha is enabled, please specify 'backup_remote_gateway_ip'")
		}
		if externalDeviceConn.BackupBgpRemoteAsNumber == 0 && externalDeviceConn.ConnectionType == "bgp" {
			return fmt.Errorf("ha is enabled, and 'connection_type' is 'bgp', please specify 'backup_bgp_remote_as_number'")
		}
	} else {
		if backupDirectConnect {
			return fmt.Errorf("ha is not enabled, please set 'back_direct_connect' to false")
		}
		if externalDeviceConn.BackupPreSharedKey != "" || externalDeviceConn.BackupLocalTunnelCidr != "" || externalDeviceConn.BackupRemoteTunnelCidr != "" {
			return fmt.Errorf("ha is not enabled, please set 'backup_pre_shared_key', 'backup_local_tunnel_cidr' and 'backup_remote_tunnel_cidr' to empty")
		}
		if externalDeviceConn.BackupBgpRemoteAsNumber != 0 && externalDeviceConn.ConnectionType == "bgp" {
			return fmt.Errorf("ha is not enabled, and 'connection_type' is 'bgp', please specify 'backup_bgp_remote_as_number' to empty")
		}
	}

	err = client.CreateExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix external device connection: %s", err)
	}

	d.SetId(externalDeviceConn.ConnectionName + "~" + externalDeviceConn.VpcID)
	return resourceAviatrixTransitExternalDeviceConnRead(d, meta)
}

func resourceAviatrixTransitExternalDeviceConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	vpcID := d.Get("vpc_id").(string)
	if connectionName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no 'connection_name' or 'vpc_id' received. Import Id is %s", id)
		d.Set("connection_name", strings.Split(id, "~")[0])
		d.Set("vpc_id", strings.Split(id, "~")[1])
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
		d.Set("remote_gateway_ip", conn.RemoteGatewayIP)
		d.Set("connection_type", conn.ConnectionType)
		if conn.BgpLocalAsNumber != 0 {
			d.Set("bgp_local_as_number", strconv.Itoa(conn.BgpLocalAsNumber))
		}
		if conn.BgpLocalAsNumber != 0 {
			d.Set("bgp_remote_as_number", strconv.Itoa(conn.BgpRemoteAsNumber))
		}
		d.Set("remote_subnet", conn.RemoteSubnet)
		if conn.DirectConnect == "enabled" {
			d.Set("direct_connect", true)
		} else {
			d.Set("direct_connect", false)
		}

		d.Set("local_tunnel_cidr", conn.LocalTunnelCidr)
		d.Set("remote_tunnel_cidr", conn.RemoteTunnelCidr)
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

		d.Set("backup_remote_gateway_ip", conn.BackupRemoteGatewayIP)
		if conn.BackupBgpRemoteAsNumber != 0 {
			d.Set("backup_bgp_remote_as_number", strconv.Itoa(conn.BackupBgpRemoteAsNumber))
		}
		d.Set("backup_local_tunnel_cidr", conn.BackupLocalTunnelCidr)
		d.Set("backup_remote_tunnel_cidr", conn.BackupRemoteTunnelCidr)
		if conn.BackupDirectConnect == "enabled" {
			d.Set("backup_direct_connect", true)
		} else {
			d.Set("backup_direct_connect", false)
		}

		if conn.EnableEdgeSegmentation == "enabled" {
			d.Set("enable_edge_segmentation", true)
		} else {
			d.Set("enable_edge_segmentation", false)
		}
		if conn.HAEnabled == "enabled" {
			d.Set("ha_enabled", true)
		} else {
			d.Set("ha_enabled", false)
		}
	}

	d.SetId(conn.ConnectionName + "~" + conn.VpcID)
	return nil
}

func resourceAviatrixTransitExternalDeviceConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:          d.Get("vpc_id").(string),
		ConnectionName: d.Get("connection_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix external device connection: %#v", externalDeviceConn)

	err := client.DeleteExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix external device connection: %s", err)
	}

	return nil
}
