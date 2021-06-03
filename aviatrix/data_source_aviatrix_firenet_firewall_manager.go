package aviatrix

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixFireNetFirewallManager() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixFireNetFirewallManagerRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FireNet gateway name.",
			},
			"vendor_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Vendor type. Valid values: 'Generic' and 'Palo Alto Networks Panorama'.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The public IP address of the Panorama instance.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Panorama login name for API calls from the Controller.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Panorama login password for API calls.",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Panorama template for each FireNet gateway.",
			},
			"template_stack": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Panorama template stack for each FireNet gateway.",
			},
			"route_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of firewall virtual router to program. If left unspecified, the Controller programs the Panorama template’s first router.",
			},
			"number_of_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of retries for 'save' or 'synchronize'.",
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Retry interval in seconds for `save` or `synchronize`.",
			},
			"save": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to save or not.",
			},
			"synchronize": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to sync or not.",
			},
		},
	}
}

func dataSourceAviatrixFireNetFirewallManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	firewallManager := &goaviatrix.FirewallManager{
		VpcID:         d.Get("vpc_id").(string),
		GatewayName:   d.Get("gateway_name").(string),
		VendorType:    d.Get("vendor_type").(string),
		PublicIP:      d.Get("public_ip").(string),
		Username:      d.Get("username").(string),
		Password:      d.Get("password").(string),
		Template:      d.Get("template").(string),
		TemplateStack: d.Get("template_stack").(string),
		RouteTable:    d.Get("route_table").(string),
		Save:          d.Get("save").(bool),
		Synchronize:   d.Get("synchronize").(bool),
	}

	if firewallManager.Save && firewallManager.Synchronize {
		return diag.Errorf("can't do 'save' and 'synchronize' at the same time for firewall manager integration")
	}

	if firewallManager.VendorType == "Palo Alto Networks Panorama" && (firewallManager.PublicIP == "" ||
		firewallManager.Username == "" || firewallManager.Password == "" || firewallManager.Template == "" || firewallManager.TemplateStack == "") {
		return diag.Errorf("'public_ip', 'username', 'password', 'template' and 'template_stack' are required for vendor type 'Palo Alto Networks Panorama'")
	}

	numberOfRetries := d.Get("number_of_retries").(int)
	retryInterval := d.Get("retry_interval").(int)

	if firewallManager.Save {
		var err error
		for i := 0; ; i++ {
			err = client.EditFireNetFirewallManagerVendorInfo(ctx, firewallManager)
			if err == nil {
				break
			}
			if i < numberOfRetries {
				time.Sleep(time.Duration(retryInterval) * time.Second)
			} else {
				d.SetId("")
				return diag.Errorf("failed to 'save' FireNet Firewall Manager Vendor Info: %s", err)
			}
		}
	}

	if firewallManager.Synchronize {
		var err error
		for i := 0; ; i++ {
			err = client.SyncFireNetFirewallManagerVendorConfig(ctx, firewallManager)
			if err == nil {
				break
			}
			if i < numberOfRetries {
				time.Sleep(time.Duration(retryInterval) * time.Second)
			} else {
				d.SetId("")
				return diag.Errorf("failed to 'synchronize' FireNet Firewall Manager Vendor Info: %s", err)
			}
		}
	}

	d.SetId(firewallManager.VpcID)
	return nil
}
