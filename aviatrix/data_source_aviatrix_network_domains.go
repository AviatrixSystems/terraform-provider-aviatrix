package aviatrix

import (
	"fmt"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixNetworkDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixNetworkDomainsRead,

		Schema: map[string]*schema.Schema{
			"network_domain_list": {
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

func dataSourceAviatrixNetworkDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainList, err := client.GetAllSecurityDomains()
	if err != nil {
		return err
	}
	var result []map[string]interface{}
	for i := range domainList {
		domain := domainList[i]
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
	if err = d.Set("network_domain_list", result); err != nil {
		return fmt.Errorf("couldn't set network_domain_list: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}
