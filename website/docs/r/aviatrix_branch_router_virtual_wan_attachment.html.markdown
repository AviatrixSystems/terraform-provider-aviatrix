---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_branch_router_virtual_wan_attachment"
description: |-
  Creates and manages a branch router and Azure Virtual WAN attachment
---

# aviatrix_branch_router_virtual_wan_attachment

The **aviatrix_branch_router_virtual_wan_attachment** resource allows the creation and management of a branch router and Azure Virtual WAN attachment

## Example Usage

```hcl
# Create an Aviatrix Branch Router and Azure Virtual WAN attachment
resource "aviatrix_branch_router_virtual_wan_attachment" "test_branch_router_virtual_wan_attachment" {
  connection_name       = "test-conn"
  branch_name           = "branch-router"
  account_name          = "azure-devops"
  resource_group        = "aviatrix-rg"
  hub_name              = "aviatrix-hub"
  branch_router_bgp_asn = 65001
}
```

## Argument Reference

The following arguments are supported:

### Required
* `connection_name` - Connection name.
* `branch_name` - Branch router name.
* `account_name` - Azure access account name.
* `resource_group` - Azure Resource Manager resource group name.
* `hub_name` - Azure Virtual WAN vHub name.
* `branch_router_bgp_asn` - Branch Router AS Number. Integer between 1-4294967294.

## Import

**branch_router_virtual_wan_attachment** can be imported using the `connection_name`, e.g.

```
$ terraform import aviatrix_branch_router_virtual_wan_attachment.test connection_name
```
