package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixConfigFeature() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixConfigFeatureCreate,
		ReadWithoutTimeout:   resourceAviatrixConfigFeatureRead,
		DeleteWithoutTimeout: resourceAviatrixConfigFeatureDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"feature_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the feature to enable or disable.",
				ValidateFunc: validation.StringInSlice([]string{
					"microseg",
					"cost_iq",
					"cai",
					"ipv6",
					"nfq_enforce_tls",
					"dcf_on_s2c",
					"dcf_on_psf",
					"dcf_stats_obs_sink",
					"dcf_logs_obs_sink",
					"k8s",
					"sre_metrics_export",
					"k8s_dcf_policies",
					"dcf_on_firenet",
					"primary_gateway_deletion",
					"enable_k8s",
					"enable_dcf_policies",
				}, true),
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Status of the feature (true for enabled, false for disabled).",
			},
		},
	}
}

func resourceAviatrixConfigFeatureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	featureName := getString(d, "feature_name")
	isEnabled := getBool(d, "is_enabled")
	if isEnabled {
		err := client.EnableFeature(ctx, featureName)
		if err != nil {
			return diag.Errorf("failed to enable feature %s: %s", featureName, err)
		}
	} else {
		err := client.DisableFeature(ctx, featureName)
		if err != nil {
			return diag.Errorf("failed to disable feature %s: %s", featureName, err)
		}
	}

	d.SetId(featureName)
	return resourceAviatrixConfigFeatureRead(ctx, d, meta)
}

func resourceAviatrixConfigFeatureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	featureName := strings.ToLower(getString(d, "feature_name"))
	if featureName == "" {
		featureName = strings.ToLower(d.Id())
		if featureName == "" {
			return diag.Errorf("missing feature_name and ID")
		}
		mustSet(d, "feature_name", featureName)
	}

	featureStatus, err := client.GetFeatureStatus(ctx, featureName)
	if err != nil {
		return diag.Errorf("failed to read feature %s status: %s", featureName, err)
	}
	mustSet(d, "is_enabled", featureStatus.Enabled)
	return nil
}

func resourceAviatrixConfigFeatureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	featureName := getString(d, "feature_name")
	err := client.DisableFeature(ctx, featureName)
	if err != nil {
		return diag.Errorf("failed to delete feature %s: %s", featureName, err)
	}
	return nil
}
