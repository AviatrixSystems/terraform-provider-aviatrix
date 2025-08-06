# Aviatrix Transit HA Gateway Resource

The `aviatrix_transit_ha_gateway` resource provides a dedicated way to manage High Availability (HA) for Aviatrix Transit Gateways. This resource allows you to create, configure, and manage HA gateways independently from the primary transit gateway.

## Key Features

- **Dedicated HA Management**: Manage HA gateways as separate resources for better organization and control
- **Multi-Cloud Support**: Support for AWS, GCP, Azure, OCI, AliCloud, and Edge (Equinix/Megaport) deployments
- **Insane Mode Support**: Full support for high-performance Insane Mode configurations
- **Edge Transit Support**: Specialized support for Edge transit gateways with custom interfaces and configurations
- **Comprehensive Configuration**: Support for all HA-related attributes including EIP, zones, availability domains, and more

## Resource Schema

### Required Arguments

- `primary_gw_name` - (Required) Name of the primary transit gateway for which HA is being enabled
- `cloud_type` - (Required) Type of cloud service provider (1=AWS, 4=GCP, 8=Azure, 16=OCI, etc.)
- `account_name` - (Required) Name of the Cloud-Account in Aviatrix controller
- `vpc_id` - (Required) VPC-ID/VNet-Name/Site-ID of cloud provider
- `ha_subnet` - (Required) HA Subnet CIDR for most cloud providers (optional for GCP)
- `ha_gw_size` - (Required) HA Gateway instance size

### Optional Arguments

#### Basic Configuration
- `ha_gw_name` - (Optional) Name of the HA gateway. Defaults to `<primary_gw_name>-hagw`
- `ha_zone` - (Optional) HA Zone. Required for GCP, optional for Azure
- `insane_mode` - (Optional) Enable Insane Mode for high performance. Default: false

#### AWS Specific
- `ha_insane_mode_az` - (Optional) Availability Zone for Insane Mode HA Gateway (required if insane_mode is enabled on AWS)
- `enable_encrypt_volume` - (Optional) Enable EBS volume encryption (AWS only)
- `customer_managed_keys` - (Optional) Customer managed key ID for encryption

#### Networking
- `ha_eip` - (Optional) Public IP address for the HA gateway
- `ha_oob_management_subnet` - (Optional) Out-of-band management subnet
- `ha_oob_availability_zone` - (Optional) Out-of-band availability zone

#### Azure Specific
- `ha_azure_eip_name_resource_group` - (Optional) Azure EIP name and resource group in format "name:resource_group"

#### OCI Specific
- `ha_availability_domain` - (Optional) HA availability domain for OCI
- `ha_fault_domain` - (Optional) HA fault domain for OCI

#### GCP Specific
- `ha_bgp_lan_interfaces` - (Optional) BGP LAN interfaces for GCP HA transit
  - `vpc_id` - (Required) VPC ID for the BGP LAN interface
  - `subnet` - (Required) Subnet CIDR for the BGP LAN interface

#### Edge Transit Specific
- `ha_device_id` - (Optional) Device ID for AEP EAT HA gateway
- `ztp_file_download_path` - (Optional) ZTP file download path for Edge transit gateways
- `interfaces` - (Optional) WAN/LAN/MANAGEMENT interfaces for Edge transit gateways
  - `ifname` - (Required) Interface name
  - `type` - (Required) Interface type (WAN, LAN, MANAGEMENT)
  - `bandwidth` - (Optional) Interface bandwidth in Mbps
  - `public_ip` - (Optional) Interface public IP
  - `tag` - (Optional) Interface tag
  - `dhcp` - (Optional) Enable DHCP on interface
  - `cidr` - (Optional) Interface CIDR
  - `gateway_ip` - (Optional) Interface gateway IP
- `interface_mapping` - (Optional) Interface mapping for ESXI devices
  - `name` - (Required) Interface name
  - `type` - (Required) Interface type (MANAGEMENT, WAN)
  - `index` - (Required) Interface index
- `management_egress_ip_prefix_list` - (Optional) Set of management egress IP prefixes

#### Private Mode
- `private_mode_subnet_zone` - (Optional) Private Mode HA subnet availability zone

#### Software Management
- `ha_software_version` - (Optional) Desired software version for the HA gateway
- `ha_image_version` - (Optional) Desired image version for the HA gateway

### Computed Attributes

- `ha_security_group_id` - Security group ID used for the HA transit gateway
- `ha_cloud_instance_id` - Cloud instance ID of the HA transit gateway
- `ha_private_ip` - Private IP address of the HA transit gateway

## Usage Examples

### Basic AWS HA Gateway

```hcl
resource "aviatrix_transit_gateway" "primary" {
  cloud_type   = 1
  account_name = "aws-account"
  gw_name      = "transit-gw"
  vpc_id       = "vpc-123456"
  vpc_reg      = "us-west-2"
  gw_size      = "t3.medium"
  subnet       = "10.0.1.0/24"
}

resource "aviatrix_transit_ha_gateway" "ha" {
  primary_gw_name = aviatrix_transit_gateway.primary.gw_name
  cloud_type      = 1
  account_name    = "aws-account"
  vpc_id          = "vpc-123456"
  ha_subnet       = "10.0.2.0/24"
  ha_gw_size      = "t3.medium"
}
```

