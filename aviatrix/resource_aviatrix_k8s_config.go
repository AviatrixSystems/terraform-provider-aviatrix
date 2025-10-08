package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixK8sConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixK8sConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixK8sConfigRead,
		UpdateWithoutTimeout: resourceAviatrixK8sConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixK8sConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_k8s": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable K8s.",
			},
			"enable_dcf_policies": {
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

	enableK8s := d.Get("enable_k8s").(bool)
	enableDcfPolicies := d.Get("enable_dcf_policies").(bool)

	// Validate that enable_dcf_policies can only be true if enable_k8s is true
	if enableDcfPolicies && !enableK8s {
		return diag.Errorf("enable_dcf_policies can only be true when enable_k8s is also true")
	}

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
	if enable {
		diags := client.EnableK8s(ctx)
		if diags != nil {
			return diag.Errorf("failed to enable K8s: %s", diags)
		}
	} else {
		diags := client.DisableK8s(ctx)
		if diags != nil {
			return diag.Errorf("failed to disable K8s: %s", diags)
		}
	}
	return nil
}

func setK8sDcfPoliciesFeature(ctx context.Context, client *goaviatrix.Client, enable bool) diag.Diagnostics {
	if enable {
		diags := client.EnableK8sDcfPolicies(ctx)
		if diags != nil {
			return diag.Errorf("failed to enable K8s DCF policies: %s", diags)
		}
	} else {
		diags := client.DisableK8sDcfPolicies(ctx)
		if diags != nil {
			return diag.Errorf("failed to disable K8s DCF policies: %s", diags)
		}
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
		return diag.Errorf("failed to read K8s status: %s", err)
	}
	d.Set("enable_k8s", k8sConfig.EnableK8s)
	d.Set("enable_dcf_policies", k8sConfig.EnableDcfPolicies)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixK8sConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enableK8s := d.Get("enable_k8s").(bool)
	enableDcfPolicies := d.Get("enable_dcf_policies").(bool)

	// Validate that enable_dcf_policies can only be true if enable_k8s is true
	if enableDcfPolicies && !enableK8s {
		return diag.Errorf("enable_dcf_policies can only be true when enable_k8s is also true")
	}

	if d.HasChange("enable_k8s") {
		if err := setK8sFeature(ctx, client, enableK8s); err != nil {
			return err
		}
	}

	if d.HasChange("enable_dcf_policies") {
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
