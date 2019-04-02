---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_tunnel"
sidebar_current: "docs-aviatrix-resource-tunnel"
description: |-
  Creates and manages an Aviatrix tunnel.
---

# aviatrix_tunnel

The Account resource allows the creation and management of an Aviatrix tunnel.

## Example Usage

```hcl
# Create Aviatrix AWS tunnel
resource "aviatrix_tunnel" "test_tunnel1" {
  vpc_name1 = "avtxgw1"
  vpc_name2 = "avtxgw2"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_name1` - (Required) The first VPC Container name to make a peer pair
* `vpc_name2` - (Required) The second VPC Container name to make a peer pair
* `over_aws_peering` - (Optional) (Deprecated) Whether the peering should be done over AWS. Use aws_peer resource instead
* `cluster` - (Optional) Whether cluster peering is enabled ( Valid inputs: "yes", "no")
* `peering_link` - (Computed) Name of the peering link
* `enable_ha` - (Optional) Whether Peering HA is enabled ( Valid inputs: "yes", "no")


The following arguments are computed - please do not edit in the resource file:

* `peering_state` - (Computed) Status of the tunnel
* `peering_hastatus` - (Computed) Status of the HA tunnel

## Import

Instance tunnel can be imported using the vpc_name1 and vpc_name2, e.g.

```
$ terraform import aviatrix_tunnel.test vpc_name1~vpc_name2
```