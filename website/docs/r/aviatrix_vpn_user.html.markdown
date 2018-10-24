---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_user"
sidebar_current: "docs-aviatrix-resource-vpn_user"
description: |-
  Manages the Aviatrix VPN Users
---

# aviatrix_vpn_user

The AviatrixVPNUser resource manages the VPN Users

## Example Usage

```hcl
# Manage Aviatrix Controller Upgrade process
resource "aviatrix_vpn_user" "test_vpn_user" {
  vpc_id = "vpc-abcd1234"
  gw_name = "gw1"
  user_name = "username1"
  user_email = "user@aviatrix.com"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) VPC Id of Aviatrix VPN gateway. Example: "vpc-abcd1234"
* `gw_name` - (Required) If ELB is enabled, this will be the name of the ELB, else it will be the name of the Aviatrix VPN gateway. Example: "gw1"
* `user_name` - (Required) VPN user name. Example: "user"
* `user_email` - (Required) VPN User's email. Example: "abc@xyz.com"
