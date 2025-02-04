package aviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	defaultLearnedCidrApprovalMode      = "gateway"
	defaultBgpPollingTime               = 50
	defaultBgpNeighborStatusPollingTime = 5
	defaultBgpHoldTime                  = 180
)

var typeMapping = map[string]string{
	"MANAGEMENT": "mgmt",
	"WAN":        "wan",
}

func resourceAviatrixTransitGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitGatewayCreate,
		Read:   resourceAviatrixTransitGatewayRead,
		Update: resourceAviatrixTransitGatewayUpdate,
		Delete: resourceAviatrixTransitGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixTransitGatewayMigrateState,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Type of cloud service provider, requires an integer value. Use 1 for AWS.",
				ValidateFunc: validateCloudType,
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
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
				Description:      "VPC-ID/VNet-Name/Site-ID of cloud provider.",
				DiffSuppressFunc: DiffSuppressFuncGatewayVpcId,
			},
			"device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Device ID for AEP EAT gateway.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Region of cloud provider.",
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
			"subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "Public Subnet Name.",
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
				Description: "AZ of subnet being created for Insane Mode Transit Gateway. Required for AWS if insane_mode is enabled.",
			},
			"ztp_file_download_path": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, old, _ string, _ *schema.ResourceData) bool {
					return old != ""
				},
				Description: "The location where the ZTP file will be stored locally.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("enable_private_oob").(bool) || d.Get("private_mode_lb_vpc_id").(string) != ""
				},
				Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"ha_subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
				Description: "HA Subnet. Required for enabling HA for AWS/AWSGov/AWSChina/Azure/OCI/Alibaba Cloud. " +
					"Optional for enabling HA for GCP gateway.",
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
				Description: "AZ of subnet being created for Insane Mode Transit HA Gateway. Required for AWS if insane_mode is enabled and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set).",
			},
			"ha_device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Device ID for HA AEP EAT gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable or disable Source NAT feature in 'single_ip' mode for this container.",
			},
			"enable_hybrid_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Sign of readiness for TGW connection.",
			},
			"connected_transit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify Connected Transit status.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable Insane Mode for Transit. Valid values: true, false. Supported for AWS/AWSGov, GCP, Azure and OCI. If insane mode is enabled, gateway size has to at least be c5 size for AWS and Standard_D3_v2 size for Azure.",
			},
			"enable_firenet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable firenet interfaces or not.",
			},
			"enable_gateway_load_balancer": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable firenet interfaces with AWS Gateway Load Balancer. Only valid when `enable_firenet` or `enable_transit_firenet`" +
					" are set to true and `cloud_type` = 1 (AWS). Currently AWS Gateway Load Balancer is only supported " +
					"in AWS regions us-west-2 and us-east-1. Valid values: true or false. Default value: false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Valid values: true, false.",
			},
			"enable_advertise_transit_cidr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable advertise transit VPC network CIDR.",
			},
			"bgp_manual_spoke_advertise_cidrs": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Intended CIDR list to be advertised to external bgp router.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS and AWSGov providers. Valid values: true, false. Default value: false.",
			},
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
					"filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s " +
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
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			"enable_egress_transit_firenet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable egress transit firenet interfaces or not.",
			},
			"enable_preserve_as_path": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable preserve as_path when advertising manual summary cidrs on transit gateway.",
			},
			"enable_transit_firenet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable transit firenet interfaces or not.",
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
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable encrypted transit approval for transit Gateway. Valid values: true, false.",
			},
			"learned_cidrs_approval_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      defaultLearnedCidrApprovalMode,
				ValidateFunc: validation.StringInSlice([]string{"gateway", "connection"}, false),
				Description: "Set the learned CIDRs approval mode. Only valid when 'enable_learned_cidrs_approval' is " +
					"set to true. If set to 'gateway', learned CIDR approval applies to ALL connections. If set to " +
					"'connection', learned CIDR approval is configured on a per connection basis. When configuring per " +
					"connection, use the enable_learned_cidrs_approval attribute within the connection resource to " +
					"toggle learned CIDR approval. Valid values: 'gateway' or 'connection'. Default value: 'gateway'.",
			},
			"approved_learned_cidrs": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional:    true,
				Description: "Approved learned CIDRs. Available as of provider version R2.21+.",
			},
			"bgp_polling_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpPollingTime,
				ValidateFunc: validation.IntBetween(10, 50),
				Description:  "BGP route polling time. Unit is in seconds. Valid values are between 10 and 50.",
			},
			"bgp_neighbor_status_polling_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpNeighborStatusPollingTime,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "BGP neighbor status polling time for BGP Transit Gateway. Unit is in seconds. Valid values are between 1 and 10.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "0"
				},
			},
			"prepend_as_path": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"local_as_number"},
				Description:  "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
			},
			"bgp_ecmp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Equal Cost Multi Path (ECMP) routing for the next hop.",
			},
			"enable_s2c_rx_balancing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable S2C receive packet CPU re-balancing on transit gateway.",
			},
			"enable_segmentation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable segmentation to allow association of transit gateway to security domains.",
			},
			"enable_active_standby": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Active-Standby Mode, available only with HA enabled.",
			},
			"enable_active_standby_preemptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.",
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
			"enable_bgp_over_lan": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Pre-allocate a network interface(eth4) for \"BGP over LAN\" functionality. Only valid for cloud_type = 4 (GCP) and 8 (Azure). Valid values: true or false. Default value: false. Available as of provider version R2.18+. Updatable as of provider version 3.0.3+.",
			},
			"bgp_lan_interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP Transit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: DiffSuppressFuncGCPVpcId,
							Description:      "VPC-ID of GCP cloud provider.",
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
							Description:  "Subnet Info.",
						},
					},
				},
			},
			"ha_bgp_lan_interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP HA Transit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: DiffSuppressFuncGCPVpcId,
							Description:      "VPC-ID of GCP cloud provider.",
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
							Description:  "Subnet Info.",
						},
					},
				},
			},
			"bgp_lan_interfaces_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Number of interfaces that will be created for BGP over LAN enabled Azure transit. Applies on HA Transit as well if enabled. Updatable as of provider version 3.0.3+.",
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
				Description: "Enable jumbo frame support for transit gateway. Valid values: true or false. Default value: true.",
			},
			"enable_gro_gso": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specify whether to disable GRO/GSO or not.",
			},
			"bgp_hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpHoldTime,
				ValidateFunc: validation.IntBetween(12, 360),
				Description:  "BGP Hold Time.",
			},
			"enable_transit_summarize_cidr_to_tgw": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable summarize CIDR to TGW.",
			},
			"enable_multi_tier_transit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Multi-tier Transit mode on transit gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A map of tags to assign to the transit gateway.",
			},
			"enable_spot_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(bool)
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
			"ha_private_mode_subnet_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: " Private Mode HA subnet availability zone.",
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
				Description:  "Public IP address that you want assigned to the HA Transit Gateway.",
			},
			"azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to this Transit Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"ha_azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to the HA Transit Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Changes the Aviatrix Transit Gateway ASN number before you setup Aviatrix Transit Gateway connection configurations.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"tunnel_detection_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(20, 600),
				Description:  "The IPSec tunnel down detection time for the transit gateway.",
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
				Description: "Security group used for the transit gateway.",
			},
			"ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA security group used for the transit gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the transit gateway.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the transit gateway created.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of HA transit gateway.",
			},
			"ha_gw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix transit gateway unique name of HA transit gateway.",
			},
			"ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of HA transit gateway.",
			},
			"lan_interface_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit gateway lan interface cidr.",
			},
			"ha_lan_interface_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit gateway lan interface cidr for the HA gateway.",
			},
			"bgp_lan_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "List of available BGP LAN interface IPs for transit external device connection creation. " +
					"Only supports GCP and Azure. Available as of provider version R2.21.0+.",
			},
			"ha_bgp_lan_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "List of available BGP LAN interface IPs for transit external device HA connection creation. " +
					"Only supports GCP and Azure. Available as of provider version R2.21.0+.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Transit Gateway created.",
			},
			"ha_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the HA Transit Gateway.",
			},
			"interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of WAN/Management interfaces for EAT gateway, each represented as a map.",
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
					},
				},
			},
			"ha_interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of WAN/Management interfaces for HA EAT gateway, each represented as a map. Required for enabling HA for edge gateway.",
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
					},
				},
			},
			"interface_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of interface names mapped to interface types and indices. Only required for ESXI.",
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
			"peer_backup_logical_ifname": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Peer backup logical interface name for the edge transit gateway (e.g., 'wan0', 'wan1').",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"peer_connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Connection type for the edge transit gateway (e.g., 'public', 'private').",
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
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
			"ha_management_egress_ip_prefix_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of management egress gateway IP/prefix for HA EAT.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAviatrixTransitGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	cloudType := d.Get("cloud_type").(int)
	flag := false
	// create edge transit gateway for AEP & Equinix
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		err := createEdgeTransitGateway(d, client, cloudType)
		if err != nil {
			return err
		}
		defer resourceAviatrixTransitGatewayReadIfRequired(d, meta, &flag)
	} else {
		gateway := &goaviatrix.TransitVpc{
			CloudType:                d.Get("cloud_type").(int),
			AccountName:              d.Get("account_name").(string),
			GwName:                   d.Get("gw_name").(string),
			VpcID:                    d.Get("vpc_id").(string),
			VpcSize:                  d.Get("gw_size").(string),
			EnableHybridConnection:   d.Get("enable_hybrid_connection").(bool),
			EnableSummarizeCidrToTgw: d.Get("enable_transit_summarize_cidr_to_tgw").(bool),
			Subnet:                   d.Get("subnet").(string),
			AvailabilityDomain:       d.Get("availability_domain").(string),
			FaultDomain:              d.Get("fault_domain").(string),
			ApprovedLearnedCidrs:     getStringSet(d, "approved_learned_cidrs"),
			Transit:                  true,
		}

		enableNAT := d.Get("single_ip_snat").(bool)
		if enableNAT {
			gateway.EnableNAT = "yes"
		} else {
			gateway.EnableNAT = "no"
		}

		singleAZ := d.Get("single_az_ha").(bool)
		if singleAZ {
			gateway.SingleAzHa = "enabled"
		} else {
			gateway.SingleAzHa = "disabled"
		}

		connectedTransit := d.Get("connected_transit").(bool)
		if connectedTransit {
			gateway.ConnectedTransit = "yes"
		} else {
			gateway.ConnectedTransit = "no"
		}

		zone := d.Get("zone").(string)
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && zone != "" {
			return fmt.Errorf("attribute 'zone' is only for use with Azure (8), Azure GOV (32) and Azure CHINA (2048)")
		}
		if zone != "" {
			// The API uses the same string field to hold both subnet and zone
			// parameters.
			gateway.Subnet = fmt.Sprintf("%s~~%s~~", d.Get("subnet").(string), zone)
		}

		if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			gateway.VpcID = d.Get("vpc_id").(string)
			if gateway.VpcID == "" {
				return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw")
			}
		} else {
			return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
		}

		if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			gateway.VpcRegion = d.Get("vpc_reg").(string)
		} else if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			// for gcp, rest api asks for "zone" rather than vpc region
			gateway.Zone = d.Get("vpc_reg").(string)
		} else {
			return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
		}

		insaneMode := d.Get("insane_mode").(bool)
		if insaneMode {
			// Insane Mode encryption is not supported in China regions
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|
				goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
				return fmt.Errorf("insane_mode is only supported for AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
			}
			if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
				if d.Get("insane_mode_az").(string) == "" {
					return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768)")
				}
				if d.Get("ha_subnet").(string) != "" && d.Get("ha_insane_mode_az").(string) == "" {
					return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768) clouds and ha_subnet is set")
				}
				// Append availability zone to subnet
				var strs []string
				insaneModeAz := d.Get("insane_mode_az").(string)
				strs = append(strs, gateway.Subnet, insaneModeAz)
				gateway.Subnet = strings.Join(strs, "~~")
			}
			gateway.InsaneMode = "yes"
		} else {
			gateway.InsaneMode = "no"
		}

		if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain == "" || gateway.FaultDomain == "") {
			return fmt.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
		}
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain != "" || gateway.FaultDomain != "") {
			return fmt.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
		}

		haSubnet := d.Get("ha_subnet").(string)
		haZone := d.Get("ha_zone").(string)
		haAvailabilityDomain := d.Get("ha_availability_domain").(string)
		haFaultDomain := d.Get("ha_fault_domain").(string)

		if haZone != "" && !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("'ha_zone' is only valid for GCP and Azure providers when enabling HA")
		}
		if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) && haSubnet != "" && haZone == "" {
			return fmt.Errorf("'ha_zone' must be set to enable HA on GCP, cannot enable HA with only 'ha_subnet'")
		}
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" && haZone != "" {
			return fmt.Errorf("'ha_subnet' must be provided to enable HA on Azure, cannot enable HA with only 'ha_zone'")
		}
		haGwSize := d.Get("ha_gw_size").(string)
		if haSubnet == "" && haZone == "" && haGwSize != "" {
			return fmt.Errorf("'ha_gw_size' is only required if enabling HA")
		}
		if haGwSize == "" && haSubnet != "" {
			return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet is set")
		}
		if haSubnet != "" {
			if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain == "" || haFaultDomain == "") {
				return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are required to enable Peering HA on OCI")
			}
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain != "" || haFaultDomain != "") {
				return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are only valid for OCI")
			}
		}

		enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
		customerManagedKeys := d.Get("customer_managed_keys").(string)
		if enableEncryptVolume && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768) providers")
		}
		if customerManagedKeys != "" {
			if !enableEncryptVolume {
				return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
			}
			gateway.CustomerManagedKeys = customerManagedKeys
		}
		if !enableEncryptVolume && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			gateway.EncVolume = "no"
		}

		enableFireNet := d.Get("enable_firenet").(bool)
		enableGatewayLoadBalancer := d.Get("enable_gateway_load_balancer").(bool)
		enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
		if enableFireNet && enableTransitFireNet {
			return fmt.Errorf("can't enable firenet function and transit firenet function at the same time")
		}
		lanVpcID := d.Get("lan_vpc_id").(string)
		lanPrivateSubnet := d.Get("lan_private_subnet").(string)
		// Transit FireNet function is not supported for AWS China or Azure China
		if enableFireNet && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSChina|goaviatrix.AzureChina) {
			return fmt.Errorf("'enable_firenet' is not supported in AWSChina (1024) or AzureChina (2048)")
		}
		if enableTransitFireNet {
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
				return fmt.Errorf("'enable_transit_firenet' is only supported in AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), Azure China (2048)")
			}
			if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
				gateway.EnableTransitFireNet = true
			}
			if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
				if lanVpcID == "" || lanPrivateSubnet == "" {
					return fmt.Errorf("'lan_vpc_id' and 'lan_private_subnet' are required when 'cloud_type' = 4 (GCP) and 'enable_transit_firenet' = true")
				}
				gateway.LanVpcID = lanVpcID
				gateway.LanPrivateSubnet = lanPrivateSubnet
			}
		}
		if (!enableTransitFireNet || !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes)) && (lanVpcID != "" || lanPrivateSubnet != "") {
			return fmt.Errorf("'lan_vpc_id' and 'lan_private_subnet' are only valid when 'cloud_type' = 4 (GCP) and 'enable_transit_firenet' = true")
		}
		if enableGatewayLoadBalancer && !enableFireNet && !enableTransitFireNet {
			return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
		}
		if enableGatewayLoadBalancer && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWS) {
			return fmt.Errorf("'enable_gateway_load_balancer' is only supported by AWS (1)")
		}
		enableEgressTransitFireNet := d.Get("enable_egress_transit_firenet").(bool)
		if enableEgressTransitFireNet && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			return fmt.Errorf("'enable_egress_transit_firenet' is only supported by AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), Azure China(2048)")
		}
		if enableEgressTransitFireNet && !enableTransitFireNet {
			return fmt.Errorf("'enable_egress_transit_firenet' requires 'enable_transit_firenet' to be set to true")
		}
		if enableEgressTransitFireNet && connectedTransit {
			return fmt.Errorf("'enable_egress_transit_firenet' requires 'connected_transit' to be set to false")
		}

		enableTransitPreserveAsPath := d.Get("enable_preserve_as_path").(bool)

		learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
		if learnedCidrsApproval {
			gateway.LearnedCidrsApproval = "yes"
		}

		if learnedCidrsApproval && d.Get("learned_cidrs_approval_mode").(string) == "connection" {
			return fmt.Errorf("'enable_learned_cidrs_approval' must be false if 'learned_cidrs_approval_mode' is set to 'connection'")
		}

		if !learnedCidrsApproval && len(gateway.ApprovedLearnedCidrs) != 0 {
			return fmt.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
		}

		enableMonitorSubnets := d.Get("enable_monitor_gateway_subnets").(bool)
		var excludedInstances []string
		for _, v := range d.Get("monitor_exclude_list").(*schema.Set).List() {
			excludedInstances = append(excludedInstances, v.(string))
		}
		// Enable monitor gateway subnets does not work with AWSChina
		if enableMonitorSubnets && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes^goaviatrix.AWSChina) {
			return fmt.Errorf("'enable_monitor_gateway_subnets' is only valid for AWS (1), AWSGov (256), AWS Top Secret (16384) or AWS Secret (32768)")
		}
		if !enableMonitorSubnets && len(excludedInstances) != 0 {
			return fmt.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
		}

		bgpOverLan := d.Get("enable_bgp_over_lan").(bool)
		if bgpOverLan && !(goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCP)) {
			return fmt.Errorf("'enable_bgp_over_lan' is only valid for GCP (4), Azure (8), AzureGov (32) or AzureChina (2048)")
		}
		bgpLanInterfacesCount, isCountSet := d.GetOk("bgp_lan_interfaces_count")
		if isCountSet && (!bgpOverLan || !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes)) {
			return fmt.Errorf("'bgp_lan_interfaces_count' is only valid for BGP over LAN enabled transit for Azure (8), AzureGov (32) or AzureChina (2048)")
		} else if !isCountSet && bgpOverLan && goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("please specify 'bgp_lan_interfaces_count' for BGP over LAN enabled Azure transit: %s", gateway.GwName)
		}
		var bgpLanVpcID []string
		var bgpLanSpecifySubnet []string
		var haBgpLanVpcID []string
		var haBgpLanSpecifySubnet []string
		for _, bgpInterface := range d.Get("bgp_lan_interfaces").([]interface{}) {
			item := bgpInterface.(map[string]interface{})
			bgpLanVpcID = append(bgpLanVpcID, item["vpc_id"].(string))
			bgpLanSpecifySubnet = append(bgpLanSpecifySubnet, item["subnet"].(string))
		}
		for _, haBgpInterface := range d.Get("ha_bgp_lan_interfaces").([]interface{}) {
			item := haBgpInterface.(map[string]interface{})
			haBgpLanVpcID = append(haBgpLanVpcID, item["vpc_id"].(string))
			haBgpLanSpecifySubnet = append(haBgpLanSpecifySubnet, item["subnet"].(string))
		}
		if !bgpOverLan && len(bgpLanVpcID) != 0 {
			return fmt.Errorf("'bgp_lan_interfaces' is only valid with enable_bgp_over_lan being set true")
		}
		if (!bgpOverLan || haSubnet == "") && len(haBgpLanVpcID) != 0 {
			return fmt.Errorf("'ha_bgp_lan_interfaces' is only valid with enable_bgp_over_lan is set true and HA is enabled")
		}
		if (len(bgpLanVpcID) != 0 || len(haBgpLanVpcID) != 0) && !(goaviatrix.IsCloudType(cloudType, goaviatrix.GCP)) {
			return fmt.Errorf("'bgp_lan_interfaces' and 'ha_bgp_lan_interfaces' are only valid for GCP (4)")
		}
		if bgpOverLan {
			gateway.BgpOverLan = true
			if goaviatrix.IsCloudType(cloudType, goaviatrix.GCP) {
				if len(bgpLanVpcID) == 0 {
					return fmt.Errorf("missing bgp_lan_interfaces for creating GCP transit gateway with BGP over LAN enabled")
				}
				if (haSubnet != "" || haZone != "") && len(haBgpLanVpcID) == 0 {
					return fmt.Errorf("missing ha_bgp_lan_interfaces for creating GCP HA transit gateway with BGP over LAN enabled")
				}
				gateway.BgpLanVpcID = strings.Join(bgpLanVpcID, ",")
				gateway.BgpLanSpecifySubnet = strings.Join(bgpLanSpecifySubnet, ",")
			} else if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				gateway.BgpLanInterfacesCount = bgpLanInterfacesCount.(int)
			}
		}

		enablePrivateOob := d.Get("enable_private_oob").(bool)
		oobManagementSubnet := d.Get("oob_management_subnet").(string)
		oobAvailabilityZone := d.Get("oob_availability_zone").(string)
		haOobManagementSubnet := d.Get("ha_oob_management_subnet").(string)
		haOobAvailabilityZone := d.Get("ha_oob_availability_zone").(string)

		if enablePrivateOob {
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_private_oob' is only valid for AWS (1), AWSGov (256) AWSChina (1024),, AWS Top Secret (16384) or AWS Secret (32768)")
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
					return fmt.Errorf("\"ha_oob_management_sbunet\" must be empty if \"ha_subnet\" is empty")
				}
			}

			gateway.EnablePrivateOob = "on"
			gateway.Subnet = gateway.Subnet + "~~" + oobAvailabilityZone
			gateway.OobManagementSubnet = oobManagementSubnet + "~~" + oobAvailabilityZone
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
				return fmt.Errorf("\"ha_oob_management_sbunet\" must be empty if \"enable_private_oob\" is false")
			}
		}

		enableMultitierTransit := d.Get("enable_multi_tier_transit").(bool)
		if enableMultitierTransit {
			if d.Get("local_as_number") == "" {
				return fmt.Errorf("local_as_number required to enable multi tier transit")
			}
		}

		_, tagsOk := d.GetOk("tags")
		if tagsOk {
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
				return errors.New("error creating transit gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
			}
			tagsMap, err := extractTags(d, gateway.CloudType)
			if err != nil {
				return fmt.Errorf("error creating tags for transit gateway: %w", err)
			}
			tagJson, err := TagsMapToJson(tagsMap)
			if err != nil {
				return fmt.Errorf("failed to add tags when creating transit gateway: %w", err)
			}
			gateway.TagJson = tagJson
		}

		enableActiveStandby := d.Get("enable_active_standby").(bool)
		if haSubnet == "" && haZone == "" && enableActiveStandby {
			return fmt.Errorf("could not configure Active-Standby as HA is not enabled")
		}
		enableActiveStandbyPreemptive := d.Get("enable_active_standby_preemptive").(bool)
		if !enableActiveStandby && enableActiveStandbyPreemptive {
			return fmt.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
		}

		enableSpotInstance := d.Get("enable_spot_instance").(bool)
		spotPrice := d.Get("spot_price").(string)
		deleteSpot := d.Get("delete_spot").(bool)
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

		rxQueueSize := d.Get("rx_queue_size").(string)
		if rxQueueSize != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("rx_queue_size only supports AWS related cloud types")
		}

		privateModeInfo, _ := client.GetPrivateModeInfo(context.Background())
		if !enablePrivateOob && !privateModeInfo.EnablePrivateMode {
			allocateNewEip := d.Get("allocate_new_eip").(bool)
			if allocateNewEip {
				gateway.ReuseEip = "off"
			} else {
				gateway.ReuseEip = "on"

				if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
					return fmt.Errorf("failed to create transit gateway: 'allocate_new_eip' can only be set to 'false' when cloud_type is AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048) or AWS Top Secret (16384)")
				}
				if _, ok := d.GetOk("eip"); !ok {
					return fmt.Errorf("failed to create transit gateway: 'eip' must be set when 'allocate_new_eip' is false")
				}
				azureEipName, azureEipNameOk := d.GetOk("azure_eip_name_resource_group")
				if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
					// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
					if !azureEipNameOk {
						return fmt.Errorf("failed to create transit gateway: 'azure_eip_name_resource_group' must be set when 'allocate_new_eip' is false and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
					}
					gateway.Eip = fmt.Sprintf("%s:%s", azureEipName.(string), d.Get("eip").(string))
				} else {
					if azureEipNameOk {
						return fmt.Errorf("failed to create transit gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
					}
					gateway.Eip = d.Get("eip").(string)
				}
			}
		}

		if privateModeInfo.EnablePrivateMode {
			if privateModeSubnetZone, ok := d.GetOk("private_mode_subnet_zone"); ok {
				gateway.Subnet = fmt.Sprintf("%s~~%s", gateway.Subnet, privateModeSubnetZone.(string))
			} else {
				if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
					return fmt.Errorf("%q must be set when creating a Transit Gateway in AWS with Private Mode enabled on the Controller", "private_mode_subnet_zone")
				}
			}

			if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
				if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
					return fmt.Errorf("private mode is only supported in AWS and Azure. %q must be empty", "private_mode_lb_vpc_id")
				}

				gateway.LbVpcID, ok = d.Get("private_mode_lb_vpc_id").(string)
				if !ok {
					return fmt.Errorf("%q must be set when creating a Transit Gateway in AWS or Azure with Private Mode enabled on the Controller", "private_mode_lb_vpc_id")
				}
			}
		} else {
			if _, ok := d.GetOk("private_mode_subnet_zone"); ok {
				return fmt.Errorf("%q is only valid when Private Mode is enabled on the Controller", "private_mode_subnet_zone")
			}
			if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
				return fmt.Errorf("%q is only valid on when Private Mode is enabled", "private_mode_lb_vpc_id")
			}
		}

		log.Printf("[INFO] Creating Aviatrix Transit Gateway: %#v", gateway)

		d.SetId(gateway.GwName)
		defer resourceAviatrixTransitGatewayReadIfRequired(d, meta, &flag)
		err := client.LaunchTransitVpc(gateway)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Transit Gateway: %s", err)
		}

		if !singleAZ {
			singleAZGateway := &goaviatrix.Gateway{
				GwName:   d.Get("gw_name").(string),
				SingleAZ: "disabled",
			}

			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)

			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
			}
		}

		if haSubnet != "" || haZone != "" {
			// Enable HA
			transitHaGw := &goaviatrix.TransitHaGateway{
				PrimaryGwName: d.Get("gw_name").(string),
				GwName:        d.Get("gw_name").(string) + "-hagw",
				Subnet:        haSubnet,
				Zone:          haZone,
				Eip:           d.Get("ha_eip").(string),
				InsaneMode:    "no",
				BgpLanVpcID:   strings.Join(haBgpLanVpcID, ","),
				BgpLanSubnet:  strings.Join(haBgpLanSpecifySubnet, ","),
			}

			if insaneMode {
				if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
					var haStrs []string
					insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
					haStrs = append(haStrs, haSubnet, insaneModeHaAz)
					haSubnet = strings.Join(haStrs, "~~")
					transitHaGw.Subnet = haSubnet
				}
				transitHaGw.InsaneMode = "yes"
			}

			if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
				if haZone == "" {
					return fmt.Errorf("no ha_zone is provided for enabling Transit HA gateway: %s", transitHaGw.GwName)
				}
			} else if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				if haZone != "" {
					transitHaGw.Subnet = fmt.Sprintf("%s~~%s~~", haSubnet, haZone)
				}
			} else if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) {
				transitHaGw.AvailabilityDomain = haAvailabilityDomain
				transitHaGw.FaultDomain = haFaultDomain
			}

			if privateModeInfo.EnablePrivateMode {
				haPrivateModeSubnetZone := d.Get("ha_private_mode_subnet_zone").(string)
				if haPrivateModeSubnetZone == "" && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
					return fmt.Errorf("%q must be set when creating a Transit HA Gateway in AWS with Private Mode enabled on the Controller", "ha_private_mode_subnet_zone")
				}
				transitHaGw.Subnet = haSubnet + "~~" + haPrivateModeSubnetZone
			}

			haAzureEipName, haAzureEipNameOk := d.GetOk("ha_azure_eip_name_resource_group")
			if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				if transitHaGw.Eip != "" {
					// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
					if !haAzureEipNameOk {
						return fmt.Errorf("failed to create HA Transit Gateway: 'ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
					}
					transitHaGw.Eip = fmt.Sprintf("%s:%s", haAzureEipName.(string), transitHaGw.Eip)
				} else {
					if haAzureEipNameOk {
						return fmt.Errorf("failed to create HA Transit Gateway: 'ha_azure_eip_name_resource_group' must be empty when 'ha_eip' is empty")
					}
				}
			} else {
				if haAzureEipNameOk {
					return fmt.Errorf("failed to create HA Transit Gateway: 'ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
				}
			}

			if bgpOverLan && goaviatrix.IsCloudType(cloudType, goaviatrix.GCP) {
				transitHaGw.BgpLanVpcID = strings.Join(haBgpLanVpcID, ",")
				transitHaGw.BgpLanSubnet = strings.Join(haBgpLanSpecifySubnet, ",")
			}

			log.Printf("[INFO] Enabling HA on Transit Gateway: %#v", haSubnet)

			_, err := client.CreateTransitHaGw(transitHaGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
			}

			// Resize HA Gateway
			log.Printf("[INFO]Resizing Transit HA Gateway: %#v", haGwSize)

			if haGwSize != gateway.VpcSize {
				if haGwSize == "" {
					return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
						"ha_subnet is set")
				}

				haGateway := &goaviatrix.Gateway{
					CloudType: d.Get("cloud_type").(int),
					GwName:    d.Get("gw_name").(string) + "-hagw",
					VpcSize:   d.Get("ha_gw_size").(string),
				}

				log.Printf("[INFO] Resizing Transit HA GAteway size to: %s ", haGateway.VpcSize)

				err = client.UpdateGateway(haGateway)
				if err != nil {
					return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
				}
			}
		}

		enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
		if enableHybridConnection {
			if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_hybrid_connection' is only supported by AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}

			err := client.AttachTransitGWForHybrid(gateway)
			if err != nil {
				return fmt.Errorf("failed to enable transit GW for Hybrid: %s", err)
			}
		}

		if connectedTransit {
			err := client.EnableConnectedTransit(gateway)
			if err != nil {
				return fmt.Errorf("failed to enable connected transit: %s", err)
			}
		}

		if enableFireNet {
			if enableGatewayLoadBalancer {
				err := client.EnableGatewayFireNetInterfacesWithGWLB(gateway)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for FireNet Interfaces with Gateway Load Balancer enabled: %s", err)
				}
			} else {
				err := client.EnableGatewayFireNetInterfaces(gateway)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
				}
			}
		}

		enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && enableVpcDnsServer {
			gwVpcDnsServer := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}

			log.Printf("[INFO] Enable VPC DNS Server: %#v", gwVpcDnsServer)

			err := client.EnableVpcDnsServer(gwVpcDnsServer)
			if err != nil {
				return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
			}
		} else if enableVpcDnsServer {
			return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
		}

		enableAdvertiseTransitCidr := d.Get("enable_advertise_transit_cidr").(bool)
		if enableAdvertiseTransitCidr {
			err := client.EnableAdvertiseTransitCidr(gateway)
			if err != nil {
				return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
			}
		}

		bgpManualSpokeAdvertiseCidrs := d.Get("bgp_manual_spoke_advertise_cidrs").(string)
		if bgpManualSpokeAdvertiseCidrs != "" {
			gateway.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
			err := client.SetBgpManualSpokeAdvertisedNetworks(gateway)
			if err != nil {
				return fmt.Errorf("failed to set BGP Manual Spoke Advertise Cidrs: %s", err)
			}
		}

		if customizedSpokeVpcRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedSpokeVpcRoutes != "" {
			transitGateway := &goaviatrix.Gateway{
				GwName:                   d.Get("gw_name").(string),
				CustomizedSpokeVpcRoutes: strings.Split(customizedSpokeVpcRoutes, ","),
			}
			for i := 0; ; i++ {
				log.Printf("[INFO] Editing customized routes of transit gateway: %s ", transitGateway.GwName)
				err := client.EditGatewayCustomRoutes(transitGateway)
				if err == nil {
					break
				}
				if i <= 10 && strings.Contains(err.Error(), "when it is down") {
					time.Sleep(10 * time.Second)
				} else {
					return fmt.Errorf("failed to customize spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
				}
			}
		}

		if filteredSpokeVpcRoutes := d.Get("filtered_spoke_vpc_routes").(string); filteredSpokeVpcRoutes != "" {
			transitGateway := &goaviatrix.Gateway{
				GwName:                 d.Get("gw_name").(string),
				FilteredSpokeVpcRoutes: strings.Split(filteredSpokeVpcRoutes, ","),
			}
			for i := 0; ; i++ {
				log.Printf("[INFO] Editing filtered routes of transit gateway: %s ", transitGateway.GwName)
				err := client.EditGatewayFilterRoutes(transitGateway)
				if err == nil {
					break
				}
				if i <= 10 && strings.Contains(err.Error(), "when it is down") {
					time.Sleep(10 * time.Second)
				} else {
					return fmt.Errorf("failed to edit filtered spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
				}
			}
		}

		if advertisedSpokeRoutesExclude := d.Get("excluded_advertised_spoke_routes").(string); advertisedSpokeRoutesExclude != "" {
			transitGateway := &goaviatrix.Gateway{
				GwName:                d.Get("gw_name").(string),
				AdvertisedSpokeRoutes: strings.Split(advertisedSpokeRoutesExclude, ","),
			}
			for i := 0; ; i++ {
				log.Printf("[INFO] Editing customized routes advertisement of transit gateway: %s ", transitGateway.GwName)
				err := client.EditGatewayAdvertisedCidr(transitGateway)
				if err == nil {
					break
				}
				if i <= 10 && strings.Contains(err.Error(), "when it is down") {
					time.Sleep(10 * time.Second)
				} else {
					return fmt.Errorf("failed to edit advertised spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
				}
			}
		}

		if enableTransitFireNet && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			if enableGatewayLoadBalancer {
				err := client.EnableTransitFireNetWithGWLB(gwTransitFireNet)
				if err != nil {
					return fmt.Errorf("failed to enable transit firenet with Gateway Load Balancer enabled: %w", err)
				}
			} else {
				err := client.EnableTransitFireNet(gwTransitFireNet)
				if err != nil {
					return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
				}
			}
		}

		if val, ok := d.GetOk("bgp_polling_time"); ok {
			err := client.SetBgpPollingTime(gateway, val.(int))
			if err != nil {
				return fmt.Errorf("could not set bgp polling time: %w", err)
			}
		}

		if val, ok := d.GetOk("bgp_neighbor_status_polling_time"); ok {
			bgp_neighbor_status_polling_time := val.(int)
			if bgp_neighbor_status_polling_time >= 1 && bgp_neighbor_status_polling_time != defaultBgpNeighborStatusPollingTime {
				err := client.SetBgpBfdPollingTime(gateway, bgp_neighbor_status_polling_time)
				if err != nil {
					return fmt.Errorf("could not set bgp neighbor status polling time: %w", err)
				}
			}
		}

		if val, ok := d.GetOk("local_as_number"); ok {
			err := client.SetLocalASNumber(gateway, val.(string))
			if err != nil {
				return fmt.Errorf("could not set local_as_number: %w", err)
			}
		}

		if val, ok := d.GetOk("prepend_as_path"); ok {
			var prependASPath []string
			slice := val.([]interface{})
			for _, v := range slice {
				prependASPath = append(prependASPath, v.(string))
			}
			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path: %w", err)
			}
		}

		if val, ok := d.GetOk("bgp_ecmp"); ok {
			err := client.SetBgpEcmp(gateway, val.(bool))
			if err != nil {
				return fmt.Errorf("could not set bgp_ecmp: %w", err)
			}
		}

		if d.Get("enable_segmentation").(bool) {
			if err := client.EnableSegmentation(gateway); err != nil {
				return fmt.Errorf("could not enable segmentation: %w", err)
			}
		}

		if enableEgressTransitFireNet {
			err := client.EnableEgressTransitFirenet(gateway)
			if err != nil {
				return fmt.Errorf("could not enable egress transit firenet: %w", err)
			}
		}

		if enableTransitPreserveAsPath {
			err := client.EnableTransitPreserveAsPath(gateway)
			if err != nil {
				return fmt.Errorf("could not enable transit preserve as path: %w", err)
			}
		}

		if enableActiveStandby {
			if enableActiveStandbyPreemptive {
				if err := client.EnableActiveStandbyPreemptive(gateway); err != nil {
					return fmt.Errorf("could not enable Preemptive Mode for Active-Standby: %w", err)
				}
			} else {
				if err := client.EnableActiveStandby(gateway); err != nil {
					return fmt.Errorf("could not enable Active-Standby: %w", err)
				}
			}
		}

		approvalMode := d.Get("learned_cidrs_approval_mode").(string)
		if approvalMode != defaultLearnedCidrApprovalMode {
			err := client.SetTransitLearnedCIDRsApprovalMode(gateway, approvalMode)
			if err != nil {
				return fmt.Errorf("could not set learned CIDRs approval mode to %q: %v", approvalMode, err)
			}
		}

		if learnedCidrsApproval && len(gateway.ApprovedLearnedCidrs) != 0 {
			err = client.UpdateTransitPendingApprovedCidrs(gateway)
			if err != nil {
				return fmt.Errorf("could not update approved CIDRs: %w", err)
			}
		}

		var customizedTransitVpcRoutes []string
		for _, v := range d.Get("customized_transit_vpc_routes").(*schema.Set).List() {
			customizedTransitVpcRoutes = append(customizedTransitVpcRoutes, v.(string))
		}
		if len(customizedTransitVpcRoutes) != 0 {
			err := client.UpdateTransitGatewayCustomizedVpcRoute(gateway.GwName, customizedTransitVpcRoutes)
			if err != nil {
				return fmt.Errorf("couldn't update transit gateway customized vpc route: %s", err)
			}
		}

		if enableMonitorSubnets {
			err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
			if err != nil {
				return fmt.Errorf("could not enable monitor gateway subnets: %w", err)
			}
		}

		if !d.Get("enable_jumbo_frame").(bool) {
			gw := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}

			err := client.DisableJumboFrame(gw)
			if err != nil {
				return fmt.Errorf("could not disable jumbo frame for transit gateway: %w", err)
			}
		}

		if !d.Get("enable_gro_gso").(bool) {
			gw := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.DisableGroGso(gw)
			if err != nil {
				return fmt.Errorf("couldn't disable GRO/GSO on transit gateway: %s", err)
			}
		}

		if holdTime := d.Get("bgp_hold_time").(int); holdTime != defaultBgpHoldTime {
			err := client.ChangeBgpHoldTime(gateway.GwName, holdTime)
			if err != nil {
				return fmt.Errorf("could not change BGP Hold Time after Transit Gateway creation: %w", err)
			}
		}

		if gateway.EnableSummarizeCidrToTgw {
			err = client.EnableSummarizeCidrToTgw(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable summarize cidr to tgw: %w", err)
			}
		}

		if enableMultitierTransit {
			err = client.EnableMultitierTransit(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable multi tier transit: %w", err)
			}
		}

		enableS2CRxBalancing := d.Get("enable_s2c_rx_balancing").(bool)
		if enableS2CRxBalancing {
			err = client.EnableS2CRxBalancing(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable S2C receive packet CPU re-balancing on transit %s: %v", gateway.GwName, err)
			}
		}

		if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
			err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime.(int))
			if err != nil {
				return fmt.Errorf("could not set tunnel detection time during Transit Gateway creation: %w", err)
			}
		}

		if rxQueueSize != "" {
			gwRxQueueSize := &goaviatrix.Gateway{
				GwName:      d.Get("gw_name").(string),
				RxQueueSize: rxQueueSize,
			}
			err := client.SetRxQueueSize(gwRxQueueSize)
			if err != nil {
				return fmt.Errorf("failed to set rx queue size for transit %s: %s", gateway.GwName, err)
			}
			if haSubnet != "" || haZone != "" {
				haGwRxQueueSize := &goaviatrix.Gateway{
					GwName:      d.Get("gw_name").(string) + "-hagw",
					RxQueueSize: rxQueueSize,
				}
				err := client.SetRxQueueSize(haGwRxQueueSize)
				if err != nil {
					return fmt.Errorf("failed to set rx queue size for transit ha %s : %s", haGwRxQueueSize.GwName, err)
				}
			}
		}
	}
	return resourceAviatrixTransitGatewayReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitGatewayReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitGatewayRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	ignoreTagsConfig := client.IgnoreTagsConfig

	var isImport bool
	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
		gwName = id
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      gwName,
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit Gateway: %s", err)
	}

	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)

	d.Set("cloud_type", gw.CloudType)
	d.Set("account_name", gw.AccountName)
	d.Set("gw_name", gw.GwName)
	d.Set("gw_size", gw.GwSize)

	// edge cloud type
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		d.Set("vpc_id", gw.VpcID)
		d.Set("bgp_lan_ip_list", nil)
		d.Set("ha_bgp_lan_ip_list", nil)
		if gw.DeviceID != "" {
			_ = d.Set("device_id", gw.DeviceID)
		}
		// set interfaces
		if len(gw.Interfaces) != 0 {
			interfaces := setInterfaceDetails(gw.Interfaces)
			if err = d.Set("interfaces", interfaces); err != nil {
				return fmt.Errorf("could not set interfaces into state: %w", err)
			}
		}
		// set interface mapping
		if len(gw.InterfaceMapping) != 0 {
			interfaceMapping := setInterfaceMappingDetails(gw.InterfaceMapping)
			if err = d.Set("interface_mapping", interfaceMapping); err != nil {
				return fmt.Errorf("could not set interface mapping into state: %w", err)
			}
		}
		// set eip map
		if gw.EipMap != nil {
			log.Printf("[TRACE] eip map: %#v", gw.EipMap)
			eipMap, err := setEipMapDetails(gw.EipMap, gw.IfNamesTranslation)
			if err != nil {
				return fmt.Errorf("could not set eip map details: %w", err)
			}
			if err = d.Set("eip_map", eipMap); err != nil {
				return fmt.Errorf("could not set eip map into state: %w", err)
			}
		}

		// set management egress ip prefix list
		if gw.ManagementEgressIPPrefix == "" {
			_ = d.Set("management_egress_ip_prefix_list", nil)
		} else {
			_ = d.Set("management_egress_ip_prefix_list", strings.Split(gw.ManagementEgressIPPrefix, ","))
		}

		// set the ha_attributes to empty if ha_interfaces is not provided for edge
		if len(gw.HaGw.Interfaces) == 0 {
			d.Set("ha_availability_domain", "")
			d.Set("ha_azure_eip_name_resource_group", "")
			d.Set("ha_cloud_instance_id", "")
			d.Set("ha_eip", "")
			d.Set("ha_fault_domain", "")
			d.Set("ha_gw_name", "")
			d.Set("ha_gw_size", "")
			d.Set("ha_image_version", "")
			d.Set("ha_insane_mode_az", "")
			d.Set("ha_lan_interface_cidr", "")
			d.Set("ha_oob_availability_zone", "")
			d.Set("ha_oob_management_subnet", "")
			d.Set("ha_private_ip", "")
			d.Set("ha_security_group_id", "")
			d.Set("ha_software_version", "")
			d.Set("ha_subnet", "")
			d.Set("ha_zone", "")
			d.Set("ha_public_ip", "")
			d.Set("ha_private_mode_subnet_zone", "")
			d.Set("ha_interfaces", nil)
			d.Set("ha_device_id", "")
			return nil
		}
		// set the backup link info details
		haGwName := gw.HaGw.GwName
		backupLinkInfo := gw.BackupLinkInfo
		if gw.HaGw.GwName == "" && backupLinkInfo != nil {
			backupLink, ok := backupLinkInfo[haGwName]
			if !ok {
				return fmt.Errorf("BackupLinkInfo does not contain HA Gateway: %s", haGwName)
			}
			peerBackupLogicalName := backupLink.PeerLogicalIntfName
			if err = d.Set("peer_backup_logical_ifname", peerBackupLogicalName); err != nil {
				return fmt.Errorf("could not set peer backup logical interface name into state: %w", err)
			}
			connectionTypePrivate := backupLink.ConnectionTypePublic
			if connectionTypePrivate {
				d.Set("peer_connection_type", "private")
			} else {
				d.Set("peer_connection_type", "public")
			}
		}
		// set ha interfaces
		if len(gw.HaGw.Interfaces) != 0 {
			ha_interfaces := setInterfaceDetails(gw.HaGw.Interfaces)
			if err = d.Set("ha_interfaces", ha_interfaces); err != nil {
				return fmt.Errorf("could not set interfaces into state: %w", err)
			}
		}
		// set device id
		if gw.HaGw.DeviceID != "" {
			d.Set("ha_device_id", gw.HaGw.DeviceID)
		}
		// set ha management egress ip prefix list
		if gw.HaGw.ManagementEgressIPPrefix == "" {
			_ = d.Set("ha_management_egress_ip_prefix_list", nil)
		} else {
			_ = d.Set("ha_management_egress_ip_prefix_list", strings.Split(gw.HaGw.ManagementEgressIPPrefix, ","))
		}
		if err := setGatewayResourceData(d, gw); err != nil {
			log.Printf("[ERROR] %v", err)
			return err
		}
	} else {
		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("eip", gw.PublicIP)
		d.Set("public_ip", gw.PublicIP)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)
		d.Set("single_ip_snat", gw.EnableNat == "yes" && gw.SnatMode == "primary")
		d.Set("single_az_ha", gw.SingleAZ == "yes")
		d.Set("enable_hybrid_connection", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) && gw.EnableHybridConnection)
		d.Set("connected_transit", gw.ConnectedTransit == "yes")
		d.Set("bgp_hold_time", gw.BgpHoldTime)
		d.Set("bgp_polling_time", gw.BgpPollingTime)
		d.Set("bgp_neighbor_status_polling_time", gw.BgpBfdPollingTime)
		d.Set("image_version", gw.ImageVersion)
		d.Set("software_version", gw.SoftwareVersion)
		d.Set("rx_queue_size", gw.RxQueueSize)
		d.Set("subnet", gw.VpcNet)

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
		d.Set("local_as_number", gw.LocalASNumber)
		d.Set("bgp_ecmp", gw.BgpEcmp)
		d.Set("enable_active_standby", gw.EnableActiveStandby)
		d.Set("enable_active_standby_preemptive", gw.EnableActiveStandbyPreemptive)
		d.Set("enable_s2c_rx_balancing", gw.EnableS2CRxBalancing)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan {
			d.Set("bgp_lan_interfaces_count", gw.BgpLanInterfacesCount)
		} else {
			d.Set("bgp_lan_interfaces_count", nil)
		}

		d.Set("enable_bgp_over_lan", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) && gw.EnableBgpOverLan {
			if len(gw.BgpLanInterfaces) != 0 {
				var interfaces []map[string]interface{}
				for _, bgpLanInterface := range gw.BgpLanInterfaces {
					interfaceDict := make(map[string]interface{})
					interfaceDict["vpc_id"] = bgpLanInterface.VpcID
					interfaceDict["subnet"] = bgpLanInterface.Subnet
					interfaces = append(interfaces, interfaceDict)
				}
				if err = d.Set("bgp_lan_interfaces", interfaces); err != nil {
					return fmt.Errorf("could not set bgp_lan_interfaces into state: %w", err)
				}
			}

			if len(gw.HaGw.HaBgpLanInterfaces) != 0 {
				var haInterfaces []map[string]interface{}
				for _, haBgpLanInterface := range gw.HaGw.HaBgpLanInterfaces {
					interfaceDict := make(map[string]interface{})
					interfaceDict["vpc_id"] = haBgpLanInterface.VpcID
					interfaceDict["subnet"] = haBgpLanInterface.Subnet
					haInterfaces = append(haInterfaces, interfaceDict)
				}
				if err = d.Set("ha_bgp_lan_interfaces", haInterfaces); err != nil {
					return fmt.Errorf("could not set ha_bgp_lan_interfaces into state: %w", err)
				}
			}

			bgpLanIpInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not get BGP LAN IP info for GCP transit gateway %s: %v", gateway.GwName, err)
			}
			if err = d.Set("bgp_lan_ip_list", bgpLanIpInfo.BgpLanIpList); err != nil {
				return fmt.Errorf("could not set bgp_lan_ip_list into state: %w", err)
			}
			if len(bgpLanIpInfo.HaBgpLanIpList) != 0 {
				if err = d.Set("ha_bgp_lan_ip_list", bgpLanIpInfo.HaBgpLanIpList); err != nil {
					return fmt.Errorf("could not set ha_bgp_lan_ip_list into state: %w", err)
				}
			} else {
				d.Set("ha_bgp_lan_ip_list", nil)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.EnableBgpOverLan {
			bgpLanIpInfo, err := client.GetBgpLanIPList(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not get BGP LAN IP info for Azure transit gateway %s: %v", gateway.GwName, err)
			}
			if err = d.Set("bgp_lan_ip_list", bgpLanIpInfo.AzureBgpLanIpList); err != nil {
				return fmt.Errorf("could not set bgp_lan_ip_list into state: %w", err)
			}
			if len(bgpLanIpInfo.AzureHaBgpLanIpList) != 0 {
				if err = d.Set("ha_bgp_lan_ip_list", bgpLanIpInfo.AzureHaBgpLanIpList); err != nil {
					return fmt.Errorf("could not set ha_bgp_lan_ip_list into state: %w", err)
				}
			} else {
				d.Set("ha_bgp_lan_ip_list", nil)
			}
		} else {
			d.Set("bgp_lan_ip_list", nil)
			d.Set("ha_bgp_lan_ip_list", nil)
		}
		d.Set("enable_transit_summarize_cidr_to_tgw", gw.EnableTransitSummarizeCidrToTgw)
		d.Set("enable_segmentation", gw.EnableSegmentation)
		d.Set("learned_cidrs_approval_mode", gw.LearnedCidrsApprovalMode)
		d.Set("enable_jumbo_frame", gw.JumboFrame)
		d.Set("enable_private_oob", gw.EnablePrivateOob)
		if gw.EnablePrivateOob {
			d.Set("oob_management_subnet", strings.Split(gw.OobManagementSubnet, "~~")[0])
			d.Set("oob_availability_zone", gw.GatewayZone)
		}
		d.Set("enable_firenet", gw.EnableFirenet)
		d.Set("enable_gateway_load_balancer", gw.EnableGatewayLoadBalancer)
		d.Set("enable_egress_transit_firenet", gw.EnableEgressTransitFirenet)
		d.Set("enable_preserve_as_path", gw.EnablePreserveAsPath)
		d.Set("customized_transit_vpc_routes", gw.CustomizedTransitVpcRoutes)
		d.Set("enable_transit_firenet", gw.EnableTransitFirenet)
		if gw.EnableTransitFirenet && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("lan_vpc_id", gw.BundleVpcInfo.LAN.VpcID)
			d.Set("lan_private_subnet", strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0])
		}

		if _, zoneIsSet := d.GetOk("zone"); goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || zoneIsSet) &&
			gw.GatewayZone != "AvailabilitySet" && gw.LbVpcId == "" {
			d.Set("zone", "az-"+gw.GatewayZone)
		}
		d.Set("enable_vpc_dns_server", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled")
		d.Set("enable_advertise_transit_cidr", gw.EnableAdvertiseTransitCidr)
		d.Set("enable_learned_cidrs_approval", gw.EnableLearnedCidrsApproval)
		if gw.EnableLearnedCidrsApproval {
			transitAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(&goaviatrix.TransitVpc{GwName: gw.GwName})
			if err != nil {
				return fmt.Errorf("could not get advanced config for transit gateway: %w", err)
			}

			if err = d.Set("approved_learned_cidrs", transitAdvancedConfig.ApprovedLearnedCidrs); err != nil {
				return fmt.Errorf("could not set approved_learned_cidrs into state: %w", err)
			}
		} else {
			d.Set("approved_learned_cidrs", nil)
		}
		d.Set("enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
		if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
			return fmt.Errorf("setting 'monitor_exclude_list' to state: %w", err)
		}
		d.Set("enable_multi_tier_transit", gw.EnableMultitierTransit)
		d.Set("tunnel_detection_time", gw.TunnelDetectionTime)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.ReuseEip, ":")
			if len(azureEip) == 3 {
				d.Set("azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the Transit Gateway %s", gw.GwName)
			}
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			if gw.AllocateNewEipRead && !gw.EnablePrivateOob {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			// gcp vpc_id returns as <vpc name>~-~<project name>
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.GatewayZone)
			d.Set("allocate_new_eip", gw.AllocateNewEipRead)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", gw.AllocateNewEipRead)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", gw.AllocateNewEipRead)
		} else if gw.CloudType == goaviatrix.AliCloud {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		}

		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				d.Set("insane_mode_az", gw.GatewayZone)
			} else {
				d.Set("insane_mode_az", "")
			}
		} else {
			d.Set("insane_mode", false)
			d.Set("insane_mode_az", "")
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			if customizedRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedRoutes != "" {
				customizedRoutesArray := strings.Split(customizedRoutes, ",")
				if len(goaviatrix.Difference(customizedRoutesArray, gw.CustomizedSpokeVpcRoutes)) == 0 &&
					len(goaviatrix.Difference(gw.CustomizedSpokeVpcRoutes, customizedRoutesArray)) == 0 {
					d.Set("customized_spoke_vpc_routes", customizedRoutes)
				} else {
					d.Set("customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
				}
			} else {
				d.Set("customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
			}
		} else {
			d.Set("customized_spoke_vpc_routes", "")
		}

		if len(gw.FilteredSpokeVpcRoutes) != 0 {
			if filteredSpokeVpcRoutes := d.Get("filtered_spoke_vpc_routes").(string); filteredSpokeVpcRoutes != "" {
				filteredSpokeVpcRoutesArray := strings.Split(filteredSpokeVpcRoutes, ",")
				if len(goaviatrix.Difference(filteredSpokeVpcRoutesArray, gw.FilteredSpokeVpcRoutes)) == 0 &&
					len(goaviatrix.Difference(gw.FilteredSpokeVpcRoutes, filteredSpokeVpcRoutesArray)) == 0 {
					d.Set("filtered_spoke_vpc_routes", filteredSpokeVpcRoutes)
				} else {
					d.Set("filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
				}
			} else {
				d.Set("filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
			}
		} else {
			d.Set("filtered_spoke_vpc_routes", "")
		}

		if len(gw.ExcludeCidrList) != 0 {
			if advertisedSpokeRoutes := d.Get("excluded_advertised_spoke_routes").(string); advertisedSpokeRoutes != "" {
				advertisedSpokeRoutesArray := strings.Split(advertisedSpokeRoutes, ",")
				if len(goaviatrix.Difference(advertisedSpokeRoutesArray, gw.ExcludeCidrList)) == 0 &&
					len(goaviatrix.Difference(gw.ExcludeCidrList, advertisedSpokeRoutesArray)) == 0 {
					d.Set("excluded_advertised_spoke_routes", advertisedSpokeRoutes)
				} else {
					d.Set("excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
				}
			} else {
				d.Set("excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
			}
		} else {
			d.Set("excluded_advertised_spoke_routes", "")
		}

		var bgpManualSpokeAdvertiseCidrs []string
		if _, ok := d.GetOk("bgp_manual_spoke_advertise_cidrs"); ok {
			bgpManualSpokeAdvertiseCidrs = strings.Split(d.Get("bgp_manual_spoke_advertise_cidrs").(string), ",")
		}
		if len(goaviatrix.Difference(bgpManualSpokeAdvertiseCidrs, gw.BgpManualSpokeAdvertiseCidrs)) != 0 ||
			len(goaviatrix.Difference(gw.BgpManualSpokeAdvertiseCidrs, bgpManualSpokeAdvertiseCidrs)) != 0 {
			bgpMSAN := ""
			for i := range gw.BgpManualSpokeAdvertiseCidrs {
				if i == 0 {
					bgpMSAN = bgpMSAN + gw.BgpManualSpokeAdvertiseCidrs[i]
				} else {
					bgpMSAN = bgpMSAN + "," + gw.BgpManualSpokeAdvertiseCidrs[i]
				}
			}
			d.Set("bgp_manual_spoke_advertise_cidrs", bgpMSAN)
		} else {
			d.Set("bgp_manual_spoke_advertise_cidrs", d.Get("bgp_manual_spoke_advertise_cidrs").(string))
		}

		lanCidr, err := client.GetTransitGatewayLanCidr(gw.GwName)
		if err != nil && err != goaviatrix.ErrNotFound {
			log.Printf("[WARN] Error getting lan cidr for transit gateway %s due to %s", gw.GwName, err)
		}
		d.Set("lan_interface_cidr", lanCidr)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			tags := goaviatrix.KeyValueTags(gw.Tags).IgnoreConfig(ignoreTagsConfig)
			if err := d.Set("tags", tags); err != nil {
				log.Printf("[WARN] Error setting tags for (%s): %s", d.Id(), err)
			}
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			if gw.GatewayZone != "" {
				d.Set("availability_domain", gw.GatewayZone)
			} else {
				d.Set("availability_domain", d.Get("availability_domain").(string))
			}
			d.Set("fault_domain", gw.FaultDomain)
		}

		if gw.EnableSpotInstance {
			d.Set("enable_spot_instance", true)
			d.Set("spot_price", gw.SpotPrice)
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && gw.DeleteSpot {
				d.Set("delete_spot", gw.DeleteSpot)
			}
		}

		d.Set("private_mode_lb_vpc_id", gw.LbVpcId)
		if gw.LbVpcId != "" && gw.GatewayZone != "AvailabilitySet" {
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				d.Set("private_mode_subnet_zone", gw.GatewayZone)
			} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				d.Set("private_mode_subnet_zone", "az-"+gw.GatewayZone)
			}
		} else {
			d.Set("private_mode_subnet_zone", nil)
		}

		enableGroGso, err := client.GetGroGsoStatus(gw)
		if err != nil {
			return fmt.Errorf("failed to get GRO/GSO status of transit gateway %s: %v", gw.GwName, err)
		}
		d.Set("enable_gro_gso", enableGroGso)
		if gw.HaGw.GwSize == "" {
			d.Set("ha_availability_domain", "")
			d.Set("ha_azure_eip_name_resource_group", "")
			d.Set("ha_cloud_instance_id", "")
			d.Set("ha_eip", "")
			d.Set("ha_fault_domain", "")
			d.Set("ha_gw_name", "")
			d.Set("ha_gw_size", "")
			d.Set("ha_image_version", "")
			d.Set("ha_insane_mode_az", "")
			d.Set("ha_lan_interface_cidr", "")
			d.Set("ha_oob_availability_zone", "")
			d.Set("ha_oob_management_subnet", "")
			d.Set("ha_private_ip", "")
			d.Set("ha_security_group_id", "")
			d.Set("ha_software_version", "")
			d.Set("ha_subnet", "")
			d.Set("ha_zone", "")
			d.Set("ha_public_ip", "")
			d.Set("ha_private_mode_subnet_zone", "")
			return nil
		}
		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			d.Set("ha_subnet", gw.HaGw.VpcNet)
			if zone := d.Get("ha_zone"); goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || zone.(string) != "") {
				if gw.LbVpcId == "" && gw.HaGw.GatewayZone != "AvailabilitySet" {
					d.Set("ha_zone", "az-"+gw.HaGw.GatewayZone)
				} else {
					d.Set("ha_zone", "")
				}
			} else {
				d.Set("ha_zone", "")
			}
		} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("ha_zone", gw.HaGw.GatewayZone)
			if d.Get("ha_subnet") != "" || isImport {
				d.Set("ha_subnet", gw.HaGw.VpcNet)
			}
		}

		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			if gw.HaGw.GatewayZone != "" {
				d.Set("ha_availability_domain", gw.HaGw.GatewayZone)
			} else {
				d.Set("ha_availability_domain", d.Get("ha_availability_domain").(string))
			}
			d.Set("ha_fault_domain", gw.HaGw.FaultDomain)
		}

		d.Set("ha_eip", gw.HaGw.PublicIP)
		d.Set("ha_public_ip", gw.HaGw.PublicIP)
		d.Set("ha_gw_size", gw.HaGw.GwSize)
		d.Set("ha_cloud_instance_id", gw.HaGw.CloudnGatewayInstID)
		d.Set("ha_gw_name", gw.HaGw.GwName)
		d.Set("ha_private_ip", gw.HaGw.PrivateIP)
		d.Set("ha_software_version", gw.HaGw.SoftwareVersion)
		d.Set("ha_image_version", gw.HaGw.ImageVersion)
		d.Set("ha_security_group_id", gw.HaGw.GwSecurityGroupID)
		lanCidr, err = client.GetTransitGatewayLanCidr(gw.HaGw.GwName)
		if err != nil && err != goaviatrix.ErrNotFound {
			log.Printf("[WARN] Error getting lan cidr for HA transit gateway %s due to %s", gw.HaGw.GwName, err)
		}
		d.Set("ha_lan_interface_cidr", lanCidr)

		if gw.HaGw.EnablePrivateOob {
			d.Set("ha_oob_management_subnet", strings.Split(gw.HaGw.OobManagementSubnet, "~~")[0])
			d.Set("ha_oob_availability_zone", gw.HaGw.GatewayZone)
		}

		if gw.LbVpcId != "" && gw.GatewayZone != "AvailabilitySet" {
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				d.Set("ha_private_mode_subnet_zone", gw.HaGw.GatewayZone)
			} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				d.Set("ha_private_mode_subnet_zone", "az-"+gw.HaGw.GatewayZone)
			}
		} else {
			d.Set("ha_private_mode_subnet_zone", "")
		}

		if gw.HaGw.InsaneMode == "yes" && goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			d.Set("ha_insane_mode_az", gw.HaGw.GatewayZone)
		} else {
			d.Set("ha_insane_mode_az", "")
		}

		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
			if len(azureEip) == 3 {
				d.Set("ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the HA Gateway %s", gw.GwName)
			}
		}
	}
	return nil
}

func resourceAviatrixTransitGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	haGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
		VpcSize:   d.Get("ha_gw_size").(string),
	}
	log.Printf("[INFO] Updating Aviatrix Transit Gateway: %#v", gateway)

	// Clarification : Can the user update EAT interface after its created. Add/Delete EAT interface
	d.Partial(true)
	if d.HasChange("ha_zone") {
		haZone := d.Get("ha_zone").(string)
		if haZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("'ha_zone' is only valid for GCP and Azure providers when enabling HA")
		}
	}
	if d.HasChange("ha_zone") || d.HasChange("ha_subnet") {
		haZone := d.Get("ha_zone").(string)
		haSubnet := d.Get("ha_subnet").(string)
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && haSubnet != "" && haZone == "" {
			return fmt.Errorf("'ha_zone' must be set to enable HA on GCP, cannot enable HA with only 'ha_subnet'")
		}
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" && haZone != "" {
			return fmt.Errorf("'ha_subnet' must be provided to enable HA on Azure, cannot enable HA with only 'ha_zone'")
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
		if o.(string) != "" && n.(string) != "" {
			return fmt.Errorf("updating ha_eip is not allowed")
		}
	}
	if d.HasChange("azure_eip_name_resource_group") {
		return fmt.Errorf("failed to update transit gateway: changing 'azure_eip_name_resource_group' is not allowed")
	}
	if d.HasChange("ha_azure_eip_name_resource_group") {
		o, n := d.GetChange("ha_azure_eip_name_resource_group")
		if o.(string) != "" && n.(string) != "" {
			return fmt.Errorf("failed to update transit gateway: changing 'ha_azure_eip_name_resource_group' is not allowed")
		}
	}
	if d.HasChange("lan_vpc_id") {
		return fmt.Errorf("updating lan_vpc_id is not allowed")
	}
	if d.HasChange("lan_private_subnet") {
		return fmt.Errorf("updating lan_private_subnet is not allowed")
	}

	// Transit FireNet function is not supported for AWS China and Azure China
	if d.HasChange("enable_firenet") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSChina|goaviatrix.AzureChina) {
		return fmt.Errorf("editing 'enable_firenet' in AWSChina (1024) and AzureChina (2048) is not supported")
	}
	// Transit FireNet function is not supported for Azure China
	if d.HasChange("enable_transit_firenet") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("editing 'enable_transit_firenet' in GCP (4), Azure (8), AzureGov (32) and AzureChina (2048) is not supported")
	}
	if d.Get("enable_egress_transit_firenet").(bool) && !d.Get("enable_transit_firenet").(bool) {
		return fmt.Errorf("'enable_egress_transit_firenet' requires 'enable_transit_firenet' to be set to true")
	}
	// Transit FireNet function is not supported for Azure China
	if d.Get("enable_egress_transit_firenet").(bool) && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes^goaviatrix.AzureChina|goaviatrix.OCIRelatedCloudTypes) {
		return fmt.Errorf("'enable_egress_transit_firenet' is currently only supported in AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
	learnedCidrsApprovedMode := d.Get("learned_cidrs_approval_mode").(string)
	approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")

	if learnedCidrsApproval && learnedCidrsApprovedMode == "connection" {
		return fmt.Errorf("'enable_learned_cidrs_approval' must be false if 'learned_cidrs_approval_mode' is set to 'connection'")
	}
	if !learnedCidrsApproval && len(approvedLearnedCidrs) != 0 {
		return fmt.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if d.HasChange("enable_private_oob") {
		return fmt.Errorf("updating enable_private_oob is not allowed")
	}

	enablePrivateOob := d.Get("enable_private_oob").(bool)
	privateModeInfo, _ := client.GetPrivateModeInfo(context.Background())

	if !enablePrivateOob {
		if d.HasChange("ha_oob_management_subnet") {
			return fmt.Errorf("updating ha_oob_manage_subnet is not allowed if private oob is disabled")
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

	var haBgpLanVpcID []string
	var haBgpLanSpecifySubnet []string
	for _, haBgpInterface := range d.Get("ha_bgp_lan_interfaces").([]interface{}) {
		item := haBgpInterface.(map[string]interface{})
		haBgpLanVpcID = append(haBgpLanVpcID, item["vpc_id"].(string))
		haBgpLanSpecifySubnet = append(haBgpLanSpecifySubnet, item["subnet"].(string))
	}
	if d.HasChange("ha_bgp_lan_interfaces") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			return fmt.Errorf("updating ha_bgp_lan_interfaces is not allowed for non-GCP Transit Gateway")
		}
		if !d.Get("enable_bgp_over_lan").(bool) {
			return fmt.Errorf("updating ha_bgp_lan_interfaces is not allowed for non-BGP over LAN GCP Transit Gateway")
		}
		if !d.HasChange("ha_subnet") {
			return fmt.Errorf("ha_bgp_lan_interfaces can only be updated along with BGP over LAN enabled GCP HA Transit Gateway's creation or deletion")
		}

		newHaSubnet := d.Get("ha_subnet").(string)
		if len(haBgpLanVpcID) == 0 && newHaSubnet != "" {
			return fmt.Errorf("missing ha_bgp_lan_interfaces for creating GCP HA transit gateway with BGP over LAN enabled in update")
		} else if len(haBgpLanVpcID) != 0 && newHaSubnet == "" {
			return fmt.Errorf("ha_bgp_lan_interfaces should be set empty for destroying GCP HA transit gateway with BGP over LAN enabled in update")
		}
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") ||
		(enablePrivateOob && (d.HasChange("ha_oob_management_subnet") || d.HasChange("ha_oob_availability_zone"))) ||
		(privateModeInfo.EnablePrivateMode && d.HasChange("ha_private_mode_subnet_zone")) ||
		d.HasChange("ha_availability_domain") || d.HasChange("ha_fault_domain") {
		transitHaGw := &goaviatrix.TransitHaGateway{
			PrimaryGwName: d.Get("gw_name").(string),
			GwName:        d.Get("gw_name").(string) + "-hagw",
			InsaneMode:    "no",
		}

		if d.Get("enable_bgp_over_lan").(bool) {
			// transitGw.BgpOverLan = "on"
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
				transitHaGw.BgpLanVpcID = strings.Join(haBgpLanVpcID, ",")
				transitHaGw.BgpLanSubnet = strings.Join(haBgpLanSpecifySubnet, ",")
			}
		}

		haEip := d.Get("ha_eip").(string)
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			transitHaGw.Eip = haEip
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		haGwSize := d.Get("ha_gw_size").(string)
		deleteHaGw := false
		changeHaGw := false
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			transitHaGw.Subnet = d.Get("ha_subnet").(string)
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && d.Get("ha_zone").(string) != "" {
				transitHaGw.Subnet = fmt.Sprintf("%s~~%s~~", d.Get("ha_subnet").(string), d.Get("ha_zone").(string))
			}

			haAvailabilityDomain := d.Get("ha_availability_domain").(string)
			haFaultDomain := d.Get("ha_fault_domain").(string)
			if newSubnet != "" {
				if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain == "" || haFaultDomain == "") {
					return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are required to enable HA on OCI")
				}
				if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain != "" || haFaultDomain != "") {
					return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are only valid for OCI")
				}
			}
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) {
				transitHaGw.Subnet = d.Get("ha_subnet").(string)
				transitHaGw.AvailabilityDomain = haAvailabilityDomain
				transitHaGw.FaultDomain = haFaultDomain
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
			transitHaGw.Zone = d.Get("ha_zone").(string)
			transitHaGw.Subnet = d.Get("ha_subnet").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}

		// TODO. Re-using EIP for HA GW on update needs test coverage
		haAzureEipName, haAzureEipNameOk := d.GetOk("ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if haEip != "" && (newHaGwEnabled || changeHaGw) {
				// No change will be detected when ha_eip is set to the empty string because it is computed.
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create HA Transit Gateway: 'ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				haAzureEipNameString, assertOk := haAzureEipName.(string)
				if !assertOk {
					return fmt.Errorf("failed to create HA Transit Gateway: 'ha_azure_eip_name_resource_group' must be a string")
				}
				transitHaGw.Eip = fmt.Sprintf("%s:%s", haAzureEipNameString, haEip)
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create HA Spoke Gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if d.Get("insane_mode").(bool) {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				var haStrs []string
				insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
				haSubnet := d.Get("ha_subnet").(string)

				if insaneModeHaAz == "" && haSubnet != "" {
					return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set for " +
						"AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
				} else if insaneModeHaAz != "" && haSubnet == "" {
					return fmt.Errorf("ha_subnet needed if insane_mode is enabled and ha_insane_mode_az is set for " +
						"AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
				}

				haStrs = append(haStrs, transitHaGw.Subnet, insaneModeHaAz)
				transitHaGw.Subnet = strings.Join(haStrs, "~~")
			}
			transitHaGw.InsaneMode = "yes"
		}

		if (newHaGwEnabled || changeHaGw) && haGwSize == "" {
			return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set")
		} else if deleteHaGw && haGwSize != "" {
			return fmt.Errorf("ha_gw_size must be empty if transit HA gateway is deleted")
		}

		haOobManagementSubnet := d.Get("ha_oob_management_subnet").(string)
		haOobAvailabilityZone := d.Get("ha_oob_availability_zone").(string)

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
					return fmt.Errorf("%q is required when creating a Transit HA Gateway in AWS if private mode is enabled and %q is provided", "ha_private_mode_subnet_zone", "ha_subnet")
				}

				privateModeSubnetZone := d.Get("ha_private_mode_subnet_zone").(string)
				transitHaGw.Subnet += "~~" + privateModeSubnetZone
			}
		}

		if newHaGwEnabled {
			_, err := client.CreateTransitHaGw(transitHaGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
			}
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				if d.Get("rx_queue_size").(string) != "" && !d.HasChange("rx_queue_size") {
					haGwRxQueueSize := &goaviatrix.Gateway{
						GwName:      d.Get("gw_name").(string) + "-hagw",
						RxQueueSize: d.Get("rx_queue_size").(string),
					}
					err := client.SetRxQueueSize(haGwRxQueueSize)
					if err != nil {
						return fmt.Errorf("could not set rx queue size for transit ha: %s during gateway update: %v", haGwRxQueueSize.GwName, err)
					}
				}
			}
		} else if deleteHaGw {
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Transit HA gateway: %s", err)
			}
		} else if changeHaGw {
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Transit HA gateway: %s", err)
			}

			transitHaGw.Eip = ""
			_, err = client.CreateTransitHaGw(transitHaGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
			}
			newHaGwEnabled = true
		}
	}
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	haEnabled := haSubnet != "" || haZone != ""
	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		singleAZ := d.Get("single_az_ha").(bool)
		if singleAZ {
			singleAZGateway.SingleAZ = "enabled"
		} else {
			singleAZGateway.SingleAZ = "disabled"
		}

		if singleAZ {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA for %s: %s", singleAZGateway.GwName, err)
			}

			if haEnabled {
				singleAZGatewayHA := &goaviatrix.Gateway{
					GwName: d.Get("gw_name").(string) + "-hagw",
				}
				err := client.EnableSingleAZGateway(singleAZGatewayHA)
				if err != nil {
					return fmt.Errorf("failed to enable single AZ GW HA for %s: %s", singleAZGatewayHA.GwName, err)
				}
			}
		} else {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA for %s: %s", singleAZGateway.GwName, err)
			}

			if haEnabled {
				singleAZGatewayHA := &goaviatrix.Gateway{
					GwName: d.Get("gw_name").(string) + "-hagw",
				}
				err := client.DisableSingleAZGateway(singleAZGatewayHA)
				if err != nil {
					return fmt.Errorf("failed to disable single AZ GW HA for %s: %s", singleAZGatewayHA.GwName, err)
				}
			}
		}
	}

	if d.HasChange("tags") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("failed to update transit gateway: adding tags is only supported for AWS (1), Azure (8), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			CloudType:    gateway.CloudType,
		}

		if d.HasChange("tags") {
			tagsMap, err := extractTags(d, gateway.CloudType)
			if err != nil {
				return fmt.Errorf("failed to update tags for transit gateway: %w", err)
			}
			tags.Tags = tagsMap
			tagJson, err := TagsMapToJson(tagsMap)
			if err != nil {
				return fmt.Errorf("failed to update tags for transit gateway: %w", err)
			}
			tags.TagJson = tagJson
			err = client.UpdateTags(tags)
			if err != nil {
				return fmt.Errorf("failed to update tags for transit gateway: %w", err)
			}
		}
	}

	if d.HasChange("connected_transit") {
		transitGateway := &goaviatrix.TransitVpc{
			CloudType:   d.Get("cloud_type").(int),
			AccountName: d.Get("account_name").(string),
			GwName:      d.Get("gw_name").(string),
			VpcID:       d.Get("vpc_id").(string),
			VpcRegion:   d.Get("vpc_reg").(string),
		}
		connectedTransit := d.Get("connected_transit").(bool)
		if connectedTransit {
			err := client.EnableConnectedTransit(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to enable connected transit: %s", err)
			}
		} else {
			err := client.DisableConnectedTransit(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to disable connected transit: %s", err)
			}
		}
	}

	if d.Get("enable_transit_firenet").(bool) {
		primaryGwSize := d.Get("gw_size").(string)
		if d.HasChange("gw_size") {
			old, _ := d.GetChange("gw_size")
			primaryGwSize = old.(string)
			gateway.VpcSize = d.Get("gw_size").(string)
			err := client.UpdateGateway(gateway)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit Gateway: %s", err)
			}
		}

		if d.HasChange("ha_gw_size") || newHaGwEnabled {
			newHaGwSize := d.Get("ha_gw_size").(string)
			if !newHaGwEnabled || (newHaGwSize != primaryGwSize) {
				// MODIFIES HA GW SIZE if
				// Ha gateway wasn't newly configured
				// OR
				// newly configured Ha gateway is set to be different size than primary gateway
				// (when ha gateway is enabled, it's size is by default the same as primary gateway)
				_, err := client.GetGateway(haGateway)
				if err != nil {
					// If HA gateway does not exist, don't try to change HA gateway size and continue with the rest of the updates
					// to the gateway
					if err != goaviatrix.ErrNotFound {
						return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway while trying to update HA Gw size: %s", err)
					}
				} else {
					if haGateway.VpcSize == "" {
						return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
							"ha_subnet or ha_zone is set")
					}
					err = client.UpdateGateway(haGateway)
					log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.VpcSize)
					if err != nil {
						return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
					}
				}
			}
		}
	}

	if d.HasChange("single_ip_snat") {
		gw := &goaviatrix.Gateway{
			CloudType:   d.Get("cloud_type").(int),
			GatewayName: d.Get("gw_name").(string),
		}
		enableNat := d.Get("single_ip_snat").(bool)

		if enableNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable 'single_ip' mode SNAT feature: %s", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable 'single_ip' mode SNAT: %s", err)
			}
		}

	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		// get wan count from interfaces
		interfaceList, ok := d.Get("interfaces").([]interface{})
		if !ok {
			return fmt.Errorf("invalid interfaces for EAT gateway")
		}
		wanCount, err := countInterfaceTypes(interfaceList)
		if err != nil {
			return fmt.Errorf("failed to get wan interface count: %w", err)
		}
		// get the cloud type
		cloudType, ok := d.Get("cloud_type").(int)
		if !ok {
			return fmt.Errorf("failed to get cloud type for EAT gateway")
		}

		// update the edge transit gateway interfaces
		if d.HasChanges("interfaces", "management_egress_ip_prefix_list") {
			interfaceList, ok := d.Get("interfaces").([]interface{})
			if !ok {
				return fmt.Errorf("invalid interfaces for EAT gateway")
			}
			interfaces, err := getInterfaceDetails(interfaceList, cloudType)
			if err != nil {
				return fmt.Errorf("failed to get interface details: %s", err)
			}

			gateway := &goaviatrix.TransitVpc{
				GwName:     d.Get("gw_name").(string),
				Interfaces: interfaces,
			}

			managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
			if len(managementEgressIPPrefixList) > 0 {
				gateway.ManagementEgressIPPrefix = strings.Join(managementEgressIPPrefixList, ",")
			}
			err = client.UpdateEdgeGateway(gateway)
			if err != nil {
				return fmt.Errorf("failed to update edge gateway: %s", err)
			}
		}

		// updating device id for edge transit gateway is not supported
		if d.HasChanges("device_id") {
			return fmt.Errorf("updating device_id is not supported for edge transit gateway")
		}

		// updating gw_size for edge transit gateway is not supported
		if d.HasChange("gw_size") {
			return fmt.Errorf("updating gw_size is not supported for edge transit gateway")
		}

		// update eip mapping for edge transit gateway
		if d.HasChange("eip_map") {
			eipMap, ok := d.Get("eip_map").([]interface{})
			if !ok {
				return fmt.Errorf("failed to get eip_map detail for Edge Transit Gateway")
			}
			if len(eipMap) > 0 {
				eipMapList, err := getEipMapDetails(eipMap, wanCount, cloudType)
				if err != nil {
					return fmt.Errorf("failed to get the eip map details: %w", err)
				}
				gateway := &goaviatrix.TransitVpc{
					GwName: d.Get("gw_name").(string),
				}

				if cloudType == goaviatrix.EDGEMEGAPORT {
					// print eip map for edge mega port
					log.Printf("[INFO] EIP Map for Edge Mega Port: %#v", eipMapList)
					gateway.LogicalEipMap = eipMapList
				} else {
					eipMapJSON, err := json.Marshal(eipMapList)
					if err != nil {
						return fmt.Errorf("failed to marshal eip_map to JSON: %w", err)
					}
					gateway.EipMap = string(eipMapJSON)
				}
				err = client.UpdateEdgeGateway(gateway)
				if err != nil {
					return fmt.Errorf("failed to update edge gateway: %s", err)
				}
			}
		}

		if d.HasChanges("peer_backup_port_type", "peer_backup_port_index", "peer_connection_type") {
			// if transitHaGateway already exists then cannot update the backup link info
			haGateway := &goaviatrix.Gateway{
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string) + "-hagw",
			}
			resultHaGw, err := client.GetGateway(haGateway)
			if err == nil && resultHaGw != nil {
				return fmt.Errorf("cannot update the backup link info for edge transit gateway when HA gateway already exists")
			}
			cloudType := d.Get("cloud_type").(int)
			// create transit ha gw with backup link info
			transitHaGw, err := getTransitHaGatewayDetails(d, wanCount, cloudType)
			if err != nil {
				return fmt.Errorf("failed to get transit ha gateway details: %w", err)
			}
			log.Printf("[INFO] Enabling HA on Transit Gateway")
			_, err = client.CreateTransitHaGw(transitHaGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
			}
		}

		// Update ha interfaces for EAT gateway
		if d.HasChanges("ha_interfaces", "ha_management_egress_ip_prefix_list") {
			haInterfaceList, ok := d.Get("ha_interfaces").([]interface{})
			if !ok {
				return fmt.Errorf("invalid ha_interfaces for EAT HA gateway")
			}
			if len(haInterfaceList) > 0 {
				haInterfaces, err := getInterfaceDetails(haInterfaceList, cloudType)
				if err != nil {
					return fmt.Errorf("failed to get interface details: %s", err)
				}
				gateway := &goaviatrix.TransitVpc{
					GwName:     d.Get("gw_name").(string) + "-hagw",
					Interfaces: haInterfaces,
				}
				haManagementEgressIPPrefixList := getStringSet(d, "ha_management_egress_ip_prefix_list")
				if len(haManagementEgressIPPrefixList) > 0 {
					gateway.ManagementEgressIPPrefix = strings.Join(haManagementEgressIPPrefixList, ",")
				}
				err = client.UpdateEdgeGateway(gateway)
				if err != nil {
					return fmt.Errorf("failed to update edge gateway: %s", err)
				}
			} else {
				// delete the HA gateway if ha_interfaces is empty
				err := client.DeleteGateway(haGateway)
				if err != nil {
					return fmt.Errorf("failed to delete HA gateway: %s", err)
				}
			}
		}
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		if d.HasChange("enable_hybrid_connection") {
			transitGateway := &goaviatrix.TransitVpc{
				CloudType:   d.Get("cloud_type").(int),
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string),
				VpcID:       d.Get("vpc_id").(string),
				VpcRegion:   d.Get("vpc_reg").(string),
			}
			enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
			if enableHybridConnection {
				err := client.AttachTransitGWForHybrid(transitGateway)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for Hybrid: %s", err)
				}
			} else {
				err := client.DetachTransitGWForHybrid(transitGateway)
				if err != nil {
					return fmt.Errorf("failed to disable transit GW for Hybrid: %s", err)
				}
			}
		}
	} else {
		if d.HasChange("enable_hybrid_connection") {
			return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS/AWSGov providers")
		}
	}

	if d.HasChange("learned_cidrs_approval_mode") && d.HasChange("enable_learned_cidrs_approval") {
		gw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		currentMode, _ := d.GetChange("learned_cidrs_approval_mode")
		// API calls need to be in a specific order depending on the current mode
		if currentMode.(string) == "gateway" {
			if learnedCidrsApproval {
				err := client.EnableTransitLearnedCidrsApproval(gw)
				if err != nil {
					return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
				}
			} else {
				err := client.DisableTransitLearnedCidrsApproval(gw)
				if err != nil {
					return fmt.Errorf("failed to disable learned cidrs approval: %s", err)
				}
			}
			mode := d.Get("learned_cidrs_approval_mode").(string)
			err := client.SetTransitLearnedCIDRsApprovalMode(gw, mode)
			if err != nil {
				return fmt.Errorf("could not set learned CIDRs approval mode to %q: %v", mode, err)
			}
		} else {
			mode := d.Get("learned_cidrs_approval_mode").(string)
			err := client.SetTransitLearnedCIDRsApprovalMode(gw, mode)
			if err != nil {
				return fmt.Errorf("could not set learned CIDRs approval mode to %q: %v", mode, err)
			}
			if learnedCidrsApproval {
				err = client.EnableTransitLearnedCidrsApproval(gw)
				if err != nil {
					return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
				}
			} else {
				err = client.DisableTransitLearnedCidrsApproval(gw)
				if err != nil {
					return fmt.Errorf("failed to disable learned cidrs approval: %s", err)
				}
			}
		}
	} else if d.HasChange("learned_cidrs_approval_mode") {
		gw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		mode := d.Get("learned_cidrs_approval_mode").(string)
		err := client.SetTransitLearnedCIDRsApprovalMode(gw, mode)
		if err != nil {
			return fmt.Errorf("could not set learned CIDRs approval mode to %q: %v", mode, err)
		}
	} else if d.HasChange("enable_learned_cidrs_approval") {
		gw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		if learnedCidrsApproval {
			gw.LearnedCidrsApproval = "yes"
			err := client.EnableTransitLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
			}
		} else {
			gw.LearnedCidrsApproval = "no"
			err := client.DisableTransitLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to disable learned cidrs approval: %s", err)
			}
		}
	}

	if learnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gw := &goaviatrix.TransitVpc{
			GwName:               d.Get("gw_name").(string),
			ApprovedLearnedCidrs: approvedLearnedCidrs,
		}

		err := client.UpdateTransitPendingApprovedCidrs(gw)
		if err != nil {
			return fmt.Errorf("could not update approved CIDRs: %w", err)
		}
	}

	enableFireNet := d.Get("enable_firenet").(bool)
	enableGatewayLoadBalancer := d.Get("enable_gateway_load_balancer").(bool)
	enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
	if enableGatewayLoadBalancer && !enableFireNet && !enableTransitFireNet {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
	}
	if enableGatewayLoadBalancer && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWS) {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'cloud_type' = 1 (AWS)")
	}
	if enableFireNet && enableTransitFireNet {
		return fmt.Errorf("can't enable firenet function and transit firenet function at the same time")
	}

	if d.HasChange("enable_egress_transit_firenet") {
		enableEgressTransitFirenet := d.Get("enable_egress_transit_firenet").(bool)
		if !enableEgressTransitFirenet {
			err := client.DisableEgressTransitFirenet(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not disable egress transit firenet: %w", err)
			}
		}
	}

	if d.HasChanges("enable_preserve_as_path") {
		enableTransitPreserveAsPath := d.Get("enable_preserve_as_path").(bool)
		if !enableTransitPreserveAsPath {
			err := client.DisableTransitPreserveAsPath(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not disable Preserve AS Path during Transit Gateway update: %w", err)
			}
		} else {
			err := client.EnableTransitPreserveAsPath(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not enable Preserve AS Path during Transit Gateway update: %w", err)
			}
		}
	}

	if d.HasChange("enable_firenet") && d.HasChange("enable_transit_firenet") {
		transitGW := &goaviatrix.TransitVpc{
			GwName: gateway.GwName,
			VpcID:  d.Get("vpc_id").(string),
		}
		if !enableFireNet {
			err := client.DisableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to disable transit GW for FireNet Interfaces: %s", err)
			}
		}
		if !enableTransitFireNet {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.DisableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to disable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		}
		if enableFireNet {
			if enableGatewayLoadBalancer {
				err := client.EnableGatewayFireNetInterfacesWithGWLB(transitGW)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for FireNet Interfaces with Gateway Load Balancer enabled: %s", err)
				}
			} else {
				err := client.EnableGatewayFireNetInterfaces(transitGW)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
				}
			}
		}
		if enableTransitFireNet {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			if enableGatewayLoadBalancer {
				err := client.EnableTransitFireNetWithGWLB(gwTransitFireNet)
				if err != nil {
					return fmt.Errorf("failed to enable transit firenet with Gateway Load Balancer for %s due to %s", gwTransitFireNet.GwName, err)
				}
			} else {
				err := client.EnableTransitFireNet(gwTransitFireNet)
				if err != nil {
					return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
				}
			}
		}
	} else if d.HasChange("enable_firenet") {
		transitGW := &goaviatrix.TransitVpc{
			GwName: gateway.GwName,
			VpcID:  d.Get("vpc_id").(string),
		}
		if enableFireNet {
			if enableGatewayLoadBalancer {
				err := client.EnableGatewayFireNetInterfacesWithGWLB(transitGW)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for FireNet Interfaces with Gateway Load Balancer enabled: %s", err)
				}
			} else {
				err := client.EnableGatewayFireNetInterfaces(transitGW)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
				}
			}
		} else {
			err := client.DisableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to disable transit GW for FireNet Interfaces: %s", err)
			}
		}
	} else if d.HasChange("enable_transit_firenet") {
		if enableTransitFireNet {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			if enableGatewayLoadBalancer {
				err := client.EnableTransitFireNetWithGWLB(gwTransitFireNet)
				if err != nil {
					return fmt.Errorf("failed to enable transit firenet with Gateway Load Balancer for %s due to %s", gwTransitFireNet.GwName, err)
				}
			} else {
				err := client.EnableTransitFireNet(gwTransitFireNet)
				if err != nil {
					return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
				}
			}
		} else {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.DisableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to disable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		}
	} else if d.HasChange("enable_gateway_load_balancer") {
		// In this branch we know that neither 'enable_transit_firenet' or 'enable_firenet' HasChange.
		// Due to the backend design it is not possible to disable or enable 'enable_gateway_load_balancer' without
		// also disabling or enabling FireNet, so we force the user to disable or enable both at the same time.
		if enableGatewayLoadBalancer {
			return fmt.Errorf("can not enable 'enable_gateway_load_balancer' when 'enable_firenet' or 'enable_transit_firenet' is " +
				"already enabled. Changing from non-GWLB FireNet to GWLB FireNet requires 2 separate " +
				"`terraform apply` steps, once to disable non-GWLB FireNet, then again to enable GWLB FireNet")
		} else {
			return fmt.Errorf("can not disable 'enable_gateway_load_balancer' when 'enable_firenet' or 'enable_transit_firenet' is " +
				"still enabled. Changing from GWLB FireNet to non-GWLB FireNet requires 2 separate " +
				"`terraform apply` steps, once to disable GWLB FireNet, then again to enable non-GWLB FireNet")
		}
	}

	if !d.Get("enable_transit_firenet").(bool) && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		primaryGwSize := d.Get("gw_size").(string)
		if d.HasChange("gw_size") {
			old, _ := d.GetChange("gw_size")
			primaryGwSize = old.(string)
			gateway.VpcSize = d.Get("gw_size").(string)
			err := client.UpdateGateway(gateway)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit Gateway: %s", err)
			}
		}

		if d.HasChange("ha_gw_size") || newHaGwEnabled {
			newHaGwSize := d.Get("ha_gw_size").(string)
			if !newHaGwEnabled || (newHaGwSize != primaryGwSize) {
				// MODIFIES HA GW SIZE if
				// Ha gateway wasn't newly configured
				// OR
				// newly configured Ha gateway is set to be different size than primary gateway
				// (when ha gateway is enabled, it's size is by default the same as primary gateway)
				_, err := client.GetGateway(haGateway)
				if err != nil {
					// If HA gateway does not exist, don't try to change gateway size and continue with the rest of the updates
					// to the gateway
					if err != goaviatrix.ErrNotFound {
						return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway while trying to update HA Gw size: %s", err)
					}
				} else {
					if haGateway.VpcSize == "" {
						return fmt.Errorf("a valid non empty ha_gw_size parameter is mandatory for this resource if " +
							"ha_subnet or ha_zone is set")
					}
					err = client.UpdateGateway(haGateway)
					log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.VpcSize)
					if err != nil {
						return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
					}
				}
			}
		}
	}

	if d.HasChange("enable_egress_transit_firenet") {
		enableEgressTransitFirenet := d.Get("enable_egress_transit_firenet").(bool)
		if enableEgressTransitFirenet {
			err := client.EnableEgressTransitFirenet(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not enable egress transit firenet: %w", err)
			}
		}
	}

	if d.HasChange("enable_vpc_dns_server") && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gw := &goaviatrix.Gateway{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}

		enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
		if enableVpcDnsServer {
			err := client.EnableVpcDnsServer(gw)
			if err != nil {
				return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
			}
		} else {
			err := client.DisableVpcDnsServer(gw)
			if err != nil {
				return fmt.Errorf("failed to disable VPC DNS Server: %s", err)
			}
		}
	} else if d.HasChange("enable_vpc_dns_server") {
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if d.HasChange("enable_advertise_transit_cidr") {
		transitGw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		enableAdvertiseTransitCidr := d.Get("enable_advertise_transit_cidr").(bool)
		if enableAdvertiseTransitCidr {
			transitGw.EnableAdvertiseTransitCidr = true
			err := client.EnableAdvertiseTransitCidr(transitGw)
			if err != nil {
				return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
			}
		} else {
			transitGw.EnableAdvertiseTransitCidr = false
			err := client.DisableAdvertiseTransitCidr(transitGw)
			if err != nil {
				return fmt.Errorf("failed to disable advertise transit CIDR: %s", err)
			}
		}
	}

	if d.HasChange("bgp_manual_spoke_advertise_cidrs") {
		transitGw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		bgpManualSpokeAdvertiseCidrs := d.Get("bgp_manual_spoke_advertise_cidrs").(string)
		transitGw.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
		err := client.SetBgpManualSpokeAdvertisedNetworks(transitGw)
		if err != nil {
			return fmt.Errorf("failed to set bgp manual spoke advertise CIDRs: %s", err)
		}
	}

	if d.HasChange("enable_encrypt_volume") {
		if d.Get("enable_encrypt_volume").(bool) {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_encrypt_volume' is only supported by AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              d.Get("gw_name").(string),
				CustomerManagedKeys: d.Get("customer_managed_keys").(string),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
			}

			haSubnet := d.Get("ha_subnet").(string)
			haZone := d.Get("ha_zone").(string)
			haEnabled := haSubnet != "" || haZone != ""
			if haEnabled {
				gwHAEncVolume := &goaviatrix.Gateway{
					GwName:              d.Get("gw_name").(string) + "-hagw",
					CustomerManagedKeys: d.Get("customer_managed_keys").(string),
				}
				err := client.EnableEncryptVolume(gwHAEncVolume)
				if err != nil {
					return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwHAEncVolume.GwName, err)
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
		oldRouteList := strings.Split(o.(string), ",")
		newRouteList := strings.Split(n.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                   d.Get("gw_name").(string),
				CustomizedSpokeVpcRoutes: newRouteList,
			}
			err := client.EditGatewayCustomRoutes(transitGateway)
			log.Printf("[INFO] Customizeing routes of transit gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to customize spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if d.HasChange("filtered_spoke_vpc_routes") {
		o, n := d.GetChange("filtered_spoke_vpc_routes")
		oldRouteList := strings.Split(o.(string), ",")
		newRouteList := strings.Split(n.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                 d.Get("gw_name").(string),
				FilteredSpokeVpcRoutes: newRouteList,
			}
			err := client.EditGatewayFilterRoutes(transitGateway)
			log.Printf("[INFO] Editing filtered spoke vpc routes of transit gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if d.HasChange("excluded_advertised_spoke_routes") {
		o, n := d.GetChange("excluded_advertised_spoke_routes")
		oldRouteList := strings.Split(o.(string), ",")
		newRouteList := strings.Split(n.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                d.Get("gw_name").(string),
				AdvertisedSpokeRoutes: newRouteList,
			}
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			log.Printf("[INFO] Editing excluded advertised spoke vpc routes of transit gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit excluded advertised spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		bgpPollingTime := d.Get("bgp_polling_time").(int)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpPollingTime(gateway, bgpPollingTime)
		if err != nil {
			return fmt.Errorf("could not update bgp polling time: %w", err)
		}
	}

	if d.HasChange("bgp_neighbor_status_polling_time") {
		bgpBfdPollingTime := d.Get("bgp_neighbor_status_polling_time").(int)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpBfdPollingTime(gateway, bgpBfdPollingTime)
		if err != nil {
			return fmt.Errorf("could not update bgp neighbor status polling time: %w", err)
		}
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}

		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(prependASPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gateway, nil)
			if err != nil {
				return fmt.Errorf("could not delete prepend_as_path during Transit Gateway update: %w", err)
			}
		}

		if d.HasChange("local_as_number") {
			localAsNumber := d.Get("local_as_number").(string)
			err := client.SetLocalASNumber(gateway, localAsNumber)
			if err != nil {
				return fmt.Errorf("could not set local_as_number during Transit Gateway update: %w", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path during Transit Gateway update: %w", err)
			}
		}
	}

	if d.HasChange("bgp_ecmp") {
		enabled := d.Get("bgp_ecmp").(bool)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpEcmp(gateway, enabled)
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp: %w", err)
		}
	}

	if d.HasChange("enable_segmentation") {
		enabled := d.Get("enable_segmentation").(bool)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		if enabled {
			if err := client.EnableSegmentation(gateway); err != nil {
				return fmt.Errorf("could not enable segmentation: %w", err)
			}
		} else {
			if err := client.DisableSegmentation(gateway); err != nil {
				return fmt.Errorf("could not disable segmentation: %w", err)
			}
		}
	}

	if d.HasChange("enable_active_standby") || d.HasChange("enable_active_standby_preemptive") {
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		if d.Get("enable_active_standby").(bool) {
			if d.Get("enable_active_standby_preemptive").(bool) {
				if err := client.EnableActiveStandbyPreemptive(gateway); err != nil {
					return fmt.Errorf("could not enable Preemptive Mode for Active-Standby during Transit Gateway update: %w", err)
				}
			} else {
				if err := client.EnableActiveStandby(gateway); err != nil {
					return fmt.Errorf("could not enable Active-Standby during Transit Gateway update: %w", err)
				}
			}
		} else {
			if d.Get("enable_active_standby_preemptive").(bool) {
				return fmt.Errorf("could not enable Preemptive Mode with Active-Standby disabled")
			}
			if err := client.DisableActiveStandby(gateway); err != nil {
				return fmt.Errorf("could not disable Active-Standby during Transit Gateway update: %w", err)
			}
		}
	}

	if d.HasChange("customized_transit_vpc_routes") {
		var customizedTransitVpcRoutes []string
		for _, v := range d.Get("customized_transit_vpc_routes").(*schema.Set).List() {
			customizedTransitVpcRoutes = append(customizedTransitVpcRoutes, v.(string))
		}

		err := client.UpdateTransitGatewayCustomizedVpcRoute(gateway.GwName, customizedTransitVpcRoutes)
		if err != nil {
			return fmt.Errorf("couldn't update transit gateway customized vpc route: %s", err)
		}
	}

	monitorGatewaySubnets := d.Get("enable_monitor_gateway_subnets").(bool)
	var excludedInstances []string
	for _, v := range d.Get("monitor_exclude_list").(*schema.Set).List() {
		excludedInstances = append(excludedInstances, v.(string))
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
		if d.Get("enable_jumbo_frame").(bool) {
			err := client.EnableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not enable jumbo frame for transit gateway when updating: %w", err)
			}
		} else {
			err := client.DisableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not disable jumbo frame for transit gateway when updating: %w", err)
			}
		}
	}

	if d.HasChange("enable_gro_gso") {
		if d.Get("enable_gro_gso").(bool) {
			err := client.EnableGroGso(gateway)
			if err != nil {
				return fmt.Errorf("couldn't enable GRO/GSO on transit gateway when updating: %s", err)
			}
		} else {
			err := client.DisableGroGso(gateway)
			if err != nil {
				return fmt.Errorf("couldn't disable GRO/GSO on transit gateway when updating: %s", err)
			}
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(gateway.GwName, d.Get("bgp_hold_time").(int))
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time during Transit Gateway update: %w", err)
		}
	}

	if d.HasChange("enable_transit_summarize_cidr_to_tgw") {
		if d.Get("enable_transit_summarize_cidr_to_tgw").(bool) {
			err := client.EnableSummarizeCidrToTgw(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable summarize cidr to tgw when updating: %w", err)
			}
		} else {
			err := client.DisableSummarizeCidrToTgw(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not disable summarize cidr to tgw when updating: %w", err)
			}
		}
	}

	if d.HasChange("enable_multi_tier_transit") {
		if d.Get("enable_multi_tier_transit").(bool) {
			if d.Get("local_as_number") == "" {
				return fmt.Errorf("local_as_number required to enable multi tier transit")
			}
			err := client.EnableMultitierTransit(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable multi tier transit when updating: %w", err)
			}
		} else {
			err := client.DisableMultitierTransit(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not disable multi tier transit when updating: %w", err)
			}
		}
	}

	if d.HasChange("enable_s2c_rx_balancing") {
		if d.Get("enable_s2c_rx_balancing").(bool) {
			err := client.EnableS2CRxBalancing(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable S2C receive packet CPU re-balancing during Transit Gateway update: %w", err)
			}
		} else {
			err := client.DisableS2CRxBalancing(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not disable S2C receive packet CPU re-balancing during Transit Gateway update: %w", err)
			}
		}
	}

	if d.HasChange("tunnel_detection_time") {
		detectionTimeInterface, ok := d.GetOk("tunnel_detection_time")
		var detectionTime int
		if ok {
			detectionTime = detectionTimeInterface.(int)
		} else {
			var err error
			detectionTime, err = client.GetTunnelDetectionTime("Controller")
			if err != nil {
				return fmt.Errorf("could not get default tunnel detection time during Transit Gateway update: %w", err)
			}
		}
		err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime)
		if err != nil {
			return fmt.Errorf("could not modify tunnel detection time during Transit Gateway update: %w", err)
		}
	}

	if d.HasChange("rx_queue_size") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("could not update rx_queue_size since it only supports AWS related cloud types")
		}
		gw := &goaviatrix.Gateway{
			GwName:      gateway.GwName,
			RxQueueSize: d.Get("rx_queue_size").(string),
		}
		err := client.SetRxQueueSize(gw)
		if err != nil {
			return fmt.Errorf("could not modify rx queue size for transit: %s during gateway update: %v", gw.GatewayName, err)
		}
		if haSubnet != "" || haZone != "" {
			haGwRxQueueSize := &goaviatrix.Gateway{
				GwName:      d.Get("gw_name").(string) + "-hagw",
				RxQueueSize: d.Get("rx_queue_size").(string),
			}
			err := client.SetRxQueueSize(haGwRxQueueSize)
			if err != nil {
				return fmt.Errorf("could not modify rx queue size for transit ha: %s during gateway update: %v", haGwRxQueueSize.GwName, err)
			}
		}
	}

	if d.HasChanges("enable_bgp_over_lan", "bgp_lan_interfaces_count") {
		if d.HasChange("enable_bgp_over_lan") {
			if !d.Get("enable_bgp_over_lan").(bool) {
				return fmt.Errorf("disabling BGP over LAN during update is not supported for transit: %s", gateway.GwName)
			}
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				return fmt.Errorf("enabling BGP over LAN during update is only supported for Azure transit")
			}
			if _, ok := d.GetOk("bgp_lan_interfaces_count"); !ok {
				return fmt.Errorf("please specify 'bgp_lan_interfaces_count' to enable BGP over LAN during update for Azure transit: %s", gateway.GwName)
			}
		}
		if d.HasChange("bgp_lan_interfaces_count") {
			if !d.Get("enable_bgp_over_lan").(bool) || !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				return fmt.Errorf("could not update BGP over LAN interface count since it only supports BGP over LAN enabled transit for Azure (8), AzureGov (32) or AzureChina (2048)")
			}
			oldCount, newCount := d.GetChange("bgp_lan_interfaces_count")
			if oldCount.(int) > newCount.(int) {
				return fmt.Errorf("deleting BGP over LAN interface during update is not supported for transit: %s", gateway.GwName)
			}
		}
		gw := &goaviatrix.Gateway{
			GwName:                gateway.GwName,
			BgpLanInterfacesCount: d.Get("bgp_lan_interfaces_count").(int),
		}
		err := client.ChangeBgpOverLanIntfCnt(gw)
		if err != nil {
			return fmt.Errorf("could not modify BGP over LAN interface count for transit: %s during gateway update: %v", gw.GatewayName, err)
		}
		if haSubnet != "" {
			haGw := &goaviatrix.Gateway{
				GwName:                d.Get("gw_name").(string) + "-hagw",
				BgpLanInterfacesCount: d.Get("bgp_lan_interfaces_count").(int),
			}
			err := client.ChangeBgpOverLanIntfCnt(haGw)
			if err != nil {
				return fmt.Errorf("could not modify BGP over LAN interface count for transit ha: %s during gateway update: %v", haGw.GwName, err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixTransitGatewayRead(d, meta)
}

func resourceAviatrixTransitGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	cloudType := d.Get("cloud_type").(int)

	gateway := &goaviatrix.Gateway{
		CloudType: cloudType,
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Gateway: %#v", gateway)

	enableEgressTransitFirenet := d.Get("enable_egress_transit_firenet").(bool)
	if enableEgressTransitFirenet {
		err := client.DisableEgressTransitFirenet(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return fmt.Errorf("could not disable egress transit firenet: %w", err)
		}
	}

	enableTransitPreserveAsPath := d.Get("enable_preserve_as_path").(bool)
	if enableTransitPreserveAsPath {
		err := client.DisableTransitPreserveAsPath(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return fmt.Errorf("could not disable transit preserve as path: %w", err)
		}
	}

	enableFireNet := d.Get("enable_firenet").(bool)
	if enableFireNet {
		gw := &goaviatrix.TransitVpc{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}

		err := client.DisableGatewayFireNetInterfaces(gw)
		if err != nil {
			return fmt.Errorf("failed to disable Aviatrix Transit Gateway for FireNet Interfaces: %s", err)
		}
	}

	enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
	if enableTransitFireNet && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		err := client.DisableTransitFireNet(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable transit firenet for %s due to %s", gateway.GwName, err)
		}
	} else if enableTransitFireNet && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		err := client.IsTransitFireNetReadyToBeDisabled(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable transit firenet for %s due to %s", gateway.GwName, err)
		}
	}

	// If HA is enabled, delete HA GW first.
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	ha_interfaces, _ := d.Get("ha_interfaces").([]interface{})

	if haSubnet != "" || haZone != "" || (goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) && len(ha_interfaces) > 0) {
		gateway.GwName += "-hagw"

		try, maxTries, backoff := 0, 2, 500*time.Millisecond

		for {
			try++
			err := client.DeleteGateway(gateway)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					break
				}

				if !strings.Contains(err.Error(), "EOF") {
					return fmt.Errorf("failed to delete Aviatrix Edge Transit Gateway HA gateway: %s", err)
				}
			} else {
				break
			}

			if try == maxTries {
				return fmt.Errorf("failed to delete Aviatrix Edge Transit Gateway HA gateway: %s", err)
			}
			time.Sleep(backoff)
			// Double the backoff time after each failed try
			backoff *= 2
		}
		if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT) {
			vpcID, ok := d.Get("vpc_id").(string)
			if !ok {
				return fmt.Errorf("vpc_id is not a string")
			}
			ztpFileDownloadPath, ok := d.Get("ztp_file_download_path").(string)
			if !ok {
				return fmt.Errorf("ztp_file_download_path is not a string")
			}
			err := deleteZtpFile(gateway.GwName, vpcID, ztpFileDownloadPath)
			if err != nil {
				return fmt.Errorf("failed to delete ZTP file: %w", err)
			}
		}
	}
	gateway.GwName = d.Get("gw_name").(string)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Edge Transit Gateway: %s", err)
	}
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT) {
		vpcID, ok := d.Get("vpc_id").(string)
		if !ok {
			return fmt.Errorf("vpc_id is not a string")
		}
		ztpFileDownloadPath, ok := d.Get("ztp_file_download_path").(string)
		if !ok {
			return fmt.Errorf("ztp_file_download_path is not a string")
		}
		err := deleteZtpFile(gateway.GwName, vpcID, ztpFileDownloadPath)
		if err != nil {
			return fmt.Errorf("failed to delete ZTP file: %w", err)
		}
	}

	return nil
}

