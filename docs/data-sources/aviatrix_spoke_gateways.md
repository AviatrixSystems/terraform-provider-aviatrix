---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateways"
description: |-
  Gets a list of all Aviatrix spoke gateway's details.
---

# aviatrix_spoke_gateways

The **aviatrix_spoke_gateways** data source provides details about all spoke gateways created by the Aviatrix Controller.

## Example Usage

```hcl
# Aviatrix Spoke Gateways Data Source
data "aviatrix_spoke_gateways" "foo" {}
```

## Attribute Reference

The following attributes are exported:

* `gateway_list` - The list of all spoke gateways
    * `cloud_type` - Type of cloud service provider.
    * `account_name` - Aviatrix account name.
    * `gw_name` - Aviatrix spoke gateway name.
    * `gw_size` - Size of spoke gateway instance.
    * `vpc_id` - VPC-ID/VNet-Name of cloud provider.
    * `vpc_reg` - Region of cloud provider.
    * `zone` - Availability Zone. Only available for cloud_type = 8 (Azure). Must be in the form 'az-n', for example, 'az-2'.
    * `subnet` - A VPC Network address range selected from one of the available network ranges.
    * `insane_mode_az` - AZ of subnet being created for Insane Mode spoke gateway.
    * `single_ip_snat` - Status of Single IP Source Nat mode of the spoke gateway.
    * `allocate_new_eip` - When value is false, an idle address in Elastic IP pool is reused for this gateway. Otherwise, a new Elastic IP is allocated and used for this gateway.
    * `public_ip` - Public IP address of the Spoke Gateway.
    * `single_az_ha` - Status of Single AZ HA of spoke gateway.
    * `transit_gw` - Transit Gateways this spoke has joined.
    * `insane_mode` - Status of Insane Mode of the spoke gateway.
    * `enable_vpc_dns_server` - Status of Vpc Dns Server of the spoke Gateway.
    * `enable_encrypt_volume` - Status of Encrypt Gateway EBS Volume of the spoke gateway.
    * `customized_spoke_vpc_routes` - A list of comma separated CIDRs to be customized for the spoke VPC routes.
    * `filtered_spoke_vpc_routes` - A list of comma separated CIDRs to be filtered from the spoke VPC route table.
    * `included_advertised_spoke_routes` - A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC.
    * `security_group_id` - Security group used for the spoke gateway.
    * `cloud_instance_id` - Instance ID of the spoke gateway.
    * `private_ip` - Private IP address of the spoke gateway created.
    * `enable_private_oob` - Status of private OOB for the spoke gateway.
    * `oob_management_subnet` - OOB management subnet.
    * `oob_availability_zone` - OOB availability zone.
    * `tunnel_detection_time` - The IPSec tunnel down detection time for the spoke gateway.
    * `availability_domain` - Availability domain for OCI.
    * `fault_domain` - Fault domain for OCI.
    * `software_version` - The software version of the gateway.
    * `image_version` - The image version of the gateway.
    * `enable_monitor_gateway_subnets` - Enable monitor gateway subnets. Only valid for cloud_type = 1 (AWS) or 256 (AWSGov).
    * `monitor_exclude_list` - A set of monitored instance ids. Only valid when 'enable_monitor_gateway_subnets' = true.
    * `enable_jumbo_frame` - Enable jumbo frame support for spoke gateway.
    * `enable_private_vpc_default_route` - Enable Private VPC Default Route.
    * `enable_skip_public_route_table_update` - Skip Public Route Table Update.
    * `enable_auto_advertise_s2c_cidrs` - Automatically advertise remote CIDR to Aviatrix Spoke Gateway when route based Site2Cloud Tunnel is created.
    * `spoke_bgp_manual_advertise_cidrs` - Intended CIDR list to be advertised to external BGP router.
    * `enable_bgp` - Enable BGP.
    * `enable_learned_cidrs_approval` - Switch to enable/disable encrypted transit approval for BGP Spoke Gateway.
    * `learned_cidrs_approval_mode` - Set the learned CIDRs approval mode.
    * `bgp_ecmp` - Enable Equal Cost Multi Path (ECMP) routing for the next hop.
    * `enable_active_standby` - Enables Active-Standby Mode, available only with HA enabled.
    * `enable_active_standby_preemptive` - Enables Preemptive Mode for Active-Standby, available only with Active-Standby enabled.
    * `disable_route_propagation` - Disables route propagation on BGP Spoke to attached Transit Gateway.
    * `local_as_number` - Changes the Aviatrix Spoke Gateway ASN number before you setup Aviatrix Spoke Gateway connection configurations.
    * `prepend_as_path` - List of AS numbers to populate BGP AP_PATH field when it advertises to VGW or peer devices.
    * `bgp_polling_time` - BGP route polling time. Unit is in seconds.
    * `bgp_hold_time` - BGP Hold Time.
    * `enable_spot_instance` - Enable spot instance. NOT supported for production deployment.
    * `spot_price` - Price for spot instance. NOT supported for production deployment.
    * `azure_eip_name_resource_group` - The name of the public IP address and its resource group in Azure to assign to this Spoke Gateway.
    * `eip` - The EIP address of the Spoke Gateway.
