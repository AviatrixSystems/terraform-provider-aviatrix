---
subcategory: "Gateway"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
description: |-
  Gets an Aviatrix gateway's details.
---

# aviatrix_gateway

The **aviatrix_gateway** data source provides details about a specific gateway created by the Aviatrix Controller.

This data source can prove useful when a module accepts a gateway's detail as an input variable. For example, requiring the gateway's name configuring a site2cloud connection.

## Example Usage

```hcl
# Aviatrix Gateway Data Source
data "aviatrix_gateway" "foo" {
  gw_name = "gatewayname"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Gateway name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `account_name` - Aviatrix account name.
* `additional_cidrs` - A list of destination CIDR ranges that will also go through the VPN tunnel when Split Tunnel Mode is enabled.
* `additional_cidrs_designated_gateway` - A list of CIDR ranges separated by comma to configure when 'designated_gateway' feature is enabled.
* `allocate_new_eip` - When value is false, an idle address in Elastic IP pool is reused for this gateway. Otherwise, a new Elastic IP is allocated and used for this gateway.
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
* `enable_vpc_dns_server` - Status of VPC Dns Server for Gateway.
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
* `peering_ha_subnet` - Public Subnet Information while creating Peering HA Gateway, only subnet is accepted. Required to create peering ha gateway if cloud_type = 1 or 8 (AWS or Azure).
* `peering_ha_zone` - Zone information for creating Peering HA Gateway. Required to create peering ha gateway if cloud_type = 4 (GCP).
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
* `vpc_id` - VPC-ID/VNet-Name of cloud provider.
* `vpc_reg` - Region of cloud provider.
* `vpn_access` - Status of user access through VPN to the container.
* `vpn_cidr` - VPN CIDR block for the container.
* `vpn_protocol` - ELB protocol for VPN gateway with ELB enabled.
* `availability_domain` - Availability domain for OCI.
* `fault_domain` - Fault domain for OCI.
* `peering_ha_availability_domain` - HA gateway availability domain for OCI.
* `peering_ha_fault_domain` - HA gateway fault domain for OCI.
* `software_version` - The software version of the gateway.
* `image_version` - The image version of the gateway.
* `peering_ha_software_version` - The software version of the HA gateway.
* `peering_ha_image_version` - The image version of the HA gateway.

The following argument is deprecated:

* `tag_list` - Instance tag of cloud provider.
