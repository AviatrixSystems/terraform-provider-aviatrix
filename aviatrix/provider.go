package aviatrix

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const supportedVersion = "5.2"

// Provider returns a schema.Provider for Aviatrix.
func Provider() terraform.ResourceProvider {
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
			"aviatrix_account":                 resourceAviatrixAccount(),
			"aviatrix_account_user":            resourceAviatrixAccountUser(),
			"aviatrix_arm_peer":                resourceAviatrixARMPeer(),
			"aviatrix_aws_peer":                resourceAviatrixAWSPeer(),
			"aviatrix_aws_tgw":                 resourceAviatrixAWSTgw(),
			"aviatrix_aws_tgw_vpc_attachment":  resourceAviatrixAwsTgwVpcAttachment(),
			"aviatrix_aws_tgw_vpn_conn":        resourceAviatrixAwsTgwVpnConn(),
			"aviatrix_controller_config":       resourceAviatrixControllerConfig(),
			"aviatrix_firenet":                 resourceAviatrixFireNet(),
			"aviatrix_firewall":                resourceAviatrixFirewall(),
			"aviatrix_firewall_instance":       resourceAviatrixFirewallInstance(),
			"aviatrix_firewall_tag":            resourceAviatrixFirewallTag(),
			"aviatrix_fqdn":                    resourceAviatrixFQDN(),
			"aviatrix_gateway":                 resourceAviatrixGateway(),
			"aviatrix_saml_endpoint":           resourceAviatrixSamlEndpoint(),
			"aviatrix_site2cloud":              resourceAviatrixSite2Cloud(),
			"aviatrix_spoke_gateway":           resourceAviatrixSpokeGateway(),
			"aviatrix_spoke_vpc":               resourceAviatrixSpokeVpc(),
			"aviatrix_trans_peer":              resourceAviatrixTransPeer(),
			"aviatrix_transit_gateway":         resourceAviatrixTransitGateway(),
			"aviatrix_transit_gateway_peering": resourceAviatrixTransitGatewayPeering(),
			"aviatrix_transit_vpc":             resourceAviatrixTransitVpc(),
			"aviatrix_tunnel":                  resourceAviatrixTunnel(),
			"aviatrix_vgw_conn":                resourceAviatrixVGWConn(),
			"aviatrix_vpc":                     resourceAviatrixVpc(),
			"aviatrix_vpn_profile":             resourceAviatrixProfile(),
			"aviatrix_vpn_user":                resourceAviatrixVPNUser(),
			"aviatrix_vpn_user_accelerator":    resourceAviatrixVPNUserAccelerator(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"aviatrix_account":                    dataSourceAviatrixAccount(),
			"aviatrix_caller_identity":            dataSourceAviatrixCallerIdentity(),
			"aviatrix_firenet_vendor_integration": dataSourceAviatrixFireNetVendorIntegration(),
			"aviatrix_gateway":                    dataSourceAviatrixGateway(),
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

	err = client.ControllerVersionValidation(supportedVersion)
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
