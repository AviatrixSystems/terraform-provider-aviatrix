package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSite2Cloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSite2CloudCreate,
		Read:   resourceAviatrixSite2CloudRead,
		Update: resourceAviatrixSite2CloudUpdate,
		Delete: resourceAviatrixSite2CloudDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixSite2CloudMigrateState,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC Id of the cloud gateway.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site2Cloud Connection Name.",
			},
			"remote_gateway_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Remote gateway type. Valid values: 'generic', 'avx', 'aws', 'azure', 'sonicwall' and 'oracle'.",
			},
			"connection_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connection Type. Valid values: 'mapped' and 'unmapped'.",
			},
			"tunnel_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site2Cloud Tunnel Type. Valid values: 'udp' and 'tcp'",
			},
			"primary_cloud_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Primary Cloud Gateway Name.",
			},
			"remote_gateway_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Remote Gateway IP.",
			},
			"remote_subnet_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote Subnet CIDR.",
			},
			"backup_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup gateway name.",
			},
			"pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Pre-Shared Key.",
			},
			"local_subnet_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Local Subnet CIDR.",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specify whether enabling HA or not.",
			},
			"backup_remote_gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup remote remote gateway IP.",
			},
			"backup_pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Backup Pre-Shared Key.",
			},
			"remote_subnet_virtual": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote Subnet CIDR (Virtual).",
			},
			"local_subnet_virtual": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local Subnet CIDR (Virtual).",
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
			"private_route_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Private route encryption switch.",
			},
			"route_table_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "Route tables to modify.",
			},
			"remote_gateway_latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Latitude of remote gateway.",
			},
			"remote_gateway_longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Longitude of remote gateway.",
			},
			"backup_remote_gateway_latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Latitude of backup remote gateway.",
			},
			"backup_remote_gateway_longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Longitude of backup remote gateway.",
			},
			"ssl_server_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specify ssl_server_pool for tunnel_type 'tcp'. Default value is '192.168.44.0/24'",
			},
			"enable_dead_peer_detection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Switch to Enable/Disable Deed Peer Detection for an existing site2cloud connection.",
			},
		},
	}
}

