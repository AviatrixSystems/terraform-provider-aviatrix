package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var customMappedAttributeNames = []string{
	"remote_source_real_cidrs",
	"remote_source_virtual_cidrs",
	"remote_destination_real_cidrs",
	"remote_destination_virtual_cidrs",
	"local_source_real_cidrs",
	"local_source_virtual_cidrs",
	"local_destination_real_cidrs",
	"local_destination_virtual_cidrs",
}

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"policy", "route"}, false),
				Description:  "Site2Cloud Tunnel Type. Valid values: 'policy' and 'route'.",
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Remote Subnet CIDR.",
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
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Local Subnet CIDR.",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Specify whether enabling HA or not.",
			},
			"backup_pre_shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Backup Pre-Shared Key.",
			},
			"remote_subnet_virtual": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Remote Subnet CIDR (Virtual).",
			},
			"local_subnet_virtual": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Local Subnet CIDR (Virtual).",
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
				Description: "Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and " +
					"'AES-256-CBC'.",
				ValidateFunc: validation.StringInSlice([]string{
					"3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC",
				}, false),
			},
			"phase_2_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', " +
					"'AES-256-CBC', 'AES-128-GCM-64', 'AES-128-GCM-96', 'AES-128-GCM-128', and 'NULL-ENCR'.",
				ValidateFunc: validation.StringInSlice([]string{
					"3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "AES-128-GCM-64", "AES-128-GCM-96",
					"AES-128-GCM-128", "NULL-ENCR",
				}, false),
			},
			"enable_ikev2": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Switch to enable IKEv2 for policy based site2cloud.",
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
			"enable_active_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable active_active_ha for an existing site2cloud connection.",
			},
			"forward_traffic_to_transit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable spoke gateway with mapped site2cloud configurations to forward traffic from site2cloud connection to Aviatrix Transit Gateway.",
			},
			"custom_mapped": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable custom mapped.",
			},
			"enable_single_ip_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable single IP HA on a site2cloud connection.",
			},
			"remote_source_real_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Remote Initiated Traffic Source Real CIDRs.",
			},
			"remote_source_virtual_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Remote Initiated Traffic Source Virtual CIDRs.",
			},
			"remote_destination_real_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Remote Initiated Traffic Destination Real CIDRs.",
			},
			"remote_destination_virtual_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Remote Initiated Traffic Destination Virtual CIDRs.",
			},
			"local_source_real_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Local Initiated Traffic Source Real CIDRs.",
			},
			"local_source_virtual_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Local Initiated Traffic Source Virtual CIDRs.",
			},
			"local_destination_real_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Local Initiated Traffic Destination Real CIDRs.",
			},
			"local_destination_virtual_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				Description: "Local Initiated Traffic Destination Virtual CIDRs.",
			},
			"enable_event_triggered_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Event Triggered HA.",
			},
			"local_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local tunnel IP address.",
			},
			"remote_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote tunnel IP address.",
			},
			"backup_local_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup local tunnel IP address.",
			},
			"backup_remote_tunnel_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup remote tunnel IP address.",
			},
			"backup_remote_gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup remote remote gateway IP.",
			},
			"phase1_remote_identifier": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsIPv4Address},
				DiffSuppressFunc: goaviatrix.S2CPh1RemoteIdDiffSuppressFunc,
				Description:      "Phase 1 remote identifier of the IPsec tunnel.",
			},
		},
	}
}

func getCSVFromStringList(d *schema.ResourceData, attributeName string) string {
	s := d.Get(attributeName).([]interface{})
	expandedList := goaviatrix.ExpandStringList(s)
	return strings.Join(expandedList, ",")
}

