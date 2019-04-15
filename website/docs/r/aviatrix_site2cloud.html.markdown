---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_site2cloud"
sidebar_current: "docs-aviatrix-resource-site2cloud"
description: |-
  Creates and Manages Aviatrix Site2Cloud connection
---

# aviatrix_site2cloud

The Site2Cloud resource Creates and Manages Aviatrix Site2Cloud connection

## Example Usage

```hcl
# Create Aviatrix site2cloud
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

* `primary_cloud_gateway_name` - (Required) Primary Cloud Gateway Name
* `backup_gateway_name` - (Optional)
* `vpc_id` - (Required) VPC Id of the cloud gateway
* `connection_name` - (Required) Site2Cloud Connection Name
* `connection_type` - (Required) Connection Type. Valid Value(s): "mapped", "unmapped"
* `tunnel_type` - (Required) Site2Cloud Tunnel Type. Valid Value(s): "udp", "tcp"
* `remote_gateway_type` - (Required) Remote Gateway Type. Valid Value(s): "generic", "avx", "aws", "azure", "sonicwall", "oracle".
* `remote_gateway_ip` - (Required) Remote Gateway IP
* `backup_remote_gateway_ip` - (Optional)
* `pre_shared_key` - (Optional) Pre-Shared Key
* `backup_pre_shared_key` - (Optional) Backup Pre-Shared Key
* `remote_subnet_cidr` - (Required) Remote Subnet CIDR
* `local_subnet_cidr` - (Optional) Local Subnet CIDR
* `ha_enabled` - (Optional) Specify whether enabling HA or not. Valid Value(s): "yes", "no"
* `ssl_server_pool` - (Optional) Specify ssl_server_pool for tunnel_type "tcp". Default value is "192.168.44.0/24".

## Import

Instance site2cloud can be imported using the connection_name and vpc_id, e.g.

```
$ terraform import aviatrix_site2cloud.test connection_name~vpc_id
```