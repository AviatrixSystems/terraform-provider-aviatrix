---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall"
description: |-
  Creates and manages Aviatrix Firewall Policies
---

# aviatrix_firewall

The aviatrix_firewall resource allows the creation and management of Aviatrix Firewall policies.

## Example Usage

```hcl
# Create an Aviatrix Firewall
resource "aviatrix_firewall" "test_firewall" {
  gw_name          = "gateway-1"
  base_policy      = "allow-all"
  base_log_enabled = true

  policy {
    protocol    = "tcp"
    src_ip      = "10.15.0.224/32"
    log_enabled = false
    dst_ip      = "10.12.0.172/32"
    action      = "allow"
    port        = "0:65535"
    description = "This is policy no.1"
  }

  policy {
    protocol    = "tcp"
    src_ip      = "10.15.1.224/32"
    log_enabled = false
    dst_ip      = "10.12.1.172/32"
    action      = "deny"
    port        = "0:65535"
    description = "This is policy no.2"
  }

  policy {
    protocol    = "tcp"
    src_ip      = "10.15.2.224/32"
    log_enabled = false
    dst_ip      = "10.12.3.172/32"
    action      = "force-drop"
    port        = "0:65535"
    description = "This is policy no.3"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) The name of gateway.
* `base_policy` - (Optional) New base policy. Valid Values: "allow-all", "deny-all".
* `base_log_enabled` - (Optional) Indicates whether enable logging or not. Valid Values: true, false.
* `policy` - (Optional) New access policy for the gateway. Type: String (valid JSON). Seven fields are required for each policy item: src_ip, dst_ip, protocol, port, allow_deny, log_enabled and description.
  * `src_ip` - (Required) CIDRs separated by comma or tag names such "HR" or "marketing" etc. Example: "10.30.0.0/16,10.45.0.0/20". The aviatrix_firewall_tag resource should be created prior to using the tag name.
  * `dst_ip` - (Required) CIDRs separated by comma or tag names such "HR" or "marketing" etc. Example: "10.30.0.0/16,10.45.0.0/20". The aviatrix_firewall_tag resource should be created prior to using the tag name.
  * `protocol`- (Optional): "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp".
  * `port` - (Required) a single port or a range of port numbers. Example: "25", "25:1024".
  * `action`- (Required) Valid values: "allow", "deny" and "force-drop" (in stateful firewall rule to allow immediate packet dropping on established sessions).
  * `log_enabled`- (Optional) Valid values: true, false. Default value: false.
  * `description`- (Optional) Description of the policy. Example: "This is policy no.1".

## Import

Instance firewall can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_firewall.test gw_name
```
