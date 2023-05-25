---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_copilot_fault_tolerant_deployment"
description: |-
  Creates CoPilot Fault Tolerant Deployment
---

# aviatrix_copilot_fault_tolerant_deployment

The **aviatrix_copilot_fault_tolerant_deployment** resource allows deployment of fault tolerant copilot.

## Example Usage

```hcl
# Deploy a Fault Tolerant CoPilot 
resource "aviatrix_copilot_fault_tolerant_deployment" "test" {
    cloud_type                          = 1
    account_name                        = "test-aws"
    region                              = "us-west-2"
    main_copilot_vpc_id                 = "vpc-012345"
    main_copilot_subnet                 = "10.0.1.0/24"
    controller_service_account_username = "admin"
    controller_service_account_password = "password"
    
    cluster_data_nodes {
        vpc_id = "vpc-012345"
        subnet = "10.0.1.0/24"
    }
    
    cluster_data_nodes {
        vpc_id = "vpc-012345"
        subnet = "10.0.1.0/24"
    }
    
    cluster_data_nodes {
        vpc_id = "vpc-012345"
        subnet = "10.0.1.0/24"
    }
}
```

## Argument Reference
The following arguments are supported:

### Required
* `cloud_type` - (Required) Cloud type, requires an integer value.
* `account_name` - (Required) Aviatrix access account name.
* `region` - (Required) Region name.
* `main_copilot_vpc_id` - (Required) Main copilot VPC ID.
* `main_copilot_subnet` - (Required) Main copilot subnet CIDR.
* `controller_service_account_username` - (Required) Controller service account username.
* `controller_service_account_password` - (Required) Controller service account password
* `cluster_data_nodes` - (Required) Cluster data nodes.
  * `vpc_id` - (Required) VPC ID.
  * `subnet` - (Required) Subnet CIDR.
  * `instance_size` - (Optional) Instance size. Default value: "t3.2xlarge".
  * `data_volume_size` - (Optional) Copilot data volume size, requires an integer value. Default value: 100.
  
### Optional
* `main_copilot_instance_size` - (Optional) Main copilot instance size. Default value: "t3.2xlarge".

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `main_copilot_public_id` - Main copilot public IP.
* `main_copilot_private_id` - Main copilot private IP.
