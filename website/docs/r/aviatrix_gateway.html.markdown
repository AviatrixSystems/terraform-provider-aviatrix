---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
description: |-
  Creates and manages Aviatrix Gateways
---

# aviatrix_gateway

The aviatrix_gateway resource allows the creation and management of Aviatrix gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Gateway
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtxgw1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-west-1"
  gw_size      = "t2.micro"
  subnet       = "10.0.0.0/24"
  tag_list     = [
    "k1:v1",
    "k2:v2",
  ]
}
```
```hcl
# Create an Aviatrix AWS Gateway with VPN enabled
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtxgw1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-west-1"
  gw_size      = "t2.micro"
  subnet       = "10.0.0.0/24"
  vpn_access   = "yes"
  vpn_cidr     = "192.168.43.0/24"
  max_vpn_conn = "100"
}
```
```hcl
# Create an Aviatrix GCP Gateway
resource "aviatrix_gateway" "test_gateway_gcp" {
  cloud_type   = 4
  account_name = "devops-gcp"
  gw_name      = "avtxgw-gcp"
  vpc_id       = "gcp-gw-vpc"
  vpc_reg      = "us-west1-b"
  gw_size      = "n1-standard-1"
  subnet       = "10.12.0.0/24"
}
```
```hcl
# Create an Aviatrix ARM Gateway
resource "aviatrix_gateway" "test_gateway_arm" {
  cloud_type   = 8
  account_name = "devops-arm"
  gw_name      = "avtxgw-arm"
  vpc_id       = "gateway:test-gw-123"
  vpc_reg      = "West US"
  gw_size      = "Standard_D2"
  subnet       = "10.13.0.0/24"
}
```
```hcl
# Create an Aviatrix Oracle Gateway
resource "aviatrix_gateway" "test_gateway_oracle" {
  cloud_type   = 16
  account_name = "devops-oracle"
  gw_name      = "avtxgw-oracle"
  vpc_id       = "vpc-oracle-test"
  vpc_reg      = "us-ashburn-1"
  gw_size      = "VM.Standard2.2"
  subnet       = "10.7.0.0/16"
}
```
```hcl
# Create an Aviatrix AWSGov Gateway
resource "aviatrix_gateway" "test_gateway_awsgov" {
  cloud_type   = 256
  account_name = "devops-awsgov"
  gw_name      = "avtxgw-awsgov"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-gov-west-1"
  gw_size      = "t2.micro"
  subnet       = "10.0.0.0/24"
  tag_list     = [
    "k1:v1",
    "k2:v2",
  ]
}
```
```hcl
# Create an Aviatrix AWS Gateway with Peering HA enabled
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type         = 1
  account_name       = "devops"
  gw_name            = "avtxgw1"
  vpc_id             = "vpc-abcdef"
  vpc_reg            = "us-west-1"
  gw_size            = "t2.micro"
  subnet             = "10.0.0.0/24"
  peering_ha_subnet  = "10.0.0.0/24"
  peering_ha_gw_size = "t2.micro"
}
```
```hcl
# Create an Aviatrix GCP Gateway with Peering HA enabled
resource "aviatrix_gateway" "test_gateway_gcp" {
  cloud_type         = 4
  account_name       = "devops-gcp"
  gw_name            = "avtxgw-gcp"
  vpc_id             = "gcp-gw-vpc"
  vpc_reg            = "us-west1-b"
  gw_size            = "n1-standard-1"
  subnet             = "10.12.0.0/24"
  peering_ha_zone    = "us-west1-c"
  peering_ha_gw_size = "n1-standard-1"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently only AWS(1), GCP(4), ARM(8), OCI(16), and AWSGov(256) are supported.
* `account_name` - (Required) Account name. This account will be used to launch Aviatrix gateway.
* `gw_name` - (Required) Aviatrix gateway unique name.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Example: AWS: "vpc-abcd1234", GCP: "vpc-gcp-test", ARM: "vnet1:hello", OCI: "vpc-oracle-test1".
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", GCP: "us-west2-a", ARM: "East US 2", Oracle: "us-ashburn-1".
* `gw_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", ARM: "Standard_B1s", Oracle: "VM.Standard2.2", GCP: "n1-standard-1".
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane-mode).**
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Gateway. Required for AWS and AWSGov if insane_mode is set. Example: AWS: "us-west-1a".
* `enable_snat` - (Optional) Enable Source NAT for this container. Valid values: true, false. Default value: false. **NOTE: If using SNAT for FQDN use-case, please see notes [here](#fqdn).**
* `vpn_access` - (Optional) Enable user access through VPN to this container. Valid values: true, false.
* `vpn_cidr` - (Optional) VPN CIDR block for the container. Required if "vpn_access" is true. Example: "192.168.43.0/24".
* `max_vpn_conn` - (Optional) Maximum number of active VPN users allowed to be connected to this gateway. Required if vpn_access is true. Make sure the number is smaller than the VPN CIDR block. Example: 100. **NOTE: Please see notes [here](#max_vpn_conn) in regards to any deltas found in your state with the addition of this argument in R1.14.**
* `enable_elb` - (Optional) Specify whether to enable ELB or not. Not supported for Oracle gateways. Valid values: true, false.
* `elb_name` - (Optional) A name for the ELB that is created. If it is not specified, a name is generated automatically.
* `split_tunnel` - (Optional) Specify split tunnel mode. Valid values: true, false.
* `name_servers` - (Optional) A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `search_domains` - (Optional) A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `additional_cidrs` - (Optional) A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `otp_mode` - (Optional) Two step authentication mode. "2": DUO, "3": Okta.
* `saml_enabled` - (Optional) This field indicates whether enabling SAML or not. This field is available in controller version 3.3 or later release. Valid values: true, false.
* `enable_vpn_nat` - (Optional) This field indicates whether enabling VPN NAT or not. Only supported for VPN gateway. Valid values: true, false. Default value: true.
* `okta_token` - (Optional) Token for Okta auth mode. Required if otp_mode is "3".
* `okta_url` - (Optional) URL for Okta auth mode. Required if otp_mode is "3".
* `okta_username_suffix` - (Optional) Username suffix for Okta auth mode. Example: "aviatrix.com".
* `duo_integration_key` - (Optional) Integration key for DUO auth mode. Required if otp_mode is "2".
* `duo_secret_key` - (Optional) Secret key for DUO auth mode. Required if otp_mode is "2".
* `duo_api_hostname` - (Optional) API hostname for DUO auth mode. Required: Yes if otp_mode is "2".
* `duo_push_mode` - (Optional) Push mode for DUO auth. Required if otp_mode is "2". Valid values: "auto", "selective" and "token".
* `enable_ldap` - (Optional) Specify whether to enable LDAP or not. Valid values: true, false.
* `ldap_server` - (Optional) LDAP server address. Required if enable_ldap is true.
* `ldap_bind_dn` - (Optional) LDAP bind DN. Required if enable_ldap is true.
* `ldap_password` - (Optional) LDAP password. Required if enable_ldap is true.
* `ldap_base_dn` - (Optional) LDAP base DN. Required if enable_ldap is true.
* `ldap_username_attribute` - (Optional) LDAP user attribute. Required if enable_ldap is true.
* `peering_ha_subnet` - (Optional) Public subnet CIDR to create Peering HA Gateway in. Required for AWS/ARM if enabling Peering HA. Example: AWS: "10.0.0.0/16".
* `peering_ha_zone` - (Optional) Zone information for creating Peering HA Gateway, only zone is accepted. Required for GCP if enabling Peering HA. Example: GCP: "us-west1-c".
* `peering_ha_insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Peering HA Gateway. Required for AWS if insane_mode is set and peering_ha_subnet is set. Example: AWS: "us-west-1a".
* `peering_ha_eip` - (Optional) Public IP address that you want assigned to the HA peering instance. Only available for AWS.
* `peering_ha_gw_size` - (Optional) Size of the Peering HA Gateway to be created. **NOTE: Please see notes [here](#peering_ha_gw_size) in regards to any deltas found in your state with the addition of this argument in R1.8.**
* `single_az_ha` (Optional) When value is true, Controller monitors the health of the gateway and restarts the gateway if it becomes unreachable. Valid values: true, false.
* `allocate_new_eip` - (Optional) When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in controller 2.7 or later release. Valid values: true, false. Default: true. Option not available for GCP, ARM and Oracle gateways, they will automatically allocate new eip's.
* `eip` - (Optional) Required when allocate_new_eip is false. It uses specified EIP for this gateway. Available in controller 3.5 or later release. Only available for AWS.
* `tag_list` - (Optional) Instance tag of cloud provider. Only available for AWS and AWSGov. Example: ["key1:value1", "key2:value2"].
* `insane_mode` - (Optional) Enable Insane Mode for Gateway. Insane Mode Gateway size must be at least c5 (AWS) or Standard_D3_v2 (ARM). If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet. Only supported for AWS, AWSGov or ARM. Valid values: true, false.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server for Gateway. Currently only supports AWS and AWSGov. Valid values: true, false. Default value: false.
* `enable_designated_gateway` - (Optional) Enable 'designated_gateway' feature for Gateway. Only supports AWS. Valid values: true, false. Default value: false.
* `additional_cidrs_designated_gateway` - (Optional) A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `public_ip` - Public IP address of the Gateway created.
* `backup_public_ip` - Private IP address of the Gateway created.
* `public_dns_server` - DNS server used by the gateway. Default is "8.8.8.8", can be overridden with the VPC's setting.
* `security_group_id` - Security group used for the gateway.
* `cloud_instance_id` - Instance ID of the gateway.
* `cloudn_bkup_gateway_inst_id` - Instance ID of the backup gateway.

The following arguments are deprecated:

* `dns_server` - Specify the DNS IP, only required while using a custom private DNS for the VPC.


## Import

Instance gateway can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_gateway.test gw_name
```


## Notes
### FQDN
In order for the FQDN feature to be enabled for the specified gateway, `enable_snat` must be set to true. If it is not set at gateway creation, creation of FQDN resource will automatically enable SNAT and users must rectify the diff in the Terraform state by setting `enable_snat = true` in their config file.

### Insane Mode
If `insane_mode` is enabled, you must specify a valid /26 CIDR segment of the VPC specified for the `subnet`. This will then create a new subnet to be used for the corresponding gateway. You cannot specify an existing /26 subnet.

### max_vpn_conn
If you are using/upgraded to Aviatrix Terraform Provider R1.14+, and a gateway with VPN enabled was originally created with a provider version <R1.14, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also input this attribute and its value to "100" in your `.tf` file.

### peering_ha_gw_size
If you are using/upgraded to Aviatrix Terraform Provider R1.8+, and a peering-HA gateway was originally created with a provider version <R1.8, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also input this attribute and its value to its corresponding gateway resource in your `.tf` file.
