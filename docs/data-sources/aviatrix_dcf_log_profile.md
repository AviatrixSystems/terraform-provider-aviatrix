---
subcategory: "Distributed Cloud Firewall"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_log_profile"
description: |-
  Gets details about a specific DCF log profile.
---

# aviatrix_dcf_log_profile

The **aviatrix_dcf_log_profile** data source provides details about a specific Distributed Cloud Firewall (DCF) log profile.

There are 3 system defined log_profiles that can be referenced:
1. start - Log profile for logging session start only
2. end - Log profile for logging session end only
3. start/end - Log profile for logging session start and end


## Example Usage

```hcl
# Aviatrix DCF Log Profile Data Source
data "aviatrix_dcf_log_profile" "example" {
  profile_name = "my-log-profile"
}

# Use the log profile ID in a DCF rule
resource "aviatrix_dcf_ruleset" "example" {
  name = "Example ruleset with custom log profile"

  rules {
    name             = "rule-with-custom-logging"
    action           = "PERMIT"
    priority         = 1
    protocol         = "TCP"
    logging          = true
    # Use the profile_id to refer to a log_profile here.
    log_profile      = data.aviatrix_dcf_log_profile.example.profile_id
    src_smart_groups = [
      "f15c9890-c8c4-4c1a-a2b5-ef0ab34d2e30"
    ]
    dst_smart_groups = [
      "82e50c85-82bf-4b3b-b9da-aaed34a3aa53"
    ]
    port_ranges {
      lo = 80
      hi = 80
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `profile_name` - (Required) (String) Name of the Log Profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `profile_id` - (String) The unique identifier for the Log Profile which can be referenced in a DCF Rule
* `session_end` - (Boolean) Tells us if the logging of session end is enabled.
* `session_start` - (Boolean) Tells us if the logging of session start is enabled.
