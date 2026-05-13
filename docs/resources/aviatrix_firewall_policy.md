---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_policy"
description: |-
  Manages Aviatrix Stateful Firewall Policies
---

# aviatrix_firewall_policy

The **aviatrix_firewall_policy** resource manages a single Stateful Firewall policy resource.

~> **NOTE on Firewall and Firewall Policy resources:** Terraform currently provides both a standalone Firewall Policy resource and a Firewall resource with policies defined in-line. At this time, you cannot use a Firewall resource with in-line rules in conjunction with any Firewall Policy resources. Doing so will cause a conflict of policy settings and will overwrite policies. In order to use this resource, please set `manage_firewall_policies` in the **aviatrix_firewall** resource to false.

## Example Usage

```hcl
# Create an Aviatrix Stateful Firewall Policy
resource "aviatrix_firewall_policy" "test_firewall_policy" {
  gw_name     = aviatrix_firewall.test_firewall.gw_name
  src_ip      = "10.15.0.224/32"
  dst_ip      = "10.12.0.172/32"
  protocol    = "tcp"
  port        = "0:65535"
  action      = "allow"
  log_enabled = true
  description = "Test policy."
}
```
```hcl
# Create an Aviatrix Stateful Firewall Policy and insert it to a specific position
resource "aviatrix_firewall_policy" "test_firewall_policy" {
  gw_name     = aviatrix_firewall.test_firewall.gw_name
  src_ip      = "10.15.0.225/32"
  dst_ip      = "10.12.0.173/32"
  protocol    = "tcp"
  port        = "0:65535"
  action      = "allow"
  log_enabled = true
  description = "Test policy."
  position    = 2
}
```
## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Gateway name to attach firewall policy to.
* `src_ip` - (Required) CIDRs separated by comma or tag names such "HR" or "marketing" etc. Example: "10.30.0.0/16,10.45.0.0/20". The **aviatrix_firewall_tag** resource should be created prior to using the tag name.
* `dst_ip` - (Required) CIDRs separated by comma or tag names such "HR" or "marketing" etc. Example: "10.30.0.0/16,10.45.0.0/20". The **aviatrix_firewall_tag** resource should be created prior to using the tag name.
* `protocol`- (Optional): "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp".
* `port` - (Required) A single port or a range of port numbers. Example: "25", "25:1024".
* `action`- (Required) Valid values: "allow", "deny" and "force-drop" (in stateful firewall rule to allow immediate packet dropping on established sessions).
* `log_enabled`- (Optional) Valid values: true, false. Default value: false.
* `description`- (Optional) Description of the policy. Example: "This is policy no.1".
* `position`- (Optional) Position in the policy list, where the firewall policy will be inserted to. Valid values: any positive integer. Example: 2. If it is larger than the size of policy list, the policy will be inserted to the end.

## Import

**firewall_policy** can be imported using the `gw_name`, `src_ip`, `dst_ip`, `protocol`, `port` and `action` separated by `~`, e.g.

```
$ terraform import aviatrix_firewall_policy.test "gw_name~src_ip~dst_ip~protocol~port~action"
```
