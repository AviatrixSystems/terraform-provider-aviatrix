---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_access_account_attachment"
description: |-
  Creates and manages Aviatrix rbac group access account attachments
---

# aviatrix_rbac_group_access_account_attachment

The **aviatrix_rbac_group_access_account_attachment** resource allows the creation and management of Aviatrix rbac group access account attachments.

## Example Usage

```hcl
# Create an Aviatrix Rbac Group Access Account Attachment
resource "aviatrix_rbac_group_access_account_attachment" "test_attachment" {
  group_name          = "write_only"
  access_account_name = "account_name"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) This parameter represents the name of a rbac group.
* `access_account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Import

**rbac_group_access_account_attachment** can be imported using the `group_name` and `access_account_name`, e.g.

```
$ terraform import aviatrix_rbac_group_access_account_attachment.test group_name~access_account_name
```