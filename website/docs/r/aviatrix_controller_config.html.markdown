---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_config"
sidebar_current: "docs-aviatrix-resource-controller-config"
description: |-
  Creates and manages an Aviatrix controller config resource.
---

# aviatrix_controller_config

The Account resource allows the creation and management of an Aviatrix controller config.

## Example Usage

```hcl
# Create Aviatrix Controller Config
resource "aviatrix_controller_config" "test_controller_config" {
  sg_management_account_name = "username"
  http_access                = true
  fqdn_exception_rule        = false
  security_group_management  = true
}
```

## Argument Reference

The following arguments are supported:

* `sg_management_account_name` - (Optional) Cloud account name of user.
* `http_access` - (Optional) Switch for http access. Default: false.
* `fqdn_exception_rule` - (Optional) A system-wide mode. Default: true.
* `security_group_management` - (Optional) Used to manage the Controller instance’s inbound rules from gateways. Default: false.

## Import

Instance controller_config can be imported using ControllerIP, e.g.

```
$ terraform import aviatrix_controller_config.test ControllerIP
```