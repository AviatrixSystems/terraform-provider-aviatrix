---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway"
description: |-
  Creates and manages the Aviatrix Transit Network gateways
---

# aviatrix_transit_gateway

The **aviatrix_transit_gateway** resource allows the creation and management of [Aviatrix Transit Network](https://docs.aviatrix.com/HowTos/transitvpc_faq.html#) gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
  cloud_type               = 1
  account_name             = "devops_aws"
  gw_name                  = "transit"
  vpc_id                   = "vpc-abcd1234"
  vpc_reg                  = "us-east-1"
  gw_size                  = "t2.micro"
  subnet                   = "10.1.0.0/24"
  ha_subnet                = "10.1.0.0/24"
  ha_gw_size               = "t2.micro"
  tags                     = {
    name = "value"
  }
  enable_hybrid_connection = true
  connected_transit        = true
}
```
```hcl
# Create an Aviatrix GCP Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_gcp" {
  cloud_type   = 4
  account_name = "devops-gcp"
  gw_name      = "avtxgw-gcp"
  vpc_id       = "vpc-gcp-test~-~project-id"
  vpc_reg      = "us-west2-a"
  gw_size      = "n1-standard-1"
  subnet       = "10.8.0.0/16"
  ha_zone      = "us-west2-b"
  ha_subnet    = "10.8.0.0/16"
  ha_gw_size   = "n1-standard-1"
}
```
```hcl
# Create an Aviatrix GCP Transit Network Gateway with HA enabled and BGP over LAN enabled
resource "aviatrix_transit_gateway" "test_transit_gateway_gcp" {
  cloud_type        = 4
  account_name      = "devops_gcp"
  gw_name           = "avtxgw-gcp"
  vpc_id            = "vpc-gcp-test"
  vpc_reg           = "us-west1-a"
  gw_size           = "n1-highcpu-16"
  subnet            = "10.1.0.0/24"
  ha_subnet         = "10.1.0.0/24"
  ha_zone           = "us-west1-b"
  ha_gw_size        = "n1-highcpu-16"
  bgp_lan_interfaces {
    vpc_id = "gcp-vpc-bgp"
    subnet = "172.16.0.0/16"
  }
  bgp_lan_interfaces {
    vpc_id = "gcp-vpc-bgp1"
    subnet = "173.16.0.0/16"
  }
  bgp_lan_interfaces {
    vpc_id = "gcp-vpc-bgp2"
    subnet = "174.16.0.0/16"
  }
  ha_bgp_lan_interfaces {
    vpc_id = "gcp-vpc-bgp3"
    subnet = "175.16.0.0/16"
  }
  ha_bgp_lan_interfaces {
    vpc_id = "gcp-vpc-bgp4"
    subnet = "176.16.0.0/16"
  }
  ha_bgp_lan_interfaces {
    vpc_id = "gcp-vpc-bgp5"
    subnet = "177.16.0.0/16"
  }
  connected_transit = true
}
```
```hcl
# Create an Aviatrix Azure Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_azure" {
  cloud_type        = 8
  account_name      = "devops_azure"
  gw_name           = "transit"
  vpc_id            = "vnet_name:rg_name:resource_guid"
  vpc_reg           = "West US"
  gw_size           = "Standard_B1ms"
  subnet            = "10.30.0.0/24"
  zone              = "az-1"
  ha_subnet         = "10.30.0.0/24"
  ha_zone           = "az-2"
  ha_gw_size        = "Standard_B1ms"
  connected_transit = true
}
```
```hcl
# Create an Aviatrix Azure Transit Network Gateway with HA enabled and BGP over LAN enabled with multiple interfaces
resource "aviatrix_transit_gateway" "test_transit_gateway_azure" {
  cloud_type                  = 8
  account_name                = "devops_azure"
  gw_name                     = "transit"
  vpc_id                      = "vnet_name:rg_name:resource_guid"
  vpc_reg                     = "West US"
  gw_size                     = "Standard_B1ms"
  subnet                      = "10.30.0.0/24"
  zone                        = "az-1"
  ha_subnet                   = "10.30.0.0/24"
  ha_zone                     = "az-2"
  ha_gw_size                  = "Standard_B1ms"
  connected_transit           = true
  learned_cidrs_approval_mode = "connection"
  single_az_ha                = true
  enable_bgp_over_lan         = true
  bgp_lan_interfaces_count    = 2
}
```
```hcl
# Create an Aviatrix OCI Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_oracle" {
  cloud_type          = 16
  account_name        = "devops-oracle"
  gw_name             = "avtxgw-oracle"
  vpc_id              = "ocid1.vcn.oc1.iad.aaaaaaaaba3pv6wkcr4jqae5f44n2b2m2yt2j6rx32uzr4h25vqstifsfdsq"
  vpc_reg             = "us-ashburn-1"
  gw_size             = "VM.Standard2.2"
  subnet              = "10.7.0.0/16"
  availability_domain = aviatrix_vpc.oci_vpc.availability_domains[0]
  fault_domain        = aviatrix_vpc.oci_vpc.fault_domains[0]
}
```
```hcl
# Create an Aviatrix AzureGov Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_azuregov" {
  cloud_type        = 32
  account_name      = "devops_azuregov"
  gw_name           = "transit"
  vpc_id            = "vnet_name:rg_name:resource_guid"
  vpc_reg           = "USGov Arizona"
  gw_size           = "Standard_B1ms"
  subnet            = "10.30.0.0/24"
  ha_subnet         = "10.30.0.0/24"
  ha_gw_size        = "Standard_B1ms"
  connected_transit = true
}
```
```hcl
# Create an Aviatrix AWSGov Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_awsgov" {
  cloud_type               = 256
  account_name             = "devops_awsgov"
  gw_name                  = "transit"
  vpc_id                   = "vpc-abcd12345678"
  vpc_reg                  = "us-gov-east-1"
  gw_size                  = "t2.micro"
  subnet                   = "10.1.0.0/24"
  ha_subnet                = "10.1.0.0/24"
  ha_gw_size               = "t2.micro"
  enable_hybrid_connection = true
  connected_transit        = true
}
```
```hcl
# Create an Aviatrix AWS China Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_aws_china" {
  cloud_type        = 1024
  account_name      = "devops_aws_china"
  gw_name           = "transit"
  vpc_id            = "vpc-abcd12345678"
  vpc_reg           = "cn-north-1"
  gw_size           = "t2.micro"
  subnet            = "10.1.0.0/24"
  ha_subnet         = "10.1.0.0/24"
  ha_gw_size        = "t2.micro"
  tags              = {
    name  = "value",
    name1 = "value1",
    name2 = "value2",
  }
  connected_transit = true
}
```
```hcl
# Create an Aviatrix Azure China Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_azure_china" {
  cloud_type   = 2048
  account_name = "devops_azure_china"
  gw_name      = "transit"
  vpc_id       = "vnet_name:rg_name:resource_guid"
  vpc_reg      = "China North"
  gw_size      = "Standard_A0"
  subnet       = "10.30.0.0/24"
  ha_subnet    = "10.30.0.0/24"
  ha_zone      = "az-2"
  ha_gw_size   = "Standard_A0"
}
```
```hcl
# Create an OOB Aviatrix AWS Transit Network Gateway
resource "aviatrix_transit_gateway" "test_oob_transit" {
  cloud_type               = 1
  account_name             = "devops-aws"
  gw_name                  = "oob-transit"
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
# Create an Aviatrix Alibaba Cloud Transit Network Gateway with HA enabled
resource "aviatrix_transit_gateway" "test_transit_gateway_alibaba" {
  cloud_type   = 8192
  account_name = "devops"
  gw_name      = "avtx-gw-1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "acs-us-west-1 (Silicon Valley)"
  gw_size      = "ecs.g5ne.large"
  subnet       = "10.0.0.0/24"
  ha_subnet    = "10.0.0.0/24"
  ha_gw_size   = "ecs.g5ne.large"
}
```
```hcl
# Create an Aviatrix AWS Top Secret Region Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_aws_top_secret" {
  cloud_type        = 16384
  account_name      = "devops_aws_top_secret"
  gw_name           = "transit"
  vpc_id            = "vpc-abcd12345678"
  vpc_reg           = "us-iso-east-1"
  gw_size           = "t2.micro"
  subnet            = "10.1.0.0/24"
  ha_subnet         = "10.1.0.0/24"
  ha_gw_size        = "t2.micro"
  tags              = {
    name  = "value",
    name1 = "value1",
    name2 = "value2",
  }
  connected_transit = true
}
```
```hcl
# Create an Aviatrix AWS Secret Region Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_aws_secret" {
  cloud_type        = 32768
  account_name      = "devops_aws_secret"
  gw_name           = "transit"
  vpc_id            = "vpc-abcd12345678"
  vpc_reg           = "us-isob-east-1"
  gw_size           = "t2.micro"
  subnet            = "10.1.0.0/24"
  ha_subnet         = "10.1.0.0/24"
  ha_gw_size        = "t2.micro"
  tags              = {
    name  = "value",
    name1 = "value1",
    name2 = "value2",
  }
  connected_transit = true
}
```

## Argument Reference

The following arguments are supported:

### Required
* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently only AWS(1), GCP(4), Azure(8), OCI(16), AzureGov(32), AWSGov(256), AWSChina(1024), AzureChina(2048), Alibaba Cloud(8192), AWS Top Secret(16384) and AWS Secret (32768) are supported.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.

!> When creating a Transit Gateway with an Azure VNet created in Controller version 6.4 or earlier or with an Azure VNet created out of band, referencing `vpc_id` in anothe resource on the same apply that creates this Transit Gateway will cause Terraform to throw an error. Please use the Transit Gateway data source to reference the `vpc_id` of this Transit Gateway in other resources.

~> As of Provider version R2.21.2+, the `vpc_id` of an OCI VCN has been changed from its name to its OCID.

!> As of Provider version R.22.0+, the `vpc_id` of a GCP VPC has been updated to include the project ID, e.g. vpc_name~-~project_id. When creating a Transit Gateway using the old format, referencing `vpc_id` in another resource on the same apply that creates this Transit Gateway will cause Terraform to throw an error. Please use the Transit Gateway data source to reference the `vpc_id` of this Transit Gateway in other resources.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Example: AWS/AWSGov/AWSChina: "vpc-abcd1234", GCP: "vpc-gcp-test~-~project-id", Azure/AzureGov/AzureChina: "vnet_name:rg_name:resource_guid", OCI: "ocid1.vcn.oc1.iad.aaaaaaaaba3pv6wkcr4jqae5f44n2b2m2yt2j6rx32uzr4h25vqstifsfdsq".
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", GCP: "us-west2-a", Azure: "East US 2", OCI: "us-ashburn-1", AzureGov: "USGov Arizona", AWSGov: "us-gov-west-1", AWSChina: "cn-north-1", AzureChina: "China North", AWS Top Secret: "us-iso-east-1", AWS Secret: "us-isob-east-1".
* `gw_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", Azure/AzureGov: "Standard_B1s", OCI: "VM.Standard2.2", GCP: "n1-standard-1", AWSGov: "t2.large", AWSChina: "t2.large", AzureChina: "Standard_A0".
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane_mode).**
* `availability_domain` - (Optional) Availability domain. Required and valid only for OCI. Available as of provider version R2.19.3.
* `fault_domain` - (Optional) Fault domain. Required and valid only for OCI. Available as of provider version R2.19.3.

