---
subcategory: "CloudN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_cloudn_edge_gateway"
description: |-
  Creates Aviatrix Cloundn Edge Gateways
---

# aviatrix_cloudn_edge_gateway

The **aviatrix_cloudn_edge_gateway** resource creates the Aviatrix CloudN Edge Gateway. This resource is available as of provider version R2.22+.

## Example Usage

```hcl
# Create a DHCP CloudN Edge Gateway
resource "aviatrix_cloudn_edge_gateway" "test" {
  gw_name                    = "edge-test"
  management_connection_type = "DHCP"
  wan_interface_ip           = "10.60.0.0/24"
  wan_default_gateway        = "10.60.0.0"
  lan_interface_ip           = "10.60.0.0/24"
  image_download_path        = "/image/download/path"
  local_as_number            = "65000"
  prepend_as_path = [
    "65000",
    "65000",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `gw_name` - (Required) Edge gateway name.
* `management_connection_type` - (Required) Management connection type. Valid values: "DHCP", "Static". 
* `wan_interface_ip` - (Required) WAN interface IP.
* `wan_default_gateway` - (Required) WAN default gateway IP.
* `lan_interface_ip` - (Required) LAN interface IP.
* `image_download_path` - (Required) The folder path where the image file will be downloaded.

### Optional
* `over_private_network` - (Optional) Indicates whether it is public or private connection between controller and gateway. Valid values: true, false. Default value: false.
* `management_interface_ip` - (Optional) Management interface IP.
* `default_gateway_ip` - (Optional) Management default gateway IP.
* `dns_server` - (Optional) Management default gateway IP.
* `secondary_dns_server` - (Optional) Management default gateway IP.
* `local_as_number` - (Optional) BGP AS Number to assign to the gateway.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection. Requires local_as_number to be set. Type: List.

## Import

**cloudn_edge_gateway** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_cloudn_edge_gateway.test gw_name
```
