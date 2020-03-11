---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_account"
description: |-
  Creates and manages Aviatrix cloud accounts
---

# aviatrix_account

The **aviatrix_account** resource allows the creation and management of Aviatrix cloud accounts.

## Example Usage

```hcl
# Create an Aviatrix AWS Account with IAM roles
resource "aviatrix_account" "temp_acc_aws" {
  account_name       = "username"
  cloud_type         = 1
  aws_account_number = "123456789012"
  aws_iam            = true
  aws_role_app       = "arn:aws:iam::123456789012:role/aviatrix-role-app"
  aws_role_ec2       = "arn:aws:iam::123456789012:role/aviatrix-role-ec2"
}
```
```hcl
# Or you can create an Aviatrix AWS Account with access_key/secret key
resource "aviatrix_account" "temp_acc_aws" {
  account_name       = "username"
  cloud_type         = 1
  aws_iam            = false
  aws_account_number = "123456789012"
  aws_access_key     = "ABCDEFGHIJKL"
  aws_secret_key     = "ABCDEFGHIJKLabcdefghijkl"
}
```
```hcl
# Create an Aviatrix GCP Account
resource "aviatrix_account" "temp_acc_gcp" {
  account_name                        = "username"
  cloud_type                          = 4
  gcloud_project_id                   = "aviatrix-123456"
  gcloud_project_credentials_filepath = "/home/ubuntu/test_gcp/aviatrix-abc123.json"
}
```
```hcl
# Create an Aviatrix Azure Account
resource "aviatrix_account" "temp_acc_azure" {
  account_name        = "username"
  cloud_type          = 8
  arm_subscription_id = "12345678-abcd-efgh-ijkl-123456789abc"
  arm_directory_id    = "abcdefgh-1234-5678-9100-abc123456789"
  arm_application_id  = "1234abcd-12ab-34cd-56ef-abcdef123456"
  arm_application_key = "213df1SDF1231Gsaf/fa23-4A/324j12390801+FSwe="
}
```
```hcl
# Create an Aviatrix Oracle OCI Account
resource "aviatrix_account" "temp_acc_oci" {
  account_name                 = "username"
  cloud_type                   = 16
  oci_tenancy_id               = "ocid1.tenancy.oc1..aaaaaaaa"
  oci_user_id                  = "ocid1.user.oc1..aaaaaaaazly"
  oci_compartment_id           = "ocid1.tenancy.oc1..aaaaaaaaxo"
  oci_api_private_key_filepath = "/Users/public/Documents/oci_api_key.pem"
}
```
```hcl
# Create an Aviatrix AWS Gov Account
resource "aviatrix_account" "temp_acc_awsgov" {
  account_name          = "username"
  cloud_type            = 256
  awsgov_account_number = "123456789012"
  awsgov_access_key     = "ABCDEFGHIJKL"
  awsgov_secret_key     = "ABCDEFGHIJKLabcdefghijkl"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.
* `cloud_type` - (Required) Type of cloud service provider. Only AWS, GCP, AZURE, OCI, and AWS Gov are supported currently. Enter 1 for AWS, 4 for GCP, 8 for AZURE, 16 for OCI, 256 for AWS Gov.

### AWS
* `aws_account_number` - (Optional) AWS Account number to associate with Aviatrix account. Required when creating an account for AWS.
* `aws_iam` - (Optional) AWS IAM-role based flag, this option is for UserConnect.
* `aws_access_key` - (Optional) AWS Access Key. Required when `aws_iam` is "false" and when creating an account for AWS.
* `aws_secret_key` - (Optional) AWS Secret Key. Required when `aws_iam` is "false" and when creating an account for AWS.
* `aws_role_app` - (Optional) AWS App role ARN, this option is for UserConnect. Required when `aws_iam` is "true" and when creating an account for AWS.
* `aws_role_ec2` - (Optional) AWS EC2 role ARN, this option is for UserConnect. Required when `aws_iam` is "true" and when creating an account for AWS.

### Azure
* `arm_subscription_id` - (Optional) Azure ARM Subscription ID. Required when creating an account for Azure.
* `arm_directory_id` - (Optional) Azure ARM Directory ID. Required when creating an account for Azure.
* `arm_application_id` - (Optional) Azure ARM Application ID. Required when creating an account for Azure.
* `arm_application_key` - (Optional) Azure ARM Application key. Required when creating an account for Azure.

### Google Cloud
* `gcloud_project_id` - (Optional) GCloud Project ID.
* `gcloud_project_credentials_filepath` - (Optional) GCloud Project Credentials [local filepath].json. Required when creating an account for GCP.

### Oracle Cloud
* `oci_tenancy_id` - (Optional) Oracle OCI Tenancy ID. Required when creating an account for OCI.
* `oci_user_id` - (Optional) Oracle OCI User ID. Required when creating an account for OCI.
* `oci_compartment_id` - (Optional) Oracle OCI Compartment ID. Required when creating an account for OCI.
* `oci_api_private_key_filepath` - (Optional) Oracle OCI API Private Key local file path. Required when creating an account for OCI.

### AWS GovCloud
* `awsgov_account_number` - (Optional) AWS Gov Account number to associate with Aviatrix account. Required when creating an account for AWS Gov.
* `awsgov_access_key` - (Optional) AWS Access Key. Required when creating an account for AWS Gov.
* `awsgov_secret_key` - (Optional) AWS Secret Key. Required when creating an account for AWS Gov.

-> **NOTE:** Please make sure that the IAM roles/profiles have already been created before running this, if `aws_iam = true`. More information on the IAM roles is at https://docs.aviatrix.com/HowTos/iam_policies.html and https://docs.aviatrix.com/HowTos/HowTo_IAM_role.html

## Import

**account** can be imported using the `account_name` (when doing import, need to leave `aws_secret_key` blank), e.g.

```
$ terraform import aviatrix_account.test account_name
```