### HA
* `single_az_ha` (Optional) Set to true if this [feature](https://docs.aviatrix.com/Solutions/gateway_ha.html#single-az-gateway) is desired. Valid values: true, false.
* `ha_subnet` - (Optional) HA Subnet CIDR. Required only if enabling HA for AWS, Azure, AzureGov, AWSGov, AWSChina, AzureChina, OCI, Alibaba Cloud, AWS Top Secret or AWS Secret gateways. Optional for GCP. Setting to empty/unsetting will disable HA. Setting to a valid subnet CIDR will create an HA gateway on the subnet. Example: "10.12.0.0/24".
* `ha_zone` - (Optional) HA Zone. Required if enabling HA for GCP gateway. Optional if enabling HA for Azure gateway. For GCP, setting to empty/unsetting will disable HA and setting to a valid zone will create an HA gateway in the zone. Example: "us-west1-c". For Azure, this is an optional parameter to place the HA gateway in a specific availability zone. Valid values for Azure gateways are in the form "az-n". Example: "az-2". Available for Azure as of provider version R2.17+.
* `ha_insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit HA Gateway. Required for AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret if `insane_mode` is enabled and `ha_subnet` is set. Example: AWS: "us-west-1a".
* `ha_eip` - (Optional) Public IP address that you want to assign to the HA peering instance. If no value is given, a new EIP will automatically be allocated. Only available for AWS, GCP, Azure, OCI, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret.
* `ha_azure_eip_name_resource_group` - (Optional) Name of public IP Address resource and its resource group in Azure to be assigned to the HA Transit Gateway instance. Example: "IP_Name:Resource_Group_Name". Required if `ha_eip` is set and `cloud_type` is Azure, AzureGov or AzureChina. Available as of provider version 2.20+.
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if enabling HA. Example: "t2.micro".
* `ha_availability_domain` - (Optional) HA gateway availability domain. Required and valid only for OCI. Available as of provider version R2.19.3.
* `ha_fault_domain` - (Optional) HA gateway fault domain. Required and valid only for OCI. Available as of provider version R2.19.3.

### Insane Mode
* `insane_mode` - (Optional) Specify true for [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) high performance gateway. Insane Mode gateway size must be at least c5 size (AWS, AWSGov, AWS China, AWS Top Secret and AWS Secret) or Standard_D3_v2 (Azure and AzureGov); for GCP only four size are supported: "n1-highcpu-4", "n1-highcpu-8", "n1-highcpu-16" and "n1-highcpu-32". If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet for AWS, Azure, AzureGov, AWSGov, AWS Top Secret and AWS Secret. Only available for AWS, GCP/OCI, Azure, AzureGov, AWSGov, AWS Top Secret and AWS Secret. Valid values: true, false. Default value: false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit Gateway. Required for AWS, AWSGov, AWS China, AWS Top Secret or AWS Secret if `insane_mode` is enabled. Example: AWS: "us-west-1a".

### SNAT
* `single_ip_snat` - (Optional) Enable "single_ip" mode Source NAT for this container. Valid values: true, false. **NOTE: Please see notes [here](#enable_snat) in regards to changes to this argument in R2.10.**

### Segmentation
* `enable_segmentation` - (Optional) Enable transit gateway for segmentation. Valid values: true, false. Default: false.

### Advanced Options
* `connected_transit` - (Optional) Specify Connected Transit status. If enabled, it allows spokes to run traffics to other spokes via transit gateway. Valid values: true, false. Default value: false.
* `enable_advertise_transit_cidr` - (Optional) Switch to enable/disable advertise transit VPC network CIDR for a VGW connection. Available as of R2.6. **NOTE: If previously enabled through vgw_conn resource prior to provider version R2.6, please see notes [here](#cidr-advertising).**
* `bgp_manual_spoke_advertise_cidrs` - (Optional) Intended CIDR list to be advertised to external BGP router. Example: "10.2.0.0/16,10.4.0.0/16". Available as of R2.6. **NOTE: If previously enabled through vgw_conn resource prior to provider version R2.6, please see notes [here](#cidr-advertising).**
* `enable_hybrid_connection` - (Optional) Sign of readiness for AWS TGW connection. Only supported for AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret. Example: false.
* `enable_firenet` - (Optional) Set to true to use gateway for legacy [AWS TGW-based FireNet](https://docs.aviatrix.com/HowTos/firewall_network_faq.html) connection. Valid values: true, false. Default value: false. **NOTE: If previously using an older provider version R2.5 where attribute name was `enable_firenet_interfaces`, please see notes [here](#enable_firenet-1).**
* `enable_transit_summarize_cidr_to_tgw` - (Optional) Enable summarize CIDR to TGW. Valid values: true, false. Default value: false.
* `enable_active_standby` - (Optional) Enables [Active-Standby Mode](https://docs.aviatrix.com/HowTos/transit_advanced.html#active-standby). Available only with HA enabled. Valid values: true, false. Default value: false. Available in provider version R2.17.1+.
* `enable_active_standby_preemptive` - (Optional) Enables Preemptive Mode for Active-Standby. Available only with BGP enabled, HA enabled and Active-Standby enabled. Valid values: true, false. Default value: false.
* `bgp_polling_time` - (Optional) BGP route polling time. Unit is in seconds. Valid values are between 10 and 50. Default value: "50".
* `bgp_hold_time` - (Optional) BGP hold time. Unit is in seconds. Valid values are between 12 and 360. Default value: 180.
* `prepend_as_path` - (Optional) List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.
* `local_as_number` - (Optional) Changes the Aviatrix Transit Gateway ASN number before you setup Aviatrix Transit Gateway connection configurations.
* `bgp_ecmp` - (Optional) Enable Equal Cost Multi Path (ECMP) routing for the next hop. Default value: false.
* `enable_multi_tier_transit` - (Optional) Enable Multi-tier Transit mode on transit gateway. When enabled, transit gateway will propagate routes it receives from its transit peering peer to other transit peering peers. `local_as_number` is required. Default value: false. Available as of provider version R2.19+.
* `enable_s2c_rx_balancing` - (Optional) Enable S2C receive packet CPU re-balancing on transit gateway. Valid values: true, false. Default value: false. Available in provider version R2.21.2+.
* `enable_preserve_as_path` - (Optional) Enable preserve as_path when advertising manual summary cidrs on transit gateway. Valid values: true, false. Default value: false. Available as of provider version R.2.22.1+
  },

-> **NOTE:** Enabling FireNet will automatically enable hybrid connection. If `enable_firenet` is set to true, please set `enable_hybrid_connection` to true in the respective **aviatrix_transit_gateway** as well.

* `enable_transit_firenet` - (Optional) Set to true to use gateway for [Transit FireNet](https://docs.aviatrix.com/HowTos/transit_firenet_faq.html) connection. Valid values: true, false. Default value: false. Available in provider version R2.12+.
* `lan_vpc_id` - (Optional) LAN VPC ID. Only valid when enabling Transit FireNet on GCP. Available as of provider version R2.18.1+.
* `lan_private_subnet` - (Optional) LAN Private Subnet. Only valid when enabling Transit FireNet on GCP. Available as of provider version R2.18.1+.
* `enable_egress_transit_firenet` - (Optional) Enable [Egress Transit FireNet](https://docs.aviatrix.com/HowTos/transit_firenet_workflow.html#b-enable-transit-firenet-on-aviatrix-egress-transit-gateway). Valid values: true, false. Default value: false. Available in provider version R2.16.3+.

-> **NOTE:** Enabling or disabling `enable_gateway_load_balancer` requires that the FireNet interfaces also be disabled or enabled. For example, if the transit gateway currently has `enable_firenet` = true and `enable_gateway_load_balancer` = false, to enable `enable_gateway_load_balancer` you would first set `enable_firenet` = false and apply to disable the FireNet interfaces. Then you would set `enable_firenet` = true and `enable_gateway_load_balancer` = true and apply to reach the desired configuration.

* `enable_gateway_load_balancer` - (Optional) Enable FireNet interfaces with AWS Gateway Load Balancer. Only valid when `enable_firenet` or `enable_transit_firenet` are set to true and `cloud_type` = 1 (AWS). Currently, AWS Gateway Load Balancer is only supported in AWS regions: us-west-2, us-east-1, eu-west-1, ap-southeast-2 and sa-east-1. Valid values: true or false. Default value: false. Available as of provider version R2.18+.

### BGP over LAN
* `enable_bgp_over_lan` - (Optional) Pre-allocate a network interface(eth4) for "BGP over LAN" functionality. Must be enabled to create a BGP over LAN `aviatrix_transit_external_device_conn` resource with this Transit Gateway. Only valid for GCP (4), Azure (8), AzureGov (32) or AzureChina (2048). Valid values: true or false. Default value: false. Available as of provider version R2.18+.
* `bgp_lan_interfaces` - (Optional) Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP Transit. Each interface has the following attributes:
  * `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider.
  * `subnet` - (Required) A VPC Network address range selected from one of the available network ranges.
* `ha_bgp_lan_interfaces` - (Optional) Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP Transit HA. Each interface has the following attributes:
  * `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider.
  * `subnet` - (Required) A VPC Network address range selected from one of the available network ranges.
* `bgp_lan_interfaces_count` - (Optional) Number of interfaces that will be created for BGP over LAN enabled Azure transit. Valid value: 1~5 for FireNet case, 1~7 for Non-FireNet case. Default value: 1. Available as of provider version R2.22+.

### Encryption
* `enable_encrypt_volume` - (Optional) Enable EBS volume encryption for Gateway. Only supports AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret. Valid values: true, false. Default value: false.
* `customer_managed_keys` - (Optional and Sensitive) Customer managed key ID.

### Route Customization
* `customized_spoke_vpc_routes` - (Optional) A list of comma-separated CIDRs to be customized for the spoke VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. It applies to all spoke gateways attached to this transit gateway. Example: "10.0.0.0/16,10.2.0.0/16".
* `filtered_spoke_vpc_routes` - (Optional) A list of comma-separated CIDRs to be filtered from the spoke VPC route table. When configured, filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s routing table. It applies to all spoke gateways attached to this transit gateway. Example: "10.2.0.0/16,10.3.0.0/16".
* `excluded_advertised_spoke_routes` - (Optional) A list of comma-separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'. When configured, it inspects all the advertised CIDRs from its spoke gateways and remove those included in the 'Excluded CIDR List'. Example: "10.4.0.0/16,10.5.0.0/16".
* `customized_transit_vpc_routes` - (Optional) A list of CIDRs to be customized for the transit VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. To be effective, `enable_advertise_transit_cidr` or firewall management access for a Transit FireNet gateway must be enabled. Example: ["10.0.0.0/16", "10.2.0.0/16"].

### [Learned CIDRs Approval](https://docs.aviatrix.com/HowTos/transit_approval.html)

-> **NOTE:** `enable_learned_cidrs_approval` can be set to true only if `learned_cidrs_approval_mode` is set to 'gateway'. If `learned_cidrs_approval_mode` is set to 'connection' then enabling learned CIDRs approval is handled within each individual connection resource.

* `enable_learned_cidrs_approval` - (Optional) Switch to enable/disable encrypted transit approval for transit gateway. Valid values: true, false. Default value: false.
* `learned_cidrs_approval_mode` - (Optional) Learned CIDRs approval mode. Either "gateway" (approval on a per gateway basis) or "connection" (approval on a per connection basis). Default value: "gateway". Available as of provider version R2.18+.
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
* `allocate_new_eip` - (Optional) When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in Controller 4.7+. Valid values: true, false. Default: true.
* `eip` - (Optional) Required when `allocate_new_eip` is false. It uses the specified EIP for this gateway. Available in Controller version 4.7+. Only available for AWS, GCP, Azure, OCI, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret.
* `azure_eip_name_resource_group` - (Optional) Name of public IP Address resource and its resource group in Azure to be assigned to the Transit Gateway instance. Example: "IP_Name:Resource_Group_Name". Required if `allocate_new_eip` is false and `cloud_type` is Azure, AzureGov or AzureChina. Available as of provider version 2.20+.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server for Gateway. Currently only supported for AWS, Azure, AzureGov, AWSGov, AWSChina, AzureChina, Alibaba Cloud, AWS Top Secret and AWS Secret gateways. Valid values: true, false. Default value: false.
* `zone` - (Optional) Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'. Available in provider version R2.17+.
* `enable_jumbo_frame` - (Optional) Enable jumbo frames for this transit gateway. Default value is true.
* `tags` - (Optional) Map of tags to assign to the gateway. Only available for AWS, Azure, AzureGov, AWSGov, AWSChina, AzureChina, AWS Top Secret and AWS Secret gateways. Allowed characters vary by cloud type but always include: letters, spaces, and numbers. AWS, AWSGov, AWSChina, AWS Top Secret and AWS Secret allow the use of any character.  Azure, AzureGov and AzureChina allows the following special characters: + - = . _ : @. Example: {"key1" = "value1", "key2" = "value2"}.
* `tunnel_detection_time` - (Optional) The IPsec tunnel down detection time for the Transit Gateway in seconds. Must be a number in the range [20-600]. The default value is set by the controller (60 seconds if nothing has been changed). **NOTE: The controller UI has an option to set the tunnel detection time for all gateways. To achieve the same functionality in Terraform, use the same TF_VAR to manage the tunnel detection time for all gateways.** Available in provider R2.19+.
* `rx_queue_size` - (Optional) Gateway ethernet interface RX queue size. Once set, can't be deleted or disabled. Available for AWS as of provider version R2.22+.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `eip` - Public IP address assigned to the gateway.
* `ha_eip` - Public IP address assigned to the HA gateway.
* `security_group_id` - Security group used for the transit gateway.
* `ha_security_group_id` - HA security group used for the transit gateway.
* `cloud_instance_id` - Cloud instance ID of the transit gateway.
* `private_ip` - Private IP address of the transit gateway created.
* `ha_cloud_instance_id` - Cloud instance ID of the HA transit gateway.
* `ha_gw_name` - Aviatrix transit gateway unique name of HA transit gateway.
* `ha_private_ip` - Private IP address of the HA transit gateway created.
* `lan_interface_cidr` - LAN interface CIDR of the transit gateway created (will be used when enabling FQDN Firenet in Azure). Available in provider version R2.17.1+.
* `ha_lan_interface_cidr` - LAN interface CIDR of the HA transit gateway created (will be used when enabling FQDN Firenet in Azure). Available in provider version R2.18+.
* `bgp_lan_ip_list` - List of available BGP LAN interface IPs for transit external device connection creation. Only supports GCP. Available as of provider version R2.21.0+.
* `ha_bgp_lan_ip_list` - List of available BGP LAN interface IPs for transit external device HA connection creation. Only supports GCP. Available as of provider version R2.21.0+.

The following arguments are deprecated:

* `enable_firenet_interfaces` - (Optional) Sign of readiness for FireNet connection. Valid values: true, false. Default value: false.
* `enable_snat` - (Optional) Enable Source NAT for this container. Valid values: true, false.
* `tag_list` - (Optional) Instance tag of cloud provider. Only supported for AWS, Azure, AzureGov, AWSGov, AWSChina, AzureChina. Example: ["key1:value1","key2:value2"].
* `enable_active_mesh` - (Optional) Switch to enable/disable [Active Mesh Mode](https://docs.aviatrix.com/HowTos/activemesh_faq.html) for Transit Gateway. Valid values: true, false. Default value: false.
* `storage_name` (Optional) Specify a storage account. Required if `cloud_type` is 2048 (AzureChina). Removed in Provider version 2.21.0+.

## Import

**transit_gateway** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_transit_gateway.test gw_name
```

## Notes
### CIDR advertising
`enable_advertise_transit_cidr` and `bgp_manual_spoke_advertise_cidrs` functionality has been migrated over to **aviatrix_transit_gateway** as of Aviatrix Terraform Provider R2.6. If you are using/upgraded to Aviatrix Terraform Provider R2.6+, and a **vgw_conn** resource was originally created with a provider version <R2.6, you must cut and paste these two arguments (and values) into the corresponding transit gateway resource referenced in the **vgw_conn**. A 'terraform refresh' will then successfully complete the migration and rectify the state file.

### enable_firenet
If you are using/upgraded to Aviatrix Terraform Provider R2.5+/UserConnect-5.0+ , and an AWS transit_gateway resource with `enable_firenet_interfaces` enabled was created with a provider version < R2.5/ UserConnect-5.0, you must replace `enable_firenet_interfaces` with `enable_firenet` in your configuration file, and do ‘terraform refresh’ to set its value to `enable_firenet` and apply it into the state file.

### insane_mode
If `insane_mode` is enabled, you must specify a valid /26 CIDR segment of the VPC specified for the `subnet`. This will then create a new subnet to be used for the corresponding gateway. You cannot specify an existing /26 subnet.

### enable_snat
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a transit gateway with `enable_snat` set to true was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also change this attribute to `single_ip_snat` in your `.tf` file.

### ha_subnet
If you are using Aviatrix Terraform Provider R2.15+, and import a Google Cloud transit gateway with HA enabled then you must set a value for `ha_subnet` in your Terraform config.
