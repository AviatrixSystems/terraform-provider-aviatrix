---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 9.0.0
### Notes:
- Supported Controller version: **9.0.0**

### Enhancements:
| Enhancement | Description |
| :--- | :--- |
| AVX-65775 | Added new **aviatrix_spoke_gateway_group** and **aviatrix_transit_gateway_group** resources for creating and managing Spoke and Transit Gateway Groups for CSP and Edge gateways, identified by `group_uuid` |
| AVX-65776 | Added new **aviatrix_spoke_instance** resource for creating Spoke gateways linked to an existing gateway group via `group_uuid`, supporting CSP and Edge deployments |
| AVX-67518 | Added new **aviatrix_transit_instance** resource for creating Transit gateways linked to an existing gateway group via `group_uuid`, supporting CSP and Edge deployments |
| AVX-69545 | Added support for enabling and disabling DCF controller features via the **aviatrix_controller_config_feature** resource |
| AVX-70660 | Added support for setting the default IPS Profile in DCF via the **aviatrix_dcf_ips_default_profile** resource |
| AVX-72121 | Added IPv6 support for Equinix Edge as a Service (EAS) gateways via **aviatrix_edge_equinix**, enabling dual-stack deployments |
| AVX-72330 | Added `proxy_id` argument to **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn** for Site2Cloud proxy identity configuration |
| AVX-72382 | Added new **aviatrix_dcf_tls_profile** data source to look up a DCF TLS Profile by name and retrieve its UUID for use in policy references |
| AVX-72741 | Added `private_route_table_config` argument to **aviatrix_vpc**, **aviatrix_spoke_gateway**, **aviatrix_transit_gateway**, **aviatrix_spoke_gateway_group**, and **aviatrix_transit_gateway_group** for Azure SNAT private route table configuration |
| AVX-72896 | Added **aviatrix_dcf_mitm_ca** resource and data source, and **aviatrix_dcf_mitm_ca_selection** resource for MITM CA management |
| AVX-73054 | Added `config_mode` and `firewall_template_config` arguments to **aviatrix_firenet_firewall_manager** for per-firewall template configuration in advanced FireNet deployments |
| AVX-73155 | Updated **aviatrix_controller_config_feature** to dynamically validate feature names against the controller API, replacing the previously hardcoded list |
| AVX-73936 | Added `never_drop_sids` argument to **aviatrix_dcf_ips_default_profile** for configuring SID exceptions in the DCF IPS default profile |

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-67977 | Fixed **aviatrix_spoke_transit_attachment** not being written to state after an attachment timeout; running `terraform refresh` now correctly recovers state without requiring manual deletion and re-attachment |
| AVX-69674 | Fixed `terraform destroy` failing on **aviatrix_edge_megaport_ha** when `ztp_file_download_path` was null in state, causing an "invalid or missing value" error |
| AVX-69679 | Fixed `rx_queue_size` on Edge gateway resources to be applied after the gateway reaches UP state, avoiding RPC errors during initial provisioning |
| AVX-70531 | Fixed `enable_ipv6` on **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn** to be `ForceNew`, as in-place IPv6 edits on BGP connections are not supported by the controller |
| AVX-72574 | Fixed **aviatrix_site2cloud** custom mapped NAT CIDR updates (`remote_destination_real_cidrs`, `remote_destination_virtual_cidrs`, and related fields) not taking effect in the controller when modified in-place; changes now propagate correctly without requiring resource recreation |
| AVX-72621 | Fixed misleading error message on **aviatrix_transit_gateway** that incorrectly referred to a "spoke gateway" when IPv6 enablement failed due to no IPv6 CIDR on the VPC |
| AVX-72822 | Fixed **aviatrix_spoke_transit_attachment** failing when a spoke or transit was already part of a gateway group, where the group name and gateway name differ |
| AVX-72968 | Fixed **aviatrix_dcf_ruleset** recalculating and re-applying all rules when only a single rule was added or removed, causing excessive diffs |
| AVX-73557 | Fixed spoke-transit attachment not being propagated to controller internal state, which caused downstream failures in FireNet management access and SNAT configuration |
| AVX-74429 | Fixed **aviatrix_vpc** forcing resource replacement on upgrade from 8.2 to 9.0 due to `ipv6_access_type` defaulting to `EXTERNAL` on non-GCP VPCs where the field is not applicable |
| AVX-75065 | Fixed **aviatrix_transit_gateway** ignoring `enable_jumbo_frame = false` during creation, causing the gateway to be created with jumbo frame enabled and producing configuration drift on subsequent plan |
| AVX-75145 | Fixed an issue where using the Terraform attachment point data source with the names `TERRAFORM_AFTER_UI_MANAGED` or `TERRAFORM_BEFORE_UI_MANAGED`, before Distributed Cloud Firewall (DCF) was enabled, stored incorrect UUIDs for these system attachment points. After DCF was enabled, any rulesets or policy groups attached to those attachment points were orphaned and the rules within them were not enforced. The data source APIs now check whether DCF is enabled and fail when it is not, preventing the misconfigured state. Users who already encountered this in an earlier version can recover by clearing all DCF policies, disabling DCF, and re-enabling it before using the system attachment points again. |
| AVX-75634 | Fixed `ha_ipv6_ip` on **aviatrix_spoke_gateway** producing a persistent diff on re-apply when the HA IPv6 IP was not explicitly set |
| AVX-76197 | Fixed **aviatrix_dcf_ruleset** regenerating rule UUIDs when rule attributes were modified, causing unexpected rule recreation |

### Deprecations:
| Issue | Description |
| :--- | :--- |
| AVX-70018 | Removed **aviatrix_gateway_certificate_config**; the underlying gateway CA upload API was removed from the controller |
| AVX-73294 | Deprecated **aviatrix_distributed_firewalling_proxy_ca_config** and **aviatrix_distributed_firewalling_origin_cert_enforcement_config**; use **aviatrix_dcf_mitm_ca** and **aviatrix_dcf_tls_profile** instead |
| AVX-73772 | Removed **aviatrix_dcf_ips_profile_vpc**; IPS Profile is now configured globally via **aviatrix_dcf_ips_default_profile** |
