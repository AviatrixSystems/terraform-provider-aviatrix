package aviatrix

import (
	"context"
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
			"filters": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ip", "tag"}, false),
							Description:  "Type of filter. Must be one of ip or tag.",
						},
						"ips": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
							Description: "List of CIDRs to filter the app domain.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Map of key value pairs to filter the app domain.",
						},
						"resources": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of resources to apply the tag filters to.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixAppDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	filters := d.Get("filters").([]interface{})
	appDomain, err := formatAppDomainFilters(filters)
	if err != nil {
		return diag.Errorf("failed to format filters when creating App Domain: %s", err)
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
	appDomain, err := client.GetAppDomain(ctx, uuid)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read App Domain: %s", err)
	}

	var appDomainFilters []map[string]interface{}

	for _, appDomainFilter := range appDomain.Filters {
		filter := make(map[string]interface{})
		filter["type"] = appDomainFilter.Type
		if appDomainFilter.Type == "ip" {
			filter["ips"] = appDomainFilter.Ips
		} else if appDomainFilter.Type == "tag" {
			filter["tags"] = appDomainFilter.Tags
			filter["resources"] = appDomainFilter.Resources
		}

		appDomainFilters = append(appDomainFilters, filter)
	}

	if err := d.Set("filters", appDomainFilters); err != nil {
		return diag.Errorf("failed to set filters for App Domain during read: %s", err)
	}

	return nil
}

func resourceAviatrixAppDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChange("filters") {
		filters := d.Get("filters").([]interface{})
		appDomain, err := formatAppDomainFilters(filters)
		if err != nil {
			return diag.Errorf("failed to format filters when updating App Domain: %s", err)
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

func formatAppDomainFilters(filters []interface{}) (*goaviatrix.AppDomain, error) {
	appDomain := &goaviatrix.AppDomain{}

	for _, filterInterface := range filters {
		filter := filterInterface.(map[string]interface{})
		filterType := filter["type"].(string)
		appDomainFilter := &goaviatrix.AppDomainFilter{
			Type: filterType,
		}
		// TODO Check invalid ips or tags
		if filterType == "ip" {
			for _, ip := range filter["ips"].([]interface{}) {
				appDomainFilter.Ips = append(appDomainFilter.Ips, ip.(string))
			}
		} else if filterType == "tag" {
			appDomainFilter.Tags = map[string]string{}
			for key, value := range filter["tags"].(map[string]interface{}) {
				appDomainFilter.Tags[key] = value.(string)
			}

			for _, resource := range filter["resources"].([]interface{}) {
				appDomainFilter.Resources = append(appDomainFilter.Ips, resource.(string))
			}
		}
		appDomain.Filters = append(appDomain.Filters, appDomainFilter)
	}

	return appDomain, nil
}
