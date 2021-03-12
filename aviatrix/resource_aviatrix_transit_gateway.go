package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	defaultLearnedCidrApprovalMode = "gateway"
	defaultBgpHoldTime             = 180
)

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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC-ID/VNet-Name of cloud provider.",
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
				Description:  "Public Subnet Name.",
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureAZ,
				Description:  "Availability Zone. Only available for cloud_type = 8 (AZURE). Must be in the form 'az-n', for example, 'az-2'.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Transit Gateway. Required for AWS if insane_mode is enabled.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("enable_private_oob").(bool)
				},
				Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"ha_subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
				Description: "HA Subnet. Required for enabling HA for AWS/AZURE/AWSGOV gateway. " +
					"Optional for enabling HA for GCP gateway.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Zone. Required if enabling HA for GCP. Optional for AZURE.",
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
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "Instance tag of cloud provider.",
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
				Description: "Enable Insane Mode for Transit. Valid values: true, false. If insane mode is enabled, gateway size has to at least be c5 size for AWS and Standard_D3_v2 size for AZURE.",
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
			"enable_active_mesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable Active Mesh Mode for Transit Gateway. Valid values: true, false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Only supports AWS/AWSGOV. Valid values: true, false.",
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
				Description:      "Intended CIDR list to advertise to VGW.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS and AWSGOV providers. Valid values: true, false. Default value: false.",
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
			"enable_transit_firenet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable transit firenet interfaces or not.",
			},
			"lan_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LAN VPC ID. Only used for GCP Transit FireNet.",
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
			"bgp_polling_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "50",
				Description: "BGP route polling time. Unit is in seconds. Valid values are between 10 and 50.",
			},
			"prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
			},
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Changes the Aviatrix Transit Gateway ASN number before you setup Aviatrix Transit Gateway connection configurations.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"bgp_ecmp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Equal Cost Multi Path (ECMP) routing for the next hop.",
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
				Description: "Enables Active-Standby Mode, available only with Active Mesh Mode and HA enabled.",
			},
			"enable_monitor_gateway_subnets": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Enable [monitor gateway subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet). " +
					"Only valid for cloud_type = 1 (AWS) or 256 (AWSGOV). Valid values: true, false. Default value: false.",
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
				ForceNew:    true,
				Description: "Pre-allocate a network interface(eth4) for \"BGP over LAN\" functionality. Only valid for cloud_type = 8 (AZURE). Valid values: true or false. Default value: false. Available as of provider version R2.18+",
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
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the transit gateway.",
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
		},
	}
}

func resourceAviatrixTransitGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.TransitVpc{
		CloudType:                d.Get("cloud_type").(int),
		AccountName:              d.Get("account_name").(string),
		GwName:                   d.Get("gw_name").(string),
		VpcID:                    d.Get("vpc_id").(string),
		VpcSize:                  d.Get("gw_size").(string),
		Subnet:                   d.Get("subnet").(string),
		EnableHybridConnection:   d.Get("enable_hybrid_connection").(bool),
		EnableSummarizeCidrToTgw: d.Get("enable_transit_summarize_cidr_to_tgw").(bool),
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

	enablePrivateOob := d.Get("enable_private_oob").(bool)

	if !enablePrivateOob {
		allocateNewEip := d.Get("allocate_new_eip").(bool)
		if allocateNewEip {
			gateway.ReuseEip = "off"
		} else {
			gateway.ReuseEip = "on"
			gateway.Eip = d.Get("eip").(string)
		}
	}

	cloudType := d.Get("cloud_type").(int)
	zone := d.Get("zone").(string)
	if cloudType != goaviatrix.AZURE && zone != "" {
		return fmt.Errorf("attribute 'zone' is only for use with cloud_type = 8 (AZURE)")
	}
	if zone != "" {
		// The API uses the same string field to hold both subnet and zone
		// parameters.
		gateway.Subnet = fmt.Sprintf("%s~~%s~~", d.Get("subnet").(string), zone)
	}

	if cloudType == goaviatrix.AWS || cloudType == goaviatrix.GCP || cloudType == goaviatrix.OCI || cloudType == goaviatrix.AWSGOV {
		gateway.VpcID = d.Get("vpc_id").(string)
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw")
		}
	} else if cloudType == goaviatrix.AZURE {
		gateway.VNetNameResourceGroup = d.Get("vpc_id").(string)
		if gateway.VNetNameResourceGroup == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw")
		}
	}

	if gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AZURE || gateway.CloudType == goaviatrix.OCI || gateway.CloudType == goaviatrix.AWSGOV {
		gateway.VpcRegion = d.Get("vpc_reg").(string)
	} else if gateway.CloudType == goaviatrix.GCP {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = d.Get("vpc_reg").(string)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), AZURE (8), OCI (16), or AWSGOV (256)")
	}

	insaneMode := d.Get("insane_mode").(bool)
	if insaneMode {
		if cloudType != goaviatrix.AWS && cloudType != goaviatrix.GCP && cloudType != goaviatrix.AZURE && cloudType != goaviatrix.AWSGOV {
			return fmt.Errorf("insane_mode is only supported for AWS, GCP, AZURE, and AWSGOV (cloud_type = 1, 4, 8 or 256)")
		}
		if cloudType == goaviatrix.AWS || cloudType == goaviatrix.AWSGOV {
			if d.Get("insane_mode_az").(string) == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for AWS/AWSGOV cloud")
			}
			if d.Get("ha_subnet").(string) != "" && d.Get("ha_insane_mode_az").(string) == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for AWS/AWSGOV cloud and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			insaneModeAz := d.Get("insane_mode_az").(string)
			strs = append(strs, gateway.Subnet, insaneModeAz)
			gateway.Subnet = strings.Join(strs, "~~")
		}
		if cloudType == goaviatrix.GCP && !d.Get("enable_active_mesh").(bool) {
			return fmt.Errorf("insane_mode is supported for GCP provder only if active mesh 2.0 is enabled")
		}
		gateway.InsaneMode = "on"
	} else {
		gateway.InsaneMode = "off"
	}

	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haZone != "" && gateway.CloudType != goaviatrix.GCP && gateway.CloudType != goaviatrix.AZURE {
		return fmt.Errorf("'ha_zone' is only valid for GCP and AZURE providers when enabling HA")
	}
	if gateway.CloudType == goaviatrix.GCP && haSubnet != "" && haZone == "" {
		return fmt.Errorf("'ha_zone' must be set to enable HA on GCP, cannot enable HA with only 'ha_subnet'")
	}
	if gateway.CloudType == goaviatrix.AZURE && haSubnet == "" && haZone != "" {
		return fmt.Errorf("'ha_subnet' must be provided to enable HA on AZURE, cannot enable HA with only 'ha_zone'")
	}
	haGwSize := d.Get("ha_gw_size").(string)
	if haSubnet == "" && haZone == "" && haGwSize != "" {
		return fmt.Errorf("'ha_gw_size' is only required if enabling HA")
	}
	if haGwSize == "" && haSubnet != "" {
		return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
			"ha_subnet is set")
	}

	enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
	customerManagedKeys := d.Get("customer_managed_keys").(string)
	if enableEncryptVolume && d.Get("cloud_type").(int) != goaviatrix.AWS && d.Get("cloud_type").(int) != goaviatrix.AWSGOV {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS and AWSGOV providers")
	}
	if customerManagedKeys != "" {
		if !enableEncryptVolume {
			return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
		}
		gateway.EncVolume = "no"
	}
	if !enableEncryptVolume && (gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AWSGOV) {
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
	if enableTransitFireNet {
		if !intInSlice(gateway.CloudType, []int{goaviatrix.AWS, goaviatrix.AWSGOV, goaviatrix.GCP, goaviatrix.AZURE}) {
			return fmt.Errorf("'enable_transit_firenet' is only supported in AWS, AWSGOV, GCP and AZURE providers")
		}
		if intInSlice(gateway.CloudType, []int{goaviatrix.AZURE, goaviatrix.GCP}) {
			gateway.EnableTransitFireNet = "on"
		}
		if gateway.CloudType == goaviatrix.GCP {
			if lanVpcID == "" || lanPrivateSubnet == "" {
				return fmt.Errorf("'lan_vpc_id' and 'lan_private_subnet' are required when 'cloud_type' = 4 (GCP) and 'enable_transit_firenet' = true")
			}
			gateway.LanVpcID = lanVpcID
			gateway.LanPrivateSubnet = lanPrivateSubnet
		}
	}
	if (!enableTransitFireNet || gateway.CloudType != goaviatrix.GCP) && (lanVpcID != "" || lanPrivateSubnet != "") {
		return fmt.Errorf("'lan_vpc_id' and 'lan_private_subnet' are only valid when 'cloud_type' = 4 (GCP) and 'enable_transit_firenet' = true")
	}
	if enableGatewayLoadBalancer && !enableFireNet && !enableTransitFireNet {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
	}
	if enableGatewayLoadBalancer && gateway.CloudType != goaviatrix.AWS {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'cloud_type' = 1 (AWS)")
	}

	enableEgressTransitFireNet := d.Get("enable_egress_transit_firenet").(bool)
	if enableEgressTransitFireNet && gateway.CloudType != goaviatrix.AWS && gateway.CloudType != goaviatrix.AZURE && gateway.CloudType != goaviatrix.AWSGOV {
		return fmt.Errorf("'enable_egress_transit_firenet' is only supported in AWS, AZURE and AWSGOV cloud providers")
	}
	if enableEgressTransitFireNet && !enableTransitFireNet {
		return fmt.Errorf("'enable_egress_transit_firenet' requires 'enable_transit_firenet' to be set to true")
	}
	if enableEgressTransitFireNet && connectedTransit {
		return fmt.Errorf("'enable_egress_transit_firenet' requires 'connected_transit' to be set to false")
	}

	learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if learnedCidrsApproval {
		gateway.LearnedCidrsApproval = "on"
	}

	if learnedCidrsApproval && d.Get("learned_cidrs_approval_mode").(string) == "connection" {
		return fmt.Errorf("'enable_learned_cidrs_approval' must be false if 'learned_cidrs_approval_mode' is set to 'connection'")
	}

	enableMonitorSubnets := d.Get("enable_monitor_gateway_subnets").(bool)
	var excludedInstances []string
	for _, v := range d.Get("monitor_exclude_list").(*schema.Set).List() {
		excludedInstances = append(excludedInstances, v.(string))
	}
	if enableMonitorSubnets && cloudType != goaviatrix.AWS && cloudType != goaviatrix.AWSGOV {
		return fmt.Errorf("'enable_monitor_gateway_subnets' is only valid for cloud_type = 1 (AWS) or 256 (AWSGOV)")
	}
	if !enableMonitorSubnets && len(excludedInstances) != 0 {
		return fmt.Errorf("'monitor_exclude_list' must be empty if 'enable_monitor_gateway_subnets' is false")
	}

	bgpOverLan := d.Get("enable_bgp_over_lan").(bool)
	if bgpOverLan && cloudType != goaviatrix.AZURE {
		return fmt.Errorf("'enable_bgp_over_lan' is only valid for cloud_type = 8 (AZURE)")
	}
	if bgpOverLan {
		gateway.BgpOverLan = "on"
	}

	oobManagementSubnet := d.Get("oob_management_subnet").(string)
	oobAvailabilityZone := d.Get("oob_availability_zone").(string)
	haOobManagementSubnet := d.Get("ha_oob_management_subnet").(string)
	haOobAvailabilityZone := d.Get("ha_oob_availability_zone").(string)

	if enablePrivateOob {
		if cloudType != goaviatrix.AWS && cloudType != goaviatrix.AWSGOV {
			return fmt.Errorf("'enable_private_oob' is only valid for cloud_type = 1 (AWS) or 256 (AWSGOV)")
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
			return fmt.Errorf("\"oob_mangeemnt_sbunet\" must be empty if \"enable_private_oob\" is false")
		}

		if haOobAvailabilityZone != "" {
			return fmt.Errorf("\"ha_oob_availability_zone\" must be empty if \"enable_private_oob\" is false")
		}

		if haOobManagementSubnet != "" {
			return fmt.Errorf("\"ha_oob_management_sbunet\" must be empty if \"enable_private_oob\" is false")
		}
	}

	log.Printf("[INFO] Creating Aviatrix Transit Gateway: %#v", gateway)

	err := client.LaunchTransitVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit Gateway: %s", err)
	}

	d.SetId(gateway.GwName)

	flag := false
	defer resourceAviatrixTransitGatewayReadIfRequired(d, meta, &flag)

	if customerManagedKeys != "" && enableEncryptVolume {
		gwEncVolume := &goaviatrix.Gateway{
			GwName:              d.Get("gw_name").(string),
			CustomerManagedKeys: d.Get("customer_managed_keys").(string),
		}
		err := client.EnableEncryptVolume(gwEncVolume)
		if err != nil {
			return fmt.Errorf("failed to enable encrypt gateway volume when creating transit gateway: %s due to %s", gwEncVolume.GwName, err)
		}
	}

	if enableActiveMesh := d.Get("enable_active_mesh").(bool); !enableActiveMesh {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		gw.EnableActiveMesh = "no"

		err := client.DisableActiveMesh(gw)
		if err != nil {
			return fmt.Errorf("couldn't disable Active Mode for Aviatrix Transit Gateway: %s", err)
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
		transitGateway := &goaviatrix.TransitVpc{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
			HASubnet:  haSubnet,
			Eip:       d.Get("ha_eip").(string),
		}

		if insaneMode && (transitGateway.CloudType == goaviatrix.AWS || transitGateway.CloudType == goaviatrix.AWSGOV) {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			haStrs = append(haStrs, haSubnet, insaneModeHaAz)
			haSubnet = strings.Join(haStrs, "~~")
			transitGateway.HASubnet = haSubnet
		}

		if transitGateway.CloudType == goaviatrix.GCP && haZone == "" {
			return fmt.Errorf("no ha_zone is provided for enabling Transit HA gateway: %s", transitGateway.GwName)
		} else if transitGateway.CloudType == goaviatrix.GCP {
			transitGateway.HAZone = haZone
			transitGateway.HASubnetGCP = haSubnet
		}

		if transitGateway.CloudType == goaviatrix.AZURE && haZone != "" {
			transitGateway.HASubnet = fmt.Sprintf("%s~~%s~~", haSubnet, haZone)
		}

		if enablePrivateOob {
			transitGateway.HASubnet = transitGateway.HASubnet + "~~" + haOobAvailabilityZone
			transitGateway.HAOobManagementSubnet = haOobManagementSubnet + "~~" + haOobAvailabilityZone
		}

		log.Printf("[INFO] Enabling HA on Transit Gateway: %#v", haSubnet)

		if transitGateway.CloudType == goaviatrix.GCP {
			err = client.EnableHaTransitGateway(transitGateway)
		} else {
			err = client.EnableHaTransitVpc(transitGateway)
		}
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
		}

		//Resize HA Gateway
		log.Printf("[INFO]Resizing Transit HA Gateway: %#v", haGwSize)

		if haGwSize != gateway.VpcSize {
			if haGwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet is set")
			}

			haGateway := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string) + "-hagw",
			}
			haGateway.GwSize = d.Get("ha_gw_size").(string)

			log.Printf("[INFO] Resizing Transit HA GAteway size to: %s ", haGateway.GwSize)

			err := client.UpdateGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
			}
		}
	}

	if _, ok := d.GetOk("tag_list"); ok {
		if cloudType != goaviatrix.AWS && cloudType != goaviatrix.AWSGOV && cloudType != goaviatrix.AZURE {
			return fmt.Errorf("'tag_list' is only supported for AWS/AWSGOV/AZURE providers")
		}
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
		tagListStr = goaviatrix.TagListStrColon(tagListStr)
		gateway.TagList = strings.Join(tagListStr, ",")
		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			TagList:      gateway.TagList,
			CloudType:    gateway.CloudType,
		}

		if tags.CloudType == goaviatrix.AZURE {
			err := client.AzureUpdateTags(tags)
			if err != nil {
				return fmt.Errorf("failed to add tags : %s", err)
			}
		} else {
			err := client.AddTags(tags)
			if err != nil {
				return fmt.Errorf("failed to add tags: %s", err)
			}
		}
	}

	enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
	if enableHybridConnection {
		if cloudType != goaviatrix.AWS && cloudType != goaviatrix.AWSGOV {
			return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS/AWSGOV providers")
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

	if enableNAT {
		gw := &goaviatrix.Gateway{
			GatewayName: gateway.GwName,
		}

		err := client.EnableSNat(gw)
		if err != nil {
			return fmt.Errorf("failed to enable SNAT: %s", err)
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
	if (d.Get("cloud_type").(int) == goaviatrix.AWS || d.Get("cloud_type").(int) == goaviatrix.AWSGOV) && enableVpcDnsServer {
		gwVpcDnsServer := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		log.Printf("[INFO] Enable VPC DNS Server: %#v", gwVpcDnsServer)

		err := client.EnableVpcDnsServer(gwVpcDnsServer)
		if err != nil {
			return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
		}
	} else if enableVpcDnsServer {
		return fmt.Errorf("'enable_vpc_dns_server' only supports AWS/AWSGOV providers")
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

	if enableTransitFireNet && (gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AWSGOV) {
		enableActiveMesh := d.Get("enable_active_mesh").(bool)
		if !enableActiveMesh {
			return fmt.Errorf("active_mesh needs to be enabled to enable transit firenet")
		}
		gwTransitFireNet := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		if enableGatewayLoadBalancer {
			err := client.EnableTransitFireNetWithGWLB(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to enable transit firenet with Gateway Load Balancer enabled: %v", err)
			}
		} else {
			err := client.EnableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		}
	}

	if val, ok := d.GetOk("bgp_polling_time"); ok {
		err := client.SetBgpPollingTime(gateway, val.(string))
		if err != nil {
			return fmt.Errorf("could not set bgp polling time: %v", err)
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
			return fmt.Errorf("could not set prepend_as_path: %v", err)
		}
	}

	if val, ok := d.GetOk("local_as_number"); ok {
		err := client.SetLocalASNumber(gateway, val.(string))
		if err != nil {
			return fmt.Errorf("could not set local_as_number: %v", err)
		}
	}

	if val, ok := d.GetOk("bgp_ecmp"); ok {
		err := client.SetBgpEcmp(gateway, val.(bool))
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp: %v", err)
		}
	}

	if d.Get("enable_segmentation").(bool) {
		if err := client.EnableSegmentation(gateway); err != nil {
			return fmt.Errorf("could not enable segmentation: %v", err)
		}
	}

	if enableEgressTransitFireNet {
		err := client.EnableEgressTransitFirenet(gateway)
		if err != nil {
			return fmt.Errorf("could not enable egress transit firenet: %v", err)
		}
	}

	enableActiveStandby := d.Get("enable_active_standby").(bool)
	if enableActiveStandby {
		if err := client.EnableActiveStandby(gateway); err != nil {
			return fmt.Errorf("could not enable Active Standby Mode: %v", err)
		}
	}

	approvalMode := d.Get("learned_cidrs_approval_mode").(string)
	if approvalMode != defaultLearnedCidrApprovalMode {
		err := client.SetTransitLearnedCIDRsApprovalMode(gateway, approvalMode)
		if err != nil {
			return fmt.Errorf("could not set learned CIDRs approval mode to %q: %v", approvalMode, err)
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
			return fmt.Errorf("could not enable monitor gateway subnets: %v", err)
		}
	}

	if !d.Get("enable_jumbo_frame").(bool) {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		err := client.DisableJumboFrame(gw)
		if err != nil {
			return fmt.Errorf("could not disable jumbo frame for transit gateway: %v", err)
		}
	}

	if holdTime := d.Get("bgp_hold_time").(int); holdTime != defaultBgpHoldTime {
		err := client.ChangeBgpHoldTime(gateway.GwName, holdTime)
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time after Transit Gateway creation: %v", err)
		}
	}

	if gateway.EnableSummarizeCidrToTgw {
		err = client.EnableSummarizeCidrToTgw(gateway.GwName)
		if err != nil {
			return fmt.Errorf("could not enable summarize cidr to tgw: %v", err)
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
		GwName:      d.Get("gw_name").(string),
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

	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("subnet", gw.VpcNet)

		if gw.CloudType == goaviatrix.AWS || gw.CloudType == goaviatrix.AWSGOV {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			if gw.AllocateNewEipRead && !gw.EnablePrivateOob {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == goaviatrix.GCP {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0])
			d.Set("vpc_reg", gw.GatewayZone)
			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == goaviatrix.AZURE || gw.CloudType == goaviatrix.OCI {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		}

		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("eip", gw.PublicIP)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			d.Set("single_ip_snat", true)
		} else {
			d.Set("single_ip_snat", false)
		}

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if gw.CloudType == goaviatrix.AWS || gw.CloudType == goaviatrix.AWSGOV {
			d.Set("enable_hybrid_connection", gw.EnableHybridConnection)
		} else {
			d.Set("enable_hybrid_connection", false)
		}

		if gw.ConnectedTransit == "yes" {
			d.Set("connected_transit", true)
		} else {
			d.Set("connected_transit", false)
		}

		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			if gw.CloudType == goaviatrix.AWS || gw.CloudType == goaviatrix.AWSGOV {
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

		d.Set("enable_private_oob", gw.EnablePrivateOob)
		if gw.EnablePrivateOob {
			d.Set("oob_management_subnet", strings.Split(gw.OobManagementSubnet, "~~")[0])
			d.Set("oob_availability_zone", gw.GatewayZone)
		}

		gwDetail, err := client.GetGatewayDetail(gw)
		if err != nil {
			return fmt.Errorf("couldn't get Aviatrix Transit Gateway: %s", err)
		}

		d.Set("enable_firenet", gwDetail.EnableFireNet)
		d.Set("enable_gateway_load_balancer", gwDetail.EnabledGatewayLoadBalancer)
		d.Set("enable_egress_transit_firenet", gwDetail.EnableEgressTransitFireNet)
		d.Set("customized_transit_vpc_routes", gwDetail.CustomizedTransitVpcRoutes)

		d.Set("enable_transit_firenet", gwDetail.EnableTransitFireNet)
		if gwDetail.EnableTransitFireNet && gw.CloudType == goaviatrix.GCP {
			d.Set("lan_vpc_id", gwDetail.BundleVpcInfo.LAN.VpcID)
			d.Set("lan_private_subnet", strings.Split(gwDetail.BundleVpcInfo.LAN.Subnet, "~~")[0])
		}

		if _, zoneIsSet := d.GetOk("zone"); gw.CloudType == goaviatrix.AZURE && (isImport || zoneIsSet) &&
			gwDetail.GwZone != "AvailabilitySet" {
			d.Set("zone", "az-"+gwDetail.GwZone)
		}

		if gw.EnableActiveMesh == "yes" {
			d.Set("enable_active_mesh", true)
		} else {
			d.Set("enable_active_mesh", false)
		}

		if (gw.CloudType == goaviatrix.AWS || gw.CloudType == goaviatrix.AWSGOV) && gw.EnableVpcDnsServer == "Enabled" {
			d.Set("enable_vpc_dns_server", true)
		} else {
			d.Set("enable_vpc_dns_server", false)
		}

		if gwDetail.EnableAdvertiseTransitCidr == "yes" {
			d.Set("enable_advertise_transit_cidr", true)
		} else {
			d.Set("enable_advertise_transit_cidr", false)
		}

		if gwDetail.LearnedCidrsApproval == "yes" {
			d.Set("enable_learned_cidrs_approval", true)
		} else {
			d.Set("enable_learned_cidrs_approval", false)
		}

		var bgpManualSpokeAdvertiseCidrs []string
		if _, ok := d.GetOk("bgp_manual_spoke_advertise_cidrs"); ok {
			bgpManualSpokeAdvertiseCidrs = strings.Split(d.Get("bgp_manual_spoke_advertise_cidrs").(string), ",")
		}
		if len(goaviatrix.Difference(bgpManualSpokeAdvertiseCidrs, gwDetail.BgpManualSpokeAdvertiseCidrs)) != 0 ||
			len(goaviatrix.Difference(gwDetail.BgpManualSpokeAdvertiseCidrs, bgpManualSpokeAdvertiseCidrs)) != 0 {
			bgpMSAN := ""
			for i := range gwDetail.BgpManualSpokeAdvertiseCidrs {
				if i == 0 {
					bgpMSAN = bgpMSAN + gwDetail.BgpManualSpokeAdvertiseCidrs[i]
				} else {
					bgpMSAN = bgpMSAN + "," + gwDetail.BgpManualSpokeAdvertiseCidrs[i]
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

		d.Set("enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
		if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
			return fmt.Errorf("setting 'monitor_exclude_list' to state: %v", err)
		}
	}

	if gw.CloudType == goaviatrix.AWS || gw.CloudType == goaviatrix.AWSGOV || gw.CloudType == goaviatrix.AZURE {
		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			CloudType:    gw.CloudType,
		}

		tagList, err := client.GetTags(tags)
		if err != nil {
			return fmt.Errorf("unable to read tag_list for gateway: %v due to %v", gateway.GwName, err)
		}

		var tagListStr []string
		if _, ok := d.GetOk("tag_list"); ok {
			tagList1 := d.Get("tag_list").([]interface{})
			tagListStr = goaviatrix.ExpandStringList(tagList1)
		}
		if len(goaviatrix.Difference(tagListStr, tagList)) != 0 || len(goaviatrix.Difference(tagList, tagListStr)) != 0 {
			if err := d.Set("tag_list", tagList); err != nil {
				log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
			}
		} else {
			if err := d.Set("tag_list", tagListStr); err != nil {
				log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
			}
		}
	}

	transitGateway := &goaviatrix.TransitVpc{GwName: gwName}
	advancedConfig, err := client.GetTransitGatewayAdvancedConfig(transitGateway)
	if err != nil {
		return fmt.Errorf("could not get advanced config: %v", err)
	}

	d.Set("bgp_hold_time", advancedConfig.BgpHoldTime)
	d.Set("bgp_polling_time", advancedConfig.BgpPollingTime)
	err = d.Set("prepend_as_path", advancedConfig.PrependASPath)
	if err != nil {
		return fmt.Errorf("could not set prepend_as_path: %v", err)
	}
	d.Set("local_as_number", advancedConfig.LocalASNumber)
	d.Set("bgp_ecmp", advancedConfig.BgpEcmpEnabled)
	d.Set("enable_active_standby", advancedConfig.ActiveStandbyEnabled)
	if gw.CloudType == goaviatrix.AZURE {
		d.Set("enable_bgp_over_lan", advancedConfig.TunnelAddrLocal != "")
	} else {
		d.Set("enable_bgp_over_lan", false)
	}
	d.Set("enable_transit_summarize_cidr_to_tgw", advancedConfig.EnableSummarizeCidrToTgw)

	isSegmentationEnabled, err := client.IsSegmentationEnabled(transitGateway)
	if err != nil {
		return fmt.Errorf("could not read if segmentation is enabled: %v", err)
	}
	d.Set("enable_segmentation", isSegmentationEnabled)

	d.Set("learned_cidrs_approval_mode", advancedConfig.LearnedCIDRsApprovalMode)

	jumboFrameStatus, err := client.GetJumboFrameStatus(gateway)
	if err != nil {
		return fmt.Errorf("could not get jumbo frame status for transit gateway: %v", err)
	}
	d.Set("enable_jumbo_frame", jumboFrameStatus)

	haGateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string) + "-hagw",
	}
	haGw, err := client.GetGateway(haGateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.Set("ha_gw_size", "")
			d.Set("ha_subnet", "")
			d.Set("ha_zone", "")
			d.Set("ha_insane_mode_az", "")
			d.Set("ha_eip", "")
			d.Set("ha_oob_management_subnet", "")
			d.Set("ha_oob_availability_zone", "")
		} else {
			return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway: %s", err)
		}
	} else {
		if haGw.CloudType == goaviatrix.AWS || haGw.CloudType == goaviatrix.AZURE || haGw.CloudType == goaviatrix.OCI || haGw.CloudType == goaviatrix.AWSGOV {
			d.Set("ha_subnet", haGw.VpcNet)
			if zone := d.Get("ha_zone"); haGw.CloudType == goaviatrix.AZURE && (isImport || zone.(string) != "") {
				haGwDetail, err := client.GetGatewayDetail(haGateway)
				if err != nil {
					return fmt.Errorf("could not get HA transit gateway details: %v", err)
				}
				if haGwDetail.GwZone != "AvailabilitySet" {
					d.Set("ha_zone", "az-"+haGwDetail.GwZone)
				} else {
					d.Set("ha_zone", "")
				}
			} else {
				d.Set("ha_zone", "")
			}
		} else if haGw.CloudType == goaviatrix.GCP {
			d.Set("ha_zone", haGw.GatewayZone)
			if d.Get("ha_subnet") != "" || isImport {
				d.Set("ha_subnet", haGw.VpcNet)
			}
		}
		d.Set("ha_eip", haGw.PublicIP)
		d.Set("ha_gw_size", haGw.GwSize)
		d.Set("ha_cloud_instance_id", haGw.CloudnGatewayInstID)
		d.Set("ha_gw_name", haGw.GwName)
		d.Set("ha_private_ip", haGw.PrivateIP)
		lanCidr, err := client.GetTransitGatewayLanCidr(haGw.GwName)
		if err != nil && err != goaviatrix.ErrNotFound {
			log.Printf("[WARN] Error getting lan cidr for HA transit gateway %s due to %s", haGw.GwName, err)
		}
		d.Set("ha_lan_interface_cidr", lanCidr)

		if haGw.EnablePrivateOob {
			d.Set("ha_oob_management_subnet", strings.Split(haGw.OobManagementSubnet, "~~")[0])
			d.Set("ha_oob_availability_zone", haGw.GatewayZone)
		}

		if haGw.InsaneMode == "yes" && (haGw.CloudType == goaviatrix.AWS || haGw.CloudType == goaviatrix.AWSGOV) {
			d.Set("ha_insane_mode_az", haGw.GatewayZone)
		} else {
			d.Set("ha_insane_mode_az", "")
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
	}
	log.Printf("[INFO] Updating Aviatrix Transit Gateway: %#v", gateway)

	d.Partial(true)
	if d.HasChange("ha_zone") {
		haZone := d.Get("ha_zone").(string)
		if haZone != "" && gateway.CloudType != goaviatrix.GCP && gateway.CloudType != goaviatrix.AZURE {
			return fmt.Errorf("'ha_zone' is only valid for GCP and AZURE providers when enabling HA")
		}
	}
	if d.HasChange("ha_zone") || d.HasChange("ha_subnet") {
		haZone := d.Get("ha_zone").(string)
		haSubnet := d.Get("ha_subnet").(string)
		if gateway.CloudType == goaviatrix.GCP && haSubnet != "" && haZone == "" {
			return fmt.Errorf("'ha_zone' must be set to enable HA on GCP, cannot enable HA with only 'ha_subnet'")
		}
		if gateway.CloudType == goaviatrix.AZURE && haSubnet == "" && haZone != "" {
			return fmt.Errorf("'ha_subnet' must be provided to enable HA on AZURE, cannot enable HA with only 'ha_zone'")
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
	if d.HasChange("lan_vpc_id") {
		return fmt.Errorf("updating lan_vpc_id is not allowed")
	}
	if d.HasChange("lan_private_subnet") {
		return fmt.Errorf("updating lan_private_subnet is not allowed")
	}

	if d.HasChange("enable_transit_firenet") && intInSlice(d.Get("cloud_type").(int), []int{goaviatrix.AZURE, goaviatrix.GCP}) {
		return fmt.Errorf("editing 'enable_transit_firenet' in GCP and AZURE is not supported")
	}
	if d.Get("enable_egress_transit_firenet").(bool) && !d.Get("enable_transit_firenet").(bool) {
		return fmt.Errorf("'enable_egress_transit_firenet' requires 'enable_transit_firenet' to be set to true")
	}
	if d.Get("enable_egress_transit_firenet").(bool) && gateway.CloudType != goaviatrix.AZURE && gateway.CloudType != goaviatrix.AWS && gateway.CloudType != goaviatrix.AWSGOV {
		return fmt.Errorf("'enable_egress_transit_firenet' is currently only supported on AWS, AZURE and AWSGOV cloud providers")
	}

	if d.Get("enable_learned_cidrs_approval").(bool) && d.Get("learned_cidrs_approval_mode").(string) == "connection" {
		return fmt.Errorf("'enable_learned_cidrs_approval' must be false if 'learned_cidrs_approval_mode' is set to 'connection'")
	}

	if d.HasChange("enable_private_oob") {
		return fmt.Errorf("updating enable_private_oob is not allowed")
	}

	enablePrivateOob := d.Get("enable_private_oob").(bool)

	if !enablePrivateOob {
		if d.HasChange("ha_oob_management_subnet") {
			return fmt.Errorf("updating ha_oob_manage_subnet is not allowed if private oob is disabled")
		}

		if d.HasChange("ha_oob_availability_zone") {
			return fmt.Errorf("updating ha_oob_availability_zone is not allowed if private oob is disabled")
		}
	}

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

		if singleAZGateway.SingleAZ == "enabled" {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA: %s", err)
			}
		} else if singleAZGateway.SingleAZ == "disabled" {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
			}
		}

	}

	primaryGwSize := d.Get("gw_size").(string)
	if d.HasChange("gw_size") {
		old, _ := d.GetChange("gw_size")
		primaryGwSize = old.(string)
		gateway.GwSize = d.Get("gw_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit Gateway: %s", err)
		}
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") ||
		(enablePrivateOob && (d.HasChange("ha_oob_management_subnet") || d.HasChange("ha_oob_availability_zone"))) {
		transitGw := &goaviatrix.TransitVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
		}

		if transitGw.CloudType == goaviatrix.AWS || transitGw.CloudType == goaviatrix.GCP || transitGw.CloudType == goaviatrix.AWSGOV {
			transitGw.Eip = d.Get("ha_eip").(string)
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}
		if d.Get("insane_mode").(bool) && (transitGw.CloudType == goaviatrix.AWS || transitGw.CloudType == goaviatrix.AWSGOV) {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			if insaneModeHaAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set")
			}
			haStrs = append(haStrs, transitGw.HASubnet, insaneModeHaAz)
			transitGw.HASubnet = strings.Join(haStrs, "~~")
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false
		if transitGw.CloudType == goaviatrix.AWS || transitGw.CloudType == goaviatrix.AZURE || transitGw.CloudType == goaviatrix.AWSGOV {
			transitGw.HASubnet = d.Get("ha_subnet").(string)
			if transitGw.CloudType == goaviatrix.AZURE && d.Get("ha_zone").(string) != "" {
				transitGw.HASubnet = fmt.Sprintf("%s~~%s~~", d.Get("ha_subnet").(string), d.Get("ha_zone").(string))
			}
			if !enablePrivateOob {
				if oldSubnet == "" && newSubnet != "" {
					newHaGwEnabled = true
				} else if oldSubnet != "" && newSubnet == "" {
					deleteHaGw = true
				} else if oldSubnet != "" && newSubnet != "" {
					changeHaGw = true
				} else if d.HasChange("ha_zone") {
					changeHaGw = true
				}
			} else {
				if oldSubnet == "" && newSubnet != "" {
					newHaGwEnabled = true
				} else if newSubnet == "" {
					deleteHaGw = true
				} else if oldSubnet != newSubnet || (oldSubnet == newSubnet && (d.HasChange("ha_oob_management_subnet") || d.HasChange("ha_oob_availability_zone"))) {
					changeHaGw = true
				}
			}
		} else if transitGw.CloudType == goaviatrix.GCP {
			transitGw.HAZone = d.Get("ha_zone").(string)
			transitGw.HASubnetGCP = d.Get("ha_subnet").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
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

				transitGw.HASubnet = transitGw.HASubnet + "~~" + haOobAvailabilityZone
				transitGw.HAOobManagementSubnet = haOobManagementSubnet + "~~" + haOobAvailabilityZone
			} else if deleteHaGw {
				if haOobAvailabilityZone != "" {
					return fmt.Errorf("\"ha_oob_availability_zone\" must be empty if \"ha_subnet\" is empty")
				}

				if haOobManagementSubnet != "" {
					return fmt.Errorf("\"ha_oob_mangeemnt_sbunet\" must be empty if \"ha_subnet\" is empty")
				}
			}
		}

		if newHaGwEnabled {
			if transitGw.CloudType == goaviatrix.GCP {
				err := client.EnableHaTransitGateway(transitGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
				}
			} else {
				err := client.EnableHaTransitVpc(transitGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
				}
			}
		} else if deleteHaGw {
			if d.Get("ha_gw_size").(string) != "" {
				return fmt.Errorf("\"ha_gw_size\" must be empty if transit HA gateway is deleted")
			}

			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Transit HA gateway: %s", err)
			}
		} else if changeHaGw {
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Transit HA gateway: %s", err)
			}

			transitGw.Eip = ""

			if transitGw.CloudType == goaviatrix.GCP {
				err := client.EnableHaTransitGateway(transitGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
				}
			} else {
				err := client.EnableHaTransitVpc(transitGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
				}
			}

			newHaGwEnabled = true
		}
	}

	if gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AWSGOV || gateway.CloudType == goaviatrix.AZURE {
		if d.HasChange("tag_list") {
			tags := &goaviatrix.Tags{
				ResourceType: "gw",
				ResourceName: d.Get("gw_name").(string),
				CloudType:    gateway.CloudType,
			}

			if tags.CloudType == goaviatrix.AZURE {
				tagList := goaviatrix.ExpandStringList(d.Get("tag_list").([]interface{}))
				tagList = goaviatrix.TagListStrColon(tagList)
				tags.TagList = strings.Join(tagList, ",")
				err := client.AzureUpdateTags(tags)
				if err != nil {
					return fmt.Errorf("failed to update tags : %s", err)
				}
			} else {
				o, n := d.GetChange("tag_list")
				if o == nil {
					o = new([]interface{})
				}
				if n == nil {
					n = new([]interface{})
				}
				os := o.([]interface{})
				ns := n.([]interface{})
				oldList := goaviatrix.ExpandStringList(os)
				newList := goaviatrix.ExpandStringList(ns)
				oldTagList := goaviatrix.Difference(oldList, newList)
				newTagList := goaviatrix.Difference(newList, oldList)
				if len(oldTagList) != 0 || len(newTagList) != 0 {
					if len(oldTagList) != 0 {
						oldTagList = goaviatrix.TagListStrColon(oldTagList)
						tags.TagList = strings.Join(oldTagList, ",")
						err := client.DeleteTags(tags)
						if err != nil {
							return fmt.Errorf("failed to delete tags : %s", err)
						}
					}
					if len(newTagList) != 0 {
						newTagList = goaviatrix.TagListStrColon(newTagList)
						tags.TagList = strings.Join(newTagList, ",")
						err := client.AddTags(tags)
						if err != nil {
							return fmt.Errorf("failed to add tags : %s", err)
						}
					}
				}
			}
		}
	} else {
		if d.HasChange("tag_list") {
			return fmt.Errorf("'tag_list' is only supported for AWS/AWSGOV/AZURE providers")
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
				if err == goaviatrix.ErrNotFound {
					d.Set("ha_gw_size", "")
					d.Set("ha_subnet", "")
					d.Set("ha_zone", "")
					d.Set("ha_insane_mode_az", "")
					return nil
				}
				return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway while trying to update HA Gw size: %s", err)
			}
			haGateway.GwSize = d.Get("ha_gw_size").(string)
			if haGateway.GwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set")
			}
			err = client.UpdateGateway(haGateway)
			log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
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

	if gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AWSGOV {
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
			return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS/AWSGOV providers")
		}
	}

	if d.HasChange("enable_active_mesh") {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		enableActiveMesh := d.Get("enable_active_mesh").(bool)
		if enableActiveMesh {
			gw.EnableActiveMesh = "yes"
			err := client.EnableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to enable Active Mesh Mode: %s", err)
			}
		} else {
			gw.EnableActiveMesh = "no"
			err := client.DisableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to disable Active Mesh Mode: %s", err)
			}
		}
	}

	if d.HasChange("learned_cidrs_approval_mode") && d.HasChange("enable_learned_cidrs_approval") {
		gw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		currentMode, _ := d.GetChange("learned_cidrs_approval_mode")
		// API calls need to be in a specific order depending on the current mode
		if currentMode.(string) == "gateway" {
			learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
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
			learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
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
		learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
		if learnedCidrsApproval {
			gw.LearnedCidrsApproval = "on"
			err := client.EnableTransitLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
			}
		} else {
			gw.LearnedCidrsApproval = "off"
			err := client.DisableTransitLearnedCidrsApproval(gw)
			if err != nil {
				return fmt.Errorf("failed to disable learned cidrs approval: %s", err)
			}
		}
	}

	enableFireNet := d.Get("enable_firenet").(bool)
	enableGatewayLoadBalancer := d.Get("enable_gateway_load_balancer").(bool)
	enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
	if enableGatewayLoadBalancer && !enableFireNet && !enableTransitFireNet {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'enable_firenet' or 'enable_transit_firenet' is set to true")
	}
	if enableGatewayLoadBalancer && gateway.CloudType != goaviatrix.AWS {
		return fmt.Errorf("'enable_gateway_load_balancer' is only valid when 'cloud_type' = 1 (AWS)")
	}
	if enableFireNet && enableTransitFireNet {
		return fmt.Errorf("can't enable firenet function and transit firenet function at the same time")
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
			enableActiveMesh := d.Get("enable_active_mesh").(bool)
			if !enableActiveMesh {
				return fmt.Errorf("active_mesh needs to be enabled to enable transit firenet")
			}
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
			enableActiveMesh := d.Get("enable_active_mesh").(bool)
			if !enableActiveMesh {
				return fmt.Errorf("active_mesh needs to be enabled to enable transit firenet")
			}
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

	if d.HasChange("enable_egress_transit_firenet") {
		enableEgressTransitFirenet := d.Get("enable_egress_transit_firenet").(bool)
		if enableEgressTransitFirenet {
			err := client.EnableEgressTransitFirenet(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not enable egress transit firenet: %v", err)
			}
		} else {
			err := client.DisableEgressTransitFirenet(&goaviatrix.TransitVpc{GwName: gateway.GwName})
			if err != nil {
				return fmt.Errorf("could not disable egress transit firenet: %v", err)
			}
		}
	}

	if d.HasChange("enable_vpc_dns_server") && (d.Get("cloud_type").(int) == goaviatrix.AWS || d.Get("cloud_type").(int) == goaviatrix.AWSGOV) {
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
		return fmt.Errorf("'enable_vpc_dns_server' only supports AWS/AWSGOV providers")
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
			if d.Get("cloud_type").(int) != goaviatrix.AWS && d.Get("cloud_type").(int) != goaviatrix.AWSGOV {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS/AWSGOV providers")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              d.Get("gw_name").(string),
				CustomerManagedKeys: d.Get("customer_managed_keys").(string),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
			}
		} else {
			return fmt.Errorf("can't disable Encrypt Volume for gateway: %s", gateway.GwName)
		}
	} else if d.HasChange("customer_managed_keys") {
		return fmt.Errorf("updating customer_managed_keys only is not allowed")
	}

	if d.HasChange("customized_spoke_vpc_routes") {
		o, n := d.GetChange("customized_spoke_vpc_routes")
		if o == nil {
			o = new(interface{})
		}
		if n == nil {
			n = new(interface{})
		}
		os := o.(interface{})
		ns := n.(interface{})
		oldRouteList := strings.Split(os.(string), ",")
		newRouteList := strings.Split(ns.(string), ",")
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
		if o == nil {
			o = new(interface{})
		}
		if n == nil {
			n = new(interface{})
		}
		os := o.(interface{})
		ns := n.(interface{})
		oldRouteList := strings.Split(os.(string), ",")
		newRouteList := strings.Split(ns.(string), ",")
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
		if o == nil {
			o = new(interface{})
		}
		if n == nil {
			n = new(interface{})
		}
		os := o.(interface{})
		ns := n.(interface{})
		oldRouteList := strings.Split(os.(string), ",")
		newRouteList := strings.Split(ns.(string), ",")
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
		bgpPollingTime := d.Get("bgp_polling_time").(string)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpPollingTime(gateway, bgpPollingTime)
		if err != nil {
			return fmt.Errorf("could not update bgp polling time: %v", err)
		}
	}

	if d.HasChange("prepend_as_path") {
		var prependASPath []string
		slice := d.Get("prepend_as_path").([]interface{})
		for _, v := range slice {
			prependASPath = append(prependASPath, v.(string))
		}
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetPrependASPath(gateway, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path: %v", err)
		}
	}

	if d.HasChange("local_as_number") {
		localAsNumber := d.Get("local_as_number").(string)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetLocalASNumber(gateway, localAsNumber)
		if err != nil {
			return fmt.Errorf("could not set local_as_number: %v", err)
		}
	}

	if d.HasChange("bgp_ecmp") {
		enabled := d.Get("bgp_ecmp").(bool)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		err := client.SetBgpEcmp(gateway, enabled)
		if err != nil {
			return fmt.Errorf("could not set bgp_ecmp: %v", err)
		}
	}

	if d.HasChange("enable_segmentation") {
		enabled := d.Get("enable_segmentation").(bool)
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		if enabled {
			if err := client.EnableSegmentation(gateway); err != nil {
				return fmt.Errorf("could not enable segmentation: %v", err)
			}
		} else {
			if err := client.DisableSegmentation(gateway); err != nil {
				return fmt.Errorf("could not disable segmentation: %v", err)
			}
		}
	}

	if d.HasChange("enable_active_standby") {
		gateway := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		if d.Get("enable_active_standby").(bool) {
			if err := client.EnableActiveStandby(gateway); err != nil {
				return fmt.Errorf("could not enable active standby mode: %v", err)
			}
		} else {
			if err := client.DisableActiveStandby(gateway); err != nil {
				return fmt.Errorf("could not disable active standby mode: %v", err)
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

	if d.HasChange("enable_jumbo_frame") {
		if d.Get("enable_jumbo_frame").(bool) {
			err := client.EnableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not enable jumbo frame for transit gateway when updating: %v", err)
			}
		} else {
			err := client.DisableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("could not disable jumbo frame for transit gateway when updating: %v", err)
			}
		}
	}

	if d.HasChange("bgp_hold_time") {
		err := client.ChangeBgpHoldTime(gateway.GwName, d.Get("bgp_hold_time").(int))
		if err != nil {
			return fmt.Errorf("could not change BGP Hold Time during Transit Gateway update: %v", err)
		}
	}

	if d.HasChange("enable_transit_summarize_cidr_to_tgw") {
		if d.Get("enable_transit_summarize_cidr_to_tgw").(bool) {
			err := client.EnableSummarizeCidrToTgw(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not enable summarize cidr to tgw when updating: %v", err)
			}
		} else {
			err := client.DisableSummarizeCidrToTgw(gateway.GwName)
			if err != nil {
				return fmt.Errorf("could not disable summarize cidr to tgw when updating: %v", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixTransitGatewayRead(d, meta)
}

func resourceAviatrixTransitGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Gateway: %#v", gateway)

	enableEgressTransitFirenet := d.Get("enable_egress_transit_firenet").(bool)
	if enableEgressTransitFirenet {
		err := client.DisableEgressTransitFirenet(&goaviatrix.TransitVpc{GwName: gateway.GwName})
		if err != nil {
			return fmt.Errorf("could not disable egress transit firenet: %v", err)
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
	if enableTransitFireNet && (gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AWSGOV) {
		err := client.DisableTransitFireNet(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable transit firenet for %s due to %s", gateway.GwName, err)
		}
	} else if enableTransitFireNet && gateway.CloudType == goaviatrix.AZURE {
		err := client.IsTransitFireNetReadyToBeDisabled(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable transit firenet for %s due to %s", gateway.GwName, err)
		}
	}

	//If HA is enabled, delete HA GW first.
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haSubnet != "" || haZone != "" {
		gateway.GwName += "-hagw"

		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix Transit Gateway HA gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transit Gateway: %s", err)
	}

	return nil
}
