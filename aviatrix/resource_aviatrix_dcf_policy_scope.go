package aviatrix

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

var dcfPolicyScopeSelectorKeys = []string{
	"csps",
	"vpc_ids",
	"csp_account_ids",
	"region_names",
	"gateway_group_names",
	"gateway_group_ids",
	"preset",
}

var dcfPolicyScopeCSPValues = []string{
	"AWS",
	"AWS-GOV",
	"AWS-CHINA",
	"AZURE",
	"AZURE-GOV",
	"AZURE-CHINA",
	"GCP",
	"OCI",
	"ALICLOUD",
}

// resourceAviatrixDCFPolicyScope defines the Terraform schema and lifecycle operations.
func resourceAviatrixDCFPolicyScope() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFPolicyScopeCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFPolicyScopeRead,
		UpdateWithoutTimeout: resourceAviatrixDCFPolicyScopeUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFPolicyScopeDelete,
		CustomizeDiff:        resourceAviatrixDCFPolicyScopeCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
				Description:  "Name of the DCF policy scope.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the DCF policy scope.",
			},
			"included": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Selectors whose matches are unioned to form the included policy enforcement points.",
				Elem:        dcfPolicyScopeSelectorsSchema(),
			},
			"excluded": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Selectors whose unioned matches are subtracted from the included set. When omitted, the controller uses preset = [\"none\"].",
				Elem:        dcfPolicyScopeSelectorsSchema(),
			},
		},
	}
}

// dcfPolicyScopeSelectorsSchema builds the selector fields for one policy-scope block.
func dcfPolicyScopeSelectorsSchema() *schema.Resource {
	selectorSchema := make(map[string]*schema.Schema, len(dcfPolicyScopeSelectorKeys))
	for _, key := range dcfPolicyScopeSelectorKeys {
		valueSchema := &schema.Schema{
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringIsNotEmpty},
			Description: fmt.Sprintf("Values for the %s selector. Matches from all selectors in this block are unioned.", key),
		}
		switch key {
		case "preset":
			valueSchema.MaxItems = 1
			valueSchema.Elem = &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"all", "none"}, false),
			}
		case "csps":
			valueSchema.Elem = &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(dcfPolicyScopeCSPValues, false),
			}
		}
		selectorSchema[key] = valueSchema
	}
	return &schema.Resource{Schema: selectorSchema}
}

// resourceAviatrixDCFPolicyScopeCustomizeDiff rejects configured selector blocks with no populated selector.
func resourceAviatrixDCFPolicyScopeCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	// SDK schema relationships do not run when every optional field in a nested block is absent.
	if d.NewValueKnown("included") {
		if err := validateDCFPolicyScopeSelectorConfig("included", d.Get("included")); err != nil {
			return err
		}
	}
	if d.NewValueKnown("excluded") {
		if excluded, ok := d.GetOk("excluded"); ok {
			if err := validateDCFPolicyScopeSelectorConfig("excluded", excluded); err != nil {
				return err
			}
		}
	}
	// excluded is Optional+Computed, so removing it from configuration would otherwise
	// keep the prior state value and never send the update.
	if dcfPolicyScopeExcludedOmitted(d.GetRawConfig()) {
		if err := d.SetNew("excluded", dcfPolicyScopeDefaultExcluded()); err != nil {
			return fmt.Errorf("failed to plan default excluded selectors: %w", err)
		}
	}
	return nil
}

// dcfPolicyScopeDefaultExcluded returns the controller default for an omitted excluded block.
func dcfPolicyScopeDefaultExcluded() []any {
	return []any{map[string]any{"preset": []any{"none"}}}
}

// dcfPolicyScopeExcludedOmitted reports whether excluded is known to be absent from configuration.
func dcfPolicyScopeExcludedOmitted(rawConfig cty.Value) bool {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return false
	}
	excluded := rawConfig.GetAttr("excluded")
	return excluded.IsKnown() && (excluded.IsNull() || excluded.LengthInt() == 0)
}

// validateDCFPolicyScopeSelectorConfig verifies that exactly one nonempty selector block is present.
func validateDCFPolicyScopeSelectorConfig(name string, value any) error {
	blocks, ok := value.([]any)
	if !ok || len(blocks) != 1 {
		return fmt.Errorf("%s must contain exactly one selector block", name)
	}
	block, ok := blocks[0].(map[string]any)
	if !ok {
		return fmt.Errorf("%s selector block has an invalid type", name)
	}
	for _, key := range dcfPolicyScopeSelectorKeys {
		if values, ok := block[key].(*schema.Set); ok && values.Len() > 0 {
			return nil
		}
	}
	return fmt.Errorf("%s must configure at least one selector", name)
}

