package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewayCreate,
		Read:   resourceAviatrixGatewayRead,
		Update: resourceAviatrixGatewayUpdate,
		Delete: resourceAviatrixGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Description:  "Type of cloud service provider.",
				ValidateFunc: validateCloudType,
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account name. This account will be used to launch Aviatrix gateway.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Aviatrix gateway unique name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of legacy VPC/Vnet to be connected.",
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
				Description: "Availability Zone. Only available for Azure and Public Subnet Filtering gateway",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
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
				Default:     "",
				Description: "Maximum connection of VPN access.",
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
				Description: "Specify whether to enable LDAP or not. Supported values: 'yes' and 'no'.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Customer managed key ID.",
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
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "LAN VPC ID. Only used for GCP FQDN Gateway.",
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
			"tags": {
				Type:          schema.TypeMap,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Description:   "A map of tags to assign to the gateway.",
				ConflictsWith: []string{"tag_list"},
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of storage account with gateway images. Only valid for Azure China (2048)",
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
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType:          d.Get("cloud_type").(int),
		GwName:             d.Get("gw_name").(string),
		AccountName:        d.Get("account_name").(string),
		VpcID:              d.Get("vpc_id").(string),
		VpcNet:             d.Get("subnet").(string),
		VpcSize:            d.Get("gw_size").(string),
		VpnCidr:            d.Get("vpn_cidr").(string),
		ElbName:            d.Get("elb_name").(string),
		MaxConn:            d.Get("max_vpn_conn").(string),
		OtpMode:            d.Get("otp_mode").(string),
		OktaToken:          d.Get("okta_token").(string),
		OktaURL:            d.Get("okta_url").(string),
		OktaUsernameSuffix: d.Get("okta_username_suffix").(string),
		DuoIntegrationKey:  d.Get("duo_integration_key").(string),
		DuoSecretKey:       d.Get("duo_secret_key").(string),
		DuoAPIHostname:     d.Get("duo_api_hostname").(string),
		DuoPushMode:        d.Get("duo_push_mode").(string),
		LdapServer:         d.Get("ldap_server").(string),
		LdapBindDn:         d.Get("ldap_bind_dn").(string),
		LdapPassword:       d.Get("ldap_password").(string),
		LdapBaseDn:         d.Get("ldap_base_dn").(string),
		LdapUserAttr:       d.Get("ldap_username_attribute").(string),
		AdditionalCidrs:    d.Get("additional_cidrs").(string),
		NameServers:        d.Get("name_servers").(string),
		SearchDomains:      d.Get("search_domains").(string),
		Eip:                d.Get("eip").(string),
		SaveTemplate:       "no",
		AvailabilityDomain: d.Get("availability_domain").(string),
		FaultDomain:        d.Get("fault_domain").(string),
	}

	err := checkPublicSubnetFilteringConfig(d)
	if err != nil {
		return err
	}
	if d.Get("enable_public_subnet_filtering").(bool) {
		var routeTables []string
		for _, v := range d.Get("public_subnet_filtering_route_tables").(*schema.Set).List() {
			routeTables = append(routeTables, v.(string))
		}
		gateway.RouteTable = strings.Join(routeTables, ",")
		gateway.VpcNet = fmt.Sprintf("%s~~%s", d.Get("subnet").(string), d.Get("zone").(string))
	}

	fqdnLanCidr := d.Get("fqdn_lan_cidr").(string)
	fqdnLanVpcID := d.Get("fqdn_lan_vpc_id").(string)
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

	if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && !d.Get("enable_public_subnet_filtering").(bool) && d.Get("zone").(string) != "" {
		return fmt.Errorf("attribute 'zone' is only valid for Azure and Public Subnet Filtering Gateways")
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && d.Get("zone").(string) != "" {
		gateway.VpcNet = fmt.Sprintf("%s~~%s~~", d.Get("subnet").(string), d.Get("zone").(string))
	}

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		gateway.VpcRegion = d.Get("vpc_reg").(string)
	} else if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = d.Get("vpc_reg").(string)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) or AWS Secret (32768)")
	}

	singleIpNat := d.Get("single_ip_snat").(bool)
	if singleIpNat {
		gateway.EnableNat = "yes"
	} else {
		gateway.EnableNat = "no"
	}

	allocateNewEip := d.Get("allocate_new_eip").(bool)
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
			gateway.Eip = fmt.Sprintf("%s:%s", azureEipName.(string), d.Get("eip").(string))
		} else {
			if azureEipNameOk {
				return fmt.Errorf("failed to create gateway: 'azure_eip_name_resource_group' must be empty when Azure is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
			}
			gateway.Eip = d.Get("eip").(string)
		}
	}

	insaneMode := d.Get("insane_mode").(bool)
	if insaneMode {
		// Insane Mode encryption is not supported in China regions
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes^goaviatrix.AWSChina|goaviatrix.AzureArmRelatedCloudTypes^goaviatrix.AzureChina) {
			return fmt.Errorf("insane_mode is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			if d.Get("insane_mode_az").(string) == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS Top Secret (16384) or AWS Secret (32768)")
			}
			if d.Get("peering_ha_subnet").(string) != "" && d.Get("peering_ha_insane_mode_az").(string) == "" {
				return fmt.Errorf("peering_ha_insane_mode_az needed if insane_mode is enabled for AWS (1), AWSGov (256), AWS Top Secret (16384) or AWS Secret (32768) and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			insaneModeAz := d.Get("insane_mode_az").(string)
			strs = append(strs, gateway.VpcNet, insaneModeAz)
			gateway.VpcNet = strings.Join(strs, "~~")
		}
		gateway.InsaneMode = "on"
	} else {
		gateway.InsaneMode = "off"
	}

	samlEnabled := d.Get("saml_enabled").(bool)
	if samlEnabled {
		gateway.SamlEnabled = "yes"
	} else {
		gateway.SamlEnabled = "no"
	}

	splitTunnel := d.Get("split_tunnel").(bool)
	if splitTunnel {
		gateway.SplitTunnel = "yes"
	} else {
		gateway.SplitTunnel = "no"
	}

	enableElb := d.Get("enable_elb").(bool)
	if enableElb {
		gateway.EnableElb = "yes"
	} else {
		gateway.EnableElb = "no"
	}

	enableLdap := d.Get("enable_ldap").(bool)
	if enableLdap {
		gateway.EnableLdap = "yes"
	} else {
		gateway.EnableLdap = "no"
	}

	vpnStatus := d.Get("vpn_access").(bool)
	vpnProtocol := d.Get("vpn_protocol").(string)
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
			if gateway.EnableLdap == "yes" || gateway.OtpMode != "" {
				return fmt.Errorf("ldap and mfa can't be configured if saml is enabled")
			}
		}

		if gateway.OtpMode != "" && gateway.OtpMode != "2" && gateway.OtpMode != "3" {
			return fmt.Errorf("otp_mode can only be '2' or '3' or empty string")
		}

		if gateway.EnableLdap == "yes" && gateway.OtpMode == "3" {
			return fmt.Errorf("ldap can't be configured along with okta authentication")
		}
		if gateway.EnableLdap == "yes" {
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
		} else if gateway.OtpMode == "3" {
			if gateway.OktaToken == "" {
				return fmt.Errorf("okta token must be set if otp_mode is set to 3")
			}
			if gateway.OktaURL == "" {
				return fmt.Errorf("okta url must be set if otp_mode is set to 3")
			}
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

	peeringHaGwSize := d.Get("peering_ha_gw_size").(string)
	peeringHaSubnet := d.Get("peering_ha_subnet").(string)
	peeringHaZone := d.Get("peering_ha_zone").(string)
	peeringHaAvailabilityDomain := d.Get("peering_ha_availability_domain").(string)
	peeringHaFaultDomain := d.Get("peering_ha_fault_domain").(string)

	if peeringHaZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) && !d.Get("enable_public_subnet_filtering").(bool) {
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

	enableDesignatedGw := d.Get("enable_designated_gateway").(bool)
	if enableDesignatedGw {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'designated_gateway' feature is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
		}
		if peeringHaSubnet != "" || peeringHaZone != "" {
			return fmt.Errorf("can't enable HA for gateway with 'designated_gateway' enabled")
		}
		gateway.EnableDesignatedGateway = "true"
	}

	enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
	customerManagedKeys := d.Get("customer_managed_keys").(string)
	if enableEncryptVolume && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) providers")
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

	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureChina) {
		storageName, storageNameOk := d.GetOk("storage_name")
		if storageNameOk {
			gateway.StorageName = storageName.(string)
		} else {
			return fmt.Errorf("storage_name is required when creating a Gateway in AzureChina (2048)")
		}
	}

	_, tagListOk := d.GetOk("tag_list")
	_, tagsOk := d.GetOk("tags")
	if tagListOk || tagsOk {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return errors.New("failed to create gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}

		if tagListOk {
			tagList := d.Get("tag_list").([]interface{})
			tagListStr := goaviatrix.ExpandStringList(tagList)
			tagListStr = goaviatrix.TagListStrColon(tagListStr)
			gateway.TagList = strings.Join(tagListStr, ",")
		} else {
			tagsMap, err := extractTags(d, gateway.CloudType)
			if err != nil {
				return fmt.Errorf("error creating tags for gateway: %v", err)
			}
			tagJson, err := TagsMapToJson(tagsMap)
			if err != nil {
				return fmt.Errorf("failed to add tags when creating gateway: %v", err)
			}
			gateway.TagJson = tagJson
		}
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

	log.Printf("[INFO] Creating Aviatrix gateway: %#v", gateway)

	if d.Get("enable_public_subnet_filtering").(bool) {
		err := client.CreatePublicSubnetFilteringGateway(gateway)
		if err != nil {
			log.Printf("[INFO] failed to create public subnet filtering gateway: %#v", gateway)
			return fmt.Errorf("could not create public subnet filtering gateway: %v", err)
		}
		if !d.Get("public_subnet_filtering_guard_duty_enforced").(bool) {
			err = client.DisableGuardDutyEnforcement(gateway)
			if err != nil {
				return fmt.Errorf("could not disable guard duty enforcement for public subnet filtering gateway: %v", err)
			}
		}
	} else {
		err := client.CreateGateway(gateway)
		if err != nil {
			log.Printf("[INFO] failed to create Aviatrix gateway: %#v", gateway)
			return fmt.Errorf("failed to create Aviatrix gateway: %s", err)
		}
	}

	d.SetId(gateway.GwName)

	flag := false
	defer resourceAviatrixGatewayReadIfRequired(d, meta, &flag)

	if customerManagedKeys != "" && enableEncryptVolume {
		gwEncVolume := &goaviatrix.Gateway{
			GwName:              d.Get("gw_name").(string),
			CustomerManagedKeys: d.Get("customer_managed_keys").(string),
		}
		err := client.EnableEncryptVolume(gwEncVolume)
		if err != nil {
			return fmt.Errorf("failed to enable encrypt gateway volume when creating gateway: %s due to %s", gwEncVolume.GwName, err)
		}
	}

	enableVpnNat := d.Get("enable_vpn_nat").(bool)
	if vpnStatus {
		if !enableVpnNat {
			err := client.DisableVpnNat(gateway)
			if err != nil {
				return fmt.Errorf("failed to disable VPN NAT: %s", err)
			}
		}
	} else if !enableVpnNat {
		return fmt.Errorf("'enable_vpc_nat' is only supported for vpn gateway. Can't modify it for non-vpn gateway")
	}

	singleAZ := d.Get("single_az_ha").(bool)
	if singleAZ && !d.Get("enable_public_subnet_filtering").(bool) {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: "enabled",
		}

		log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

		err := client.EnableSingleAZGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to create single AZ GW HA: %s", err)
		}
	} else if !singleAZ && d.Get("enable_public_subnet_filtering").(bool) {
		// Public Subnet Filtering Gateways are created with single_az_ha=true by default.
		// Thus, if user set single_az_ha=false, we need to disable.
		err := client.DisableSingleAZGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable single AZ : %v", err)
		}
	}

	if enableDesignatedGw {
		additionalCidrsDesignatedGw := d.Get("additional_cidrs_designated_gateway").(string)
		if additionalCidrsDesignatedGw != "" {
			designatedGw := &goaviatrix.Gateway{
				GwName:                      d.Get("gw_name").(string),
				AdditionalCidrsDesignatedGw: additionalCidrsDesignatedGw,
			}
			err := client.EditDesignatedGateway(designatedGw)
			if err != nil {
				return fmt.Errorf("failed to edit additional cidrs for 'designated_gateway' feature due to %s", err)
			}
		}
	}

	// peering_ha_subnet is for Peering HA Gateway. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if peeringHaSubnet != "" || peeringHaZone != "" {
		if peeringHaGwSize == "" && !d.Get("enable_public_subnet_filtering").(bool) {
			return fmt.Errorf("A valid non empty peering_ha_gw_size parameter is mandatory for " +
				"this resource if peering_ha_subnet or peering_ha_zone is set. Example: t2.micro")
		}
		peeringHaGateway := &goaviatrix.Gateway{
			Eip:       d.Get("peering_ha_eip").(string),
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
		}

		if goaviatrix.IsCloudType(peeringHaGateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			peeringHaGateway.PeeringHASubnet = peeringHaSubnet
			if insaneMode {
				var peeringHaStrs []string
				peeringHaInsaneModeAz := d.Get("peering_ha_insane_mode_az").(string)
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
				peeringHaGateway.Eip = fmt.Sprintf("%s:%s", haAzureEipName.(string), peeringHaGateway.Eip)
			} else if haAzureEipNameOk {
				return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be empty when 'peering_ha_eip' is empty")
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		if d.Get("enable_public_subnet_filtering").(bool) {
			log.Printf("[INFO] Enable public subnet filtering HA: %#v", peeringHaGateway)
			var haRouteTables []string
			for _, v := range d.Get("public_subnet_filtering_ha_route_tables").(*schema.Set).List() {
				haRouteTables = append(haRouteTables, v.(string))
			}
			peeringHaGateway.RouteTable = strings.Join(haRouteTables, ",")
			peeringHaGateway.PeeringHASubnet = fmt.Sprintf("%s~~%s", peeringHaSubnet, peeringHaZone)
			err := client.EnablePublicSubnetFilteringHAGateway(peeringHaGateway)
			if err != nil {
				return fmt.Errorf("could not create public subnet filtering gateway HA: %v", err)
			}
		} else {
			log.Printf("[INFO] Enable peering HA: %#v", peeringHaGateway)
			err := client.EnablePeeringHaGateway(peeringHaGateway)
			if err != nil {
				return fmt.Errorf("failed to create peering HA: %s", err)
			}
		}

		log.Printf("[INFO] Resizing Peering HA Gateway: %#v", peeringHaGwSize)
		if peeringHaGwSize != gateway.VpcSize {
			if peeringHaGwSize == "" {
				return fmt.Errorf("A valid non empty peering_ha_gw_size parameter is mandatory for " +
					"this resource if peering_ha_subnet is set. Example: t2.micro")
			}
			peeringHaGateway := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string) + "-hagw", //CHECK THE NAME of peering ha gateway in
				// controller, test out first. just assuming it has that suffix
			}
			peeringHaGateway.GwSize = peeringHaGwSize
			err := client.UpdateGateway(peeringHaGateway)
			log.Printf("[INFO] Resizing Peering Ha Gateway size to: %s,", peeringHaGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Peering HA Gateway size: %s", err)
			}
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
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384) and AWS Secret (32768)")
	}

	if enableMonitorSubnets {
		log.Printf("[INFO] Enable Monitor Gateway Subnets")
		err := client.EnableMonitorGatewaySubnets(gateway.GwName, excludedInstances)
		if err != nil {
			return fmt.Errorf("could not enable monitor gateway subnets: %v", err)
		}
	}

	gatewayServer := &goaviatrix.Gateway{
		VpcID: d.Get("vpc_id").(string),
	}

	idleTimeoutValue := d.Get("idle_timeout").(int)
	if idleTimeoutValue != -1 {
		if d.Get("enable_elb").(bool) {
			gw, err := client.GetGateway(&goaviatrix.Gateway{
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string),
			})
			if err != nil {
				return fmt.Errorf("couldn't find Aviatrix Gateway for idle timeout : %s", gw.GwName)
			}
			gatewayServer.GwName = gw.ElbName
		} else {
			gatewayServer.GwName = d.Get("gw_name").(string)
		}
		enableVPNServer := &goaviatrix.VPNConfig{
			Name:  "Idle timeout",
			Value: strconv.Itoa(idleTimeoutValue),
		}
		log.Printf("[INFO] Enable Modify VPN Config (Idle Timeout)")
		err := client.EnableVPNConfig(gatewayServer, enableVPNServer)
		if err != nil {
			return fmt.Errorf("fail to enable idle timeout: %s", err)
		}
	}

	renegoIntervalValue := d.Get("renegotiation_interval").(int)
	if renegoIntervalValue != -1 {
		if d.Get("enable_elb").(bool) {
			gw, err := client.GetGateway(&goaviatrix.Gateway{
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string),
			})
			if err != nil {
				return fmt.Errorf("couldn't find Aviatrix Gateway renegotiation interval : %s", gw.GwName)
			}
			gatewayServer.GwName = gw.ElbName
		} else {
			gatewayServer.GwName = d.Get("gw_name").(string)
		}
		enableVPNServer := &goaviatrix.VPNConfig{
			Name:  "Renegotiation interval",
			Value: strconv.Itoa(renegoIntervalValue),
		}
		log.Printf("[INFO] Enable Modify VPN Config (Renegotiation Interval)")
		err := client.EnableVPNConfig(gatewayServer, enableVPNServer)
		if err != nil {
			return fmt.Errorf("fail to enable renegotiation interval: %s", err)
		}
	}

	if !d.Get("enable_jumbo_frame").(bool) {
		err := client.DisableJumboFrame(gateway)
		if err != nil {
			return fmt.Errorf("couldn't disable jumbo frames for Gateway: %s", err)
		}
	}

	if detectionTime, ok := d.GetOk("tunnel_detection_time"); ok {
		err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime.(int))
		if err != nil {
			return fmt.Errorf("could not set tunnel detection time during Gateway creation: %v", err)
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
	client := meta.(*goaviatrix.Client)

	var isImport bool
	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
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
		return fmt.Errorf("couldn't find Aviatrix Gateway %s: %v", gwName, err)
	}

	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)

	d.Set("cloud_type", gw.CloudType)
	d.Set("account_name", gw.AccountName)
	d.Set("gw_name", gw.GwName)
	d.Set("subnet", gw.VpcNet)
	d.Set("single_ip_snat", gw.EnableNat == "yes" && gw.SnatMode == "primary")
	d.Set("enable_ldap", gw.EnableLdapRead)
	d.Set("vpn_cidr", gw.VpnCidr)
	d.Set("saml_enabled", gw.SamlEnabled == "yes")
	d.Set("okta_url", gw.OktaURL)
	d.Set("okta_username_suffix", gw.OktaUsernameSuffix)
	d.Set("duo_integration_key", gw.DuoIntegrationKey)
	d.Set("duo_api_hostname", gw.DuoAPIHostname)
	d.Set("duo_push_mode", gw.DuoPushMode)
	d.Set("ldap_server", gw.LdapServer)
	d.Set("ldap_bind_dn", gw.LdapBindDn)
	d.Set("ldap_base_dn", gw.LdapBaseDn)
	d.Set("ldap_username_attribute", gw.LdapUserAttr)
	d.Set("single_az_ha", gw.SingleAZ == "yes")
	d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
	d.Set("eip", gw.PublicIP)
	d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
	d.Set("public_dns_server", gw.PublicDnsServer)
	d.Set("security_group_id", gw.GwSecurityGroupID)
	d.Set("peering_ha_security_group_id", gw.HaGw.GwSecurityGroupID)
	d.Set("private_ip", gw.PrivateIP)
	d.Set("enable_jumbo_frame", gw.JumboFrame)
	d.Set("enable_vpc_dns_server", goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled")
	d.Set("tunnel_detection_time", gw.TunnelDetectionTime)
	d.Set("image_version", gw.ImageVersion)
	d.Set("software_version", gw.SoftwareVersion)

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureChina) {
		d.Set("storage_name", gw.StorageName)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			d.Set("azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Gateway %s", gw.GwName)
		}
	}

	if gw.IdleTimeout != "NA" {
		idleTimeout, err := strconv.Atoi(gw.IdleTimeout)
		if err != nil {
			return fmt.Errorf("couldn't get idle timeout for the gateway %s: %v", gw.GwName, err)
		}
		d.Set("idle_timeout", idleTimeout)
	} else {
		d.Set("idle_timeout", -1)
	}

	if gw.RenegotiationInterval != "NA" {
		renegotiationInterval, err := strconv.Atoi(gw.RenegotiationInterval)
		if err != nil {
			return fmt.Errorf("couldn't get renegotiation interval for the gateway %s: %v", gw.GwName, err)
		}
		d.Set("renegotiation_interval", renegotiationInterval)
	} else {
		d.Set("renegotiation_interval", -1)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
		// AWS vpc_id returns as <vpc_id>~~<other vpc info>
		d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
		d.Set("vpc_reg", gw.VpcRegion)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		// gcp vpc_id returns as <vpc_id>~-~<other vpc info>
		d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0])
		d.Set("vpc_reg", gw.GatewayZone)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
		d.Set("allocate_new_eip", gw.AllocateNewEipRead)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AliCloudRelatedCloudTypes) {
		d.Set("allocate_new_eip", true)
	}

	if gw.EnableDesignatedGateway == "Yes" || gw.EnableDesignatedGateway == "yes" {
		d.Set("enable_designated_gateway", true)
		cidrsTF := strings.Split(d.Get("additional_cidrs_designated_gateway").(string), ",")
		cidrsRESTAPI := strings.Split(gw.AdditionalCidrsDesignatedGw, ",")
		if len(goaviatrix.Difference(cidrsTF, cidrsRESTAPI)) == 0 && len(goaviatrix.Difference(cidrsRESTAPI, cidrsTF)) == 0 {
			d.Set("additional_cidrs_designated_gateway", d.Get("additional_cidrs_designated_gateway").(string))
		} else {
			d.Set("additional_cidrs_designated_gateway", gw.AdditionalCidrsDesignatedGw)
		}
	} else {
		d.Set("enable_designated_gateway", false)
		d.Set("additional_cidrs_designated_gateway", "")
	}

	_, zoneIsSet := d.GetOk("zone")
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && (isImport || zoneIsSet) && gw.GatewayZone != "AvailabilitySet" {
		d.Set("zone", "az-"+gw.GatewayZone)
	}

	if gw.VpnStatus != "" {
		if gw.VpnStatus == "disabled" {
			d.Set("vpn_access", false)
			d.Set("enable_vpn_nat", true)
			d.Set("vpn_protocol", "")
			d.Set("split_tunnel", true)
			d.Set("max_vpn_conn", "")
		} else if gw.VpnStatus == "enabled" {
			d.Set("vpn_access", true)
			d.Set("split_tunnel", gw.SplitTunnel == "yes")
			d.Set("max_vpn_conn", gw.MaxConn)
			d.Set("enable_vpn_nat", gw.EnableVpnNat)
			if gw.ElbState == "enabled" {
				if strings.ToUpper(gw.VpnProtocol) == "UDP" {
					d.Set("vpn_protocol", "UDP")
				} else {
					d.Set("vpn_protocol", "TCP")
				}
			} else {
				d.Set("vpn_protocol", "UDP")
			}
		}
	}

	if gw.ElbState == "enabled" {
		d.Set("enable_elb", true)
		d.Set("elb_name", gw.ElbName)
		d.Set("elb_dns_name", gw.ElbDNSName)
	} else {
		d.Set("enable_elb", false)
		d.Set("elb_name", "")
	}

	if gw.AuthMethod == "duo_auth" || gw.AuthMethod == "duo_auth+LDAP" {
		d.Set("otp_mode", "2")
	} else if gw.AuthMethod == "okta_auth" {
		d.Set("otp_mode", "3")
	} else {
		d.Set("otp_mode", "")
	}

	if gw.NewZone != "" {
		d.Set("zone", gw.NewZone)
	}

	// Though go_aviatrix Gateway struct declares VpcSize as only used on gateway creation
	// it is the attribute receiving the instance size of an existing gateway instead of
	// GwSize. (at least in v3.5)
	if gw.GwSize != "" {
		d.Set("gw_size", gw.GwSize)
	} else {
		if gw.VpcSize != "" {
			d.Set("gw_size", gw.VpcSize)
		}
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
			if err := d.Set("tags", gw.Tags); err != nil {
				log.Printf("[WARN] Error setting tags for (%s): %s", d.Id(), err)
			}
		}
	}

	if gw.VpnStatus == "enabled" && gw.SplitTunnel == "yes" {
		d.Set("name_servers", gw.NameServers)
		d.Set("search_domains", gw.SearchDomains)
		d.Set("additional_cidrs", gw.AdditionalCidrs)
	} else {
		d.Set("name_servers", "")
		d.Set("search_domains", "")
		d.Set("additional_cidrs", "")
	}

	d.Set("enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
	if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
		return fmt.Errorf("setting 'monitor_exclude_list' to state: %v", err)
	}

	fqdnGatewayLanInterface := fmt.Sprintf("av-nic-%s_eth1", gw.GwName)
	fqdnLanCidr, ok := gw.ArmFqdnLanCidr[gw.GwName]
	if ok && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("fqdn_lan_interface", fqdnGatewayLanInterface)
		d.Set("fqdn_lan_cidr", fqdnLanCidr)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("fqdn_lan_vpc_id", gw.BundleVpcInfo.LAN.VpcID)
		d.Set("fqdn_lan_cidr", strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0])
	} else {
		d.Set("fqdn_lan_interface", "")
		d.Set("fqdn_lan_cidr", "")
	}

	if !gw.IsPsfGateway {
		d.Set("enable_public_subnet_filtering", false)
		d.Set("public_subnet_filtering_route_tables", []string{})
		d.Set("public_subnet_filtering_ha_route_tables", []string{})
		d.Set("public_subnet_filtering_guard_duty_enforced", true)
	} else {
		d.Set("enable_public_subnet_filtering", true)
		if err := d.Set("public_subnet_filtering_route_tables", gw.PsfDetails.RouteTableList); err != nil {
			return fmt.Errorf("could not set public_subnet_filtering_route_tables into state: %v", err)
		}
		d.Set("public_subnet_filtering_guard_duty_enforced", gw.PsfDetails.GuardDutyEnforced == "yes")
		d.Set("subnet", gw.PsfDetails.GwSubnetCidr)
		d.Set("zone", gw.PsfDetails.GwSubnetAz)
		if gw.HaGw.GwSize == "" {
			err := d.Set("public_subnet_filtering_ha_route_tables", []string{})
			if err != nil {
				return fmt.Errorf("could not set public_subnet_filtering_ha_route_tables into state: %v", err)
			}
		} else {
			if err := d.Set("public_subnet_filtering_ha_route_tables", gw.PsfDetails.HaRouteTableList); err != nil {
				return fmt.Errorf("could not set public_subnet_filtering_ha_route_tables into state: %v", err)
			}
			d.Set("peering_ha_subnet", gw.PsfDetails.HaGwSubnetCidr)
			d.Set("peering_ha_zone", gw.PsfDetails.HaGwSubnetAz)
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

	if gw.HaGw.GwSize == "" {
		d.Set("peering_ha_cloud_instance_id", "")
		d.Set("peering_ha_subnet", "")
		d.Set("peering_ha_zone", "")
		d.Set("peering_ha_eip", "")
		d.Set("peering_ha_azure_eip_name_resource_group", "")
		d.Set("peering_ha_gw_size", "")
		d.Set("peering_ha_insane_mode_az", "")
		return nil
	}

	d.Set("peering_ha_cloud_instance_id", gw.HaGw.CloudnGatewayInstID)
	d.Set("peering_ha_gw_name", gw.HaGw.GwName)
	d.Set("peering_ha_eip", gw.HaGw.PublicIP)
	d.Set("peering_ha_gw_size", gw.HaGw.GwSize)
	d.Set("peering_ha_private_ip", gw.HaGw.PrivateIP)
	d.Set("peering_ha_software_version", gw.HaGw.SoftwareVersion)
	d.Set("peering_ha_image_version", gw.HaGw.ImageVersion)
	if gw.IsPsfGateway {
		// For PSF gateway, peering_ha_subnet and peering_ha_zone are
		// set above. Return early.
		return nil
	}

	if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("peering_ha_subnet", gw.HaGw.VpcNet)
		d.Set("peering_ha_zone", "")
		if gw.HaGw.InsaneMode == "yes" {
			d.Set("peering_ha_insane_mode_az", gw.HaGw.GatewayZone)
		}
	} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		d.Set("peering_ha_subnet", gw.HaGw.VpcNet)
		d.Set("peering_ha_zone", "")
		if gw.HaGw.GatewayZone != "" {
			d.Set("peering_ha_availability_domain", gw.HaGw.GatewayZone)
		} else {
			d.Set("peering_ha_availability_domain", d.Get("peering_ha_availability_domain").(string))
		}
		d.Set("peering_ha_fault_domain", gw.HaGw.FaultDomain)
	} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("peering_ha_zone", gw.HaGw.GatewayZone)
		// only set peering_ha_subnet if the user has explicitly set it.
		if d.Get("peering_ha_subnet").(string) != "" || isImport {
			d.Set("peering_ha_subnet", gw.HaGw.VpcNet)
		}
	} else if goaviatrix.IsCloudType(gw.HaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("peering_ha_subnet", gw.HaGw.VpcNet)
		if _, haZoneIsSet := d.GetOk("peering_ha_zone"); isImport || haZoneIsSet {
			if gw.GatewayZone != "AvailabilitySet" {
				d.Set("peering_ha_zone", "az-"+gw.GatewayZone)
			}
		}

		azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
		if len(azureEip) == 3 {
			d.Set("peering_ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Peering HA Gateway %s", gw.GwName)
		}
	} else if gw.HaGw.CloudType == goaviatrix.AliCloud {
		d.Set("peering_ha_subnet", gw.HaGw.VpcNet)
		d.Set("peering_ha_zone", "")
	}

	return nil
}

func resourceAviatrixGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", d.Get("gw_name").(string))

	d.Partial(true)
	if d.HasChange("cloud_type") {
		return fmt.Errorf("updating cloud_type is not allowed")
	}
	if d.HasChange("account_name") {
		return fmt.Errorf("updating account_name is not allowed")
	}
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
	if d.HasChange("insane_mode") {
		return fmt.Errorf("updating insane_mode is not allowed")
	}
	if d.HasChange("insane_mode_az") {
		return fmt.Errorf("updating insane_mode_az is not allowed")
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
		if o.(string) != "" && n.(string) != "" {
			return fmt.Errorf("failed to update gateway: changing 'peering_ha_azure_eip_name_resource_group' is not allowed")
		}
	}
	if d.HasChange("enable_spot_instance") {
		return fmt.Errorf("updating enable_spot_instance is not allowed")
	}
	if d.HasChange("spot_price") {
		return fmt.Errorf("updating spot_price is not allowed")
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
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
		GwSize:    d.Get("gw_size").(string),
	}
	vpnAccess := d.Get("vpn_access").(bool)
	enableElb := false
	geoVpnDnsName := ""
	if vpnAccess {
		enableElb = d.Get("enable_elb").(bool)
		if enableElb {
			gateway.ElbDNSName = d.Get("elb_dns_name").(string)
			geoVpn, err := client.GetGeoVPNName(gateway)
			if err == nil {
				geoVpnDnsName = geoVpn.ServiceName
			}
		}
	}
	if d.HasChange("peering_ha_zone") {
		peeringHaZone := d.Get("peering_ha_zone").(string)
		if peeringHaZone != "" && !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) && !d.Get("enable_public_subnet_filtering").(bool) {
			return fmt.Errorf("'peering_ha_zone' is only valid for GCP, Azure and Public Subnet Filtering Gateway if enabling Peering HA")
		}
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		peeringHaSubnet := d.Get("peering_ha_subnet").(string)
		peeringHaZone := d.Get("peering_ha_zone").(string)
		if peeringHaZone == "" && peeringHaSubnet != "" {
			return fmt.Errorf("'peering_ha_zone' must be set to enable Peering HA on GCP, " +
				"cannot enable Peering HA with only 'peering_ha_subnet' enabled")
		}
	}
	if goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		peeringHaSubnet := d.Get("peering_ha_subnet").(string)
		peeringHaZone := d.Get("peering_ha_zone").(string)
		if peeringHaZone != "" && peeringHaSubnet == "" {
			return fmt.Errorf("'peering_ha_subnet' must be set to enable Peering HA on Azure, " +
				"cannot enable Peering HA with only 'peering_ha_zone' enabled")
		}
	}

	singleAZ := d.Get("single_az_ha").(bool)
	if singleAZ {
		gateway.SingleAZ = "enabled"
	} else {
		gateway.SingleAZ = "disabled"
	}

	peeringHaGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
		GwSize:    d.Get("peering_ha_gw_size").(string),
	}

	// Get primary gw size if gw_size changed, to be used later on for peering ha gw size update
	primaryGwSize := d.Get("gw_size").(string)
	if d.HasChange("gw_size") {
		old, _ := d.GetChange("gw_size")
		primaryGwSize = old.(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Gateway: %s", err)
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
			VpcID:              d.Get("vpc_id").(string),
			OtpMode:            d.Get("otp_mode").(string),
			OktaToken:          d.Get("okta_token").(string),
			OktaURL:            d.Get("okta_url").(string),
			OktaUsernameSuffix: d.Get("okta_username_suffix").(string),
			DuoIntegrationKey:  d.Get("duo_integration_key").(string),
			DuoSecretKey:       d.Get("duo_secret_key").(string),
			DuoAPIHostname:     d.Get("duo_api_hostname").(string),
			DuoPushMode:        d.Get("duo_push_mode").(string),
			LdapServer:         d.Get("ldap_server").(string),
			LdapBindDn:         d.Get("ldap_bind_dn").(string),
			LdapPassword:       d.Get("ldap_password").(string),
			LdapBaseDn:         d.Get("ldap_base_dn").(string),
			LdapUserAttr:       d.Get("ldap_username_attribute").(string),
		}

		samlEnabled := d.Get("saml_enabled").(bool)
		if samlEnabled {
			vpn_gw.SamlEnabled = "yes"
		} else {
			vpn_gw.SamlEnabled = "no"
		}

		enableLdap := d.Get("enable_ldap").(bool)
		if enableLdap {
			vpn_gw.EnableLdap = "yes"
		} else {
			vpn_gw.EnableLdap = "no"
		}

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
			if vpn_gw.EnableLdap == "yes" || vpn_gw.OtpMode != "" {
				return fmt.Errorf("ldap and mfa can't be configured if saml is enabled")
			}
		}
		if vpn_gw.EnableLdap == "yes" && vpn_gw.OtpMode == "3" {
			return fmt.Errorf("ldap can't be configured along with okta authentication")
		}
		if vpn_gw.EnableLdap == "yes" {
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
			if vpn_gw.EnableLdap == "yes" {
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
			if vpn_gw.EnableLdap == "yes" {
				vpn_gw.AuthType = "ldap_auth"
			} else if vpn_gw.SamlEnabled == "yes" {
				vpn_gw.AuthType = "saml_auth"
			} else {
				vpn_gw.AuthType = "none"
			}
		}
		if enableElb := d.Get("enable_elb").(bool); enableElb {
			vpn_gw.LbOrGatewayName = d.Get("elb_name").(string)
		} else {
			vpn_gw.LbOrGatewayName = d.Get("gw_name").(string)
		}

		err := client.SetVpnGatewayAuthentication(vpn_gw)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix VPN Gateway Authentication: %s", err)
		}
	}

	if d.HasChange("tag_list") || d.HasChange("tags") {
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("failed to update gateway: adding tags is only supported for AWS (1), Azure (8), AzureGov (32), AWSGov(256) AWSChina (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
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
				return fmt.Errorf("failed to update tag_list for gateway: %s", err)
			}
		}
		if d.HasChange("tags") && len(tagList) == 0 {
			tagsMap, err := extractTags(d, gateway.CloudType)
			if err != nil {
				return fmt.Errorf("failed to update tags for gateway: %v", err)
			}
			tags.Tags = tagsMap
			tagJson, err := TagsMapToJson(tagsMap)
			if err != nil {
				return fmt.Errorf("failed to update tags for gateway: %v", err)
			}
			tags.TagJson = tagJson
			err = client.UpdateTags(tags)
			if err != nil {
				return fmt.Errorf("failed to update tags for gateway: %v", err)
			}
		}
	}

	if d.HasChange("split_tunnel") || d.HasChange("additional_cidrs") ||
		d.HasChange("name_servers") || d.HasChange("search_domains") {
		splitTunnel := d.Get("split_tunnel").(bool)
		sTunnel := &goaviatrix.SplitTunnel{
			VpcID:   d.Get("vpc_id").(string),
			ElbName: d.Get("elb_name").(string),
		}
		if sTunnel.ElbName == "" {
			sTunnel.ElbName = d.Get("gw_name").(string)
		}

		if splitTunnel && (d.HasChange("additional_cidrs") || d.HasChange("name_servers") || d.HasChange("search_domains")) {
			sTunnel.AdditionalCidrs = d.Get("additional_cidrs").(string)
			sTunnel.NameServers = d.Get("name_servers").(string)
			sTunnel.SearchDomains = d.Get("search_domains").(string)
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
				return fmt.Errorf("failed to modify split tunnel: %s", err)
			}
		} else if !splitTunnel && (d.Get("additional_cidrs").(string) != "" || d.Get("name_servers").(string) != "" || d.Get("search_domains").(string) != "") {
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
				return fmt.Errorf("failed to disable split tunnel: %s", err)
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
			gw.EnableNat = "yes"
		} else {
			gw.EnableNat = "no"
		}

		if enableNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable SNAT: %s", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable SNAT: %s", err)
			}
		}

	}
	if d.HasChange("additional_cidrs_designated_gateway") {
		if !d.Get("enable_designated_gateway").(bool) {
			return fmt.Errorf("failed to edit additional cidrs for 'designated_gateway' since it is not enabled")
		}
		if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("'designated_gateway' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768)")
		}
		designatedGw := &goaviatrix.Gateway{
			GwName:                      d.Get("gw_name").(string),
			AdditionalCidrsDesignatedGw: d.Get("additional_cidrs_designated_gateway").(string),
		}
		err := client.EditDesignatedGateway(designatedGw)
		if err != nil {
			return fmt.Errorf("failed to edit additional cidrs for 'designated_gateway' feature due to %s", err)
		}
	}
	if d.HasChange("vpn_cidr") {
		if d.Get("vpn_access").(bool) {
			gw := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string),
				VpnCidr:   d.Get("vpn_cidr").(string),
			}

			err := client.UpdateVpnCidr(gw)
			if err != nil {
				return fmt.Errorf("failed to update vpn cidr: %s", err)
			}
		} else {
			log.Printf("[INFO] can't update vpn cidr because vpn_access is disabled for gateway: %#v", gateway.GwName)
		}

	}
	if d.HasChange("max_vpn_conn") {
		if vpnAccess {
			gw := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string),
				VpcID:     d.Get("vpc_id").(string),
				ElbName:   d.Get("elb_name").(string),
			}

			if gw.ElbName == "" {
				gw.ElbName = d.Get("gw_name").(string)
			}

			_, n := d.GetChange("max_vpn_conn")
			gw.MaxConn = n.(string)
			if enableElb && geoVpnDnsName != "" {
				gw.ElbName = geoVpnDnsName
				gw.Dns = "true"
			}
			err := client.UpdateMaxVpnConn(gw)
			if err != nil {
				return fmt.Errorf("failed to update max vpn connections: %s", err)
			}
		} else {
			log.Printf("[INFO] can't update max vpn connections because vpn is disabled for gateway: %#v", gateway.GwName)
		}

	}

	newHaGwEnabled := false
	if d.HasChange("peering_ha_subnet") || d.HasChange("peering_ha_zone") || d.HasChange("peering_ha_insane_mode_az") ||
		d.HasChange("peering_ha_availability_domain") || d.HasChange("peering_ha_fault_domain") {
		if d.Get("enable_designated_gateway").(bool) {
			return fmt.Errorf("can't update HA status for gateway with 'designated_gateway' enabled")
		}
		gw := &goaviatrix.Gateway{
			Eip:       d.Get("peering_ha_eip").(string),
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
			GwSize:    d.Get("peering_ha_gw_size").(string),
		}

		haAzureEipName, haAzureEipNameOk := d.GetOk("peering_ha_azure_eip_name_resource_group")
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			if gw.Eip != "" && gw.GwSize != "" {
				// No change will be detected when peering_ha_eip is set to the empty string because it is computed.
				// Instead, check peering_ha_gw_size to detect when HA gateway is being deleted.
				if !haAzureEipNameOk {
					return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be set when a custom EIP is provided and cloud_type is Azure (8), AzureGov (32) or AzureChina (2048)")
				}
				// AVX-9874 Azure EIP has a different format e.g. 'test_ip:rg:104.45.186.20'
				gw.Eip = fmt.Sprintf("%s:%s", haAzureEipName.(string), gw.Eip)
			}
		} else if haAzureEipNameOk {
			return fmt.Errorf("failed to create Peering HA Gateway: 'peering_ha_azure_eip_name_resource_group' must be empty when cloud_type is not one of Azure (8), AzureGov (32) or AzureChina (2048)")
		}

		oldSubnet, newSubnet := d.GetChange("peering_ha_subnet")
		oldZone, newZone := d.GetChange("peering_ha_zone")
		deleteHaGw := false
		changeHaGw := false

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			gw.PeeringHASubnet = d.Get("peering_ha_subnet").(string)
			if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && newZone != "" {
				gw.PeeringHASubnet = fmt.Sprintf("%s~~%s~~", newSubnet, newZone)
			}
			peeringHaAvailabilityDomain := d.Get("peering_ha_availability_domain").(string)
			peeringHaFaultDomain := d.Get("peering_ha_fault_domain").(string)
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
			gw.NewZone = d.Get("peering_ha_zone").(string)
			gw.NewSubnet = d.Get("peering_ha_subnet").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}

		if d.Get("insane_mode").(bool) && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			var haStrs []string
			peeringHaInsaneModeAz := d.Get("peering_ha_insane_mode_az").(string)
			if peeringHaInsaneModeAz == "" {
				return fmt.Errorf("peering_ha_insane_mode_az needed if insane_mode is enabled and peering_ha_subnet is set")
			}
			haStrs = append(haStrs, gw.PeeringHASubnet, peeringHaInsaneModeAz)
			gw.PeeringHASubnet = strings.Join(haStrs, "~~")
		}

		if (newHaGwEnabled || changeHaGw) && gw.GwSize == "" {
			return fmt.Errorf("A valid non empty peering_ha_gw_size parameter is mandatory for this resource if " +
				"peering_ha_subnet or peering_ha_zone is set")
		} else if deleteHaGw && gw.GwSize != "" {
			return fmt.Errorf("peering_ha_gw_size must be empty if transit HA gateway is deleted")
		}

		if d.Get("enable_public_subnet_filtering").(bool) {
			var haRouteTables []string
			for _, v := range d.Get("public_subnet_filtering_ha_route_tables").(*schema.Set).List() {
				haRouteTables = append(haRouteTables, v.(string))
			}
			gw.RouteTable = strings.Join(haRouteTables, ",")
			gw.PeeringHASubnet = fmt.Sprintf("%s~~%s", d.Get("peering_ha_subnet"), d.Get("peering_ha_zone"))
			if newHaGwEnabled {
				err := client.EnablePublicSubnetFilteringHAGateway(gw)
				if err != nil {
					return fmt.Errorf("failed to enable Aviatrix public subnet filtering HA gateway: %s", err)
				}
			} else if deleteHaGw {
				err := client.DeletePublicSubnetFilteringGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix public subnet filtering HA gateway: %s", err)
				}
			} else if changeHaGw {
				err := client.DeletePublicSubnetFilteringGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix public subnet filtering HA gateway: %s", err)
				}

				gw.Eip = ""

				gateway.GwName = d.Get("gw_name").(string)
				err = client.EnablePublicSubnetFilteringHAGateway(gw)
				if err != nil {
					return fmt.Errorf("failed to enable Aviatrix public subnet filtering HA gateway: %s", err)
				}

				newHaGwEnabled = true
			}
		} else {
			if newHaGwEnabled {
				err := client.EnablePeeringHaGateway(gw)
				if err != nil {
					return fmt.Errorf("failed to enable Aviatrix peering HA gateway: %s", err)
				}
			} else if deleteHaGw {
				err := client.DeleteGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix peering HA gateway: %s", err)
				}
			} else if changeHaGw {
				err := client.DeleteGateway(peeringHaGateway)
				if err != nil {
					return fmt.Errorf("failed to delete Aviatrix peering HA gateway: %s", err)
				}

				gw.Eip = ""

				gateway.GwName = d.Get("gw_name").(string)
				haErr := client.EnablePeeringHaGateway(gw)
				if haErr != nil {
					return fmt.Errorf("failed to enable Aviatrix peering HA gateway: %s", haErr)
				}

				newHaGwEnabled = true
			}
		}
	}
	haSubnet := d.Get("peering_ha_subnet").(string)
	haZone := d.Get("peering_ha_zone").(string)
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

	if d.HasChange("peering_ha_gw_size") || newHaGwEnabled {
		newHaGwSize := d.Get("peering_ha_gw_size").(string)
		if !newHaGwEnabled || (newHaGwSize != primaryGwSize) {
			// MODIFIES Peering HA GW SIZE if
			// Ha gateway wasn't newly configured
			// OR
			// newly configured peering HA gateway is set to be different size than primary gateway
			// (when peering ha gateway is enabled, it's size is by default the same as primary gateway)
			_, err := client.GetGateway(peeringHaGateway)
			if err != nil {
				if err != goaviatrix.ErrNotFound {
					return fmt.Errorf("couldn't find Aviatrix Peering HA Gateway while trying to update HA Gw "+
						"size: %s", err)
				}
			} else {
				if peeringHaGateway.GwSize == "" {
					return fmt.Errorf("A valid non empty peering_ha_gw_size parameter is mandatory for this resource if " +
						"peering_ha_subnet or peering_ha_zone is set. Example: t2.micro or us-west1-b respectively")
				}
				err = client.UpdateGateway(peeringHaGateway)
				log.Printf("[INFO] Updating Peering HA Gateway size to: %s ", peeringHaGateway.GwSize)
				if err != nil {
					return fmt.Errorf("failed to update Aviatrix Peering HA Gw size: %s", err)
				}
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
		return fmt.Errorf("'enable_vpc_dns_server' only supported by AWS (1), Azure (8), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192)")
	}

	if d.HasChange("enable_vpn_nat") {
		if !vpnAccess {
			return fmt.Errorf("'enable_vpc_nat' is only supported for vpn gateway. Can't updated it for Non VPN Gateway")
		} else {
			gw := &goaviatrix.Gateway{
				CloudType:    d.Get("cloud_type").(int),
				GwName:       d.Get("gw_name").(string),
				VpcID:        d.Get("vpc_id").(string),
				ElbName:      d.Get("elb_name").(string),
				EnableVpnNat: true,
			}
			if enableElb && geoVpnDnsName != "" {
				gw.ElbName = geoVpnDnsName
				gw.Dns = "true"
			}

			if d.Get("enable_vpn_nat").(bool) {
				err := client.EnableVpnNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable VPN NAT: %s", err)
				}
			} else if !d.Get("enable_vpn_nat").(bool) {
				err := client.DisableVpnNat(gw)
				if err != nil {
					return fmt.Errorf("failed to disable VPN NAT: %s", err)
				}
			}
		}
	}

	if d.HasChange("enable_encrypt_volume") {
		if d.Get("enable_encrypt_volume").(bool) {
			if !goaviatrix.IsCloudType(gateway.CloudType, goaviatrix.AWSRelatedCloudTypes) {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS (1), AWSGov (256), AWSChina (1024), AWS Top Secret (16384) and AWS Secret (32768) provider")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              d.Get("gw_name").(string),
				CustomerManagedKeys: d.Get("customer_managed_keys").(string),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
			}

			haSubnet := d.Get("peering_ha_subnet").(string)
			haZone := d.Get("peering_ha_zone").(string)
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

	gatewayServer := &goaviatrix.Gateway{
		VpcID: d.Get("vpc_id").(string),
	}

	if d.HasChange("idle_timeout") {
		idleTimeoutValue := d.Get("idle_timeout").(int)
		VPNServer := &goaviatrix.VPNConfig{
			Name: "Idle timeout",
		}
		if d.Get("enable_elb").(bool) {
			gatewayServer.GwName = d.Get("elb_name").(string)
		} else {
			gatewayServer.GwName = d.Get("gw_name").(string)
		}
		if idleTimeoutValue != -1 {
			VPNServer.Value = strconv.Itoa(idleTimeoutValue)
			log.Printf("[INFO] Modify VPN Config (update idle timeout value)")
			err := client.EnableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to update idle timeout value due to : %s", err)
			}
		} else {
			log.Printf("[INFO] Modify VPN Config (disable idle timeout)")
			err := client.DisableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to disable idle timeout due to : %s", err)
			}
		}
	}

	if d.HasChange("renegotiation_interval") {
		renegoIntervalValue := d.Get("renegotiation_interval").(int)
		VPNServer := &goaviatrix.VPNConfig{
			Name: "Renegotiation interval",
		}
		if d.Get("enable_elb").(bool) {
			gatewayServer.GwName = d.Get("elb_name").(string)
		} else {
			gatewayServer.GwName = d.Get("gw_name").(string)
		}
		if renegoIntervalValue != -1 {
			VPNServer.Value = strconv.Itoa(renegoIntervalValue)
			log.Printf("[INFO] Modify VPN Config (update renegotiation interval value)")
			err := client.EnableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to enable renegotiation interval due to : %s", err)
			}
		} else {
			log.Printf("[INFO] Modify VPN Config (disable renegotiation interval)")
			err := client.DisableVPNConfig(gatewayServer, VPNServer)
			if err != nil {
				return fmt.Errorf("fail to disable renegotiation interval due to: %s", err)
			}
		}
	}

	gatewayServer = &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
	}

	if d.HasChange("public_subnet_filtering_route_tables") {
		var routeTables []string
		for _, v := range d.Get("public_subnet_filtering_route_tables").(*schema.Set).List() {
			routeTables = append(routeTables, v.(string))
		}
		if len(routeTables) == 0 {
			return fmt.Errorf("attribute 'public_subnet_filtering_route_tables' must not be empty if 'enable_public_subnet_filtering' is set to true")
		}
		err := client.EditPublicSubnetFilteringRouteTableList(gatewayServer, routeTables)
		if err != nil {
			return fmt.Errorf("could not edit public subnet filtering route table rules: %v", err)
		}
	}
	if d.HasChange("public_subnet_filtering_ha_route_tables") && !d.HasChange("peering_ha_subnet") && d.Get("peering_ha_subnet").(string) != "" {
		var haRouteTables []string
		for _, v := range d.Get("public_subnet_filtering_ha_route_tables").(*schema.Set).List() {
			haRouteTables = append(haRouteTables, v.(string))
		}
		peeringHaGateway.RouteTable = strings.Join(haRouteTables, ",")
		err := client.EditPublicSubnetFilteringRouteTableList(peeringHaGateway, haRouteTables)
		if err != nil {
			return fmt.Errorf("could not edit HA public subnet filtering route table rules: %v", err)
		}
	}
	if d.HasChange("public_subnet_filtering_guard_duty_enforced") {
		if d.Get("public_subnet_filtering_guard_duty_enforced").(bool) {
			err := client.EnableGuardDutyEnforcement(gatewayServer)
			if err != nil {
				return fmt.Errorf("could not enable public subnet filtering guard duty enforcement: %v", err)
			}
		} else {
			err := client.DisableGuardDutyEnforcement(gatewayServer)
			if err != nil {
				return fmt.Errorf("could not disable public subnet filtering guard duty enforcement: %v", err)
			}
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if d.Get("enable_jumbo_frame").(bool) {
			err := client.EnableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("couldn't enable jumbo frames for Gateway when updating: %s", err)
			}
		} else {
			err := client.DisableJumboFrame(gateway)
			if err != nil {
				return fmt.Errorf("couldn't disable jumbo frames for Gateway when updating: %s", err)
			}
		}
	}

	if d.HasChange("tunnel_detection_time") {
		detectionTimeInterface, ok := d.GetOk("tunnel_detection_time")
		var detectionTime int
		if ok {
			detectionTime = detectionTimeInterface.(int)
		} else {
			detectionTime, err = client.GetTunnelDetectionTime("Controller")
			if err != nil {
				return fmt.Errorf("could not get default tunnel detection time during Gateway update: %v", err)
			}
		}
		err := client.ModifyTunnelDetectionTime(gateway.GwName, detectionTime)
		if err != nil {
			return fmt.Errorf("could not modify tunnel detection time during Gateway update: %v", err)
		}
	}

	primaryHasVersionChange := d.HasChanges("software_version", "image_version")
	haHasVersionChange := haEnabled && d.HasChanges("peering_ha_software_version", "peering_ha_image_version")
	primaryHasImageVersionChange := d.HasChange("image_version")
	haHasImageVersionChange := d.HasChange("peering_ha_image_version")
	if primaryHasVersionChange || haHasVersionChange {
		if primaryHasVersionChange && haHasVersionChange && !primaryHasImageVersionChange && !haHasImageVersionChange {
			// Both Primary and HA have changed just their software_version
			// so we can perform upgrade in parallel.
			swVersion := d.Get("software_version").(string)
			gw := &goaviatrix.Gateway{
				GwName:          d.Get("gw_name").(string),
				SoftwareVersion: swVersion,
			}
			haSwVersion := d.Get("peering_ha_software_version").(string)
			hagw := &goaviatrix.Gateway{
				GwName:          d.Get("gw_name").(string) + "-hagw",
				SoftwareVersion: haSwVersion,
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
				return fmt.Errorf("could not upgrade primary and HA gateway "+
					"software_version=%s peering_ha_software_version=%s: primaryErr: %v haErr: %v",
					swVersion, haSwVersion, primaryErr, haErr)
			} else if primaryErr != nil {
				return fmt.Errorf("could not upgrade primary gateway software_version=%s: %v", swVersion, primaryErr)
			} else if haErr != nil {
				return fmt.Errorf("could not upgrade HA gateway peering_ha_software_version=%s: %v", haSwVersion, primaryErr)
			}
		} else { // Only primary or only HA has changed, or they have changed image_version
			if primaryHasVersionChange {
				swVersion := d.Get("software_version").(string)
				imageVersion := d.Get("image_version").(string)
				gw := &goaviatrix.Gateway{
					GwName:          d.Get("gw_name").(string),
					SoftwareVersion: swVersion,
					ImageVersion:    imageVersion,
				}
				err := client.UpgradeGateway(gw)
				if err != nil {
					return fmt.Errorf("could not upgrade gateway during update image_version=%s software_version=%s: %v", gw.ImageVersion, gw.SoftwareVersion, err)
				}
			}
			if haHasVersionChange {
				haSwVersion := d.Get("peering_ha_software_version").(string)
				haImageVersion := d.Get("peering_ha_image_version").(string)
				hagw := &goaviatrix.Gateway{
					GwName:          d.Get("gw_name").(string) + "-hagw",
					SoftwareVersion: haSwVersion,
					ImageVersion:    haImageVersion,
				}
				err := client.UpgradeGateway(hagw)
				if err != nil {
					return fmt.Errorf("could not upgrade HA gateway during update image_version=%s software_version=%s: %v", hagw.ImageVersion, hagw.SoftwareVersion, err)
				}
			}
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixGatewayRead(d, meta)
}

func resourceAviatrixGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	var err error
	isPublicSubnetFilteringGateway := d.Get("enable_public_subnet_filtering").(bool)
	// peering_ha_subnet is for Peering HA
	peeringHaSubnet := d.Get("peering_ha_subnet").(string)
	peeringHaZone := d.Get("peering_ha_zone").(string)
	if peeringHaSubnet != "" || peeringHaZone != "" {
		//Delete backup gateway first
		gateway.GwName += "-hagw"
		log.Printf("[INFO] Deleting Aviatrix Backup Gateway [-hagw]: %#v", gateway)

		if isPublicSubnetFilteringGateway {
			err = client.DeletePublicSubnetFilteringGateway(gateway)
		} else {
			err = client.DeleteGateway(gateway)
		}

		if err != nil {
			return fmt.Errorf("failed to delete backup [-hgw] gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)

	log.Printf("[INFO] Deleting Aviatrix gateway: %#v", gateway)

	if isPublicSubnetFilteringGateway {
		err = client.DeletePublicSubnetFilteringGateway(gateway)
	} else {
		err = client.DeleteGateway(gateway)
	}
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Gateway: %s", err)
	}

	return nil
}

func checkPublicSubnetFilteringConfig(d *schema.ResourceData) error {
	var routeTables, haRouteTables []string
	for _, v := range d.Get("public_subnet_filtering_route_tables").(*schema.Set).List() {
		routeTables = append(routeTables, v.(string))
	}
	for _, v := range d.Get("public_subnet_filtering_ha_route_tables").(*schema.Set).List() {
		haRouteTables = append(haRouteTables, v.(string))
	}
	isPublicSubnetFilteringGw := d.Get("enable_public_subnet_filtering").(bool)
	// Public subnet filtering only supported for AWS and AWSGov
	if isPublicSubnetFilteringGw && !goaviatrix.IsCloudType(d.Get("cloud_type").(int), goaviatrix.AWS|goaviatrix.AWSGov) {
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
		if d.Get("enable_public_subnet_filtering").(bool) && !d.Get("enable_encrypt_volume").(bool) {
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
