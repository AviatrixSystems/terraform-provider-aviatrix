---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_spoke_external_device_conn"
description: |-
  Creates and manages Edge as a Spoke external device connections
---

# aviatrix_edge_spoke_external_device_conn

The **aviatrix_edge_spoke_external_device_conn** resource creates and manages the connection between Edge as a Spoke and an External Device. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Create an Edge as a Spoke External Device Connection
resource "aviatrix_edge_spoke_external_device_conn" "test" {
  site_id           = "site-abcd1234"
  connection_name   = "conn"
  gw_name           = "eaas"
  bgp_local_as_num  = "123"
  bgp_remote_as_num = "345"
  local_lan_ip      = "10.230.3.23"
  remote_lan_ip     = "10.0.60.1"
}
```

## Argument Reference

The following arguments are supported:

### Required

-> **NOTE:** As of Controller version 6.8/provider version R2.23, only BGP over LAN connection is supported.

* `site_id` - (Required) Edge as a Spoke site iD.
* `connection_name` - (Required) Connection name.
* `gw_name` - (Required) Edge as a Spoke name.
* `bgp_local_as_num` - (Required) BGP local AS number.
* `bgp_remote_as_num` - (Required) BGP remote AS number.
* `local_lan_ip` - (Required) Local LAN IP.
* `remote_lan_ip` - (Required) Remote LAN IP.

### Optional
* `connection_type` - (Optional) Connection type. Valid value: 'bgp'. Default value: 'bgp'.
* `tunnel_protocol` - (Optional) Tunnel protocol. Valid value: 'LAN'. Default value: 'LAN'. Case insensitive.

## Import

**edge_spoke_external_device_conn** can be imported using the `connection_name` and `site_id`, e.g.

```
$ terraform import aviatrix_edge_spoke_external_device_conn.test connection_name~site_id
```
