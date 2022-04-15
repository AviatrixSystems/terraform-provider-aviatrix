---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_microseg_policy_list"
description: |-
  Creates and manages an Aviatrix Micro-segmentation Policy List
---

# aviatrix_microseg_policy_list

The **aviatrix_microseg_policy_list** resource handles the creation and management of Micro-segmentation Policies. Available as of Provider R2.22.0+.

## Example Usage

```hcl
# Create an Aviatrix Microseg Policy
resource "aviatrix_microseg_policy_list" "test" {
  policies {
    name            = "microseg-policy-1"
    action          = "DENY"
    priority        = 1
    protocol        = "ICMP"
    logging         = false
    watch           = false
    src_app_domains = [
      "f15c9890-c8c4-4c1a-a2b5-ef0ab34d2e30"
    ]
    dst_app_domains = [
      "82e50c85-82bf-4b3b-b9da-aaed34a3aa53"
    ]
  }
  
  policies {
    name            = "microseg-policy"
    action          = "PERMIT"
    priority        = 0
    protocol        = "TCP"
    src_app_domains = [
      "7e7d1573-7a7a-4a53-bcb5-1ad5041961e0"
    ]
    dst_app_domains = [
      "f05b0ad7-d2d7-4d16-b2f6-48492319414c"
    ]
    
    port_ranges {
      hi = 50000
      lo = 49000
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `policies` - (Required) List of policies.
    * `name` - (Required) Name of the policy.
    * `action` - (Required) Action for the policy. Must be one of PERMIT or DENY.
    * `priority` - (Optional)  Priority for the policy. Default: 0. Type: Integer.
    * `protocol` - (Required) Protocol for the policy. Must be one of TCP, UDP or ICMP.
    * `src_app_domains` - (Required) List of App Domain UUIDs for the source for the policy.
    * `dst_app_domains` - (Required) List of App Domain UUIDs for the destination for the policy.
    * `port_ranges` - (Optional) List of port ranges for the policy. Cannot be used when `protocol` is "ICMP".
      * `lo` - (Required) Lower bound for the range of ports.
      * `hi` - (Optional) Upper bound for the range of ports. When not set, `lo` is the only port that matches the policy.
    * `uuid` - (Computed) UUID for the Policy.

## Import

**aviatrix_microseg_policy_list** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_microseg_policy_list.test 10-11-12-13
```
