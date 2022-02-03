---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway_subnet_group"
description: |-
  Creates and manages Aviatrix spoke gateway subnet groups
---

# aviatrix_spoke_gateway_subnet_group

The **aviatrix_spoke_gateway_subnet_group** resource creates and manages the spoke gateway subnet groups.

-> **NOTE:** This feature is only valid for Azure.

## Example Usage

```hcl
# Create an Aviatrix Spoke Inspection Subnet Group
resource "aviatrix_spoke_gateway_subnet_group" "test" {
  name               = "subnet_group"
  spoke_gateway_name = "spoke"
  subnets            = ["10.2.48.0/20~~subnet1", "10.2.64.0/20~~subnet2"]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of spoke gateway subnet group.
* `spoke_gateway_name` - (Required) Aviatrix spoke gateway name.
* `subnets` - (Optional) A set of subnets in the subnet group. The format of each subnet must be "CIDR~~subnet name". Example: `["10.2.48.0/20~~subnet1", "10.2.64.0/20~~subnet2"]`

## Import

**spoke_gateway_subnet_group** can be imported using the `spoke_gateway_name` and `name`, e.g.

```
$ terraform import aviatrix_spoke_gateway_subnet_group.test spoke_gateway_name~name
```
