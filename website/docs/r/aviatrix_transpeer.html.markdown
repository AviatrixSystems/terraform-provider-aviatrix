---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transpeer"
sidebar_current: "docs-aviatrix-resource-transpeer"
description: |-
  Creates and manages an Aviatrix transitive peering.
---

# aviatrix_trans_peer

The Account resource allows the creation and management of an Aviatrix transitive peering.

## Example Usage

```hcl
# Create Aviatrix AWS transitive peering.
resource "aviatrix_trans_peer" "test_transpeer" {
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