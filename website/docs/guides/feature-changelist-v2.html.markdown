---
layout: "aviatrix"
page_title: "Guides: Aviatrix R2.x Feature Changelist"
description: |-
  The Aviatrix provider R2.x Feature Changelist
---

# Aviatrix R2.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R2.0 to present. New resources may be tracked in the [Release Notes](https://www.terraform.io/docs/providers/aviatrix/guides/release-notes.html).

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

Otherwise, this list does not really apply. Please view the below list for details regarding this:

1. If a customer is transitioning to use Terraform to manage existing infrastructure, it is recommended to start from the latest Controller and Aviatrix Terraform provider version, and use the Terraform Export feature and import their infrastructure for a quick and easy migration.
  - **Please note that "Export" is still a beta feature (and only up-to-date for U4.7 at the moment)**
  - Customer can still choose to manually write their config file to their own specifications and use ``terraform import`` to bring their infrastructure into Terraform state management
2. If a customer is adopting Terraform for the first time, with no pre-existing infrastructure/ for managing new infrastructure, this changelist does not apply whatsoever. They simply need to follow our [Terraform tutorial](https://docs.aviatrix.com/HowTos/tf_aviatrix_howto.html) to setup their Terraform environment and use the Terraform documentation featured on the left sidebar.

We **highly** recommend customers that are starting to adopt Terraform to manage their infrastructure, to start with our Release 2.0+ (R2.0+), which is compatible with our Controller 4.7+.

---

``Last updated: R2.8 (UserConnect-5.2)``


---

## R2.0 (UserConnect-4.7-patch) (Terraform v0.12)
**NOTICE:** With the Release of Aviatrix v.2, there is major restructuring of our code as well as major changes such as renaming of attributes, resources, and attribute values. All these changes are all in the name of standardization of naming conventions and resources. Although we recognize that it is a major inconvenience to customers, we believe that these changes will benefit everyone in the long-term not only for customer clarity but ease of future feature-implementations and code maintenance.

For most changes, unless stated otherwise in the tables below, after editing the respective .tf files, a simple ``terraform refresh`` should rectify the state of the infrastructure.

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | spoke_vpc  | spoke_gateway  | **Yes**; **spoke_vpc** will be deprecated and replaced with **spoke_gateway**. It will be kept for now for backward compatibility and will be removed in the future. Please use **spoke_gateway** instead <br><br> You will need to remove any pre-existing **spoke_vpc** resources from the statefile, update your .tf files to **spoke_gateway**, as well as any attributes (and corresponding values) as necessary, and ``terraform import`` as **aviatrix_spoke_gateway** <br><br> Attributes that are renamed include: ``vpc_size``, ``enable_nat`` <br><br>Attributes whose accepted values are changed to booleans include: ``enable_nat``, ``single_az_ha``|
|(deprecated) | transit_vpc | transit_gateway | **Yes**; **transit_vpc** will be deprecated and replaced with **transit_gateway**. It will be kept for now for backward compatibility and will be removed in the future. Please use **transit_gateway** instead <br><br> You will need to remove any pre-existing **transit_vpc** resources from the statefile, update your .tf files to **transit_gateway**, as well as any attributes (and corresponding values) as necessary, and ``terraform import`` as **aviatrix_transit_gateway** <br><br> Attributes that are renamed include: ``vpc_size``, ``enable_nat``<br><br> Attributes whose accepted values are changed to booleans include: ``enable_nat``, ``connected_transit`` |

### Attribute Renaming
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | firewall   | base_allow_deny   | **Yes**; renamed to ```base_policy``` |
|(changed) | --   | base_log_enable   | **Yes**; renamed to ```base_log_enabled``` |
|(changed) | --   | policy: allow_deny | **Yes**; renamed to ```action``` |
|(changed) | --   | policy: log_enable | **Yes**; renamed to ```log_enabled``` |
|(changed) | fqdn       | fqdn_status       | **Yes**; renamed to ```fqdn_enabled``` |
|(changed) | gateway    | vpc_size          | **Yes**; change any gateways' ```vpc_size``` to ```gw_size``` |
|(changed) | gateway    | vpc_net           | **Yes**; change any gateways' ```vpc_net``` to ```subnet``` |
|(changed) | gateway    | enable_nat        | **Yes**; change any gateways' ```enable_nat``` to ```enable_snat``` |
|(changed) | spoke_vpc --> spoke_gateway  | vpc_size          | **Yes**, if migrating to **spoke_gateway** resource; change any **spoke_vpc** resources' ```vpc_size``` to ```gw_size``` <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | --         | enable_nat        | **Yes**, if migrating to **spoke_gateway** resource; change any **spoke_vpc** resources' ```enable_nat``` to ```enable_snat``` <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | transit_vpc --> transit_gateway | vpc_size          | **Yes**, if migrating to **transit_gateway** resource; change any **transit_vpc** resources' ```vpc_size``` to ```gw_size``` <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | --         | enable_nat        | **Yes**, if migrating to **transit_gateway** resource; change any **transit_vpc** resources' ```enable_nat``` to ```enable_snat``` <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | tunnel     | vpc_name1, vpc_name2 | **Yes**; change ```vpc_name1``` and ```vpc_name2``` to ```gw_name1``` and ```gw_name2``` respectively |

### Boolean Standardization
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | account    | aws_iam           | **Yes**; Accepted values are changed to **true** or **false** booleans rather than string |
|(changed) | firewall   | base_log_enable   | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | --   | policy: log_enable | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | fqdn       | fqdn_status       | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | gateway    | enable_nat     | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | --         | vpn_access        | **Yes**; Accepted values are changed to **true** or **false** |
|(changed) | --         | enable_elb        | **Yes**; see above for details |
|(changed) | --         | split_tunnel      | **Yes**; see above for details |
|(changed) | --         | saml_enabled      | **Yes**; see above for details |
|(changed) | --         | enable_ldap       | **Yes**; see above for details |
|(changed) | --         | single_az_ha      | **Yes**; see above for details |
|(changed) | --         | allocate_new_eip  | **Yes**; see above for details |
|(changed) | site2cloud | ha_enabled        | **Yes**; see above for details |
|(changed) | spoke_vpc --> spoke_gateway | enable_nat | **Yes**, if migrating to **spoke_gateway** resource; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | --         | single_az_ha      | **Yes**, if migrating to **spoke_gateway** resource; Accepted values are changed to **true** or **false** <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | transit_vpc --> transit_gateway | enable_nat | **Yes**, if migrating to **transit_gateway** resource; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | --         | connected_transit | **Yes**, if migrating to **transit_gateway** resource; Accepted values are changed to **true** or **false** <br><br> Please see **Resource Renaming** table for full instructions |
|(changed) | tunnel     | enable_ha         | **Yes**; Accepted values are changed to **true** or **false** |


## R2.1 (UserConnect-4.7-patch) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | transit_gateway| allocate_new_eip, eip, ha_eip | **No**; Terraform support added to feature already available for regular gateway resource. By default, ``allocate_new_eip`` will be set to **true**. Only works for AWS |
|(new) | spoke_gateway  | allocate_new_eip, eip, ha_eip | **No**; see above for details |


## R2.3 (UserConnect-5.0) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | vgw_conn       | vgw_account, vgw_region | **Yes**; if customer has a **vgw_conn** resource, they must add these two attributes and their respective values to their config file and perform a ```terraform refresh```
|(new) | gateway        | insane_mode, insane_mode_az, peering_ha_insane_mode_az | **No**; Terraform added support for insane mode for Azure-related resources |
|(new) | spoke_gateway, transit_gateway | insane_mode, insane_mode_az, ha_insane_mode_az | **No**; see above for details |
|(new) | aws_tgw        | aviatrix_firewall, native_egress, native_firewall | **No**; Terraform support added to feature already available through Controller. These are new optional attributes for Terraform |
|(new) | spoke_gateway, transit_gateway | enable_active_mesh | **No**; Terraform support released alongside new feature available in Controller 5.0. New optional feature. Default value: **false** |


## R2.4 (UserConnect-5.0) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | firewall       | description       | **Yes**; Terraform support added to feature already available for stateful firewall policies. This attribute is optional by default. If customer already has a description set through the Controller, they must add the description attribute and its corresponding value to their Terraform file and perform a ``terraform refresh`` to rectify the diff  |
|(new) | account        | oci_tenancy_id, oci_user_id, oci_compartment_id, oci_api_private_key_filepath | **No**; Terraform now supports Oracle Cloud, hence the new attributes needed to create an Oracle access account |


## R2.5 (UserConnect-5.1) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | gateway, spoke_gateway, transit_gateway | enable_vpc_dns_server | **No**; Terraform added support for enabling/ disabling VPC DNS server. This attribute is optional, and set 'false' by default. This feature is only available on AWS |


## R2.6 (UserConnect-5.1) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | vgw_conn | enable_advertise_transit_cidr, bgp_manual_spoke_advertise_cidrs | **Yes**; functionality migrated over to **transit_gateway** resource. In order to maintain same functionality, customer must cut-paste the these two attributes and their respective values into the corresponding **transit_gateway** and perform a ```terraform refresh``` |
|(new) | transit_gateway | enable_advertise_transit_cidr, bgp_manual_spoke_advertise_cidrs | **No**; action required depends on above stated for **vgw_conn**. If customer does not have **vgw_conn** that originally advertised any sort of CIDR before this release, no action is required
|(changed) | -- | enable_firenet_interfaces | **Yes**; if customer is using **transit_gateway** and has set ``enable_firenet_interfaces``, attribute must be renamed to ``enable_firenet`` and a ``terraform refresh`` must be performed |
|(new) | -- | single_az_ha | **Yes**; Terraform support added for enabling/ disabling ``single_az_ha`` for **transit_gateway**. If customer has originally enabled ``single_az_ha`` through Controller prior to this release, then this attribute must be set to **true**, and a ``terraform refresh`` must be performed to rectify the state |


## R2.7 (UserConnect-5.1) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | vpc            | subnets           | **No**; Terraform support added for creating a GCP VPC. If customer wants to manage their existing GCP VPC through Terraform, they must write a matching configuration for the existing GCP VPC and perform a ``terraform import`` to bring it into the state |
|(deprecated) | --      | public_subnets, private_subnets | **Yes**; if customers have referenced these 2 attributes that were added in the previous release R2.6, they must change reference back to whichever ``subnets`` it corresponds to and perform a ``terraform refresh`` |
|(new) | saml_endpoint  | custom_saml_request_template | **No**; Terraform support added for using custom SAML templates for the endpoint |
|(new) | aws_tgw        | customized_routes, disable_local_route_propagation | **No**; Terraform support added for Controller 5.0 feature. Customers may now use customized routes and disable local route propagation when attaching a VPC to their TGW |
|(new) | gateway        | enable_vpn_nat    | **Yes**; Terraform support added for Controller 5.0 feature. If customers have enabled/ disabled VPN NAT feature through the Controller for their VPN gateway, they must specify this attribute and its corresponding value in Terraform and perform a ``terraform refresh`` |


## R2.8 (UserConnect-5.1, 5.2) (Terraform v0.12)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | account        | awsgov_account_number, awsgov_access_key, awsgov_secret_key | **No**; Terraform now supports AWS GovCloud accounts, hence the new attributes |
|(new) | aws_tgw_vpc_attachment | customized_routes, disable_local_route_propagation | **No**; Terraform support added for Controller 5.0 feature. Customers may now use customized routes and disable local route propagation when attaching a VPC to their TGW, managing them separately outside the TGW resource |
