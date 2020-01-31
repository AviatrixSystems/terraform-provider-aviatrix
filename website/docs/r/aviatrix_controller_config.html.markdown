---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_config"
description: |-
  Creates and manages an Aviatrix controller config resource
---

# aviatrix_controller_config

The aviatrix_controller_config resource allows management of an Aviatrix Controller's configurations.

## Example Usage

```hcl
# Create an Aviatrix Controller Config
resource "aviatrix_controller_config" "test_controller_config" {
  sg_management_account_name = "username"
  http_access                = true
  fqdn_exception_rule        = false
  security_group_management  = true
}
```
```hcl
# Create an Aviatrix Controller Config with Controller Upgrade
resource "aviatrix_controller_config" "test_controller_config" {
  sg_management_account_name = "username"
  http_access                = true
  fqdn_exception_rule        = false
  security_group_management  = true
  target_version             = "latest"
}
```

## Argument Reference

The following arguments are supported:

* `sg_management_account_name` - (Optional) Cloud account name of user.
* `http_access` - (Optional) Switch for http access. Valid values: true, false. Default value: false.
* `fqdn_exception_rule` - (Optional) A system-wide mode. Valida values: true, false. Defaultvalue: true.
* `security_group_management` - (Optional) Used to manage the Controller instance’s inbound rules from gateways. Valid values: true, false. Default value: false.
* `target_version` - (Optional) The release version number to which the controller will be upgraded to. If not specified, controller will not be upgraded. If set to "latest", controller will be upgraded to the latest release. Please look at https://docs.aviatrix.com/HowTos/inline_upgrade.html for more information.
* `backup_config_enabled` - (Optional) Switch to enable/disable controller cloudn backup config. Valid values: true, false. Default value: false.
* `cloud_type` - (Optional) Type of cloud service provider, requires an integer value. Use 1 for AWS.
* `account_name` - (Optional) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `bucket_name` - (Optional) S3 Bucket Name for AWS.
* `multiple` - (Optional) Switch to enable the controller to backup up to a maximum of 3 rotating backups.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the controller.

## Import

Instance controller_config can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_config.test 10-11-12-13
```
