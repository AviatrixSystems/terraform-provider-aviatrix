---
layout: "aviatrix"
page_title: "Aviatrix R2.x Feature Changelist"
description: |-
  The Aviatrix provider R2.x Feature Changelist
---

# Aviatrix R2.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R2.0 to present. New resources may be tracked in the [Release Notes](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-notes).

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

Otherwise, this list does not really apply. Please view the below list for details regarding this:

1. If a customer is transitioning to use Terraform to manage existing infrastructure, it is recommended to start from the latest Controller and Aviatrix Terraform provider version, and use the Terraform Export feature and import their infrastructure for a quick and easy migration.
  - Customer can still choose to manually write their config file to their own specifications and use ``terraform import`` to bring their infrastructure into Terraform state management
2. If a customer is adopting Terraform for the first time, with no pre-existing infrastructure/ for managing new infrastructure, this changelist does not apply whatsoever. They simply need to setup their Terraform environment and use the Terraform documentation featured on the left sidebar.

We **highly** recommend customers that are starting to adopt Terraform to manage their infrastructure, to start with our Release 2.0+ (R2.0+), which is compatible with our Controller 4.7+.

---

``Last updated: R2.23.0 (UserConnect-6.8)``


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


## R2.1 (UserConnect-4.7-patch)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | transit_gateway| allocate_new_eip, eip, ha_eip | **No**; Terraform support added to feature already available for regular gateway resource. By default, ``allocate_new_eip`` will be set to **true**. Only works for AWS |
|(new) | spoke_gateway  | allocate_new_eip, eip, ha_eip | **No**; see above for details |


## R2.3 (UserConnect-5.0)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | vgw_conn       | vgw_account, vgw_region | **Yes**; if customer has a **vgw_conn** resource, they must add these two attributes and their respective values to their config file and perform a ```terraform refresh```
|(new) | gateway        | insane_mode, insane_mode_az, peering_ha_insane_mode_az | **No**; Terraform added support for insane mode for Azure-related resources |
|(new) | spoke_gateway, transit_gateway | insane_mode, insane_mode_az, ha_insane_mode_az | **No**; see above for details |
|(new) | aws_tgw        | aviatrix_firewall, native_egress, native_firewall | **No**; Terraform support added to feature already available through Controller. These are new optional attributes for Terraform |
|(new) | spoke_gateway, transit_gateway | enable_active_mesh | **No**; Terraform support released alongside new feature available in Controller 5.0. New optional feature. Default value: **false** |


## R2.4 (UserConnect-5.0)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | firewall       | description       | **Yes**; Terraform support added to feature already available for stateful firewall policies. This attribute is optional by default. <br><br> If customer already has a description set through the Controller, they must add the description attribute and its corresponding value to their Terraform file and perform a ``terraform refresh`` to rectify the diff  |
|(new) | account        | oci_tenancy_id, oci_user_id, oci_compartment_id, oci_api_private_key_filepath | **No**; Terraform now supports Oracle Cloud, hence the new attributes needed to create an Oracle access account |


## R2.5 (UserConnect-5.1)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | gateway, spoke_gateway, transit_gateway | enable_vpc_dns_server | **No**; Terraform added support for enabling/ disabling VPC DNS server. <br><br> This attribute is optional, and set 'false' by default. This feature is only available on AWS |


