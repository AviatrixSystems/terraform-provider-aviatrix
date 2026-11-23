---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_policy_scope"
description: |-
  Creates and manages an Aviatrix DCF policy scope
---

# aviatrix_dcf_policy_scope

The **aviatrix_dcf_policy_scope** resource creates and manages a Distributed Cloud Firewall (DCF) policy scope. A policy scope selects the policy enforcement points where DCF policies apply. DCF must be enabled on the controller before using this resource.

## Example Usage

```hcl
resource "aviatrix_dcf_policy_scope" "production" {
  name = "production-scope"

  included {
    csps         = ["AWS", "AZURE"]
    region_names = ["us-east-1"]
  }

  excluded {
    gateway_group_names = ["maintenance"]
  }
}
```

To include every policy enforcement point without exclusions:

```hcl
resource "aviatrix_dcf_policy_scope" "all" {
  name = "all-gateways"

  included {
    preset = ["all"]
  }
}
```

## Scope Evaluation

Every value across all populated selectors in a block is **unioned (ORed)**.
Selector fields are not intersected with one another. The controller then
subtracts the union of the `excluded` selectors from the union of the
`included` selectors:

`effective scope = union(included selectors) - union(excluded selectors)`

For example, the first configuration includes every enforcement point in AWS,
plus every enforcement point in `us-east-1`, regardless of cloud provider. It
then removes enforcement points in the `maintenance` gateway group. It does
not mean "AWS enforcement points that are also in `us-east-1`."

## Argument Reference

### Required

* `name` - (Required) Name of the policy scope. Must contain 1–256 characters.
* `included` - (Required) A single selector block identifying policy enforcement points included in the scope. At least one selector must be configured. Matches from every configured selector are unioned.

### Optional

* `excluded` - (Optional) A single selector block whose unioned matches are subtracted from the included set. If omitted, the controller supplies `preset = ["none"]`, which excludes no policy enforcement points. Removing a previously configured `excluded` block resets it to `preset = ["none"]`.

The `included` and `excluded` blocks support these unordered sets:

* `csps` - Cloud service providers. Valid values (case-sensitive): `AWS`, `AWS-GOV`, `AWS-CHINA`, `AZURE`, `AZURE-GOV`, `AZURE-CHINA`, `GCP`, `OCI`, `ALICLOUD`.
* `vpc_ids` - VPC identifiers.
* `csp_account_ids` - Cloud account identifiers.
* `region_names` - Cloud region names.
* `gateway_group_names` - Gateway group names.
* `gateway_group_ids` - Gateway group identifiers.
* `preset` - Exactly one of `all` or `none`.

All values and selector fields within a block are unioned. `preset = ["all"]` selects all policy enforcement points, while `preset = ["none"]` selects none.

### Computed

* `uuid` - UUID of the policy scope. Use this value when another configuration requires a policy-scope UUID.

## Operational Notes

* The name and selector blocks can be updated without replacing the Terraform resource.
* A policy scope referenced by a policy group or ruleset cannot be deleted until those references are removed.
* The controller-managed `Global Scope` cannot be updated or deleted and should not be managed with this resource.

## Import

**aviatrix_dcf_policy_scope** can be imported using its UUID:

```
$ terraform import aviatrix_dcf_policy_scope.example 123e4567-e89b-12d3-a456-426614174000
```
