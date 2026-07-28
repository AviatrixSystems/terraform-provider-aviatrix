---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_ha_gateway"
description: |-
  Creates and manages Aviatrix spoke ha gateways
---

# aviatrix_spoke_ha_gateway

The **aviatrix_spoke_ha_gateway** resource allows the creation and management of Aviatrix spoke ha gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Spoke HA Gateway
resource "aviatrix_spoke_ha_gateway" "test_spoke_ha_aws" {
  primary_gw_name = aviatrix_spoke_gateway.primary_spoke.id
  subnet          = "10.11.0.0/24"
}
```
```hcl
# Create an Aviatrix GCP Spoke HA Gateway
resource "aviatrix_spoke_ha_gateway" "test_spoke_ha_gcp" {
  primary_gw_name = aviatrix_spoke_gateway.primary_spoke.id
  gw_name         = "spoke-gw-gcp-ha"
  zone            = "us-west1-b"
  gw_size         = "n1-standard-1"
  subnet          = "10.12.0.0/24"
}
```
```hcl
# Create an Aviatrix Azure Spoke HA Gateway
resource "aviatrix_spoke_ha_gateway" "test_spoke_ha_azure" {
  primary_gw_name = aviatrix_spoke_gateway.primary_spoke.id
  gw_name         = "spoke-gw-azure-ha"
  gw_size         = "Standard_B1ms"
  subnet          = "10.13.0.0/24"
}
```
```hcl
# Create an Aviatrix OCI Spoke HA Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_oracle" {
  primary_gw_name     = aviatrix_spoke_gateway.primary_spoke.id
  gw_name             = "spoke-gw-oci-ha"
  gw_size             = "VM.Standard2.2"
  subnet              = "10.7.0.0/16"
  availability_domain = aviatrix_vpc.oci_vpc.availability_domains[0]
  fault_domain        = aviatrix_vpc.oci_vpc.fault_domains[0]
}
```


## Argument Reference

The following arguments are supported:

### Required
* `primary_gw_name` - (Required) Name of the primary gateway which is already or will be created before this Spoke HA Gateway.
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane_mode).**
* `zone` - (Optional) Availability Zone. Required for GCP gateway, example: "us-west1-c". Optional for Azure gateway in the form "az-n", example: "az-2".
* `availability_domain` - (Optional) Availability domain. Required and valid only for OCI.
* `fault_domain` - (Optional) Fault domain. Required and valid only for OCI.

### Insane Mode
* `insane_mode` - (Optional) Enable [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) for Spoke HA Gateway. Insane Mode gateway size must be at least c5 size (AWS, AWSGov, AWS China, AWS Top Secret and AWS Secret) or Standard_D3_v2 (Azure and AzureGov); for GCP only four size are supported: "n1-highcpu-4", "n1-highcpu-8", "n1-highcpu-16" and "n1-highcpu-32". If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet for AWS, Azure, AzureGov, AWSGov, AWS Top Secret and AWS Secret. Only available for AWS, GCP/OCI, Azure, AzureGov, AzureChina, AWSGov, AWS Top Secret and AWS Secret. Valid values: true, false. Default value: false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Spoke HA Gateway. Required for AWS, AWSGov, AWS China, AWS Top Secret or AWS Secret if `insane_mode` is enabled. Example: AWS: "us-west-1a".

### Misc.
* `gw_name` - (Optional) Name of the Spoke HA Gateway which is going to be created. If not set, controller will auto generate a name for this gateway.
* `gw_size` - (Optional) Size of the Spoke HA Gateway instance. If not set, controller will use the same value as primary gateway's. Example: AWS/AWSGov/AWSChina: "t2.large", Azure/AzureGov/AzureChina: "Standard_B1s", OCI: "VM.Standard2.2", GCP: "n1-standard-1".
* `eip` - (Optional) If set, the set IP will be used for this gateway.
* `azure_eip_name_resource_group` - (Optional) Name of public IP Address resource and its resource group in Azure to be assigned to the Spoke Gateway instance. Example: "IP_Name:Resource_Group_Name". Required if `eip` is set and `cloud_type` is Azure, AzureGov or AzureChina. Available as of provider version 3.0+.
* `single_az_ha` - (Optional) Set to true if this [feature](https://docs.aviatrix.com/Solutions/gateway_ha.html#single-az-gateway) is desired. Valid values: true, false. If not set, the value is derived from primary gateway during launch. After launch, the setting can be edited using this parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cloud_type` - Type of cloud service provider.
* `account_name` - Name of a Cloud-Account in Aviatrix controller.
* `software_version` - Software version of the gateway.
* `image_version` - Image version of the gateway.
* `vpc_reg` - Region in which the Spoke HA Gateway was created.
* `security_group_id` - Security group used for the Spoke HA Gateway.
* `cloud_instance_id` - Cloud instance ID of the Spoke HA Gateway.
* `private_ip` - Private IP address of the Spoke HA Gateway created.
* `public_ip` - Public IP address of the Spoke HA Gateway created.

## Import

**spoke_ha_gateway** can be imported using the `gw_name`, e.g.
****
```
$ terraform import aviatrix_spoke_ha_gateway.test gw_name
```
