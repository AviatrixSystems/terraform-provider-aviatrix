---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_user_accelerator"
sidebar_current: "docs-aviatrix-resource-vpn_user_accelerator"
description: |-
  Manages the Aviatrix VPN User Accelerator
---

# aviatrix_vpn_user_accelerator

The AviatrixVPNUserAccelerator resource manages the VPN Users

## Example Usage

```hcl
# Manage the Aviatrix Vpn User Accelerator
resource "aviatrix_vpn_user_accelerator" "test_xlr" {
  elb_name = "Aviatrix-vpc-abcd2134"
}
```

## Argument Reference

The following arguments are supported:

* `elb_name` - (Required) Name of ELB to be added to VPN User Accelerator. Example: "Aviatrix-vpc-abcd2134", "Aviatrix-vpc-abcd1234"

## Import

```
$ terraform import aviatrix_vpn_user_acclerator.test Aviatrix-vpc-abcd1234
```