// marshalDCFPolicyScopeInput converts Terraform state into the API request model.
func marshalDCFPolicyScopeInput(d *schema.ResourceData) *goaviatrix.DCFPolicyScopeInput {
	input := &goaviatrix.DCFPolicyScopeInput{
		Name:     getString(d, "name"),
		Included: expandDCFPolicyScopeSelectors(d.Get("included")),
	}
	if dcfPolicyScopeExcludedConfigured(d) {
		excluded := expandDCFPolicyScopeSelectors(d.Get("excluded"))
		input.Excluded = &excluded
	}
	return input
}

// dcfPolicyScopeExcludedConfigured reports whether excluded was present in configuration.
func dcfPolicyScopeExcludedConfigured(d *schema.ResourceData) bool {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		_, configured := d.GetOk("excluded")
		return configured
	}
	excluded := rawConfig.GetAttr("excluded")
	return excluded.IsKnown() && !excluded.IsNull() && excluded.LengthInt() > 0
}

// expandDCFPolicyScopeSelectors converts a Terraform selector block into API selectors.
func expandDCFPolicyScopeSelectors(value any) goaviatrix.DCFPolicyScopeSelectors {
	blocks := mustSlice(value)
	// blocks[0] is a map[string]any like {"csps": {"aws", "azure"}}, we only allow 1 with MaxItems: 1 so blocks[0] is the only one.
	block := mustMap(blocks[0])
	selectors := make([]goaviatrix.DCFPolicyScopeSelector, 0, len(dcfPolicyScopeSelectorKeys))
	for _, key := range dcfPolicyScopeSelectorKeys {
		values, ok := block[key].(*schema.Set)
		if !ok || values.Len() == 0 {
			continue
		}
		selectors = append(selectors, goaviatrix.DCFPolicyScopeSelector{
			key: expandStringList(values.List()),
		})
	}
	return goaviatrix.DCFPolicyScopeSelectors{Selectors: selectors}
}

// flattenDCFPolicyScopeSelectors converts API selectors into one stable Terraform block.
func flattenDCFPolicyScopeSelectors(selectors goaviatrix.DCFPolicyScopeSelectors) []any {
	if len(selectors.Selectors) == 0 {
		return []any{}
	}
	block := make(map[string]any)
	for _, selector := range selectors.Selectors {
		for key, values := range selector {
			// The API may return repeated keys; merge them into Terraform's canonical set field.
			existing, _ := block[key].([]string)
			block[key] = append(existing, values...)
		}
	}
	return []any{block}
}

// resourceAviatrixDCFPolicyScopeCreate creates a policy scope and refreshes its state.
func resourceAviatrixDCFPolicyScopeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	scope, err := mustClient(meta).CreateDCFPolicyScope(ctx, marshalDCFPolicyScopeInput(d))
	if err != nil {
		return diag.Errorf("failed to create DCF policy scope: %s", err)
	}
	d.SetId(scope.UUID)
	return resourceAviatrixDCFPolicyScopeRead(ctx, d, meta)
}

// resourceAviatrixDCFPolicyScopeRead refreshes a policy scope or removes missing state.
func resourceAviatrixDCFPolicyScopeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	scope, err := mustClient(meta).GetDCFPolicyScope(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF policy scope: %s", err)
	}

	mustSet(d, "name", scope.Name)
	mustSet(d, "uuid", scope.UUID)
	mustSet(d, "included", flattenDCFPolicyScopeSelectors(scope.Included))
	if scope.Excluded == nil {
		mustSet(d, "excluded", []any{})
	} else {
		mustSet(d, "excluded", flattenDCFPolicyScopeSelectors(*scope.Excluded))
	}
	d.SetId(scope.UUID)
	return nil
}

// resourceAviatrixDCFPolicyScopeUpdate replaces a policy scope and refreshes its state.
func resourceAviatrixDCFPolicyScopeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := mustClient(meta).UpdateDCFPolicyScope(ctx, d.Id(), marshalDCFPolicyScopeInput(d)); err != nil {
		return diag.Errorf("failed to update DCF policy scope: %s", err)
	}
	return resourceAviatrixDCFPolicyScopeRead(ctx, d, meta)
}

// resourceAviatrixDCFPolicyScopeDelete deletes a policy scope and tolerates an absent object.
func resourceAviatrixDCFPolicyScopeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := mustClient(meta).DeleteDCFPolicyScope(ctx, d.Id())
	if errors.Is(err, goaviatrix.ErrNotFound) {
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to delete DCF policy scope: %s", err)
	}
	return nil
}