## R2.6 (UserConnect-5.1)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | vgw_conn | enable_advertise_transit_cidr, bgp_manual_spoke_advertise_cidrs | **Yes**; functionality migrated over to **transit_gateway** resource. <br><br> In order to maintain same functionality, customer must cut-paste the these two attributes and their respective values into the corresponding **transit_gateway** and perform a ```terraform refresh``` |
|(new) | transit_gateway | enable_advertise_transit_cidr, bgp_manual_spoke_advertise_cidrs | **No**; action required depends on above stated for **vgw_conn**. <br><br> If customer does not have **vgw_conn** that originally advertised any sort of CIDR before this release, no action is required
|(changed) | -- | enable_firenet_interfaces | **Yes**; if customer is using **transit_gateway** and has set ``enable_firenet_interfaces``, attribute must be renamed to ``enable_firenet`` and a ``terraform refresh`` must be performed |
|(new) | -- | single_az_ha | **Yes**; Terraform support added for enabling/ disabling ``single_az_ha`` for **transit_gateway**. <br><br> If customer has originally enabled ``single_az_ha`` through Controller prior to this release, then this attribute must be set to **true**, and a ``terraform refresh`` must be performed to rectify the state |


## R2.7 (UserConnect-5.1)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | vpc            | subnets           | **No**; Terraform support added for creating a GCP VPC. <br><br> If customer wants to manage their existing GCP VPC through Terraform, they must write a matching configuration for the existing GCP VPC and perform a ``terraform import`` to bring it into the state |
|(deprecated) | --      | public_subnets, private_subnets | **Yes**; if customers have referenced these 2 attributes that were added in the previous release **R2.6**, they must change reference back to whichever ``subnets`` it corresponds to and perform a ``terraform refresh`` |
|(new) | saml_endpoint  | custom_saml_request_template | **No**; Terraform support added for using custom SAML templates for the endpoint |
|(new) | aws_tgw        | customized_routes, disable_local_route_propagation | **No**; Terraform support added for Controller 5.0 feature. <br><br> Customers may now use customized routes and disable local route propagation when attaching a VPC to their TGW |
|(new) | gateway        | enable_vpn_nat    | **Yes**; Terraform support added for Controller 5.0 feature. <br><br> If customers have enabled/ disabled VPN NAT feature through the Controller for their VPN gateway, they must specify this attribute and its corresponding value in Terraform and perform a ``terraform refresh`` |


## R2.8 (UserConnect-5.1, 5.2)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | account        | awsgov_account_number, awsgov_access_key, awsgov_secret_key | **No**; Terraform now supports AWS GovCloud accounts, hence the new attributes |
|(new) | aws_tgw_vpc_attachment | customized_routes, disable_local_route_propagation | **No**; Terraform support added for Controller 5.0 feature. <br><br> Customers may now use customized routes and disable local route propagation when attaching a VPC to their TGW, managing them separately outside the TGW resource |


## R2.9 (UserConnect-5.2)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | gateway        | enable_designated_gateway, additional_cidrs_designated_gateway | **No**; Terraform now supports the "designated gateway" feature. <br><br> If customer wants to manage their existing gateway with this feature enabled, they must write a matching configuration for the existing gateway and perform a ``terraform import`` to bring it into the state |
|(new) | --             | enable_encrypt_volume | **Yes**; Terraform now supports encrypting EBS volumes. <br><br> If an existing gateway in Terraform state has had their EBS volume encrypted through the GUI, this attribute must be set to **true**, and a ``terraform refresh`` must be performed to rectify the state |
|(new) | --             | dnat_policy       | **Yes**; Terraform now supports DNAT. <br><br> If an existing gateway in Terraform state has had DNAT policies set through the GUI, corresponding ``dnat_policy`` blocks must be added and a ``terraform refresh`` must be performed to rectify the state |
|(new) | spoke_gateway  | snat_policy       | **Yes**; Terraform now supports custom SNAT or multi-IP SNAT. <br><br> If an existing spoke gateway in Terraform state has had SNAT policies set through the GUI, corresponding ``snat_policy`` blocks must be added and a ``terraform refresh`` must be performed to rectify the state |


