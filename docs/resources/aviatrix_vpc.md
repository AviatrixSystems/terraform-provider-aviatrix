---
subcategory: "Useful Tools"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpc"
description: |-
  Creates and manages Aviatrix-created VPCs
---

# aviatrix_vpc

The **aviatrix_vpc** resource allows the creation and management of Aviatrix-created VPCs of various cloud types.

## Example Usage

```hcl
# Create an AWS VPC
resource "aviatrix_vpc" "aws_vpc" {
  cloud_type           = 1
  account_name         = "devops"
  region               = "us-west-1"
  name                 = "aws-vpc"
  cidr                 = "10.0.0.0/16"
  aviatrix_transit_vpc = false
  aviatrix_firenet_vpc = false
}
```
```hcl
# Create a GCP VPC
resource "aviatrix_vpc" "gcp_vpc" {
  cloud_type           = 4
  account_name         = "devops"
  name                 = "gcp-vpc"

  subnets {
    name   = "subnet-1"
    region = "us-west1"
    cidr   = "10.10.0.0/24"
  }

  subnets {
    name   = "subnet-2"
    region = "us-west2"
    cidr  = "10.11.0.0/24"
  }
}
```
```hcl
# Create an Azure VNet
resource "aviatrix_vpc" "azure_vnet" {
  cloud_type           = 8
  account_name         = "devops"
  region               = "Central US"
  name                 = "azure-vnet"
  cidr                 = "12.0.0.0/16"
  aviatrix_firenet_vpc = false
}
```
```hcl
# Create an OCI VPC
resource "aviatrix_vpc" "oci_vpc" {
  cloud_type   = 16
  account_name = "devops"
  region       = "us-ashburn-1"
  name         = "oci-vpc"
  cidr         = "10.0.0.0/24"
}
```
```hcl
# Create an AWSGov VPC
resource "aviatrix_vpc" "awsgov_vnet" {
  cloud_type           = 256
  account_name         = "devops"
  region               = "us-gov-west-1"
  name                 = "awsgov-vpc"
  cidr                 = "12.0.0.0/20"
  aviatrix_transit_vpc = false
  aviatrix_firenet_vpc = false
}
```

## Argument Reference

The following arguments are supported:

### Required
* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently only AWS(1), GCP(4), AZURE(8), OCI(16) and AWSGov(256) are supported.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `name` - (Required) Name of the VPC to be created.
* `region` - (Optional) Region of cloud provider. **Required to be empty for GCP provider, and non-empty for other providers.** Example: AWS: "us-east-1", AZURE: "East US 2", OCI: "us-ashburn-1", AWSGov: "us-gov-east-1".
* `cidr` - (Optional) VPC CIDR. **Required to be empty for GCP provider, and non-empty for other providers.** Example: "10.11.0.0/24".

### Google Cloud
* `subnets` - (Optional) List of subnets to be specify for GCP provider. Required to be non-empty for GCP provider, and empty for other providers.
  * `region` - Region of this subnet.
  * `cidr` - CIDR block.
  * `name` - Name of this subnet.

### Advanced Options
* `subnet_size` - (Optional) Subnet size. Only supported for AWS, Azure provider. Example: 24. Available in provider version R2.17+.
* `num_of_subnet_pairs` - (Optional) Number of public subnet and private subnet pair created. Only supported for AWS, Azure provider. Example: 1. Available in provider version R2.17+.

### Misc.
* `aviatrix_transit_vpc` - (Optional) Specify whether it is an [Aviatrix Transit VPC](https://docs.aviatrix.com/HowTos/create_vpc.html#aviatrix-transit-vpc) to be used for [Transit Network](https://docs.aviatrix.com/HowTos/transitvpc_faq.html) or [TGW](https://docs.aviatrix.com/HowTos/tgw_faq.html) solutions. **Only AWS and AWSGov are supported. Required to be false for other providers.** Valid values: true, false. Default: false.
* `aviatrix_firenet_vpc` - (Optional) Specify whether it is an Aviatrix FireNet VPC to be used for [Aviatrix FireNet](https://docs.aviatrix.com/HowTos/firewall_network_faq.html) and [Transit FireNet](https://docs.aviatrix.com/HowTos/transit_firenet_faq.html) solutions. **Only AWS, AWSGov and Azure are supported. Required to be false for other providers.** Valid values: true, false. Default: false.
* `enable_native_gwlb` - (Optional) Enable Native AWS GWLB for FireNet Function. Only valid with cloud_type = 1 (AWS). Valid values: true or false. Default value: false. Available as of provider version R2.18+.


-> **NOTE:** `aviatrix_firenet_vpc` - If you are using/ upgraded to Aviatrix Terraform Provider R1.8+, and a VPC resource was originally created with a provider version <R1.8, you must do 'terraform refresh' to update and apply the attribute’s default value (false) into the state file.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_id` - ID of the vpc to be created.
* `resource_group` - Resource group of the Azure VPC created.
* `azure_vnet_resource_id` - Azure vnet resource ID.  
* `subnets` - List of subnet of the VPC to be created.
  * `cidr` - CIDR block.
  * `name` - Name of this subnet.
  * `subnet_id` - ID of this subnet.
* `private_subnets` - List of private subnet of the VPC(AWS, Azure) to be created.
  * `cidr` - CIDR block.
  * `name` - Name of this subnet.
  * `subnet_id` - ID of this subnet.
* `public_subnets` - List of public subnet of the VPC(AWS, Azure) to be created.
  * `cidr` - CIDR block.
  * `name` - Name of this subnet.
  * `subnet_id` - ID of this subnet.  

-> **NOTE:** `subnets` - If created as a FireNet VPC, four public subnets will be created in the following order: subnet for firewall-mgmt in the first zone, subnet for ingress-egress in the first zone, subnet for firewall-mgmt in the second zone, and subnet for ingress-egress in the second zone.

## Import

**vpc** can be imported using the VPC's `name`, e.g.

```
$ terraform import aviatrix_vpc.test name
```
