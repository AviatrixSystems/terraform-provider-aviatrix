package aviatrix

import (
	"context"
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixAppDomain() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixAppDomainCreate,
		ReadWithoutTimeout:   resourceAviatrixAppDomainRead,
		UpdateWithoutTimeout: resourceAviatrixAppDomainUpdate,
		DeleteWithoutTimeout: resourceAviatrixAppDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the App Domain.",
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
										ValidateFunc: validation.IsCIDR,
										Description:  "CIDR block this expression matches.",
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
				Description: "List of match expressions for the App Domain.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the App Domain.",
			},
		},
	}
}

func marshalAppDomainInput(d *schema.ResourceData) (*goaviatrix.AppDomain, error) {
	appDomain := &goaviatrix.AppDomain{
		Name: d.Get("name").(string),
	}

	for _, selectorInterface := range d.Get("selector.0.match_expressions").([]interface{}) {
		if selectorInterface == nil {
			return nil, fmt.Errorf("match expressions block cannot be empty")
		}
		selectorInfo := selectorInterface.(map[string]interface{})
		filter := &goaviatrix.AppDomainMatchExpression{
			CIDR:        selectorInfo["cidr"].(string),
			Type:        selectorInfo["type"].(string),
			ResId:       selectorInfo["res_id"].(string),
			AccountId:   selectorInfo["account_id"].(string),
			AccountName: selectorInfo["account_name"].(string),
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

		appDomain.Selector.Expressions = append(appDomain.Selector.Expressions, filter)
	}

	return appDomain, nil
}

func resourceAviatrixAppDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	appDomain, err := marshalAppDomainInput(d)
	if err != nil {
		return diag.Errorf("failed to marshal inputs for App Domain during create: %s", err)
	}

	flag := false
	defer resourceAviatrixAppDomainReadIfRequired(ctx, d, meta, &flag)

	uuid, err := client.CreateAppDomain(ctx, appDomain)
	if err != nil {
		return diag.Errorf("failed to create App Domain: %s", err)
	}
	d.SetId(uuid)
	return resourceAviatrixAppDomainReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixAppDomainReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAppDomainRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixAppDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Set("uuid", uuid)

	appDomain, err := client.GetAppDomain(ctx, uuid)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read App Domain: %s", err)
	}

	d.Set("name", appDomain.Name)

	var expressions []interface{}

	for _, filter := range appDomain.Selector.Expressions {
		filterMap := map[string]interface{}{
			"type":         filter.Type,
			"cidr":         filter.CIDR,
			"res_id":       filter.ResId,
			"account_id":   filter.AccountId,
			"account_name": filter.AccountName,
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
		return diag.Errorf("failed to set selector during App Domain read: %s", err)
	}

	return nil
}

func resourceAviatrixAppDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("selector") {
		appDomain, err := marshalAppDomainInput(d)
		if err != nil {
			return diag.Errorf("failed to marshal inputs for App Domain during update: %s", err)
		}

		err = client.UpdateAppDomain(ctx, appDomain, uuid)
		if err != nil {
			return diag.Errorf("failed to update App Domain filters: %s", err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceAviatrixAppDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	err := client.DeleteAppDomain(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete App Domain: %v", err)
	}

	return nil
}
