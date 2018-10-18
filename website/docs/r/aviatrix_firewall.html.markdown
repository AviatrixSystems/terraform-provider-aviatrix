---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall"
sidebar_current: "docs-aviatrix-resource-firewall"
description: |-
  Creates and manages Aviatrix Firewall Policy
---

# aviatrix_firewall

The Firewall resource allows the creation and management of Aviatrix Firewall Policy

## Example Usage

```hcl
# Create Aviatrix Firewall Policy
resource "aviatrix_firewall" "test_firewall" {
  gw_name = "gateway-1"
  base_allow_deny =  "allow-all"
  base_log_enable = "off"
  policy = [
            {
              protocol = "tcp"
              src_ip = "10.15.0.224/32"
              log_enable = "on"
              dst_ip = "10.12.0.172/32"
              allow_deny = "deny"
              port = "0-65535"
            },
            {
              protocol = "tcp"
              src_ip = "test_tag"
              log_enable = "off"
              dst_ip = "10.12.1.172/32"
              allow_deny = "deny"
              port = "0-65535"
            }
          ]
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) The name of gateway.
* `base_allow_deny` - (Optional) New base policy. Valid Value: "allow-all", "deny-all".
* `base_log_enable` - (Optional) Indicates whether enable logging or not. Valid Value: "on", "off"
* `policy` - (Optional) New access policy for the gateway. Type: String (valid JSON). 6 fields are required for each policy item: src_ip, dst_ip, protocol, port, allow_deny, log_enable. Valid values are 
  * `src_ip` - CIDRs separated by comma or tag names such "HR" or "marketing" etc.  e.g.: "10.30.0.0/16,10.45.0.0/20".
  * `dst_ip` - CIDRs separated by comma or tag names such "HR" or "marketing" etc.  e.g.: "10.30.0.0/16,10.45.0.0/20".
  * `protocol`: "all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp".
  * `port`: a single port or a range of port numbers. e.g.: "25", "25:1024".
  * `allow_deny`: "allow" and "deny"
  * `log_enable`: "on" and "off"
