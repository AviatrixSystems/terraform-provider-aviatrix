---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_site2cloud"
description: |-
  Creates and Manages Aviatrix Site2Cloud connections
---

# aviatrix_site2cloud

The aviatrix_site2cloud resource Creates and Manages Aviatrix Site2Cloud connections.

## Example Usage

```hcl
# Create an Aviatrix Site2cloud
resource "aviatrix_site2cloud" "test_s2c" {
  vpc_id                     = "vpc-abcd1234"
  connection_name            = "my_conn"
  connection_type            = "unmapped"
  remote_gateway_type        = "generic"
  tunnel_type                = "udp"
  primary_cloud_gateway_name = "gw1"
  remote_gateway_ip          = "5.5.5.5"
  remote_subnet_cidr         = "10.23.0.0/24"
  local_subnet_cidr          = "10.20.1.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `primary_cloud_gateway_name` - (Required) Primary Cloud Gateway Name.
* `backup_gateway_name` - (Optional) Backup gateway name.
* `vpc_id` - (Required) VPC Id of the cloud gateway.
* `connection_name` - (Required) Site2Cloud Connection Name.
* `connection_type` - (Required) Connection Type. Valid Values: "mapped", "unmapped".
* `tunnel_type` - (Required) Site2Cloud Tunnel Type. Valid Values: "udp", "tcp".
* `remote_gateway_type` - (Required) Remote Gateway Type. Valid Values: "generic", "avx", "aws", "azure", "sonicwall", "oracle".
* `remote_gateway_ip` - (Required) Remote Gateway IP.
* `backup_remote_gateway_ip` - (Optional) Backup Remote Gateway IP.
* `pre_shared_key` - (Optional) Pre-Shared Key.
* `backup_pre_shared_key` - (Optional) Backup Pre-Shared Key.
* `remote_subnet_cidr` - (Required) Remote Subnet CIDR.
* `local_subnet_cidr` - (Optional) Local Subnet CIDR. Required for connection type "mapped".
* `remote_subnet_virtual` - Remote Subnet CIDR (Virtual). Required for connection type "mapped" only.
* `local_subnet_virtual` - Local Subnet CIDR (Virtual). Required for connection type "mapped" only.
* `ha_enabled` - (Optional) Specify whether enabling HA or not. Valid Values: true, false.
* `custom_algorithms` - (Optional) Switch to enable custom/non-default algorithms for IPSec Authentication/Encryption. Valid values: true, false.
* `phase_1_authentication` - (Optional) Phase one Authentication. Valid values: 'SHA-1', 'SHA-256', 'SHA-384' and 'SHA-512'. Default value: 'SHA-1'.
* `phase_2_authentication` - (Optional) Phase two Authentication. Valid values: 'NO-AUTH', 'HMAC-SHA-1', 'HMAC-SHA-256', 'HMAC-SHA-384' and 'HMAC-SHA-512'. Default value: 'HMAC-SHA-1'.
* `phase_1_dh_groups` - (Optional) Phase one DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'. Default value: '2'.
* `phase_2_dh_groups` - (Optional) Phase two DH Groups. Valid values: '1', '2', '5', '14', '15', '16', '17' and '18'. Default value: '2'.
* `phase_1_encryption` - (Optional) Phase one Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC' and 'AES-256-CBC'. Default value: 'AES-256-CBC'.
* `phase_2_encryption` - (Optional) Phase two Encryption. Valid values: '3DES', 'AES-128-CBC', 'AES-192-CBC', 'AES-256-CBC', 'AES-128-GCM-64', 'AES-128-GCM-96' and 'AES-128-GCM-128'. Default value: 'AES-256-CBC'.
* `private_route_encryption` - (Optional) Private route encryption switch. Valid values: true, false.
* `route_table_list` - (Optional) Route tables to modify.
* `remote_gateway_latitude` - (Optional) Latitude of remote gateway. Does not support refresh.
* `remote_gateway_longitude` - (Optional) Longitude of remote gateway. Does not support refresh.
* `backup_remote_gateway_latitude` - (Optional) Latitude of backup remote gateway. Does not support refresh.
* `backup_remote_gateway_longitude` - (Optional) Longitude of backup remote gateway. Does not support refresh.	 
* `ssl_server_pool` - (Optional) Specify ssl_server_pool for tunnel_type "tcp". Default value: "192.168.44.0/24".
* `enable_dead_peer_detection` - (Optional) Switch to Enable/Disable Deed Peer Detection for an existing site2cloud connection. Default value: true.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `local_subnet_cidr` - Local subnet CIDR.

-> **NOTE:** `custom_algorithms` - Only supported for 'udp' tunnel type. If set to true, the six algorithm arguments cannot all be default value. If set to false, default values will be used for all six algorithm arguments.

-> **NOTE:** `ssl_server_pool` - Only supported for 'tcp' tunnel type. If not set, default value will be used. If set, needs to be set to a different value than default value.

-> **NOTE:** `enable_dead_peer_detection` - If you are using/upgraded to Aviatrix Terraform Provider R1.9+, and an site2cloud resource was originally created with a provider version <R1.9, you must do ‘terraform refresh’ to update and apply the attribute’s default value (true) into the state file.

## Import

Instance site2cloud can be imported using the connection_name and vpc_id, e.g.

```
$ terraform import aviatrix_site2cloud.test connection_name~vpc_id
```