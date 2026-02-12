---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_group"
description: |-
  Creates and manages Aviatrix spoke gateway groups
---

# aviatrix_spoke_group

The **aviatrix_spoke_group** resource allows the creation and management of Aviatrix spoke gateway groups.

~> **NOTE:** This resource is available as of provider version R3.0+ and requires controller version 7.0+.

## What is Gateway Group

The Gateway Group feature allows for better horizontal scaling by moving from the "primary + HA" gateway-pair model to support a gateway grouping, represented by a "primary + N number of HA gateways". The Gateway Group feature allows users to create N number of HA gateways under a primary gateway.

For more information, see the [Introduction to Gateway Group](../guides/introduction_to_gateway_group.html) guide.

## Example Usage

```hcl
# Create an Aviatrix AWS Spoke Group
resource "aviatrix_spoke_group" "test_spoke_group_aws" {
  group_name          = "my-spoke-group"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "spoke"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  vpc_region          = "us-west-1"
}
```

```hcl
# Create an Aviatrix AWS Spoke Group with BGP enabled
resource "aviatrix_spoke_group" "test_spoke_group_bgp" {
  group_name          = "my-spoke-group-bgp"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "spoke"
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
# Create an Aviatrix GCP Spoke Group
resource "aviatrix_spoke_group" "test_spoke_group_gcp" {
  group_name          = "my-gcp-spoke-group"
  cloud_type          = 4
  account_name        = "my-gcp-account"
  gw_type             = "spoke"
  group_instance_size = "n1-standard-1"
  vpc_id              = "gcp-vpc-name~-~project-id"
  vpc_region          = "us-west1-b"
}
```

```hcl
# Create an Aviatrix Azure Spoke Group
resource "aviatrix_spoke_group" "test_spoke_group_azure" {
  group_name          = "my-azure-spoke-group"
  cloud_type          = 8
  account_name        = "my-azure-account"
  gw_type             = "spoke"
  group_instance_size = "Standard_B2ms"
  vpc_id              = "vnet_name:rg_name:resource_guid"
  vpc_region          = "West US"
}
```

