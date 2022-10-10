---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_microseg_intra_vpc"
description: |-
  Creates and manages an Aviatrix Micro-segmentation Intra VPC List
---

# aviatrix_microseg_intra_vpc

The **aviatrix_microseg_intra_vpc** resource handles the creation and management of Micro-segmentation Intra VPCs. Available as of Provider R3.0.0+.

## Example Usage

```hcl
# Create an Aviatrix Microseg Policy
resource "aviatrix_microseg_policy_list" "test" {
  vpcs {
    account_name = "azure-account"
    vpc_id       = "azure-vpc-0:rg-av-azure-vpc-0-808200:8168668b-a646-45b9-b88b-d756e60cf130"
    region       = "Central US"
  }
  
  vpcs {
    account_name = "azure-account"
    vpc_id       = "azure-vpc-1:rg-av-azure-vpc-1-562104:622a2277-5c57-4149-bcb9-c00d9284ee18"
    region       = "Central US"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `vpcs` - (Required) List of micro-seg enabled intra VPCs.
    * `account_name` - (Required) Account Name of the VPC.
    * `vpc_id` - (Required) vpc_id of the VPC.
    * `region` - (Required) Region of the VPC..

## Import

**aviatrix_microseg_intra_vpc** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_microseg_intra_vpc.test 10-11-12-13
```
