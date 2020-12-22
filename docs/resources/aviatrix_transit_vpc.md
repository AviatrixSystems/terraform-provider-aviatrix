---
subcategory: "Deprecated"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_vpc"
description: |-
  Creates and Manages the Aviatrix Transit Network Gateways
---

# aviatrix_transit_vpc

The aviatrix_transit_vpc resource creates and manages the Aviatrix Transit Network Gateways.

!> **WARNING:** The `aviatrix_transit_vpc` resource is deprecated as of **Release 2.0**. It is currently kept for backward-compatibility and will be removed in the future. Please use the transit gateway resource instead. If this is already in the state, please remove it from state file and import as `aviatrix_transit_gateway`.

## Example Usage

```hcl
# Manage Aviatrix Transit Network Gateways in aws
resource "aviatrix_transit_vpc" "test_transit_gw_aws" {
  cloud_type               = 1
  account_name             = "devops_aws"
  gw_name                  = "transit"
  vpc_id                   = "vpc-abcd1234"
  vpc_reg                  = "us-east-1"
  vpc_size                 = "t2.micro"
  subnet                   = "10.1.0.0/24"
  ha_subnet                = "10.1.0.0/24"
  ha_gw_size               = "t2.micro"
  tag_list                 = [
    "name:value",
    "name1:value1",
    "name2:value2"
  ]
  enable_hybrid_connection = true
  connected_transit        = "yes"
}

# Manage Aviatrix Transit Network Gateways in azure
resource "aviatrix_transit_vpc" "test_transit_gw_azure" {
  cloud_type        = 8
  account_name      = "devops_azure"
  gw_name           = "transit"
  vpc_id            = "vnet1:hello"
  vpc_reg           = "West US"
  vpc_size          = "Standard_B1s"
  subnet            = "10.30.0.0/24"
  ha_subnet         = "10.30.0.0/24"
  ha_gw_size        = "Standard_B1s"
  connected_transit = "yes"
}

```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Use 1 for AWS.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Required if for aws. Example: AWS: "vpc-abcd1234", GCP: "mygooglecloudvpcname", etc...
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2", etc...
* `vpc_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", etc...
* `subnet` - (Required) Public Subnet CIDR. Example: AWS: "10.0.0.0/24". Copy/paste from AWS Console to get the right subnet CIDR.
* `ha_subnet` - (Optional) HA Subnet CIDR. Example: "10.12.0.0/24".Setting to empty/unset will disable HA. Setting to a valid subnet CIDR will create an HA gateway on the subnet.
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set). Example: "t2.micro".
* `enable_snat` - (Optional) Enable Source NAT for this container. Supported values: true, false. Default value: false.
* `tag_list` - (Optional) Instance tag of cloud provider. Only supported for aws. Example: ["key1:value1","key002:value002"]
* `enable_hybrid_connection` - (Optional) Sign of readiness for TGW connection. Only supported for aws. Example: false.
* `enable_firenet_interfaces` - (Optional) Sign of readiness for FireNet connection. Valid values: true and false. Default: false.
* `connected_transit` - (Optional) Specify Connected Transit status. Supported values: true, false.
* `insane_mode` - (Optional) Specify Insane Mode high performance gateway. Insane Mode gateway size must be at least c5 size. If enabled, will look for spare /26 segment to create a new subnet. Only available for AWS. Supported values: true, false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit Gateway. Required if insane_mode is enabled.
* `ha_insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit HA Gateway. Required if insane_mode is enabled and ha_subnet is set.

The following arguments are deprecated:

* `dns_server` - Specify the DNS IP, only required while using a custom private DNS for the VPC.
* `vnet_name_resource_group` - (Optional) VPC-ID/VNet-Name of cloud provider. Required if for azure. ARM: "VNet_Name:Resource_Group_Name". It is replaced by "vpc_id".

-> **NOTE:** `enable_firenet_interfaces` - If you are using/upgraded to Aviatrix Terraform Provider R1.8+, and a transit_vpc resource was originally created with a provider version < R1.8, you must do ‘terraform refresh’ to update and apply the attribute’s default value (false) into the state file.

-> **NOTE:** `vnet_name_resource_group` - If you are using/upgraded to Aviatrix Terraform Provider R1.10+, and an ARM transit_vpc resource was originally created with a provider version < R1.10, you must replace "vnet_name_resource_group" with "vpc_id" in your configuration file, and do ‘terraform refresh’ to set its value to "vpc_id" and apply it into the state file.

## Import

Instance transit_vpc can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_transit_vpc.test gw_name
```
