---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_access_account_attachment"
description: |-
  Creates and manages Aviatrix RBAC group access account attachments
---

# aviatrix_rbac_group_access_account_attachment

The **aviatrix_rbac_group_access_account_attachment** resource allows the creation and management of access account attachments to Aviatrix (Role-Based Access Control) RBAC groups.

## Example Usage

```hcl
# Create an Aviatrix RBAC Group Access Account Attachment
resource "aviatrix_rbac_group_access_account_attachment" "test_attachment" {
  group_name          = "write_only"
  access_account_name = "account_name"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) This parameter represents the name of a RBAC group.
* `access_account_name` - (Required) Account name. This can be used for logging in to CloudN console or UserConnect controller.

-> **NOTE:** If "all" is specified as the value for `access_account_name`, all existing access accounts will be attached to the specified RBAC group. If "all" is set, there is no need to specify any more access accounts attachments for that RBAC group.

## Import

**rbac_group_access_account_attachment** can be imported using the `group_name` and `access_account_name`, e.g.

```
$ terraform import aviatrix_rbac_group_access_account_attachment.test group_name~access_account_name
```
