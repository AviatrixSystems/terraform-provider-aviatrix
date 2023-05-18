---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_spoke_external_device_conn"
description: |-
  Creates and manages Edge as a Spoke external device connections
---

# aviatrix_edge_spoke_external_device_conn

The **aviatrix_edge_spoke_external_device_conn** resource creates and manages the connection between Edge as a Spoke and an External Device. This resource is available as of provider version R2.23+.

~> **NOTE:** Please use a separate **aviatrix_edge_spoke_external_device_conn** to create WAN underlay connection for Edge HA.

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
* `gw_name` - (Required) Edge as a Spoke name.
* `bgp_local_as_num` - (Required) BGP local AS number.
* `bgp_remote_as_num` - (Required) BGP remote AS number.
* `local_lan_ip` - (Required) Local LAN IP.
* `remote_lan_ip` - (Required) Remote LAN IP.

### Optional

-> **NOTE:** When `enable_edge_underlay` is false, `connection_name` is required. When `enable_edge_underlay` is true, `connection_name` must be empty. 

* `connection_name` - (Optional) Connection name.
* `connection_type` - (Optional) Connection type. Valid value: 'bgp'. Default value: 'bgp'.
* `tunnel_protocol` - (Optional) Tunnel protocol. Valid value: 'LAN'. Default value: 'LAN'. Case insensitive.
* `enable_edge_underlay` - (Optional) Enable BGP over WAN underlay. Valid values: true, false. Default value: false.
* `remote_cloud_type` - (Optional) Remote cloud type. Valid values: "AWS", "AZURE".
* `number_of_retries` - (Optional) Number of retries. Default value: 0.
* `retry_interval` - (Optional) Retry interval in seconds. Default value: 300.
* `ha_enabled` - (Optional) Set as true if there are two external devices.
* `backup_remote_lan_ip` - (Optional) Backup Remote LAN IP. Required for HA BGP over LAN connection.
* `backup_local_lan_ip` - (Optional) Backup Local LAN IP. Required for HA BGP over LAN connection.
* `backup_bgp_remote_as_num` - (Optional) Backup BGP remote ASN (Autonomous System Number). Integer between 1-4294967294. Required if HA enabled for 'bgp' connection.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection.
* `manual_bgp_advertised_cidrs` - (Optional) Configure manual BGP advertised CIDRs for this connection.

## Import

**edge_spoke_external_device_conn** can be imported using the `connection_name`, `site_id` and `gw_name`, e.g.

```
$ terraform import aviatrix_edge_spoke_external_device_conn.test connection_name~site_id~gw_name
```
