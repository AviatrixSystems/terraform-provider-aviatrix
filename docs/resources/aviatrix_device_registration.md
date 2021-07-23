---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_registration"
description: |-
  Creates and manages device registration for CloudWAN
---

# aviatrix_device_registration

The **aviatrix_device_registration** resource allows the registration and management of devices for use in CloudWAN.

~> **NOTE:** Before this device can be attached to any Aviatrix Transit Gateway, AWS TGW or Azure Virtual WAN you must configure its WAN interface and IP via the `aviatrix_device_interface_config` resource.

## Example Usage

```hcl
# Register a device with private key authentication
resource "aviatrix_device_registration" "test_device" {
  name      = "test-device"
  public_ip = "58.151.114.231"
  username  = "ec2-user"
  key_file  = "/path/to/key_file.pem"
}
```

```hcl
# Register a device with password authentication
resource "aviatrix_device_registration" "test_device" {
  name      = "test-device"
  public_ip = "58.151.114.231"
  username  = "ec2-user"
  password  = "secret"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of the device.
* `public_ip` - (Required) Public IP address of the device.
* `username` - (Required) Username for SSH into the device.
* `key_file` - (Optional) Path to private key file for SSH into the device. Either `key_file` or `password` must be set to register a device successfully.
* `password` - (Optional) Password for SSH into the router. Either `key_file` or `password` must be set to register a device successfully. This attribute can also be set via environment variable 'AVIATRIX_DEVICE_PASSWORD'. If both are set, the value in the config file will be used.

### Optional
* `host_os` - (Optional) Device host OS. Default value is 'ios'. Valid values are 'ios' or 'aviatrix'.
* `ssh_port` - (Optional) SSH port for connecting to the device. Default value is 22.
* `address_1` - (Optional) Address line 1.
* `address_2` - (Optional) Address line 2.
* `city` - (Optional) City.
* `state` - (Optional) State.
* `country` - (Optional) ISO two-letter country code.
* `zip_code` - (Optional) Zip code.
* `description` - (Optional) Description.

### Managed CloudN (CaaG) Upgrade
* `software_version` - (Optional/Computed) The desired software version of the CaaG. If set, we will attempt to update the CaaG to the specified version. If left blank, the software version will continue to be managed through the aviatrix_controller_config resource. Type: String. Example: "6.5.892". Available as of provider version R2.20.0.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `is_caag` - Is this device a Managed CloudN (CaaG). Type: Boolean. Available as of provider version R2.20.0.

## Import

**device_registration** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_device_registration.test name
```
