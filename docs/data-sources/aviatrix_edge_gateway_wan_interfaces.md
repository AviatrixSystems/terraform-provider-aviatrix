---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_gateway_wan_interfaces"
description: |-
  Gets the Aviatrix Edge gateway WAN interfaces.
---

# aviatrix_edge_gateway_wan_interfaces

Use this data source to get the Edge gateway WAN interfaces for use in other resources.

## Example Usage

```hcl
# Aviatrix Edge Gateway WAN Interfaces Data Source
data "aviatrix_edge_gateway_wan_interfaces" "test" {
  gw_name = "test-gw"
}
```

## Argument Reference

The following argument are supported:

* `gw_name` - (Required) Edge gateway name.

## Attribute Reference

In addition to the argument above, the following attributes is exported:

* `wan_interfaces` - Set of the WAN interfaces.
