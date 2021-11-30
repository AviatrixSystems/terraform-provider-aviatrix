---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_cloudn_conn"
description: |-
  Creates and manages Aviatrix Transit Gateway to Aviatrix CloudN Connection
---

# aviatrix_transit_cloudn_conn

The **aviatrix_transit_cloudn_conn** resource creates and manages the connection between an Aviatrix Transit Gateway and an Aviatrix CloudN device. Available as of provider version R2.21.0+.

## Example Usage

```hcl
# Create an Aviatrix Transit Gateway to CloudN Connection
resource "aviatrix_transit_cloudn_conn" "test" {
  vpc_id                  = "vpc-abcd1234"
  connection_name         = "my_conn"
  gw_name                 = "transitGw"
  bgp_local_as_num        = "123"
  cloudn_as_num           = "345"
  cloudn_remote_ip        = "172.12.13.14"
  cloudn_neighbor_ip      = "182.1.2.3"
  cloudn_neighbor_as_num  = "65005"
}
```
```hcl
# Create an Aviatrix Transit Gateway to CloudN Connection with HA
resource "aviatrix_transit_cloudn_conn" "test" {
  vpc_id                        = "vpc-abcd1234"
  connection_name               = "my_conn"
  gw_name                       = "transitGw"
  bgp_local_as_num              = "123"
  cloudn_as_num                 = "345"
  cloudn_remote_ip              = "1.2.3.4"
  cloudn_neighbor_ip            = "1.2.3.5"
  cloudn_neighbor_as_num        = "65005"
  enable_ha                     = true
  backup_cloudn_ip              = "1.2.3.6"
  backup_cloudn_as_num          = "123"
  backup_cloudn_neighbor_ip     = "1.2.3.7"
  backup_cloudn_neighbor_as_num = "345"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Required) VPC ID of the Aviatrix Transit Gateway. Type: String.
* `connection_name` - (Required) Name of the Transit Gateway to CloudN Connection. Type: String.
* `gw_name` - (Required) Name of the Transit Gateway. Type: String.
* `bgp_local_as_num` - (Optional) BGP AS Number of the Transit Gateway. Type: String.
* `cloudn_as_num` - (Required) BGP AS Number of the Aviatrix CloudN. Type: String.
* `cloudn_remote_ip` - (Required) IP Address of Aviatrix CloudN. Type: String.
* `cloudn_neighbor_as_num` - (Required) BGP AS Number of the Aviatrix CloudN neighbor. Type: String.
* `cloudn_neighbor_ip` - (Required) IP Address of Aviatrix CloudN neighbor. Type: String.


### HA
* `enable_ha` - (Optional) Enable connection to HA CloudN. Type: Boolean.
* `backup_cloudn_ip` - (Optional) IP Address of HA CloudN. Required when `enable_ha` is true. Type: String.
* `backup_cloudn_as_num` (Optional) BGP AS Number of HA CloudN. Type: String.
* `backup_cloudn_neighbor_ip` - (Optional) IP Address of HA CloudN Neighbor. Type: String.
* `backup_cloudn_neighbor_as_num` - (Optional) BGP AS Number of HA CloudN Neighbor. Type: String.
* `backup_insane_mode` - (Optional) Flag to enable insane mode connection to HA CloudN. Type: Boolean.
* `backup_direct_connect` - (Optional) Flag to enable direct connect over private network to HA CloudN. Type: Boolean.
* `enable_load_balancing` - (Optional) Flag to enable load balancing between CloudN and HA CloudN. Type: Boolean.


### MISC
* `insane_mode` - (Optional) Enable insane mode connection. Type: Boolean.
* `direct_connect` - (Optional) Enable direct connect over private network. Type: Boolean. Default: true. 
* `enable_learned_cidrs_approval` - (Optional) Enable encrypted transit approval for connection. Type: Boolean.
* `approved_cidrs` - (Optional/Computed) Set of approved CIDRs. Requires `enable_learned_cidrs_approval` to be true. Type: Set(String).

## Import

**transit_cloudn_conn** can be imported using the `connection_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_transit_cloudn_conn.test connection_name~vpc_id
```
