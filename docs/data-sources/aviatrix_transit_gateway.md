---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway"
description: |-
  Gets an Aviatrix transit gateway's details.
---

# aviatrix_transit_gateway

The **aviatrix_transit_gateway** data source provides details about a specific transit gateway created by the Aviatrix Controller.

This data source can prove useful when a module accepts a transit gateway's detail as an input variable. For example, requiring the transit gateway's name for a spoke gateway's attachment.

## Example Usage

```hcl
# Aviatrix Transit Gateway Data Source
data "aviatrix_transit_gateway" "foo" {
  gw_name = "gatewayname"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Transit gateway name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `account_name` - Aviatrix account name.
* `allocate_new_eip` - When value is false, an idle address in Elastic IP pool is reused for this gateway. Otherwise, a new Elastic IP is allocated and used for this gateway.
* `bgp_manual_spoke_advertise_cidrs` - Intended CIDR list to advertise to VGW.
* `cloud_instance_id` - Instance ID of the transit gateway.
* `cloud_type` - Type of cloud service provider.
* `connected_transit"` -  Status of Connected Transit of transit gateway.
* `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes.
* `customized_transit_vpc_routes` - A list of CIDRs to be customized for the transit VPC routes.
* `gw_size` - Size of transit gateway instance.
* `gw_name` - Aviatrix transit gateway name.
* `insane_mode_az` - AZ of subnet being created for Insane Mode transit gateway.
* `enable_active_mesh` - Status of Active Mesh mode of the transit gateway.
* `enable_advertise_transit_cidr` - Status of Advertise Transit VPC network CIDR of the transit gateway.
* `enable_encrypt_volume` - Status of Encrypt Gateway EBS Volume of the transit gateway.
* `enable_firenet` - Status of FireNet Interfaces of the transit gateway.
* `enable_hybrid_connection` - Sign of readiness for TGW connection.
* `enable_learned_cidrs_approval` - Status of Encrypted Transit Approval for transit gateway.
* `enable_vpc_dns_server` - Status of Vpc Dns Server of the transit Gateway.
* `enable_transit_firenet` - Status of Transit FireNet Interfaces of the transit gateway.
* `enable_egress_transit_firenet` - Status of Egress Transit FireNet being enabled on the gateway.
* `excluded_advertised_spoke_routes` - A list of comma separated CIDRs to be advertised to on-prem as "Excluded CIDR List".
* `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table.
* `ha_insane_mode_az` - AZ of subnet being created for Insane Mode Transit HA Gateway.
* `ha_cloud_instance_id` - Cloud instance ID of HA transit gateway.
* `ha_gw_name` - Aviatrix transit gateway unique name of HA transit gateway.
* `ha_gw_size"` - HA Gateway Size.
* `ha_private_ip` - Private IP address that assigned to the HA Transit Gateway.
* `ha_public_ip` - Public IP address that assigned to the HA Transit Gateway.
* `ha_subnet` - HA Subnet.
* `ha_zone` - HA Zone.
* `insane_mode` - Status of Insane Mode of the transit gateway.
* `private_ip` - Private IP address of the transit gateway created.
* `public_ip` - Public IP address of the Transit Gateway created.
* `security_group_id` - Security group used for the transit gateway.
* `single_az_ha` - Status of Single AZ HA of transit gateway.
* `single_ip_snat` - Status of Single IP Source Nat mode of the transit gateway.
* `subnet` - A VPC Network address range selected from one of the available network ranges.
* `tag_list` - Instance tag of cloud provider.
* `vpc_id` - VPC-ID/VNet-Name of cloud provider.
* `vpc_reg` - Region of cloud provider.