func deleteZtpFile(gatewayName, vpcID, ztpFileDownloadPath string) error {
	fileName := ztpFileDownloadPath + "/" + gatewayName + "-" + vpcID + "-cloud-init.txt"
	if err := os.Remove(fileName); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("could not remove the ztp file: %w", err)
	}
	return nil
}

func createEdgeTransitGateway(d *schema.ResourceData, client *goaviatrix.Client, cloudType int) error {
	gateway := &goaviatrix.TransitVpc{
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
		VpcID:       d.Get("vpc_id").(string),
		VpcSize:     d.Get("gw_size").(string),
		Transit:     true,
	}
	// get the interface config details
	interfaces, ok := d.Get("interfaces").([]interface{})
	if !ok || len(interfaces) == 0 {
		return fmt.Errorf("interfaces attribute is required for Edge Transit Gateway")
	}
	interfacesList, err := getInterfaceDetails(interfaces, cloudType)
	if err != nil {
		return fmt.Errorf("failed to get the interface details: %w", err)
	}
	gateway.Interfaces = interfacesList
	// get the count of WAN inetrfaces to map interface type and index to interface name
	wanCount, err := countInterfaceTypes(interfaces)
	if err != nil {
		return fmt.Errorf("failed to get the interface count: %w", err)
	}
	// interface mapping and device_id are required only for AEP edge gateway
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGENEO) {
		/*
			TODO: Use the device_id to determine the interface mapping. This change will provide support for other device models interface mapping. For now, we will use the user provided interface mapping for ESXI devices and default values for Dell devices.
		*/
		interfaceMappingInput := d.Get("interface_mapping").([]interface{})
		interfaceMapping, err := getInterfaceMappingDetails(interfaceMappingInput)
		if err != nil {
			return fmt.Errorf("failed to get the interface mapping details: %w", err)
		}
		gateway.InterfaceMapping = interfaceMapping
		// set the device_id for AEP gateway
		gateway.DeviceID, ok = d.Get("device_id").(string)
		if !ok {
			return fmt.Errorf("device_id attribute is required for Edge Transit Gateway")
		}
	}

	// ztp file download path is required for Equinix and Megaport edge gateways
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGEMEGAPORT) {
		gateway.ZtpFileDownloadPath, ok = d.Get("ztp_file_download_path").(string)
		if !ok {
			return fmt.Errorf("ztp_file_download_path attribute is required for Edge Transit Gateway")
		}
		managementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
		if len(managementEgressIPPrefixList) > 0 {
			gateway.ManagementEgressIPPrefix = strings.Join(managementEgressIPPrefixList, ",")
		}
	}

	// create the transit gateway
	log.Printf("[INFO] Creating Aviatrix Transit Gateway: %#v", gateway)
	d.SetId(gateway.GwName)
	err = client.LaunchTransitVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit Gateway: %w", err)
	}
	// create ha transit gateway if ha_interfaces are provided
	haInterfaces, ok := d.Get("ha_interfaces").([]interface{})
	if ok && len(haInterfaces) > 0 {
		transitHaGw, err := getTransitHaGatewayDetails(d, wanCount, cloudType)
		if err != nil {
			return fmt.Errorf("failed to get the HA gateway details: %w", err)
		}
		// log transit ha gateway details
		log.Printf("[INFO] Creating HA Aviatrix Transit Gateway: %#v", transitHaGw)
		_, err = client.CreateTransitHaGw(transitHaGw)
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %w", err)
		}
	}

	// for AEP EAT gateway, set the eip_map, bgp polling time, bgp neighbor status polling time, local_as_number and prepend_as_path after the transit is created. AEP gateways take ~15 mins to be up and running. Updates to the gateway while is in the process of being created will fail.
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGENEO) {
		return nil
	}

	// eip map is updated after the transit is created
	eipMap, ok := d.Get("eip_map").([]interface{})
	if !ok {
		return fmt.Errorf("failed to get eip_map detail for Edge Transit Gateway")
	}
	if len(eipMap) > 0 {
		eipMapList, err := getEipMapDetails(eipMap, wanCount, cloudType)
		if err != nil {
			return fmt.Errorf("failed to get the eip map details: %w", err)
		}

		if cloudType == goaviatrix.EDGEMEGAPORT {
			// print eip map for edge mega port
			log.Printf("[INFO] EIP Map for Edge Mega Port: %#v", eipMapList)
			gateway.LogicalEipMap = eipMapList
		} else {
			eipMapJSON, err := json.Marshal(eipMapList)
			if err != nil {
				return fmt.Errorf("failed to marshal eip_map to JSON: %w", err)
			}
			gateway.EipMap = string(eipMapJSON)
		}
		// update EIP map
		err = client.UpdateEdgeGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update edge gateway: %s", err)
		}
	}

	if val, ok := d.GetOk("bgp_polling_time"); ok {
		err := client.SetBgpPollingTime(gateway, val.(int))
		if err != nil {
			return fmt.Errorf("could not set bgp polling time: %w", err)
		}
	}

	if val, ok := d.GetOk("bgp_neighbor_status_polling_time"); ok {
		err := client.SetBgpBfdPollingTime(gateway, val.(int))
		if err != nil {
			return fmt.Errorf("could not set bgp neighbor status polling time: %w", err)
		}
	}

	if val, ok := d.GetOk("local_as_number"); ok {
		err := client.SetLocalASNumber(gateway, val.(string))
		if err != nil {
			return fmt.Errorf("could not set local_as_number: %w", err)
		}
	}

	if val, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		slice := val.([]interface{})
		for _, v := range slice {
			prependASPath = append(prependASPath, v.(string))
		}
		err := client.SetPrependASPath(gateway, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %w", err)
		}
	}
	return nil
}

