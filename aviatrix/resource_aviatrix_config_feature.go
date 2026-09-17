package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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

func resourceAviatrixConfigFeatureCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	featureName := getString(d, "feature_name")
	err := goaviatrix.ValidateFeatureName(ctx, client, featureName)
	if err != nil {
		return diag.Errorf("failed to validate feature name %s: %s", featureName, err)
	}

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

func resourceAviatrixConfigFeatureRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceAviatrixConfigFeatureDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	featureName := getString(d, "feature_name")
	err := client.DisableFeature(ctx, featureName)
	if err != nil {
		return diag.Errorf("failed to delete feature %s: %s", featureName, err)
	}
	return nil
}
