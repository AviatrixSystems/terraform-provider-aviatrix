---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_cloudwatch_agent"
description: |-
  Enables and disables cloudwatch_agent
---

# aviatrix_cloudwatch_agent

The **aviatrix_cloudwatch_agent** resource allows the enabling and disabling of cloudwatch agent.

## Example Usage

```hcl
# Enable cloudwatch agent
resource "aviatrix_cloudwatch_agent" "test_cloudwatch_agent" {
  cloudwatch_role_arn = "arn:aws:iam::469550033836:role/aviatrix-role-cloudwatch"
  region              = "us-east-1"
  excluded_gateways   = ["a", "b"]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `cloudwatch_role_arn` (Required) CloudWatch role ARN.
* `region` (Required) Name of AWS region.

### Optional
* `log_group_name` (Optional) Log group name. "AVIATRIX-CLOUDWATCH-LOG" by default.
* `excluded_gateways` (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of cloudwatch agent.

## Import

**cloudwatch_agent** can be imported using "cloudwatch_agent", e.g.

```
$ terraform import aviatrix_cloudwatch_agent.test cloudwatch_agent
```