// count the number of WAN interfaces in the interface config
func countInterfaceTypes(interfaceList []interface{}) (int, error) {
	wanCount := 0
	// Iterate over each interface in the list
	for _, item := range interfaceList {
		interfaceMap, ok := item.(map[string]interface{})
		if !ok {
			return 0, fmt.Errorf("invalid interface entry, expected a map")
		}
		logicalIfName, ok := interfaceMap["logical_ifname"].(string)
		if !ok {
			return 0, fmt.Errorf("logical interface name must be a string")
		}
		// Count WAN interfaces; check if logicalIfName starts with "wan"
		if logicalIfName != "" && strings.HasPrefix(logicalIfName, "wan") {
			wanCount++
		}
	}
	return wanCount, nil
}

// Extract the interface type and index from the logical interface name
func extractInterfaceTypeAndIndex(logicalIfName string) (string, int, error) {
	// Regular expression to match "wan0", "mgmt1", etc.
	re := regexp.MustCompile(`^(wan|mgmt)(\d+)$`)
	matches := re.FindStringSubmatch(logicalIfName)

	if matches == nil {
		return "", 0, fmt.Errorf("invalid logical interface name: %s", logicalIfName)
	}

	// Map the interface type to its corresponding full form
	typeMap := map[string]string{
		"wan":  "WAN",
		"mgmt": "MANAGEMENT",
	}

	// Extract type and index from regex match
	intfType, ok := typeMap[matches[1]]
	if !ok {
		return "", 0, fmt.Errorf("unsupported interface type: %s", matches[1])
	}
	var intfIndex int
	_, err := fmt.Sscanf(matches[2], "%d", &intfIndex)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse index from: %s", logicalIfName)
	}

	return intfType, intfIndex, nil
}

