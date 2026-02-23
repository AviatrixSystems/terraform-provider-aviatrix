package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PSK",
				ValidateFunc: validation.StringInSlice([]string{"PSK", "Cert"}, false),
				Description:  "Authentication Type. Valid values: 'PSK' and 'Cert'. Default value: 'PSK'.",
			},
			"ca_cert_tag_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of Remote CA Certificate Tag for creating Site2Cloud tunnels. Required for Cert based authentication type.",
			},
			"remote_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remote identifier. Required for Cert based authentication type.",
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
			"backup_remote_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup remote identifier. Required for Cert based authentication type with HA enabled.",
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
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRemoteGwLatitude,
				Description:      "Latitude of remote gateway.",
			},
			"remote_gateway_longitude": {
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRemoteGwLongitude,
				Description:      "Longitude of remote gateway.",
			},
			"backup_remote_gateway_latitude": {
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncBackupRemoteGwLatitude,
				Description:      "Latitude of backup remote gateway.",
			},
			"backup_remote_gateway_longitude": {
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncBackupRemoteGwLongitude,
				Description:      "Longitude of backup remote gateway.",
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
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRemoteSourceRealCIDRs,
				Description:      "Remote Initiated Traffic Source Real CIDRs.",
			},
			"remote_source_virtual_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRemoteSourceVirtualCIDRs,
				Description:      "Remote Initiated Traffic Source Virtual CIDRs.",
			},
			"remote_destination_real_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRemoteDestinationRealCIDRs,
				Description:      "Remote Initiated Traffic Destination Real CIDRs.",
			},
			"remote_destination_virtual_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncRemoteDestinationVirtualCIDRs,
				Description:      "Remote Initiated Traffic Destination Virtual CIDRs.",
			},
			"local_source_real_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncLocalSourceRealCIDRs,
				Description:      "Local Initiated Traffic Source Real CIDRs.",
			},
			"local_source_virtual_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncLocalSourceVirtualCIDRs,
				Description:      "Local Initiated Traffic Source Virtual CIDRs.",
			},
			"local_destination_real_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncLocalDestinationRealCIDRs,
				Description:      "Local Initiated Traffic Destination Real CIDRs.",
			},
			"local_destination_virtual_cidrs": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsCIDR},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncLocalDestinationVirtualCIDRs,
				Description:      "Local Initiated Traffic Destination Virtual CIDRs.",
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
				DiffSuppressFunc: goaviatrix.S2CPh1RemoteIdDiffSuppressFunc,
				Description:      "List of phase 1 remote identifier of the IPsec tunnel. This can be configured as a list of any string, including empty string.",
			},
			"proxy_id_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable proxy ID for site2cloud connection.",
			},
		},
	}
}

func getCSVFromStringList(d *schema.ResourceData, attributeName string) string {
	s := getList(d, attributeName)
	expandedList := goaviatrix.ExpandStringList(s)
	return strings.Join(expandedList, ",")
}

