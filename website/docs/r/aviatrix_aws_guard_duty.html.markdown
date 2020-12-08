---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_guard_duty"
description: |-
  Manage AWS GuardDuty configuration
---

# aviatrix_aws_guard_duty

The **aviatrix_aws_guard_duty** resource allows the configuration of [AWS GuardDuty](https://docs.aviatrix.com/HowTos/guardduty.html).


## Example Usage

```hcl
# Configure AWS GuardDuty 
resource "aviatrix_aws_guard_duty" "test_aws_guard_duty" {
  scanning_interval = 30
  enabled_accounts {
    account_name = aviatrix_account.account_1.account_name
    region       = "us-west-1"
    excluded_ips = ["127.0.0.1", "10.0.0.1"]
  }
  enabled_accounts {
    account_name = aviatrix_account.account_2.account_name
    region       = "us-east-1"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `scanning_interval` - (Optional) Scanning interval.
* `enabled_accounts` - (Required) Set of accounts to enable GuardDuty in.
  * `account_name` - (Required) Account name.
  * `region` - (Required) Region.
  * `excluded_ips` - (Optional) Set of excluded IPs.


## Import

Since there is only 1 logical **aws_guard_duty** resource, the ID is always the same "aviatrix_aws_guard_duty", e.g.

```
$ terraform import aws_guard_duty.test_aws_guard_duty aviatrix_aws_guard_duty
```
