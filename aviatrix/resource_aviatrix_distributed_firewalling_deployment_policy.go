package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDistributedFirewallingDeploymentPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingDeploymentPolicyCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingDeploymentPolicyRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingDeploymentPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"providers": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "List of CSPs to apply the DCF policies to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"set_defaults": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
				Description: "Set to true to reset the deployment policy to default values. " +
					"Set to false to create a new deployment policy with the specified providers.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingDeploymentPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}
	setDefaults, ok := d.Get("set_defaults").(bool)
	if !ok {
		return diag.Errorf("failed to assert 'set_defaults' as bool")
	}
	providers, ok := d.Get("providers").(*schema.Set)
	if !ok {
		return diag.Errorf("failed to assert 'providers' as array of strings")
	}
	validDiag, isValid := validateDeploymentPolicyInput(providers, setDefaults)
	if !isValid {
		return validDiag
	} else if validDiag != nil {
		diags = append(diags, validDiag...)
	}

	providersList := []string{}
	for _, v := range providers.List() {
		if str, ok := v.(string); ok {
			providersList = append(providersList, str)
		} else {
			return diag.Errorf("failed to convert provider value %v to string", v)
		}
	}

	deploymentPolicy := &goaviatrix.DistributedFirewallingDeploymentPolicy{
		Providers:   providersList,
		SetDefaults: setDefaults,
	}

	if err := client.CreateDistributedFirewallingDeploymentPolicy(ctx, deploymentPolicy); err != nil {
		return diag.Errorf("failed to create Aviatrix Distributed Firewalling Deployment Policy: %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))
	return append(diags, resourceAviatrixDistributedFirewallingDeploymentPolicyRead(ctx, d, meta)...)
}

func validateDeploymentPolicyInput(providers *schema.Set, setDefaults bool) (diag.Diagnostics, bool) {
	if providers.Len() == 0 && !setDefaults {
		return diag.Errorf("providers must be specified if set_defaults is false"), false
	}
	if providers.Len() > 0 && setDefaults {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Providers will be ignored",
			Detail:   "When 'set_defaults' is set to true, providers that are set will be ignored.",
		}}, true
	}
	return nil, true
}

func resourceAviatrixDistributedFirewallingDeploymentPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)

	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.Id() != strings.ReplaceAll(client.ControllerIP, ".", "-") {
		return diag.Errorf("ID: %s does not match controller IP %q: please provide correct ID for importing", d.Id(), client.ControllerIP)
	}

	deploymentPolicy, err := client.GetDistributedFirewallingDeploymentPolicy(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Aviatrix Distributed Firewalling Deployment Policy: %s", err)
	}

	if err := d.Set("providers", deploymentPolicy.Providers); err != nil {
		return diag.Errorf("failed to set 'providers': %v", err)
	}

	if err := d.Set("set_defaults", deploymentPolicy.SetDefaults); err != nil {
		return diag.Errorf("failed to set 'set_defaults': %v", err)
	}

	d.SetId(strings.ReplaceAll(client.ControllerIP, ".", "-"))
	return nil
}

func resourceAviatrixDistributedFirewallingDeploymentPolicyDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	// These dummy values are required but will be ignored by the API when SetDefaults=true
	dummyProviders := []string{
		"GCP",
		"AWS",
	}
	deploymentPolicy := &goaviatrix.DistributedFirewallingDeploymentPolicy{
		Providers:   dummyProviders,
		SetDefaults: true, // Reset to default to delete
	}

	if err := client.CreateDistributedFirewallingDeploymentPolicy(ctx, deploymentPolicy); err != nil {
		return diag.Errorf("failed to delete the current deployment policy and reset to default: %v", err)
	}

	return nil
}
