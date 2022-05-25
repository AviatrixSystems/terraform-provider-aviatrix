package aviatrix

import (
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixAllNetworkDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixAllNetworkDomainsRead,

		Schema: map[string]*schema.Schema{
			"network_domain_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of network domains and attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network domain name.",
						},
						"tgw_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS TGW name.",
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"egress_inspection_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Egress inspection name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of network domain.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixAllNetworkDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainList, err := client.GetAllSecurityDomains()
	if err != nil {
		return err
	}
	var result []map[string]interface{}
	for i := range domainList {
		domain := domainList[i]
		tempDomain := make(map[string]interface{})
		tempDomain["name"] = domain.Name
		tempDomain["tgw_name"] = domain.TgwName
		tempDomain["route_table_id"] = domain.RouteTableId
		tempDomain["account"] = domain.Account
		tempDomain["cloud_type"] = domain.CouldType
		tempDomain["region"] = domain.Region
		tempDomain["intra_domain_inspection"] = domain.IntraDomainInspectionEnabled
		tempDomain["egress_inspection"] = domain.EgressInspection
		tempDomain["inspection_policy"] = domain.InspectionPolicy
		tempDomain["intra_domain_inspection_name"] = domain.IntraDomainInspectionName
		tempDomain["egress_inspection_name"] = domain.EgressInspectionName
		tempDomain["type"] = domain.Type
		result = append(result, tempDomain)
	}
	if err = d.Set("network_domain_list", result); err != nil {
		return fmt.Errorf("couldn't set network_domain_list: %s", err)
	}
	d.SetId("network_domain_list-id")
	return nil
}
