---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_data_account"
sidebar_current: "docs-aviatrix-data_source-account"
description: |-
  Gets the an Aviatrix cloud account.
---

# aviatrix_account

Use this data source to get the Aviatrix cloud account for use in other resources.

## Example Usage

```hcl
# Create Aviatrix account data source
data "aviatrix_account" "foo" {
  account_name = "username"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Attribute Reference

* `cloud_type` - Type of cloud service provider. (Only AWS is supported currently. Value of 1 for AWS.)
* `aws_account_number` - AWS Account number to associate with Aviatrix account.
* `aws_access_key` - AWS Access Key. 
* `aws_role_app` - AWS App role ARN.
* `aws_role_ec2` - AWS EC2 role ARN.
