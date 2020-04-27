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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gw_name` - Aviatrix spoke gateway name.
* `cloud_type` - Type of cloud service provider.
* `account_name` - The name of a Cloud-Account in Aviatrix controller.
* `vpc_id` - VPC-ID/VNet-Name of cloud provider.
* `vpc_reg` - Region of cloud provider.
* `gw_size` - Size of the gateway instance.
* `subnet` - Public Subnet Info.
* `insane_mode_az` - AZ of subnet being created for Insane Mode Spoke Gateway.
* `single_ip_snat` - If Source NAT feature in 'single_ip' mode is on the gateway.
* `allocate_new_eip` - If allocating a new ip for this gateway.
* `public_ip` - Public IP of spoke gateway.
* `ha_subnet` - HA Subnet. Required if enabling HA for AWS/AZURE.
* `ha_zone` - HA Zone. Required if enabling HA for GCP.
* `ha_insane_mode_az` - AZ of subnet being created for Insane Mode Spoke HA Gateway. Required if insane_mode is true and ha_subnet is set.
* `ha_gw_size` - HA Gateway Size.
* `ha_public_ip` - Public IP address of the HA Spoke Gateway.
* `single_az_ha` - If this feature is desired.
* `transit_gw` - The transit gateway to attach this spoke gateway to.
* `tag_list` - Instance tag of cloud provider.
* `insane_mode` - If Insane Mode is enabled for Spoke Gateway. If is enabled, gateway size has to at least be c5 size.
* `enable_active_mesh` - If Active Mesh Mode for spoke gateway.
* `enable_vpc_dns_server` - If vpc_dns_server is enabled for spoke gateway. Only supports AWS.
* `enable_encrypt_volume` - If encrypt gateway EBS volume is enabled for spoke gateway. Only supported for AWS provider. 
* `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes. 
* `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table. 
* `included_advertised_spoke_routes` - A list of comma separated CIDRs to be advertised to on-prem as "Included CIDR List".​
* `security_group_id` - Security group used for the spoke gateway.
* `cloud_instance_id` - Cloud instance ID.
* `private_ip` - Private IP address of the spoke gateway.
* `ha_cloud_instance_id` - Cloud instance ID of HA spoke gateway.
* `ha_gw_name` - Aviatrix spoke gateway unique name of HA spoke gateway.