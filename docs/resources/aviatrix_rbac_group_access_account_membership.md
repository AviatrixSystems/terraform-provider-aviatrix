---
subcategory: "Accounts"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_access_account_membership"
description: |-
Creates and manages Aviatrix RBAC group access-account membership
---
# aviatrix_rbac_group_access_account_membership

The **_aviatrix_rbac_group_access_account_membership** resource allows the creation and management of access-account membership for Aviatrix RBAC (Role-Based Access Control) groups. This resource is authoritative for a group’s access-account membership: it manages the complete set of access accounts that belong to the specified group and will add/remove memberships to match your configuration.

> **Note:** There is another related resource, [`aviatrix_rbac_group_access_account_attachment`](./aviatrix_rbac_group_access_account_attachment.md), which manages a single access-account-to-group attachment per resource. While that resource is still available, you should prefer **`aviatrix_rbac_group_access_account_membership`** when managing an entire group’s membership as a single source of truth, especially for larger sets of access accounts or groups.

## Example Usage

### Basic Usage

```hcl
# Manage access-account membership for group "ops_team"
resource "aviatrix_rbac_group_access_account_membership" "ops_membership" {
  group_name = "ops_team"

  access_account_names = [
    "aws-prod",
    "azure-shared",
    "gcp-analytics",
  ]

  # When this resource is destroyed, remove these accounts from the group.
  remove_access_accounts_on_destroy = true
}
```

### Advanced Usage with access-account References
```hcl
# Examples of onboarded accounts
resource "aviatrix_account" "aws_account" {
  account_name       = "aws1"
  cloud_type         = 1
  aws_account_number = var.aws_account_number
  aws_iam            = false
  aws_access_key     = var.aws_access_key
  aws_secret_key     = var.aws_secret_key
}

resource "aviatrix_account" "azure_account" {
  account_name        = "azure1"
  cloud_type          = 8
  arm_subscription_id = var.azure_subscription_id
  arm_directory_id    = var.azure_tenant_id
  arm_application_id  = var.azure_app_id
  arm_application_key = var.azure_app_secret
}

# RBAC group itself
resource "aviatrix_rbac_group" "ops" {
  group_name = "rbac_ops"
}

# Authoritatively manage the group's access-account set using references
resource "aviatrix_rbac_group_access_account_membership" "ops_membership" {
  group_name = aviatrix_rbac_group.ops.group_name

  access_account_names = [
    aviatrix_account.aws_account.account_name,
    aviatrix_account.azure_account.account_name,
  ]

  remove_access_accounts_on_destroy = true
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) The RBAC group name whose access-account membership you want to manage. This resource is authoritative for that group’s access-account set.
* `access_account_names` - (Required) Complete set of access-account names that must be members of the group. At least one name must be specified.

### Optional
* `remove_access_accounts_on_destroy` - (Optional) If true, deleting this resource will remove all listed access accounts from the group. Defaults to false (memberships are left as-is on destroy).

## Import

**rbac_group_access_account_membership** can be imported using the group name, for example:

```
$ terraform import aviatrix_rbac_group_access_account_membership.ops_membership rbac_ops
```
