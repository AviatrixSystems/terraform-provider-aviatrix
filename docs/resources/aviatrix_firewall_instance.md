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

Valid `firewall_image` values:

**AWS**
1. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1
- 10.1.3
- 10.1.1
- 10.1.0
- 10.0.9
- 10.0.8-h8
- 10.0.7
- 10.0.6
- 10.0.5
- 10.0.4
- 10.0.3
2. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1 [VM-300]
- 9.1.2
- 9.1.0-h3
- 9.0.9.xfr
- 9.0.9-h1.xfr
- 9.0.6
- 9.0.5.xfr
- 9.0.3.xfr
- 9.0.1
- 8.1.15
- 8.1.9
3. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 2
- 10.1.3
- 10.1.1
- 10.1.0
- 10.0.9
- 10.0.8-h8
- 10.0.7
- 10.0.6
- 10.0.5
- 10.0.4
- 10.0.3
4. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 2 [VM-300]
- 9.0.9.xfr
- 9.0.9-h1.xfr
- 8.1.15
5. Palo Alto Networks VM-Series Next-Generation Firewall (BYOL)
- 10.1.3
- 10.1.1
- 10.1.0
- 10.0.9
- 10.0.8-h8
- 10.0.7
- 10.0.6
- 10.0.5
- 10.0.4
- 10.0.3
6. Check Point CloudGuard IaaS Next-Gen Firewall w. Threat Prevention & SandBlast BYOL
- R81.10-335.940
- R81-392.938
- R81-392.917
- R81-392.897
- R81-392.840
- R81-392.836
- R81-392.788
- R81-392.753
- R81-392.718
- R80.40-294.936
7. Check Point CloudGuard IaaS Next-Gen Firewall with Threat Prevention
- R81.10-335.986
- R81.10-335.940
- R81.10-335.883
- R81-392.983
- R81-392.938
- R81-392.917
- R81-392.897
- R81-392.840
- R81-392.836
- R81-392.788
8. Check Point CloudGuard IaaS All-In-One
- R81.10-335.986
- R81.10-335.942
- R81.10-335.883
- R81-392.983
- R81-392.942
- R81-392.715
- R80.40-294.983
- R80.40-294.947
- R80.40-294.942
- R80.40-294.774
9. Fortinet FortiGate Next-Generation Firewall
- (7.0.5)
- (7.0.3)
- (6.4.8)
- (6.4.7)
10. Fortinet FortiGate (BYOL) Next-Generation Firewall
- (7.0.5)
- (7.0.3)
- (6.4.8)
- (6.4.7)

**Azure**
1. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1
- 9.1.0
- 9.0.9
- 9.0.6
- 9.0.4
- 9.0.1
2. Palo Alto Networks VM-Series Next-Generation Firewall Bundle 2
- 9.1.0
- 9.0.9
- 9.0.6
- 9.0.4
- 9.0.1
3. Palo Alto Networks VM-Series Next-Generation Firewall (BYOL)
- 9.1.0
- 9.0.9
- 9.0.6
- 9.0.4
- 9.0.1
4. Palo Alto Networks VM-Series Flex Next-Generation Firewall Bundle 1
- 10.1.4
- 10.1.3
- 10.1.2
- 10.1.0
- 10.0.9
5. Palo Alto Networks VM-Series Flex Next-Generation Firewall Bundle 2
- 10.1.4
- 10.1.3
- 10.1.2
- 10.1.0
- 10.0.9
6. Palo Alto Networks VM-Series Flex Next-Generation Firewall (BYOL)
- 10.1.4
- 10.1.3
- 10.1.0
- 10.0.9
- 10.0.7
7. Check Point CloudGuard IaaS Single Gateway R80.30 - Bring Your Own License
- 8030.900273.0978 
- 8030.900273.0950 
- 8030.900273.0907 
- 8030.900273.0819 
- 8030.900273.0801
8. Check Point CloudGuard IaaS Single Gateway R80.30 - Pay As You Go (NGTP)
- 8030.900273.0590
- 8030.900273.0562
9. Check Point CloudGuard IaaS Single Gateway R80.30 - Pay As You Go (NGTX)
- 8030.900273.0590
- 8030.900273.0562
10. Check Point CloudGuard IaaS Single Gateway R80.40 - Bring Your Own License
- 8040.900294.0978
- 8040.900294.0907
- 8040.900294.0801
- 8040.900294.0728
- 8040.900294.0682
11. Check Point CloudGuard IaaS Single Gateway R80.40 - Pay As You Go (NGTP)
- 8040.900294.0978
- 8040.900294.0907
- 8040.900294.0801
- 8040.900294.0728
- 8040.900294.0682
12. Check Point CloudGuard IaaS Single Gateway R80.40 - Pay As You Go (NGTX)
- 8040.900294.0978
- 8040.900294.0907
- 8040.900294.0801
- 8040.900294.0728
- 8040.900294.0682
13. Check Point CloudGuard IaaS Standalone (gateway + management) R80.40 - Bring Your Own License
- 8040.900294.0978
- 8040.900294.0752
- 8040.900294.0586
14. Fortinet FortiGate (BYOL) Next-Generation Firewall
- 7.0.4
- 7.0.3
- 7.0.2
- 7.0.1
- 7.0.0 
15. Fortinet FortiGate (PAYG) Next-Generation Firewall
- 6.0.4
- 6.0.3
- 6.0.02
- 5.6.6
- 5.6.5
16. Fortinet FortiGate (PAYG_20190624) Next-Generation Firewall Latest Release
- 7.0.4
- 7.0.3
- 7.0.2
- 7.0.1
- 7.0.0

