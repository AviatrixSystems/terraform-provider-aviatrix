---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_smart_groups"
description: |-
  Gets a list of all Smart Groups.
---

# aviatrix_smart_groups

The **aviatrix_smart_groups** data source provides details about all Smart Groups created by the Aviatrix Controller. Available as of provider version 3.1.2+.

## Example Usage

 ```hcl
 # Aviatrix Smart Groups Data Source
 data "aviatrix_smart_groups" "foo" {}
 ```


## Attribute Reference

The following attributes are exported:
* `smart_groups` - The list of all Smart Groups.
    * `name` - Name of Smart Group.
    * `uuid` - UUID of Smart Group.
    * `selector` - Block containing match expressions to filter the Smart Group.
        * `match_expressions` - List of match expressions. The Smart Group is a union of all resources matched by each `match_expressions`.
            * `cidr` - CIDR block or IP Address this expression matches.
            * `fqdn` - FQDN address this expression matches.
            * `site` - Edge Site-ID this expression matches.
            * `type` - Type of resource this expression matches.
            * `res_id` - Resource ID this expression matches.
            * `account_id` - Account ID this expression matches.
            * `account_name` - Account name this expression matches.
            * `name` - Name this expression matches.
            * `region` - Region this expression matches.
            * `zone` - Zone this expression matches.
            * `tags` - Map of tags this expression matches.