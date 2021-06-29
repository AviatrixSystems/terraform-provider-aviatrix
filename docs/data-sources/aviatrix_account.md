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
* `aws_role_arn` - AWS App role ARN.
* `aws_role_ec2` - AWS EC2 role ARN.
* `aws_gateway_role_app` - A separate AWS App role ARN to assign to gateways created by the controller. Available as of provider version R2.19+.
* `aws_gateway_role_ec2` - A separate AWS EC2 role ARN to assign to gateways created by the controller. Available as of provider version R2.19+.
  
### Azure
* `arm_subscription_id` - Azure ARM Subscription ID.

### GCP
* `gcloud_project_id` - GCloud Project ID.

### AzureGov Cloud
* `azuregov_subscription_id` - AzureGov ARM Subscription ID.

### AWSGov Cloud
* `awsgov_account_number` - AWSGov Account number.
* `awsgov_iam` - If enabled, `awsgov_role_app` and `awschina_role_ec2` will be set. Available as of provider version R2.19+.
* `awsgov_role_app` - AWSGov App role ARN. Available as of provider version R2.19+.
* `awsgov_role_ec2` - AWSGov EC2 role ARN. Available as of provider version R2.19+.

### AWSChina Cloud
* `awschina_account_number` - AWSChina Account number. Available as of provider version R2.19+.
* `awschina_iam` - If enabled, `awschina_role_app` and `awschina_role_ec2` will be set. Available as of provider version R2.19+.
* `awschina_role_app` - AWSChina App role ARN. Available as of provider version R2.19+.
* `awschina_role_ec2` - AWSChina EC2 role ARN. Available as of provider version R2.19+.

### AzureChina Cloud
* `azurechina_subscription_id` - AzureChina ARM Subscription ID. Available as of provider version R2.19+.

### Alibaba Cloud
* `alicloud_account_id` - Alibaba Cloud Account ID.

### AWS Top Secret Cloud
* `awsts_account_number` - AWS Top Secret Region Account Number. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_url` - AWS Top Secret Region CAP Url. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_agency` - AWS Top Secret Region CAP Agency. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_mission` - AWS Top Secret Region Mission. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_role_name` - AWS Top Secret Region Role Name. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
  `awsts_cap_cert_path` - AWS Top Secret Region CAP Certificate file name on the controller. Available as of provider R2.19.5+.
* `awsts_cap_cert_key_path` - AWS Top Secret Region CAP Certificate Key file name on the controller. Available as of provider R2.19.5+.
* `aws_ca_cert_path` - AWS Top Secret Region or Secret Region Custom Certificate Authority file name on the controller. Available as of provider R2.19.5+.

### AWS Secret Cloud
* `awss_account_number` - AWS Secret Region Account Number. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_url` - AWS Secret Region CAP Url. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_agency` - AWS Secret Region CAP Agency. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_account_name` - AWS Secret Region Account Name. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_role_name` - AWS Secret Region Role Name. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_cert_path` - AWS Secret Region CAP Certificate file name on the controller. Available as of provider R2.19.5+.
* `awss_cap_cert_key_path` - AWS Secret Region CAP Certificate Key file name on the controller. Available as of provider R2.19.5+.
