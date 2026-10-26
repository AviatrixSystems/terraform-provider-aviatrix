---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_attachment_vrf_status"
description: |-
  Read and (optionally) toggle per-peering VRF attachment status.
---

# aviatrix_attachment_vrf_status

The **aviatrix_attachment_vrf_status** data source reads the per-peering `vrf_attachment_enabled` flag for eligible transit↔transit and edge-as-spoke↔transit peerings, and can optionally toggle it before reading.

Eligibility is enforced by the controller: both gateways involved in a peering must be at version 9.1 or later, and the topology must be transit↔transit or edge-as-spoke↔transit. The controller-wide VRF feature must also be enabled (see [`aviatrix_config_feature`](../resources/aviatrix_config_feature.md) with `feature_name = "vrf"`).

## Example Usage

### Read status for a single peering

```hcl
data "aviatrix_attachment_vrf_status" "pair" {
  gateway1 = "transit-us-east"
  gateway2 = "transit-us-west"
}

output "vrf_enabled" {
  value = data.aviatrix_attachment_vrf_status.pair.attachments[0].vrf_attachment_enabled
}
```

### Toggle and then read status

```hcl
data "aviatrix_attachment_vrf_status" "enable" {
  gateway1              = "transit-us-east"
  gateway2              = "transit-us-west"
  enable_vrf_attachment = "yes"
}
```

### List all eligible peerings on the controller

```hcl
data "aviatrix_attachment_vrf_status" "all" {
  all = true
}
```

### Toggle every eligible peering (bulk rollout)

~> **Warning:** This flips `vrf_attachment_enabled` on every eligible peering managed by the controller. Make sure that is what you want before applying.

```hcl
data "aviatrix_attachment_vrf_status" "bulk_enable" {
  all                   = true
  enable_vrf_attachment = "yes"
}
```

## Argument Reference

* `gateway1` - (Optional) Source gateway or gateway-group name. Required unless `all = true`.
* `gateway2` - (Optional) Destination gateway or gateway-group name. Required unless `all = true`.
* `all` - (Optional) When `true`, operate on every eligible peering. `gateway1` and `gateway2` must be empty in this mode. Defaults to `false`.
* `enable_vrf_attachment` - (Optional) Set to `"yes"` or `"no"` to toggle `vrf_attachment_enabled` before reading status. When empty (the default), the data source only reads the current status and does not mutate the controller.

## Attribute Reference

* `attachments` - List of per-peering VRF attachment records. Each entry contains:
  * `peering_name` - Controller-side peering identifier (`<gateway1>.<gateway2>`).
  * `gateway1` - First gateway name.
  * `gateway2` - Second gateway name.
  * `attachment_type` - `transit-transit` or `edge-spoke-transit`.
  * `vrf_attachment_enabled` - Current per-peering VRF attachment flag.

## Notes

* Removing this data source from configuration does **not** revert `vrf_attachment_enabled`. To turn the flag off, set `enable_vrf_attachment = "no"` (or call the controller API directly).
* Re-applying with the same value for `enable_vrf_attachment` is idempotent on the controller side, but the toggle call is still issued on every `terraform plan`/`apply` while the field is non-empty.
