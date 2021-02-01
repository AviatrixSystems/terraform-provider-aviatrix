---
layout: "aviatrix"
page_title: "Release Compatibility Chart"
description: |-
  The Aviatrix provider Release Compatibility Guide
---

# Aviatrix Provider: Release Compatibility Chart

## USAGE:
Quick at-a-glance access to Aviatrix Terraform provider's release compatibility with the most up-to-date Controller release versions. New resources and features may be tracked in the Release Notes.

-> **NOTE:** This only provides a quick glance at version compatibility between platforms. Please check the release notes for full details on new features, changes and deprecations [here](https://www.terraform.io/docs/providers/aviatrix/guides/release-notes.html).


---

``Last updated: R2.18 (UserConnect-6.3)``


---


| Terraform Version (v) | Aviatrix Provider Version (R) | Minimum Supported Controller Version |
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
| 0.12              | 2.6                      | UserConnect-5.1              |
| 0.12              | 2.7                      | UserConnect-5.1              |
| 0.12              | 2.8                      | UserConnect-5.1, 5.2         |
| 0.12              | 2.9                      | UserConnect-5.2              |
| 0.12              | 2.9.1                    | UserConnect-5.2.2122         |
| 0.12              | 2.10                     | UserConnect-5.2.2122         |
| 0.12              | 2.11                     | UserConnect-5.3.1391         |
| 0.12              | 2.12                     | UserConnect-5.3.1491         |
| 0.12              | 2.13                     | UserConnect-5.4.1074         |
| 0.12              | 2.14                     | UserConnect-5.4.1201         |
| 0.12              | 2.14.1                   | UserConnect-5.4.1232         |
| 0.12              | 2.15                     | **UserConnect-6.0 <sup>3</sup>** |
| 0.12              | 2.15.1                   | UserConnect-6.0              |
| 0.12              | 2.16                     | UserConnect-6.1              |
| 0.12              | 2.16.1                   | UserConnect-6.1              |
| 0.12              | 2.16.2                   | UserConnect-6.1.1280         |
| 0.12              | 2.16.3                   | UserConnect-6.1.1309         |
| 0.12              | 2.17                     | UserConnect-6.2              |
| 0.12, 0.13        | 2.17.1                   | UserConnect-6.2.1891         |
| 0.12, 0.13        | 2.17.2                   | UserConnect-6.2.1914         |
| 0.12, 0.13        | 2.18                     | UserConnect-6.3              |

**<sup>1</sup>** : Note that Terraform v0.12 is not backwards-compatible with previous Terraform versions. For R1.16, there will be a need to change some syntax in the Terraform configuration files. Please see Hashicorp's [announcement](https://www.hashicorp.com/blog/announcing-terraform-0-12) for more information

**<sup>2</sup>** : With R2.0, there is major code restructuring that includes attribute/resource renaming and changes to attribute values. We *highly* recommend customers reference the [R2.0 upgrade guide](https://www.terraform.io/docs/providers/aviatrix/guides/v2-upgrade-guide.html) for detailed instructions before upgrading to R2.0

**<sup>3</sup>** : Terraform R2.15 will only support Controller 6.0+ and will not be backwards-compatible with any older Controller versions due to major changes and improvements on the Controller backend.

## Example:
If your Aviatrix Controller is on UserConnect-5.0.x, you should be using Hashicorp Terraform v0.12.x and our Aviatrix provider R2.4.

Although R2.3 is also compatible with UserConnect-5.0, we recommend using the latest compatible provider version corresponding to the Controller.

```hcl

provider "aviatrix" {
  controller_ip           = "1.2.3.4"
  username                = "admin"
  password                = "password"
  skip_version_validation = false
  version                 = "2.8.0"
}
```
