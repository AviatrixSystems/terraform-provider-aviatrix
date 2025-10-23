---
subcategory: "Distributed Cloud Firewall"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_log_profile"
description: |-
  Gets details about a specific DCF log profile.
---

# aviatrix_dcf_log_profile

The **aviatrix_dcf_log_profile** data source provides details about a specific Distributed Cloud Firewall (DCF) log profile created by the Aviatrix Controller.

## Example Usage

```hcl
# Aviatrix DCF Log Profile Data Source
data "aviatrix_dcf_log_profile" "example" {
  profile_name = "my-log-profile"
}
```

## Argument Reference

The following arguments are supported:

* `profile_name` - (Required) (String) Name of the Log Profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `profile_id` - (String) The unique identifier for the Log Profile.
* `session_end` - (Boolean) Toggle to enable logging of session end.
* `session_start` - (Boolean) Toggle to enable logging of session start.
