---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_config"
description: |-
  Creates and manages an Aviatrix controller config resource
---

# aviatrix_controller_config

The **aviatrix_controller_config** resource allows management of an Aviatrix Controller's configurations.

## Example Usage

```hcl
# Create an Aviatrix Controller Config
resource "aviatrix_controller_config" "test_controller_config" {
  fqdn_exception_rule = false
}
```
```hcl
# Create an Aviatrix Controller Config with Controller Upgrade Without Upgrading Gateways
resource "aviatrix_controller_config" "test_controller_config" {
  fqdn_exception_rule     = false
  target_version          = "latest"
  manage_gateway_upgrades = false
}
```
```hcl
# Create an Aviatrix Controller Config with Controller Upgrade + Upgrade All Gateways
resource "aviatrix_controller_config" "test_controller_config" {
  fqdn_exception_rule = false
  target_version      = "latest"
}
```
```hcl
# Create an Aviatrix Controller Config with Cloudn Backup Configuration Enabled
resource "aviatrix_controller_config" "test_controller_config" {
  backup_configuration = true
  backup_cloud_type    = 1
  backup_account_name  = "account_example"
  backup_bucket_name   = "bucket_example"
}
```
```hcl
# Create an Aviatrix Controller Config and import HTTPS certificates
resource "aviatrix_controller_config" "test_controller_config" {
  ca_certificate_file_path            = "/path/to/ca_certificate.pem"
  server_private_key_file_path        = "/path/to/server.key"
  server_public_certificate_file_path = "/path/to/server.crt"
}
```
```hcl
# Create an Aviatrix Controller Config and configure the AWS Guard Duty Scanning Interval
resource "aviatrix_controller_config" "test_controller_config" {
  aws_guard_duty_scanning_interval = 10
}
```


## Argument Reference

The following arguments are supported:

### Controller and Gateway Upgrade

-> **NOTE:** To selectively upgrade your gateways, you MUST set `manage_gateway_upgrades` to false. Gateway upgrades can then be managed via the software_version and image_version attributes of the gateway resources. If you do not wish to selectively upgrade gateways, `manage_gateway_upgrades` can be left as the default true value.

* `target_version` - (Optional) The release version number to which the controller will be upgraded to. If not specified, controller will not be upgraded. If set to "latest", controller will be upgraded to the latest release. Please see the [Controller upgrade guide](https://docs.aviatrix.com/HowTos/inline_upgrade.html) for more information.
* `manage_gateway_upgrades` - (Optional) If true, aviatrix_controller_config will upgrade all gateways when target_version is set. If false, only the controller will be upgraded when target_version is set. In that case gateway upgrades should be handled in each gateway resource individually using the software_version and image_version attributes. Type: boolean. Default: true. Available as of provider version R2.20.0+.

### Security Options
* `fqdn_exception_rule` - (Optional) Enable/disable packets without an SNI field to pass through gateway(s). Valid values: true, false. Default value: true. For more information on this setting, please see [here](https://docs.aviatrix.com/HowTos/FQDN_Whitelists_Ref_Design.html#exception-rule)
* `aws_guard_duty_scanning_interval` - (Optional) Configure the AWS Guard Duty scanning interval. Valid values: 5, 10, 15, 30 or 60. Default value: 60. Available as of provider version R2.18+.

### Backup
* `backup_configuration` - (Optional) Switch to enable/disable controller CloudN backup config. Valid values: true, false. Default value: false.
* `backup_cloud_type` - (Optional) Type of cloud service provider, requires an integer value. Use 1 for AWS, 4 for GCP, 8 for Azure, 16 for OCI, and 256 for AWSGov.
* `backup_account_name` - (Optional) Name of the cloud account in the Aviatrix controller.
* `backup_bucket_name` - (Optional) Bucket Name. Required to enable configuration backup for AWS, AWSGov, GCP and OCI.
* `backup_storage_name` - (Optional) Storage name. Required to enable configuration backup for Azure.
* `backup_container_name` - (Optional) Container name. Required to enable configuration backup for Azure.
* `backup_region` - (Optional) Name of region. Required to enable configuration backup for Azure and OCI.
* `multiple_backups` - (Optional) Switch to enable the Controller to backup up to a maximum of 3 rotating backups. Valid values: true, false. Default value: false.

-> **NOTE:** `backup_bucket_name` is required for AWS, AWSGov and GCP. `backup_storage_name`, `backup_container_name` and `backup_region` are required for Azure. `backup_bucket_name` and `backup_region` are required for OCI.

### [TLS Certificate Import](https://docs.aviatrix.com/HowTos/controller_certificate.html)

~> **NOTE:** Please use either the combination of `ca_certificate_file_path`, `server_public_certificate_file_path` and `server_private_key_file_path` or the combination of `ca_certificate_file`, `server_public_certificate_file` and `server_private_key_file`.

* `ca_certificate_file_path` - (Optional) File path to CA certificate. Available as of provider version R2.18+.
* `server_public_certificate_file_path` - (Optional) File path to the server public certificate. Available as of provider version R2.18+.
* `server_private_key_file_path` - (Optional) File path to server private key. Available as of provider version R2.18+.
* `ca_certificate_file` - (Optional) CA certificate. To read certificate file from a file, please use the built-in `file` function. Available as of provider version R2.21.2+.
* `server_public_certificate_file` - (Optional) Server public certificate. To read certificate file from a file, please use the built-in `file` function. Available as of provider version R2.21.2+.
* `server_private_key_file` - (Optional) Server private key. To read the private key from a file, please use the built-in `file` function. Available as of provider version R2.21.2+.

### Misc.
* `enable_vpc_dns_server` - (Optional) Enable VPC/VNET DNS Server for the controller. Valid values: true, false. Default value: false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the controller without build number. Example: "6.5"
* `previous_version` - Previous version of the controller including the build number. Example: "6.5.123". Available as of provider version R2.20.0+.
* `current_version` - Current version of the controller including the build number. Example: "6.5.123". Available as of provider version R2.20.0+.

~> **NOTE:** The following attributes are deprecated and removed. Please use **aviatrix_controller_security_group_management_config** resource to manage controller's security group management settings.

* `sg_management_account_name` - (Optional) Select the [primary access account](https://docs.aviatrix.com/HowTos/aviatrix_account.html#setup-primary-access-account-for-aws-cloud).
* `security_group_management` - (Optional) Enable to allow Controller to automatically manage inbound rules from gateways. Valid values: true, false. Default value: false.


## Import

Instance controller_config can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_config.test 10-11-12-13
```
