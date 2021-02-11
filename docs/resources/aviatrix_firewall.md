---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall"
description: |-
  Creates and manages Aviatrix Stateful Firewall Policies
---

# aviatrix_firewall

The **aviatrix_firewall** resource allows the creation and management of [Aviatrix Stateful Firewall](https://docs.aviatrix.com/HowTos/stateful_firewall_faq.html) policies.

~> **NOTE on Firewall and Firewall Policy resources:** Terraform currently provides both a standalone Firewall Policy resource and a Firewall resource with policies defined in-line. At this time, you cannot use a Firewall resource with in-line rules in conjunction with any Firewall Policy resources. Doing so will cause a conflict of policy settings and will overwrite policies. In order to use the **aviatrix_firewall_policy** resource, `manage_firewall_policies` must be set to false in this resource.

## Example Usage

```hcl
# Create an Aviatrix Firewall
resource "aviatrix_firewall" "stateful_firewall_1" {
  gw_name                  = "gateway-1"
  base_policy              = "allow-all"
  base_log_enabled         = true
  manage_firewall_policies = false
}
```

## Argument Reference

!> **WARNING:** Attribute `policy` has been deprecated as of provider version R2.18.1+ and will not receive further updates. Please use the standalone `aviatrix_firewall_policy` resource instead, and set `manage_firewall_policies` to false.

The following arguments are supported:

* `gw_name` - (Required) Gateway name to attach firewall policy to.
* `base_policy` - (Optional) New base policy. Valid Values: "allow-all", "deny-all". Default value: "deny-all"
* `base_log_enabled` - (Optional) Indicates whether enable logging or not. Valid Values: true, false. Default value: false.
* `manage_firewall_policies` - (Optional) Enable to manage firewall policies via in-line rules. If false, policies must be managed using `aviatrix_firewall_policy` resources. Default: true. Valid values: true, false. Available in provider version R2.17+.
* `policy` - (Optional) New access policy for the gateway. Type: String (valid JSON). Seven fields are required for each policy item: `src_ip`, `dst_ip`, `protocol`, `port`, `allow_deny`, `log_enabled` and `description`.
  * `src_ip` - (Required) CIDRs separated by comma or tag names such "HR" or "marketing" etc. Example: "10.30.0.0/16,10.45.0.0/20". The **aviatrix_firewall_tag** resource should be created prior to using the tag name.
  * `dst_ip` - (Required) CIDRs separated by comma or tag names such "HR" or "marketing" etc. Example: "10.30.0.0/16,10.45.0.0/20". The **aviatrix_firewall_tag** resource should be created prior to using the tag name.
  * `protocol`- (Optional): "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp".
  * `port` - (Required) a single port or a range of port numbers. Example: "25", "25:1024".
  * `action`- (Required) Valid values: "allow", "deny" and "force-drop" (in stateful firewall rule to allow immediate packet dropping on established sessions).
  * `log_enabled`- (Optional) Valid values: true, false. Default value: false.
  * `description`- (Optional) Description of the policy. Example: "This is policy no.1".

## Import

**firewall** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_firewall.test gw_name
```
