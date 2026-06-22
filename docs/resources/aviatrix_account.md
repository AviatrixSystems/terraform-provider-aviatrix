---
subcategory: "Accounts"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_account"
description: |-
  Creates and manages Aviatrix cloud accounts
---

# aviatrix_account

The **aviatrix_account** resource allows the creation and management of Aviatrix cloud accounts.

~> **NOTE:** With the release of Controller 5.4 (compatible with Aviatrix provider R2.13), Role-Based Access Control (RBAC) is now integrated into the Accounts workflow. Any **aviatrix_account** created in 5.3 by default will have admin privileges (attached to the 'admin' RBAC permission group). In 5.4, any new accounts created will not be attached to any RBAC group unless otherwise specified through the **aviatrix_rbac_group_access_account_attachment** resource.

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
# Create an Aviatrix AzureGov Account
resource "aviatrix_account" "temp_acc_azure_gov" {
  account_name             = "username"
  cloud_type               = 32
  azuregov_subscription_id = "12345678-abcd-efgh-ijkl-123456789abc"
  azuregov_directory_id    = "abcdefgh-1234-5678-9100-abc123456789"
  azuregov_application_id  = "1234abcd-12ab-34cd-56ef-abcdef123456"
  azuregov_application_key = "213df1SDF1231Gsaf/fa23-4A/324j12390801+FSwe="
}
```
```hcl
# Create an Aviatrix AWSGov Account
resource "aviatrix_account" "temp_acc_awsgov" {
  account_name          = "username"
  cloud_type            = 256
  awsgov_account_number = "123456789012"
  awsgov_access_key     = "ABCDEFGHIJKL"
  awsgov_secret_key     = "ABCDEFGHIJKLabcdefghijkl"
}
```
```hcl
# Create an Aviatrix AWS China Account with IAM roles
resource "aviatrix_account" "temp_acc_aww_china_iam" {
  account_name            = "username"
  cloud_type              = 1024
  awschina_account_number = "123456789012"
  awschina_iam            = true
  awschina_role_app       = "arn:aws-cn:iam::123456789012:role/aviatrix-role-app"
  awschina_role_ec2       = "arn:aws-cn:iam::123456789012:role/aviatrix-role-ec2"
}
```
```hcl
# Or you can create an Aviatrix AWS China Account with access_key/secret key
resource "aviatrix_account" "temp_acc_awschina" {
  account_name            = "username"
  cloud_type              = 1024
  awschina_account_number = "123456789012"
  awschina_iam            = false
  awschina_access_key     = "ABCDEFGHIJKL"
  awschina_secret_key     = "ABCDEFGHIJKLabcdefghijkl"
}
```
```hcl
# Create an Aviatrix Azure China Account
resource "aviatrix_account" "temp_acc_azurechina" {
  account_name               = "username"
  cloud_type                 = 2048
  azurechina_subscription_id = "12345678-abcd-efgh-ijkl-123456789abc"
  azurechina_directory_id    = "abcdefgh-1234-5678-9100-abc123456789"
  azurechina_application_id  = "1234abcd-12ab-34cd-56ef-abcdef123456"
  azurechina_application_key = "213df1SDF1231Gsaf/fa23-4A/324j12390801+FSwe="
}
```
```hcl
# Create an Alibaba Cloud Account
resource "aviatrix_account" "temp_acc_alibaba" {
  account_name        = "username"
  cloud_type          = 8192
  alicloud_account_id = "123456789012"
  alicloud_access_key = "ABCDEFGHIJKL"
  alicloud_secret_key = "ABCDEFGHIJKLabcdefghijkl"
}
 ```

```hcl
# Create an Aviatrix AWS Account with IAM roles and RBAC group: group-c
resource "aviatrix_account" "temp_acc_aws" {
  account_name       = "username"
  cloud_type         = 1
  aws_account_number = "123456789012"
  aws_iam            = true
  aws_role_app       = "arn:aws:iam::123456789012:role/aviatrix-role-app"
  aws_role_ec2       = "arn:aws:iam::123456789012:role/aviatrix-role-ec2"
  rbac_groups        = ["group-c"]
}
```


## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.
* `cloud_type` - (Required) Type of cloud service provider. Only AWS, GCP, Azure, OCI, AzureGov, AWSGov, AWSChina, AzureChina, Alibaba Cloud, Edge CSP/Zededa, Edge Equinix and Edge NEO/Platform are supported currently. Enter 1 for AWS, 4 for GCP, 8 for Azure, 16 for OCI, 32 for AzureGov, 256 for AWSGov, 1024 for AWSChina or 2048 for AzureChina, 8192 for Alibaba Cloud, 65536 for Edge CSP/Zededa, 524288 for Edge Equinix, 262144 for Edge NEO/Platform, 1048576 for Edge Megaport.

### AWS
~> **NOTE:** As of Aviatrix provider version R2.19+, the Aviatrix Controller supports the use of custom IAM roles through the `aws_role_app` and `aws_role_ec2` attributes. If the Controller's IAM role is changed through the AWS console, please run `terraform apply -refresh=false` in order to update `aws_role_app` and `aws_role_ec2`. `audit_account` must be set to "false" when using custom IAM roles.
* `aws_account_number` - (Optional) AWS Account number to associate with Aviatrix account. Required when creating an account for AWS.
* `aws_iam` - (Optional) AWS IAM-role based flag, this option is for UserConnect.
* `aws_access_key` - (Optional) AWS Access Key. Required when `aws_iam` is "false" and when creating an account for AWS.
* `aws_secret_key` - (Optional) AWS Secret Key. Required when `aws_iam` is "false" and when creating an account for AWS.
* `aws_role_app` - (Optional) AWS App role ARN, this option is for UserConnect. Required when `aws_iam` is "true" and when creating an account for AWS.
* `aws_role_ec2` - (Optional) AWS EC2 role ARN, this option is for UserConnect. Required when `aws_iam` is "true" and when creating an account for AWS.
* `aws_gateway_role_app` - (Optional) A separate AWS App role ARN to assign to gateways created by the controller. Required when `aws_gateway_role_ec2` is set. Only allowed when `aws_iam`, `awsgov_iam`, or `awschina_iam` is "true" when creating an account for AWS, AWSGov or AWSChina, respectively. Available as of provider version R2.19+.
* `aws_gateway_role_ec2` - (Optional) A separate AWS EC2 role ARN to assign to gateways created by the controller. Required when `aws_gateway_role_app` is set. Only allowed when `aws_iam`, `awsgov_iam`, or `awschina_iam` is "true" when creating an account for AWS, AWSGov or AWSChina, respectively. Available as of provider version R2.19+.

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

### AzureGov Cloud
* `azuregov_subscription_id` - (Optional) AzureGov ARM Subscription ID. Required when creating an account for AzureGov. Available as of provider version R2.19+.
* `azuregov_directory_id` - (Optional) AzureGov ARM Directory ID. Required when creating an account for AzureGov. Available as of provider version R2.19+.
* `azuregov_application_id` - (Optional) AzureGov ARM Application ID. Required when creating an account for AzureGov. Available as of provider version R2.19+.
* `azuregov_application_key` - (Optional) AzureGov ARM Application key. Required when creating an account for AzureGov. Available as of provider version R2.19+.

### AWSGov Cloud
* `awsgov_account_number` - (Optional) AWSGov Account number to associate with Aviatrix account. Required when creating an account for AWSGov.
* `awsgov_iam` - (Optional) AWSGov IAM-role based flag. Available as of provider version 2.19+.
* `awsgov_access_key` - (Optional) AWS Access Key. Required when creating an account for AWSGov.
* `awsgov_secret_key` - (Optional) AWS Secret Key. Required when creating an account for AWSGov.
* `awsgov_role_app` - (Optional) AWSGov App role ARN. Available when `awsgov_iam` is "true" and when creating an account for AWSGov. If left empty, the ARN will be computed. Available as of provider version 2.19+.
* `awsgov_role_ec2` - (Optional) AWSGov EC2 role ARN. Available when `awsgov_iam` is "true" and when creating an account for AWSGov. If left empty, the ARN will be computed. Available as of provider version 2.19+.

### AWSChina Cloud
* `awschina_account_number` - (Optional) AWSChina Account number to associate with Aviatrix account. Required when creating an account for AWSChina. Available as of provider version 2.19+.
* `awschina_iam` - (Optional) AWSChina IAM-role based flag. Available as of provider version 2.19+.
* `awschina_role_app` - (Optional) AWSChina App role ARN. Available when `awschina_iam` is "true" and when creating an account for AWSChina. If left empty, the ARN will be computed. Available as of provider version 2.19+.
* `awschina_role_ec2` - (Optional) AWSChina EC2 role ARN. Available when `awschina_iam` is "true" and when creating an account for AWSChina. If left empty, the ARN will be computed. Available as of provider version 2.19+.
* `awschina_access_key` - (Optional) AWSChina Access Key. Required when `awschina_iam` is "false" and when creating an account for AWSChina. Available as of provider version 2.19+.
* `awschina_secret_key` - (Optional) AWSChina Secret Key. Required when `awschina_iam` is "false" and when creating an account for AWSChina. Available as of provider version 2.19+.

### AzureChina Cloud
* `azurechina_subscription_id` - (Optional) AzureChina ARM Subscription ID. Required when creating an account for AzureChina. Available as of provider version 2.19+.
* `azurechina_directory_id` - (Optional) AzureChina ARM Directory ID. Required when creating an account for AzureChina. Available as of provider version 2.19+.
* `azurechina_application_id` - (Optional) AzureChina ARM Application ID. Required when creating an account for AzureChina. Available as of provider version 2.19+.
* `azurechina_application_key` - (Optional) AzureChina ARM Application key. Required when creating an account for AzureChina. Available as of provider version 2.19+.

### Alibaba Cloud
* `alicloud_account_id` - (Optional) Alibaba Cloud Account number to associate with Aviatrix account. Required when creating an account for Alibaba Cloud.
* `alicloud_access_key` - (Optional) Alibaba Cloud Access Key. Required when creating an account for Alibaba Cloud.
* `alicloud_secret_key` - (Optional) Alibaba Cloud Secret Key. Required when creating an account for Alibaba Cloud.

### AWS Top Secret Region
* `awsts_account_number` - (Optional) AWS Top Secret Region Account Number. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_url` - (Optional) AWS Top Secret Region CAP Url. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_agency` - (Optional) AWS Top Secret Region CAP Agency. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_mission` - (Optional) AWS Top Secret Region Mission. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_role_name` - (Optional) AWS Top Secret Region Role Name. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_cert` - (Optional) AWS Top Secret Region CAP Certificate local file path. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_cap_cert_key` - (Optional) AWS Top Secret Region CAP Certificate Key local file path. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.
* `awsts_ca_chain_cert` - (Optional) AWS Top Secret Region Custom Certificate Authority local file path. Required when creating an account in AWS Top Secret Region. Available as of provider version R2.19.5+.

### AWS Secret Region
* `awss_account_number` - (Optional) AWS Secret Region Account Number. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_url` - (Optional) AWS Secret Region CAP Url. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_agency` - (Optional) AWS Secret Region CAP Agency. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_account_name` - (Optional) AWS Secret Region Account Name. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_role_name` - (Optional) AWS Secret Region Role Name. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_cert` - (Optional) AWS Secret Region CAP Certificate local file path. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_cap_cert_key` - (Optional) AWS Secret Region CAP Certificate Key local file path. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.
* `awss_ca_chain_cert` - (Optional) AWS Secret Region Custom Certificate Authority local file path. Required when creating an account in AWS Secret Region. Available as of provider version R2.19.5+.

### Edge CSP
~> **NOTE:** Since V3.1.1+, please use `edge_zededa_username` and `edge_zededa_password` instead, `edge_csp_username` and `edge_csp_password` will be deprecated in the V3.2.0 release.
* `edge_csp_username` - (Optional) Edge CSP username. Required when creating an Edge CSP account.
* `edge_csp_password` - (Optional) Edge CSP password. Required when creating an Edge CSP account.

### Edge Zededa
* `edge_zededa_username` - (Optional) Edge Zededa username. Required when creating an Edge Zededa account.
* `edge_zededa_password` - (Optional) Edge Zededa password. Required when creating an Edge Zededa account.

### Misc.
~> **NOTE:** On Terraform versions 0.12.x, 0.13.x, and 0.14.x, Terraform will not detect any changes to the account when the account audit fail warning is given. In order to apply changes or set `audit_account = false`, please run `terraform apply -refresh=false`.
* `audit_account` - (Optional) Specify whether to enable the audit account feature. If this feature is enabled, terraform will give a warning if there is an issue with the account credentials. Changing `audit_account` to "false" will not prevent the Controller from performing account audits. It will only prevent Terraform from displaying a warning. Valid values: true, false. Default: false. Available as of provider version 2.19+. **Note: The warning may still appear for a few hours after fixing the underlying issue.**
* `rbac_groups` - (Optional) A list of existing RBAC group names. This attribute should only be used when creating an account. Updating this attribute will have no effect. Available as of provider version R2.23.0+.

-> **NOTE:** Please make sure that the IAM roles/profiles have already been created before running this, if `aws_iam = true`. More information on the IAM roles is at https://docs.aviatrix.com/HowTos/iam_policies.html and https://docs.aviatrix.com/HowTos/HowTo_IAM_role.html

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `awsts_cap_cert_path` - (Optional) AWS Top Secret Region CAP Certificate file name on the controller. Available as of provider R2.19.5+.
* `awsts_cap_cert_key_path` - (Optional) AWS Top Secret Region CAP Certificate Key file name on the controller. Available as of provider R2.19.5+.
* `aws_ca_cert_path` - (Optional) AWS Top Secret Region or Secret Region Custom Certificate Authority file name on the controller. Available as of provider R2.19.5+.
* `awss_cap_cert_path` - (Optional) AWS Secret Region CAP Certificate file name on the controller. Available as of provider R2.19.5+.
* `awss_cap_cert_key_path` - (Optional) AWS Secret Region CAP Certificate Key file name on the controller. Available as of provider R2.19.5+.

## Import

**account** can be imported using the `account_name` (when doing import, need to leave sensitive attributes blank), e.g.

```
$ terraform import aviatrix_account.test account_name
```
