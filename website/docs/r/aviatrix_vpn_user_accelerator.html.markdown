---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_user_accelerator"
description: |-
  Manages the Aviatrix VPN User Accelerator
---

# aviatrix_vpn_user_accelerator

The aviatrix_vpn_user_accelerator resource manages the Aviatrix VPN User Accelerator.

## Example Usage

```hcl
# Create an Aviatrix Vpn User Accelerator
resource "aviatrix_vpn_user_accelerator" "test_vpc_accelerator" {
  elb_name = "Aviatrix-vpc-abcd2134"
}
```

## Argument Reference

The following arguments are supported:

* `elb_name` - (Required) Name of ELB to be added to VPN User Accelerator. Example: "Aviatrix-vpc-abcd2134".

## Import

```
$ terraform import aviatrix_vpn_user_acclerator.test Aviatrix-vpc-abcd1234
```