```hcl
# Create an Aviatrix Spoke Group with advanced features
resource "aviatrix_spoke_group" "test_spoke_group_advanced" {
  group_name          = "my-spoke-group-advanced"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "spoke"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  vpc_region          = "us-west-1"

  # BGP Configuration
  enable_bgp                       = true
  local_as_number                  = "65001"
  prepend_as_path                  = ["65001", "65001"]
  spoke_bgp_manual_advertise_cidrs = ["10.10.0.0/16", "10.20.0.0/16"]
  enable_preserve_as_path          = true
  enable_auto_advertise_s2c_cidrs  = true
  disable_route_propagation        = false

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
  enable_private_vpc_default_route    = true
  enable_skip_public_route_table_update = false

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
* `gw_type` - (Required) Gateway type for the group. Valid values: "SPOKE", "EDGESPOKE", "STANDALONE". Case-insensitive.
* `group_instance_size` - (Required) Instance size for gateways in the group. Example: "t3.medium" (AWS), "n1-standard-1" (GCP), "Standard_B2ms" (Azure).
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider.
  * AWS/AWSGov/AWSChina: VPC ID (e.g., "vpc-abcd1234")
  * GCP: VPC name and project ID separated by "~-~" (e.g., "vpc-name~-~project-id")
  * Azure/AzureGov/AzureChina: VNet name, resource group, and resource GUID separated by ":" (e.g., "vnet_name:rg_name:resource_guid")
* `account_name` - (Required) Name of the Cloud-Account in Aviatrix controller.
* `vpc_region` - (Optional) Region of cloud provider. Required for CSP (Cloud Service Provider) deployments.

### Optional - General Settings

* `customized_cidr_list` - (Optional) Set of customized CIDRs for the spoke group.
* `domain` - (Optional) Network domain for the spoke group.
* `include_cidr` - (Optional) Set of CIDRs to include for the spoke group.

### Optional - Feature Flags

* `enable_nat` - (Optional) Enable NAT. Valid values: true, false. Default: false.
* `enable_jumbo_frame` - (Optional) Enable jumbo frame support. Valid values: true, false. Default: true.
* `enable_ipv6` - (Optional) Enable IPv6. Valid values: true, false. Default: false.
* `enable_gro_gso` - (Optional) Enable GRO/GSO. Valid values: true, false. Default: true.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server. Valid values: true, false. Default: false.
* `enable_private_vpc_default_route` - (Optional) Enable private VPC default route. Valid values: true, false. Default: false.
* `enable_skip_public_route_table_update` - (Optional) Skip updating public route tables. Valid values: true, false. Default: false.

### Optional - BGP Configuration

* `enable_bgp` - (Optional) Enable BGP. Valid values: true, false. Default: false.
* `local_as_number` - (Optional) BGP local AS number. Required when `enable_bgp` is set to true.
* `prepend_as_path` - (Optional) List of AS numbers to prepend to the AS_Path field. Valid only when `local_as_number` is set. Example: ["65001", "65001"].
* `disable_route_propagation` - (Optional) Disable route propagation. Valid values: true, false. Default: false.
* `spoke_bgp_manual_advertise_cidrs` - (Optional) Set of intended CIDRs to be advertised to external BGP router. Example: ["10.1.0.0/16", "10.2.0.0/16"].
* `enable_preserve_as_path` - (Optional) Preserve AS path when advertising manual summary CIDRs. Valid values: true, false. Default: false.
* `enable_auto_advertise_s2c_cidrs` - (Optional) Auto advertise Site2Cloud CIDRs. Valid values: true, false. Default: false.
* `enable_bgp_ecmp` - (Optional) Enable BGP ECMP. Valid values: true, false. Default: false.

### Optional - BGP Timers

* `bgp_polling_time` - (Optional) BGP route polling time in seconds. Valid values: 10-50. Default: 50.
* `bgp_neighbor_status_polling_time` - (Optional) BGP neighbor status polling time in seconds. Valid values: 1-10. Default: 5.
* `bgp_hold_time` - (Optional) BGP hold time in seconds. Valid values: 12-360. Default: 180.

### Optional - BGP Communities

* `bgp_send_communities` - (Optional) Send BGP communities. Valid values: true, false. Default: false.
* `bgp_accept_communities` - (Optional) Accept BGP communities. Valid values: true, false. Default: false.

### Optional - BGP over LAN

* `enable_bgp_over_lan` - (Optional) Enable BGP over LAN. Valid values: true, false. Default: false.

### Optional - Learned CIDR Approval

* `enable_learned_cidrs_approval` - (Optional) Enable learned CIDR approval. Valid values: true, false. Default: false.
* `learned_cidrs_approval_mode` - (Optional) Learned CIDRs approval mode. Valid values: "gateway". Default: "gateway".
* `approved_learned_cidrs` - (Optional) Set of approved learned CIDRs. Valid only when `enable_learned_cidrs_approval` is true. Example: ["10.1.0.0/16", "10.2.0.0/16"].

### Optional - Active-Standby

* `enable_active_standby` - (Optional) Enable Active-Standby mode. Valid values: true, false. Default: false.
* `enable_active_standby_preemptive` - (Optional) Enable Active-Standby Preemptive mode. Valid only when `enable_active_standby` is true. Valid values: true, false. Default: false.

### Optional - AWS Specific

* `insane_mode` - (Optional) Enable Insane Mode (High Performance Encryption). Valid values: true, false. Default: false.

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

**spoke_group** can be imported using the `group_uuid`, e.g.

```
$ terraform import aviatrix_spoke_group.test group_uuid
```

## Notes

* The `group_name` is used to identify the gateway group and must be unique within the controller.
* BGP configuration options are only applicable when `enable_bgp` is set to true.
* The `bgp_polling_time`, `bgp_neighbor_status_polling_time`, and `bgp_hold_time` fields will default to their respective values (50, 5, and 180 seconds) when the API returns 0.
* The `learned_cidrs_approval_mode` will default to "gateway" when the API returns an empty string.
* Active-Standby Preemptive mode can only be enabled when Active-Standby mode is already enabled.
