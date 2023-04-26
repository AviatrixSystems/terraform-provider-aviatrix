package aviatrix

import (
	"context"
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixWebGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixWebGroupCreate,
		ReadWithoutTimeout:   resourceAviatrixWebGroupRead,
		UpdateWithoutTimeout: resourceAviatrixWebGroupUpdate,
		DeleteWithoutTimeout: resourceAviatrixWebGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Web Group.",
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
									"snifilter": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Server name indicator this expression matches.",
									},
									"urlfilter": {
										Type:     schema.TypeString,
										Optional: true,
										//ConflictsWith: []string{"snifilter"},
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "URL address this expression matches.",
									},
								},
							},
						},
					},
				},
				Description: "List of match expressions for the Web Group.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the Web Group.",
			},
		},
	}
}

func marshalWebGroupInput(d *schema.ResourceData) (*goaviatrix.WebGroup, error) {
	webGroup := &goaviatrix.WebGroup{
		Name: d.Get("name").(string),
	}

	for _, selectorInterface := range d.Get("selector.0.match_expressions").([]interface{}) {
		if selectorInterface == nil {
			return nil, fmt.Errorf("match expressions block cannot be empty")
		}
		selectorInfo := selectorInterface.(map[string]interface{})
		var filter *goaviatrix.WebGroupMatchExpression

		if mapContains(selectorInfo, "snifilter") && mapContains(selectorInfo, "urlfilter") {
			return nil, fmt.Errorf("snifilter and urlfilter can't be set at the same time under the same match_expressions")
		}

		filter = &goaviatrix.WebGroupMatchExpression{
			SniFilter: selectorInfo["snifilter"].(string),
			UrlFilter: selectorInfo["urlfilter"].(string),
		}

		webGroup.Selector.Expressions = append(webGroup.Selector.Expressions, filter)
	}

	return webGroup, nil
}

func resourceAviatrixWebGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	webGroup, err := marshalWebGroupInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for Web Group during create: %s", err)
	}

	flag := false
	defer resourceAviatrixWebGroupReadIfRequired(ctx, d, meta, &flag)

	uuid, err := client.CreateWebGroup(ctx, webGroup)
	if err != nil {
		return diag.Errorf("failed to create Web Group: %s", err)
	}
	d.SetId(uuid)
	return resourceAviatrixWebGroupReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixWebGroupReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixWebGroupRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixWebGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Set("uuid", uuid)

	webGroup, err := client.GetWebGroup(ctx, uuid)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Web Group: %s", err)
	}

	d.Set("name", webGroup.Name)

	var expressions []interface{}

	for _, filter := range webGroup.Selector.Expressions {
		filterMap := map[string]interface{}{
			"snifilter": filter.SniFilter,
			"urlfilter": filter.UrlFilter,
		}

		expressions = append(expressions, filterMap)
	}

	selector := []interface{}{
		map[string]interface{}{
			"match_expressions": expressions,
		},
	}
	if err := d.Set("selector", selector); err != nil {
		return diag.Errorf("failed to set selector during Web Group read: %s", err)
	}

	return nil
}

func resourceAviatrixWebGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("name", "selector") {
		webGroup, err := marshalWebGroupInput(d)
		if err != nil {
			return diag.Errorf("invalid inputs for Web Group during update: %s", err)
		}

		err = client.UpdateWebGroup(ctx, webGroup, uuid)
		if err != nil {
			return diag.Errorf("failed to update Web Group selector: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixWebGroupRead(ctx, d, meta)
}

func resourceAviatrixWebGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	err := client.DeleteWebGroup(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete Web Group: %v", err)
	}

	return nil
}
