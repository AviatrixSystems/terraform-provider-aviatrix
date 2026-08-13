---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_centralized_transit_firenet"
description: |-
  Creates and manages the centralized Transit FireNet
---

# aviatrix_centralized_transit_firenet

The **aviatrix_centralized_transit_firenet** resource allows the creation and management of the centralized Transit FireNet.

-> **NOTE:** Before creating a centralized Transit FireNet, please make sure both primary FireNet and secondary FireNet meet the required conditions. Please refer to this [link](https://docs.aviatrix.com/documentation/latest/firewall-and-security/firenet-centralized.html?expand=true) for more details.

## Example Usage

```hcl
# Create a Centralized Transit FireNet
resource "aviatrix_centralized_transit_firenet" "test" {
  primary_firenet_gw_name = "primary-transit"
  secondary_firenet_gw_name = "secondary-transit"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_firenet_gw_name` - (Required) Primary FireNet gateway name.
* `secondary_firenet_gw_name` - (Required) Secondary FireNet gateway name.

## Import

**aviatrix_centralized_transit_firenet** can be imported using the `primary_firenet_gw_name` and `secondary_firenet_gw_name`, e.g.

```
$ terraform import aviatrix_centralized_firenet.test primary-transit~secondary-transit
```
