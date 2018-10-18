---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_profile"
sidebar_current: "docs-aviatrix-resource-vpn-profile"
description: |-
  Creates and manages an Aviatrix VPN User Profile.
---

# aviatrix_vpn_profile

The Profile resource allows the creation and management of an Aviatrix VPN VPN User Profile.

## Example Usage

```hcl
# Create Aviatrix AWS VPN User Profile
resource "aviatrix_vpn_profile" "test_profile1" {
  name = "my_profile"
  base_rule = "allow_all"
  users = ["user1", "user2"]
  policy = [
    {
     action = "deny"
     proto = "tcp"
     port = "443"
     target = "10.0.0.0/32"
    },
    {
     action = "deny"
     proto = "tcp"
     port = "443"
     target = "10.0.0.1/32"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Enter any name for the VPN profile
* `base_rule` - (Optional) Base policy rule of  the profile to be added. Enter "allow_all" or "deny_all", based on whether you want a white list or black list
* `users` - (Optional) List of VPN users to attach to this profile
* `policy` - (Optional) New security policy for the profile. Each policy has the following attributes:
    * `action` - (Optional) Should be the opposite of the base rule for correct behaviour. Valid values for action: "allow" and "deny"
    * `proto` - (Optional) Protocol to allow or deny. Valid values for protocol: "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp"
    * `port` - (Optional) Port to be allowed or denied. Valid values for port: a single port or a range of port numbers
e.g.: "25", "25:1024"
    * `target` - (Optional) CIDR to be allowed or denied. Valid values for target: CIDRs separated by comma. e.g.: "10.30.0.0/16,10.45.0.0/20"
