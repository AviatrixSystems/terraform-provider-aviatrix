---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_security_group_management_config"
description: |-
  Creates and manages an Aviatrix controller security group management config resource
---

# aviatrix_controller_security_group_management_config

The **aviatrix_controller_security_group_management_config** resource allows management of an Aviatrix Controller's security group management configurations. This resource is available as of v2.20.1.

## Example Usage

```hcl
# Create an Aviatrix Controller Security Group Management Config to Enable Security Group Management
resource "aviatrix_controller_security_group_management_config" "test_sqm_config" {
  account_name                     = "devops"
  enable_security_group_management = true
}
```
```hcl
# Create an Aviatrix Controller Security Group Management Config to Disable Security Group Management
resource "aviatrix_controller_security_group_management_config" "test_sqm_config" {
  enable_security_group_management = false
}
```


## Argument Reference

The following arguments are supported:

* `account_name` - (Optional) Select the [primary access account](https://docs.aviatrix.com/HowTos/aviatrix_account.html#setup-primary-access-account-for-aws-cloud).
* `enable_security_group_management` - (Required) Enable to allow Controller to automatically manage inbound rules from gateways. Valid values: true, false.


## Import

Instance controller_security_group_management_config can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_security_group_management_config.test 10-11-12-13
```
