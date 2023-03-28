---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_access_allow_list_config"
description: |-
  Creates Aviatrix Controller Access Allow List Config
---

# aviatrix_controller_access_allow_list_config

!> **WARNING:** If any of the Allow List IPs are incorrect, the Controller will be inaccessible.

The **aviatrix_controller_access_allow_list_config** resource creates the Aviatrix Controller Access Allow List Config.

## Example Usage

```hcl
# Create a Controller Access Allow List Config
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
* `allow_list` - (Required) List of allowed IPs.
    * `ip_address` - (Required) IP address.
    * `description` - (Optional) Description.

### Optional
* `enable_enforce` - (Optional) Enable enforce. Valid values: true, false. Default: false.

## Import

**controller_access_allow_list_config** can be imported using "allow_list_config", e.g.

```
$ terraform import aviatrix_controller_access_allow_list_config.test allow_list_config
```
