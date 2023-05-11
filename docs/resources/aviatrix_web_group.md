---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_web_group"
description: |-
  Creates and manages an Aviatrix Web Group
---

# aviatrix_web_group

The **aviatrix_web_group** resource handles the creation and management of Web Groups. Available as of Provider R3.1.0+.

## Example Usage

```hcl
# Create an Aviatrix SNI Web Group
resource "aviatrix_web_group" "test_web_group_ip" {
  name = "web-group"
  selector {
    match_expressions {
      snifilter = "aviatrix.com"
    }
  }
}
```
```hcl
# Create an Aviatrix URL Web Group
resource "aviatrix_web_group" "test_web_group_ip" {
  name = "web-group"
  selector {
    match_expressions {
      urlfilter = "https://aviatrix.com/test"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Web Group.
* `selector` - (Required) Block containing match expressions to filter the Web Group.
    * `match_expressions` - (Required) List of match expressions. The Web Group will be a union of all resources matched by each `match_expressions`.`match_expressions` blocks cannot be empty.
        * `snifilter` - (Optional) - Server name indicator this expression matches. `snifilter` cannot be used with `urlfilter` filters in the same `match_expressions` block.
        * `urlfilter` - (Optional) - URL address this expression matches. `urlfilter` cannot be used with `snifilter` filters in the same `match_expressions` block.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the Web Group.

## Import

**aviatrix_web_group** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_web_group.test 41984f8b-5a37-4272-89b3-57c79e9ff77c
```
