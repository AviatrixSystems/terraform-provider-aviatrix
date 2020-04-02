---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_account_user"
description: |-
  Creates and manages Aviatrix user accounts
---

# aviatrix_account_user

The **aviatrix_account_user** resource allows the creation and management of Aviatrix user accounts.

~> **NOTE:** With the release of Controller 5.4 (compatible with Aviatrix provider R2.13), Role-Based Access Control (RBAC) is now integrated into the Accounts workflow. Any **aviatrix_account_user** created in 5.3 by default will have admin privileges (attached to the 'admin' RBAC permission group). In 5.4, any new account users created will no longer have the option to specify an `account_name`, but rather have the option to attach the user to specific RBAC groups through the **aviatrix_rbac_group_user_attachment** resource for more granular security control. Account users created in 5.4 will have minimal access (read_only) unless otherwise specified in the RBAC group permissions that the users are attached to.

## Example Usage

```hcl
# Create an Aviatrix User Account
resource "aviatrix_account_user" "test_accountuser" {
  username     = "username1"
  email        = "username1@testdomain.com"
  password     = "passwordforuser1-1234"
}
```

## Argument Reference

The following arguments are supported for creating user account:

### Required
* `username` - (Required) Name of account user to be created.
* `email` - (Required) Email of address of account user to be created.
* `password` - (Required) Login password for the account user to be created. If password is changed, current account will be destroyed and a new account will be created.

The following arguments are deprecated:

* `account_name` - (Required) Cloud account name of user to be created. Deprecated as of Aviatrix provider R2.13 (Controller 5.4) due to RBAC implementation.

-> **NOTE:** `account_name` - If you are using/upgraded to Aviatrix Terraform Provider R2.13+, and an **aviatrix_account_user** resource was originally created with a provider version <R2.13, you must remove this attribute and perform a 'terraform refresh' to rectify the state file.

## Import

**account_user** can be imported using the `username` (when doing import, need to leave `password` argument blank), e.g.

```
$ terraform import aviatrix_account_user.test username
```
