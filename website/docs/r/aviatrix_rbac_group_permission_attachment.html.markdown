---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_permission_attachment"
description: |-
  Creates and manages Aviatrix rbac group permission attachments
---

# aviatrix_rbac_group_permission_attachment

The **aviatrix_rbac_group_permission_attachment** resource allows the creation and management of Aviatrix rbac group permission attachments.

## Example Usage

```hcl
# Create an Aviatrix Rbac Group Permission Attachment
resource "aviatrix_rbac_group_permission_attachment" "test_attachment" {
  group_name      = "write_only"
  permission_name = "all_write"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) This parameter represents the name of a rbac group.
* `permission_name` - (Required) This parameter represents the name of a Permission.

## Import

**rbac_group_permission_attachment** can be imported using the `group_name` and `permission_name`, e.g.

```
$ terraform import aviatrix_rbac_group_permission_attachment.test group_name~permission_name
```