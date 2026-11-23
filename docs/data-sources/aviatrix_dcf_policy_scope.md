---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_policy_scope"
description: |-
  Gets details about a DCF policy scope.
---

# aviatrix_dcf_policy_scope

The **aviatrix_dcf_policy_scope** data source provides details about an existing Distributed Cloud Firewall (DCF) policy scope, such as a scope created in the UI or the controller-managed `Global Scope`. DCF must be enabled on the controller before using this data source.

## Example Usage

```hcl
# Look up the controller-managed Global Scope
data "aviatrix_dcf_policy_scope" "global" {
  name = "Global Scope"
}

# Look up a scope created in the UI
data "aviatrix_dcf_policy_scope" "production" {
  name = "production-scope"
}

resource "aviatrix_dcf_ruleset" "example" {
  name = "production-ruleset"

  scope_config {
    scope_uuids = [data.aviatrix_dcf_policy_scope.production.uuid]
  }

  # rules { ... }
}
```

## Argument Reference

* `name` - (Required) Name of the policy scope. Policy scope names are unique on the controller.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the policy scope. Use this value in `scope_config.scope_uuids` or `aviatrix_dcf_root_policy_scope.scope_uuids`.
* `included` - The selector block identifying policy enforcement points included in the scope.
* `excluded` - The selector block whose matches are subtracted from the included set.

Each selector block exports the unordered sets `csps`, `vpc_ids`, `csp_account_ids`, `region_names`, `gateway_group_names`, `gateway_group_ids`, and `preset`. See the [aviatrix_dcf_policy_scope](../resources/aviatrix_dcf_policy_scope.md) resource for how selectors are evaluated.
