package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSpokeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeGatewayCreate,
		Read:   resourceAviatrixSpokeGatewayRead,
		Update: resourceAviatrixSpokeGatewayUpdate,
		Delete: resourceAviatrixSpokeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixSpokeGatewayResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixSpokeGatewayStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Type of cloud service provider.",
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureAZ,
				Description:  "Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
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
					return d.Get("enable_private_oob").(bool) || d.Get("private_mode_lb_vpc_id").(string) != ""
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
				Default:     false,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"transit_gw": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Deprecated:  "Please set `manage_transit_gateway_attachment` to false, and use the standalone aviatrix_spoke_transit_attachment resource instead.",
				Description: "Specify the transit Gateways to attach to this spoke. Format is a comma-separated list of transit gateway names. For example, 'transit-gw1,transit-gw2'.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldGws := strings.Split(old, ",")
					newGws := strings.Split(new, ",")
					return goaviatrix.Equivalent(oldGws, newGws)
				},
			},
			"manage_transit_gateway_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "This parameter is a switch used to determine whether or not to manage attaching this spoke gateway to transit gateways " +
					"using the aviatrix_spoke_gateway resource. If this is set to false, attaching this spoke gateway to " +
					"transit gateways must be done using the aviatrix_spoke_transit_attachment resource. " +
					"Valid values: true, false. Default value: true.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Deprecated:  "Use tags instead.",
				Description: "Instance tag of cloud provider.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Customer managed key ID.",
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
			"tags": {
				Type:          schema.TypeMap,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Description:   "A map of tags to assign to the spoke gateway.",
				ConflictsWith: []string{"tag_list"},
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
					ValidateFunc: validation.IsCIDR,
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
			"bgp_hold_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultBgpHoldTime,
				ValidateFunc: validation.IntBetween(12, 360),
				Description:  "BGP Hold Time for BGP Spoke Gateway. Unit is in seconds. Valid values are between 12 and 360.",
			},
			"enable_spot_instance": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "Enable spot instance. NOT supported for production deployment.",
				RequiredWith: []string{"spot_price"},
			},
			"spot_price": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Price for spot instance. NOT supported for production deployment.",
				RequiredWith: []string{"enable_spot_instance"},
			},
			"rx_queue_size": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1K", "2K", "4K", "8K", "16K"}, false),
				Description:  "Gateway ethernet interface RX queue size. Supported for AWS related clouds only.",
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
		},
	}
}

func resourceAviatrixSpokeGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.SpokeVpc{
		CloudType:            d.Get("cloud_type").(int),
		AccountName:          d.Get("account_name").(string),
		GwName:               d.Get("gw_name").(string),
		VpcSize:              d.Get("gw_size").(string),
		Subnet:               d.Get("subnet").(string),
		HASubnet:             d.Get("ha_subnet").(string),
		AvailabilityDomain:   d.Get("availability_domain").(string),
		FaultDomain:          d.Get("fault_domain").(string),
		ApprovedLearnedCidrs: getStringSet(d, "approved_learned_cidrs"),
	}

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)

	if d.Get("enable_private_vpc_default_route").(bool) && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_private_vpc_default_route is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if d.Get("enable_skip_public_route_table_update").(bool) && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_skip_public_route_update is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if _, hasSetZone := d.GetOk("zone"); !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.Azure) && hasSetZone {
		return fmt.Errorf("attribute 'zone' is only valid for Azure (8)")
	}

	if _, hasSetZone := d.GetOk("zone"); goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && hasSetZone {
		gateway.Subnet = fmt.Sprintf("%s~~%s~~", d.Get("subnet").(string), d.Get("zone").(string))
	}

	enableSNat := d.Get("single_ip_snat").(bool)
	if enableSNat {
		gateway.EnableNat = "yes"
	}

	singleAZ := d.Get("single_az_ha").(bool)
	if singleAZ {
		gateway.SingleAzHa = "enabled"
	} else {
		gateway.SingleAzHa = "disabled"
	}

	enableBgp := d.Get("enable_bgp").(bool)
	disableRoutePropagation := d.Get("disable_route_propagation").(bool)
	if enableBgp {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWS|goaviatrix.Azure) {
			return fmt.Errorf("enabling BGP is only supported for AWS (1) and Azure (8)")
		}
		gateway.EnableBgp = "on"
	} else {
		if disableRoutePropagation {
			return fmt.Errorf("disable route propagation is not supported on Non-BGP Spoke")
		}
	}

	learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if !learnedCidrsApproval && len(gateway.ApprovedLearnedCidrs) != 0 {
		return fmt.Errorf("'approved_learned_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gateway.VpcID = d.Get("vpc_id").(string)
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		gateway.VNetNameResourceGroup = d.Get("vpc_id").(string)
		if gateway.VNetNameResourceGroup == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gateway.VpcRegion = d.Get("vpc_reg").(string)
	} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = d.Get("vpc_reg").(string)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192) or, AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain == "" || gateway.FaultDomain == "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
	}
	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain != "" || gateway.FaultDomain != "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
	}

	insaneMode := d.Get("insane_mode").(bool)
	insaneModeAz := d.Get("insane_mode_az").(string)
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	haAvailabilityDomain := d.Get("ha_availability_domain").(string)
	haFaultDomain := d.Get("ha_fault_domain").(string)

	if haZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("'ha_zone' is only valid for GCP (4), Azure (8), AzureGov (32) and AzureChina (2048) providers if enabling HA")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && haSubnet != "" && haZone == "" {
		return fmt.Errorf("'ha_zone' must be set to enable HA on GCP (4), cannot enable HA with only 'ha_subnet'")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" && haZone != "" {
		return fmt.Errorf("'ha_subnet' must be provided to enable HA on Azure (4), AzureGov (32) or AzureChina (2048), cannot enable HA with only 'ha_zone'")
	}
	haGwSize := d.Get("ha_gw_size").(string)
	if haSubnet == "" && haZone == "" && haGwSize != "" {
		return fmt.Errorf("'ha_gw_size' is only required if enabling HA")
	}
	haInsaneModeAz := d.Get("ha_insane_mode_az").(string)
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
			gateway.Subnet = strings.Join(strs, "~~")
		}
		gateway.InsaneMode = "on"
	} else {
		gateway.InsaneMode = "off"
	}
	if haZone != "" || haSubnet != "" {
		if haGwSize == "" {
			return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
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

	enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
	customerManagedKeys := d.Get("customer_managed_keys").(string)
	if enableEncryptVolume && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}
	if customerManagedKeys != "" {
		if !enableEncryptVolume {
			return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
		}
		gateway.EncVolume = "no"
	}
	if !enableEncryptVolume && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		gateway.EncVolume = "no"
	}

	enableMonitorSubnets := d.Get("enable_monitor_gateway_subnets").(bool)
	var excludedInstances []string
	for _, v := range d.Get("monitor_exclude_list").(*schema.Set).List() {
		excludedInstances = append(excludedInstances, v.(string))
	}
	// Enable monitor gateway subnets does not work with AWSChina
	if enableMonitorSubnets && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes^goaviatrix.AWSChina) {
		return fmt.Errorf("'enable_monitor_gateway_subnets' is only valid for AWS (1), AWSGov (256), AWS Top Secret (16384) or AWS Secret (32768)")
	}
	if !enableMonitorSubnets && len(excludedInstances) != 0 {
		return fmt.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
	}

	enablePrivateOob := d.Get("enable_private_oob").(bool)
	oobManagementSubnet := d.Get("oob_management_subnet").(string)
	oobAvailabilityZone := d.Get("oob_availability_zone").(string)
	haOobManagementSubnet := d.Get("ha_oob_management_subnet").(string)
	haOobAvailabilityZone := d.Get("ha_oob_availability_zone").(string)

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
			return fmt.Errorf("\"ha_oob_management_subnet\" must be empty if \"enable_private_oob\" is false")
		}
	}

	_, tagListOk := d.GetOk("tag_list")
	_, tagsOk := d.GetOk("tags")
	if tagListOk || tagsOk {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return errors.New("failed to create spoke gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) or AWS Secret (32768)")
		}

		if tagListOk {
			tagList := d.Get("tag_list").([]interface{})
			tagListStr := goaviatrix.ExpandStringList(tagList)
			tagListStr = goaviatrix.TagListStrColon(tagListStr)
			gateway.TagList = strings.Join(tagListStr, ",")
		} else {
			tagsMap, err := extractTags(d, gateway.CloudType)
			if err != nil {
				return fmt.Errorf("error creating tags for spoke gateway: %v", err)
			}
			tagJson, err := TagsMapToJson(tagsMap)
			if err != nil {
				return fmt.Errorf("failed to add tags whenc creating spoke gateway: %v", err)
			}
			gateway.TagJson = tagJson
		}
	}

	enableActiveStandby := d.Get("enable_active_standby").(bool)
	if haSubnet == "" && haZone == "" && enableActiveStandby {
		return fmt.Errorf("could not configure Active-Standby as HA is not enabled")
	}
	if !enableBgp && enableActiveStandby {
		return fmt.Errorf("could not configure Active-Standby as it is not BGP capable gateway")
	}
	enableActiveStandbyPreemptive := d.Get("enable_active_standby_preemptive").(bool)
	if !enableActiveStandby && enableActiveStandbyPreemptive {
		return fmt.Errorf("could not configure Preemptive Mode with Active-Standby disabled")
	}

	enableSpotInstance := d.Get("enable_spot_instance").(bool)
	spotPrice := d.Get("spot_price").(string)
	if enableSpotInstance {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("enable_spot_instance only supports AWS related cloud types")
		}
		gateway.EnableSpotInstance = true
		gateway.SpotPrice = spotPrice
	} else {
		if spotPrice != "" {
			return fmt.Errorf("spot_price is set for enabling spot instance. Please set enable_spot_instance to true")
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
				gateway.Eip = fmt.Sprintf("%s:%s", azureEipName.(string), d.Get("eip").(string))
			} else {
				if azureEipNameOk {
					return fmt.Errorf("failed to create spoke gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				gateway.Eip = d.Get("eip").(string)
			}
		}
	}

	if privateModeInfo.EnablePrivateMode {
		if privateModeSubnetZone, ok := d.GetOk("private_mode_subnet_zone"); ok {
			gateway.Subnet = fmt.Sprintf("%s~~%s", gateway.Subnet, privateModeSubnetZone.(string))
		} else {
			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("%q must be set when creating a Spoke Gateway in AWS with Private Mode enabled on the Controller", "private_mode_subnet_zone")
			}
		}

		if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
				return fmt.Errorf("private mode is only supported in AWS and Azure. %q must be empty", "private_mode_lb_vpc_id")
			}

			gateway.LbVpcId = d.Get("private_mode_lb_vpc_id").(string)
		}
	} else {
		if _, ok := d.GetOk("private_mode_subnet_zone"); ok {
			return fmt.Errorf("%q is only valid when Private Mode is enabled on the Controller", "private_mode_subnet_zone")
		}
		if _, ok := d.GetOk("private_mode_lb_vpc_id"); ok {
			return fmt.Errorf("%q is only valid when Private Mode is enabled", "private_mode_lb_vpc_id")
		}
	}

	log.Printf("[INFO] Creating Aviatrix Spoke Gateway: %#v", gateway)

	d.SetId(gateway.GwName)
	flag := false
	defer resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag)

	err := client.LaunchSpokeVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Spoke Gateway: %s", err)
	}

	if customerManagedKeys != "" && enableEncryptVolume {
		gwEncVolume := &goaviatrix.Gateway{
			GwName:              d.Get("gw_name").(string),
			CustomerManagedKeys: d.Get("customer_managed_keys").(string),
		}
		err := client.EnableEncryptVolume(gwEncVolume)
		if err != nil {
			return fmt.Errorf("failed to enable encrypt gateway volume when creating spoke gateway: %s due to %s", gwEncVolume.GwName, err)
		}
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
		//Enable HA
		haGateway := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
			HASubnet:  haSubnet,
			HAZone:    haZone,
			Eip:       d.Get("ha_eip").(string),
		}

		if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			haGateway.HASubnetGCP = haSubnet
		}

		if insaneMode && goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			var haStrs []string
			haStrs = append(haStrs, haSubnet, haInsaneModeAz)
			haSubnet = strings.Join(haStrs, "~~")
			haGateway.HASubnet = haSubnet
		}

		if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haZone != "" {
			haGateway.HASubnet = fmt.Sprintf("%s~~%s~~", haSubnet, haZone)
		}

		if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			haGateway.Subnet = haSubnet
			haGateway.AvailabilityDomain = haAvailabilityDomain
			haGateway.FaultDomain = haFaultDomain
		}

		if enablePrivateOob {
			haGateway.HASubnet = haGateway.HASubnet + "~~" + haOobAvailabilityZone
			haGateway.HAOobManagementSubnet = haOobManagementSubnet + "~~" + haOobAvailabilityZone
		}

		if privateModeInfo.EnablePrivateMode {
			haPrivateModeSubnetZone := d.Get("ha_private_mode_subnet_zone").(string)
			haGateway.HASubnet = haSubnet + "~~" + haPrivateModeSubnetZone
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if haGateway.Eip != "" {
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				haGateway.Eip = fmt.Sprintf("%s:%s", haAzureEipName.(string), haGateway.Eip)
			} else if haAzureEipNameOk {
				return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be empty when 'ha_eip' is empty")
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			err = client.EnableHaSpokeGateway(haGateway)
		} else {
			err = client.EnableHaSpokeVpc(haGateway)
		}
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
		}

		log.Printf("[INFO]Resizing Spoke HA Gateway: %#v", haGwSize)

		if haGwSize != gateway.VpcSize {
			if haGwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set")
			}

			haGateway := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string) + "-hagw",
				GwSize:    d.Get("ha_gw_size").(string),
			}

			log.Printf("[INFO] Resizing Spoke HA Gateway size to: %s ", haGateway.GwSize)

			err := client.UpdateGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %s", err)
			}

			d.Set("ha_gw_size", haGwSize)
		}
	}

	enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && enableVpcDnsServer {
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

	if customizedSpokeVpcRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                   d.Get("gw_name").(string),
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
				return fmt.Errorf("failed to customize spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if filteredSpokeVpcRoutes := d.Get("filtered_spoke_vpc_routes").(string); filteredSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                 d.Get("gw_name").(string),
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
				return fmt.Errorf("failed to edit filtered spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if includedAdvertisedSpokeRoutes := d.Get("included_advertised_spoke_routes").(string); includedAdvertisedSpokeRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                d.Get("gw_name").(string),
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
				return fmt.Errorf("failed to edit advertised spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if enableMonitorSubnets {
		err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
		if err != nil {
			return fmt.Errorf("could not enable monitor gateway subnets: %v", err)
		}
	}

	if transitGwName := d.Get("transit_gw").(string); transitGwName != "" {
		if manageTransitGwAttachment {
			gws := strings.Split(d.Get("transit_gw").(string), ",")
			for _, gw := range gws {
				gateway.TransitGateway = gw
				try, maxTries, backoff := 0, 8, 1000*time.Millisecond
				for {
					try++
					err := client.SpokeJoinTransit(gateway)
					if err != nil {
						if strings.Contains(err.Error(), "is not up") {
							if try == maxTries {
								return fmt.Errorf("spoke gateway %s couldn't join transit gateway %q: %v", gateway.GwName, gw, err)
							}
							time.Sleep(backoff)
							// Double the backoff time after each failed try
							backoff *= 2
							continue
						}
						return fmt.Errorf("spoke gateway %s failed to join transit gateway %q: %v", gateway.GwName, gw, err)
					}
					break
				}
			}
		} else {
			return fmt.Errorf("'manage_transit_gateway_attachment' is set to false. Please set it to true, or use " +
				"'aviatrix_spoke_transit_attachment' to attach this spoke to transit gateways")
		}
	}

	if !d.Get("enable_jumbo_frame").(bool) {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		err := client.DisableJumboFrame(gw)
		if err != nil {
			return fmt.Errorf("could not disable jumbo frame for spoke gateway: %v", err)
		}
	}

	if d.Get("enable_private_vpc_default_route").(bool) {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		err := client.EnablePrivateVpcDefaultRoute(gw)
		if err != nil {
			return fmt.Errorf("could not enable private vpc default route after spoke gateway creation: %v", err)
		}
	}

	if d.Get("enable_skip_public_route_table_update").(bool) {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		err := client.EnableSkipPublicRouteUpdate(gw)
		if err != nil {
			return fmt.Errorf("could not enable skip public route update after spoke gateway creation: %v", err)
		}
	}

	if d.Get("enable_auto_advertise_s2c_cidrs").(bool) {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		err := client.EnableAutoAdvertiseS2CCidrs(gw)
		if err != nil {
			return fmt.Errorf("could not enable auto advertise s2c cidrs after spoke gateaway creation: %v", err)
		}
	}

	if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
		err := client.ModifyTunnelDetectionTime(d.Get("gw_name").(string), detectionTime.(int))
		if err != nil {
			return fmt.Errorf("could not set tunnel detection time during Spoke Gateway creation: %v", err)
		}
	}

	if learnedCidrsApproval {
		gateway.LearnedCidrsApproval = "on"
		err := client.EnableSpokeLearnedCidrsApproval(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
		}
	}
	if len(gateway.ApprovedLearnedCidrs) != 0 {
		err := client.UpdateSpokePendingApprovedCidrs(gateway)
		if err != nil {
			return fmt.Errorf("failed to update approved CIDRs: %v", err)
		}
	}

	if val, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		var spokeBgpManualSpokeAdvertiseCidrs []string
		slice := val.([]interface{})
		for _, v := range slice {
			spokeBgpManualSpokeAdvertiseCidrs = append(spokeBgpManualSpokeAdvertiseCidrs, v.(string))
		}
		gateway.BgpManualSpokeAdvertiseCidrs = strings.Join(spokeBgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetSpokeBgpManualAdvertisedNetworks(gateway)
		if err != nil {
			return fmt.Errorf("failed to set spoke BGP Manual Advertise Cidrs: %s", err)
		}
	}

	if val, ok := d.GetOk("bgp_ecmp"); ok {
		err := client.SetBgpEcmpSpoke(gateway, val.(bool))
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp: %v", err)
		}
	}

	if enableActiveStandby {
		if enableActiveStandbyPreemptive {
			if err := client.EnableActiveStandbyPreemptiveSpoke(gateway); err != nil {
				return fmt.Errorf("could not enable Preemptive Mode for Active-Standby: %v", err)
			}
		} else {
			if err := client.EnableActiveStandbySpoke(gateway); err != nil {
				return fmt.Errorf("could not enable Active-Standby: %v", err)
			}
		}
	}

	if disableRoutePropagation {
		if err := client.DisableSpokeOnpremRoutePropagation(gateway); err != nil {
			return fmt.Errorf("could not disable route propagation for Spoke %s : %v", gateway.GwName, err)
		}
	}

	if val, ok := d.GetOk("local_as_number"); ok {
		err := client.SetLocalASNumberSpoke(gateway, val.(string))
		if err != nil {
			return fmt.Errorf("could not set local_as_number: %v", err)
		}
	}

	if val, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		slice := val.([]interface{})
		for _, v := range slice {
			prependASPath = append(prependASPath, v.(string))
		}
		err := client.SetPrependASPathSpoke(gateway, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %v", err)
		}
	}

	if val, ok := d.GetOk("bgp_polling_time"); ok {
		err := client.SetBgpPollingTimeSpoke(gateway, strconv.Itoa(val.(int)))
		if err != nil {
			return fmt.Errorf("could not set bgp polling time: %v", err)
		}
	}

	if holdTime := d.Get("bgp_hold_time").(int); holdTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gateway.GwName, holdTime)
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time after Spoke Gateway creation: %v", err)
		}
	}

	enableSpokePreserveAsPath := d.Get("enable_preserve_as_path").(bool)
	if enableSpokePreserveAsPath {
		if enableBgp {
			err := client.EnableSpokePreserveAsPath(gateway)
			if err != nil {
				return fmt.Errorf("could not enable spoke preserve as path: %v", err)
			}
		} else {
			return fmt.Errorf("enable_preserve_as_path is not supported for Non-BGP Spoke Gateways")
		}
	}

	if rxQueueSize != "" {
		gwRxQueueSize := &goaviatrix.Gateway{
			GwName:      d.Get("gw_name").(string),
			RxQueueSize: rxQueueSize,
		}
		err := client.SetRxQueueSize(gwRxQueueSize)
		if err != nil {
			return fmt.Errorf("failed to set rx queue size for spoke %s: %s", gateway.GwName, err)
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
	client := meta.(*goaviatrix.Client)
	ignoreTagsConfig := client.IgnoreTagsConfig

	var isImport bool
	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.Set("manage_transit_gateway_attachment", true)
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Spoke Gateway: %s", err)
	}

	log.Printf("[TRACE] reading spoke gateway %s: %#v", d.Get("gw_name").(string), gw)

	d.Set("cloud_type", gw.CloudType)
	d.Set("account_name", gw.AccountName)
	d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
	d.Set("enable_private_vpc_default_route", gw.PrivateVpcDefaultEnabled)
	d.Set("enable_skip_public_route_table_update", gw.SkipPublicVpcUpdateEnabled)
	d.Set("enable_auto_advertise_s2c_cidrs", gw.AutoAdvertiseCidrsEnabled)
	d.Set("eip", gw.PublicIP)
	d.Set("subnet", gw.VpcNet)
	d.Set("gw_size", gw.GwSize)
	d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
	d.Set("security_group_id", gw.GwSecurityGroupID)
	d.Set("private_ip", gw.PrivateIP)
	d.Set("single_az_ha", gw.SingleAZ == "yes")
	d.Set("enable_vpc_dns_server", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled")
	d.Set("single_ip_snat", gw.EnableNat == "yes" && gw.SnatMode == "primary")
	d.Set("enable_jumbo_frame", gw.JumboFrame)
	d.Set("enable_bgp", gw.EnableBgp)
	d.Set("enable_learned_cidrs_approval", gw.EnableLearnedCidrsApproval)
	d.Set("enable_preserve_as_path", gw.EnablePreserveAsPath)
	d.Set("rx_queue_size", gw.RxQueueSize)
	d.Set("public_ip", gw.PublicIP)

	if gw.EnableLearnedCidrsApproval {
		spokeAdvancedConfig, err := client.GetSpokeGatewayAdvancedConfig(&goaviatrix.SpokeVpc{GwName: gw.GwName})
		if err != nil {
			return fmt.Errorf("could not get advanced config for spoke gateway: %v", err)
		}

		if err = d.Set("approved_learned_cidrs", spokeAdvancedConfig.ApprovedLearnedCidrs); err != nil {
			return fmt.Errorf("could not set approved_learned_cidrs into state: %v", err)
		}
	} else {
		d.Set("approved_learned_cidrs", nil)
	}
	d.Set("local_as_number", gw.LocalASNumber)
	d.Set("bgp_ecmp", gw.BgpEcmp)
	d.Set("enable_active_standby", gw.EnableActiveStandby)
	d.Set("enable_active_standby_preemptive", gw.EnableActiveStandbyPreemptive)
	d.Set("disable_route_propagation", gw.DisableRoutePropagation)
	var prependAsPath []string
	for _, p := range strings.Split(gw.PrependASPath, " ") {
		if p != "" {
			prependAsPath = append(prependAsPath, p)
		}
	}
	err = d.Set("prepend_as_path", prependAsPath)
	if err != nil {
		return fmt.Errorf("could not set prepend_as_path: %v", err)
	}
	if gw.EnableBgp {
		d.Set("learned_cidrs_approval_mode", gw.LearnedCidrsApprovalMode)
		d.Set("bgp_polling_time", gw.BgpPollingTime)
		d.Set("bgp_hold_time", gw.BgpHoldTime)
	} else {
		d.Set("learned_cidrs_approval_mode", "gateway")
		d.Set("bgp_polling_time", 50)
		d.Set("bgp_hold_time", 180)
	}
	d.Set("tunnel_detection_time", gw.TunnelDetectionTime)
	d.Set("image_version", gw.ImageVersion)
	d.Set("software_version", gw.SoftwareVersion)

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			d.Set("azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Spoke Gateway %s", gw.GwName)
		}
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0]) //AWS vpc_id returns as <vpc_id>~~<other vpc info> in rest api
		d.Set("vpc_reg", gw.VpcRegion)                    //AWS vpc_reg returns as vpc_region in rest api

		if gw.AllocateNewEipRead && !gw.EnablePrivateOob {
			d.Set("allocate_new_eip", true)
		} else {
			d.Set("allocate_new_eip", false)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		// gcp vpc_id returns as <vpc name>~-~<project name>
		d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.GatewayZone) //gcp vpc_reg returns as gateway_zone in json

		d.Set("allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0]) // oci vpc_id returns as <vpc_id>~~<vpc_name> in rest api
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
		if customizedSpokeVpcRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedSpokeVpcRoutes != "" {
			customizedRoutesArray := strings.Split(customizedSpokeVpcRoutes, ",")
			if len(goaviatrix.Difference(customizedRoutesArray, gw.CustomizedSpokeVpcRoutes)) == 0 &&
				len(goaviatrix.Difference(gw.CustomizedSpokeVpcRoutes, customizedRoutesArray)) == 0 {
				d.Set("customized_spoke_vpc_routes", customizedSpokeVpcRoutes)
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

	if len(gw.IncludeCidrList) != 0 {
		if includedAdvertisedSpokeRoutes := d.Get("included_advertised_spoke_routes").(string); includedAdvertisedSpokeRoutes != "" {
			advertisedSpokeRoutesArray := strings.Split(includedAdvertisedSpokeRoutes, ",")
			if len(goaviatrix.Difference(advertisedSpokeRoutesArray, gw.IncludeCidrList)) == 0 &&
				len(goaviatrix.Difference(gw.IncludeCidrList, advertisedSpokeRoutesArray)) == 0 {
				d.Set("included_advertised_spoke_routes", includedAdvertisedSpokeRoutes)
			} else {
				d.Set("included_advertised_spoke_routes", strings.Join(gw.IncludeCidrList, ","))
			}
		} else {
			d.Set("included_advertised_spoke_routes", strings.Join(gw.AdvertisedSpokeRoutes, ","))
		}
	} else {
		d.Set("included_advertised_spoke_routes", "")
	}

	d.Set("enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
	if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
		return fmt.Errorf("setting 'monitor_exclude_list' to state: %v", err)
	}

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)
	if manageTransitGwAttachment {
		if gw.SpokeVpc == "yes" {
			var transitGws []string
			if gw.EgressTransitGwName != "" {
				transitGws = append(transitGws, gw.EgressTransitGwName)
			}
			if gw.TransitGwName != "" {
				transitGws = append(transitGws, gw.TransitGwName)
			}
			d.Set("transit_gw", strings.Join(transitGws, ","))
		} else {
			d.Set("transit_gw", "")
		}
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		if _, ok := d.GetOk("tag_list"); ok {
			tagList := make([]string, 0, len(gw.Tags))
			for key, val := range gw.Tags {
				str := key + ":" + val
				tagList = append(tagList, str)
			}

			tagListFromUserConfig := d.Get("tag_list").([]interface{})
			tagListStr := goaviatrix.ExpandStringList(tagListFromUserConfig)

			if len(goaviatrix.Difference(tagListStr, tagList)) != 0 || len(goaviatrix.Difference(tagList, tagListStr)) != 0 {
				if err := d.Set("tag_list", tagList); err != nil {
					log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
				}
			} else {
				if err := d.Set("tag_list", tagListStr); err != nil {
					log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
				}
			}
		} else {
			tags := goaviatrix.KeyValueTags(gw.Tags).IgnoreConfig(ignoreTagsConfig)
			if err := d.Set("tags", tags); err != nil {
				log.Printf("[WARN] Error setting tags for (%s): %s", d.Id(), err)
			}
		}
	}

	var spokeBgpManualAdvertiseCidrs []string
	if val, ok := d.GetOk("spoke_bgp_manual_advertise_cidrs"); ok {
		slice := val.([]interface{})
		for _, v := range slice {
			spokeBgpManualAdvertiseCidrs = append(spokeBgpManualAdvertiseCidrs, v.(string))
		}
	}
	if len(goaviatrix.Difference(spokeBgpManualAdvertiseCidrs, gw.BgpManualSpokeAdvertiseCidrs)) != 0 ||
		len(goaviatrix.Difference(gw.BgpManualSpokeAdvertiseCidrs, spokeBgpManualAdvertiseCidrs)) != 0 {
		d.Set("spoke_bgp_manual_advertise_cidrs", gw.BgpManualSpokeAdvertiseCidrs)
	} else {
		d.Set("spoke_bgp_manual_advertise_cidrs", spokeBgpManualAdvertiseCidrs)
	}

	d.Set("enable_private_oob", gw.EnablePrivateOob)
	if gw.EnablePrivateOob {
		d.Set("oob_management_subnet", strings.Split(gw.OobManagementSubnet, "~~")[0])
		d.Set("oob_availability_zone", gw.GatewayZone)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		_, zoneIsSet := d.GetOk("zone")
		if (isImport || zoneIsSet) && gw.GatewayZone != "AvailabilitySet" {
			d.Set("zone", "az-"+gw.GatewayZone)
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
	}

	d.Set("private_mode_lb_vpc_id", gw.LbVpcId)
	if gw.LbVpcId != "" && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("private_mode_subnet_zone", gw.GatewayZone)
	} else {
		d.Set("private_mode_subnet_zone", nil)
	}

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

	log.Printf("[INFO] Spoke HA Gateway size: %s", gw.HaGw.GwSize)
	if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		d.Set("ha_subnet", gw.HaGw.VpcNet)
		if zone := d.Get("ha_zone"); goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || zone.(string) != "") {
			if gw.HaGw.GatewayZone != "AvailabilitySet" {
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
		} else {
			d.Set("ha_subnet", "")
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
	d.Set("ha_gw_size", gw.HaGw.GwSize)
	d.Set("ha_cloud_instance_id", gw.HaGw.CloudnGatewayInstID)
	d.Set("ha_gw_name", gw.HaGw.GwName)
	d.Set("ha_private_ip", gw.HaGw.PrivateIP)
	d.Set("ha_software_version", gw.HaGw.SoftwareVersion)
	d.Set("ha_image_version", gw.HaGw.ImageVersion)
	d.Set("ha_security_group_id", gw.HaGw.GwSecurityGroupID)
	d.Set("ha_public_ip", gw.HaGw.PublicIP)
	if gw.HaGw.InsaneMode == "yes" && goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("ha_insane_mode_az", gw.HaGw.GatewayZone)
	} else {
		d.Set("ha_insane_mode_az", "")
	}
	if gw.HaGw.EnablePrivateOob {
		d.Set("ha_oob_management_subnet", strings.Split(gw.HaGw.OobManagementSubnet, "~~")[0])
		d.Set("ha_oob_availability_zone", gw.HaGw.GatewayZone)
	}
	if gw.LbVpcId != "" && goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("ha_private_mode_subnet_zone", gw.HaGw.GatewayZone)
	} else {
		d.Set("ha_private_mode_subnet_zone", "")
	}
	if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
		if len(azureEip) == 3 {
			d.Set("ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the HA Gateway %s", gw.GwName)
		}
	}

	return nil
}

func resourceAviatrixSpokeGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	haGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
		GwSize:    d.Get("ha_gw_size").(string),
	}

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	d.Partial(true)
	if d.Get("enable_private_vpc_default_route").(bool) && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_private_vpc_default_route is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}
	if d.Get("enable_skip_public_route_table_update").(bool) && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("enable_skip_public_route_update is only valid for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if d.HasChange("ha_zone") {
		haZone := d.Get("ha_zone").(string)
		if haZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("'ha_zone' is only valid for GCP (4), Azure (8), AzureGov (32) and AzureChina (2048) providers if enabling HA")
		}
	}
	if d.HasChange("ha_zone") || d.HasChange("ha_subnet") {
		haZone := d.Get("ha_zone").(string)
		haSubnet := d.Get("ha_subnet").(string)
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && haSubnet != "" && haZone == "" {
			return fmt.Errorf("'ha_zone' must be set to enable HA on GCP (4), cannot enable HA with only 'ha_subnet'")
		}
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" && haZone != "" {
			return fmt.Errorf("'ha_subnet' must be provided to enable HA for Azure (8), AzureGov (32) or AzureChina (2048), cannot enable HA with only 'ha_zone'")
		}
	}
	if d.HasChange("cloud_type") {
		return fmt.Errorf("updating cloud_type is not allowed")
	}
	if d.HasChange("account_name") {
		return fmt.Errorf("updating account_name is not allowed")
	}
	if d.HasChange("insane_mode") {
		return fmt.Errorf("updating insane_mode is not allowed")
	}
	if d.HasChange("insane_mode_az") {
		return fmt.Errorf("updating insane_mode_az is not allowed")
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
		return fmt.Errorf("failed to update spoke gateway: changing 'azure_eip_name_resource_group' is not allowed")
	}
	if d.HasChange("ha_azure_eip_name_resource_group") {
		o, n := d.GetChange("ha_azure_eip_name_resource_group")
		if o.(string) != "" && n.(string) != "" {
			return fmt.Errorf("failed to update spoke gateway: changing 'ha_azure_eip_name_resource_group' is not allowed")
		}
	}
	if d.HasChange("enable_spot_instance") {
		return fmt.Errorf("updating enable_spot_instance is not allowed")
	}
	if d.HasChange("spot_price") {
		return fmt.Errorf("updating spot_price is not allowed")
	}

	learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
	approvedLearnedCidrs := getStringSet(d, "approved_learned_cidrs")
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

	if d.HasChange("enable_preserve_as_path") {
		enableBgp := d.Get("enable_bgp").(bool)
		enableSpokePreserveAsPath := d.Get("enable_preserve_as_path").(bool)
		if enableSpokePreserveAsPath && !enableBgp {
			return fmt.Errorf("enable_preserve_as_path is not supported for Non-BGP Spoke during Spoke Gateway update")
		}
		if !enableSpokePreserveAsPath {
			err := client.DisableSpokePreserveAsPath(&goaviatrix.SpokeVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not disable Preserve AS Path during Spoke Gateway update: %v", err)
			}
		} else {
			err := client.EnableSpokePreserveAsPath(&goaviatrix.SpokeVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not enable Preserve AS Path during Spoke Gateway update: %v", err)
			}
		}
	}

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)
	if d.HasChange("manage_transit_gateway_attachment") {
		_, nMTGA := d.GetChange("manage_transit_gateway_attachment")
		newManageTransitGwAttachment := nMTGA.(bool)
		if newManageTransitGwAttachment {
			d.Set("manage_transit_gateway_attachment", true)
		} else {
			d.Set("manage_transit_gateway_attachment", false)
		}
	}
	if !manageTransitGwAttachment && d.Get("transit_gw").(string) != "" {
		return fmt.Errorf("'manage_transit_gateway_attachment' is set to false. Please set it to true, or use " +
			"'aviatrix_spoke_transit_attachment' to attach this spoke to transit gateways")
	}

	if d.HasChange("tag_list") || d.HasChange("tags") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("error updating spoke gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			CloudType:    gateway.CloudType,
		}
		tagList := goaviatrix.ExpandStringList(d.Get("tag_list").([]interface{}))

		if d.HasChange("tag_list") {
			tagList = goaviatrix.TagListStrColon(tagList)
			tags.TagList = strings.Join(tagList, ",")
			err := client.UpdateTags(tags)
			if err != nil {
				return fmt.Errorf("failed to update tags for spoke gateway: %s", err)
			}
		}
		if d.HasChange("tags") && len(tagList) == 0 {
			tagsMap, err := extractTags(d, gateway.CloudType)
			if err != nil {
				return fmt.Errorf("failed to update tags for spoke gateway: %v", err)
			}
			tags.Tags = tagsMap
			tagJson, err := TagsMapToJson(tagsMap)
			if err != nil {
				return fmt.Errorf("failed to update tags for spoke gateway: %v", err)
			}
			tags.TagJson = tagJson
			err = client.UpdateTags(tags)
			if err != nil {
				return fmt.Errorf("failed to update tags for spoke gateway: %v", err)
			}
		}
	}

	//Get primary gw size if gw_size changed, to be used later on for ha gateway size update
	primaryGwSize := d.Get("gw_size").(string)
	if d.HasChange("gw_size") {
		old, _ := d.GetChange("gw_size")
		primaryGwSize = old.(string)
		gateway.GwSize = d.Get("gw_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Spoke Gateway: %s", err)
		}
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") ||
		(enablePrivateOob && (d.HasChange("ha_oob_management_subnet") || d.HasChange("ha_oob_availability_zone"))) ||
		(privateModeInfo.EnablePrivateMode && d.HasChange("ha_private_mode_subnet_zone")) ||
		d.HasChange("ha_availability_domain") || d.HasChange("ha_fault_domain") {
		spokeGw := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
			GwSize:    d.Get("ha_gw_size").(string),
		}

		haEip := d.Get("ha_eip").(string)
		if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			spokeGw.Eip = haEip
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if haEip != "" && spokeGw.GwSize != "" {
				// No change will be detected when ha_eip is set to the empty string because it is computed.
				// Instead, check ha_gw_size to detect when HA gateway is being deleted.
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create HA Spoke Gateway: 'ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				spokeGw.Eip = fmt.Sprintf("%s:%s", haAzureEipName.(string), haEip)
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create HA Spoke Gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false
		if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			spokeGw.HASubnet = d.Get("ha_subnet").(string)
			if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && d.Get("ha_zone").(string) != "" {
				spokeGw.HASubnet = fmt.Sprintf("%s~~%s~~", d.Get("ha_subnet").(string), d.Get("ha_zone").(string))
			}

			haAvailabilityDomain := d.Get("ha_availability_domain").(string)
			haFaultDomain := d.Get("ha_fault_domain").(string)
			if newSubnet != "" {
				if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain == "" || haFaultDomain == "") {
					return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are required to enable HA on OCI")
				}
				if !goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.OCIRelatedCloudTypes) && (haAvailabilityDomain != "" || haFaultDomain != "") {
					return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are only valid for OCI")
				}
			}
			if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
				spokeGw.Subnet = d.Get("ha_subnet").(string)
				spokeGw.AvailabilityDomain = haAvailabilityDomain
				spokeGw.FaultDomain = haFaultDomain
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
		} else if goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			spokeGw.HAZone = d.Get("ha_zone").(string)
			spokeGw.HASubnetGCP = d.Get("ha_subnet").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}

		if d.Get("insane_mode").(bool) && goaviatrix.IsCloudType(spokeGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			haSubnet := d.Get("ha_subnet").(string)

			if insaneModeHaAz == "" && haSubnet != "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set " +
					"for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			} else if insaneModeHaAz != "" && haSubnet == "" {
				return fmt.Errorf("ha_subnet needed if insane_mode is enabled and ha_insane_mode_az is set " +
					"for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}

			haStrs = append(haStrs, spokeGw.HASubnet, insaneModeHaAz)
			spokeGw.HASubnet = strings.Join(haStrs, "~~")
		}

		if (newHaGwEnabled || changeHaGw) && spokeGw.GwSize == "" {
			return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set")
		} else if deleteHaGw && spokeGw.GwSize != "" {
			return fmt.Errorf("ha_gw_size must be empty if spoke HA gateway is deleted")

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

				spokeGw.HASubnet = spokeGw.HASubnet + "~~" + haOobAvailabilityZone
				spokeGw.HAOobManagementSubnet = haOobManagementSubnet + "~~" + haOobAvailabilityZone
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
				if _, ok := d.GetOk("ha_private_mode_subnet_zone"); !ok {
					return fmt.Errorf("%q is required if private mode is enabled and %q is provided", "ha_private_mode_subnet_zone", "ha_subnet")
				}

				privateModeSubnetZone := d.Get("ha_private_mode_subnet_zone").(string)
				spokeGw.HASubnet += "~~" + privateModeSubnetZone
			}
		}

		if newHaGwEnabled {
			//New configuration to enable HA
			if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
				err := client.EnableHaSpokeGateway(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			} else {
				err := client.EnableHaSpokeVpc(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			}
		} else if deleteHaGw {
			//Ha configuration has been deleted
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
			}
		} else if changeHaGw {
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
			}

			spokeGw.Eip = ""

			//New configuration to enable HA
			if goaviatrix.IsCloudType(haGateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
				err := client.EnableHaSpokeGateway(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			} else {
				err := client.EnableHaSpokeVpc(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
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
					return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway while trying to update HA Gw size: %s", err)
				}
			} else {
				if haGateway.GwSize == "" {
					return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
						"ha_subnet or ha_zone is set")
				}
				err = client.UpdateGateway(haGateway)
				log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.GwSize)
				if err != nil {
					return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %s", err)
				}
			}
		}
	}

	if d.HasChange("single_ip_snat") {
		enableSNat := d.Get("single_ip_snat").(bool)
		gw := &goaviatrix.Gateway{
			CloudType:   d.Get("cloud_type").(int),
			GatewayName: d.Get("gw_name").(string),
		}
		if enableSNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable single_ip' mode SNAT: %s", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable 'single_ip' mode SNAT: %s", err)
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

	if d.HasChange("enable_learned_cidrs_approval") {
		gw := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}
		if learnedCidrsApproval {
			gw.LearnedCidrsApproval = "on"
			err := client.EnableSpokeLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
			}
		} else {
			gw.LearnedCidrsApproval = "off"
			err := client.DisableSpokeLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to disable learned cidrs approval: %s", err)
			}
		}
	}

	if learnedCidrsApproval && d.HasChange("approved_learned_cidrs") {
		gw := &goaviatrix.SpokeVpc{
			GwName:               d.Get("gw_name").(string),
			ApprovedLearnedCidrs: approvedLearnedCidrs,
		}

		err := client.UpdateSpokePendingApprovedCidrs(gw)
		if err != nil {
			return fmt.Errorf("could not update approved CIDRs: %v", err)
		}
	}

	if d.HasChange("enable_encrypt_volume") {
		if d.Get("enable_encrypt_volume").(bool) {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
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
			log.Printf("[INFO] Customizeing routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to customize spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
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
			log.Printf("[INFO] Editing filtered spoke vpc routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if d.HasChange("included_advertised_spoke_routes") {
		o, n := d.GetChange("included_advertised_spoke_routes")
		oldRouteList := strings.Split(o.(string), ",")
		newRouteList := strings.Split(n.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                d.Get("gw_name").(string),
				AdvertisedSpokeRoutes: newRouteList,
			}
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			log.Printf("[INFO] Editing included advertised spoke vpc routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit included advertised spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
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
				return fmt.Errorf("could not enable monitor gateway subnets: %v", err)
			}
		} else {
			err := client.DisableMonitorGatewaySubnets(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not disable monitor gateway subnets: %v", err)
			}
		}
	} else if d.HasChange("monitor_exclude_list") {
		err := client.DisableMonitorGatewaySubnets(gateway.GwName)
		if err != nil {
			return fmt.Errorf("could not disable monitor gateway subnets: %v", err)
		}
		err = client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
		if err != nil {
			return fmt.Errorf("could not enable monitor gateway subnets: %v", err)
		}
	}

	if d.HasChange("transit_gw") {
		if manageTransitGwAttachment {
			spokeVPC := &goaviatrix.SpokeVpc{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string),
				HASubnet:  d.Get("ha_subnet").(string),
			}

			o, n := d.GetChange("transit_gw")
			oldTransitGws := strings.Split(o.(string), ",")
			newTransitGws := strings.Split(n.(string), ",")
			if len(oldTransitGws) > 0 && oldTransitGws[0] != "" {
				for _, gw := range oldTransitGws {
					// Leave any transit gateways that are in the old list but not in the new.
					if goaviatrix.Contains(newTransitGws, gw) {
						continue
					}
					spokeVPC.TransitGateway = gw
					err := client.SpokeLeaveTransit(spokeVPC)
					if err != nil {
						return fmt.Errorf("failed to leave Transit Gateway %q: %v", gw, err)
					}
				}
			}
			if len(newTransitGws) > 0 && newTransitGws[0] != "" {
				for _, gw := range newTransitGws {
					// Join any transit gateways that are in the new list but not in the old.
					if goaviatrix.Contains(oldTransitGws, gw) {
						continue
					}
					spokeVPC.TransitGateway = gw
					err := client.SpokeJoinTransit(spokeVPC)
					if err != nil {
						return fmt.Errorf("failed to join Transit Gateway %q: %v", gw, err)
					}
				}
			}
		} else {
			if !d.HasChange("manage_transit_gateway_attachment") {
				return fmt.Errorf("not able to update transit_gw since transit attachment is managed elsewhere")
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if d.Get("enable_jumbo_frame").(bool) {
			err := client.EnableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not enable jumbo frame for spoke gateway when updating: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not disable jumbo frame for spoke gateway when updating: %v", err)
			}
		}
	}

	if d.HasChange("enable_private_vpc_default_route") {
		if d.Get("enable_private_vpc_default_route").(bool) {
			err := client.EnablePrivateVpcDefaultRoute(gateway)
			if err != nil {
				return fmt.Errorf("could not enable private vpc default route during spoke gateway update: %v", err)
			}
		} else {
			err := client.DisablePrivateVpcDefaultRoute(gateway)
			if err != nil {
				return fmt.Errorf("could not disable private vpc default route during spoke gateway update: %v", err)
			}
		}
	}

	if d.HasChange("enable_skip_public_route_table_update") {
		if d.Get("enable_skip_public_route_table_update").(bool) {
			err := client.EnableSkipPublicRouteUpdate(gateway)
			if err != nil {
				return fmt.Errorf("could not enable skip public route update during spoke gateway update: %v", err)
			}
		} else {
			err := client.DisableSkipPublicRouteUpdate(gateway)
			if err != nil {
				return fmt.Errorf("could not disable skip public route update during spoke gateway update: %v", err)
			}
		}
	}

	if d.HasChange("enable_auto_advertise_s2c_cidrs") {
		if d.Get("enable_auto_advertise_s2c_cidrs").(bool) {
			err := client.EnableAutoAdvertiseS2CCidrs(gateway)
			if err != nil {
				return fmt.Errorf("could not enable auto advertise s2c cidrs during spoke gateway update: %v", err)
			}
		} else {
			err := client.DisableAutoAdvertiseS2CCidrs(gateway)
			if err != nil {
				return fmt.Errorf("could not disable auto advertise s2c cidrs during spoke gateway update: %v", err)
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
				return fmt.Errorf("could not get default tunnel detection time during Spoke Gateway update: %v", err)
			}
		}
		err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime)
		if err != nil {
			return fmt.Errorf("could not modify tunnel detection time during Spoke Gateway update: %v", err)
		}
	}

	primaryHasVersionChange := d.HasChanges("software_version", "image_version")
	haHasVersionChange := haEnabled && d.HasChanges("ha_software_version", "ha_image_version")
	primaryHasImageVersionChange := d.HasChange("image_version")
	haHasImageVersionChange := d.HasChange("ha_image_version")
	if primaryHasVersionChange || haHasVersionChange {
		// To determine if this is an attempted software rollback, we check if
		// old is a higher version than new. Or, the new version is the
		// special string "previous".
		oldPrimarySoftwareVersion, newPrimarySoftwareVersion := d.GetChange("software_version")
		comparePrimary, err := goaviatrix.CompareSoftwareVersions(oldPrimarySoftwareVersion.(string), newPrimarySoftwareVersion.(string))
		primaryRollbackSoftwareVersion := (err == nil && comparePrimary > 0) || newPrimarySoftwareVersion == "previous"

		oldHaSoftwareVersion, newHaSoftwareVersion := d.GetChange("ha_software_version")
		compareHa, err := goaviatrix.CompareSoftwareVersions(oldHaSoftwareVersion.(string), newHaSoftwareVersion.(string))
		haRollbackSoftwareVersion := (err == nil && compareHa > 0) || newHaSoftwareVersion == "previous"

		if primaryHasVersionChange && haHasVersionChange &&
			!primaryHasImageVersionChange && !haHasImageVersionChange &&
			!primaryRollbackSoftwareVersion && !haRollbackSoftwareVersion {
			// Both Primary and HA have upgraded just their software_version
			// so we can perform upgrade in parallel.
			log.Printf("[INFO] Upgrading spoke gateway gw_name=%s ha/primary pair in parallel", gateway.GwName)
			swVersion := d.Get("software_version").(string)
			imageVersion := d.Get("image_version").(string)
			gw := &goaviatrix.Gateway{
				GwName:          gateway.GwName,
				SoftwareVersion: swVersion,
				ImageVersion:    imageVersion,
			}
			haSwVersion := d.Get("ha_software_version").(string)
			haImageVersion := d.Get("ha_image_version").(string)
			hagw := &goaviatrix.Gateway{
				GwName:          gateway.GwName + "-hagw",
				SoftwareVersion: haSwVersion,
				ImageVersion:    haImageVersion,
			}
			var wg sync.WaitGroup
			wg.Add(2)
			var primaryErr, haErr error
			go func() {
				primaryErr = client.UpgradeGateway(gw)
				wg.Done()
			}()
			go func() {
				haErr = client.UpgradeGateway(hagw)
				wg.Done()
			}()
			wg.Wait()
			if primaryErr != nil && haErr != nil {
				return fmt.Errorf("could not upgrade primary and HA spoke gateway "+
					"software_version=%s ha_software_version=%s image_version=%s ha_image_version=%s:"+
					"\n primaryErr: %v\n haErr: %v",
					swVersion, haSwVersion, imageVersion, haImageVersion, primaryErr, haErr)
			} else if primaryErr != nil {
				return fmt.Errorf("could not upgrade primary spoke gateway software_version=%s: %v", swVersion, primaryErr)
			} else if haErr != nil {
				return fmt.Errorf("could not upgrade HA spoke gateway ha_software_version=%s: %v", haSwVersion, haErr)
			}
		} else { // Only primary or only HA has changed, or image_version changed, or it is a software rollback
			log.Printf("[INFO] Upgrading spoke gateway gw_name=%s ha or primary in serial", gateway.GwName)
			if primaryHasVersionChange {
				swVersion := d.Get("software_version").(string)
				imageVersion := d.Get("image_version").(string)
				gw := &goaviatrix.Gateway{
					GwName:          gateway.GwName,
					SoftwareVersion: swVersion,
					ImageVersion:    imageVersion,
				}
				err := client.UpgradeGateway(gw)
				if err != nil {
					return fmt.Errorf("could not upgrade spoke gateway during update image_version=%s software_version=%s: %v", gw.ImageVersion, gw.SoftwareVersion, err)
				}
			}
			if haHasVersionChange {
				haSwVersion := d.Get("ha_software_version").(string)
				haImageVersion := d.Get("ha_image_version").(string)
				hagw := &goaviatrix.Gateway{
					GwName:          gateway.GwName + "-hagw",
					SoftwareVersion: haSwVersion,
					ImageVersion:    haImageVersion,
				}
				err := client.UpgradeGateway(hagw)
				if err != nil {
					return fmt.Errorf("could not upgrade HA spoke gateway during update image_version=%s software_version=%s: %v", hagw.ImageVersion, hagw.SoftwareVersion, err)
				}
			}
		}
	}

	if d.HasChange("spoke_bgp_manual_advertise_cidrs") {
		spokeGw := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}
		var spokeBgpManualSpokeAdvertiseCidrs []string
		for _, v := range d.Get("spoke_bgp_manual_advertise_cidrs").([]interface{}) {
			spokeBgpManualSpokeAdvertiseCidrs = append(spokeBgpManualSpokeAdvertiseCidrs, v.(string))
		}
		spokeGw.BgpManualSpokeAdvertiseCidrs = strings.Join(spokeBgpManualSpokeAdvertiseCidrs, ",")
		err := client.SetSpokeBgpManualAdvertisedNetworks(spokeGw)
		if err != nil {
			return fmt.Errorf("failed to set spoke bgp manual advertise CIDRs during Spoke Gateway update: %s", err)
		}
	}

	if d.HasChange("bgp_ecmp") {
		enabled := d.Get("bgp_ecmp").(bool)
		gateway := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpEcmpSpoke(gateway, enabled)
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp during Spoke Gateway update: %v", err)
		}
	}

	if d.HasChange("enable_active_standby") || d.HasChange("enable_active_standby_preemptive") {
		gateway := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}
		if d.Get("enable_active_standby").(bool) {
			if d.Get("enable_active_standby_preemptive").(bool) {
				if err := client.EnableActiveStandbyPreemptiveSpoke(gateway); err != nil {
					return fmt.Errorf("could not enable Preemptive Mode for Active-Standby during Spoke Gatway update: %v", err)
				}
			} else {
				if err := client.EnableActiveStandbySpoke(gateway); err != nil {
					return fmt.Errorf("could not enable Active-Standby during Spoke Gateway update: %v", err)
				}
			}
		} else {
			if d.Get("enable_active_standby_preemptive").(bool) {
				return fmt.Errorf("could not enable Preemptive Mode with Active-Standby disabled")
			}
			if err := client.DisableActiveStandbySpoke(gateway); err != nil {
				return fmt.Errorf("could not disable Active-Standby during Spoke Gateway update: %v", err)
			}
		}
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		gateway := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}

		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(prependASPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPathSpoke(gateway, nil)
			if err != nil {
				return fmt.Errorf("could not delete prepend_as_path during Spoke Gateway update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			localAsNumber := d.Get("local_as_number").(string)
			err := client.SetLocalASNumberSpoke(gateway, localAsNumber)
			if err != nil {
				return fmt.Errorf("could not set local_as_number: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err := client.SetPrependASPathSpoke(gateway, prependASPath)
			if err != nil {
				return fmt.Errorf("could not set prepend_as_path during Spoke Gateway update: %v", err)
			}
		}
	}

	if d.HasChange("bgp_polling_time") {
		bgpPollingTime := d.Get("bgp_polling_time")
		gateway := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpPollingTimeSpoke(gateway, strconv.Itoa(bgpPollingTime.(int)))
		if err != nil {
			return fmt.Errorf("could not update bgp polling time during Spoke Gateway update: %v", err)
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(gateway.GwName, d.Get("bgp_hold_time").(int))
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time during Spoke Gateway update: %v", err)
		}
	}

	if d.HasChange("disable_route_propagation") {
		disableRoutePropagation := d.Get("disable_route_propagation").(bool)
		enableBgp := d.Get("enable_bgp").(bool)
		if disableRoutePropagation && !enableBgp {
			return fmt.Errorf("disable route propagation is not supported for Non-BGP Spoke during Spoke Gateway update")
		}
		gw := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}
		if disableRoutePropagation {
			err := client.DisableSpokeOnpremRoutePropagation(gw)
			if err != nil {
				return fmt.Errorf("failed to disable route propagation for Spoke %s during Spoke Gateway update: %v", gw.GwName, err)
			}
		} else {
			err := client.EnableSpokeOnpremRoutePropagation(gw)
			if err != nil {
				return fmt.Errorf("failed to enable route propagation for Spoke %s during Spoke Gateway update: %v", gw.GwName, err)
			}
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
			return fmt.Errorf("could not modify rx queue size for spoke: %s during gateway update: %v", gw.GatewayName, err)
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeGatewayRead(d, meta)
}

func resourceAviatrixSpokeGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke Gateway: %#v", gateway)

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)
	if manageTransitGwAttachment {
		if transitGw := d.Get("transit_gw").(string); transitGw != "" {
			spokeVPC := &goaviatrix.SpokeVpc{
				GwName: d.Get("gw_name").(string),
			}

			gws := strings.Split(transitGw, ",")
			for _, gw := range gws {
				spokeVPC.TransitGateway = gw
				err := client.SpokeLeaveTransit(spokeVPC)
				if err != nil {
					return fmt.Errorf("failed to leave transit gateway %q: %v", gw, err)
				}
			}
		}
	}

	//If HA is enabled, delete HA GW first.
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haSubnet != "" || haZone != "" {
		//Delete HA Gw too
		gateway.GwName += "-hagw"
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Spoke Gateway: %s", err)
	}

	return nil
}
