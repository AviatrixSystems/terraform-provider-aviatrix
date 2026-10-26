---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_gateway_keepalive_config"
description: |-
  Creates and manages an Aviatrix Controller Gateway Keepalive for gateways
---

# aviatrix_controller_gateway_keepalive_config

The **aviatrix_controller_gateway_keepalive_config** resource allows management of an Aviatrix Controller's gateway keepalive template configuration. This resource is available as of provider version R2.19.2+.

## Example Usage

```hcl
# Create an Aviatrix Controller Gateway Keepalive config
resource "aviatrix_controller_gateway_keepalive_config" "test_gateway_keepalive" {
  keepalive_speed = "medium"
}
```


## Argument Reference

The following arguments are supported:

### Required
* `keepalive_speed` - The gateway keepalive template name. Must be one of "slow", "medium" or "fast". Visit [here](https://docs.aviatrix.com/HowTos/gateway.html#gateway-keepalives) for the complete documentation about the gateway keepalive configuration.

## Import

**aviatrix_controller_gateway_keepalive_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_gateway_keepalive_config.test_gateway_keepalive 10-11-12-13
```