## R2.10 (UserConnect-5.2.2122+)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | gateway, spoke_gateway, transit_gateway | enable_snat | **Yes**; renamed to ``single_ip_snat``. <br><br> Update the attribute in .tf files and perform a ``terraform refresh`` |
|(deprecated) | --      | snat_mode, snat_policy {}, dnat_policy {} | **Yes**; customized SNAT and DNAT support has been deprecated from **R2.9** and functionality has been moved to **aviatrix_gateway_snat** and **aviatrix_gateway_dnat** respectively, to improve policy management. <br><br> If customized SNAT/DNAT policies were created in **R2.9**, copy the existing policy block and paste into the new corresponding resource. Specify the corresponding gateway name value for ``gw_name`` and change the attributes that were renamed as necessary (ex. ``src_ip`` to ``src_cidr``). Full documentation for the resources may be viewed under the "Gateway" tab


## R2.11 (UserConnect-5.3.1391)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | aws_tgw_vpn_conn | connection_type | **Yes**; New attribute added for distinction between two possible TGW VPN connections. <br><br> If an existing aws_tgw_vpn_conn (static VPN connection) has been created prior to this release, customers must add ``connection_type = static`` into their configuration file and perform a ``terraform refresh`` to update and apply the attribute's value (static) into the state file.
|(new) | gateway        | vpn_protocol      | **No**; Terraform now supports specifying the VPN protocol to use for the VPN gateway/ELB. <br><br> Customers managing existing gateways and ELBs do not need to change their configurations, but are recommended for future-created infrastructure to specify a protocol for full control as well as clarity |


## R2.12 (UserConnect-5.3.1491)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | arm_peer | --               | **Yes**; **arm_peer** will be deprecated and replaced with **azure_peer**. It will be kept for now for backward-compatibility and will be removed in the future. Please use **azure_peer** instead <br><br> You will need to remove any pre-existing **arm_peer** resources from the statefile, update your .tf files to **azure_peer**, and perform a ``terraform import`` |
|(changed) | gateway    | cloudn_bkup_gateway_inst_id | **No**; the following computed attribute will changed to ``peering_ha_cloud_instance_id``. The next ``terraform refresh`` will automatically rectify the state with the new name |


## R2.13 (UserConnect-5.4.1060)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | account_user | account_name | **Yes**; ``account_name`` has been deprecated due to the new Role-based Access Control feature (RBAC) implementation in 5.4. This attribute must be removed from the .tf files and a ``terraform refresh`` will rectify the state file. |
|(new) | aws_tgw        | manage_transit_gateway_attachment | **No**; New attribute flag added to support managing transit gateway attachments to the AWS TGW outside of the resource, using the **aviatrix_aws_tgw_transit_gateway_attachment**. A simple ``terraform refresh`` will apply this attribute's default value (true) and rectify the state |


## R2.14 (UserConnect-5.4.1201)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | gateway | public_ip, peering_ha_public_ip | **Yes**; these attributes have been deprecated in favor of `eip` and `peering_ha_eip` to maintain consistency between gateway types. If any Terraform configuration of other resources reference this attribute, replace them with the aforementioned attributes respectively. |
|(new) | --             | peering_ha_gw_name, peering_ha_private_ip| **No**; these are two new attributes that will be exported as output values in the resource, and can be referenced as necessary in other resources or modules. |
|(new) | spoke_gateway, transit_gateway | ha_gw_name, ha_private_ip | **No**; see above for details |


## R2.15 (UserConnect-6.0)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | site2cloud | tunnel_type       | **Yes**; with the release of Controller 6.0, **site2cloud** has been changed to now be policy-based or route-based. With that change, the new accepted values are "policy" and "route". These values must be updated in their respective .tf files and a ``terraform refresh`` will rectify the state file. |
|(changed) | gateway    | enable_encrypt_volume | **No**; nothing needs to be changed. This is just a note to bring to attention that with gateways created by Terraform, volumes will NOT be encrypted, unless otherwise specified. Please see Release Notes for details |
|(changed) | gateway_snat | connection      | **Yes**; any **gateway_snat** configurations where the `connection` value is specified as an empty string must be changed to either null or "eth0" due to new default values. |


