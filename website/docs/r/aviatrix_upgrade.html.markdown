---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_upgrade"
sidebar_current: "docs-aviatrix-resource-upgrade"
description: |-
  Manages the Aviatrix Controller Upgrade Process.
---

# aviatrix_upgrade

The AviatrixUpgrade resource manages the controller upgrade process

## Example Usage

```hcl
# Manage Aviatrix Controller Upgrade process
resource "aviatrix_upgrade" "test_upgrade" {
  version = "3.5"
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Optional) The release version number to which the controller will be upgraded to. If not specified, it automatically will be upgraded to the latest release. Please look at https://docs.aviatrix.com/HowTos/inline_upgrade.html for more information.
