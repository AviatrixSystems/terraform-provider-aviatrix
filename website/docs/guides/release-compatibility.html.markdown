---
layout: "aviatrix"
page_title: "Guides: Aviatrix Provider Release Compatibility"
description: |-
  The Aviatrix provider Release Compatibility Guide
---

# Aviatrix Provider: Release Compatibility Chart

## USAGE:
Quick at-a-glance access to Aviatrix Terraform provider's release compatibility with the Controller release versions. New resources and features may be tracked in the Release Notes.

-> **NOTE:** This only provides a quick glance at version compatibility between platforms. Please check the release notes for full details on new features, changes and deprecations [here](https://www.terraform.io/docs/providers/aviatrix/guides/release-notes.html).


---

``Last updated: R2.5 (UserConnect-5.1)``


---


| Terraform Version (v) | Aviatrix Provider Version (R) | Supported Controller Version |
|:-----------------:|:-------------------------:|:----------------------------:|
| 0.11              | 1.0                      | UserConnect-4.0              |
| 0.11              | 1.1                      | UserConnect-4.1              |
| 0.11              | 1.2                      | UserConnect-4.1              |
| 0.11              | 1.3                      | UserConnect-4.2              |
| 0.11              | 1.4                      | UserConnect-4.2              |
| 0.11              | 1.5                      | UserConnect-4.2              |
| 0.11              | 1.6                      | UserConnect-4.2              |
| 0.11              | 1.7                      | UserConnect-4.3              |
| 0.11              | 1.8                      | UserConnect-4.3              |
| 0.11              | 1.9                      | UserConnect-4.6              |
| 0.11              | 1.10                     | UserConnect-4.6              |
| 0.11              | 1.11                     | UserConnect-4.7              |
| 0.11              | 1.12                     | UserConnect-4.7              |
| 0.11              | 1.13                     | UserConnect-4.7              |
| 0.11              | 1.14                     | UserConnect-4.7              |
| 0.11              | 1.15                     | UserConnect-4.7              |
| **0.12 <sup>1</sup>** | 1.16                 | UserConnect-4.7              |
| 0.12              | **2.0 <sup>2</sup>**     | UserConnect-4.7              |
| 0.12              | 2.1                      | UserConnect-4.7              |
| 0.12              | 2.2                      | UserConnect-4.7              |
| 0.12              | 2.3                      | UserConnect-5.0              |
| 0.12              | 2.4                      | UserConnect-5.0              |
| 0.12              | 2.5                      | UserConnect-5.1              |

**<sup>1</sup>** : Note that Terraform v0.12 is not backwards-compatible with previous Terraform versions. For R1.16, there will be a need to change some syntax in the Terraform configuration files. Please see Hashicorp's [announcement](https://www.hashicorp.com/blog/announcing-terraform-0-12) for more information

**<sup>2</sup>** : With R2.0, there is major code restructuring that includes attribute/resource renaming and changes to attribute values. We *highly* recommend customers reference the [R2.0 upgrade guide](https://www.terraform.io/docs/providers/aviatrix/guides/v2-upgrade-guide.html) for detailed instructions before upgrading to R2.0
