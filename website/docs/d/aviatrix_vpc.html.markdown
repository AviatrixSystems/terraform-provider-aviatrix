---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpc"
description: |-
  Gets the Aviatrix Vpc.
---

# aviatrix_vpc

Use this data source to get the Aviatrix vpc for use in other resources.

## Example Usage

```hcl
# Aviatrix Vpc Data Source
data "aviatrix_vpc" "test" {
  name = "vpc-test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Aviatrix VPC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cloud_type` - Type of cloud service provider.
* `account_name` - Account name of the VPC created.
* `region` - Region of the VPC created.
* `cidr` - Subnet of the VPC created.
* `region` - Region of the VPC created.
* `aviatrix_transit_vpc` - Switch if the VPC created is an Aviatrix Transit VPC or not.
* `aviatrix_firenet_vpc` - Switch if the VPC created is an Aviatrix FireNet VPC or not.
* `vpc_id` - ID of the VPC created.
* `subnets` - List of subnet of the VPC created.
  * `cidr` - Subnet cidr.
  * `name` - Subnet name.
  * `subnet_id` - Subnet ID.
* `public_subnets` - List of public subnet of the VPC created.
  * `cidr` - Public subnet cidr.
  * `name` - Public subnet name.
  * `subnet_id` - Public subnet ID.
* `private_subnets` - List of private subnet of the VPC created.
  * `cidr` - Private subnet cidr.
  * `name` - Private subnet name.
  * `subnet_id` - Private subnet ID.