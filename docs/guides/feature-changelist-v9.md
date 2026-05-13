---
layout: "aviatrix"
page_title: "Aviatrix R9.x Feature Changelist"
description: |-
  The Aviatrix provider R9.x Feature Changelist
---

# Aviatrix R9.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R9.0 to present. Only releases with customer-impacting changes are listed. Releases without such changes are omitted. New resources may be tracked in the [Release Notes](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-notes).

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and are planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

---

## R9.0.0 (9.0.0)

### Resource Renaming
| Diff | Resource | New Resource Name | Action Required? |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource | Action Required? |
|:--------------:|:--------------------------:|
| aviatrix_gateway_certificate_config | **Yes**; this resource has been removed. Remove it from your config. |
| aviatrix_dcf_ips_profile_vpc | **Yes**; this resource has been removed. Migrate to **aviatrix_dcf_ips_default_profile** for global IPS Profile configuration. |
| aviatrix_distributed_firewalling_proxy_ca_config | **Yes**; migrate to **aviatrix_dcf_mitm_ca** and **aviatrix_dcf_mitm_ca_selection**. |
| aviatrix_distributed_firewalling_origin_cert_enforcement_config | **Yes**; migrate to **aviatrix_dcf_tls_profile**. |

### Attribute Deprecations

| Diff | Resource | Attribute | Action Required? |
|:----:|----------------|:-----------------:|----------------------------|
| (ForceNew) | aviatrix_spoke_external_device_conn | enable_ipv6 | **Yes**; if `enable_ipv6` was previously set to `true` on an existing connection, changing this value will now destroy and recreate the resource. IPv6 cannot be toggled in-place. |
| (ForceNew) | aviatrix_transit_external_device_conn | enable_ipv6 | **Yes**; same as above — in-place IPv6 edits are not supported. |
| (removed default) | aviatrix_vpc | ipv6_access_type | **Yes**; upgrading from 8.2 to 9.0 may cause forced replacement of `aviatrix_vpc` resources on non-GCP clouds if `ipv6_access_type` is not explicitly set in your config. Add `ipv6_access_type = ""` to prevent replacement. |
