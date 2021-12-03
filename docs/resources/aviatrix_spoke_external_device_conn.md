---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_external_device_conn"
description: |-
  Creates and manages Aviatrix spoke external device connections
---

# aviatrix_spoke_external_device_conn

The **aviatrix_spoke_external_device_conn** resource creates and manages the connection between the Aviatrix BGP enabled spoke gateway and an External Device for purposes of Transit Network.

## Example Usage

```hcl
# Create an Aviatrix Spoke External Device Connection
resource "aviatrix_spoke_external_device_conn" "test" {
  vpc_id            = "vpc-abcd1234"
  connection_name   = "my_conn"
  gw_name           = "spokeGw"
  connection_type   = "bgp"
  bgp_local_as_num  = "123"
  bgp_remote_as_num = "345"
  remote_gateway_ip = "172.12.13.14"
}
```
```hcl
# Create an Aviatrix Spoke External Device Connection with HA enabled
resource "aviatrix_spoke_external_device_conn" "test" {
  vpc_id                   = "vpc-abcd1234"
  connection_name          = "my_conn"
  gw_name                  = "spokeGw"
  connection_type          = "static"
  remote_subnet            = "12.0.0.0/24"
  remote_gateway_ip        = "172.12.13.14"
  enable_ha                = true
  backup_remote_gateway_ip = "172.12.13.15"
}
```
```hcl
# Create an Aviatrix Spoke External Device Connection with Connection AS Path Prepend set
resource "aviatrix_spoke_external_device_conn" "test" {
  vpc_id            = "vpc-abcd1234"
  connection_name   = "my_conn"
  gw_name           = "spokeGw"
  connection_type   = "bgp"
  bgp_local_as_num  = "123"
  bgp_remote_as_num = "345"
  remote_gateway_ip = "172.12.13.14"
  prepend_as_path   = [
    "123",
    "123"
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Required) VPC ID of the Aviatrix spoke gateway.
* `connection_name` - (Required) Spoke external device connection name.
* `gw_name` - (Required) Aviatrix spoke gateway name.
* `remote_gateway_ip` - (Optional) Remote gateway IP.
* `connection_type` - (Required) Connection type. Valid values: 'bgp', 'static'. Default value: 'bgp'.

### HA
* `ha_enabled` - (Optional) Set as true if there are two external devices.
* `backup_remote_gateway_ip ` - (Optional) Backup remote gateway IP. Required if HA enabled.
* `backup_bgp_remote_as_num` - (Optional) Backup BGP remote ASN (Autonomous System Number). Integer between 1-4294967294. Required if HA enabled for 'bgp' connection.
* `backup_pre_shared_key` - (Optional) Backup Pre-Shared Key.
* `backup_local_tunnel_cidr` - (Optional) Source CIDR for the tunnel from the backup Aviatrix spoke gateway.
* `backup_remote_tunnel_cidr` - (Optional) Destination CIDR for the tunnel to the backup external device.
* `backup_direct_connect` - (Optional) Backup direct connect for backup external device.

### Custom Algorithms
* `custom_algorithms` - (Optional) Switch to enable custom/non-default algorithms for IPSec Authentication/Encryption. Valid values: true, false. **NOTE: Please see notes [here](#custom_algorithms-1) for more information.**
* `phase_1_authentication` - (Optional) Phase one Authentication. Valid values: 'SHA-1', 'SHA-256', 'SHA-384' and 'SHA-512'. Default value: 'SHA-256'.
* `phase_2_authentication` - (Optional) Phase two Authentication. Valid values: 'NO-AUTH', 'HMAC-SHA-1', 'HMAC-SHA-256', 'HMAC-SHA-384' and 'HMAC-SHA-512'. Default value: 'HMAC-SHA-256'.
* `phase_1_dh_groups` - (Optional) Phase one DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17', '18', '19', '20' and '21'. Default value: '14'.
* `phase_2_dh_groups` - (Optional) Phase two DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17', '18', '19', '20' and '21'. Default value: '14'.
* `phase_1_encryption` - (Optional) Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and 'AES-256-CBC'. Default value: 'AES-256-CBC'.
* `phase_2_encryption` - (Optional) Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', 'AES-256-CBC', 'AES-128-GCM-64', 'AES-128-GCM-96' and 'AES-128-GCM-128'. Default value: 'AES-256-CBC'.

### Misc.
* `tunnel_protocol` - (Optional) Tunnel protocol, only valid with `connection_type` = 'bgp'. Valid values: 'IPsec'. Default value: 'IPsec'. Case insensitive.
* `bgp_local_as_num` - (Optional) BGP local ASN (Autonomous System Number). Integer between 1-4294967294. Required for 'bgp' connection.
* `bgp_remote_as_num` - (Optional) BGP remote ASN (Autonomous System Number). Integer between 1-4294967294. Required for 'bgp' connection.
* `remote_subnet` - (Optional) Remote CIDRs joined as a string with ','. Required for a 'static' type connection.
* `direct_connect` - (Optional) Set true for private network infrastructure.
* `pre_shared_key` - (Optional) Pre-Shared Key.
* `local_tunnel_cidr` - (Optional) Source CIDR for the tunnel from the Aviatrix spoke gateway.
* `remote_tunnel_cidr` - (Optional) Destination CIDR for the tunnel to the external device.
* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval for the connection. Only valid with `connection_type` = 'bgp'. Requires the spoke_gateway's `learned_cidrs_approval_mode` attribute be set to 'connection'. Valid values: true, false. Default value: false.
* `approved_cidrs` - (Optional/Computed) Set of approved CIDRs. Requires `enable_learned_cidrs_approval` to be true. Type: Set(String).
* `enable_ikev2` - (Optional) Set as true to enable IKEv2 protocol.
* `manual_bgp_advertised_cidrs` - (Optional) Configure manual BGP advertised CIDRs for this connection. Only valid with `connection_type`= 'bgp'.
* `enable_event_triggered_ha` - (Optional) Enable Event Triggered HA. Default value: false. Valid values: true or false.
* `phase1_remote_identifier` - (Optional) Phase 1 remote identifier of the IPsec tunnel. This can be configured to be either the public IP address or the private IP address of the peer terminating the IPsec tunnel. Example: ["1.2.3.4"] when HA is disabled, ["1.2.3.4", "5.6.7.8"] when HA is enabled.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection.

## Import

**spoke_external_device_conn** can be imported using the `connection_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_spoke_external_device_conn.test connection_name~vpc_id
```

## Notes
### custom_algorithms
If set to true, the six algorithm arguments cannot all be default value. If set to false, default values will be used for all six algorithm arguments.
