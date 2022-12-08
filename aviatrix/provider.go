package aviatrix

import (
	"errors"
	"os"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var supportedVersions = []string{"7.0"}

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
			"verify_ssl_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"path_to_ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignore_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration block with settings to ignore tags across all resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keys": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Tag keys to ignore across all resources.",
						},
						"key_prefixes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Tag key prefixes to ignore across all resources.",
						},
					},
				},
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"aviatrix_account":                                        resourceAviatrixAccount(),
			"aviatrix_account_user":                                   resourceAviatrixAccountUser(),
			"aviatrix_smart_group":                                    resourceAviatrixSmartGroup(),
			"aviatrix_aws_peer":                                       resourceAviatrixAWSPeer(),
			"aviatrix_aws_guard_duty":                                 resourceAviatrixAwsGuardDuty(),
			"aviatrix_aws_tgw":                                        resourceAviatrixAWSTgw(),
			"aviatrix_aws_tgw_connect":                                resourceAviatrixAwsTgwConnect(),
			"aviatrix_aws_tgw_connect_peer":                           resourceAviatrixAwsTgwConnectPeer(),
			"aviatrix_aws_tgw_directconnect":                          resourceAviatrixAWSTgwDirectConnect(),
			"aviatrix_aws_tgw_intra_domain_inspection":                resourceAviatrixAwsTgwIntraDomainInspection(),
			"aviatrix_aws_tgw_network_domain":                         resourceAviatrixAwsTgwNetworkDomain(),
			"aviatrix_aws_tgw_peering":                                resourceAviatrixAWSTgwPeering(),
			"aviatrix_aws_tgw_peering_domain_conn":                    resourceAviatrixAWSTgwPeeringDomainConn(),
			"aviatrix_aws_tgw_transit_gateway_attachment":             resourceAviatrixAwsTgwTransitGatewayAttachment(),
			"aviatrix_aws_tgw_vpc_attachment":                         resourceAviatrixAwsTgwVpcAttachment(),
			"aviatrix_aws_tgw_vpn_conn":                               resourceAviatrixAwsTgwVpnConn(),
			"aviatrix_azure_peer":                                     resourceAviatrixAzurePeer(),
			"aviatrix_azure_spoke_native_peering":                     resourceAviatrixAzureSpokeNativePeering(),
			"aviatrix_azure_vng_conn":                                 resourceAviatrixAzureVngConn(),
			"aviatrix_centralized_transit_firenet":                    resourceAviatrixCentralizedTransitFireNet(),
			"aviatrix_cloudn_registration":                            resourceAviatrixCloudnRegistration(),
			"aviatrix_cloudn_transit_gateway_attachment":              resourceAviatrixCloudnTransitGatewayAttachment(),
			"aviatrix_cloudwatch_agent":                               resourceAviatrixCloudwatchAgent(),
			"aviatrix_controller_bgp_max_as_limit_config":             resourceAviatrixControllerBgpMaxAsLimitConfig(),
			"aviatrix_controller_cert_domain_config":                  resourceAviatrixControllerCertDomainConfig(),
			"aviatrix_controller_config":                              resourceAviatrixControllerConfig(),
			"aviatrix_controller_email_config":                        resourceAviatrixControllerEmailConfig(),
			"aviatrix_controller_email_exception_notification_config": resourceAviatrixControllerEmailExceptionNotificationConfig(),
			"aviatrix_controller_gateway_keepalive_config":            resourceAviatrixControllerGatewayKeepaliveConfig(),
			"aviatrix_controller_private_mode_config":                 resourceAviatrixControllerPrivateModeConfig(),
			"aviatrix_controller_private_oob":                         resourceAviatrixControllerPrivateOob(),
			"aviatrix_controller_security_group_management_config":    resourceAviatrixControllerSecurityGroupManagementConfig(),
			"aviatrix_copilot_association":                            resourceAviatrixCopilotAssociation(),
			"aviatrix_copilot_security_group_management_config":       resourceAviatrixCopilotSecurityGroupManagementConfig(),
			"aviatrix_datadog_agent":                                  resourceAviatrixDatadogAgent(),
			"aviatrix_device_interface_config":                        resourceAviatrixDeviceInterfaceConfig(),
			"aviatrix_distributed_firewalling_config":                 resourceAviatrixDistributedFirewallingConfig(),
			"aviatrix_distributed_firewalling_intra_vpc":              resourceAviatrixDistributedFirewallingIntraVpc(),
			"aviatrix_distributed_firewalling_policy_list":            resourceAviatrixDistributedFirewallingPolicyList(),
			"aviatrix_edge_csp":                                       resourceAviatrixEdgeCSP(),
			"aviatrix_edge_spoke":                                     resourceAviatrixEdgeSpoke(),
			"aviatrix_edge_spoke_external_device_conn":                resourceAviatrixEdgeSpokeExternalDeviceConn(),
			"aviatrix_edge_spoke_transit_attachment":                  resourceAviatrixEdgeSpokeTransitAttachment(),
			"aviatrix_filebeat_forwarder":                             resourceAviatrixFilebeatForwarder(),
			"aviatrix_firenet":                                        resourceAviatrixFireNet(),
			"aviatrix_firewall":                                       resourceAviatrixFirewall(),
			"aviatrix_firewall_instance":                              resourceAviatrixFirewallInstance(),
			"aviatrix_firewall_instance_association":                  resourceAviatrixFirewallInstanceAssociation(),
			"aviatrix_firewall_management_access":                     resourceAviatrixFirewallManagementAccess(),
			"aviatrix_firewall_policy":                                resourceAviatrixFirewallPolicy(),
			"aviatrix_firewall_tag":                                   resourceAviatrixFirewallTag(),
			"aviatrix_fqdn":                                           resourceAviatrixFQDN(),
			"aviatrix_fqdn_global_config":                             resourceAviatrixFQDNGlobalConfig(),
			"aviatrix_fqdn_pass_through":                              resourceAviatrixFQDNPassThrough(),
			"aviatrix_fqdn_tag_rule":                                  resourceAviatrixFQDNTagRule(),
			"aviatrix_gateway":                                        resourceAviatrixGateway(),
			"aviatrix_gateway_certificate_config":                     resourceAviatrixGatewayCertificateConfig(),
			"aviatrix_gateway_dnat":                                   resourceAviatrixGatewayDNat(),
			"aviatrix_gateway_snat":                                   resourceAviatrixGatewaySNat(),
			"aviatrix_geo_vpn":                                        resourceAviatrixGeoVPN(),
			"aviatrix_netflow_agent":                                  resourceAviatrixNetflowAgent(),
			"aviatrix_periodic_ping":                                  resourceAviatrixPeriodicPing(),
			"aviatrix_private_mode_lb":                                resourceAviatrixPrivateModeLb(),
			"aviatrix_private_mode_multicloud_endpoint":               resourceAviatrixPrivateModeMulticloudEndpoint(),
			"aviatrix_proxy_config":                                   resourceAviatrixProxyConfig(),
			"aviatrix_rbac_group":                                     resourceAviatrixRbacGroup(),
			"aviatrix_rbac_group_access_account_attachment":           resourceAviatrixRbacGroupAccessAccountAttachment(),
			"aviatrix_rbac_group_permission_attachment":               resourceAviatrixRbacGroupPermissionAttachment(),
			"aviatrix_rbac_group_user_attachment":                     resourceAviatrixRbacGroupUserAttachment(),
			"aviatrix_remote_syslog":                                  resourceAviatrixRemoteSyslog(),
			"aviatrix_saml_endpoint":                                  resourceAviatrixSamlEndpoint(),
			"aviatrix_segmentation_network_domain":                    resourceAviatrixSegmentationNetworkDomain(),
			"aviatrix_segmentation_network_domain_association":        resourceAviatrixSegmentationNetworkDomainAssociation(),
			"aviatrix_segmentation_network_domain_connection_policy":  resourceAviatrixSegmentationNetworkDomainConnectionPolicy(),
			"aviatrix_site2cloud":                                     resourceAviatrixSite2Cloud(),
			"aviatrix_site2cloud_ca_cert_tag":                         resourceAviatrixSite2CloudCaCertTag(),
			"aviatrix_splunk_logging":                                 resourceAviatrixSplunkLogging(),
			"aviatrix_spoke_gateway":                                  resourceAviatrixSpokeGateway(),
			"aviatrix_spoke_ha_gateway":                               resourceAviatrixSpokeHaGateway(),
			"aviatrix_spoke_gateway_subnet_group":                     resourceAviatrixSpokeGatewaySubnetGroup(),
			"aviatrix_spoke_external_device_conn":                     resourceAviatrixSpokeExternalDeviceConn(),
			"aviatrix_spoke_transit_attachment":                       resourceAviatrixSpokeTransitAttachment(),
			"aviatrix_sumologic_forwarder":                            resourceAviatrixSumologicForwarder(),
			"aviatrix_transit_external_device_conn":                   resourceAviatrixTransitExternalDeviceConn(),
			"aviatrix_trans_peer":                                     resourceAviatrixTransPeer(),
			"aviatrix_transit_firenet_policy":                         resourceAviatrixTransitFireNetPolicy(),
			"aviatrix_transit_gateway":                                resourceAviatrixTransitGateway(),
			"aviatrix_transit_gateway_peering":                        resourceAviatrixTransitGatewayPeering(),
			"aviatrix_tunnel":                                         resourceAviatrixTunnel(),
			"aviatrix_vgw_conn":                                       resourceAviatrixVGWConn(),
			"aviatrix_vpc":                                            resourceAviatrixVpc(),
			"aviatrix_vpn_cert_download":                              resourceAviatrixVPNCertDownload(),
			"aviatrix_vpn_profile":                                    resourceAviatrixProfile(),
			"aviatrix_vpn_user":                                       resourceAviatrixVPNUser(),
			"aviatrix_vpn_user_accelerator":                           resourceAviatrixVPNUserAccelerator(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"aviatrix_account":                          dataSourceAviatrixAccount(),
			"aviatrix_caller_identity":                  dataSourceAviatrixCallerIdentity(),
			"aviatrix_device_interfaces":                dataSourceAviatrixDeviceInterfaces(),
			"aviatrix_firenet":                          dataSourceAviatrixFireNet(),
			"aviatrix_firenet_firewall_manager":         dataSourceAviatrixFireNetFirewallManager(),
			"aviatrix_firenet_vendor_integration":       dataSourceAviatrixFireNetVendorIntegration(),
			"aviatrix_gateway":                          dataSourceAviatrixGateway(),
			"aviatrix_gateway_image":                    dataSourceAviatrixGatewayImage(),
			"aviatrix_network_domains":                  dataSourceAviatrixNetworkDomains(),
			"aviatrix_spoke_gateway":                    dataSourceAviatrixSpokeGateway(),
			"aviatrix_spoke_gateways":                   dataSourceAviatrixSpokeGateways(),
			"aviatrix_spoke_gateway_inspection_subnets": dataSourceAviatrixSpokeGatewayInspectionSubnets(),
			"aviatrix_transit_gateway":                  dataSourceAviatrixTransitGateway(),
			"aviatrix_transit_gateways":                 dataSourceAviatrixTransitGateways(),
			"aviatrix_vpc":                              dataSourceAviatrixVpc(),
			"aviatrix_vpc_tracker":                      dataSourceAviatrixVpcTracker(),
			"aviatrix_firewall":                         dataSourceAviatrixFirewall(),
			"aviatrix_firewall_instance_images":         dataSourceAviatrixFirewallInstanceImages(),
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
		VerifyCert:   d.Get("verify_ssl_certificate").(bool),
		PathToCACert: d.Get("path_to_ca_certificate").(string),
		IgnoreTags:   expandProviderIgnoreTags(d.Get("ignore_tags").([]interface{})),
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
		VerifyCert:   d.Get("verify_ssl_certificate").(bool),
		PathToCACert: d.Get("path_to_ca_certificate").(string),
		IgnoreTags:   expandProviderIgnoreTags(d.Get("ignore_tags").([]interface{})),
	}

	return config.Client()
}

func expandProviderIgnoreTags(l []interface{}) *goaviatrix.IgnoreTagsConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	ignoreConfig := &goaviatrix.IgnoreTagsConfig{}
	m := l[0].(map[string]interface{})

	if v, ok := m["keys"].(*schema.Set); ok {
		ignoreConfig.Keys = goaviatrix.NewIgnoreTags(v.List())
	}

	if v, ok := m["key_prefixes"].(*schema.Set); ok {
		ignoreConfig.KeyPrefixes = goaviatrix.NewIgnoreTags(v.List())
	}

	return ignoreConfig
}
