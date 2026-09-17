package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func normalizeResourceIDName(name string) string {
	return strings.ReplaceAll(name, ".", "-")
}

// Details on the BGP MED To SDN Metric feature and related APIs can be found here:
// Design doc: https://docs.google.com/document/d/1h-gxgwZ6OxNuLNLFgKymldqdoRU3kY1zEshbZs13C0A/edit?usp=sharing
// APIs: https://aviatrix.atlassian.net/wiki/spaces/AVXENG/pages/3136946177/BGP+MED+to+SDN+Metric+APIs

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
				Description: "BGP Multi-Exit Discriminator to SDN metric global configuration",
			},
		},
	}
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	bgpMedToSdnMetric := getBool(d, "bgp_med_to_sdn_metric_global")
	if bgpMedToSdnMetric {
		err := client.EnableControllerBgpMedToSdnMetricGlobal(ctx)
		if err != nil {
			return diag.Errorf("failed to enable controller BGP Multi-Exit Discriminator to SDN metric global config: %v", err)
		}
	} else {
		err := client.DisableControllerBgpMedToSdnMetricGlobal(ctx)
		if err != nil {
			return diag.Errorf("failed to disable controller BGP Multi-Exit Discriminator to SDN metric global config: %v", err)
		}
	}

	d.SetId(normalizeResourceIDName(client.ControllerIP))
	return resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if d.Id() != normalizeResourceIDName(client.ControllerIP) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	commGlobal, err := client.GetControllerBgpMedToSdnMetricGlobal(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get controller BGP Multi-Exit Discriminator to SDN metric global config: %v", err)
	}

	mustSet(d, "bgp_med_to_sdn_metric_global", commGlobal)
	d.SetId(normalizeResourceIDName(client.ControllerIP))
	return nil
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if d.HasChange("bgp_med_to_sdn_metric_global") {
		bgpMedToSdnMetric := getBool(d, "bgp_med_to_sdn_metric_global")
		if bgpMedToSdnMetric {
			err := client.EnableControllerBgpMedToSdnMetricGlobal(ctx)
			if err != nil {
				return diag.Errorf("failed to enable controller BGP Multi-Exit Discriminator to SDN metric global config: %v", err)
			}
		} else {
			err := client.DisableControllerBgpMedToSdnMetricGlobal(ctx)
			if err != nil {
				return diag.Errorf("failed to disable controller BGP Multi-Exit Discriminator to SDN metric global config: %v", err)
			}
		}
	}

	d.SetId(normalizeResourceIDName(client.ControllerIP))
	return resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpMedToSdnMetricGlobalConfigDelete(ctx context.Context, _ *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DisableControllerBgpMedToSdnMetricGlobal(ctx)
	if err != nil {
		return diag.Errorf("failed to delete controller BGP Multi-Exit Discriminator to SDN metric global config: %v", err)
	}

	return nil
}
