---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_csp_ha"
description: |-
  Creates Aviatrix Edge CSP HA
---

# aviatrix_edge_csp_ha

The **aviatrix_edge_csp_ha** resource creates the Aviatrix Edge CSP HA.

## Example Usage

```hcl
# Create an Edge CSP HA
resource "aviatrix_edge_csp_ha" "test" {
  primary_gw_name   = "primary_edge_csp"
  compute_node_uuid = "abcde12345"

  interfaces {
    name       = "eth1"
    type       = "LAN"
    ip_address = "10.220.11.20/24"
    gateway_ip = "10.220.11.1"
  }

  interfaces {
    name        = "eth2"
    type        = "MANAGEMENT"
    enable_dhcp = true
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_gw_name` - (Required) Edge CSP name.
* `compute_node_uuid` - (Required) Edge CSP compute node UUID.

-> **NOTE:** At least one LAN interface is required.
* `interfaces` - (Required) WAN/LAN interfaces.
  * `name` - (Required) Interface name.
  * `type` - (Required) Type.
  * `bandwidth` - (Optional) Bandwidth.
  * `enable_dhcp` - (Optional) Enable DHCP. Valid values: true, false. Default value: false.
  * `wan_public_ip` - (Optional) WAN public IP.
  * `ip_address` - (Optional) Interface static IP address.
  * `gateway_ip` - (Optional) Gateway IP.
  * `dns_server_ip` - (Optional) Primary DNS server IP.
  * `secondary_dns_server_ip` - (Optional) Secondary DNS server IP.
  * `tag` - (Optional) Tag.
    
## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `account_name` - Edge CSP account name.

## Import

**edge_csp_ha** can be imported using the `primary_gw_name` in the form `primary_gw_name` + "_hagw" e.g.

```
$ terraform import aviatrix_edge_csp_ha.test primary_gw_name_hagw
```
