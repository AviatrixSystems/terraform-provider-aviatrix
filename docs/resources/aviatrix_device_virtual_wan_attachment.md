---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_virtual_wan_attachment"
description: |-
  Creates and manages a device and Azure Virtual WAN attachment
---

# aviatrix_device_virtual_wan_attachment

The **aviatrix_device_virtual_wan_attachment** resource allows the creation and management of a device and Azure Virtual WAN attachment for use in CloudWAN.

~> **NOTE:** Before creating this attachment the device must have its WAN interface and IP configured via the `aviatrix_device_interface_config` resource. To avoid attempting to create the attachment before the interface and IP are configured use a `depends_on` meta-argument so that the `aviatrix_device_interface_config` resource is created before the attachment.

## Example Usage

```hcl
# Create an Device and Azure Virtual WAN attachment
resource "aviatrix_device_virtual_wan_attachment" "test_device_virtual_wan_attachment" {
  connection_name = "test-conn"
  device_name     = "device-a"
  account_name    = "azure-devops"
  resource_group  = "aviatrix-rg"
  hub_name        = "aviatrix-hub"
  device_bgp_asn  = 65001

  depends_on = [aviatrix_device_interface_config.test_device_interface_config]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `connection_name` - Connection name.
* `device_name` - Device name.
* `account_name` - Azure access account name.
* `resource_group` - Azure Resource Manager resource group name.
* `hub_name` - Azure Virtual WAN vHub name.
* `device_bgp_asn` - Device AS Number. Integer between 1-4294967294.


## Import

**device_virtual_wan_attachment** can be imported using the `connection_name`, e.g.

```
$ terraform import aviatrix_device_virtual_wan_attachment.test connection_name
```
