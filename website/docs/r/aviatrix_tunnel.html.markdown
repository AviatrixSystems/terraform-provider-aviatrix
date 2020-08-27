---
subcategory: "Peering"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_tunnel"
description: |-
  Creates and manages Aviatrix Encrypted Peering tunnels
---

# aviatrix_tunnel

The **aviatrix_tunnel** resource allows the creation and management of Aviatrix Encrypted Peering tunnels.

## Example Usage

```hcl
# Create an Aviatrix AWS Tunnel
resource "aviatrix_tunnel" "test_tunnel" {
  gw_name1 = "avtx-gw1"
  gw_name2 = "avtx-gw2"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `gw_name1` - (Required) The first VPC Container name to make a peer pair.
* `gw_name2` - (Required) The second VPC Container name to make a peer pair.

### HA
* `enable_ha` - (Optional) Enable this attribute if peering-HA is enabled on the gateways. Valid values: true, false. Default value: false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `peering_state` - (Computed) Status of the tunnel.
* `peering_hastatus` - (Computed) Status of the HA tunnel.
* `peering_link` - (Computed) Name of the peering link.
* `gw_original_name1` - Name of the first VPC Container when it was created.
* `gw_original_name2` - Name of the second VPC Container when it was created.

## Import

**tunnel** can be imported using the `gw_name1` and `gw_name2`, e.g.

```
$ terraform import aviatrix_tunnel.test gw_name1~gw_name2
```
