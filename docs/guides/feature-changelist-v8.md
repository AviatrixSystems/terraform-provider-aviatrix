---
layout: "aviatrix"
page_title: "Aviatrix R8.x Feature Changelist"
description: |-
  The Aviatrix provider R8.x Feature Changelist
---

# Aviatrix R8.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R8.0 to present. Only releases with customer-impacting changes are listed. Releases without such changes are omitted.

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

## R8.2.10

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
|-|-|-|-|

---

## R8.0.0 (8.0.0-1000.2432)
**NOTICE:** Starting with this release, the Terraform provider will synchronize its version with the Aviatrix Controller version. This means the provider version has jumped from v3.2.2 to v8.0.0 to align with the Controller’s latest major version. This change makes it easier to determine which provider version is compatible with which Controller version.
Moving forward, the provider will follow semantic versioning (major.minor.patch).

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

The following resources are removed:

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
aviatrix_cloudn_registration | CloudN is no longer supported. Make sure to migrate to Edge before upgrading. |
aviatrix_cloudn_transit_gateway_attachment | CloudN is no longer supported. Make sure to migrate to Edge before upgrading. |

### Attribute Deprecations

The following attributes are removed:

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated)|aviatrix_edge_csp|aviatrix_dns_profile|**Yes**; please remove this attribute from the config.|
|(deprecated)|aviatrix_edge_equinix|aviatrix_dns_profile|**Yes**; please remove this attribute from the config.|
|(deprecated)|aviatrix_edge_equinix|aviatrix_dns_profile|**Yes**; please remove this attribute from the config.|
|(deprecated)|aviatrix_edge_platform|aviatrix_dns_profile|**Yes**; please remove this attribute from the config.|
|(deprecated)|aviatrix_edge_zededa|aviatrix_dns_profile|**Yes**; please remove this attribute from the config.|
