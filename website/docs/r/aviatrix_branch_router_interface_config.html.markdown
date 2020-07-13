---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_branch_router_interface_config"
description: |-
  Configures primary WAN interface and IP for a branch router.
---

# aviatrix_branch_router_interface_config

The **aviatrix_branch_router_interface_config** resource allows the configuration of the primary WAN interface and IP for a branch router.

## Example Usage

```hcl
# Configure the primary WAN interface and IP for a branch router.
resource "aviatrix_branch_router_interface_config" "test_branch_router_interface_config" {
  branch_router_name              = "router-name"
  wan_primary_interface           = "GigabitEthernet1"
  wan_primary_interface_public_ip = "181.12.43.21"
}
```

## Argument Reference

The following arguments are supported:

* `branch_router_name` - (Required) Name of the branch router.
* `wan_primary_interface` - (Required) Name of the WAN Primary Interface.
* `wan_primary_interface_public_ip` - (Required) IP of the WAN Primary IP.

## Import

**branch_router_interface_config** can be imported using the `branch_router_name`, e.g.

```
$ terraform import aviatrix_branch_router_interface_config.test branch_router_name
```
