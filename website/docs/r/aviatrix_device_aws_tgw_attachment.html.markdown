---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_aws_tgw_attachment"
description: |-
  Creates and manages a device and AWS TGW attachment
---

# aviatrix_device_aws_tgw_attachment

The **aviatrix_device_aws_tgw_attachment** resource allows the creation and management of a device and AWS TGW attachment

~> **NOTE:** Before creating this attachment the device must have its WAN interface and IP configured via the `aviatrix_device_interface_config` resource. To avoid attempting to create the attachment before the interface and IP are configured use a `depends_on` meta-argument so that the `aviatrix_device_interface_config` resource is created before the attachment.

## Example Usage

```hcl
# Create an Device and AWS TGW attachment
resource "aviatrix_device_aws_tgw_attachment" "test_device_aws_tgw_attachment" {
  connection_name      = "test-conn"
  device_name          = "device-a"
  aws_tgw_name         = "tgw-test"
  device_bgp_asn       = 65001
  security_domain_name = "Default_Domain"

  depends_on = [aviatrix_device_interface_config.test_device_interface_config]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `connection_name` - Connection name.
* `device_name` - Device name.
* `aws_tgw_name` - AWS TGW name.
* `device_bgp_asn` - BGP AS Number for the device.
* `security_domain_name` - Security Domain Name for the attachment.

## Import

**device_aws_tgw_attachment** can be imported using the `connection_name`, `device_name` and `aws_tgw_name`, e.g.

```
$ terraform import aviatrix_device_aws_tgw_attachment.test connection-name~device-name~aws-tgw-name
```
