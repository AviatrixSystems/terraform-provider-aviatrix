---
subcategory: "OpenVPN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_profile"
description: |-
  Creates and manages Aviatrix VPN User Profiles
---

# aviatrix_vpn_profile

The **aviatrix_vpn_profile** resource allows the creation and management of Aviatrix VPN user profiles.

~> **NOTE:** As of R2.15, management of user/profile attachment can be set using `manage_user_attachment`. This argument must be set to *true* in either **aviatrix_vpn_user** or **aviatrix_vpn_profile**. If attachment is managed in the **aviatrix_vpn_profile** (set to *true*), it must be set to *false* in the **aviatrix_vpn_user** resource and vice versa.

## Example Usage

```hcl
# Create an Aviatrix AWS VPN User Profile
resource "aviatrix_vpn_profile" "test_vpn_profile" {
  name      = "my_profile"
  base_rule = "allow_all"
  users     = [
    "user1",
    "user2"
  ]

  policy {
    action = "deny"
    proto  = "tcp"
    port   = "443"
    target = "10.0.0.0/32"
  }

  policy {
    action = "deny"
    proto  = "tcp"
    port   = "443"
    target = "10.0.0.1/32"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Enter any name for the VPN profile.
* `base_rule` - (Optional) Base policy rule of the profile to be added. Enter "allow_all" or "deny_all", based on whether you want a whitelist or blacklist.

### Policy Options
* `policy` - (Optional) New security policy for the profile. Each policy has the following attributes:
  * `action` - (Required) Should be the opposite of the base rule for correct behavior. Valid values for action: "allow", "deny".
  * `proto` - (Required) Protocol to allow or deny. Valid values for protocol: "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp".
  * `port` - (Required) Port to be allowed or denied. Valid values for port: a single port or a range of port numbers e.g.: "25", "25:1024". For "all" and "icmp", port should only be "0:65535".
  * `target` - (Required) CIDR to be allowed or denied. Valid values for target: IPv4 CIDRs. Example: "10.30.0.0/16".

### Misc.
* `manage_user_attachment` - (Optional) This parameter is a switch used to determine whether or not to manage VPN user attachments to the VPN profile using this resource. If this is set to false, attachment must be managed using the **aviatrix_vpn_user** resource. Valid values: true, false. Default value: true.
* `users` - (Optional) List of VPN users to attach to this profile. This should be set to null if `manage_user_attachment` is set to false.


## Import

**vpn_profile** can be imported using the VPN profile's `name`, e.g.

```
$ terraform import aviatrix_vpn_profile.test name
```