func isValidInterfaceName(name string) bool {
	return name != "" && (strings.HasPrefix(name, "wan") || strings.HasPrefix(name, "mgmt"))
}

func calculateInterfaceName(intfType string, intfIndex, wanCount int) (string, error) {
	var num int

	switch intfType {
	case "WAN":
		num = calculateWANInterfaceNumber(intfIndex)
	case "MANAGEMENT":
		num = calculateMgmtInterfaceNumber(intfIndex, wanCount)
	default:
		return "", errors.New("unexpected interface type")
	}

	interfaceName := fmt.Sprintf("eth%d", num)
	log.Printf("Mapping interface %s%d to port %s", intfType, intfIndex, interfaceName)
	return interfaceName, nil
}

func calculateWANInterfaceNumber(intfIndex int) int {
	if intfIndex == 0 {
		return 0
	}
	if intfIndex == 1 {
		return 1
	}
	return intfIndex + 1
}

func calculateMgmtInterfaceNumber(intfIndex, wanCount int) int {
	if intfIndex == 0 {
		return 2
	}
	return wanCount + intfIndex
}

// get the interface name (ethX) from logical interface name and counts
func getInterfaceName(logicalIfName string, wanCount int) (string, error) {
	// check if the logical interface name is valid
	if !isValidInterfaceName(logicalIfName) {
		return "", fmt.Errorf("invalid logical interface name: %s", logicalIfName)
	}

	intfType, intfIndex, err := extractInterfaceTypeAndIndex(logicalIfName)
	if err != nil {
		return "", fmt.Errorf("failed to extract interface type and index: %w", err)
	}

	interfaceName, err := calculateInterfaceName(intfType, intfIndex, wanCount)
	if err != nil {
		return "", err
	}
	return interfaceName, nil
}

