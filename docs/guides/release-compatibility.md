---
layout: "aviatrix"
page_title: "Release Compatibility Chart"
description: |-
  The Aviatrix provider Release Compatibility Guide
---

# Aviatrix Provider: Release Compatibility Chart

## USAGE:
Quick at-a-glance access to Aviatrix Terraform provider's release compatibility with the most up-to-date Controller release versions. New resources and features may be tracked in the Release Notes.

-> **NOTE:** This only provides a quick glance at version compatibility between platforms.


NOTICE: Starting with 8.0.0 release, the Terraform provider will synchronize its version with the Aviatrix Controller version. The provider and controller will share the same version number, both following semantic versioning (major.minor.patch) format.


| Terraform Version (v) | Aviatrix Provider Version (R) |   Supported Controller Version   |
| :-------------------: | :---------------------------: | :------------------------------: |
|          1.0          |             x                 |       x                          |
|          1.0          |             8.2.0             |       8.2.0                      |
|          1.0          |             8.1.20            |       8.1.20                     |
|          1.0          |             8.0.40            |       8.0.40                     |
|          1.0          |             8.1.10            |       8.1.10                     |
|          1.0          |             8.0.30            |       8.0.30                     |
|          1.0          |             8.1.1             |       8.1.1                      |
|          1.0          |             8.1.0             |       8.1.0                      |
|          1.0          |             8.0.10            |       8.0.10                     |
|          1.0          |             8.0.0             |       8.0.0                      |
|          1.0          |             3.2.2             |       UserConnect-7.2.5090       |
|          1.0          |             3.2.1             |       UserConnect-7.2.4996       |
|          1.0          |             3.2.0             |       UserConnect-7.2.4820       |


## Example:
If your Aviatrix Controller is on 8.1.10, you should be using Hashicorp Terraform v1.0 and our Aviatrix provider R8.1.10.

We recommend always using a provider version that matches your Controller version, as they are synchronized and share the same version number.

```hcl

provider "aviatrix" {
  controller_ip           = "1.2.3.4"
  username                = "admin"
  password                = "password"
  skip_version_validation = false
  version                 = "2.19.0"
}
```
