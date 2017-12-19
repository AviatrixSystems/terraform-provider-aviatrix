---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transpeer"
sidebar_current: "docs-aviatrix-resource-transpeer"
description: |-
  Creates and manages an Aviatrix transitive peering.
---

# aviatrix_transpeer

The Account resource allows the creation and management of an Aviatrix transitive peering.

## Example Usage

```hcl
# Create Aviatrix AWS transitive peering.
resource "aviatrix_transpeer" "test_transpeer" {
  source = "avtxuseastgw1"
  nexthop = "avtxuseastgw2"
  reachable_cidr = "10.152.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) Name of Source gateway.
* `nexthop` - (Required) Name of nexthop gateway.
* `reachable_cidr` - (Required) Destination CIDR.
