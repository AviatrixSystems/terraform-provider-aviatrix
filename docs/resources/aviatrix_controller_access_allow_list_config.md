---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_access_allow_list_config"
description: |-
  Creates and manages an Aviatrix Controller's Access Allow List Config
---

# aviatrix_controller_access_allow_list_config

!> **WARNING:** If any of the IPs in `allow_list {}` are incorrect, the Controller will be inaccessible.

The **aviatrix_controller_access_allow_list_config** resource enables configuration of a set of IP addresses allowed HTTP(s) access to an Aviatrix Controller.

## Example Usage

```hcl
# Create an Aviatrix Controller Access Allow List Config
resource "aviatrix_controller_access_allow_list_config" "test" {
  allow_list {
    ip_address  = "<< allowed IP address 1 >>"
    description = "allowed IP address 1"
  }

  allow_list {
    ip_address  = "<< allowed IP address 2 >>"
    description = "allowed IP address 2"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `allow_list` - (Required) Set of IPs allowed access to the Controller.
  * `ip_address` - (Required) IP address allowed access to the Controller.
  * `description` - (Optional) Description of the IP address.

### Optional
* `enable_enforce` - (Optional) Set to true to enable enforcement of the `allow_list {}`'s IPs. Valid values: true, false. Default: false.

## Import

**controller_access_allow_list_config** can be imported using controller IP, e.g.

```
$ terraform import aviatrix_controller_access_allow_list_config.test 10-11-12-13
```
