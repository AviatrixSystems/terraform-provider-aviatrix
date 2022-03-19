---
subcategory: "Deprecated"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_transit_gateway_attachment"
description: |-
  Creates and manages a device and Aviatrix Transit Gateway attachment
---

# aviatrix_device_transit_gateway_attachment

The **aviatrix_device_transit_gateway_attachment** resource allows the creation and management of a device and Aviatrix Transit Gateway attachment for use in CloudWAN.

~> **NOTE:** Before creating this attachment the device must have its WAN interface and IP configured via the `aviatrix_device_interface_config` resource. To avoid attempting to create the attachment before the interface and IP are configured use a `depends_on` meta-argument so that the `aviatrix_device_interface_config` resource is created before the attachment.  

## Example Usage

```hcl
# Create an Device and Transit Gateway attachment
resource "aviatrix_device_transit_gateway_attachment" "test_device_transit_gateway_attachment" {
  device_name             = "device-a"
  transit_gateway_name    = "transit-gw"
  connection_name         = "test-conn"
  transit_gateway_bgp_asn = 65000
  device_bgp_asn          = 65001

  depends_on              = [
    aviatrix_device_interface_config.test_device_interface_config
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `device_name` - Device name.
* `transit_gateway_name` - Aviatrix Transit Gateway name.
* `connection_name` - Connection name.
* `transit_gateway_bgp_asn` - BGP AS Number for transit gateway.
* `device_bgp_asn` - BGP AS Number for the device.

### Optional
* `phase1_authentication` - Phase 1 authentication algorithm. Default "SHA-256".
* `phase1_dh_groups` - Number of phase 1 Diffie-Hellman groups. Default "14".
* `phase1_encryption` - Phase 1 encryption algorithm. Default "AES-256-CBC".
* `phase2_authentication` - Phase 2 authentication algorithm. Default "HMAC-SHA-256".
* `phase2_dh_groups` - Number of phase 2 Diffie-Hellman groups. Default "14".
* `phase2_encryption` - Phase 2 encryption algorithm. Default "AES-256-CBC".
* `enable_global_accelerator` - Boolean enable AWS Global Accelerator. Default "false".
* `pre_shared_key` - Pre-shared Key.
* `local_tunnel_ip` - Local tunnel IP.
* `remote_tunnel_ip` - Remote tunnel IP.
* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval for the connection. Requires the transit_gateway's 'learned_cidrs_approval_mode' attribute be set to 'connection'. Valid values: true, false. Default value: false. Available as of provider version R2.18+.
* `manual_bgp_advertised_cidrs` - (Optional) Configure manual BGP advertised CIDRs for this connection. Available as of provider version R2.18+.
* `enable_event_triggered_ha` - (Optional) Enable Event Triggered HA. Default value: false. Valid values: true or false. Available as of provider version R2.19+.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_id` - VPC ID.


## Import

**device_transit_gateway_attachment** can be imported using the `connection_name`, e.g.

```
$ terraform import aviatrix_device_transit_gateway_attachment.test connection-name
```
