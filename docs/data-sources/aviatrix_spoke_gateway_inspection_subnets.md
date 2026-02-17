---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway_inspection_subnets"
description: |-
  Gets all subnets available for the subnet inspection feature.
---

# aviatrix_spoke_gateway_inspection_subnets

The **aviatrix_spoke_gateway_inspection_subnets** data source provides all subnets available for the subnet inspection feature.

## Example Usage

```hcl
# Aviatrix Spoke Gateway Inspection Subnets Data Source
data "aviatrix_spoke_gateway_inspection_subnets" "foo" {
  gw_name = "gatewayname"
}
```

## Argument Reference

The following argument is supported:

* `gw_name` - (Required) Spoke gateway name.

## Attribute Reference

In addition to the argument above, the following attribute is exported:

* `subnets_for_inspection` - The list of all subnets available for the subnet inspection feature. This attribute is only supported for Azure.
