---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway"
description: |-
  Gets an Aviatrix spoke gateway's details.
---

# aviatrix_spoke_gateway

The **aviatrix_spoke_gateway** data source provides details about a specific spoke gateway created by the Aviatrix Controller.

This data source can prove useful when a module accepts a spoke gateway's detail as an input variable.

## Example Usage

```hcl
# Aviatrix Spoke Gateway Data Source
data "aviatrix_spoke_gateway" "foo" {
  gw_name = "gatewayname"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Spoke gateway name. It can be used for getting spoke gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `account_name` - Aviatrix account name.
* `allocate_new_eip` - When value is false, an idle address in Elastic IP pool is reused for this gateway. Otherwise, a new Elastic IP is allocated and used for this gateway.
* `cloud_instance_id` - Cloud instance ID.
* `cloud_type` - Type of cloud service provider.
* `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes.
* `enable_encrypt_volume` - Status of Encrypt Volume of spoke gateway.
* `enable_vpc_dns_server` - Status of VPC Dns Server of spoke gateway.
* `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table.
* `ha_cloud_instance_id` - Cloud instance ID of HA spoke gateway.
* `ha_insane_mode_az` - AZ of subnet being created for Insane Mode Spoke HA Gateway.
* `ha_gw_name` - Aviatrix spoke gateway unique name of HA spoke gateway.
* `ha_gw_size` - HA Gateway Size.
* `ha_private_ip` - Private IP address of HA spoke gateway.
* `ha_public_ip` - Public IP address of the HA spoke gateway.
* `ha_subnet` - HA Subnet.
* `ha_zone` - HA Zone.
* `gw_name` - Aviatrix spoke gateway name.
* `gw_size` - Size of spoke gateway instance.
* `included_advertised_spoke_routes` - A list of comma separated CIDRs to be advertised to on-prem as "Included CIDR List".
* `insane_mode` - Status of Insane Mode for Spoke Gateway.
* `insane_mode_az` - AZ of subnet being created for Insane Mode spoke gateway.
* `private_ip` - Private IP address of the spoke gateway.
* `public_ip` - Public IP of spoke gateway.
* `security_group_id` - Security group used of the spoke gateway.
* `single_az_ha` - Status of Single AZ HA of spoke gateway.
* `single_ip_snat` - Status of Single IP Source NAT mode of the spoke gateway.
* `subnet` - A VPC Network address range selected from one of the available network ranges.
* `tag_list` - Instance tag of cloud provider.
* `transit_gw` - Transit gateways attached to this spoke gateway.
* `vpc_id` - VPC-ID/VNet-Name of cloud provider.
* `vpc_reg` - Region of cloud provider.
* `enable_private_oob` - Status of private OOB for the spoke gateway.
* `oob_management_subnet` - OOB management subnet.
* `oob_availability_zone` - OOB availability zone.
* `ha_oob_management_subnet` - HA OOB management subnet.
* `ha_oob_availability_zone` - HA OOB availability zone.
* `availability_domain` - Availability domain for OCI.
* `fault_domain` - Fault domain for OCI.
* `ha_availability_domain` - HA gateway availability domain for OCI.
* `ha_fault_domain` - HA gateway fault domain for OCI.
* `software_version` - The software version of the gateway.
* `image_version` - The image version of the gateway.
* `ha_software_version` - The software version of the HA gateway.
* `ha_image_version` - The image version of the HA gateway.
* `eip` - The EIP address of the Spoke Gateway.
* `ha_eip` - The EIP address of the HA Spoke Gateway.
