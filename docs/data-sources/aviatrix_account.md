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

### AWS
* `aws_account_number` - AWS Account number.
* `aws_access_key` - AWS Access Key.
* `aws_role_arn` - AWS App role ARN.
* `aws_role_ec2` - AWS EC2 role ARN.
* `aws_gateway_role_app` - A separate AWS App role ARN to assign to gateways created by the controller. Available as of provider version R2.19+.
* `aws_gateway_role_ec2` - A separate AWS EC2 role ARN to assign to gateways created by the controller. Available as of provider version R2.19+.
  
### Azure
* `arm_subscription_id` - Azure ARM Subscription ID.

### GCP
* `gcloud_project_id` - GCloud Project ID.

### AzureGov Cloud
* `azure_gov_subscription_id` - AzureGov ARM Subscription ID.

### AWSGov Cloud
* `awsgov_account_number` - AWSGov Account number.
* `awsgov_access_key` - AWSGov Access Key.

### AWSChina Cloud
* `awschina_account_number` - AWSChina Account number. Available as of provider version R2.19+.
* `awschina_iam` - If enabled, `awschina_role_app` and `awschina_role_ec2` will be set. Otherwise, `awschina_access_key` will be set. Available as of provider version R2.19+.
* `awschina_role_app` - AWSChina App role ARN. Available as of provider version R2.19+.
* `awschina_role_ec2` - AWSChina EC2 role ARN. Available as of provider version R2.19+.
* `awschina_access_key` - AWSChina Access Key. Available as of provider version R2.19+.

### AzureChina Cloud
* `azurechina_subscription_id` - AzureChina ARM Subscription ID. Available as of provider version R2.19+.

### Alibaba Cloud
* `alicloud_account_id` - Alibaba Cloud Account ID.
