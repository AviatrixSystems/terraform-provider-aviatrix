---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_transit_attachment"
description: |-
  Creates and manages Aviatrix Spoke-to-Transit attachments
---

# aviatrix_spoke_transit_attachment

The **aviatrix_spoke_transit_attachment** resource allows the creation and management of Aviatrix Spoke-to-Transit gateway attachments.

~> **NOTE:** Terraform currently provides both a standalone spoke-to-transit attachment resource and a spoke gateway with `transit_gw` attachment(s) defined in-line. At this time, you cannot use a spoke gateway resource with transit attachments defined in conjunction with the spoke-to-transit attachment resources. Doing so will cause a conflict of settings. In order to use this resource, please set `manage_transit_gateway_attachment` in the **aviatrix_spoke_gateway** to false.

~> **NOTE:** This resource should only be used to manage the primary gateway attachments. The HA gateway attachments will be handled automatically by the backend.

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

### Required
* `spoke_gw_name` - (Required) Name of the spoke gateway to attach to transit network.
* `transit_gw_name` - (Required) Name of the transit gateway to attach the spoke gateway to.

### Advanced Options
* `route_tables` - (Optional) Learned routes will be propagated to these route tables. Example: ["rtb-212ff547","rtb-04539787"].
* `enable_max_performance` - (Optional) Indicates whether the maximum amount of HPE tunnels will be created. Only valid when transit and spoke gateways are each launched in Insane Mode and in the same cloud type. Default value: true. Available as of provider version R2.22.3+.

## Import

**spoke_transit_attachment** can be imported using the `spoke_gw_name` and `transit_gw_name`, e.g.

```
$ terraform import aviatrix_spoke_transit_attachment.test spoke_gw_name~transit_gw_name
```