func resourceAviatrixSite2CloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	s2c := &goaviatrix.Site2Cloud{
		GwName:                        d.Get("primary_cloud_gateway_name").(string),
		BackupGwName:                  d.Get("backup_gateway_name").(string),
		VpcID:                         d.Get("vpc_id").(string),
		TunnelName:                    d.Get("connection_name").(string),
		ConnType:                      d.Get("connection_type").(string),
		TunnelType:                    d.Get("tunnel_type").(string),
		RemoteGwType:                  d.Get("remote_gateway_type").(string),
		RemoteGwIP:                    d.Get("remote_gateway_ip").(string),
		RemoteGwIP2:                   d.Get("backup_remote_gateway_ip").(string),
		PreSharedKey:                  d.Get("pre_shared_key").(string),
		BackupPreSharedKey:            d.Get("backup_pre_shared_key").(string),
		RemoteSubnet:                  d.Get("remote_subnet_cidr").(string),
		LocalSubnet:                   d.Get("local_subnet_cidr").(string),
		RemoteSubnetVirtual:           d.Get("remote_subnet_virtual").(string),
		LocalSubnetVirtual:            d.Get("local_subnet_virtual").(string),
		CustomMap:                     d.Get("custom_mapped").(bool),
		RemoteSourceRealCIDRs:         getCSVFromStringList(d, "remote_source_real_cidrs"),
		RemoteSourceVirtualCIDRs:      getCSVFromStringList(d, "remote_source_virtual_cidrs"),
		RemoteDestinationRealCIDRs:    getCSVFromStringList(d, "remote_destination_real_cidrs"),
		RemoteDestinationVirtualCIDRs: getCSVFromStringList(d, "remote_destination_virtual_cidrs"),
		LocalSourceRealCIDRs:          getCSVFromStringList(d, "local_source_real_cidrs"),
		LocalSourceVirtualCIDRs:       getCSVFromStringList(d, "local_source_virtual_cidrs"),
		LocalDestinationRealCIDRs:     getCSVFromStringList(d, "local_destination_real_cidrs"),
		LocalDestinationVirtualCIDRs:  getCSVFromStringList(d, "local_destination_virtual_cidrs"),
		LocalTunnelIp:                 d.Get("local_tunnel_ip").(string),
		RemoteTunnelIp:                d.Get("remote_tunnel_ip").(string),
		BackupLocalTunnelIp:           d.Get("backup_local_tunnel_ip").(string),
		BackupRemoteTunnelIp:          d.Get("backup_remote_tunnel_ip").(string),
	}

	haEnabled := d.Get("ha_enabled").(bool)
	singleIpHA := d.Get("enable_single_ip_ha").(bool)
	if haEnabled {
		s2c.HAEnabled = "yes"
		// 22021: Remote GW IP is not required when singleIPHA is enabled as only 1 tunnel is created
		if s2c.BackupGwName == "" || (s2c.RemoteGwIP2 == "" && !singleIpHA) {
			return fmt.Errorf("'backup_gateway_name' and 'backup_remote_gateway_ip' are required when HA is enabled")
		}
	} else {
		s2c.HAEnabled = "no"
		if s2c.BackupGwName != "" || s2c.RemoteGwIP2 != "" || s2c.BackupLocalTunnelIp != "" || s2c.BackupRemoteTunnelIp != "" {
			return fmt.Errorf("'backup_gateway_name', 'backup_remote_gateway_ip', 'backup_local_tunnel_ip' " +
				"and 'backup_remote_tunnel_ip' are only valid when HA is enabled")
		}
	}

	if s2c.TunnelType == "policy" {
		if s2c.LocalTunnelIp != "" || s2c.RemoteTunnelIp != "" || s2c.BackupLocalTunnelIp != "" || s2c.BackupRemoteTunnelIp != "" {
			return fmt.Errorf("'local_tunnel_ip', 'remote_tunnel_ip', 'backup_local_tunnel_ip' " +
				"and 'backup_remote_tunnel_ip' are only valid for route based connection")
		}
	}

	activeActive := d.Get("enable_active_active").(bool)
	if activeActive && !haEnabled {
		return fmt.Errorf("active_active_ha can't be enabled if HA isn't enabled for site2cloud connection")
	}

	if s2c.ConnType != "mapped" && s2c.ConnType != "unmapped" {
		return fmt.Errorf("'connection_type' should be 'mapped' or 'unmapped'")
	}

	gateway := &goaviatrix.Gateway{
		GwName: s2c.GwName,
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		return fmt.Errorf("couldn't find Aviatrix Gateway %s: %v", s2c.GwName, err)
	}

	if gw.TransitVpc == "yes" {
		if s2c.ConnType == "unmapped" && s2c.TunnelType == "policy" && haEnabled && !activeActive {
			return fmt.Errorf("active_active_ha must be enabled if HA is enabled for transit gateway unmapped policy based site2cloud connection")
		}
	}

	if singleIpHA {
		if !haEnabled {
			return fmt.Errorf("'enable_single_ip_ha' can't be enabled if HA isn't enabled for site2cloud connection")
		}
		s2c.EnableSingleIpHA = true
	}

	if !s2c.CustomMap && s2c.RemoteSubnet == "" {
		return fmt.Errorf("'remote_subnet_cidr' is required unless you are using 'custom_mapped'")
	}
	if s2c.CustomMap {
		if s2c.RemoteSubnet != "" {
			return fmt.Errorf("'remote_subnet_cidr' is not valid for 'custom_mapped' connection")
		}
		if s2c.RemoteSubnetVirtual != "" {
			return fmt.Errorf("'remote_subnet_virtual' is not valid for 'custom_mapped' connection")
		}
		if s2c.LocalSubnet != "" {
			return fmt.Errorf("'local_subnet_cidr' is not valid for 'custom_mapped' connection")
		}
		if s2c.LocalSubnetVirtual != "" {
			return fmt.Errorf("'local_subnet_virtual' is not valid for 'custom_mapped' connection")
		}
	}
	if s2c.ConnType == "mapped" && !s2c.CustomMap && (s2c.RemoteSubnetVirtual == "" || s2c.LocalSubnetVirtual == "") {
		return fmt.Errorf("'remote_subnet_virtual' and 'local_subnet_virtual' are both required for " +
			"connection type: mapped, unless 'custom_mapped' is enabled")
	} else if s2c.ConnType == "unmapped" && (s2c.RemoteSubnetVirtual != "" || s2c.LocalSubnetVirtual != "") {
		return fmt.Errorf("'remote_subnet_virtual' and 'local_subnet_virtual' both should be empty for " +
			"connection type: ummapped")
	}
	hasSetAnyCustomMapAttribute := s2c.RemoteSourceRealCIDRs != "" || s2c.RemoteSourceVirtualCIDRs != "" ||
		s2c.RemoteDestinationRealCIDRs != "" || s2c.RemoteDestinationVirtualCIDRs != "" ||
		s2c.LocalSourceRealCIDRs != "" || s2c.LocalSourceVirtualCIDRs != "" ||
		s2c.LocalDestinationRealCIDRs != "" || s2c.LocalDestinationVirtualCIDRs != ""
	if !s2c.CustomMap && hasSetAnyCustomMapAttribute {
		return fmt.Errorf("attributes %v are only valid with 'custom_mapped' enabled", customMappedAttributeNames)
	}
	if s2c.CustomMap && (s2c.ConnType != "mapped" || s2c.TunnelType != "route") {
		return fmt.Errorf("'connection_type' should be 'mapped' and 'tunnel_type' should be 'route' for 'custom_mapped' enabled connection")
	}
	hasSetAllCustomRemoteCIDRs := s2c.RemoteSourceRealCIDRs != "" && s2c.RemoteSourceVirtualCIDRs != "" && s2c.RemoteDestinationRealCIDRs != "" && s2c.RemoteDestinationVirtualCIDRs != ""
	hasSetAllCustomLocalCIDRs := s2c.LocalSourceRealCIDRs != "" && s2c.LocalSourceVirtualCIDRs != "" && s2c.LocalDestinationRealCIDRs != "" && s2c.LocalDestinationVirtualCIDRs != ""
	if s2c.CustomMap && !hasSetAllCustomLocalCIDRs && !hasSetAllCustomRemoteCIDRs {
		return fmt.Errorf("'custom_mapped' enabled connection requires either all Remote Initiated CIDRs or all Local Initated CIDRs be provided")
	}

	s2c.Phase1Auth = d.Get("phase_1_authentication").(string)
	s2c.Phase1DhGroups = d.Get("phase_1_dh_groups").(string)
	s2c.Phase1Encryption = d.Get("phase_1_encryption").(string)
	s2c.Phase2Auth = d.Get("phase_2_authentication").(string)
	s2c.Phase2DhGroups = d.Get("phase_2_dh_groups").(string)
	s2c.Phase2Encryption = d.Get("phase_2_encryption").(string)

	customAlgorithms := d.Get("custom_algorithms").(bool)
	if customAlgorithms {
		if s2c.Phase1Auth == "" ||
			s2c.Phase2Auth == "" ||
			s2c.Phase1DhGroups == "" ||
			s2c.Phase2DhGroups == "" ||
			s2c.Phase1Encryption == "" ||
			s2c.Phase2Encryption == "" {
			return fmt.Errorf("custom_algorithms is enabled, please set all of the algorithm parameters")
		} else if s2c.Phase1Auth == goaviatrix.Phase1AuthDefault &&
			s2c.Phase2Auth == goaviatrix.Phase2AuthDefault &&
			s2c.Phase1DhGroups == goaviatrix.Phase1DhGroupDefault &&
			s2c.Phase2DhGroups == goaviatrix.Phase2DhGroupDefault &&
			s2c.Phase1Encryption == goaviatrix.Phase1EncryptionDefault &&
			s2c.Phase2Encryption == goaviatrix.Phase2EncryptionDefault {
			return fmt.Errorf("custom_algorithms is enabled, cannot use default values for " +
				"all six algorithm parameters. Please change value of one or multiple of the six algorithm parameters")
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

	enableIKEv2 := d.Get("enable_ikev2").(bool)
	if enableIKEv2 {
		s2c.EnableIKEv2 = "true"
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

	phase1RemoteIdentifier := d.Get("phase1_remote_identifier").([]interface{})
	ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)
	if haEnabled && !singleIpHA && len(ph1RemoteIdList) != 0 && len(ph1RemoteIdList) != 2 {
		return fmt.Errorf("please either set two phase 1 remote IDs or none, when HA is enabled and single IP HA is disabled")
	} else if (!haEnabled || singleIpHA) && len(phase1RemoteIdentifier) > 1 {
		return fmt.Errorf("please either set one phase 1 remote ID or none, when HA is disabled or single IP HA is enabled")
	}

	log.Printf("[INFO] Creating Aviatrix Site2Cloud: %#v", s2c)

	d.SetId(s2c.TunnelName + "~" + s2c.VpcID)
	flag := false
	defer resourceAviatrixSite2CloudReadIfRequired(d, meta, &flag)

	err = client.CreateSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed Site2Cloud create: %s", err)
	}

	enableDeadPeerDetection := d.Get("enable_dead_peer_detection").(bool)
	if !enableDeadPeerDetection {
		err := client.DisableDeadPeerDetection(s2c)
		if err != nil {
			return fmt.Errorf("failed to disable dead peer detection: %s", err)
		}
	}

	if activeActive {
		err := client.EnableSite2cloudActiveActive(s2c)
		if err != nil {
			return fmt.Errorf("failed to enable active active HA for site2cloud: %s: %s", s2c.TunnelName, err)
		}
	} else {
		if gw.TransitVpc == "no" && s2c.ConnType == "unmapped" && s2c.TunnelType == "route" && haEnabled {
			err := client.DisableSite2cloudActiveActive(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable active active HA for site2cloud: %s: %s", s2c.TunnelName, err)
			}
		}
	}

	forwardToTransit := d.Get("forward_traffic_to_transit").(bool)
	if forwardToTransit {
		err := client.EnableSpokeMappedSite2CloudForwarding(s2c)
		if err != nil {
			return fmt.Errorf("failed to enable traffic forwarding to transit: %v", err)
		}
	}

	if d.Get("enable_event_triggered_ha").(bool) {
		err := client.EnableSite2CloudEventTriggeredHA(s2c.VpcID, s2c.TunnelName)
		if err != nil {
			return fmt.Errorf("could not enable event triggered HA for site2cloud after creation: %v", err)
		}
	}

	if len(ph1RemoteIdList) == 1 && ph1RemoteIdList[0] != s2c.RemoteGwIP {
		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:                 s2c.GwName,
			VpcID:                  s2c.VpcID,
			ConnName:               s2c.TunnelName,
			Phase1RemoteIdentifier: ph1RemoteIdList[0],
		}

		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud phase 1 remote identifier: %s", err)
		}
	}

	if len(ph1RemoteIdList) == 2 && (ph1RemoteIdList[0] != s2c.RemoteGwIP || ph1RemoteIdList[1] != s2c.RemoteGwIP2) {
		editSite2cloud := &goaviatrix.EditSite2Cloud{
			GwName:                 s2c.GwName,
			VpcID:                  s2c.VpcID,
			ConnName:               s2c.TunnelName,
			Phase1RemoteIdentifier: strings.Join(ph1RemoteIdList, ","),
		}

		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud phase 1 remote identifier: %s", err)
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
		parts := strings.Split(id, "~")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid import ID format")
		}
		d.Set("connection_name", parts[0])
		d.Set("vpc_id", parts[1])
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

		d.Set("remote_gateway_ip", s2c.RemoteGwIP)
		d.Set("primary_cloud_gateway_name", s2c.GwName)
		d.Set("local_tunnel_ip", s2c.LocalTunnelIp)
		d.Set("remote_tunnel_ip", s2c.RemoteTunnelIp)

		if s2c.HAEnabled == "enabled" {
			d.Set("backup_remote_gateway_ip", s2c.RemoteGwIP2)
			d.Set("backup_gateway_name", s2c.BackupGwName)
			d.Set("backup_local_tunnel_ip", s2c.BackupLocalTunnelIp)
			d.Set("backup_remote_tunnel_ip", s2c.BackupRemoteTunnelIp)
		}

		// Custom Mapped is a sub-type of Mapped
		if s2c.ConnType == "custom_mapped" {
			d.Set("custom_mapped", true)
			s2c.ConnType = "mapped"
		} else {
			d.Set("custom_mapped", false)
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
		d.Set("enable_active_active", s2c.EnableActiveActive)
		d.Set("forward_traffic_to_transit", s2c.ForwardToTransit)
		d.Set("enable_event_triggered_ha", s2c.EventTriggeredHA)
		d.Set("enable_single_ip_ha", s2c.EnableSingleIpHA)

		if s2c.EnableIKEv2 == "true" {
			d.Set("enable_ikev2", true)
		} else {
			d.Set("enable_ikev2", false)
		}

		if s2c.RemoteSourceRealCIDRs != "" {
			if err := d.Set("remote_source_real_cidrs", strings.Split(s2c.RemoteSourceRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_source_real_cidrs' to state: %v", err)
			}
		}
		if s2c.RemoteSourceVirtualCIDRs != "" {
			if err := d.Set("remote_source_virtual_cidrs", strings.Split(s2c.RemoteSourceVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_source_virtual_cidrs' to state: %v", err)
			}
		}
		if s2c.RemoteDestinationRealCIDRs != "" {
			if err := d.Set("remote_destination_real_cidrs", strings.Split(s2c.RemoteDestinationRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_destination_real_cidrs' to state: %v", err)
			}
		}
		if s2c.RemoteDestinationVirtualCIDRs != "" {
			if err := d.Set("remote_destination_virtual_cidrs", strings.Split(s2c.RemoteDestinationVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_destination_virtual_cidrs' to state: %v", err)
			}
		}
		if s2c.LocalSourceRealCIDRs != "" {
			if err := d.Set("local_source_real_cidrs", strings.Split(s2c.LocalSourceRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_source_real_cidrs' to state: %v", err)
			}
		}
		if s2c.LocalSourceVirtualCIDRs != "" {
			if err := d.Set("local_source_virtual_cidrs", strings.Split(s2c.LocalSourceVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_source_virtual_cidrs' to state: %v", err)
			}
		}
		if s2c.LocalDestinationRealCIDRs != "" {
			if err := d.Set("local_destination_real_cidrs", strings.Split(s2c.LocalDestinationRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_destination_real_cidrs' to state: %v", err)
			}
		}
		if s2c.LocalDestinationVirtualCIDRs != "" {
			if err := d.Set("local_destination_virtual_cidrs", strings.Split(s2c.LocalDestinationVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_destination_virtual_cidrs' to state: %v", err)
			}
		}

		ph1RemoteId := strings.Split(s2c.Phase1RemoteIdentifier, ",")
		for i, v := range ph1RemoteId {
			ph1RemoteId[i] = strings.TrimSpace(v)
		}

		d.Set("phase1_remote_identifier", ph1RemoteId)
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

	if d.HasChange("local_subnet_cidr") {
		if d.Get("custom_mapped").(bool) && d.Get("local_subnet_cidr").(string) != "" {
			return fmt.Errorf("'local_subnet_cidr' is not valid when 'custom_mapped' is enabled")
		}
		editSite2cloud.CloudSubnetCidr = d.Get("local_subnet_cidr").(string)
		editSite2cloud.CloudSubnetVirtual = d.Get("local_subnet_virtual").(string)
		editSite2cloud.NetworkType = "1"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud local_subnet_cidr: %s", err)
		}
	}

	if d.HasChange("local_subnet_virtual") {
		if d.Get("custom_mapped").(bool) && d.Get("local_subnet_virtual").(string) != "" {
			return fmt.Errorf("'local_subnet_virtual' is not valid when 'custom_mapped' is enabled")
		}
		if d.Get("connection_type").(string) == "mapped" && d.Get("local_subnet_virtual").(string) == "" {
			return fmt.Errorf("'local_subnet_virtual' is required for connection type: mapped, unless 'custom_mapped' is enabled")
		}
		if d.Get("connection_type").(string) == "unmapped" && d.Get("local_subnet_virtual").(string) != "" {
			return fmt.Errorf("'local_subnet_virtual' should be empty for connection type: ummapped")
		}
		editSite2cloud.CloudSubnetCidr = d.Get("local_subnet_cidr").(string)
		editSite2cloud.CloudSubnetVirtual = d.Get("local_subnet_virtual").(string)
		editSite2cloud.NetworkType = "1"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud local_subnet_virtual: %s", err)
		}
	}

	if d.HasChange("remote_subnet_cidr") {
		if d.Get("custom_mapped").(bool) && d.Get("remote_subnet_cidr").(string) != "" {
			return fmt.Errorf("'remote_subnet_cidr' is not valid when 'custom_mapped' is enabled")
		}
		editSite2cloud.CloudSubnetCidr = d.Get("remote_subnet_cidr").(string)
		editSite2cloud.CloudSubnetVirtual = d.Get("remote_subnet_virtual").(string)
		editSite2cloud.NetworkType = "2"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud remote_subnet_cidr: %s", err)
		}
	}

	if d.HasChange("remote_subnet_virtual") {
		if d.Get("custom_mapped").(bool) && d.Get("remote_subnet_virtual").(string) != "" {
			return fmt.Errorf("'remote_subnet_virtual' is not valid when 'custom_mapped' is enabled")
		}
		if d.Get("connection_type").(string) == "mapped" && d.Get("remote_subnet_virtual").(string) == "" {
			return fmt.Errorf("'remote_subnet_virtual' is required for connection type: mapped, unless 'custom_mapped' is enabled")
		}
		if d.Get("connection_type").(string) == "unmapped" && d.Get("remote_subnet_virtual").(string) != "" {
			return fmt.Errorf("'remote_subnet_virtual' should be empty for connection type: ummapped")
		}
		editSite2cloud.CloudSubnetCidr = d.Get("remote_subnet_cidr").(string)
		editSite2cloud.CloudSubnetVirtual = d.Get("remote_subnet_virtual").(string)
		editSite2cloud.NetworkType = "2"
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud remote_subnet_virtual: %s", err)
		}
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
	}

	if d.HasChange("enable_active_active") {
		s2c := &goaviatrix.Site2Cloud{
			VpcID:      d.Get("vpc_id").(string),
			TunnelName: d.Get("connection_name").(string),
		}
		activeActive := d.Get("enable_active_active").(bool)
		if activeActive {
			if haEnabled := d.Get("ha_enabled").(bool); !haEnabled {
				return fmt.Errorf("active_active_ha can't be enabled if HA isn't enabled for site2cloud connection")
			}
			err := client.EnableSite2cloudActiveActive(s2c)
			if err != nil {
				return fmt.Errorf("failed to enable active active HA for site2cloud: %s: %s", s2c.TunnelName, err)
			}
		} else {
			err := client.DisableSite2cloudActiveActive(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable active active HA for site2cloud: %s: %s", s2c.TunnelName, err)
			}
		}
	}

	if d.HasChange("forward_traffic_to_transit") {
		s2c := &goaviatrix.Site2Cloud{
			VpcID:      d.Get("vpc_id").(string),
			TunnelName: d.Get("connection_name").(string),
		}
		forwardToTransit := d.Get("forward_traffic_to_transit").(bool)
		if forwardToTransit {
			err := client.EnableSpokeMappedSite2CloudForwarding(s2c)
			if err != nil {
				return fmt.Errorf("failed to enable traffic forwarding to transit for site2cloud: %s: %s", s2c.TunnelName, err)
			}
		} else {
			err := client.DisableSpokeMappedSite2CloudForwarding(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable traffic forwarding to transit for site2cloud: %s: %s", s2c.TunnelName, err)
			}
		}
	}

	if d.HasChange("enable_event_triggered_ha") {
		if d.Get("enable_event_triggered_ha").(bool) {
			err := client.EnableSite2CloudEventTriggeredHA(editSite2cloud.VpcID, editSite2cloud.ConnName)
			if err != nil {
				return fmt.Errorf("could not enable event triggered HA for site2cloud during update: %v", err)
			}
		} else {
			err := client.DisableSite2CloudEventTriggeredHA(editSite2cloud.VpcID, editSite2cloud.ConnName)
			if err != nil {
				return fmt.Errorf("could not disable event triggered HA for site2cloud during update: %v", err)
			}
		}
	}

	if d.HasChanges(customMappedAttributeNames...) {
		if !d.Get("custom_mapped").(bool) {
			return fmt.Errorf("attributes %v are not valid when 'custom_mapped' is disabled", customMappedAttributeNames)
		}
		s2c := &goaviatrix.EditSite2Cloud{
			GwName:                        d.Get("primary_cloud_gateway_name").(string),
			VpcID:                         d.Get("vpc_id").(string),
			ConnName:                      d.Get("connection_name").(string),
			NetworkType:                   "3",
			RemoteSourceRealCIDRs:         getCSVFromStringList(d, "remote_source_real_cidrs"),
			RemoteSourceVirtualCIDRs:      getCSVFromStringList(d, "remote_source_virtual_cidrs"),
			RemoteDestinationRealCIDRs:    getCSVFromStringList(d, "remote_destination_real_cidrs"),
			RemoteDestinationVirtualCIDRs: getCSVFromStringList(d, "remote_destination_virtual_cidrs"),
			LocalSourceRealCIDRs:          getCSVFromStringList(d, "local_source_real_cidrs"),
			LocalSourceVirtualCIDRs:       getCSVFromStringList(d, "local_source_virtual_cidrs"),
			LocalDestinationRealCIDRs:     getCSVFromStringList(d, "local_destination_real_cidrs"),
			LocalDestinationVirtualCIDRs:  getCSVFromStringList(d, "local_destination_virtual_cidrs"),
		}
		hasSetAllCustomRemoteCIDRs := s2c.RemoteSourceRealCIDRs != "" && s2c.RemoteSourceVirtualCIDRs != "" && s2c.RemoteDestinationRealCIDRs != "" && s2c.RemoteDestinationVirtualCIDRs != ""
		hasSetAllCustomLocalCIDRs := s2c.LocalSourceRealCIDRs != "" && s2c.LocalSourceVirtualCIDRs != "" && s2c.LocalDestinationRealCIDRs != "" && s2c.LocalDestinationVirtualCIDRs != ""
		if !hasSetAllCustomLocalCIDRs && !hasSetAllCustomRemoteCIDRs {
			return fmt.Errorf("'custom_mapped' enabled connection requires either all Remote Initiated CIDRs or all Local Initated CIDRs be provided")
		}
		err := client.UpdateSite2Cloud(s2c)
		if err != nil {
			return fmt.Errorf("could not update site2cloud connection Remote or Local CIDRs: %v", err)
		}
	}

	if d.HasChange("phase1_remote_identifier") {
		haEnabled := d.Get("ha_enabled").(bool)
		singleIpHA := d.Get("enable_single_ip_ha").(bool)
		ip := d.Get("remote_gateway_ip").(string)
		haIp := d.Get("backup_remote_gateway_ip").(string)
		phase1RemoteIdentifier := d.Get("phase1_remote_identifier").([]interface{})
		ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)

		if haEnabled && !singleIpHA && len(ph1RemoteIdList) != 0 && len(ph1RemoteIdList) != 2 {
			return fmt.Errorf("please either set two phase 1 remote IDs or none, when HA is enabled and single IP HA is disabled")
		} else if (!haEnabled || singleIpHA) && len(phase1RemoteIdentifier) > 1 {
			return fmt.Errorf("please either set one phase 1 remote ID or none, when HA is disabled or single IP HA is enabled")
		}

		if len(ph1RemoteIdList) == 0 && haEnabled && !singleIpHA {
			editSite2cloud.Phase1RemoteIdentifier = ip + "," + haIp
		} else if len(ph1RemoteIdList) == 0 && (!haEnabled || singleIpHA) {
			editSite2cloud.Phase1RemoteIdentifier = ip
		} else {
			editSite2cloud.Phase1RemoteIdentifier = strings.Join(ph1RemoteIdList, ",")
		}

		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud phase 1 remote identifier: %s", err)
		}
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

	forwardToTransit := d.Get("forward_traffic_to_transit").(bool)
	if forwardToTransit {
		err := client.DisableSpokeMappedSite2CloudForwarding(s2c)
		if err != nil {
			log.Println("[WARN] Failed to disable forwarding to transit:", err)
		}
	}

	err := client.DeleteSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Site2Cloud: %s", err)
	}

	return nil
}
