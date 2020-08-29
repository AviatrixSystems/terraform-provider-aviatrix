---
subcategory: "Peering"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_trans_peer"
description: |-
  Creates and manages Aviatrix transitive peerings
---

# aviatrix_trans_peer

The **aviatrix_trans_peer** resource allows the creation and management of Aviatrix [Encrypted Transitive Peering](https://docs.aviatrix.com/HowTos/TransPeering.html).

## Example Usage

```hcl
# Create an Aviatrix AWS Transitive Peering
resource "aviatrix_trans_peer" "test_trans_peer" {
  source         = "avtx-us-east-gw1"
  nexthop        = "avtx-us-east-gw2"
  reachable_cidr = "10.152.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `source` - (Required) Name of Source gateway.
* `nexthop` - (Required) Name of nexthop gateway.
* `reachable_cidr` - (Required) Destination CIDR.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `source_original_name` - Name of the source gateway when it was created.
* `nexthop_original_name` - Name of the nexthop gateway when it was created.

## Import

**trans_peer** can be imported using the `source`, `nexthop` and `reachable_cidr`, e.g.

```
$ terraform import aviatrix_trans_peer.test source~nexthop~reachable_cidr
```
