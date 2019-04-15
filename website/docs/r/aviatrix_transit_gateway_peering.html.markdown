---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway_peering"
sidebar_current: "docs-aviatrix-resource-transit_gateway_peering"
description: |-
  Creates and manages an Aviatrix transit gateway peering.
---

# aviatrix_transit_gateway_peering

The Account resource allows the creation and management of an Aviatrix transit gateway peering.

## Example Usage

```hcl
# Create Aviatrix transit gateway peering
resource "aviatrix_transit_gateway_peering" "foo" {
  transit_gateway_name1 = "transitGw1"
  transit_gateway_name2 = "transitGw2"
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_name1` - (Required) The first transit gateway name to make a peer pair
* `transit_gateway_name2` - (Required) The second transit gateway name to make a peer pair

## Import

Instance transit_vpc can be imported using the transit_gateway_name1 and transit_gateway_name2, e.g.

```
$ terraform import aviatrix_transit_gateway_peering.test transit_gateway_name1~transit_gateway_name2
```


