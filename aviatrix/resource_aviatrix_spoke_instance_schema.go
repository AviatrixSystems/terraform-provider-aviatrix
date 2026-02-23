package aviatrix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ============================================================================
// REQUIRED SCHEMA
// ============================================================================

// spokeInstanceRequiredSchema returns the required schema attributes for spoke instance resource.
func spokeInstanceRequiredSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "UUID of the gateway group this spoke gateway belongs to.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Basic Configuration
// ============================================================================

// spokeInstanceOptionalBasicSchema returns the optional basic schema attributes for spoke instance resource.
//
//nolint:funlen
func spokeInstanceOptionalBasicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
			Description:  "Public Subnet CIDR for the gateway. Required for CSP spoke instances, not applicable for Edge.",
		},
		"gw_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Size of the gateway instance. Required for CSP spoke instances, not applicable for Edge.",
		},
		"zone": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateAzureAZ,
			Description:  "Availability Zone. Only available for Azure (8), Azure GOV (32) and Azure CHINA (2048). Must be in the form 'az-n', for example, 'az-2'.",
		},
		"allocate_new_eip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "If false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway.",
		},
		"single_az_ha": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable single AZ HA for the spoke gateway.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "A map of tags to assign to the spoke gateway.",
		},
		"private_mode_lb_vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Private Mode controller load balancer VPC ID. Required when private mode is enabled for the Controller.",
		},
		"private_mode_subnet_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Subnet availability zone for Private Mode. Required when Private Mode is enabled on the Controller and cloud_type is AWS.",
		},
		"tunnel_detection_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(20, 600),
			Description:  "The IPSec tunnel down detection time for the Spoke Gateway. Valid values: 20-600 seconds.",
		},
		"filtered_spoke_vpc_routes": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			Description: "A list of comma-separated CIDRs to be filtered from the spoke VPC route table. " +
				"When configured, filtering CIDR(s) or its subnet will be deleted from VPC routing tables as well as from spoke gateway's routing table.",
		},
		"included_advertised_spoke_routes": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "A list of comma-separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC.",
		},
		"insane_mode": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable Insane Mode for spoke gateway. Supported for AWS/AWSGov, GCP, Azure and OCI.",
		},

		// ============================================================================
		// SPOT INSTANCE (AWS and Azure)
		// ============================================================================
		"enable_spot_instance": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v := mustBool(val)
				if !v {
					errs = append(errs, fmt.Errorf("expected %s to be true to enable spot instance, got: %v", key, val))
					return warns, errs
				}
				return
			},
			Description:  "Enable spot instance. NOT supported for production deployment. Only valid for AWS and Azure.",
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

		// ============================================================================
		// BGP OVER LAN
		// ============================================================================
		"enable_bgp_over_lan": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
			Description: "Pre-allocate a network interface(eth4) for 'BGP over LAN' functionality. " +
				"Only valid for Azure (8), AzureGov (32) or AzureChina (2048).",
		},
		"bgp_lan_interfaces_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "Number of interfaces that will be created for BGP over LAN enabled Azure spoke. " +
				"Valid values: 1-5. Default value: 1.",
		},

		// ============================================================================
		// ROUTE CONFIGURATION
		// ============================================================================
		"enable_private_vpc_default_route": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Program default route in VPC private route table.",
		},
		"enable_skip_public_route_table_update": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Skip programming VPC public route table.",
		},
		"enable_monitor_gateway_subnets": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable monitor gateway subnet. This feature allows users to be notified " +
				"when an instance is deployed in the same subnet where the gateway has been deployed.",
		},
		"monitor_exclude_list": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "A set of monitored instance IDs to exclude.",
		},

		// ============================================================================
		// ENCRYPTION
		// ============================================================================
		"enable_encrypt_volume": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable EBS volume encryption for Gateway. Only supports AWS and AWSGov.",
		},
		"customer_managed_keys": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Customer managed key ID for EBS volume encryption.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - AWS Specific
// ============================================================================

