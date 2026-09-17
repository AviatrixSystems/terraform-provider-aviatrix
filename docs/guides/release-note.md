---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 10.1.0
### Notes:
- Supported Controller version: **10.1.0**

### Breaking Changes:
| Change | Description |
| :--- | :--- |
| AVX-78304 | The `single_az_ha` default for the **`aviatrix_gateway`** resource changed from `false` to `true` to match the Controller's behavior of enabling Single-AZ HA on every new gateway. Configurations that relied on the old `false` default without setting `single_az_ha` explicitly will now create gateways with Single-AZ HA enabled. To preserve the previous behavior, set `single_az_ha = false` explicitly. |

### Enhancements:
| Enhancement | Description |
| :--- | :--- |
| AVX-32589 | Added support for updating route table associations on **`aviatrix_spoke_gateway`** without requiring gateway detachment and reattachment |
| AVX-49318 | Added `enforcement` attribute to rules in **`aviatrix_dcf_ruleset`**, supporting three states: `enforce` (fully applied), `monitor` (matches traffic but does not enforce the action), and `disable` (stored only, never pushed to gateways). The existing `watch` attribute is now deprecated — see Deprecations. |
| AVX-58362 | Added `enable_bfd` support to **`aviatrix_transit_external_device_conn`**, **`aviatrix_spoke_external_device_conn`**, and **`aviatrix_edge_spoke_external_device_conn`** for BGP external device connections |
| AVX-75736 | Added `enable_symmetric_routing` to **`aviatrix_spoke_group`** for AWS gateways; added `enable_az_affinity` to **`aviatrix_spoke_transit_attachment`** and **`aviatrix_transit_gateway_peering`** for AWS gateways |
| AVX-75929 | Added `egress_path` attribute to rules in **`aviatrix_dcf_ruleset`**, controlling whether traffic egresses via the default path (`EGRESS_PATH_DEFAULT`) or is redirected through a local spoke gateway (`EGRESS_PATH_LOCAL`) |
| AVX-76307 | Added support for per-cluster feature flags to **`aviatrix_kubernetes_cluster`** |
| AVX-76490 | Added new **`aviatrix_attachment_vrf_status`** data source for querying and optionally toggling the `vrf_attachment_enabled` flag on eligible Aviatrix peerings |
| AVX-76515 | Added `private_network` attribute (`ForceNew`, AWS/Azure only) to **`aviatrix_spoke_group`** and **`aviatrix_transit_group`**, enabling gateway groups without public IPs; added `private_subnet_egress_target` attribute (`ForceNew`, requires `insane_mode`) to **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`** for specifying the HPE gateway egress target |
| AVX-79044 | Added support for positioning a Terraform-managed **`aviatrix_dcf_ruleset`** after Kubernetes CRD-managed rulesets |

### Deprecations:
| Issue | Description |
| :--- | :--- |
| AVX-49318 | The `watch` attribute on rules in **`aviatrix_dcf_ruleset`** is deprecated. Use the new `enforcement` attribute instead. The `watch` attribute will be removed in a future release. |
| AVX-76516 | Removed **`aviatrix_controller_private_mode_config`**, **`aviatrix_private_mode_lb`**, and **`aviatrix_private_mode_multicloud_endpoint`** resources. Private Mode is fully deprecated in Controller 10.0.0+. Remove these resources from any existing Terraform configurations before upgrading. |
| AVX-77577 | Deprecated the **`aviatrix_gateway_image`** data source. |

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-53161 | Fixed **`aviatrix_transit_firenet_policy`** `terraform import` failing in Azure when `inspected_resource_name` references a spoke subnet group with a name containing `~~` |
| AVX-55065 | Fixed **`aviatrix_edge_spoke_transit_attachment`** silently ignoring changes to `insane_mode_tunnel_number`; the attribute is now `ForceNew` to trigger destroy and recreate |
| AVX-59367 | Fixed **`aviatrix_transit_external_device_conn`** failing when attempting to remove `bgp_md5_key` after it was initially set |
| AVX-62273 | Fixed `gw_size` being incorrectly required on **`aviatrix_transit_gateway`** for Megaport Edge gateways |
| AVX-68466 | Fixed BGP community configuration being excluded from the Terraform export feature |
| AVX-74243 | Fixed **`aviatrix_distributed_firewalling_policy_list`** deployment policies for Azure not being applied due to a provider name case mismatch |
| AVX-75379 | Fixed **`aviatrix_spoke_transit_attachment`** failing with a not-found error when the spoke uses the gateway group/instance model |
| AVX-75380 | Fixed **`aviatrix_spoke_instance`** allowing `gw_size` to be removed from configuration without producing an error or plan diff |
| AVX-75652 | Fixed **`aviatrix_transit_instance`** post-create configuration not being executed, leaving the gateway in an incomplete state after `terraform apply` |
| AVX-75737 | Fixed **`aviatrix_vpc`** in GCP producing spurious subnet replacements on re-apply due to non-deterministic subnet ordering |
| AVX-76261 | Fixed **`aviatrix_spoke_ha_gateway`** returning an inconsistent result error ("Root object was present, but now absent") when `manage_ha_gateway = false` |
| AVX-77280 | Fixed **`aviatrix_dcf_ruleset`** in-place rule updates failing with `AVXERR-DFW-0004` when `port_ranges` was specified |
| AVX-78064 | Fixed **`aviatrix_transit_external_device_conn`**, **`aviatrix_spoke_external_device_conn`**, and **`aviatrix_edge_spoke_external_device_conn`** being read back as HA-enabled when configured as a single non-HA connection with comma-separated multi-tunnel values (`remote_gateway_ip`, `local_tunnel_cidr`, `remote_tunnel_cidr`). On edge-as-transit gateways this caused a spurious destroy/recreate on the second `terraform plan`. |
| AVX-78171 | Fixed **`aviatrix_edge_spoke`** failing to apply configuration on Equinix Edge-as-Transit gateways when cloud-assigned management interfaces were present |
| AVX-78300 | Fixed **`aviatrix_gateway`** and **`aviatrix_spoke_gateway`** allocating a new EIP even when `ha_eip` or `peering_ha_eip` was pre-specified |
| AVX-78304 | Fixed a perpetual `terraform plan` diff (`single_az_ha = true -> false`) on **`aviatrix_gateway`** resources that did not set `single_az_ha` explicitly. The schema default now matches the Controller default. |
| AVX-78638 | Fixed **`aviatrix_transit_group`** becoming permanently unreadable after a failed **`aviatrix_transit_instance`** creation. The instance-launch rollback no longer deletes the group's backing VPC object while the group still references it, so the group is no longer left with a dangling reference. The provider Read also now treats a "VPC not found" error as a missing resource and marks the group for recreation, allowing automatic recovery on the next `terraform apply`. |
| AVX-78640 | Fixed `enable_transit_firenet = true` on **`aviatrix_transit_group`** failing at `terraform apply` (`[AVXERR-TRANSIT-0121] <group-name> is not primary gateway`) when the group was created before any gateway existed. FireNet (`enable_transit_firenet`, `enable_firenet`, `enable_gateway_load_balancer`) is managed at the group level — the Controller applies the group's desired state when the primary instance launches. |
| AVX-78748 | Fixed **`aviatrix_gateway`** showing a persistent plan diff on `allocate_new_eip` for private-network gateways |
| AVX-78758 | Fixed **`aviatrix_spoke_external_device_conn`** writing Terraform state before the API call completed, leaving phantom state on creation failure |
| AVX-78760 | Fixed **`aviatrix_spoke_instance`** not populating the `eip` attribute after creation when `allocate_new_eip = true` |
| AVX-78765 | Fixed **`aviatrix_spoke_group`** silently ignoring `include_cidr` and `customized_spoke_vpc_routes` on update |
| AVX-79035 | Fixed **`aviatrix_transit_instance`** not saving tags to state, causing a persistent plan diff on every re-apply |
| AVX-79062 | Fixed **`aviatrix_tunnel`** returning "Root object was present, but now absent" for spoke peerings identified by group name |
| AVX-79383 | Fixed self-managed **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`** failing to decode the ZTP ISO file path during `terraform apply` |
| AVX-79490 | Fixed **`aviatrix_transit_external_device_conn`** reading back `local_tunnel_cidr` and `remote_tunnel_cidr` as empty when the local gateway has HA enabled, causing a spurious destroy/recreate on re-plan |
| AVX-79795 | Fixed **`aviatrix_firenet`** returning an error when associating firewall instances to a FireNet gateway on Controller 10.1.0 |
| AVX-79942 | Fixed **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`** failing to create GCP instances using the n4 machine family due to sending `zone` instead of `region` to the Controller |
| AVX-80061 | Fixed **`aviatrix_spoke_instance`** and **`aviatrix_transit_instance`** `terraform import` not working correctly |
| AVX-80246 | Fixed **`aviatrix_spoke_instance`** creation failing when `allocate_new_eip = false` |
| AVX-80429 | Fixed **`aviatrix_spoke_instance`** incorrectly accepting `customized_spoke_vpc_routes`, which is a group-level attribute on **`aviatrix_spoke_group`** |
| AVX-80581 | Fixed **`aviatrix_transit_instance`** and **`aviatrix_spoke_instance`** showing a spurious diff on `rx_queue_size` on every re-apply |
| AVX-80646 | Fixed **`aviatrix_transit_group`** `enable_hybrid_connection` always reading back `false` due to the provider reading from the wrong Controller response field |
| AVX-80656 | Fixed `vendor_name` being missing from the **`aviatrix_spoke_group`** resource schema |
| AVX-80657 | Fixed **`aviatrix_transit_instance`** `terraform destroy` failing when the gateway belongs to a FireNet-enabled transit group |
