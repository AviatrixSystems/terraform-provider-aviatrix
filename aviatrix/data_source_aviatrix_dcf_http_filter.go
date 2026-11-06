package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixDCFHTTPFilter() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixDCFHTTPFilterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the HTTP filter.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the HTTP filter.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the HTTP filter.",
			},
			"filters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Clauses of the HTTP filter. The filter matches when any clause matches.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Header matches that must all hold for the clause to match.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP header name to match.",
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Header values to match.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixDCFHTTPFilterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	name := getString(d, "name")
	filter, err := client.GetDCFHTTPFilterByName(ctx, name)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return diag.Errorf("DCF HTTP filter with name %q not found", name)
		}
		return diag.Errorf("failed to read DCF HTTP filter: %v", err)
	}

	desc := ""
	if filter.Desc != nil {
		desc = *filter.Desc
	}
	mustSet(d, "uuid", filter.UUID)
	mustSet(d, "description", desc)
	mustSet(d, "filters", flattenDCFHTTPFilterClauses(filter.Filters))
	d.SetId(filter.UUID)

	return nil
}
