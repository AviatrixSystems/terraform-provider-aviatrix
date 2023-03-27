---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_smart_group"
description: |-
  Creates and manages an Aviatrix Smart Group
---

# aviatrix_smart_group

The **aviatrix_smart_group** resource handles the creation and management of Smart Groups. Available as of Provider R2.22.0+.

## Example Usage

```hcl
# Create an Aviatrix Smart Group
resource "aviatrix_smart_group" "test_smart_group_ip" {
  name = "smart-group"
  selector {
    match_expressions {
      type         = "vm"
      account_name = "devops"
      region       = "us-west-2"
      tags         = {
        k3 = "v3"
      }
    }

    match_expressions {
      cidr = "10.0.0.0/16"
    }

    match_expressions {
      fqdn = "www.aviatrix.com"
    }

    match_expressions {
      site = "site-test-0"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Smart Group.
* `selector` - (Required) Block containing match expressions to filter the Smart Group.
  * `match_expressions` - (Required) List of match expressions. The Smart Group will be a union of all resources matched by each `match_expressions`.`match_expressions` blocks cannot be empty.
    * `cidr` - (Optional) - CIDR block or IP Address this expression matches. `cidr` cannot be used with any other filters in the same `match_expressions` block.
    * `fqdn` - (Optional) - FQDN address this expression matches. `fqdn` cannot be used with any other filters in the same `match_expressions` block.
    * `site` - (Optional) - Edge Site-ID this expression matches. `site` cannot be used with any other filters in the same `match_expressions` block.
    * `type` - (Optional) - Type of resource this expression matches. Must be one of "vm", "vpc" or "subnet". `type` is required when `cidr`, `fqdn` and `site` are all not used.
    * `res_id` - (Optional) - Resource ID this expression matches.
    * `account_id` - (Optional) - Account ID this expression matches.
    * `account_name` - (Optional) - Account name this expression matches.
    * `name` - (Optional) - Name this expression matches.
    * `region` - (Optional) - Region this expression matches.
    * `zone` - (Optional) - Zone this expression matches.
    * `tags` - (Optional) - Map of tags this expression matches.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the Smart Group.

## Import

**aviatrix_smart_group** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_smart_group.test 41984f8b-5a37-4272-89b3-57c79e9ff77c
```
