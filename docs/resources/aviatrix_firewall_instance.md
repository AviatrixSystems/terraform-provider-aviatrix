---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_instance"
description: |-
  Creates and deletes Aviatrix Firewall Instances
---

# aviatrix_firewall_instance

The **aviatrix_firewall_instance** resource allows the creation and management of Aviatrix Firewall Instances.

This resource is used in [Aviatrix FireNet](https://docs.aviatrix.com/HowTos/firewall_network_faq.html) and [Aviatrix Transit FireNet](https://docs.aviatrix.com/HowTos/transit_firenet_faq.html) solutions, in conjunction with other resources that may include, and are not limited to: **firenet**, **firewall_instance_association**, **aws_tgw** and **transit_gateway** resources.

## Example Usage

```hcl
# Create an Aviatrix Firewall Instance
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id            = "vpc-032005cc371"
  firenet_gw_name   = "avx-firenet-gw"
  firewall_name     = "avx-firewall-instance"
  firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_size     = "m5.xlarge"
  management_subnet = "10.4.0.16/28"
  egress_subnet     = "10.4.0.32/28"
}
```
```hcl
# Create an Aviatrix Firewall Instance with Native GWLB Enabled VPC
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id            = "vpc-032005cc371"
  firewall_name     = "avx-firewall-instance"
  firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_size     = "m5.xlarge"
  management_subnet = "10.4.0.16/28"
  egress_subnet     = "10.4.0.32/28"
  zone              = "us-east-1a"
}
```
```hcl
# Create an Aviatrix Firewall Instance on GCP
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id                 = format("%s~-~%s", aviatrix_transit_gateway.test_transit_gateway.vpc_id, aviatrix_account.gcp.gcloud_project_id)
  firenet_gw_name        = aviatrix_transit_gateway.test_transit_gateway.gw_name
  firewall_name          = "gcp-firewall-instance"
  firewall_image         = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_image_version = "9.0.9"
  firewall_size          = "n1-standard-4"
  management_vpc_id      = aviatrix_vpc.management_vpc.vpc_id
  management_subnet      = format("%s~~%s~~%s", aviatrix_vpc.management_vpc.subnets[0].cidr, aviatrix_vpc.management_vpc.subnets[0].region, aviatrix_vpc.management_vpc.subnets[0].name)
  egress_vpc_id          = aviatrix_vpc.egress_vpc.vpc_id
  egress_subnet          = format("%s~~%s~~%s", aviatrix_vpc.egress_vpc.subnets[0].cidr, aviatrix_vpc.egress_vpc.subnets[0].region, aviatrix_vpc.egress_vpc.subnets[0].name)
  zone                   = aviatrix_transit_gateway.test_transit_gateway.vpc_reg
}
```

## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Required) VPC ID of the Security VPC. For GCP, `vpc_id` must be in the form vpc_id~-~gcloud_project_id.
* `firenet_gw_name` - (Optional) Name of the primary FireNet gateway. **Required for all FireNet deployments that do not utilize the TGW-Integrated FireNet with AWS Native GWLB VPC.**
* `firewall_name` - (Required) Name of the firewall instance to be created.
* `firewall_image` - (Required) One of the AWS/Azure/GCP AMIs from various vendors such as Palo Alto Networks.
* `firewall_image_id` - (Optional) Firewall image ID. Applicable to AWS and Azure only. For AWS, please use AMI ID. For Azure, the format is “Publisher:Offer:Plan:Version”. Available as of provider version R2.19+.
* `firewall_size` - (Required) Instance size of the firewall. Example: "m5.xlarge".  
* `management_subnet` - (Optional) Management Interface Subnet. Select the subnet whose name contains “gateway and firewall management”. For GCP, `management_subnet` must be in the form `cidr~~region~~name`. Required for Palo Alto Networks VM-Series and OCI Check Point firewalls. Otherwise, it must be empty.
* `egress_subnet` - (Required) Egress Interface Subnet. Select the subnet whose name contains “FW-ingress-egress”. For GCP, `egress_subnet` must be in the form `cidr~~region~~name`.
* `firewall_image_version` - (Optional) Version of firewall image. If not specified, Controller will automatically select the latest version available.
* `zone` - (Optional) Availability Zone. Required if creating a Firewall Instance with a Native AWS GWLB-enabled VPC. Applicable to AWS, Azure, and GCP only. Available as of provider version R2.17+.
* `management_vpc_id` - (Optional) Management VPC ID. Only used for GCP firewall. Required for Palo Alto Networks VM-Series, and required to be empty for Check Point or Fortinet series. Available as of provider version R2.18.1+.
* `egress_vpc_id` - (Optional) Egress VPC ID. Required for GCP. Available as of provider version R2.18.1+.
* `availability_domain` - (Optional) Availability domain. Required and valid only for OCI. Available as of provider version R2.19.3.
* `fault_domain` - (Optional) Fault domain. Required and valid only for OCI. Available as of provider version R2.19.3.

-> **NOTE:** Please use the data source `aviatrix_firewall_instance_images` to get the information for `firewall_image`, `firewall_image_version` and `firewall_size`.

### Authentication method
* `key_name`- (Optional) Applicable to AWS deployment only. AWS Key Pair name. If not provided a Key Pair will be generated.
* `username`- (Optional) Applicable to Azure or AzureGov deployment only. "admin" as a username is not accepted.
* `password`- (Optional) Applicable to Azure or AzureGov deployment only.
* `ssh_public_key` - (Optional) Applicable to Azure or AzureGov deployment only.

### Advanced Options
* `iam_role` - (Optional) Only available for AWS. In advanced mode, create an IAM Role on the AWS account that launched the FireNet gateway. Create a policy to attach to the role. The policy is to allow access to "Bootstrap Bucket".
* `bootstrap_bucket_name`- (Optional) Only available for AWS and GCP. For GCP, only Palo Alto Networks VM-Series deployment can use this attribute. In advanced mode, specify a bootstrap bucket name where the initial configuration and policy file is stored.
* `bootstrap_storage_name` - (Optional) Advanced option. Bootstrap storage name. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series/Fortinet Series deployment only. Available as of provider version R2.17.1+.
* `storage_access_key` - (Optional) Advanced option. Storage access key. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series deployment only. Available as of provider version R2.17.1+.
* `file_share_folder` - (Optional) Advanced option. File share folder. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series deployment only. Available as of provider version R2.17.1+.
* `share_directory` - (Optional) Advanced option. Share directory. Applicable to Azure or AzureGov and Palo Alto Networks VM-Series deployment only. Available as of provider version R2.17.1+.
* `sic_key` - (Optional) Advanced option. Sic key. Applicable to Check Point Series deployment only.
* `container_folder` - (Optional) Advanced option. Container folder. Applicable to Azure or AzureGov and Fortinet Series deployment only.
* `sas_url_config` - (Optional) Advanced option. SAS URL Config. Applicable to Azure or AzureGov and Fortinet Series deployment only.
* `sas_url_license` - (Optional) Advanced option. SAS URL License. Applicable to Azure or AzureGov and Fortinet Series deployment only.
* `user_data` - (Optional) Advanced option. User Data. Applicable to Check Point Series and Fortinet Series deployment only. Type: String.

### Misc.
* `tags` - (Optional) Mapping of key value pairs of tags for a firewall instance. Only available for AWS, AWSGov, GCP and Azure firewall instances. For AWS, AWSGov and Azure allowed characters are: letters, spaces, and numbers plus the following special characters: + - = . _ : @. For GCP allowed characters are: lowercase letters, numbers, "-" and "_". Example: {"key1" = "value1", "key2" = "value2"}.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `instance_id`- ID of the firewall instance created.
* `lan_interface`- ID of Lan Interface created.
* `management_interface`- ID of Management Interface created.
* `egress_interface`- ID of Egress Interface created.
* `public_ip`- Management Public IP.
* `cloud_type` - Cloud Type.
* `gcp_vpc_id` - GCP Only. The current VPC ID.

## Import

**firewall_instance** can be imported using the `instance_id`. For Azure or AzureGov FireNet instances, the value will be the `firewall_name` concatenated with a ":" and the Resource Group of the `vpc_id` set for that instance. e.g.

```
$ terraform import aviatrix_firewall_instance.test instance_id
```
