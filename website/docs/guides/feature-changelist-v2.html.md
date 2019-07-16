# Terraform Feature Changelist (v2)
**USAGE:** Tracks customer-impacting changes to Terraform environment (existing resources) throughout releases from R2.0 to present. New resources may be tracked in the Release Notes.

**Note:** This changelist assumes that customers have an existing Terraform configuration and is planning to upgrade their Controller (and subsequently must also update their Aviatrix Terraform provider version to the appropriate corresponding one).

Otherwise, this list does not really apply. Please view the below list for details regarding this:
1. If a customer is transitioning to use Terraform to manage existing infrastructure, it is recommended to start from the latest Controller and Aviatrix Terraform provider version, and use the Terraform Export feature and import their infrastructure for a quick and easy migration.
  - **Please note that "Export" is still a beta feature (and only up-to-date for 4.3 at the moment)**
  - Customer can still choose to manually write their config file to their own specifications and use ``terraform import`` to bring their infrastructure into Terraform state management
2. If a customer is adopting Terraform for the first time, clean slate, this changelist does not apply whatsoever. They simply need to follow the appropriate doc to setup their Terraform environment and use the Terraform documentation corresponding to their Controller version to begin using Terraform

---

``Last updated: R2.0 (UserConnect-4.7-patch)``


---

## R2.0 (UserConnect-4.7.494) (Terraform v0.12)
**NOTICE:** With the Release of Aviatrix v2, there is major restructuring of our code as well as major changes such as renaming of attributes, resources, and attribute values. All these changes are all in the name of standardization  of naming conventions and resources. Although we recognize that it is a major inconvenience to customers, we believe that these changes will benefit everyone in the long-term not only for customer clarity but ease of future feature-implementations and code maintenance.

### Attribute Renaming
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | firewall   | base_allow_deny   | **Yes**; renamed to ```base_policy``` |
|(changed) | --   | base_log_enable   | **Yes**; renamed to ```base_log_enabled``` |
|(changed) | --   | policy: allow_deny | **Yes**; renamed to ```action``` |
|(changed) | --   | policy: log_enable | **Yes**; renamed to ```log_enabled``` |
|(changed) | fqdn       | fqdn_status       | **Yes**; renamed to ```fqdn_enabled``` |
|(changed) | gateway    | vpc_size          | **Yes**; change any gateway's ```vpc_size``` to ```gw_size``` |
|(changed) | gateway    | vpc_net           | **Yes**; change any gateway's ```vpc_net``` to ```subnet``` |
|(changed) | spoke_vpc  | vpc_size          | **Yes**; change any spoke gateway's ```vpc_size``` to ```gw_size``` |
|(changed) | transit_vpc| vpc_size          | **Yes**; change any transit gateway's ```vpc_size``` to ```gw_size``` |
|(changed) | tunnel     | vpc_name1, vpc_name2 | **Yes**; change ```vpc_name1``` and ```vpc_name2``` to ```gw_name1``` and ```gw_name2``` respectively |

### Resource Renaming
| Diff | Resource       | New Resource Name | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(deprecated) | spoke_vpc  | spoke_gateway  | **Yes**; **spoke_vpc** will be deprecated and replaced with **spoke_gateway**. It will be kept for now for backward compatibility and will be removed in the future. Please use **spoke_gateway** instead. You will need to remove it from statefile and import as **aviatrix_spoke_gateway** if it is already in state |
|(deprecated) | transit_vpc | transit_gateway | **Yes**; **transit_vpc** will be deprecated and replaced with **transit_gateway**. It will be kept for now for backward compatibility and will be removed in the future. Please use **transit_gateway** instead. You will need to remove it from statefile and import as **aviatrix_transit_gateway** if it is already in state |

### Boolean Standardization
| Diff | Resource       | Attribute         | Action Required?           |
|:----:|----------------|:-----------------:|----------------------------|
|(changed) | firewall   | base_log_enable   | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | --   | policy: log_enable | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | fqdn       | fqdn_status       | **Yes**; Accepted values are changed to **true** or **false**. Please note this attribute was also renamed. See **Attribute Renaming** table above |
|(changed) | gateway    | enable_nat     | **Yes**; Accepted values are changed to **true** or **false** |
|(changed) | --         | vpn_access        | **Yes**; see above for details |
|(changed) | --         | enable_elb        | **Yes**; see above for details |
|(changed) | --         | split_tunnel      | **Yes**; see above for details |
|(changed) | --         | saml_enabled      | **Yes**; see above for details |
|(changed) | --         | enable_ldap       | **Yes**; see above for details |
|(changed) | --         | single_az_ha      | **Yes**; see above for details |
|(changed) | --         | allocate_new_eip  | **Yes**; see above for details |
|(changed) | site2cloud | ha_enabled        | **Yes**; see above for details |
|(changed) | spoke_vpc / spoke_gateway | enable_nat | **Yes**; see above for details |
|(changed) | --         | single_az_ha      | **Yes**; see above for details |
|(changed) | transit_vpc / transit_gateway | enable_nat | **Yes**; see above for details |
|(changed) | --         | connected_transit | **Yes**; see above for details |
|(changed) | tunnel     | enable_ha         | **Yes**; see above for details |
