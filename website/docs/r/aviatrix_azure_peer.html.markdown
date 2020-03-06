---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_azure_peer"
description: |-
  Creates and manages Aviatrix Azure peerings
---

# aviatrix_azure_peer

The aviatrix_azure_peer resource allows the creation and management of Aviatrix Azure peerings.

## Example Usage

```hcl
# Create an Aviatrix Azure Peering
resource "aviatrix_azure_peer" "test_azurepeer" {
  account_name1             = "test1-account"
  account_name2             = "test2-account"
  vnet_name_resource_group1 = "Foo_VNet1:Bar_RG1"
  vnet_name_resource_group2 = "Foo_VNet2:Bar_RG2"
  vnet_reg1                 = "Central US"
  vnet_reg2                 = "East US"
}
```

## Argument Reference

The following arguments are supported:

* `account_name1` - (Required) This parameter represents the name of an Azure Cloud-Account in Aviatrix controller.
* `account_name2` - (Required) This parameter represents the name of an Azure Cloud-Account in Aviatrix controller.
* `vnet_name_resource_group1` - (Required) VNet-Name of Azure cloud. Example: "VNet_Name:Resource_Group_Name".
* `vnet_name_resource_group2` - (Required) VNet-Name of Azure cloud. Example: "VNet_Name:Resource_Group_Name".
* `vnet_reg1` - (Required) Region of Azure cloud. Example: "East US 2".
* `vnet_reg2` - (Required) Region of Azure cloud. Example: "East US 2".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vnet_cidr1` - List of VNet CIDR of vnet_name_resource_group1.
* `vnet_cidr2` - List of VNet CIDR of vnet_name_resource_group2.

## Import

Instance azure_peer can be imported using the vnet_name_resource_group1 and vnet_name_resource_group2, e.g.

```
$ terraform import aviatrix_azure_peer.test vnet_name_resource_group1~vnet_name_resource_group2
```
