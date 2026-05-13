---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_distributed_firewalling_default_action_rule"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling default action rule
---

# aviatrix_distributed_firewalling_default_action_rule

The **aviatrix_distributed_firewalling_default_action_rule** resource handles the creation and management of Distributed-firewalling default action rule. Available as of Provider 3.2.2+.

Once the Distributed Cloud Firewall (DCF) is enabled, the system will generate a default action rule (a DCF rule) with the following settings: action=PERMIT, logging=false, src=Anywhere, and dst=Anywhere. This configuration ensures that no traffic will be dropped.

However, if the default rule action is set to DENY, all traffic in the applied gateways (GWs) will be blocked. To allow traffic to pass, additional PERMIT rules must be created based on the specific requirements of the architecture design.

## Example Usage

```hcl
# Create an Aviatrix Distributed Firewalling Default Action Rule
resource "aviatrix_distributed_firewalling_default_action_rule" "test" {
  action = "DENY"
  logging = true
}
```

```hcl
# Create an Aviatrix Distributed Firewalling Default Action Rule with custom log profile
data "aviatrix_dcf_log_profile" "all" {
  profile_name = "start/end"
}

<!-- data "aviatrix_dcf_log_profile" "start" {
  profile_name = "start"
} -->

<!-- data "aviatrix_dcf_log_profile" "end" {
  profile_name = "end"
} -->


resource "aviatrix_distributed_firewalling_default_action_rule" "test_with_log_profile" {
  action      = "DENY"
  logging     = true
  log_profile = data.aviatrix_dcf_log_profile.all.profile_id
}
```

## Argument Reference

The following arguments are supported:

### Required
    * `action` - (Required) Action for the rule. Must be one of PERMIT or DENY. Type: String.
    * `logging` - (Required) Whether to enable logging for packets that match the rule. Type: Boolean.

### Optional
    * `log_profile` - (Optional) Logging profile UUID. There are 3 system defined log profiles that can be referenced using the `aviatrix_dcf_log_profile` data source:
        1. `start` - Log profile for logging session start only. If no log profile is provided, this will be used by default.
        2. `end` - Log profile for logging session end only
        3. `start/end` - Log profile for logging session start and end

## Import

**aviatrix_distributed_firewalling_default_action_rule** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_distributed_firewalling_default_action_rule.test 10-11-12-13
```
