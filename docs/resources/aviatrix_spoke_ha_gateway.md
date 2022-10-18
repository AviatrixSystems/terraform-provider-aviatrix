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
  cloud_type      = 1
  primary_gw_name = "spoke-gw-aws"
  subnet          = "10.11.0.0/24"
}
```
```hcl
# Create an Aviatrix GCP Spoke HA Gateway
resource "aviatrix_spoke_ha_gateway" "test_spoke_ha_gcp" {
  cloud_type      = 4
  primary_gw_name = "spoke-gw-gcp"
  gw_name         = "spoke-gw-gcp-ha"
  zone            = "us-west1-b"
  gw_size         = "n1-standard-1"
  subnet          = "10.12.0.0/24"
}
```
```hcl
# Create an Aviatrix Azure Spoke HA Gateway
resource "aviatrix_spoke_ha_gateway" "test_spoke_ha_azure" {
  cloud_type      = 8
  primary_gw_name = "spoke-gw-azure"
  gw_name         = "spoke-gw-azure-ha"
  gw_size         = "Standard_B1ms"
  subnet          = "10.13.0.0/24"
}
```
```hcl
# Create an Aviatrix OCI Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_oracle" {
  cloud_type          = 16
  primary_gw_name     = "spoke-gw-oci"
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
* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently, only AWS(1), GCP(4), Azure(8), OCI(16), AzureGov(32), AWSGov(256), AWSChina(1024), AzureChina(2048), Alibaba Cloud(8192), AWS Top Secret(16384) and AWS Secret (32768) are supported.
* `primary_gw_name` - (Required) Name of the primary gateway which is already or will be created before this Spoke HA Gateway.
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane_mode).**
* `zone` - (Optional) Availability Zone. Required for GCP gateway. Example: "us-west1-c".
* `availability_domain` - (Optional) Availability domain. Required and valid only for OCI.
* `fault_domain` - (Optional) Fault domain. Required and valid only for OCI.

### Insane Mode
* `insane_mode` - (Optional) Enable [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) for Spoke HA Gateway. Insane Mode gateway size must be at least c5 size (AWS, AWSGov, AWS China, AWS Top Secret and AWS Secret) or Standard_D3_v2 (Azure and AzureGov); for GCP only four size are supported: "n1-highcpu-4", "n1-highcpu-8", "n1-highcpu-16" and "n1-highcpu-32". If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet for AWS, Azure, AzureGov, AWSGov, AWS Top Secret and AWS Secret. Only available for AWS, GCP/OCI, Azure, AzureGov, AzureChina, AWSGov, AWS Top Secret and AWS Secret. Valid values: true, false. Default value: false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Spoke HA Gateway. Required for AWS, AWSGov, AWS China, AWS Top Secret or AWS Secret if `insane_mode` is enabled. Example: AWS: "us-west-1a".

### Misc.
* `gw_name` - (Optional) Name of the Spoke HA Gateway which is going to be created. If not set, controller will auto generate a name for this gateway.
* `gw_size` - (Optional) Size of the Spoke HA Gateway instance. If not set, controller will use the same value as primary gateway's. Example: AWS/AWSGov/AWSChina: "t2.large", Azure/AzureGov/AzureChina: "Standard_B1s", OCI: "VM.Standard2.2", GCP: "n1-standard-1".
* `eip` - (Optional) Required when `allocate_new_eip` is false. It uses the specified EIP for this gateway. Available in Controller 4.7+. Only available for AWS, GCP, Azure, OCI, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

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
