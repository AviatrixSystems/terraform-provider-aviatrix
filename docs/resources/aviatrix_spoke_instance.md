---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_instance"
description: |-
  Creates and manages Aviatrix spoke instances within a gateway group
---

# aviatrix_spoke_instance

The **aviatrix_spoke_instance** resource allows the creation and management of Aviatrix spoke instances within a gateway group. Unlike the standalone `aviatrix_spoke_gateway` resource, spoke instances are created within the context of an existing `aviatrix_spoke_group`, inheriting common configuration like cloud type, account, VPC, and region from the group.

## Example Usage

### AWS Spoke Instance

```hcl
# First, create a spoke group
resource "aviatrix_spoke_group" "aws_spoke_group" {
  group_name          = "my-aws-spoke-group"
  cloud_type          = 1
  gw_type             = "SPOKE"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  account_name        = "my-aws-account"
  vpc_region          = "us-west-2"
}

# Create a spoke instance in the group
resource "aviatrix_spoke_instance" "aws_spoke" {
  group_uuid = aviatrix_spoke_group.aws_spoke_group.group_uuid
  gw_name    = "my-aws-spoke-instance"
  subnet     = "10.0.1.0/24"
  gw_size    = "t3.medium"
}
```

### AWS Spoke Instance with Insane Mode

```hcl
resource "aviatrix_spoke_instance" "aws_spoke_insane" {
  group_uuid     = aviatrix_spoke_group.aws_spoke_group.group_uuid
  gw_name        = "my-aws-spoke-insane"
  subnet         = "10.0.2.0/24"
  gw_size        = "c5.xlarge"
  insane_mode    = true
  insane_mode_az = "us-west-2a"
}
```

### Azure Spoke Instance

```hcl
resource "aviatrix_spoke_group" "azure_spoke_group" {
  group_name          = "my-azure-spoke-group"
  cloud_type          = 8
  gw_type             = "SPOKE"
  group_instance_size = "Standard_B2ms"
  vpc_id              = "vnet_name:rg_name:resource_guid"
  account_name        = "my-azure-account"
  vpc_region          = "West US 2"
}

resource "aviatrix_spoke_instance" "azure_spoke" {
  group_uuid = aviatrix_spoke_group.azure_spoke_group.group_uuid
  gw_name    = "my-azure-spoke-instance"
  subnet     = "10.1.0.0/24"
  gw_size    = "Standard_B2ms"
  zone       = "az-1"
}
```

### Azure Spoke Instance with BGP over LAN

```hcl
resource "aviatrix_spoke_instance" "azure_spoke_bgp" {
  group_uuid              = aviatrix_spoke_group.azure_spoke_group.group_uuid
  gw_name                 = "my-azure-spoke-bgp"
  subnet                  = "10.1.1.0/24"
  gw_size                 = "Standard_B2ms"
  zone                    = "az-2"
  enable_bgp_over_lan     = true
  bgp_lan_interfaces_count = 2
}
```

### OCI Spoke Instance

```hcl
resource "aviatrix_spoke_group" "oci_spoke_group" {
  group_name          = "my-oci-spoke-group"
  cloud_type          = 16
  gw_type             = "SPOKE"
  group_instance_size = "VM.Standard2.2"
  vpc_id              = "ocid1.vcn.oc1.iad.xxxx"
  account_name        = "my-oci-account"
  vpc_region          = "us-ashburn-1"
}

resource "aviatrix_spoke_instance" "oci_spoke" {
  group_uuid          = aviatrix_spoke_group.oci_spoke_group.group_uuid
  gw_name             = "my-oci-spoke-instance"
  subnet              = "10.2.0.0/24"
  gw_size             = "VM.Standard2.2"
  availability_domain = "AD-1"
  fault_domain        = "FD-1"
}
```

### Spoke Instance with Tags and Route Configuration

