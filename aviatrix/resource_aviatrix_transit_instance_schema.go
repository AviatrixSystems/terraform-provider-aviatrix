package aviatrix

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ============================================================================
// REQUIRED SCHEMA
// ============================================================================

// transitInstanceRequiredSchema returns the required schema attributes for transit instance resource.
func transitInstanceRequiredSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "UUID of the transit group this instance belongs to. The cloud_type, account_name and vpc_id are derived from this group.",
		},
		"gw_size": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Size of the gateway instance.",
			DiffSuppressFunc: func(_, old, _ string, _ *schema.ResourceData) bool {
				// Suppress the diff if the old value is "UNKNOWN"
				return old == "UNKNOWN"
			},
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Basic Configuration
// ============================================================================

// transitInstanceOptionalBasicSchema returns the optional basic schema attributes for transit instance resource.
//
//nolint:funlen
func transitInstanceOptionalBasicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gw_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Name of the gateway which is going to be created. If not provided, it will be auto-generated.",
		},
		"subnet": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
			Description:  "Public Subnet Name. Required for CSP transit gateways.",
		},
		"allocate_new_eip": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// Suppress diff when private mode is enabled (allocate_new_eip is not applicable)
				return getString(d, "private_mode_lb_vpc_id") != ""
			},
			Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
				"Otherwise, allocate a new Elastic IP and use it for this gateway.",
		},
		"eip": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsIPAddress,
			Description:  "Required when allocate_new_eip is false. It uses specified EIP for this gateway.",
		},
		"single_az_ha": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Set to 'enabled' if this feature is desired.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "A map of tags to assign to the transit gateway.",
		},
		"tunnel_detection_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(20, 600),
			Description:  "The IPSec tunnel down detection time for the transit gateway.",
		},
		"private_mode_lb_vpc_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Description:   "Private Mode Controller load balancer VPC ID. Required when private mode is enabled for the Controller.",
			ConflictsWith: []string{"allocate_new_eip"},
		},
		"private_mode_subnet_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Private Mode subnet availability zone.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Route Configuration
// ============================================================================

// transitInstanceOptionalRouteSchema returns route-related optional schema attributes.
func transitInstanceOptionalRouteSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customized_spoke_vpc_routes": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			Description: "A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, " +
				"it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. " +
				"It applies to all spoke gateways attached to this transit gateway.",
		},
		"filtered_spoke_vpc_routes": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, " +
				"filtering CIDR(s) or it's subnet will be deleted from VPC routing tables as well as from spoke gateway's " +
				"routing table. It applies to all spoke gateways attached to this transit gateway.",
		},
		"excluded_advertised_spoke_routes": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			Description: "A list of comma separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'. " +
				"When configured, it inspects all the advertised CIDRs from its spoke gateways and " +
				"remove those included in the 'Excluded CIDR List'.",
		},
		"customized_transit_vpc_routes": {
			Type:     schema.TypeSet,
			Optional: true,
			Description: "A list of CIDRs to be customized for the transit VPC routes. " +
				"When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs." +
				"To be effective, `enable_advertise_transit_cidr` or firewall management access for a transit firenet gateway must be enabled.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"bgp_manual_spoke_advertise_cidrs": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "",
			DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
			Description:      "Intended CIDR list to be advertised to external bgp router. Does not require enable_bgp = true.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Feature Flags
// ============================================================================

// transitInstanceOptionalFeatureSchema returns feature flag optional schema attributes.
func transitInstanceOptionalFeatureSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_transit_firenet": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Specify whether to enable transit firenet interfaces or not.",
		},
		"enable_firenet": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Specify whether to enable firenet interfaces or not.",
		},
		"lan_vpc_id": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: DiffSuppressFuncGCPVpcId,
			Description:      "LAN VPC ID. Only used for GCP Transit FireNet.",
		},
		"lan_private_subnet": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "LAN Private Subnet. Only used for GCP Transit FireNet.",
		},
		"enable_gateway_load_balancer": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable firenet interfaces with AWS Gateway Load Balancer. Only valid when `enable_firenet` or `enable_transit_firenet`" +
				" are set to true and `cloud_type` = 1 (AWS). Currently AWS Gateway Load Balancer is only supported " +
				"in AWS regions us-west-2 and us-east-1. Valid values: true or false. Default value: false.",
		},
		"enable_bgp_over_lan": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Pre-allocate a network interface(eth4) for \"BGP over LAN\" functionality. Only valid for cloud_type = 4 (GCP) and 8 (Azure). Valid values: true or false. Default value: false. Available as of provider version R2.18+. Updatable as of provider version 3.0.3+.",
		},
		"bgp_lan_interfaces_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Number of interfaces that will be created for BGP over LAN enabled Azure transit. Applies on HA Transit as well if enabled. Updatable as of provider version 3.0.3+.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Spot Instance
