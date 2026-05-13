---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 8.2.10
### Notes:
- Supported Controller version: **8.2.10**

### Enhancements:
| Enhancement | Description |
| :--- | :--- |
| AVX-69545 | Added support for enabling and disabling DCF controller features via the **aviatrix_controller_config_feature** resource |
| AVX-70660 | Added support for setting the default IPS Profile in DCF via the **aviatrix_dcf_ips_default_profile** resource |
| AVX-72382 | Added new **aviatrix_dcf_tls_profile** data source to look up a DCF TLS Profile by name and retrieve its UUID for use in policy references |
| AVX-73936 | Added `never_drop_sids` argument to **aviatrix_dcf_ips_default_profile** for configuring SID exceptions in the DCF IPS default profile |

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-72574 | Fixed **aviatrix_site2cloud** custom mapped NAT CIDR updates (`remote_destination_real_cidrs`, `remote_destination_virtual_cidrs`, and related fields) not taking effect in the controller when modified in-place; changes now propagate correctly without requiring resource recreation |
| AVX-72621 | Fixed misleading error message on **aviatrix_transit_gateway** that incorrectly referred to a "spoke gateway" when IPv6 enablement failed due to no IPv6 CIDR on the VPC |
| AVX-72822 | Fixed **aviatrix_spoke_transit_attachment** failing when a spoke or transit was already part of a gateway group, where the group name and gateway name differ |
| AVX-72968 | Fixed **aviatrix_dcf_ruleset** recalculating and re-applying all rules when only a single rule was added or removed, causing excessive diffs |
| AVX-75065 | Fixed **aviatrix_transit_gateway** ignoring `enable_jumbo_frame = false` during creation, causing the gateway to be created with jumbo frame enabled and producing configuration drift on subsequent plan |
| AVX-75145 | Fixed an issue where using the Terraform attachment point data source with the names `TERRAFORM_AFTER_UI_MANAGED` or `TERRAFORM_BEFORE_UI_MANAGED`, before Distributed Cloud Firewall (DCF) was enabled, stored incorrect UUIDs for these system attachment points. After DCF was enabled, any rulesets or policy groups attached to those attachment points were orphaned and the rules within them were not enforced. The data source APIs now check whether DCF is enabled and fail when it is not, preventing the misconfigured state. Users who already encountered this in an earlier version can recover by clearing all DCF policies, disabling DCF, and re-enabling it before using the system attachment points again. |
| AVX-76197 | Fixed **aviatrix_dcf_ruleset** regenerating rule UUIDs when rule attributes were modified, causing unexpected rule recreation |

### Deprecations:
| Issue | Description |
| :--- | :--- |
| AVX-70018 | Removed **aviatrix_gateway_certificate_config**; the underlying gateway CA upload API was removed from the controller |
| AVX-73294 | Deprecated **aviatrix_distributed_firewalling_proxy_ca_config** and **aviatrix_distributed_firewalling_origin_cert_enforcement_config**; use **aviatrix_dcf_mitm_ca** and **aviatrix_dcf_tls_profile** instead |
| AVX-73772 | Removed **aviatrix_dcf_ips_profile_vpc**; IPS Profile is now configured globally via **aviatrix_dcf_ips_default_profile** |
