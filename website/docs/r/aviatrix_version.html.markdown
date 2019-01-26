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
  version = "3.5"
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Optional) The release version number to which the controller will be upgraded to. If not specified, it automatically will be upgraded to the latest release. Please look at https://docs.aviatrix.com/HowTos/inline_upgrade.html for more information.
