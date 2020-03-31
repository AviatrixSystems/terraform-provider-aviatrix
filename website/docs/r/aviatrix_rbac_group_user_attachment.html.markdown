---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_user_attachment"
description: |-
  Creates and manages Aviatrix rbac group user attachments
---

# aviatrix_rbac_group_user_attachment

The **aviatrix_rbac_group_user_attachment** resource allows the creation and management of Aviatrix rbac group user attachments.

## Example Usage

```hcl
# Create an Aviatrix Rbac Group User Attachment
resource "aviatrix_rbac_group_user_attachment" "test_attachment" {
  group_name = "write_only"
  user_name  = "user_name"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) This parameter represents the name of a rbac group.
* `user_name` - (Required) User name.

## Import

**rbac_group_user_attachment** can be imported using the `group_name` and `user_name`, e.g.

```
$ terraform import aviatrix_rbac_group_user_attachment.test group_name~user_name
```