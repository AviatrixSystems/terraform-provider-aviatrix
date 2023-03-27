package aviatrix

import (
	"context"
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixSmartGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixSmartGroupCreate,
		ReadWithoutTimeout:   resourceAviatrixSmartGroupRead,
		UpdateWithoutTimeout: resourceAviatrixSmartGroupUpdate,
		DeleteWithoutTimeout: resourceAviatrixSmartGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Smart Group.",
			},
			"selector": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_expressions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPAddress),
										Description:  "CIDR block or IP Address this expression matches.",
									},
									"fqdn": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "FQDN address this expression matches.",
									},
									"site": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Edge Site-ID this expression matches.",
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"vm", "vpc", "subnet"}, false),
										Description:  "Type of resource this expression matches.",
									},
									"res_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource ID this expression matches.",
									},
									"account_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account ID this expression matches.",
									},
									"account_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account name this expression matches.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name this expression matches.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region this expression matches.",
									},
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Zone this expression matches.",
									},
									"tags": {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Map of tags this expression matches.",
									},
								},
							},
						},
					},
				},
				Description: "List of match expressions for the Smart Group.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the Smart Group.",
			},
		},
	}
}

func marshalSmartGroupInput(d *schema.ResourceData) (*goaviatrix.SmartGroup, error) {
	smartGroup := &goaviatrix.SmartGroup{
		Name: d.Get("name").(string),
	}

	for _, selectorInterface := range d.Get("selector.0.match_expressions").([]interface{}) {
		if selectorInterface == nil {
			return nil, fmt.Errorf("match expressions block cannot be empty")
		}
		selectorInfo := selectorInterface.(map[string]interface{})
		var filter *goaviatrix.SmartGroupMatchExpression

		if mapContains(selectorInfo, "cidr") || mapContains(selectorInfo, "fqdn") || mapContains(selectorInfo, "site") {
			for _, key := range []string{"type", "res_id", "account_id", "account_name", "name", "region", "zone", "tags"} {
				if mapContains(selectorInfo, key) {
					return nil, fmt.Errorf("%q must be empty when %q is set", key, "cidr")
				}
			}

			filter = &goaviatrix.SmartGroupMatchExpression{
				CIDR: selectorInfo["cidr"].(string),
				FQDN: selectorInfo["fqdn"].(string),
				Site: selectorInfo["site"].(string),
			}
		} else {
			if !mapContains(selectorInfo, "type") {
				return nil, fmt.Errorf("%q is required when %q, %q and %q are all empty", "type", "cidr", "fqdn", "site")
			}
			filter = &goaviatrix.SmartGroupMatchExpression{
				Type:        selectorInfo["type"].(string),
				ResId:       selectorInfo["res_id"].(string),
				AccountId:   selectorInfo["account_id"].(string),
				AccountName: selectorInfo["account_name"].(string),
				Name:        selectorInfo["name"].(string),
				Region:      selectorInfo["region"].(string),
				Zone:        selectorInfo["zone"].(string),
			}

			if _, ok := selectorInfo["tags"]; ok {
				tags := make(map[string]string)
				for key, value := range selectorInfo["tags"].(map[string]interface{}) {
					tags[key] = value.(string)
				}
				filter.Tags = tags
			}
		}

		smartGroup.Selector.Expressions = append(smartGroup.Selector.Expressions, filter)
	}

	return smartGroup, nil
}

func resourceAviatrixSmartGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	smartGroup, err := marshalSmartGroupInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for Smart Group during create: %s", err)
	}

	flag := false
	defer resourceAviatrixSmartGroupReadIfRequired(ctx, d, meta, &flag)

	uuid, err := client.CreateSmartGroup(ctx, smartGroup)
	if err != nil {
		return diag.Errorf("failed to create Smart Group: %s", err)
	}
	d.SetId(uuid)
	return resourceAviatrixSmartGroupReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixSmartGroupReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSmartGroupRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixSmartGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Set("uuid", uuid)

	smartGroup, err := client.GetSmartGroup(ctx, uuid)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Smart Group: %s", err)
	}

	d.Set("name", smartGroup.Name)

	var expressions []interface{}

	for _, filter := range smartGroup.Selector.Expressions {
		filterMap := map[string]interface{}{
			"type":         filter.Type,
			"cidr":         filter.CIDR,
			"fqdn":         filter.FQDN,
			"site":         filter.Site,
			"res_id":       filter.ResId,
			"account_id":   filter.AccountId,
			"account_name": filter.AccountName,
			"name":         filter.Name,
			"region":       filter.Region,
			"zone":         filter.Zone,
			"tags":         filter.Tags,
		}

		expressions = append(expressions, filterMap)
	}

	selector := []interface{}{
		map[string]interface{}{
			"match_expressions": expressions,
		},
	}
	if err := d.Set("selector", selector); err != nil {
		return diag.Errorf("failed to set selector during Smart Group read: %s", err)
	}

	return nil
}

func resourceAviatrixSmartGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("name", "selector") {
		smartGroup, err := marshalSmartGroupInput(d)
		if err != nil {
			return diag.Errorf("invalid inputs for Smart Group during update: %s", err)
		}

		err = client.UpdateSmartGroup(ctx, smartGroup, uuid)
		if err != nil {
			return diag.Errorf("failed to update Smart Group selector: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixSmartGroupRead(ctx, d, meta)
}

func resourceAviatrixSmartGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	err := client.DeleteSmartGroup(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete Smart Group: %v", err)
	}

	return nil
}
