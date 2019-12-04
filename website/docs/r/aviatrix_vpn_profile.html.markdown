---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_profile"
description: |-
  Creates and manages Aviatrix VPN User Profiles
---

# aviatrix_vpn_profile

The aviatrix_vpn_profile resource allows the creation and management of Aviatrix VPN user profiles.

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

* `name` - (Required) Enter any name for the VPN profile.
* `base_rule` - (Optional) Base policy rule of  the profile to be added. Enter "allow_all" or "deny_all", based on whether you want a white list or black list.
* `users` - (Optional) List of VPN users to attach to this profile.
* `policy` - (Optional) New security policy for the profile. Each policy has the following attributes:
  * `action` - (Required) Should be the opposite of the base rule for correct behaviour. Valid values for action: "allow", "deny".
  * `proto` - (Required) Protocol to allow or deny. Valid values for protocol: "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp".
  * `port` - (Required) Port to be allowed or denied. Valid values for port: a single port or a range of port numbers e.g.: "25", "25:1024". For "all" and "icmp", port should only be "0:65535".
  * `target` - (Required) CIDR to be allowed or denied. Valid values for target: IPv4 CIDRs. Example: "10.30.0.0/16".

## Import

Instance vpn_profile can be imported using the name, e.g.

```
$ terraform import aviatrix_vpn_profile.test name
```
