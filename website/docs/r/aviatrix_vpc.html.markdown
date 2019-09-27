---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpc"
description: |-
  Creates and manages VPCs.
---

# aviatrix_vpc

The aviatrix_vpc resource allows the creation and management of VPCs.

## Example Usage

```hcl
# Create an AWS VPC
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

* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently only AWS(1) is supported.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `name` - (Required) Name of the vpc which is going to be created.
* `region` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2".
* `cidr` - (Required) VPC cidr.
* `aviatrix_transit_vpc` - (Optional) Specify whether it is an Aviatrix Transit VPC. Supported values: true, false. Default: false.
* `aviatrix_firenet_vpc` - (Optional) Specify whether it is an Aviatrix Firenet VPC. Supported values: true, false. Default: false.

-> **NOTE:** `aviatrix_firenet_vpc` - If you are using/upgraded to Aviatrix Terraform Provider R1.8+, and an VPC resource was originally created with a provider version < R1.8, you must do ‘terraform refresh’ to update and apply the attribute’s default value (false) into the state file.

## Import

Instance vpc can be imported using the name, e.g.

```
$ terraform import aviatrix_vpc.test name
```
