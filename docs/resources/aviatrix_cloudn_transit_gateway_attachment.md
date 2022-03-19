---
subcategory: "CloudN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_cloudn_transit_gateway_attachment"
description: |-
  Create and manage CloudN Transit Gateway Attachments
---

# aviatrix_cloudn_transit_gateway_attachment

The **aviatrix_cloudn_transit_gateway_attachment** resource allows the creation and management of CloudN Transit Gateway Attachments. This resource is available as of provider version R2.19+.

## Example Usage

```hcl
# Create a CloudN Transit Gateway Attachment
resource "aviatrix_cloudn_transit_gateway_attachment" "test" {
  device_name                           = aviatrix_device_registration.test_device.name
  transit_gateway_name                  = aviatrix_transit_gateway.aws_transit.gw_name
  connection_name                       = "cloudn-transit-attachment-test" 
  transit_gateway_bgp_asn               = "65000"
  cloudn_bgp_asn                        = "65046"
  cloudn_lan_interface_neighbor_ip      = "10.210.38.100"
  cloudn_lan_interface_neighbor_bgp_asn = "65219"
  enable_over_private_network           = false 
  enable_jumbo_frame                    = true 
  enable_dead_peer_detection            = true 
}
```

## Argument Reference

The following arguments are supported:

### Required
* `device_name` - (Required) CloudN device name. Type: String.
* `transit_gateway_name` - (Required) Transit Gateway Name. Type: String.
* `connection_name` - (Required) Connection Name. Type: String.
* `transit_gateway_bgp_asn` - (Required) Transit Gateway BGP AS Number. Type: String.
* `cloudn_bgp_asn` - (Required) CloudN BGP AS Number. Type: String.
* `cloudn_lan_interface_neighbor_ip` - (Required) CloudN LAN Interface Neighbor's IP Address. Type: String.
* `cloudn_lan_interface_neighbor_bgp_asn` - (Required) CloudN LAN Interface Neighbor's AS Number. Type: String.
* `enable_over_private_network` - (Required) Enable connection over private network. Type: Boolean.

### Optional
* `enable_jumbo_frame` - (Optional) Enable Jumbo Frame support for the connection. Type: Boolean. Default: false.
* `enable_dead_peer_detection` - (Optional) Enable Dead Peer Detection. Type: Boolean. Default: true.
* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval. Type: Boolean. Default: false. Available as of provider version R2.21.0+.
* `approved_cidrs` - (Optional/Computed) Set of approved CIDRs. Requires `enable_learned_cidrs_approval` to be true. Type: Set(String). Available as of provider version R2.21.0+.
* `prepend_as_path` - (Optional)  Connection AS Path Prepend customized by specifying AS PATH for a BGP connection. Requires transit_gateway_bgp_asn to be set. Type: List. Available as of provider version R2.21.0+.

## Import

**aviatrix_cloudn_transit_gateway_attachment** can be imported using the `connection_name`, e.g.

```
$ terraform import aviatrix_cloudn_transit_gateway_attachment.test connection_name
```
