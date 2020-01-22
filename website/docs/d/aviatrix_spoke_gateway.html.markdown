---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway"
description: |-
  Gets the Aviatrix spoke gateway.
---

# aviatrix_spoke_gateway

Use this data source to get the Aviatrix spoke gateway for use in other resources.

## Example Usage

```hcl
# Aviatrix Spoke Gateway Data Source
data "aviatrix_spoke_gateway" "foo" {
  gw_name      = "gatewayname"
  account_name = "username"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Spoke gateway name. This can be used for getting spoke gateway.
* `account_name` - (Optional) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gw_name` - Aviatrix spoke gateway name.
* `account_name` - Aviatrix account name.
* `cloud_type` - Type of cloud service provider.
* `vpc_id` - VPC ID.
* `vpc_reg` - VPC Region.
* `gw_size` - Instance type.
* `subnet` - Range of the subnet where the spoke gateway is launched.
* `public_ip` - Public IP address of the spoke gateway created.
* `allocate_new_eip` - Description: "Whether the eip is newly allocated or not.
* `single_az_ha* `Enable/Disable this feature.
* `transit_gw` - The transit gateway that the spoke gateway is attached to.
* `tag_list` - Instance tag of cloud provider. Only supported for AWS provider.
* `insane_mode` - Enable/Disable Insane Mode for Spoke Gateway.
* `insane_mode_az` - AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.
* `enable_active_mesh` - Enable/Disable Active Mesh Mode for Spoke Gateway.
* `enable_vpc_dns_server` - Enable/Disalbe vpc_dns_server for Gateway.
* `enable_encrypt_volume` - Enable encrypt gateway EBS volume. Only supported for AWS provider.
* `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. It applies to this spoke gateway only​.
* `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s routing table. It applies to this spoke gateway only.
* `advertised_spoke_routes_include` - A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC.
* `cloud_instance_id` - Cloud instance ID

