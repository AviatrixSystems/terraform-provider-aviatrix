---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group"
description: |-
  Creates and manages Aviatrix rbac groups
---

# aviatrix_rbac_group

The **aviatrix_rbac_group** resource allows the creation and management of Aviatrix rbac groups.

## Example Usage

```hcl
# Create an Aviatrix Rbac Group
resource "aviatrix_rbac_group" "test_group" {
  group_name = "write_only"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) This parameter represents the name of a rbac group to be created.

## Import

**rbac_group** can be imported using the `group_name`, e.g.

```
$ terraform import aviatrix_rbac_group.test group_name
```
