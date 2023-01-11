---
subcategory: "Secured Networking" ???????????????????????????????????????/
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dns_profile_list"
description: |-
  Creates and manages DNS profiles
---

# aviatrix_dns_profile_list

The **aviatrix_dns_profile_list** resource handles the creation and management of DNS profiles.

## Example Usage

```hcl
# Create a DNS profile
resource "aviatrix_dns_profile_list" "test" {
  profiles {
    name = "profileA"
    global = ["8.8.8.8", "8.8.4.4"]
    lan = ["1.2.3.4", "5.6.7.8"]
    local_domain_names = ["avx.internal.com", "avx.media.com"]
    wan = ["2.3.4.5", "6.7.8.9"]
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `profiles` - (Required) List of policies.
    * `name` - (Required) Name of the policy.
    * `action` - (Required) Action for the policy. Must be one of PERMIT or DENY.
    * `priority` - (Optional)  Priority for the policy. Default: 0. Type: Integer.
    * `protocol` - (Required) Protocol for the policy. Must be one of TCP, UDP, ICMP or ANY.
    * `src_smart_groups` - (Required) List of Smart Group UUIDs for the source for the policy.

## Import

**aviatrix_dns_profile_list** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_dns_profile_list.test 10-11-12-13
```
