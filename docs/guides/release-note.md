---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 10.0.10
### Notes:
- Supported Controller version: **10.0.10**

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-80735 | Fixed **`aviatrix_dcf_ruleset`** in-place rule updates failing with `AVXERR-DFW-0004` (`src_ads not a list`) when a rule change caused a phantom entry with empty `src_smart_groups` to be included in the PUT payload |
| AVX-82151 | Fixed **`aviatrix_dcf_ruleset`** showing the entire policy list as changed on every `terraform plan`/apply even when no rules were modified. The hash function now normalizes `watch` and `enforcement` fields before hashing, preventing spurious diffs. |