// ============================================================================

// transitInstanceOptionalSpotSchema returns spot instance optional schema attributes.
func transitInstanceOptionalSpotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_spot_instance": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(bool)
				if !ok {
					errs = append(errs, fmt.Errorf("expected %s to be a bool, got: %T", key, val))
					return warns, errs
				}
				if !v {
					errs = append(errs, fmt.Errorf("expected %s to true to enable spot instance, got: %v", key, val))
					return warns, errs
				}
				return
			},
			Description:  "Enable spot instance. NOT supported for production deployment. Supported for AWS and Azure.",
			RequiredWith: []string{"spot_price"},
		},
		"spot_price": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Price for spot instance. NOT supported for production deployment.",
			RequiredWith: []string{"enable_spot_instance"},
		},
		"delete_spot": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Description: "If set true, the spot instance will be deleted on eviction. Otherwise, the instance will be deallocated on eviction. Only supports Azure.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - AWS Specific
// ============================================================================

// transitInstanceOptionalAWSSchema returns AWS-specific optional schema attributes.
func transitInstanceOptionalAWSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"insane_mode_az": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "AZ of subnet being created for Insane Mode Transit Gateway. Required for AWS if insane_mode is enabled.",
		},
		"rx_queue_size": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"1K", "2K", "4K", "8K", "16K"}, false),
			Description:  "Gateway ethernet interface RX queue size. Supported for AWS related clouds only. Applies on HA as well if enabled.",
		},
		"enable_monitor_gateway_subnets": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable monitor gateway subnets. " +
				"Only valid for cloud_type = 1 (AWS) or 256 (AWSGov). Valid values: true, false. Default value: false.",
		},
		"monitor_exclude_list": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A set of monitored instance ids. Only valid when 'enable_monitor_gateway_subnets' = true.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Azure Specific
// ============================================================================

// transitInstanceOptionalAzureSchema returns Azure-specific optional schema attributes.
func transitInstanceOptionalAzureSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateAzureAZ,
			Description:  "Availability Zone. Required for Azure (8), Azure GOV (32) and Azure CHINA (2048). Must be in the form 'az-n', for example, 'az-2'.",
		},
		"azure_eip_name_resource_group": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "The name of the public IP address and its resource group in Azure to assign to this Transit Gateway.",
			ValidateFunc: validateAzureEipNameResourceGroup,
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - OCI Specific
// ============================================================================

// transitInstanceOptionalOCISchema returns OCI-specific optional schema attributes.
func transitInstanceOptionalOCISchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"availability_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Availability domain for OCI.",
		},
		"fault_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Fault domain for OCI.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Edge Specific
// ============================================================================

