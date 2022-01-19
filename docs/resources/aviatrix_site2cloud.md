---
subcategory: "Site2Cloud"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_site2cloud"
description: |-
  Create and manage Aviatrix Site2Cloud connections
---

# aviatrix_site2cloud

The **aviatrix_site2cloud** resource creates and manages Aviatrix-created Site2Cloud connections.

## Example Usage

```hcl
# Create an Aviatrix Site2cloud Connection
resource "aviatrix_site2cloud" "test_s2c" {
  vpc_id                     = "vpc-abcd1234"
  connection_name            = "my_conn"
  connection_type            = "unmapped"
  remote_gateway_type        = "generic"
  tunnel_type                = "policy"
  primary_cloud_gateway_name = "gw1"
  remote_gateway_ip          = "5.5.5.5"
  remote_subnet_cidr         = "10.23.0.0/24"
  local_subnet_cidr          = "10.20.1.0/24"
}
```
```hcl
# Create an Aviatrix Site2cloud Route Based Custom Mapped Connection
resource "aviatrix_site2cloud" "test_s2c" {
  vpc_id                           = "vpc-abcd1234"
  connection_name                  = "my_conn"
  connection_type                  = "mapped"
  remote_gateway_type              = "generic"
  tunnel_type                      = "route"
  primary_cloud_gateway_name       = "gw1"
  remote_gateway_ip                = "5.5.5.5"
  custom_mapped                    = true
  remote_source_real_cidrs         = ["10.10.0.0/24"]
  remote_source_virtual_cidrs      = ["10.10.1.0/24"]
  remote_destination_real_cidrs    = ["10.10.2.0/24"]
  remote_destination_virtual_cidrs = ["10.10.4.0/24"]
  local_source_real_cidrs          = ["10.11.0.0/24"]
  local_source_virtual_cidrs       = ["10.11.1.0/24"]
  local_destination_real_cidrs     = ["10.11.2.0/24"]
  local_destination_virtual_cidrs  = ["10.11.4.0/24"]
}
```

## Argument Reference

The following arguments are supported:

### Required
-> **NOTE:** As of Controller version 6.5+/provider version R2.20+, the `vpc_id` for Gateways in Azure should be in the format "vnet_name:rg_name:resource_guid".
* `vpc_id` - (Required) VPC ID of the cloud gateway.
* `connection_name` - (Required) Site2Cloud connection name.
* `remote_gateway_type` - (Required) Remote gateway type. Valid Values: "generic", "avx", "aws", "azure", "sonicwall", "oracle".
* `connection_type` - (Required) Connection type. Valid Values: "mapped", "unmapped".
* `tunnel_type` - (Required) Site2Cloud tunnel type. Valid Values: "policy", "route".
* `primary_cloud_gateway_name` - (Required) Primary cloud gateway name.
* `remote_gateway_ip` - (Required) Remote gateway IP.
* `remote_subnet_cidr` - (Required) Remote subnet CIDR. **Not required for custom_mapped connection.**
* `remote_subnet_virtual` - Remote subnet CIDR (Virtual). **Required for connection type "mapped", except for `custom_mapped` connection.**
* `local_subnet_cidr` - (Optional) Local subnet CIDR. **Required for connection type "mapped", except for `custom_mapped` connection.**
* `local_subnet_virtual` - Local subnet CIDR (Virtual). **Required for connection type "mapped", except for `custom_mapped` connection.**

