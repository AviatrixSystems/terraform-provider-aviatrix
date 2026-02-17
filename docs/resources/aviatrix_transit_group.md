---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_group"
description: |-
  Creates and manages Aviatrix transit gateway groups
---

# aviatrix_transit_group

The **aviatrix_transit_group** resource allows the creation and management of Aviatrix transit gateway groups.

~> **NOTE:** This resource is available as of provider version R3.0+ and requires controller version 7.0+.

## What is Gateway Group

The Gateway Group feature allows for better horizontal scaling by moving from the "primary + HA" gateway-pair model to support a gateway grouping, represented by a "primary + N number of HA gateways". The Gateway Group feature allows users to create N number of HA gateways under a primary gateway.

For more information, see the [Introduction to Gateway Group](../guides/introduction_to_gateway_group.html) guide.

## Example Usage

```hcl
# Create an Aviatrix AWS Transit Group
resource "aviatrix_transit_group" "test_transit_group_aws" {
  group_name          = "my-transit-group"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "transit"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  vpc_region          = "us-west-1"
}
```

```hcl
# Create an Aviatrix AWS Transit Group with BGP enabled
resource "aviatrix_transit_group" "test_transit_group_bgp" {
  group_name          = "my-transit-group-bgp"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "transit"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  vpc_region          = "us-west-1"

  enable_bgp             = true
  local_as_number        = "65001"
  prepend_as_path        = ["65001", "65001"]
  bgp_polling_time       = 30
  bgp_hold_time          = 120
  enable_bgp_ecmp        = true
}
```

```hcl
# Create an Aviatrix GCP Transit Group
resource "aviatrix_transit_group" "test_transit_group_gcp" {
  group_name          = "my-gcp-transit-group"
  cloud_type          = 4
  account_name        = "my-gcp-account"
  gw_type             = "transit"
  group_instance_size = "n1-standard-1"
  vpc_id              = "gcp-vpc-name~-~project-id"
  vpc_region          = "us-west1-b"
}
```

```hcl
# Create an Aviatrix Azure Transit Group
resource "aviatrix_transit_group" "test_transit_group_azure" {
  group_name          = "my-azure-transit-group"
  cloud_type          = 8
  account_name        = "my-azure-account"
  gw_type             = "transit"
  group_instance_size = "Standard_B2ms"
  vpc_id              = "vnet_name:rg_name:resource_guid"
  vpc_region          = "West US"
}
```

```hcl
# Create an Aviatrix Transit Group with advanced features
resource "aviatrix_transit_group" "test_transit_group_advanced" {
  group_name          = "my-transit-group-advanced"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "transit"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  vpc_region          = "us-west-1"

  # Transit-specific features
  enable_connected_transit      = true
  enable_segmentation           = true
  enable_advertise_transit_cidr = true

  # BGP Configuration
  enable_bgp                       = true
  local_as_number                  = "65001"
  prepend_as_path                  = ["65001", "65001"]
  bgp_manual_spoke_advertise_cidrs = ["10.10.0.0/16", "10.20.0.0/16"]
  enable_preserve_as_path          = true

  # BGP Timers
  bgp_polling_time                    = 30
  bgp_neighbor_status_polling_time    = 3
  bgp_hold_time                       = 120

  # BGP Communities
  bgp_send_communities   = true
  bgp_accept_communities = true

  # Active-Standby
  enable_active_standby            = true
  enable_active_standby_preemptive = false

  # Feature Flags
  enable_nat                          = true
  enable_jumbo_frame                  = true
  enable_gro_gso                      = true
  enable_ipv6                         = false
  enable_vpc_dns_server               = true

  # Learned CIDR Approval
  enable_learned_cidrs_approval = true
  learned_cidrs_approval_mode   = "gateway"
  approved_learned_cidrs        = ["192.168.1.0/24"]

  # Insane Mode
  insane_mode = false
}
```

## Argument Reference

The following arguments are supported:

### Required

* `group_name` - (Required) Name of the gateway group.
* `cloud_type` - (Required) Type of cloud service provider. Valid values: 1 (AWS), 4 (GCP), 8 (Azure), 16 (OCI), 32 (AzureGov), 256 (AWSGov), 1024 (AWSChina), 2048 (AzureChina), 8192 (Alibaba Cloud), 16384 (AWSTop Secret), 32768 (AWSSecret).
* `gw_type` - (Required) Gateway type for the group. Valid values: "TRANSIT", "EDGETRANSIT", "STANDALONE". Case-insensitive.
* `group_instance_size` - (Required) Instance size for gateways in the group. Example: "t3.medium" (AWS), "n1-standard-1" (GCP), "Standard_B2ms" (Azure).
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider.
  * AWS/AWSGov/AWSChina: VPC ID (e.g., "vpc-abcd1234")
  * GCP: VPC name and project ID separated by "~-~" (e.g., "vpc-name~-~project-id")
  * Azure/AzureGov/AzureChina: VNet name, resource group, and resource GUID separated by ":" (e.g., "vnet_name:rg_name:resource_guid")
* `account_name` - (Required) Name of the Cloud-Account in Aviatrix controller.
* `vpc_region` - (Required) VPC region for CSP deployments. Optional for edge deployments (gw_type = "EDGETRANSIT").

### Optional - General Settings

* `customized_cidr_list` - (Optional) Set of customized CIDRs for the transit group.
* `domain` - (Optional) Network domain for the transit group.

### Optional - Feature Flags

