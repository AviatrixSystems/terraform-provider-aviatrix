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
| AVX-78304 | The `single_az_ha` default for the `aviatrix_gateway` resource changed from `false` to `true` to match the Controller's behavior of enabling Single-AZ HA on every new gateway. Configurations that relied on the old `false` default without setting `single_az_ha` explicitly will now create gateways with Single-AZ HA enabled. To preserve the previous behavior, set `single_az_ha = false` explicitly. |

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-78304 | Fixed a perpetual `terraform plan` diff (`single_az_ha = true -> false`) on `aviatrix_gateway` resources that did not set `single_az_ha` explicitly. The schema default now matches the Controller default. |
| AVX-78064 | Fixed `aviatrix_transit_external_device_conn`, `aviatrix_spoke_external_device_conn`, and `aviatrix_edge_spoke_external_device_conn` being read back as HA-enabled when configured as a single non-HA connection with comma-separated multi-tunnel values (`remote_gateway_ip`, `local_tunnel_cidr`, `remote_tunnel_cidr`). On edge-as-transit gateways this caused a spurious destroy/recreate on the second `terraform plan`. |
| AVX-78638 | Fixed `aviatrix_transit_group` becoming permanently unreadable after a failed `aviatrix_transit_instance` creation. The instance-launch rollback no longer deletes the group's backing VPC object while the group still references it, so the group is no longer left with a dangling reference. The provider Read also now treats a "VPC not found" error as a missing resource and marks the group for recreation, allowing automatic recovery on the next `terraform apply`. |

## 9.1.0
### Notes:
- Supported Controller version: **9.1.0**

### Enhancements:
| Enhancement | Description |
| :--- | :--- |
| AVX-xxxx | Added the yyyy |

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-xxxx | Fixed an issue yyyy |
