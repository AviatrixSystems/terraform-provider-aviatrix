---
subcategory: "Gateway"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
description: |-
  Creates and manages Aviatrix gateways
---

# aviatrix_gateway

The **aviatrix_gateway** resource allows the creation and management of Aviatrix gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Gateway
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtx-gw-1"
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
# Create an Aviatrix AWS Gateway with Peering HA enabled
resource "aviatrix_gateway" "test_gateway_aws" {
  cloud_type         = 1
  account_name       = "devops"
  gw_name            = "avtx-gw-1"
  vpc_id             = "vpc-abcdef"
  vpc_reg            = "us-west-1"
  gw_size            = "t2.micro"
  subnet             = "10.0.0.0/24"
  peering_ha_subnet  = "10.0.0.0/24"
  peering_ha_gw_size = "t2.micro"
}
```
```hcl
# Create an Aviatrix AWS Gateway with VPN enabled
resource "aviatrix_gateway" "test_vpn_gateway_aws" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtx-gw-1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-west-1"
  gw_size      = "t2.micro"
  subnet       = "10.0.0.0/24"
  vpn_access   = true
  vpn_cidr     = "192.168.43.0/24"
  max_vpn_conn = "100"
}
```
```hcl
# Create an Aviatrix AWS Public Subnet Filtering Gateway
resource "aviatrix_gateway" "test_psf_gateway_aws" {
  cloud_type                                  = 1
  account_name                                = "devops"
  gw_name                                     = "avtx-psf-gw-1"
  vpc_id                                      = "vpc-abcdef"
  vpc_reg                                     = "us-west-1"
  gw_size                                     = "t2.micro"
  subnet                                      = "10.0.0.0/24"
  zone                                        = "us-west-1b"
  enable_public_subnet_filtering              = true
  public_subnet_filtering_route_tables        = [data.aviatrix_vpc.test_vpc.route_tables[0]]
  public_subnet_filtering_guard_duty_enforced = true
  single_az_ha                                = true
  enable_encrypt_volume                       = true
}
```
```hcl
# Create an Aviatrix AWS Public Subnet Filtering Gateway with HA enabled
resource "aviatrix_gateway" "test_psf_gateway_aws" {
  cloud_type                                  = 1
  account_name                                = "devops"
  gw_name                                     = "avtx-psf-gw-1"
  vpc_id                                      = "vpc-abcdef"
  vpc_reg                                     = "us-west-1"
  gw_size                                     = "t2.micro"
  subnet                                      = "10.0.0.0/24"
  zone                                        = "us-west-1b"
  enable_public_subnet_filtering              = true
  public_subnet_filtering_route_tables        = [data.aviatrix_vpc.test_vpc.route_tables[0]]
  peering_ha_subnet                           = "10.10.0.64/26"
  peering_ha_zone                             = "us-west-1b"
  public_subnet_filtering_ha_route_tables     = [data.aviatrix_vpc.test_vpc.route_tables[1]]
  public_subnet_filtering_guard_duty_enforced = true
  single_az_ha                                = true
  enable_encrypt_volume                       = true
}
```
```hcl
# Create an Aviatrix GCP Gateway
resource "aviatrix_gateway" "test_gateway_gcp" {
  cloud_type   = 4
  account_name = "devops-gcp"
  gw_name      = "avtx-gw-gcp"
  vpc_id       = "gcp-gw-vpc"
  vpc_reg      = "us-west1-b"
  gw_size      = "n1-standard-1"
  subnet       = "10.12.0.0/24"
}
```
```hcl
# Create an Aviatrix GCP Gateway with Peering HA enabled
resource "aviatrix_gateway" "test_gateway_gcp" {
  cloud_type         = 4
  account_name       = "devops-gcp"
  gw_name            = "avtx-gw-gcp"
  vpc_id             = "gcp-gw-vpc"
  vpc_reg            = "us-west1-b"
  gw_size            = "n1-standard-1"
  subnet             = "10.12.0.0/24"
  peering_ha_zone    = "us-west1-c"
  peering_ha_subnet  = "10.12.0.0/24"
  peering_ha_gw_size = "n1-standard-1"
}
```
```hcl
# Create an Aviatrix Azure Gateway
resource "aviatrix_gateway" "test_gateway_azure" {
  cloud_type   = 8
  account_name = "devops-azure"
  gw_name      = "avtx-gw-azure"
  vpc_id       = "gateway:test-gw-123"
  vpc_reg      = "West US"
  gw_size      = "Standard_D2"
  subnet       = "10.13.0.0/24"
  zone         = "az-2"
}
```
```hcl
# Create an Aviatrix Oracle Gateway
resource "aviatrix_gateway" "test_gateway_oci" {
  cloud_type   = 16
  account_name = "devops-oci"
  gw_name      = "avtx-gw-oci"
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
  gw_name      = "avtx-gw-awsgov"
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

## Argument Reference

The following arguments are supported:

### Required
* `cloud_type` - (Required) Cloud service provider to use to launch the gateway. Requires an integer value. Currently supports AWS(1), GCP(4), AZURE(8), OCI(16), and AWSGov(256).
* `account_name` - (Required) Account name. This account will be used to launch Aviatrix gateway.
* `gw_name` - (Required) Name of the Aviatrix gateway to be created.
* `vpc_id` - (Required) VPC ID/VNet name of cloud provider. Example: AWS: "vpc-abcd1234", GCP: "vpc-gcp-test", AZURE: "vnet1:hello", OCI: "vpc-oracle-test1".
* `vpc_reg` - (Required) VPC region the gateway will be created in. Example: AWS: "us-east-1", GCP: "us-west2-a", AZURE: "East US 2", OCI: "us-ashburn-1".
* `gw_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", GCP: "n1-standard-1", AZURE: "Standard_B1s", OCI: "VM.Standard2.2".
* `subnet` - (Required) A VPC network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane_mode-1).**

### HA
* `single_az_ha` (Optional) If enabled, Controller monitors the health of the gateway and restarts the gateway if it becomes unreachable. Valid values: true, false. Default value: false.
* `peering_ha_subnet` - (Optional) Public subnet CIDR to create Peering HA Gateway in. Required if enabling Peering HA for AWS/AZURE. Optional if enabling Peering HA for GCP. Example: AWS: "10.0.0.0/16".
* `peering_ha_zone` - (Optional) Zone to create Peering HA Gateway in. Required if enabling Peering HA for GCP. Example: GCP: "us-west1-c". Optional for AZURE. Valid values for AZURE gateways are in the form "az-n". Example: "az-2". Available for AZURE as of provider version R2.17+.
* `peering_ha_insane_mode_az` - (Optional) Region + Availability Zone of subnet being created for Insane Mode-enabled Peering HA Gateway. Required for AWS only if `insane_mode` is set and `peering_ha_subnet` is set. Example: AWS: "us-west-1a".
* `peering_ha_eip` - (Optional) Public IP address to be assigned to the HA peering instance. Only available for AWS and GCP.
* `peering_ha_gw_size` - (Optional) Size of the Peering HA Gateway to be created. Required if enabling Peering HA. **NOTE: Please see notes [here](#peering_ha_gw_size-1) in regards to any deltas found in your state with the addition of this argument in R1.8.**

### Insane Mode
* `insane_mode` - (Optional) Enable [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) for Gateway. Insane Mode gateway size must be at least c5 series (AWS) or Standard_D3_v2 (AZURE). If enabled, a valid /26 CIDR segment of the VPC must be specified to create a new subnet. Only supported for AWS, AWSGov or Azure. Valid values: true, false.
* `insane_mode_az` - (Optional) Region + Availability Zone of subnet being created for Insane Mode gateway. Required for AWS and AWSGov if `insane_mode` is set. Example: AWS: "us-west-1a".

### SNAT/DNAT
* `single_ip_snat` - (Optional) Enable Source NAT in "single ip" mode for this gateway. Valid values: true, false. Default value: false. **NOTE: If using SNAT for FQDN use-case, please see notes [here](#fqdn).**

-> **NOTE:** `enable_snat` has been renamed to `single_ip_snat` in provider version R2.10. Please see notes [here](#enable_snat-1) for more information.

~> **NOTE:** Custom DNAT support has been deprecated and functionality has been moved to **aviatrix_gateway_dnat** in provider version R2.10. Please see notes [here](#dnat_policy-1).

### VPN Access
~> **NOTE:** If the ELB/gateway is being managed by a Geo VPN, in order to update VPN configurations of the Geo VPN, all the VPN configurations of the ELBs/gateways must be updated simultaneously and share the same values. This can be achieved by managing the VPN configurations through variables and updating their values accordingly.

* `vpn_access` - (Optional) Enable [user access through VPN](https://docs.aviatrix.com/HowTos/gateway.html#vpn-access) to this gateway. Valid values: true, false.
* `vpn_cidr` - (Optional) VPN CIDR block for the gateway. Required if `vpn_access` is true. Example: "192.168.43.0/24".
* `max_vpn_conn` - (Optional) Maximum number of active VPN users allowed to be connected to this gateway. Required if `vpn_access` is true. Make sure the number is smaller than the VPN CIDR block. Example: 100. **NOTE: Please see notes [here](#max_vpn_conn-1) in regards to any deltas found in your state with the addition of this argument in R1.14.**
* `enable_elb` - (Optional) Specify whether to enable ELB or not. Not supported for OCI gateways. Valid values: true, false.
* `elb_name` - (Optional) A name for the ELB that is created. If it is not specified, a name is generated automatically.
* `vpn_protocol` - (Optional) Transport mode for VPN connection. All `cloud_types` support TCP with ELB, and UDP without ELB. AWS(1) additionally supports UDP with ELB. Valid values: "TCP", "UDP". If not specified, "TCP" will be used.

#### Split Tunnel
* `split_tunnel` - (Optional) Enable/disable Split Tunnel Mode. Valid values: true, false. Default value: true. Please see [here](https://docs.aviatrix.com/HowTos/gateway.html#split-tunnel-mode) for more information on split tunnel.
* `name_servers` - (Optional) A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `search_domains` - (Optional) A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `additional_cidrs` - (Optional) A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.

#### MFA Authentication
* `otp_mode` - (Optional) Two step authentication mode. Valid values: "2" for DUO, "3" for Okta.
* `saml_enabled` - (Optional) Enable/disable SAML. This field is available in Controller version 3.3 or later release. Valid values: true, false. Default value: false.
* `enable_vpn_nat` - (Optional) Enable/disable VPN NAT. Only supported for VPN gateway. Valid values: true, false. Default value: true.
* `okta_token` - (Optional) Token for Okta auth mode. Required if `otp_mode` is "3".
* `okta_url` - (Optional) URL for Okta auth mode. Required if `otp_mode` is "3".
* `okta_username_suffix` - (Optional) Username suffix for Okta auth mode. Example: "aviatrix.com".
* `duo_integration_key` - (Optional) Integration key for DUO auth mode. Required if `otp_mode` is "2".
* `duo_secret_key` - (Optional) Secret key for DUO auth mode. Required if `otp_mode` is "2".
* `duo_api_hostname` - (Optional) API hostname for DUO auth mode. Required: Yes if `otp_mode` is "2".
* `duo_push_mode` - (Optional) Push mode for DUO auth. Required if `otp_mode` is "2". Valid values: "auto", "selective" and "token".
* `enable_ldap` - (Optional) Enable/disable LDAP. Valid values: true, false. Default value: false.
* `ldap_server` - (Optional) LDAP server address. Required if `enable_ldap` is true.
* `ldap_bind_dn` - (Optional) LDAP bind DN. Required if `enable_ldap` is true.
* `ldap_password` - (Optional) LDAP password. Required if `enable_ldap` is true.
* `ldap_base_dn` - (Optional) LDAP base DN. Required if `enable_ldap` is true.
* `ldap_username_attribute` - (Optional) LDAP user attribute. Required if `enable_ldap` is true.

#### Modify VPN Configuration
* `idle_timeout` - (Optional) It sets the value (seconds) of the [idle timeout](https://docs.aviatrix.com/HowTos/openvpn_faq.html#how-do-i-fix-the-aviatrix-vpn-timing-out-too-quickly). This idle timeout feature is enable only if this attribute is set, otherwise it is disabled. The entered value must be an integer number greater than 300.  Available in provider version R2.17.1+.
* `renegotiation_interval` - (Optional) It sets the value (seconds) of the [renegotiation interval](https://docs.aviatrix.com/HowTos/openvpn_faq.html#how-do-i-fix-the-aviatrix-vpn-timing-out-too-quickly). This renegotiation interval feature is enable only if this attribute is set, otherwise it is disabled. The entered value must be an integer number greater than 300. Available in provider version R2.17.1+.

### Designated Gateway
* `enable_designated_gateway` - (Optional) Enable Designated Gateway feature for Gateway. Only supported for AWS and AWSGov gateways. Valid values: true, false. Default value: false. Please view documentation [here](https://docs.aviatrix.com/HowTos/gateway.html#designated-gateway) for more information on this feature.
* `additional_cidrs_designated_gateway` - (Optional) A list of CIDR ranges separated by comma to configure when "Designated Gateway" feature is enabled. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".

### Encryption
* `enable_encrypt_volume` - (Optional) Enable EBS volume encryption for the gateway. Only supported for AWS and AWSGov gateways. Valid values: true, false. Default value: false.
* `customer_managed_keys` - (Optional and Sensitive) Customer-managed key ID.

### Monitor Gateway Subnets
~> **NOTE:** This feature is only available for AWS gateways.

* `enable_monitor_gateway_subnets` - (Optional) If set to true, the [Monitor Gateway Subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet) feature is enabled. Default value is false. Available in provider version R2.17.1+.
* `monitor_exclude_list` - (Optional) A list of monitored instance IDs separated by comma when monitoring feature is enabled. Default value is "". Available in provider version R2.17.1+.

### FQDN Gateway
~> **NOTE:** This attribute is only to be used in Azure FQDN FireNet workflows.

* `fqdn_lan_cidr` - (Optional) If `fqdn_lan_cidr` is set, the FQDN gateway will be created with an additional LAN interface using the provided CIDR. This attribute is required when enabling FQDN gateway FireNet in Azure. Available in provider version R2.17.1+.

### Misc.
* `allocate_new_eip` - (Optional) If set to false, use an available address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in Controller 2.7+. Valid values: true, false. Default: true. Option not available for Azure and OCI gateways, they will automatically allocate new EIPs.
* `eip` - (Optional) Specified EIP to use for gateway creation. Required when `allocate_new_eip` is false.  Available in Controller version 3.5+. Only supported for AWS and GCP gateways.
* `tag_list` - (Optional) Tag list of the gateway instance. Only available for AWS and AWSGov gateways. Example: ["key1:value1", "key2:value2"].
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server for gateway. Currently only supported for AWS and AWSGov gateways. Valid values: true, false. Default value: false.
* `zone` - (Optional) Availability Zone. Only available for Azure and Public Subnet Filtering gateway. Available for Azure as of provider version R2.17+.

### Public Subnet Filtering Gateway

~> **NOTE:** When `enable_public_subnet_filtering` is set to true the following attributes cannot be used and doing so will result in a plan time error: "additional_cidrs", "additional_cidrs_designated_gateway", "allocate_new_eip", "customer_managed_keys", "duo_api_hostname", "duo_integration_key", "duo_push_mode", "duo_secret_key", "eip", "elb_name", "enable_designated_gateway", "enable_elb", "enable_ldap", "enable_monitor_gateway_subnets", "enable_vpc_dns_server", "enable_vpn_nat", "fqdn_lan_cidr", "idle_timeout", "insane_mode", "insane_mode_az", "ldap_base_dn", "ldap_bind_dn", "ldap_password", "ldap_server", "ldap_username_attribute", "max_vpn_conn", "monitor_exclude_list", "name_servers", "okta_token", "okta_url", "okta_username_suffix", "otp_mode", "peering_ha_eip", "peering_ha_insane_mode_az", "renegotiation_interval", "saml_enabled", "search_domains", "single_ip_snat", "split_tunnel", "vpn_access", "vpn_cidr", "vpn_protocol".

* `enable_public_subnet_filtering` - (Optional) Create a [Public Subnet Filtering gateway](https://docs.aviatrix.com/HowTos/public_subnet_filtering_faq.html). Valid values: true or false. Default value: false. Available as of provider version R2.18+.
* `public_subnet_filtering_route_tables` - (Optional) Route tables whose associated public subnets are protected. Only valid when `enable_public_subnet_filtering` attribute is true. Available as of provider version R2.18+.
* `public_subnet_filtering_ha_route_tables` - (Optional) Route tables whose associated public subnets are protected for the HA PSF gateway. Required when `enable_public_subnet_filtering` and `peering_ha_subnet` are set. Available as of provider version R2.18+.
* `public_subnet_filtering_guard_duty_enforced` - (Optional) Whether to enforce Guard Duty IP blocking.  Only valid when `enable_public_subnet_filtering` attribute is true. Valid values: true or false. Default value: true. Available as of provider version R2.18+.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `elb_dns_name` - ELB DNS name.
* `public_dns_server` - DNS server used by the gateway. Default is "8.8.8.8", can be overridden with the VPC's setting.
* `security_group_id` - Security group used for the gateway.
* `cloud_instance_id` - Cloud instance ID of the gateway.
* `private_ip` - Private IP address of the gateway created.
* `peering_ha_cloud_instance_id` - Cloud instance ID of the HA gateway.
* `peering_ha_gw_name` - Aviatrix gateway unique name of HA gateway.
* `peering_ha_private_ip` - Private IP address of HA gateway.
* `fqdn_lan_interface` - The lan interface id of the of FQDN gateway with additional LAN interface. This attribute will be exported when enabling FQDN gateway firenet in Azure. Available in provider version R2.17.1+.

The following arguments are deprecated:

* `dns_server` - Specify the DNS IP, only required while using a custom private DNS for the VPC.
* `enable_snat` - (Optional) Enable Source NAT for this gateway. Valid values: true, false. Default value: false. **NOTE: If using SNAT for FQDN use-case, please see notes [here](#fqdn).**

* `dnat_policy` - (Optional) Policy rule applied for enabling Destination NAT (DNAT), which allows you to change the destination to a virtual address range. Currently only supports AWS(1) and AZURE(8).
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
* `cloudn_bkup_gateway_inst_id` - Instance ID of the backup gateway.
* `public_ip` - Public IP address of the gateway created.
* `peering_ha_public_ip` - Public IP address of the peering HA Gateway created.

## Import

**gateway** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_gateway.test gw_name
```


## Notes
### FQDN
In order for the FQDN feature to be enabled for the specified gateway, `single_ip_snat` must be set to true. If it is not set at gateway creation, creation of FQDN resource will automatically enable SNAT and users must rectify the diff in the Terraform state by setting `single_ip_snat = true` in their config file.

### insane_mode
If `insane_mode` is enabled, you must specify a valid /26 CIDR segment of the VPC specified for the `subnet`. This will then create a new subnet to be used for the corresponding gateway. You **cannot** specify an existing /26 subnet.

### max_vpn_conn
If you are using/upgraded to Aviatrix Terraform Provider R1.14+, and a gateway with VPN enabled was originally created with a provider version <R1.14, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also input this attribute and its value to "100" in your `.tf` file.

### peering_ha_gw_size
If you are using/upgraded to Aviatrix Terraform Provider R1.8+, and a peering-HA gateway was originally created with a provider version <R1.8, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also input this attribute and its value to its corresponding gateway resource in your `.tf` file.

### enable_snat
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a gateway with `enable_snat` set to true was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also change this attribute to `single_ip_snat` in your `.tf` file.

### dnat_policy
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a gateway with `dnat_policy` was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to remove attribute’s value from the state. In addition, you must transfer its corresponding values to the **aviatrix_gateway_dnat** resource in your `.tf` file and perform a 'terraform import' to rectify the state file.

### peering_ha_subnet
If you are using Aviatrix Terraform Provider R2.15+, and import a Google Cloud gateway with HA enabled then you must set a value for `peering_ha_subnet` in your Terraform config.
