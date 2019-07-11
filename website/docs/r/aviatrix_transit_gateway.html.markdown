---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway"
sidebar_current: "docs-aviatrix-resource-transit_gateway"
description: |-
  Creates and Manages the Aviatrix Transit Network Gateways
---

# aviatrix_transit_gateway

The AviatrixTransitGateway resource creates and manages the Aviatrix Transit Network Gateways.

## Example Usage

```hcl
# Create an Aviatrix Transit Network Gateway in AWS
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
  cloud_type               = 1
  account_name             = "devops_aws"
  gw_name                  = "transit"
  vpc_id                   = "vpc-abcd1234"
  vpc_reg                  = "us-east-1"
  gw_size                  = "t2.micro"
  subnet                   = "10.1.0.0/24"
  ha_subnet                = "10.1.0.0/24"
  ha_gw_size               = "t2.micro"
  tag_list                 = ["name:value", "name1:value1", "name2:value2"]
  enable_hybrid_connection = true
  connected_transit        = "yes"
}

# Create an Aviatrix Transit Network Gateway in ARM
resource "aviatrix_transit_gateway" "test_transit_gateway_azure" {
  cloud_type               = 8
  account_name             = "devops_azure"
  gw_name                  = "transit"
  vpc_id                   = "vnet1:hello"
  vpc_reg                  = "West US"
  gw_size                  = "Standard_B1s"
  subnet                   = "10.30.0.0/24"
  ha_subnet                = "10.30.0.0/24"
  ha_gw_size               = "Standard_B1s"
  connected_transit        = "yes"
}

```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Use 1 for AWS.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Required if for aws. Example: AWS: "vpc-abcd1234", GCP: "mygooglecloudvpcname", etc...
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2", etc...
* `gw_size` - (Required) Size of the gateway instance.  Example: AWS: "t2.large", etc...
* `subnet` - (Required) Public Subnet CIDR. Copy/paste from AWS Console to get the right subnet CIDR. Example: AWS: "10.0.0.0/24".
* `ha_subnet` - (Optional) HA Subnet CIDR. Setting to empty/unset will disable HA. Setting to a valid subnet CIDR will create an HA gateway on the subnet. Example: "10.12.0.0/24".
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set). Example: "t2.micro".
* `enable_nat` - (Optional) Enable NAT for this container. Supported values: true, false.
* `tag_list` - (Optional) Instance tag of cloud provider. Only supported for aws. Example: ["key1:value1","key002:value002"].
* `enable_hybrid_connection` - (Optional) Sign of readiness for TGW connection. Only supported for aws. Example: false.
* `enable_firenet_interfaces` - (Optional) Sign of readiness for FireNet connection. Valid values: true, false. Default: false.
* `connected_transit` - (Optional) Specify Connected Transit status. Supported values: true, false.
* `insane_mode` - (Optional) Specify Insane Mode high performance gateway. Insane Mode gateway size must be at least c5 size. If enabled, will look for spare /26 segment to create a new subnet. (Only available for AWS.) Supported values: true, false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit Gateway. Required if insane_mode is enabled.
* `ha_insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit HA Gateway. Required if insane_mode is enabled and ha_subnet is set.

## Import

Instance transit_gateway can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_transit_gateway.test gw_name
```
