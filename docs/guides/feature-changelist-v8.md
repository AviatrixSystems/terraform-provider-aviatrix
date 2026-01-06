---
layout: "aviatrix"
page_title: "Aviatrix R8.x Feature Changelist"
description: |-
  The Aviatrix provider R8.x Feature Changelist
---

# Aviatrix R8.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R8.0 to present. New resources may be tracked in the [Release Notes](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-notes).

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

---

``Last updated: R8.2.0 (8.2.0)``

---
## R8.2.0 (8.2.0)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

## R8.1.20 (8.1.20)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

## R8.0.40 (8.0.40)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

## R8.1.10 (8.1.10)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|


## R8.0.30 (8.0.30)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

## R8.1.1 (8.1.1)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|


## R8.1.0 (8.1.0)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|


## R8.0.10 (8.0.10)

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|

### Resource Deprecations

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
|-|-|

### Attribute Deprecations

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|-|-|-|-|


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
