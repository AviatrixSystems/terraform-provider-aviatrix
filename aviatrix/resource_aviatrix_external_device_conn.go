package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixExternalDeviceConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixExternalDeviceConnCreate,
		Read:   resourceAviatrixExternalDeviceConnRead,
		Delete: resourceAviatrixExternalDeviceConnDelete,
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
			"conn_name": {
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
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "AWS side as a number. Integer between 1-65535. Example: '12'. Required for a dynamic VPN connection.",
			},
			"bgp_remote_as_number": {
				Type:        schema.TypeInt,
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
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"local_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"remote_tunnel_ip": {
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
			"enable_ha": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_remote_gateway_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_bgp_remote_as_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_pre_shared_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_local_tunnel_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_remote_tunnel_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"backup_direct_connect": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
			"enable_edge_segmentation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
		},
	}
}

func resourceAviatrixExternalDeviceConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("zjin00 nothing is wrong here")

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:                   d.Get("vpc_id").(string),
		ConnName:                d.Get("conn_name").(string),
		GwName:                  d.Get("gw_name").(string),
		ConnType:                d.Get("connection_type").(string),
		BgpLocalAsNumber:        d.Get("bgp_local_as_number").(int),
		BgpRemoteAsNumber:       d.Get("bgp_remote_as_number").(int),
		RemoteGatewayIP:         d.Get("remote_gateway_ip").(string),
		RemoteSubnet:            d.Get("remote_subnet").(string),
		PreSharedKey:            d.Get("pre_shared_key").(string),
		LocalTunnelIP:           d.Get("local_tunnel_ip").(string),
		RemoteTunnelIP:          d.Get("remote_tunnel_ip").(string),
		Phase1Auth:              d.Get("phase_1_authentication").(string),
		Phase1DhGroups:          d.Get("phase_1_dh_groups").(string),
		Phase1Encryption:        d.Get("phase_1_encryption").(string),
		Phase2Auth:              d.Get("phase_2_authentication").(string),
		Phase2DhGroups:          d.Get("phase_2_dh_groups").(string),
		Phase2Encryption:        d.Get("phase_2_encryption").(string),
		BackupRemoteGatewayIP:   d.Get("backup_remote_gateway_ip").(string),
		BackupBgpRemoteAsNumber: d.Get("backup_bgp_remote_as_number").(int),
		BackupPreSharedKey:      d.Get("backup_pre_shared_key").(string),
		BackupLocalTunnelIP:     d.Get("backup_local_tunnel_ip").(string),
		BackupRemoteTunnelIP:    d.Get("backup_remote_tunnel_ip").(string),
	}

	log.Printf("zjin01 nothing is wrong here")

	directConnect := d.Get("direct_connect").(bool)
	if directConnect {
		externalDeviceConn.DirectConnect = "true"
	}
	log.Printf("zjin020 nothing is wrong here")

	haEnabled := d.Get("enable_ha").(bool)
	if haEnabled {
		externalDeviceConn.HAEnabled = "true"
	}
	log.Printf("zjin021 nothing is wrong here")

	backupDirectConnect := d.Get("backup_direct_connect").(bool)
	if backupDirectConnect {
		externalDeviceConn.BackupDirectConnect = "true"
	}
	log.Printf("zjin022 nothing is wrong here")

	enableEdgeSegmentation := d.Get("enable_edge_segmentation").(bool)
	if enableEdgeSegmentation {
		externalDeviceConn.EnableEdgeSegmentation = "true"
	}
	log.Printf("zjin023 nothing is wrong here")

	if externalDeviceConn.ConnType == "bgp" && externalDeviceConn.RemoteSubnet != "" {
		return fmt.Errorf("'remote_subnet' is needed for connection type of 'static' not 'bpg'")
	} else if externalDeviceConn.ConnType == "static" && (externalDeviceConn.BgpLocalAsNumber != 0 || externalDeviceConn.BgpRemoteAsNumber != 0) {
		return fmt.Errorf("'bgp_local_as_number' and 'bgp_remote_as_number' are needed for connection type of 'bgp' not 'static'")
	}
	log.Printf("zjin024 nothing is wrong here")

	customAlgorithms := d.Get("custom_algorithms").(bool)
	if !customAlgorithms && (externalDeviceConn.Phase1Auth != "" || externalDeviceConn.Phase1DhGroups != "" ||
		externalDeviceConn.Phase1Encryption != "" || externalDeviceConn.Phase2Auth != "" ||
		externalDeviceConn.Phase2DhGroups != "" || externalDeviceConn.Phase2Encryption != "") {
		return fmt.Errorf("'custom_algorithms' is not enabled, all algorithms fields should be left empty")
	}

	log.Printf("zjin025 nothing is wrong here")

	if haEnabled {
		if externalDeviceConn.BackupRemoteGatewayIP == "" || externalDeviceConn.BackupBgpRemoteAsNumber == 0 {
			return fmt.Errorf("ha is enabled, please specify 'backup_remote_gateway_ip' and 'backup_bgp_remote_as_number'")
		}
	} else {
		if backupDirectConnect {
			return fmt.Errorf("ha is not enabled, please set 'back_direct_connect' to false")
		}
		if externalDeviceConn.BackupPreSharedKey != "" || externalDeviceConn.BackupLocalTunnelIP != "" || externalDeviceConn.BackupRemoteTunnelIP != "" {
			return fmt.Errorf("ha is not enabled, please set 'backup_pre_shared_key', 'backup_local_tunnel_ip' and 'backup_remote_tunnel_ip' to empty")
		}
	}

	log.Printf("zjin03 nothing is wrong here")

	err := client.CreateExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix external device connection: %s", err)
	}

	d.SetId(externalDeviceConn.ConnName + "~" + externalDeviceConn.VpcID)
	return resourceAviatrixExternalDeviceConnRead(d, meta)
}

func resourceAviatrixExternalDeviceConnRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*goaviatrix.Client)
	//
	//connName := d.Get("conn_name").(string)
	//vpcID := d.Get("vpc_id").(string)
	//if connName == "" || vpcID == "" {
	//	id := d.Id()
	//	log.Printf("[DEBUG] Looks like an import, no 'conn_name' or 'vpc_id' received. Import Id is %s", id)
	//	d.Set("conn_name", strings.Split(id, "~")[0])
	//	d.Set("vpc_id", strings.Split(id, "~")[1])
	//	d.SetId(id)
	//}

	//externalDeviceConn := &goaviatrix.ExternalDeviceConn{
	//	VpcID:    d.Get("vpc_id").(string),
	//	ConnName: d.Get("conn_name").(string),
	//}
	//conn, err := client.GetSite2CloudConnDetail(externalDeviceConn)
	//if err != nil {
	//	if err == goaviatrix.ErrNotFound {
	//		d.SetId("")
	//		return nil
	//	}
	//	return fmt.Errorf("couldn't find Aviatrix Site2Cloud: %s, %#v", err, s2c)
	//}
	//
	//if s2c != nil {
	//	d.Set("vpc_id", s2c.VpcID)
	//	d.Set("remote_gateway_type", s2c.RemoteGwType)
	//	d.Set("tunnel_type", s2c.TunnelType)
	//	d.Set("local_subnet_cidr", s2c.LocalSubnet)
	//	d.Set("remote_subnet_cidr", s2c.RemoteSubnet)
	//	if s2c.HAEnabled == "enabled" {
	//		d.Set("ha_enabled", true)
	//	} else {
	//		d.Set("ha_enabled", false)
	//	}
	//
	//	if s2c.HAEnabled == "enabled" {
	//		d.Set("remote_gateway_ip", s2c.RemoteGwIP)
	//		d.Set("backup_remote_gateway_ip", s2c.RemoteGwIP2)
	//		d.Set("primary_cloud_gateway_name", s2c.GwName)
	//		d.Set("backup_gateway_name", s2c.BackupGwName)
	//	} else {
	//		d.Set("remote_gateway_ip", s2c.RemoteGwIP)
	//		d.Set("primary_cloud_gateway_name", s2c.GwName)
	//	}
	//
	//	d.Set("connection_type", s2c.ConnType)
	//	if s2c.ConnType == "mapped" {
	//		d.Set("remote_subnet_virtual", s2c.RemoteSubnetVirtual)
	//		d.Set("local_subnet_virtual", s2c.LocalSubnetVirtual)
	//	}
	//
	//	if s2c.CustomAlgorithms {
	//		d.Set("custom_algorithms", true)
	//		d.Set("phase_1_authentication", s2c.Phase1Auth)
	//		d.Set("phase_2_authentication", s2c.Phase2Auth)
	//		d.Set("phase_1_dh_groups", s2c.Phase1DhGroups)
	//		d.Set("phase_2_dh_groups", s2c.Phase2DhGroups)
	//		d.Set("phase_1_encryption", s2c.Phase1Encryption)
	//		d.Set("phase_2_encryption", s2c.Phase2Encryption)
	//	} else {
	//		d.Set("custom_algorithms", false)
	//	}
	//
	//	if s2c.PrivateRouteEncryption == "true" {
	//		d.Set("private_route_encryption", true)
	//
	//		if err := d.Set("route_table_list", s2c.RouteTableList); err != nil {
	//			log.Printf("[WARN] Error setting route_table_list for (%s): %s", d.Id(), err)
	//		}
	//	} else {
	//		d.Set("private_route_encryption", false)
	//	}
	//
	//	if s2c.SslServerPool != "" {
	//		d.Set("ssl_server_pool", s2c.SslServerPool)
	//	}
	//
	//	d.Set("enable_dead_peer_detection", s2c.DeadPeerDetection)
	//	d.Set("enable_active_active", s2c.EnableActiveActive)
	//}
	//
	//log.Printf("[TRACE] Reading Aviatrix Site2Cloud %s: %#v", d.Get("connection_name").(string), site2cloud)
	//log.Printf("[TRACE] Reading Aviatrix Site2Cloud connection_type: [%s]", d.Get("connection_type").(string))
	//
	//d.SetId(site2cloud.TunnelName + "~" + site2cloud.VpcID)
	//return nil
	return nil
}

func resourceAviatrixExternalDeviceConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	externalDeviceConn := &goaviatrix.ExternalDeviceConn{
		VpcID:    d.Get("vpc_id").(string),
		ConnName: d.Get("conn_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix external device connection: %#v", externalDeviceConn)

	err := client.DeleteExternalDeviceConn(externalDeviceConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix external device connection: %s", err)
	}

	return nil
}
