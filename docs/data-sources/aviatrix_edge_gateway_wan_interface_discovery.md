---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_gateway_wan_interface_discovery"
description: |-
  Gets the Aviatrix Edge gateway WAN interface public IP address.
---

# aviatrix_edge_gateway_wan_interface_discovery

Use this data source to get the Edge gateway WAN interface public IP address for use in other resources.

## Example Usage

```hcl
# Aviatrix Edge Gateway WAN Interface Discovery Data Source
data "aviatrix_edge_gateway_wan_interface_discovery" "test" {
  gw_name            = "test-gw"
  wan_interface_name = "eth0"
}
```

## Argument Reference

The following argument are supported:

* `gw_name` - (Required) Edge gateway name.
* `wan_interface_name` - (Required) Name of the WAN interface to be discovered.

## Attribute Reference

In addition to the argument above, the following attributes is exported:

* `ip_address` - Public IP of the Edge gateway's WAN interface.
