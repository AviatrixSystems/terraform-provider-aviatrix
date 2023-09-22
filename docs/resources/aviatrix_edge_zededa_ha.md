---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_zededa_ha"
description: |-
  Creates Aviatrix Edge Zededa HA
---

# aviatrix_edge_zededa_ha

The **aviatrix_edge_zededa_ha** resource creates the Aviatrix Edge Zededa HA.

-> **NOTE:** A primary **aviatrix_edge_zededa** is required to create **aviatrix_edge_zededa_ha**.

## Example Usage

```hcl
# Create an Edge Zededa HA
resource "aviatrix_edge_zededa_ha" "test" {
  primary_gw_name   = "primary_edge_zededa"
  compute_node_uuid = "abcde12345"

  interfaces {
    name       = "eth1"
    type       = "LAN"
    ip_address = "10.220.11.20/24"
    gateway_ip = "10.220.11.1"
  }

  interfaces {
    name        = "eth2"
    type        = "MANAGEMENT"
    enable_dhcp = true
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_gw_name` - (Required) Primary Edge Zededa name.
* `compute_node_uuid` - (Required) Edge Zededa compute node UUID.

-> **NOTE:** At least one LAN interface is required.
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

* `account_name` - Edge Zededa account name.

## Import

**edge_zededa_ha** can be imported using the `primary_gw_name` in the form `primary_gw_name` + "-hagw" e.g.

```
$ terraform import aviatrix_edge_zededa_ha.test primary_gw_name-hagw
```
