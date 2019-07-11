---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_trans_peer"
sidebar_current: "docs-aviatrix-resource-trans-peer"
description: |-
  Creates and manages Aviatrix Transitive Peerings
---

# aviatrix_trans_peer

The Account resource allows the creation and management of Aviatrix Transitive Peerings.

## Example Usage

```hcl
# Create an Aviatrix AWS transitive peering
resource "aviatrix_trans_peer" "test_trans_peer" {
  source         = "avtxuseastgw1"
  nexthop        = "avtxuseastgw2"
  reachable_cidr = "10.152.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) Name of Source gateway.
* `nexthop` - (Required) Name of nexthop gateway.
* `reachable_cidr` - (Required) Destination CIDR.

## Import

Instance trans_peer can be imported using the source, nexthop and reachable_cidr, e.g.

```
$ terraform import aviatrix_trans_peer.test source~nexthop~reachable_cidr
```
