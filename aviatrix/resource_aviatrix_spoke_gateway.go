package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const subnetSeparator = "~~"

func resourceAviatrixSpokeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeGatewayCreate,
		Read:   resourceAviatrixSpokeGatewayRead,
		Update: resourceAviatrixSpokeGatewayUpdate,
		Delete: resourceAviatrixSpokeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		// CustomizeDiff handles custom diff logic during plan operations:
		// - Forces resource recreation when IPv6 subnet fields change (if previously set and enable_ipv6 is true)
		CustomizeDiff: resourceAviatrixSpokeGatewayCustomizeDiff,

		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixSpokeGatewayResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixSpokeGatewayStateUpgradeV0,
				Version: 0,
			},
			{
				Type:    resourceAviatrixSpokeGatewayResourceV1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixSpokeGatewayStateUpgradeV1,
				Version: 1,
			},
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "Type of cloud service provider.",
				ValidateFunc: validateCloudType,
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the gateway which is going to be created.",
			},
			"vpc_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "VPC-ID/VNet-Name of cloud provider.",
				DiffSuppressFunc: DiffSuppressFuncGatewayVpcId,
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of cloud provider.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size of the gateway instance.",
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "Public Subnet Info.",
			},
			"subnet_ipv6_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIPv6CIDR,
				// DiffSuppressFunc ignores changes to this field when enable_ipv6 is false or cloud_type is GCP
				// This prevents unnecessary diffs for a field that is not used in that configuration
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					return !getBool(d, "enable_ipv6") || goaviatrix.IsCloudType(getInt(d, "cloud_type"), goaviatrix.GCPRelatedCloudTypes)
				},
				Description: "IPv6 CIDR for the subnet. Only used if enable_ipv6 flag is set. Currently only supported on Azure and AWS Cloud.",
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureAZ,
				Description:  "Availability Zone. Only available for Azure (8), Azure GOV (32) and Azure CHINA (2048). Must be in the form 'az-n', for example, 'az-2'.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for AWS cloud.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable Source NAT feature in 'single_ip' mode on the gateway or not.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return getBool(d, "enable_private_oob") || getString(d, "private_mode_lb_vpc_id") != ""
				},
				Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"ha_subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "HA Subnet. Required if enabling HA for AWS/AWSGov/AWSChina/Azure/AzureChina/OCI/Alibaba Cloud. Optional if enabling HA for GCP.",
			},
			"ha_subnet_ipv6_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIPv6CIDR,
				// DiffSuppressFunc ignores changes to this field when enable_ipv6 is false or cloud_type is GCP
				// This prevents unnecessary diffs for a field that is not used in that configuration
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					return !getBool(d, "enable_ipv6") || goaviatrix.IsCloudType(getInt(d, "cloud_type"), goaviatrix.GCPRelatedCloudTypes)
				},
				Description: "IPv6 CIDR for the HA subnet. Only used if enable_ipv6 flag is set. Currently only supported on Azure and AWS Cloud.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Zone. Required if enabling HA for GCP. Optional for Azure.",
			},
			"ha_insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Spoke HA Gateway. Required for AWS if insane_mode is true and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Gateway Size.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"manage_ha_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "This parameter is a switch used to determine whether or not to manage spoke ha gateway " +
					"using the aviatrix_spoke_gateway resource. If this is set to false, managing spoke ha gateway " +
					"must be done using the aviatrix_spoke_ha_gateway resource. Valid values: true, false. Default value: true.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable Insane Mode for Spoke Gateway. Valid values: true, false. Supported for AWS/AWSGov, GCP, Azure and OCI. If insane mode is enabled, gateway size has to at least be c5 size for AWS and Standard_D3_v2 size for Azure.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Valid values: true, false.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.",
			},
			"enable_preserve_as_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable preserve as_path when advertising manual summary cidrs on BGP spoke gateway.",
			},
			"customized_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, " +
					"it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. " +
					"It applies to this spoke gateway only.",
			},
			"filtered_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, " +
					"filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s " +
					"routing table. It applies to this spoke gateway only.",
			},
			"included_advertised_spoke_routes": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC.",
			},
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			"enable_monitor_gateway_subnets": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable [monitor gateway subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet). " +
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
			"enable_private_oob": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable private OOB.",
			},
			"oob_management_subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "OOB management subnet.",
			},
			"oob_availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OOB subnet availability zone.",
			},
			"ha_oob_management_subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "OOB HA management subnet.",
			},
			"ha_oob_availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OOB HA availability zone.",
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable jumbo frame support for spoke gateway. Valid values: true or false. Default value: true.",
			},
			"enable_gro_gso": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specify whether to disable GRO/GSO or not.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A map of tags to assign to the spoke gateway.",
			},
			"enable_private_vpc_default_route": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Config Private VPC Default Route.",
			},
			"enable_skip_public_route_table_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip Public Route Table Update.",
			},
			"enable_auto_advertise_s2c_cidrs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically advertise remote CIDR to Aviatrix Transit Gateway when route based Site2Cloud Tunnel is created.",
			},
			"spoke_bgp_manual_advertise_cidrs": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "Intended CIDR list to be advertised to external BGP router.",
			},
			"enable_bgp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable BGP. Default: false.",
			},
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable learned CIDR approval for BGP Spoke Gateway. Valid values: true, false.",
			},
			"learned_cidrs_approval_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      defaultLearnedCidrApprovalMode,
				ValidateFunc: validation.StringInSlice([]string{"gateway"}, false),
				Description: "Set the learned CIDRs approval mode for BGP Spoke Gateway. Only valid when 'enable_learned_cidrs_approval' is " +
					"set to true. Currently, only 'gateway' is supported: learned CIDR approval applies to " +
					"ALL connections. Default value: 'gateway'.",
			},
			"approved_learned_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: ValidateCIDRRule,
				},
				Optional:    true,
				Description: "Approved learned CIDRs for BGP Spoke Gateway. Available as of provider version R2.21+.",
			},
			"bgp_ecmp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Equal Cost Multi Path (ECMP) routing for the next hop for BGP Spoke Gateway.",
			},
			"enable_active_standby": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Active-Standby Mode, available only with HA enabled for BGP Spoke Gateway.",
			},
			"enable_active_standby_preemptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.",
			},
			"disable_route_propagation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disables route propagation on BGP Spoke to attached Transit Gateway. Default: false.",
			},
			"private_mode_lb_vpc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "Private Mode controller load balancer vpc_id.  Required when private mode is enabled for the Controller.",
				ConflictsWith: []string{"allocate_new_eip"},
			},
			"private_mode_subnet_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Subnet availability zone. Required when Private Mode is enabled on the Controller and cloud_type is AWS.",
			},
			"ha_private_mode_subnet_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: " Private Mode HA subnet availability zone.",
			},
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Changes the Aviatrix BGP Spoke Gateway ASN number before you setup Aviatrix BGP Spoke Gateway connection configurations.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"prepend_as_path": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"local_as_number"},
				Description:  "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices. Only valid for BGP Spoke Gateway",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
			},
			"bgp_polling_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validation.IntBetween(10, 50),
				Description:  "BGP route polling time for BGP Spoke Gateway. Unit is in seconds. Valid values are between 10 and 50.",
			},
			"bgp_neighbor_status_polling_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpNeighborStatusPollingTime,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "BGP neighbor status polling time. Unit is in seconds. Valid values are between 1 and 10.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "0"
				},
			},
			"bgp_hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpHoldTime,
				ValidateFunc: validation.IntBetween(12, 360),
				Description:  "BGP Hold Time for BGP Spoke Gateway. Unit is in seconds. Valid values are between 12 and 360.",
			},
			"enable_bgp_over_lan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "Pre-allocate a network interface(eth4) for \"BGP over LAN\" functionality. " +
					"Only valid for 8 (Azure), 32 (AzureGov) or AzureChina (2048). Valid values: true or false. " +
					"Default value: false. Available as of provider version R3.0.2+.",
			},
			"bgp_lan_interfaces_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description: "Number of interfaces that will be created for BGP over LAN enabled Azure spoke. " +
					"Only valid for 8 (Azure), 32 (AzureGov) or AzureChina (2048). Default value: 1. " +
					"Available as of provider version R3.0.2+.",
			},
			"enable_spot_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := mustBool(val)
					if !v {
						errs = append(errs, fmt.Errorf("expected %s to true to enable spot instance, got: %v", key, val))
						return warns, errs
					}
					return
				},
				Description:  "Enable spot instance. NOT supported for production deployment.",
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
			"rx_queue_size": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1K", "2K", "4K", "8K", "16K"}, false),
				Description:  "Gateway ethernet interface RX queue size. Supported for AWS related clouds only. Applies on HA as well if enabled.",
			},
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
			"ha_availability_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "HA availability domain for OCI.",
			},
			"ha_fault_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "HA fault domain for OCI.",
			},
			"eip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Required when allocate_new_eip is false. It uses specified EIP for this gateway.",
			},
			"ha_eip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Public IP address that you want assigned to the HA Spoke Gateway.",
			},
			"azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to this Spoke Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"ha_azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to the HA Spoke Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"tunnel_detection_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(20, 600),
				Description:  "The IPSec tunnel down detection time for the Spoke Gateway.",
			},
			"software_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "software_version can be used to set the desired software version of the gateway. " +
					"If set, we will attempt to update the gateway to the specified version. " +
					"If left blank, the gateway software version will continue to be managed through the aviatrix_controller_config resource.",
			},
			"ha_software_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "ha_software_version can be used to set the desired software version of the HA gateway. " +
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
			"ha_image_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "ha_image_version can be used to set the desired image version of the HA gateway. " +
					"If set, we will attempt to update the gateway to the specified version.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the spoke gateway.",
			},
			"ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA security group used for the spoke gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the spoke gateway created.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of HA spoke gateway.",
			},
			"ha_gw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix spoke gateway unique name of HA spoke gateway.",
			},
			"ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the spoke gateway created.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Spoke Gateway created.",
			},
			"ha_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the HA Spoke Gateway.",
			},
			"bgp_lan_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "List of available BGP LAN interface IPs for spoke external device connection creation. " +
					"Only supports 8 (Azure), 32 (AzureGov) or AzureChina (2048). Available as of provider version R3.0.2+.",
			},
			"ha_bgp_lan_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "List of available BGP LAN interface IPs for spoke external device HA connection creation. " +
					"Only supports 8 (Azure), 32 (AzureGov) or AzureChina (2048). Available as of provider version R3.0.2+.",
			},
			"enable_global_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true to enable global VPC. Only supported for GCP.",
			},
			"bgp_send_communities": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "BGP communities gateway send configuration.",
				Default:     false,
			},
			"bgp_accept_communities": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "BGP communities gateway accept configuration.",
				Default:     false,
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable IPv6 for the gateway. Only supported for AWS (1), Azure (8).",
			},
			"insertion_gateway": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				Description:   "Enable insertion gateway mode.",
				ConflictsWith: []string{"insane_mode"},
			},
			"insertion_gateway_az": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				Description:  "AZ of subnet being created for Insertion Gateway. Required if insertion_gateway is enabled.",
				RequiredWith: []string{"insertion_gateway"},
			},
			"tunnel_encryption_cipher": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Encryption ciphers for gateway peering tunnels. Config options are default (AES-126-GCM-96) or strong (AES-256-GCM-96).",
				ValidateFunc: validation.StringInSlice([]string{"default", "strong"}, false),
				Default:      "default",
			},
			"tunnel_forward_secrecy": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Perfect Forward Secrecy (PFS) for gateway peering tunnels. Config Options are enable/disable.",
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
				Default:      "disable",
			},
			"private_route_table_config": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Set of Azure route table selectors to treat as private route tables for the spoke VNet. Each entry is in the format \"<route_table_name>:<resource_group_name>\". Only applicable for Azure (8), AzureGov (32) and AzureChina (2048).",
			},
		},
	}
}

