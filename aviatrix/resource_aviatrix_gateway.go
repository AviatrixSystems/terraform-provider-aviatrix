package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewayCreate,
		Read:   resourceAviatrixGatewayRead,
		Update: resourceAviatrixGatewayUpdate,
		Delete: resourceAviatrixGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixGatewayResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixGatewayStateUpgradeV0,
				Version: 0,
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
				Description: "Account name. This account will be used to launch Aviatrix gateway.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Aviatrix gateway unique name.",
			},
			"vpc_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "ID of legacy VPC/Vnet to be connected.",
				DiffSuppressFunc: DiffSuppressFuncGatewayVpcId,
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region where this gateway will be launched.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size of Gateway Instance.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A VPC Network address range selected from one of the available network ranges.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Availability Zone. Only available for Azure (8), Azure GOV (32), Azure CHINA (2048) and Public Subnet Filtering gateway. Must be in the form 'az-n', for example, 'az-2'.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "AZ of subnet being created for Insane Mode Gateway. Required if insane_mode is set.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Source NAT for this container.",
			},
			"vpn_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable user access through VPN to this container.",
			},
			"vpn_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "VPN CIDR block for the container.",
			},
			"enable_elb": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable ELB or not.",
			},
			"split_tunnel": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specify split tunnel mode.",
			},
			"max_vpn_conn": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Computed:    true,
				Description: "Maximum connection of VPN access. Valid for VPN gateway only. If not set, '100' will be default value.",
			},
			"name_servers": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of DNS servers used to resolve domain names by " +
					"a connected VPN user when Split Tunnel Mode is enabled.",
			},
			"search_domains": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of domain names that will use the NameServer " +
					"when a specific name is not in the destination when Split Tunnel Mode is enabled.",
			},
			"additional_cidrs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of destination CIDR ranges that will also go through the VPN tunnel " +
					"when Split Tunnel Mode is enabled.",
			},
			"otp_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Two step authentication mode.",
			},
			"saml_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "This field indicates whether to enable SAML or not.",
			},
			"enable_vpn_nat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "This field indicates whether to enable VPN NAT or not. Only supported for VPN gateway. Valid values: true, false. Default value: true.",
			},
			"okta_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Token for Okta auth mode.",
			},
			"okta_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "URL for Okta auth mode.",
			},
			"okta_username_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Username suffix for Okta auth mode.",
			},
			"duo_integration_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Integration key for DUO auth mode.",
			},
			"duo_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "Secret key for DUO auth mode.",
			},
			"duo_api_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "API hostname for DUO auth mode.",
			},
			"duo_push_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Push mode for DUO auth.",
			},
			"enable_ldap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable LDAP or not.",
			},
			"ldap_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP server address. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_bind_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP bind DN. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "LDAP password. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_base_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP base DN. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_username_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP user attribute. Required: Yes if enable_ldap is 'yes'.",
			},
			"peering_ha_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. " +
					"Required to create peering ha gateway if cloud_type = 1 or 8 (AWS or Azure). Optional if cloud_type = 4 (GCP)",
			},
			"peering_ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Zone information for creating Peering HA Gateway. Required to create peering ha gateway if cloud_type = 4 (GCP). Optional for cloud_type = 8 (Azure).",
			},
			"peering_ha_insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Peering HA Gateway. Required if insane_mode is set.",
			},
			"peering_ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Peering HA Gateway Size.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true if this feature is desired.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "When value is false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable Insane Mode for Gateway. Valid values: true, false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Valid values: true, false.",
			},
			"enable_designated_gateway": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable 'designated_gateway' feature for Gateway. Valid values: true, false.",
			},
			"additional_cidrs_designated_gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.",
			},
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			"enable_monitor_gateway_subnets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable monitor gateway subnets. Valid values: true, false. Default value: false.",
			},
			"monitor_exclude_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A set of monitored instance ids. Only valid when 'enable_monitor_gateway_subnets' = true.",
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntAtLeast(301),
				Description:  "Typed value when modifying idle_timeout. If it's -1, this feature is disabled.",
			},
			"renegotiation_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntAtLeast(301),
				Description:  "Typed value when modifying renegotiation_interval. If it's -1, this feature is disabled.",
			},
			"fqdn_lan_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "FQDN gateway lan interface cidr.",
			},
			"fqdn_lan_vpc_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: DiffSuppressFuncGCPVpcId,
				Description:      "LAN VPC ID. Only used for GCP FQDN Gateway.",
			},
			"enable_public_subnet_filtering": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				RequiredWith: []string{
					"public_subnet_filtering_route_tables",
					"public_subnet_filtering_guard_duty_enforced",
				},
				ConflictsWith: conflictingPublicSubnetFilteringGatewayConfigKeys,
				Description:   "Create a [Public Subnet Filtering gateway](https://docs.aviatrix.com/HowTos/public_subnet_filtering_faq.html).",
			},
			"public_subnet_filtering_route_tables": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Route tables whose associated public subnets are protected. Required when `enable_public_subnet_filtering` attribute is true.",
			},
			"public_subnet_filtering_ha_route_tables": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Route tables whose associated public subnets are protected for the HA PSF gateway. Required when enable_public_subnet_filtering and peering_ha_subnet are set.",
			},
			"public_subnet_filtering_guard_duty_enforced": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Whether to enforce Guard Duty IP blocking. Required when `enable_public_subnet_filtering` attribute is true. Valid values: true or false. Default value: true.",
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Enable jumbo frame support for Gateway. Valid values: true or false. Default value: true.",
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
				Description: "A map of tags to assign to the gateway.",
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
			"peering_ha_availability_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Peering HA availability domain for OCI.",
			},
			"peering_ha_fault_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Peering HA fault domain for OCI.",
			},
			"eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required when allocate_new_eip is false. It uses specified EIP for this gateway.",
			},
			"peering_ha_eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA peering instance.",
			},
			"azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to this Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"peering_ha_azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to the Peering HA Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"elb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A name for the ELB that is created.",
			},
			"vpn_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
				Description: "Elb protocol for VPN gateway with elb enabled. Only supports AWS provider. " +
					"Valid values: 'TCP', 'UDP'. If not specified, 'TCP'' will be used.",
			},
			"tunnel_detection_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(20, 600),
				Description:  "The IPSec tunnel down detection time for the Gateway.",
			},
			"software_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "software_version can be used to set the desired software version of the gateway. " +
					"If set, we will attempt to update the gateway to the specified version. " +
					"If left blank, the gateway software version will continue to be managed through the aviatrix_controller_config resource.",
			},
			"peering_ha_software_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "peering_ha_software_version can be used to set the desired software version of the HA gateway. " +
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
			"peering_ha_image_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "peering_ha_image_version can be used to set the desired image version of the HA gateway. " +
					"If set, we will attempt to update the gateway to the specified version.",
			},
			"elb_dns_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ELB DNS Name.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the gateway.",
			},
			"peering_ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peering HA security group used for the gateway.",
			},
			"public_dns_server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NS server used by the gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the gateway.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the Gateway created.",
			},
			"peering_ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the peering HA gateway.",
			},
			"peering_ha_gw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix gateway unique name of HA gateway.",
			},
			"peering_ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of HA gateway.",
			},
			"fqdn_lan_interface": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "FQDN gateway lan interface id.",
			},
		},
	}
}

func resourceAviatrixGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType:          getInt(d, "cloud_type"),
		GwName:             getString(d, "gw_name"),
		AccountName:        getString(d, "account_name"),
		VpcID:              getString(d, "vpc_id"),
		VpcNet:             getString(d, "subnet"),
		VpcSize:            getString(d, "gw_size"),
		VpnCidr:            getString(d, "vpn_cidr"),
		ElbName:            getString(d, "elb_name"),
		MaxConn:            getString(d, "max_vpn_conn"),
		OtpMode:            getString(d, "otp_mode"),
		OktaToken:          getString(d, "okta_token"),
		OktaURL:            getString(d, "okta_url"),
		OktaUsernameSuffix: getString(d, "okta_username_suffix"),
		DuoIntegrationKey:  getString(d, "duo_integration_key"),
		DuoSecretKey:       getString(d, "duo_secret_key"),
		DuoAPIHostname:     getString(d, "duo_api_hostname"),
		DuoPushMode:        getString(d, "duo_push_mode"),
		LdapServer:         getString(d, "ldap_server"),
		LdapBindDn:         getString(d, "ldap_bind_dn"),
		LdapPassword:       getString(d, "ldap_password"),
		LdapBaseDn:         getString(d, "ldap_base_dn"),
		LdapUserAttr:       getString(d, "ldap_username_attribute"),
		AdditionalCidrs:    getString(d, "additional_cidrs"),
		NameServers:        getString(d, "name_servers"),
		SearchDomains:      getString(d, "search_domains"),
		Eip:                getString(d, "eip"),
		SaveTemplate:       "no",
		AvailabilityDomain: getString(d, "availability_domain"),
		FaultDomain:        getString(d, "fault_domain"),
	}

	err := checkPublicSubnetFilteringConfig(d)
	if err != nil {
		return err
	}
	if getBool(d, "enable_public_subnet_filtering") {
		var routeTables []string
		for _, v := range getSet(d, "public_subnet_filtering_route_tables").List() {
			routeTables = append(routeTables, mustString(v))
		}
		gateway.RouteTable = strings.Join(routeTables, ",")
		gateway.VpcNet = fmt.Sprintf("%s~~%s", getString(d, "subnet"), getString(d, "zone"))
	}

	fqdnLanCidr := getString(d, "fqdn_lan_cidr")
	fqdnLanVpcID := getString(d, "fqdn_lan_vpc_id")
	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && fqdnLanVpcID != "" {
		return fmt.Errorf("attribute 'fqdn_lan_vpc_id' is only valid for GCP FQDN Gateways")
	}
	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) && fqdnLanCidr != "" {
		return fmt.Errorf("attribute 'fqdn_lan_cidr' is only valid for GCP and Azure FQDN Gateways")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		if (fqdnLanCidr != "" && fqdnLanVpcID == "") || (fqdnLanCidr == "" && fqdnLanVpcID != "") {
			return fmt.Errorf("to create a GCP FQDN gateway, both 'fqdn_lan_cidr' and 'fqdn_lan_vpc_id' must be set")
		}
		if fqdnLanCidr != "" && fqdnLanVpcID != "" {
			gateway.LanVpcID = fqdnLanVpcID
			gateway.LanPrivateSubnet = fqdnLanCidr
			gateway.CreateFQDNGateway = true
		}
	} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		if fqdnLanCidr != "" {
			gateway.FqdnLanCidr = fqdnLanCidr
		}
	}

	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && !getBool(d, "enable_public_subnet_filtering") && getString(d, "zone") != "" {
		return fmt.Errorf("attribute 'zone' is only valid for Azure, Azure GOV, Azure China or Public Subnet Filtering Gateways")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && getString(d, "zone") != "" {
		gateway.VpcNet = fmt.Sprintf("%s~~%s~~", getString(d, "subnet"), getString(d, "zone"))
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gateway.VpcRegion = getString(d, "vpc_reg")
	} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = getString(d, "vpc_reg")
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
	}

	singleIpNat := getBool(d, "single_ip_snat")
	if singleIpNat {
		gateway.EnableNat = "yes"
	} else {
		gateway.EnableNat = "no"
	}

	allocateNewEip := getBool(d, "allocate_new_eip")
	if allocateNewEip {
		gateway.AllocateNewEip = "on"
	} else {
		gateway.AllocateNewEip = "off"

		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			return fmt.Errorf("failed to create transit gateway: 'allocate_new_eip' can only be set to 'false' when cloud_type is AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048) or AWS Top Secret (16384)")
		}
		if _, ok := d.GetOk("eip"); !ok {
			return fmt.Errorf("failed to create gateway: 'eip' must be set when 'allocate_new_eip' is false")
		}
		azureEipName, azureEipNameOk := d.GetOk("azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
			if !azureEipNameOk {
				return fmt.Errorf("failed to create gateway: 'azure_eip_name_resource_group' must be set when 'allocate_new_eip' is false and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
			}
			gateway.Eip = fmt.Sprintf("%s:%s", mustString(azureEipName), getString(d, "eip"))
		} else {
			if azureEipNameOk {
				return fmt.Errorf("failed to create gateway: 'azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
			}
			gateway.Eip = getString(d, "eip")
		}
	}

	insaneMode := getBool(d, "insane_mode")
	if insaneMode {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("insane_mode is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWS China (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			if getString(d, "insane_mode_az") == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China(1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}
			if getString(d, "peering_ha_subnet") != "" && getString(d, "peering_ha_insane_mode_az") == "" {
				return fmt.Errorf("peering_ha_insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS China(1024), AWS Top Secret (16384) or AWS Secret (32768) and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			insaneModeAz := getString(d, "insane_mode_az")
			strs = append(strs, gateway.VpcNet, insaneModeAz)
			gateway.VpcNet = strings.Join(strs, "~~")
		}
		gateway.InsaneMode = "yes"
	} else {
		gateway.InsaneMode = "no"
	}

	samlEnabled := getBool(d, "saml_enabled")
	if samlEnabled {
		gateway.SamlEnabled = "yes"
	} else {
		gateway.SamlEnabled = "no"
	}

	splitTunnel := getBool(d, "split_tunnel")
	if splitTunnel {
		gateway.SplitTunnel = "yes"
	} else {
		gateway.SplitTunnel = "no"
	}

	enableElb := getBool(d, "enable_elb")
	if enableElb {
		gateway.EnableElb = "yes"
	}

	gateway.EnableLdap = getBool(d, "enable_ldap")

	vpnStatus := getBool(d, "vpn_access")
	vpnProtocol := getString(d, "vpn_protocol")
	if vpnStatus {
		gateway.VpnStatus = "yes"

		if enableElb && goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			gateway.VpnProtocol = vpnProtocol
		} else if enableElb && vpnProtocol == "UDP" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'UDP' for VPN gateway with ELB is only supported by AWS provider")
		} else if !enableElb && vpnProtocol == "TCP" {
			return fmt.Errorf("'vpn_protocol' should be left empty or set to 'UDP' for vpn gateway of AWS provider without elb enabled")
		}

		if gateway.SamlEnabled == "yes" {
			if gateway.EnableLdap || gateway.OtpMode != "" {
				return fmt.Errorf("ldap and mfa can't be configured if saml is enabled")
			}
		}

		if gateway.OtpMode != "" && gateway.OtpMode != "2" && gateway.OtpMode != "3" {
			return fmt.Errorf("otp_mode can only be '2' or '3' or empty string")
		}

		if gateway.EnableLdap && gateway.OtpMode == "3" {
			return fmt.Errorf("ldap can't be configured along with okta authentication")
		}
		if gateway.EnableLdap {
			if gateway.LdapServer == "" {
				return fmt.Errorf("ldap server must be set if ldap is enabled")
			}
			if gateway.LdapBindDn == "" {
				return fmt.Errorf("ldap bind dn must be set if ldap is enabled")
			}
			if gateway.LdapPassword == "" {
				return fmt.Errorf("ldap password must be set if ldap is enabled")
			}
			if gateway.LdapBaseDn == "" {
				return fmt.Errorf("ldap base dn must be set if ldap is enabled")
			}
			if gateway.LdapUserAttr == "" {
				return fmt.Errorf("ldap user attribute must be set if ldap is enabled")
			}
		}
		if gateway.OtpMode == "2" {
			if gateway.DuoIntegrationKey == "" {
				return fmt.Errorf("duo integration key required if otp_mode set to 2")
			}
			if gateway.DuoSecretKey == "" {
				return fmt.Errorf("duo secret key required if otp_mode set to 2")
			}
			if gateway.DuoAPIHostname == "" {
				return fmt.Errorf("duo api hostname required if otp_mode set to 2")
			}
			if gateway.DuoPushMode != "auto" && gateway.DuoPushMode != "token" && gateway.DuoPushMode != "selective" {
				return fmt.Errorf("duo push mode must be set to a valid value (auto, selective, or token)")
			}
			gateway.AuthMethod = "DUO"
		} else if gateway.OtpMode == "3" {
			if gateway.OktaToken == "" {
				return fmt.Errorf("okta token must be set if otp_mode is set to 3")
			}
			if gateway.OktaURL == "" {
				return fmt.Errorf("okta url must be set if otp_mode is set to 3")
			}
			gateway.AuthMethod = "okta"
		}

	} else {
		gateway.VpnStatus = "no"
		if gateway.EnableElb == "yes" {
			return fmt.Errorf("can not enable elb without VPN access enabled")
		}
		if vpnProtocol != "" {
			return fmt.Errorf("'vpn_protocol' should be left empty for non-vpn gateway")
		}
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain == "" || gateway.FaultDomain == "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
	}
	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (gateway.AvailabilityDomain != "" || gateway.FaultDomain != "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
	}

	peeringHaGwSize := getString(d, "peering_ha_gw_size")
	peeringHaSubnet := getString(d, "peering_ha_subnet")
	peeringHaZone := getString(d, "peering_ha_zone")
	peeringHaAvailabilityDomain := getString(d, "peering_ha_availability_domain")
	peeringHaFaultDomain := getString(d, "peering_ha_fault_domain")

	if peeringHaZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) && !getBool(d, "enable_public_subnet_filtering") {
		return fmt.Errorf("'peering_ha_zone' is only valid for GCP, Azure and Public Subnet Filtering Gateway if enabling Peering HA")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) && peeringHaZone == "" && peeringHaSubnet != "" {
		return fmt.Errorf("'peering_ha_zone' must be set to enable Peering HA on GCP, " +
			"cannot enable Peering HA with only 'peering_ha_subnet' enabled")
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && peeringHaZone != "" && peeringHaSubnet == "" {
		return fmt.Errorf("'peering_ha_subnet' must be provided to enable HA on Azure, " +
			"cannot enable HA with only 'peering_ha_zone'")
	}
	if peeringHaSubnet == "" && peeringHaZone == "" && peeringHaGwSize != "" {
		return fmt.Errorf("'peering_ha_gw_size' is only required if enabling Peering HA")
	}
	if peeringHaSubnet != "" {
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (peeringHaAvailabilityDomain == "" || peeringHaFaultDomain == "") {
			return fmt.Errorf("'peering_ha_availability_domain' and 'peering_ha_fault_domain' are required to enable Peering HA on OCI")
		}
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (peeringHaAvailabilityDomain != "" || peeringHaFaultDomain != "") {
			return fmt.Errorf("'peering_ha_availability_domain' and 'peering_ha_fault_domain' are only valid for OCI")
		}
	}

	enableDesignatedGw := getBool(d, "enable_designated_gateway")
	if enableDesignatedGw {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'designated_gateway' feature is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
		}
		if peeringHaSubnet != "" || peeringHaZone != "" {
			return fmt.Errorf("can't enable HA for gateway with 'designated_gateway' enabled")
		}
		gateway.EnableDesignatedGateway = "true"
	}

	enableEncryptVolume := getBool(d, "enable_encrypt_volume")
	customerManagedKeys := getString(d, "customer_managed_keys")
	if enableEncryptVolume && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
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

	_, tagsOk := d.GetOk("tags")
	if tagsOk {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return errors.New("failed to create gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		tagsMap, err := extractTags(d, gateway.CloudType)
		if err != nil {
			return fmt.Errorf("error creating tags for gateway: %w", err)
		}
		tagJson, err := TagsMapToJson(tagsMap)
		if err != nil {
			return fmt.Errorf("failed to add tags when creating gateway: %w", err)
		}
		gateway.TagJson = tagJson
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
	if rxQueueSize != "" {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("rx_queue_size only supports AWS related cloud types")
		} else {
			gateway.RxQueueSize = rxQueueSize
		}
	}

	log.Printf("[INFO] Creating Aviatrix gateway: %#v", gateway)

	d.SetId(gateway.GwName)
	flag := false
	defer func() { _ = resourceAviatrixGatewayReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if getBool(d, "enable_public_subnet_filtering") {
		err := client.CreatePublicSubnetFilteringGateway(gateway)
		if err != nil {
			log.Printf("[INFO] failed to create public subnet filtering gateway: %#v", gateway)
			return fmt.Errorf("could not create public subnet filtering gateway: %w", err)
		}
		if !getBool(d, "public_subnet_filtering_guard_duty_enforced") {
			err = client.DisableGuardDutyEnforcement(gateway)
			if err != nil {
				return fmt.Errorf("could not disable guard duty enforcement for public subnet filtering gateway: %w", err)
			}
		}
	} else {
		err := client.CreateGateway(gateway)
		if err != nil {
			log.Printf("[INFO] failed to create Aviatrix gateway: %#v", gateway)
			return fmt.Errorf("failed to create Aviatrix gateway: %w", err)
		}
	}

	enableVpnNat := getBool(d, "enable_vpn_nat")
	if vpnStatus {
		if !enableVpnNat {
			err := client.DisableVpnNat(gateway)
			if err != nil {
				return fmt.Errorf("failed to disable VPN NAT: %w", err)
			}
		}
	} else if !enableVpnNat {
		return fmt.Errorf("'enable_vpc_nat' is only supported for vpn gateway. Can't modify it for non-vpn gateway")
	}

	singleAZ := getBool(d, "single_az_ha")
	if singleAZ && !getBool(d, "enable_public_subnet_filtering") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   getString(d, "gw_name"),
			SingleAZ: "yes",
		}

		log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

		err := client.EnableSingleAZGateway(singleAZGateway)
		if err != nil {
			return fmt.Errorf("failed to create single AZ GW HA: %w", err)
		}
	} else if !singleAZ && getBool(d, "enable_public_subnet_filtering") {
		// Public Subnet Filtering Gateways are created with single_az_ha=true by default.
		// Thus, if user set single_az_ha=false, we need to disable.
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   getString(d, "gw_name"),
			SingleAZ: "no",
		}
		err := client.DisableSingleAZGateway(singleAZGateway)
		if err != nil {
			return fmt.Errorf("failed to disable single AZ : %w", err)
		}
	}

	if enableDesignatedGw {
		additionalCidrsDesignatedGw := getString(d, "additional_cidrs_designated_gateway")
		if additionalCidrsDesignatedGw != "" {
			designatedGw := &goaviatrix.Gateway{
				GwName:                      getString(d, "gw_name"),
				AdditionalCidrsDesignatedGw: additionalCidrsDesignatedGw,
			}
			err := client.EditDesignatedGateway(designatedGw)
			if err != nil {
				return fmt.Errorf("failed to edit additional cidrs for 'designated_gateway' feature due to %w", err)
			}
		}
	}

	// peering_ha_subnet is for Peering HA Gateway. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if peeringHaSubnet != "" || peeringHaZone != "" {
		if peeringHaGwSize == "" && !getBool(d, "enable_public_subnet_filtering") {
			return fmt.Errorf("a valid non empty peering_ha_gw_size parameter is mandatory for " +
				"this resource if peering_ha_subnet or peering_ha_zone is set. Example: t2.micro")
		}
		peeringHaGateway := &goaviatrix.Gateway{
			Eip:       getString(d, "peering_ha_eip"),
			GwName:    getString(d, "gw_name"),
			CloudType: getInt(d, "cloud_type"),
		}

		if goaviatrix.IsCloudType(peeringHaGateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			peeringHaGateway.PeeringHASubnet = peeringHaSubnet
			if insaneMode {
				var peeringHaStrs []string
				peeringHaInsaneModeAz := getString(d, "peering_ha_insane_mode_az")
				peeringHaStrs = append(peeringHaStrs, peeringHaSubnet, peeringHaInsaneModeAz)
				peeringHaSubnet = strings.Join(peeringHaStrs, "~~")
				peeringHaGateway.PeeringHASubnet = peeringHaSubnet
			}
		} else if goaviatrix.IsCloudType(peeringHaGateway.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			peeringHaGateway.PeeringHASubnet = peeringHaSubnet
			peeringHaGateway.AvailabilityDomain = peeringHaAvailabilityDomain
			peeringHaGateway.FaultDomain = peeringHaFaultDomain
		} else if goaviatrix.IsCloudType(peeringHaGateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			peeringHaGateway.NewZone = peeringHaZone
			if peeringHaSubnet != "" {
				peeringHaGateway.NewSubnet = peeringHaSubnet
			}
		} else if goaviatrix.IsCloudType(peeringHaGateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			peeringHaGateway.PeeringHASubnet = peeringHaSubnet
			if peeringHaZone != "" {
				peeringHaGateway.PeeringHASubnet = fmt.Sprintf("%s~~%s~~", peeringHaSubnet, peeringHaZone)
			}
		} else if peeringHaGateway.CloudType == goaviatrix.AliCloud {
			peeringHaGateway.PeeringHASubnet = peeringHaSubnet
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("peering_ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(peeringHaGateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if peeringHaGateway.Eip != "" {
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				peeringHaGateway.Eip = fmt.Sprintf("%s:%s", mustString(haAzureEipName), peeringHaGateway.Eip)
			} else if haAzureEipNameOk {
				return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be empty when 'peering_ha_eip' is empty")
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if getBool(d, "enable_public_subnet_filtering") {
			log.Printf("[INFO] Enable public subnet filtering HA: %#v", peeringHaGateway)
			var haRouteTables []string
			for _, v := range getSet(d, "public_subnet_filtering_ha_route_tables").List() {
				haRouteTables = append(haRouteTables, mustString(v))
			}
			peeringHaGateway.RouteTable = strings.Join(haRouteTables, ",")
			peeringHaGateway.PeeringHASubnet = fmt.Sprintf("%s~~%s", peeringHaSubnet, peeringHaZone)
			err := client.EnablePublicSubnetFilteringHAGateway(peeringHaGateway)
			if err != nil {
				return fmt.Errorf("could not create public subnet filtering gateway HA: %w", err)
			}
		} else {
			log.Printf("[INFO] Enable peering HA: %#v", peeringHaGateway)
			err := client.EnablePeeringHaGateway(peeringHaGateway)
			if err != nil {
				return fmt.Errorf("failed to create peering HA: %w", err)
			}
		}

		log.Printf("[INFO] Resizing Peering HA Gateway: %#v", peeringHaGwSize)
		if peeringHaGwSize != gateway.VpcSize {
			if peeringHaGwSize == "" {
				return fmt.Errorf("a valid non empty peering_ha_gw_size parameter is mandatory for " +
					"this resource if peering_ha_subnet is set. Example: t2.micro")
			}
			peeringHaGateway := &goaviatrix.Gateway{
				CloudType: getInt(d, "cloud_type"),
				GwName:    getString(d, "gw_name") + "-hagw", // CHECK THE NAME of peering ha gateway in
				// controller, test out first. just assuming it has that suffix
			}
			peeringHaGateway.VpcSize = peeringHaGwSize
			err := client.UpdateGateway(peeringHaGateway)
			log.Printf("[INFO] Resizing Peering Ha Gateway size to: %s,", peeringHaGateway.VpcSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Peering HA Gateway size: %w", err)
			}
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
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if enableMonitorSubnets {
		log.Printf("[INFO] Enable Monitor Gateway Subnets")
		err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
		if err != nil {
			return fmt.Errorf("could not enable monitor gateway subnets: %w", err)
		}
	}

	gatewayServer := &goaviatrix.Gateway{
		VpcID: getString(d, "vpc_id"),
	}

	idleTimeoutValue := getInt(d, "idle_timeout")
	if idleTimeoutValue != -1 {
		if getBool(d, "enable_elb") {
			gwName := getString(d, "gw_name")

			gw, err := client.GetGateway(&goaviatrix.Gateway{
				AccountName: getString(d, "account_name"),
				GwName:      gwName,
			})
			if err != nil {
				return fmt.Errorf("couldn't find Aviatrix Gateway for idle timeout : %s. Error: %w", gwName, err)
			}
			gatewayServer.GwName = gw.ElbName
		} else {
			gatewayServer.GwName = getString(d, "gw_name")
		}
		enableVPNServer := &goaviatrix.VPNConfig{
			Name:  "Idle timeout",
			Value: strconv.Itoa(idleTimeoutValue),
		}
		log.Printf("[INFO] Enable Modify VPN Config (Idle Timeout)")
		err := client.EnableVPNConfig(gatewayServer, enableVPNServer)
		if err != nil {
			return fmt.Errorf("fail to enable idle timeout: %w", err)
		}
	}

	renegoIntervalValue := getInt(d, "renegotiation_interval")
	if renegoIntervalValue != -1 {
		if getBool(d, "enable_elb") {
			gwName := getString(d, "gw_name")

			gw, err := client.GetGateway(&goaviatrix.Gateway{
				AccountName: getString(d, "account_name"),
				GwName:      gwName,
			})
			if err != nil {
				return fmt.Errorf("couldn't find Aviatrix Gateway renegotiation interval : %s", gwName)
			}
			gatewayServer.GwName = gw.ElbName
		} else {
			gatewayServer.GwName = getString(d, "gw_name")
		}
		enableVPNServer := &goaviatrix.VPNConfig{
			Name:  "Renegotiation interval",
			Value: strconv.Itoa(renegoIntervalValue),
		}
		log.Printf("[INFO] Enable Modify VPN Config (Renegotiation Interval)")
		err := client.EnableVPNConfig(gatewayServer, enableVPNServer)
		if err != nil {
			return fmt.Errorf("fail to enable renegotiation interval: %w", err)
		}
	}

	if !getBool(d, "enable_jumbo_frame") {
		err := client.DisableJumboFrame(gateway)
		if err != nil {
			return fmt.Errorf("couldn't disable jumbo frames for Gateway: %w", err)
		}
	}

	if !getBool(d, "enable_gro_gso") {
		err := client.DisableGroGso(gateway)
		if err != nil {
			return fmt.Errorf("couldn't disable GRO/GSO on gateway: %w", err)
		}
	}

	if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
		err := client.ModifyTunnelDetectionTime(gateway.GwName, mustInt(detectionTime))
		if err != nil {
			return fmt.Errorf("could not set tunnel detection time during Gateway creation: %w", err)
		}
	}

	if getBool(d, "enable_public_subnet_filtering") && len(gateway.TagJson) > 0 {
		// Workaround for setting tags during creation of a Public Subnet Filtering Gateway in R2.21.1
		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: getString(d, "gw_name"),
			CloudType:    gateway.CloudType,
		}
		if len(gateway.TagJson) > 0 {
			tags.TagJson = gateway.TagJson
			err := client.UpdateTags(tags)
			if err != nil {
				return fmt.Errorf("failed to set tags for gateway during creation: %w", err)
			}
		}
	}

	if rxQueueSize != "" {
		err := client.SetRxQueueSize(gateway)
		if err != nil {
			return fmt.Errorf("failed to set rx queue size for gateway %s: %w", gateway.GwName, err)
		}
		if peeringHaSubnet != "" || peeringHaZone != "" {
			haGwRxQueueSize := &goaviatrix.Gateway{
				GwName:      getString(d, "gw_name") + "-hagw",
				RxQueueSize: rxQueueSize,
			}
			err := client.SetRxQueueSize(haGwRxQueueSize)
			if err != nil {
				return fmt.Errorf("failed to set rx queue size for gateway ha %s : %w", haGwRxQueueSize.GwName, err)
			}
		}
	}

	return resourceAviatrixGatewayReadIfRequired(d, meta, &flag)
}

func resourceAviatrixGatewayReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixGatewayRead(d, meta)
	}
	return nil
}

func resourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	ignoreTagsConfig := client.IgnoreTagsConfig

	var isImport bool
	gwName := getString(d, "gw_name")
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		mustSet(d, "gw_name", id)
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
		return fmt.Errorf("couldn't find Aviatrix Gateway %s: %w", gwName, err)
	}

	log.Printf("[TRACE] reading gateway %s: %#v", getString(d, "gw_name"), gw)
	mustSet(d, "cloud_type", gw.CloudType)
	mustSet(d, "account_name", gw.AccountName)
	mustSet(d, "gw_name", gw.GwName)
	mustSet(d, "subnet", gw.VpcNet)
	mustSet(d, "single_ip_snat", gw.EnableNat == "yes" && gw.SnatMode == "primary")
	mustSet(d, "enable_ldap", gw.EnableLdap)
	mustSet(d, "vpn_cidr", gw.VpnCidr)
	mustSet(d, "saml_enabled", gw.SamlEnabled == "yes")
	mustSet(d, "okta_url", gw.OktaURL)
	mustSet(d, "okta_username_suffix", gw.OktaUsernameSuffix)
	mustSet(d, "duo_integration_key", gw.DuoIntegrationKey)
	mustSet(d, "duo_api_hostname", gw.DuoAPIHostname)
	mustSet(d, "duo_push_mode", gw.DuoPushMode)
	mustSet(d, "ldap_server", gw.LdapServer)
	mustSet(d, "ldap_bind_dn", gw.LdapBindDn)
	mustSet(d, "ldap_base_dn", gw.LdapBaseDn)
	mustSet(d, "ldap_username_attribute", gw.LdapUserAttr)
	mustSet(d, "single_az_ha", gw.SingleAZ == "yes")
	mustSet(d, "enable_encrypt_volume", gw.EnableEncryptVolume)
	mustSet(d, "eip", gw.PublicIP)
	mustSet(d, "cloud_instance_id", gw.CloudnGatewayInstID)
	mustSet(d, "public_dns_server", gw.PublicDnsServer)
	mustSet(d, "security_group_id", gw.GwSecurityGroupID)
	mustSet(d, "private_ip", gw.PrivateIP)
	mustSet(d, "enable_jumbo_frame", gw.JumboFrame)
	mustSet(d, "enable_vpc_dns_server", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled")
	mustSet(d, "tunnel_detection_time", gw.TunnelDetectionTime)
	mustSet(d, "image_version", gw.ImageVersion)
	mustSet(d, "software_version", gw.SoftwareVersion)
	mustSet(d, "rx_queue_size", gw.RxQueueSize)

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			mustSet(d, "azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Gateway %s", gw.GwName)
		}
	}

	if gw.IdleTimeout != "NA" {
		idleTimeout, err := strconv.Atoi(gw.IdleTimeout)
		if err != nil {
			return fmt.Errorf("couldn't get idle timeout for the gateway %s: %w", gw.GwName, err)
		}
		mustSet(d, "idle_timeout", idleTimeout)
	} else {
		mustSet(d, "idle_timeout", -1)
	}

	if gw.RenegotiationInterval != "NA" {
		renegotiationInterval, err := strconv.Atoi(gw.RenegotiationInterval)
		if err != nil {
			return fmt.Errorf("couldn't get renegotiation interval for the gateway %s: %w", gw.GwName, err)
		}
		mustSet(d, "renegotiation_interval", renegotiationInterval)
	} else {
		mustSet(d, "renegotiation_interval", -1)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		mustSet(
			// AWS vpc_id returns as <vpc_id>~~<other vpc info>
			d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
		mustSet(d, "vpc_reg", gw.VpcRegion)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(
			// gcp vpc_id returns as <vpc name>~-~<project name>
			d, "vpc_id", gw.VpcID)
		mustSet(d, "vpc_reg", gw.GatewayZone)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		mustSet(d, "vpc_id", gw.VpcID)
		mustSet(d, "vpc_reg", gw.VpcRegion)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		mustSet(d, "allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AliCloudRelatedCloudTypes) {
		mustSet(d, "allocate_new_eip", true)
	}

	if gw.EnableDesignatedGateway == "Yes" || gw.EnableDesignatedGateway == "yes" {
		mustSet(d, "enable_designated_gateway", true)
		cidrsTF := strings.Split(getString(d, "additional_cidrs_designated_gateway"), ",")
		cidrsRESTAPI := strings.Split(gw.AdditionalCidrsDesignatedGw, ",")
		if len(goaviatrix.Difference(cidrsTF, cidrsRESTAPI)) == 0 && len(goaviatrix.Difference(cidrsRESTAPI, cidrsTF)) == 0 {
			mustSet(d, "additional_cidrs_designated_gateway", getString(d, "additional_cidrs_designated_gateway"))
		} else {
			mustSet(d, "additional_cidrs_designated_gateway", gw.AdditionalCidrsDesignatedGw)
		}
	} else {
		mustSet(d, "enable_designated_gateway", false)
		mustSet(d, "additional_cidrs_designated_gateway", "")
	}

	_, zoneIsSet := d.GetOk("zone")
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || zoneIsSet) && gw.GatewayZone != "AvailabilitySet" {
		mustSet(d, "zone", "az-"+gw.GatewayZone)
	}

	if gw.VpnStatus != "" {
		if gw.VpnStatus == "disabled" {
			mustSet(d, "vpn_access", false)
			mustSet(d, "enable_vpn_nat", true)
			mustSet(d, "vpn_protocol", "")
			mustSet(d, "split_tunnel", true)
			mustSet(d, "max_vpn_conn", "")
		} else if gw.VpnStatus == "enabled" {
			mustSet(d, "vpn_access", true)
			mustSet(d, "split_tunnel", gw.SplitTunnel == "yes")
			mustSet(d, "max_vpn_conn", gw.MaxConn)
			mustSet(d, "enable_vpn_nat", gw.EnableVpnNat)
			if gw.ElbState == "enabled" {
				if strings.ToUpper(gw.VpnProtocol) == "UDP" {
					mustSet(d, "vpn_protocol", "UDP")
				} else {
					mustSet(d, "vpn_protocol", "TCP")
				}
			} else {
				mustSet(d, "vpn_protocol", "UDP")
			}
		}
	}

	if gw.ElbState == "enabled" {
		mustSet(d, "enable_elb", true)
		mustSet(d, "elb_name", gw.ElbName)
		mustSet(d, "elb_dns_name", gw.ElbDNSName)
	} else {
		mustSet(d, "enable_elb", false)
		mustSet(d, "elb_name", "")
	}

	if gw.AuthMethod == "duo_auth" || gw.AuthMethod == "duo_auth+LDAP" {
		mustSet(d, "otp_mode", "2")
	} else if gw.AuthMethod == "okta_auth" {
		mustSet(d, "otp_mode", "3")
	} else {
		mustSet(d, "otp_mode", "")
	}

	if gw.NewZone != "" {
		mustSet(d, "zone", gw.NewZone)
	}

	// Though go_aviatrix Gateway struct declares VpcSize as only used on gateway creation
	// it is the attribute receiving the instance size of an existing gateway instead of
	// GwSize. (at least in v3.5)
	if gw.GwSize != "" {
		mustSet(d, "gw_size", gw.GwSize)
	} else {
		if gw.VpcSize != "" {
			mustSet(d, "gw_size", gw.VpcSize)
		}
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

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		tags := goaviatrix.KeyValueTags(gw.Tags).IgnoreConfig(ignoreTagsConfig)
		if err := d.Set("tags", tags); err != nil {
			log.Printf("[WARN] Error setting tags for (%s): %s", d.Id(), err)
		}
	}

	if gw.VpnStatus == "enabled" && gw.SplitTunnel == "yes" {
		mustSet(d, "name_servers", gw.NameServers)
		mustSet(d, "search_domains", gw.SearchDomains)
		mustSet(d, "additional_cidrs", gw.AdditionalCidrs)
	} else {
		mustSet(d, "name_servers", "")
		mustSet(d, "search_domains", "")
		mustSet(d, "additional_cidrs", "")
	}
	mustSet(d, "enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
	if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
		return fmt.Errorf("setting 'monitor_exclude_list' to state: %w", err)
	}

	fqdnLanCidr, ok := gw.ArmFqdnLanCidr[gw.GwName]
	if ok && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		mustSet(d, "fqdn_lan_cidr", fqdnLanCidr)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "fqdn_lan_vpc_id", gw.BundleVpcInfo.LAN.VpcID)
		mustSet(d, "fqdn_lan_cidr", strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0])
	} else {
		mustSet(d, "fqdn_lan_cidr", "")
	}

	if !gw.IsPsfGateway {
		mustSet(d, "enable_public_subnet_filtering", false)
		mustSet(d, "public_subnet_filtering_route_tables", []string{})
		mustSet(d, "public_subnet_filtering_ha_route_tables", []string{})
		mustSet(d, "public_subnet_filtering_guard_duty_enforced", true)
	} else {
		mustSet(d, "enable_public_subnet_filtering", true)
		if err := d.Set("public_subnet_filtering_route_tables", gw.PsfDetails.RouteTableList); err != nil {
			return fmt.Errorf("could not set public_subnet_filtering_route_tables into state: %w", err)
		}
		mustSet(d, "public_subnet_filtering_guard_duty_enforced", gw.PsfDetails.GuardDutyEnforced == "yes")
		mustSet(d, "subnet", gw.PsfDetails.GwSubnetCidr)
		mustSet(d, "zone", gw.PsfDetails.GwSubnetAz)
		if gw.HaGw.GwSize == "" {
			err := d.Set("public_subnet_filtering_ha_route_tables", []string{})
			if err != nil {
				return fmt.Errorf("could not set public_subnet_filtering_ha_route_tables into state: %w", err)
			}
		} else {
			if err := d.Set("public_subnet_filtering_ha_route_tables", gw.PsfDetails.HaRouteTableList); err != nil {
				return fmt.Errorf("could not set public_subnet_filtering_ha_route_tables into state: %w", err)
			}
			mustSet(d, "peering_ha_subnet", gw.PsfDetails.HaGwSubnetCidr)
			mustSet(d, "peering_ha_zone", gw.PsfDetails.HaGwSubnetAz)
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

	enableGroGso, err := client.GetGroGsoStatus(gw)
	if err != nil {
		return fmt.Errorf("failed to get GRO/GSO status of gateway %s: %w", gw.GwName, err)
	}
	mustSet(d, "enable_gro_gso", enableGroGso)

	if gw.HaGw.GwSize == "" {
		mustSet(d, "peering_ha_availability_domain", "")
		mustSet(d, "peering_ha_azure_eip_name_resource_group", "")
		mustSet(d, "peering_ha_cloud_instance_id", "")
		mustSet(d, "peering_ha_eip", "")
		mustSet(d, "peering_ha_fault_domain", "")
		mustSet(d, "peering_ha_gw_name", "")
		mustSet(d, "peering_ha_gw_size", "")
		mustSet(d, "peering_ha_image_version", "")
		mustSet(d, "peering_ha_insane_mode_az", "")
		mustSet(d, "peering_ha_private_ip", "")
		mustSet(d, "peering_ha_security_group_id", "")
		mustSet(d, "peering_ha_software_version", "")
		mustSet(d, "peering_ha_subnet", "")
		mustSet(d, "peering_ha_zone", "")
		return nil
	}
	mustSet(d, "peering_ha_cloud_instance_id", gw.HaGw.CloudnGatewayInstID)
	mustSet(d, "peering_ha_gw_name", gw.HaGw.GwName)
	mustSet(d, "peering_ha_eip", gw.HaGw.PublicIP)
	mustSet(d, "peering_ha_gw_size", gw.HaGw.GwSize)
	mustSet(d, "peering_ha_private_ip", gw.HaGw.PrivateIP)
	mustSet(d, "peering_ha_software_version", gw.HaGw.SoftwareVersion)
	mustSet(d, "peering_ha_image_version", gw.HaGw.ImageVersion)
	mustSet(d, "peering_ha_security_group_id", gw.HaGw.GwSecurityGroupID)

	if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		if gw.HaGw.InsaneMode == "yes" {
			mustSet(d, "peering_ha_insane_mode_az", gw.HaGw.GatewayZone)
		}
	} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		if gw.HaGw.GatewayZone != "" {
			mustSet(d, "peering_ha_availability_domain", gw.HaGw.GatewayZone)
		} else {
			mustSet(d, "peering_ha_availability_domain", getString(d, "peering_ha_availability_domain"))
		}
		mustSet(d, "peering_ha_fault_domain", gw.HaGw.FaultDomain)
	} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
		if len(azureEip) == 3 {
			mustSet(d, "peering_ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Peering HA Gateway %s", gw.GwName)
		}
	}

	if !gw.IsPsfGateway {
		// For PSF gateway, peering_ha_subnet and peering_ha_zone are set above.
		// This block is only to set peering_ha_subnet and peering_ha_zone.
		if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "peering_ha_subnet", gw.HaGw.VpcNet)
			mustSet(d, "peering_ha_zone", "")
		} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			mustSet(d, "peering_ha_subnet", gw.HaGw.VpcNet)
			mustSet(d, "peering_ha_zone", "")
		} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			mustSet(d, "peering_ha_zone", gw.HaGw.GatewayZone)
			// only set peering_ha_subnet if the user has explicitly set it.
			if getString(d, "peering_ha_subnet") != "" || isImport {
				mustSet(d, "peering_ha_subnet", gw.HaGw.VpcNet)
			}
		} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			mustSet(d, "peering_ha_subnet", gw.HaGw.VpcNet)
			if _, haZoneIsSet := d.GetOk("peering_ha_zone"); isImport || haZoneIsSet {
				if gw.HaGw.GatewayZone != "AvailabilitySet" {
					mustSet(d, "peering_ha_zone", "az-"+gw.HaGw.GatewayZone)
				}
			}
		} else if gw.HaGw.CloudType == goaviatrix.AliCloud {
			mustSet(d, "peering_ha_subnet", gw.HaGw.VpcNet)
			mustSet(d, "peering_ha_zone", "")
		}
	}

	return nil
}

func resourceAviatrixGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", getString(d, "gw_name"))

	d.Partial(true)
	if d.HasChange("vpn_access") {
		return fmt.Errorf("updating vpn_access is not allowed")
	}
	if d.HasChange("enable_elb") {
		return fmt.Errorf("updating enable_elb is not allowed")
	}
	if d.HasChange("elb_name") {
		return fmt.Errorf("updating elb_name is not allowed")
	}
	if d.HasChange("vpn_protocol") {
		return fmt.Errorf("updating vpn_protocol is not allowed")
	}
	if d.HasChange("allocate_new_eip") {
		return fmt.Errorf("updating allocate_new_eip is not allowed")
	}
	if d.HasChange("eip") {
		return fmt.Errorf("updating eip is not allowed")
	}
	if d.HasChange("peering_ha_eip") {
		o, n := d.GetChange("peering_ha_eip")
		if o != "" && n != "" {
			return fmt.Errorf("updating peering_ha_eip is not allowed")
		}
	}
	if d.HasChange("azure_eip_name_resource_group") {
		return fmt.Errorf("failed to update gateway: changing 'azure_eip_name_resource_group' is not allowed")
	}
	if d.HasChange("peering_ha_azure_eip_name_resource_group") {
		o, n := d.GetChange("peering_ha_azure_eip_name_resource_group")
		if mustString(o) != "" && mustString(n) != "" {
			return fmt.Errorf("failed to update gateway: changing 'peering_ha_azure_eip_name_resource_group' is not allowed")
		}
	}
	if d.HasChange("enable_designated_gateway") {
		return fmt.Errorf("updating enable_designated_gateway is not allowed")
	}
	if d.HasChange("enable_public_subnet_filtering") {
		return fmt.Errorf("updating enable_public_subnet_filtering is not allowed")
	}
	err := checkPublicSubnetFilteringConfig(d)
	if err != nil {
		return err
	}

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
		VpcSize:   getString(d, "gw_size"),
	}
	vpnAccess := getBool(d, "vpn_access")
	enableElb := false
	geoVpnDnsName := ""
	if vpnAccess {
		enableElb = getBool(d, "enable_elb")
		if enableElb {
			gateway.ElbDNSName = getString(d, "elb_dns_name")
			geoVpn, err := client.GetGeoVPNName(gateway)
			if err == nil {
				geoVpnDnsName = geoVpn.ServiceName
			}
		}
	}
	if d.HasChange("peering_ha_zone") {
		peeringHaZone := getString(d, "peering_ha_zone")
		if peeringHaZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) && !getBool(d, "enable_public_subnet_filtering") {
			return fmt.Errorf("'peering_ha_zone' is only valid for GCP, Azure and Public Subnet Filtering Gateway if enabling Peering HA")
		}
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		peeringHaSubnet := getString(d, "peering_ha_subnet")
		peeringHaZone := getString(d, "peering_ha_zone")
		if peeringHaZone == "" && peeringHaSubnet != "" {
			return fmt.Errorf("'peering_ha_zone' must be set to enable Peering HA on GCP, " +
				"cannot enable Peering HA with only 'peering_ha_subnet' enabled")
		}
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		peeringHaSubnet := getString(d, "peering_ha_subnet")
		peeringHaZone := getString(d, "peering_ha_zone")
		if peeringHaZone != "" && peeringHaSubnet == "" {
			return fmt.Errorf("'peering_ha_subnet' must be set to enable Peering HA on Azure, " +
				"cannot enable Peering HA with only 'peering_ha_zone' enabled")
		}
	}

	singleAZ := getBool(d, "single_az_ha")
	if singleAZ {
		gateway.SingleAZ = "yes"
	} else {
		gateway.SingleAZ = "no"
	}

	peeringHaGateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name") + "-hagw",
		VpcSize:   getString(d, "peering_ha_gw_size"),
	}

	// Get primary gw size if gw_size changed, to be used later on for peering ha gw size update
	primaryGwSize := getString(d, "gw_size")
	if d.HasChange("gw_size") {
		old, _ := d.GetChange("gw_size")
		primaryGwSize = mustString(old)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Gateway: %w", err)
		}
	}

	if d.HasChange("otp_mode") || d.HasChange("enable_ldap") || d.HasChange("saml_enabled") ||
		d.HasChange("okta_token") || d.HasChange("okta_url") || d.HasChange("okta_username_suffix") ||
		d.HasChange("duo_integration_key") || d.HasChange("duo_secret_key") || d.HasChange("duo_api_hostname") ||
		d.HasChange("duo_push_mode") || d.HasChange("ldap_server") || d.HasChange("ldap_bind_dn") ||
		d.HasChange("ldap_password") || d.HasChange("ldap_base_dn") || d.HasChange("ldap_username_attribute") {

		if !vpnAccess {
			return fmt.Errorf("vpn_access must be set to yes to modify vpn authentication")
		}

		vpn_gw := &goaviatrix.VpnGatewayAuth{
			VpcID:              getString(d, "vpc_id"),
			OtpMode:            getString(d, "otp_mode"),
			OktaToken:          getString(d, "okta_token"),
			OktaURL:            getString(d, "okta_url"),
			OktaUsernameSuffix: getString(d, "okta_username_suffix"),
			DuoIntegrationKey:  getString(d, "duo_integration_key"),
			DuoSecretKey:       getString(d, "duo_secret_key"),
			DuoAPIHostname:     getString(d, "duo_api_hostname"),
			DuoPushMode:        getString(d, "duo_push_mode"),
			LdapServer:         getString(d, "ldap_server"),
			LdapBindDn:         getString(d, "ldap_bind_dn"),
			LdapPassword:       getString(d, "ldap_password"),
			LdapBaseDn:         getString(d, "ldap_base_dn"),
			LdapUserAttr:       getString(d, "ldap_username_attribute"),
		}

		samlEnabled := getBool(d, "saml_enabled")
		if samlEnabled {
			vpn_gw.SamlEnabled = "yes"
		} else {
			vpn_gw.SamlEnabled = "no"
		}

		vpn_gw.EnableLdap = getBool(d, "enable_ldap")

		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			// GCP vpn gw rest api call needs gcloud project id included in vpc id
			gw := &goaviatrix.Gateway{
				GwName: gateway.GwName,
			}
			gw1, err := client.GetGateway(gw)
			if err != nil {
				return fmt.Errorf("couldn't find Aviatrix Gateway: %s", gw.GwName)
			}
			vpn_gw.VpcID = gw1.VpcID
		}

		if vpn_gw.OtpMode != "" && vpn_gw.OtpMode != "2" && vpn_gw.OtpMode != "3" {
			return fmt.Errorf("otp_mode can only be '2' or '3' or empty string")
		}
		if vpn_gw.SamlEnabled == "yes" {
			if vpn_gw.EnableLdap || vpn_gw.OtpMode != "" {
				return fmt.Errorf("ldap and mfa can't be configured if saml is enabled")
			}
		}
		if vpn_gw.EnableLdap && vpn_gw.OtpMode == "3" {
			return fmt.Errorf("ldap can't be configured along with okta authentication")
		}
		if vpn_gw.EnableLdap {
			if vpn_gw.LdapServer == "" {
				return fmt.Errorf("ldap server must be set if ldap is enabled")
			}
			if vpn_gw.LdapBindDn == "" {
				return fmt.Errorf("ldap bind dn must be set if ldap is enabled")
			}
			if vpn_gw.LdapPassword == "" {
				return fmt.Errorf("ldap password must be set if ldap is enabled")
			}
			if vpn_gw.LdapBaseDn == "" {
				return fmt.Errorf("ldap base dn must be set if ldap is enabled")
			}
			if vpn_gw.LdapUserAttr == "" {
				return fmt.Errorf("ldap user attribute must be set if ldap is enabled")
			}
		}
		if vpn_gw.OtpMode == "2" {
			if vpn_gw.DuoIntegrationKey == "" {
				return fmt.Errorf("duo integration key required if otp_mode set to 2")
			}
			if vpn_gw.DuoSecretKey == "" {
				return fmt.Errorf("duo secret key required if otp_mode set to 2")
			}
			if vpn_gw.DuoAPIHostname == "" {
				return fmt.Errorf("duo api hostname required if otp_mode set to 2")
			}
			if vpn_gw.DuoPushMode != "auto" && vpn_gw.DuoPushMode != "token" && vpn_gw.DuoPushMode != "selective" {
				return fmt.Errorf("duo push mode must be set to a valid value (auto, selective, or token)")
			}
			if vpn_gw.EnableLdap {
				vpn_gw.AuthType = "duo_ldap_auth"
			} else {
				vpn_gw.AuthType = "duo_auth"
			}
		} else if vpn_gw.OtpMode == "3" {
			if vpn_gw.OktaToken == "" {
				return fmt.Errorf("okta token must be set if otp_mode is set to 3")
			}
			if vpn_gw.OktaURL == "" {
				return fmt.Errorf("okta url must be set if otp_mode is set to 3")
			}
			vpn_gw.AuthType = "okta_auth"
		} else {
			if vpn_gw.EnableLdap {
				vpn_gw.AuthType = "ldap_auth"
			} else if vpn_gw.SamlEnabled == "yes" {
				vpn_gw.AuthType = "saml_auth"
			} else {
				vpn_gw.AuthType = "none"
			}
		}
		if enableElb := getBool(d, "enable_elb"); enableElb {
			vpn_gw.LbOrGatewayName = getString(d, "elb_name")
		} else {
			vpn_gw.LbOrGatewayName = getString(d, "gw_name")
		}

		err := client.SetVpnGatewayAuthentication(vpn_gw)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix VPN Gateway Authentication: %w", err)
		}
	}

	if d.HasChange("tags") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("failed to update gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov(256) AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}

		tags := &goaviatrix.Tags{
			ResourceType: "gw",
			ResourceName: getString(d, "gw_name"),
			CloudType:    gateway.CloudType,
		}

		tagsMap, err := extractTags(d, gateway.CloudType)
		if err != nil {
			return fmt.Errorf("failed to update tags for gateway: %w", err)
		}
		tags.Tags = tagsMap
		tagJson, err := TagsMapToJson(tagsMap)
		if err != nil {
			return fmt.Errorf("failed to update tags for gateway: %w", err)
		}
		tags.TagJson = tagJson
		err = client.UpdateTags(tags)
		if err != nil {
			return fmt.Errorf("failed to update tags for gateway: %w", err)
		}
	}

	if d.HasChange("split_tunnel") || d.HasChange("additional_cidrs") ||
		d.HasChange("name_servers") || d.HasChange("search_domains") {
		splitTunnel := getBool(d, "split_tunnel")
		sTunnel := &goaviatrix.SplitTunnel{
			VpcID:   getString(d, "vpc_id"),
			ElbName: getString(d, "elb_name"),
		}
		if sTunnel.ElbName == "" {
			sTunnel.ElbName = getString(d, "gw_name")
		}

		if splitTunnel && (d.HasChange("additional_cidrs") || d.HasChange("name_servers") || d.HasChange("search_domains")) {
			sTunnel.AdditionalCidrs = getString(d, "additional_cidrs")
			sTunnel.NameServers = getString(d, "name_servers")
			sTunnel.SearchDomains = getString(d, "search_domains")
			sTunnel.SaveTemplate = "no"
			sTunnel.SplitTunnel = "yes"

			if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
				// ELB name is computed, search for gw to get elb name
				gw := &goaviatrix.Gateway{
					GwName: gateway.GwName,
				}

				gw1, err := client.GetGateway(gw)
				if err != nil {
					return fmt.Errorf("couldn't find Aviatrix Gateway: %s", gw.GwName)
				}
				if gw1.ElbState != "enabled" {
					sTunnel.ElbName = gw1.GwName
				} else {
					sTunnel.ElbName = gw1.ElbName
				}
				// VPC ID for gcp needs to include gcloud project ID
				sTunnel.VpcID = gw1.VpcID
			}

			err := client.ModifySplitTunnel(sTunnel)
			if err != nil {
				return fmt.Errorf("failed to modify split tunnel: %w", err)
			}
		} else if !splitTunnel && (getString(d, "additional_cidrs") != "" || getString(d, "name_servers") != "" || getString(d, "search_domains") != "") {
			return fmt.Errorf("to disable split_tunnel, following three attributes should be null: " +
				"'additional_cidrs', 'name_servers', and 'search_domains'")
		} else if !splitTunnel {
			sTunnel.SplitTunnel = "no"
			if vpnAccess && enableElb && geoVpnDnsName != "" {
				sTunnel.ElbName = geoVpnDnsName
				sTunnel.Dns = "true"
			}
			err := client.ModifySplitTunnel(sTunnel)
			if err != nil {
				return fmt.Errorf("failed to disable split tunnel: %w", err)
			}
		}
	}

	if d.HasChange("single_ip_snat") {
		gw := &goaviatrix.Gateway{
			CloudType:   getInt(d, "cloud_type"),
			GatewayName: getString(d, "gw_name"),
		}

		enableNat := getBool(d, "single_ip_snat")
		if enableNat {
			gw.EnableNat = "yes"
		} else {
			gw.EnableNat = "no"
		}

		if enableNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable SNAT: %w", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable SNAT: %w", err)
			}
		}

	}
	if d.HasChange("additional_cidrs_designated_gateway") {
		if !getBool(d, "enable_designated_gateway") {
			return fmt.Errorf("failed to edit additional cidrs for 'designated_gateway' since it is not enabled")
		}
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'designated_gateway' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		designatedGw := &goaviatrix.Gateway{
			GwName:                      getString(d, "gw_name"),
			AdditionalCidrsDesignatedGw: getString(d, "additional_cidrs_designated_gateway"),
		}
		err := client.EditDesignatedGateway(designatedGw)
		if err != nil {
			return fmt.Errorf("failed to edit additional cidrs for 'designated_gateway' feature due to %w", err)
		}
	}
	if d.HasChange("vpn_cidr") {
		if getBool(d, "vpn_access") {
			gw := &goaviatrix.Gateway{
				CloudType: getInt(d, "cloud_type"),
				GwName:    getString(d, "gw_name"),
				VpnCidr:   getString(d, "vpn_cidr"),
			}

			err := client.UpdateVpnCidr(gw)
			if err != nil {
				return fmt.Errorf("failed to update vpn cidr: %w", err)
			}
		} else {
			log.Printf("[INFO] can't update vpn cidr because vpn_access is disabled for gateway: %#v", gateway.GwName)
		}
	}
	if d.HasChange("max_vpn_conn") {
		if vpnAccess {
			gw := &goaviatrix.Gateway{
				CloudType: getInt(d, "cloud_type"),
				GwName:    getString(d, "gw_name"),
				VpcID:     getString(d, "vpc_id"),
				ElbName:   getString(d, "elb_name"),
			}

			if gw.ElbName == "" {
				gw.ElbName = getString(d, "gw_name")
			}

			_, n := d.GetChange("max_vpn_conn")
			gw.MaxConn = mustString(n)
			if enableElb && geoVpnDnsName != "" {
				gw.ElbName = geoVpnDnsName
				gw.Dns = "true"
			}
			err := client.UpdateMaxVpnConn(gw)
			if err != nil {
				return fmt.Errorf("failed to update max vpn connections: %w", err)
			}
		} else {
			log.Printf("[INFO] can't update max vpn connections because vpn is disabled for gateway: %#v", gateway.GwName)
		}
	}

	newHaGwEnabled := false
	if d.HasChange("peering_ha_subnet") || d.HasChange("peering_ha_zone") || d.HasChange("peering_ha_insane_mode_az") ||
		d.HasChange("peering_ha_availability_domain") || d.HasChange("peering_ha_fault_domain") {
		if getBool(d, "enable_designated_gateway") {
			return fmt.Errorf("can't update HA status for gateway with 'designated_gateway' enabled")
		}
		gw := &goaviatrix.Gateway{
			Eip:       getString(d, "peering_ha_eip"),
			GwName:    getString(d, "gw_name"),
			CloudType: getInt(d, "cloud_type"),
			VpcSize:   getString(d, "peering_ha_gw_size"),
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("peering_ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if gw.Eip != "" && gw.VpcSize != "" {
				// No change will be detected when peering_ha_eip is set to the empty string because it is computed.
				// Instead, check peering_ha_gw_size to detect when HA gateway is being deleted.
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				gw.Eip = fmt.Sprintf("%s:%s", mustString(haAzureEipName), gw.Eip)
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		oldSubnet, newSubnet := d.GetChange("peering_ha_subnet")
		oldZone, newZone := d.GetChange("peering_ha_zone")
		deleteHaGw := false
		changeHaGw := false

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			gw.PeeringHASubnet = getString(d, "peering_ha_subnet")
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && newZone != "" {
				gw.PeeringHASubnet = fmt.Sprintf("%s~~%s~~", newSubnet, newZone)
			}
			peeringHaAvailabilityDomain := getString(d, "peering_ha_availability_domain")
			peeringHaFaultDomain := getString(d, "peering_ha_fault_domain")
			if newSubnet != "" {
				if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (peeringHaAvailabilityDomain == "" || peeringHaFaultDomain == "") {
					return fmt.Errorf("'peering_ha_availability_domain' and 'peering_ha_fault_domain' are required to enable Peering HA on OCI")
				}
				if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.OCIRelatedCloudTypes) && (peeringHaAvailabilityDomain != "" || peeringHaFaultDomain != "") {
					return fmt.Errorf("'peering_ha_availability_domain' and 'peering_ha_fault_domain' are only valid for OCI")
				}
			}
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
				gw.AvailabilityDomain = peeringHaAvailabilityDomain
				gw.FaultDomain = peeringHaFaultDomain
			}
			if oldSubnet == "" && newSubnet != "" {
				newHaGwEnabled = true
			} else if oldSubnet != "" && newSubnet == "" {
				deleteHaGw = true
			} else if oldSubnet != "" && newSubnet != "" {
				changeHaGw = true
			} else if d.HasChange("peering_ha_zone") || d.HasChange("peering_ha_availability_domain") || d.HasChange("peering_ha_fault_domain") {
				changeHaGw = true
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			gw.NewZone = getString(d, "peering_ha_zone")
			gw.NewSubnet = getString(d, "peering_ha_subnet")
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}

		if getBool(d, "insane_mode") && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			var haStrs []string
			peeringHaInsaneModeAz := getString(d, "peering_ha_insane_mode_az")
			peeringHaSubnet := getString(d, "peering_ha_subnet")

			if peeringHaInsaneModeAz == "" && peeringHaSubnet != "" {
				return fmt.Errorf("peering_ha_insane_mode_az needed if insane_mode is enabled and peering_ha_subnet " +
					"is set for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			} else if peeringHaInsaneModeAz != "" && peeringHaSubnet == "" {
				return fmt.Errorf("peering_ha_subnet needed if insane_mode is enabled and peering_ha_insane_mode_az " +
					"is set for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}

			haStrs = append(haStrs, gw.PeeringHASubnet, peeringHaInsaneModeAz)
			gw.PeeringHASubnet = strings.Join(haStrs, "~~")
		}

		if (newHaGwEnabled || changeHaGw) && gw.VpcSize == "" {
			return fmt.Errorf("a valid non empty peering_ha_gw_size parameter is mandatory for this resource if " +
				"peering_ha_subnet or peering_ha_zone is set")
		} else if deleteHaGw && gw.VpcSize != "" {
			return fmt.Errorf("peering_ha_gw_size must be empty if transit HA gateway is deleted")
		}

		if getBool(d, "enable_public_subnet_filtering") {
			var haRouteTables []string
			for _, v := range getSet(d, "public_subnet_filtering_ha_route_tables").List() {
				haRouteTables = append(haRouteTables, mustString(v))
			}
			gw.RouteTable = strings.Join(haRouteTables, ",")
			gw.PeeringHASubnet = fmt.Sprintf("%s~~%s", getString(d, "peering_ha_subnet"), getString(d, "peering_ha_zone"))
			if newHaGwEnabled {
				err := client.EnablePublicSubnetFilteringHAGateway(gw)
				if err != nil {
					return fmt.Errorf("failed to enable Aviatrix public subnet filtering HA gateway: %w", err)
				}
			} else if deleteHaGw {
				err := client.DeletePublicSubnetFilteringGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix public subnet filtering HA gateway: %w", err)
				}
			} else if changeHaGw {
				err := client.DeletePublicSubnetFilteringGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix public subnet filtering HA gateway: %w", err)
				}

				gw.Eip = ""

				gateway.GwName = getString(d, "gw_name")
				err = client.EnablePublicSubnetFilteringHAGateway(gw)
				if err != nil {
					return fmt.Errorf("failed to enable Aviatrix public subnet filtering HA gateway: %w", err)
				}

				newHaGwEnabled = true
			}
		} else {
			if newHaGwEnabled {
				err := client.EnablePeeringHaGateway(gw)
				if err != nil {
					return fmt.Errorf("failed to enable Aviatrix peering HA gateway: %w", err)
				}
				if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
					if getString(d, "rx_queue_size") != "" && !d.HasChange("rx_queue_size") {
						haGwRxQueueSize := &goaviatrix.Gateway{
							GwName:      getString(d, "gw_name") + "-hagw",
							RxQueueSize: getString(d, "rx_queue_size"),
						}
						err := client.SetRxQueueSize(haGwRxQueueSize)
						if err != nil {
							return fmt.Errorf("could not set rx queue size for gateway ha: %s during gateway update: %w", haGwRxQueueSize.GwName, err)
						}
					}
				}
			} else if deleteHaGw {
				err := client.DeleteGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix peering HA gateway: %w", err)
				}
			} else if changeHaGw {
				err := client.DeleteGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix peering HA gateway: %w", err)
				}

				gw.Eip = ""

				gateway.GwName = getString(d, "gw_name")
				haErr := client.EnablePeeringHaGateway(gw)
				if haErr != nil {
					return fmt.Errorf("failed to enable Aviatrix peering HA gateway: %w", haErr)
				}

				newHaGwEnabled = true
			}
		}
	}
	haSubnet := getString(d, "peering_ha_subnet")
	haZone := getString(d, "peering_ha_zone")
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

			if haEnabled {
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

			if haEnabled {
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

	if d.HasChange("peering_ha_gw_size") || newHaGwEnabled {
		newHaGwSize := getString(d, "peering_ha_gw_size")
		if !newHaGwEnabled || (newHaGwSize != primaryGwSize) {
			// MODIFIES Peering HA GW SIZE if
			// Ha gateway wasn't newly configured
			// OR
			// newly configured peering HA gateway is set to be different size than primary gateway
			// (when peering ha gateway is enabled, it's size is by default the same as primary gateway)
			_, err := client.GetGateway(peeringHaGateway)
			if err != nil {
				if !errors.Is(err, goaviatrix.ErrNotFound) {
					return fmt.Errorf("couldn't find Aviatrix Peering HA Gateway while trying to update HA Gw "+
						"size: %s", err)
				}
			} else {
				if peeringHaGateway.VpcSize == "" {
					return fmt.Errorf("a valid non empty peering_ha_gw_size parameter is mandatory for this resource if " +
						"peering_ha_subnet or peering_ha_zone is set. Example: t2.micro or us-west1-b respectively")
				}
				err = client.UpdateGateway(peeringHaGateway)
				log.Printf("[INFO] Updating Peering HA Gateway size to: %s ", peeringHaGateway.VpcSize)
				if err != nil {
					return fmt.Errorf("failed to update Aviatrix Peering HA Gw size: %w", err)
				}
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
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192)")
	}

	if d.HasChange("enable_vpn_nat") {
		if !vpnAccess {
			return fmt.Errorf("'enable_vpc_nat' is only supported for vpn gateway. Can't updated it for Non VPN Gateway")
		} else {
			gw := &goaviatrix.Gateway{
				CloudType:    getInt(d, "cloud_type"),
				GwName:       getString(d, "gw_name"),
				VpcID:        getString(d, "vpc_id"),
				ElbName:      getString(d, "elb_name"),
				EnableVpnNat: true,
			}
			if enableElb && geoVpnDnsName != "" {
				gw.ElbName = geoVpnDnsName
				gw.Dns = "true"
			}

			if getBool(d, "enable_vpn_nat") {
				err := client.EnableVpnNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable VPN NAT: %w", err)
				}
			} else if !getBool(d, "enable_vpn_nat") {
				err := client.DisableVpnNat(gw)
				if err != nil {
					return fmt.Errorf("failed to disable VPN NAT: %w", err)
				}
			}
		}
	}

	if d.HasChange("enable_encrypt_volume") {
		if getBool(d, "enable_encrypt_volume") {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) provider")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              getString(d, "gw_name"),
				CustomerManagedKeys: getString(d, "customer_managed_keys"),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %w", gwEncVolume.GwName, err)
			}

			haSubnet := getString(d, "peering_ha_subnet")
			haZone := getString(d, "peering_ha_zone")
			haEnabled := haSubnet != "" || haZone != ""
			if haEnabled {
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

	gatewayServer := &goaviatrix.Gateway{
		VpcID: getString(d, "vpc_id"),
	}

	if d.HasChange("idle_timeout") {
		idleTimeoutValue := getInt(d, "idle_timeout")
		VPNServer := &goaviatrix.VPNConfig{
			Name: "Idle timeout",
		}
		if getBool(d, "enable_elb") {
			gatewayServer.GwName = getString(d, "elb_name")
		} else {
			gatewayServer.GwName = getString(d, "gw_name")
		}
		if idleTimeoutValue != -1 {
			VPNServer.Value = strconv.Itoa(idleTimeoutValue)
			log.Printf("[INFO] Modify VPN Config (update idle timeout value)")
			err := client.EnableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to update idle timeout value due to : %w", err)
			}
		} else {
			log.Printf("[INFO] Modify VPN Config (disable idle timeout)")
			err := client.DisableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to disable idle timeout due to : %w", err)
			}
		}
	}

	if d.HasChange("renegotiation_interval") {
		renegoIntervalValue := getInt(d, "renegotiation_interval")
		VPNServer := &goaviatrix.VPNConfig{
			Name: "Renegotiation interval",
		}
		if getBool(d, "enable_elb") {
			gatewayServer.GwName = getString(d, "elb_name")
		} else {
			gatewayServer.GwName = getString(d, "gw_name")
		}
		if renegoIntervalValue != -1 {
			VPNServer.Value = strconv.Itoa(renegoIntervalValue)
			log.Printf("[INFO] Modify VPN Config (update renegotiation interval value)")
			err := client.EnableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to enable renegotiation interval due to : %w", err)
			}
		} else {
			log.Printf("[INFO] Modify VPN Config (disable renegotiation interval)")
			err := client.DisableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to disable renegotiation interval due to: %w", err)
			}
		}
	}

	gatewayServer = &goaviatrix.Gateway{
		GwName: getString(d, "gw_name"),
	}

	if d.HasChange("public_subnet_filtering_route_tables") {
		var routeTables []string
		for _, v := range getSet(d, "public_subnet_filtering_route_tables").List() {
			routeTables = append(routeTables, mustString(v))
		}
		if len(routeTables) == 0 {
			return fmt.Errorf("attribute 'public_subnet_filtering_route_tables' must not be empty if 'enable_public_subnet_filtering' is set to true")
		}
		err := client.EditPublicSubnetFilteringRouteTableList(gatewayServer, routeTables)
		if err != nil {
			return fmt.Errorf("could not edit public subnet filtering route table rules: %w", err)
		}
	}
	if d.HasChange("public_subnet_filtering_ha_route_tables") && !d.HasChange("peering_ha_subnet") && getString(d, "peering_ha_subnet") != "" {
		var haRouteTables []string
		for _, v := range getSet(d, "public_subnet_filtering_ha_route_tables").List() {
			haRouteTables = append(haRouteTables, mustString(v))
		}
		peeringHaGateway.RouteTable = strings.Join(haRouteTables, ",")
		err := client.EditPublicSubnetFilteringRouteTableList(peeringHaGateway, haRouteTables)
		if err != nil {
			return fmt.Errorf("could not edit HA public subnet filtering route table rules: %w", err)
		}
	}
	if d.HasChange("public_subnet_filtering_guard_duty_enforced") {
		if getBool(d, "public_subnet_filtering_guard_duty_enforced") {
			err := client.EnableGuardDutyEnforcement(gatewayServer)
			if err != nil {
				return fmt.Errorf("could not enable public subnet filtering guard duty enforcement: %w", err)
			}
		} else {
			err := client.DisableGuardDutyEnforcement(gatewayServer)
			if err != nil {
				return fmt.Errorf("could not disable public subnet filtering guard duty enforcement: %w", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if getBool(d, "enable_jumbo_frame") {
			err := client.EnableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("couldn't enable jumbo frames for Gateway when updating: %w", err)
			}
		} else {
			err := client.DisableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("couldn't disable jumbo frames for Gateway when updating: %w", err)
			}
		}
	}

	if d.HasChange("enable_gro_gso") {
		if getBool(d, "enable_gro_gso") {
			err := client.EnableGroGso(gateway)
			if err != nil {
				return fmt.Errorf("couldn't enable GRO/GSO on gateway when updating: %w", err)
			}
		} else {
			err := client.DisableGroGso(gateway)
			if err != nil {
				return fmt.Errorf("couldn't disable GRO/GSO on gateway when updating: %w", err)
			}
		}
	}

	if d.HasChange("tunnel_detection_time") {
		detectionTimeInterface, ok := d.GetOk("tunnel_detection_time")
		var detectionTime int
		if ok {
			detectionTime = mustInt(detectionTimeInterface)
		} else {
			detectionTime, err = client.GetTunnelDetectionTime("Controller")
			if err != nil {
				return fmt.Errorf("could not get default tunnel detection time during Gateway update: %w", err)
			}
		}
		err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime)
		if err != nil {
			return fmt.Errorf("could not modify tunnel detection time during Gateway update: %w", err)
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
			return fmt.Errorf("could not modify rx queue size for gateway: %s during gateway update: %w", gw.GatewayName, err)
		}
		if haSubnet != "" || haZone != "" {
			haGwRxQueueSize := &goaviatrix.Gateway{
				GwName:      getString(d, "gw_name") + "-hagw",
				RxQueueSize: getString(d, "rx_queue_size"),
			}
			err := client.SetRxQueueSize(haGwRxQueueSize)
			if err != nil {
				return fmt.Errorf("could not modify rx queue size for gateway ha: %s during gateway update: %w", haGwRxQueueSize.GwName, err)
			}
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixGatewayRead(d, meta)
}

func resourceAviatrixGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}
	var err error
	isPublicSubnetFilteringGateway := getBool(d, "enable_public_subnet_filtering")
	// peering_ha_subnet is for Peering HA
	peeringHaSubnet := getString(d, "peering_ha_subnet")
	peeringHaZone := getString(d, "peering_ha_zone")
	if peeringHaSubnet != "" || peeringHaZone != "" {
		// Delete backup gateway first
		gateway.GwName += "-hagw"
		log.Printf("[INFO] Deleting Aviatrix Backup Gateway [-hagw]: %#v", gateway)

		if isPublicSubnetFilteringGateway {
			err = client.DeletePublicSubnetFilteringGateway(gateway)
		} else {
			err = client.DeleteGateway(gateway)
		}

		if err != nil {
			return fmt.Errorf("failed to delete backup [-hgw] gateway: %w", err)
		}
	}

	gateway.GwName = getString(d, "gw_name")

	log.Printf("[INFO] Deleting Aviatrix gateway: %#v", gateway)

	if isPublicSubnetFilteringGateway {
		err = client.DeletePublicSubnetFilteringGateway(gateway)
	} else {
		err = client.DeleteGateway(gateway)
	}
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Gateway: %w", err)
	}

	return nil
}

