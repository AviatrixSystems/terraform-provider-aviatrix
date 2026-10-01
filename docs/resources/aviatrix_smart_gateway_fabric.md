---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_smart_gateway_fabric"
description: |-
  Turns the Smart Gateway underlay mesh on or off across every eligible gateway group.
---

# aviatrix_smart_gateway_fabric

The **aviatrix_smart_gateway_fabric** resource turns the Smart Gateway
underlay mesh on or off across every eligible gateway group.

Its `feature_name` and `is_enabled` arguments mirror
[`aviatrix_config_feature`](aviatrix_config_feature.md), which is where the
`smart_gateway` controller feature itself is enabled. The route resolver is
set per gateway group with `enable_route_resolver`, which is available on
both resource models:

- [`aviatrix_transit_gateway`](aviatrix_transit_gateway.md) and
  [`aviatrix_spoke_gateway`](aviatrix_spoke_gateway.md), where it applies to
  the gateway's group.
- [`aviatrix_transit_group`](aviatrix_transit_group.md) and
  [`aviatrix_spoke_group`](aviatrix_spoke_group.md), where it sits on the
  group. The `aviatrix_transit_instance` and `aviatrix_spoke_instance`
  resources in the group do not have it.

## Example Usage

```hcl
resource "aviatrix_config_feature" "smart_gateway" {
  feature_name = "smart_gateway"
  is_enabled   = true
}

resource "aviatrix_smart_gateway_fabric" "underlay_mesh" {
  feature_name = "underlay_mesh"
  is_enabled   = true

  depends_on = [aviatrix_config_feature.smart_gateway]
}

resource "aviatrix_transit_gateway" "transit" {
  # ...
  enable_route_resolver = true

  depends_on = [aviatrix_smart_gateway_fabric.underlay_mesh]
}

resource "aviatrix_transit_group" "transit_group" {
  # ...
  enable_route_resolver = true

  depends_on = [aviatrix_smart_gateway_fabric.underlay_mesh]
}
```

## Argument Reference

* `feature_name` - (Required, ForceNew) Which Smart Gateway fabric feature
  this instance owns. Must be `underlay_mesh`.
* `is_enabled` - (Required) Turn the underlay mesh on or off across every
  eligible gateway group. Requires the `smart_gateway` controller feature.

## Attribute Reference

* `id` - The value of `feature_name`.

## Import

Import by feature name:

```
$ terraform import aviatrix_smart_gateway_fabric.underlay_mesh underlay_mesh
```

## Notes

- Eligible groups are spoke, transit, edge-spoke and edge-transit groups.
  Standalone, VPN and non-BGP cloud spoke groups are skipped.
- Only one instance per controller. A second one would compete with the
  first on every apply.
- The controller refuses two ordering breaks, so the provider does too:
  - Turning `underlay_mesh` off is refused while any group has the route
    resolver on. Set `enable_route_resolver = false` on those gateways first.
  - Creating a gateway or group with `enable_route_resolver = true` is
    refused while any group has the underlay off.
- Rollback takes one of two paths:
  - Two applies: set `enable_route_resolver = false` on every gateway and
    apply, then set `underlay_mesh` `is_enabled = false` and apply.
    Terraform updates `underlay_mesh` before the gateways that depend on it,
    so doing both in one apply is refused.
  - One apply: set `is_enabled = false` on the `aviatrix_config_feature` for
    `smart_gateway`, on `underlay_mesh`, and `enable_route_resolver = false`
    on every gateway. The controller turns every per-group flag off with the
    feature, so the rest are no-ops.
- Turning off only the `smart_gateway` feature turns the underlay and the
  route resolver off in every group. Configuration that still says `true`
  then shows a diff, and the apply is refused while the feature is off.
- A gateway or group created without `enable_route_resolver` gets the route
  resolver only if every other eligible group already has it on; otherwise
  it is off. Without the attribute, Terraform keeps whatever the controller
  has, so after the feature is turned off and on again such groups stay off.
  Set `enable_route_resolver` explicitly on every gateway that must have it.
- A group created with `enable_route_resolver = true` has no gateway yet, so
  the controller cannot set the route resolver in that apply. The apply shows
  a warning, and the next apply sets it once the group has an instance.
- A group changed outside Terraform so that only part of the fabric has the
  underlay on shows as a diff on the next plan, and apply sweeps it back.
- Destroying this resource only removes it from state. Controller state is
  left as it is.
