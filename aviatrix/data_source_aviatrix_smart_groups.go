package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixSmartGroups() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixSmartGroupsRead,

		Schema: map[string]*schema.Schema{
			"smart_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Smart Groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Smart Group.",
						},
						"selector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_expressions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												goaviatrix.CidrKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "CIDR block or IP Address this expression matches.",
												},
												goaviatrix.FqdnKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "FQDN address this expression matches.",
												},
												goaviatrix.SiteKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Edge Site-ID this expression matches.",
												},
												goaviatrix.TypeKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of resource this expression matches.",
												},
												goaviatrix.ResIdKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource ID this expression matches.",
												},
												goaviatrix.AccountIdKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Account ID this expression matches.",
												},
												goaviatrix.AccountNameKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Account name this expression matches.",
												},
												goaviatrix.NameKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name this expression matches.",
												},
												goaviatrix.RegionKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region this expression matches.",
												},
												goaviatrix.ZoneKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Zone this expression matches.",
												},
												goaviatrix.TagsPrefix: {
													Type:        schema.TypeMap,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Map of tags this expression matches.",
												},
												goaviatrix.K8sClusterIdKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Kubernetes Cluster ID this expression matches.",
												},
												goaviatrix.K8sPodNameKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the Kubernetes Pod this expression matches.",
												},
												goaviatrix.K8sServiceKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the Kubernetes Service this expression matches.",
												},
												goaviatrix.K8sNamespaceKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the Kubernetes Namespace this expression matches.",
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
				},
			},
		},
	}
}

func dataSourceAviatrixSmartGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	smartGroups, err := client.GetSmartGroups(ctx)
	if err != nil {
		return diag.Errorf("could not get Aviatrix Smart Groups: %s", err)
	}

	var result []map[string]interface{}
	for _, smartGroup := range smartGroups {
		var expressions []interface{}

		for _, filter := range smartGroup.Selector.Expressions {
			filterMap := map[string]interface{}{
				goaviatrix.TypeKey:         filter.Type,
				goaviatrix.CidrKey:         filter.CIDR,
				goaviatrix.FqdnKey:         filter.FQDN,
				goaviatrix.SiteKey:         filter.Site,
				goaviatrix.ResIdKey:        filter.ResId,
				goaviatrix.AccountIdKey:    filter.AccountId,
				goaviatrix.AccountNameKey:  filter.AccountName,
				goaviatrix.NameKey:         filter.Name,
				goaviatrix.RegionKey:       filter.Region,
				goaviatrix.ZoneKey:         filter.Zone,
				goaviatrix.TagsPrefix:      filter.Tags,
				goaviatrix.K8sClusterIdKey: filter.K8sClusterId,
				goaviatrix.K8sNamespaceKey: filter.K8sNamespace,
				goaviatrix.K8sServiceKey:   filter.K8sService,
				goaviatrix.K8sPodNameKey:   filter.K8sPodName,
			}

			expressions = append(expressions, filterMap)
		}

		selector := []interface{}{
			map[string]interface{}{
				"match_expressions": expressions,
			},
		}

		smtGroup := map[string]interface{}{
			"name":     smartGroup.Name,
			"uuid":     smartGroup.UUID,
			"selector": selector,
		}

		result = append(result, smtGroup)
	}
	if err = d.Set("smart_groups", result); err != nil {
		return diag.Errorf("couldn't set smart_groups: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}
