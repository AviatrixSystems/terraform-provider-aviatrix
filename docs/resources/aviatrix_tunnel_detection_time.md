---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_tunnel_detection_time"
description: |-
Creates and manages Aviatrix IPSec Tunnel Down Detection Time for gateways
---

# aviatrix_tunnel_detection_time

The **aviatrix_tunnel_detection_time** resource allows the creation and management of the IPSec Tunnel Down Detection Time for gateways.

## Example Usage

```hcl
# Create an Aviatrix Tunnel Detection Time Resource
resource "aviatrix_tunnel_detection_time" "test_tunnel_detection_time" {
  aviatrix_entity = "avtx-gw1"
  detection_time = 120
}
```

## Argument Reference

The following arguments are supported:

### Required
* `detection_time` - (Required) Set this attribute to the IPSec tunnel down detection time in seconds. The minimum is 20 seconds. The maximum is 600 seconds.

### Optional
* `aviatrix_entity` - (Optional) Gateway name to change IPSec tunnel down detection time for. If empty or set to \"Controller\", this resource will update all gateways to share the same tunnel down detection time.

## Import

**tunnel_detection_time** can be imported using the `aviatrix_entity`, e.g.

```
$ terraform import aviatrix_tunnel_detection_time.test_tunnel_detection_time aviatrix_entity
```
