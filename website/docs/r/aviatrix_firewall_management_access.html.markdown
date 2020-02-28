---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_management_access"
description: |-
  Creates and manages Aviatrix firewall management accesses
---

# aviatrix_firewall_management_access

The aviatrix_firewall_management_access resource allows the creation and management of Aviatrix firewall management accesses.

## Example Usage

```hcl
# Create an Aviatrix Firewall Management Access
resource "aviatrix_firewall_management_access" "test_firewall_management_access" {
  transit_firenet_gateway_name    = "transitGw1"
  management_access_resource_name = "SPOKE:spokeGw1"
}
```

## Argument Reference

The following arguments are supported:

* `transit_firenet_gateway_name` - (Required) Name of the transit gateway with transit firenet enabled. Currently supports AWS and AZURE providers.
* `management_access_resource_name` - (Required) Name of the resource to be enabled firewall management access.

## Import

Instance firewall_management_access can be imported using the transit_firenet_gateway_name, e.g.

```
$ terraform import aviatrix_firewall_management_access.test transit_firenet_gateway_name
```
