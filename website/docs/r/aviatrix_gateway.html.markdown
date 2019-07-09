---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
sidebar_current: "docs-aviatrix-resource-gateway"
description: |-
  Creates and manages an Aviatrix gateway.
---

# aviatrix_gateway

The Account resource allows the creation and management of an Aviatrix gateway.

## Example Usage

```hcl
# Create Aviatrix AWS gateway
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtxgw1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-west-1"
  vpc_size     = "t2.micro"
  vpc_net      = "10.0.0.0/24"
  tag_list     = [
    "k1:v1",
    "k2:v2",
  ]
}

# Create Aviatrix AWS gateway with VPN enabled
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtxgw1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-west-1"
  gw_size      = "t2.micro"
  subnet       = "10.0.0.0/24"
  vpn_acess    = "yes"
  vpn_cidr     = "192.168.43.0/24"
  max_vpn_conn = "100"
}

# Create Aviatrix GCP gateway
resource "aviatrix_gateway" "test_gateway_gcp" {
  cloud_type   = 4
  account_name = "devops-gcp"
  gw_name      = "avtxgw-gcp"
  vpc_id       = "gcp-gw-vpc"
  vpc_reg      = "us-west1-b"
  gw_size      = "f1-micro"
  subnet       = "10.12.0.0/24"
}

# Create Aviatrix ARM gateway
resource "aviatrix_gateway" "test_gateway_arm" {
  cloud_type   = 8
  account_name = "devops-arm"
  gw_name      = "avtxgw-arm"
  vpc_id       = "gateway:test-gw-123"
  vpc_reg      = "West US"
  gw_size      = "Standard_D2"
  subnet       = "10.13.0.0/24"
}

# Create Aviatrix AWS gateway with Peering HA enabled
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type        = 1
  account_name      = "devops"
  gw_name           = "avtxgw1"
  vpc_id            = "vpc-abcdef"
  vpc_reg           = "us-west-1"
  gw_size           = "t2.micro"
  subnet            = "10.0.0.0/24"
  peering_ha_subnet = "10.0.0.0/24"
}
# Create Aviatrix GCP gateway with Peering HA enabled
resource "aviatrix_gateway" "test_gateway_gcp" {
  cloud_type      = 4
  account_name    = "devops-gcp"
  gw_name         = "avtxgw-gcp"
  vpc_id          = "gcp-gw-vpc"
  vpc_reg         = "us-west1-b"
  gw_size         = "f1-micro"
  subnet          = "10.12.0.0/24"
  peering_ha_zone = "us-west1-c"
}

```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. (Only AWS is supported currently. Enter 1 for AWS.)
* `account_name` - (Required) Account name. This account will be used to launch Aviatrix gateway.
* `gw_name` - (Required) Aviatrix gateway unique name.
* `vpc_id` - (Required) ID of legacy VPC/Vnet to be connected. A string that is consisted of VPC/Vnet name and cloud provider's resource name. Please check the "Gateway" page on Aviatrix controller GUI for the precise value if needed. (Example:  "vpc-abcd1234" )
* `vpc_reg` - (Required) Region where this gateway will be launched. (Example: us-east-1) If creating GCP gateway, enter a valid zone for vpc_reg. (Example: us-west1-c)
* `gw_size` - (Required) Size of Gateway Instance. e.g.: "t2.micro"
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. ( Example: "172.31.0.0/20")
* `enable_nat` - (Optional) Enable NAT for this container. (Supported values: true, false)
* `vpn_access` - (Optional) Enable user access through VPN to this container. (Supported values: true, false)
* `vpn_cidr` - (Optional) VPN CIDR block for the container. (Required if vpn_access is "yes", Example: "192.168.43.0/24")
* `max_vpn_conn` - (Optional) Maximum number of active VPN users allowed to be connected to this gateway. (Required if vpn_access is "yes". Make sure the number is smaller than the VPN CIDR block, e.g. 100)
* `enable_elb` - (Optional) Specify whether to enable ELB or not. (Required: Yes when cloud_type is "1", "4", "256" or "1024", supported values: true, false)
* `elb_name` - (Optional) A name for the ELB that is created. If it is not specified a name is generated automatically
* `split_tunnel` - (Optional) Specify split tunnel mode. (Supported values: true, false)
* `name_servers` - (Optional) A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `search_domains` - (Optional) A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `additional_cidrs` - (Optional) A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `otp_mode` - (Optional) Two step authentication mode. "2": DUO, "3": Okta.
* `saml_enabled` - (Optional) This field indicates whether enabling SAML or not. (This field is available in version 3.3 or later release.) (Supported values: true, false)
* `okta_token` - (Optional) Token for Okta auth mode. (Required: Yes if otp_mode is "3")
* `okta_url` - (Optional) URL for Okta auth mode. (Required: Yes if otp_mode is "3")
* `okta_username_suffix` - (Optional) Username suffix for Okta auth mode. (Example: "aviatrix.com")
* `duo_integration_key` - (Optional) Integration key for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_secret_key` - (Optional) Secret key for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_api_hostname` - (Optional) API hostname for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_push_mode` - (Optional) Push mode for DUO auth. Valid values: "auto", "selective" and "token". (Required: Yes if otp_mode is "2")
* `enable_ldap` - (Optional) Specify whether to enable LDAP or not. (Supported values: true, false)
* `ldap_server` - (Optional) LDAP server address. (Required: Yes if enable_ldap is "yes")
* `ldap_bind_dn` - (Optional) LDAP bind DN. (Required: Yes if enable_ldap is "yes)
* `ldap_password` - (Optional) LDAP password. (Required: Yes if enable_ldap is "yes")
* `ldap_base_dn` - (Optional) LDAP base DN. (Required: Yes if enable_ldap is "yes")
* `ldap_username_attribute` - (Optional) LDAP user attribute. (Required: Yes if enable_ldap is "yes")
* `peering_ha_subnet` - (Optional) Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Required for AWS/ARM if enabling Peering HA. Example: AWS: "10.0.0.0/16".
* `peering_ha_zone` - (Optional) Zone information for creating Peering HA Gateway, only zone is accepted. Required for GCP if enabling Peering HA. (Example: GCP: "us-west1-c")
* `peering_ha_eip` - (Optional) Public IP address that you want assigned to the HA peering instance. Only available for AWS.
* `peering_ha_gw_size` - (Optional) Size of the Peering HA Gateway.
* `single_az_ha` (Optional) Set to true if this feature is desired. (Supported values: true, false)
* `allocate_new_eip` - (Optional) When value is off, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in 2.7 or later release. (Supported values : true, false) (Default: true) Option not available for GCP and ARM gateways, they will automatically allocate new eip's.
* `eip` - (Optional) Required when allocate_new_eip is "off". It uses specified EIP for this gateway. Available in 3.5 or later release eip. Only available for AWS.
* `tag_list` - (Optional) Instance tag of cloud provider. Example: ["key1:value1", "key002:value002"] Only available for AWS.

The following arguments are computed - please do not edit in the resource file:

* `public_ip` - Public IP address of the Gateway created.
* `backup_public_ip` - Private IP address of the Gateway created.
* `public_dns_server` - DNS server used by the gateway. Default is "8.8.8.8", can be overridden with the VPC's setting.
* `security_group_id` - Security group used for the gateway.
* `cloud_instance_id` - Instance ID of the gateway.
* `cloudn_bkup_gateway_inst_id` - Instance ID of the backup gateway.

The following arguments are deprecated:

* `dns_server` - Specify the DNS IP, only required while using a custom private DNS for the VPC.

-> **NOTE:** 

* `peering_ha_gw_size` - If you are using/upgraded to Aviatrix Terraform Provider v4.3+, and a peering-HA gateway was originally created with a provider version <4.3, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also input this attribute and its value to its corresponding gateway resource in your `.tf` file. 
* `enable_nat` - In order for the FQDN feature to be enabled for the specified gateway, "enable_nat" must be set to “yes”. If it is not set at gateway creation, creation of FQDN resource will automatically enable SNAT and users must rectify the diff in the Terraform state by setting "enable_nat = 'yes'" in their config file.
* `max_vpn_conn` - If you are using/upgraded to Aviatrix Terraform Provider v4.7+, and a gateway with VPN enabled was originally created with a provider version <4.7, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also input this attribute and its value to "100" in your `.tf` file.

## Import

Instance gateway can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_gateway.test gw_name
```