### AWS with Insane Mode

```hcl
resource "aviatrix_transit_ha_gateway" "insane_mode" {
  primary_gw_name   = aviatrix_transit_gateway.primary.gw_name
  cloud_type        = 1
  account_name      = "aws-account"
  vpc_id            = "vpc-123456"
  ha_subnet         = "10.0.2.0/24"
  ha_gw_size        = "c5n.large"
  insane_mode       = true
  ha_insane_mode_az = "us-west-2b"
}
```

### GCP with BGP LAN Interfaces

```hcl
resource "aviatrix_transit_ha_gateway" "gcp" {
  primary_gw_name = aviatrix_transit_gateway.primary.gw_name
  cloud_type      = 4
  account_name    = "gcp-account"
  vpc_id          = "gcp-vpc"
  ha_subnet       = "10.0.2.0/24"
  ha_zone         = "us-west1-b"
  ha_gw_size      = "n1-standard-1"
  
  ha_bgp_lan_interfaces {
    vpc_id = "bgp-vpc"
    subnet = "172.16.1.0/24"
  }
}
```

### Azure with Custom EIP

```hcl
resource "aviatrix_transit_ha_gateway" "azure" {
  primary_gw_name                   = aviatrix_transit_gateway.primary.gw_name
  cloud_type                        = 8
  account_name                      = "azure-account"
  vpc_id                            = "azure-vnet:rg"
  ha_subnet                         = "10.0.2.0/24"
  ha_zone                           = "1"
  ha_gw_size                        = "Standard_B2s"
  ha_eip                            = "52.1.2.3"
  ha_azure_eip_name_resource_group = "eip-name:resource-group"
}
```

### Edge Transit (Equinix)

```hcl
resource "aviatrix_transit_ha_gateway" "edge" {
  primary_gw_name        = aviatrix_transit_gateway.primary.gw_name
  cloud_type             = 65536
  account_name           = "equinix-account"
  vpc_id                 = "site-123"
  ha_gw_size             = "Standard"
  ztp_file_download_path = "/tmp/ztp-files"
  
  interfaces {
    ifname = "eth0"
    type   = "MANAGEMENT"
    dhcp   = true
  }
  
  interfaces {
    ifname     = "eth1"
    type       = "WAN"
    bandwidth  = 1000
    public_ip  = "192.168.1.10"
    cidr       = "192.168.1.0/24"
    gateway_ip = "192.168.1.1"
  }
  
  management_egress_ip_prefix_list = ["10.0.0.0/8"]
}
```

## Import

Transit HA gateways can be imported using the HA gateway name:

```bash
terraform import aviatrix_transit_ha_gateway.example transit-gw-hagw
```

## Cloud Type Values

- `1` - AWS
- `4` - GCP
- `8` - Azure
- `16` - OCI
- `32` - Azure Gov
- `256` - AWS Gov
- `1024` - AWS China
- `2048` - Azure China
- `8192` - Alibaba Cloud
- `16384` - AWS Top Secret
- `32768` - AWS Secret
- `65536` - Edge Equinix
- `131072` - Edge Megaport
- `262144` - Edge NEO

## Validation Rules

### Cloud-Specific Requirements

1. **GCP**: `ha_zone` is required when enabling HA
2. **Azure**: `ha_subnet` is required when enabling HA
3. **AWS**: When `insane_mode` is enabled, `ha_insane_mode_az` is required
4. **OCI**: Both `ha_availability_domain` and `ha_fault_domain` are required
5. **Edge**: Specific interface configurations are required

### Size Requirements

- HA gateway size (`ha_gw_size`) is mandatory when creating HA
- For Insane Mode on AWS, minimum gateway size is c5 series
- For Insane Mode on Azure, minimum gateway size is Standard_D3_v2

## Benefits of Using This Resource

1. **Cleaner Resource Management**: Separate HA configuration from primary gateway
2. **Better State Management**: Independent lifecycle management for HA components
3. **Improved Modularity**: Easier to create reusable modules
4. **Enhanced Visibility**: Clear separation of primary and HA gateway configurations
5. **Simplified Troubleshooting**: Isolated management of HA-specific issues

## Migration from transit_gateway Resource

If you're currently using the `ha_subnet` parameter in the `aviatrix_transit_gateway` resource, you can migrate to this dedicated resource:

1. Remove HA-related parameters from your `aviatrix_transit_gateway` resource
2. Create a new `aviatrix_transit_ha_gateway` resource with the HA configuration
3. Run `terraform plan` to ensure the changes are correct
4. Apply the changes to migrate to the new structure

This approach provides better separation of concerns and more granular control over your transit gateway HA configuration.
