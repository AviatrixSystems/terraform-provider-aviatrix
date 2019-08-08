# Aviatrix Terraform Provider Upgrade Guide (R1.x to R2.x)
**USAGE:** For customers who are currently already using Terraform to manage their infrastructure prior to Controller 4.7 and are looking to upgrade their Controller to 4.7.501+ due to Field Notice 005 or other reasons, please follow along for guidance on the upgrade process to ensure a smooth transition.

---
In summary:

**Current Setup/ Configuration:**
- Controller <4.7
- Terraform v0.11.x
- Aviatrix Terraform provider R1.xx
- Looking to upgrade Controller to 4.7+

---
## Context
With the Aviatrix Terraform provider R2.0, there is major restructuring of our code as well as major changes such as renaming of attributes, resources, and attribute values. All these changes are all in the name of standardization of naming conventions and resources. Although we recognize that it is a major inconvenience to customers, we believe that these changes will benefit everyone in the long-term not only for customer clarity but ease of future feature-implementations and code maintenance.

For most changes, unless stated otherwise, after editing the respective (.tf) files, a simple ``terraform refresh`` should rectify the state of the infrastructure.

**NOTE:** For future releases, Aviatrix Terraform provider R2.1+ is only compatible with Terraform v0.12 and Aviatrix Controller 5.0+.

In this document, we broke down the upgrade process in multiple phases to make it easier to follow along.

In summary:

| Phase | Notes           | From (Controller / Provider / Terraform) | To (Controller / Provider / Terraform)  | Effort   |
|:-----:|:---------------:|:----------------------------------------:|:---------------------------------------:|:--------:|
| 1     | Upgrade Controller | <4.7 / R1.xx / v0.11.x | **4.7** / R1.xx / v0.11.x | low |
| 2     | Upgrade Terraform | 4.7 / R1.xx / v0.11.x | 4.7 / R1.xx / **v0.12.x** | low   |
| 3     | Upgrade Provider  | 4.7 / R1.xx / v0.12.x   | 4.7 / **R2.0** / v0.12.x | non-trivial |
|FUTURE | Upgrade both      | 4.7 / R2.0 / v0.12.x    | **5.0+** / **R2.0+** / v0.12.x | low |

---
## Phase 1: Upgrading Aviatrix Controller to 4.7
**Summary:**
This phase involves first upgrading the Controller to 4.7. Afterwards, customers must update their (.tf) files as necessary, then upgrade their Aviatrix Terraform provider version, and then perform ``terraform refresh`` to rectify their state.

1. Upgrade Controller through Aviatrix Console Web GUI. Please see below links for references:
  - Controller Upgrade documentation:
    - https://docs.aviatrix.com/HowTos/inline_upgrade.html
  - Controller Release Notes:
    - https://docs.aviatrix.com/HowTos/UCC_Release_Notes.html
- Update Terraform files (.tf) as necessary. Please reference documentation linked below to note specific changes to any resource attributes that you may be using in your configuration:
  - Please refer **up to** R1.14's table (for Terraform 0.11) (compatible with Controller 4.7):
  - https://github.com/AviatrixSystems/terraform-provider-aviatrix/blob/master/website/docs/guides/feature-changelist.html.md
  - Aviatrix Terraform Provider Release Notes:
    - https://github.com/AviatrixSystems/terraform-provider-aviatrix/releases
- Update Aviatrix Terraform provider:
  - Navigate to Aviatrix Terraform provider directory:
    - Mac:
      - ``cd $GOPATH/src/github.com/terraform-providers``
    - Windows:
      - ``cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-aviatrix``
  - Update repository by pulling changes:
    - ``git pull``
  - Change branch to **UserConnect-4.7-TF.11**:
    - ``git checkout UserConnect-4.7-TF.11``
  - Build the provider:
    - Mac:
      - ``make fmt``
      - ``make build``
    - Windows:
      - ``go fmt``
      - ``go install``
- Navigate to your Terraform directory/ directories and refresh Terraform state to update:
  - ``terraform refresh``
  - Please perform a plan command to note any deltas:
    - ``terraform plan``
    - If there are any deltas still, you may fix and run refresh/ plan again accordingly

---
## Phase 2: Upgrading Hashicorp's Terraform to v0.12
**Summary:**
This phase involves upgrading Hashicorp's Terraform from v0.11 to v0.12. As far as the Aviatrix Terraform provider is concerned, Hashicorp's Terraform's v0.12 only involves syntactical changes. In addition, customers must also update their (.tf) files as necessary, then upgrade their Aviatrix Terraform provider version, and then perform ``terraform refresh`` to rectify their state.

1. Upgrade Terraform version to v0.12:
  - Hashicorp Terraform Upgrade documentation:
    - https://www.terraform.io/upgrade-guides/0-12.html
