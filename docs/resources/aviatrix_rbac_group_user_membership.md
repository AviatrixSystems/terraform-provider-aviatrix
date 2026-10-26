---
subcategory: "Accounts"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_user_membership"
description: |-
  Creates and manages Aviatrix RBAC group user membership
---
# aviatrix_rbac_group_user_membership
The **aviatrix_rbac_group_user_membership** resource allows the creation and management of user membership for Aviatrix (Role-Based Access Control) RBAC groups. This resource is authoritative for the group's user membership, meaning it manages the complete set of users that belong to a specific group.

> **Note:** There is another related resource, [`aviatrix_rbac_group_user_attachment`](./aviatrix_rbac_group_user_attachment.md), which manages a single user-to-group attachment per resource. While that resource is still available, you should prefer **`aviatrix_rbac_group_user_membership`** when managing an entire group’s membership as a single source of truth, especially for larger sets of users or groups.

## Example Usage

### Basic Usage
```hcl
# Create an Aviatrix RBAC Group User Membership
resource "aviatrix_rbac_group_user_membership" "test_membership" {
  group_name = "write_only"
  user_names = [
    "user1",
    "user2",
    "admin_user"
  ]
  remove_users_on_destroy = true
}
```

### Advanced Usage with User References
```hcl
# Create Aviatrix account users
resource "aviatrix_account_user" "rbac_user_1" {
  username = "rbac_user_1"
  email    = "rbac_user_1@aviatrix.com"
  password = "Rbac_user1"
}

resource "aviatrix_account_user" "rbac_user_2" {
  username = "rbac_user_2"
  email    = "rbac_user_2@aviatrix.com"
  password = "Rbac_user2"
}

resource "aviatrix_account_user" "rbac_user_3" {
  username = "rbac_user_3"
  email    = "rbac_user_3@aviatrix.com"
  password = "Rbac_user3"
}

# Create RBAC group membership with referenced users
resource "aviatrix_rbac_group_user_membership" "rbac_grp_membership1" {
  group_name = "rbac_local_login_1"
  user_names = [
    aviatrix_account_user.rbac_user_1.username,
    aviatrix_account_user.rbac_user_2.username,
    aviatrix_account_user.rbac_user_3.username,
  ]
  remove_users_on_destroy = true
}
```

## Argument Reference
The following arguments are supported:

### Required
* `group_name` - (Required) RBAC permission group name. This resource is authoritative for the group's user membership.
* `user_names` - (Required) Complete set of user names that must be members of the group (authoritative). At least one user name must be specified.

### Optional
* `remove_users_on_destroy` - (Optional) If true, deleting this resource will remove all users from the group. Default is false (the users are left in place).

## Import
**rbac_group_user_membership** can be imported using the `group_name`, e.g.

```
$ terraform import aviatrix_rbac_group_user_membership.test write_only
```
