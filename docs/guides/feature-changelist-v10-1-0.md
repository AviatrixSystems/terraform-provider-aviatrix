---
layout: "aviatrix"
page_title: "Aviatrix R10.1.0 Feature Changelist"
description: |-
  The Aviatrix provider R10.1.0 Feature Changelist
---

# Aviatrix R10.1.0 Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) for the R10.1.0 release.

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and are planning to upgrade their Controller to 10.1.0 (and subsequently must also update their Aviatrix Terraform provider to version 10.1.0).


## R10.1.0

### Resource Renaming
| Diff | Resource | New Resource Name | Action Required? |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

The following resources are removed:

| Resource | Action Required? |
|:--------------:|:--------------------------:|
| aviatrix_controller_private_mode_config | **Yes**; Private Mode is fully deprecated in Controller 10.0.0+. Remove this resource from your Terraform configuration before upgrading. |
| aviatrix_private_mode_lb | **Yes**; Private Mode is fully deprecated in Controller 10.0.0+. Remove this resource from your Terraform configuration before upgrading. |
| aviatrix_private_mode_multicloud_endpoint | **Yes**; Private Mode is fully deprecated in Controller 10.0.0+. Remove this resource from your Terraform configuration before upgrading. |

### Attribute Deprecations

The following attributes are changed or removed:

| Diff | Resource | Attribute | Action Required? |
|:----:|----------------|:-----------------:|----------------------------|
| (default changed) | aviatrix_gateway | single_az_ha | **Yes**; the default changed from `false` to `true`. Configurations that do not set `single_az_ha` explicitly will now create gateways with Single-AZ HA enabled. Add `single_az_ha = false` to preserve the previous behavior. |
| (ForceNew) | aviatrix_edge_spoke_transit_attachment | insane_mode_tunnel_number | **Yes**; changing this attribute now triggers destroy and recreate. Plan for downtime if you need to update `insane_mode_tunnel_number` on an existing attachment. |
| (removed) | aviatrix_spoke_instance | customized_spoke_vpc_routes | **Yes**; this attribute has been removed from `aviatrix_spoke_instance`. Set `customized_spoke_vpc_routes` on `aviatrix_spoke_group` instead. Remove it from any `aviatrix_spoke_instance` blocks in your configuration. |
| (deprecated) | aviatrix_dcf_ruleset | watch | **Yes**; the `watch` attribute on DCF rules is deprecated. Replace it with the new `enforcement` attribute (`enforce`, `monitor`, or `disable`). The `watch` attribute will be removed in a future release. |