### HA
* `ha_enabled` - (Optional) Specify whether or not to enable HA. Valid Values: true, false. **NOTE: Please see notes [here](#ha-enabled) regarding HA requirements.**
* `backup_gateway_name` - (Optional) Backup gateway name. **NOTE: Please see notes [here](#ha-enabled) regarding HA requirements.**
* `backup_remote_gateway_ip` - (Optional) Backup Remote Gateway IP. **NOTE: Please see notes [here](#ha-enabled) regarding HA requirements.**
* `backup_pre_shared_key` - (Optional) Backup Pre-Shared Key.
* `local_tunnel_ip` - (Optional) Local tunnel IP address. Only valid for route based connection. Available as of provider version R2.19+.
* `remote_tunnel_ip` - (Optional) Remote tunnel IP address. Only valid for route based connection. Available as of provider version R2.19+.
* `backup_local_tunnel_ip` - (Optional) Backup local tunnel IP address. Only valid when HA enabled route based connection. Available as of provider version R2.19+.
* `backup_remote_tunnel_ip` - (Optional) Backup remote tunnel IP address. Only valid when HA enabled route based connection. Available as of provider version R2.19+.
* `enable_single_ip_ha` - (Optional) Enable single IP HA feature. Available as of provider version 2.19+.

### Custom Algorithms
* `custom_algorithms` - (Optional) Switch to enable custom/non-default algorithms for IPSec Authentication/Encryption. Valid values: true, false. **NOTE: Please see notes [here](#custom_algorithms-1) for more information.**
* `phase_1_authentication` - (Optional) Phase one Authentication. Valid values: "SHA-1", "SHA-256", "SHA-384" and "SHA-512". Default value: "SHA-256".
* `phase_2_authentication` - (Optional) Phase two Authentication. Valid values: "NO-AUTH", "HMAC-SHA-1", "HMAC-SHA-256", "HMAC-SHA-384" and "HMAC-SHA-512". Default value: "HMAC-SHA-256".
* `phase_1_dh_groups` - (Optional) Phase one DH Groups. Valid values: "1", "2", "5", "14", "15", "16", "17", "18", "19", "20" and "21". Default value: "14".
* `phase_2_dh_groups` - (Optional) Phase two DH Groups. Valid values: "1", "2", "5", "14", "15", "16", "17", "18", "19", "20" and "21". Default value: "14".
* `phase_1_encryption` - (Optional) Phase one Encryption. Valid values: "3DES", "AES-128-CBC", "AES-192-CBC" and "AES-256-CBC". Default value: "AES-256-CBC".
* `phase_2_encryption` - (Optional) Phase two Encryption. Valid values: "3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "AES-128-GCM-64", "AES-128-GCM-96", "AES-128-GCM-128" and "NULL-ENCR". Default value: "AES-256-CBC".

### Encryption over ExpressRoute/DirectConnect
* `private_route_encryption` - (Optional) Private route encryption switch. Valid values: true, false.
* `route_table_list` - (Optional) Route tables to modify.
* `remote_gateway_latitude` - (Optional) Latitude of remote gateway. Does not support refresh.
* `remote_gateway_longitude` - (Optional) Longitude of remote gateway. Does not support refresh.
* `backup_remote_gateway_latitude` - (Optional) Latitude of backup remote gateway. Does not support refresh.
* `backup_remote_gateway_longitude` - (Optional) Longitude of backup remote gateway. Does not support refresh.

### Custom Mapped

~> **NOTE:** To enable custom mapped connection, 'connection_type' must be 'mapped' and 'tunnel_type' must be 'route'. All remote CIDR attributes or all local CIDR attributes must be set when using `custom_mapped`. Setting all CIDR attributes is also valid.

* `custom_mapped` - (Optional) Enable custom mapped connection. Default value: false. Valid values: true/false. Available in provider version R2.17.1+.
* `remote_source_real_cidrs` - (Optional) List of Remote Initiated Traffic Source Real CIDRs.
* `remote_source_virtual_cidrs` - (Optional) List of Remote Initiated Traffic Source Virtual CIDRs.
* `remote_destination_real_cidrs` - (Optional) List of  Remote Initiated Traffic Destination Real CIDRs.
* `remote_destination_virtual_cidrs` - (Optional) List of Remote Initiated Traffic Destination Virtual CIDRs.
* `local_source_real_cidrs` - (Optional) List of Local Initiated Traffic Source Real CIDRs.
* `local_source_virtual_cidrs` - (Optional) List of Local Initiated Traffic Source Virtual CIDRs.
* `local_destination_real_cidrs` - (Optional) List of Local Initiated Traffic Destination Real CIDRs.
* `local_destination_virtual_cidrs` - (Optional) List of Local Initiated Traffic Destination Virtual CIDRs.

### Misc.
* `pre_shared_key` - (Optional) Pre-Shared Key.
* `ssl_server_pool` - (Optional) Specify ssl_server_pool. Default value: "192.168.44.0/24". **NOTE: Please see notes [here](#ssl_server_pool-1) for more information.**
* `enable_dead_peer_detection` - (Optional) Enable/disable Deed Peer Detection for an existing site2cloud connection. Default value: true. **NOTE: Please see notes [here](#enable_dead_peer_detection-1) in regards to any deltas found in your state with the addition of this argument in R1.9**
* `enable_active_active` - (Optional) Enable/disable active active HA for an existing site2cloud connection. Valid values: true, false. Default value: false.
* `enable_ikev2` - (Optional) Switch to enable IKEv2. Valid values: true, false. Default value: false.
* `forward_traffic_to_transit` - (Optional) Enable spoke gateway with mapped site2cloud configurations to forward traffic from site2cloud connection to Aviatrix Transit Gateway. Default value: false. Valid values: true or false. Available in provider version 2.17.2+.
* `enable_event_triggered_ha` - (Optional) Enable Event Triggered HA. Default value: false. Valid values: true or false. Available as of provider version R2.19+.
* `phase1_remote_identifier` - (Optional) Phase 1 remote identifier of the IPsec tunnel. This can be configured to be either the public IP address or the private IP address of the peer terminating the IPsec tunnel. Example: ["1.2.3.4"] when HA is disabled, ["1.2.3.4", "5.6.7.8"] when HA is enabled. Available as of provider version R2.19+.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `local_subnet_cidr` - Local subnet CIDR.


## Import

**site2cloud** can be imported using the `connection_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_site2cloud.test connection_name~vpc_id
```


## Notes
### custom_algorithms
If set to true, the six algorithm arguments cannot all be default value. If set to false, default values will be used for all six algorithm arguments.

### enable_dead_peer_detection
If you are using/upgraded to Aviatrix Terraform Provider R1.9+, and a site2cloud resource was originally created with a provider version <R1.9, you must do ‘terraform refresh’ to update and apply the attribute’s default value (true) into the state file.

### HA Enabled
The following arguments are only supported if the backup gateway is set up by enabling peering HA through the primary gateway resource by specifying a `peering_ha_subnet` and `peering_ha_gw_size`. For more information on site2cloud, please see the doc site [here](https://docs.aviatrix.com/HowTos/site2cloud.html):

* `backup_gateway_name`
* `backup_remote_gateway_ip`
* `ha_enabled`

### ssl_server_pool
If not set, default value will be used. If set, needs to be set to a different value than the default value.
