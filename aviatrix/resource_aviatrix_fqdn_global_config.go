package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFQDNGlobalConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixFQDNGlobalConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixFQDNGlobalConfigRead,
		UpdateWithoutTimeout: resourceAviatrixFQDNGlobalConfigUpdate,
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
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "Customized packet destination address ranges not to be filtered by FQDN. " +
					"Can be selected from pre-defined RFC 1918 range, or own network range.",
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
		},
	}
}

func resourceAviatrixFQDNGlobalConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	defer resourceAviatrixFQDNGlobalConfigReadIfRequired(ctx, d, meta, &flag)

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
		var configIPs []string
		for _, v := range d.Get("configured_ips").([]interface{}) {
			configIPs = append(configIPs, v.(string))
		}
		err := client.SetFQDNCustomNetwork(ctx, configIPs)
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

	return resourceAviatrixFQDNGlobalConfigReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixFQDNGlobalConfigReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFQDNGlobalConfigRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixFQDNGlobalConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceAviatrixFQDNGlobalConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)

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

		if d.HasChanges("enable_custom_network_filtering", "configured_ips") {
			if enableCustomNetworkFiltering {
				var configIPs []string
				for _, v := range d.Get("configured_ips").([]interface{}) {
					configIPs = append(configIPs, v.(string))
				}
				err := client.SetFQDNCustomNetwork(ctx, configIPs)
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

		if d.HasChange("enable_private_network_filtering") {
			if enablePrivateNetworkFiltering {
				err := client.EnableFQDNPrivateNetworks(ctx)
				if err != nil {
					return diag.Errorf("failed to enable private network filtering in update: %s", err)
				}
			} else {
				if !enableCustomNetworkFiltering {
					err := client.DisableFQDNPrivateNetwork(ctx)
					if err != nil {
						return diag.Errorf("failed to disable private network filtering in update: %s", err)
					}
				}
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

	d.Partial(false)
	return resourceAviatrixFQDNGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixFQDNGlobalConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client) // default enabled
	if !ok {
		return diag.Errorf("meta is not a valid Client pointer")
	}

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