## R2.16 (UserConnect-6.1)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | transit_gateway_peering | gateway1_excluded_network_cidrs, gateway1_excluded_tgw_connections, gateway2_excluded_network_cidrs, gateway2_excluded_tgw_connections | **No**; these are new attributes to support additional configurations for transit gateway peerings. A simple ``terraform refresh`` will apply this attributes' default values into the state |
|(new) | aviatrix_transit_gateway | bgp_polling_time, prepend_as_path, local_as_number, bgp_ecmp | **No**; see reason above |
|(new) | controller_config | enable_vpc_dns_server | **No**; see reason above |
|(new) | gateway_dnat, gateway_snat | sync_to_ha | **No**; see reason above |


## R2.16.2 (UserConnect-6.1.1280)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | aws_tgw    | connected_domains | **No**; nothing needs to be changed. This is a note to bring to attention that in Controller 6.1.1276+, if there are peered TGWs and **aws_tgw_peering_domain_conn**s in one's Terraform topology, users must upgrade to R2.16.2. Please see Release Notes for details |


## R2.16.3 (UserConnect-6.1.1309)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | spoke_gateway | transit_gw     | **No**; nothing needs to be changed. This is a note to bring to attention new support for allowing multiple transit gateway attachments, in the same format of a comma-separated string |


## R2.17 (UserConnect-6.2)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | fqdn           | manage_domain_names | **No**; nothing needs to be changed. This is a new attribute introduced that allows decoupling of FQDN rules management out of the resource. A simple refresh will set this attribute into the state as default true |
|(new) | firewall       | manage_firewall_policies | **No**; nothing needs to be changed. This is a new attribute introduced that allows decoupling of Stateful firewall rules management out of the resource. A simple refresh will set this attribute into the state as default true |
|(new) | spoke_gateway  | manage_transit_gateway_attachment | **No**; nothing needs to be changed. This is a new attribute introduced that allows decoupling of spoke to transit attachments out of the resource. A simple refresh will set this attribute into the state as default true |


## R2.17.2 (UserConnect-6.2.1914)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | transit_gateway_peering | gateway1_excluded_network_cidrs, gateway1_excluded_tgw_connections, gateway2_excluded_network_cidrs, gateway2_excluded_tgw_connections | **Yes**; if customers are experiencing deltas in their state in any of the aforementioned attributes due to Terraform re-ordering prior to R2.17.2, the workaround is to reorder the list in the config to match what the returned 'expected' output is and run a `terraform apply`. The fix in R2.17.2 will correct this behavior for any **new** future peering resources created. Note that if manual re-ordering is not a viable option, one may also destroy and re-create the resource as necessary |


## R2.18 (UserConnect-6.3)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | gateway    | monitor_exclude_list | **Yes**; attribute has changed from string-type to a list. If this attribute is already being used in a comma-separated string, please convert values to a list of strings. Eg. ``"foo,bar"`` to ``["foo", "bar"]`` |


## R2.18.1 (UserConnect-6.3.2364)
**NOTICE:** Starting in R2.18.1, as part of the initiative to improve future support, development and performance of the provider, certain attributes will be deprecated in order to avoid cross-resource behavior dependencies and improve scalability. In **R3.0**, we will fully remove these attributes.

