---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_global_vpc_excluded_instance"
description: |-
  Manages the instance to be excluded for Aviatrix global VPC tagging 
---

# aviatrix_global_vpc_excluded_instance

The **aviatrix_global_vpc_excluded_instance** resource manages the instance to be excluded for Aviatrix global VPC tagging.

## Example Usage

```hcl
# Add an Aviatrix Global VPC Excluded Instance
resource "aviatrix_global_vpc_excluded_instance" "test" {
  account_name  = "test-account"
  instance_name = "test-instance"
  region        = "us-west1"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Aviatrix GCP access account name.
* `instance_name` - (Required) Name of the instance to be excluded for tagging.
* `region` - (Required) Region of the instance.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `uuid` - UUID of the global VPC exclude list object.

## Import

**global_vpc_excluded_instance** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_global_vpc_excluded_instance.test uuid
```
