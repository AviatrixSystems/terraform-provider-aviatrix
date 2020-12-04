---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway_peering"
description: |-
  Creates and manages Aviatrix transit gateway peerings
---

# aviatrix_transit_gateway_peering

The **aviatrix_transit_gateway_peering** resource allows the creation and management of peerings between Aviatrix transit gateways.

## Example Usage

```hcl
# Create an Aviatrix Transit Gateway Peering
resource "aviatrix_transit_gateway_peering" "test_transit_gateway_peering" {
  transit_gateway_name1               = "transit-Gw1"
  transit_gateway_name2               = "transit-Gw2"
  gateway1_excluded_network_cidrs     = ["10.0.0.48/28"]
  gateway2_excluded_network_cidrs     = ["10.0.0.48/28"]
  gateway1_excluded_tgw_connections   = ["vpn_connection_a"]
  gateway2_excluded_tgw_connections   = ["vpn_connection_b"]
  prepend_as_path1                    = [
    "111",
    "222"
  ]
  prepend_as_path2                    = [
    "333",
    "444"
  ]
  enable_peering_over_private_network = false
}
```

## Argument Reference

The following arguments are supported:

### Required
* `transit_gateway_name1` - (Required) The first transit gateway name to make a peer pair.
* `transit_gateway_name2` - (Required) The second transit gateway name to make a peer pair.

### Optional
* `gateway1_excluded_network_cidrs` - (Optional) List of excluded network CIDRs for the first transit gateway.
* `gateway2_excluded_network_cidrs` - (Optional) List of excluded network CIDRs for the second transit gateway.
* `gateway1_excluded_tgw_connections` - (Optional) List of excluded TGW connections for the first transit gateway.
* `gateway2_excluded_tgw_connections` - (Optional) List of excluded TGW connections for the second transit gateway.
* `prepend_as_path1` - (Optional) AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit_gateway_name1.
* `prepend_as_path2` - (Optional) AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit_gateway_name2.
* `enable_peering_over_private_network` - (Optional) Enable peering over private network. ActiveMesh and Insane Mode is required on both transit gateways. Available in provider version R2.17.1+.

## Import

**transit_gateway_peering** can be imported using the `transit_gateway_name1` and `transit_gateway_name2`, e.g.

```
$ terraform import aviatrix_transit_gateway_peering.test transit_gateway_name1~transit_gateway_name2
```