Note there are standalone resources already in place to be used and one only needs to set the specified `manage_x` attribute to false and import the specific resource(s) into the state. For more information on import, please see the Hashicorp import docs [here](https://www.terraform.io/docs/cli/commands/import.html)

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | aws_tgw | attached_vpc      | **Yes**; Please set `manage_vpc_attachment` to false, and use the standalone **aviatrix_aws_tgw_vpc_attachment** resource instead. |
|(deprecated) | firenet | firewall_instance_association | **Yes**; Please set `manage_firewall_instance_association` to false, and use the standalone **aviatrix_firewall_instance_association** resource instead. |
|(deprecated) | firewall | policy           | **Yes**; Please set `manage_firewall_policies` to false, and use the standalone **aviatrix_firewall_policy** resource instead. |
|(deprecated) | fqdn    | domain_names      | **Yes**; Please set `manage_domain_names` to false, and use the standalone **aviatrix_fqdn_tag_rule** resource instead. |
|(deprecated) | spoke_gateway | transit_gw  | **Yes**; Please set `manage_transit_gateway_attachment` to false, and use the standalone **aviatrix_spoke_transit_attachment** resource instead. |


## R2.19.0 (UserConnect-6.4)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | aws_tgw | security_domains, security_domain_name, connected_domains, aviatrix_firewall, native_egress, native_firewall | **Yes**; Please set `manage_security_domain` to false, and use the standalone **aviatrix_aws_tgw_security_domain** and **aviatrix_aws_tgw_security_domain_connection** resources instead. |
|(deprecated) | gateway, spoke_gateway, transit_gateway | tag_list | **Yes**; Please migrate `tag_list` data values into a map-type format and use `tags` instead. Example: tags = {"key1" = "value1", "key2" = "value2"} |
|(new) | remote_syslog  | name              | **Yes**; this is a new required attribute introduced in Controller 6.4. For customers who are already managing Remote Syslog in Terraform, please update the config files to include this, and a simple ``terraform apply`` will fix this issue. |


## R2.19.3 (UserConnect-6.4.2672)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(new) | gateway        | availability_domain, fault_domain, peering_ha_availability_domain, peering_ha_fault_domain | **Yes**; while existing OCI infrastructure will not be affected due to attribute being computed value, future OCI gateways are required to launch with an `availability_domain` and `fault_domain`. The `peering_ha` variants will be required if the gateways are to HA-enabled. |
|(new) | spoke_gateway, transit_gateway | availability_domain, fault_domain, ha_availability_domain, ha_fault_domain | **Yes**; while existing OCI infrastructure will not be affected as described above, all future OCI spokes and transit are required to launch with `availability_domain` and `fault_domain` specified. The `ha` variants will be required if the spokes/transits are to be HA-enabled. |
|(new) | firewall_instance | availability_domain, fault_domain | **Yes**; while existing OCI infrastructure will not be affected as described above, all future OCI firewall instances are required to launch an `availability_domain` and `fault_domain`. |


## R2.20 (UserConnect-6.5)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | vpc        | vpc_id            | **Yes**; while existing Azure VNets created will not be affected, all future VNets created in 6.5/R2.20 will be in the following new 3-tuple format `<vnet-name>:<resource-group-name>:<GUID>`. Any resources such as **gateways** that reference this attribute will be unaffected. |
|(changed) | external_device_conn | remote_vpc_name | **Yes**; while existing **external_device_conn** resources created prior to 6.5/R2.20 will not experience any changes, all future **external_device_conn** must be created with the following new 3-tuple format `<vnet-name>:<resource-group-name>:<subscription-id>`. |


## R2.20.1 (UserConnect-6.5.2608)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | spoke_gateway, transit_gateway | enable_active_mesh | **Yes**; if customers are using spoke or transit gateways with `enable_active_mesh` set to **false**, action is required. Non-ActiveMesh features will no longer be supported in R2.21.0. Please follow the guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_to_active_mesh_transit_network) to migrate from Classic Aviatrix Encrypted Transit Network to Aviatrix ActiveMesh Transit Network. |
|(deprecated) | controller_config | sg_management_account_name, security_group_management | **Yes**; if Security Group Management was enabled through the **aviatrix_controller_config** resource, these attributes must be removed from the config file and a ``terraform refresh`` must be performed. Please use the standalone **aviatrix_controller_security_group_management_config** resource instead and perform a ``terraform import``. |


## R2.21.0 (UserConnect-6.6)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | spoke_gateway, transit_gateway | enable_active_mesh | **Yes**; this attribute is now fully deprecated. Please see above note in the R2.20.1 table for details on how to migrate to ActiveMesh 2.0. If the gateways are already ActiveMesh-enabled, this attribute may be safely removed from configuration, and a ``terraform refresh`` may be run. |
|(deprecated) | spoke_gateway, transit_gateway | storage_name | **Yes**; this attribute is fully deprecated. If there are existing AzureChina gateways, this attribute may be safely removed from configuration, and a ``terraform refresh`` may be run. New AzureChina gateways will no longer require this attribute. |


## R2.21.1 (UserConnect-6.6.5404)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | rbac_group_permission_attachment | permission_name | **Yes**; one of values of this attribute has been renamed:  `all_cloudwan_write` to `all_cloudn_write`. If this permission was set previously, it will be deleted and must be remade. The permission can be safely renamed in the configuration and re-applied. |


## R2.21.2 (UserConnect-6.6.5544)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | gateway <br> spoke_gateway <br> transit_gateway <br> firenet <br> firewall_instance <br> site2cloud <br> transit_external_device_conn <br> spoke_external_device_conn <br> transit_cloudn_conn <br> vpn_user  | `vpc_id` | **Yes**; any of the listed resources built on OCI that references the `name` attribute from the **aviatrix_vpc** resource, should be changed to the `vpc_id` attribute. Running ``terraform refresh`` after this config change is made, will rectify any deltas in the state. This change was made to standardize the behavior between all CSPs (clouds). |


## R2.22.0 (UserConnect-6.7)
### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | aviatrix_aws_tgw_security_domain  | aviatrix_aws_tgw_network_domain  | **Yes**; **aws_tgw_security_domain** will be deprecated and replaced with **aws_tgw_network_domain**. It will be kept for now for backward compatibility and will be removed in the future. Please use **aws_tgw_network_domain** instead <br><br> You will need to remove any pre-existing **aws_tgw_security_domain** resources from the statefile, update your .tf files to **aws_tgw_network_domain**, as well as any attributes (and corresponding values) as necessary, and ``terraform import`` as **aviatrix_aws_tgw_network_domain** <br>|
|(deprecated) | segmentation_security_domain | segmentation_network_domain | **Yes**; please see above for details or follow the guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) for migration steps |
|(deprecated) | segmentation_security_domain_association | segmentation_network_domain_association | **Yes**; please see above |
|(deprecated) | segmentation_security_domain_connection_policy | segmentation_network_domain_connection_policy | **Yes**; please see above |