**GCP**
1. Palo Alto Networks VM-Series Next-Generation Firewall BUNDLE1
- 9.0.9
- 8.1.15
- 9.1.2
- 9.0.4
- 8.1.9
- 9.0.1
- 814
- 810
2. Palo Alto Networks VM-Series Next-Generation Firewall BUNDLE2
- 9.0.9
- 8.1.15
- 9.1.2
- 9.0.4
- 8.1.9
- 9.0.1
- 814
- 810
3. Palo Alto Networks VM-Series Next-Generation Firewall BYOL
- 9.1.3
- 9.0.9
- 8.1.15
- 9.0.8
- 9.1.2
- 9.0.6
- 9.0.4
- 8.1.9
- 9.0.1
- 814
4. Fortinet FortiGate Next-Generation Firewall
- 7.0.5
- 7.0.3
- 6.4.8
- 7.0.2
- 6.4.7
- 7.0.1
- 6.4.6
- 7.0.0
- 6.4.5
- 6.4.4
5. Fortinet FortiGate Next-Generation Firewall (BYOL)
- 7.0.5
- 7.0.3
- 6.4.8
- 7.0.2
- 6.4.7
- 7.0.1
- 6.4.6
- 7.0.0
- 6.4.5
- 6.4.4
6. Check Point CloudGuard IaaS Firewall & Threat Prevention (Gateway only)
- R80.40-294.688
- R80.40-294.904
- R81-344.139
- R81-351.146
- R81-392.710
- R81-392.751
- R81-392.906
- R81-392.915
- R81.10-335.883
- R81.10-335.985
7. Check Point CloudGuard IaaS Firewall & Threat Prevention (Gateway only) (BYOL)
- R80.40-294.904
- R81-383.704
- R81-386.706
- R81-390.708
- R81-392.710
- R81-392.751
- R81-392.906
- R81-392.915
- R81.10-335.883
- R81.10-335.985
8. Check Point CloudGuard IaaS Firewall & Threat Prevention (Standalone)
- R80.40-294.127
- R80.40-294.587
- R80.40-294.735
- R80.40-294.759
- R81-392.758
- R81-392.957
- R81-394.735
- R81.10-335.883
9. Check Point CloudGuard IaaS Firewall & Threat Prevention (Standalone) (BYOL)
- R80.40-294.127
- R80.40-294.687
- R80.40-294.759
- R81-344.139
- R81-351.146
- R81-373.150
- R81-386.706
- R81-392.710
- R81-392.957
- R81.10-335.883

**OCI**
1. Palo Alto Networks VM-Series Bundle1 - 4 OCPUs
- PA-VM-10.0.9
- PA-VM-10.0.6
- PA-VM-10.1.0
- PA-VM-10.0.3
2. Palo Alto Networks VM-Series Bundle1 - 8 OCPUs
- PA-VM-10.0.9
- PA-VM-10.0.6
- PA-VM-10.1.0
- PA-VM-10.0.3-012120210118
3. Palo Alto Networks VM-Series Bundle2 - 4 OCPUs
- PA-VM-10.0.9
- PA-VM-10.0.6
- PA-VM-10.1.0
- PA-VM-10.0.3
4. Palo Alto Networks VM-Series Bundle2 - 8 OCPUs
- PA-VM-10.0.9
- PA-VM-10.0.6
- PA-VM-10.1.0
- PA-VM-10.0.3
5. Palo Alto Networks VM-Series Next Generation Firewall
- 10.0.9
- 10.1.3
- 10.0.6
- 10.1.0-latest
- 10.0.4
- 10.0.3
- 9.1.6
- 10.0.0
- 9.1.3-070120202138
6. CloudGuard Next-Gen Firewall w\/ Threat Prevention - 4 OCPUs
- R80.40_rev1.3
- R81.10_rev1.1
- R81_rev1.2
- R81.10_rev1.0
- R81_rev1.1
- R81_rev1.0
- R80.40_rev1.2
- R80.40_rev1.1
- R80.40_rev1.0-042420202345
7. CloudGuard Next-Gen Firewall w\/ Threat Prevention and SandBlast - 4 OCPUs
- R80.40_rev1.3
- R81_rev1.2
- R81.10_rev1.1
- R81.10_rev1.0
- R81_rev1.0
- R80.40_rev1.2
8. CloudGuard Next-Gen Firewall with Threat Prevention and SandBlast - BYOL
- R80.40_rev1.4
- R81.10_rev1.1
- R81_rev1.2
- R80.40_rev1.3
- R81.10_rev1.0
- R81_rev1.1
- R80.30_rev_1.4

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
