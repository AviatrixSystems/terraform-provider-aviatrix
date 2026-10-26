---
subcategory: "Peering"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_azure_peer"
description: |-
  Creates and manages of the Aviatrix peerings between Azure VNets
---

# aviatrix_azure_peer

The **aviatrix_azure_peer** resource allows the creation and management of the Aviatrix-created peerings between Azure VNets.

## Example Usage

```hcl
# Create an Aviatrix Azure Peering
resource "aviatrix_azure_peer" "test_azurepeer" {
  account_name1             = "test1-account"
  account_name2             = "test2-account"
  vnet_name_resource_group1 = "Foo_VNet1:Bar_RG1:GUID1"
  vnet_name_resource_group2 = "Foo_VNet2:Bar_RG2:GUID2"
  vnet_reg1                 = "Central US"
  vnet_reg2                 = "East US"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name1` - (Required) Name of the Azure cloud account in the Aviatrix controller for VNet 1.
* `account_name2` - (Required) Name of the Azure cloud account in the Aviatrix controller for VNet 2.

-> **NOTE:** As of Controller version 6.5+/provider version R2.20+, the `vnet_name_resource_group1` and `vnet_name_resource_group2` attributes should be in the format "vnet_name:rg_name:resource_guid".
* `vnet_name_resource_group1` - (Required) Azure VNet 1's name. Example: "VNet_Name1:Resource_Group_Name1:GUID1".
* `vnet_name_resource_group2` - (Required) Azure VNet 2's name. Example: "VNet_Name2:Resource_Group_Name2:GUID2".
* `vnet_reg1` - (Required) Region of Azure VNet 1. Example: "East US 2".
* `vnet_reg2` - (Required) Region of Azure VNet 2. Example: "East US 2".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vnet_cidr1` - List of VNet CIDR of vnet_name_resource_group1.
* `vnet_cidr2` - List of VNet CIDR of vnet_name_resource_group2.

## Import

**azure_peer** can be imported using the `vnet_name_resource_group1` and `vnet_name_resource_group2`, e.g.

```
$ terraform import aviatrix_azure_peer.test vnet_name_resource_group1~vnet_name_resource_group2
```
