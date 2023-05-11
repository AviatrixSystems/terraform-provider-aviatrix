---
layout: "aviatrix"
page_title: "Aviatrix R3.x Feature Changelist"
description: |-
  The Aviatrix provider R3.x Feature Changelist
---

# Aviatrix R3.x Feature Changelist

## USAGE:
Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R3.0 to present. New resources may be tracked in the [Release Notes](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-notes).

-> **NOTE:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

---

``Last updated: R3.1.0 (UserConnect-7.1)``


---

## R3.0.0 (UserConnect-7.0)
**NOTICE:** With the Release of Aviatrix v3.0.0, we have made more sweeping changes such as renaming/removal of attributes and resources. All these changes are all in the name of standardization of naming conventions and resources to match the Controller. Although we recognize that it may be a major inconvenience, we believe that these changes will benefit everyone in the long-term not only for clarity but ease of future feature-implementations and code maintenance.

For most changes, unless stated otherwise in the tables below, after editing the respective .tf files, a simple ``terraform refresh`` should rectify the state of the infrastructure.

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | app_domain | smart_group | **Yes**; **app_domain** will be deprecated and replaced with **smart_group**. Please use **smart_group** instead <br><br> You will need to remove any pre-existing **app_domain** resources from the statefile, update your .tf files to **smart_group**, and ``terraform import`` as **aviatrix_smart_group** <br><br> |
|(changed) | microseg_policy_list | distributed_firewalling_policy_list | **Yes**; **microseg_policy_list** will be replaced with **distributed_firewalling_policy_list**. Please use **distributed_firewalling_policy_list** instead <br><br> You will need to remove any pre-existing **microseg_policy_list** resources from the statefile, update your .tf files to **distributed_firewalling_policy_list** and ``terraform import`` <br><br> |

### Resource Deprecations

The following resources are removed:

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
| arm_peer | **Yes**; please use **aviatrix_azure_peer** instead |
| aws_tgw_security_domain | **Yes**; please use **aviatrix_aws_tgw_network_domain** instead. Please see the migration guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) for further instructions |
| aws_tgw_security_domain_connection | **Yes**; please use **aviatrix_aws_tgw_peering_network_domain** instead |
| edge_caag | **Yes**; please remove this from your config |
| segmentation_security_domain | **Yes**; please use **aviatrix_segmentation_network_domain** instead. Please see the migration guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) for further instructions |
| segmentation_security_domain_association | **Yes**; please use **aviatrix_segmentation_network_domain_association** instead. See above for instructions |
| segmentation_security_domain_connection_policy | **Yes**; please use **aviatrix_segmentation_network_domain_connection_policy** instead. See above for instructions |
| spoke_vpc | **Yes**; please use **aviatrix_spoke_gateway** instead |
| transit_vpc | **Yes**; please use **aviatrix_transit_gateway** instead |
| trans_peer | **Yes**; this resource will be removed in the following 3.0.1 provider release |

### Attribute Deprecations

The following attributes are removed:

| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | gateway_snat, gateway_dnat | sync_to_ha | **Yes**; please remove this attribute from the gateway_dnat/snat configs. SNAT/DNAT configs are now synced to HA gateways by default |
|(deprecated) | gateway, spoke_gateway, transit_gateway | tag_list | **Yes**; please migrate `tag_list` data values into a map-type format and use `tags` instead. Example: tags = {"key1" = "value1", "key2" = "value2"} |
|(deprecated) | spoke_gateway | manage_transit_gateway_attachment, transit_gw | **Yes**; please use the standalone **aviatrix_spoke_transit_attachment** resource instead |
|(deprecated) | firenet | manage_firewall_instance_association, firewall_instance_association | **Yes**; please use the standalone **aviatrix_firewall_instance_association** resource instead |
|(deprecated) | aws_tgw | manage_security_domain, security_domains, manage_vpc_attachment, attached_vpc, manage_transit_gateway_attachment, attached_transit_gateway | **Yes**; please see the [migration guide](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) if `security_domain` is still managed in-line. Once migrated to use network domains, please remove any `manage_x` attributes from the **aws_tgw** resource |
|(deprecated) | aws_tgw_connect, aws_tgw_direct_connect, aws_tgw_vpc_attachment | security_domain_name | **Yes**; please rename the attribute to `network_domain_name` instead. See the migration guide above for more information |


## R3.1.0 (UserConnect-7.1)
### Resource Deprecations

The following logging resources are removed:

| Resource       | Action Required?           |
|:--------------:|:--------------------------:|
| splunk_logging | **Yes**; please remove these resources from TF configuration |
| filebeat_forwarder | **Yes**; please remove these resources from TF configuration |
| sumologic_forwarder | **Yes**; please remove these resources from TF configuration |