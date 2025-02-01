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
```hcl
# Create an Aviatrix Edge Transit Gateway Peering
resource "aviatrix_transit_gateway_peering" "test_edge_transit_gateway_peering" {
  transit_gateway_name1   = "test-edge-transit-1"
  transit_gateway_name2   = "test-edge-transit-2"
  enable_peering_over_private_network  = true
  jumbo_frame             = false
  insane_mode             = true
  gateway1_logical_ifnames      = ["wan1"]
  gateway2_logical_ifnames      = ["wan1"]
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
* `enable_peering_over_private_network` - (Optional) Advanced option. Enable peering over private network. Only appears and applies to when the two Multi-cloud Transit Gateways are each launched in Insane Mode and in a different cloud type. Conflicts with `enable_insane_mode_encryption_over_internet` and `tunnel_count`. Type: Boolean. Default: false. Available in provider version R2.17.1+.
* `enable_single_tunnel_mode` - (Optional) Advanced option. Enable peering with Single-Tunnel mode. Only appears and applies to when the two Multi-cloud Transit Gateways are each launched in Insane Mode and in a different cloud type. Required with `enable_peering_over_private_network`. Conflicts with `enable_insane_mode_encryption_over_internet` and `tunnel_count`. Type: Boolean. Default: false. Available as of provider version R2.18+.
* `enable_insane_mode_encryption_over_internet` - (Optional) Advanced option. Enable Insane Mode Encryption over Internet. Transit gateways must be in Insane Mode. Currently, only inter-cloud connections between AWS and Azure are supported. Required with valid `tunnel_count`. Conflicts with `enable_peering_over_private_network` and `enable_single_tunnel_mode`. Type: Boolean. Default: false. Available as of provider version R2.19+.
* `tunnel_count` - (Optional) Advanced option. Number of public tunnels. Required with `enable_insane_mode_encryption_over_internet`. Conflicts with `enable_peering_over_private_network` and `enable_single_tunnel_mode`. Type: Integer. Valid Range: 2-20. Available as of provider version R2.19+.
* `enable_max_performance` - (Optional) Indicates whether the maximum amount of HPE tunnels will be created. Only valid when the two transit gateways are each launched in Insane Mode and in the same cloud type. Default value: true. Available as of provider version R2.22.2+.
* `jumbo_frame` - (Optional) Enable jumbo frame for over private peering with Edge Transit. Required only for edge transit peerings.
* `insane_mode` - (Optional) Enable HPE mode for peering with Edge Transit. Required only for edge transit peerings.
* `gateway1_logical_ifnames` - (Optional) Logical source WAN interfaces for edge gateways where the peering originates. Required only for edge transit attachments.
* `gateway2_logical_ifnames` - (Optional) Logical destination WAN interface for edge gateways where the peering terminates. Required only for edge transit attachments.


~> **NOTE:** `enable_single_tunnel_mode` is only valid when `enable_peering_over_private_network` is set to `true`. Private Transit Gateway Peering with Single-Tunnel Mode expands the existing Insane Mode Transit Gateway Peering Over Private Network to apply it to single IPSec tunnel. One use case is for low speed encryption between cloud networks.

## Import

**transit_gateway_peering** can be imported using the `transit_gateway_name1` and `transit_gateway_name2`, e.g.

```
$ terraform import aviatrix_transit_gateway_peering.test transit_gateway_name1~transit_gateway_name2
```
