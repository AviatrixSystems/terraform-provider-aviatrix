package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bgp_med_to_sdn_metric_global": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "BGP MED to SDN metric global configuration",
			},
		},
	}
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	bgpMedToSdnMetric, ok := d.Get("bgp_med_to_sdn_metric_global").(bool)
	if !ok {
		return diag.Errorf("failed to assert bgp_med_to_sdn_metric_global as bool")
	}
	if bgpMedToSdnMetric {
		err := client.EnableControllerBgpMedToSdnMetricGlobal(ctx)
		if err != nil {
			return diag.Errorf("failed to enable controller BGP MED to SDN metric global config: %v", err)
		}
	} else {
		err := client.DisableControllerBgpMedToSdnMetricGlobal(ctx)
		if err != nil {
			return diag.Errorf("failed to disable controller BGP MED to SDN metric global config: %v", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	commGlobal, err := client.GetControllerBgpMedToSdnMetricGlobal(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get controller BGP MED to SDN metric global config: %v", err)
	}

	err = d.Set("bgp_med_to_sdn_metric_global", commGlobal)
	if err != nil {
		return diag.Errorf("failed to set bgp_med_to_sdn_metric_global: %v", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.HasChange("bgp_med_to_sdn_metric_global") {
		bgpMedToSdnMetric, ok := d.Get("bgp_med_to_sdn_metric_global").(bool)
		if !ok {
			return diag.Errorf("failed to assert bgp_med_to_sdn_metric_global as bool")
		}
		if bgpMedToSdnMetric {
			err := client.EnableControllerBgpMedToSdnMetricGlobal(ctx)
			if err != nil {
				return diag.Errorf("failed to enable controller BGP MED to SDN metric global config: %v", err)
			}
		} else {
			err := client.DisableControllerBgpMedToSdnMetricGlobal(ctx)
			if err != nil {
				return diag.Errorf("failed to disable controller BGP MED to SDN metric global config: %v", err)
			}
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	err := client.DisableControllerBgpMedToSdnMetricGlobal(ctx)
	if err != nil {
		return diag.Errorf("failed to delete controller BGP MED to SDN metric global config: %v", err)
	}

	return nil
}
