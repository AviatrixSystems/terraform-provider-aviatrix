---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_interface_config"
description: |-
  Configures primary WAN interface and IP for a device.
---

# aviatrix_device_interface_config

The **aviatrix_device_interface_config** resource allows the configuration of the primary WAN interface and IP for a device.

## Example Usage

```hcl
# Configure the primary WAN interface and IP for a device.
resource "aviatrix_device_interface_config" "test_device_interface_config" {
  device_name                     = "device-name"
  wan_primary_interface           = "GigabitEthernet1"
  wan_primary_interface_public_ip = "181.12.43.21"
}
```

## Argument Reference

The following arguments are supported:

* `device_name` - (Required) Name of the device.
* `wan_primary_interface` - (Required) Name of the WAN Primary Interface.
* `wan_primary_interface_public_ip` - (Required) IP of the WAN Primary IP.

## Import

**device_interface_config** can be imported using the `device_name`, e.g.

```
$ terraform import aviatrix_device_interface_config.test device_name
```
