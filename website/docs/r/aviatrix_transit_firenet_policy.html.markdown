---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_firenet_policy"
description: |-
  Creates and manages Aviatrix Transit FireNet policies
---

# aviatrix_transit_firenet_policy

The **aviatrix_transit_firenet_policy** resource allows the creation and management of Aviatrix Transit FireNet policies that determine which resources should be inspected in the Transit FireNet solution.

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

### Required
* `transit_firenet_gateway_name` - (Required) Name of the Transit FireNet-enabled transit gateway. Currently supports AWS and Azure.
* `inspected_resource_name` - (Required) The name of the resource which will be inspected.

## Import

**transit_firenet_policy** can be imported using the `transit_firenet_gateway_name` and `inspected_resource_name`, e.g.

```
$ terraform import aviatrix_transit_firenet_policy.test transit_firenet_gateway_name~inspected_resource_name
```
