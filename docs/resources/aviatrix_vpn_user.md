---
subcategory: "OpenVPN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_user"
description: |-
  Creates and Manages Aviatrix VPN Users
---

# aviatrix_vpn_user

The **aviatrix_vpn_user** resource creates and manages Aviatrix VPN users.

~> **NOTE:** As of R2.15, management of user/profile attachment can be set using `manage_user_attachment`. This argument must be to *true* in either **aviatrix_vpn_user** or **aviatrix_vpn_profile**. If attachment is managed in the **aviatrix_vpn_user** (set to *true*), it must be set to *false* in the **aviatrix_vpn_profile** resource and vice versa.

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
```hcl
# Create an Aviatrix VPN User under Geo VPN
resource "aviatrix_vpn_user" "test_vpn_user" {
  dns_name   = "vpn.testuser.com"
  user_name  = "username1"
  user_email = "user@aviatrix.com"
}
```
```hcl
# Create an Aviatrix VPN User on GCP
# See note below about vpc_id formatting for GCP
resource "aviatrix_vpn_user" "test_vpn_user" {
  vpc_id     = "${aviatrix_vpc.test_vpc.vpc_id}~-~${aviatrix_account.test_account.gcloud_project_id}"
  gw_name    = "gw1"
  user_name  = "username1"
  user_email = "user@aviatrix.com"
}
```

## Argument Reference

The following arguments are supported:

### Required

~> **NOTE:** For GCP, the vpc_id must be in the form `vpc_id~-~gcloud_project_id`. For example, `"${aviatrix_vpc.test_vpc.vpc_id}~-~${aviatrix_account.test_account.gcloud_project_id}"`.

-> As of Provider version R2.21.2+, the `vpc_id` of an Oracle VCN has been changed from its name to its OCID.
* `vpc_id` - (Optional) VPC ID of Aviatrix VPN gateway. Used together with `gw_name`. Example: "vpc-abcd1234".
* `gw_name` - (Optional) If ELB is enabled, this will be the name of the ELB, else it will be the name of the Aviatrix VPN gateway. Used together with `vpc_id`. Example: "gw1".
* `dns_name` - (Optional) FQDN of a DNS based VPN service such as GeoVPN or UDP load balancer. Example: "vpn.testuser.com".
* `user_name` - (Required) VPN user name. Example: "user".
* `user_email` - (Optional) VPN user's email. Example: "abc@xyz.com".

### SAML
* `saml_endpoint` - (Optional) This is the name of the SAML endpoint to which the user is to be associated. This is required if adding user to a SAML gateway/LB.

### Misc.
* `manage_user_attachment` - (Optional) This parameter is a switch to determine whether or not to manage VPN user attachments to the VPN profile using this resource. If this is set to false, attachment must be managed using the **aviatrix_vpn_profile** resource. Valid values: true, false. Default value: false.
* `profiles` - (Optional) List of VPN profiles for user to attach to. This should be set to null if `manage_user_attachment` is set to false.


## Import

**vpn_user** can be imported using the `user_name`, e.g.

```
$ terraform import aviatrix_vpn_user.test user_name
```
