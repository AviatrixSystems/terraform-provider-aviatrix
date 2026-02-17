package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixGatewayRead,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway name. This can be used for getting gateway.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account name. This account will be used to launch Aviatrix gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of legacy VPC/Vnet to be connected.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region where gateway is launched.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Size of Gateway Instance.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A VPC Network address range selected from one of the available network ranges.",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability Zone. Only set for Azure and Public Subnet Filtering gateway",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Gateway. Required if insane_mode is set.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Source NAT for this container.",
			},
			"vpn_access": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable user access through VPN to this container.",
			},
			"vpn_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPN CIDR block for the container.",
			},
			"enable_elb": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specify whether to enable ELB or not.",
			},
			"elb_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A name for the ELB that is created.",
			},
			"vpn_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Elb protocol for VPN gateway with elb enabled. Only supports AWS provider. Valid values: 'TCP', 'UDP'. If not specified, 'TCP'' will be used.",
			},
			"split_tunnel": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specify split tunnel mode.",
			},
			"max_vpn_conn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Maximum connection of VPN access.",
			},
			"name_servers": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of DNS servers used to resolve domain names by " +
					"a connected VPN user when Split Tunnel Mode is enabled.",
			},
			"search_domains": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of domain names that will use the NameServer " +
					"when a specific name is not in the destination when Split Tunnel Mode is enabled.",
			},
			"additional_cidrs": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "A list of destination CIDR ranges that will also go through the VPN tunnel " +
					"when Split Tunnel Mode is enabled.",
			},
			"otp_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Two step authentication mode.",
			},
			"saml_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This field indicates whether to enable SAML or not.",
			},
			"enable_vpn_nat": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This field indicates whether to enable VPN NAT or not. Only supported for VPN gateway. Valid values: true, false. Default value: true.",
			},
			"okta_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for Okta auth mode.",
			},
			"okta_username_suffix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username suffix for Okta auth mode.",
			},
			"duo_integration_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Integration key for DUO auth mode.",
			},
			"duo_api_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API hostname for DUO auth mode.",
			},
			"duo_push_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Push mode for DUO auth.",
			},
			"enable_ldap": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether enable LDAP or not.",
			},
			"ldap_server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LDAP server address. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_bind_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LDAP bind DN. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_base_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LDAP base DN. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_username_attribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LDAP user attribute. Required: Yes if enable_ldap is 'yes'.",
			},
			"peering_ha_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Required to create peering ha gateway if cloud_type = 1 or 8 (AWS or Azure)",
			},
			"peering_ha_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone information for creating Peering HA Gateway. Required to create peering ha gateway if cloud_type = 4 (gcp)",
			},
			"peering_ha_insane_mode_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AZ of subnet being created for Insane Mode Peering HA Gateway. Required if insane_mode is set.",
			},
			"peering_ha_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA peering instance.",
			},
			"peering_ha_gw_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peering HA Gateway Size.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set to true if this feature is desired.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "When value is false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Insane Mode for Gateway. Valid values: true, false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable vpc_dns_server for Gateway. Valid values: true, false.",
			},
			"enable_designated_gateway": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable 'designated_gateway' feature for Gateway. Only supports AWS. Valid values: true, false.",
			},
			"additional_cidrs_designated_gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.",
			},
			"elb_dns_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ELB DNS Name.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Gateway created.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the gateway.",
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
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A map of tags assigned to the gateway.",
			},
			"tunnel_detection_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IPSec tunnel down detection time for the gateway.",
			},
			"availability_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability domain for OCI.",
			},
			"fault_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fault domain for OCI.",
			},
			"peering_ha_availability_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peering HA availability domain for OCI.",
			},
			"peering_ha_fault_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peering HA fault domain for OCI.",
			},
			"software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software version of the gateway.",
			},
			"peering_ha_software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software version of the HA gateway.",
			},
			"image_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image version of the gateway.",
			},
			"peering_ha_image_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image version of the HA gateway.",
			},
			"enable_monitor_gateway_subnets": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable monitor gateway subnets.",
			},
			"monitor_exclude_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A set of monitored instance ids. Only set when 'enable_monitor_gateway_subnets' = true",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Typed value when modifying idle_timeout. If it's -1, this feature is disabled.",
			},
			"renegotiation_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Typed value when modifying renegotiation_interval. If it's -1, this feature is disabled.",
			},
			"fqdn_lan_interface": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "FQDN gateway lan interface id.",
			},
			"fqdn_lan_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "FQDN gateway lan interface cidr.",
			},
			"fqdn_lan_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LAN VPC ID. Only used for GCP FQDN Gateway.",
			},
			"enable_public_subnet_filtering": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Create a [Public Subnet Filtering gateway](https://docs.aviatrix.com/HowTos/public_subnet_filtering_faq.html).",
			},
			"public_subnet_filtering_route_tables": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Route tables whose associated public subnets are protected. Only set when `enable_public_subnet_filtering` attribute is true.",
			},
			"public_subnet_filtering_ha_route_tables": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Route tables whose associated public subnets are protected for the HA PSF gateway. Only set when enable_public_subnet_filtering and peering_ha_subnet are set.",
			},
			"public_subnet_filtering_guard_duty_enforced": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enforce Guard Duty IP blocking. Only set when `enable_public_subnet_filtering` attribute is true.",
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable jumbo frame support for Gateway.",
			},
			"enable_spot_instance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable spot instance. NOT supported for production deployment.",
			},
			"spot_price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Price for spot instance. NOT supported for production deployment.",
			},
			"azure_eip_name_resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the public IP address and its resource group in Azure to assign to this Gateway.",
			},
			"peering_ha_azure_eip_name_resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the public IP address and its resource group in Azure to assign to the Peering HA Gateway.",
			},
			"peering_ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peering HA security group used for the gateway.",
			},
		},
	}
}

func dataSourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		GwName: getString(d, "gw_name"),
	}

	if getString(d, "account_name") != "" {
		gateway.AccountName = getString(d, "account_name")
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Gateway: %w", err)
	}
	if gw != nil {
		mustSet(d, "cloud_type", gw.CloudType)
		mustSet(d, "account_name", gw.AccountName)
		mustSet(d, "gw_name", gw.GwName)

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			mustSet(d, "vpc_id", strings.Split(gw.VpcID, "~~")[0])
			mustSet(d, "vpc_reg", gw.VpcRegion)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			mustSet(d, "vpc_id", gw.VpcID)
			mustSet(d, "vpc_reg", gw.GatewayZone)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			mustSet(d, "vpc_id", gw.VpcID)
			mustSet(d, "vpc_reg", gw.VpcRegion)
		}

		_, zoneIsSet := d.GetOk("zone")
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) && zoneIsSet && gw.GatewayZone != "AvailabilitySet" {
			mustSet(d, "zone", "az-"+gw.GatewayZone)
		}
		if gw.NewZone != "" {
			mustSet(d, "zone", gw.NewZone)
		}
		mustSet(d, "subnet", gw.VpcNet)

		if gw.EnableNat == "yes" {
			if gw.SnatMode == "primary" {
				mustSet(d, "single_ip_snat", true)
			} else {
				mustSet(d, "single_ip_snat", false)
			}
		} else {
			mustSet(d, "single_ip_snat", false)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes) {
			if gw.AllocateNewEipRead {
				mustSet(d, "allocate_new_eip", true)
			} else {
				mustSet(d, "allocate_new_eip", false)
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) {
			mustSet(d, "allocate_new_eip", true)
		}

		if gw.EnableDesignatedGateway == "Yes" || gw.EnableDesignatedGateway == "yes" {
			mustSet(d, "enable_designated_gateway", true)
			mustSet(d, "additional_cidrs_designated_gateway", gw.AdditionalCidrsDesignatedGw)
		} else {
			mustSet(d, "enable_designated_gateway", false)
			mustSet(d, "additional_cidrs_designated_gateway", "")
		}

		if gw.EnableLdap {
			mustSet(d, "enable_ldap", true)
		} else {
			mustSet(d, "enable_ldap", false)
		}

		gwDetail, err := client.GetGatewayDetail(gateway)
		if err != nil {
			return fmt.Errorf("couldn't get Detail info for VPN gateway: %s due to: %w", gateway.GwName, err)
		}
		if gw.VpnStatus != "" {
			if gw.VpnStatus == "disabled" {
				mustSet(d, "vpn_access", false)
				mustSet(d, "enable_vpn_nat", true)
				mustSet(d, "vpn_protocol", "")
			} else if gw.VpnStatus == "enabled" {
				mustSet(d, "vpn_access", true)
				gateway.VpcID = getString(d, "vpc_id")
				if gwDetail.VpnNat {
					mustSet(d, "enable_vpn_nat", true)
				} else {
					mustSet(d, "enable_vpn_nat", false)
				}

				if gw.ElbState == "enabled" {
					if gwDetail.Elb.VpnProtocol == "udp" || gwDetail.Elb.VpnProtocol == "UDP" {
						mustSet(d, "vpn_protocol", "UDP")
					} else {
						mustSet(d, "vpn_protocol", "TCP")
					}
				} else {
					mustSet(d, "vpn_protocol", "")
				}
			}
		}

		vpnAccess := getBool(d, "vpn_access")
		if !vpnAccess {
			mustSet(d, "split_tunnel", true)
			mustSet(d, "max_vpn_conn", "")
		} else {
			if gw.SplitTunnel == "yes" {
				mustSet(d, "split_tunnel", true)
			} else {
				mustSet(d, "split_tunnel", false)
			}
			mustSet(d, "max_vpn_conn", gw.MaxConn)
		}
		mustSet(d, "vpn_cidr", gw.VpnCidr)

		if gw.ElbState == "enabled" {
			mustSet(d, "enable_elb", true)
			mustSet(d, "elb_name", gw.ElbName)
			mustSet(d, "elb_dns_name", gw.ElbDNSName)
		} else {
			mustSet(d, "enable_elb", false)
			mustSet(d, "elb_name", "")
		}

		if gw.SamlEnabled == "yes" {
			mustSet(d, "saml_enabled", true)
		} else {
			mustSet(d, "saml_enabled", false)
		}

		if gw.AuthMethod == "duo_auth" || gw.AuthMethod == "duo_auth+LDAP" {
			mustSet(d, "otp_mode", "2")
		} else if gw.AuthMethod == "okta_auth" {
			mustSet(d, "otp_mode", "3")
		} else {
			mustSet(d, "otp_mode", "")
		}
		mustSet(d, "okta_url", gw.OktaURL)
		mustSet(d, "okta_username_suffix", gw.OktaUsernameSuffix)
		mustSet(d, "duo_integration_key", gw.DuoIntegrationKey)
		mustSet(d, "duo_api_hostname", gw.DuoAPIHostname)
		mustSet(d, "duo_push_mode", gw.DuoPushMode)
		mustSet(d, "ldap_server", gw.LdapServer)
		mustSet(d, "ldap_bind_dn", gw.LdapBindDn)
		mustSet(d, "ldap_base_dn", gw.LdapBaseDn)
		mustSet(d, "ldap_username_attribute", gw.LdapUserAttr)

		if gw.NewZone != "" {
			mustSet(d, "zone", gw.NewZone)
		}

		if gw.SingleAZ != "" {
			if gw.SingleAZ == "yes" {
				mustSet(d, "single_az_ha", true)
			} else {
				mustSet(d, "single_az_ha", false)
			}
		}
		mustSet(d, "enable_encrypt_volume", gw.EnableEncryptVolume)

		if gw.GwSize != "" {
			mustSet(d, "gw_size", gw.GwSize)
		} else {
			if gw.VpcSize != "" {
				mustSet(d, "gw_size", gw.VpcSize)
			}
		}
		mustSet(d, "public_ip", gw.PublicIP)
		mustSet(d, "cloud_instance_id", gw.CloudnGatewayInstID)
		mustSet(d, "public_dns_server", gw.PublicDnsServer)
		mustSet(d, "security_group_id", gw.GwSecurityGroupID)
		mustSet(d, "private_ip", gw.PrivateIP)
		mustSet(d, "image_version", gw.ImageVersion)
		mustSet(d, "software_version", gw.SoftwareVersion)

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

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.AliCloudRelatedCloudTypes) && gw.EnableVpcDnsServer == "Enabled" {
			mustSet(d, "enable_vpc_dns_server", true)
		} else {
			mustSet(d, "enable_vpc_dns_server", false)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
			mustSet(d, "availability_domain", gw.GatewayZone)
			mustSet(d, "fault_domain", gw.FaultDomain)
		}

		if !gw.IsPsfGateway {
			mustSet(d, "enable_public_subnet_filtering", false)
			mustSet(d, "public_subnet_filtering_route_tables", nil)
			mustSet(d, "public_subnet_filtering_ha_route_tables", nil)
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
				err := d.Set("public_subnet_filtering_ha_route_tables", nil)
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

		peeringHaGateway := &goaviatrix.Gateway{
			AccountName: getString(d, "account_name"),
			GwName:      getString(d, "gw_name") + "-hagw",
		}
		gwHaGw, _ := client.GetGateway(peeringHaGateway)
		if gwHaGw != nil {
			mustSet(d, "peering_ha_cloud_instance_id", gwHaGw.CloudnGatewayInstID)
			mustSet(d, "peering_ha_gw_name", gwHaGw.GwName)
			mustSet(d, "peering_ha_public_ip", gwHaGw.PublicIP)
			mustSet(d, "peering_ha_gw_size", gwHaGw.GwSize)
			mustSet(d, "peering_ha_private_ip", gwHaGw.PrivateIP)
			mustSet(d, "peering_ha_image_version", gwHaGw.ImageVersion)
			mustSet(d, "peering_ha_software_version", gwHaGw.SoftwareVersion)
			mustSet(d, "peering_ha_security_group_id", gw.HaGw.GwSecurityGroupID)
			if goaviatrix.IsCloudType(gwHaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
				azureEip := strings.Split(gw.HaGw.ReuseEip, ":")
				if len(azureEip) == 3 {
					mustSet(d, "peering_ha_azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
				} else {
					log.Printf("[WARN] could not get Azure EIP name and resource group for the Peering HA Gateway %s", gw.GwName)
				}
			}
			if !gw.IsPsfGateway {
				// For PSF gateway, peering_ha_subnet and peering_ha_zone are set above
				if goaviatrix.IsCloudType(gwHaGw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
					mustSet(d, "peering_ha_subnet", gwHaGw.VpcNet)
					mustSet(d, "peering_ha_zone", "")
					if gwHaGw.InsaneMode == "yes" {
						mustSet(d, "peering_ha_insane_mode_az", gwHaGw.GatewayZone)
					}
				} else if goaviatrix.IsCloudType(gwHaGw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
					mustSet(d, "peering_ha_subnet", gwHaGw.VpcNet)
					if _, haZoneIsSet := d.GetOk("peering_ha_zone"); haZoneIsSet {
						if gw.GatewayZone != "AvailabilitySet" {
							mustSet(d, "peering_ha_zone", "az-"+gw.GatewayZone)
						}
					}
				} else if goaviatrix.IsCloudType(gwHaGw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
					mustSet(d, "peering_ha_subnet", gw.HaGw.VpcNet)
					mustSet(d, "peering_ha_zone", gwHaGw.GatewayZone)
				} else if goaviatrix.IsCloudType(gwHaGw.CloudType, goaviatrix.AliCloudRelatedCloudTypes) {
					mustSet(d, "peering_ha_subnet", gwHaGw.VpcNet)
					mustSet(d, "peering_ha_zone", "")
				} else if goaviatrix.IsCloudType(gwHaGw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
					mustSet(d, "peering_ha_subnet", gwHaGw.VpcNet)
					mustSet(d, "peering_ha_zone", "")
					mustSet(d, "peering_ha_availability_domain", gwHaGw.GatewayZone)
					mustSet(d, "peering_ha_fault_domain", gwHaGw.FaultDomain)
				}
			}
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			tags := &goaviatrix.Tags{
				ResourceType: "gw",
				ResourceName: getString(d, "gw_name"),
				CloudType:    gw.CloudType,
			}

			_, err := client.GetTags(tags)
			if err != nil {
				log.Printf("[WARN] Failed to get tags for gateway %s: %v", tags.ResourceName, err)
			}
			if len(tags.Tags) > 0 {
				if err := d.Set("tags", tags.Tags); err != nil {
					log.Printf("[WARN] Error setting tags for gateway %s: %v", tags.ResourceName, err)
				}
			}
		}

		if gw.VpnStatus == "enabled" && gw.SplitTunnel == "yes" {
			splitTunnel := &goaviatrix.SplitTunnel{
				VpcID: gw.VpcID,
			}
			if gw.ElbState == "enabled" {
				splitTunnel.ElbName = gw.ElbName
			} else {
				splitTunnel.ElbName = gw.GwName
			}
			splitTunnel1, _ := client.GetSplitTunnel(splitTunnel)
			if splitTunnel1 != nil {
				mustSet(d, "name_servers", splitTunnel1.NameServers)
				mustSet(d, "search_domains", splitTunnel1.SearchDomains)
				mustSet(d, "additional_cidrs", splitTunnel1.AdditionalCidrs)
			}
		} else {
			mustSet(d, "name_servers", "")
			mustSet(d, "search_domains", "")
			mustSet(d, "additional_cidrs", "")
		}
		mustSet(d, "tunnel_detection_time", gw.TunnelDetectionTime)
		mustSet(d, "enable_jumbo_frame", gw.JumboFrame)
		mustSet(d, "enable_monitor_gateway_subnets", gw.MonitorSubnetsAction == "enable")
		if err := d.Set("monitor_exclude_list", gw.MonitorExcludeGWList); err != nil {
			return fmt.Errorf("setting 'monitor_exclude_list' to state: %w", err)
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

		fqdnLanCidr, ok := gw.ArmFqdnLanCidr[gw.GwName]
		if ok && goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			mustSet(d, "fqdn_lan_cidr", fqdnLanCidr)
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			mustSet(d, "fqdn_lan_vpc_id", gw.BundleVpcInfo.LAN.VpcID)
			mustSet(d, "fqdn_lan_cidr", strings.Split(gw.BundleVpcInfo.LAN.Subnet, "~~")[0])
		} else {
			mustSet(d, "fqdn_lan_cidr", "")
		}

		if gw.EnableSpotInstance {
			mustSet(d, "enable_spot_instance", true)
			mustSet(d, "spot_price", gw.SpotPrice)
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			azureEip := strings.Split(gw.ReuseEip, ":")
			if len(azureEip) == 3 {
				mustSet(d, "azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
			} else {
				log.Printf("[WARN] could not get Azure EIP name and resource group for the Gateway %s", gw.GwName)
			}
		}
	}

	d.SetId(gateway.GwName)
	return nil
}
