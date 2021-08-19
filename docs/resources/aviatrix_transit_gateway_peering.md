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
  transit_gateway_name1                       = "transit-Gw1"
  transit_gateway_name2                       = "transit-Gw2"
  gateway1_excluded_network_cidrs             = ["10.0.0.48/28"]
  gateway2_excluded_network_cidrs             = ["10.0.0.48/28"]
  gateway1_excluded_tgw_connections           = ["vpn_connection_a"]
  gateway2_excluded_tgw_connections           = ["vpn_connection_b"]
  prepend_as_path1                            = [
    "65001",
    "65001",
    "65001"
  ]
  prepend_as_path2                            = [
    "65002"
  ]
  enable_peering_over_private_network         = false
  enable_insane_mode_encryption_over_internet = false
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
* `prepend_as_path1` - (Optional) AS Path Prepend for BGP connection. Can only use the transit's own local AS number, repeated up to 25 times. Applies on transit_gateway_name1. Available in provider version R2.17.2+.
* `prepend_as_path2` - (Optional) AS Path Prepend for BGP connection. Can only use the transit's own local AS number, repeated up to 25 times. Applies on transit_gateway_name2. Available in provider version R2.17.2+.
* `enable_peering_over_private_network` - (Optional) Enable peering over private network. ActiveMesh and Insane Mode is required on both transit gateways. Available in provider version R2.17.1+.
* `enable_single_tunnel_mode` - (Optional) Enable peering with Single Tunnel mode. False by default. Available as of provider version R2.18+.
* `enable_insane_mode_encryption_over_internet` - (Optional) Enable Insane Mode Encryption over Internet. Type: Boolean. Default: false. Required with `tunnel_count`. Conflicts with `enable_peering_over_private_network` and `enable_single_tunnel_mode`. Available as of provider version R2.19+.
* `tunnel_count` - (Optional) Number of public tunnels. Type: Integer. Valid Range: 2-20. Required with `enable_insane_mode_encryption_over_internet`. Conflicts with `enable_peering_over_private_network` and `enable_single_tunnel_mode`. Available as of provider version R2.19+.
  
~> **NOTE:** `enable_single_tunnel_mode` is only valid when `enable_peering_over_private_network` is set to `true`. Private Transit Gateway Peering with Single-Tunnel Mode expands the existing Insane Mode Transit Gateway Peering Over Private Network to apply it to single IPSec tunnel. One use case is for low speed encryption between cloud networks.

## Import

**transit_gateway_peering** can be imported using the `transit_gateway_name1` and `transit_gateway_name2`, e.g.

```
$ terraform import aviatrix_transit_gateway_peering.test transit_gateway_name1~transit_gateway_name2
```
