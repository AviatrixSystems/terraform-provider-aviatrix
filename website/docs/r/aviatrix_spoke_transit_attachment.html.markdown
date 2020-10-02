---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_transit_attachment"
description: |-
  Creates and manages Aviatrix Spoke Transit attachments
---

# aviatrix_spoke_transit_attachment

The **aviatrix_spoke_transit_attachment** resource allows the creation and management of Aviatrix Spoke Transit attachments.

## Example Usage

```hcl
# Create an Aviatrix Spoke Transit Attachment
resource "aviatrix_spoke_transit_attachment" "test_attachment" {
  spoke_gw_name   = "spoke-gw"
  transit_gw_name = "transit-gw"
  route_tables    = [
    "rtb-737d540c",
    "rtb-626d045c"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `spoke_gw_name` - (Required) Name of the spoke gateway to attach to transit network.
* `transit_gw_name` - (Required) Name of the transit gateway to attach the spoke gateway to.
* `route_tables` - (Optional) Advanced option. Learned routes will be propagated to these route tables. Example: ["rtb-212ff547","rtb-04539787"].



## Import

**spoke_transit_attachment** can be imported using the `spoke_gw_name` and `transit_gw_name`, e.g.

```
$ terraform import aviatrix_spoke_transit_attachment.test spoke_gw_name~transit_gw_name
```
