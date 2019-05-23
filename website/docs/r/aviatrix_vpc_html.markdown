---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpc"
sidebar_current: "docs-aviatrix-resource-vpc"
description: |-
  Creates and manages an VPC.
---

# aviatrix_vpc

The vpc resource allows the creation and management of an VPC.

## Example Usage

```hcl
# Create a new VPC
resource "aviatrix_vpc" "test_vpc" {
  cloud_type           = 1
  account_name         = "devops"
  region               = "us-west-1"
  name                 = "vpcTest"
  cidr                 = "10.0.0.0/16"
  aviatrix_transit_vpc = false
  aviatrix_firenet_vpc = false
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Use 1 for AWS.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `name` - (Required) Name of the vpc which is going to be created.
* `region` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2", etc...
* `cidr` - (Required) VPC cidr.
* `aviatrix_transit_vpc` - (Optional) Specify whether it is an aviatrix transit vpc. (Supported values: true, false. Default: false)
* `aviatrix_firenet_vpc` - (Optional) Specify whether it is an aviatrix firenet vpc. (Supported values: true, false. Default: false)

-> **NOTE:** 

* `aviatrix_firenet_vpc` - If you are using/upgraded to Aviatrix Terraform Provider R1.8+/UserConnect-4.6 , and an vpc resource was originally created with a provider version < R1.8/UserConnect-4.6, you must do ‘terraform refresh’ to update and apply the attribute’s default value (“false”) into the state file.

## Import

Instance vpc can be imported using the name, e.g.

```
$ terraform import aviatrix_vpc.test name
```