---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_firenet_policy"
description: |-
  Creates and manages Aviatrix transit firenet policies
---

# aviatrix_transit_firenet_policy

The aviatrix_transit_firenet_policy resource allows the creation and management of Aviatrix transit firenet policies.

## Example Usage

```hcl
# Create an Aviatrix Transit FireNet Policy
resource "aviatrix_transit_firenet_policy" "test_transit_firenet_policy" {
  transit_firenet_gateway_name = "transitGw1"
  inspected_resource_name      = "SPOKE:spokeGw1"
}
```

## Argument Reference

The following arguments are supported:

* `transit_firenet_gateway_name` - (Required) The name of a transit gateway with transit firenet enabled. Currently supports AWS and AZURE providers.
* `inspected_resource_name` - (Required) The name of the resource which will be inspected.

## Import

Instance transit_firenet_policy can be imported using the transit_firenet_gateway_name and inspected_resource_name, e.g.

```
$ terraform import aviatrix_transit_firenet_policy.test transit_firenet_gateway_name~inspected_resource_name
```
