---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_copilot_security_group_management"
description: |-
  Creates and manages CoPilot Security Group Management 
---

# aviatrix_copilot_security_group_management

The **aviatrix_copilot_security_group_management** resource allows management of controller CoPilot Security Group Management. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Enable the CoPilot Security Group Management
resource "aviatrix_copilot_security_group_management" "test" {
    cloud_type   = 1
    account_name = "aws-account"
    region       = "us-east-1"
    vpc_id       = "vpc-1234567890"
    instance_id  = "i-1234567890"
}
```


## Argument Reference

The following arguments are supported:

### Required
* `cloud_type` - (Required) Cloud type. The type of this attribute is Integer. Only support AWS, Azure and OCI.
* `account_name` - (Required) Aviatrix access account name.
* `vpc_id` - (Required) VPC ID.
* `instance_id` - (Required) CoPilot instance ID.

### Optional
* `region` - (Optional) Region where CoPilot is deployed. Required and valid for AWS and Azure.
* `zone` - (Optional) Zone where CoPilot is deployed. Required and valid for GCP.

## Import

**aviatrix_copilot_security_group_management** can be imported using instance ID, e.g.

```
$ terraform import aviatrix_copilot_security_group_management.test i-1234567890
```
