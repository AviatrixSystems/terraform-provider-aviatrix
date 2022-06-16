---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_network_domains"
description: |-
  Gets a list of all Network Domains.
---

# aviatrix_all_network_domains

The **aviatrix_network_domains** data source provides details about all Network Domains created by the Aviatrix Controller. Available as of provider version 2.23+.

## Example Usage

 ```hcl
 # Aviatrix All Network Domains Data Source
 data "aviatrix_network_domains" "foo" {}
 ```


## Attribute Reference

The following attributes are exported:
* `network_domains` - The list of all Network Domains
    * `name` - Network Domain name.
    * `tgw_name` - AWS TGW name.
    * `account` - Access Account name.
    * `route_table_id` - Route table's id.
    * `cloud_type` - Type of cloud service provider.
    * `region` - Region of cloud provider.
    * `intra_domain_inspection` - Firewall inspection for traffic within one Security Domain.
    * `egress_inspection` - Egress inspection is enable or not.
    * `inspection_policy` - Inspection policy name.
    * `intra_domain_inspection_name` - Intra domain inspection name.
    * `egress_inspection_name` - Egress inspection name.
    * `type` - Type of network domain.