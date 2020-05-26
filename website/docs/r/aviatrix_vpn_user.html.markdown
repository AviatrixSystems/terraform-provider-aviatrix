---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_user"
description: |-
  Creates and Manages Aviatrix VPN Users
---

# aviatrix_vpn_user

The **aviatrix_vpn_user** resource creates and manages Aviatrix VPN users.

## Example Usage

```hcl
# Create an Aviatrix VPN User
resource "aviatrix_vpn_user" "test_vpn_user" {
  vpc_id     = "vpc-abcd1234"
  gw_name    = "gw1"
  user_name  = "username1"
  user_email = "user@aviatrix.com"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Optional) VPC ID of Aviatrix VPN gateway. Used together with `gw_name`. Example: "vpc-abcd1234".
* `gw_name` - (Optional) If ELB is enabled, this will be the name of the ELB, else it will be the name of the Aviatrix VPN gateway. Used together with `vpc_id`. Example: "gw1".
* `dna_name` - (Optional) If DNS is enabled, this will be the name of the DNS. Example: "vpn.testuser.com".
* `user_name` - (Required) VPN user name. Example: "user".
* `user_email` - (Optional) VPN user's email. Example: "abc@xyz.com".

### SAML
* `saml_endpoint` - (Optional) This is the name of the SAML endpoint to which the user is to be associated. This is required if adding user to a SAML gateway/LB.

## Import

**vpn_user** can be imported using the `user_name`, e.g.

```
$ terraform import aviatrix_vpn_user.test user_name
```
