---
subcategory: "Private Mode"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_private_mode_lb"
description: |-
  Creates and manages a Private Mode load balancer
---

# aviatrix_private_mode_lb

The **aviatrix_private_mode_lb** resource allows management of a Private Mode load balancer. This resource is available as of provider version R2.23+.

## Example Usage

```hcl
# Create a Private Mode load balancer
resource "aviatrix_private_mode_lb" "test" {
  // TODO
}
```


## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Name of the access account.
* `vpc_id` - (Required) ID of the VPC for the load balancer.
* `region` - (Required) Name of the region containing the VPC.
* `lb_type` - (Required) Type of load balancer to create. Must be one of controller or multicloud.
* `multicloud_access_vpc_id` - (Optional) ID of the VPC with a multicloud endpoint. Required when `lb_type` is multicloud.
* `proxies` - (Optional) List of multicloud proxies. Only valid when `lb_type` is multicloud.
  * `instance_id` - (Required) Instance ID of the proxy.
  * `proxy_type` - (Required) Type of load balancer. Must be one of controller or multicloud.

## Import

**aviatrix_private_mode_lb** can be imported using the `vpc_id`, e.g.

```
$ terraform import aviatrix_private_mode_lb.test vpc-1234567
```
