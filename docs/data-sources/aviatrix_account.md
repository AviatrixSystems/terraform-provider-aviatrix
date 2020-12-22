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

* `cloud_type` - Type of cloud service provider. (Only AWS is supported currently. Value of 1 for AWS.)
* `aws_account_number` - AWS Account number to associate with Aviatrix account.
* `aws_access_key` - AWS Access Key.
* `aws_role_app` - AWS App role ARN.
* `aws_role_ec2` - AWS EC2 role ARN.
* `gcloud_project_id` - GCloud Project ID.
* `gcloud_project_credentials_filepath` - GCloud Project Credentials.
* `arm_subscription_id` - Azure ARM Subscription ID.
* `arm_directory_id` - Azure ARM Directory ID.
* `arm_application_id` - Azure ARM Application ID.
* `arm_application_key` - Azure ARM Application key.
* `oci_tenancy_id` - Oracle OCI Tenancy ID.
* `oci_user_id` - Oracle OCI User ID.
* `oci_compartment_id` - Oracle OCI Compartment ID.
* `oci_api_private_key_filepath` - Oracle OCI API Private Key local file path.
