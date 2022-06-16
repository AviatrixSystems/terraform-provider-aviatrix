package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixNetworkDomains() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAviatrixNetworkDomainsRead,

		Schema: map[string]*schema.Schema{
			"network_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Network Domains.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network Domain name.",
						},
						"tgw_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS TGW name.",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route table id.",
						},
						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access Account name.",
						},
						"cloud_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of cloud service provider.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of cloud provider.",
						},
						"intra_domain_inspection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Firewall inspection for traffic within one Security Domain.",
						},
						"egress_inspection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Egress inspection is enable or not.",
						},
						"inspection_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Inspection policy name.",
						},
						"intra_domain_inspection_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intra domain inspection name.",
						},
						"egress_inspection_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Egress inspection name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Network Domain.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixNetworkDomainsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	domainList, err := client.GetAllNetworkDomains(ctx)
	if err != nil {
		return diag.Errorf("could not get Aviatrix Network Domains: %s", err)
	}
	var result []map[string]interface{}
	for i, domain := range domainList {
		domainList[i] = domain
		tempDomain := map[string]interface{}{
			"name":                         domain.Name,
			"tgw_name":                     domain.TgwName,
			"route_table_id":               domain.RouteTableId,
			"account":                      domain.Account,
			"cloud_type":                   domain.CouldType,
			"region":                       domain.Region,
			"intra_domain_inspection":      domain.IntraDomainInspectionEnabled,
			"egress_inspection":            domain.EgressInspection,
			"inspection_policy":            domain.InspectionPolicy,
			"intra_domain_inspection_name": domain.IntraDomainInspectionName,
			"egress_inspection_name":       domain.EgressInspectionName,
			"type":                         domain.Type,
		}
		result = append(result, tempDomain)
	}
	if err = d.Set("network_domains", result); err != nil {
		return diag.Errorf("couldn't set network_domains: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}
