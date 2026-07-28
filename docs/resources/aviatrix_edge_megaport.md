---
subcategory: "Edge"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_megaport"
description: |-
  Creates Aviatrix Edge Megaport
---

# aviatrix_edge_megaport

The **aviatrix_edge_megaport** resource creates the Aviatrix Edge Megaport.

## Example Usage

```hcl
# Create an Edge Megaport
resource "aviatrix_edge_megaport" "test" {
  account_name           = "edge_megaport-account"
  gw_name                = "megaport-test"
  site_id                = "site-123"
  ztp_file_download_path = "/ztp/file/download/path"
  interfaces {
    gateway_ip     = "10.220.14.1"
    ip_address     = "10.220.14.10/24"
    logical_ifname = "lan0"
    dns_server_ip  = "10.220.14.1"
  }

  interfaces {
    enable_dhcp    = true
    logical_ifname = "mgmt0"
  }

  interfaces {
    gateway_ip     = "192.168.99.1"
    ip_address     = "192.168.99.14/24"
    logical_ifname = "wan0"
    wan_public_ip  = "67.207.104.19"
    dns_server_ip  = "192.168.99.1"
  }

  interfaces {
    gateway_ip     = "192.168.88.1"
    ip_address     = "192.168.88.14/24"
    logical_ifname = "wan1"
    wan_public_ip  = "67.71.12.148"
    dns_server_ip  = "192.168.88.1"
  }

  interfaces {
    gateway_ip     = "192.168.77.1"
    ip_address     = "192.168.77.14/24"
    logical_ifname = "wan2"
    wan_public_ip  = "67.72.12.149"
    dns_server_ip  = "192.168.77.1"
  }
  management_egress_ip_prefix_list = ["162.43.147.137/31"]

  included_advertised_spoke_routes = [
    "10.10.0.0/16",
    "172.16.0.0/12"
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Edge Megaport account name.
* `gw_name` - (Required) Edge Megaport gateway name.
* `site_id` - (Required) Site ID.
* `ztp_file_download_path` - (Required) The folder path where the ZTP file will be downloaded.
* `interfaces` - (Required) WAN/LAN/MANAGEMENT interfaces.
  * `logical_ifname` - (Required) Logical interface name e.g., wan0, lan0, mgmt0.
  * `enable_dhcp` - (Optional) Enable DHCP. Valid values: true, false. Default value: false.
  * `wan_public_ip` - (Optional) WAN public IP.
  * `ip_address` - (Optional) Interface static IP address.
  * `gateway_ip` - (Optional) Gateway IP.
  * `dns_server_ip` - (Optional) Primary DNS server IP.
  * `secondary_dns_server_ip` - (Optional) Secondary DNS server IP.
  * `enable_vrrp` - (Optional) Enable VRRP. Valid values: true, false. Default value: false.
  * `vrrp_virtual_ip` - (Optional) VRRP virtual IP.
  * `tag` - (Optional) Tag.

### Optional
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP and subnet prefix. Example: ["67.207.104.16/29", "64.71.12.144/29"].
* `enable_management_over_private_network` - (Optional) Switch to enable management over the private network. Valid values: true, false. Default value: false.
* `enable_edge_active_standby` - (Optional) Switch to enable Edge Active-Standby mode. Valid values: true, false. Default value: false.
* `enable_edge_active_standby_preemptive` - (Optional) Switch to enable Preemptive Mode for Edge Active-Standby. Valid values: true, false. Default value: false.
* `dns_server_ip` - (Optional) DNS server IP. Required and valid when `management_interface_config` is "Static".
* `secondary_dns_server_ip` - (Optional) Secondary DNS server IP. Required and valid when `management_interface_config` is "Static".
* `local_as_number` - (Optional) BGP AS Number to assign to Edge Megaport.
* `prepend_as_path` - (Optional) List of AS numbers to prepend gateway BGP AS_Path field. Valid only when `local_as_number` is set. Example: ["65023", "65023"].
* `enable_learned_cidrs_approval` - (Optional) Switch to enable learned CIDR approval. Valid values: true, false. Default value: false.
* `approved_learned_cidrs` - (Optional) Set of approved learned CIDRs. Valid only when `enable_learned_cidrs_approval` is set to true. Example: ["10.1.0.0/16", "10.2.0.0/16"].
* `spoke_bgp_manual_advertise_cidrs` - (Optional) Set of intended CIDRs to be advertised to external BGP router. Example: ["10.1.0.0/16", "10.2.0.0/16"].
* `enable_preserve_as_path` - (Optional) Switch to enable preserve as_path when advertising manual summary CIDRs. Valid values: true, false. Default value: false.
* `bgp_polling_time` - (Optional) BGP route polling time in seconds. Valid values are between 10 and 50. Default value: 50.
* `bgp_neighbor_status_polling_time` - (Optional) BGP neighbor status polling time in seconds. Valid values are between 1 and 10. Default value: 5.
* `bgp_hold_time` - (Optional) BGP hold time in seconds. Valid values are between 12 and 360. Default value: 180.
* `enable_edge_transitive_routing` - (Optional) Switch to enable Edge transitive routing. Valid values: true, false. Default value: false.
* `enable_jumbo_frame` - (Optional) Switch to enable jumbo frame. Valid values: true, false. Default value: false.
* `latitude` - (Optional) Latitude of Edge Megaport. Valid values are between -90 and 90. Example: "47.7511".
* `longitude` - (Optional) Longitude of Edge Megaport. Valid values are between -180 and 180. Example: "120.7401".
* `rx_queue_size` - (Optional) Ethernet interface RX queue size. Once set, can't be deleted or disabled. Valid values: "1K", "2K", "4K".
* `vlan` - (Optional) VLAN configuration.
  * `parent_logical_interface_name` - (Required) Parent logical interface name e.g. lan0.
  * `vlan_id` - (Required) VLAN ID.
  * `ip_address` - (Optional) LAN sub-interface IP address.
  * `gateway_ip` - (Optional) LAN sub-interface gateway IP.
  * `peer_ip_address` - (Optional) LAN sub-interface IP address on HA gateway.
  * `peer_gateway_ip` - (Optional) LAN sub-interface gateway IP on HA gateway.
  * `vrrp_virtual_ip` - (Optional) LAN sub-interface virtual IP.
  * `tag` - (Optional) Tag.
* `enable_single_ip_snat` - (Optional) Enable Single IP SNAT. Valid values: true, false. Default value: false.
* `enable_auto_advertise_lan_cidrs` - (Optional) Enable auto advertise LAN CIDRs. Valid values: true, false. Default value: true.
* `included_advertised_spoke_routes` - (Optional) A list of CIDRs to be advertised to on-prem gateways as Included CIDR List. When configured, it will replace all advertised routes from this VPC.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `state` - State of Edge Megaport.

## Import

**edge_megaport** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_edge_megaport.test gw_name
```

## Deprecations
* Deprecated ``dns_server_ip`` and ``secondary_dns_server_ip``. These configuration values have no effect and have been replaced with ``dns_server_ip`` and  ``secondary_dns_server_ip`` present in **WAN/LAN/MGMT interfaces**. It will be removed from the Aviatrix provider in the 8.1.0 release.
