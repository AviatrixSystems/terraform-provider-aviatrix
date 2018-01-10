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
