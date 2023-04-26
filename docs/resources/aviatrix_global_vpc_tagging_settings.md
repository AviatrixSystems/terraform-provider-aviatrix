---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_global_vpc_tagging_settings"
description: |-
  Manages how to tag newly found instances for global VPC.
---

# aviatrix_global_vpc_tagging_settings

The **aviatrix_global_vpc_tagging_settings** resource manages how to tag newly found instances for global VPC.

-> **NOTE:** The default service state is "semi_automatic". Therefore, after the resource is destroyed, the service state will be reset to "semi_automatic".

## Example Usage

```hcl
# Enable Automatic Tagging
resource "aviatrix_global_vpc_tagging_settings" "test" {
    service_state = "automatic"
    enable_alert  = false
}
```

## Argument Reference

The following arguments are supported:

### Required
* `service_state` - (Required) Tagging service state. Valid values: "semi_automatic", "automatic", "disabled".
* `enable_alert` - (Required) Set to true to enable alert. Valid values: true, false.

## Import

**global_vpc_tagging_settings** can be imported using controller IP, e.g.

```
$ terraform import aviatrix_global_vpc_tagging_settings.test 10-11-12-13
```
