---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_external_connection"
description: |-
  Creates and manages Aviatrix external connections
---

# aviatrix_external_connection

The **aviatrix_external_connection** resource creates and manages external connections between an Aviatrix gateway group and an external device. It supports full lifecycle management of both the connection and its individual tunnels, including in-place updates for mutable settings.

~> **NOTE:** This resource is the next-generation replacement for `aviatrix_transit_external_device_conn`. It uses the v2.5 API and identifies connections by UUID rather than name+vpc_id.

~> **NOTE:** In the first phase, this resource only supports **BGP over LAN** connections (`routing_protocol = "bgp"` and `tunnel_protocol = "LAN"`). Support for IPsec and GRE tunnel protocols will be added in future releases.

## Example Usage

### Basic BGP over LAN

```hcl
resource "aviatrix_external_connection" "basic" {
  name             = "my-ext-conn"
  group_uuid       = aviatrix_transit_gateway.gw.group_uuid
  routing_protocol = "bgp"
  tunnel_protocol  = "LAN"

  tunnels {
    local_gw_name    = "transit-gw"
    remote_device_ip = "10.0.4.10"
  }
}
```

### Full Options with Multiple Tunnels

```hcl
resource "aviatrix_external_connection" "full" {
  name                        = "my-ext-conn-full"
  group_uuid                  = aviatrix_transit_gateway.gw.group_uuid
  routing_protocol            = "bgp"
  tunnel_protocol             = "LAN"
  bgp_local_as_number         = 65000
  bgp_remote_as_number        = 65001
  enable_bgp_multihop         = true
  enable_bfd                  = true
  bfd_transmit_interval       = 300
  bfd_receive_interval        = 300
  bfd_detect_multiplier       = 3
  bgp_lan_activemesh          = true
  conn_learned_cidrs_approval = false

  tunnels {
    local_gw_name        = "transit-gw"
    remote_device_ip     = "10.0.4.10"
    bgp_remote_as_number = 65002
    bgp_md5_key          = "my-secret-key"
  }

  tunnels {
    local_gw_name    = "transit-gw-ha"
    remote_device_ip = "10.0.4.20"
    direct_connect   = true
  }
}
```

### Azure BGP over LAN with Peer VNet

```hcl
resource "aviatrix_external_connection" "azure" {
  name               = "azure-ext-conn"
  group_uuid         = aviatrix_transit_gateway.azure_gw.group_uuid
  routing_protocol   = "bgp"
  tunnel_protocol    = "LAN"
  bgp_local_as_number  = 65100
  bgp_remote_as_number = 65200
  peer_vnet_id       = "peer-vnet-name:resource-group:subscription-id"

  tunnels {
    local_gw_name    = "azure-transit-gw"
    remote_device_ip = "10.1.0.10"
  }
}
```

### IPv6 Enabled

```hcl
resource "aviatrix_external_connection" "ipv6" {
  name               = "ipv6-ext-conn"
  group_uuid         = aviatrix_transit_gateway.gw.group_uuid
  routing_protocol   = "bgp"
  tunnel_protocol    = "LAN"
  bgp_local_as_number  = 65000
  bgp_remote_as_number = 65001
  ipv6_enabled       = true

  tunnels {
    local_gw_name      = "transit-gw"
    remote_device_ip   = "10.0.4.10"
    local_device_ipv6  = "2001:db8::1/128"
    remote_device_ipv6 = "2001:db8::2/128"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) External connection name. Must be unique. Changing this forces a new resource.
* `group_uuid` - (Required) UUID of the gateway group this connection belongs to. Changing this forces a new resource.
* `routing_protocol` - (Required) Routing protocol for the connection. Valid values: `bgp`, `static`. Changing this forces a new resource. In the first phase, only `bgp` is accepted.
* `tunnel_protocol` - (Required) Tunnel protocol. Valid values: `LAN`, `IPsec`, `GRE`. Changing this forces a new resource. In the first phase, onely `LAN` is accepted.
* `tunnels` - (Required) One or more tunnel blocks. At least one tunnel must be specified. See [Tunnels](#tunnels) below.

### BGP

* `bgp_local_as_number` - (Optional) BGP local ASN. If not specified, the gateway group's local ASN is used. Changing this forces a new resource.
* `bgp_remote_as_number` - (Optional) Connection-level BGP remote ASN. If not specified, each tunnel must provide its own `bgp_remote_as_number`. Changing this forces a new resource.

### Mutable Connection Settings

These attributes can be updated in-place without recreating the connection.

* `enable_bgp_multihop` - (Optional) Enable BGP multihop. Valid values: true, false. Default: `true`.
* `jumbo_frame` - (Optional) Enable jumbo frames. Valid values: true, false. Default: `false`.
* `enable_bfd` - (Optional) Enable BFD (Bidirectional Forwarding Detection). Valid values: true, false. Default: `false`.
* `bfd_transmit_interval` - (Optional) BFD transmit interval in milliseconds. Valid values: 10-60000.
* `bfd_receive_interval` - (Optional) BFD receive interval in milliseconds. Valid values: 10-60000.
* `bfd_detect_multiplier` - (Optional) BFD detect multiplier. Valid values: 2-255.
* `conn_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval for this connection. Valid values: true, false. Default: `false`.

