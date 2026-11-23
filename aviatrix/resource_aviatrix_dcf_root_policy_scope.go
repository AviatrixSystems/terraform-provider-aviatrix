package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	dcfRootPolicyScopeID          = "dcf_root_policy_scope"
	dcfSystemRootPolicyBlockUUID  = "defa11a1-3000-0000-0000-000000000000"
	dcfDefaultGlobalPolicyScopeID = "def000ad-8000-0000-0000-000000000001"
)

// resourceAviatrixDCFRootPolicyScope defines the system root policy scope configuration.
func resourceAviatrixDCFRootPolicyScope() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFRootPolicyScopeCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFRootPolicyScopeRead,
		UpdateWithoutTimeout: resourceAviatrixDCFRootPolicyScopeUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFRootPolicyScopeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"scope_uuids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsUUID},
				Description: "Set of DCF policy scope UUIDs whose effective enforcement points are unioned to form the root scope level. Descendant scope levels are intersected with this level.",
			},
		},
	}
}

// resourceAviatrixDCFRootPolicyScopeCreate applies the configured scopes to the system root policy block.
func resourceAviatrixDCFRootPolicyScopeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	if _, err := client.PatchDCFPolicyScopeConfig(ctx, dcfSystemRootPolicyBlockUUID, getStringSet(d, "scope_uuids")); err != nil {
		return diag.Errorf("failed to set DCF root policy scopes: %s", err)
	}
	d.SetId(dcfRootPolicyScopeID)
	return resourceAviatrixDCFRootPolicyScopeRead(ctx, d, meta)
}

// resourceAviatrixDCFRootPolicyScopeRead refreshes the scopes applied to the system root policy block.
func resourceAviatrixDCFRootPolicyScopeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.Id() != dcfRootPolicyScopeID {
		return diag.Errorf("ID %q does not match expected import ID %q", d.Id(), dcfRootPolicyScopeID)
	}

	policyBlock, err := mustClient(meta).GetDCFPolicyBlock(ctx, dcfSystemRootPolicyBlockUUID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF root policy scopes: %s", err)
	}
	if policyBlock.ScopeConfig == nil || len(policyBlock.ScopeConfig.ScopeUUIDs) == 0 {
		return diag.Errorf("system root policy block response contains no policy scopes")
	}

	mustSet(d, "scope_uuids", policyBlock.ScopeConfig.ScopeUUIDs)
	d.SetId(dcfRootPolicyScopeID)
	return nil
}

// resourceAviatrixDCFRootPolicyScopeUpdate replaces the scopes applied to the system root policy block.
func resourceAviatrixDCFRootPolicyScopeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if _, err := mustClient(meta).PatchDCFPolicyScopeConfig(ctx, dcfSystemRootPolicyBlockUUID, getStringSet(d, "scope_uuids")); err != nil {
		return diag.Errorf("failed to update DCF root policy scopes: %s", err)
	}
	return resourceAviatrixDCFRootPolicyScopeRead(ctx, d, meta)
}

// resourceAviatrixDCFRootPolicyScopeDelete restores the system root policy block to the Global Scope.
func resourceAviatrixDCFRootPolicyScopeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if _, err := mustClient(meta).PatchDCFPolicyScopeConfig(
		ctx,
		dcfSystemRootPolicyBlockUUID,
		[]string{dcfDefaultGlobalPolicyScopeID},
	); err != nil {
		return diag.Errorf("failed to reset DCF root policy scopes: %s", err)
	}
	d.SetId("")
	return nil
}
