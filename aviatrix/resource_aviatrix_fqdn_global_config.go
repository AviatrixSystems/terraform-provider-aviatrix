package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFQDNGlobalConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixFQDNGlobalConfigsCreate,
		ReadWithoutTimeout:   resourceAviatrixFQDNGlobalConfigsRead,
		UpdateWithoutTimeout: resourceAviatrixFQDNGlobalConfigsUpdate,
		DeleteWithoutTimeout: resourceAviatrixFQDNGlobalConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_exception_rule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If enabled, it allows packets passing through the gateway without an SNI field. Only applies to whitelist.",
			},
			"enable_private_network_filtering": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If enabled, destination FQDN names that translate to private IP address range (RFC 1918) " +
					"are subject to FQDN whitelist filtering function.",
			},
			"enable_custom_network_filtering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, it customizes packet destination address ranges not to be filtered by FQDN.",
			},
			"configured_ips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Source IPs configured for a specific tag. Can be subnet CIDRs or host IP addressesConfig IP addresses.",
			},
			"enable_caching": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If enabled, it caches the resolved IP address from FQDN filter.",
			},
			"enable_exact_match": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If enabled, the resolved IP address from FQDN filter is cached so that " +
					"if subsequent TCP session matches the cached IP address list, " +
					"FQDN domain name is not checked and the session is allowed to pass.",
			},
			"rfc_1918": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Unconfigured Private IP address range.",
			},
		},
	}
}

func resourceAviatrixFQDNGlobalConfigsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enablePrivateNetworkFiltering := d.Get("enable_private_network_filtering").(bool)
	enableCustomNetworkFiltering := d.Get("enable_custom_network_filtering").(bool)

	if enablePrivateNetworkFiltering && enableCustomNetworkFiltering {
		return diag.Errorf("enable_private_network_filtering and enable_custom_network_filtering can't be set true at the same time")
	}
	if enableCustomNetworkFiltering {
		if _, ok := d.GetOk("configured_ips"); !ok {
			return diag.Errorf("configured_ips is required to enable custom network filtering")
		}
	} else {
		if _, ok := d.GetOk("configured_ips"); ok {
			return diag.Errorf("configured_ips is required to be empty to disable custom network filtering")
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixFQDNGlobalConfigsReadIfRequired(ctx, d, meta, &flag)

	if d.Get("enable_exception_rule").(bool) {
		err := client.EnableFQDNExceptionRule(ctx)
		if err != nil {
			return diag.Errorf("failed to enable fqdn exception rule: %s", err)
		}
	} else {
		err := client.DisableFQDNExceptionRule(ctx)
		if err != nil {
			return diag.Errorf("failed to disable fqdn exception rule: %s", err)
		}
	}

	if d.Get("enable_private_network_filtering").(bool) {
		err := client.EnableFQDNPrivateNetworks(ctx)
		if err != nil {
			return diag.Errorf("failed to enable private network filtering: %s", err)
		}
	} else {
		err := client.DisableFQDNPrivateNetwork(ctx)
		if err != nil {
			return diag.Errorf("failed to disable private network filtering: %s", err)
		}
	}

	if enableCustomNetworkFiltering {
		configIpString := strings.Join(goaviatrix.ExpandStringList(d.Get("configured_ips").([]interface{})), ",")
		err := client.SetFQDNCustomNetwork(ctx, configIpString)
		if err != nil {
			return diag.Errorf("failed to customize network filtering: %s", err)
		}
	}

	if d.Get("enable_caching").(bool) {
		err := client.EnableFQDNCache(ctx)
		if err != nil {
			return diag.Errorf("failed to enable fqdn cache: %s", err)
		}
	} else {
		err := client.DisableFQDNCache(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable fqdn cache: %s", err)
		}
	}

	if d.Get("enable_exact_match").(bool) {
		err := client.EnableFQDNExactMatch(ctx)
		if err != nil {
			return diag.Errorf("failed to enable fqdn exact match: %s", err)
		}
	} else {
		err := client.DisableFQDNExactMatch(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable fqdn exact match: %s", err)
		}
	}

	return resourceAviatrixFQDNGlobalConfigsReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixFQDNGlobalConfigsReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFQDNGlobalConfigsRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixFQDNGlobalConfigsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	exceptionRuleStatus, err := client.GetFQDNExceptionRuleStatus(ctx)
	if err != nil {
		return diag.Errorf("filed to get FQDN exception rule status: %s", err)
	}
	if *exceptionRuleStatus == "enabled" {
		d.Set("enable_exception_rule", true)
	} else if *exceptionRuleStatus == "disabled" {
		d.Set("enable_exception_rule", false)
	}

	privateSubFilter, err := client.GetFQDNPrivateNetworkFilteringStatus(ctx)
	if err != nil {
		return diag.Errorf("failed to get FQDN network filter status: %s", err)
	}
	if privateSubFilter.PrivateSubFilter == "enabled" {
		d.Set("enable_private_network_filtering", true)
	} else if privateSubFilter.PrivateSubFilter == "disabled" {
		d.Set("enable_private_network_filtering", false)
	} else if privateSubFilter.PrivateSubFilter == "custom" {
		d.Set("enable_custom_network_filtering", true)
		d.Set("configured_ips", privateSubFilter.ConfiguredIps)
	} else {
		d.Set("enable_custom_network_filtering", false)
	}
	d.Set("rfc_1918", privateSubFilter.Rfc1918)

	cacheGlobalStatus, err := client.GetFQDNCacheGlobalStatus(ctx)
	if err != nil {
		return diag.Errorf("failed to get FQDN cache global status: %s", err)
	}
	if *cacheGlobalStatus == "enabled" {
		d.Set("enable_caching", true)
	} else if *cacheGlobalStatus == "disabled" {
		d.Set("enable_caching", false)
	}

	exactMatchStatus, err := client.GetFQDNExactMatchStatus(ctx)
	if err != nil {
		return diag.Errorf("failed to get FQDN exact match status: %s", err)
	}
	if *exactMatchStatus == "enabled" {
		d.Set("enable_exact_match", true)
	} else if *exactMatchStatus == "disabled" {
		d.Set("enable_exact_match", false)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixFQDNGlobalConfigsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("enable_exception_rule") {
		if d.Get("enable_exception_rule").(bool) {
			err := client.EnableFQDNExceptionRule(ctx)
			if err != nil {
				return diag.Errorf("failed to disable FQDN exception rule in update: %s", err)
			}
		} else {
			err := client.DisableFQDNExceptionRule(ctx)
			if err != nil {
				return diag.Errorf("failed to disable FQDN exception rule in update: %s", err)
			}
		}
	}

	if d.HasChanges("enable_private_network_filtering", "enable_custom_network_filtering", "configured_ips") {
		enablePrivateNetworkFiltering := d.Get("enable_private_network_filtering").(bool)
		enableCustomNetworkFiltering := d.Get("enable_custom_network_filtering").(bool)

		if enablePrivateNetworkFiltering && enableCustomNetworkFiltering {
			return diag.Errorf("enable_private_network_filtering and enable_custom_network_filtering can't be set true at the same time")
		}
		if enableCustomNetworkFiltering {
			if _, ok := d.GetOk("configured_ips"); !ok {
				return diag.Errorf("configured_ips is required to enable custom network filtering")
			}
		} else {
			if _, ok := d.GetOk("configured_ips"); ok {
				return diag.Errorf("configured_ips is required to be empty to disable custom network filtering")
			}
		}

		if d.HasChange("enable_private_network_filtering") && !enablePrivateNetworkFiltering {
			err := client.DisableFQDNPrivateNetwork(ctx)
			if err != nil {
				return diag.Errorf("failed to disable private network filtering in update: %s", err)
			}
		}

		if d.HasChanges("enable_custom_network_filtering", "configured_ips") {
			if enableCustomNetworkFiltering {
				configIpString := strings.Join(goaviatrix.ExpandStringList(d.Get("configured_ips").([]interface{})), ",")
				err := client.SetFQDNCustomNetwork(ctx, configIpString)
				if err != nil {
					return diag.Errorf("failed to customize network filtering in update: %s", err)
				}
			} else {
				if !enablePrivateNetworkFiltering {
					err := client.DisableFQDNPrivateNetwork(ctx)
					if err != nil {
						return diag.Errorf("failed to disable private network filtering in update: %s", err)
					}
				}
			}
		}

		if d.HasChange("enable_private_network_filtering") && enablePrivateNetworkFiltering {
			err := client.EnableFQDNPrivateNetworks(ctx)
			if err != nil {
				return diag.Errorf("failed to enable private network filtering in update: %s", err)
			}

		}
	}

	if d.HasChange("enable_caching") {
		if d.Get("enable_caching").(bool) {
			err := client.EnableFQDNCache(ctx)
			if err != nil {
				return diag.Errorf("failed to enable FQDN cache in update: %s", err)
			}
		} else {
			err := client.DisableFQDNCache(ctx)
			if err != nil {
				return diag.Errorf("failed to disable FQDN cache in update: %s", err)
			}
		}
	}

	if d.HasChange("enable_exact_match") {
		if d.Get("enable_exact_match").(bool) {
			err := client.EnableFQDNExactMatch(ctx)
			if err != nil {
				return diag.Errorf("failed to enable FQDN exact match in update: %s", err)
			}
		} else {
			err := client.DisableFQDNExactMatch(ctx)
			if err != nil {
				return diag.Errorf("failed to enable FQDN exact match in update: %s", err)
			}

		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixFQDNGlobalConfigsRead(ctx, d, meta)
}

func resourceAviatrixFQDNGlobalConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client) //default enabled

	err := client.EnableFQDNExceptionRule(ctx)
	if err != nil {
		return diag.Errorf("failed to set fqdn exception rule to default in destroy: %s", err)
	}

	err = client.DisableFQDNPrivateNetwork(ctx)
	if err != nil {
		return diag.Errorf("failed to set private network filtering to default in destroy: %s", err)
	}

	err = client.EnableFQDNCache(ctx)
	if err != nil {
		return diag.Errorf("failed to set fqdn cache to default in destroy: %s", err)
	}

	err = client.DisableFQDNExactMatch(ctx)
	if err != nil {
		return diag.Errorf("failed to set fqdn exact match to default in destroy: %s", err)
	}

	return nil
}
