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
resource "aviatrix_gateway" "test_gateway1" {
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtxgw1"
  vpc_id       = "vpc-abcdef"
  vpc_reg      = "us-west-1"
  vpc_size     = "t2.micro"
  vpc_net      = "10.0.0.0/24"
  tag_list     = ["k1:v1","k2:v2"]
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. (Only AWS is supported currently. Enter 1 for AWS.)
* `account_name` - (Required) Account name. This account will be used to launch Aviatrix gateway.
* `gw_name` - (Required) Aviatrix gateway unique name.
* `vpc_id` - (Required) ID of legacy VPC/Vnet to be connected. A string that is consisted of VPC/Vnet name and cloud provider's resource name. Please check the "Gateway" page on Aviatrix controller GUI for the precise value if needed. (Example:  "vpc-abcd1234" )
* `vpc_reg` - (Required) Region where this gateway will be launched. (Example: us-east-1)
* `vpc_size` - (Required) Size of Gateway Instance. e.g.: "t2.micro"
* `vpc_net` - (Required) A VPC Network address range selected from one of the available network ranges. ( Example: "172.31.0.0/20")
* `enable_nat` - (Optional) Enable NAT for this container. (Supported values: "yes", "no")
* `vpn_access` - (Optional) Enable user access through VPN to this container. (Supported values: "yes", "no")
* `vpn_cidr` - (Optional) VPN CIDR block for the container. (Required if vpn_access is "yes", Example: "192.168.43.0/24")
* `enable_elb` - (Optional) Specify whether to enable ELB or not. (Required: Yes when cloud_type is "1", "4", "256" or "1024", supported values "yes" and "no")
* `elb_name` - (Optional) A name for the ELB that is created. If it is not specified a name is generated automatically
* `split_tunnel` - (Optional) Specify split tunnel mode. (Supported values: "yes", "no")
* `name_servers` - (Optional) A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `search_domains` - (Optional) A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `additional_cidrs` - (Optional) A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `otp_mode` - (Optional) Two step authentication mode. "2": DUO, "3": Okta.
* `saml_enabled` - (Optional) This field indicates whether enabling SAML or not. (This field is available in version 3.3 or later release.) (Supported values: "yes", "no")
* `okta_token` - (Optional) Token for Okta auth mode. (Required: Yes if otp_mode is "3")
* `okta_url` - (Optional) URL for Okta auth mode. (Required: Yes if otp_mode is "3")
* `okta_username_suffix` - (Optional) Username suffix for Okta auth mode. (Example: "aviatrix.com")
* `duo_integration_key` - (Optional) Integration key for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_secret_key` - (Optional) Secret key for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_api_hostname` - (Optional) API hostname for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_push_mode` - (Optional) Push mode for DUO auth. Valid values: "auto", "selective" and "token". (Required: Yes if otp_mode is "2")
* `enable_ldap` - (Optional) Specify whether to enable LDAP or not. (Supported values: "yes", "no")
* `ldap_server` - (Optional) LDAP server address. (Required: Yes if enable_ldap is "yes")
* `ldap_bind_dn` - (Optional) LDAP bind DN. (Required: Yes if enable_ldap is "yes)
* `ldap_password` - (Optional) LDAP password. (Required: Yes if enable_ldap is "yes")
* `ldap_base_dn` - (Optional) LDAP base DN. (Required: Yes if enable_ldap is "yes")
* `ldap_username_attribute` - (Optional) LDAP user attribute. (Required: Yes if enable_ldap is "yes")
* `ha_subnet` - (Optional) This is for Gateway HA. Deprecated. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
* `peering_ha_subnet` - Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Example: AWS: "10.0.0.0/16"
* `peering_ha_eip` - (Optional) Public IP address that you want assigned to the HA peering instance.
* `zone` - (Optional) A GCE zone where this gateway will be launched. (Required when cloud_type is 4)
* `single_az_ha` (Optional) Set to "enabled" if this feature is desired.
* `allocate_new_eip` - (Optional) When value is off, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in 2.7 or later release. (Supported values : "on", "off") (Default: "on")
* `eip` - (Optional) Required when allocate_new_eip is "off". It uses specified EIP for this gateway. Available in 3.5 or later release eip.
* `tag_list` - (Optional) Instance tag of cloud provider. Example: ["key1:value1", "key002:value002"]

The following arguments are computed - please do not edit in the resource file:

* `public_ip` - Public IP address of the Gateway created.
* `backup_public_ip` - Private IP address of the Gateway created.
* `public_dns_server` - DNS server used by the gateway. Default is "8.8.8.8", can be overridden with the VPC's setting.
* `security_group_id` - Security group used for the gateway.
* `cloud_instance_id` - Instance ID of the gateway.
* `cloudn_bkup_gateway_inst_id` - Instance ID of the backup gateway.

-> **NOTE:** The following arguments are deprecated:

* `dns_server` - Specify the DNS IP, only required while using a custom private DNS for the VPC.

## Import

Instance gateway can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_gateway.test gw_name
```