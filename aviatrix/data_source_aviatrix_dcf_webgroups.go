package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixDcfWebgroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDcfWebgroupsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Web Group.",
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
									"snifilter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server name indicator this expression matches.",
									},
									"urlfilter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL address this expression matches.",
									},
									goaviatrix.ExternalURLKey: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the external domain feed whose URLs this expression matches.",
									},
									goaviatrix.ExternalSNIKey: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the external domain feed whose server names this expression matches.",
									},
									goaviatrix.ExtArgsPrefix: {
										Type:        schema.TypeMap,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Map of arguments passed to the external feed.",
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

func dataSourceAviatrixDcfWebgroupsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	name := getString(d, "name")

	if name == "" {
		return diag.Errorf("name must be specified")
	}

	webGroup, err := client.GetWebGroupByName(ctx, name)
	if err != nil {
		return diag.Errorf("could not get DCF webgroups: %s", err)
	}

	err = d.Set("name", webGroup.Name)
	if err != nil {
		return diag.Errorf("could not set name: %s", err)
	}
	err = d.Set("uuid", webGroup.UUID)
	if err != nil {
		return diag.Errorf("could not set uuid: %s", err)
	}
	if err := d.Set("selector", flattenWebGroupSelector(webGroup)); err != nil {
		return diag.Errorf("failed to set selector during Web Group read: %s", err)
	}

	d.SetId(webGroup.UUID)

	return nil
}
