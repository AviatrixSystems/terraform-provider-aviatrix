---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
description: |-
  Gets the Aviatrix gateway.
---

# aviatrix_gateway

Use this data source to get the Aviatrix gateway for use in other resources.

## Example Usage

```hcl
# Aviatrix Gateway Data Source
data "aviatrix_gateway" "foo" {
  gw_name = "gatewayname"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Gateway name. This can be used for getting gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gw_name` - Aviatrix gateway name.
* `cloud_type` - Type of cloud service provider.
* `account_name` - Account name.
* `vpc_id` - ID of legacy VPC/Vnet to be connected.
* `vpc_reg"` - Region where gateway is launched.
* `gw_size` - Size of Gateway Instance.
* `subnet` - A VPC Network address range selected from one of the available network ranges.
* `insane_mode_az` - AZ of subnet being created for Insane Mode Gateway. Required if insane_mode is set.
* `single_ip_snat` - Enable Source NAT for this container.
* `vpn_access` - Enable user access through VPN to this container.
* `vpn_cidr` - VPN CIDR block for the container.
* `enable_elb` - Specify whether to enable ELB or not.
* `elb_name` - A name for the ELB that is created.
* `vpn_protocol` - Elb protocol for VPN gateway with elb enabled. Only supports AWS provider. Valid values: 'TCP', 'UDP'. If not specified, 'TCP'' will be used.
* `split_tunnel` - Specify split tunnel mode.
* `max_vpn_conn` - Maximum connection of VPN access.
* `name_servers` - A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `search_domains` - A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `additional_cidrs` - A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `otp_mode` - Two step authentication mode.
* `saml_enabled` - This field indicates whether to enable SAML or not.
* `enable_vpn_nat` - This field indicates whether to enable VPN NAT or not. Only supported for VPN gateway. Valid values: true, false. Default value: true.
* `okta_url` - URL for Okta auth mode.
* `okta_username_suffix` - Username suffix for Okta auth mode.
* `duo_integration_key` - Integration key for DUO auth mode.
* `duo_api_hostname` - API hostname for DUO auth mode.
* `duo_push_mode` - Push mode for DUO auth.
* `enable_ldap` - Specify whether to enable LDAP or not. Supported values: 'yes' and 'no'.
* `ldap_server` - LDAP server address. Required: Yes if enable_ldap is 'yes'.
* `ldap_bind_dn` - LDAP bind DN. Required: Yes if enable_ldap is 'yes'.
* `ldap_base_dn` - LDAP base DN. Required: Yes if enable_ldap is 'yes'.
* `ldap_username_attribute` - LDAP user attribute. Required: Yes if enable_ldap is 'yes'.
* `peering_ha_subnet` - Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Required to create peering ha gateway if cloud_type = 1 or 8 (AWS or AZURE).
* `peering_ha_zone` - Zone information for creating Peering HA Gateway. Required to create peering ha gateway if cloud_type = 4 (gcp).
* `peering_ha_insane_mode_az` - AZ of subnet being created for Insane Mode Peering HA Gateway. Required if insane_mode is set.
* `peering_ha_public_ip` - Public IP address that you want assigned to the HA peering instance.
* `peering_ha_gw_size` - Peering HA Gateway Size.
* `single_az_ha` - Set to true if this feature is desired.
* `allocate_new_eip` - When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway.
* `tag_list` - Instance tag of cloud provider.",
* `insane_mode` - Enable Insane Mode for Gateway. Valid values: true, false.
* `enable_vpc_dns_server` - Enable vpc_dns_server for Gateway. Only supports AWS. Valid values: true, false.
* `enable_designated_gateway` - Enable 'designated_gateway' feature for Gateway. Only supports AWS. Valid values: true, false.
* `additional_cidrs_designated_gateway` - A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled.
* `enable_encrypt_volume` - Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.
* `elb_dns_name` - ELB DNS Name.
* `public_ip` - Public IP address of the Gateway created.
* `security_group_id` - Security group used for the gateway.
* `public_dns_server` - NS server used by the gateway.
* `cloud_instance_id` - Instance ID of the gateway.
* `private_ip` - Private IP address of the Gateway created.
* `peering_ha_cloud_instance_id` - Instance ID of the peering HA gateway.
* `peering_ha_gw_name` - Aviatrix gateway unique name of HA gateway.