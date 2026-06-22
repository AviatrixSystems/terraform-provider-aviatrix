---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 9.0.10
### Notes:
- Supported Controller version: **9.0.10**

### Enhancements:
| Enhancement | Description |
| :--- | :--- |
| AVX-76307 | Added support for per-cluster feature flags to **`aviatrix_kubernetes_cluster`** |

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-52485 | Fixed an issue where destroying **`aviatrix_spoke_gateway`** failed when the HA gateway had not yet been deleted, due to a missing deletion dependency |
| AVX-67977 | Fixed an issue where **`aviatrix_spoke_transit_attachment`** state was not updated during `terraform refresh` if the attachment had previously timed out |
| AVX-72968 | Fixed an issue where **`aviatrix_dcf_ruleset`** recalculated all rules in the diff when a single rule was added or removed |
| AVX-75065 | Fixed an issue where **`aviatrix_transit_gateway`** ignored `enable_jumbo_frame = false` during creation, causing configuration drift on subsequent plans |
| AVX-75190 | Fixed an issue where **`aviatrix_transit_instance`** did not support `insane_mode` and related HPE attributes |
| AVX-75206 | Fixed an issue where **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`** rejected the GCP full zone format (e.g., `us-east1-a`) |
| AVX-75210 | Fixed issues where **`aviatrix_spoke_group`** failed to update `group_instance_size`, and `gw_size` was incorrectly required on **`aviatrix_spoke_instance`** when the group already defines a size |
| AVX-75493 | Fixed an issue where the `zone` attribute was ignored when creating GCP gateway instances via **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`**, causing the gateway to be created in the wrong zone |
| AVX-75494 | Fixed an issue where `private_route_table_config` was missing from the **`aviatrix_spoke_gateway`** data source output |
| AVX-75634 | Fixed an issue where `ha_ipv6_ip` caused a persistent plan diff on **`aviatrix_spoke_gateway`** after apply |
| AVX-75778 | Fixed an issue where **`aviatrix_transit_gateway`** and **`aviatrix_spoke_gateway`** returned an error for `accept_bgp_med_to_sdn_metric` on controllers that do not have the feature enabled |
| AVX-76197 | Fixed an issue where updating a rule attribute in **`aviatrix_dcf_ruleset`** caused the rule UUID to change unexpectedly |
| AVX-76340 | Fixed an issue where Equinix and Megaport edge gateway instance creation via **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`** failed with an incorrect gateway size error |
| AVX-76494 | Fixed an issue where **`aviatrix_dcf_ruleset`** allowed creation of multiple rules with the same priority |
| AVX-76528 | Fixed an issue where Terraform operations on **`aviatrix_account`** failed with a JSON decode error when the controller response was truncated for large account lists |
