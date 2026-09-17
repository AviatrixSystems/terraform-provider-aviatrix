---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall"
description: |-
  Gets the Aviatrix Firewall.
---

# aviatrix_firewall

Use this data source to get the Aviatrix stateful firewall for use in other resources.

## Example Usage

```hcl
# Aviatrix Firewall Data Source
data "aviatrix_firewall" "foo" {
  gw_name = "gw-abcd"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Name of the gateway associated with the firewall.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `base_policy` - The firewall's base policy.
* `base_log_enabled` - Indicates whether logging is enabled or not.
* `policies` - List of policies associated with the firewall.
  * `src_ip` - CIDRs separated by a comma or tag names such 'HR' or 'marketing' etc.
  * `dst_ip` - CIDRs separated by a comma or tag names such 'HR' or 'marketing' etc.
  * `protocol` - `all`, `tcp`, `udp`, `icmp`, `sctp`, `rdp` or `dccp`.
  * `port` - A single port or a range of port numbers.
  * `action`- `allow`, `deny` or `force-drop`(allow immediate packet dropping on established sessions).
  * `log_enabled` - Indicates whether logging is enabled or not.
  * `description`- Policy description.
