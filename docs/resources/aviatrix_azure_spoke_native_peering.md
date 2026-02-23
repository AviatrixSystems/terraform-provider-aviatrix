---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_azure_spoke_native_peering"
description: |-
  Creates and manages Aviatrix Azure spoke native peerings
---

# aviatrix_azure_spoke_native_peering

The **aviatrix_azure_spoke_native_peering** resource allows the creation and management of Aviatrix-created Azure Spoke VNet attachments via Native Peering.

## Example Usage

```hcl
# Create an Aviatrix Azure spoke native peering
resource "aviatrix_azure_spoke_native_peering" "test" {
  transit_gateway_name = "transit-gw-azure"
  spoke_account_name   = "devops-azure"
  spoke_region         = "West US"
  spoke_vpc_id         = "Foo_VNet:Bar_RG:GUID"
}
```
```hcl
# Create an Aviatrix Azure spoke native peering with private route table configured
resource "aviatrix_azure_spoke_native_peering" "test" {
  transit_gateway_name = "transit-gw-azure"
  spoke_account_name   = "devops-azure"
  spoke_region         = "West US"
  spoke_vpc_id         = "Foo_VNet:Bar_RG:GUID"
  private_route_table_config            = [
    "Foo_VNet_RTB_1:Bar_RG",
    "Foo_VNet_RTB_2:Bar_RG",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `transit_gateway_name` - (Required) Name of an Transit FireNet-enabled Azure transit gateway.
* `spoke_account_name` - (Required) An Aviatrix account that corresponds to a subscription in Azure.
* `spoke_region` - (Required) Spoke VNet region. Example: "West US".

-> **NOTE:** As of Controller version 6.5+/provider version R2.20+, the `spoke_vpc_id` for Spoke Gateways in Azure should be in the format "vnet_name:rg_name:resource_guid".
* `spoke_vpc_id` - (Required) Combination of the Spoke's VNet name, resource group and GUID. Example: "Foo_VNet:Bar_RG:GUID".

### Optional
* `private_route_table_config` - (Optional) Set of Azure route table selectors to treat as private route tables for the spoke VNet. Each entry in the list is in the format of "<route_table_name>:<resource_group_name>" (for example: "Foo_VNet_RTB_1:Bar_RG").

## Import

**azure_spoke_native_peering** can be imported using the `transit_gateway_name`, `spoke_account_name` and `spoke_vpc_id`, e.g.

```
$ terraform import aviatrix_azure_spoke_native_peering.test transit_gateway_name~spoke_account_name~spoke_vpc_id
```
