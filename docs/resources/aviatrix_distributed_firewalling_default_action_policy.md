---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_distributed_firewalling_default_action_policy"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling default action policy
---

# aviatrix_distributed_firewalling_default_action_policy

The **aviatrix_distributed_firewalling_default_action_policy** resource handles the creation and management of Distributed-firewalling default action policy. Available as of Provider 3.2.2+.

## Example Usage

```hcl
# Create an Aviatrix Distributed Firewalling Default Action Policy
resource "aviatrix_distributed_firewalling_default_action_policy" "test" {
  action = "DENY"
  logging = true
}
```

## Argument Reference

The following arguments are supported:

### Required
    * `action` - (Required) Action for the policy. Must be one of PERMIT or DENY. Type: String.
    * `logging` - (Required) Whether to enable logging for packets that match the policy. Type: Boolean.

## Import

**aviatrix_distributed_firewalling_default_action_policy** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_distributed_firewalling_default_action_policy.test 10-11-12-13
```