func getEipMapDetails(eipMap []interface{}, wanCount int, cloudType int) (map[string][]goaviatrix.EipMap, error) {
	// Create a map to structure the EIP data
	eipMapStructured := make(map[string][]goaviatrix.EipMap)

	// Populate the structured map
	for _, eipDetails := range eipMap {
		eip, ok := eipDetails.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type: expected map[string]interface{}, got %T", eipDetails)
		}
		logicalIfName, ok := eip["logical_ifname"].(string)
		if !ok {
			return nil, fmt.Errorf("logical interface name must be a string")
		}

		var interfaceName string
		if cloudType == goaviatrix.EDGEMEGAPORT {
			interfaceName = logicalIfName
		} else {
			var err error
			interfaceName, err = getInterfaceName(logicalIfName, wanCount)
			if err != nil {
				return nil, fmt.Errorf("failed to get the interface name using type and index for eip_map: %w", err)
			}
		}

		privateIP, ok := eip["private_ip"].(string)
		if !ok {
			return nil, fmt.Errorf("private_ip must be a string")
		}
		publicIP, ok := eip["public_ip"].(string)
		if !ok {
			return nil, fmt.Errorf("public_ip must be a string")
		}

		// Append the EIP entry to the corresponding interface
		eipEntry := goaviatrix.EipMap{
			PrivateIP: privateIP,
			PublicIP:  publicIP,
		}
		eipMapStructured[interfaceName] = append(eipMapStructured[interfaceName], eipEntry)
	}

	return eipMapStructured, nil
}

