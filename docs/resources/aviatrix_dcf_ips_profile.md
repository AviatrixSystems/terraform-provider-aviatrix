---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_ips_profile"
description: |-
  Creates and manages an Aviatrix DCF IPS profile
---

# aviatrix_dcf_ips_profile

The **aviatrix_dcf_ips_profile** resource allows you to create and manage IPS profiles for Distributed Cloud Firewall (DCF) Intrusion Prevention System. Available as of Provider 3.2.2+.

## Example Usage

```hcl
# Create an IPS profile using uploaded rule feeds
resource "aviatrix_dcf_ips_profile" "custom_profile" {
  profile_name = "Custom Profile 2"

  rule_feeds {
    custom_feeds_ids   = [aviatrix_dcf_ips_rule_feed.custom_feed.uuid]
    external_feeds_ids = ["suricata-rules"]
    ignored_sids       = [100001, 100002]
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
- `profile_name` – (Required) Name of the IPS profile. Type: String.

### Optional
- `rule_feeds` – (Required) Rule feeds configuration block. If omitted, the profile will have no rule feeds by default.
    - `custom_feeds_ids` – (Required) List of custom rule feed UUIDs. Type: List(String). Can be empty.
    - `external_feeds_ids` – (Required) List of external rule feed IDs. Type: List(String). Can be empty.
    - `ignored_sids` – (Required) List of rule SIDs to ignore. Type: List(Number). Can be empty.
- `intrusion_actions` – (Required) Actions for different severity levels. Type: Map(String). Valid values: `alert`, `alert_and_drop`.
    - Keys: `informational`, `minor`, `major`, `critical` – (Required) Action for each severity level. Type: String.

### Computed
- `uuid` – UUID of the IPS profile. Type: String.

> **Notes:**
> - All lists inside `rule_feeds` can be empty if you do not want to specify any values.
> - If `rule_feeds` is omitted, the profile will not have any rule feeds.
> - You can update only the fields you want; for example, you can update just the `profile_name` or `intrusion_actions` without specifying `rule_feeds`.

## Import

**aviatrix_dcf_ips_profile** can be imported using the profile UUID:

```
$ terraform import aviatrix_dcf_ips_profile.example 74b8ed97-d07d-41c8-982c-33a645e1723e