func checkPublicSubnetFilteringConfig(d *schema.ResourceData) error {
	var routeTables, haRouteTables []string
	for _, v := range getSet(d, "public_subnet_filtering_route_tables").List() {
		routeTables = append(routeTables, mustString(v))
	}
	for _, v := range getSet(d, "public_subnet_filtering_ha_route_tables").List() {
		haRouteTables = append(haRouteTables, mustString(v))
	}
	isPublicSubnetFilteringGw := getBool(d, "enable_public_subnet_filtering")
	// Public subnet filtering only supported for AWS and AWSGov
	if isPublicSubnetFilteringGw && !goaviatrix.IsCloudType(getInt(d, "cloud_type"), goaviatrix.AWS|goaviatrix.AWSGov) {
		return fmt.Errorf("enable_public_subnet_filtering is only valid for AWS (1) or AWSGov (256)")
	}
	if isPublicSubnetFilteringGw && len(routeTables) == 0 {
		return fmt.Errorf("public_subnet_filtering_route_tables can not be empty when 'enable_public_subnet_filtering' is enabled. Please supply at least one route table ID")
	}
	if !isPublicSubnetFilteringGw && len(routeTables) != 0 {
		return fmt.Errorf("use of public_subnet_filtering_route_tables is not valid if enable_public_subnet_filtering is false")
	}
	if !isPublicSubnetFilteringGw && len(haRouteTables) != 0 {
		return fmt.Errorf("use of public_subnet_filtering_ha_route_tables is not valid if enable_public_subnet_filtering is false")
	}
	if d.IsNewResource() {
		if getBool(d, "enable_public_subnet_filtering") && !getBool(d, "enable_encrypt_volume") {
			return fmt.Errorf("enable_encrypt_volume must be set to true when 'enable_public_subnet_filtering' is enabled")
		}
	}
	return nil
}