* `enable_nat` - (Optional) Enable NAT (aka single_ip_snat). Valid values: true, false. Default: false.
* `enable_jumbo_frame` - (Optional) Enable jumbo frame support. Valid values: true, false. Default: true.
* `enable_ipv6` - (Optional) Enable IPv6. Only valid for AWS and Azure. Valid values: true, false. Default: false.
* `enable_gro_gso` - (Optional) Enable GRO/GSO. Valid values: true, false. Default: true.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server. Valid values: true, false. Default: false.
* `enable_s2c_rx_balancing` - (Optional) Enable S2C receive balancing. Valid values: true, false. Default: false.

### Optional - Transit-Specific Features

* `enable_hybrid_connection` - (Optional) Enable hybrid connection (TGW connection readiness). Valid values: true, false. Default: false.
* `enable_connected_transit` - (Optional) Enable connected transit. Valid values: true, false. Default: false.
* `enable_firenet` - (Optional) Enable FireNet. Valid values: true, false. Default: false.
* `enable_transit_firenet` - (Optional) Enable Transit FireNet. Valid values: true, false. Default: false.
* `enable_advertise_transit_cidr` - (Optional) Enable advertise transit CIDR. Valid values: true, false. Default: false.
* `enable_transit_summarize_cidr_to_tgw` - (Optional) Enable transit summarize CIDR to TGW. Valid values: true, false. Default: false.
* `enable_multi_tier_transit` - (Optional) Enable multi-tier transit. Valid values: true, false. Default: false.
* `enable_segmentation` - (Optional) Enable segmentation (LAN segmentation). Valid values: true, false. Default: false.
* `enable_gateway_load_balancer` - (Optional) Enable AWS Gateway Load Balancer. Only valid for AWS. Valid values: true, false. Default: false.

### Optional - BGP Configuration

* `enable_bgp` - (Optional) Enable BGP. Valid values: true, false. Default: false. **Note:** Changing this forces a new resource to be created.
* `local_as_number` - (Optional) BGP local AS number. Required when `enable_bgp` is set to true.
* `prepend_as_path` - (Optional) List of AS numbers to prepend to the AS_Path field. Valid only when `local_as_number` is set. Example: ["65001", "65001"].
* `bgp_manual_spoke_advertise_cidrs` - (Optional) Set of intended CIDRs to be advertised to spoke gateways via BGP. Example: ["10.1.0.0/16", "10.2.0.0/16"].
* `enable_preserve_as_path` - (Optional) Preserve AS path when advertising manual summary CIDRs. Valid values: true, false. Default: false.
* `enable_bgp_ecmp` - (Optional) Enable BGP ECMP. Valid values: true, false. Default: false.

### Optional - BGP Timers

* `bgp_polling_time` - (Optional) BGP route polling time in seconds. Valid values: 10-50. Default: 50.
* `bgp_neighbor_status_polling_time` - (Optional) BGP neighbor status polling time in seconds. Valid values: 1-10. Default: 5.
* `bgp_hold_time` - (Optional) BGP hold time in seconds. Valid values: 12-360. Default: 180.

### Optional - BGP Communities

* `bgp_send_communities` - (Optional) Send BGP communities. Valid values: true, false. Default: false.
* `bgp_accept_communities` - (Optional) Accept BGP communities. Valid values: true, false. Default: false.

### Optional - BGP over LAN

* `enable_bgp_over_lan` - (Optional) Enable BGP over LAN. Valid values: true, false. Default: false. **Note:** Changing this forces a new resource to be created.

### Optional - Learned CIDR Approval

* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDR approval. Valid values: true, false. Default: false.
* `learned_cidrs_approval_mode` - (Optional) Learned CIDRs approval mode. Valid values: "gateway", "connection". Default: "gateway".
* `approved_learned_cidrs` - (Optional) Set of approved learned CIDRs. Valid only when `enable_learned_cidrs_approval` is true. Example: ["10.1.0.0/16", "10.2.0.0/16"].

### Optional - Active-Standby

* `enable_active_standby` - (Optional) Enable Active-Standby mode. Valid values: true, false. Default: false.
* `enable_active_standby_preemptive` - (Optional) Enable Active-Standby Preemptive mode. Valid only when `enable_active_standby` is true. Valid values: true, false. Default: false.

### Optional - AWS Specific

* `insane_mode` - (Optional) Enable Insane Mode (High Performance Encryption). Valid values: true, false. Default: false. **Note:** Changing this forces a new resource to be created.

### Optional - GCP Specific

* `enable_global_vpc` - (Optional) Enable global VPC. Valid values: true, false. Default: false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `group_uuid` - Gateway group UUID.
* `gw_uuid_list` - List of gateway UUIDs in the group.
* `vpc_uuid` - VPC UUID.
* `vendor_name` - Cloud vendor name (e.g., "AWS", "GCP", "Azure").
* `explicitly_created` - Indicates if the group was explicitly created.

## Import

**transit_group** can be imported using the `group_uuid`, e.g.

```
$ terraform import aviatrix_transit_group.test <group_uuid>
```

## Notes

* The `group_name` is used to identify the gateway group and must be unique within the controller.
* BGP configuration options are only applicable when `enable_bgp` is set to true.
* The `bgp_polling_time`, `bgp_neighbor_status_polling_time`, and `bgp_hold_time` fields will default to their respective values (50, 5, and 180 seconds) when the API returns 0.
* The `learned_cidrs_approval_mode` will default to "gateway" when the API returns an empty string.
* Active-Standby Preemptive mode can only be enabled when Active-Standby mode is already enabled.
* Transit groups support additional features compared to spoke groups, including connected transit, segmentation, and advertise transit CIDR.
