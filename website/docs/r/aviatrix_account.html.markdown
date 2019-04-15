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
  account_name       = "username"
  cloud_type         = 1
  aws_account_number = "123456789012"
  aws_iam            = "true"
  aws_role_app       = "arn:aws:iam::123456789012:role/aviatrix-role-app"
  aws_role_ec2       = "arn:aws:iam::123456789012:role/aviatrix-role-ec2"
}

# Or you can create Aviatrix AWS account with access_key/secret key
resource "aviatrix_account" "tempacc" {
  account_name       = "username"
  cloud_type         = 1
  aws_iam            = "false"
  aws_account_number = "123456789012"
  aws_access_key     = "ABCDEFGHIJKL"
  aws_secret_key     = "ABCDEFGHIJKLabcdefghijkl"
}

# Create Aviatrix GCP account
resource "aviatrix_account" "tempacc_gcp" {
  account_name = "username"
  cloud_type = 4
  gcloud_project_id = "aviatrix-123456"
  gcloud_project_credentials_filepath = "/home/ubuntu/test_gcp/aviatrix-abc123.json"
}

# Create Aviatrix Azure ARM account 
resource "aviatrix_account" "tempacc_arm" {
  account_name = "username"
  cloud_type = 8
  arm_subscription_id  =  "12345678-abcd-efgh-ijkl-123456789abc"
  arm_directory_id     =  "abcdefgh-1234-5678-9100-abc123456789"
  arm_application_id   =  "1234abcd-12ab-34cd-56ef-abcdef123456"
  arm_application_key  =  "213df1SDF1231Gsaf/fa23-4A/324j12390801+FSwe=" 
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.
* `cloud_type` - (Required) Type of cloud service provider (Only AWS, GCP, and ARM supported currently. Enter 1 for AWS, 4 for GCP, 8 for ARM).
* `aws_account_number` - (Optional) AWS Account number to associate with Aviatrix account (Required when creating an account for AWS).
* `aws_iam` - (Optional) AWS IAM-role based flag, this option is for UserConnect.
* `aws_access_key` - (Optional) AWS Access Key (Required when aws_iam is "false" and when creating an account for AWS).
* `aws_secret_key` - (Optional) AWS Secret Key (Required when aws_iam is "false" and when creating an account for AWS).
* `aws_role_app` - (Optional) AWS App role ARN, this option is for UserConnect (Required when aws_iam is "true" and when creating an account for AWS).
* `aws_role_ec2` - (Optional) AWS EC2 role ARN, this option is for UserConnect (Required when aws_iam is "true" and when creating an account for AWS).
* `gcloud_project_id` - (Optional) GCloud Project ID
* `gcloud_project_credentials_filepath` - (Optional) GCloud Project Credentials [local filepath].json (Required when creating an account for GCP).
* `arm_subscription_id` = (Optional) Azure RM Subscription ID (Required when creating an account for ARM).
* `arm_directory_id` = (Optional) Azure RM Directory ID (Required when creating an account for ARM).
* `arm_application_id` = (Optional) Azure RM Application ID (Required when creating an account for ARM).
* `arm_application_key` = (Optional) Azure RM Application key (Required when creating an account for ARM).

Note: Please make sure that the IAM roles/profiles have already been created before running this, if aws_iam="true". More information on the IAM roles is at https://docs.aviatrix.com/HowTos/iam_policies.html and https://docs.aviatrix.com/HowTos/HowTo_IAM_role.html

## Import

Instance account can be imported using the account_name (when doing import, needs to leave aws_secret_key blank), e.g.

```
$ terraform import aviatrix_account.test account_name
```
