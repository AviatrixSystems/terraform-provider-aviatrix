---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_tunnel"
description: |-
  Creates and manages Aviatrix Tunnels.
---

# aviatrix_tunnel

The aviatrix_tunnel resource allows the creation and management of Aviatrix tunnels.

## Example Usage

```hcl
# Create an Aviatrix AWS Tunnel
resource "aviatrix_tunnel" "test_tunnel1" {
  gw_name1 = "avtxgw1"
  gw_name2 = "avtxgw2"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name1` - (Required) The first VPC Container name to make a peer pair.
* `gw_name2` - (Required) The second VPC Container name to make a peer pair.
* `enable_ha` - (Optional) Whether Peering HA is enabled. Valid values: true, false.

The following arguments are computed - please do not edit in the resource file:

* `peering_state` - (Computed) Status of the tunnel.
* `peering_hastatus` - (Computed) Status of the HA tunnel.
* `peering_link` - (Computed) Name of the peering link.

## Import

Instance tunnel can be imported using the gw_name1 and gw_name2, e.g.

```
$ terraform import aviatrix_tunnel.test gw_name1~gw_name2
```