package aviatrix

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// parseFirewallTemplateConfig converts terraform set data to FirewallTemplateConfig map
func parseFirewallTemplateConfig(rawSet *schema.Set) map[string]goaviatrix.FirewallTemplateConfig {
	if rawSet == nil || rawSet.Len() == 0 {
		return nil
	}

	result := make(map[string]goaviatrix.FirewallTemplateConfig)

	for _, item := range rawSet.List() {
		itemMap := mustMap(item)

		firewallID := strings.TrimSpace(mustString(itemMap["firewall_id"]))
		config := goaviatrix.FirewallTemplateConfig{}
		if template, exists := itemMap["template"]; exists {
			config.Template = strings.TrimSpace(mustString(template))
		}

		if templateStack, exists := itemMap["template_stack"]; exists {
			config.TemplateStack = strings.TrimSpace(mustString(templateStack))
		}

		if routeTable, exists := itemMap["route_table"]; exists {
			config.RouteTable = strings.TrimSpace(mustString(routeTable))
		}

		result[firewallID] = config

	}

	return result
}

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
			"config_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     goaviatrix.FirewallManagerConfigModeDefault,
				Description: "Config mode type, DEFAULT or ADVANCE",
			},
			"firewall_template_config": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The firewall template configuration.",
				Set: func(v interface{}) int {
					m := mustMap(v)
					idRaw, ok := m["firewall_id"]
					if !ok {
						panic("internal error: firewall_template_config element missing firewall_id")
					}
					id := strings.TrimSpace(mustString(idRaw))
					return schema.HashString(id)
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"firewall_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Firewall instance ID.",
						},
						"template": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Panorama template for each firewall.",
						},
						"template_stack": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Panorama template stack for each firewall.",
						},
						"route_table": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of firewall virtual router to program. If left unspecified, the Controller programs the Panorama template's first router.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixFireNetFirewallManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	firewallManager := &goaviatrix.FirewallManager{
		VpcID:       getString(d, "vpc_id"),
		GatewayName: getString(d, "gateway_name"),
		VendorType:  getString(d, "vendor_type"),
		PublicIP:    getString(d, "public_ip"),
		Username:    getString(d, "username"),
		Password:    getString(d, "password"),
		Save:        getBool(d, "save"),
		Synchronize: getBool(d, "synchronize"),
		ConfigMode:  getString(d, "config_mode"),
	}
	if firewallManager.ConfigMode == goaviatrix.FirewallManagerConfigModeAdvance {
		firewallManager.FirewallTemplateConfig = parseFirewallTemplateConfig(getSet(d, "firewall_template_config"))
	} else {
		firewallManager.Template = getString(d, "template")
		firewallManager.TemplateStack = getString(d, "template_stack")
		firewallManager.RouteTable = getString(d, "route_table")
	}

	if firewallManager.Save && firewallManager.Synchronize {
		return diag.Errorf("can't do 'save' and 'synchronize' at the same time for firewall manager integration")
	}

	if firewallManager.VendorType == "Palo Alto Networks Panorama" {
		if firewallManager.PublicIP == "" ||
			firewallManager.Username == "" || firewallManager.Password == "" {
			return diag.Errorf("'public_ip', 'username', and 'password' are required for vendor type 'Palo Alto Networks Panorama'")
		}
		if firewallManager.ConfigMode == goaviatrix.FirewallManagerConfigModeAdvance {
			if len(firewallManager.FirewallTemplateConfig) == 0 {
				return diag.Errorf("'firewall_template_config' must be specified in ADVANCE config mode")
			}
			for fwID, config := range firewallManager.FirewallTemplateConfig {
				if config.Template == "" || config.TemplateStack == "" {
					return diag.Errorf("firewall_template_config[%q]: both 'template' and 'template_stack' are required in ADVANCE config mode", fwID)
				}
			}
		} else {
			if firewallManager.Template == "" || firewallManager.TemplateStack == "" {
				return diag.Errorf("'template' and 'template_stack' must be specified in DEFAULT config mode")
			}
		}
	}

	numberOfRetries := getInt(d, "number_of_retries")
	retryInterval := getInt(d, "retry_interval")

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
