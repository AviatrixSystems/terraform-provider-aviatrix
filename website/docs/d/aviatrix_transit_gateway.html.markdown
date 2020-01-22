---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway"
description: |-
  Gets the Aviatrix transit gateway.
---

# aviatrix_transit_gateway

Use this data source to get the Aviatrix transit gateway for use in other resources.

## Example Usage

```hcl
# Aviatrix Transit Gateway Data Source
data "aviatrix_transit_gateway" "foo" {
  gw_name      = "gatewayname"
  account_name = "username"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Transit gateway name. This can be used for getting transit gateway.
* `account_name` - (Optional) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gw_name` - Aviatrix transit gateway name.
* `account_name` - Aviatrix account name.
* `cloud_type` - Type of cloud service provider.
* `vpc_id` - VPC ID.
* `vpc_reg` - VPC Region.
* `gw_size` - Instance type.
* `subnet` - Range of the subnet where the transit gateway is launched.
* `public_ip` - Public IP address of the Gateway created.
* `allocate_new_eip` - Whether the eip is newly allocated or not.
* `single_az_ha` - Enable/Disable this feature.
* `tag_list` - Instance tag of cloud provider. Only supported for AWS provider.
* `enable_hybrid_connection` - Sign of readiness for TGW connection.
* `connected_transit` - Connected Transit status.
* `insane_mode` - Enable/Disable Insane Mode for Spoke Gateway.
* `insane_mode_az` - AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.
* `enable_firenet` - Whether firenet interfaces is enabled.
* `enable_active_mesh` - Enable/Disable active mesh mode for Transit Gateway.
* `enable_vpc_dns_server` - Enable/Disable vpc_dns_server for Gateway. Only supports AWS.
* `enable_advertise_transit_cidr` - Enable/Disable advertise transit VPC network CIDR.
* `bgp_manual_spoke_advertise_cidrs` - Intended CIDR list to advertise to VGW.
* `enable_encrypt_volume` - Enable/Disable encrypt gateway EBS volume. Only supported for AWS provider.
* `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. It applies to all spoke gateways attached to this transit gateway.
* `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s routing table. It applies to all spoke gateways attached to this transit gateway.
* `advertised_spoke_routes_exclude` - A list of comma separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'. When configured, it inspects all the advertised CIDRs from its spoke gateways and remove those included in the 'Excluded CIDR List'.


