---
layout: "aviatrix"
page_title: "Aviatrix: Upgrade Guide"
description: |-
  The Aviatrix provider Upgrade Guide
---

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
2. Update Terraform files (.tf) as necessary. Please reference documentation linked below to note specific changes to any resource attributes that you may be using in your configuration:
  - Please refer **up to** R1.14's table (for Terraform 0.11) (compatible with Controller 4.7):
    - https://www.terraform.io/docs/providers/aviatrix/guides/feature-changelist.html
  - Aviatrix Terraform Provider Release Notes:
    - https://github.com/terraform-providers/terraform-provider-aviatrix/releases
3. Update Aviatrix Terraform provider:
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
4. Navigate to your Terraform directory/ directories and refresh Terraform state to update:
  - ``terraform init`` - to reinitialise a working Terraform environment for the current directory
  - ``terraform refresh``
  - Please perform a plan command to note any deltas:
    - ``terraform plan``
    - If there are any deltas still, you may fix and run refresh/ plan again accordingly

### Example Walk-through
Say you are upgrading from Controller 4.1, Terraform v0.11, Aviatrix Terraform provider R1.1. And for example, you are using Terraform to manage your topology consisting of single AWS VPN gateway that you use to connect VPN users. Your Terraform file(s) might look something like this:

```
# VPN_setup.tf

...

resource "aviatrix_gateway" "aws_vpn_gw" {
  cloud_type    = 1
  account_name  = "devops"
  gw_name       = "aws_vpn_gw"
  vpc_id        = "vpc-abc123"
  vpc_reg       = "us-east-1"
  vpc_size      = "t2.micro"
  vpc_net       = "10.0.0.0/24"

  dns_server    = "8.8.8.8"
  vpn_access    = "yes"
  vpn_cidr      = "192.168.43.0/24"
  enable_elb    = "yes"
  elb_name      = "example-elb-name"

  ...

}

resource "aviatrix_vpn_user" "vpn_user_1" {
  ...
}

...
```
According to the table **R1.5** for UserConnect-4.2, the attribute ``dns_server`` for gateway resource is deprecated and to remove it from your file(s) if you are using it. In addition, according to the table **R1.14** for UserConnect-4.7, the attribute ``max_vpn_conn`` is a new and required for any vpn-gateway from this release moving forward. So in summary, in this example, you will have to remove ``dns_server`` and add ``max_vpn_conn`` into your gateway resource. Your file(s) should now look something like this (only relevant section shown):

```
# VPN_setup.tf

...

resource "aviatrix_gateway" "aws_vpn_gw" {
  cloud_type    = 1
  account_name  = "devops"
  ...
  vpc_net       = "10.0.0.0/24"

  # note that dns_server is removed
  vpn_access    = "yes"
  vpn_cidr      = "192.168.43.0/24"
  max_vpn_conn  = 100 # newly added
  enable_elb    = "yes"
  elb_name      = "example-elb-name"

  ...

}

...
```
After completion of editing your Terraform files, the rest of this phase is as simple as updating your Aviatrix Terraform provider repository and using the correct corresponding release.