// spokeInstanceOptionalAWSSchema returns AWS-specific optional schema attributes.
func spokeInstanceOptionalAWSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"insane_mode_az": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			ForceNew:    true,
			Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for AWS cloud.",
		},
		"insertion_gateway": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable to create an insertion gateway. Only valid for AWS.",
		},
		"insertion_gateway_az": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "AZ for insertion gateway. Required when insertion_gateway is enabled.",
		},
		"rx_queue_size": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"1K", "2K", "4K"}, false),
			Description:  "Gateway ethernet interface RX queue size. Valid values: 1K, 2K, 4K. Applies on HA different than spoke_gateway resource.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Azure Specific
// ============================================================================

// spokeInstanceOptionalAzureSchema returns Azure-specific optional schema attributes.
func spokeInstanceOptionalAzureSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_eip_name_resource_group": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validateAzureEipNameResourceGroup,
			Description:  "The name of the public IP address and its resource group in Azure to assign to this Spoke Gateway.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - OCI Specific
// ============================================================================

// spokeInstanceOptionalOCISchema returns OCI-specific optional schema attributes.
func spokeInstanceOptionalOCISchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"availability_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Availability domain for OCI gateway. Required for OCI.",
		},
		"fault_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Fault domain for OCI gateway. Required for OCI.",
		},
	}
}

// ============================================================================
// OPTIONAL SCHEMA - Edge Specific
// ============================================================================

// spokeInstanceOptionalEdgeSchema returns Edge-specific optional schema attributes.
func spokeInstanceOptionalEdgeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Edge interfaces
		"interfaces": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "WAN/LAN/MANAGEMENT interfaces for Edge spoke gateway.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"logical_ifname": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Logical interface name (e.g., wan0, wan1, lan0).",
					},
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Interface type: WAN, LAN, or MANAGEMENT.",
						ValidateFunc: validation.StringInSlice([]string{"WAN", "LAN", "MANAGEMENT"}, false),
					},
					"gateway_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Gateway IP address for the interface.",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Interface static IP address.",
					},
					"public_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "WAN interface public IP.",
					},
					"dhcp": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Enable DHCP for the interface.",
					},
				},
			},
		},
		// Edge interface mapping (AEP/NEO only)
		"interface_mapping": {
			Type:        schema.TypeSet,
			Optional:    true,
			ForceNew:    true,
			Description: "Interface mapping for AEP/NEO edge gateways.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Interface name.",
					},
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Interface type: WAN, LAN, or MANAGEMENT.",
						ValidateFunc: validation.StringInSlice([]string{"WAN", "LAN", "MANAGEMENT"}, false),
					},
					"index": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "Interface index.",
					},
				},
			},
		},
		// Edge ZTP file
		"ztp_file_download_path": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The location where the ZTP file will be stored. Required for Equinix, Megaport, and Self-managed edge gateways.",
		},
		"ztp_file_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"iso", "cloud-init"}, false),
			Description:  "ZTP file type. Required for Self-managed edge gateways. Valid values: 'iso', 'cloud-init'.",
		},
		// Edge device ID (AEP/NEO only)
		"device_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Device ID for AEP/NEO edge gateway.",
		},
		// Edge management egress IP prefix list
		"management_egress_ip_prefix_list": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Set of management egress gateway IP/prefix for edge gateway.",
		},
	}
}

// ============================================================================
// COMPUTED SCHEMA
// ============================================================================

// spokeInstanceComputedSchema returns the computed schema attributes for spoke instance resource.
func spokeInstanceComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gw_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Name of the spoke gateway. If not specified, a name will be auto-generated.",
		},
		"eip": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsIPAddress,
			Description:  "Elastic IP address. Required when allocate_new_eip is false.",
		},
		"security_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Security group used for the spoke gateway.",
		},
		"cloud_instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cloud instance ID of the spoke gateway.",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private IP address of the spoke gateway.",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Public IP address of the spoke gateway.",
		},
		"azure_bgp_lan_ip_list": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "List of available BGP LAN interface IPs for spoke external device connection creation. Only valid for Azure.",
		},
		"software_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			Description: "Software version of the gateway. " +
				"If set, we will attempt to update the gateway to the specified version.",
		},
		"image_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			Description: "Image version of the gateway. " +
				"If set, we will attempt to update the gateway to the specified version.",
		},
	}
}
