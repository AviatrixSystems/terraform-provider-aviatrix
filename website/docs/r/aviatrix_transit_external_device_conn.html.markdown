---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_external_device_conn"
description: |-
  Creates and manages Aviatrix transit external device connections
---

# aviatrix_transit_external_device_conn

The **aviatrix_transit_external_device_conn** resource creates and manages the connection between the Aviatrix transit gateway and an External Device for purposes of Transit Network.

For more information on this feature, please see documentation [here](https://docs.aviatrix.com/HowTos/transitgw_external.html).

## Example Usage

```hcl
# Create an Aviatrix Transit External Device Connection
resource "aviatrix_transit_external_device_conn" "test" {
  vpc_id            = "vpc-abcd1234"
  connection_name   = "my_conn"
  gw_name           = "transitGw"
  connection_type   = "bgp"
  bgp_local_as_num  = "123"
  bgp_remote_as_num = "345"
  remote_gateway_ip = "172.12.13.14"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Required) VPC ID of the Aviatrix transit gateway.
* `connection_name` - (Required) Transit external device connection name.
* `gw_name` - (Required) Aviatrix transit gateway name.
* `remote_gateway_ip` - (Required) Remote gateway IP.
* `connection_type` - (Required) Connection type. Valid values: 'bpg', 'static'. Default value: 'bgp'.
* `bgp_local_as_num` - (Optional) BGP local ASN (Autonomous System Number). Integer between 1-4294967294. Required for 'bgp' connection.
* `bgp_remote_as_num` - (Optional) BGP remote ASN (Autonomous System Number). Integer between 1-4294967294. Required for 'bgp' connection.
* `remote_subnet` - (Optional) Remote CIDRs joined as a string with ','. Required for a 'static' type connection.

### HA
* `ha_enabled` - (Optional) Set as true if there are two external devices.
* `backup_remote_gateway_ip ` - (Optional) Backup remote gateway IP. Required if HA enabled.
* `backup_bgp_remote_as_num` - (Optional) Backup BGP remote ASN (Autonomous System Number). Integer between 1-4294967294. Required if HA enabled for 'bgp' connection.
* `backup_pre_shared_key` - (Optional) Backup Pre-Shared Key.
* `backup_local_tunnel_cidr` - (Optional) Source CIDR for the tunnel from the backup Aviatrix transit gateway.
* `backup_remote_tunnel_cidr` - (Optional) Destination CIDR for the tunnel to the backup external device.
* `backup_direct_connect` - (Optional) Backup direct connect for backup external device.

### Custom Algorithms
* `custom_algorithms` - (Optional) Switch to enable custom/non-default algorithms for IPSec Authentication/Encryption. Valid values: true, false. **NOTE: Please see notes [here](#custom_algorithms-1) for more information.**
* `phase_1_authentication` - (Optional) Phase one Authentication. Valid values: 'SHA-1', 'SHA-256', 'SHA-384' and 'SHA-512'. Default value: 'SHA-256'.
* `phase_2_authentication` - (Optional) Phase two Authentication. Valid values: 'NO-AUTH', 'HMAC-SHA-1', 'HMAC-SHA-256', 'HMAC-SHA-384' and 'HMAC-SHA-512'. Default value: 'HMAC-SHA-256'.
* `phase_1_dh_groups` - (Optional) Phase one DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'. Default value: '14'.
* `phase_2_dh_groups` - (Optional) Phase two DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'. Default value: '14'.
* `phase_1_encryption` - (Optional) Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and 'AES-256-CBC'. Default value: 'AES-256-CBC'.
* `phase_2_encryption` - (Optional) Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', 'AES-256-CBC', 'AES-128-GCM-64', 'AES-128-GCM-96' and 'AES-128-GCM-128'. Default value: 'AES-256-CBC'.

### Misc.
* `direct_connect` - (Optional) Set true for private network infrastructure.
* `pre_shared_key` - (Optional) Pre-Shared Key.
* `local_tunnel_cidr` - (Optional) Source CIDR for the tunnel from the Aviatrix transit gateway.
* `remote_tunnel_cidr` - (Optional) Destination CIDR for the tunnel to the external device.
* `enable_edge_segmentation` - (Optional) Switch to allow this connection to communicate with a Security Domain via Connection Policy.
* `switch_to_ha_standby_gateway` - (Optional) Switch to HA Standby Transit Gateway connection. Only valid with Transit Gateway that has [Active-Standby Mode](https://docs.aviatrix.com/HowTos/transit_advanced.html#active-standby) enabled and for non-HA external device. Valid values: true, false. Default: false. Available in provider version R2.17.1+.
* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval for the connection. Only valid with 'connection_type' = 'bgp'. Requires the transit_gateway's 'learned_cidrs_approval_mode' attribute be set to 'connection'. Valid values: true, false. Default value: false. Available as of provider version R2.18+.
* `manual_bgp_advertised_cidrs` - (Optional) Configure manual BGP advertised CIDRs for this connection. Only valid with 'connection_type'= 'bgp'. Available as of provider version R2.18+.

## Import

**transit_external_device_conn** can be imported using the `connection_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_transit_external_device_conn.test connection_name~vpc_id
```

## Notes
### custom_algorithms
If set to true, the six algorithm arguments cannot all be default value. If set to false, default values will be used for all six algorithm arguments.
