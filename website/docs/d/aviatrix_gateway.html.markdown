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

* `gw_name` - Aviatrix gateway name.
* `cloud_type` - Type of cloud service provider.
* `account_name` - Account name.
* `vpc_id` - VPC-ID/VNet-Name of cloud provider.
* `vpc_reg` - Region of cloud provider.
* `gw_size` - Size of gateway Instance.
* `subnet` - A VPC Network address range selected from one of the available network ranges.
* `insane_mode_az` - AZ of subnet being created for Insane Mode gateway.
* `single_ip_snat` - Single IP Source NAT status for the container.
* `vpn_access` - Status of user access through VPN to the container.
* `vpn_cidr` - VPN CIDR block for the container.
* `enable_elb` - Status of ELB for the gateway.
* `elb_name` - Name of the ELB created.
* `vpn_protocol` - Elb protocol for VPN gateway with elb enabled.
* `split_tunnel` - Status of split tunnel mode.
* `max_vpn_conn` - Maximum connection of VPN access.
* `name_servers` - A list of DNS servers used to resolve domain names by a connected VPN user when Split Tunnel Mode is enabled.
* `search_domains` - A list of domain names that will use the NameServer when a specific name is not in the destination when Split Tunnel Mode is enabled.
* `additional_cidrs` - A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `otp_mode` - Two step authentication mode.
* `saml_enabled` - Status of SAML.
* `enable_vpn_nat` - Status of VPN NAT.
* `okta_url` - URL for Okta auth mode.
* `okta_username_suffix` - Username suffix for Okta auth mode.
* `duo_integration_key` - Integration key for DUO auth mode.
* `duo_api_hostname` - API hostname for DUO auth mode.
* `duo_push_mode` - Push mode for DUO auth.
* `enable_ldap` - Status LDAP or not.
* `ldap_server` - LDAP server address.
* `ldap_bind_dn` - LDAP bind DN.
* `ldap_base_dn` - LDAP base DN.
* `ldap_username_attribute` - LDAP user attribute.
* `peering_ha_subnet` - Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Required to create peering ha gateway if cloud_type = 1 or 8 (AWS or AZURE).
* `peering_ha_zone` - Zone information for creating Peering HA Gateway. Required to create peering ha gateway if cloud_type = 4 (gcp).
* `peering_ha_insane_mode_az` - AZ of subnet being created for Insane Mode Peering HA Gateway. Required if insane_mode is set.
* `peering_ha_public_ip` - Public IP address that you want assigned to the HA peering instance.
* `peering_ha_gw_size` - Peering HA Gateway Size.
* `single_az_ha` - Status of Single AZ HA.
* `allocate_new_eip` - When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway.
* `tag_list` - Instance tag of cloud provider.
* `insane_mode` - Status of Insane Mode for Gateway.
* `enable_vpc_dns_server` - Status of Vpc Dns Server for Gateway.
* `enable_designated_gateway` - Status of Designated Gateway feature for Gateway.
* `additional_cidrs_designated_gateway` - A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled.
* `enable_encrypt_volume` - Enable encrypt gateway EBS volume. Only supported for AWS provider.
* `elb_dns_name` - ELB DNS Name.
* `public_ip` - Public IP address of the Gateway created.
* `security_group_id` - Security group used for the gateway.
* `public_dns_server` - NS server used by the gateway.
* `cloud_instance_id` - Instance ID of the gateway.
* `private_ip` - Private IP address of the Gateway created.
* `peering_ha_cloud_instance_id` - Instance ID of the peering HA gateway.
* `peering_ha_gw_name` - Aviatrix gateway unique name of HA gateway.