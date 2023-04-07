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
# Create a bgp/GRE Aviatrix Spoke External Device Connection with jumbo frame enabled and ha enabled
resource "aviatrix_transit_external_device_conn" "test" {
  vpc_id                    = "vpc-abcd1234"
  connection_name           = "my_conn"
  gw_name                   = "spokeGw"
  connection_type           = "bgp"
  tunnel_protocol           = "GRE"
  bgp_local_as_num          = "65001"
  bgp_remote_as_num         = "65000"
  remote_gateway_ip         = "172.12.13.14"
  ha_enabled                = true
  local_tunnel_cidr         = "169.254.29.2/30,169.254.30.2/30"
  remote_tunnel_cidr        = "169.254.29.1/30,169.254.30.1/30"
  backup_local_tunnel_cidr  = "169.254.39.2/30,169.254.40.2/30"
  backup_remote_tunnel_cidr = "169.254.39.1/30,169.254.40.1/30"
  backup_bgp_remote_as_num  = "65000"
  backup_remote_gateway_ip  = "172.12.13.15"
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
```hcl
# Create a BGP over LAN Aviatrix Spoke External Device Connection with an Azure Spoke Gateway
resource "aviatrix_spoke_external_device_conn" "ex-conn" {
  vpc_id            = aviatrix_spoke_gateway.spoke-gateway.vpc_id
  connection_name   = "my_conn"
  gw_name           = aviatrix_spoke_gateway.spoke-gateway.gw_name
  connection_type   = "bgp"
  tunnel_protocol   = "LAN"
  bgp_local_as_num  = "123"
  bgp_remote_as_num = "345"
  remote_lan_ip     = "172.12.13.14"
  local_lan_ip      = "172.12.13.15"
  remote_vpc_name   = "vnet-name:resource-group-name:subscription-id"
}
```
```hcl
# Create a BGP over LAN Aviatrix HA Spoke External Device Connection with an Azure Spoke Gateway
resource "aviatrix_spoke_external_device_conn" "ex-conn" {
  vpc_id                   = aviatrix_spoke_gateway.spoke-gateway.vpc_id
  connection_name          = "my_conn"
  gw_name                  = aviatrix_spoke_gateway.spoke-gateway.gw_name
  connection_type          = "bgp"
  tunnel_protocol          = "LAN"
  bgp_local_as_num         = "123"
  bgp_remote_as_num        = "345"
  remote_lan_ip            = "172.12.13.14"
  local_lan_ip             = "172.12.13.15"
  remote_vpc_name          = "vnet-name:resource-group-name:subscription-id"
  ha_enabled               = true
  backup_bgp_remote_as_num = "678"
  backup_remote_lan_ip     = "172.12.13.16"
  backup_local_lan_ip      = "172.12.13.17"
}
```

## Argument Reference

The following arguments are supported:

### Required
-> **NOTE:** As of Controller version 6.5+/provider version R2.20+, the `vpc_id` for Gateways in Azure should be in the format "vnet_name:rg_name:resource_guid".

~> As of Provider version R2.21.2+, the `vpc_id` of an OCI VCN has been changed from its name to its OCID.
* `vpc_id` - (Required) VPC ID of the Aviatrix spoke gateway.
* `connection_name` - (Required) Spoke external device connection name.
* `gw_name` - (Required) Aviatrix spoke gateway name.
* `remote_gateway_ip` - (Optional) Remote gateway IP.
* `connection_type` - (Required) Connection type. Valid values: 'bgp', 'static'. Default value: 'bgp'.

~> **NOTE:** To create a BGP over LAN connection with an Azure Spoke Gateway, the Spoke Gateway must have its `enable_bgp` and `enable_bgp_over_lan` attributes set to true.

* `tunnel_protocol` - (Optional) Tunnel protocol, only valid with `connection_type` = 'bgp'. Valid values: 'IPsec', 'GRE' or 'LAN'. Default value: 'IPsec'. Case insensitive.
* `bgp_local_as_num` - (Optional) BGP local ASN (Autonomous System Number). Integer between 1-4294967294. Required for 'bgp' connection.
* `bgp_remote_as_num` - (Optional) BGP remote ASN (Autonomous System Number). Integer between 1-4294967294. Required for 'bgp' connection.
* `remote_subnet` - (Optional) Remote CIDRs joined as a string with ','. Required for a 'static' type connection.

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
* `phase_1_encryption` - (Optional) Phase one Encryption. Valid values: "3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "AES-128-GCM-64", "AES-128-GCM-96", "AES-128-GCM-128", "AES-256-GCM-64", "AES-256-GCM-96", and "AES-256-GCM-128". Default value: "AES-256-CBC".
* `phase_2_encryption` - (Optional) Phase two Encryption. Valid values: "3DES", "AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "AES-128-GCM-64", "AES-128-GCM-96", "AES-128-GCM-128", "AES-256-GCM-64", "AES-256-GCM-96", "AES-256-GCM-128" and "NULL-ENCR". Default value: "AES-256-CBC".

### BGP over LAN (Available as of provider version R2.18+)

~> **NOTE:** BGP over LAN attributes are only valid with `tunnel_protocol` = 'LAN'.

* `remote_vpc_name` - (Optional) Name of the remote VPC for a LAN BGP connection with an Azure Spoke Gateway. Required when `connection_type` = 'bgp' and `tunnel_protocol` = 'LAN' with an Azure spoke gateway. Must be in the format "<vnet-name>:<resource-group-name>:<subscription-id>". Available as of provider version R3.0.2+.
* `remote_lan_ip` - (Optional) Remote LAN IP. Required for BGP over LAN connection.
* `local_lan_ip` - (Optional) Local LAN IP. Required for GCP BGP over LAN connection.
* `backup_remote_lan_ip` - (Optional) Backup Remote LAN IP. Required for HA BGP over LAN connection.
* `backup_local_lan_ip` - (Optional) Backup Local LAN IP. Required for GCP HA BGP over LAN connection.
* `enable_bgp_lan_activemesh` - (Optional) Switch to enable BGP LAN ActiveMesh mode. Only valid for GCP and Azure with Remote Gateway HA enabled. Requires Azure Remote Gateway insane mode enabled. Valid values: true, false. Default: false. Available as of provider version R3.0.2+.

### BGP MD5 Authentication (Available as of provider version R2.21.1+)
~> **NOTE:** BGP MD5 Authentication is only valid with `connection_type` = 'bgp'.

* `bgp_md5_key` - (Optional) BGP MD5 Authentication Key. Example: 'avx01,avx02'.
* `backup_bgp_md5_key` - (Optional) Backup BGP MD5 Authentication Key. Valid with HA enabled for connection. Example: 'avx03,avx04'.

### Misc.
* `direct_connect` - (Optional) Set true for private network infrastructure.
* `pre_shared_key` - (Optional) Pre-Shared Key.
* `local_tunnel_cidr` - (Optional) Source CIDR for the tunnel from the Aviatrix spoke gateway.
* `remote_tunnel_cidr` - (Optional) Destination CIDR for the tunnel to the external device.
* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval for the connection. Only valid with `connection_type` = 'bgp'. Requires the spoke_gateway's `learned_cidrs_approval_mode` attribute be set to 'connection'. Valid values: true, false. Default value: false.
* `approved_cidrs` - (Optional/Computed) Set of approved CIDRs. Requires `enable_learned_cidrs_approval` to be true. Type: Set(String).
* `enable_ikev2` - (Optional) Set as true to enable IKEv2 protocol.
* `manual_bgp_advertised_cidrs` - (Optional) Configure manual BGP advertised CIDRs for this connection. Only valid with `connection_type`= 'bgp'.
* `enable_event_triggered_ha` - (Optional) Enable Event Triggered HA. Default value: false. Valid values: true or false.
* `enable_jumbo_frame` - (Optional) Enable Jumbo Frame for the transit external device connection. Only valid with 'GRE' tunnels under 'bgp' connection. Requires spoke to be jumbo frame and insane mode enabled. Valid values: true, false. Default value: false. Available as of provider version R3.0.2+.

-> **NOTE:** If you are using/upgraded to Aviatrix Terraform Provider R3.1.0+, and a **spoke_external_device_conn** resource was originally created with a provider version <R3.1.0 with "private_ip" for `phase1_local_identifier`, you must paste "phase1_local_identifier = 'private_ip'" into the corresponding **spoke_external_device_conn** resource to avoid ‘terraform plan‘ from showing delta.

* `phase1_local_identifier` - (Optional) Phase 1 local identifier. By default, gateway’s public IP is configured as the Local Identifier. Available as of provider version R3.1.0+.
* `phase1_remote_identifier` - (Optional) List of phase 1 remote identifier of the IPsec tunnel. This can be configured as a list of any string, including emtpy string. Example: ["1.2.3.4"] when HA is disabled, ["1.2.3.4", "abcd"] when HA is enabled. Available as of provider version R2.19+.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection.

## Import

**spoke_external_device_conn** can be imported using the `connection_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_spoke_external_device_conn.test connection_name~vpc_id
```

## Notes
### custom_algorithms
If set to true, the six algorithm arguments cannot all be default value. If set to false, default values will be used for all six algorithm arguments.
