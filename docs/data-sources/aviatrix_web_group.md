---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_web_group"
description: |-
    Gets details about a DCF Web Group.
---

# aviatrix_web_group

The **aviatrix_web_group** data source provides details about a specific DCF Web Group created by the Aviatrix Controller.

## Example Usage

```hcl
# Aviatrix Web Group Data Source
data "aviatrix_web_group" "example" {
    name = "my-web-group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Web Group.

## Attribute Reference

The following attributes are exported:

* `name` - Name of the Web Group.
* `uuid` - UUID of the Web Group.
* `selector` - Block containing match expressions to filter the Web Group.
        * `match_expressions` - List of match expressions for the Web Group.
                * `snifilter` - Server name indication this expression matches.
                * `urlfilter` - URL address this expression matches.