### Attribute Renaming
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | aws_tgw_connect <br> aws_tgw_directconnect <br> aws_tgw_vpc_attachment | security_domain_name | **Yes**; renamed to ``network_domain_name``. Please follow the guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) for migration steps |

### Misc.
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | gateway <br> spoke_gateway <br> transit_gateway <br> vpc | vpc_id | **No**; the attribute value for these resources created in GCP now include the project ID. The new format is `<vpc_name>~-~<project_id>`; if some resources referenced the `vpc_id` and concatenated the project ID, it is no longer necessary |
|(deprecated) | device_registration <br> device_tag <br> device_transit_gateway_attachment <br> device_aws_tgw_attachment <br> device_virtual_wan_attachment | -- | **Yes**; these resources are removed from Terraform as CloudWAN is no longer a supported feature |
|(deprecated) | cloudn_transit_gateway_attachment | enable_dead_peer_detection <br> enable_learned_cidrs_approval <br> approved_cirs | **Yes**; these attributes are no longer supported in this resource. They can be safely removed from .tf config |


## R2.23.0 (UserConnect-6.8)
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | firenet | fail_close_enabled | **Yes**; if customers are using FireNet with `fail_close_enabled` set to **true/false**, action is required. This attribute is fully deprecated as Fail Close will be enabled by default in Controller 6.8+. This attribute may be safely removed from configuration, and a ``terraform refresh`` may be run. |