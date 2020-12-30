---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_private_oob"
description: |-
  Creates and manages an Aviatrix controller private OOB config resource
---

# aviatrix_controller_private_oob

The **aviatrix_controller_private_oob** resource allows management of an Aviatrix Controller's private OOB configurations.

## Example Usage

```hcl
# Create an Aviatrix Controller Private OOB
resource "aviatrix_controller_private_oob" "test_private_oob" {
  enable_private_oob = true
}
```


## Argument Reference

The following arguments are supported:

* `enable_private_oob` - (Optional) Switch to enable/disable Aviatrix controller private OOB. Valid values: true, false. Default value: false.

## Import

**controller_private_oob** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_private_oob.test 10-11-12-13
```