func getTransitHaGatewayDetails(d *schema.ResourceData, wanCount int, cloudType int) (*goaviatrix.TransitHaGateway, error) {
	gwName, ok := d.Get("gw_name").(string)
	if !ok {
		return nil, fmt.Errorf("gw_name is required for HA Edge Transit Gateway")
	}
	haGwName := gwName + "-hagw"
	transitHaGw := &goaviatrix.TransitHaGateway{
		PrimaryGwName:       gwName,
		GwName:              haGwName,
		VpcID:               d.Get("vpc_id").(string),
		ZtpFileDownloadPath: d.Get("ztp_file_download_path").(string),
		CloudType:           d.Get("cloud_type").(int),
		InsaneMode:          "yes",
	}
	peerBackupLogicalName, ok := d.Get("peer_backup_logical_ifname").([]interface{})
	if !ok {
		return nil, fmt.Errorf("peer backup logical interface name is required for HA Edge Transit Gateway")
	}
	connectionType, ok := d.Get("peer_connection_type").(string)
	if !ok {
		return nil, fmt.Errorf("peer_connection_type is required for HA Edge Transit Gateway")
	}
	// get the backup link configuration
	backupLinkConfig, err := createBackupLinkConfig(gwName, peerBackupLogicalName, connectionType, wanCount, cloudType)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup link configuration: %s", err)
	}
	transitHaGw.BackupLinkConfig = b64.StdEncoding.EncodeToString([]byte(backupLinkConfig))
	// get the HA interfaces
	haInterfaces, ok := d.Get("ha_interfaces").([]interface{})
	if !ok {
		return nil, fmt.Errorf("ha_interfaces is not a list")
	}
	interfacesList, err := getInterfaceDetails(haInterfaces, cloudType)
	if err != nil {
		return nil, fmt.Errorf("failed to get the interface details for HA Edge Transit Gateway: %w", err)
	}
	transitHaGw.Interfaces = interfacesList
	// interface mapping and device id is only required for AEP ha gateway
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGENEO) {
		// interface mapping is same as the primary gateway
		interfaceMappingInput := d.Get("interface_mapping").([]interface{})
		interfaceMapping, err := getInterfaceMappingDetails(interfaceMappingInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get the interface mapping details: %w", err)
		}
		transitHaGw.InterfaceMapping = interfaceMapping
		// set device_id for AEP ha gateway
		transitHaGw.DeviceID, ok = d.Get("ha_device_id").(string)
		if !ok {
			return nil, fmt.Errorf("ha_device_id is required for AEP HA Edge Transit Gateway")
		}
	}
	haManagementEgressIPPrefixList := getStringSet(d, "ha_management_egress_ip_prefix_list")
	if len(haManagementEgressIPPrefixList) > 0 {
		transitHaGw.ManagementEgressIPPrefix = strings.Join(haManagementEgressIPPrefixList, ",")
	}
	return transitHaGw, nil
}

