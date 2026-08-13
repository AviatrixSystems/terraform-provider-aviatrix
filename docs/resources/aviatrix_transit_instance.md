---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_instance"
description: |-
  Creates and manages Aviatrix transit instances within a gateway group
---

# aviatrix_transit_instance

The **aviatrix_transit_instance** resource allows the creation and management of Aviatrix transit instances within a gateway group. Unlike the standalone `aviatrix_transit_gateway` resource, transit instances are created within the context of an existing `aviatrix_transit_group`, inheriting common configuration like cloud type, account, VPC, and region from the group.

## Example Usage

### AWS Transit Instance

```hcl
# First, create a transit group
resource "aviatrix_transit_group" "aws_transit_group" {
  group_name          = "my-aws-transit-group"
  cloud_type          = 1
  gw_type             = "TRANSIT"
  group_instance_size = "c5.xlarge"
  vpc_id              = "vpc-abcd1234"
  account_name        = "my-aws-account"
  vpc_region          = "us-west-2"
}

# Create a transit instance in the group
resource "aviatrix_transit_instance" "aws_transit" {
  group_uuid = aviatrix_transit_group.aws_transit_group.group_uuid
  gw_name    = "my-aws-transit-instance"
  subnet     = "10.0.1.0/24"
  gw_size    = "c5.xlarge"
}
```

### Azure Transit Instance

```hcl
resource "aviatrix_transit_group" "azure_transit_group" {
  group_name          = "my-azure-transit-group"
  cloud_type          = 8
  gw_type             = "TRANSIT"
  group_instance_size = "Standard_B2ms"
  vpc_id              = "vnet_name:rg_name:resource_guid"
  account_name        = "my-azure-account"
  vpc_region          = "West US 2"
}

resource "aviatrix_transit_instance" "azure_transit" {
  group_uuid = aviatrix_transit_group.azure_transit_group.group_uuid
  gw_name    = "my-azure-transit-instance"
  subnet     = "10.1.0.0/24"
  gw_size    = "Standard_B2ms"
  zone       = "az-1"
}
```

### Azure Transit Instance with BGP over LAN

```hcl
resource "aviatrix_transit_instance" "azure_transit_bgp" {
  group_uuid               = aviatrix_transit_group.azure_transit_group.group_uuid
  gw_name                  = "my-azure-transit-bgp"
  subnet                   = "10.1.1.0/24"
  gw_size                  = "Standard_B2ms"
  zone                     = "az-2"
  enable_bgp_over_lan      = true
  bgp_lan_interfaces_count = 2
}
```

### OCI Transit Instance

```hcl
resource "aviatrix_transit_group" "oci_transit_group" {
  group_name          = "my-oci-transit-group"
  cloud_type          = 16
  gw_type             = "TRANSIT"
  group_instance_size = "VM.Standard2.2"
  vpc_id              = "ocid1.vcn.oc1.iad.xxxx"
  account_name        = "my-oci-account"
  vpc_region          = "us-ashburn-1"
}

resource "aviatrix_transit_instance" "oci_transit" {
  group_uuid          = aviatrix_transit_group.oci_transit_group.group_uuid
  gw_name             = "my-oci-transit-instance"
  subnet              = "10.2.0.0/24"
  gw_size             = "VM.Standard2.2"
  availability_domain = "AD-1"
  fault_domain        = "FD-1"
}
```

### Transit Instance with Tags and Route Configuration

```hcl
resource "aviatrix_transit_instance" "aws_transit_full" {
  group_uuid                       = aviatrix_transit_group.aws_transit_group.group_uuid
  gw_name                          = "my-aws-transit-full"
  subnet                           = "10.0.3.0/24"
  gw_size                          = "c5.xlarge"

  # Route configuration
  customized_spoke_vpc_routes      = "10.10.0.0/16,10.20.0.0/16"
  filtered_spoke_vpc_routes        = "10.30.0.0/16"
  excluded_advertised_spoke_routes = "10.40.0.0/16"

  # Monitoring
  enable_monitor_gateway_subnets   = true
  monitor_exclude_list             = ["i-1234567890abcdef0"]

  # Tunnel settings
  tunnel_detection_time            = 60

  # Tags
  tags = {
    Environment = "production"
    Team        = "networking"
  }
}
```

### Edge Equinix Transit Instance

```hcl
resource "aviatrix_transit_group" "edge_equinix_transit_group" {
  group_name   = "my-edge-equinix-transit-group"
  cloud_type   = 524288  # EDGEEQUINIX
  gw_type      = "EDGETRANSIT"
  vpc_id       = "edge-equinix-site-id"
  account_name = "my-edge-equinix-account"
}

resource "aviatrix_transit_instance" "edge_equinix_transit" {
  group_uuid             = aviatrix_transit_group.edge_equinix_transit_group.group_uuid
  gw_name                = "my-edge-equinix-transit"
  gw_size                = "UNKNOWN"
  ztp_file_download_path = "/path/to/ztp/files"

  interfaces {
    logical_ifname = "wan0"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "mgmt0"
    dhcp           = true
  }

  management_egress_ip_prefix_list = ["10.0.0.0/8"]
}
```

