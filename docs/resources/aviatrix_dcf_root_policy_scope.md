---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_root_policy_scope"
description: |-
  Manages the policy scopes applied to the DCF system root policy block
---

# aviatrix_dcf_root_policy_scope

The **aviatrix_dcf_root_policy_scope** resource manages the policy scopes applied to the Distributed Cloud Firewall (DCF) system root policy block. DCF must be enabled on the controller before using this resource.

This resource manages only the root block's scope configuration. It does not create, replace, or delete the system root policy block.

~> **NOTE:** This is a singleton resource. Declare at most one `aviatrix_dcf_root_policy_scope` per controller. See [Operational Notes](#operational-notes).

## Example Usage

```hcl
resource "aviatrix_dcf_policy_scope" "production" {
  name = "production-scope"

  included {
    csps = ["AWS"]
  }
}

resource "aviatrix_dcf_root_policy_scope" "root" {
  scope_uuids = [aviatrix_dcf_policy_scope.production.uuid]
}
```

## Scope Evaluation

The effective enforcement points of all `scope_uuids` are **unioned (ORed)**
to form the first level of every DCF scope chain. Scope levels on descendant
policy groups and rulesets narrow this root level by
being **intersected (ANDed)**. For example, root scopes `R1` and `R2` followed by
group scope `G1` produce `(R1 union R2) intersect G1`.

## Argument Reference

### Required

* `scope_uuids` - (Required) Unordered set of DCF policy scope UUIDs unioned to form the root scope level. At least one UUID is required. Every UUID must identify an existing policy scope.

## Operational Notes

* This resource is a singleton: the controller has exactly one system root policy block, so declare at most one `aviatrix_dcf_root_policy_scope` per controller, across all Terraform configurations and workspaces. The provider does not detect duplicates; multiple instances overwrite each other's scopes on every apply, and destroying any one of them resets the root to the Global Scope.
* Multiple UUIDs at the root level are unioned into one scope level.
* Updating `scope_uuids` patches only the root policy block's scope configuration.
* Destroying this resource resets the system root policy block to the controller-managed Global Scope (`def000ad-8000-0000-0000-000000000001`). The system root scope can never be empty.

## Import

**aviatrix_dcf_root_policy_scope** can be imported using its singleton ID:

```
$ terraform import aviatrix_dcf_root_policy_scope.root dcf_root_policy_scope
```
