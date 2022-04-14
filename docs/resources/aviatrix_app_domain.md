---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_app_domain"
description: |-
  Creates and manages an Aviatrix App Domain
---

# aviatrix_app_domain

The **aviatrix_app_domain** resource handles the creation and management of App Domains. Available as of Provider R2.22.0+.

## Example Usage

```hcl
# Create an Aviatrix App Domain
resource "aviatrix_app_domain" "test_app_domain_ip" {
  name = "app-domain"
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
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the App Domain.
* `selector` - (Required) Block containing match expressions to filter the App Domain.
  * `match_expressions` - (Required) List of match expressions. The App Domain will be a union of all resources matched by each `match_expressions`.`match_expressions` blocks cannot be empty.
      * `cidr` - (Optional) - CIDR block or IP Address this expression matches. `cidr` cannot be used with any other filters in the same `match_expressions` block.
      * `type` - (Optional) - Type of resource this expression matches. Must be one of "vm", "vpc" or "subnet". `type` is required when `cidr` is not used.
      * `res_id` - (Optional) - Resource ID this expression matches.
      * `account_id` - (Optional) - Account ID this expression matches.
      * `account_name` - (Optional) - Account name this expression matches.
      * `region` - (Optional) - Region this expression matches.
      * `zone` - (Optional) - Zone this expression matches.
      * `tags` - (Optional) - Map of tags this expression matches.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the App Domain.

## Import

**aviatrix_app_domain** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_app_domain.test 41984f8b-5a37-4272-89b3-57c79e9ff77c
```