func resourceAviatrixSite2CloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	s2c := &goaviatrix.Site2Cloud{
		GwName:              d.Get("primary_cloud_gateway_name").(string),
		BackupGwName:        d.Get("backup_gateway_name").(string),
		VpcID:               d.Get("vpc_id").(string),
		TunnelName:          d.Get("connection_name").(string),
		ConnType:            d.Get("connection_type").(string),
		TunnelType:          d.Get("tunnel_type").(string),
		RemoteGwType:        d.Get("remote_gateway_type").(string),
		RemoteGwIP:          d.Get("remote_gateway_ip").(string),
		RemoteGwIP2:         d.Get("backup_remote_gateway_ip").(string),
		PreSharedKey:        d.Get("pre_shared_key").(string),
		BackupPreSharedKey:  d.Get("backup_pre_shared_key").(string),
		RemoteSubnet:        d.Get("remote_subnet_cidr").(string),
		LocalSubnet:         d.Get("local_subnet_cidr").(string),
		RemoteSubnetVirtual: d.Get("remote_subnet_virtual").(string),
		LocalSubnetVirtual:  d.Get("local_subnet_virtual").(string),
	}

	haEnabled := d.Get("ha_enabled").(bool)
	if haEnabled {
		s2c.HAEnabled = "yes"
	} else {
		s2c.HAEnabled = "no"
	}

	if s2c.ConnType != "mapped" && s2c.ConnType != "unmapped" {
		return fmt.Errorf("'connection_type' should be 'mapped' or 'unmapped'")
	}
	if s2c.ConnType == "mapped" && (s2c.RemoteSubnetVirtual == "" || s2c.LocalSubnetVirtual == "") {
		return fmt.Errorf("'remote_subnet_virtual' and 'local_subnet_virtual' are both required for " +
			"connection type: mapped")
	} else if s2c.ConnType == "unmapped" && (s2c.RemoteSubnetVirtual != "" || s2c.LocalSubnetVirtual != "") {
		return fmt.Errorf("'remote_subnet_virtual' and 'local_subnet_virtual' both should be empty for " +
			"connection type: ummapped")
	}

	s2c.Phase1Auth = d.Get("phase_1_authentication").(string)
	s2c.Phase1DhGroups = d.Get("phase_1_dh_groups").(string)
	s2c.Phase1Encryption = d.Get("phase_1_encryption").(string)
	s2c.Phase2Auth = d.Get("phase_2_authentication").(string)
	s2c.Phase2DhGroups = d.Get("phase_2_dh_groups").(string)
	s2c.Phase2Encryption = d.Get("phase_2_encryption").(string)

	customAlgorithms := d.Get("custom_algorithms").(bool)
	if s2c.TunnelType == "tcp" && customAlgorithms {
		return fmt.Errorf("custom_algorithms is not supported for tunnel type 'tcp'")
	}
	if customAlgorithms {
		if s2c.Phase1Auth == goaviatrix.Phase1AuthDefault &&
			s2c.Phase2Auth == goaviatrix.Phase2AuthDefault &&
			s2c.Phase1DhGroups == goaviatrix.Phase1DhGroupDefault &&
			s2c.Phase2DhGroups == goaviatrix.Phase2DhGroupDefault &&
			s2c.Phase1Encryption == goaviatrix.Phase1EncryptionDefault &&
			s2c.Phase2Encryption == goaviatrix.Phase2EncryptionDefault {
			return fmt.Errorf("custom_algorithms is enabled, cannot use default values for " +
				"all six algorithm parameters")
		}
		err := client.Site2CloudAlgorithmCheck(s2c)
		if err != nil {
			return fmt.Errorf("algorithm values check failed: %s", err)
		}
	} else {
		if s2c.Phase1Auth != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_1_authentication should be empty")
		} else {
			s2c.Phase1Auth = goaviatrix.Phase1AuthDefault
		}
		if s2c.Phase1DhGroups != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_1_dh_groups should be empty")
		} else {
			s2c.Phase1DhGroups = goaviatrix.Phase1DhGroupDefault
		}
		if s2c.Phase1Encryption != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_1_encryption should be empty")
		} else {
			s2c.Phase1Encryption = goaviatrix.Phase1EncryptionDefault
		}
		if s2c.Phase2Auth != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_2_authentication should be empty")
		} else {
			s2c.Phase2Auth = goaviatrix.Phase2AuthDefault
		}
		if s2c.Phase2DhGroups != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_2_dh_groups should be empty")
		} else {
			s2c.Phase2DhGroups = goaviatrix.Phase2DhGroupDefault
		}
		if s2c.Phase2Encryption != "" {
			return fmt.Errorf("custom_algorithms is disabled, phase_2_encryption should be empty")
		} else {
			s2c.Phase2Encryption = goaviatrix.Phase2EncryptionDefault
		}
	}

	privateRouteEncryption := d.Get("private_route_encryption").(bool)
	var routeTableList []string
	rTList := d.Get("route_table_list").([]interface{})
	for i := range rTList {
		routeTableList = append(routeTableList, rTList[i].(string))
	}

	if privateRouteEncryption && len(routeTableList) == 0 {
		return fmt.Errorf("private_route_encryption is enabled, route_table_list cannot be empty")
	} else if privateRouteEncryption {
		s2c.PrivateRouteEncryption = "true"
		s2c.RouteTableList = routeTableList
		s2c.RemoteGwLatitude = d.Get("remote_gateway_latitude").(float64)
		s2c.RemoteGwLongitude = d.Get("remote_gateway_longitude").(float64)
		if s2c.HAEnabled == "yes" {
			s2c.BackupRemoteGwLatitude = d.Get("backup_remote_gateway_latitude").(float64)
			s2c.BackupRemoteGwLongitude = d.Get("backup_remote_gateway_longitude").(float64)
		}
	} else {
		s2c.PrivateRouteEncryption = "false"
		if len(routeTableList) != 0 {
			return fmt.Errorf("private route encryption is disabled, route_table_list should be empty")
		}
	}

	sslServerPool := d.Get("ssl_server_pool").(string)

	if s2c.TunnelType == "udp" && sslServerPool != "" {
		return fmt.Errorf("ssl_server_pool only supports tunnel type 'tcp'")
	}
	if s2c.TunnelType == "tcp" && s2c.RemoteGwType != "avx" {
		return fmt.Errorf("only 'avx' remote gateway type is supported for tunnel type 'tcp'")
	}
	if s2c.TunnelType == "tcp" {
		if sslServerPool == goaviatrix.SslServerPoolDefault {
			return fmt.Errorf("'192.168.44.0/24' is default, please specify a different value for ssl_server_pool")
		} else if sslServerPool == "" {
			s2c.SslServerPool = goaviatrix.SslServerPoolDefault
		} else {
			s2c.SslServerPool = sslServerPool
		}
	}

	log.Printf("[INFO] Creating Aviatrix Site2Cloud: %#v", s2c)

	err := client.CreateSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed Site2Cloud create: %s", err)
	}

	d.SetId(s2c.TunnelName + "~" + s2c.VpcID)

	flag := false
	defer resourceAviatrixSite2CloudReadIfRequired(d, meta, &flag)

	enableDeadPeerDetection := d.Get("enable_dead_peer_detection").(bool)
	if !enableDeadPeerDetection {
		err := client.DisableDeadPeerDetection(s2c)
		if err != nil {
			return fmt.Errorf("failed to disable dead peer detection: %s", err)
		}
	}

	return resourceAviatrixSite2CloudReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSite2CloudReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSite2CloudRead(d, meta)
	}
	return nil
}

