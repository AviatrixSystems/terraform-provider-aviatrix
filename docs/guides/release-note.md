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
