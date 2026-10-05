---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_default_ips_profile"
description: |-
  Creates and manages the Aviatrix DCF default IPS profile
---

# aviatrix_dcf_default_ips_profile

The **aviatrix_dcf_default_ips_profile** resource manages the default IPS profile for Distributed Cloud Firewall (DCF) Intrusion Prevention System. Available as of Provider 3.2.2+.

## Example Usage

```hcl
resource "aviatrix_dcf_ips_profile" "custom" {
  profile_name = "Custom Profile"

  rule_feeds {
    custom_feeds_ids   = []
    external_feeds_ids = ["suricata-rules"]
    ignored_sids       = []
  }

  intrusion_actions = {
    informational = "alert"
    minor         = "alert"
    major         = "alert_and_drop"
    critical      = "alert_and_drop"
  }
}

resource "aviatrix_dcf_default_ips_profile" "default" {
  default_ips_profile = [aviatrix_dcf_ips_profile.custom.uuid]
}
```

## Argument Reference

The following arguments are supported:

### Required
- `default_ips_profile` – (Required) List of default IPS profile UUIDs. Only one IPS Profile is supported at the current version. Type: List(String).

## Import

**aviatrix_dcf_default_ips_profile** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_dcf_default_ips_profile.default 10-11-12-13
```
