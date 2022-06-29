---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateways"
description: |-
Gets a list of all Aviatrix transit gateway's details.
---


# aviatrix_all_transit_gateways

The **aviatrix_transit_gateways** data source provides details about all transit gateways created by the Aviatrix Controller.

## Example Usage

```hcl
# Aviatrix All Transit Gateways Data Source
data "aviatrix_transit_gateways" "foo" {}
```


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_list` - The list of all transit gateways
  * `account_name` - Aviatrix account name.
  * `availability_domain` - Availability domain for OCI.
  * `azure_eip_name_resource_group` - The name of the public IP address and its resource group in Azure to assign to this Transit Gateway.
  * `allocate_new_eip` - When value is false, an idle address in Elastic IP pool is reused for this gateway. Otherwise, a new Elastic IP is allocated and used for this gateway.
  * `bgp_ecmp` - Enable Equal Cost Multi Path (ECMP) routing for the next hop.
  * `bgp_hold_time` - BGP Hold Time.
  * `bgp_lan_interfaces` - Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP Transit.
  * `bgp_lan_ip_list` - List of available BGP LAN interface IPs for transit external device connection creation. Only supports GCP. Available as of provider version R2.21.0+.
  * `bgp_polling_time` - BGP route polling time. Unit is in seconds.
  * `cloud_instance_id` - Instance ID of the transit gateway.
  * `cloud_type` - Type of cloud service provider.
  * `connected_transit"` -  Status of Connected Transit of transit gateway.
  * `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes.
  * `enable_active_standby` - Enables Active-Standby Mode, available only with HA enabled.
  * `enable_active_standby_preemptive` - Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.
  * `enable_bgp_over_lan` - Pre-allocate a network interface(eth4) for \"BGP over LAN\" functionality. Only valid for cloud_type = 4 (GCP) and 8 (Azure). Available as of provider version R2.18+
  * `enable_gateway_load_balancer` - Enable firenet interfaces with AWS Gateway Load Balancer. 
  * `enable_jumbo_frame` - Enable jumbo frame support for transit gateway.
  * `enable_monitor_gateway_subnets` - Enable [monitor gateway subnets](https://docs.aviatrix.com/HowTos/gateway.html#monitor-gateway-subnet). Only valid for cloud_type = 1 (AWS) or 256 (AWSGov).
  * `enable_segmentation` - Enable segmentation to allow association of transit gateway to security domains.
  * `enable_spot_instance` - Enable spot instance. NOT supported for production deployment.
  * `enable_transit_summarize_cidr_to_tgw` - Enable summarize CIDR to TGW.
  * `enable_encrypt_volume` - Status of Encrypt Gateway EBS Volume of the transit gateway.
  * `enable_hybrid_connection` - Sign of readiness for TGW connection.
  * `enable_vpc_dns_server` - Status of Vpc Dns Server of the transit Gateway.
  * `enable_private_oob` - Status of private OOB for the transit gateway.
  * `enable_multi_tier_transit` - Status of multi-tier transit mode on transit gateway.
  * `excluded_advertised_spoke_routes` - A list of comma separated CIDRs to be advertised to on-prem as "Excluded CIDR List".
  * `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table.
  * `fault_domain` - Fault domain for OCI.
  * `gw_size` - Size of transit gateway instance.
  * `gw_name` - Aviatrix transit gateway name.
  * `ha_bgp_lan_interfaces` - Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP HA Transit.
  * `ha_bgp_lan_ip_list` - List of available BGP LAN interface IPs for transit external device HA connection creation. Only supports GCP. Available as of provider version R2.21.0+.
  * `insane_mode_az` - AZ of subnet being created for Insane Mode transit gateway.
  * `insane_mode` - Status of Insane Mode of the transit gateway.
  * `image_version` - The image version of the gateway.
  * `lan_private_subnet` - LAN Private Subnet. Only used for GCP Transit FireNet.
  * `lan_vpc_id` - LAN VPC ID. Only used for GCP Transit FireNet.
  * `learned_cidrs_approval_mode` - Set the learned CIDRs approval mode.
  * `local_as_number` - Changes the Aviatrix Transit Gateway ASN number before you setup Aviatrix Transit Gateway connection configurations.
  * `monitor_exclude_list` - A set of monitored instance ids. Only valid when 'enable_monitor_gateway_subnets' = true.
  * `oob_management_subnet` - OOB management subnet.
  * `oob_availability_zone` - OOB availability zone.
  * `private_ip` - Private IP address of the transit gateway created.
  * `public_ip` - Public IP address of the Transit Gateway created.
  * `prepend_as_path` - List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.
  * `security_group_id` - Security group used for the transit gateway.
  * `single_az_ha` - Status of Single AZ HA of transit gateway.
  * `single_ip_snat` - Status of Single IP Source Nat mode of the transit gateway.
  * `subnet` - A VPC Network address range selected from one of the available network ranges.
  * `spot_price` - Price for spot instance. NOT supported for production deployment.
  * `software_version` - The software version of the gateway.
  * `tunnel_detection_time` - The IPSec tunnel down detection time for the transit gateway.
  * `vpc_id` - VPC-ID/VNet-Name of cloud provider.
  * `vpc_reg` - Region of cloud provider.
  * `zone` - Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'.

