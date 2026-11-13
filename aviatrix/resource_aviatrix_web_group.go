package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
										// ConflictsWith: []string{"snifilter"},
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "URL address this expression matches.",
									},
									goaviatrix.ExternalURLKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Name of the external domain feed whose URLs this expression matches.",
									},
									goaviatrix.ExternalSNIKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Name of the external domain feed whose server names this expression matches.",
									},
									goaviatrix.ExtArgsPrefix: {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Map of arguments passed to the external feed. Only valid with external_url or external_sni.",
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
		Name: getString(d, "name"),
	}

	for _, selectorInterface := range getList(d, "selector.0.match_expressions") {
		if selectorInterface == nil {
			return nil, fmt.Errorf("match expressions block cannot be empty")
		}
		filter, err := marshalWebGroupMatchExpression(mustMap(selectorInterface))
		if err != nil {
			return nil, err
		}

		webGroup.Selector.Expressions = append(webGroup.Selector.Expressions, filter)
	}

	return webGroup, nil
}

// webGroupFilterKeys are the keys that each select what a match_expressions block matches.
var webGroupFilterKeys = []string{
	goaviatrix.SniFilterKey,
	goaviatrix.UrlFilterKey,
	goaviatrix.ExternalURLKey,
	goaviatrix.ExternalSNIKey,
}

// marshalWebGroupMatchExpression builds one match expression from a match_expressions block.
func marshalWebGroupMatchExpression(selectorInfo map[string]any) (*goaviatrix.WebGroupMatchExpression, error) {
	filterCount := 0
	for _, key := range webGroupFilterKeys {
		if goaviatrix.MapContains(selectorInfo, key) {
			filterCount++
		}
	}
	if filterCount > 1 {
		return nil, fmt.Errorf("only one of snifilter, urlfilter, external_url and external_sni can be set in the same match_expressions block")
	}

	filter := &goaviatrix.WebGroupMatchExpression{
		SniFilter:   mustString(selectorInfo[goaviatrix.SniFilterKey]),
		UrlFilter:   mustString(selectorInfo[goaviatrix.UrlFilterKey]),
		ExternalURL: mustString(selectorInfo[goaviatrix.ExternalURLKey]),
		ExternalSNI: mustString(selectorInfo[goaviatrix.ExternalSNIKey]),
	}

	extArgsMap, _ := selectorInfo[goaviatrix.ExtArgsPrefix].(map[string]any)
	if len(extArgsMap) == 0 {
		return filter, nil
	}
	if filter.ExternalURL == "" && filter.ExternalSNI == "" {
		return nil, fmt.Errorf("ext_args can only be set together with external_url or external_sni")
	}
	filter.ExtArgs = make(map[string]string)
	for key, value := range extArgsMap {
		// ext_args share the API "all" map with the filter keys, so a colliding key would replace the filter.
		if slices.Contains(webGroupFilterKeys, key) {
			return nil, fmt.Errorf("ext_args cannot contain the key %q", key)
		}
		filter.ExtArgs[key] = mustString(value)
	}
	return filter, nil
}

// flattenWebGroupSelector converts a Web Group's match expressions into the selector block shared by the resource and data source.
func flattenWebGroupSelector(webGroup *goaviatrix.WebGroup) []any {
	var expressions []any
	for _, filter := range webGroup.Selector.Expressions {
		filterMap := map[string]any{
			goaviatrix.SniFilterKey:   filter.SniFilter,
			goaviatrix.UrlFilterKey:   filter.UrlFilter,
			goaviatrix.ExternalURLKey: filter.ExternalURL,
			goaviatrix.ExternalSNIKey: filter.ExternalSNI,
		}
		if len(filter.ExtArgs) > 0 {
			filterMap[goaviatrix.ExtArgsPrefix] = filter.ExtArgs
		}
		expressions = append(expressions, filterMap)
	}

	return []any{
		map[string]any{
			"match_expressions": expressions,
		},
	}
}

func resourceAviatrixWebGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

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

func resourceAviatrixWebGroupReadIfRequired(ctx context.Context, d *schema.ResourceData, meta any, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixWebGroupRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixWebGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()
	mustSet(d, "uuid", uuid)

	webGroup, err := client.GetWebGroup(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Web Group: %s", err)
	}
	mustSet(d, "name", webGroup.Name)

	if err := d.Set("selector", flattenWebGroupSelector(webGroup)); err != nil {
		return diag.Errorf("failed to set selector during Web Group read: %s", err)
	}

	return nil
}

func resourceAviatrixWebGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

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

func resourceAviatrixWebGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()
	err := client.DeleteWebGroup(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete Web Group: %v", err)
	}

	return nil
}
