package aviatrix

import (
	"context"
	"errors"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixDCFIpsRuleFeed() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFIpsRuleFeedCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFIpsRuleFeedRead,
		UpdateWithoutTimeout: resourceAviatrixDCFIpsRuleFeedUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFIpsRuleFeedDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"feed_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for the rule feed.",
			},
			"file_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPS rule feed file content containing Suricata rules.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the IPS rule feed.",
			},
			"content_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 hash of the file content.",
			},
			"ips_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of IPS rules extracted from the file.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

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
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Rule feeds configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_feeds_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of custom rule feed UUIDs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"external_feeds_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of external rule feed IDs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ignored_sids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of rule SIDs to ignore.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"never_drop_sids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of rule SIDs to never drop.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"intrusion_actions": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Actions for different severity levels.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"alert", "alert_and_drop"}, false),
				},
			},
		},
	}
}

// IPS Rule Feed CRUD operations

func resourceAviatrixDCFIpsRuleFeedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	ruleFeed := &goaviatrix.IpsRuleFeed{
		FeedName:    d.Get("feed_name").(string),
		FileContent: d.Get("file_content").(string),
	}

	response, err := client.CreateIpsRuleFeed(ctx, ruleFeed)
	if err != nil {
		return diag.Errorf("failed to create IPS rule feed: %v", err)
	}

	d.SetId(response.UUID)
	return resourceAviatrixDCFIpsRuleFeedRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsRuleFeedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	ruleFeed, err := client.GetIpsRuleFeed(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read IPS rule feed: %v", err)
	}

	d.Set("uuid", ruleFeed.UUID)
	d.Set("feed_name", ruleFeed.FeedName)
	d.Set("content_hash", ruleFeed.ContentHash)
	d.Set("ips_rules", ruleFeed.IpsRules)

	return nil
}

func resourceAviatrixDCFIpsRuleFeedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	ruleFeed := &goaviatrix.IpsRuleFeed{
		FeedName:    d.Get("feed_name").(string),
		FileContent: d.Get("file_content").(string),
	}

	_, err := client.UpdateIpsRuleFeed(ctx, d.Id(), ruleFeed)
	if err != nil {
		return diag.Errorf("failed to update IPS rule feed: %v", err)
	}

	return resourceAviatrixDCFIpsRuleFeedRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsRuleFeedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteIpsRuleFeed(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete IPS rule feed: %v", err)
	}

	return nil
}

// IPS Profile CRUD operations

func resourceAviatrixDCFIpsProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	profile := &goaviatrix.IpsProfile{
		ProfileName:      d.Get("profile_name").(string),
		RuleFeeds:        expandRuleFeeds(d.Get("rule_feeds").([]interface{})),
		IntrusionActions: expandIntrusionActions(d.Get("intrusion_actions").(map[string]interface{})),
	}

	response, err := client.CreateIpsProfile(ctx, profile)
	if err != nil {
		return diag.Errorf("failed to create IPS profile: %v", err)
	}

	d.SetId(response.UUID)
	return resourceAviatrixDCFIpsProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	profile, err := client.GetIpsProfile(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read IPS profile: %v", err)
	}

	d.Set("uuid", profile.UUID)
	d.Set("profile_name", profile.ProfileName)
	d.Set("rule_feeds", flattenRuleFeeds(profile.RuleFeeds))
	d.Set("intrusion_actions", profile.IntrusionActions)

	return nil
}

func resourceAviatrixDCFIpsProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	profile := &goaviatrix.IpsProfile{
		ProfileName:      d.Get("profile_name").(string),
		RuleFeeds:        expandRuleFeeds(d.Get("rule_feeds").([]interface{})),
		IntrusionActions: expandIntrusionActions(d.Get("intrusion_actions").(map[string]interface{})),
	}

	_, err := client.UpdateIpsProfile(ctx, d.Id(), profile)
	if err != nil {
		return diag.Errorf("failed to update IPS profile: %v", err)
	}

	return resourceAviatrixDCFIpsProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

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
			NeverDropSids:    []int{},
		}
	}

	ruleFeedsMap := ruleFeeds[0].(map[string]interface{})

	return goaviatrix.IpsRuleFeeds{
		CustomFeedsIds:   expandStringList(ruleFeedsMap["custom_feeds_ids"].([]interface{})),
		ExternalFeedsIds: expandStringList(ruleFeedsMap["external_feeds_ids"].([]interface{})),
		IgnoredSids:      expandIntList(ruleFeedsMap["ignored_sids"].([]interface{})),
		NeverDropSids:    expandIntList(ruleFeedsMap["never_drop_sids"].([]interface{})),
	}
}

func flattenRuleFeeds(ruleFeeds goaviatrix.IpsRuleFeeds) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"custom_feeds_ids":   ruleFeeds.CustomFeedsIds,
			"external_feeds_ids": ruleFeeds.ExternalFeedsIds,
			"ignored_sids":       ruleFeeds.IgnoredSids,
			"never_drop_sids":    ruleFeeds.NeverDropSids,
		},
	}
}

func expandIntrusionActions(actions map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range actions {
		result[k] = v.(string)
	}
	return result
}
