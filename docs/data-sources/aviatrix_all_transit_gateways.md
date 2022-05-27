---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_all_transit_gateways"
description: |-
Gets a list of all Aviatrix transit gateway's details.
---


# aviatrix_all_transit_gateways

The **aviatrix_all_transit_gateways** data source provides details about all transit gateways created by the Aviatrix Controller.

## Example Usage

```hcl
# Aviatrix All Transit Gateways Data Source
data "aviatrix_all_transit_gateways" "foo" {}
```


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `transit_gateway_list` - The list of all transit gateways
  * `account_name` - Aviatrix account name.
  * `allocate_new_eip` - When value is false, an idle address in Elastic IP pool is reused for this gateway. Otherwise, a new Elastic IP is allocated and used for this gateway.
  * `cloud_instance_id` - Instance ID of the transit gateway.
  * `cloud_type` - Type of cloud service provider.
  * `connected_transit"` -  Status of Connected Transit of transit gateway.
  * `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes.
  * `gw_size` - Size of transit gateway instance.
  * `gw_name` - Aviatrix transit gateway name.
  * `insane_mode_az` - AZ of subnet being created for Insane Mode transit gateway.
  * `enable_encrypt_volume` - Status of Encrypt Gateway EBS Volume of the transit gateway.
  * `enable_hybrid_connection` - Sign of readiness for TGW connection.
  * `enable_vpc_dns_server` - Status of Vpc Dns Server of the transit Gateway.
  * `excluded_advertised_spoke_routes` - A list of comma separated CIDRs to be advertised to on-prem as "Excluded CIDR List".
  * `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table.
  * `insane_mode` - Status of Insane Mode of the transit gateway.
  * `private_ip` - Private IP address of the transit gateway created.
  * `public_ip` - Public IP address of the Transit Gateway created.
  * `security_group_id` - Security group used for the transit gateway.
  * `single_az_ha` - Status of Single AZ HA of transit gateway.
  * `single_ip_snat` - Status of Single IP Source Nat mode of the transit gateway.
  * `subnet` - A VPC Network address range selected from one of the available network ranges.
  * `vpc_id` - VPC-ID/VNet-Name of cloud provider.
  * `vpc_reg` - Region of cloud provider.
  * `enable_private_oob` - Status of private OOB for the transit gateway.
  * `oob_management_subnet` - OOB management subnet.
  * `oob_availability_zone` - OOB availability zone.
  * `enable_multi_tier_transit` - Status of multi-tier transit mode on transit gateway.
  * `availability_domain` - Availability domain for OCI.
  * `fault_domain` - Fault domain for OCI.
  * `software_version` - The software version of the gateway.
  * `image_version` - The image version of the gateway.

