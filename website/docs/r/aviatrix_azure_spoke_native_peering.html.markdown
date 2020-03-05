---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_azure_spoke_native_peering"
description: |-
  Creates and manages Aviatrix Azure spoke native peerings
---

# aviatrix_azure_spoke_native_peering

The aviatrix_azure_spoke_native_peering resource allows the creation and management of Aviatrix Azure spoke native peerings.

## Example Usage

```hcl
# Create an Aviatrix Azure spoke native peering
resource "aviatrix_azure_spoke_native_peering" "test" {
  transit_gateway_name = "transit-gw-azure"
  spoke_account_name   = "devops-azure"
  spoke_region         = "West US"
  spoke_vpc_id         = "vnet:hello"
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_name` - (Required) Name of an azure transit gateway with transit firenet enabled.
* `spoke_account_name` - (Required) An Aviatrix account that corresponds to a subscription in Azure.
* `spoke_region` - (Required) Spoke VNet region. Example: "West US".
* `spoke_vpc_id` - (Required) Combination of the Spoke VNet name and resource group. Example: "vnet:hello".

## Import

Instance azure_spoke_native_peering can be imported using the transit_gateway_name, spoke_account_name and spoke_vpc_id, e.g.

```
$ terraform import aviatrix_azure_spoke_native_peering.test transit_gateway_name~spoke_account_name~spoke_vpc_id
```
