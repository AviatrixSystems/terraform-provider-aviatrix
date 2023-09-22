---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_platform_ha"
description: |-
  Creates an Aviatrix Edge Platform HA
---

# aviatrix_edge_platform_ha

-> **NOTE:** A primary **aviatrix_edge_platform** is required to create **aviatrix_edge_platform_ha**.

The **aviatrix_edge_platform_ha** resource creates the Aviatrix Edge Platform HA.

## Example Usage

```hcl
# Create an Edge Platform HA
resource "aviatrix_edge_platform_ha" "test" {
  primary_gw_name = "primary_edge_platform"
  device_id       = "device5678"

  interfaces {
    name          = "eth0"
    type          = "WAN"
    ip_address    = "10.230.5.32/24"
    gateway_ip    = "10.230.5.100"
    wan_public_ip = "64.71.24.221"
  }

  interfaces {
    name       = "eth1"
    type       = "LAN"
    ip_address = "10.230.3.32/24"
  }

  interfaces {
    name        = "eth2"
    type        = "MANAGEMENT"
    enable_dhcp = false
    ip_address  = "172.16.15.162/20"
    gateway_ip  = "172.16.0.1"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_gw_name` - (Required) Primary Edge Platform name.
* `device_id` - (Required) Device ID.
* `interfaces` - (Required) WAN/LAN/MANAGEMENT interfaces.
    * `name` - (Required) Interface name.
    * `type` - (Required) Type.
    * `bandwidth` - (Optional) The rate of data can be moved through the interface, requires an integer value. Unit is in Mb/s.
    * `enable_dhcp` - (Optional) Enable DHCP. Valid values: true, false. Default value: false.
    * `wan_public_ip` - (Optional) WAN public IP.
    * `ip_address` - (Optional) Interface static IP address.
    * `gateway_ip` - (Optional) Gateway IP.
    * `dns_server_ip` - (Optional) Primary DNS server IP.
    * `secondary_dns_server_ip` - (Optional) Secondary DNS server IP.
    * `tag` - (Optional) Tag.

### Optional
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP and subnet prefix. Example: ["67.207.104.16/29", "64.71.12.144/29"].

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `account_name` - Edge Platform account name.

## Import

**edge_platform_ha** can be imported using the `primary_gw_name` in the form `primary_gw_name` + "-hagw" e.g.

```
$ terraform import aviatrix_edge_platform_ha.test primary_gw_name-hagw
```
