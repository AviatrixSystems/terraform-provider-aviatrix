---
subcategory: "Edge"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_equinix"
description: |-
  Creates Aviatrix Edge Equinix
---

# aviatrix_edge_equinix

The **aviatrix_edge_equinix** resource creates the Aviatrix Edge Equinix.

## Example Usage - Static management port

```hcl
# Create an Edge Equinix
resource "aviatrix_edge_equinix" "test" {
  account_name           = "edge_equinix-account"
  gw_name                = "equinix-test"
  site_id                = "site-123"
  ztp_file_download_path = "/ztp/file/download/path"

  interfaces {
    name          = "eth0"
    type          = "WAN"
    ip_address    = "10.230.5.32/24"
    gateway_ip    = "10.230.5.100"
    wan_public_ip = "64.71.24.221"
  }

  interfaces {
    name       = "eth1"
    type       = "LAN"
    ip_address = "10.230.3.32/24"
  }

  interfaces {
    name        = "eth2"
    type        = "MANAGEMENT"
    enable_dhcp = false
    ip_address  = "172.16.15.162/20"
    gateway_ip  = "172.16.0.1"
  }

  included_advertised_spoke_routes = [
    "10.10.0.0/16",
    "172.16.0.0/12"
  ]
}
```

## Example Usage - DHCP management port

