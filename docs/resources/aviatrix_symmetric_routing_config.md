---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_symmetric_routing_config"
description: |-
  Creates and manages the configuration of symmetric routing
---

# aviatrix_symmetric_routing_config

The **aviatrix_symmetric_routing_config** resource allows management of symmetric routing configuration.

## Example Usage

```hcl
# Create an Aviatrix Symmetric Routing Config to Enable Symmetric Routing
resource "aviatrix_symmetric_routing_config" "test" {
  gw_name                  = "gateway"
  enable_symmetric_routing = true
}
```


## Argument Reference

The following arguments are supported:

### Required
* `gw_name` - (Required) Gateway name.
* `enable_symmetric_routing` - (Required) Enable symmetric routing. Valid values: true, false.


## Import

**symmetric_routing_config** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_symmetric_routing_config.test gw_name
```