func handleIPv6SubnetForceNew(d *schema.ResourceDiff, fieldName string) error {
	if !d.HasChange(fieldName) || !getBool(d, "enable_ipv6") || goaviatrix.IsCloudType(getInt(d, "cloud_type"), goaviatrix.GCPRelatedCloudTypes) {
		return nil
	}

	oldSubnet, newSubnet := d.GetChange(fieldName)
	oldSubnetStr, newSubnetStr := mustString(oldSubnet), mustString(newSubnet)

	if oldSubnetStr != "" && oldSubnetStr != newSubnetStr {
		return d.ForceNew(fieldName)
	}

	return nil
}

func resourceAviatrixSpokeGatewayCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Only force recreation for primary gateway's IPv6 CIDR changes
	// HA gateway IPv6 CIDR changes are handled by Update function (recreates only HA gateway)
	if err := handleIPv6SubnetForceNew(d, "subnet_ipv6_cidr"); err != nil {
		return err
	}

	return nil
}

func resourceAviatrixSpokeGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.SpokeVpc{
		CloudType:              getInt(d, "cloud_type"),
		AccountName:            getString(d, "account_name"),
		GwName:                 getString(d, "gw_name"),
		VpcSize:                getString(d, "gw_size"),
		Subnet:                 getString(d, "subnet"),
		HASubnet:               getString(d, "ha_subnet"),
		AvailabilityDomain:     getString(d, "availability_domain"),
		FaultDomain:            getString(d, "fault_domain"),
		ApprovedLearnedCidrs:   getStringSet(d, "approved_learned_cidrs"),
		EnableGlobalVpc:        getBool(d, "enable_global_vpc"),
		TunnelEncryptionCipher: getString(d, "tunnel_encryption_cipher"),
		TunnelForwardSecrecy:   getString(d, "tunnel_forward_secrecy"),
	}

	if !getBool(d, "manage_ha_gateway") {
		haSubnet := getString(d, "ha_subnet")
		haZone := getString(d, "ha_zone")
		haInsaneModeAz := getString(d, "ha_insane_mode_az")
		haEip := getString(d, "ha_eip")
		haAzureEipNameResourceGroup := getString(d, "ha_azure_eip_name_resource_group")
		haGwSize := getString(d, "ha_gw_size")
		haAvailabilityDomain := getString(d, "ha_availability_domain")
		haFaultDomain := getString(d, "ha_fault_domain")
		haOobManagementSubnet := d.Get("ha_oob_management_subnet")
		haPrivateModeSubnetZone := d.Get("ha_private_mode_subnet_zone")
		haOobAvailabilityZone := d.Get("ha_oob_availability_zone")
		haSoftwareVersion := d.Get("ha_software_version")
		haOobImageVersion := d.Get("ha_image_version")
		if haSubnet != "" || haZone != "" || haInsaneModeAz != "" || haEip != "" || haAzureEipNameResourceGroup != "" ||
			haGwSize != "" || haAvailabilityDomain != "" || haFaultDomain != "" || haOobManagementSubnet != "" ||
			haPrivateModeSubnetZone != "" || haOobAvailabilityZone != "" || haSoftwareVersion != "" || haOobImageVersion != "" {
			return fmt.Errorf("'manage_ha_gateway' is set to false. Please set it to true, or use 'aviatrix_spoke_ha_gateway' to manage spoke ha gateway")
		}
	}

	if getBool(d, "enable_private_vpc_default_route") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_private_vpc_default_route is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if getBool(d, "enable_skip_public_route_table_update") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_skip_public_route_update is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if _, hasSetZone := d.GetOk("zone"); !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && hasSetZone {
		return fmt.Errorf("attribute 'zone' is only valid for Azure (8), Azure GOV (32) and Azure CHINA (2048)")
	}

	if _, hasSetZone := d.GetOk("zone"); goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && hasSetZone {
		gateway.Subnet = fmt.Sprintf("%s~~%s~~", getString(d, "subnet"), getString(d, "zone"))
	}

	enableSNat := getBool(d, "single_ip_snat")
	if enableSNat {
		gateway.EnableNat = "yes"
	}

	singleAZ := getBool(d, "single_az_ha")
	if singleAZ {
		gateway.SingleAzHa = "enabled"
	} else {
		gateway.SingleAzHa = "disabled"
	}

	enableBgp := getBool(d, "enable_bgp")
	disableRoutePropagation := getBool(d, "disable_route_propagation")
	if enableBgp {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWS|goaviatrix.Azure) {
			return fmt.Errorf("enabling BGP is only supported for AWS (1) and Azure (8)")
		}
		gateway.EnableBgp = "yes"
	} else {
		if disableRoutePropagation {
			return fmt.Errorf("disable route propagation is not supported on Non-BGP Spoke")
		}
	}

	learnedCidrsApproval := getBool(d, "enable_learned_cidrs_approval")
	if !learnedCidrsApproval && len(gateway.ApprovedLearnedCidrs) != 0 {
		return fmt.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		gateway.VpcID = getString(d, "vpc_id")
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gateway.VpcRegion = getString(d, "vpc_reg")
	} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = getString(d, "vpc_reg")
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192) or, AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain == "" || gateway.FaultDomain == "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
	}
	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain != "" || gateway.FaultDomain != "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
	}

	insaneMode := getBool(d, "insane_mode")
	insaneModeAz := getString(d, "insane_mode_az")
	haSubnet := getString(d, "ha_subnet")
	haZone := getString(d, "ha_zone")
	haAvailabilityDomain := getString(d, "ha_availability_domain")
	haFaultDomain := getString(d, "ha_fault_domain")

	if haZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("'ha_zone' is only valid for GCP (4), Azure (8), AzureGov (32) and AzureChina (2048) providers if enabling HA")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && haSubnet != "" && haZone == "" {
		return fmt.Errorf("'ha_zone' must be set to enable HA on GCP (4), cannot enable HA with only 'ha_subnet'")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" && haZone != "" {
		return fmt.Errorf("'ha_subnet' must be provided to enable HA on Azure (4), AzureGov (32) or AzureChina (2048), cannot enable HA with only 'ha_zone'")
	}
	haGwSize := getString(d, "ha_gw_size")
	if haSubnet == "" && haZone == "" && haGwSize != "" {
		return fmt.Errorf("'ha_gw_size' is only required if enabling HA")
	}
	haInsaneModeAz := getString(d, "ha_insane_mode_az")
	if insaneMode {
		// Insane Mode encryption is not supported in Azure China regions
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|
			goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			return fmt.Errorf("insane_mode is only supported for AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}

		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			if insaneModeAz == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}
			if haSubnet != "" && haInsaneModeAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768) provider and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			strs = append(strs, gateway.Subnet, insaneModeAz)
			gateway.Subnet = strings.Join(strs, subnetSeparator)
		}
		gateway.InsaneMode = "yes"
	} else {
		gateway.InsaneMode = "no"
	}
	if haZone != "" || haSubnet != "" {
		if haGwSize == "" {
			return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set")
		}
	}
	if haSubnet != "" {
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain == "" || haFaultDomain == "") {
			return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are required to enable Peering HA on OCI")
		}
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain != "" || haFaultDomain != "") {
			return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are only valid for OCI")
		}
	}

	enableEncryptVolume := getBool(d, "enable_encrypt_volume")
	customerManagedKeys := getString(d, "customer_managed_keys")
	if enableEncryptVolume && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}
	if customerManagedKeys != "" {
		if !enableEncryptVolume {
			return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
		}
		gateway.CustomerManagedKeys = customerManagedKeys
	}
	if !enableEncryptVolume && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		gateway.EncVolume = "no"
	}

	enableMonitorSubnets := getBool(d, "enable_monitor_gateway_subnets")
	var excludedInstances []string
	for _, v := range getSet(d, "monitor_exclude_list").List() {
		excludedInstances = append(excludedInstances, mustString(v))
	}
	// Enable monitor gateway subnets does not work with AWSChina
	if enableMonitorSubnets && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes^goaviatrix.AWSChina) {
		return fmt.Errorf("'enable_monitor_gateway_subnets' is only valid for AWS (1), AWSGov (256), AWS Top Secret (16384) or AWS Secret (32768)")
	}
	if !enableMonitorSubnets && len(excludedInstances) != 0 {
		return fmt.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
	}

	bgpOverLan := getBool(d, "enable_bgp_over_lan")
	if bgpOverLan && !enableBgp {
		return fmt.Errorf("'enable_bgp' is required to be true to enable bgp over lan")
	}
	if bgpOverLan && !(goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes)) {
		return fmt.Errorf("'enable_bgp_over_lan' is only valid for Azure (8), AzureGov (32) or AzureChina (2048)")
	}
	bgpLanInterfacesCount, isCountSet := d.GetOk("bgp_lan_interfaces_count")
	if isCountSet && (!bgpOverLan || !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes)) {
		return fmt.Errorf("'bgp_lan_interfaces_count' is only valid for BGP over LAN enabled spoke for Azure (8), AzureGov (32) or AzureChina (2048)")
	} else if !isCountSet && bgpOverLan && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("please specify 'bgp_lan_interfaces_count' for BGP over LAN enabled Azure spoke: %s", gateway.GwName)
	}
	if bgpOverLan {
		gateway.BgpOverLan = true
		gateway.BgpLanInterfacesCount = mustInt(bgpLanInterfacesCount)
	}

	enablePrivateOob := getBool(d, "enable_private_oob")
	oobManagementSubnet := getString(d, "oob_management_subnet")
	oobAvailabilityZone := getString(d, "oob_availability_zone")
	haOobManagementSubnet := getString(d, "ha_oob_management_subnet")
	haOobAvailabilityZone := getString(d, "ha_oob_availability_zone")

	if enablePrivateOob {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'enable_private_oob' is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
		}

		if oobAvailabilityZone == "" {
			return fmt.Errorf("\"oob_availability_zone\" is required if \"enable_private_oob\" is true")
		}

		if oobManagementSubnet == "" {
			return fmt.Errorf("\"oob_management_subnet\" is required if \"enable_private_oob\" is true")
		}

		if haSubnet != "" {
			if haOobAvailabilityZone == "" {
				return fmt.Errorf("\"ha_oob_availability_zone\" is required if \"enable_private_oob\" is true and \"ha_subnet\" is provided")
			}

			if haOobManagementSubnet == "" {
				return fmt.Errorf("\"ha_oob_management_subnet\" is required if \"enable_private_oob\" is true and \"ha_subnet\" is provided")
			}
		} else {
			if haOobAvailabilityZone != "" {
				return fmt.Errorf("\"ha_oob_availability_zone\" must be empty if \"ha_subnet\" is empty")
			}

			if haOobManagementSubnet != "" {
				return fmt.Errorf("\"ha_oob_management_subnet\" must be empty if \"ha_subnet\" is empty")
			}
		}

		gateway.EnablePrivateOob = "on"
		gateway.Subnet = gateway.Subnet + subnetSeparator + oobAvailabilityZone
		gateway.OobManagementSubnet = oobManagementSubnet + subnetSeparator + oobAvailabilityZone
	} else {
		if oobAvailabilityZone != "" {
			return fmt.Errorf("\"oob_availability_zone\" must be empty if \"enable_private_oob\" is false")
		}

		if oobManagementSubnet != "" {
			return fmt.Errorf("\"oob_management_subnet\" must be empty if \"enable_private_oob\" is false")
		}

		if haOobAvailabilityZone != "" {
			return fmt.Errorf("\"ha_oob_availability_zone\" must be empty if \"enable_private_oob\" is false")
		}

		if haOobManagementSubnet != "" {
			return fmt.Errorf("\"ha_oob_management_subnet\" must be empty if \"enable_private_oob\" is false")
		}
	}

	_, tagsOk := d.GetOk("tags")
	if tagsOk {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return errors.New("failed to create spoke gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) or AWS Secret (32768)")
		}

		tagsMap, err := extractTags(d, gateway.CloudType)
		if err != nil {
			return fmt.Errorf("error creating tags for spoke gateway: %w", err)
		}
		tagJson, err := TagsMapToJson(tagsMap)
		if err != nil {
			return fmt.Errorf("failed to add tags whenc creating spoke gateway: %w", err)
		}
		gateway.TagJson = tagJson
	}

	enableActiveStandby := getBool(d, "enable_active_standby")
	if haSubnet == "" && haZone == "" && enableActiveStandby {
		return fmt.Errorf("could not configure Active-Standby as HA is not enabled")
	}
	if !enableBgp && enableActiveStandby {
		return fmt.Errorf("could not configure Active-Standby as it is not BGP capable gateway")
	}
	enableActiveStandbyPreemptive := getBool(d, "enable_active_standby_preemptive")
	if !enableActiveStandby && enableActiveStandbyPreemptive {
		return fmt.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	enableSpotInstance := getBool(d, "enable_spot_instance")
	spotPrice := getString(d, "spot_price")
	deleteSpot := getBool(d, "delete_spot")
	if enableSpotInstance {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("enable_spot_instance only supports AWS and Azure related cloud types")
		}

		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && deleteSpot {
			return fmt.Errorf("delete_spot only supports Azure")
		}

		gateway.EnableSpotInstance = true
		gateway.SpotPrice = spotPrice
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			gateway.DeleteSpot = deleteSpot
		}
	}

	rxQueueSize := getString(d, "rx_queue_size")
	if rxQueueSize != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("rx_queue_size only supports AWS related cloud types")
	}

	privateModeInfo, _ := client.GetPrivateModeInfo(context.Background())
	if !enablePrivateOob && !privateModeInfo.EnablePrivateMode {
		allocateNewEip := getBool(d, "allocate_new_eip")
		if allocateNewEip {
			gateway.ReuseEip = "off"
		} else {
			gateway.ReuseEip = "on"

			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
				return fmt.Errorf("failed to create spoke gateway: 'allocate_new_eip' can only be set to 'false' when cloud_type is AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048) or AWS Top Secret (16384)")
			}
			if _, ok := d.GetOk("eip"); !ok {
				return fmt.Errorf("failed to create spoke gateway: 'eip' must be set when 'allocate_new_eip' is false")
			}
			azureEipName, azureEipNameOk := d.GetOk("azure_eip_name_resource_group")
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				if !azureEipNameOk {
					return fmt.Errorf("failed to create spoke gateway: 'azure_eip_name_resource_group' must be set when 'allocate_new_eip' is false and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				gateway.Eip = fmt.Sprintf("%s:%s", mustString(azureEipName), getString(d, "eip"))
			} else {
				if azureEipNameOk {
					return fmt.Errorf("failed to create spoke gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				gateway.Eip = getString(d, "eip")
			}
		}
	}

	if privateModeInfo.EnablePrivateMode {
		if privateModeSubnetZone, ok := d.GetOk("private_mode_subnet_zone"); ok {
			gateway.Subnet = fmt.Sprintf("%s~~%s", gateway.Subnet, mustString(privateModeSubnetZone))
		} else {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("%q must be set when creating a Spoke Gateway in AWS with Private Mode enabled on the Controller", "private_mode_subnet_zone")
			}
		}

		if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
				return fmt.Errorf("private mode is only supported in AWS and Azure. %q must be empty", "private_mode_lb_vpc_id")
			}

			gateway.LbVpcId = getString(d, "private_mode_lb_vpc_id")
		}
	} else {
		if _, ok := d.GetOk("private_mode_subnet_zone"); ok {
			return fmt.Errorf("%q is only valid when Private Mode is enabled on the Controller", "private_mode_subnet_zone")
		}
		if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
			return fmt.Errorf("%q is only valid when Private Mode is enabled", "private_mode_lb_vpc_id")
		}
	}

	if gateway.EnableGlobalVpc && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("'enable_global_vpc' is only valid for GCP")
	}

	insertionGateway := getBool(d, "insertion_gateway")
	insertionGatewayAz := getString(d, "insertion_gateway_az")

	// Validation: insertion_gateway and insane_mode cannot both be true
	if insertionGateway && insaneMode {
		return fmt.Errorf("insertion_gateway and insane_mode cannot both be enabled")
	}

	// Validation: insertion_gateway is only supported on AWS
	if insertionGateway && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("insertion_gateway is only supported for AWS")
	}

	if insertionGateway {
		if insertionGatewayAz == "" {
			return fmt.Errorf("insertion_gateway_az needed if insertion_gateway is enabled")
		}
		// Append availability zone to subnet
		var strs []string
		strs = append(strs, gateway.Subnet, insertionGatewayAz)
		gateway.Subnet = strings.Join(strs, subnetSeparator)
		gateway.InsertionGateway = true
	}

	if getBool(d, "enable_ipv6") {
		if err := IPv6SupportedOnCloudType(gateway.CloudType); err != nil {
			return fmt.Errorf("error creating gateway: enable_ipv6 is not supported, %w", err)
		}
		gateway.EnableIPv6 = true

		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			subnetIPv6Cidr := getString(d, "subnet_ipv6_cidr")
			if subnetIPv6Cidr == "" {
				return fmt.Errorf("error creating gateway: subnet_ipv6_cidr must be set when enable_ipv6 is true and is enabled on %d", gateway.CloudType)
			}
			gatewaySubnet := gateway.Subnet
			// Trim any trailing '~' to normalize it first
			gatewaySubnet = strings.TrimRight(gatewaySubnet, "~")

			// Append IPv6 subnet CIDR
			gateway.Subnet = gatewaySubnet + subnetSeparator + subnetIPv6Cidr
		}
	}

	log.Printf("[INFO] Creating Aviatrix Spoke Gateway: %#v", gateway)

	d.SetId(gateway.GwName)
	flag := false
	defer func() { _ = resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.LaunchSpokeVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Spoke Gateway: %w", err)
	}

	if !singleAZ {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   getString(d, "gw_name"),
			SingleAZ: "no",
		}

		log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)

		err := client.DisableSingleAZGateway(singleAZGateway)
		if err != nil {
			return fmt.Errorf("failed to disable single AZ GW HA: %w", err)
		}
	}

	commSendCurr, commAcceptCurr, err := client.GetGatewayBgpCommunities(gateway.GwName)
	if err != nil {
		return fmt.Errorf("failed to get BGP communities for gateway %s: %w", gateway.GwName, err)
	}

	acceptComm := getBool(d, "bgp_accept_communities")
	sendComm := getBool(d, "bgp_send_communities")

	if acceptComm != commAcceptCurr {
		if err := client.SetGatewayBgpCommunitiesAccept(gateway.GwName, acceptComm); err != nil {
			return fmt.Errorf("failed to set accept BGP communities for gateway %s: %w", gateway.GwName, err)
		}
	}

	if sendComm != commSendCurr {
		if err := client.SetGatewayBgpCommunitiesSend(gateway.GwName, sendComm); err != nil {
			return fmt.Errorf("failed to set send BGP communities for gateway %s: %w", gateway.GwName, err)
		}
	}

	if haSubnet != "" || haZone != "" {
		spokeHaGw := &goaviatrix.SpokeHaGateway{
			PrimaryGwName: getString(d, "gw_name"),
			GwName:        getString(d, "gw_name") + "-hagw",
			Subnet:        haSubnet,
			Zone:          haZone,
			Eip:           getString(d, "ha_eip"),
			InsaneMode:    "no",
		}

		if insaneMode {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				var haStrs []string
				haStrs = append(haStrs, haSubnet, haInsaneModeAz)
				haSubnet = strings.Join(haStrs, subnetSeparator)
				spokeHaGw.Subnet = haSubnet
			}
			spokeHaGw.InsaneMode = "yes"
		}

		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haZone != "" {
			spokeHaGw.Subnet = fmt.Sprintf("%s~~%s~~", haSubnet, haZone)
		}

		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			spokeHaGw.Subnet = haSubnet
			spokeHaGw.AvailabilityDomain = haAvailabilityDomain
			spokeHaGw.FaultDomain = haFaultDomain
		}

		if privateModeInfo.EnablePrivateMode {
			haPrivateModeSubnetZone := getString(d, "ha_private_mode_subnet_zone")
			if haPrivateModeSubnetZone == "" && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("%q must be set when creating a Spoke HA Gateway in AWS with Private Mode enabled on the Controller", "ha_private_mode_subnet_zone")
			}
			spokeHaGw.Subnet = haSubnet + subnetSeparator + haPrivateModeSubnetZone
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if spokeHaGw.Eip != "" {
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				spokeHaGw.Eip = fmt.Sprintf("%s:%s", mustString(haAzureEipName), spokeHaGw.Eip)
			} else if haAzureEipNameOk {
				return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be empty when 'ha_eip' is empty")
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if insertionGateway {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				var haStrs []string
				haStrs = append(haStrs, haSubnet, insertionGatewayAz)
				haSubnet = strings.Join(haStrs, subnetSeparator)
				spokeHaGw.Subnet = haSubnet
			}
			spokeHaGw.InsertionGateway = true
		}

		if getBool(d, "enable_ipv6") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			haSubnetIPv6Cidr := getString(d, "ha_subnet_ipv6_cidr")
			if haSubnetIPv6Cidr == "" {
				return fmt.Errorf("error creating HA gateway: ha_subnet_ipv6_cidr must be set when enable_ipv6 is true")
			}

			haSubnet := spokeHaGw.Subnet
			haSubnetTrimmed := strings.TrimRight(haSubnet, "~")
			spokeHaGw.Subnet = haSubnetTrimmed + subnetSeparator + haSubnetIPv6Cidr
		}

		_, err := client.CreateSpokeHaGw(spokeHaGw)
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %w", err)
		}

		log.Printf("[INFO]Resizing Spoke HA Gateway: %#v", haGwSize)

		if haGwSize != gateway.VpcSize {
			if haGwSize == "" {
				return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set")
			}

			haGateway := &goaviatrix.Gateway{
				CloudType: getInt(d, "cloud_type"),
				GwName:    getString(d, "gw_name") + "-hagw",
				VpcSize:   getString(d, "ha_gw_size"),
			}

			log.Printf("[INFO] Resizing Spoke HA Gateway size to: %s ", haGateway.VpcSize)

			err := client.UpdateGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %w", err)
			}
			mustSet(d, "ha_gw_size", haGwSize)
		}
	}

	enableVpcDnsServer := getBool(d, "enable_vpc_dns_server")
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && enableVpcDnsServer {
		gwVpcDnsServer := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}

		log.Printf("[INFO] Enable VPC DNS Server: %#v", gwVpcDnsServer)

		err := client.EnableVpcDNSServer(gwVpcDnsServer)
		if err != nil {
			return fmt.Errorf("failed to enable VPC DNS Server: %w", err)
		}
	} else if enableVpcDnsServer {
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
	}

	if customizedSpokeVpcRoutes := getString(d, "customized_spoke_vpc_routes"); customizedSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                   getString(d, "gw_name"),
			CustomizedSpokeVpcRoutes: strings.Split(customizedSpokeVpcRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing customized routes of spoke gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayCustomRoutes(transitGateway)
			if err == nil {
				break
			}
			if i <= 18 && (strings.Contains(err.Error(), "when it is down") || strings.Contains(err.Error(), "hagw is down") ||
				strings.Contains(err.Error(), "gateway is down")) {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to customize spoke vpc routes of spoke gateway: %s due to: %w", transitGateway.GwName, err)
			}
		}
	}

	if filteredSpokeVpcRoutes := getString(d, "filtered_spoke_vpc_routes"); filteredSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                 getString(d, "gw_name"),
			FilteredSpokeVpcRoutes: strings.Split(filteredSpokeVpcRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing filtered routes of spoke gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayFilterRoutes(transitGateway)
			if err == nil {
				break
			}
			if i <= 18 && (strings.Contains(err.Error(), "when it is down") || strings.Contains(err.Error(), "hagw is down") ||
				strings.Contains(err.Error(), "gateway is down")) {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of spoke gateway: %s due to: %w", transitGateway.GwName, err)
			}
		}
	}

	if includedAdvertisedSpokeRoutes := getString(d, "included_advertised_spoke_routes"); includedAdvertisedSpokeRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                getString(d, "gw_name"),
			AdvertisedSpokeRoutes: strings.Split(includedAdvertisedSpokeRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing customized routes advertisement of spoke gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			if err == nil {
				break
			}
			if i <= 30 && (strings.Contains(err.Error(), "when it is down") || strings.Contains(err.Error(), "hagw is down") ||
				strings.Contains(err.Error(), "gateway is down")) {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to edit advertised spoke vpc routes of spoke gateway: %s due to: %w", transitGateway.GwName, err)
			}
		}
	}

	if enableMonitorSubnets {
		err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
		if err != nil {
			return fmt.Errorf("could not enable monitor gateway subnets: %w", err)
		}
	}

	if !getBool(d, "enable_jumbo_frame") {
		gw := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}

		err := client.DisableJumboFrame(gw)
		if err != nil {
			return fmt.Errorf("could not disable jumbo frame for spoke gateway: %w", err)
		}
	}

	if !getBool(d, "enable_gro_gso") {
		gw := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}
		err := client.DisableGroGso(gw)
		if err != nil {
			return fmt.Errorf("couldn't disable GRO/GSO on spoke gateway: %w", err)
		}
	}

	if getBool(d, "enable_private_vpc_default_route") {
		gw := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}
		err := client.EnablePrivateVpcDefaultRoute(gw)
		if err != nil {
			return fmt.Errorf("could not enable private vpc default route after spoke gateway creation: %w", err)
		}
	}

	if getBool(d, "enable_skip_public_route_table_update") {
		gw := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}
		err := client.EnableSkipPublicRouteUpdate(gw)
		if err != nil {
			return fmt.Errorf("could not enable skip public route update after spoke gateway creation: %w", err)
		}
	}

	if getBool(d, "enable_auto_advertise_s2c_cidrs") {
		gw := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}
		err := client.EnableAutoAdvertiseS2CCidrs(gw)
		if err != nil {
			return fmt.Errorf("could not enable auto advertise s2c cidrs after spoke gateaway creation: %w", err)
		}
	}

	if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
		err := client.ModifyTunnelDetectionTime(getString(d, "gw_name"), mustInt(detectionTime))
		if err != nil {
			return fmt.Errorf("could not set tunnel detection time during Spoke Gateway creation: %w", err)
		}
	}

	if learnedCidrsApproval {
		gateway.LearnedCidrsApproval = "on"
		err := client.EnableSpokeLearnedCidrsApproval(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable learned cidrs approval: %w", err)
		}
	}
	if len(gateway.ApprovedLearnedCidrs) != 0 {
		err := client.UpdateSpokePendingApprovedCidrs(gateway)
		if err != nil {
			return fmt.Errorf("failed to update approved CIDRs: %w", err)
		}
	}

	if val, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		var spokeBgpManualSpokeAdvertiseCidrs []string
		slice := mustSlice(val)
		for _, v := range slice {
			spokeBgpManualSpokeAdvertiseCidrs = append(spokeBgpManualSpokeAdvertiseCidrs, mustString(v))
		}
		gateway.BgpManualSpokeAdvertiseCidrs = strings.Join(spokeBgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetSpokeBgpManualAdvertisedNetworks(gateway)
		if err != nil {
			return fmt.Errorf("failed to set spoke BGP Manual Advertise Cidrs: %w", err)
		}
	}

	if val, ok := d.GetOk("bgp_ecmp"); ok {
		err := client.SetBgpEcmpSpoke(gateway, mustBool(val))
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp: %w", err)
		}
	}

	if enableActiveStandby {
		if enableActiveStandbyPreemptive {
			if err := client.EnableActiveStandbyPreemptiveSpoke(gateway); err != nil {
				return fmt.Errorf("could not enable Preemptive Mode for Active-Standby: %w", err)
			}
		} else {
			if err := client.EnableActiveStandbySpoke(gateway); err != nil {
				return fmt.Errorf("could not enable Active-Standby: %w", err)
			}
		}
	}

	if disableRoutePropagation {
		if err := client.DisableSpokeOnpremRoutePropagation(gateway); err != nil {
			return fmt.Errorf("could not disable route propagation for Spoke %s : %w", gateway.GwName, err)
		}
	}

	if val, ok := d.GetOk("local_as_number"); ok {
		err := client.SetLocalASNumberSpoke(gateway, mustString(val))
		if err != nil {
			return fmt.Errorf("could not set local_as_number: %w", err)
		}
	}

	if val, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		slice := mustSlice(val)
		for _, v := range slice {
			prependASPath = append(prependASPath, mustString(v))
		}
		err := client.SetPrependASPathSpoke(gateway, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %w", err)
		}
	}

	if val, ok := d.GetOk("bgp_polling_time"); ok {
		bgp_polling_time := mustInt(val)
		if bgp_polling_time >= 10 && bgp_polling_time != defaultBgpPollingTime {
			err := client.SetBgpPollingTimeSpoke(gateway, bgp_polling_time)
			if err != nil {
				return fmt.Errorf("could not set bgp polling time: %w", err)
			}
		}
	}

	if val, ok := d.GetOk("bgp_neighbor_status_polling_time"); ok {
		bgp_neighbor_status_polling_time := mustInt(val)
		if bgp_neighbor_status_polling_time >= 1 && bgp_neighbor_status_polling_time != defaultBgpNeighborStatusPollingTime {
			err := client.SetBgpBfdPollingTimeSpoke(gateway, mustInt(val))
			if err != nil {
				return fmt.Errorf("could not set bgp neighbor status polling time: %w", err)
			}
		}
	}

	if holdTime := getInt(d, "bgp_hold_time"); holdTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gateway.GwName, holdTime)
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time after Spoke Gateway creation: %w", err)
		}
	}

	enableSpokePreserveAsPath := getBool(d, "enable_preserve_as_path")
	if enableSpokePreserveAsPath {
		if enableBgp {
			err := client.EnableSpokePreserveAsPath(gateway)
			if err != nil {
				return fmt.Errorf("could not enable spoke preserve as path: %w", err)
			}
		} else {
			return fmt.Errorf("enable_preserve_as_path is not supported for Non-BGP Spoke Gateways")
		}
	}

	if rxQueueSize != "" {
		gwRxQueueSize := &goaviatrix.Gateway{
			GwName:      getString(d, "gw_name"),
			RxQueueSize: rxQueueSize,
		}
		err := client.SetRxQueueSize(gwRxQueueSize)
		if err != nil {
			return fmt.Errorf("failed to set rx queue size for spoke %s: %w", gateway.GwName, err)
		}
		if haSubnet != "" || haZone != "" {
			haGwRxQueueSize := &goaviatrix.Gateway{
				GwName:      getString(d, "gw_name") + "-hagw",
				RxQueueSize: rxQueueSize,
			}
			err := client.SetRxQueueSize(haGwRxQueueSize)
			if err != nil {
				return fmt.Errorf("failed to set rx queue size for spoke ha %s : %w", haGwRxQueueSize.GwName, err)
			}
		}
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		routeTables := getStringSet(d, "private_route_table_config")
		fmt.Println("######## routeTables", routeTables)
		if len(routeTables) > 0 {
			gw := &goaviatrix.Gateway{GwName: getString(d, "gw_name")}
			err := client.EditPrivateRouteTableConfig(gw, routeTables)
			if err != nil {
				return fmt.Errorf("could not edit private route table config: %w", err)
			}
		}
	}

	return resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSpokeGatewayReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeGatewayRead(d, meta)
	}
	return nil
}

func resourceAviatrixSpokeGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	ignoreTagsConfig := client.IgnoreTagsConfig

	var isImport bool
	gwName := getString(d, "gw_name")
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		mustSet(d, "gw_name", id)
		mustSet(d, "manage_ha_gateway", true)
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: getString(d, "account_name"),
		GwName:      getString(d, "gw_name"),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Spoke Gateway: %w", err)
	}

	log.Printf("[TRACE] reading spoke gateway %s: %#v", getString(d, "gw_name"), gw)
	mustSet(d, "cloud_type", gw.CloudType)
	mustSet(d, "account_name", gw.AccountName)
	mustSet(d, "enable_encrypt_volume", gw.EnableEncryptVolume)
	mustSet(d, "enable_private_vpc_default_route", gw.PrivateVpcDefaultEnabled)
	mustSet(d, "enable_skip_public_route_table_update", gw.SkipPublicVpcUpdateEnabled)
	mustSet(d, "private_route_table_config", gw.PrivateRouteTableConfig)
	mustSet(d, "enable_auto_advertise_s2c_cidrs", gw.AutoAdvertiseCidrsEnabled)
	mustSet(d, "eip", gw.PublicIP)
	mustSet(d, "subnet", gw.VpcNet)
	mustSet(d, "gw_size", gw.GwSize)
	mustSet(d, "cloud_instance_id", gw.CloudnGatewayInstID)
	mustSet(d, "security_group_id", gw.GwSecurityGroupID)
	mustSet(d, "private_ip", gw.PrivateIP)
	mustSet(d, "single_az_ha", gw.SingleAZ == "yes")
	mustSet(d, "enable_vpc_dns_server", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled")
	mustSet(d, "single_ip_snat", gw.EnableNat == "yes" && gw.SnatMode == "primary")
	mustSet(d, "enable_jumbo_frame", gw.JumboFrame)
	mustSet(d, "enable_bgp", gw.EnableBgp)
	mustSet(d, "enable_bgp_over_lan", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan)
	mustSet(d, "enable_ipv6", gw.EnableIPv6)
	mustSet(d, "insertion_gateway", gw.InsertionGateway)
	mustSet(d, "subnet_ipv6_cidr", gw.SubnetIPv6Cidr)

	if gw.InsertionGateway && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		mustSet(d, "insertion_gateway_az", gw.GatewayZone)
	} else {
		mustSet(d, "insertion_gateway_az", "")
	}
	mustSet(d, "tunnel_encryption_cipher", gw.TunnelEncryptionCipher)
	mustSet(d, "tunnel_forward_secrecy", gw.TunnelForwardSecrecy)

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan {
		bgpLanIpInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return fmt.Errorf("could not get BGP LAN IP info for Azure spoke gateway %s: %w", gateway.GwName, err)
		}
		if err = d.Set("bgp_lan_ip_list", bgpLanIpInfo.AzureBgpLanIpList); err != nil {
			log.Printf("[WARN] could not set bgp_lan_ip_list into state: %s", err)
		}
		if len(bgpLanIpInfo.AzureHaBgpLanIpList) != 0 {
			if err = d.Set("ha_bgp_lan_ip_list", bgpLanIpInfo.AzureHaBgpLanIpList); err != nil {
				log.Printf("[WARN] could not set ha_bgp_lan_ip_list into state: %s", err)
			}
		} else {
			mustSet(d, "ha_bgp_lan_ip_list", nil)
		}
	} else {
		mustSet(d, "bgp_lan_ip_list", nil)
		mustSet(d, "ha_bgp_lan_ip_list", nil)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan {
		mustSet(d, "bgp_lan_interfaces_count", gw.BgpLanInterfacesCount)
	} else {
		mustSet(d, "bgp_lan_interfaces_count", nil)
	}
	mustSet(d, "enable_learned_cidrs_approval", gw.EnableLearnedCidrsApproval)
	mustSet(d, "enable_preserve_as_path", gw.EnablePreserveAsPath)
	mustSet(d, "rx_queue_size", gw.RxQueueSize)
	mustSet(d, "public_ip", gw.PublicIP)
	mustSet(d, "enable_global_vpc", gw.EnableGlobalVpc)

	if gw.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: gw.GwName})
		if err != nil {
			return fmt.Errorf("could not get advanced config for spoke gateway: %w", err)
		}

		if err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs); err != nil {
			return fmt.Errorf("could not set approved_learned_cidrs into state: %w", err)
		}
	} else {
		mustSet(d, "approved_learned_cidrs", nil)
	}
	mustSet(d, "local_as_number", gw.LocalASNumber)
	mustSet(d, "bgp_ecmp", gw.BgpEcmp)
	mustSet(d, "enable_active_standby", gw.EnableActiveStandby)
	mustSet(d, "enable_active_standby_preemptive", gw.EnableActiveStandbyPreemptive)
	mustSet(d, "disable_route_propagation", gw.DisableRoutePropagation)
	var prependAsPath []string
	for _, p := range strings.Split(gw.PrependASPath, " ") {
		if p != "" {
			prependAsPath = append(prependAsPath, p)
		}
	}
	err = d.Set("prepend_as_path", prependAsPath)
	if err != nil {
		return fmt.Errorf("could not set prepend_as_path: %w", err)
	}
	if gw.EnableBgp {
		mustSet(d, "learned_cidrs_approval_mode", gw.LearnedCidrsApprovalMode)
		mustSet(d, "bgp_polling_time", gw.BgpPollingTime)
		mustSet(d, "bgp_neighbor_status_polling_time", gw.BgpBfdPollingTime)
		mustSet(d, "bgp_hold_time", gw.BgpHoldTime)
	} else {
		mustSet(d, "learned_cidrs_approval_mode", "gateway")
		mustSet(d, "bgp_polling_time", 50)
		mustSet(d, "bgp_neighbor_status_polling_time", defaultBgpNeighborStatusPollingTime)
		mustSet(d, "bgp_hold_time", 180)
	}
	mustSet(d, "tunnel_detection_time", gw.TunnelDetectionTime)
	mustSet(d, "image_version", gw.ImageVersion)
	mustSet(d, "software_version", gw.SoftwareVersion)

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			mustSet(d, "azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Spoke Gateway %s", gw.GwName)
		}
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, subnetSeparator)[0])
		mustSet( // AWS vpc_id returns as <vpc_id>~~<other vpc info> in rest api
			d, "vpc_reg", gw.VpcRegion) // AWS vpc_reg returns as vpc_region in rest api

		if gw.AllocateNewEipRead && !gw.EnablePrivateOob {
			mustSet(d, "allocate_new_eip", true)
		} else {
			mustSet(d, "allocate_new_eip", false)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(
			// gcp vpc_id returns as <vpc name>~-~<project name>
			d, "vpc_id", gw.VpcID)
		mustSet(d, "vpc_reg", gw.GatewayZone)
		mustSet( // gcp vpc_reg returns as gateway_zone in json

			d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		mustSet(d, "vpc_reg", gw.VpcRegion)
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, subnetSeparator)[0])
		mustSet( // oci vpc_id returns as <vpc_id>~~<vpc_name> in rest api
			d, "vpc_reg", gw.VpcRegion)
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if gw.CloudType == goaviatrix.AliCloud {
		mustSet(d, "vpc_id", strings.Split(gw.VpcID, subnetSeparator)[0])
		mustSet(d, "vpc_reg", gw.VpcRegion)
		mustSet(d, "allocate_new_eip", true)
	}

	if gw.InsaneMode == "yes" {
		mustSet(d, "insane_mode", true)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "insane_mode_az", gw.GatewayZone)
		} else {
			mustSet(d, "insane_mode_az", "")
		}
	} else {
		mustSet(d, "insane_mode", false)
		mustSet(d, "insane_mode_az", "")
	}

	if len(gw.CustomizedSpokeVpcRoutes) != 0 {
		if customizedSpokeVpcRoutes := getString(d, "customized_spoke_vpc_routes"); customizedSpokeVpcRoutes != "" {
			customizedRoutesArray := strings.Split(customizedSpokeVpcRoutes, ",")
			if len(goaviatrix.Difference(customizedRoutesArray, gw.CustomizedSpokeVpcRoutes)) == 0 &&
				len(goaviatrix.Difference(gw.CustomizedSpokeVpcRoutes, customizedRoutesArray)) == 0 {
				mustSet(d, "customized_spoke_vpc_routes", customizedSpokeVpcRoutes)
			} else {
				mustSet(d, "customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
			}
		} else {
			mustSet(d, "customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
		}
	} else {
		mustSet(d, "customized_spoke_vpc_routes", "")
	}

	if len(gw.FilteredSpokeVpcRoutes) != 0 {
		if filteredSpokeVpcRoutes := getString(d, "filtered_spoke_vpc_routes"); filteredSpokeVpcRoutes != "" {
			filteredSpokeVpcRoutesArray := strings.Split(filteredSpokeVpcRoutes, ",")
			if len(goaviatrix.Difference(filteredSpokeVpcRoutesArray, gw.FilteredSpokeVpcRoutes)) == 0 &&
				len(goaviatrix.Difference(gw.FilteredSpokeVpcRoutes, filteredSpokeVpcRoutesArray)) == 0 {
				mustSet(d, "filtered_spoke_vpc_routes", filteredSpokeVpcRoutes)
			} else {
				mustSet(d, "filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
			}
		} else {
			mustSet(d, "filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
		}
	} else {
		mustSet(d, "filtered_spoke_vpc_routes", "")
	}

	if len(gw.IncludeCidrList) != 0 {
		if includedAdvertisedSpokeRoutes := getString(d, "included_advertised_spoke_routes"); includedAdvertisedSpokeRoutes != "" {
			advertisedSpokeRoutesArray := strings.Split(includedAdvertisedSpokeRoutes, ",")
			if len(goaviatrix.Difference(advertisedSpokeRoutesArray, gw.IncludeCidrList)) == 0 &&
				len(goaviatrix.Difference(gw.IncludeCidrList, advertisedSpokeRoutesArray)) == 0 {
				mustSet(d, "included_advertised_spoke_routes", includedAdvertisedSpokeRoutes)
			} else {
				mustSet(d, "included_advertised_spoke_routes", strings.Join(gw.IncludeCidrList, ","))
			}
		} else {
			mustSet(d, "included_advertised_spoke_routes", strings.Join(gw.AdvertisedSpokeRoutes, ","))
		}
	} else {
		mustSet(d, "included_advertised_spoke_routes", "")
	}
	mustSet(d, "enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
	if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
		return fmt.Errorf("setting 'monitor_exclude_list' to state: %w", err)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		tags := goaviatrix.KeyValueTags(gw.Tags).IgnoreConfig(ignoreTagsConfig)
		if err := d.Set("tags", tags); err != nil {
			log.Printf("[WARN] Error setting tags for (%s): %s", d.Id(), err)
		}
	}

	var spokeBgpManualAdvertiseCidrs []string
	if val, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		slice := mustSlice(val)
		for _, v := range slice {
			spokeBgpManualAdvertiseCidrs = append(spokeBgpManualAdvertiseCidrs, mustString(v))
		}
	}
	if len(goaviatrix.Difference(spokeBgpManualAdvertiseCidrs, gw.BgpManualSpokeAdvertiseCidrs)) != 0 ||
		len(goaviatrix.Difference(gw.BgpManualSpokeAdvertiseCidrs, spokeBgpManualAdvertiseCidrs)) != 0 {
		mustSet(d, "spoke_bgp_manual_advertise_cidrs", gw.BgpManualSpokeAdvertiseCidrs)
	} else {
		mustSet(d, "spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertiseCidrs)
	}
	mustSet(d, "enable_private_oob", gw.EnablePrivateOob)
	if gw.EnablePrivateOob {
		mustSet(d, "oob_management_subnet", strings.Split(gw.OobManagementSubnet, subnetSeparator)[0])
		mustSet(d, "oob_availability_zone", gw.GatewayZone)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		_, zoneIsSet := d.GetOk("zone")
		if (isImport || zoneIsSet) && gw.GatewayZone != "AvailabilitySet" && gw.LbVpcId == "" {
			mustSet(d, "zone", "az-"+gw.GatewayZone)
		}
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		if gw.GatewayZone != "" {
			mustSet(d, "availability_domain", gw.GatewayZone)
		} else {
			mustSet(d, "availability_domain", getString(d, "availability_domain"))
		}
		mustSet(d, "fault_domain", gw.FaultDomain)
	}

	if gw.EnableSpotInstance {
		mustSet(d, "enable_spot_instance", true)
		mustSet(d, "spot_price", gw.SpotPrice)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.DeleteSpot {
			mustSet(d, "delete_spot", gw.DeleteSpot)
		}
	}
	mustSet(d, "private_mode_lb_vpc_id", gw.LbVpcId)
	if gw.LbVpcId != "" && gw.GatewayZone != "AvailabilitySet" {
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "private_mode_subnet_zone", gw.GatewayZone)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			mustSet(d, "private_mode_subnet_zone", "az-"+gw.GatewayZone)
		}
	} else {
		mustSet(d, "private_mode_subnet_zone", nil)
	}

	enableGroGso, err := client.GetGroGsoStatus(gw)
	if err != nil {
		return fmt.Errorf("failed to get GRO/GSO status of spoke gateway %s: %w", gw.GwName, err)
	}
	mustSet(d, "enable_gro_gso", enableGroGso)

	if getBool(d, "manage_ha_gateway") {
		if gw.HaGw.GwSize == "" {
			mustSet(d, "ha_availability_domain", "")
			mustSet(d, "ha_azure_eip_name_resource_group", "")
			mustSet(d, "ha_cloud_instance_id", "")
			mustSet(d, "ha_eip", "")
			mustSet(d, "ha_fault_domain", "")
			mustSet(d, "ha_gw_name", "")
			mustSet(d, "ha_gw_size", "")
			mustSet(d, "ha_image_version", "")
			mustSet(d, "ha_insane_mode_az", "")
			mustSet(d, "ha_oob_availability_zone", "")
			mustSet(d, "ha_oob_management_subnet", "")
			mustSet(d, "ha_private_ip", "")
			mustSet(d, "ha_security_group_id", "")
			mustSet(d, "ha_software_version", "")
			mustSet(d, "ha_subnet", "")
			mustSet(d, "ha_subnet_ipv6_cidr", "")
			mustSet(d, "ha_zone", "")
			mustSet(d, "ha_public_ip", "")
			mustSet(d, "ha_private_mode_subnet_zone", "")
			mustSet(d, "ha_bgp_lan_ip_list", nil)
			return nil
		}

		log.Printf("[INFO] Spoke HA Gateway size: %s", gw.HaGw.GwSize)
		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			mustSet(d, "ha_subnet", gw.HaGw.VpcNet)
			if zone := d.Get("ha_zone"); goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || mustString(zone) != "") {
				if gw.LbVpcId == "" && gw.HaGw.GatewayZone != "AvailabilitySet" {
					mustSet(d, "ha_zone", "az-"+gw.HaGw.GatewayZone)
				} else {
					mustSet(d, "ha_zone", "")
				}
			} else {
				mustSet(d, "ha_zone", "")
			}
		} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			mustSet(d, "ha_zone", gw.HaGw.GatewayZone)
			if d.Get("ha_subnet") != "" || isImport {
				mustSet(d, "ha_subnet", gw.HaGw.VpcNet)
			} else {
				mustSet(d, "ha_subnet", "")
			}
		}
		mustSet(d, "ha_subnet_ipv6_cidr", gw.HaGw.SubnetIPv6Cidr)

		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			if gw.HaGw.GatewayZone != "" {
				mustSet(d, "ha_availability_domain", gw.HaGw.GatewayZone)
			} else {
				mustSet(d, "ha_availability_domain", getString(d, "ha_availability_domain"))
			}
			mustSet(d, "ha_fault_domain", gw.HaGw.FaultDomain)
		}
		mustSet(d, "ha_eip", gw.HaGw.PublicIP)
		mustSet(d, "ha_gw_size", gw.HaGw.GwSize)
		mustSet(d, "ha_cloud_instance_id", gw.HaGw.CloudnGatewayInstID)
		mustSet(d, "ha_gw_name", gw.HaGw.GwName)
		mustSet(d, "ha_private_ip", gw.HaGw.PrivateIP)
		mustSet(d, "ha_software_version", gw.HaGw.SoftwareVersion)
		mustSet(d, "ha_image_version", gw.HaGw.ImageVersion)
		mustSet(d, "ha_security_group_id", gw.HaGw.GwSecurityGroupID)
		mustSet(d, "ha_public_ip", gw.HaGw.PublicIP)
		if gw.HaGw.InsaneMode == "yes" && goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "ha_insane_mode_az", gw.HaGw.GatewayZone)
		} else {
			mustSet(d, "ha_insane_mode_az", "")
		}
		if gw.HaGw.EnablePrivateOob {
			mustSet(d, "ha_oob_management_subnet", strings.Split(gw.HaGw.OobManagementSubnet, subnetSeparator)[0])
			mustSet(d, "ha_oob_availability_zone", gw.HaGw.GatewayZone)
		}
		if gw.LbVpcId != "" && gw.GatewayZone != "AvailabilitySet" {
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				mustSet(d, "ha_private_mode_subnet_zone", gw.HaGw.GatewayZone)
			} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				mustSet(d, "ha_private_mode_subnet_zone", "az-"+gw.HaGw.GatewayZone)
			}
		} else {
			mustSet(d, "ha_private_mode_subnet_zone", "")
		}
		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
			if len(azureEip) == 3 {
				mustSet(d, "ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the HA Gateway %s", gw.GwName)
			}
		}
	}

	sendComm, acceptComm, err := client.GetGatewayBgpCommunities(gateway.GwName)
	if err != nil {
		return fmt.Errorf("failed to get BGP communities for gateway %s: %w", gateway.GwName, err)
	}
	err = d.Set("bgp_send_communities", sendComm)
	if err != nil {
		return fmt.Errorf("failed to set bgp_send_communities: %w", err)
	}
	err = d.Set("bgp_accept_communities", acceptComm)
	if err != nil {
		return fmt.Errorf("failed to set bgp_accept_communities: %w", err)
	}

	return nil
}

func resourceAviatrixSpokeGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}

	manageHaGw := getBool(d, "manage_ha_gateway")
	if d.HasChange("manage_ha_gateway") {
		_, nMHG := d.GetChange("manage_ha_gateway")
		newManageHaGw := mustBool(nMHG)
		if newManageHaGw {
			mustSet(d, "manage_ha_gateway", true)
		} else {
			mustSet(d, "manage_ha_gateway", false)
		}
	}

	if !manageHaGw && !d.HasChange("manage_ha_gateway") {
		if d.HasChanges("ha_subnet", "ha_zone", "ha_gw_size", "ha_insane_mode_az", "ha_eip",
			"ha_azure_eip_name_resource_group", "ha_availability_domain", "ha_fault_domain", "ha_oob_management_subnet",
			"ha_private_mode_subnet_zone", "ha_oob_availability_zone", "ha_software_version", "ha_image_version") {
			return fmt.Errorf("'manage_ha_gateway' is set to false. Please set it to true, or use 'aviatrix_spoke_ha_gateway' to manage editing spoke ha gateway")
		}
	}

	haGateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name") + "-hagw",
		VpcSize:   getString(d, "ha_gw_size"),
	}

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	d.Partial(true)
	commSendCurr, commAcceptCurr, err := client.GetGatewayBgpCommunities(gateway.GwName)
	if err != nil {
		return fmt.Errorf("failed to get BGP communities for gateway %s: %w", gateway.GwName, err)
	}

	if d.HasChange("bgp_accept_communities") {
		acceptComm := getBool(d, "bgp_accept_communities")

		if acceptComm != commAcceptCurr {
			if err := client.SetGatewayBgpCommunitiesAccept(gateway.GwName, acceptComm); err != nil {
				return fmt.Errorf("failed to set accept BGP communities for gateway %s: %w", gateway.GwName, err)
			}
		}
	}
	if d.HasChange("bgp_send_communities") {
		sendComm := getBool(d, "bgp_send_communities")

		if sendComm != commSendCurr {
			if err := client.SetGatewayBgpCommunitiesSend(gateway.GwName, sendComm); err != nil {
				return fmt.Errorf("failed to set send BGP communities for gateway %s: %w", gateway.GwName, err)
			}
		}
	}

	if d.HasChange("private_route_table_config") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		routeTables := getStringSet(d, "private_route_table_config")
		err := client.EditPrivateRouteTableConfig(gateway, routeTables)
		if err != nil {
			return fmt.Errorf("could not edit private route table config: %w", err)
		}
	}

	if getBool(d, "enable_private_vpc_default_route") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_private_vpc_default_route is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}
	if getBool(d, "enable_skip_public_route_table_update") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_skip_public_route_update is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if d.HasChange("ha_zone") {
		haZone := getString(d, "ha_zone")
		if haZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("'ha_zone' is only valid for GCP (4), Azure (8), AzureGov (32) and AzureChina (2048) providers if enabling HA")
		}
	}
	if d.HasChange("ha_zone") || d.HasChange("ha_subnet") {
		haZone := getString(d, "ha_zone")
		haSubnet := getString(d, "ha_subnet")
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && haSubnet != "" && haZone == "" {
			return fmt.Errorf("'ha_zone' must be set to enable HA on GCP (4), cannot enable HA with only 'ha_subnet'")
		}
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" && haZone != "" {
			return fmt.Errorf("'ha_subnet' must be provided to enable HA for Azure (8), AzureGov (32) or AzureChina (2048), cannot enable HA with only 'ha_zone'")
		}
	}
	if d.HasChange("allocate_new_eip") {
		return fmt.Errorf("updating allocate_new_eip is not allowed")
	}
	if d.HasChange("eip") {
		return fmt.Errorf("updating eip is not allowed")
	}
	if d.HasChange("ha_eip") {
		o, n := d.GetChange("ha_eip")
		if mustString(o) != "" && mustString(n) != "" {
			return fmt.Errorf("updating ha_eip is not allowed")
		}
	}
	if d.HasChange("azure_eip_name_resource_group") {
		return fmt.Errorf("failed to update spoke gateway: changing 'azure_eip_name_resource_group' is not allowed")
	}
	if d.HasChange("ha_azure_eip_name_resource_group") {
		o, n := d.GetChange("ha_azure_eip_name_resource_group")
		if mustString(o) != "" && mustString(n) != "" {
			return fmt.Errorf("failed to update spoke gateway: changing 'ha_azure_eip_name_resource_group' is not allowed")
		}
	}

	learnedCidrsApproval := getBool(d, "enable_learned_cidrs_approval")
	approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")
	if !learnedCidrsApproval && len(approvedLearnedCidrs) != 0 {
		return fmt.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if d.HasChange("enable_private_oob") {
		return fmt.Errorf("updating enable_private_oob is not allowed")
	}
	enablePrivateOob := getBool(d, "enable_private_oob")
	privateModeInfo, _ := client.GetPrivateModeInfo(context.Background())
	if !enablePrivateOob {
		if d.HasChange("ha_oob_management_subnet") {
			return fmt.Errorf("updating ha_oob_management_subnet is not allowed if private oob is disabled")
		}

		if d.HasChange("ha_oob_availability_zone") {
			return fmt.Errorf("updating ha_oob_availability_zone is not allowed if private oob is disabled")
		}
	}
	if !privateModeInfo.EnablePrivateMode {
		if d.HasChange("ha_private_mode_subnet_zone") {
			return fmt.Errorf("updating %q is not allowed if private mode is disabled", "ha_private_mode_subnet_zone")
		}
	}

	if d.HasChange("enable_global_vpc") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("global vpc can only be enabled for GCP")
	}

	if d.HasChange("enable_preserve_as_path") {
		enableBgp := getBool(d, "enable_bgp")
		enableSpokePreserveAsPath := getBool(d, "enable_preserve_as_path")
		if enableSpokePreserveAsPath && !enableBgp {
			return fmt.Errorf("enable_preserve_as_path is not supported for Non-BGP Spoke during Spoke Gateway update")
		}
		if !enableSpokePreserveAsPath {
			err := client.DisableSpokePreserveAsPath(&goaviatrix.SpokeVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not disable Preserve AS Path during Spoke Gateway update: %w", err)
			}
		} else {
			err := client.EnableSpokePreserveAsPath(&goaviatrix.SpokeVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not enable Preserve AS Path during Spoke Gateway update: %w", err)
			}
		}
	}

	if d.HasChange("tags") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("error updating spoke gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: getString(d, "gw_name"),
			CloudType:    gateway.CloudType,
		}

		tagsMap, err := extractTags(d, gateway.CloudType)
		if err != nil {
			return fmt.Errorf("failed to update tags for spoke gateway: %w", err)
		}
		tags.Tags = tagsMap
		tagJson, err := TagsMapToJson(tagsMap)
		if err != nil {
			return fmt.Errorf("failed to update tags for spoke gateway: %w", err)
		}
		tags.TagJson = tagJson
		err = client.UpdateTags(tags)
		if err != nil {
			return fmt.Errorf("failed to update tags for spoke gateway: %w", err)
		}
	}

	if d.HasChange("gw_size") {
		gateway.VpcSize = getString(d, "gw_size")
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Spoke Gateway: %w", err)
		}
	}

	newHaGwEnabled := false
	if manageHaGw && (d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") || d.HasChange("ha_subnet_ipv6_cidr") ||
		(enablePrivateOob && (d.HasChange("ha_oob_management_subnet") || d.HasChange("ha_oob_availability_zone"))) ||
		(privateModeInfo.EnablePrivateMode && d.HasChange("ha_private_mode_subnet_zone")) ||
		d.HasChange("ha_availability_domain") || d.HasChange("ha_fault_domain")) {
		haGwSize := getString(d, "ha_gw_size")
		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false

		spokeHaGw := &goaviatrix.SpokeHaGateway{
			PrimaryGwName: getString(d, "gw_name"),
			GwName:        getString(d, "gw_name") + "-hagw",
			GwSize:        haGwSize,
			InsaneMode:    "no",
		}

		haEip := getString(d, "ha_eip")
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			spokeHaGw.Eip = haEip
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if haEip != "" && newSubnet != "" {
				// No change will be detected when ha_eip is set to the empty string because it is computed.
				// Instead, check ha_gw_size to detect when HA gateway is being deleted.
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				spokeHaGw.Eip = fmt.Sprintf("%s:%s", mustString(haAzureEipName), haEip)
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create HA Spoke Gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}

		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			spokeHaGw.Subnet = getString(d, "ha_subnet")
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && getString(d, "ha_zone") != "" {
				spokeHaGw.Subnet = fmt.Sprintf("%s~~%s~~", getString(d, "ha_subnet"), getString(d, "ha_zone"))
			}

			haAvailabilityDomain := getString(d, "ha_availability_domain")
			haFaultDomain := getString(d, "ha_fault_domain")
			if newSubnet != "" {
				if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain == "" || haFaultDomain == "") {
					return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are required to enable HA on OCI")
				}
				if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain != "" || haFaultDomain != "") {
					return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are only valid for OCI")
				}
			}
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) {
				spokeHaGw.Subnet = getString(d, "ha_subnet")
				spokeHaGw.AvailabilityDomain = haAvailabilityDomain
				spokeHaGw.FaultDomain = haFaultDomain
			}

			if oldSubnet == "" && newSubnet != "" {
				newHaGwEnabled = true
			} else if oldSubnet != "" && newSubnet == "" {
				deleteHaGw = true
			} else if oldSubnet != "" && newSubnet != "" {
				changeHaGw = true
			} else if enablePrivateOob && d.HasChanges("ha_oob_management_subnet", "ha_oob_availability_zone") ||
				privateModeInfo.EnablePrivateMode && d.HasChange("ha_private_mode_subnet_zone") ||
				d.HasChanges("ha_zone", "ha_availability_domain", "ha_fault_domain") {
				changeHaGw = true
			}
		} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			spokeHaGw.Zone = getString(d, "ha_zone")
			spokeHaGw.Subnet = getString(d, "ha_subnet")
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}

		if getBool(d, "insane_mode") {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				var haStrs []string
				insaneModeHaAz := getString(d, "ha_insane_mode_az")
				haSubnet := getString(d, "ha_subnet")

				if insaneModeHaAz == "" && haSubnet != "" {
					return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set " +
						"for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
				} else if insaneModeHaAz != "" && haSubnet == "" {
					return fmt.Errorf("ha_subnet needed if insane_mode is enabled and ha_insane_mode_az is set " +
						"for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
				}

				haStrs = append(haStrs, spokeHaGw.Subnet, insaneModeHaAz)
				spokeHaGw.Subnet = strings.Join(haStrs, subnetSeparator)
			}
			spokeHaGw.InsaneMode = "yes"
		}

		if (newHaGwEnabled || changeHaGw) && haGwSize == "" {
			return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set")
		} else if deleteHaGw && haGwSize != "" {
			return fmt.Errorf("ha_gw_size must be empty if spoke HA gateway is deleted")
		}

		haOobManagementSubnet := getString(d, "ha_oob_management_subnet")
		haOobAvailabilityZone := getString(d, "ha_oob_availability_zone")

		if enablePrivateOob {
			if newHaGwEnabled || changeHaGw {
				if haOobAvailabilityZone == "" {
					return fmt.Errorf("\"ha_oob_availability_zone\" is required if \"enable_private_oob\" is true and \"ha_subnet\" is provided")
				}

				if haOobManagementSubnet == "" {
					return fmt.Errorf("\"ha_oob_management_subnet\" is required if \"enable_private_oob\" is true and \"ha_subnet\" is provided")
				}
			} else if deleteHaGw {
				if haOobAvailabilityZone != "" {
					return fmt.Errorf("\"ha_oob_availability_zone\" must be empty if \"ha_subnet\" is empty")
				}

				if haOobManagementSubnet != "" {
					return fmt.Errorf("\"ha_oob_management_subnet\" must be empty if \"ha_subnet\" is empty")
				}
			}
		}

		if privateModeInfo.EnablePrivateMode {
			if newHaGwEnabled || changeHaGw {
				if _, ok := d.GetOk("ha_private_mode_subnet_zone"); !ok && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
					return fmt.Errorf("%q is required when creating a Spoke HA Gateway in AWS if private mode is enabled and %q is provided", "ha_private_mode_subnet_zone", "ha_subnet")
				}

				privateModeSubnetZone := getString(d, "ha_private_mode_subnet_zone")
				spokeHaGw.Subnet += subnetSeparator + privateModeSubnetZone
			}
		}

		insertionGateway := getBool(d, "insertion_gateway")
		insertionGatewayAz := getString(d, "insertion_gateway_az")

		if insertionGateway {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				var haStrs []string
				haStrs = append(haStrs, spokeHaGw.Subnet, insertionGatewayAz)
				spokeHaGw.Subnet = strings.Join(haStrs, subnetSeparator)
			}
			spokeHaGw.InsertionGateway = true
		}

		if getBool(d, "enable_ipv6") && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			haSubnetIPv6Cidr := getString(d, "ha_subnet_ipv6_cidr")
			if haSubnetIPv6Cidr == "" {
				return fmt.Errorf("error creating HA gateway: ha_subnet_ipv6_cidr must be set when enable_ipv6 is true")
			}

			haSubnet := spokeHaGw.Subnet
			haSubnetTrimmed := strings.TrimRight(haSubnet, "~")
			spokeHaGw.Subnet = haSubnetTrimmed + subnetSeparator + haSubnetIPv6Cidr
		}

		if newHaGwEnabled {
			// New configuration to enable HA
			_, err := client.CreateSpokeHaGw(spokeHaGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %w", err)
			}
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				if getString(d, "rx_queue_size") != "" && !d.HasChange("rx_queue_size") {
					haGwRxQueueSize := &goaviatrix.Gateway{
						GwName:      getString(d, "gw_name") + "-hagw",
						RxQueueSize: getString(d, "rx_queue_size"),
					}
					err := client.SetRxQueueSize(haGwRxQueueSize)
					if err != nil {
						return fmt.Errorf("could not set rx queue size for spoke ha: %s during gateway update: %w", haGwRxQueueSize.GwName, err)
					}
				}
			}
			//}
		} else if deleteHaGw {
			// Ha configuration has been deleted
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %w", err)
			}
		} else if changeHaGw {
			// HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %w", err)
			}

			spokeHaGw.Eip = ""

			_, err = client.CreateSpokeHaGw(spokeHaGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %w", err)
			}
			newHaGwEnabled = true
		}
	}

	haSubnet := getString(d, "ha_subnet")
	haZone := getString(d, "ha_zone")
	haEnabled := haSubnet != "" || haZone != ""

	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}

		singleAZ := getBool(d, "single_az_ha")
		if singleAZ {
			singleAZGateway.SingleAZ = "yes"
		} else {
			singleAZGateway.SingleAZ = "no"
		}

		if singleAZ {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA for %s: %w", singleAZGateway.GwName, err)
			}

			if haEnabled && manageHaGw {
				singleAZGatewayHA := &goaviatrix.Gateway{
					GwName: getString(d, "gw_name") + "-hagw",
				}
				err := client.EnableSingleAZGateway(singleAZGatewayHA)
				if err != nil {
					return fmt.Errorf("failed to enable single AZ GW HA for %s: %w", singleAZGatewayHA.GwName, err)
				}
			}
		} else {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA for %s: %w", singleAZGateway.GwName, err)
			}

			if haEnabled && manageHaGw {
				singleAZGatewayHA := &goaviatrix.Gateway{
					GwName: getString(d, "gw_name") + "-hagw",
				}
				err := client.DisableSingleAZGateway(singleAZGatewayHA)
				if err != nil {
					return fmt.Errorf("failed to disable single AZ GW HA for %s: %w", singleAZGatewayHA.GwName, err)
				}
			}
		}
	}

	if d.HasChange("ha_gw_size") && !newHaGwEnabled && manageHaGw {
		_, err := client.GetGateway(haGateway)
		if err != nil {
			// If HA gateway does not exist, don't try to change gateway size and continue with the rest of the updates
			// to the gateway
			if !errors.Is(err, goaviatrix.ErrNotFound) {
				return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway while trying to update HA Gw size: %w", err)
			}
		} else {
			if haGateway.VpcSize == "" {
				return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set")
			}
			err = client.UpdateGateway(haGateway)
			log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.VpcSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %w", err)
			}
		}
	}

	if d.HasChange("single_ip_snat") {
		enableSNat := getBool(d, "single_ip_snat")
		gw := &goaviatrix.Gateway{
			CloudType:   getInt(d, "cloud_type"),
			GatewayName: getString(d, "gw_name"),
		}
		if enableSNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable single_ip' mode SNAT: %w", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable 'single_ip' mode SNAT: %w", err)
			}
		}
	}

	if d.HasChange("enable_vpc_dns_server") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gw := &goaviatrix.Gateway{
			CloudType: getInt(d, "cloud_type"),
			GwName:    getString(d, "gw_name"),
		}

		enableVpcDnsServer := getBool(d, "enable_vpc_dns_server")
		if enableVpcDnsServer {
			err := client.EnableVpcDNSServer(gw)
			if err != nil {
				return fmt.Errorf("failed to enable VPC DNS Server: %w", err)
			}
		} else {
			err := client.DisableVpcDNSServer(gw)
			if err != nil {
				return fmt.Errorf("failed to disable VPC DNS Server: %w", err)
			}
		}

	} else if d.HasChange("enable_vpc_dns_server") {
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		gw := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		if learnedCidrsApproval {
			gw.LearnedCidrsApproval = "on"
			err := client.EnableSpokeLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to enable learned cidrs approval: %w", err)
			}
		} else {
			gw.LearnedCidrsApproval = "off"
			err := client.DisableSpokeLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to disable learned cidrs approval: %w", err)
			}
		}
	}

	if learnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gw := &goaviatrix.SpokeVpc{
			GwName:               getString(d, "gw_name"),
			ApprovedLearnedCidrs: approvedLearnedCidrs,
		}

		err := client.UpdateSpokePendingApprovedCidrs(gw)
		if err != nil {
			return fmt.Errorf("could not update approved CIDRs: %w", err)
		}
	}

	if d.HasChange("enable_encrypt_volume") {
		if getBool(d, "enable_encrypt_volume") {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              getString(d, "gw_name"),
				CustomerManagedKeys: getString(d, "customer_managed_keys"),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %w", gwEncVolume.GwName, err)
			}

			haSubnet := getString(d, "ha_subnet")
			haZone := getString(d, "ha_zone")
			haEnabled := haSubnet != "" || haZone != ""
			if haEnabled && manageHaGw {
				gwHAEncVolume := &goaviatrix.Gateway{
					GwName:              getString(d, "gw_name") + "-hagw",
					CustomerManagedKeys: getString(d, "customer_managed_keys"),
				}
				err := client.EnableEncryptVolume(gwHAEncVolume)
				if err != nil {
					return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %w", gwHAEncVolume.GwName, err)
				}
			}
		} else {
			return fmt.Errorf("can't disable Encrypt Volume for gateway: %s", gateway.GwName)
		}
	} else if d.HasChange("customer_managed_keys") {
		return fmt.Errorf("updating customer_managed_keys only is not allowed")
	}

	if d.HasChange("customized_spoke_vpc_routes") {
		o, n := d.GetChange("customized_spoke_vpc_routes")
		oldRouteList := strings.Split(mustString(o), ",")
		newRouteList := strings.Split(mustString(n), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                   getString(d, "gw_name"),
				CustomizedSpokeVpcRoutes: newRouteList,
			}
			err := client.EditGatewayCustomRoutes(transitGateway)
			log.Printf("[INFO] Customizeing routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to customize spoke vpc routes of spoke gateway: %s due to: %w", transitGateway.GwName, err)
			}
		}
	}

	if d.HasChange("filtered_spoke_vpc_routes") {
		o, n := d.GetChange("filtered_spoke_vpc_routes")
		oldRouteList := strings.Split(mustString(o), ",")
		newRouteList := strings.Split(mustString(n), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                 getString(d, "gw_name"),
				FilteredSpokeVpcRoutes: newRouteList,
			}
			err := client.EditGatewayFilterRoutes(transitGateway)
			log.Printf("[INFO] Editing filtered spoke vpc routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of spoke gateway: %s due to: %w", transitGateway.GwName, err)
			}
		}
	}

	if d.HasChange("included_advertised_spoke_routes") {
		o, n := d.GetChange("included_advertised_spoke_routes")
		oldRouteList := strings.Split(mustString(o), ",")
		newRouteList := strings.Split(mustString(n), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                getString(d, "gw_name"),
				AdvertisedSpokeRoutes: newRouteList,
			}
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			log.Printf("[INFO] Editing included advertised spoke vpc routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit included advertised spoke vpc routes of spoke gateway: %s due to: %w", transitGateway.GwName, err)
			}
		}
	}

	monitorGatewaySubnets := getBool(d, "enable_monitor_gateway_subnets")
	var excludedInstances []string
	for _, v := range getSet(d, "monitor_exclude_list").List() {
		excludedInstances = append(excludedInstances, mustString(v))
	}
	if !monitorGatewaySubnets && len(excludedInstances) != 0 {
		return fmt.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
	}
	if d.HasChange("enable_monitor_gateway_subnets") {
		if monitorGatewaySubnets {
			err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
			if err != nil {
				return fmt.Errorf("could not enable monitor gateway subnets: %w", err)
			}
		} else {
			err := client.DisableMonitorGatewaySubnets(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not disable monitor gateway subnets: %w", err)
			}
		}
	} else if d.HasChange("monitor_exclude_list") {
		err := client.DisableMonitorGatewaySubnets(gateway.GwName)
		if err != nil {
			return fmt.Errorf("could not disable monitor gateway subnets: %w", err)
		}
		err = client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
		if err != nil {
			return fmt.Errorf("could not enable monitor gateway subnets: %w", err)
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if getBool(d, "enable_jumbo_frame") {
			err := client.EnableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not enable jumbo frame for spoke gateway when updating: %w", err)
			}
		} else {
			err := client.DisableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not disable jumbo frame for spoke gateway when updating: %w", err)
			}
		}
	}

	if d.HasChange("enable_gro_gso") {
		if getBool(d, "enable_gro_gso") {
			err := client.EnableGroGso(gateway)
			if err != nil {
				return fmt.Errorf("couldn't enable GRO/GSO on spoke gateway when updating: %w", err)
			}
		} else {
			err := client.DisableGroGso(gateway)
			if err != nil {
				return fmt.Errorf("couldn't disable GRO/GSO on spoke gateway when updating: %w", err)
			}
		}
	}

	if d.HasChange("enable_private_vpc_default_route") {
		if getBool(d, "enable_private_vpc_default_route") {
			err := client.EnablePrivateVpcDefaultRoute(gateway)
			if err != nil {
				return fmt.Errorf("could not enable private vpc default route during spoke gateway update: %w", err)
			}
		} else {
			err := client.DisablePrivateVpcDefaultRoute(gateway)
			if err != nil {
				return fmt.Errorf("could not disable private vpc default route during spoke gateway update: %w", err)
			}
		}
	}

	if d.HasChange("enable_skip_public_route_table_update") {
		if getBool(d, "enable_skip_public_route_table_update") {
			err := client.EnableSkipPublicRouteUpdate(gateway)
			if err != nil {
				return fmt.Errorf("could not enable skip public route update during spoke gateway update: %w", err)
			}
		} else {
			err := client.DisableSkipPublicRouteUpdate(gateway)
			if err != nil {
				return fmt.Errorf("could not disable skip public route update during spoke gateway update: %w", err)
			}
		}
	}

	if d.HasChange("enable_auto_advertise_s2c_cidrs") {
		if getBool(d, "enable_auto_advertise_s2c_cidrs") {
			err := client.EnableAutoAdvertiseS2CCidrs(gateway)
			if err != nil {
				return fmt.Errorf("could not enable auto advertise s2c cidrs during spoke gateway update: %w", err)
			}
		} else {
			err := client.DisableAutoAdvertiseS2CCidrs(gateway)
			if err != nil {
				return fmt.Errorf("could not disable auto advertise s2c cidrs during spoke gateway update: %w", err)
			}
		}
	}

	if d.HasChange("tunnel_detection_time") {
		detectionTimeInterface, ok := d.GetOk("tunnel_detection_time")
		var detectionTime int
		if ok {
			detectionTime = mustInt(detectionTimeInterface)
		} else {
			var err error
			detectionTime, err = client.GetTunnelDetectionTime("Controller")
			if err != nil {
				return fmt.Errorf("could not get default tunnel detection time during Spoke Gateway update: %w", err)
			}
		}
		err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime)
		if err != nil {
			return fmt.Errorf("could not modify tunnel detection time during Spoke Gateway update: %w", err)
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		spokeGw := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		var spokeBgpManualSpokeAdvertiseCidrs []string
		for _, v := range getList(d, "spoke_bgp_manual_advertise_cidrs") {
			spokeBgpManualSpokeAdvertiseCidrs = append(spokeBgpManualSpokeAdvertiseCidrs, mustString(v))
		}
		spokeGw.BgpManualSpokeAdvertiseCidrs = strings.Join(spokeBgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetSpokeBgpManualAdvertisedNetworks(spokeGw)
		if err != nil {
			return fmt.Errorf("failed to set spoke bgp manual advertise CIDRs during Spoke Gateway update: %w", err)
		}
	}

	if d.HasChange("bgp_ecmp") {
		enabled := getBool(d, "bgp_ecmp")
		gateway := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		err := client.SetBgpEcmpSpoke(gateway, enabled)
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp during Spoke Gateway update: %w", err)
		}
	}

	if d.HasChange("enable_active_standby") || d.HasChange("enable_active_standby_preemptive") {
		gateway := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		if getBool(d, "enable_active_standby") {
			if getBool(d, "enable_active_standby_preemptive") {
				if err := client.EnableActiveStandbyPreemptiveSpoke(gateway); err != nil {
					return fmt.Errorf("could not enable Preemptive Mode for Active-Standby during Spoke Gateway update: %w", err)
				}
			} else {
				if err := client.EnableActiveStandbySpoke(gateway); err != nil {
					return fmt.Errorf("could not enable Active-Standby during Spoke Gateway update: %w", err)
				}
			}
		} else {
			if getBool(d, "enable_active_standby_preemptive") {
				return fmt.Errorf("could not enable Preemptive Mode with Active-Standby disabled")
			}
			if err := client.DisableActiveStandbySpoke(gateway); err != nil {
				return fmt.Errorf("could not disable Active-Standby during Spoke Gateway update: %w", err)
			}
		}
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		var prependASPath []string
		for _, v := range getList(d, "prepend_as_path") {
			prependASPath = append(prependASPath, mustString(v))
		}
		gateway := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}

		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(prependASPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPathSpoke(gateway, nil)
			if err != nil {
				return fmt.Errorf("could not delete prepend_as_path during Spoke Gateway update: %w", err)
			}
		}

		if d.HasChange("local_as_number") {
			localAsNumber := getString(d, "local_as_number")
			err := client.SetLocalASNumberSpoke(gateway, localAsNumber)
			if err != nil {
				return fmt.Errorf("could not set local_as_number: %w", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err := client.SetPrependASPathSpoke(gateway, prependASPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path during Spoke Gateway update: %w", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		bgpPollingTime := getInt(d, "bgp_polling_time")
		gateway := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		err := client.SetBgpPollingTimeSpoke(gateway, bgpPollingTime)
		if err != nil {
			return fmt.Errorf("could not update bgp polling time during Spoke Gateway update: %w", err)
		}
	}

	if d.HasChange("bgp_neighbor_status_polling_time") {
		bgpBfdPollingTime := getInt(d, "bgp_neighbor_status_polling_time")
		gateway := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		err := client.SetBgpBfdPollingTimeSpoke(gateway, bgpBfdPollingTime)
		if err != nil {
			return fmt.Errorf("could not update bgp neighbor status polling time during Spoke Gateway update: %w", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(gateway.GwName, getInt(d, "bgp_hold_time"))
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time during Spoke Gateway update: %w", err)
		}
	}

	if d.HasChange("disable_route_propagation") {
		disableRoutePropagation := getBool(d, "disable_route_propagation")
		enableBgp := getBool(d, "enable_bgp")
		if disableRoutePropagation && !enableBgp {
			return fmt.Errorf("disable route propagation is not supported for Non-BGP Spoke during Spoke Gateway update")
		}
		gw := &goaviatrix.SpokeVpc{
			GwName: getString(d, "gw_name"),
		}
		if disableRoutePropagation {
			err := client.DisableSpokeOnpremRoutePropagation(gw)
			if err != nil {
				return fmt.Errorf("failed to disable route propagation for Spoke %s during Spoke Gateway update: %w", gw.GwName, err)
			}
		} else {
			err := client.EnableSpokeOnpremRoutePropagation(gw)
			if err != nil {
				return fmt.Errorf("failed to enable route propagation for Spoke %s during Spoke Gateway update: %w", gw.GwName, err)
			}
		}
	}

	if d.HasChange("rx_queue_size") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("could not update rx_queue_size since it only supports AWS related cloud types")
		}
		gw := &goaviatrix.Gateway{
			GwName:      gateway.GwName,
			RxQueueSize: getString(d, "rx_queue_size"),
		}
		err := client.SetRxQueueSize(gw)
		if err != nil {
			return fmt.Errorf("could not modify rx queue size for spoke: %s during gateway update: %w", gw.GatewayName, err)
		}
		if haSubnet != "" || haZone != "" {
			haGwRxQueueSize := &goaviatrix.Gateway{
				GwName:      getString(d, "gw_name") + "-hagw",
				RxQueueSize: getString(d, "rx_queue_size"),
			}
			err := client.SetRxQueueSize(haGwRxQueueSize)
			if err != nil {
				return fmt.Errorf("could not modify rx queue size for spoke ha: %s during gateway update: %w", haGwRxQueueSize.GwName, err)
			}
		}
	}

	if d.HasChange("enable_global_vpc") {
		if getBool(d, "enable_global_vpc") {
			err := client.EnableGlobalVpc(gateway)
			if err != nil {
				return fmt.Errorf("could not enable global vpc during spoke gateway update: %w", err)
			}
		} else {
			err := client.DisableGlobalVpc(gateway)
			if err != nil {
				return fmt.Errorf("could not disable global vpc during spoke gateway update: %w", err)
			}
		}
	}

	if d.HasChange("enable_ipv6") {
		if getBool(d, "enable_ipv6") {
			err := client.EnableIPv6(gateway)
			if err != nil {
				return fmt.Errorf("couldn't enable IPv6 on spoke gateway when updating: %w", err)
			}
		} else {
			err := client.DisableIPv6(gateway)
			if err != nil {
				return fmt.Errorf("couldn't disable IPv6 on spoke gateway when updating: %w", err)
			}
		}
	}

	if d.HasChange("tunnel_encryption_cipher") || d.HasChange("tunnel_forward_secrecy") {
		encPolicy := getString(d, "tunnel_encryption_cipher")

		pfsPolicy := getString(d, "tunnel_forward_secrecy")

		err := client.SetGatewayPhase2Policy(gateway.GwName, encPolicy, pfsPolicy)
		if err != nil {
			return fmt.Errorf("could not set tunnel cipher settings during gateway update: %w", err)
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeGatewayRead(d, meta)
}

func resourceAviatrixSpokeGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke Gateway: %#v", gateway)

	// If HA is enabled, delete HA GW first.
	if getBool(d, "manage_ha_gateway") {
		haSubnet := getString(d, "ha_subnet")
		haZone := getString(d, "ha_zone")
		if haSubnet != "" || haZone != "" {
			// Delete HA Gw too
			gateway.GwName += "-hagw"
			err := client.DeleteGateway(gateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %w", err)
			}
		}
	}
	gateway.GwName = getString(d, "gw_name")

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Spoke Gateway: %w", err)
	}

	return nil
}
