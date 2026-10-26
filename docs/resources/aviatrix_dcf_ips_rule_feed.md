---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_ips_rule_feed"
description: |-
  Creates and manages an Aviatrix DCF IPS rule feed
---

# aviatrix_dcf_ips_rule_feed

The **aviatrix_dcf_ips_rule_feed** resource allows you to upload and manage IPS rule feeds containing Suricata rules for Distributed Cloud Firewall (DCF) Intrusion Prevention System. Available as of Provider 3.2.2+.

## Example Usage

```hcl
# Upload a custom IPS rule feed
resource "aviatrix_dcf_ips_rule_feed" "custom_feed" {
  feed_name    = "malware_detection_rules"
  file_content = file("./malware_rules.rules")
}
```

## Argument Reference

The following arguments are supported:

### Required
- `feed_name` - (Required) Name for the rule feed. Type: String.
- `file_content` - (Required) IPS rule feed file content containing Suricata rules. Type: String.

### Computed
- `uuid` - UUID of the IPS rule feed. Type: String.
- `content_hash` - SHA-256 hash of the file content. Type: String.
- `ips_rules` - List of IPS rules extracted from the file. Type: List(String).

## Import

**aviatrix_dcf_ips_rule_feed** can be imported using the rule feed UUID:

```
$ terraform import aviatrix_dcf_ips_rule_feed.example 550e8400-e29b-41d4-a716-446655440000
```
