---
subcategory: "Private Mode"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_private_mode_multicloud_endpoint"
description: |-
  Creates and manages a Private Mode multicloud endpoint
---

# aviatrix_private_mode_multicloud_endpoint

The **aviatrix_private_mode_multicloud_endpoint** resource allows management of a Private Mode multicloud endpoint. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Create an Aviatrix Controller Private Mode config
resource "aviatrix_private_mode_multicloud_endpoint" "test" {
  enable_private_mode = true
}
```


## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Name of the access account.
* `vpc_id` - (Required) ID of the VPC to create the endpoint in.
* `region` - (Required) Region of the VPC.
* `controller_lb_vpc_id` - (Required) ID of the VPC containing a Private Mode controller load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:
* `dns_entry` - DNS entry of the endpoint.

## Import

**aviatrix_private_mode_multicloud_endpoint** can be imported using the `vpc_id`, e.g.

```
$ terraform import aviatrix_private_mode_multicloud_endpoint.test vpc-1234567
```