func resourceAviatrixSite2CloudRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tunnelName := d.Get("connection_name").(string)
	vpcID := d.Get("vpc_id").(string)
	if tunnelName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no tunnel name or vpc id names received. Import Id is %s", id)
		d.Set("connection_name", strings.Split(id, "~")[0])
		d.Set("vpc_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	site2cloud := &goaviatrix.Site2Cloud{
		TunnelName: d.Get("connection_name").(string),
		VpcID:      d.Get("vpc_id").(string),
	}
	s2c, err := client.GetSite2CloudConnDetail(site2cloud)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Site2Cloud: %s, %#v", err, s2c)
	}

	if s2c != nil {
		d.Set("vpc_id", s2c.VpcID)
		d.Set("remote_gateway_type", s2c.RemoteGwType)
		d.Set("tunnel_type", s2c.TunnelType)
		d.Set("local_subnet_cidr", s2c.LocalSubnet)
		d.Set("remote_subnet_cidr", s2c.RemoteSubnet)
		if s2c.HAEnabled == "enabled" {
			d.Set("ha_enabled", true)
		} else {
			d.Set("ha_enabled", false)
		}

		if s2c.HAEnabled == "enabled" {
			d.Set("remote_gateway_ip", s2c.RemoteGwIP)
			d.Set("backup_remote_gateway_ip", s2c.RemoteGwIP2)
			d.Set("primary_cloud_gateway_name", s2c.GwName)
			d.Set("backup_gateway_name", s2c.BackupGwName)
		} else {
			d.Set("remote_gateway_ip", s2c.RemoteGwIP)
			d.Set("primary_cloud_gateway_name", s2c.GwName)
		}

		d.Set("connection_type", s2c.ConnType)
		if s2c.ConnType == "mapped" {
			d.Set("remote_subnet_virtual", s2c.RemoteSubnetVirtual)
			d.Set("local_subnet_virtual", s2c.LocalSubnetVirtual)
		}

		if s2c.CustomAlgorithms {
			d.Set("custom_algorithms", true)
			d.Set("phase_1_authentication", s2c.Phase1Auth)
			d.Set("phase_2_authentication", s2c.Phase2Auth)
			d.Set("phase_1_dh_groups", s2c.Phase1DhGroups)
			d.Set("phase_2_dh_groups", s2c.Phase2DhGroups)
			d.Set("phase_1_encryption", s2c.Phase1Encryption)
			d.Set("phase_2_encryption", s2c.Phase2Encryption)
		} else {
			d.Set("custom_algorithms", false)
		}

		if s2c.PrivateRouteEncryption == "true" {
			d.Set("private_route_encryption", true)

			if err := d.Set("route_table_list", s2c.RouteTableList); err != nil {
				log.Printf("[WARN] Error setting route_table_list for (%s): %s", d.Id(), err)
			}
		} else {
			d.Set("private_route_encryption", false)
		}

		if s2c.SslServerPool != "" {
			d.Set("ssl_server_pool", s2c.SslServerPool)
		}

		d.Set("enable_dead_peer_detection", s2c.DeadPeerDetection)
	}

	log.Printf("[TRACE] Reading Aviatrix Site2Cloud %s: %#v", d.Get("connection_name").(string), site2cloud)
	log.Printf("[TRACE] Reading Aviatrix Site2Cloud connection_type: [%s]", d.Get("connection_type").(string))

	d.SetId(site2cloud.TunnelName + "~" + site2cloud.VpcID)
	return nil
}

func resourceAviatrixSite2CloudUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	editSite2cloud := &goaviatrix.EditSite2Cloud{
		GwName:   d.Get("primary_cloud_gateway_name").(string),
		VpcID:    d.Get("vpc_id").(string),
		ConnName: d.Get("connection_name").(string),
	}

	d.Partial(true)
	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", editSite2cloud)

	if ok := d.HasChange("local_subnet_cidr"); ok {
		editSite2cloud.CloudSubnetCidr = d.Get("local_subnet_cidr").(string)
		editSite2cloud.NetworkType = "1"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud local_subnet_cidr: %s", err)
		}
		d.SetPartial("local_subnet_cidr")
	}

	if ok := d.HasChange("remote_subnet_cidr"); ok {
		editSite2cloud.CloudSubnetCidr = d.Get("remote_subnet_cidr").(string)
		editSite2cloud.NetworkType = "2"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud remote_subnet_cidr: %s", err)
		}
		d.SetPartial("remote_subnet_cidr")
	}

	if d.HasChange("enable_dead_peer_detection") {
		s2c := &goaviatrix.Site2Cloud{
			VpcID:      d.Get("vpc_id").(string),
			TunnelName: d.Get("connection_name").(string),
		}
		enableDeadPeerDetection := d.Get("enable_dead_peer_detection").(bool)
		if enableDeadPeerDetection {
			err := client.EnableDeadPeerDetection(s2c)
			if err != nil {
				return fmt.Errorf("failed to enable deed peer detection: %s", err)
			}
		} else {
			err := client.DisableDeadPeerDetection(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable deed peer detection: %s", err)
			}
		}
		d.SetPartial("enable_dead_peer_detection")
	}

	d.Partial(false)
	d.SetId(editSite2cloud.ConnName + "~" + editSite2cloud.VpcID)
	return resourceAviatrixSite2CloudRead(d, meta)
}

func resourceAviatrixSite2CloudDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	s2c := &goaviatrix.Site2Cloud{
		VpcID:      d.Get("vpc_id").(string),
		TunnelName: d.Get("connection_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix s2c: %#v", s2c)

	err := client.DeleteSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Site2Cloud: %s", err)
	}

	return nil
}
