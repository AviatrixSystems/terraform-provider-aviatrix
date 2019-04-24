---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_config"
sidebar_current: "docs-aviatrix-resource-controller-config"
description: |-
  Creates and manages an Aviatrix controller config resource.
---

# aviatrix_controller_config

The Account resource allows the creation and management of an Aviatrix controller config.

## Example Usage

```hcl
# Create Aviatrix Controller Config
resource "aviatrix_controller_config" "test_controller_config" {
  http_access         = "disabled"
  fqdn_exception_rule = "disabled"
}
```

## Argument Reference

The following arguments are supported:

* `http_access` - (Optional) Switch for http access. Default: "disabled".
* `fqdn_exception_rule` - (Optional) A system-wide mode. Default: "enabled".

## Import

Instance controller_config can be imported using ControllerIP, e.g.

```
$ terraform import aviatrix_controller_config.test ControllerIP
```