Navigate to your local Aviatrix Terraform provider repository, which by default, if setup according to our initial setup doc [here](https://github.com/terraform-providers/terraform-provider-aviatrix/blob/master/README.md), is ``$GOPATH/src/github.com/terraform-providers``. We will then use Git to pull latest changes from our remote repository (``git pull``), and switch to the branch that corresponds with our Controller version (``git checkout UserConnect-4.7-TF.11``). Build the provider as according to Step 3 depending on your OS. Finally, navigate back to where your Terraform files reside, and perform a ``terraform init`` to reinitialise the Terraform environment for the current directory based on the new provider, and run a ``terraform refresh``. A ``terraform plan`` may be performed to catch any deltas. If there are still any deltas, you may fix and repeat the refresh/ plan steps again.

---
## Phase 2: Upgrading Hashicorp's Terraform to v0.12
**Summary:**
This phase involves upgrading Hashicorp's Terraform from v0.11 to v0.12. As far as the Aviatrix Terraform provider is concerned, Hashicorp's Terraform's v0.12 only involves syntactical changes. In addition, customers must also update their (.tf) files as necessary, then upgrade their Aviatrix Terraform provider version, and then perform ``terraform refresh`` to rectify their state.

1. Upgrade Terraform version to v0.12:
  - Hashicorp Terraform Upgrade documentation:
    - https://www.terraform.io/upgrade-guides/0-12.html
  - Hashicorp Terraform 0.12 Release Notes:
    - https://www.hashicorp.com/blog/announcing-terraform-0-12
2. Update Terraform files (.tf) as necessary. Please reference documentation linked below to note specific changes to any resource attributes that you may be using in your configuration:
  - Please follow the **R1.16's table** (for Terraform 0.12) (compatible with Controller 4.7):
    - https://github.com/terraform-providers/terraform-provider-aviatrix/blob/master/website/docs/guides/feature-changelist.html.#r11620-userconnect-47520-terraform-v012
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
4. Navigate to your Terraform directory/ directories and refresh Terraform state to update:
  - ``terraform init`` - to reinitialise a working Terraform environment for the current directory
  - ``terraform refresh``
  - Please perform a plan command to note any deltas:
    - ``terraform plan``
    - If there are any deltas still, you may fix and run refresh/ plan again accordingly

### Example Walk-through
As stated previously, as far as the Aviatrix Terraform provider is concerned, Hashicorp's Terraform's v0.12 only involves syntactical changes in regards to certain resources. You will need to first follow Hashicorp's [upgrade instructions](https://www.terraform.io/upgrade-guides/0-12.html) as seen in Step 1 of this phase; please see the [Hashicorp Terraform v0.12 release notes](https://www.hashicorp.com/blog/announcing-terraform-0-12) as well. Upon doing so, you will need to update your Terraform file(s) to accommodate for the new syntactical changes.

For example, if your Terraform manages a VPN-configuration that includes a vpn-gateway, vpn-users and some profiles, it might currently look something like this in v0.11:

```
# v0.11 Example .tf file
...

# VPN-related resources
resource "aviatrix_gateway" "vpn-gw" {
  ...
}

resource "aviatrix_vpn_user" "vpn-user-1" {
  ...
}

resource "aviatrix_vpn_profile" "vpn-profile-1" {
  name      = "vpn_profile_1"
  base_rule = "allow_all"
  users     = ["vpn-user-1"]
  policy = [
    {
      action  = "deny"
      proto   = "tcp"
      port    = "443"
      target  = "10.0.0.0/32"
    },
    {
      action  = "deny"
      proto   = "tcp"
      port    = "443"
      target  = "10.0.0.1/32"
    }
  ]
  depends_on = ["aviatrix_vpn_user.vpn-user-1"]
}
```
With the new Terraform v0.12, the vpn_profile resource, along with others, as documented in the **R1.16** table above, the map attributes are now written as separate blocks, as seen below (only relevant section shown):
```
# v0.12 Example .tf file
...

# VPN-related resources
...

resource "aviatrix_vpn_profile" "vpn-profile-1" {
  name      = "vpn_profile_1"
  base_rule = "allow_all"
  users     = ["vpn-user-1"]

  # Note the change: each policy is written as a separate block now
  policy {
    action  = "deny"
    proto   = "tcp"
    port    = "443"
    target  = "10.0.0.0/32"
  }
  policy {
    action  = "deny"
    proto   = "tcp"
    port    = "443"
    target  = "10.0.0.1/32"
  }
  depends_on = ["aviatrix_vpn_user.vpn-user-1"]
}
```
Once again, as in Phase 1, you will have to update your Aviatrix Terraform provider repository and use the correct corresponding release.

Navigate to your local Aviatrix Terraform provider repository, which by default, if setup according to our initial setup doc [here](https://github.com/terraform-providers/terraform-provider-aviatrix/blob/master/README.md), is ``$GOPATH/src/github.com/terraform-providers``. We will then use Git to pull latest changes from our remote repository (``git pull``), and switch to the branch that corresponds with our Controller version (``git checkout UserConnect-4.7-TF.12-v1``). Build the provider as according to Step 3 depending on your OS. Finally, navigate back to where your Terraform files reside, and perform a ``terraform init`` to reinitialise the Terraform environment for the current directory based on the new provider, and run a ``terraform refresh``. A ``terraform plan`` may be performed to catch any deltas. If there are still any deltas, you may fix and repeat the refresh/ plan steps again.

---
## Phase 3: Upgrading Aviatrix Terraform Provider to R2.0
**Summary:** This will be the largest phase in terms of the upgrade process. While this phase only involves upgrading the customer's (.tf) files as necessary, the amount of changes from R1.xx to R2.0+ are not small. Afterwards, customers must upgrade their Aviatrix Terraform provider version, and then perform ``terraform refresh`` or ``terraform import`` as necessary.
1. Update Terraform files (.tf) as necessary. Please reference documentation linked below to note specific changes to any resource attributes that you may be using in your configuration:
  - https://github.com/terraform-providers/terraform-provider-aviatrix/blob/master/website/docs/guides/feature-changelist-v2.html
2. Update Aviatrix Terraform provider:
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
3. Navigate to your Terraform directory/ directories and refresh/ import to update:
  - ``terraform init`` - to reinitialise a working Terraform environment for the current directory
  - **Rules to determine whether to refresh or import:**
  - **NOTE:** Rules are listed in terms of priority. Import takes precedence over refresh; meaning if your resource has one of the below import rules apply, that resource will require an import no matter how many of the refresh rules apply to it
    - If you are using **transit_vpc** or **spoke_vpc** resources, note they are deprecated and support will eventually be removed. **transit_gateway** and **spoke_gateway** will replace them respectively. ``terraform import`` for these resources must be done to rectify the state
      - **Please refer to the documentation in step 1 for more detailed instructions**
    - If your resource uses an attribute whose accepted values changed, ``terraform import`` must be done to rectify the state
    - If your resource uses an attribute that has been renamed, ``terraform refresh`` is sufficient to rectify the state
    - If your resource uses an attribute that has been renamed **AND** had its accepted values changed, ``terraform refresh`` is sufficient to rectify the state

### Example Walk-through
Phase 3 will definitely be a larger task, but is still nothing too different from what has been done in the previous 2 phases. Here we will use an example to demonstrate the refresh/ import rules more clearly.

Let's say you have finished Phase 2, and are currently on Controller 4.7, Terraform v0.12 and now you need to upgrade your Aviatrix Terraform provider to R2.0. Let's take the completed vpn-configuration example from Phase 1:
```
# VPN_setup.tf

...

resource "aviatrix_gateway" "aws_vpn_gw" {
  cloud_type    = 1
  account_name  = "devops"
  gw_name       = "aws_vpn_gw"
  vpc_id        = "vpc-abc123"
  vpc_reg       = "us-east-1"
  vpc_size      = "t2.micro"
  vpc_net       = "10.0.0.0/24"

  # note that dns_server is removed
  vpn_access    = "yes"
  vpn_cidr      = "192.168.43.0/24"
  max_vpn_conn  = 100 # newly added
  enable_elb    = "yes"
  elb_name      = "example-elb-name"

  enable_nat    = "yes"
  ...

}
```

We will now need to update our Terraform file(s). In this case, note that due to the attribute re-naming, we will have to rename attributes such as ``vpc_size`` and ``vpc_net`` to ``gw_size`` and ``subnet``, respectively, as it is much more clear as to what these attributes refer to. ``enable_nat`` must also be changed to ``enable_snat`` here. In addition, due to boolean standardization, ``enable_snat``'s accepted value is changed from 'yes'/'no' to true/ false. ``vpn_access`` here did not get renamed, but the accepted value has also changed to boolean.

The updated file should now look something like:
```
# Updated VPN_setup.tf
...

resource "aviatrix_gateway" "aws_vpn_gw" {
  cloud_type    = 1
  account_name  = "devops"
  gw_name       = "aws_vpn_gw"
  vpc_id        = "vpc-abc123"
  vpc_reg       = "us-east-1"
  gw_size       = "t2.micro" # vpc_size -> gw_size
  subnet        = "10.0.0.0/24" # vpc_net -> subnet

  # note that dns_server is removed
  vpn_access    = true # "yes" -> true
  vpn_cidr      = "192.168.43.0/24"
  max_vpn_conn  = 100 # newly added
  enable_elb    = "yes"
  elb_name      = "example-elb-name"

  enable_snat   = true # enable_nat = "yes" -> enable_snat = true
  ...

}
```

Once again, as in the previous two phases, you will have to update your Aviatrix Terraform provider repository and use the correct corresponding release.

Navigate to your local Aviatrix Terraform provider repository, which by default, if setup according to our initial setup doc [here](https://github.com/terraform-providers/terraform-provider-aviatrix/blob/master/README.md), is ``$GOPATH/src/github.com/terraform-providers``. We will then use Git to pull latest changes from our remote repository (``git pull``), and switch to the branch that corresponds with our Controller version (``git checkout UserConnect-4.7-TF.12-v2``). Build the provider as according to Step 3 depending on your OS. Finally, navigate back to where your Terraform files reside, and perform a ``terraform init`` to reinitialise the Terraform environment for the current directory based on the new provider.

Now, to determine whether or not to do a ``terraform refresh`` or ``terraform import``, refer to the rules in Step 3 of this phase. In this example, according to the stated rules, because attributes ``vpc_size`` and ``vpc_net`` got renamed, and ``enable_nat`` got renamed and changed to accept boolean, it would seem this gateway resource would only need a ``terraform refresh``. However, because this vpn-gateway uses the attribute ``vpn_access``, whose only change is that the accepted value got changed to boolean, this resource would need to a ``terraform import`` instead, because the import rule would take precedence over refresh. After performing the necessary ``terraform refresh`` / ``terraform import`` (per resource), a ``terraform plan`` may be performed to catch any deltas. If there are still any deltas, you may fix and repeat again.

---
## Beyond Phase 3: to infinity and beyond ~
- Any updates/ future releases for the Aviatrix Terraform provider will continue to be documented here:
  - https://github.com/terraform-providers/terraform-provider-aviatrix/releases
- Any updates/ future releases for R2.0+ that might impact customers will continue to be documented here:
  - https://www.terraform.io/docs/providers/aviatrix/guides/feature-changelist-v2.html
  - Any future necessary changes will only be simple and only require small tweaks and a ``terraform refresh``