func resourceAviatrixSite2CloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	s2c := &goaviatrix.Site2Cloud{
		GwName:                        getString(d, "primary_cloud_gateway_name"),
		BackupGwName:                  getString(d, "backup_gateway_name"),
		VpcID:                         getString(d, "vpc_id"),
		TunnelName:                    getString(d, "connection_name"),
		ConnType:                      getString(d, "connection_type"),
		AuthType:                      getString(d, "auth_type"),
		TunnelType:                    getString(d, "tunnel_type"),
		CaCertTagName:                 getString(d, "ca_cert_tag_name"),
		RemoteIdentifier:              getString(d, "remote_identifier"),
		RemoteGwType:                  getString(d, "remote_gateway_type"),
		RemoteGwIP:                    getString(d, "remote_gateway_ip"),
		RemoteGwIP2:                   getString(d, "backup_remote_gateway_ip"),
		PreSharedKey:                  getString(d, "pre_shared_key"),
		BackupPreSharedKey:            getString(d, "backup_pre_shared_key"),
		BackupRemoteIdentifier:        getString(d, "backup_remote_identifier"),
		RemoteSubnet:                  getString(d, "remote_subnet_cidr"),
		LocalSubnet:                   getString(d, "local_subnet_cidr"),
		RemoteSubnetVirtual:           getString(d, "remote_subnet_virtual"),
		LocalSubnetVirtual:            getString(d, "local_subnet_virtual"),
		CustomMap:                     getBool(d, "custom_mapped"),
		RemoteSourceRealCIDRs:         getCSVFromStringList(d, "remote_source_real_cidrs"),
		RemoteSourceVirtualCIDRs:      getCSVFromStringList(d, "remote_source_virtual_cidrs"),
		RemoteDestinationRealCIDRs:    getCSVFromStringList(d, "remote_destination_real_cidrs"),
		RemoteDestinationVirtualCIDRs: getCSVFromStringList(d, "remote_destination_virtual_cidrs"),
		LocalSourceRealCIDRs:          getCSVFromStringList(d, "local_source_real_cidrs"),
		LocalSourceVirtualCIDRs:       getCSVFromStringList(d, "local_source_virtual_cidrs"),
		LocalDestinationRealCIDRs:     getCSVFromStringList(d, "local_destination_real_cidrs"),
		LocalDestinationVirtualCIDRs:  getCSVFromStringList(d, "local_destination_virtual_cidrs"),
		LocalTunnelIp:                 getString(d, "local_tunnel_ip"),
		RemoteTunnelIp:                getString(d, "remote_tunnel_ip"),
		BackupLocalTunnelIp:           getString(d, "backup_local_tunnel_ip"),
		BackupRemoteTunnelIp:          getString(d, "backup_remote_tunnel_ip"),
	}

	s2c.ProxyIdEnabled = getBool(d, "proxy_id_enabled")

	haEnabled := getBool(d, "ha_enabled")
	if s2c.AuthType == "Cert" {
		if s2c.CaCertTagName == "" || s2c.RemoteIdentifier == "" {
			return fmt.Errorf("'ca_cert_tag_name' and 'remote_identifier' are both required for Cert based authentication type")
		}
		if haEnabled && s2c.BackupRemoteIdentifier == "" {
			return fmt.Errorf("'backup_remote_identifier' is required for Cert based authentication type with HA enabled")
		}
		s2c.AuthType = "pubkey"
	} else {
		if s2c.CaCertTagName != "" || s2c.RemoteIdentifier != "" || s2c.BackupRemoteIdentifier != "" {
			return fmt.Errorf("'ca_cert_tag_name', 'remote_identifier' and 'backup_remote_identifier' are required to be empty for PSK(Pubkey) based authentication type")
		}
	}

	singleIpHA := getBool(d, "enable_single_ip_ha")
	if haEnabled {
		s2c.HAEnabled = "yes"
		// 22021: Remote GW IP is not required when singleIPHA is enabled as only 1 tunnel is created
		if s2c.BackupGwName == "" || (s2c.RemoteGwIP2 == "" && !singleIpHA) {
			return fmt.Errorf("'backup_gateway_name' and 'backup_remote_gateway_ip' are required when HA is enabled")
		} else if s2c.RemoteGwIP2 != "" && singleIpHA {
			return fmt.Errorf("'backup_remote_gateway_ip' is not required when HA is enabled and single ip ha is enabled")
		}
		if s2c.RemoteGwIP2 != "" {
			s2c.RemoteGwIP = s2c.RemoteGwIP + "," + s2c.RemoteGwIP2
			s2c.RemoteGwIP2 = ""
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

	activeActive := getBool(d, "enable_active_active")
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
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("couldn't find Aviatrix Gateway %s", s2c.GwName)
		} else {
			return fmt.Errorf("couldn't find Aviatrix Gateway %s: %w", s2c.GwName, err)
		}
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
		if s2c.AuthType == "pubkey" {
			return fmt.Errorf("single IP HA is only supported for PSK authentication type based site2cloud connection")
		}
		if s2c.RemoteGwIP2 != "" && s2c.RemoteGwIP != s2c.RemoteGwIP2 {
			return fmt.Errorf("'backup_remote_gateway_ip' is required to be empty or the same as 'remote_gateway_ip' when single IP HA is enabled")
		}
		if s2c.BackupPreSharedKey != "" && s2c.PreSharedKey != s2c.BackupPreSharedKey {
			return fmt.Errorf("'backup_pre_shared_key' is required to be empty or the same as 'pre_shared_key' when single IP HA is enabled")
		}
		if s2c.BackupLocalTunnelIp != "" || s2c.BackupRemoteTunnelIp != "" {
			return fmt.Errorf("'backup_local_tunnel_ip' and 'backup_remote_tunnel_ip' are required to be empty when single IP HA is enabled")
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
		return fmt.Errorf("'custom_mapped' enabled connection requires either all Remote Initiated CIDRs or all Local Initiated CIDRs be provided")
	}

	s2c.Phase1Auth = getString(d, "phase_1_authentication")
	s2c.Phase1DhGroups = getString(d, "phase_1_dh_groups")
	s2c.Phase1Encryption = getString(d, "phase_1_encryption")
	s2c.Phase2Auth = getString(d, "phase_2_authentication")
	s2c.Phase2DhGroups = getString(d, "phase_2_dh_groups")
	s2c.Phase2Encryption = getString(d, "phase_2_encryption")

	customAlgorithms := getBool(d, "custom_algorithms")
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
				"all six algorithm parameters. Please change the value of at least one of the six algorithm parameters")
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

	enableIKEv2 := getBool(d, "enable_ikev2")
	if enableIKEv2 {
		s2c.EnableIKEv2 = "true"
	}

	privateRouteEncryption := getBool(d, "private_route_encryption")
	var routeTableList []string
	rTList := getList(d, "route_table_list")
	for i := range rTList {
		routeTableList = append(routeTableList, mustString(rTList[i]))
	}
	remoteGwLatitude := getFloat64(d, "remote_gateway_latitude")
	remoteGwLongitude := getFloat64(d, "remote_gateway_longitude")
	backupRemoteGwLatitude := getFloat64(d, "backup_remote_gateway_latitude")
	backupRemoteGwLongitude := getFloat64(d, "backup_remote_gateway_longitude")

	if privateRouteEncryption && len(routeTableList) == 0 {
		return fmt.Errorf("private_route_encryption is enabled, route_table_list cannot be empty")
	} else if privateRouteEncryption {
		s2c.PrivateRouteEncryption = "true"
		s2c.RouteTableList = routeTableList
		if remoteGwLatitude == 0 || remoteGwLongitude == 0 {
			return fmt.Errorf("private_route_encryption is enabled, please set remote_gateway_latitude and remote_gateway_longitude")
		}
		s2c.RemoteGwLatitude = remoteGwLatitude
		s2c.RemoteGwLongitude = remoteGwLongitude
		if s2c.HAEnabled == "yes" {
			if backupRemoteGwLatitude == 0 || backupRemoteGwLongitude == 0 {
				return fmt.Errorf("private_route_encryption is enabled and ha is enabled, please set backup_remote_gateway_latitude and backup_remote_gateway_longitude")
			}
			s2c.BackupRemoteGwLatitude = backupRemoteGwLatitude
			s2c.BackupRemoteGwLongitude = backupRemoteGwLongitude
		}
	} else {
		s2c.PrivateRouteEncryption = "false"
		if len(routeTableList) != 0 {
			return fmt.Errorf("private route encryption is disabled, route_table_list should be empty")
		}
		if remoteGwLatitude != 0 || remoteGwLongitude != 0 || backupRemoteGwLatitude != 0 || backupRemoteGwLongitude != 0 {
			return fmt.Errorf("private route encryption is disabled, all of remote_gateway_latitude, " +
				"backup_remote_gateway_latitude and backup_remote_gateway_longitude should be empty")
		}
	}

	sslServerPool := getString(d, "ssl_server_pool")

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

	phase1RemoteIdentifier := getList(d, "phase1_remote_identifier")
	ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)
	if haEnabled && !singleIpHA && len(phase1RemoteIdentifier) != 0 && len(phase1RemoteIdentifier) != 2 {
		if !(s2c.RemoteGwIP != "" && s2c.RemoteGwIP2 != "" && s2c.RemoteGwIP == s2c.RemoteGwIP2 && len(phase1RemoteIdentifier) == 1) {
			return fmt.Errorf("please either set two phase 1 remote IDs or none, when HA is enabled and single IP HA is disabled")
		}
	} else if (!haEnabled || singleIpHA) && len(phase1RemoteIdentifier) > 1 {
		return fmt.Errorf("please either set one phase 1 remote ID or none, when HA is disabled or single IP HA is enabled")
	}

	log.Printf("[INFO] Creating Aviatrix Site2Cloud: %#v", s2c)

	d.SetId(s2c.TunnelName + "~" + s2c.VpcID)
	flag := false
	defer func() { _ = resourceAviatrixSite2CloudReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err = client.CreateSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed Site2Cloud create: %w", err)
	}

	enableDeadPeerDetection := getBool(d, "enable_dead_peer_detection")
	if !enableDeadPeerDetection {
		err := client.DisableDeadPeerDetection(s2c)
		if err != nil {
			return fmt.Errorf("failed to disable dead peer detection: %w", err)
		}
	}

	if activeActive {
		err := client.EnableSite2cloudActiveActive(s2c)
		if err != nil {
			return fmt.Errorf("failed to enable active active HA for site2cloud: %s: %w", s2c.TunnelName, err)
		}
	} else {
		if gw.TransitVpc == "no" && s2c.ConnType == "unmapped" && s2c.TunnelType == "route" && haEnabled {
			err := client.DisableSite2cloudActiveActive(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable active active HA for site2cloud: %s: %w", s2c.TunnelName, err)
			}
		}
	}

	forwardToTransit := getBool(d, "forward_traffic_to_transit")
	if forwardToTransit {
		err := client.EnableSpokeMappedSite2CloudForwarding(s2c)
		if err != nil {
			return fmt.Errorf("failed to enable traffic forwarding to transit: %w", err)
		}
	}

	if getBool(d, "enable_event_triggered_ha") {
		err := client.EnableSite2CloudEventTriggeredHA(s2c.VpcID, s2c.TunnelName)
		if err != nil {
			return fmt.Errorf("could not enable event triggered HA for site2cloud after creation: %w", err)
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
			GwName:                 s2c.GwName,
			VpcID:                  s2c.VpcID,
			ConnName:               s2c.TunnelName,
			Phase1RemoteIdentifier: ph1RemoteId,
		}

		err = client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud phase 1 remote identifier: %w", err)
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
			GwName:                 s2c.GwName,
			VpcID:                  s2c.VpcID,
			ConnName:               s2c.TunnelName,
			Phase1RemoteIdentifier: ph1RemoteId,
		}

		err = client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud phase 1 remote identifier: %w", err)
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
	client := mustClient(meta)

	tunnelName := getString(d, "connection_name")
	vpcID := getString(d, "vpc_id")
	if tunnelName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no tunnel name or vpc id names received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid import ID format")
		}
		mustSet(d, "connection_name", parts[0])
		mustSet(d, "vpc_id", parts[1])
		d.SetId(id)
	}

	site2cloud := &goaviatrix.Site2Cloud{
		TunnelName: getString(d, "connection_name"),
		VpcID:      getString(d, "vpc_id"),
	}
	s2c, err := client.GetSite2CloudConnDetail(site2cloud)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Site2Cloud: %w, %#v", err, s2c)
	}

	if s2c != nil {
		mustSet(d, "vpc_id", s2c.VpcID)
		mustSet(d, "remote_gateway_type", s2c.RemoteGwType)
		mustSet(d, "tunnel_type", s2c.TunnelType)
		if s2c.AuthType == "pubkey" {
			mustSet(d, "auth_type", "Cert")
			mustSet(d, "ca_cert_tag_name", s2c.CaCertTagName)
			mustSet(d, "remote_identifier", s2c.RemoteIdentifier)
			if s2c.HAEnabled == "enabled" {
				mustSet(d, "backup_remote_identifier", s2c.BackupRemoteIdentifier)
			}
		} else {
			mustSet(d, "auth_type", "PSK")
		}
		mustSet(d, "local_subnet_cidr", s2c.LocalSubnet)
		mustSet(d, "remote_subnet_cidr", s2c.RemoteSubnet)
		if s2c.HAEnabled == "enabled" {
			mustSet(d, "ha_enabled", true)
		} else {
			mustSet(d, "ha_enabled", false)
		}
		mustSet(d, "remote_gateway_ip", s2c.RemoteGwIP)
		mustSet(d, "primary_cloud_gateway_name", s2c.GwName)
		mustSet(d, "local_tunnel_ip", s2c.LocalTunnelIp)
		mustSet(d, "remote_tunnel_ip", s2c.RemoteTunnelIp)
		mustSet(d, "phase1_local_identifier", s2c.Phase1LocalIdentifier)

		if s2c.HAEnabled == "enabled" {
			mustSet(d, "backup_remote_gateway_ip", s2c.RemoteGwIP2)
			mustSet(d, "backup_gateway_name", s2c.BackupGwName)
			mustSet(d, "backup_local_tunnel_ip", s2c.BackupLocalTunnelIp)
			mustSet(d, "backup_remote_tunnel_ip", s2c.BackupRemoteTunnelIp)
		}

		// Custom Mapped is a sub-type of Mapped
		if s2c.ConnType == "custom_mapped" {
			mustSet(d, "custom_mapped", true)
			s2c.ConnType = "mapped"
		} else {
			mustSet(d, "custom_mapped", false)
		}
		mustSet(d, "connection_type", s2c.ConnType)
		if s2c.ConnType == "mapped" {
			mustSet(d, "remote_subnet_virtual", s2c.RemoteSubnetVirtual)
			mustSet(d, "local_subnet_virtual", s2c.LocalSubnetVirtual)
		}

		if s2c.CustomAlgorithms {
			mustSet(d, "custom_algorithms", true)
			mustSet(d, "phase_1_authentication", s2c.Phase1Auth)
			mustSet(d, "phase_2_authentication", s2c.Phase2Auth)
			mustSet(d, "phase_1_dh_groups", s2c.Phase1DhGroups)
			mustSet(d, "phase_2_dh_groups", s2c.Phase2DhGroups)
			mustSet(d, "phase_1_encryption", s2c.Phase1Encryption)
			mustSet(d, "phase_2_encryption", s2c.Phase2Encryption)
		} else {
			mustSet(d, "custom_algorithms", false)
		}

		if s2c.PrivateRouteEncryption == "true" {
			mustSet(d, "private_route_encryption", true)
			if err := d.Set("route_table_list", s2c.RouteTableList); err != nil {
				log.Printf("[WARN] Error setting route_table_list for (%s): %s", d.Id(), err)
			}
			mustSet(d, "remote_gateway_latitude", s2c.RemoteGwLatitude)
			mustSet(d, "remote_gateway_longitude", s2c.RemoteGwLongitude)
			if s2c.HAEnabled == "enabled" {
				mustSet(d, "backup_remote_gateway_latitude", s2c.BackupRemoteGwLatitude)
				mustSet(d, "backup_remote_gateway_longitude", s2c.BackupRemoteGwLongitude)
			}
		} else {
			mustSet(d, "private_route_encryption", false)
		}

		if s2c.SslServerPool != "" {
			mustSet(d, "ssl_server_pool", s2c.SslServerPool)
		}
		mustSet(d, "enable_dead_peer_detection", s2c.DeadPeerDetection)
		mustSet(d, "enable_active_active", s2c.EnableActiveActive)
		mustSet(d, "forward_traffic_to_transit", s2c.ForwardToTransit)
		mustSet(d, "enable_event_triggered_ha", s2c.EventTriggeredHA)
		mustSet(d, "enable_single_ip_ha", s2c.EnableSingleIpHA)
		mustSet(d, "proxy_id_enabled", s2c.ProxyIdEnabled)

		if s2c.EnableIKEv2 == "true" {
			mustSet(d, "enable_ikev2", true)
		} else {
			mustSet(d, "enable_ikev2", false)
		}

		if s2c.RemoteSourceRealCIDRs != "" {
			if err := d.Set("remote_source_real_cidrs", strings.Split(s2c.RemoteSourceRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_source_real_cidrs' to state: %w", err)
			}
		}
		if s2c.RemoteSourceVirtualCIDRs != "" {
			if err := d.Set("remote_source_virtual_cidrs", strings.Split(s2c.RemoteSourceVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_source_virtual_cidrs' to state: %w", err)
			}
		}
		if s2c.RemoteDestinationRealCIDRs != "" {
			if err := d.Set("remote_destination_real_cidrs", strings.Split(s2c.RemoteDestinationRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_destination_real_cidrs' to state: %w", err)
			}
		}
		if s2c.RemoteDestinationVirtualCIDRs != "" {
			if err := d.Set("remote_destination_virtual_cidrs", strings.Split(s2c.RemoteDestinationVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'remote_destination_virtual_cidrs' to state: %w", err)
			}
		}
		if s2c.LocalSourceRealCIDRs != "" {
			if err := d.Set("local_source_real_cidrs", strings.Split(s2c.LocalSourceRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_source_real_cidrs' to state: %w", err)
			}
		}
		if s2c.LocalSourceVirtualCIDRs != "" {
			if err := d.Set("local_source_virtual_cidrs", strings.Split(s2c.LocalSourceVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_source_virtual_cidrs' to state: %w", err)
			}
		}
		if s2c.LocalDestinationRealCIDRs != "" {
			if err := d.Set("local_destination_real_cidrs", strings.Split(s2c.LocalDestinationRealCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_destination_real_cidrs' to state: %w", err)
			}
		}
		if s2c.LocalDestinationVirtualCIDRs != "" {
			if err := d.Set("local_destination_virtual_cidrs", strings.Split(s2c.LocalDestinationVirtualCIDRs, ",")); err != nil {
				return fmt.Errorf("could not write 'local_destination_virtual_cidrs' to state: %w", err)
			}
		}

		ph1RemoteId := strings.Split(s2c.Phase1RemoteIdentifier, ",")
		for i, v := range ph1RemoteId {
			ph1RemoteId[i] = strings.TrimSpace(v)
		}

		haEnabled := getBool(d, "ha_enabled")
		singleIpHA := getBool(d, "enable_single_ip_ha")
		ip := getString(d, "remote_gateway_ip")
		haIp := getString(d, "backup_remote_gateway_ip")

		if haEnabled && !singleIpHA && !(ip != "" && haIp != "" && ip == haIp) && len(ph1RemoteId) == 1 && ph1RemoteId[0] == "" {
			ph1RemoteId = append(ph1RemoteId, "")
		}
		mustSet(d, "phase1_remote_identifier", ph1RemoteId)
	}

	log.Printf("[TRACE] Reading Aviatrix Site2Cloud %s: %#v", getString(d, "connection_name"), site2cloud)
	log.Printf("[TRACE] Reading Aviatrix Site2Cloud connection_type: [%s]", getString(d, "connection_type"))

	d.SetId(site2cloud.TunnelName + "~" + site2cloud.VpcID)
	return nil
}

func resourceAviatrixSite2CloudUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	editSite2cloud := &goaviatrix.EditSite2Cloud{
		GwName:   getString(d, "primary_cloud_gateway_name"),
		VpcID:    getString(d, "vpc_id"),
		ConnName: getString(d, "connection_name"),
	}

	d.Partial(true)
	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", editSite2cloud)

	if d.HasChange("local_subnet_cidr") {
		if getBool(d, "custom_mapped") && getString(d, "local_subnet_cidr") != "" {
			return fmt.Errorf("'local_subnet_cidr' is not valid when 'custom_mapped' is enabled")
		}
		editSite2cloud.CloudSubnetCidr = getString(d, "local_subnet_cidr")
		editSite2cloud.CloudSubnetVirtual = getString(d, "local_subnet_virtual")
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud local_subnet_cidr: %w", err)
		}
	}

	if d.HasChange("local_subnet_virtual") {
		if getBool(d, "custom_mapped") && getString(d, "local_subnet_virtual") != "" {
			return fmt.Errorf("'local_subnet_virtual' is not valid when 'custom_mapped' is enabled")
		}
		if getString(d, "connection_type") == "mapped" && getString(d, "local_subnet_virtual") == "" {
			return fmt.Errorf("'local_subnet_virtual' is required for connection type: mapped, unless 'custom_mapped' is enabled")
		}
		if getString(d, "connection_type") == "unmapped" && getString(d, "local_subnet_virtual") != "" {
			return fmt.Errorf("'local_subnet_virtual' should be empty for connection type: ummapped")
		}
		editSite2cloud.CloudSubnetCidr = getString(d, "local_subnet_cidr")
		editSite2cloud.CloudSubnetVirtual = getString(d, "local_subnet_virtual")
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud local_subnet_virtual: %w", err)
		}
	}

	if d.HasChange("remote_subnet_cidr") {
		if getBool(d, "custom_mapped") && getString(d, "remote_subnet_cidr") != "" {
			return fmt.Errorf("'remote_subnet_cidr' is not valid when 'custom_mapped' is enabled")
		}
		editSite2cloud.RemoteSubnet = getString(d, "remote_subnet_cidr")
		editSite2cloud.RemoteSubnetVirtual = getString(d, "remote_subnet_virtual")
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud remote_subnet_cidr: %w", err)
		}
	}

	if d.HasChange("remote_subnet_virtual") {
		if getBool(d, "custom_mapped") && getString(d, "remote_subnet_virtual") != "" {
			return fmt.Errorf("'remote_subnet_virtual' is not valid when 'custom_mapped' is enabled")
		}
		if getString(d, "connection_type") == "mapped" && getString(d, "remote_subnet_virtual") == "" {
			return fmt.Errorf("'remote_subnet_virtual' is required for connection type: mapped, unless 'custom_mapped' is enabled")
		}
		if getString(d, "connection_type") == "unmapped" && getString(d, "remote_subnet_virtual") != "" {
			return fmt.Errorf("'remote_subnet_virtual' should be empty for connection type: ummapped")
		}
		editSite2cloud.RemoteSubnet = getString(d, "remote_subnet_cidr")
		editSite2cloud.RemoteSubnetVirtual = getString(d, "remote_subnet_virtual")
		err := client.UpdateSite2Cloud(editSite2cloud)
		if err != nil {
			return fmt.Errorf("failed to update Site2Cloud remote_subnet_virtual: %w", err)
		}
	}

	if d.HasChange("enable_dead_peer_detection") {
		s2c := &goaviatrix.Site2Cloud{
			VpcID:      getString(d, "vpc_id"),
			TunnelName: getString(d, "connection_name"),
		}
		enableDeadPeerDetection := getBool(d, "enable_dead_peer_detection")
		if enableDeadPeerDetection {
			err := client.EnableDeadPeerDetection(s2c)
			if err != nil {
				return fmt.Errorf("failed to enable deed peer detection: %w", err)
			}
		} else {
			err := client.DisableDeadPeerDetection(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable deed peer detection: %w", err)
			}
		}
	}

	if d.HasChange("enable_active_active") {
		s2c := &goaviatrix.Site2Cloud{
			VpcID:      getString(d, "vpc_id"),
			TunnelName: getString(d, "connection_name"),
		}
		activeActive := getBool(d, "enable_active_active")
		if activeActive {
			if haEnabled := getBool(d, "ha_enabled"); !haEnabled {
				return fmt.Errorf("active_active_ha can't be enabled if HA isn't enabled for site2cloud connection")
			}
			err := client.EnableSite2cloudActiveActive(s2c)
			if err != nil {
				return fmt.Errorf("failed to enable active active HA for site2cloud: %s: %w", s2c.TunnelName, err)
			}
		} else {
			err := client.DisableSite2cloudActiveActive(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable active active HA for site2cloud: %s: %w", s2c.TunnelName, err)
			}
		}
	}

	if d.HasChange("forward_traffic_to_transit") {
		s2c := &goaviatrix.Site2Cloud{
			VpcID:      getString(d, "vpc_id"),
			TunnelName: getString(d, "connection_name"),
		}
		forwardToTransit := getBool(d, "forward_traffic_to_transit")
		if forwardToTransit {
			err := client.EnableSpokeMappedSite2CloudForwarding(s2c)
			if err != nil {
				return fmt.Errorf("failed to enable traffic forwarding to transit for site2cloud: %s: %w", s2c.TunnelName, err)
			}
		} else {
			err := client.DisableSpokeMappedSite2CloudForwarding(s2c)
			if err != nil {
				return fmt.Errorf("failed to disable traffic forwarding to transit for site2cloud: %s: %w", s2c.TunnelName, err)
			}
		}
	}

	if d.HasChange("enable_event_triggered_ha") {
		if getBool(d, "enable_event_triggered_ha") {
			err := client.EnableSite2CloudEventTriggeredHA(editSite2cloud.VpcID, editSite2cloud.ConnName)
			if err != nil {
				return fmt.Errorf("could not enable event triggered HA for site2cloud during update: %w", err)
			}
		} else {
			err := client.DisableSite2CloudEventTriggeredHA(editSite2cloud.VpcID, editSite2cloud.ConnName)
			if err != nil {
				return fmt.Errorf("could not disable event triggered HA for site2cloud during update: %w", err)
			}
		}
	}

	if d.HasChanges(customMappedAttributeNames...) {
		if !getBool(d, "custom_mapped") {
			return fmt.Errorf("attributes %v are not valid when 'custom_mapped' is disabled", customMappedAttributeNames)
		}
		s2c := &goaviatrix.EditSite2Cloud{
			GwName:                        getString(d, "primary_cloud_gateway_name"),
			VpcID:                         getString(d, "vpc_id"),
			ConnName:                      getString(d, "connection_name"),
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
			return fmt.Errorf("'custom_mapped' enabled connection requires either all Remote Initiated CIDRs or all Local Initiated CIDRs be provided")
		}
		err := client.UpdateSite2Cloud(s2c)
		if err != nil {
			return fmt.Errorf("could not update site2cloud connection Remote or Local CIDRs: %w", err)
		}
	}

	if d.HasChange("phase1_remote_identifier") {
		haEnabled := getBool(d, "ha_enabled")
		singleIpHA := getBool(d, "enable_single_ip_ha")
		ip := getString(d, "remote_gateway_ip")
		haIp := getString(d, "backup_remote_gateway_ip")
		phase1RemoteIdentifier := getList(d, "phase1_remote_identifier")
		ph1RemoteIdList := goaviatrix.ExpandStringList(phase1RemoteIdentifier)

		if haEnabled && !singleIpHA && len(phase1RemoteIdentifier) != 2 {
			if !(ip != "" && haIp != "" && ip == haIp && len(phase1RemoteIdentifier) == 1) {
				return fmt.Errorf("please set two phase 1 remote IDs, when HA is enabled and single IP HA is disabled")
			}
		} else if (!haEnabled || singleIpHA) && len(phase1RemoteIdentifier) != 1 {
			return fmt.Errorf("please set one phase 1 remote ID, when HA is disabled or single IP HA is enabled")
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
			return fmt.Errorf("failed to update Site2Cloud phase 1 remote identifier: %w", err)
		}
	}

	if d.HasChanges("remote_identifier", "backup_remote_identifier") {
		haEnabled := getBool(d, "ha_enabled")
		authType := getString(d, "auth_type")
		remoteIdentifier := getString(d, "remote_identifier")
		backupRemoteIdentifier := getString(d, "backup_remote_identifier")

		if authType == "Cert" {
			if remoteIdentifier == "" {
				return fmt.Errorf("'ca_cert_tag_name' and 'remote_identifier' are both required for Cert based authentication type")
			}
			if haEnabled && backupRemoteIdentifier == "" {
				return fmt.Errorf("'backup_remote_identifier' is required for Cert based authentication type with HA enabled")
			}
		} else {
			if remoteIdentifier != "" || backupRemoteIdentifier != "" {
				return fmt.Errorf("'remote_identifier' and 'backup_remote_identifier' are both required to be empty for PSK(Pubkey) based authentication type")
			}
		}
		if d.HasChanges("remote_identifier") {
			s2c := &goaviatrix.EditSite2Cloud{
				VpcID:            getString(d, "vpc_id"),
				ConnName:         getString(d, "connection_name"),
				RemoteIdentifier: remoteIdentifier,
			}

			err := client.UpdateSite2Cloud(s2c)
			if err != nil {
				return fmt.Errorf("failed to update remote identifier: %w", err)
			}
		}
		if d.HasChanges("backup_remote_identifier") {
			s2c := &goaviatrix.EditSite2Cloud{
				VpcID:                  getString(d, "vpc_id"),
				ConnName:               getString(d, "connection_name"),
				BackupRemoteIdentifier: backupRemoteIdentifier,
			}

			err := client.UpdateSite2Cloud(s2c)
			if err != nil {
				return fmt.Errorf("failed to update backup remote identifier: %w", err)
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

	if d.HasChange("proxy_id_enabled") {
		s2c := &goaviatrix.EditSite2Cloud{
			VpcID:    getString(d, "vpc_id"),
			ConnName: getString(d, "connection_name"),
		}
		if getBool(d, "proxy_id_enabled") {
			s2c.ProxyIdEnabled = "true"
		} else {
			s2c.ProxyIdEnabled = "false"
		}
		err := client.UpdateSite2Cloud(s2c)
		if err != nil {
			return fmt.Errorf("failed to update proxy_id_enabled: %w", err)
		}
	}

	d.Partial(false)
	d.SetId(editSite2cloud.ConnName + "~" + editSite2cloud.VpcID)
	return resourceAviatrixSite2CloudRead(d, meta)
}

func resourceAviatrixSite2CloudDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	s2c := &goaviatrix.Site2Cloud{
		VpcID:      getString(d, "vpc_id"),
		TunnelName: getString(d, "connection_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix s2c: %#v", s2c)

	forwardToTransit := getBool(d, "forward_traffic_to_transit")
	if forwardToTransit {
		err := client.DisableSpokeMappedSite2CloudForwarding(s2c)
		if err != nil {
			log.Println("[WARN] Failed to disable forwarding to transit:", err)
		}
	}

	err := client.DeleteSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Site2Cloud: %w", err)
	}

	return nil
}
