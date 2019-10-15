---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vgw_conn"
description: |-
  Manages the Aviatrix Transit Gateway to VGW Connection
---

# aviatrix_vgw_conn

The aviatrix_vgw_conn resource managers the Aviatrix Transit Gateway to VGW Connection.

## Example Usage

```hcl
# Create an Aviatrix Vgw Connection
resource "aviatrix_vgw_conn" "test_vgw_conn" {
  conn_name        = "my-connection-vgw-to-tgw"
  gw_name          = "my-transit-gw"
  vpc_id           = "vpc-abcd1234"
  bgp_vgw_id       = "vgw-abcd1234"
  bgp_vgw_account  = "dev-account-1"  
  bgp_vgw_region   = "us-east-1"
  bgp_local_as_num = "65001"
}
```

## Argument Reference

The following arguments are supported:

* `conn_name` - (Required) The name of for Transit GW to VGW connection connection which is going to be created. Example: "my-connection-vgw-to-tgw".
* `gw_name` - (Required) Name of the Transit Gateway. Example: "my-transit-gw".
* `vpc_id` - (Required) VPC-ID where the Transit Gateway is located. Example: AWS: "vpc-abcd1234".
* `bgp_vgw_id` - (Required) Id of AWS's VGW that is used for this connection. Example: "vgw-abcd1234".
* `bgp_vgw_account` - (Required) Account of AWS's VGW that is used for this connection. Example: "dev-account-1".
* `bgp_vgw_region` - (Required) Region of AWS's VGW that is used for this connection. Example: "us-east-1".
* `bgp_local_as_num` - (Required) BGP Local ASN (Autonomous System Number). Integer between 1-65535. Example: "65001".

The following arguments are deprecated:

* `enable_advertise_transit_cidr` - (Optional) Switch to Enable/Disable advertise transit VPC network CIDR for a vgw connection.
* `bgp_manual_spoke_advertise_cidrs` - (Optional) Intended CIDR list to advertise to VGW. Example: "10.2.0.0/16,10.4.0.0/16".

-> **NOTE:** `enable_advertise_transit_cidr` - If you are using/upgraded to Aviatrix Terraform Provider R1.9+, and a vgw_conn resource was originally created with a provider version <R1.9, you must do ‘terraform refresh’ to update and apply the attribute’s default value (false) into the state file.

-> **NOTE:** `enable_advertise_transit_cidr` and `bgp_manual_spoke_advertise_cidrs` functionality has been migrated over to **aviatrix_transit_gateway** as of Aviatrix Terraform Provider R2.6. If you are using/upgraded to Aviatrix Terraform Provider R2.6+, and a vgw_conn resource was originally created with a provider version <R2.6, you must cut and paste these two arguments (and values) into the corresponding transit gateway resource referenced in this **vgw_conn**. A 'terraform refresh' will then successfully complete the migration and rectify the state file.


## Import

Instance vgw_conn can be imported using the conn_name and vpc_id, e.g.

```
$ terraform import aviatrix_vgw_conn.test conn_name~vpc_id
```