### Edge Megaport Transit Instance

```hcl
resource "aviatrix_transit_group" "edge_megaport_transit_group" {
  group_name   = "my-edge-megaport-transit-group"
  cloud_type   = 1048576  # EDGEMEGAPORT
  gw_type      = "EDGETRANSIT"
  vpc_id       = "edge-megaport-site-id"
  account_name = "my-edge-megaport-account"
}

resource "aviatrix_transit_instance" "edge_megaport_transit" {
  group_uuid             = aviatrix_transit_group.edge_megaport_transit_group.group_uuid
  gw_name                = "my-edge-megaport-transit"
  gw_size                = "UNKNOWN"
  ztp_file_download_path = "/path/to/ztp/files"

  interfaces {
    logical_ifname = "wan0"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "mgmt0"
    dhcp           = true
  }
}
```

### Edge Self-Managed (ESXI/KVM) Transit Instance

```hcl
resource "aviatrix_transit_group" "edge_selfmanaged_transit_group" {
  group_name   = "my-edge-selfmanaged-transit-group"
  cloud_type   = 4096  # EDGESELFMANAGED
  gw_type      = "EDGETRANSIT"
  vpc_id       = "edge-selfmanaged-site-id"
  account_name = "my-edge-selfmanaged-account"
}

resource "aviatrix_transit_instance" "edge_selfmanaged_transit" {
  group_uuid             = aviatrix_transit_group.edge_selfmanaged_transit_group.group_uuid
  gw_name                = "my-edge-selfmanaged-transit"
  gw_size                = "UNKNOWN"
  ztp_file_download_path = "/path/to/ztp/files"
  ztp_file_type          = "iso"

  interfaces {
    logical_ifname = "wan0"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "mgmt0"
    dhcp           = true
  }

  interface_mapping {
    name  = "eth0"
    type  = "WAN"
    index = 0
  }

  interface_mapping {
    name  = "eth1"
    type  = "MANAGEMENT"
    index = 1
  }
}
```

### Edge Platform (AEP/NEO) Transit Instance

