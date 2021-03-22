---
subcategory: "Accounts"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_account"
description: |-
  Gets an Aviatrix cloud account's details.
---

# aviatrix_account

The **aviatrix_account** data source provides details about a specific cloud account created by the Aviatrix Controller.

This data source can prove useful when a module accepts an account's detail as an input variable.

## Example Usage

```hcl
# Aviatrix Account Data Source
data "aviatrix_account" "foo" {
  account_name = "username"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cloud_type` - Type of cloud service provider.
* `aws_account_number` - AWS Account number.
* `aws_access_key` - AWS Access Key.
* `aws_role_arn` - AWS App role ARN.
* `aws_role_ec2` - AWS EC2 role ARN.
* `awsgov_account_number` - AWS Gov Account number.
* `awsgov_access_key` - AWS Gov Access Key.
* `gcloud_project_id` - GCloud Project ID.
* `arm_subscription_id` - Azure ARM Subscription ID.
