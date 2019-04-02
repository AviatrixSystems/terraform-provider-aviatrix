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
  username     = "username1"
  account_name = "test-accountname"
  email        = "username1@testdomain.com"
  password     = "passwordforuser1-1234"
}
```

## Argument Reference

The following arguments are supported for creating user account:

* `username` - (Required) Name of account user to be created.
* `account_name` - (Required) Cloud account name of user to be created.
* `email` - (Required) Email of address of account user to be created.
* `password` - (Required) Login password for the account user to be created. If password is changed, current account will be destroyed and a new account will be created.

## Import

Instance account_user can be imported using the username (when doing import, needs to leave password argument blank), e.g.

```
$ terraform import aviatrix_account_user.test username
```