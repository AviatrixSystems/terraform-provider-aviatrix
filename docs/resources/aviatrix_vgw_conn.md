---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vgw_conn"
description: |-
  Manages the connection between the Aviatrix Transit Gateway to VGW
---

# aviatrix_vgw_conn

The **aviatrix_vgw_conn** resource manages the connection between the Aviatrix transit gateway and AWS VGW for purposes of Transit Network.

## Example Usage

```hcl
# Create an Aviatrix VGW Connection
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

### Required
* `conn_name` - (Required) The name of for Transit GW to VGW connection connection which is going to be created. Example: "my-connection-vgw-to-tgw".
* `gw_name` - (Required) Name of the Transit Gateway. Example: "my-transit-gw".
* `vpc_id` - (Required) VPC ID where the Transit Gateway is located. Example: AWS: "vpc-abcd1234".
* `bgp_vgw_id` - (Required) ID of AWS VGW that will be used for this connection. Example: "vgw-abcd1234".
* `bgp_vgw_account` - (Required) Cloud Account used to create the AWS VGW that will be used for this connection. Example: "dev-account-1".
* `bgp_vgw_region` - (Required) Region of AWS VGW that will be used for this connection. Example: "us-east-1".
* `bgp_local_as_num` - (Required) BGP Local ASN (Autonomous System Number). Integer between 1-4294967294. Example: "65001".

### Optional
* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDRs approval for the connection. Requires the transit_gateway's 'learned_cidrs_approval_mode' attribute be set to 'connection'. Valid values: true, false. Default value: false. Available as of provider version R2.18+.
* `manual_bgp_advertised_cidrs` - (Optional) Configure manual BGP advertised CIDRs for this connection. Available as of provider version R2.18+.
* `enable_event_triggered_ha` - (Optional) Enable Event Triggered HA. Default value: false. Valid values: true or false. Available as of provider version R2.19+.

The following arguments are deprecated:

* `enable_advertise_transit_cidr` - (Optional) Switch to enable/disable advertise transit VPC network CIDR for a vgw connection.
* `bgp_manual_spoke_advertise_cidrs` - (Optional) Intended CIDR list to advertise to VGW. Example: "10.2.0.0/16,10.4.0.0/16".

-> **NOTE:** `enable_advertise_transit_cidr` - If you are using/upgraded to Aviatrix Terraform Provider R1.9+, and a vgw_conn resource was originally created with a provider version <R1.9, you must do ‘terraform refresh’ to update and apply the attribute’s default value (false) into the state file.

~> **NOTE:** `enable_advertise_transit_cidr` and `bgp_manual_spoke_advertise_cidrs` functionality has been migrated over to **aviatrix_transit_gateway** as of Aviatrix Terraform Provider R2.6. If you are using/upgraded to Aviatrix Terraform Provider R2.6+, and a **vgw_conn** resource was originally created with a provider version <R2.6, you must cut and paste these two arguments (and values) into the corresponding transit gateway resource referenced in this **vgw_conn**. A 'terraform refresh' will then successfully complete the migration and rectify the state file.


## Import

**vgw_conn** can be imported using the `conn_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_vgw_conn.test conn_name~vpc_id
```
