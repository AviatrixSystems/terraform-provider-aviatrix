---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_management_access"
description: |-
  Manage Aviatrix Firewall Management Access
---

# aviatrix_firewall_management_access

The **aviatrix_firewall_management_access** resource allows the management of which resource to permit visibility into the Transit (FireNet) VPC.

## Example Usage

```hcl
# Create an Aviatrix Firewall Management Access
resource "aviatrix_firewall_management_access" "test_firewall_management_access" {
  transit_firenet_gateway_name    = "transit-gw"
  management_access_resource_name = "SPOKE:spoke-gw"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `transit_firenet_gateway_name` - (Required) Name of the Transit FireNet-enabled transit gateway. Currently supports AWS(1) and Azure(8) providers.
* `management_access_resource_name` - (Required) Name of the resource to enable Firewall Management Access.

## Import

**firewall_management_access** can be imported using the `transit_firenet_gateway_name`, e.g.

```
$ terraform import aviatrix_firewall_management_access.test transit_firenet_gateway_name
```
