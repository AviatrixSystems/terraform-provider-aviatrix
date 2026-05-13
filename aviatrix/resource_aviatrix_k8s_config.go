package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	propertyEnableK8s         = "enable_k8s"
	propertyEnableDcfPolicies = "enable_dcf_policies"
)

func resourceAviatrixK8sConfig() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage:   "This resource is deprecated. Use aviatrix_config_feature instead.",
		CreateWithoutTimeout: resourceAviatrixK8sConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixK8sConfigRead,
		UpdateWithoutTimeout: resourceAviatrixK8sConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixK8sConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			enableK8s := d.Get(propertyEnableK8s).(bool)
			enableDcfPolicies := d.Get(propertyEnableDcfPolicies).(bool)

			if enableDcfPolicies && !enableK8s {
				return errors.New("enable_dcf_policies can only be true when enable_k8s is also true")
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			propertyEnableK8s: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable K8s.",
			},
			propertyEnableDcfPolicies: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable DCF policies in K8s clusters. Can only be true if enable_k8s is also true.",
			},
		},
	}
}

func resourceAviatrixK8sConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enableK8s := d.Get(propertyEnableK8s).(bool)
	enableDcfPolicies := d.Get(propertyEnableDcfPolicies).(bool)

	if err := setK8sFeature(ctx, client, enableK8s); err != nil {
		return err
	}

	if err := setK8sDcfPoliciesFeature(ctx, client, enableDcfPolicies); err != nil {
		return err
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixK8sConfigRead(ctx, d, meta)
}

func setK8sFeature(ctx context.Context, client *goaviatrix.Client, enable bool) diag.Diagnostics {
	err := client.ToggleControllerFeature(ctx, goaviatrix.FeatureK8s, enable)
	if err != nil {
		return diag.Errorf("failed to set K8s feature: %s", err)
	}
	return nil
}

func setK8sDcfPoliciesFeature(ctx context.Context, client *goaviatrix.Client, enable bool) diag.Diagnostics {
	err := client.ToggleControllerFeature(ctx, goaviatrix.FeatureK8sDcfPolicies, enable)
	if err != nil {
		return diag.Errorf("failed to set K8s DCF policies feature: %s", err)
	}
	return nil
}

func resourceAviatrixK8sConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	k8sConfig, err := client.GetK8sStatus(ctx)
	if err != nil {
		// The API always exposes these feature flags. GetK8sStatus should not return "not found".
		// If a feature is unknown, the controller reports it as disabled (enabled=false).
		return diag.Errorf("failed to read K8s status: %s", err)
	}
	if err := d.Set(propertyEnableK8s, k8sConfig.EnableK8s); err != nil {
		return diag.Errorf("failed to set enable_k8s on terraform state: %s", err)
	}
	if err := d.Set(propertyEnableDcfPolicies, k8sConfig.EnableDcfPolicies); err != nil {
		return diag.Errorf("failed to set enable_dcf_policies on terraform state: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixK8sConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enableK8s := d.Get(propertyEnableK8s).(bool)
	enableDcfPolicies := d.Get(propertyEnableDcfPolicies).(bool)

	if d.HasChange(propertyEnableK8s) {
		if err := setK8sFeature(ctx, client, enableK8s); err != nil {
			return err
		}
	}

	if d.HasChange(propertyEnableDcfPolicies) {
		if err := setK8sDcfPoliciesFeature(ctx, client, enableDcfPolicies); err != nil {
			return err
		}
	}

	return resourceAviatrixK8sConfigRead(ctx, d, meta)
}

func resourceAviatrixK8sConfigDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if err := setK8sDcfPoliciesFeature(ctx, client, false); err != nil {
		return err
	}

	if err := setK8sFeature(ctx, client, false); err != nil {
		return err
	}

	return nil
}
