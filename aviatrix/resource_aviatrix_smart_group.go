package aviatrix

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// dns1123FmtRe is a regular expression that matches dns label names according to rfc 1123.
// K8s resource names must adhere to this.
var dns1123FmtRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

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
			goaviatrix.NameKey: {
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
									goaviatrix.CidrKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.Any(validation.IsCIDR, validation.IsIPAddress),
										Description:  "CIDR block or IP Address this expression matches.",
									},
									goaviatrix.FqdnKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "FQDN address this expression matches.",
									},
									goaviatrix.SiteKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "Edge Site-ID this expression matches.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type of resource this expression matches.",
									},
									goaviatrix.K8sClusterIdKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Kubernetes Cluster ID this expression matches.",
									},
									goaviatrix.K8sPodNameKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(dns1123FmtRe, "must be a valid Kubernetes Pod name"),
										Description:  "Name of the Kubernetes Pod this expression matches.",
									},
									goaviatrix.K8sServiceKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(dns1123FmtRe, "must be a valid Kubernetes Service name"),
										Description:  "Name of the Kubernetes Service this expression matches.",
									},
									goaviatrix.K8sNamespaceKey: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(dns1123FmtRe, "must be a valid Kubernetes Namespace name"),
										Description:  "Name of the Kubernetes Namespace this expression matches.",
									},
									goaviatrix.ResIdKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource ID this expression matches.",
									},
									goaviatrix.AccountIdKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account ID this expression matches.",
									},
									goaviatrix.AccountNameKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account name this expression matches.",
									},
									goaviatrix.NameKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name this expression matches.",
									},
									goaviatrix.RegionKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region this expression matches.",
									},
									goaviatrix.ZoneKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Zone this expression matches.",
									},
									goaviatrix.TagsPrefix: {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Map of tags this expression matches.",
									},
									goaviatrix.ExtArgsPrefix: {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Map of external arguments this expression matches.",
									},
									goaviatrix.S2CKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of remote site.",
									},
									goaviatrix.ExternalKey: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Identifier of remote data source.",
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
		filter := goaviatrix.NewSmartGroupMatchExpression(selectorInfo)
		if goaviatrix.MapContains(selectorInfo, goaviatrix.ExternalKey) {
			// build a map out of the external arguments
			if extArgsMap, ok := selectorInfo[goaviatrix.ExtArgsPrefix]; ok {
				extArgs := make(map[string]string)
				for key, value := range extArgsMap.(map[string]interface{}) {
					extArgs[key] = value.(string)
				}
				filter.ExtArgs = extArgs
			}
		}
		if tagsMap, ok := selectorInfo[goaviatrix.TagsPrefix]; ok {
			tags := make(map[string]string)
			for key, value := range tagsMap.(map[string]interface{}) {
				tags[key] = value.(string)
			}
			filter.Tags = tags
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
		filterMap := goaviatrix.SmartGroupFilterToResource(filter)
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
