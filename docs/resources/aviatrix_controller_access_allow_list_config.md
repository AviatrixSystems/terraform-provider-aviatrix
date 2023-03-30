---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_access_allow_list_config"
description: |-
  Creates and manages an Aviatrix Controller's Access Allow List Config
---

# aviatrix_controller_access_allow_list_config

!> **WARNING:** If any of the Allow List IPs are incorrect, the Controller will be inaccessible.

The **aviatrix_controller_access_allow_list_config** resource allows management of an Aviatrix Controller's Access Allow List.

## Example Usage

```hcl
# Create an Aviatrix Controller Access Allow List Config
resource "aviatrix_controller_access_allow_list_config" "test" {
  allow_list {
    ip_address = "0.0.0.0"
  }

  allow_list {
    ip_address = "1.2.3.4"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `allow_list` - (Required) Set of allowed IPs which have the access to the controller.
  * `ip_address` - (Required) IP address which has the access to the controller.
  * `description` - (Optional) Description of the IP address.

### Optional
* `enable_enforce` - (Optional) Specify whether to enable enforce or not. Valid values: true, false. Default: false.

## Import

**controller_access_allow_list_config** can be imported using controller IP, e.g.

```
$ terraform import aviatrix_controller_access_allow_list_config.test 10-11-12-13
```
