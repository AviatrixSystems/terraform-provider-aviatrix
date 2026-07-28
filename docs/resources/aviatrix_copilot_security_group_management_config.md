---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_copilot_security_group_management_config"
description: |-
  Creates and manages CoPilot Security Group Management Configuration
---

# aviatrix_copilot_security_group_management_config

The **aviatrix_copilot_security_group_management_config** resource allows management of controller CoPilot security group management configuration. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Enable the CoPilot Security Group Management when CoPilot runs in AWS
resource "aviatrix_copilot_security_group_management_config" "test" {
    enable_copilot_security_group_management = true
    cloud_type                               = 1
    account_name                             = "aws-account"
    region                                   = "us-east-1"
    vpc_id                                   = "vpc-1234567890"
    instance_id                              = "i-1234567890"
}
```
```hcl
# Enable the CoPilot Security Group Management  when CoPilot runs in Azure
resource "aviatrix_copilot_security_group_management_config" "test" {
    enable_copilot_security_group_management = true
    cloud_type                               = 8
    account_name                             = "azure-account"
    region                                   = "West Europe"
    vpc_id                                   = "shared-vnet:shared-rg:1234567-123c-42a1-b3d7-e1234567890"
    instance_id                              = "copilot-vm:SHARED-RG"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `enable_copilot_security_group_management` - (Required) Switch to enable copilot security group management. Valid values: true, false.

### Optional
* `cloud_type` - (Optional) Cloud type. The type of this attribute is Integer. Only support AWS, Azure and OCI. Required to enable copilot security group management.
* `account_name` - (Optional) Aviatrix access account name. Required to enable copilot security group management.
* `vpc_id` - (Optional) VPC ID. Required to enable copilot security group management.
* `instance_id` - (Optional) CoPilot instance ID. Required to enable copilot security group management.
* `region` - (Optional) Region where CoPilot is deployed. Required and valid for AWS and Azure.
* `zone` - (Optional) Zone where CoPilot is deployed. Required and valid for GCP.

## Import

**aviatrix_copilot_security_group_management_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_copilot_security_group_management_config.test 10-11-12-13
```