```hcl
# Create an Edge Equinix
resource "aviatrix_edge_equinix" "test" {
  account_name           = "edge_equinix-account"
  gw_name                = "equinix-test"
  site_id                = "site-123"
  ztp_file_download_path = "/ztp/file/download/path"

  interfaces {
    name          = "eth0"
    type          = "WAN"
    ip_address    = "10.230.5.32/24"
    gateway_ip    = "10.230.5.100"
    wan_public_ip = "64.71.24.221"
  }

  interfaces {
    name       = "eth1"
    type       = "LAN"
    ip_address = "10.230.3.32/24"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Edge Equinix account name.
* `gw_name` - (Required) Edge Equinix name.
* `site_id` - (Required) Site ID.
* `ztp_file_download_path` - (Required) The folder path where the ZTP file will be downloaded.
* `interfaces` - (Required) WAN/LAN/MANAGEMENT interfaces.
  * `name` - (Required) Interface name.
  * `type` - (Required) Type. Valid values: WAN, LAN, or MANAGEMENT.
  * `bandwidth` - (Optional) The rate of data can be moved through the interface, requires an integer value. Unit is in Mb/s.
  * `enable_dhcp` - (Optional) Enable DHCP. Valid values: true, false. Default value: false.
  * `wan_public_ip` - (Optional) WAN public IP.
  * `ip_address` - (Optional) Interface static IP address.
  * `gateway_ip` - (Optional) Gateway IP.
  * `dns_server_ip` - (Optional) Primary DNS server IP.
  * `secondary_dns_server_ip` - (Optional) Secondary DNS server IP.
  * `enable_vrrp` - (Optional) Enable VRRP. Valid values: true, false. Default value: false.
  * `vrrp_virtual_ip` - (Optional) VRRP virtual IP.
  * `tag` - (Optional) Tag.

### Optional
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP and subnet prefix. Example: ["67.207.104.16/29", "64.71.12.144/29"]. This is required to open the security group of the controller, in order to allow communication with the Edge gateway. Should contain the public IP address of the Edge gateway management interface(s).
* `enable_management_over_private_network` - (Optional) Switch to enable management over the private network. Valid values: true, false. Default value: false.
* `enable_edge_active_standby` - (Optional) Switch to enable Edge Active-Standby mode. Valid values: true, false. Default value: false.
* `enable_edge_active_standby_preemptive` - (Optional) Switch to enable Preemptive Mode for Edge Active-Standby. Valid values: true, false. Default value: false.
* `dns_server_ip` - (Optional) DNS server IP. Required and valid when `management_interface_config` is "Static".
* `secondary_dns_server_ip` - (Optional) Secondary DNS server IP. Required and valid when `management_interface_config` is "Static".
* `local_as_number` - (Optional) BGP AS Number to assign to Edge Equinix.
* `prepend_as_path` - (Optional) List of AS numbers to prepend gateway BGP AS_Path field. Valid only when `local_as_number` is set. Example: ["65023", "65023"].
* `enable_learned_cidrs_approval` - (Optional) Switch to enable learned CIDR approval. Valid values: true, false. Default value: false.
* `approved_learned_cidrs` - (Optional) Set of approved learned CIDRs. Valid only when `enable_learned_cidrs_approval` is set to true. Example: ["10.1.0.0/16", "10.2.0.0/16"].
* `spoke_bgp_manual_advertise_cidrs` - (Optional) Set of intended CIDRs to be advertised to external BGP router. Example: ["10.1.0.0/16", "10.2.0.0/16"].
* `enable_preserve_as_path` - (Optional) Switch to enable preserve as_path when advertising manual summary CIDRs. Valid values: true, false. Default value: false.
* `bgp_polling_time` - (Optional) BGP route polling time in seconds. Valid values are between 10 and 50. Default value: 50.
* `bgp_neighbor_status_polling_time` - (Optional) BGP neighbor status polling time in seconds. Valid values are between 1 and 10. Default value: 5.
* `bgp_hold_time` - (Optional) BGP hold time in seconds. Valid values are between 12 and 360. Default value: 180.
* `enable_edge_transitive_routing` - (Optional) Switch to enable Edge transitive routing. Valid values: true, false. Default value: false.
* `enable_jumbo_frame` - (Optional) Switch to enable jumbo frame. Valid values: true, false. Default value: false.
* `latitude` - (Optional) Latitude of Edge Equinix. Valid values are between -90 and 90. Example: "47.7511".
* `longitude` - (Optional) Longitude of Edge Equinix. Valid values are between -180 and 180. Example: "120.7401".
* `rx_queue_size` - (Optional) Ethernet interface RX queue size. Once set, can't be deleted or disabled. Valid values: "1K", "2K", "4K".
* `vlan` - (Required) VLAN configuration.
  * `parent_interface_name` - (Required) Parent interface name.
  * `vlan_id` - (Required) VLAN ID.
  * `ip_address` - (Optional) LAN sub-interface IP address.
  * `gateway_ip` - (Optional) LAN sub-interface gateway IP.
  * `peer_ip_address` - (Optional) LAN sub-interface IP address on HA gateway.
  * `peer_gateway_ip` - (Optional) LAN sub-interface gateway IP on HA gateway.
  * `vrrp_virtual_ip` - (Optional) LAN sub-interface virtual IP.
  * `tag` - (Optional) Tag.
* `enable_single_ip_snat` - (Optional) Enable Single IP SNAT. Valid values: true, false. Default value: false.
* `enable_auto_advertise_lan_cidrs` - (Optional) Enable auto advertise LAN CIDRs. Valid values: true, false. Default value: true.
* `included_advertised_spoke_routes` - (Optional) A list of CIDRs to be advertised to on-prem gateways as Included CIDR List. When configured, it will replace all advertised routes from this VPC.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `state` - State of Edge Equinix.

## Deployment on Equinix Fabric
In order to deploy the Edge gateway on Equinix Fabric, you need to use the [`equinix_network_device`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/network_device)  and [`equinix_network_file`](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/network_file) resources. Critical argument values for these resource for deployment of the Edge gateway in Equinix Fabric are displayed in the tables below.

### equinix_network_device
| Key                | Value                                                                            |
|--------------------|----------------------------------------------------------------------------------|
| type_code          | "AVIATRIX_EDGE_10"                                                               |
| self_managed       | true                                                                             |
| byol               | true                                                                             |
| package_code       | "STD"                                                                            |
| version            | "7.1", "7.1.b", "7.1.c", "7.1.d" or "7.2", depending on your controller version. |
| cloud_init_file_id | Reference the `equinix_network_file` resource.                                   |

### equinix_network_file
| Key                | Value               |
|--------------------|---------------------|
| device_type_code   | "AVIATRIX_EDGE"     |
| process_type       | "CLOUD_INIT"        |
| self_managed       | true                |
| byol               | true                |

Make sure to use the generated cloud-init file (ztp file) for creation of the `equinix_network_file` resource and provide this to the `equinix_network_device` resource.
For a more extensive example of how to deploy Aviatrix Edge on Equinix, refer to this [Terraform module](https://github.com/terraform-aviatrix-modules/terraform-aviatrix-equinix-edge-spoke).

## Import

**edge_equinix** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_edge_equinix.test gw_name
```

## Deprecations
* Deprecated ``bandwidth`` in **WAN/LAN/MGMT interfaces**. This configuration value no longer has any effect. It will be removed from the Aviatrix provider in the 3.2.0 release.
* Deprecated ``dns_server_ip`` and ``secondary_dns_server_ip``. These configuration values have no effect and have been replaced with ``dns_server_ip`` and  ``secondary_dns_server_ip`` present in **WAN/LAN/MGMT interfaces**. It will be removed from the Aviatrix provider in the 8.1.0 release.
