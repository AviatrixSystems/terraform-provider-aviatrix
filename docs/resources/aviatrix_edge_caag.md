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
  image_download_path         = "/image/download/path"
  local_as_number             = "65000"
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
* `image_download_path` - (Required) The folder path where the image file will be downloaded.

### Optional
* `enable_over_private_network` - (Optional) Indicates whether it is public or private connection between controller and gateway. Valid values: true, false. Default value: false.
* `management_interface_ip_prefix` - (Optional) Management interface IP and subnet prefix.
* `management_default_gateway_ip` - (Optional) Management default gateway IP.
* `dns_server_ip` - (Optional) DNS server IP.
* `secondary_dns_server` - (Optional) Secondary DNS server IP.
* `local_as_number` - (Optional) BGP AS Number to assign to Edge as a CaaG.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection. Requires local_as_number to be set. Type: List.

## Import

**cloudn_edge_gateway** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_cloudn_edge_gateway.test gw_name
```
