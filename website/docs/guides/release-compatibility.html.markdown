---
layout: "aviatrix"
page_title: "Guides: Aviatrix Provider Release Compatibility"
description: |-
  The Aviatrix provider Release Compatibility Guide
---

# Aviatrix  Provider: Release Compatibility Chart

**USAGE:** Quick at-a-glance access to Aviatrix Terraform provider's release compatibility with the Controller release versions. New resources and features may be tracked in the Release Notes.

**NOTE:** This only provides a quick glance at version compatibility between platforms. Please check the release notes for full details on new features, changes and deprecations [here](https://github.com/terraform-providers/terraform-provider-aviatrix/releases).


---

``Last updated: R2.3 (UserConnect-5.0)``


---


| Terraform Version | Aviatrix Provider Version | Supported Controller Version |
|:-----------------:|:-------------------------:|:----------------------------:|
| v0.11             | R1.0                      | UserConnect-4.0              |
| v0.11             | R1.1                      | UserConnect-4.1              |
| v0.11             | R1.2                      | UserConnect-4.1              |
| v0.11             | R1.3                      | UserConnect-4.2              |
| v0.11             | R1.4                      | UserConnect-4.2              |
| v0.11             | R1.5                      | UserConnect-4.2              |
| v0.11             | R1.6                      | UserConnect-4.2              |
| v0.11             | R1.7                      | UserConnect-4.3              |
| v0.11             | R1.8                      | UserConnect-4.3              |
| v0.11             | R1.9                      | UserConnect-4.6              |
| v0.11             | R1.10                     | UserConnect-4.6              |
| v0.11             | R1.11                     | UserConnect-4.7              |
| v0.11             | R1.12                     | UserConnect-4.7              |
| v0.11             | R1.13                     | UserConnect-4.7              |
| v0.11             | R1.14                     | UserConnect-4.7              |
| v0.11             | R1.15                     | UserConnect-4.7              |
| **v0.12 <sup>1</sup>** | R1.16                | UserConnect-4.7              |
| v0.12             | **R2.0 <sup>2</sup>**     | UserConnect-4.7              |
| v0.12             | R2.1                      | UserConnect-4.7              |
| v0.12             | R2.2                      | UserConnect-4.7              |
| v0.12             | R2.3                      | UserConnect-5.0              |

**<sup>1</sup>** : Note that Terraform v0.12 is not backwards-compatible with previous Terraform versions. For R1.16, there will be a need to change some syntax in the Terraform configuration files. Please see Hashicorp's [announcement](https://www.hashicorp.com/blog/announcing-terraform-0-12) for more information

**<sup>2</sup>** : With R2.0, there is major code restructuring that includes attribute/resource renaming and changes to attribute values. We *highly* recommend customers reference the [R2.0 upgrade guide](https://www.terraform.io/docs/providers/aviatrix/guides/v2-upgrade-guide.html) for detailed instructions before upgrading to R2.0
