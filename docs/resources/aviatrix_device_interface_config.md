---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_interface_config"
description: |-
  Configures primary WAN interface and IP for a device.
---

# aviatrix_device_interface_config

The **aviatrix_device_interface_config** resource allows the configuration of the primary WAN interface and IP for a device, for use in CloudWAN.

~> **NOTE:** Before configuring WAN interface and IP, the device must be registered with the Aviatrix controller via the `aviatrix_device_registration` resource. To guarantee the correct order of resource creation please set an explicit or implicit dependency on the corresponding `aviatrix_device_registration` resource. For an example of an implicit dependency please see the Example Usage section below. For explicit dependency please utilize a `depends_on` meta-argument within this resource.

## Example Usage

```hcl
# Configure the primary WAN interface and IP for a device.
resource "aviatrix_device_interface_config" "test_device_interface_config" {
  # Set an implicit dependency on the aviatrix_device_registration resource by referencing the
  # device name from that resource.
  device_name                     = aviatrix_device_registration.test_device_registration.name
  wan_primary_interface           = "GigabitEthernet1"
  wan_primary_interface_public_ip = "181.12.43.21"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `device_name` - (Required) Name of the device.
* `wan_primary_interface` - (Required) Name of the WAN Primary Interface.
* `wan_primary_interface_public_ip` - (Required) IP of the WAN Primary IP.

## Import

**device_interface_config** can be imported using the `device_name`, e.g.

```
$ terraform import aviatrix_device_interface_config.test device_name
```
