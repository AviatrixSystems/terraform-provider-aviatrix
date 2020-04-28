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

* `gw_name` - (Required) Gateway name. It can be used for getting gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `account_name` - Account name.
* `additional_cidrs` - A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `additional_cidrs_designated_gateway` - A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled.
* `allocate_new_eip` - When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway.
* `cloud_instance_id` - Instance ID of the gateway.
* `cloud_type` - Type of cloud service provider.
* `duo_api_hostname` - API hostname for DUO auth mode.
* `duo_integration_key` - Integration key for DUO auth mode.
* `duo_push_mode` - Push mode for DUO auth.
* `elb_dns_name` - ELB DNS Name.
* `elb_name` - Name of the ELB created.
* `enable_designated_gateway` - Status of Designated Gateway feature for Gateway.
* `enable_elb` - Status of ELB for the gateway.
* `enable_encrypt_volume` - Enable encrypt gateway EBS volume. Only supported for AWS provider.
* `enable_ldap` - Status LDAP or not.
* `enable_vpc_dns_server` - Status of Vpc Dns Server for Gateway.
* `enable_vpn_nat` - Status of VPN NAT.
* `gw_size` - Size of gateway Instance.
* `gw_name` - Aviatrix gateway name.
* `insane_mode` - Status of Insane Mode for Gateway.
* `insane_mode_az` - AZ of subnet being created for Insane Mode gateway.
* `ldap_bind_dn` - LDAP bind DN.
* `ldap_base_dn` - LDAP base DN.
* `ldap_server` - LDAP server address.
* `ldap_username_attribute` - LDAP user attribute.
* `max_vpn_conn` - Maximum connection of VPN access.
* `name_servers` - A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `okta_url` - URL for Okta auth mode.
* `okta_username_suffix` - Username suffix for Okta auth mode.
* `otp_mode` - Two step authentication mode.
* `peering_ha_cloud_instance_id` - Instance ID of the peering HA gateway.
* `peering_ha_gw_name` - Aviatrix gateway unique name of HA gateway.
* `peering_ha_gw_size` - Peering HA Gateway Size.
* `peering_ha_insane_mode_az` - AZ of subnet being created for Insane Mode Peering HA Gateway. Required if insane_mode is set.
* `peering_ha_private_ip` - Private IP address of HA gateway.
* `peering_ha_public_ip` - Public IP address that you want assigned to the HA peering instance.
* `peering_ha_subnet` - Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Required to create peering ha gateway if cloud_type = 1 or 8 (AWS or AZURE).
* `peering_ha_zone` - Zone information for creating Peering HA Gateway. Required to create peering ha gateway if cloud_type = 4 (gcp).
* `private_ip` - Private IP address of the Gateway created.
* `public_dns_server` - NS server used by the gateway.
* `public_ip` - Public IP address of the Gateway created.
* `saml_enabled` - Status of SAML.
* `search_domains` - A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `security_group_id` - Security group used for the gateway.
* `single_az_ha` - Status of Single AZ HA.
* `single_ip_snat` - Single IP Source NAT status for the container.
* `split_tunnel` - Status of split tunnel mode.
* `subnet` - A VPC Network address range selected from one of the available network ranges.
* `tag_list` - Instance tag of cloud provider.
* `vpc_id` - VPC-ID/VNet-Name of cloud provider.
* `vpc_reg` - Region of cloud provider.
* `vpn_access` - Status of user access through VPN to the container.
* `vpn_cidr` - VPN CIDR block for the container.
* `vpn_protocol` - Elb protocol for VPN gateway with elb enabled.

