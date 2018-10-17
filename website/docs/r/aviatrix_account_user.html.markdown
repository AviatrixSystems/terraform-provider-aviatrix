---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_account_user"
sidebar_current: "docs-aviatrix-resource-account-user"
description: |-
  Creates and manages Aviatrix User Accounts
---

# aviatrix_account_user

The AccountUser resource allows the creation and management of Aviatrix User Accounts.

## Example Usage

```hcl
# Create Aviatrix User Account
resource "aviatrix_account_user" "test_accountuser" {
  username = "username1"
  account_name = "test-accountname"
  email = "username1@testdomain.com"
  password = "passwordforuser1-1234"
}
```

## Argument Reference

The following arguments are supported for creating user account:

* `username` - (Required) Name of account user to be created.
* `account_name` - (Required) Cloud account name of user to be created.
* `email` - (Optional) Email of address of account user to be created.
* `password` - (Optional) Login password for the account user to be created.

The following arguments are supported for editing user account:

* `what` - (Optional) Type of change, indicating what info of user to be changed. Valid values: "account_name", "email", "password"
* `old_password` - (Optional) (Required, when what is "password")
* `new_password` - (Optional) (Required, when what is "password")
