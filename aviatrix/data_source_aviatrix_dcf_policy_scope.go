package aviatrix

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// dataSourceAviatrixDCFPolicyScope looks up an existing DCF policy scope by name.
func dataSourceAviatrixDCFPolicyScope() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixDCFPolicyScopeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the DCF policy scope to look up.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the DCF policy scope.",
			},
			"included": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Selectors whose matches are unioned to form the included policy enforcement points.",
				Elem:        dcfPolicyScopeComputedSelectorsSchema(),
			},
			"excluded": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Selectors whose unioned matches are subtracted from the included set.",
				Elem:        dcfPolicyScopeComputedSelectorsSchema(),
			},
		},
	}
}

// dcfPolicyScopeComputedSelectorsSchema builds read-only selector fields for one policy-scope block.
func dcfPolicyScopeComputedSelectorsSchema() *schema.Resource {
	selectorSchema := make(map[string]*schema.Schema, len(dcfPolicyScopeSelectorKeys))
	for _, key := range dcfPolicyScopeSelectorKeys {
		selectorSchema[key] = &schema.Schema{
			Type:        schema.TypeSet,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: fmt.Sprintf("Values for the %s selector.", key),
		}
	}
	return &schema.Resource{Schema: selectorSchema}
}

func dataSourceAviatrixDCFPolicyScopeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	name := getString(d, "name")
	scope, err := mustClient(meta).GetDCFPolicyScopeByName(ctx, name)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return diag.Errorf("DCF policy scope %q not found", name)
		}
		return diag.Errorf("failed to read DCF policy scope %q: %s", name, err)
	}

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
