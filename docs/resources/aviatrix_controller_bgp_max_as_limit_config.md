---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_bgp_max_as_limit_config"
description: |-
  Creates and manages an Aviatrix controller BGP max AS limit for transit gateways
---

# aviatrix_controller_bgp_max_as_limit_config

The **aviatrix_controller_bgp_max_as_limit_config** resource allows management of an Aviatrix Controller's BGP max AS limit for transit gateways. This resource is available as of provider version R2.18.1+.

## Example Usage

```hcl
# Create an Aviatrix Controller BGP max AS limit config
resource "aviatrix_controller_bgp_max_as_limit_config" "test_max_as_limit" {
  max_as_limit                = 1
  max_as_limit_non_rfc1918    = 2
}
```


## Argument Reference

The following arguments are supported:

* `max_as_limit` - (Optional) The maximum AS path limit allowed by transit gateways when handling BGP/Peering route propagation with RFC1918 CIDRs. Must be a number in the range [1-254].
* `max_as_limit_non_rfc1918` - (Optional) The maximum AS path limit allowed by transit gateways when handling BGP/Peering route propagation with non-RFC1918 CIDRs. Must be a number in the range [1-254].

## Import

**aviatrix_controller_bgp_max_as_limit_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_bgp_max_as_limit_config.test_max_as_limit 10-11-12-13
```
