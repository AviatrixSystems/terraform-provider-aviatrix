---
subcategory: "Private Mode"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_private_mode_config"
description: |-
  Creates and manages an Aviatrix Controller's Private Mode config
---

# aviatrix_controller_private_mode_config

The **aviatrix_controller_private_mode_config** resource allows management of an Aviatrix Controller's Private Mode configuration. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Create an Aviatrix Controller Private Mode config
resource "aviatrix_controller_private_mode_config" "test" {
  enable_private_mode = true
}
```


## Argument Reference

The following arguments are supported:

### Required
* `enable_private_mode` - (Required) Whether to enable Private Mode on an Aviatrix Controller.

### Optional
* `copilot_instance_id` - (Optional) Instance ID of a copilot instance to associate with an Aviatrix Controller in Private Mode. The copilot instance must be in the same VPC as the Aviatrix Controller.
* `proxies` - (Optional) Set of Controller proxies for Private Mode.


## Import

**aviatrix_controller_private_mode_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_private_mode_config.test 10-11-12-13
```
