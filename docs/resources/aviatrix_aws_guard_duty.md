---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_guard_duty"
description: |- Manage AWS GuardDuty configuration
---

# aviatrix_aws_guard_duty

The **aviatrix_aws_guard_duty** resource allows the enabling of [AWS GuardDuty](https://docs.aviatrix.com/HowTos/guardduty.html) for an account and region. To configure the Guard Duty scanning interval, please use the `aws_guard_duty_scanning_interval` attribute within the [`aviatrix_controller_config`](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/resources/aviatrix_controller_config) resource.

## Example Usage

```hcl
# Configure AWS GuardDuty 
resource "aviatrix_aws_guard_duty" "test_aws_guard_duty" {
  account_name = aviatrix_account.account_1.account_name
  region       = "us-west-1"
  excluded_ips = ["127.0.0.1", "10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `account_name` - (Required) Account name.
* `region` - (Required) Region.
* `excluded_ips` - (Optional) Set of excluded IPs.

## Import

**aws_guard_duty** resource can be imported with the `account_name` and `region` in the form "account_name~~region", e.g.

```
$ terraform import aws_guard_duty.test_aws_guard_duty devops-acc~~us-west-1
```