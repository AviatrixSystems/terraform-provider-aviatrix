---
subcategory: "CloudN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_caag"
description: |-
  Creates Aviatrix Edge as a CaaG
---

# aviatrix_edge_caag

The **aviatrix_edge_caag** resource creates the Aviatrix Edge as a CaaG. This resource is available as of provider version R2.22+.

## Example Usage

```hcl
# Create a DHCP Edge as a CaaG
resource "aviatrix_edge_caag" "test" {
  name                        = "edge-test"
  management_interface_config = "DHCP"
  wan_interface_ip_prefix     = "10.60.0.0/24"
  wan_default_gateway_ip      = "10.60.0.0"
  lan_interface_ip_prefix     = "10.60.0.0/24"
  ztp_file_type               = "iso"
  ztp_file_download_path      = "/image/download/path"
  local_as_number             = "65000"
  prepend_as_path = [
    "65000",
    "65000",
  ]
}
```
```hcl
# Create a Static Edge as a CaaG
resource "aviatrix_edge_caag" "test" {
  name                           = "edge-test"
  management_interface_config    = "Static"
  wan_interface_ip_prefix        = "10.60.0.0/24"
  wan_default_gateway_ip         = "10.60.0.0"
  lan_interface_ip_prefix        = "10.60.0.0/24"
  management_interface_ip_prefix = "10.60.0.0/24"
  management_default_gateway_ip  = "10.60.0.0"
  dns_server_ip                  = "10.60.0.0"
  secondary_dns_server_ip        = "10.60.0.0"
  ztp_file_type                  = "iso"
  ztp_file_download_path         = "/image/download/path"
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
* `name` - (Required) Edge as a CaaG name.
* `management_interface_config` - (Required) Management interface configuration. Valid values: "DHCP", "Static".
* `wan_interface_ip_prefix` - (Required) WAN interface IP and subnet prefix.
* `wan_default_gateway_ip` - (Required) WAN default gateway IP.
* `lan_interface_ip_prefix` - (Required) LAN interface IP and subnet prefix.
* `ztp_file_type` - (Required) ZTP file type. Valid values: "iso", "cloud-init".
* `ztp_file_download_path` - (Required) The folder path where the ZTP file will be downloaded.

### Optional
* `management_egress_ip_prefix` - (Optional) Management egress gateway IP and subnet prefix.
* `enable_over_private_network` - (Optional) Indicates whether it is public or private connection between controller and gateway. Valid values: true, false. Default value: false.
* `management_interface_ip_prefix` - (Optional) Management interface IP and subnet prefix. Required and valid when `management_interface_config` is "Static". 
* `management_default_gateway_ip` - (Optional) Management default gateway IP. Required and valid when `management_interface_config` is "Static".
* `dns_server_ip` - (Optional) DNS server IP. Required and valid when `management_interface_config` is "Static".
* `secondary_dns_server_ip` - (Optional) Secondary DNS server IP. Required and valid when `management_interface_config` is "Static".
* `local_as_number` - (Optional) BGP AS Number to assign to Edge as a CaaG.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection. Requires local_as_number to be set. Type: List.

## Import

**edge_caag** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_edge_caag.test name
```
