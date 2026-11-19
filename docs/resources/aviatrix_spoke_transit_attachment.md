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

~> **DEPRECATED:** The `route_tables` argument is deprecated. Use `route_tables` argument in aviatrix_spoke_gateway or aviatrix_spoke_group resource to set route tables. `route_tables` on this resource will show a Terraform warning and **will be removed in a future major provider release**.

## Example Usage

```hcl
# Create an Aviatrix Spoke Transit Attachment
resource "aviatrix_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"
}
```

```hcl
# Selective route tables belong on the spoke gateway or spoke group.
# route_tables on aviatrix_spoke_transit_attachment is deprecated.
resource "aviatrix_spoke_gateway" "spoke" {
  cloud_type   = 1
  account_name = "my-aws"
  gw_name      = "spoke-gw"
  vpc_id       = "vpc-abcd1234"
  vpc_reg      = "us-west-1"
  gw_size      = "t3.micro"
  subnet       = "10.11.0.0/24"

  # AWS route table IDs, or Azure entries as "<route_table_name>:<resource_group_name>".
  # route_tables = [] clears selective configuration.
  route_tables = ["rtb-737d540c", "rtb-626d045c"]
}

resource "aviatrix_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = aviatrix_spoke_gateway.spoke.gw_name
  transit_gw_name = "transit-gw"
}
```

On a spoke group, set `route_tables` on [`aviatrix_spoke_group`](aviatrix_spoke_group.md) the same way.

```hcl
# Create an Aviatrix Spoke Transit Attachment for EAT
resource "aviatrix_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"
  transit_gateway_logical_ifnames = ["wan1"]
}
```

```hcl
### Directional route filtering

Use `spoke_to_transit_filter_cidrs` and `transit_to_spoke_filter_cidrs` to prevent specific IPv4 or IPv6 prefixes from being advertised across the attachment.

In this example, the spoke does not advertise `10.10.0.0/16` to the transit gateway, and the transit gateway does not advertise `172.16.0.0/12` to the spoke. The filter lists can be updated without recreating the attachment.

resource "aviatrix_spoke_transit_attachment" "filtered" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"

  spoke_to_transit_filter_cidrs = [
    "10.10.0.0/16",
  ]

  transit_to_spoke_filter_cidrs = [
    "172.16.0.0/12",
  ]
}
```

### Deprecated: `route_tables` on this resource

Setting `route_tables` on this attachment is **deprecated** (Terraform will emit a warning). Set [`aviatrix_spoke_gateway.route_tables`](aviatrix_spoke_gateway.md) or [`aviatrix_spoke_group.route_tables`](aviatrix_spoke_group.md) instead, as in the example above. This argument will be removed in a future major release.

## Argument Reference

The following arguments are supported:

### Required
* `spoke_gw_name` - (Required) Name of the spoke gateway or gateway group to attach to transit network. Switching between a gateway name and its group name (both resolving to the same gateway) is applied in place without recreating the attachment; changing to a different gateway forces recreation.
* `transit_gw_name` - (Required) Name of the transit gateway or gateway group to attach the spoke gateway to. Switching between a gateway name and its group name (both resolving to the same gateway) is applied in place without recreating the attachment; changing to a different gateway forces recreation.

### Advanced Options
* `route_tables` - (Optional, **Deprecated**) Learned routes will be propagated to these route tables. Example: `["rtb-212ff547","rtb-04539787"]`. **Deprecated:** set `route_tables` on [`aviatrix_spoke_gateway`](aviatrix_spoke_gateway.md) or [`aviatrix_spoke_group`](aviatrix_spoke_group.md) instead. This argument will be removed in a future major release.
* `tunnel_count` - (Optional) Number of public tunnels. Required with both Spoke and Transit to be insane mode enabled and max performance enabled. Type: Integer. Valid Range: 1-49. Available as of provider version R3.1.3+.
* `enable_max_performance` - (Optional) Indicates whether the maximum amount of HPE tunnels will be created. Only valid when transit and spoke gateways are each launched in Insane Mode and in the same cloud type. Default value: true. Available as of provider version R2.22.2+.
* `spoke_prepend_as_path` - (Optional) Connection based AS Path Prepend. Valid only for BGP connection. Can only use the gateway's own local AS number, repeated up to 25 times. Applies on spoke_gateway_name. Available as of provider version R2.23+.
* `transit_prepend_as_path` - (Optional) Connection based AS Path Prepend. Valid only for BGP connection. Can only use the gateway's own local AS number, repeated up to 25 times. Applies on transit_gateway_name. Available as of provider version R2.23+.
* `transit_gateway_logical_ifnames` - (Optional) Transit gateway logical interface names for edge gateways where the peering terminates. Required only for edge as a transit attachment.
* `spoke_to_transit_filter_cidrs` - (Optional) Set of IPv4 or IPv6 CIDR prefixes to block from being advertised from the spoke gateway to the transit gateway. Changes are applied in place without recreating the attachment.
* `transit_to_spoke_filter_cidrs` - (Optional) (Optional) Set of IPv4 or IPv6 CIDR prefixes to block from being advertised from the transit gateway to the spoke gateway. Changes are applied in place without recreating the attachment.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `spoke_bgp_enabled` - Indicates whether the spoke gateway is BGP enabled or not.

## Import

**spoke_transit_attachment** can be imported using the `spoke_gw_name` and `transit_gw_name`, e.g.

```
$ terraform import aviatrix_spoke_transit_attachment.test spoke_gw_name~transit_gw_name
```
