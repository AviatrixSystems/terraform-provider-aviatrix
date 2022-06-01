---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_spoke"
description: |-
  Creates Aviatrix Edge as a Spoke
---

# aviatrix_edge_spoke

The **aviatrix_edge_spoke** resource creates the Aviatrix Edge as a Spoke. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Create a DHCP Edge as a Spoke
resource "aviatrix_edge_spoke" "test" {
  gw_name                     = "edge-test"
  site_id                     = "site-123"
  management_interface_config = "DHCP"
  wan_interface_ip_prefix     = "10.60.0.0/24"
  wan_default_gateway_ip      = "10.60.0.0"
  lan_interface_ip_prefix     = "10.60.0.0/24"
  ztp_file_type               = "iso"
  ztp_file_download_path      = "/ztp/download/path"
  local_as_number             = "65000"
  prepend_as_path = [
    "65000",
    "65000",
  ]
}
```
```hcl
# Create a Static Edge as a Spoke
resource "aviatrix_edge_spoke" "test" {
  gw_name                        = "edge-test"
  site_id                        = "site-123"
  management_interface_config    = "Static"
  wan_interface_ip_prefix        = "10.60.0.0/24"
  wan_default_gateway_ip         = "10.60.0.0"
  lan_interface_ip_prefix        = "10.60.0.0/24"
  management_interface_ip_prefix = "10.60.0.0/24"
  management_default_gateway_ip  = "10.60.0.0"
  dns_server_ip                  = "10.60.0.0"
  secondary_dns_server_ip        = "10.60.0.0"
  ztp_file_type                  = "iso"
  ztp_file_download_path         = "/ztp/download/path"
  local_as_number                = "65000"
  prepend_as_path = [
    "65000",
    "65000",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `gw_name` - (Required) Edge as a Spoke name.
* `site_id` - (Required) Site ID.
* `management_interface_config` - (Required) Management interface configuration. Valid values: "DHCP", "Static".
* `wan_interface_ip_prefix` - (Required) WAN interface IP and subnet prefix.
* `wan_default_gateway_ip` - (Required) WAN default gateway IP.
* `lan_interface_ip_prefix` - (Required) LAN interface IP and subnet prefix.
* `ztp_file_type` - (Required) ZTP file type. Valid values: "iso", "cloud-init".
* `ztp_file_download_path` - (Required) The folder path where the ZTP file will be downloaded.

### Optional
* `management_egress_ip_prefix` - (Optional) Management egress gateway IP and subnet prefix.
* `enable_over_private_network` - (Optional) Switch to enable management over the private network. Valid values: true, false. Default value: false.
* `enable_active_standby` - (Optional) Switch to enable Active-Standby mode. Valid values: true, false. Default value: false.
* `enable_active_standby_preemptive` - (Optional) Switch to enable Preemptive Mode for Active-Standby. Valid values: true, false. Default value: false.
* `management_interface_ip_prefix` - (Optional) Management interface IP and subnet prefix. Required and valid when `management_interface_config` is "Static".
* `management_default_gateway_ip` - (Optional) Management default gateway IP. Required and valid when `management_interface_config` is "Static".
* `dns_server_ip` - (Optional) DNS server IP. Required and valid when `management_interface_config` is "Static".
* `secondary_dns_server_ip` - (Optional) Secondary DNS server IP. Required and valid when `management_interface_config` is "Static".

### Advanced Options
* `local_as_number` - (Optional) BGP AS Number to assign to Edge as a Spoke.
* `prepend_as_path` - (Optional) List of AS numbers to prepend gateway BGP AS_Path field. Valid only when `local_as_number` is set. Example: ["65023", "65023"].
* `edge_gateway_static_local_cidrs` - (Optional) Set of CIDRs to be advertised as 'Edge Gateway Static Local CIDRs'. Example: ["10.1.0.0/116", "10.2.0.0/16"].
* `enable_learned_cidrs_approval` - (Optional) Switch to enable learned CIDR approval. Valid values: true, false. Default value: false.
* `approved_learned_cidrs` - (Optional) Set of approved learned CIDRs. Valid only when `enable_learned_cidrs_approval` is set to true. Example: ["10.1.0.0/116", "10.2.0.0/16"].
* `spoke_bgp_manual_advertise_cidrs` - (Optional) Set of intended CIDRs to be advertised to external BGP router. Example: ["10.1.0.0/116", "10.2.0.0/16"].
* `enable_preserve_as_path` - (Optional) Switch to enable preserve as_path when advertising manual summary CIDRs. Valid values: true, false. Default value: false.
* `bgp_polling_time` - (Optional) BGP route polling time. Unit is in seconds. Valid values are between 10 and 50. Default value: 50.
* `bgp_hold_time` - (Optional) BGP hold time. Unit is in seconds. Valid values are between 12 and 360. Default value: 180.
* `enable_edge_transitive_routing` - (Optional) Switch to enable Edge transitive routing. Valid values: true, false. Default value: false.
* `enable_jumbo_frame` - (Optional) Switch to enable jumbo frame. Valid values: true, false. Default value: false.
* `latitude` - (Optional) Latitude of Edge as a Spoke. Valid values are between -90 and 90. Example: 47.7511.
* `longitude` - (Optional) Longitude of Edge as a Spoke. Valid values are between -180 and 180. Example: 120.7401.

## Import

**edge_spoke** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_edge_spoke.test gw_name
```
