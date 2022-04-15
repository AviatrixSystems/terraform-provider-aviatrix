---
subcategory: "CloudN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_interface_config"
description: |-
  Configures WAN primary interface and IP for a device.
---

# aviatrix_device_interface_config

The **aviatrix_device_interface_config** resource allows the configuration of the WAN primary interface and IP for a device, for use in CloudN.

## Example Usage

```hcl
# Configure the primary WAN interface and IP for a device.
resource "aviatrix_device_interface_config" "test_device_interface_config" {
  device_name                     = "test-device"
  wan_primary_interface           = "eth0"
  wan_primary_interface_public_ip = "181.12.43.21"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `device_name` - (Required) Name of the device.
* `wan_primary_interface` - (Required) Name of the WAN primary interface.
* `wan_primary_interface_public_ip` - (Required) The WAN Primary interface public IP.

## Import

**device_interface_config** can be imported using the `device_name`, e.g.

```
$ terraform import aviatrix_device_interface_config.test device_name
```
