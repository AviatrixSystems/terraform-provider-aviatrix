package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	log "github.com/sirupsen/logrus"
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
			"ip_filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "Set of CIDRs to filter the app domain.",
			},
			"tag_filter": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of key value pairs to filter the app domain.",
				//RequiredWith: []string{"resource_filter"},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the App Domain.",
			},
			//"resources": {
			//	Type:        schema.TypeList,
			//	Optional:    true,
			//	Elem:        &schema.Schema{Type: schema.TypeString},
			//	Description: "List of resources to apply the tag filters to.",
			//	RequiredWith: []string{"tag_filter"},
			//},
			//"filters": {
			//	Type:     schema.TypeList,
			//	Required: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"type": {
			//				Type:         schema.TypeString,
			//				Required:     true,
			//				ValidateFunc: validation.StringInSlice([]string{"ip", "tag"}, false),
			//				Description:  "Type of filter. Must be one of ip or tag.",
			//			},
			//			"ips": {
			//				Type:     schema.TypeList,
			//
			//			},
			//			"tags": {
			//				Type:        schema.TypeMap,
			//
			//			},
			//			"resources": {
			//
			//			},
			//		},
			//	},
			//},
		},
	}
}

func marshalAppDomainInput(d *schema.ResourceData) *goaviatrix.AppDomain {
	appDomain := &goaviatrix.AppDomain{
		Name: d.Get("name").(string),
	}

	if _, ok := d.GetOk("ip_filter"); ok {
		ipFilter := &goaviatrix.AppDomainIPFilter{
			Type: "ip",
		}

		for _, ip := range d.Get("ip_filter").(*schema.Set).List() {
			ipFilter.Ips = append(ipFilter.Ips, ip.(string))
		}

		appDomain.IpFilter = ipFilter
	}

	if _, ok := d.GetOk("tag_filter"); ok {
		tagFilter := &goaviatrix.AppDomainTagFilter{
			Type: "tag",
			Tags: make(map[string]string),
		}

		for key, val := range d.Get("tag_filter").(map[string]interface{}) {
			tagFilter.Tags[key] = val.(string)
		}

		//for _, resource := range d.Get("resources").([]interface{}) {
		//	tagFilter.Resources = append(tagFilter.Resources, resource.(string))
		//}
		appDomain.TagFilter = tagFilter
	}

	//for _, filterInterface := range filters {
	//	filter := filterInterface.(map[string]interface{})
	//	filterType := filter["type"].(string)
	//	appDomainFilter := &goaviatrix.AppDomainFilter{
	//		Type: filterType,
	//	}
	//	// TODO Check invalid ips or tags
	//	if filterType == "ip" {
	//		for _, ip := range filter["ips"].([]interface{}) {
	//			appDomainFilter.Ips = append(appDomainFilter.Ips, ip.(string))
	//		}
	//	} else if filterType == "tag" {
	//		appDomainFilter.Tags = map[string]string{}
	//		for key, value := range filter["tags"].(map[string]interface{}) {
	//			appDomainFilter.Tags[key] = value.(string)
	//		}
	//
	//		for _, resource := range filter["resources"].([]interface{}) {
	//			appDomainFilter.Resources = append(appDomainFilter.Ips, resource.(string))
	//		}
	//	}
	//	appDomain.Filters = append(appDomain.Filters, appDomainFilter)
	//}

	return appDomain
}

func resourceAviatrixAppDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	appDomain := marshalAppDomainInput(d)

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

	if appDomain.IpFilter != nil {
		if err := d.Set("ip_filter", appDomain.IpFilter.Ips); err != nil {
			log.Errorf("failed to set ip_filter during App Domain read: %s", err)
		}
	}

	if appDomain.TagFilter != nil {
		if err := d.Set("tag_filter", appDomain.TagFilter.Tags); err != nil {
			log.Errorf("failed to set tag_filter during App Domain read: %s", err)
		}
	}

	return nil
}

func resourceAviatrixAppDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("ip_filter", "tag_filter") {
		appDomain := marshalAppDomainInput(d)

		err := client.UpdateAppDomain(ctx, appDomain, uuid)
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
