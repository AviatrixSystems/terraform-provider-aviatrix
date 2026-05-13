---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_transit_attachment"
description: |-
  Creates and manages Aviatrix Spoke-to-Transit attachments
---

# aviatrix_spoke_transit_attachment

The **aviatrix_spoke_transit_attachment** resource allows the creation and management of Aviatrix Spoke-to-Transit gateway attachments.

~> **NOTE:** This resource should only be used to manage the primary gateway attachments. The HA gateway attachments will be handled automatically by the backend.

~> **DEPRECATED:** The `route_tables` argument is deprecated. Attach the spoke to transit first, then configure selective (managed) route tables using the Controller **`edit_gw_config`** action **`edit_managed_route_tables`** (REST API **`edit_gw_config`** / **`edit_managed_route_tables`** with `gateway_name` and `route_tables`). `route_tables` on this resource will show a Terraform warning and **will be removed in a future major provider release**.

## Example Usage

```hcl
# Create an Aviatrix Spoke Transit Attachment (managed route tables: use Controller
# edit_gw_config / edit_managed_route_tables after attach; see Deprecation above)
resource "aviatrix_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"
}
```

```hcl
# Create an Aviatrix Spoke Transit Attachment for EAT
resource "aviatrix_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"
  transit_gateway_logical_ifnames = ["wan1"]
}
```

### Deprecated: `route_tables` on this resource

The following pattern is **deprecated** (Terraform will emit a warning):

```hcl
resource "aviatrix_spoke_transit_attachment" "legacy" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"
  route_tables    = ["rtb-737d540c", "rtb-626d045c"]
}
```

Use **`edit_managed_route_tables`** on **`edit_gw_config`** after the attachment exists instead.

## Argument Reference

The following arguments are supported:

### Required
* `spoke_gw_name` - (Required) Name of the spoke gateway or gateway group to attach to transit network.
* `transit_gw_name` - (Required) Name of the transit gateway or gateway group to attach the spoke gateway to.

### Advanced Options
* `route_tables` - (Optional, **Deprecated**) Learned routes will be propagated to these route tables. Example: `["rtb-212ff547","rtb-04539787"]`. **Deprecated:** use Controller **`edit_gw_config`** with **`edit_managed_route_tables`** after attach; this argument will be removed in a future major release.
* `tunnel_count` - (Optional) Number of public tunnels. Required with both Spoke and Transit to be insane mode enabled and max performance enabled. Type: Integer. Valid Range: 1-49. Available as of provider version R3.1.3+.
* `enable_max_performance` - (Optional) Indicates whether the maximum amount of HPE tunnels will be created. Only valid when transit and spoke gateways are each launched in Insane Mode and in the same cloud type. Default value: true. Available as of provider version R2.22.2+.
* `spoke_prepend_as_path` - (Optional) Connection based AS Path Prepend. Valid only for BGP connection. Can only use the gateway's own local AS number, repeated up to 25 times. Applies on spoke_gateway_name. Available as of provider version R2.23+.
* `transit_prepend_as_path` - (Optional) Connection based AS Path Prepend. Valid only for BGP connection. Can only use the gateway's own local AS number, repeated up to 25 times. Applies on transit_gateway_name. Available as of provider version R2.23+.
* `transit_gateway_logical_ifnames` - (Optional) Transit gateway logical interface names for edge gateways where the peering terminates. Required only for edge as a transit attachment.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `spoke_bgp_enabled` - Indicates whether the spoke gateway is BGP enabled or not.

## Import

**spoke_transit_attachment** can be imported using the `spoke_gw_name` and `transit_gw_name`, e.g.

```
$ terraform import aviatrix_spoke_transit_attachment.test spoke_gw_name~transit_gw_name
```
