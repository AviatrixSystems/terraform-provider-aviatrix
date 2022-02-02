---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway"
description: |-
  Creates and manages Aviatrix spoke gateways
---

# aviatrix_spoke_gateway

The **aviatrix_spoke_gateway** resource allows the creation and management of Aviatrix spoke gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws" {
  cloud_type                        = 1
  account_name                      = "my-aws"
  gw_name                           = "spoke-gw-aws"
  vpc_id                            = "vpc-abcd123"
  vpc_reg                           = "us-west-1"
  gw_size                           = "t2.micro"
  subnet                            = "10.11.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
  tags                              = {
    name = "value"
  }
}
```
```hcl
# Create an Aviatrix AWS Spoke Gateway with BGP enabled
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws" {
  cloud_type                        = 1
  account_name                      = "my-aws"
  gw_name                           = "spoke-gw-aws"
  vpc_id                            = "vpc-abcd123"
  vpc_reg                           = "us-west-1"
  gw_size                           = "t2.micro"
  subnet                            = "10.11.0.0/24"
  single_ip_snat                    = false
  enable_bgp                        = true
  manage_transit_gateway_attachment = false
  tags                              = {
    name = "value"
  }
}
```
```hcl
# Create an Aviatrix GCP Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_gcp" {
  cloud_type                        = 4
  account_name                      = "my-gcp"
  gw_name                           = "spoke-gw-gcp"
  vpc_id                            = "gcp-spoke-vpc"
  vpc_reg                           = "us-west1-b"
  gw_size                           = "n1-standard-1"
  subnet                            = "10.12.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
}
```
```hcl
# Create an Aviatrix Azure Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_azure" {
  cloud_type                        = 8
  account_name                      = "my-azure"
  gw_name                           = "spoke-gw-01"
  vpc_id                            = "vnet_name:rg_name:resource_guid"
  vpc_reg                           = "West US"
  gw_size                           = "Standard_B1ms"
  subnet                            = "10.13.0.0/24"
  zone                              = "az-1"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
}
```
```hcl
# Create an Aviatrix OCI Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_oracle" {
  cloud_type                        = 16
  account_name                      = "my-oracle"
  gw_name                           = "avtxgw-oracle"
  vpc_id                            = "vpc-oracle-test"
  vpc_reg                           = "us-ashburn-1"
  gw_size                           = "VM.Standard2.2"
  subnet                            = "10.7.0.0/16"
  manage_transit_gateway_attachment = false
  availability_domain               = aviatrix_vpc.oci_vpc.availability_domains[0]
  fault_domain                      = aviatrix_vpc.oci_vpc.fault_domains[0]
}
```
```hcl
# Create an Aviatrix AzureGov Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_azuregov" {
  cloud_type                        = 32
  account_name                      = "my-azuregov"
  gw_name                           = "spoke-gw-01"
  vpc_id                            = "vnet_name:rg_name:resource_guid"
  vpc_reg                           = "USGov Arizona"
  gw_size                           = "Standard_B1ms"
  subnet                            = "10.13.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
}
```
```hcl
# Create an Aviatrix AWSGov Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_awsgov" {
  cloud_type                        = 256
  account_name                      = "my-awsgov"
  gw_name                           = "spoke-gw-awsgov"
  vpc_id                            = "vpc-abcd123"
  vpc_reg                           = "us-gov-west-1"
  gw_size                           = "t2.micro"
  subnet                            = "10.11.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
}
```
```hcl
# Create an Aviatrix AWS China Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws_china" {
  cloud_type                        = 1024
  account_name                      = "my-aws-china"
  gw_name                           = "spoke-gw-aws-china"
  vpc_id                            = "vpc-abcd123"
  vpc_reg                           = "cn-north-1"
  gw_size                           = "t2.micro"
  subnet                            = "10.11.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
  tags                              = {
    k1 = "v1",
    k2 = "v2",
  }
}
```
```hcl
# Create an Aviatrix Azure China Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_azure" {
  cloud_type                        = 2048
  account_name                      = "my-azure-china"
  gw_name                           = "spoke-gw-01"
  vpc_id                            = "vnet_name:rg_name:resource_guid"
  vpc_reg                           = "China North"
  gw_size                           = "Standard_A0"
  subnet                            = "10.13.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
}
```
```hcl
# Create an OOB Aviatrix AWS Spoke Gateway
resource "aviatrix_spoke_gateway" "test_oob_spoke" {
  cloud_type               = 1
  account_name             = "devops-aws"
  gw_name                  = "oob-spoke"
  vpc_id                   = "vpc-abcd1234"
  vpc_reg                  = "us-west-1"
  gw_size                  = "c5.xlarge"
  enable_private_oob       = true
  subnet                   = "11.0.0.128/26"
  oob_management_subnet    = "11.0.2.0/24"
  oob_availability_zone    = "us-west-1a"
  ha_subnet                = "11.0.3.64/26"
  ha_gw_size               = "c5.xlarge"
  ha_oob_management_subnet = "11.0.0.48/28"
  ha_oob_availability_zone = "us-west-1b"
}
```
```hcl
# Create an Aviatrix Alibaba Cloud Spoke Gateway with HA enabled
resource "aviatrix_spoke_gateway" "test_spoke_gateway_alibaba" {
  cloud_type                        = 8192
  account_name                      = "devops"
  gw_name                           = "avtx-gw-1"
  vpc_id                            = "vpc-abcdef"
  vpc_reg                           = "acs-us-west-1 (Silicon Valley)"
  gw_size                           = "ecs.g5ne.large"
  subnet                            = "10.0.0.0/24"
  ha_subnet                         = "10.0.0.0/24"
  ha_gw_size                        = "ecs.g5ne.large"
  manage_transit_gateway_attachment = false
}
```
```hcl
# Create an Aviatrix AWS Top Secret Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws_top_secret" {
  cloud_type                        = 16384
  account_name                      = "my-aws-top-secret"
  gw_name                           = "spoke-gw-aws-top-secret"
  vpc_id                            = "vpc-abcd123"
  vpc_reg                           = "us-iso-east-1"
  gw_size                           = "t2.micro"
  subnet                            = "10.11.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
  tags                              = {
    k1 = "v1",
    k2 = "v2",
  }
}
```
```hcl
# Create an Aviatrix AWS Secret Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws_secret" {
  cloud_type                        = 16384
  account_name                      = "my-aws-secret"
  gw_name                           = "spoke-gw-aws-secret"
  vpc_id                            = "vpc-abcd123"
  vpc_reg                           = "us-isob-east-1"
  gw_size                           = "t2.micro"
  subnet                            = "10.11.0.0/24"
  single_ip_snat                    = false
  manage_transit_gateway_attachment = false
  tags                              = {
    k1 = "v1",
    k2 = "v2",
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently, only AWS(1), GCP(4), Azure(8), OCI(16), AzureGov(32), AWSGov(256), AWSChina(1024), AzureChina(2048), Alibaba Cloud(8192), AWS Top Secret(16384) and AWS Secret (32768) are supported.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.

!> When creating a Spoke Gateway with an Azure VNet created in Controller version 6.4 or earlier or with an Azure VNet created out of band, referencing `vpc_id` in another resource on the same apply that creates this Spoke Gateway will cause Terraform to throw an error. Please use the Spoke Gateway data source to reference the `vpc_id` of this Spoke Gateway in other resources.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Example: AWS/AWSGov/AWSChina: "vpc-abcd1234", GCP: "vpc-gcp-test", Azure/AzureGov/AzureChina: "vnet_name:rg_name:resource_guid", OCI: "vpc-oracle-test1".
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", GCP: "us-west2-a", Azure: "East US 2", OCI: "us-ashburn-1", AzureGov: "USGov Arizona", AWSGov: "us-gov-west-1, AWSChina: "cn-north-1", AzureChina: "China North", AWS Top Secret: "us-iso-east-1", AWS Secret: "us-isob-east-1".
* `gw_size` - (Required) Size of the gateway instance. Example: AWS/AWSGov/AWSChina: "t2.large", Azure/AzureGov/AzureChina: "Standard_B1s", OCI: "VM.Standard2.2", GCP: "n1-standard-1".
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane_mode).**
* `availability_domain` - (Optional) Availability domain. Required and valid only for OCI. Available as of provider version R2.19.3.
* `fault_domain` - (Optional) Fault domain. Required and valid only for OCI. Available as of provider version R2.19.3.

