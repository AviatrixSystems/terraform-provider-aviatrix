---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_copilot_simple_deployment"
description: |-
  Creates CoPilot Simple Deployment
---

# aviatrix_copilot_simple_deployment

The **aviatrix_copilot_simple_deployment** resource allows deployment of simple copilot.

## Example Usage

```hcl
# Deploy a Simple CoPilot 
resource "aviatrix_copilot_simple_deployment" "test" {
    cloud_type                          = 1
    account_name                        = "test-aws"
    region                              = "us-west-1"
    vpc_id                              = "vpc-012345"
    subnet                              = "10.0.1.0/24"
    controller_service_account_username = "admin"
    controller_service_account_password = "password"
}
```

## Argument Reference
The following arguments are supported:

### Required
* `cloud_type` - (Required) Cloud type, requires an integer value.
* `account_name` - (Required) Aviatrix access account name.
* `region` - (Required) Region name.
* `vpc_id` - (Required) VPC ID.
* `subnet` - (Required) Subnet CIDR.
* `controller_service_account_username` - (Required) Controller service account username.
* `controller_service_account_password` - (Required) Controller service account password

### Optional
* `instance_size` - (Optional) Copilot instance size. Default value: "t3.2xlarge".
* `data_volome_size` - (Optional) Copilot data volume size, requires an integer value. Default value: 100.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `public_id` - Copilot public IP.
* `private_id` - Copilot private IP.