// transitInstanceOptionalEdgeSchema returns Edge-specific optional schema attributes.
//
//nolint:funlen
func transitInstanceOptionalEdgeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ztp_file_download_path": {
			Type:     schema.TypeString,
			Optional: true,
			DiffSuppressFunc: func(_, old, _ string, _ *schema.ResourceData) bool {
				return old != ""
			},
			Description: "The location where the ZTP file will be stored locally. For Equinix/Megaport/Self-managed edge transit.",
		},
		"ztp_file_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "ZTP file type. For Self-managed edge transit.",
			ValidateFunc: validation.StringInSlice([]string{"iso", "cloud-init"}, false),
		},
		"device_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Device ID for AEP edge transit gateway.",
		},
		"interfaces": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A set of WAN/Management interfaces for edge transit gateway.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"logical_ifname": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Logical interface name e.g., wan0, wan1, mgmt0.",
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^(wan|mgmt)[0-9]+$`),
							"Logical interface name must start with 'wan', or 'mgmt' followed by a number (e.g., 'wan0', 'mgmt0').",
						),
					},
					"gateway_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The gateway IP address associated with this interface.",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The static IP address assigned to this interface.",
					},
					"public_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The public IP address associated with this interface (if applicable).",
					},
					"dhcp": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether DHCP is enabled on this interface.",
					},
					"secondary_private_cidr_list": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "A list of secondary private CIDR blocks associated with this interface.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"underlay_cidr": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "The underlay CIDR in the format of ipaddr/netmask for this interface.",
						ValidateFunc: validation.IsCIDR,
					},
				},
			},
		},
		"interface_mapping": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of interface names mapped to interface types and indices. Only required for ESXI. For self-managed edge transit.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Interface name (e.g., 'eth0', 'eth1').",
					},
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Interface type (e.g., 'WAN', 'MANAGEMENT').",
						ValidateFunc: validation.StringInSlice([]string{"WAN", "MANAGEMENT"}, false),
					},
					"index": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "Interface index (e.g., 0, 1).",
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},
		"peer_connection_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Connection type for the edge transit gateway (e.g., 'public', 'private').",
			ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
		},
		"peer_backup_logical_ifname": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Peer backup logical interface name for the edge transit gateway (e.g., 'wan0', 'wan1').",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"eip_map": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"logical_ifname": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Logical interface name e.g., wan0, mgmt0.",
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^(wan|mgmt)[0-9]+$`),
							"Logical interface name must start with 'wan', or 'mgmt' followed by a number (e.g., 'wan0', 'mgmt0').",
						),
					},
					"private_ip": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The private IP address associated with the interface.",
					},
					"public_ip": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The public IP address associated with the interface.",
					},
				},
			},
			Description: "A list of mappings between interface names and their associated private and public IPs.",
		},
		"management_egress_ip_prefix_list": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Set of management egress gateway IP/prefix.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// ============================================================================
// COMPUTED SCHEMA
// ============================================================================

// transitInstanceComputedSchema returns the computed schema attributes for transit instance resource.
func transitInstanceComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the transit group. Derived from the transit group.",
		},
		"cloud_type": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Type of cloud service provider. Derived from the transit group.",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of a Cloud-Account in Aviatrix controller. Derived from the transit group.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "VPC-ID/VNet-Name/Site-ID of cloud provider. Derived from the transit group.",
		},
		"gateway_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "UUID of the transit gateway.",
		},
		"cloud_instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Instance ID of the transit gateway.",
		},
		"security_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Security group used for the transit gateway.",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private IP address of the transit gateway created.",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Public IP address of the Transit Gateway created.",
		},
		"lan_interface_cidr": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Transit gateway lan interface cidr.",
		},
		"bgp_lan_ip_list": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
			Description: "List of available BGP LAN interface IPs for transit external device connection creation. " +
				"Only supports GCP and Azure. Available as of provider version R2.21.0+.",
		},
		"azure_bgp_lan_ip_list": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
			Description: "List of available BGP LAN interface IPs for Azure transit external device connection creation. " +
				"Only supports Azure. Available as of provider version R2.21.0+.",
		},
		"software_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			Description: "software_version can be used to set the desired software version of the gateway. " +
				"If set, we will attempt to update the gateway to the specified version. " +
				"If left blank, the gateway software version will continue to be managed through the aviatrix_controller_config resource.",
		},
		"image_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			Description: "image_version can be used to set the desired image version of the gateway. " +
				"If set, we will attempt to update the gateway to the specified version.",
		},
	}
}

// ============================================================================
// COMBINED SCHEMA
// ============================================================================

// transitInstanceSchema returns the complete schema for transit instance resource.
func transitInstanceSchema() map[string]*schema.Schema {
	schemaMap := make(map[string]*schema.Schema)

	// Merge all schema functions
	for k, v := range transitInstanceRequiredSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalBasicSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalRouteSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalFeatureSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalSpotSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalAWSSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalAzureSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalOCISchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceOptionalEdgeSchema() {
		schemaMap[k] = v
	}
	for k, v := range transitInstanceComputedSchema() {
		schemaMap[k] = v
	}

	return schemaMap
}
