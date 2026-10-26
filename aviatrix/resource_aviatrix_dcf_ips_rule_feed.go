//revive:disable:var-naming
package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
				Sensitive:   true,
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

// IPS Rule Feed CRUD operations

func resourceAviatrixDCFIpsRuleFeedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	ruleFeed := &goaviatrix.IpsRuleFeed{
		FeedName:    getString(d, "feed_name"),
		FileContent: getString(d, "file_content"),
	}

	response, err := client.CreateIpsRuleFeed(ctx, ruleFeed)
	if err != nil {
		return diag.Errorf("failed to create IPS rule feed: %v", err)
	}

	d.SetId(response.UUID)
	return resourceAviatrixDCFIpsRuleFeedRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsRuleFeedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	ruleFeed, err := client.GetIpsRuleFeed(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read IPS rule feed: %v", err)
	}
	mustSet(d, "uuid", ruleFeed.UUID)
	mustSet(d, "feed_name", ruleFeed.FeedName)
	mustSet(d, "content_hash", ruleFeed.ContentHash)
	mustSet(d, "ips_rules", ruleFeed.IpsRules)

	return nil
}

func resourceAviatrixDCFIpsRuleFeedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChanges("feed_name", "file_content") {
		return resourceAviatrixDCFIpsRuleFeedRead(ctx, d, meta)
	}

	client := mustClient(meta)

	ruleFeed := &goaviatrix.IpsRuleFeed{
		FeedName:    getString(d, "feed_name"),
		FileContent: getString(d, "file_content"),
	}

	_, err := client.UpdateIpsRuleFeed(ctx, d.Id(), ruleFeed)
	if err != nil {
		return diag.Errorf("failed to update IPS rule feed: %v", err)
	}

	return resourceAviatrixDCFIpsRuleFeedRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsRuleFeedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteIpsRuleFeed(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete IPS rule feed: %v", err)
	}

	return nil
}
