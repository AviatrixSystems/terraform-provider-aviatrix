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

	if enableK8s {
		err := client.EnableK8s(ctx)
		if err != nil {
			return diag.Errorf("failed to enable K8s: %s", err)
		}
	} else {
		err := client.DisableK8s(ctx)
		if err != nil {
			return diag.Errorf("failed to disable K8s: %s", err)
		}
	}

	if enableDcfPolicies {
		err := client.EnableK8sDcfPolicies(ctx)
		if err != nil {
			return diag.Errorf("failed to enable K8s DCF policies: %s", err)
		}
	} else {
		err := client.DisableK8sDcfPolicies(ctx)
		if err != nil {
			return diag.Errorf("failed to disable K8s DCF policies: %s", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixK8sConfigRead(ctx, d, meta)
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
		if enableK8s {
			err := client.EnableK8s(ctx)
			if err != nil {
				return diag.Errorf("failed to enable K8s during update: %s", err)
			}
		} else {
			err := client.DisableK8s(ctx)
			if err != nil {
				return diag.Errorf("failed to disable K8s during update: %s", err)
			}
		}
	}

	if d.HasChange("enable_dcf_policies") {
		if enableDcfPolicies {
			err := client.EnableK8sDcfPolicies(ctx)
			if err != nil {
				return diag.Errorf("failed to enable K8s DCF policies during update: %s", err)
			}
		} else {
			err := client.DisableK8sDcfPolicies(ctx)
			if err != nil {
				return diag.Errorf("failed to disable K8s DCF policies during update: %s", err)
			}
		}
	}

	return resourceAviatrixK8sConfigRead(ctx, d, meta)
}

func resourceAviatrixK8sConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableK8sDcfPolicies(ctx)
	if err != nil {
		return diag.Errorf("failed to disable K8s DCF policies during delete: %s", err)
	}

	err = client.DisableK8s(ctx)
	if err != nil {
		return diag.Errorf("failed to delete K8s config: %s", err)
	}

	return nil
}