```hcl
resource "aviatrix_spoke_instance" "aws_spoke_full" {
  group_uuid                        = aviatrix_spoke_group.aws_spoke_group.group_uuid
  gw_name                           = "my-aws-spoke-full"
  subnet                            = "10.0.3.0/24"
  gw_size                           = "t3.medium"

  # Route configuration
  filtered_spoke_vpc_routes         = "10.10.0.0/16,10.20.0.0/16"
  included_advertised_spoke_routes  = "10.0.0.0/8"
  enable_private_vpc_default_route  = true
  enable_skip_public_route_table_update = true

  # Monitoring
  enable_monitor_gateway_subnets    = true
  monitor_exclude_list              = ["i-1234567890abcdef0"]

  # Tunnel settings
  tunnel_detection_time             = 60

  # Tags
  tags = {
    Environment = "production"
    Team        = "networking"
  }
}
```

### Edge Equinix Spoke Instance

```hcl
resource "aviatrix_spoke_group" "edge_equinix_spoke_group" {
  group_name   = "my-edge-equinix-spoke-group"
  cloud_type   = 524288  # EDGEEQUINIX
  gw_type      = "EDGESPOKE"
  vpc_id       = "edge-equinix-site-id"
  account_name = "my-edge-equinix-account"
}

resource "aviatrix_spoke_instance" "edge_equinix_spoke" {
  group_uuid             = aviatrix_spoke_group.edge_equinix_spoke_group.group_uuid
  gw_name                = "my-edge-equinix-spoke"
  ztp_file_download_path = "/path/to/ztp/files"

  interfaces {
    logical_ifname = "wan0"
    type           = "WAN"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "lan0"
    type           = "LAN"
    ip_address     = "10.230.6.1/24"
  }

  interfaces {
    logical_ifname = "mgmt0"
    type           = "MANAGEMENT"
    dhcp           = true
  }

  management_egress_ip_prefix_list = ["10.0.0.0/8"]
}
```

### Edge Megaport Spoke Instance

```hcl
resource "aviatrix_spoke_group" "edge_megaport_spoke_group" {
  group_name   = "my-edge-megaport-spoke-group"
  cloud_type   = 1048576  # EDGEMEGAPORT
  gw_type      = "EDGESPOKE"
  vpc_id       = "edge-megaport-site-id"
  account_name = "my-edge-megaport-account"
}

resource "aviatrix_spoke_instance" "edge_megaport_spoke" {
  group_uuid             = aviatrix_spoke_group.edge_megaport_spoke_group.group_uuid
  gw_name                = "my-edge-megaport-spoke"
  ztp_file_download_path = "/path/to/ztp/files"

  interfaces {
    logical_ifname = "wan0"
    type           = "WAN"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "lan0"
    type           = "LAN"
    ip_address     = "10.230.6.1/24"
  }

  interfaces {
    logical_ifname = "mgmt0"
    type           = "MANAGEMENT"
    dhcp           = true
  }
}
```

### Edge Self-Managed (ESXI/KVM) Spoke Instance

```hcl
resource "aviatrix_spoke_group" "edge_selfmanaged_spoke_group" {
  group_name   = "my-edge-selfmanaged-spoke-group"
  cloud_type   = 4096  # EDGESELFMANAGED
  gw_type      = "EDGESPOKE"
  vpc_id       = "edge-selfmanaged-site-id"
  account_name = "my-edge-selfmanaged-account"
}

resource "aviatrix_spoke_instance" "edge_selfmanaged_spoke" {
  group_uuid             = aviatrix_spoke_group.edge_selfmanaged_spoke_group.group_uuid
  gw_name                = "my-edge-selfmanaged-spoke"
  ztp_file_download_path = "/path/to/ztp/files"
  ztp_file_type          = "iso"

  interfaces {
    logical_ifname = "wan0"
    type           = "WAN"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "lan0"
    type           = "LAN"
    ip_address     = "10.230.6.1/24"
  }

  interfaces {
    logical_ifname = "mgmt0"
    type           = "MANAGEMENT"
    dhcp           = true
  }
}
```

### Edge Platform (AEP/NEO) Spoke Instance

