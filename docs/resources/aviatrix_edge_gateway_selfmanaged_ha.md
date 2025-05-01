---
subcategory: "Edge"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_gateway_selfmanaged_ha"
description: |-
  Creates Aviatrix Edge Gateway Selfmanaged HA
---

# aviatrix_edge_gateway_selfmanaged_ha

-> **NOTE:** A primary **aviatrix_edge_gateway_selfmanaged** is required to create **aviatrix_edge_gateway_selfmanaged_ha**.

The **aviatrix_edge_gateway_selfmanaged_ha** resource creates the Aviatrix Edge Gateway Selfmanaged HA.

## Example Usage

```hcl
# Create an Edge Gateway Selfmanaged HA
resource "aviatrix_edge_gateway_selfmanaged_ha" "test" {
  primary_gw_name         = "primary-edge-vm-selfmanaged"
  site_id                 = "site-123"
  ztp_file_type           = "iso"
  ztp_file_download_path  = "/ztp/download/path"

  interfaces {
    name          = "eth0"
    type          = "WAN"
    ip_address    = "10.230.6.32/24"
    gateway_ip    = "10.230.6.100"
    wan_public_ip = "64.71.25.221"
    dns_server_ip = "8.8.8.8"
    secondary_dns_server_ip = "8.8.6.6"
  }

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

  custom_interface_mapping {
    logical_ifname   = "wan0"
    identifier_type  = "system-assigned"
    identifier_value = "auto"
  }

  custom_interface_mapping {
    logical_ifname   = "lan0"
    identifier_type  = "mac"
    identifier_value = "00:0c:29:63:82:b2"
  }

  custom_interface_mapping {
    logical_ifname   = "mgmt0"
    identifier_type  = "pci"
    identifier_value = "pci@0000:04:00.0"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_gw_name` - (Required) Name of the primary Edge Gateway Selfmanaged.
* `site_id` - (Required) Site ID.

-> **NOTE:** At least one LAN interface is required.
* `interfaces` - (Required) WAN/LAN/MANAGEMENT interfaces.
  * `name` - (Required) Interface name.
  * `type` - (Required) Type.
  * `enable_dhcp` - (Optional) Enable DHCP. Valid values: true, false. Default value: false.
  * `wan_public_ip` - (Optional) WAN public IP.
  * `ip_address` - (Optional) Interface static IP address.
  * `gateway_ip` - (Optional) Gateway IP.
  * `dns_server_ip` - (Optional) Primary DNS server IP.
  * `secondary_dns_server_ip` - (Optional) Secondary DNS server IP.

### Optional
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP and subnet prefix. Example: ["67.207.104.16/29", "64.71.12.144/29"].
* `dns_server_ip` - (Optional) DNS server IP. Required and valid when `management_interface_config` is "Static".
* `secondary_dns_server_ip` - (Optional) Secondary DNS server IP. Required and valid when `management_interface_config` is "Static".
* `custom_interface_mapping` - (Optional) A list of custom interface mappings containing logical interfaces mapped to mac addresses or pci id's.
  * `logical_ifname` - (Required) Logical interface name must start with 'wan','lan' or 'mgmt' followed by a number (e.g., 'wan0', 'mgmt0', 'lan0').
  * `idenitifer_type` - (Required) Type of identifier used to map the logical interface to the physical interface e.g., mac, pci, system-assigned.
  * `idenitifer_value` - (Required) Value of the identifier used to map the logical interface to the physical interface. Can be a MAC address, PCI ID, or auto if system-assigned.

## Import

**edge_gateway_selfmanaged_ha** can be imported using the `primary_gw_name` in the form `primary_gw_name` + "-hagw" e.g.

```
$ terraform import aviatrix_edge_gateway_selfmanaged_ha.test primary_gw_name-hagw
```

## Deprecations
* Deprecated ``dns_server_ip`` and ``secondary_dns_server_ip``. These configuration values have no effect and have been replaced with ``dns_server_ip`` and  ``secondary_dns_server_ip`` present in **WAN/LAN/MGMT interfaces**. It will be removed from the Aviatrix provider in the 8.1.0 release.
