---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_version"
sidebar_current: "docs-aviatrix-resource-version"
description: |-
  Manages the Aviatrix Controller Version.
---

# aviatrix_version

The AviatrixVersion resource manages the controller upgrade process

## Example Usage

```hcl
# Manage Aviatrix Controller Upgrade process
resource "aviatrix_version" "test_version" {
  target_version = "latest"
}
```

```hcl
# Manage Aviatrix Controller Upgrade process
resource "aviatrix_version" "test_version" {
  target_version = "4.1"
}
```

## Argument Reference

The following arguments are supported:

* `target_version` - (Optional) The release version number to which the controller will be upgraded to. If not specified, controller will not be upgraded. If set to "latest", controller will be upgraded to the latest release. Please look at https://docs.aviatrix.com/HowTos/inline_upgrade.html for more information.


The following arguments are computed - please do not edit in the resource file:

* `version` - Current version of the controller.

## Import

Instance version can be imported using the target_version, e.g.

```
$ terraform import aviatrix_version.test target_version
```