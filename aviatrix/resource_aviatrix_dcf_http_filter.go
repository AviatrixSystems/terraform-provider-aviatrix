package aviatrix

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// dcfHTTPFilterHeaderRegexp accepts the supported HTTP/2 pseudo-headers or a
// lowercase RFC 9110 header field name token.
var dcfHTTPFilterHeaderRegexp = regexp.MustCompile("^(:(method|scheme|authority|path)|[a-z0-9!#$%&'*+.^_`|~-]+)$")

func resourceAviatrixDCFHTTPFilter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFHTTPFilterCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFHTTPFilterRead,
		UpdateWithoutTimeout: resourceAviatrixDCFHTTPFilterUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFHTTPFilterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Unique name of the HTTP filter.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the HTTP filter.",
			},
			"filters": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Description: "Clauses of the HTTP filter. The filter matches a request when any clause matches; " +
					"a clause matches when all of its header matches hold.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all": {
							Type:        schema.TypeList,
							Required:    true,
							MinItems:    1,
							Description: "Header matches that must all hold for the clause to match.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(dcfHTTPFilterHeaderRegexp,
											"must be a lowercase HTTP header name or one of the pseudo-headers :method, :scheme, :authority, :path"),
										Description: "HTTP header name to match. Must be a lowercase HTTP header name or one of the " +
											"pseudo-headers :method, :scheme, :authority, :path.",
									},
									"values": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Header values to match.",
									},
								},
							},
						},
					},
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the HTTP filter.",
			},
		},
	}
}

// marshalDCFHTTPFilterInput builds the HTTP filter API request from resource data.
func marshalDCFHTTPFilterInput(d *schema.ResourceData) *goaviatrix.DCFHTTPFilter {
	filter := &goaviatrix.DCFHTTPFilter{
		Name:    getString(d, "name"),
		Filters: expandDCFHTTPFilterClauses(getList(d, "filters")),
	}
	if desc := getString(d, "description"); desc != "" {
		filter.Desc = &desc
	}
	return filter
}

// expandDCFHTTPFilterClauses converts the filters blocks into API clauses.
func expandDCFHTTPFilterClauses(clauses []any) []goaviatrix.DCFHTTPFilterAndClause {
	result := make([]goaviatrix.DCFHTTPFilterAndClause, 0, len(clauses))
	for _, c := range clauses {
		clause := goaviatrix.DCFHTTPFilterAndClause{}
		for _, e := range mustSlice(mustMap(c)["all"]) {
			elem := mustMap(e)
			clause.All = append(clause.All, goaviatrix.DCFHTTPFilterElement{
				Hdr:    mustString(elem["header"]),
				Values: expandStringList(mustSlice(elem["values"])),
			})
		}
		result = append(result, clause)
	}
	return result
}

// flattenDCFHTTPFilterClauses converts API clauses into filters blocks.
func flattenDCFHTTPFilterClauses(clauses []goaviatrix.DCFHTTPFilterAndClause) []any {
	result := make([]any, 0, len(clauses))
	for _, clause := range clauses {
		all := make([]any, 0, len(clause.All))
		for _, elem := range clause.All {
			all = append(all, map[string]any{
				"header": elem.Hdr,
				"values": elem.Values,
			})
		}
		result = append(result, map[string]any{"all": all})
	}
	return result
}

func resourceAviatrixDCFHTTPFilterCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	filterUUID, err := client.CreateDCFHTTPFilter(ctx, marshalDCFHTTPFilterInput(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(filterUUID)
	return resourceAviatrixDCFHTTPFilterRead(ctx, d, meta)
}

func resourceAviatrixDCFHTTPFilterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	filter, err := client.GetDCFHTTPFilter(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF HTTP filter: %v", err)
	}

	desc := ""
	if filter.Desc != nil {
		desc = *filter.Desc
	}
	mustSet(d, "uuid", filter.UUID)
	mustSet(d, "name", filter.Name)
	mustSet(d, "description", desc)
	mustSet(d, "filters", flattenDCFHTTPFilterClauses(filter.Filters))

	return nil
}

func resourceAviatrixDCFHTTPFilterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if err := client.UpdateDCFHTTPFilter(ctx, d.Id(), marshalDCFHTTPFilterInput(d)); err != nil {
		return diag.FromErr(err)
	}

	return resourceAviatrixDCFHTTPFilterRead(ctx, d, meta)
}

func resourceAviatrixDCFHTTPFilterDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if err := client.DeleteDCFHTTPFilter(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