func createBackupLinkConfig(gwName string, peerBackupLogicalNames []interface{}, connectionType string, wanCount int, cloudType int) (string, error) {
	var peerBackupPorts []string
	for _, logicalName := range peerBackupLogicalNames {
		// convert the logical interface name to string
		logicalName, ok := logicalName.(string)
		if !ok {
			return "", fmt.Errorf("peer backup logical interface name is not a string")
		}
		portName, err := getInterfaceName(logicalName, wanCount) // returns eth0, eth1, etc.
		if err != nil {
			return "", fmt.Errorf("failed to get the peer backup port name for logical name %s: %w", logicalName, err)
		}
		peerBackupPorts = append(peerBackupPorts, portName)
	}
	backupLinkData := goaviatrix.BackupLinkInterface{
		PeerGwName:     gwName,
		ConnectionType: connectionType,
	}

	if cloudType == goaviatrix.EDGEMEGAPORT {
		// convert the peer backup logical interface name []interface{} to []string
		var peerBackupLogicalPort []string
		if len(peerBackupLogicalNames) != 0 {
			for _, peerPort := range peerBackupLogicalNames {
				peerPortStr, ok := peerPort.(string)
				if !ok {
					return "", fmt.Errorf("secondary_private_cidr_list contains non-string elements")
				}
				peerBackupLogicalPort = append(peerBackupLogicalPort, peerPortStr)
			}
		}
		backupLinkData.PeerBackupLogicalIfNames = peerBackupLogicalPort
		backupLinkData.SelfBackupLogicalIfNames = peerBackupLogicalPort
	} else {
		peerBackupPort := strings.Join(peerBackupPorts, ",")
		backupLinkData.PeerBackupPort = peerBackupPort
		backupLinkData.SelfBackupPort = peerBackupPort
	}

	backupLinkList := []goaviatrix.BackupLinkInterface{backupLinkData}
	backupLinkConfig, err := json.Marshal(backupLinkList)
	if err != nil {
		return "", fmt.Errorf("failed to marshal backup link configuration: %w", err)
	}

	return string(backupLinkConfig), nil
}

// set the interface mapping details for the gateway
func getInterfaceMappingDetails(interfaceMappingInput []interface{}) (string, error) {
	interfaceMapping := map[string][]string{}
	if len(interfaceMappingInput) > 0 {
		// Set the interface mapping for ESXI devices
		for _, value := range interfaceMappingInput {
			mappingMap, ok := value.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("invalid type %T for interface mapping, expected a map", value)
			}
			interfaceName, ok1 := mappingMap["name"].(string)
			interfaceType, ok2 := mappingMap["type"].(string)
			interfaceIndex, ok3 := mappingMap["index"].(int)
			if !ok1 || !ok2 || !ok3 {
				return "", fmt.Errorf("invalid interface mapping, 'name', 'type', and 'index' must be strings")
			}
			updatedInterfaceType, ok := typeMapping[interfaceType]
			if !ok {
				return "", fmt.Errorf("invalid interface type %s", interfaceType)
			}
			interfaceMapping[interfaceName] = []string{updatedInterfaceType, strconv.Itoa(interfaceIndex)}
		}
	} else {
		// Set the interface mapping for Dell devices
		interfaceMapping["eth0"] = []string{"mgmt", "0"}
		interfaceMapping["eth5"] = []string{"wan", "0"}
		interfaceMapping["eth2"] = []string{"wan", "1"}
		interfaceMapping["eth3"] = []string{"wan", "2"}
		interfaceMapping["eth4"] = []string{"wan", "3"}
	}
	// Convert interfaceMapping to JSON byte slice
	interfaceMappingJSON, err := json.Marshal(interfaceMapping)
	if err != nil {
		return "", fmt.Errorf("failed to marshal interface mapping to json: %w", err)
	}
	return string(interfaceMappingJSON), nil
}

// get the interface details from the interface config and cloud type
func getInterfaceDetails(interfaces []interface{}, cloudType int) (string, error) {
	// get the count of WAN and MANAGEMENT interfaces
	wanCount, err := countInterfaceTypes(interfaces)
	if err != nil {
		return "", fmt.Errorf("failed to get the wan interface count: %w", err)
	}
	interfaceList := []goaviatrix.EdgeTransitInterface{}
	for _, iface := range interfaces {
		ifaceInfo, ok := iface.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("interface is not a map[string]interface{}")
		}

		ifaceData, err := parseInterface(ifaceInfo, wanCount, cloudType)
		if err != nil {
			return "", fmt.Errorf("failed to parse interface: %w", err)
		}

		interfaceList = append(interfaceList, ifaceData)
	}
	interfaceListJson, err := json.Marshal(interfaceList)
	if err != nil {
		return "", err
	}
	interfacesEncoded := b64.StdEncoding.EncodeToString(interfaceListJson)
	return interfacesEncoded, nil
}

func parseInterface(ifaceInfo map[string]interface{}, wanCount, cloudType int) (goaviatrix.EdgeTransitInterface, error) {
	var (
		logicalIfName, ifaceName, ifaceType, ifaceGatewayIP, ifaceIP, ifacePublicIP string
		ifaceDHCP                                                                   bool
		secondaryCIDRs                                                              []string
	)

	logicalIfName, err := getStringAttribute(ifaceInfo, "logical_ifname")
	if err != nil {
		return goaviatrix.EdgeTransitInterface{}, fmt.Errorf("logical_ifname: %w", err)
	}

	ifaceName, err = getInterfaceName(logicalIfName, wanCount)
	if err != nil {
		return goaviatrix.EdgeTransitInterface{}, fmt.Errorf("failed to get the interface name: %w", err)
	}

	ifaceType, _, err = extractInterfaceTypeAndIndex(logicalIfName)
	if err != nil {
		return goaviatrix.EdgeTransitInterface{}, fmt.Errorf("failed to extract interface type and index: %w", err)
	}

	ifaceGatewayIP, _ = getStringAttribute(ifaceInfo, "gateway_ip")
	ifaceIP, _ = getStringAttribute(ifaceInfo, "ip_address")
	ifacePublicIP, _ = getStringAttribute(ifaceInfo, "public_ip")
	ifaceDHCP, _ = getBoolAttribute(ifaceInfo, "dhcp")
	secondaryCIDRs, _ = getStringListAttribute(ifaceInfo, "secondary_private_cidr_list")

	ifaceData := goaviatrix.EdgeTransitInterface{
		GatewayIp:      ifaceGatewayIP,
		PublicIp:       ifacePublicIP,
		Dhcp:           ifaceDHCP,
		IpAddress:      ifaceIP,
		SecondaryCIDRs: secondaryCIDRs,
	}

	if cloudType == goaviatrix.EDGEMEGAPORT {
		ifaceData.LogicalIfName = logicalIfName
	} else {
		ifaceData.Name = ifaceName
		ifaceData.Type = ifaceType
	}

	return ifaceData, nil
}

func getStringAttribute(data map[string]interface{}, key string) (string, error) {
	val, exists := data[key]
	if !exists || val == nil {
		return "", nil
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string", key)
	}
	return str, nil
}

func getBoolAttribute(data map[string]interface{}, key string) (bool, error) {
	val, exists := data[key]
	if !exists || val == nil {
		return false, nil
	}
	boolean, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("%s is not a bool", key)
	}
	return boolean, nil
}

func getStringListAttribute(data map[string]interface{}, key string) ([]string, error) {
	val, exists := data[key]
	if !exists || val == nil {
		return nil, nil
	}
	list, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s is not a valid list", key)
	}

	result := make([]string, len(list))
	for i, item := range list {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("%s contains non-string elements", key)
		}
		result[i] = str
	}

	return result, nil
}

func setInterfaceDetails(interfaces []goaviatrix.EdgeTransitInterface) []map[string]interface{} {
	interfaceList := make([]map[string]interface{}, 0)
	sortedInterfaces := sortInterfacesByCustomOrder(interfaces)
	for _, intf := range sortedInterfaces {
		interfaceDict := make(map[string]interface{})
		interfaceDict["logical_ifname"] = intf.LogicalIfName

		if intf.PublicIp != "" {
			interfaceDict["public_ip"] = intf.PublicIp
		}
		if intf.Dhcp {
			interfaceDict["dhcp"] = intf.Dhcp
		}
		if intf.IpAddress != "" {
			interfaceDict["ip_address"] = intf.IpAddress
		}
		if intf.GatewayIp != "" {
			interfaceDict["gateway_ip"] = intf.GatewayIp
		}
		if intf.SecondaryCIDRs != nil {
			secondaryCIDRs := make([]string, 0)
			for _, cidr := range intf.SecondaryCIDRs {
				if cidr != "" {
					secondaryCIDRs = append(secondaryCIDRs, cidr)
				}
			}
			interfaceDict["secondary_private_cidr_list"] = secondaryCIDRs
		}
		interfaceList = append(interfaceList, interfaceDict)
	}
	return interfaceList
}

func setInterfaceMappingDetails(interfaceMapping []goaviatrix.InterfaceMapping) []map[string]interface{} {
	interfaceMappingList := make([]map[string]interface{}, 0)
	// sort the interface mapping by interface name. This will maintain the order of the interface mappings in the state
	sortedInterfaceMappings := sortInterfaceMappingByCustomOrder(interfaceMapping)
	for _, intf := range sortedInterfaceMappings {
		interfaceMap := make(map[string]interface{})
		interfaceMap["name"] = intf.Name
		if intf.Type == "mgmt" {
			interfaceMap["type"] = "MANAGEMENT"
		} else {
			interfaceMap["type"] = strings.ToUpper(intf.Type)
		}
		interfaceMap["index"] = intf.Index
		interfaceMappingList = append(interfaceMappingList, interfaceMap)
	}
	return interfaceMappingList
}

func setEipMapDetails(eipMap map[string][]goaviatrix.EipMap, ifNameTranslation map[string]string) ([]map[string]interface{}, error) {
	eipMapList := make([]map[string]interface{}, 0)
	keys := make([]string, 0, len(eipMap))
	for key := range eipMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Process the eipMap in sorted key order
	for _, interfaceName := range keys {
		eipList := eipMap[interfaceName]
		// Get the interface type and index from the interface name
		interfaces := ifNameTranslation[interfaceName]
		interfaceDetails := strings.Split(interfaces, ".")
		if len(interfaceDetails) != 2 {
			return nil, fmt.Errorf("invalid interface format for %s: %s", interfaceName, interfaces)
		}
		interfaceType := strings.ToLower(interfaceDetails[0])
		interfaceIndex, err := strconv.Atoi(interfaceDetails[1])
		if err != nil {
			return nil, fmt.Errorf("failed to convert interface index to integer for %d: %w", interfaceIndex, err)
		}
		// Create the logical interface name using the type and index
		logicalIfName := fmt.Sprintf("%s%d", interfaceType, interfaceIndex)

		for _, eip := range eipList {
			eipMap := map[string]interface{}{
				"logical_ifname": logicalIfName,
				"private_ip":     eip.PrivateIP,
				"public_ip":      eip.PublicIP,
			}
			eipMapList = append(eipMapList, eipMap)
		}
	}
	return eipMapList, nil
}

func setGatewayResourceData(d *schema.ResourceData, gw *goaviatrix.Gateway) error {
	settings := map[string]interface{}{
		"eip":                                  gw.PublicIP,
		"cloud_instance_id":                    gw.CloudnGatewayInstID,
		"security_group_id":                    gw.GwSecurityGroupID,
		"single_ip_snat":                       gw.EnableNat == "yes" && gw.SnatMode == "primary",
		"connected_transit":                    gw.ConnectedTransit == "yes",
		"bgp_hold_time":                        gw.BgpHoldTime,
		"bgp_polling_time":                     gw.BgpPollingTime,
		"bgp_neighbor_status_polling_time":     gw.BgpBfdPollingTime,
		"image_version":                        gw.ImageVersion,
		"software_version":                     gw.SoftwareVersion,
		"local_as_number":                      gw.LocalASNumber,
		"bgp_ecmp":                             gw.BgpEcmp,
		"enable_active_standby":                gw.EnableActiveStandby,
		"enable_active_standby_preemptive":     gw.EnableActiveStandbyPreemptive,
		"enable_s2c_rx_balancing":              gw.EnableS2CRxBalancing,
		"enable_transit_summarize_cidr_to_tgw": gw.EnableTransitSummarizeCidrToTgw,
		"enable_segmentation":                  gw.EnableSegmentation,
		"learned_cidrs_approval_mode":          gw.LearnedCidrsApprovalMode,
		"enable_jumbo_frame":                   gw.JumboFrame,
		"enable_private_oob":                   gw.EnablePrivateOob,
		"enable_firenet":                       gw.EnableFirenet,
		"enable_gateway_load_balancer":         gw.EnableGatewayLoadBalancer,
		"enable_egress_transit_firenet":        gw.EnableEgressTransitFirenet,
		"enable_preserve_as_path":              gw.EnablePreserveAsPath,
		"customized_transit_vpc_routes":        gw.CustomizedTransitVpcRoutes,
		"enable_transit_firenet":               gw.EnableTransitFirenet,
		"enable_advertise_transit_cidr":        gw.EnableAdvertiseTransitCidr,
		"enable_learned_cidrs_approval":        gw.EnableLearnedCidrsApproval,
		"enable_multi_tier_transit":            gw.EnableMultitierTransit,
		"tunnel_detection_time":                gw.TunnelDetectionTime,
	}

	for key, value := range settings {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}

	var prependAsPath []string
	for _, p := range strings.Split(gw.PrependASPath, " ") {
		if p != "" {
			prependAsPath = append(prependAsPath, p)
		}
	}

	if err := d.Set("prepend_as_path", prependAsPath); err != nil {
		return fmt.Errorf("failed to set prepend_as_path: %w", err)
	}

	return nil
}
