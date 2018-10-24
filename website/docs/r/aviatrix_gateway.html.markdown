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
  cloud_type = 1
  account_name = "devops"
  gw_name = "avtxgw1"
  vpc_id = "vpc-abcdef"
  vpc_reg = "us-west-1"
  vpc_size = "t2.micro"
  vpc_net = "10.0.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. (Only AWS is supported currently. Enter 1 for AWS.)
* `account_name` - (Required) Account name. This account will be used to launch Aviatrix gateway.
* `gw_name` - (Required) Aviatrix gateway unique name.
* `vpc_id` - (Required) ID of legacy VPC/Vnet to be connected. A string that is consisted of VPC/Vnet name and cloud provider's resource name. Please check the "Gateway" page on Aviatrix controller GUI for the precise value if needed
* `vpc_reg` - (Required) Region where this gateway will be launched.
* `vpc_size` - (Required) Size of Gateway Instance. e.g.: "t2.micro"
* `vpc_net` - (Required) A VPC Network address range selected from one of the available network ranges.
* `enable_nat` - (Optional) Enable NAT for this container.
* `dns_server` - (Optional) Specify a public DNS for the gateway
* `vpn_access` - (Optional) Enable user access through VPN to this container.
* `cidr` - (Optional) VPN CIDR block for the container (Required: Yes if vpn_access is "yes")
* `enable_elb` - (Optional) Specify whether to enable ELB or not. (Required: Yes when cloud_type is "1", "4", "256" or "1024")
* `split_tunnel` - (Optional) Specify split tunnel mode.
* `otp_mode` - (Optional) Two step authentication mode. "2": DUO, "3": Okta. 
* `saml_enabled` - (Optional) This field indicates whether enabling SAML or not. (This field is available in version 3.3 or later release.)
* `okta_token` - (Optional) Token for Okta auth mode (Required: Yes if otp_mode is "3")
* `okta_url` - (Optional) URL for Okta auth mode. (Required: Yes if otp_mode is "3")
* `okta_username_suffix` - (Optional) Username suffix for Okta auth mode
* `duo_integration_key` - (Optional) Integration key for DUO auth mode (Required: Yes if otp_mode is "2")
* `duo_secret_key` - (Optional) Secret key for DUO auth mode (Required: Yes if otp_mode is "2")
* `duo_api_hostname` - (Optional) API hostname for DUO auth mode. (Required: Yes if otp_mode is "2")
* `duo_push_mode` - (Optional) Push mode for DUO auth. Valid values: "auto", "selective" and "token" (Required: Yes if otp_mode is "2")
* `enable_ldap` - (Optional) Specify whether to enable LDAP or not.
* `ldap_server` - (Optional) LDAP server address (Required: Yes if enable_ldap is "yes")
* `ldap_bind_dn` - (Optional) LDAP bind DN (Required: Yes if enable_ldap is "yes)
* `ldap_password` - (Optional) LDAP password (Required: Yes if enable_ldap is "yes")
* `ldap_base_dn` - (Optional) LDAP base DN (Required: Yes if enable_ldap is "yes")
* `ldap_username_attribute` - (Optional) LDAP user attribute (Required: Yes if enable_ldap is "yes")
* `ha_subnet` - (Optional) This is for Gateway HA. Deprecated. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
* `public_subnet` - Public Subnet Information while creating Peering HA Gateway. Example: AWS: "10.0.0.0/16\~\~ZONE\~\~SubnetName"
* `zone` - (Optional) A GCE zone where this gateway will be launched. (Required when cloud_type is 4)
* `single_az_ha` (Optional) Set to "enabled" if this feature is desired
* `allocate_new_eip` - (Optional) When value is off, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in 2.7 or later release
* `eip` - (Optional) When allocate_new_eip is off, use specified IP for this gateway. Available in 3.5 or later release
eip

The following arguments are computed - please do not edit in the resource file:

* `public_ip` - Public IP address of the Gateway created
* `backup_public_ip` - Private IP address of the Gateway created
* `public_dns_server` - DNS server used by the gateway. Default is "8.8.8.8", can be overridden with the VPC's setting
* `security_group_id` - Security group used for the gateway.
* `cloud_instance_id` - Instance ID of the gateway
* `cloudn_bkup_gateway_inst_id` - Instance ID of the backup gateway
