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

```hcl
# Create an Aviatrix Firewall with in-line rules
resource "aviatrix_firewall" "stateful_firewall_1" {
  gw_name          = "gateway-1"
  base_policy      = "allow-all"
  base_log_enabled = true

  policy {
    protocol    = "all"
    src_ip      = "10.17.0.224/32"
    log_enabled = true
    dst_ip      = "10.12.0.172/32"
    action      = "force-drop"
    port        = "0:65535"
    description = "first_policy"
  }

  policy {
    protocol    = "tcp"
    src_ip      = "10.16.0.224/32"
    log_enabled = false
    dst_ip      = "10.12.1.172/32"
    action      = "force-drop"
    port        = "325"
    description = "second_policy"
  }

  policy {
    protocol    = "udp"
    src_ip      = "10.14.0.225/32"
    log_enabled = false
    dst_ip      = "10.13.1.173/32"
    action      = "deny"
    port        = "325"
    description = "third_policy"
  }
  
  policy {
    protocol    = "tcp"
    src_ip      = aviatrix_firewall_tag.test.firewall_tag
    log_enabled = false
    dst_ip      = "10.13.1.173/32"
    action      = "deny"
    port        = "325"
    description = "fourth_policy"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Gateway name to attach firewall policy to.
* `base_policy` - (Optional) New base policy. Valid Values: "allow-all", "deny-all". Default value: "deny-all"
* `base_log_enabled` - (Optional) Indicates whether enable logging or not. Valid Values: true, false. Default value: false.
* `manage_firewall_policies` - (Optional) Enable to manage firewall policies via in-line rules. If false, policies must be managed using `aviatrix_firewall_policy` resources. Default: true. Valid values: true, false. Available in provider version R2.17+.
* `policy` - (Optional) New access policy for the gateway. Seven fields are required for each policy item: `src_ip`, `dst_ip`, `protocol`, `port`, `action`, `log_enabled` and `description`. No duplicate rules(with same `src_ip`, `dst_ip`, `protocol` and `port`) are allowed 
  * `src_ip` - (Required) Source address, a valid IPv4 address or tag name such "HR" or "marketing" etc. Example: "10.30.0.0/16". The **aviatrix_firewall_tag** resource should be created prior to using the tag name.
  * `dst_ip` - (Required) Destination address, a valid IPv4 address or tag name such "HR" or "marketing" etc. Example: "10.30.0.0/16". The **aviatrix_firewall_tag** resource should be created prior to using the tag name.
  * `protocol`- (Optional): Valid values: "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp". Default value: "all".
  * `port` - (Required) A single port or a range of port numbers. Example: "25", "25:1024".
  * `action`- (Required) Valid values: "allow", "deny" and "force-drop" (in stateful firewall rule to allow immediate packet dropping on established sessions).
  * `log_enabled`- (Optional) Valid values: true, false. Default value: false.
  * `description`- (Optional) Description of the policy. Example: "This is policy no.1".

## Import

**firewall** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_firewall.test gw_name
```
