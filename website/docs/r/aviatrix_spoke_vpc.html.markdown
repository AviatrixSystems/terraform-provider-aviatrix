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
# Set Aviatrix spoke_vpc
resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
  cloud_type = 1
  account_name = my-aws
  gw_name = spoke-gw-01
  vpc_id = vpc-abcd123~~spoke-vpc-01
  vnet_and_resource_group_names =
  vpc_reg = us-west-1
  vpc_size = t2.micro
  subnet = 10.11.0.0/24~~us-west-1b~~spoke-vpc-01-pubsub
  enable_nat = no
  dns_server = 8.8.8.8
  tag_list = k1:v1,k2:v2
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. AWS=1
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller
* `gw_name` - (Required) Name of the gateway which is going to be created.
* `vpc_id` - (Optional) VPC-ID/VNet-Name of cloud provider. Example: AWS: "vpc-abcd1234", etc...
* `vnet_and_resource_group_names` - (Optional) The string consisted of name of (Azure) VNet and name Resource-Group. Valid Value(s): Refer to Aviatrix controller GUI.          (Required if cloud_type is "8")
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2", etc...
* `vpc_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", etc...
* `subnet` - (Required) Public Subnet Info. Example: AWS: "CIDR~~ZONE~~SubnetName", etc...
* `ha_subnet` - (Optional) HA Subnet
* `enable_nat` - (Optional) Specify whether enabling NAT feature on the gateway or not. (Please disable AWS NAT instance before enabling this feature) Example: "yes", "no"
* `dns_server` - (Optional) Specify the DNS IP
* `transit_gw` - (Optional)  Specify the transit Gateway
* `tag_list` - (Optional) Instance tag of cloud provider. Example: key1:value1,key002:value002, etc...