```hcl
resource "aviatrix_spoke_group" "edge_neo_spoke_group" {
  group_name   = "my-edge-neo-spoke-group"
  cloud_type   = 262144  # EDGENEO
  gw_type      = "EDGESPOKE"
  vpc_id       = "edge-neo-site-id"
  account_name = "my-edge-neo-account"
}

resource "aviatrix_spoke_instance" "edge_neo_spoke" {
  group_uuid = aviatrix_spoke_group.edge_neo_spoke_group.group_uuid
  gw_name    = "my-edge-neo-spoke"
  device_id  = "device-12345"

  interfaces {
    logical_ifname = "wan0"
    type           = "WAN"
    ip_address     = "10.230.5.32/24"
    gateway_ip     = "10.230.5.1"
    public_ip      = "64.71.24.221"
  }

  interfaces {
    logical_ifname = "lan0"
    type           = "LAN"
    ip_address     = "10.230.6.1/24"
  }

  interfaces {
    logical_ifname = "mgmt0"
    type           = "MANAGEMENT"
    dhcp           = true
  }

  interface_mapping {
    name  = "eth0"
    type  = "WAN"
    index = 0
  }

  interface_mapping {
    name  = "eth1"
    type  = "LAN"
    index = 1
  }

  interface_mapping {
    name  = "eth2"
    type  = "MANAGEMENT"
    index = 2
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `group_uuid` - (Required) UUID of the gateway group this spoke instance belongs to. This links the instance to an existing `aviatrix_spoke_group`.

### Optional - Basic Configuration

* `gw_name` - (Optional) Name of the spoke gateway. If not specified, a name will be auto-generated.
* `subnet` - (Optional) Public Subnet CIDR for the gateway. Required for CSP spoke instances (AWS, Azure, GCP, OCI), not applicable for Edge.
* `gw_size` - (Optional) Size of the gateway instance. Required for CSP spoke instances, not applicable for Edge.
* `zone` - (Optional) Availability Zone. Only available for Azure (8), Azure GOV (32) and Azure China (2048). Must be in the form 'az-n', for example, 'az-2'.
* `allocate_new_eip` - (Optional) If false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Default: true.
* `eip` - (Optional) Elastic IP address. Required when `allocate_new_eip` is false.
* `single_az_ha` - (Optional) Enable single AZ HA for the spoke gateway. Default: true.
* `tags` - (Optional) A map of tags to assign to the spoke gateway.
* `tunnel_detection_time` - (Optional) The IPSec tunnel down detection time for the Spoke Gateway. Valid values: 20-600 seconds.
* `insane_mode` - (Optional) Enable Insane Mode for spoke gateway. Supported for AWS/AWSGov, GCP, Azure and OCI. Default: false.

### Optional - Private Mode

* `private_mode_lb_vpc_id` - (Optional) Private Mode controller load balancer VPC ID. Required when private mode is enabled for the Controller.
* `private_mode_subnet_zone` - (Optional) Subnet availability zone for Private Mode. Required when Private Mode is enabled on the Controller and cloud_type is AWS.

### Optional - Route Configuration

* `filtered_spoke_vpc_routes` - (Optional) A list of comma-separated CIDRs to be filtered from the spoke VPC route table.
* `included_advertised_spoke_routes` - (Optional) A list of comma-separated CIDRs to be advertised to on-prem as 'Included CIDR List'.
* `enable_private_vpc_default_route` - (Optional) Program default route in VPC private route table. Default: false.
* `enable_skip_public_route_table_update` - (Optional) Skip programming VPC public route table. Default: false.
* `enable_monitor_gateway_subnets` - (Optional) Enable monitor gateway subnet. Default: false.
* `monitor_exclude_list` - (Optional) A set of monitored instance IDs to exclude.

### Optional - Spot Instance (AWS and Azure)

* `enable_spot_instance` - (Optional) Enable spot instance. NOT supported for production deployment. Only valid for AWS and Azure.
* `spot_price` - (Optional) Price for spot instance. NOT supported for production deployment. Required when `enable_spot_instance` is true.
* `delete_spot` - (Optional) If set true, the spot instance will be deleted on eviction. Otherwise, the instance will be deallocated on eviction. Only supports Azure.

### Optional - BGP over LAN (Azure only)

* `enable_bgp_over_lan` - (Optional) Pre-allocate a network interface (eth4) for 'BGP over LAN' functionality. Only valid for Azure (8), AzureGov (32) or AzureChina (2048). Default: false.
* `bgp_lan_interfaces_count` - (Optional) Number of interfaces that will be created for BGP over LAN enabled Azure spoke. Valid values: 1-5.

### Optional - Encryption (AWS only)

* `enable_encrypt_volume` - (Optional) Enable EBS volume encryption for Gateway. Only supports AWS and AWSGov. Default: false.
* `customer_managed_keys` - (Optional) Customer managed key ID for EBS volume encryption.

### Optional - AWS Specific

* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Spoke Gateway. Required if `insane_mode` is enabled for AWS cloud.
* `insertion_gateway` - (Optional) Enable to create an insertion gateway. Only valid for AWS. Default: false.
* `insertion_gateway_az` - (Optional) AZ for insertion gateway. Required when `insertion_gateway` is enabled.
* `rx_queue_size` - (Optional) Gateway ethernet interface RX queue size. Valid values: "1K", "2K", "4K".

### Optional - Azure Specific

* `azure_eip_name_resource_group` - (Optional) The name of the public IP address and its resource group in Azure to assign to this Spoke Gateway.

### Optional - OCI Specific

* `availability_domain` - (Optional) Availability domain for OCI gateway. Required for OCI.
* `fault_domain` - (Optional) Fault domain for OCI gateway. Required for OCI.

### Optional - Edge Spoke Gateway (Equinix, Megaport, Self-managed, AEP/NEO)

* `interfaces` - (Optional) A set of WAN/LAN/MANAGEMENT interface configurations for edge spoke gateways. Each interface block supports:
  * `logical_ifname` - (Required) Logical interface name (e.g., wan0, wan1, lan0, mgmt0).
  * `type` - (Required) Interface type. Valid values: "WAN", "LAN", "MANAGEMENT".
  * `ip_address` - (Optional) Interface static IP address in CIDR format.
  * `gateway_ip` - (Optional) Gateway IP address for the interface.
  * `public_ip` - (Optional) WAN interface public IP address.
  * `dhcp` - (Optional) Enable DHCP for the interface. Default: false.
* `interface_mapping` - (Optional) Interface mapping for AEP/NEO edge gateways. Each block supports:
  * `name` - (Required) Physical interface name (e.g., eth0, eth1).
  * `type` - (Required) Interface type. Valid values: "WAN", "LAN", "MANAGEMENT".
  * `index` - (Required) Interface index.
* `ztp_file_download_path` - (Optional) The local path where the ZTP file will be stored. Required for Equinix, Megaport, and Self-managed edge gateways.
* `ztp_file_type` - (Optional) ZTP file type for Self-managed edge gateways. Valid values: "iso", "cloud-init".
* `device_id` - (Optional) Device ID for AEP/NEO edge gateways.
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP/prefix CIDRs for edge gateways.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `security_group_id` - Security group used for the spoke gateway.
* `cloud_instance_id` - Cloud instance ID of the spoke gateway.
* `private_ip` - Private IP address of the spoke gateway.
* `public_ip` - Public IP address of the spoke gateway.
* `eip` - Elastic IP address assigned to the spoke gateway.
* `azure_bgp_lan_ip_list` - List of available BGP LAN interface IPs for spoke external device connection creation. Only valid for Azure.
* `software_version` - Software version of the gateway.
* `image_version` - Image version of the gateway.

## Import

**aviatrix_spoke_instance** can be imported using the `gw_name`, e.g.

```shell
$ terraform import aviatrix_spoke_instance.test my-spoke-instance-name
```
