---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_vpc"
sidebar_current: "docs-aviatrix-resource-spoke_vpc"
description: |-
  Sets Aviatrix Spoke Gateway
---

# aviatrix_spoke_vpc

The SpokeVpc resource allows to manage Aviatrix Spoke Gateway

## Example Usage

```hcl
# Set Aviatrix aws spoke_vpc
resource "aviatrix_spoke_vpc" "test_spoke_vpc_aws" {
  cloud_type   = 1
  account_name = "my-aws"
  gw_name      = "spoke-gw-aws"
  vpc_id       = "vpc-abcd123~~spoke-vpc-01"
  vpc_reg      = "us-west-1"
  vpc_size     = "t2.micro"
  subnet       = "10.11.0.0/24~~us-west-1b~~spoke-vpc-01-pubsub"
  enable_nat   = "no"
  dns_server   = "8.8.8.8"
  tag_list     = ["k1:v1","k2:v2"]
}

# Set Aviatrix gcp spoke_vpc
resource "aviatrix_spoke_vpc" "test_spoke_vpc_gcp" {
  cloud_type   = 4
  account_name = "my-gcp"
  gw_name      = "spoke-gw-gcp"
  vpc_id       = "gcp-spoke-vpc"
  vpc_reg      = "us-west1-b"
  vpc_size     = "t2.micro"
  subnet       = "10.12.0.0/24"
  enable_nat   = "no"
}

# Set Aviatrix arm spoke_vpc
resource "aviatrix_spoke_vpc" "test_spoke_vpc_arm" {
  cloud_type   = 1
  account_name = "my-aws"
  gw_name      = "spoke-gw-01"
  vpc_id       = "spoke:test-spoke-gw-123"
  vpc_reg      = "West US"
  vpc_size     = "t2.micro"
  subnet       = "10.13.0.0/24"
  enable_nat   = "no"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. AWS=1, GCP=4, ARM=8
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.
* `vpc_id` - (Optional) VPC-ID/VNet-Name of cloud provider. Example: AWS: "vpc-abcd1234", etc... (Required if cloud_type is "1" or "4")
* `vnet_and_resource_group_names` - (Optional) The string consisted of name of (Azure) VNet and name Resource-Group. Valid Value(s): Refer to Aviatrix controller GUI. (Required if cloud_type is "8")
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", GCP: "us-west1-b", ARM: "East US 2", etc...
* `vpc_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", GCP: "f1.micro", ARM: "StandardD2", etc...
* `subnet` - (Required) Public Subnet Info. Example: AWS: "CIDR~~ZONE~~SubnetName", etc...
* `ha_subnet` - (Optional) HA Subnet. Setting to empty/unset will disable HA. Setting to a valid subnet (Example: 10.12.0.0/24) will create an HA gateway on the subnet. If enabling HA for a GCP gateway, enter a valid zone.
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set)(Example: "t2.micro")
* `enable_nat` - (Optional) Specify whether enabling NAT feature on the gateway or not. (Please disable AWS NAT instance before enabling this feature) Example: "yes", "no"
* `single_az_ha` (Optional) Set to "enabled" if this feature is desired.
* `transit_gw` - (Optional)  Specify the transit Gateway.
* `tag_list` - (Optional) Instance tag of cloud provider. Example: key1:value1,key002:value002, etc... Only AWS (cloud_type is "1") is supported

-> **NOTE:** The following arguments are deprecated:

* `dns_server` - Specify the DNS IP, only required while using a custom private DNS for the VPC.

## Import

Instance spoke_vpc can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_spoke_vpc.test gw_name
```
