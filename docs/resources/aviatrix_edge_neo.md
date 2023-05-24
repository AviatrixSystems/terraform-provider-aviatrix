---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_neo"
description: |-
  Creates Aviatrix Edge NEO
---

# aviatrix_edge_neo

-> **NOTE:** Since V3.1.1+, please use resource **aviatrix_edge_platform** instead. Resource **aviatrix_edge_neo** will be deprecated in the V3.2.0 release.

The **aviatrix_edge_neo** resource creates the Aviatrix Edge NEO.

## Example Usage

```hcl
# Create an Edge NEO
resource "aviatrix_edge_neo" "test" {
  account_name               = "edge_neo-account"
  gw_name                    = "neo-test"
  site_id                    = "site-123"
  device_id                  = "abcd123"
  gw_size                    = "small"
  management_interface_names = ["eth2"]
  lan_interface_names        = ["eth1"]
  wan_interface_names        = ["eth0"]

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
* `account_name` - (Required) Edge NEO account name.
* `gw_name` - (Required) Edge NEO name.
* `site_id` - (Required) Site ID.
* `device_id` - (Required) Device ID.
* `gw_size` - (Required) Gateway size. Valid values: "small", "medium", "large", "x-large".
* `wan_interface_names` - (Required) List of WAN interface names. Default value: ["eth0"].
* `lan_interface_names` - (Required) List of LAN interface names. Default value: ["eth1"].
* `management_interface_names` - (Required) List of management interface names. Default value: ["eth2"].
* `interfaces` - (Required) WAN/LAN/MANAGEMENT interfaces.
    * `name` - (Required) Interface name.
    * `type` - (Required) Type. Valid values: WAN, LAN, or MANAGEMENT.
    * `bandwidth` - (Optional) The rate of data can be moved through the interface, requires an integer value. Unit is in Mb/s.
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
* `wan_interface_names` - (Optional) List of WAN interface names. Default value: ["eth0"].
* `lan_interface_names` - (Optional) List of LAN interface names. Default value: ["eth1"].
* `management_interface_names` - (Optional) List of management interface names. Default value: ["eth2"].
* `local_as_number` - (Optional) BGP AS Number to assign to Edge NEO.
* `prepend_as_path` - (Optional) List of AS numbers to prepend gateway BGP AS_Path field. Valid only when `local_as_number` is set. Example: ["65023", "65023"].
* `enable_learned_cidrs_approval` - (Optional) Switch to enable learned CIDR approval. Valid values: true, false. Default value: false.
* `approved_learned_cidrs` - (Optional) Set of approved learned CIDRs. Valid only when `enable_learned_cidrs_approval` is set to true. Example: ["10.1.0.0/116", "10.2.0.0/16"].
* `spoke_bgp_manual_advertise_cidrs` - (Optional) Set of intended CIDRs to be advertised to external BGP router. Example: ["10.1.0.0/116", "10.2.0.0/16"].
* `enable_preserve_as_path` - (Optional) Switch to enable preserve as_path when advertising manual summary CIDRs. Valid values: true, false. Default value: false.
* `bgp_polling_time` - (Optional) BGP route polling time in seconds. Valid values are between 10 and 50. Default value: 50.
* `bgp_hold_time` - (Optional) BGP hold time in seconds. Valid values are between 12 and 360. Default value: 180.
* `enable_edge_transitive_routing` - (Optional) Switch to enable Edge transitive routing. Valid values: true, false. Default value: false.
* `enable_jumbo_frame` - (Optional) Switch to enable jumbo frame. Valid values: true, false. Default value: false.
* `latitude` - (Optional) Latitude of Edge NEO. Valid values are between -90 and 90. Example: "47.7511".
* `longitude` - (Optional) Longitude of Edge Neo. Valid values are between -180 and 180. Example: "120.7401".
* `rx_queue_size` - (Optional) Ethernet interface RX queue size. Once set, can't be deleted or disabled. Valid values: "1K", "2K", "4K".
* `vlan` - (Required) VLAN configuration.
    * `parent_interface_name` - (Required) Parent interface name.
    * `vlan_id` - (Required) VLAN ID.
    * `ip_address` - (Optional) LAN sub-interface IP address.
    * `gateway_ip` - (Optional) LAN sub-interface gateway IP.
    * `peer_ip_address` - (Optional) LAN sub-interface IP address on HA gateway.
    * `peer_gateway_ip` - (Optional) LAN sub-interface gateway IP on HA gateway.
    * `vrrp_virtual_ip` - (Optional) LAN sub-interface virtual IP.
    * `tag` - (Optional) Tag.
* `dns_profile_name` - (Optional) DNS profile to be associated with gateway, select an existing template.
* `enable_single_ip_snat` - (Optional) Enable Single IP SNAT. Valid values: true, false. Default value: false.
* `enable_auto_advertise_lan_cidrs` - (Optional) Enable auto advertise LAN CIDRs. Valid values: true, false. Default value: true.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `state` - State of Edge NEO.

## Import

**edge_neo** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_edge_equinix.test gw_name
```
