---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway"
sidebar_current: "docs-aviatrix-resource-spoke_gateway"
description: |-
  Creates and Manages Aviatrix Spoke Gateways
---

# aviatrix_spoke_gateway

The aviatrix_spoke_gateway resource allows to create and manage Aviatrix Spoke Gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws" {
  cloud_type   = 1
  account_name = "my-aws"
  gw_name      = "spoke-gw-aws"
  vpc_id       = "vpc-abcd123~~spoke-vpc-01"
  vpc_reg      = "us-west-1"
  gw_size      = "t2.micro"
  subnet       = "10.11.0.0/24~~us-west-1b~~spoke-vpc-01-pubsub"
  enable_snat  = false
  dns_server   = "8.8.8.8"
  tag_list     = ["k1:v1","k2:v2"]
}

# Create an Aviatrix GCP Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_gcp" {
  cloud_type   = 4
  account_name = "my-gcp"
  gw_name      = "spoke-gw-gcp"
  vpc_id       = "gcp-spoke-vpc"
  vpc_reg      = "us-west1-b"
  gw_size      = "t2.micro"
  subnet       = "10.12.0.0/24"
  enable_snat  = false
}

# Create an Aviatrix ARM Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_arm" {
  cloud_type   = 8
  account_name = "my-arm"
  gw_name      = "spoke-gw-01"
  vpc_id       = "spoke:test-spoke-gw-123"
  vpc_reg      = "West US"
  gw_size      = "t2.micro"
  subnet       = "10.13.0.0/24"
  enable_snat  = false
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. AWS=1, GCP=4, ARM=8.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Required if cloud_type is "1" or "4". Example: AWS: "vpc-abcd1234". 
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", GCP: "us-west1-b", ARM: "East US 2".
* `gw_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", GCP: "f1.micro", ARM: "StandardD2".
* `subnet` - (Required) Public Subnet Info. Example: AWS: "172.31.0.0/20".
* `allocate_new_eip` - (Optional) When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in 4.7 or later release. Supported values: true, false. Default: true. Option not available for GCP and ARM gateways, they will automatically allocate new eip's.
* `eip` - (Optional) Required when allocate_new_eip is false. It uses specified EIP for this gateway. Available in 4.7 or later release.
* `ha_zone` - (Optional) HA Zone. Required for enabling HA for GCP gateway. Setting to empty/unset will disable HA. Setting to a valid zone will create an HA gateway in the zone. Example: "us-west1-c".
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set). Example: "t2.micro".
* `ha_eip` - (Optional) Public IP address that you want to assign to the HA peering instance. If no value is given, a new eip will automatically allocated. Only available for AWS.
* `enable_snat` - (Optional) Specify whether enabling Source NAT feature on the gateway or not. Please disable AWS NAT instance before enabling this feature. Supported values: true, false.
* `single_az_ha` (Optional) Set to true if this feature is desired. Supported values: true, false.
* `transit_gw` - (Optional) Specify the transit Gateway.
* `tag_list` - (Optional) Instance tag of cloud provider. Only AWS, cloud_type is "1", is supported. Example: ["key1:value1", "key2:value2"]. 
* `enable_active_mesh` - (Optional) Enable Active Mesh Mode for Spoke Gateway. Valid values: true, false. Default value: false.

## Import

Instance spoke_gateway can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_spoke_gateway.test gw_name
```
