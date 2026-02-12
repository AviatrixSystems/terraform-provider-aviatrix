//revive:disable:var-naming
package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDCFIpsProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFIpsProfileCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFIpsProfileRead,
		UpdateWithoutTimeout: resourceAviatrixDCFIpsProfileUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFIpsProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"profile_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the IPS profile.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the IPS profile.",
			},
			"rule_feeds": {
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Description: "Rule feeds configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_feeds_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "List of custom rule feed UUIDs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"external_feeds_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "List of external rule feed IDs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ignored_sids": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "List of rule SIDs to ignore.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"intrusion_actions": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Actions for different severity levels.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"alert", "alert_and_drop"}, false),
				},
				ValidateFunc: validateIntrusionActionsKeys,
			},
		},
	}
}

// validateIntrusionActionsKeys ensures only allowed keys are used in the intrusion_actions map.
func validateIntrusionActionsKeys(val interface{}, _ string) (warns []string, errs []error) {
	allowedKeys := map[string]struct{}{
		"informational": {},
		"minor":         {},
		"major":         {},
		"critical":      {},
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		errs = append(errs, errors.New("intrusion_actions must be a map"))
		return
	}
	for k := range m {
		if _, found := allowedKeys[k]; !found {
			errs = append(errs, errors.New("invalid key for intrusion_actions: '"+k+"'. Allowed keys are: informational, minor, major, critical"))
		}
	}
	return
}

// IPS Profile CRUD operations

func resourceAviatrixDCFIpsProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	profile := &goaviatrix.IpsProfile{
		ProfileName:      getString(d, "profile_name"),
		RuleFeeds:        expandRuleFeeds(getSet(d, "rule_feeds").List()),
		IntrusionActions: expandIntrusionActions(mustMap(d.Get("intrusion_actions"))),
	}

	response, err := client.CreateIpsProfile(ctx, profile)
	if err != nil {
		return diag.Errorf("failed to create IPS profile: %v", err)
	}

	d.SetId(response.UUID)
	return resourceAviatrixDCFIpsProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	profile, err := client.GetIpsProfile(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read IPS profile: %v", err)
	}
	mustSet(d, "uuid", profile.UUID)
	mustSet(d, "profile_name", profile.ProfileName)
	mustSet(d, "rule_feeds", flattenRuleFeeds(profile.RuleFeeds))
	mustSet(d, "intrusion_actions", profile.IntrusionActions)

	return nil
}

func resourceAviatrixDCFIpsProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	profile := &goaviatrix.IpsProfile{
		ProfileName:      getString(d, "profile_name"),
		RuleFeeds:        expandRuleFeeds(getSet(d, "rule_feeds").List()),
		IntrusionActions: expandIntrusionActions(mustMap(d.Get("intrusion_actions"))),
	}

	_, err := client.UpdateIpsProfile(ctx, d.Id(), profile)
	if err != nil {
		return diag.Errorf("failed to update IPS profile: %v", err)
	}

	return resourceAviatrixDCFIpsProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteIpsProfile(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete IPS profile: %v", err)
	}

	return nil
}

// Helper functions

func expandRuleFeeds(ruleFeeds []interface{}) goaviatrix.IpsRuleFeeds {
	if len(ruleFeeds) == 0 {
		return goaviatrix.IpsRuleFeeds{
			CustomFeedsIds:   []string{},
			ExternalFeedsIds: []string{},
			IgnoredSids:      []int{},
		}
	}

	ruleFeedsMap := mustMap(ruleFeeds[0])

	var customFeedsIds []string
	if v, ok := ruleFeedsMap["custom_feeds_ids"]; ok && v != nil {
		customFeedsIds = expandStringList(mustSchemaSet(v).List())
	}

	var externalFeedsIds []string
	if v, ok := ruleFeedsMap["external_feeds_ids"]; ok && v != nil {
		externalFeedsIds = expandStringList(mustSchemaSet(v).List())
	}

	var ignoredSids []int
	if v, ok := ruleFeedsMap["ignored_sids"]; ok && v != nil {
		ignoredSids = expandIntList(mustSchemaSet(v).List())
	}

	return goaviatrix.IpsRuleFeeds{
		CustomFeedsIds:   customFeedsIds,
		ExternalFeedsIds: externalFeedsIds,
		IgnoredSids:      ignoredSids,
	}
}

func flattenRuleFeeds(ruleFeeds goaviatrix.IpsRuleFeeds) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"custom_feeds_ids":   ruleFeeds.CustomFeedsIds,
			"external_feeds_ids": ruleFeeds.ExternalFeedsIds,
			"ignored_sids":       ruleFeeds.IgnoredSids,
		},
	}
}

func expandIntrusionActions(actions map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range actions {
		result[k] = mustString(v)
	}
	return result
}
