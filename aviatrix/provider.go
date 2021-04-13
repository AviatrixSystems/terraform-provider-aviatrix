package aviatrix

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var supportedVersions = []string{"6.3"}

// Provider returns a schema.Provider for Aviatrix.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"controller_ip": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("AVIATRIX_CONTROLLER_IP"),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("AVIATRIX_USERNAME"),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("AVIATRIX_PASSWORD"),
			},
			"skip_version_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"aviatrix_account":                                        resourceAviatrixAccount(),
			"aviatrix_account_user":                                   resourceAviatrixAccountUser(),
			"aviatrix_arm_peer":                                       resourceAviatrixARMPeer(),
			"aviatrix_aws_peer":                                       resourceAviatrixAWSPeer(),
			"aviatrix_aws_guard_duty":                                 resourceAviatrixAwsGuardDuty(),
			"aviatrix_aws_tgw":                                        resourceAviatrixAWSTgw(),
			"aviatrix_aws_tgw_connect":                                resourceAviatrixAwsTgwConnect(),
			"aviatrix_aws_tgw_connect_peer":                           resourceAviatrixAwsTgwConnectPeer(),
			"aviatrix_aws_tgw_directconnect":                          resourceAviatrixAWSTgwDirectConnect(),
			"aviatrix_aws_tgw_peering":                                resourceAviatrixAWSTgwPeering(),
			"aviatrix_aws_tgw_peering_domain_conn":                    resourceAviatrixAWSTgwPeeringDomainConn(),
			"aviatrix_aws_tgw_security_domain":                        resourceAviatrixAwsTgwSecurityDomain(),
			"aviatrix_aws_tgw_transit_gateway_attachment":             resourceAviatrixAwsTgwTransitGatewayAttachment(),
			"aviatrix_aws_tgw_vpc_attachment":                         resourceAviatrixAwsTgwVpcAttachment(),
			"aviatrix_aws_tgw_vpn_conn":                               resourceAviatrixAwsTgwVpnConn(),
			"aviatrix_azure_peer":                                     resourceAviatrixAzurePeer(),
			"aviatrix_azure_spoke_native_peering":                     resourceAviatrixAzureSpokeNativePeering(),
			"aviatrix_azure_vng_conn":                                 resourceAviatrixAzureVngConn(),
			"aviatrix_cloudn_transit_gateway_attachment":              resourceAviatrixCloudnTransitGatewayAttachment(),
			"aviatrix_cloudwatch_agent":                               resourceAviatrixCloudwatchAgent(),
			"aviatrix_controller_bgp_max_as_limit_config":             resourceAviatrixControllerBgpMaxAsLimitConfig(),
			"aviatrix_controller_cert_domain_config":                  resourceAviatrixControllerCertDomainConfig(),
			"aviatrix_controller_config":                              resourceAviatrixControllerConfig(),
			"aviatrix_controller_email_exception_notification_config": resourceAviatrixControllerEmailExceptionNotificationConfig(),
			"aviatrix_controller_private_oob":                         resourceAviatrixControllerPrivateOob(),
			"aviatrix_copilot_association":                            resourceAviatrixCopilotAssociation(),
			"aviatrix_datadog_agent":                                  resourceAviatrixDatadogAgent(),
			"aviatrix_device_aws_tgw_attachment":                      resourceAviatrixDeviceAwsTgwAttachment(),
			"aviatrix_device_interface_config":                        resourceAviatrixDeviceInterfaceConfig(),
			"aviatrix_device_registration":                            resourceAviatrixDeviceRegistration(),
			"aviatrix_device_tag":                                     resourceAviatrixDeviceTag(),
			"aviatrix_device_transit_gateway_attachment":              resourceAviatrixDeviceTransitGatewayAttachment(),
			"aviatrix_device_virtual_wan_attachment":                  resourceAviatrixDeviceVirtualWanAttachment(),
			"aviatrix_filebeat_forwarder":                             resourceAviatrixFilebeatForwarder(),
			"aviatrix_firenet":                                        resourceAviatrixFireNet(),
			"aviatrix_firewall":                                       resourceAviatrixFirewall(),
			"aviatrix_firewall_instance":                              resourceAviatrixFirewallInstance(),
			"aviatrix_firewall_instance_association":                  resourceAviatrixFirewallInstanceAssociation(),
			"aviatrix_firewall_management_access":                     resourceAviatrixFirewallManagementAccess(),
			"aviatrix_firewall_policy":                                resourceAviatrixFirewallPolicy(),
			"aviatrix_firewall_tag":                                   resourceAviatrixFirewallTag(),
			"aviatrix_fqdn":                                           resourceAviatrixFQDN(),
			"aviatrix_fqdn_pass_through":                              resourceAviatrixFQDNPassThrough(),
			"aviatrix_fqdn_tag_rule":                                  resourceAviatrixFQDNTagRule(),
			"aviatrix_gateway":                                        resourceAviatrixGateway(),
			"aviatrix_gateway_certificate_config":                     resourceAviatrixGatewayCertificateConfig(),
			"aviatrix_gateway_dnat":                                   resourceAviatrixGatewayDNat(),
			"aviatrix_gateway_snat":                                   resourceAviatrixGatewaySNat(),
			"aviatrix_geo_vpn":                                        resourceAviatrixGeoVPN(),
			"aviatrix_netflow_agent":                                  resourceAviatrixNetflowAgent(),
			"aviatrix_periodic_ping":                                  resourceAviatrixPeriodicPing(),
			"aviatrix_proxy_config":                                   resourceAviatrixProxyConfig(),
			"aviatrix_rbac_group":                                     resourceAviatrixRbacGroup(),
			"aviatrix_rbac_group_access_account_attachment":           resourceAviatrixRbacGroupAccessAccountAttachment(),
			"aviatrix_rbac_group_permission_attachment":               resourceAviatrixRbacGroupPermissionAttachment(),
			"aviatrix_rbac_group_user_attachment":                     resourceAviatrixRbacGroupUserAttachment(),
			"aviatrix_remote_syslog":                                  resourceAviatrixRemoteSyslog(),
			"aviatrix_saml_endpoint":                                  resourceAviatrixSamlEndpoint(),
			"aviatrix_segmentation_security_domain":                   resourceAviatrixSegmentationSecurityDomain(),
			"aviatrix_segmentation_security_domain_association":       resourceAviatrixSegmentationSecurityDomainAssociation(),
			"aviatrix_segmentation_security_domain_connection_policy": resourceAviatrixSegmentationSecurityDomainConnectionPolicy(),
			"aviatrix_site2cloud":                                     resourceAviatrixSite2Cloud(),
			"aviatrix_splunk_logging":                                 resourceAviatrixSplunkLogging(),
			"aviatrix_spoke_gateway":                                  resourceAviatrixSpokeGateway(),
			"aviatrix_spoke_transit_attachment":                       resourceAviatrixSpokeTransitAttachment(),
			"aviatrix_spoke_vpc":                                      resourceAviatrixSpokeVpc(),
			"aviatrix_sumologic_forwarder":                            resourceAviatrixSumologicForwarder(),
			"aviatrix_transit_external_device_conn":                   resourceAviatrixTransitExternalDeviceConn(),
			"aviatrix_trans_peer":                                     resourceAviatrixTransPeer(),
			"aviatrix_transit_firenet_policy":                         resourceAviatrixTransitFireNetPolicy(),
			"aviatrix_transit_gateway":                                resourceAviatrixTransitGateway(),
			"aviatrix_transit_gateway_peering":                        resourceAviatrixTransitGatewayPeering(),
			"aviatrix_transit_vpc":                                    resourceAviatrixTransitVpc(),
			"aviatrix_tunnel":                                         resourceAviatrixTunnel(),
			"aviatrix_vgw_conn":                                       resourceAviatrixVGWConn(),
			"aviatrix_vpc":                                            resourceAviatrixVpc(),
			"aviatrix_vpn_cert_download":                              resourceAviatrixVPNCertDownload(),
			"aviatrix_vpn_profile":                                    resourceAviatrixProfile(),
			"aviatrix_vpn_user":                                       resourceAviatrixVPNUser(),
			"aviatrix_vpn_user_accelerator":                           resourceAviatrixVPNUserAccelerator(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"aviatrix_account":                    dataSourceAviatrixAccount(),
			"aviatrix_caller_identity":            dataSourceAviatrixCallerIdentity(),
			"aviatrix_firenet":                    dataSourceAviatrixFireNet(),
			"aviatrix_firenet_vendor_integration": dataSourceAviatrixFireNetVendorIntegration(),
			"aviatrix_gateway":                    dataSourceAviatrixGateway(),
			"aviatrix_spoke_gateway":              dataSourceAviatrixSpokeGateway(),
			"aviatrix_transit_gateway":            dataSourceAviatrixTransitGateway(),
			"aviatrix_vpc":                        dataSourceAviatrixVpc(),
			"aviatrix_vpc_tracker":                dataSourceAviatrixVpcTracker(),
			"aviatrix_firewall":                   dataSourceAviatrixFirewall(),
		},
		ConfigureFunc: aviatrixConfigure,
	}
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}

func aviatrixConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ControllerIP: d.Get("controller_ip").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}

	skipVersionValidation := d.Get("skip_version_validation").(bool)
	if skipVersionValidation {
		return config.Client()
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	err = client.ControllerVersionValidation(supportedVersions)
	if err != nil {
		return nil, errors.New("controller version validation failed: " + err.Error())
	}

	return client, nil
}

func aviatrixConfigureWithoutVersionValidation(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ControllerIP: d.Get("controller_ip").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}

	return config.Client()
}
