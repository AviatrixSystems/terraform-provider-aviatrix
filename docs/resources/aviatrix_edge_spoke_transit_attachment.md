---
subcategory: "Edge"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_spoke_transit_attachment"
description: |-
  Creates and manages Aviatrix Edge as a Spoke to Transit attachments
---

# aviatrix_edge_spoke_transit_attachment

The **aviatrix_edge_spoke_transit_attachment** resource allows the creation and management of Aviatrix Edge as a Spoke to Transit gateway attachments. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Create an Aviatrix Edge as a Spoke Transit Attachment
resource "aviatrix_edge_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = "edge-as-a-spoke"
  transit_gw_name = "transit-gw"
}
```
```hcl
# Create an Aviatrix Edge as a Spoke to Edge as a Transit Attachment
resource "aviatrix_edge_spoke_transit_attachment" "test_attachment_2" {
  spoke_gw_name               = "e2e-edge-spoke-1"
  transit_gw_name             = "test-edge-transit-1"
  enable_over_private_network = true
  enable_jumbo_frame          = false
  enable_insane_mode          = true
  edge_wan_interfaces         = ["eth0"]
  transit_gateway_logical_ifnames = ["wan1"]
}
```
```hcl
# Create an Aviatrix Edge as a Spoke to Edge as a Transit Attachment for Megaport
resource "aviatrix_edge_spoke_transit_attachment" "test_attachment_2" {
  spoke_gw_name                   = "test-edge-spoke-1"
  transit_gw_name                 = "test-edge-transit-1"
  enable_over_private_network     = true
  enable_jumbo_frame              = false
  enable_insane_mode              = true
  spoke_gateway_logical_ifnames   = ["wan1"]
  transit_gateway_logical_ifnames = ["wan1"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `spoke_gw_name` - (Required) Name of the Edge as a Spoke to attach to transit network.
* `transit_gw_name` - (Required) Name of the transit gateway to attach the Edge as a Spoke to.

### Options

* `enable_over_private_network` - (Optional) Switch to enable over the private network. Valid values: true, false. Default: true.
* `enable_jumbo_frame` - (Optional) Switch to enable jumbo frame. Valid values: true, false. Default: false.
* `enable_insane_mode` - (Optional) Switch to enable insane mode. Valid values: true, false. Default: false.
* `insane_mode_tunnel_number` - (Optional) Insane mode tunnel number, requires an integer value. Valid range for HPE over private network: 0-49. Valid range for HPE over internet: 2-20.
* `spoke_prepend_as_path` - (Optional) Connection based AS Path Prepend. Can only use the gateway's own local AS number, repeated up to 25 times. Applies on the Edge as a Spoke.
* `transit_prepend_as_path` - (Optional) Connection based AS Path Prepend. Can only use the gateway's own local AS number, repeated up to 25 times. Applies on the Transit Gateway.
* `number_of_retries` - (Optional) Number of retries. Default value: 0.
* `retry_interval` - (Optional) Retry interval in seconds. Default value: 300.
* `edge_wan_interfaces` - (Optional) Set of Edge WAN interfaces.
* `spoke_gateway_logical_ifnames` - (Optional) Spoke gateway logical interface names for edge gateways where the peering originates. Required only for Megaport edge as a transit attachment.
* `transit_gateway_logical_ifnames` - (Optional) Transit gateway logical interface names for edge gateways where the peering terminates. Required only for edge as a transit attachment.

## Import

**spoke_transit_attachment** can be imported using the `spoke_gw_name` and `transit_gw_name`, e.g.

```
$ terraform import aviatrix_spoke_transit_attachment.test spoke_gw_name~transit_gw_name
```