- Update Terraform files (.tf) as necessary. Please reference documentation linked below to note specific changes to any resource attributes that you may be using in your configuration:
  - Please follow the **R1.16's table** (for Terraform 0.12) (compatible with Controller 4.7):
    - https://github.com/AviatrixSystems/terraform-provider-aviatrix/blob/master/website/docs/guides/feature-changelist.html.md#r11620-userconnect-47520-terraform-v012
  - The same table is shown below for your convenience:

### R1.16.20 (UserConnect-4.7.520) (Terraform v0.12)

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | aws_tgw    | security_domains, attached_vpc | **Yes**; due to Hashicorp's Terraform v0.12 release, syntactical changes were introduced. Most notably, map attributes become written as separate blocks |
|(changed) | firewall   | policy            | **Yes**; see above for details |
|(changed) | firewall_tag | cidr_list       | **Yes**; see above for details |
|(changed) | fqdn       | gw_filter_tag_list, domain_names | **Yes**; see above for details |
|(changed) | vpn_profile| policy            | **Yes**; see above for details |

3. Update Aviatrix Terraform provider:
  - Navigate to Aviatrix Terraform provider directory:
    - Mac:
      - ``cd $GOPATH/src/github.com/terraform-providers``
    - Windows:
      - ``cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-aviatrix``
  - Update repository by pulling changes:
    - ``git pull``
  - Change branch to **UserConnect-4.7-TF.12-v1**:
    - ``git checkout UserConnect-4.7-TF.12-v1``
  - Build the provider:
    - Mac:
      - ``make fmt``
      - ``make build``
    - Windows:
      - ``go fmt``
      - ``go install``
- Navigate to your Terraform directory/ directories and refresh Terraform state to update:
  - ``terraform refresh``
  - Please perform a plan command to note any deltas:
    - ``terraform plan``
    - If there are any deltas still, you may fix and run refresh/ plan again accordingly

---
## Phase 3: Upgrading Aviatrix Terraform Provider to R2.0
**Summary:** This will be the largest phase in terms of the upgrade process. While this phase only involves upgrading the customer's (.tf) files as necessary, the amount of changes from R1.xx to R2.0+ are not small. Afterwards, customers must upgrade their Aviatrix Terraform provider version, and then perform ``terraform refresh`` or ``terraform import`` as necessary.
1. Update Terraform files (.tf) as necessary. Please reference documentation linked below to note specific changes to any resource attributes that you may be using in your configuration:
  - https://github.com/AviatrixSystems/terraform-provider-aviatrix/blob/master/website/docs/guides/feature-changelist-v2.html.md
- Update Aviatrix Terraform provider:
  - Navigate to Aviatrix Terraform provider directory:
    - Mac:
      - ``cd $GOPATH/src/github.com/terraform-providers``
    - Windows:
      - ``cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-aviatrix``
  - Update repository by pulling changes:
    - ``git pull``
  - Change branch to **UserConnect-4.7-TF.12-v2**:
    - ``git checkout UserConnect-4.7-TF.12-v2``
  - Build the provider:
    - Mac:
      - ``make fmt``
      - ``make build``
    - Windows:
      - ``go fmt``
      - ``go install``
- Navigate to your Terraform directory/ directories and refresh/ import to update:
  - **Rules to determine whether to refresh or import:**
  - **NOTE:** Rules listed in terms of priority. Import takes precedence over refresh; meaning if your resource has one of the below import rules apply, that resource will require an import no matter how many of the refresh rules apply to it
    - If you are using **transit_vpc** or **spoke_vpc** resources, note they are deprecated and support will eventually be removed. **transit_gateway** and **spoke_gateway** will replace them respectively. ``terraform import`` for these resources must be done to rectify the state
      - **Please refer to the documentation in step 1 for more detailed instructions**
    - If your resource uses an attribute whose accepted values changed, ``terraform import`` must be done to rectify the state
    - If your resource uses an attribute that has been renamed, ``terraform refresh`` is sufficient to rectify the state
    - If your resource uses an attribute that has been renamed **AND** had its accepted values changed, ``terraform refresh`` is sufficient to rectify the state

---
## Beyond Phase 3: to infinity and beyond ~
- Any updates/ future releases for the Aviatrix Terraform provider will continue to be documented here:
  - https://github.com/AviatrixSystems/terraform-provider-aviatrix/releases
- Any updates/ future releases for R2.0+ that might impact customers will continue to be documented here:
  - https://github.com/AviatrixSystems/terraform-provider-aviatrix/blob/master/website/docs/guides/feature-changelist-v2.html.md
  - Any future necessary changes will only be simple and only require small tweaks and a ``terraform refresh``
