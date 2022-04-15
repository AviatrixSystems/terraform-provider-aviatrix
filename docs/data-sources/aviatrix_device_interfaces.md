---
subcategory: "CloudN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_interfaces"
description: |-
  Gets the WAN primary interfaces and IPs for a device.
---

# aviatrix_device_interfaces

Use this data source to get the list of device WAN interfaces for use in other resources.

## Example Usage

```hcl
# Aviatrix Device Interfaces Data Source
data "aviatrix_device_interfaces" "test" {
  device_name = "test-device"
}
```

## Argument Reference

The following argument is supported:

* `device_name` - (Required) Device name.

## Attribute Reference

In addition to the argument above, the following attributes are exported:

* `wan_interfaces` - List of WAN interfaces.
  * `wan_primary_interface` - Name of the WAN primary interface.
  * `wan_primary_interface_public_ip` - The WAN Primary interface public IP.
