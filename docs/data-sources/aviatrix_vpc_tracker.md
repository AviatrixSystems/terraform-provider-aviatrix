---
subcategory: "Useful Tools"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpc_tracker"
description: |-
  Gets the Aviatrix VPC Tracker information.
---

# aviatrix_vpc_tracker

Use this data source to get the list of VPC's for use in other resources.

## Example Usage

```hcl
# Aviatrix VPC Tracker Data Source
data "aviatrix_vpc_tracker" "foo" {
  cloud_type   = 1
  cidr         = "10.0.0.1/24"
  region       = "us-west-1"
  account_name = "bar"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Optional) Filters VPC list by cloud provider id. For example, cloud_type = 1 will give all AWS VPCs.
* `cidr` - (Optional) Filters VPC list by CIDR (AWS/AZURE only).
* `region` - (Optional) Filters VPC list by region (AWS/AZURE only).
* `account_name` - (Optional) Filters VPC list by access account name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_list` - List of VPCs from the VPC tracker.
  * `cloud_type` - Cloud provider id hosting this VPC.
  * `vpc_id` - VPC id.
  * `account_name` - Aviatrix access account associated with the VPC.
  * `region` - VPC region (AWS/AZURE only).
  * `name` - VPC name.
  * `cidr` - VPC CIDR (AWS/AZURE only).
  * `instance_count` - Number of running instances in the VPC.
  * `subnets` - List of subnets within this VPC (GCP only).
    * `region` - Subnet region.
    * `name` - Subnet name.
    * `cidr` - Subnet CIDR.
    * `gw_ip` - Subnet gateway ip.

## Notes

* Please be aware this data source could take up to 20 minutes to refresh depending on the number of VPCs and cloud accounts.
