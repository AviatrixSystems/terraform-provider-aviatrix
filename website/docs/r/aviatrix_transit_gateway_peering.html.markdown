---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway_peering"
description: |-
  Creates and manages Aviatrix transit gateway peerings
---

# aviatrix_transit_gateway_peering

The aviatrix_transit_gateway_peering resource allows the creation and management of Aviatrix transit gateway peerings.

## Example Usage

```hcl
# Create an Aviatrix Transit Gateway Peering
resource "aviatrix_transit_gateway_peering" "test_transit_gateway_peering" {
  transit_gateway_name1 = "transitGw1"
  transit_gateway_name2 = "transitGw2"
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_name1` - (Required) The first transit gateway name to make a peer pair.
* `transit_gateway_name2` - (Required) The second transit gateway name to make a peer pair.

## Import

Instance transit_gateway_peering can be imported using the transit_gateway_name1 and transit_gateway_name2, e.g.

```
$ terraform import aviatrix_transit_gateway_peering.test transit_gateway_name1~transit_gateway_name2
```
