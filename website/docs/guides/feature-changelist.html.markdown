---
layout: "aviatrix"
page_title: "Guides: Aviatrix R1.x Feature Changelist"
description: |-
  The Aviatrix provider R1.x Feature Changelist
---

# Aviatrix R1.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R1.0 to present. New resources may be tracked in the [Release Notes](https://www.terraform.io/docs/providers/aviatrix/guides/release-notes.html).

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

Otherwise, this list does not really apply. Please view the below list for details regarding this:

1. If a customer is transitioning to use Terraform to manage existing infrastructure, it is recommended to start from the latest Controller and Aviatrix Terraform provider version, and use the Terraform Export feature and import their infrastructure for a quick and easy migration.
  - **Please note that "Export" is still a beta feature (and only up-to-date for 4.7 at the moment)**
  - Customer can still choose to manually write their config file to their own specifications and use ``terraform import`` to bring their infrastructure into Terraform state management
2. If a customer is adopting Terraform for the first time, with no pre-existing infrastructure/ for managing new infrastructure, this changelist does not apply whatsoever. They simply need to follow our [Terraform tutorial](https://docs.aviatrix.com/HowTos/tf_aviatrix_howto.html) to setup their Terraform environment and use the Terraform documentation featured on the left sidebar.

---

``Last updated: R1.16 (UserConnect-4.7-patch); Terraform v0.12``


---
## R1.0.242 (UserConnect-4.0)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | transit_vpc    | enable_hybrid_connection | **No**, Terraform support added to an existing functionality on Controller. **Unless** the customer has enabled this on that transit gw through Controller, they must rectify the diff by adding that line in their config file and doing a ``terraform refresh`` |
|(new) |  spoke_vpc     | ha_gw_size        | **Yes**; if customer has HA enabled for their spoke gw, they must also update their config file for ea/ spoke gw's respective sizes |
|(new) | transit_vpc    | ha_gw_size        | **Yes**; if customer has HA enabled for their transit gw, they must also update their config file for ea/ transit gw's respective sizes |


## R1.1.66 (UserConnect-4.1.981)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | gateway        | peering_ha_eip    | **No**, Terraform support added to an existing functionality on Controller. **Unless** the customer has enabled HA Peering through the Controller and manually specified the EIP of the HA GW, they must rectify the diff by adding that line in their config file and doing a ``terraform refresh`` |
|(deprecated)| tunnel   | over_aws_peering  | **Yes**, this attribute is deprecated. If users plan to continue to use AWS peering, please use the ``aws_peer`` resource instead |
|(new) | transit_vpc    | enable_nat        | **No**; Terraform support added to an existing functionality on Controller. **Unless** the customer has enabled (SNAT) on that transit gw through Controller, they must rectify the diff by adding that line in their config file and doing a ``terraform refresh`` |
|(new) | transit_vpc    | connected_transit | **No**; Terraform support added to an existing functionality on Controller. **Unless** the customer has enabled this on that transit gw through Controller, they must rectify the diff by adding that line in their config file and doing a ``terraform refresh`` |
|(new) | spoke_vpc      | single_az_ha      | **No**; Terraform support added to an existing 3.1 feature on Controller. **Unless** the customer has enabled this on that spoke gw (instead of HA) through the Controller, they must rectify the diff by adding that line in their config file and doing a ``terraform refresh`` |


## R1.3.12 (UserConnect-4.1.982, 4.2.634)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated)| dc_extn  | --                | **Yes**; this resource is deprecated. If this was being used in Terraform config files, remove it |


## R1.4.4 (UserConnect-4.2.634)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | aws_tgw        | manage_vpc_attachment | **Yes**; if the ``aws_tgw`` resource was originally created with provider <R1.4.4 (<Controller 4.2), you must do a ``terraform refresh`` to update and apply the attribute's default value ("true") into the state file |


## R1.5.24 (UserConnect-4.2.764)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | gateway | dns_server        | **Yes**; if customer has this in their config file, must comment it out/ remove |
|(deprecated) | transit_vpc | dns_server    | **Yes**; if customer has this in their config file, must comment it out/ remove |
|(deprecated) | spoke_vpc   | dns_server    | **Yes**; if customer has this in their config file, must comment it out/ remove |


## R1.6.29 (UserConnect-4.2.764)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | fqdn           | gw_filter_tag_list| **Yes**; new attribute implemented as a map for FQDN resource. gw_list is now implemented and nested under this attribute. From this release forward, any FQDN resource must follow the new format |
|(changed) | fqdn       | gw_list           | **Yes**; gateway list is now nested under ``gw_filter_tag_list`` |
|(new) | spoke_vpc      | ha_zone           | **No**; GCP support for spoke gateway added in Terraform. **Unless**, customer has GCP HA spoke gateway created through Controller that they want to import into Terraform, they must specify the ``ha_zone`` in that specific spoke resource in their config file, along with the originally required ``ha_subnet`` and ``ha_gw_size`` |
|(new) | spoke_vpc      | vnet_and_resource_group_names | **Yes**; Azure support for spoke gateway added in Terraform. When creating an Azure spoke gateway, this attribute must be specified in place of ``vpc_id``, which is used for the AWS-counterpart. **However**, note that this is deprecated in a future release (R1.10.10) and replaced with ``vpc_id`` |
|(new) | transit_vpc    | vnet_name_resource_group | **Yes**; Azure support for transit gateway added in Terraform. When creating an Azure transit gateway, this attribute must be specified in place of ``vpc_id``, which is used for the AWS-counterpart. **However**, note that this is deprecated in a future release (R1.10.10) and replaced with ``vpc_id`` |


## R1.8.26 (UserConnect-4.3.1275)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | gateway | ha_subnet         | **Yes**; enabling traditional gateway HA has been deprecated on Controller already. However, functionality support has not yet been removed. As such, we will be removing Terraform support for enabling this feature. You will need to remove this attribute from any gateway resource that currently has gateway HA enabled through this attribute. To disable this gateway HA, customers must do so through Controller |
|(new) | gateway        | peering_ha_gw_size| **Yes**; if customer has peering HA enabled for their gateway, they must also update their config file for ea/ gateway's respective sizes |
|(new) | gateway        | peering_ha_zone   | **No**; GCP support added for gateway. **Unless**, customer has GCP HA gateway created through Controller that they want to import into Terraform, they must specify the ``peering_ha_zone`` in that specific gateway resource in their config file, along with the originally required ``peering_ha_gw_size``. |
|(new) | transit_vpc    | insane_mode, insane_mode_az | **No**; Terraform support added to an existing 4.1 feature on Controller. **Unless**, customer has an insane mode-enabled transit gateway created through Controller that they want to import into Terraform, they must specify this attribute, along with ``insane_mode_az`` |
|(new) | transit_vpc    | ha_insane_mode_az | **No**; see above for details. This attribute is only required if the insane mode-enabled transit gateway also has HA enabled |


## R1.9.28 (UserConnect-4.6.569)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | tunnel  | cluster           | **Yes**; if ``cluster`` was being set in this resource, remove from config file |
|(deprecated) | admin_email | --            | **Yes**; if the admin email is being set using this resource, remove from .tf file. It is no longer supported |
|(deprecated) | customer_id | --            | **Yes**; if the customer id is being set using this resource, remove from .tf file. It is no longer supported |
|(new) | site2cloud     | enable_dead_peer_detection, local/remote_subnet_virtual, custom_algorithms, private_route_encryption | **No**; Terraform support added to existing S2C features on Controller. Please see site2cloud documentation for details |
|(new) | vgw_conn       | enable_advertise_transit_cidr | **No**; Terraform support added to existing 4.2 feature on Controller |
|(new) | vpc            | aviatrix_firenet_vpc | **No**; Terraform support added for new Firenet VPC creation. Original VPC resource added in R1.7.18 |
|(new) | transit_vpc    | enable_firenet_interfaces | **No**; Terraform support added for new Firenet 4.3 feature on Controller |


## R1.10.10 (UserConnect-4.6.604)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | spoke_vpc | vnet_and_resource_group_names | **Yes**; if an Azure spoke gateway is being used, the VNet and ResourceGroup is no longer being set in this attribute. Remove and replace with ``vpc_id``. This is part of the standardization in behavior across all gateway-type resources for Terraform |
|(deprecated) | transit_vpc | vnet_name_resource_group | **Yes**; if an Azure transit gateway is being used, the VNet and ResourceGroup is no longer being set in this attribute. Remove and replace with ``vpc_id``. This is part of standardization in behavior across all gateway-type resources for Terraform |


## R1.11.11 (UserConnect-4.7.378)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | version | --                | **Yes**; if Controller version was being managed using this resource, remove from .tf file. This functionality has been moved to **controller_config** resource as an attribute |
|(new) | controller_config | target_version | **No**; see above. This is an optional attribute and only needs to be set one wants/ requires Controller version management to be done through Terraform |
|(new) | vgw_conn       | bgp_manual_spoke_advertise_cidrs | **No**; Terraform support added for new 4.7 feature on Controller |


## R1.12.12 (UserConnect-4.7.378)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | aws_tgw_vpn_conn | inside_ip_cidr_tun_1, pre_shared_key_tun_1/2 | **No**; Terraform support added to a new 4.6 feature on Controller. The original new **aws_tgw_vpn_conn** resource was released in R1.11.11 |


## R1.14.15 (UserConnect-4.7.474)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | gateway        | max_vpn_conn      | **Yes**; if any vpn-gateway was created before this release, add this attribute to .tf file(s) and specify the corresponding value(s), which may be viewed on the Controller's Gateway Page. Any new vpn-gateway created on or after this release must specify this parameter. Default is 100 |


## R1.16.20 (UserConnect-4.7.520) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | aws_tgw    | security_domains, attached_vpc | **Yes**; due to Hashicorp's Terraform v0.12 release, syntactical changes were introduced. Most notably, map attributes become written as separate blocks |
|(changed) | firewall   | policy            | **Yes**; see above for details |
|(changed) | firewall_tag | cidr_list       | **Yes**; see above for details |
|(changed) | fqdn       | gw_filter_tag_list, domain_names | **Yes**; see above for details |
|(changed) | vpn_profile| policy            | **Yes**; see above for details |