### Other

* `peer_vnet_id` - (Optional) Peer VNet ID. Required for Azure BGP over LAN connections. Format: `vnet-name:resource-group:subscription-id`. Changing this forces a new resource.
* `ipv6_enabled` - (Optional) Enable IPv6 for this connection. When true, each tunnel must include `local_device_ipv6` and `remote_device_ipv6`. Valid values: true, false. Default: `false`. Changing this forces a new resource.
* `edge_underlay` - (Optional) Whether the connection uses Edge underlay. When true, all gateways in the group must have a tunnel configured. Valid values: true, false. Default: `false`. Changing this forces a new resource.
* `bgp_lan_activemesh` - (Optional) Enable BGP LAN ActiveMesh mode. Only valid for GCP, Azure, and Edge gateways. Valid values: true, false. Default: `false`. Changing this forces a new resource.

### Tunnels

Each `tunnels` block supports the following:

* `local_gw_name` - (Required) Local Aviatrix gateway name. Must be a member of the gateway group.
* `remote_device_ip` - (Required) Remote device IPv4 address.
* `local_device_ip` - (Optional/Computed) Local device IPv4 address. If not provided, assigned by the server (for LAN protocol on AWS/Azure).
* `bgp_remote_as_number` - (Optional) Per-tunnel BGP remote ASN override. If the connection-level `bgp_remote_as_number` is not set, this is required on each tunnel.
* `local_device_ipv6` - (Optional) Local device IPv6 address with prefix length. Required when `ipv6_enabled` is true.
* `remote_device_ipv6` - (Optional) Remote device IPv6 address with prefix length. Required when `ipv6_enabled` is true.
* `bgp_md5_key` - (Optional, Sensitive) BGP MD5 authentication key. Can be updated in-place without recreating the tunnel.
* `direct_connect` - (Optional) Whether this tunnel uses direct-connect underlay. Valid values: true, false. Default: `false`.

~> **NOTE:** A tunnel is identified by the combination of `local_gw_name` and `remote_device_ip`. Changing either of these values will delete the old tunnel and create a new one. Changes to other tunnel fields (except `bgp_md5_key`) also trigger a delete+recreate. `bgp_md5_key` changes are applied in-place via PATCH.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Server-assigned UUID of the external connection. This is also used as the resource ID.
* `tunnels.*.status` - Tunnel status as reported by the server. Possible values: `UP`, `DOWN`, `UNKNOWN`.

## Import

**aviatrix_external_connection** can be imported using the connection UUID:

```
$ terraform import aviatrix_external_connection.example 5f4b1a3e-2c8d-4a7e-9b1f-0a3c2d6e8f01
```

~> **NOTE:** `bgp_md5_key` is not returned by the server and will not be populated after import. You must re-add the `bgp_md5_key` values in your configuration file and run `terraform apply` to re-sync.
