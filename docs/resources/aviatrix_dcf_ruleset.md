---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_ruleset"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling Ruleset
---

# aviatrix_dcf_ruleset

The **aviatrix_dcf_ruleset** resource handles the creation and management of Aviatrix Distributed-firewalling Policies and Ruleset.
Make sure to use one of the terraform attachment points to attach your terraform objects (rulesets/groups)

## Example Usage

The two terraform attachment points are:
- TERRAFORM_BEFORE_UI_MANAGED - Rulesets will be created before the rulessets mentioned in the UI
- TERRAFORM_AFTER_UI_MANAGED - Rulesets will be created after the rulesets mentioned in the UI.

The base terraform objects created in terraform should be attached to one of the above two attachment points, using data sources.
It is best to attach a policy_group to these above attachment_points, then place any ruleset in that policy_group, for easier management.

Note: We cannot attach 2 objects to these attachment points, only one object can be attached to each. If we want to build out a tree of multiple objects, we can attach a policy group to the above attachment points, and then create child attachment points as needed under this group.

Steps to attach a ruleset to one of the above attachment points or any other attachment point:

We need to get the attachment point ID based on its name (the name should be gloabally unique for each attachment point). In this example we will use the "TERRAFORM_BEFORE_UI_MANAGED" as the attachment point name to retrieve its ID and pass to the attach_to field in our ruleset.
```hcl
data "aviatrix_dcf_attachment_point" "tf_before_ui" {
    name = "TERRAFORM_BEFORE_UI_MANAGED"
}

We can then retrieve the ID of the attachment point and attach the ruleset to it using the attach_to field.
resource "aviatrix_dcf_ruleset" "base_ruleset" {
    # attach_to field can be used to attach to any other attachment_point in another policy_group
    attach_to = data.aviatrix_dcf_attachment_point.tf_before_ui.id
    name = "example-ruleset"
}
```

```hcl
# Create an Aviatrix Distributed Firewalling Ruleset
resource "aviatrix_dcf_ruleset" "test" {
  name = "Test rule"
  rules {
    name             = "df-rule-1"
    action           = "DENY"
    priority         = 1
    protocol         = "ICMP"
    logging          = false
    watch            = false
    src_smart_groups = [
      "f15c9890-c8c4-4c1a-a2b5-ef0ab34d2e30"
    ]
    dst_smart_groups = [
      "82e50c85-82bf-4b3b-b9da-aaed34a3aa53"
    ]
    tls_profile = "def000ad-6000-0000-0000-000000000001"
  }

  rules {
    name             = "df-rule"
    action           = "PERMIT"
    priority         = 0
    protocol         = "TCP"
    src_smart_groups = [
      "7e7d1573-7a7a-4a53-bcb5-1ad5041961e0"
    ]
    dst_smart_groups = [
      "f05b0ad7-d2d7-4d16-b2f6-48492319414c"
    ]

    port_ranges {
      hi = 50000
      lo = 49000
    }
  }
}
```
```hcl
# Create an Aviatrix Distributed Firewalling Ruleset
resource "aviatrix_dcf_ruleset" "test" {
  name = "Test rule"

  rules {
    name             = "df-rule"
    action           = "PERMIT"
    priority         = 0
    protocol         = "TCP"
    src_smart_groups = [
      "7e7d1573-7a7a-4a53-bcb5-1ad5041961e0"
    ]
    dst_smart_groups = [
      "f05b0ad7-d2d7-4d16-b2f6-48492319414c"
    ]

    port_ranges {
      hi = 50000
      lo = 49000
    }
  }

  rules {
    name                     = "df-rule-1"
    action                   = "DEEP_PACKET_INSPECTION_PERMIT"
    priority                 = 1
    protocol                 = "ANY"
    logging                  = false
    watch                    = false
    exclude_sg_orchestration = true
    src_smart_groups         = [
      "f15c9890-c8c4-4c1a-a2b5-ef0ab34d2e30"
    ]
    dst_smart_groups         = [
      "82e50c85-82bf-4b3b-b9da-aaed34a3aa53"
    ]
    web_groups               = [
      "6bff3e91-3707-4582-9ea6-70e37b08760b"
    ]
    flow_app_requirement     = "TLS_REQUIRED"
    decrypt_policy           = "DECRYPT_ALLOWED"
    port_ranges {
      hi = 50000
      lo = 49000
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of the ruleset
* `rules` - (Optional) Set of rules.
    * `name` - (Required) Name of the rule.
    * `action` - (Required) Action for the rule. Must be one of PERMIT, DENY, DEEP_PACKET_INSPECTION_PERMIT or INTRUSION_DETECTION_PERMIT.
    * `priority` - (Optional)  Priority for the rule. Default: 0. Type: Integer.
    * `protocol` - (Required) Protocol for the rule. Must be one of TCP, UDP, ICMP or ANY.
    * `src_smart_groups` - (Required) Set of Smart Group UUIDs for the source for the rule.
    * `dst_smart_groups` - (Required) Set of Smart Group UUIDs for the destination for the rule.
    * `web_groups` - (Optional) Set of Web Group UUIDs for the rule.
    * `flow_app_requirement` - (Optional) Flow application requirement for the rule. Must be one of APP_UNSPECIFIED, TLS_REQUIRED or NOT_TLS_REQUIRED.
    * `decrypt_policy` - (Optional) Decryption options for the rule. Must be one of DECRYPT_UNSPECIFIED, DECRYPT_ALLOWED or DECRYPT_NOT_ALLOWED.
    * `exclude_sg_orchestration` - (Optional) If this flag is set to true, this rule will be ignored for SG orchestration. Valid values: true, false. Default: false.
    * `port_ranges` - (Optional) Set of port ranges for the rule. Cannot be used when `protocol` is "ICMP".
      * `lo` - (Required) Lower bound for the range of ports.
      * `hi` - (Optional) Upper bound for the range of ports. When not set, `lo` is the only port that matches the rule.
    * `watch` - (Optional) Whether to enforce the rule or only watch packets. If "true" packets are only watched. This allows you to observe if the traffic impacted by this rule causes any inadvertent issues (such as traffic being dropped). Type: Boolean.
    * `logging` - (Optional) Whether to enable logging for packets that match the rule. Type: Boolean.
    * `uuid` - (Computed) UUID for the Rule.
    * `tls_profile` - (Optional) TLS profile UUID for the rule.
    * `log_profile` - (Optional) Logging profile UUID. Must be one of {"def000ad-7000-0000-0000-000000000001", "def000ad-7000-0000-0000-000000000002", "def000ad-7000-0000-0000-000000000003"}. The UUIDs correspod to: def000ad-7000-0000-0000-000000000001: DEF_LOG_PROFILE_START, def000ad-7000-0000-0000-000000000002: DEF_LOG_PROFILE_END, def000ad-7000-0000-0000-000000000003: DEF_LOG_PROFILE_ALL

## Import

**aviatrix_dcf_ruleset** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_dcf_ruleset.test <ruleset_uuid>
```
