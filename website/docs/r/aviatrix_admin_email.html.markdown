---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_admin_email"
sidebar_current: "docs-aviatrix-resource-admin_email"
description: |-
  Sets Aviatrix Admin Email Address.
---

# aviatrix_admin_email

The AdminEmail resource allows to set Aviatrix Admin Email Address.

## Example Usage

```hcl
# Set Aviatrix Admin Email Address
resource "aviatrix_admin_email" "test_adminemail" {
  admin_email = "testadmingemail@testaccount.com"
}
```

## Argument Reference

The following arguments are supported:

* `admin_email` - (Required) E-mail address of admin user to be set. Valid Value: Any valid e-mail address
