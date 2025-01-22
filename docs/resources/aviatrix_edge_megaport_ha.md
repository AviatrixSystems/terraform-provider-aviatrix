---
subcategory: "Edge"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_megaport_ha"
description: |-
  Creates an Aviatrix Edge Megaport HA
---

# aviatrix_edge_megaport_ha

The **aviatrix_edge_megaport_ha** resource creates the Aviatrix Edge Megaport HA.

-> **NOTE:** A primary **aviatrix_edge_megaport** is required to create **aviatrix_edge_megaport_ha**.

## Example Usage

```hcl
# Create an Edge Megaport HA
resource "aviatrix_edge_megaport_ha" "test" {
    primary_gw_name        = "primary_edge_megaport"
    ztp_file_download_path = "/ztp/file/download/path"

    interfaces {
        gateway_ip     = "10.220.14.1"
        ip_address     = "10.220.14.11/24"
        logical_ifname = "lan0"
        dns_server_ip  = "10.220.14.1"
    }

    interfaces {
        enable_dhcp    = true
        logical_ifname = "mgmt0"
    }

    interfaces {
        gateway_ip     = "192.168.99.1"
        ip_address     = "192.168.99.15/24"
        logical_ifname = "wan0"
        wan_public_ip  = "67.207.104.20"
        dns_server_ip  = "192.168.99.1"
    }

    interfaces {
        gateway_ip     = "192.168.88.1"
        ip_address     = "192.168.88.15/24"
        logical_ifname = "wan1"
        wan_public_ip  = "67.71.12.149"
        dns_server_ip  = "192.168.88.1"
    }

    interfaces {
        gateway_ip     = "192.168.77.1"
        ip_address     = "192.168.77.15/24"
        logical_ifname = "wan2"
        wan_public_ip  = "67.72.12.150"
        dns_server_ip  = "192.168.77.1"
    }

    management_egress_ip_prefix_list = ["162.43.147.139/31"]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_gw_name` - (Required) Primary Edge Megaport name.
* `ztp_file_download_path` - (Required) The folder path where the ZTP file will be downloaded.
* `interfaces` - (Required) WAN/LAN/MANAGEMENT interfaces.
  * `logical_ifname` - (Required) Logical interface name e.g., wan0, lan0, mgmt0.
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

* `account_name` - Edge Megaport account name.

## Import

**edge_megaport_ha** can be imported using the `primary_gw_name` in the form `primary_gw_name` + "-hagw" e.g.

```
$ terraform import aviatrix_edge_megaport_ha.test primary_gw_name-hagw
```