### HA
* `single_az_ha` (Optional) Set to true if this [feature](https://docs.aviatrix.com/Solutions/gateway_ha.html#single-az-gateway) is desired. Valid values: true, false.
* `ha_subnet` - (Optional) HA Subnet. Required if enabling HA for AWS, AWSGov, AWSChina, Azure, AzureGov, AzureChina, OCI, Alibaba Cloud, AWS Top Secret or AWS Secret gateways. Optional for GCP. Setting to empty/unsetting will disable HA. Setting to a valid subnet CIDR will create an HA gateway on the subnet. Example: "10.12.0.0/24"
* `ha_zone` - (Optional) HA Zone. Required if enabling HA for GCP gateway. Optional for Azure. For GCP, setting to empty/unsetting will disable HA and setting to a valid zone will create an HA gateway in the zone. Example: "us-west1-c". For Azure, this is an optional parameter to place the HA gateway in a specific availability zone. Valid values for Azure gateways are in the form "az-n". Example: "az-2". Available for Azure as of provider version R2.17+.
* `ha_insane_mode_az` (Optional) AZ of subnet being created for Insane Mode Spoke HA Gateway. Required for AWS, AzureGov, AWSGov, AWS Top Secret and AWS Secret if `insane_mode` is enabled and `ha_subnet` is set. Example: AWS: "us-west-1a".
* `ha_eip` - (Optional) Public IP address that you want to assign to the HA peering instance. If no value is given, a new EIP will automatically be allocated. Only available for AWS, GCP, Azure, OCI, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret.
* `ha_azure_eip_name_resource_group` - (Optional) Name of public IP Address resource and its resource group in Azure to be assigned to the HA Spoke Gateway instance. Example: "IP_Name:Resource_Group_Name". Required if `ha_eip` is set and `cloud_type` is Azure, AzureGov or AzureChina. Available as of provider version 2.20+.
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if enabling HA.
* `ha_availability_domain` - (Optional) HA gateway availability domain. Required and valid only for OCI. Available as of provider version R2.19.3.
* `ha_fault_domain` - (Optional) HA gateway fault domain. Required and valid only for OCI. Available as of provider version R2.19.3.

### Insane Mode
* `insane_mode` - (Optional) Enable [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) for Spoke Gateway. Insane Mode gateway size must be at least c5 size (AWS, AWSGov, AWS Top Secret and AWS Secret) or Standard_D3_v2 (Azure and AzureGov); for GCP only four size are supported: "n1-highcpu-4", "n1-highcpu-8", "n1-highcpu-16" and "n1-highcpu-32". If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet for AWS, Azure, AzureGov, AWSGov, AWS Top Secret and AWS Secret. Only available for AWS, GCP/OCI, Azure, AzureGov, AWSGov, AWS Top Secret and AWS Secret. Valid values: true, false. Default value: false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Spoke Gateway. Required for AWS, AWSGov, AWS Top Secret or AWS Secret if `insane_mode` is enabled. Example: AWS: "us-west-1a".

### SNAT/DNAT
* `single_ip_snat` - (Optional) Specify whether to enable Source NAT feature in "single_ip" mode on the gateway or not. Please disable AWS NAT instance before enabling this feature. Currently only supports AWS(1) and Azure(8). Valid values: true, false.

-> **NOTE:** `enable_snat` has been renamed to `single_ip_snat` in provider version R2.10. Please see notes [here](#enable_snat) for more information.

~> **NOTE:** Custom SNAT and DNAT support have been deprecated and functionality has been moved to **aviatrix_gateway_snat** and **aviatrix_gateway_dnat** respectively, in provider version R2.10. Please see notes for `snat_mode`, `snat_policy` and `dnat_policy` in the Notes section below.

### Advanced Options for BGP Spoke Gateway
* `bgp_ecmp` - (Optional) Enable Equal Cost Multi Path (ECMP) routing for the next hop. Default value: false.
* `bgp_hold_time` - (Optional) BGP hold time. Unit is in seconds. Valid values are between 12 and 360. Default value: 180.
* `bgp_polling_time` - (Optional) BGP route polling time. Unit is in seconds. Valid values are between 10 and 50. Default value: "50".
* `spoke_bgp_manual_advertise_cidrs` - (Optional) Intended CIDR list to be advertised to external BGP router. Example: "10.2.0.0/16,10.4.0.0/16". 
* `enable_active_standby` - (Optional) Enables [Active-Standby Mode](https://docs.aviatrix.com/HowTos/transit_advanced.html#active-standby). Available only with HA enabled. Valid values: true, false. Default value: false.
* `enable_active_standby_preemptive` - (Optional) Enables Preemptive Mode for Active-Standby. Available only with BGP enabled, HA enabled and Active-Standby enabled. Valid values: true, false. Default value: false.
* `local_as_number` - (Optional) Changes the Aviatrix Spoke Gateway ASN number before you setup Aviatrix Spoke Gateway connection configurations.
* `prepend_as_path` - (Optional) List of AS numbers to populate BGP AS_PATH field when it advertises to VGW or peer devices.

### Encryption
* `enable_encrypt_volume` - (Optional) Enable EBS volume encryption for Gateway. Only supports AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret providers. Valid values: true, false. Default value: false.
* `customer_managed_keys` - (Optional and Sensitive) Customer managed key ID.

### Route Customization
* `customized_spoke_vpc_routes` - (Optional) A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. It applies to this spoke gateway only. Example: "10.0.0.0/116,10.2.0.0/16".
* `filtered_spoke_vpc_routes` - (Optional) A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s routing table. It applies to this spoke gateway only. Example: "10.2.0.0/116,10.3.0.0/16".
* `included_advertised_spoke_routes` - (Optional) A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC. Example: "10.4.0.0/116,10.5.0.0/16".
* `enable_private_vpc_default_route` - (Optional) Program default route in VPC private route table. Default: false. Valid values: true or false. Available as of provider version R2.19+.
* `enable_skip_public_route_table_update` - (Optional) Skip programming VPC public route table. Default: false. Valid values: true or false. Available as of provider version R2.19+.
* `enable_auto_advertise_s2c_cidrs` - (Optional) Auto Advertise Spoke Site2Cloud CIDRs. Default: false. Valid values: true or false. Available as of provider version R2.19+.

### [Learned CIDRs Approval for BGP Spoke Gateway](https://docs.aviatrix.com/HowTos/transit_approval.html)

-> **NOTE:** `enable_learned_cidrs_approval` can be set to true only if `learned_cidrs_approval_mode` is set to 'gateway'.

* `enable_learned_cidrs_approval` - (Optional) Switch to enable/disable encrypted transit approval for BGP Spoke Gateway. Valid values: true, false. Default value: false.
* `learned_cidrs_approval_mode` - (Optional) Learned CIDRs approval mode. Either "gateway" (approval on a per-gateway basis) or "connection" (approval on a per-connection basis). Only "gateway" is supported for BGP SPOKE Gateway. Default value: "gateway". Available as of provider version R2.21+.
* `approved_learned_cidrs` - (Optional) A set of approved learned CIDRs. Only valid when `enable_learned_cidrs_approval` is set to true. Example: ["10.250.0.0/16", "10.251.0.0/16"]. Available as of provider version R2.21+.

### [Monitor Gateway Subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet)
~> **NOTE:** This feature is only available for AWS gateways.

* `enable_monitor_gateway_subnets` - (Optional) If set to true, the [Monitor Gateway Subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet) feature is enabled. Default value is false. Available in provider version R2.18+.
* `monitor_exclude_list` - (Optional) Set of monitored instance ids. Only valid when 'enable_monitor_gateway_subnets' = true. Available in provider version R2.18+.

### [Private OOB](https://docs.aviatrix.com/HowTos/private_oob.html)
* `enable_private_oob` - (Optional) Enable Private OOB feature. Only available for AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret. Valid values: true, false. Default value: false.
* `oob_management_subnet` - (Optional) OOB management subnet. Required if enabling Private OOB. Example: "11.0.2.0/24".
* `oob_availability_zone` - (Optional) OOB availability zone. Required if enabling Private OOB. Example: "us-west-1a".
* `ha_oob_management_subnet` - (Optional) HA OOB management subnet. Required if enabling Private OOB and HA. Example: "11.0.0.48/28".
* `ha_oob_availability_zone` - (Optional) HA OOB availability zone. Required if enabling Private OOB and HA. Example: "us-west-1b".

### Spot Instance
* `enable_spot_instance` - (Optional) Enable spot instance. NOT supported for production deployment.
* `spot_price` - (Optional) Price for spot instance. NOT supported for production deployment.

### Gateway Upgrade
* `software_version` - (Optional/Computed) The software version of the gateway. If set, we will attempt to update the gateway to the specified version if current version is different. If left blank, the gateway upgrade can be managed with the `aviatrix_controller_config` resource. Type: String. Example: "6.5.821". Available as of provider version R2.20.0.
* `image_version` - (Optional/Computed) The image version of the gateway. Use `aviatrix_gateway_image` data source to programmatically retrieve this value for the desired `software_version`. If set, we will attempt to update the gateway to the specified version if current version is different. If left blank, the gateway upgrades can be managed with the `aviatrix_controller_config` resource. Type: String. Example: "hvm-cloudx-aws-022021". Available as of provider version R2.20.0.
* `ha_software_version` - (Optional/Computed) The software version of the HA gateway. If set, we will attempt to update the HA gateway to the specified version if current version is different. If left blank, the HA gateway upgrade can be managed with the `aviatrix_controller_config` resource. Type: String. Example: "6.5.821". Available as of provider version R2.20.0.
* `ha_image_version` - (Optional/Computed) The image version of the HA gateway. Use `aviatrix_gateway_image` data source to programmatically retrieve this value for the desired `ha_software_version`. If set, we will attempt to update the HA gateway to the specified version if current version is different. If left blank, the gateway upgrades can be managed with the `aviatrix_controller_config` resource. Type: String. Example: "hvm-cloudx-aws-022021". Available as of provider version R2.20.0.

### Misc.

!> **WARNING:** Attribute `transit_gw` has been deprecated as of provider version R2.18.1+ and will not receive further updates. Please set `manage_transit_gateway_attachment` to false, and use the standalone `aviatrix_spoke_transit_attachment` resource instead.

* `allocate_new_eip` - (Optional) When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in Controller 4.7+. Valid values: true, false. Default: true.
* `eip` - (Optional) Required when `allocate_new_eip` is false. It uses the specified EIP for this gateway. Available in Controller 4.7+. Only available for AWS, GCP, Azure, OCI, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret.
* `azure_eip_name_resource_group` - (Optional) Name of public IP Address resource and its resource group in Azure to be assigned to the Spoke Gateway instance. Example: "IP_Name:Resource_Group_Name". Required if `allocate_new_eip` is false and `cloud_type` is Azure, AzureGov or AzureChina. Available as of provider version 2.20+.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server for Gateway. Currently only supported for AWS, Azure, AzureGov, AWSGov, AWSChina, AzureChina, Alibaba Cloud, AWS Top Secret and AWS Secret gateways. Valid values: true, false. Default value: false.
* `zone` - (Optional) Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'. Available in provider version R2.17+.
* `manage_transit_gateway_attachment` - (Optional) Enable to manage spoke-to-Aviatrix transit gateway attachments using the **aviatrix_spoke_gateway** resource with the below `transit_gw` attribute. If this is set to false, attaching this spoke to transit gateways must be done using the **aviatrix_spoke_transit_attachment** resource. Valid values: true, false. Default value: true. Available in provider R2.17+.
* `transit_gw` - (Optional) Specify the Aviatrix transit gateways to attach this spoke gateway to. Format is a comma separated list of transit gateway names. For example: "transit-gw1,transit-gw2".
* `enable_jumbo_frame` - (Optional) Enable jumbo frames for this spoke gateway. Default value is true.
* `tags` - (Optional) Map of tags to assign to the gateway. Only available for AWS, Azure, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret gateways. Allowed characters vary by cloud type but always include: letters, spaces, and numbers. AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret allow the use of any character. Azure, AzureGov and AzureChina allows the following special characters: + - = . _ : @. Example: {"key1" = "value1", "key2" = "value2"}.
* `tunnel_detection_time` - (Optional) The IPsec tunnel down detection time for the Spoke Gateway in seconds. Must be a number in the range [20-600]. The default value is set by the controller (60 seconds if nothing has been changed). **NOTE: The controller UI has an option to set the tunnel detection time for all gateways. To achieve the same functionality in Terraform, use the same TF_VAR to manage the tunnel detection time for all gateways.** Available in provider R2.19+.
* `enable_bgp` - (Optional) Enable BGP for this spoke gateway. Only available for AWS and Azure. Valid values: true, false. Default value: true. Available in provider R2.21.0+.

-> **NOTE:** `manage_transit_gateway_attachment` - If you are using/upgraded to Aviatrix Terraform Provider R2.17+, and an **aviatrix_spoke_gateway** resource was originally created with a provider version <R2.17, you must do 'terraform refresh' to update and apply the attribute's default value (true) into the state file.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `eip` - Public IP address assigned to the gateway.
* `ha_eip` - Public IP address assigned to the HA gateway.
* `security_group_id` - Security group used for the spoke gateway.
* `ha_security_group_id` - HA security group used for the spoke gateway.
* `cloud_instance_id` - Cloud instance ID of the spoke gateway.
* `private_ip` - Private IP address of the spoke gateway created.
* `ha_cloud_instance_id` - Cloud instance ID of the HA spoke gateway.
* `ha_gw_name` - Aviatrix spoke gateway unique name of HA spoke gateway.
* `ha_private_ip` - Private IP address of HA spoke gateway.

The following arguments are deprecated:

* `enable_snat` - (Optional) Specify whether enabling Source NAT feature on the gateway or not. Please disable AWS NAT instance before enabling this feature. Currently only supports AWS(1), Azure(8) and AWSGov(256). Valid values: true, false.
* `snat_mode` - (Optional) Valid values: "primary", "secondary" and "custom". Default value: "primary".
* `snat_policy` - (Optional) Policy rule applied for "snat_mode" of "custom".
  * `src_ip` - (Optional) A source IP address range where the policy rule applies.
  * `src_port` - (Optional) A source port that the policy rule applies.
  * `dst_ip` - (Optional) A destination IP address range where the policy rule applies.
  * `dst_port` - (Optional) A destination port where the policy rule applies.
  * `protocol` - (Optional) A destination port protocol where the policy rule applies.
  * `interface` - (Optional) An output interface where the policy rule applies.
  * `connection` - (Optional) Default value: "None".
  * `mark` - (Optional) A tag or mark of a TCP session where the policy rule applies.
  * `new_src_ip` - (Optional) The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.
  * `new_src_port` - (Optional) The translated destination port when all specified qualifier conditions meet. One of the rule field must be specified for this rule to take effect.
  * `exclude_rtb` - (Optional) This field specifies which VPC private route table will not be programmed with the default route entry.
* `dnat_policy` - (Optional) Policy rule applied for enabling Destination NAT (DNAT), which allows you to change the destination to a virtual address range. Currently only supports AWS(1), Azure(8), and AWSGov(256).
  * `src_ip` - (Optional) A source IP address range where the policy rule applies.
  * `src_port` - (Optional) A source port that the policy rule applies.
  * `dst_ip` - (Optional) A destination IP address range where the policy rule applies.
  * `dst_port` - (Optional) A destination port where the policy rule applies.
  * `protocol` - (Optional) A destination port protocol where the policy rule applies.
  * `interface` - (Optional) An output interface where the policy rule applies.
  * `connection` - (Optional) Default value: "None".
  * `mark` - (Optional) A tag or mark of a TCP session where the policy rule applies.
  * `new_src_ip` - (Optional) The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.
  * `new_src_port` - (Optional) The translated destination port when all specified qualifier conditions meet. One of the rule field must be specified for this rule to take effect.
  * `exclude_rtb` - (Optional) This field specifies which VPC private route table will not be programmed with the default route entry.
* `tag_list` - (Optional) Instance tag of cloud provider. Only supported for AWS, Azure, AzureGov, AWSGov, AWSChina and AzureChina. Example: ["key1:value1", "key2:value2"].
* `enable_active_mesh` - (Optional) Switch to enable/disable [Active Mesh Mode](https://docs.aviatrix.com/HowTos/activemesh_faq.html) for Spoke Gateway. Valid values: true, false. Default value: false.
* `storage_name` (Optional) Specify a storage account. Required if `cloud_type` is 2048 (AzureChina). Removed in Provider version 2.21.0+.

## Import

**spoke_gateway** can be imported using the `gw_name`, e.g.
****
```
$ terraform import aviatrix_spoke_gateway.test gw_name
```
-> **NOTE:** If `manage_transit_gateway_attachment` is set to "false", import action will also import the information of the transit gateways to which this spoke is attached to into the state file. Will need to do *terraform apply* to sync `manage_transit_gateway_attachment` to "false".


## Notes
### insane_mode
If `insane_mode` is enabled, you must specify a valid /26 CIDR segment of the VPC specified for the `subnet`. This will then create a new subnet to be used for the corresponding gateway. You cannot specify an existing /26 subnet.

### enable_snat
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a spoke gateway with `enable_snat` set to true was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also change this attribute to `single_ip_snat` in your `.tf` file.

### snat_mode & snat_policy
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a spoke gateway with `snat_mode` and `snat_policy` was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to remove attribute’s value from the state. In addition, you must transfer its corresponding values to the **aviatrix_gateway_snat** resource in your `.tf` file and perform a 'terraform import' to rectify the state file.

### dnat_policy
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a spoke gateway with `dnat_policy` was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to remove attribute’s value from the state. In addition, you must its value to its corresponding **aviatrix_gateway_dnat** resource in your `.tf` file and perform a 'terraform import' to rectify the state file.

### ha_subnet
If you are using Aviatrix Terraform Provider R2.15+, and import a Google Cloud spoke gateway with HA enabled then you must set a value for `ha_subnet` in your Terraform config.