```hcl
resource "aviatrix_transit_group" "edge_neo_transit_group" {
  group_name   = "my-edge-neo-transit-group"
  cloud_type   = 262144  # EDGENEO
  gw_type      = "EDGETRANSIT"
  vpc_id       = "edge-neo-site-id"
  account_name = "my-edge-neo-account"
}

resource "aviatrix_transit_instance" "edge_neo_transit" {
  group_uuid = aviatrix_transit_group.edge_neo_transit_group.group_uuid
  gw_name    = "my-edge-neo-transit"
  gw_size    = "UNKNOWN"
  device_id  = "device-12345"

  interfaces {
    logical_ifname = "wan0"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "mgmt0"
    dhcp           = true
  }

  interface_mapping {
    name  = "eth0"
    type  = "WAN"
    index = 0
  }

  interface_mapping {
    name  = "eth1"
    type  = "MANAGEMENT"
    index = 1
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `group_uuid` - (Required) UUID of the gateway group this transit instance belongs to. This links the instance to an existing `aviatrix_transit_group`.
* `gw_size` - (Required) Size of the gateway instance. For edge transit gateways, use "UNKNOWN".

### Optional - Basic Configuration

* `gw_name` - (Optional) Name of the transit gateway. If not specified, a name will be auto-generated.
* `subnet` - (Optional) Public Subnet CIDR for the gateway. Required for CSP transit instances (AWS, Azure, GCP, OCI), not applicable for Edge.
* `allocate_new_eip` - (Optional) If false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Default: true.
* `eip` - (Optional) Elastic IP address. Required when `allocate_new_eip` is false.
* `single_az_ha` - (Optional) Enable single AZ HA for the transit gateway. Default: true.
* `tags` - (Optional) A map of tags to assign to the transit gateway.
* `tunnel_detection_time` - (Optional) The IPSec tunnel down detection time for the Transit Gateway. Valid values: 20-600 seconds.

### Optional - Private Mode

* `private_mode_lb_vpc_id` - (Optional) Private Mode controller load balancer VPC ID. Required when private mode is enabled for the Controller.
* `private_mode_subnet_zone` - (Optional) Subnet availability zone for Private Mode.

### Optional - Route Configuration

* `customized_spoke_vpc_routes` - (Optional) A list of comma-separated CIDRs to be customized for the spoke VPC routes.
* `filtered_spoke_vpc_routes` - (Optional) A list of comma-separated CIDRs to be filtered from the spoke VPC route table.
* `excluded_advertised_spoke_routes` - (Optional) A list of comma-separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'.
* `customized_transit_vpc_routes` - (Optional) A set of CIDRs to be customized for the transit VPC routes.
* `bgp_manual_spoke_advertise_cidrs` - (Optional) Intended CIDR list to be advertised to external BGP router.

### Optional - Feature Flags

* `enable_transit_firenet` - (Optional) Enable transit firenet interfaces. Default: false.
* `enable_firenet` - (Optional) Enable firenet interfaces. Default: false.
* `lan_vpc_id` - (Optional) LAN VPC ID. Only used for GCP Transit FireNet.
* `lan_private_subnet` - (Optional) LAN Private Subnet. Only used for GCP Transit FireNet.
* `enable_gateway_load_balancer` - (Optional) Enable firenet interfaces with AWS Gateway Load Balancer. Default: false.
* `enable_bgp_over_lan` - (Optional) Pre-allocate a network interface for "BGP over LAN" functionality. Only valid for GCP and Azure. Default: false.
* `bgp_lan_interfaces_count` - (Optional) Number of interfaces for BGP over LAN enabled Azure transit.

### Optional - Spot Instance (AWS and Azure)

* `enable_spot_instance` - (Optional) Enable spot instance. NOT supported for production deployment.
* `spot_price` - (Optional) Price for spot instance. Required when `enable_spot_instance` is true.
* `delete_spot` - (Optional) If true, the spot instance will be deleted on eviction. Only supports Azure.

### Optional - AWS Specific

* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit Gateway. Required for AWS if insane_mode is enabled.
* `rx_queue_size` - (Optional) Gateway ethernet interface RX queue size. Valid values: "1K", "2K", "4K", "8K", "16K".
* `enable_monitor_gateway_subnets` - (Optional) Enable monitor gateway subnets. Default: false.
* `monitor_exclude_list` - (Optional) A set of monitored instance IDs to exclude.

### Optional - Azure Specific

* `zone` - (Optional) Availability Zone. Must be in the form 'az-n', for example, 'az-2'.
* `azure_eip_name_resource_group` - (Optional) The name of the public IP address and its resource group in Azure.

### Optional - OCI Specific

* `availability_domain` - (Optional) Availability domain for OCI gateway. Required for OCI.
* `fault_domain` - (Optional) Fault domain for OCI gateway. Required for OCI.

### Optional - Edge Transit Gateway (Equinix, Megaport, Self-managed, AEP/NEO)

* `interfaces` - (Optional) A set of WAN/Management interface configurations for edge transit gateways. Each interface block supports:
  * `logical_ifname` - (Required) Logical interface name (e.g., wan0, wan1, mgmt0).
  * `ip_address` - (Optional) Interface static IP address in CIDR format.
  * `gateway_ip` - (Optional) Gateway IP address for the interface.
  * `public_ip` - (Optional) WAN interface public IP address.
  * `dhcp` - (Optional) Enable DHCP for the interface.
  * `secondary_private_cidr_list` - (Optional) A list of secondary private CIDR blocks.
  * `underlay_cidr` - (Optional) The underlay CIDR for this interface.
* `interface_mapping` - (Optional) Interface mapping for Self-managed (ESXI) edge gateways. Each block supports:
  * `name` - (Required) Physical interface name (e.g., eth0, eth1).
  * `type` - (Required) Interface type. Valid values: "WAN", "MANAGEMENT".
  * `index` - (Required) Interface index.
* `ztp_file_download_path` - (Optional) The local path where the ZTP file will be stored. Required for Equinix, Megaport, and Self-managed edge gateways.
* `ztp_file_type` - (Optional) ZTP file type for Self-managed edge gateways. Valid values: "iso", "cloud-init".
* `device_id` - (Optional) Device ID for AEP/NEO edge gateways.
* `peer_connection_type` - (Optional) Connection type for the edge transit gateway. Valid values: "public", "private".
* `peer_backup_logical_ifname` - (Optional) Peer backup logical interface names.
* `eip_map` - (Optional) A list of mappings between interface names and their associated private and public IPs.
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP/prefix CIDRs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `group_name` - Name of the transit group.
* `cloud_type` - Type of cloud service provider.
* `account_name` - Name of the Cloud-Account in Aviatrix controller.
* `vpc_id` - VPC-ID/VNet-Name/Site-ID of cloud provider.
* `gateway_uuid` - UUID of the transit gateway.
* `security_group_id` - Security group used for the transit gateway.
* `cloud_instance_id` - Cloud instance ID of the transit gateway.
* `private_ip` - Private IP address of the transit gateway.
* `public_ip` - Public IP address of the transit gateway.
* `eip` - Elastic IP address assigned to the transit gateway.
* `lan_interface_cidr` - Transit gateway LAN interface CIDR.
* `bgp_lan_ip_list` - List of available BGP LAN interface IPs for GCP and Azure.
* `azure_bgp_lan_ip_list` - List of available BGP LAN interface IPs for Azure.
* `software_version` - Software version of the gateway.
* `image_version` - Image version of the gateway.

## Import

**aviatrix_transit_instance** can be imported using the `gw_name`, e.g.

```shell
$ terraform import aviatrix_transit_instance.test my-transit-instance-name
```
