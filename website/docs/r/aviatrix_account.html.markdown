---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_account"
sidebar_current: "docs-aviatrix-resource-account"
description: |-
  Creates and manages an Aviatrix cloud account.
---

# aviatrix_account

The Account resource allows the creation and management of an Aviatrix cloud account.

## Example Usage

```hcl
# Create Aviatrix AWS account with IAM roles
resource "aviatrix_account" "tempacc" {
  account_name = "username"
  account_password = "password"
  account_email = "abc@xyz.com"
  cloud_type = 1
  aws_account_number = "123456789012"
  aws_iam = "true"
  aws_role_app = "arn:aws:iam::123456789012:role/aviatrix-role-app"
  aws_role_ec2 = "arn:aws:iam::123456789012:role/aviatrix-role-ec2"
}

# Or you can create Aviatrix AWS account with access_key/secret key
resource "aviatrix_account" "tempacc" {
  account_name = "username"
  account_password = "password"
  account_email = "abc@xyz.com"
  cloud_type = 1
  aws_account_number = "123456789012"
  aws_access_key = "ABCDEFGHIJKL"
  aws_secret_key = "ABCDEFGHIJKLabcdefghijkl"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.
* `account_password` - (Required) Password for the account.
* `account_email` - (Required) Account email address where notifications will be sent.
* `cloud_type` - (Required) Type of cloud service provider. (Only AWS is supported currently. Enter 1 for AWS.)
* `aws_account_number` - (Required) AWS Account number to associate with Aviatrix account.
* `aws_iam` - (Optional) AWS IAM-role based flag, this option is for UserConnect.
* `aws_access_key` - (Optional) AWS Access Key (Required when aws_iam is "false" and when creating an account for AWS)
* `aws_secret_key` - (Optional) AWS Secret Key (Required when aws_iam is "false" and when creating an account for AWS)
* `aws_role_app` - (Optional) AWS App role ARN, this option is for UserConnect (Required when aws_iam is "true" and when creating an account for AWS).
* `aws_role_ec2` - (Optional) AWS EC2 role ARN, this option is for UserConnect (Required when aws_iam is "true" and when creating an account for AWS).
