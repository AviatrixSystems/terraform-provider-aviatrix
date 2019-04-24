package aviatrix

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

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
		},

		ResourcesMap: map[string]*schema.Resource{
			"aviatrix_account":                 resourceAccount(),
			"aviatrix_account_user":            resourceAccountUser(),
			"aviatrix_admin_email":             resourceAdminEmail(),
			"aviatrix_aws_peer":                resourceAWSPeer(),
			"aviatrix_aws_tgw":                 resourceAWSTgw(),
			"aviatrix_aws_tgw_vpc_attachment":  resourceAwsTgwVpcAttachment(),
			"aviatrix_customer_id":             resourceCustomerID(),
			"aviatrix_controller_config":       resourceControllerConfig(),
			"aviatrix_firewall":                resourceAviatrixFirewall(),
			"aviatrix_firewall_tag":            resourceAviatrixFirewallTag(),
			"aviatrix_fqdn":                    resourceAviatrixFQDN(),
			"aviatrix_gateway":                 resourceAviatrixGateway(),
			"aviatrix_site2cloud":              resourceAviatrixSite2Cloud(),
			"aviatrix_spoke_vpc":               resourceAviatrixSpokeVpc(),
			"aviatrix_trans_peer":              resourceTransPeer(),
			"aviatrix_transit_gateway_peering": resourceTransitGatewayPeering(),
			"aviatrix_transit_vpc":             resourceAviatrixTransitVpc(),
			"aviatrix_tunnel":                  resourceTunnel(),
			"aviatrix_version":                 resourceAviatrixVersion(),
			"aviatrix_vgw_conn":                resourceAviatrixVGWConn(),
			"aviatrix_vpn_profile":             resourceAviatrixProfile(),
			"aviatrix_vpn_user":                resourceAviatrixVPNUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"aviatrix_caller_identity": dataSourceAviatrixCallerIdentity(),
			"aviatrix_account":         dataSourceAviatrixAccount(),
			"aviatrix_gateway":         dataSourceAviatrixGateway(),
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

func envDefaultFuncAllowMissing(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		v := os.Getenv(k)
		return v, nil
	}
}

func aviatrixConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ControllerIP: d.Get("controller_ip").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}
	return config.Client()
}
