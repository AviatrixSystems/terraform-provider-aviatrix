package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
												"cidr": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "CIDR block or IP Address this expression matches.",
												},
												"fqdn": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "FQDN address this expression matches.",
												},
												"site": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Edge Site-ID this expression matches.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of resource this expression matches.",
												},
												"res_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource ID this expression matches.",
												},
												"account_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Account ID this expression matches.",
												},
												"account_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Account name this expression matches.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name this expression matches.",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region this expression matches.",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Zone this expression matches.",
												},
												"tags": {
													Type:        schema.TypeMap,
													Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceAviatrixSmartGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	smartGroups, err := client.GetSmartGroups(ctx)
	if err != nil {
		return diag.Errorf("could not get Aviatrix Smart Groups: %s", err)
	}

	var result []map[string]interface{}
	for _, smartGroup := range smartGroups {
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