// Attributes that cannot be set when enabling public subnet filtering.
var conflictingPublicSubnetFilteringGatewayConfigKeys = []string{
	"additional_cidrs",
	"additional_cidrs_designated_gateway",
	"allocate_new_eip",
	"customer_managed_keys",
	"duo_api_hostname",
	"duo_integration_key",
	"duo_push_mode",
	"duo_secret_key",
	"eip",
	"elb_name",
	"enable_designated_gateway",
	"enable_elb",
	"enable_ldap",
	"enable_monitor_gateway_subnets",
	"enable_vpc_dns_server",
	"enable_vpn_nat",
	"fqdn_lan_cidr",
	"fqdn_lan_vpc_id",
	"idle_timeout",
	"insane_mode",
	"insane_mode_az",
	"ldap_base_dn",
	"ldap_bind_dn",
	"ldap_password",
	"ldap_server",
	"ldap_username_attribute",
	"max_vpn_conn",
	"monitor_exclude_list",
	"name_servers",
	"okta_token",
	"okta_url",
	"okta_username_suffix",
	"otp_mode",
	"peering_ha_eip",
	"peering_ha_insane_mode_az",
	"renegotiation_interval",
	"saml_enabled",
	"search_domains",
	"single_ip_snat",
	"split_tunnel",
	"vpn_access",
	"vpn_cidr",
	"vpn_protocol",
	"enable_jumbo_frame",
}
