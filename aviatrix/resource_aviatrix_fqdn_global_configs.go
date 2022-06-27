package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFQDNGlobalConfigs() *schema.Resource {

	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixFQDNGlobalConfigsCreate,
		ReadWithoutTimeout:   resourceAviatrixFQDNGlobalConfigsRead,
		UpdateWithoutTimeout: resourceAviatrixFQDNGlobalConfigsUpdate,
		DeleteWithoutTimeout: resourceAviatrixFQDNGlobalConfigDelete,

		Schema: map[string]*schema.Schema{
			"exception_rule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow packets passing through the gateway without an SNI field.",
			},
			"network_filtering": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Packet destination address ranges be filtered by FQDN.",
			},
			"configured_ips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Config IP address.",
			},
			"rfc_1918": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Private IP address range.",
			},
			"caching": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Cached the resolved IP address from FQDN filter.",
			},
			"exact_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Exact match in FQDN filter.",
			},
		},
	}
}

func resourceAviatrixFQDNGlobalConfigsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	exceptionRuleStatus, err := client.GetFQDNExceptionRuleStatus(ctx)
	if err != nil {
		return diag.Errorf("Filed to get FQDN exception rule status in create func: %s", err)
	}
	if d.Get("exception_rule").(bool) && *exceptionRuleStatus == "disable" {
		err := client.EnableFQDNExceptionRule(ctx)
		if err != nil {
			return diag.Errorf("Failed to enable fqdn exception rule in create func: %s", err)
		}
	} else if !d.Get("exception_rule").(bool) && *exceptionRuleStatus == "enabled" {
		err := client.DisableFQDNExceptionRule(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable fqdn exception rule in create func: %s", err)
		}
	}

	if d.Get("network_filtering").(string) == "Enable Private Network Filtering" {
		err := client.EnableFQDNPrivateNetworks(ctx)
		if err != nil {
			return diag.Errorf("Failed to enable private network filtering in create func: %s", err)
		}
	} else if d.Get("network_filtering").(string) == "Disable Private Network Filtering" {
		err := client.DisableFQDNPrivateNetwork(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable private network filtering in create func: %s", err)
		}
	} else if d.Get("network_filtering").(string) == "Customize Network Filtering" {
		if _, ok := d.GetOk("configured_ips"); ok {
			configIpString := strings.Join(goaviatrix.ExpandStringList(d.Get("configured_ips").([]interface{})), ",")
			err := client.SetFQDNCustomNetwork(ctx, configIpString)
			if err != nil {
				return diag.Errorf("Failed to customize network filtering in create func: %s", err)
			}
		}
	}

	cacheGlobalStatus, err := client.GetFQDNCacheGlobalStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN cache global status in create func: %s", err)
	}
	if d.Get("caching").(bool) && *cacheGlobalStatus == "disabled" {
		err := client.EnableFQDNCache(ctx)
		if err != nil {
			return diag.Errorf("Failed to enable fqdn cache in create func: %s", err)
		}
	} else if !d.Get("caching").(bool) && *cacheGlobalStatus == "enabled" {
		err := client.DisableFQDNCache(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable fqdn cache in create func: %s", err)
		}
	}

	exactMatchStatus, err := client.GetFQDNExactMatchStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN exact match status in create func: %s", err)
	}
	if d.Get("exact_match").(bool) && *exactMatchStatus == "disabled" {
		err := client.EnableFQDNExactMatch(ctx)
		if err != nil {
			return diag.Errorf("Failed to enable fqdn exact match in create func: %s", err)
		}
	} else if !d.Get("exact_match").(bool) && *exactMatchStatus == "enabled" {
		err := client.DisableFQDNExactMatch(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable fqdn exact match in create func: %s", err)
		}
	}

	flag := false
	defer resourceAviatrixFQDNGlobalConfigsReadIfRequired(ctx, d, meta, &flag)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
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

	exceptionRuleStatus, err := client.GetFQDNExceptionRuleStatus(ctx)
	if err != nil {
		return diag.Errorf("Filed to get FQDN exception rule status: %s", err)
	}
	if *exceptionRuleStatus == "enabled" {
		d.Set("exception_rule", true)
	} else if *exceptionRuleStatus == "disabled" {
		d.Set("exception_rule", false)
	}

	privateSubFilter, err := client.GetFQDNPrivateNetworkFilteringStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN private network filter status: %s", err)
	}
	if privateSubFilter.PrivateSubFilter == "enabled" {
		d.Set("network_filtering", "Enable Private Network Filtering")
	} else if privateSubFilter.PrivateSubFilter == "disabled" {
		d.Set("network_filtering", "Disable Private Network Filtering")
	} else if privateSubFilter.PrivateSubFilter == "custom" {
		d.Set("network_filtering", "Customize Network Filtering")
	}
	d.Set("configured_ips", privateSubFilter.ConfiguredIps)
	d.Set("rfc_1918", privateSubFilter.Rfc1918)

	cacheGlobalStatus, err := client.GetFQDNCacheGlobalStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN cache global status: %s", err)
	}
	if *cacheGlobalStatus == "enabled" {
		d.Set("caching", true)
	} else if *cacheGlobalStatus == "disabled" {
		d.Set("caching", false)
	}

	exactMatchStatus, err := client.GetFQDNExactMatchStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN exact match status: %s", err)
	}
	if *exactMatchStatus == "enabled" {
		d.Set("exact_match", true)
	} else if *exactMatchStatus == "disabled" {
		d.Set("exact_match", false)
	}

	return nil
}

func resourceAviatrixFQDNGlobalConfigsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	exceptionRuleStatus, err := client.GetFQDNExceptionRuleStatus(ctx)
	if err != nil {
		return diag.Errorf("Filed to get FQDN exception rule status in update func: %s", err)
	}
	if d.HasChanges("exception_rule") {
		n := d.Get("exception_rule")
		if n == false && *exceptionRuleStatus == "enabled" {
			err := client.DisableFQDNExceptionRule(ctx)
			if err != nil {
				return diag.Errorf("Failed to disable FQDN exception rule in update func: %s", err)
			}
		} else if n == true && *exceptionRuleStatus == "disabled" {
			err := client.EnableFQDNExceptionRule(ctx)
			if err != nil {
				return diag.Errorf("Failed to enable FQDN exception rule in update func: %s", err)
			}
		}
	}

	if d.HasChange("network_filtering") {
		n := d.Get("network_filtering")
		if n == "Disable Private Network Filtering" {
			err := client.DisableFQDNPrivateNetwork(ctx)
			if err != nil {
				return diag.Errorf("Failed to disable FQDN private network in update func: %s", err)
			}
		} else if n == "Enable Private Network Filtering" {
			err := client.EnableFQDNPrivateNetworks(ctx)
			if err != nil {
				return diag.Errorf("Failed to enable FQDN private network in update func: %s", err)
			}
		} else if n == "Customize Network Filtering" {
			if _, ok := d.GetOk("configured_ips"); ok {
				configIpString := strings.Join(goaviatrix.ExpandStringList(d.Get("configured_ips").([]interface{})), ",")
				err := client.SetFQDNCustomNetwork(ctx, configIpString)
				if err != nil {
					return diag.Errorf("Failed to set FQDN custom network in update func: %s", err)
				}
			}
		}
	}

	o, n := d.GetChange("network_filtering")
	if d.HasChange("configured_ips") && o == "Customize Network Filtering" && n == "Customize Network Filtering" {
		if _, ok := d.GetOk("configured_ips"); ok {
			configIpString := strings.Join(goaviatrix.ExpandStringList(d.Get("configured_ips").([]interface{})), ",")
			err := client.SetFQDNCustomNetwork(ctx, configIpString)
			if err != nil {
				return diag.Errorf("Failed to update the configured ips: %s", err)
			}
		}
	}

	cacheGlobalStatus, err := client.GetFQDNCacheGlobalStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN cache global status in update func: %s", err)
	}
	if d.HasChange("caching") {
		n := d.Get("caching")
		if n == false && *cacheGlobalStatus == "enabled" {
			err := client.DisableFQDNCache(ctx)
			if err != nil {
				return diag.Errorf("Failed to Disable FQDN cache in update func: %s", err)
			}
		} else if n == true && *cacheGlobalStatus == "disabled" {
			err := client.EnableFQDNCache(ctx)
			if err != nil {
				return diag.Errorf("Failed to enable FQDN cache in update func: %s", err)
			}
		}
	}

	exactMatchStatus, err := client.GetFQDNExactMatchStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN exact match status in update func: %s", err)
	}
	if d.HasChange("exact_match") {
		if _, ok := d.GetOk("exact_match"); ok {
			n := d.Get("exact_match")
			if n == false && *exactMatchStatus == "enabled" {
				err := client.DisableFQDNExactMatch(ctx)
				if err != nil {
					return diag.Errorf("Failed to disable FQDN exact match in update func: %s", err)
				}
			} else if n == true && *exactMatchStatus == "disabled" {
				err := client.EnableFQDNExactMatch(ctx)
				if err != nil {
					return diag.Errorf("Failed to enable FQDN exact match in update func: %s", err)
				}
			}
		}
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixFQDNGlobalConfigsRead(ctx, d, meta)
}

func resourceAviatrixFQDNGlobalConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	//default enabled
	exceptionRuleStatus, err := client.GetFQDNExceptionRuleStatus(ctx)
	if err != nil {
		return diag.Errorf("Filed to get FQDN exception rule status in delete func: %s", err)
	}
	if *exceptionRuleStatus == "disabled" {
		err := client.EnableFQDNExceptionRule(ctx)
		if err != nil {
			return diag.Errorf("Failed to enable fqdn exception rule in destroy: %s", err)
		}
	}

	//default not enabled Private Network Filtering
	privateSubFilter, err := client.GetFQDNPrivateNetworkFilteringStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN private network filter status in delete func: %s", err)
	}
	if privateSubFilter.PrivateSubFilter != "disabled" {
		err := client.DisableFQDNPrivateNetwork(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable private network filtering in destroy: %s", err)
		}
	}

	//default enabled
	cacheGlobalStatus, err := client.GetFQDNCacheGlobalStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN cache global status in delete func: %s", err)
	}
	if *cacheGlobalStatus == "disabled" {
		err := client.EnableFQDNCache(ctx)
		if err != nil {
			return diag.Errorf("Failed to enable fqdn cache in destroy: %s", err)
		}
	}

	//default not enabled
	exactMatchStatus, err := client.GetFQDNExactMatchStatus(ctx)
	if err != nil {
		return diag.Errorf("Failed to get FQDN exact match status in delete func: %s", err)
	}
	if *exactMatchStatus == "enabled" {
		err := client.DisableFQDNExactMatch(ctx)
		if err != nil {
			return diag.Errorf("Failed to disable fqdn exact match in destroy: %s", err)
		}
	}
	return nil
}
